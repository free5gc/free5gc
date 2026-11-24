package adapter

import (
	ausffactory "github.com/free5gc/ausf/pkg/factory"
	"github.com/free5gc/config/cfg-generator/global"
)

// AUSFAdapter applies global fields to ausfcfg.yaml (docs/global-mapping.md §1).
//
//   - plmn: plmnSupportList[0] only; further entries are additional PLMNs and
//     are kept (docs/global-mapping.md §2)
//   - nrf, sbiScheme, logger, timeouts: see applyCommon
type AUSFAdapter struct{}

var _ Adapter[ausffactory.Config] = AUSFAdapter{}

func (AUSFAdapter) NF() NFName       { return AUSF }
func (AUSFAdapter) FileName() string { return "ausfcfg.yaml" }

func (AUSFAdapter) Apply(g *global.Global, cfg *ausffactory.Config, r *Recorder) error {
	if g == nil {
		return nil
	}
	c := cfg.Configuration
	if c == nil {
		return errMissingConfiguration
	}

	if p := g.Plmn; p != nil && len(c.PlmnSupportList) > 0 {
		Set(r, AUSF, "plmnSupportList[0].mcc", &c.PlmnSupportList[0].Mcc, p.Mcc)
		Set(r, AUSF, "plmnSupportList[0].mnc", &c.PlmnSupportList[0].Mnc, p.Mnc)
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
		cfg.Logger = &ausffactory.Logger{}
	}
	if cfg.Logger != nil {
		f.LoggerEnable, f.LoggerLevel, f.LoggerReportCaller =
			&cfg.Logger.Enable, &cfg.Logger.Level, &cfg.Logger.ReportCaller
	}
	applyCommon(g, AUSF, f, r)
	return nil
}

func (AUSFAdapter) Validate(cfg *ausffactory.Config) error {
	_, err := cfg.Validate()
	return err
}
