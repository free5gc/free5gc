package adapter

import (
	"github.com/free5gc/config/cfg-generator/global"
	neffactory "github.com/free5gc/nef/pkg/factory"
)

// NEFAdapter applies global fields to nefcfg.yaml (docs/global-mapping.md §1).
//
//   - nrf, sbiScheme, logger, timeouts: see applyCommon
type NEFAdapter struct{}

var _ Adapter[neffactory.Config] = NEFAdapter{}

func (NEFAdapter) NF() NFName       { return NEF }
func (NEFAdapter) FileName() string { return "nefcfg.yaml" }

func (NEFAdapter) Apply(g *global.Global, cfg *neffactory.Config, r *Recorder) error {
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
	}
	if c.Sbi != nil {
		f.SbiScheme = &c.Sbi.Scheme
	}
	if cfg.Logger == nil && g.Logger != nil {
		cfg.Logger = &neffactory.Logger{}
	}
	if cfg.Logger != nil {
		f.LoggerEnable, f.LoggerLevel, f.LoggerReportCaller =
			&cfg.Logger.Enable, &cfg.Logger.Level, &cfg.Logger.ReportCaller
	}
	applyCommon(g, NEF, f, r)
	return nil
}

func (NEFAdapter) Validate(cfg *neffactory.Config) error {
	_, err := cfg.Validate()
	return err
}
