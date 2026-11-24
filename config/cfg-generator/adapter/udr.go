package adapter

import (
	"github.com/free5gc/config/cfg-generator/global"
	udrfactory "github.com/free5gc/udr/pkg/factory"
)

// UDRAdapter applies global fields to udrcfg.yaml (docs/global-mapping.md §1).
//
//   - nrf, sbiScheme, mongodb, logger, timeouts: see applyCommon
type UDRAdapter struct{}

var _ Adapter[udrfactory.Config] = UDRAdapter{}

func (UDRAdapter) NF() NFName       { return UDR }
func (UDRAdapter) FileName() string { return "udrcfg.yaml" }

func (UDRAdapter) Apply(g *global.Global, cfg *udrfactory.Config, r *Recorder) error {
	if g == nil {
		return nil
	}
	c := cfg.Configuration
	if c == nil {
		return errMissingConfiguration
	}

	f := commonFields{
		IdleTimeout:     &c.IdleTimeout,
		ReadIdleTimeout: &c.ReadIdleTimeout,
		PingTimeout:     &c.PingTimeout,
		NrfUri:          &c.NrfUri,
		MongoPath:       "mongodb",
	}
	if c.Sbi != nil {
		f.SbiScheme = &c.Sbi.Scheme
	}
	if c.Mongodb == nil && g.MongoDB != nil {
		c.Mongodb = &udrfactory.Mongodb{}
	}
	if c.Mongodb != nil {
		f.MongoName, f.MongoUrl = &c.Mongodb.Name, &c.Mongodb.Url
	}
	if cfg.Logger == nil && g.Logger != nil {
		cfg.Logger = &udrfactory.Logger{}
	}
	if cfg.Logger != nil {
		f.LoggerEnable, f.LoggerLevel, f.LoggerReportCaller =
			&cfg.Logger.Enable, &cfg.Logger.Level, &cfg.Logger.ReportCaller
	}
	applyCommon(g, UDR, f, r)
	return nil
}

func (UDRAdapter) Validate(cfg *udrfactory.Config) error {
	_, err := cfg.Validate()
	return err
}
