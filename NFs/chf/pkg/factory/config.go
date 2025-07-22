/*
 * CHF Configuration Factory
 */

package factory

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"sync"

	"github.com/asaskevich/govalidator"

	"github.com/free5gc/chf/internal/logger"
)

const (
	ChfDefaultTLSKeyLogPath          = "./log/chfsslkey.log"
	ChfDefaultTLSPemPath             = "./cert/chf.pem"
	ChfDefaultTLSKeyPath             = "./cert/chf.key"
	ChfDefaultConfigPath             = "./config/chfcfg.yaml"
	ChfSbiDefaultIPv4                = "127.0.0.113"
	ChfSbiDefaultPort                = 8000
	ChfSbiDefaultScheme              = "https"
	ChfDefaultNrfUri                 = "https://127.0.0.10:8000"
	CgfDefaultCdrFilePath            = "/tmp"
	ConvergedChargingResUriPrefix    = "/nchf-convergedcharging/v3"
	OfflineOnlyChargingResUriPrefix  = "/nchf-offlineonlycharging/v1"
	SpendingLimitControlResUriPrefix = "/nchf-spendinglimitcontrol/v1"
)

type Config struct {
	Info          *Info          `yaml:"info" valid:"required"`
	Configuration *Configuration `yaml:"configuration" valid:"required"`
	Logger        *Logger        `yaml:"logger" valid:"required"`
	sync.RWMutex
}

func (c *Config) Validate() (bool, error) {
	if configuration := c.Configuration; configuration != nil {
		if result, err := configuration.validate(); err != nil {
			return result, err
		}
	}

	result, err := govalidator.ValidateStruct(c)
	return result, appendInvalid(err)
}

type Info struct {
	Version     string `yaml:"version,omitempty" valid:"required,in(1.0.3)"`
	Description string `yaml:"description,omitempty" valid:"-"`
}

type Configuration struct {
	ChfName             string    `yaml:"chfName,omitempty" valid:"required, type(string)"`
	Sbi                 *Sbi      `yaml:"sbi,omitempty" valid:"required"`
	ServiceNameList     []string  `yaml:"serviceNameList,omitempty" valid:"required"`
	NrfUri              string    `yaml:"nrfUri,omitempty" valid:"required, url"`
	NrfCertPem          string    `yaml:"nrfCertPem,omitempty" valid:"optional"`
	Mongodb             *Mongodb  `yaml:"mongodb" valid:"required"`
	VolumeLimit         int32     `yaml:"volumeLimit,omitempty" valid:"optional"`
	VolumeLimitPDU      int32     `yaml:"volumeLimitPDU,omitempty" valid:"optional"`
	ReserveQuotaRatio   int32     `yaml:"reserveQuotaRatio,omitempty" valid:"optional"`
	VolumeThresholdRate float32   `yaml:"volumeThresholdRate,omitempty" valid:"optional"`
	QuotaValidityTime   int32     `yaml:"quotaValidityTime,omitempty" valid:"optional"`
	RfDiameter          *Diameter `yaml:"rfDiameter,omitempty" valid:"required"`
	AbmfDiameter        *Diameter `yaml:"abmfDiameter,omitempty" valid:"required"`
	Cgf                 *Cgf      `yaml:"cgf,omitempty" valid:"required"`
}

type Logger struct {
	Enable       bool   `yaml:"enable" valid:"type(bool)"`
	Level        string `yaml:"level" valid:"required,in(trace|debug|info|warn|error|fatal|panic)"`
	ReportCaller bool   `yaml:"reportCaller" valid:"type(bool)"`
}

func (c *Configuration) validate() (bool, error) {
	if sbi := c.Sbi; sbi != nil {
		if result, err := sbi.validate(); err != nil {
			return result, err
		}
	}

	for index, serviceName := range c.ServiceNameList {
		switch serviceName {
		case "nchf-convergedcharging":
		case "nchf-offlineonlycharging":
		case "nchf-spendinglimitcontrol":
		default:
			err := errors.New("Invalid serviceNameList[" + strconv.Itoa(index) + "]: " +
				serviceName + ", should be nchf-convergedcharging.")
			return false, err
		}
	}

	result, err := govalidator.ValidateStruct(c)
	return result, appendInvalid(err)
}

type Service struct {
	ServiceName string `yaml:"serviceName" valid:"required, service"`
	SuppFeat    string `yaml:"suppFeat,omitempty" valid:"-"`
}

type Diameter struct {
	Protocol string `yaml:"protocol" valid:"required"`
	HostIPv4 string `yaml:"hostIPv4,omitempty" valid:"required,host"`
	Port     int    `yaml:"port,omitempty" valid:"required,port"`
	Tls      *Tls   `yaml:"tls,omitempty" valid:"optional"`
}

type Cgf struct {
	Enable                   bool   `yaml:"enable,omitempty" valid:"type(bool)"`
	HostIPv4                 string `yaml:"hostIPv4,omitempty" valid:"required,host"`
	Port                     int    `yaml:"port,omitempty" valid:"required,port"`
	ListenPort               int    `yaml:"listenPort,omitempty" valid:"required,port"`
	PassiveTransferPortRange struct {
		Start int `yaml:"start" valid:"required,port"`
		End   int `yaml:"end" valid:"required,port"`
	} `yaml:"passiveTransferPortRange,omitempty" valid:"optional"`
	Tls         *Tls   `yaml:"tls,omitempty" valid:"optional"`
	CdrFilePath string `yaml:"cdrFilePath,omitempty" valid:"optional"`
}
type Sbi struct {
	Scheme       string `yaml:"scheme" valid:"required,scheme"`
	RegisterIPv4 string `yaml:"registerIPv4,omitempty" valid:"required,host"` // IP that is registered at NRF.
	BindingIPv4  string `yaml:"bindingIPv4,omitempty" valid:"required,host"`  // IP used to run the server in the node.
	Port         int    `yaml:"port,omitempty" valid:"required,port"`
	Tls          *Tls   `yaml:"tls,omitempty" valid:"optional"`
}

