package adapter

import (
	"github.com/free5gc/config/cfg-generator/global"
	scpfactory "github.com/free5gc/scp/pkg/factory"
)

// SCPAdapter applies global fields to scpcfg.yaml (docs/global-mapping.md §1).
//
//   - nrf, sbiScheme, logger, timeouts: see applyCommon
//
// The per-NF URIs (amfUri, smfUri, ...) depend on other NFs' SBI addresses
// and belong to the Topology stage, which is not implemented yet
// (docs/adapter-data-model.md §6).
type SCPAdapter struct{}

var _ Adapter[scpfactory.Config] = SCPAdapter{}

func (SCPAdapter) NF() NFName       { return SCP }
func (SCPAdapter) FileName() string { return "scpcfg.yaml" }

func (SCPAdapter) Apply(g *global.Global, cfg *scpfactory.Config, r *Recorder) error {
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
		cfg.Logger = &scpfactory.Logger{}
	}
	if cfg.Logger != nil {
		f.LoggerEnable, f.LoggerLevel, f.LoggerReportCaller =
			&cfg.Logger.Enable, &cfg.Logger.Level, &cfg.Logger.ReportCaller
	}
	applyCommon(g, SCP, f, r)
	return nil
}

func (SCPAdapter) Validate(cfg *scpfactory.Config) error {
	_, err := cfg.Validate()
	return err
}
