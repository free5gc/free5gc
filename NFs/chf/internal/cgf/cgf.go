// ftpserver allows to create your own FTP(S) server
package cgf

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime/debug"
	"strconv"
	"sync"
	"time"

	"github.com/fclairamb/ftpserver/config"
	"github.com/fclairamb/ftpserver/server"
	ftpserver "github.com/fclairamb/ftpserverlib"
	"github.com/jlaffaye/ftp"

	"github.com/free5gc/chf/internal/logger"
	"github.com/free5gc/chf/pkg/factory"
)

type Cgf struct {
	ftpServer *ftpserver.FtpServer
	driver    *server.Server
	conn      *ftp.ServerConn
	addr      string
	ftpConfig FtpConfig

	connMutex sync.Mutex
}

type Access struct {
	User   string            `json:"user"`
	Pass   string            `json:"pass"`
	Fs     string            `json:"fs"`
	Params map[string]string `json:"params"`
}

type PortRange struct {
	Start int `json:"start"`
	End   int `json:"end"`
}

type FtpConfig struct {
	Version       int      `json:"version"`
	Accesses      []Access `json:"accesses"`
	ListenAddress string   `json:"listen_address"`

	Passive_transfer_port_range PortRange `json:"passive_transfer_port_range"`
}

var cgf *Cgf

var CGFEnable bool = false

func OpenServer(ctx context.Context, wg *sync.WaitGroup) *Cgf {
	// Arguments vars
	cgf = new(Cgf)

	cgfConfig := factory.ChfConfig.Configuration.Cgf
	cgf.addr = cgfConfig.HostIPv4 + ":" + strconv.Itoa(cgfConfig.Port)
	// set default port value
	startPort := 2123
	endPort := 2130
	if cgfConfig.PassiveTransferPortRange.Start != 0 {
		startPort = cgfConfig.PassiveTransferPortRange.Start
	}
	if cgfConfig.PassiveTransferPortRange.End != 0 {
		endPort = cgfConfig.PassiveTransferPortRange.End
	}

	cgf.ftpConfig = FtpConfig{
		Version: 1,
		Accesses: []Access{
			{
				User: "admin",
				Pass: "free5gc",
				Fs:   "os",
				Params: map[string]string{
					"basePath": "/tmp",
				},
			},
		},
		Passive_transfer_port_range: PortRange{
			Start: startPort,
			End:   endPort,
		},
		ListenAddress: factory.ChfConfig.Configuration.Sbi.BindingIPv4 + ":" + strconv.Itoa(cgfConfig.ListenPort),
	}

	file, err := os.Create("/tmp/config.json")
	if err != nil {
		panic(err)
	}
	defer func() {
		if close_err := file.Close(); close_err != nil {
			logger.CfgLog.Error("Can't close file", close_err)
		}
	}()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if errEncode := encoder.Encode(cgf.ftpConfig); errEncode != nil {
		panic(errEncode)
	}

	conf, errConfig := config.NewConfig("/tmp/config.json", logger.FtpServerLog)
	if errConfig != nil {
		logger.CgfLog.Error("Can't load conf", "Err", errConfig)

		return nil
	}

	// Loading the driver
	var errNewServer error
	cgf.driver, errNewServer = server.NewServer(conf, logger.FtpServerLog)

	if errNewServer != nil {
		logger.CgfLog.Error("Could not load the driver", "err", errNewServer)

		return nil
	}

	// Instantiating the server by passing our driver implementation
	cgf.ftpServer = ftpserver.NewFtpServer(cgf.driver)

	// Setting up the ftpserver logger
	cgf.ftpServer.Logger = logger.FtpServerLog

	go cgf.Serve(ctx, wg)
	logger.CgfLog.Info("FTP server Start")

	return cgf
}

func Login() error {
	cgf.connMutex.Lock()
	defer cgf.connMutex.Unlock()

	if cgf.conn != nil {
		ping_err := cgf.conn.NoOp()
		if ping_err == nil {
			logger.CgfLog.Infof("FTP already login.")
			return nil
		}
	}

	// FTP server is for CDR transfer
	var c *ftp.ServerConn

	c, err := ftp.Dial(cgf.addr, ftp.DialWithTimeout(2*time.Second))
	if err != nil {
		return err
	}

	err = c.Login(cgf.ftpConfig.Accesses[0].User, cgf.ftpConfig.Accesses[0].Pass)
	if err != nil {
		logger.CgfLog.Warnf("Login FTP server fail")
		return err
	}

	logger.CgfLog.Info("Login FTP server succeed")
	cgf.conn = c
	return err
}