func (s *Sbi) validate() (bool, error) {
	govalidator.TagMap["scheme"] = govalidator.Validator(func(str string) bool {
		return str == "https" || str == "http"
	})

	if tls := s.Tls; tls != nil {
		if result, err := tls.validate(); err != nil {
			return result, err
		}
	}

	result, err := govalidator.ValidateStruct(s)
	return result, appendInvalid(err)
}

type Tls struct {
	Pem string `yaml:"pem,omitempty" valid:"type(string),minstringlength(1),required"`
	Key string `yaml:"key,omitempty" valid:"type(string),minstringlength(1),required"`
}

func (t *Tls) validate() (bool, error) {
	result, err := govalidator.ValidateStruct(t)
	return result, err
}

type Mongodb struct {
	Name string `yaml:"name" valid:"required, type(string)"`
	Url  string `yaml:"url" valid:"required"`
}

// Commenting the unused function
// func (m *Mongodb) validate() (bool, error) {
// 	pattern := `[-a-zA-Z0-9@:%._\+~#=]{1,256}\b([-a-zA-Z0-9()@:%_\+.~#?&//=]*)`
// 	if result := govalidator.StringMatches(m.Url, pattern); !result {
// 		err := fmt.Errorf("Invalid Url: %s", m.Url)
// 		return result, err
// 	}
// 	result, err := govalidator.ValidateStruct(m)
// 	return result, err
// }

func appendInvalid(err error) error {
	var errs govalidator.Errors

	if err == nil {
		return nil
	}

	es := err.(govalidator.Errors).Errors()
	for _, e := range es {
		errs = append(errs, fmt.Errorf("invalid %w", e))
	}

	return error(errs)
}

func (c *Config) GetVersion() string {
	c.RLock()
	defer c.RUnlock()

	if c.Info.Version != "" {
		return c.Info.Version
	}
	return ""
}

func (c *Config) SetLogEnable(enable bool) {
	c.Lock()
	defer c.Unlock()

	if c.Logger == nil {
		logger.CfgLog.Warnf("Logger should not be nil")
		c.Logger = &Logger{
			Enable: enable,
			Level:  "info",
		}
	} else {
		c.Logger.Enable = enable
	}
}

func (c *Config) SetLogLevel(level string) {
	c.Lock()
	defer c.Unlock()

	if c.Logger == nil {
		logger.CfgLog.Warnf("Logger should not be nil")
		c.Logger = &Logger{
			Level: level,
		}
	} else {
		c.Logger.Level = level
	}
}

func (c *Config) SetLogReportCaller(reportCaller bool) {
	c.Lock()
	defer c.Unlock()

	if c.Logger == nil {
		logger.CfgLog.Warnf("Logger should not be nil")
		c.Logger = &Logger{
			Level:        "info",
			ReportCaller: reportCaller,
		}
	} else {
		c.Logger.ReportCaller = reportCaller
	}
}

func (c *Config) GetLogEnable() bool {
	c.RLock()
	defer c.RUnlock()
	if c.Logger == nil {
		logger.CfgLog.Warnf("Logger should not be nil")
		return false
	}
	return c.Logger.Enable
}

func (c *Config) GetLogLevel() string {
	c.RLock()
	defer c.RUnlock()
	if c.Logger == nil {
		logger.CfgLog.Warnf("Logger should not be nil")
		return "info"
	}
	return c.Logger.Level
}

func (c *Config) GetLogReportCaller() bool {
	c.RLock()
	defer c.RUnlock()
	if c.Logger == nil {
		logger.CfgLog.Warnf("Logger should not be nil")
		return false
	}
	return c.Logger.ReportCaller
}

func (c *Config) GetSbiBindingAddr() string {
	c.RLock()
	defer c.RUnlock()
	return c.GetSbiBindingIP() + ":" + strconv.Itoa(c.GetSbiPort())
}

func (c *Config) GetSbiBindingIP() string {
	c.RLock()
	defer c.RUnlock()
	bindIP := "0.0.0.0"
	if c.Configuration == nil || c.Configuration.Sbi == nil {
		return bindIP
	}
	if c.Configuration.Sbi.BindingIPv4 != "" {
		if bindIP = os.Getenv(c.Configuration.Sbi.BindingIPv4); bindIP != "" {
			logger.CfgLog.Infof("Parsing ServerIPv4 [%s] from ENV Variable", bindIP)
		} else {
			bindIP = c.Configuration.Sbi.BindingIPv4
		}
	}
	return bindIP
}

func (c *Config) GetSbiPort() int {
	c.RLock()
	defer c.RUnlock()
	if c.Configuration != nil && c.Configuration.Sbi != nil && c.Configuration.Sbi.Port != 0 {
		return c.Configuration.Sbi.Port
	}
	return ChfSbiDefaultPort
}

func (c *Config) GetSbiScheme() string {
	c.RLock()
	defer c.RUnlock()
	if c.Configuration != nil && c.Configuration.Sbi != nil && c.Configuration.Sbi.Scheme != "" {
		return c.Configuration.Sbi.Scheme
	}
	return ChfSbiDefaultScheme
}

func (c *Config) GetCertPemPath() string {
	c.RLock()
	defer c.RUnlock()
	return c.Configuration.Sbi.Tls.Pem
}

func (c *Config) GetCertKeyPath() string {
	c.RLock()
	defer c.RUnlock()
	return c.Configuration.Sbi.Tls.Key
}
