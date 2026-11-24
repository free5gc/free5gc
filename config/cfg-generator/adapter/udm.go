package adapter

import (
	"github.com/free5gc/config/cfg-generator/global"
	udmfactory "github.com/free5gc/udm/pkg/factory"
)

// UDMAdapter applies global fields to udmcfg.yaml (docs/global-mapping.md §1).
//
//   - nrf, sbiScheme, logger, timeouts: see applyCommon
type UDMAdapter struct{}

var _ Adapter[udmfactory.Config] = UDMAdapter{}

func (UDMAdapter) NF() NFName       { return UDM }
func (UDMAdapter) FileName() string { return "udmcfg.yaml" }

func (UDMAdapter) Apply(g *global.Global, cfg *udmfactory.Config, r *Recorder) error {
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
		cfg.Logger = &udmfactory.Logger{}
	}
	if cfg.Logger != nil {
		f.LoggerEnable, f.LoggerLevel, f.LoggerReportCaller =
			&cfg.Logger.Enable, &cfg.Logger.Level, &cfg.Logger.ReportCaller
	}
	applyCommon(g, UDM, f, r)
	return nil
}

func (UDMAdapter) Validate(cfg *udmfactory.Config) error {
	_, err := cfg.Validate()
	return err
}
