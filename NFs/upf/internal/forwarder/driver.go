package forwarder

import (
	"fmt"
	"net"
	"sync"

	"github.com/pkg/errors"
	"github.com/wmnsk/go-pfcp/ie"

	"github.com/free5gc/go-upf/internal/logger"
	"github.com/free5gc/go-upf/internal/report"
	"github.com/free5gc/go-upf/pkg/factory"
)

type Driver interface {
	Close()

	CreatePDR(uint64, *ie.IE) error
	UpdatePDR(uint64, *ie.IE) error
	RemovePDR(uint64, *ie.IE) error

	CreateFAR(uint64, *ie.IE) error
	UpdateFAR(uint64, *ie.IE) error
	RemoveFAR(uint64, *ie.IE) error

	CreateQER(uint64, *ie.IE) error
	UpdateQER(uint64, *ie.IE) error
	RemoveQER(uint64, *ie.IE) error

	CreateURR(uint64, *ie.IE) error
	UpdateURR(uint64, *ie.IE) ([]report.USAReport, error)
	RemoveURR(uint64, *ie.IE) ([]report.USAReport, error)
	QueryURR(uint64, uint32) ([]report.USAReport, error)

	CreateBAR(uint64, *ie.IE) error
	UpdateBAR(uint64, *ie.IE) error
	RemoveBAR(uint64, *ie.IE) error

	HandleReport(report.Handler)
}

func NewDriver(wg *sync.WaitGroup, cfg *factory.Config) (Driver, error) {
	cfgGtpu := cfg.Gtpu
	if cfgGtpu == nil {
		return nil, errors.Errorf("no Gtpu config")
	}

	logger.MainLog.Infof("starting Gtpu Forwarder [%s]", cfgGtpu.Forwarder)
	if cfgGtpu.Forwarder == "gtp5g" {
		var gtpuAddr string
		var mtu uint32
		for _, ifInfo := range cfgGtpu.IfList {
			mtu = ifInfo.MTU
			gtpuAddr = fmt.Sprintf("%s:%d", ifInfo.Addr, factory.UpfGtpDefaultPort)
			logger.MainLog.Infof("GTP Address: %q", gtpuAddr)
			break
		}
		if gtpuAddr == "" {
			return nil, errors.Errorf("not found GTP address")
		}
		driver, err := OpenGtp5g(wg, gtpuAddr, mtu)
		if err != nil {
			return nil, errors.Wrap(err, "open Gtp5g")
		}

		link := driver.Link()
		for _, dnn := range cfg.DnnList {
			_, dst, err := net.ParseCIDR(dnn.Cidr)
			if err != nil {
				logger.MainLog.Errorln(err)
				continue
			}
			err = link.RouteAdd(dst)
			if err != nil {
				driver.Close()
				return nil, err
			}
		}
		return driver, nil
	} else if cfgGtpu.Forwarder == "dummy" {
		logger.MainLog.Infof("Using Dummy Forwarder")
		return NewDummy(wg, cfgGtpu.IfList)
	}
	return nil, errors.Errorf("not support forwarder:%q", cfgGtpu.Forwarder)
}