func SendCDR(supi string) error {
	logger.CfgLog.Debugln("SendCDR:", supi)
	if !CGFEnable {
		logger.CfgLog.Warningln("CGF Not enable: SendCDR() didn't do anything.")
		return nil
	}

	if cgf.conn == nil {
		err := Login()
		if err != nil {
			return err
		}
		logger.CgfLog.Infof("FTP Re-Login Success")
	}

	ping_err := cgf.conn.NoOp()
	if ping_err != nil {
		logger.CgfLog.Infof("Faile to ping FTP server, relogin...")
		err := Login()
		if err != nil {
			return err
		}
		logger.CgfLog.Infof("FTP Re-Login Success")
	}
	cgf.connMutex.Lock()
	defer cgf.connMutex.Unlock()

	fileName := supi + ".cdr"
	cdrByte, err := os.ReadFile("/tmp/" + fileName)
	if err != nil {
		return err
	}

	cdrReader := bytes.NewReader(cdrByte)
	stor_err := cgf.conn.Stor(fileName, cdrReader)
	if stor_err != nil {
		return err
	}

	// check file exist and verify size
	entries, err := cgf.conn.List(".")
	if err != nil {
		logger.CgfLog.Warnf("List entries error: %+v, the CDR file not check", err)
		return nil
	}

	fileUploaded := false
	for _, entry := range entries {
		if entry.Name == fileName {
			fileUploaded = true
			logger.CgfLog.Debugln("File found in remote")
			if entry.Size == uint64(len(cdrByte)) {
				logger.CgfLog.Debugln("File size matches")
			} else {
				logger.CgfLog.Warningf("File size does not match: local[%v], remote[%v]", len(cdrByte), entry.Size)
			}
			break
		}
	}
	if !fileUploaded {
		logger.CgfLog.Warningln("File upload failed.")
		return fmt.Errorf("sendCDR failed: %+v", err)
	}
	logger.CgfLog.Infof("SendCDR success: %+v", supi)
	return nil
}

const (
	FTP_LOGIN_RETRY_NUMBER       = 3
	FTP_LOGIN_RETRY_WAITING_TIME = 1 * time.Second // second
)

func (f *Cgf) Serve(ctx context.Context, wg *sync.WaitGroup) {
	go func() {
		defer func() {
			if p := recover(); p != nil {
				// Print stack for panic to log. Fatalf() will let program exit.
				logger.InitLog.Fatalf("panic: %v\n%s", p, string(debug.Stack()))
			}
		}()
		<-ctx.Done()
		f.Terminate()
		wg.Done()
	}()

	for i := 0; i < FTP_LOGIN_RETRY_NUMBER; i++ {
		if err := Login(); err != nil {
			logger.CgfLog.Warnf("Login to Webconsole FTP fail: %s, retrying [%d]\n", err, i+1)
			time.Sleep(FTP_LOGIN_RETRY_WAITING_TIME)
		} else {
			break
		}
	}

	if err := f.ftpServer.ListenAndServe(); err != nil {
		logger.CgfLog.Error("Problem listening", "err", err)
	}

	// We wait at most 1 minutes for all clients to disconnect
	if err := f.driver.WaitGracefully(time.Minute); err != nil {
		logger.CgfLog.Warn("Problem stopping server", "Err", err)
	}
}

func (f *Cgf) Terminate() {
	logger.CgfLog.Infoln("CGF Terminating...")

	f.driver.Stop()

	if err := f.ftpServer.Stop(); err != nil {
		logger.CgfLog.Error("Problem stopping server", "Err", err)
	}

	if f.conn != nil {
		// close ftp connection if exist
		if err := f.conn.Quit(); err != nil {
			logger.CgfLog.Errorf("Problem stopping connection: %+v", err)
		}
	}

	var cdrFilePath string
	if factory.ChfConfig.Configuration.Cgf.CdrFilePath == "" {
		cdrFilePath = factory.CgfDefaultCdrFilePath
	} else {
		cdrFilePath = factory.ChfConfig.Configuration.Cgf.CdrFilePath
	}
	files, err := filepath.Glob(cdrFilePath + "/*.cdr")
	if err != nil {
		logger.CgfLog.Warnln("no CDR file")
	}

	for _, file := range files {
		if _, errStat := os.Stat(file); errStat == nil {
			logger.CgfLog.Infof("Remove CDR file: %s", file)
			if errRemove := os.Remove(file); errRemove != nil {
				logger.CgfLog.Warnf("Failed to remove CDR file: %s\n", file)
			}
		}
	}
	logger.CgfLog.Infoln("CGF terminated")
}
