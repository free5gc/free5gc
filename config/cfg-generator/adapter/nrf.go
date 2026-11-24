package adapter

import (
	"fmt"
	"strconv"

	"github.com/free5gc/config/cfg-generator/global"
	nrffactory "github.com/free5gc/nrf/pkg/factory"
)

// NRFAdapter applies global fields to nrfcfg.yaml (docs/global-mapping.md §1).
//
//   - plmn:     DefaultPlmnId. NFs that register without a PlmnList (AUSF,
//     UDM, UDR, PCF, CHF, NEF, SCP) get this PLMN at the NRF.
//   - mongodb:  MongoDBName, MongoDBUrl (flat, capitalised keys)
//   - nrf:      NRF has no nrfUri. Its own SBI is not overwritten; instead the
//     global endpoint must match it, otherwise every other NF would point to
//     an address the NRF does not serve.
//   - sbiScheme, logger, timeouts: see applyCommon
type NRFAdapter struct{}

var _ Adapter[nrffactory.Config] = NRFAdapter{}

func (NRFAdapter) NF() NFName       { return NRF }
func (NRFAdapter) FileName() string { return "nrfcfg.yaml" }

func (NRFAdapter) Apply(g *global.Global, cfg *nrffactory.Config, r *Recorder) error {
	if g == nil {
		return nil
	}
	c := cfg.Configuration
	if c == nil {
		return errMissingConfiguration
	}

	// Check before changing anything. sbiScheme is applied below, so the
	// effective scheme is the global one when given.
	if n := g.Nrf; n != nil && c.Sbi != nil {
		scheme := c.Sbi.Scheme
		if g.SbiScheme != nil {
			scheme = *g.SbiScheme
		}
		port := c.Sbi.Port
		if port == 0 {
			port = nrffactory.NrfSbiDefaultPort
		}
		have := scheme + "://" + c.Sbi.RegisterIPv4 + ":" + strconv.Itoa(port)
		if scheme != n.Scheme || c.Sbi.RegisterIPv4 != n.Host || port != n.Port {
			return fmt.Errorf("global nrf %s does not match the NRF's own SBI %s "+
				"(nrf.configuration.sbi scheme/registerIPv4/port); update the nrf block too", n.URI(), have)
		}
	}

	if p := g.Plmn; p != nil {
		Set(r, NRF, "DefaultPlmnId.mcc", &c.DefaultPlmnId.Mcc, p.Mcc)
		Set(r, NRF, "DefaultPlmnId.mnc", &c.DefaultPlmnId.Mnc, p.Mnc)
	}

	f := commonFields{
		IdleTimeout:     &c.IdleTimeout,
		ReadIdleTimeout: &c.ReadIdleTimeout,
		PingTimeout:     &c.PingTimeout,
		MongoName:       &c.MongoDBName,
		MongoUrl:        &c.MongoDBUrl,
		// MongoPath "" selects the flat MongoDBName/MongoDBUrl paths.
	}
	if c.Sbi != nil {
		f.SbiScheme = &c.Sbi.Scheme
	}
	if cfg.Logger == nil && g.Logger != nil {
		cfg.Logger = &nrffactory.Logger{}
	}
	if cfg.Logger != nil {
		f.LoggerEnable, f.LoggerLevel, f.LoggerReportCaller =
			&cfg.Logger.Enable, &cfg.Logger.Level, &cfg.Logger.ReportCaller
	}
	applyCommon(g, NRF, f, r)
	return nil
}

func (NRFAdapter) Validate(cfg *nrffactory.Config) error {
	_, err := cfg.Validate()
	return err
}
