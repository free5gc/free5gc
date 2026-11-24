package adapter

import (
	chffactory "github.com/free5gc/chf/pkg/factory"
	"github.com/free5gc/config/cfg-generator/global"
)

// CHFAdapter applies global fields to chfcfg.yaml (docs/global-mapping.md §1).
//
//   - nrf, sbiScheme, mongodb, logger, timeouts: see applyCommon
//
// The cgf block pairs with webui billingServer, but their port fields are
// swapped in meaning (docs/nf-config-inventory.md §4); it is not touched.
type CHFAdapter struct{}

var _ Adapter[chffactory.Config] = CHFAdapter{}

func (CHFAdapter) NF() NFName       { return CHF }
func (CHFAdapter) FileName() string { return "chfcfg.yaml" }

func (CHFAdapter) Apply(g *global.Global, cfg *chffactory.Config, r *Recorder) error {
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
		c.Mongodb = &chffactory.Mongodb{}
	}
	if c.Mongodb != nil {
		f.MongoName, f.MongoUrl = &c.Mongodb.Name, &c.Mongodb.Url
	}
	if cfg.Logger == nil && g.Logger != nil {
		cfg.Logger = &chffactory.Logger{}
	}
	if cfg.Logger != nil {
		f.LoggerEnable, f.LoggerLevel, f.LoggerReportCaller =
			&cfg.Logger.Enable, &cfg.Logger.Level, &cfg.Logger.ReportCaller
	}
	applyCommon(g, CHF, f, r)
	return nil
}

func (CHFAdapter) Validate(cfg *chffactory.Config) error {
	_, err := cfg.Validate()
	return err
}
