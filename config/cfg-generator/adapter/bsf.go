package adapter

import (
	"errors"
	"fmt"

	"github.com/free5gc/openapi/models"

	bsffactory "github.com/free5gc/bsf/pkg/factory"
	"github.com/free5gc/config/cfg-generator/global"
)

// BSFAdapter applies global fields to bsfcfg.yaml (docs/global-mapping.md §1).
//
//   - plmn: BSF has no PLMN field; it registers PLMN 208/93 hard-coded
//     (bsf/internal/context/context.go, GetBsfProfile). A different global
//     PLMN is recorded as NotApplied.
//   - nrf, sbiScheme, mongodb, logger, timeouts: see applyCommon
type BSFAdapter struct{}

var _ Adapter[bsffactory.Config] = BSFAdapter{}

// bsfRegisteredPlmn is the PLMN BSF registers at the NRF regardless of its
// config.
var bsfRegisteredPlmn = models.PlmnId{Mcc: "208", Mnc: "93"}

func (BSFAdapter) NF() NFName       { return BSF }
func (BSFAdapter) FileName() string { return "bsfcfg.yaml" }

func (BSFAdapter) Apply(g *global.Global, cfg *bsffactory.Config, r *Recorder) error {
	if g == nil {
		return nil
	}
	c := cfg.Configuration
	if c == nil {
		return errMissingConfiguration
	}

	if p := g.Plmn; p != nil && *p != bsfRegisteredPlmn {
		r.NotApply(BSF, "NRF registration PLMN",
			bsfRegisteredPlmn.Mcc+"/"+bsfRegisteredPlmn.Mnc, p.Mcc+"/"+p.Mnc)
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
	if c.MongoDB == nil && g.MongoDB != nil {
		c.MongoDB = &bsffactory.MongoDB{}
	}
	if c.MongoDB != nil {
		f.MongoName, f.MongoUrl = &c.MongoDB.Name, &c.MongoDB.Url
	}
	if cfg.Logger == nil && g.Logger != nil {
		cfg.Logger = &bsffactory.Logger{}
	}
	if cfg.Logger != nil {
		f.LoggerEnable, f.LoggerLevel, f.LoggerReportCaller =
			&cfg.Logger.Enable, &cfg.Logger.Level, &cfg.Logger.ReportCaller
	}
	applyCommon(g, BSF, f, r)
	return nil
}

// Validate checks BSF configs in the generator because BSF itself does not:
// its ReadConfig only unmarshals, its model has no valid tags, and a version
// mismatch is only logged (docs/nf-config-inventory.md §5.1).
func (BSFAdapter) Validate(cfg *bsffactory.Config) error {
	var errs []error
	add := func(format string, args ...any) { errs = append(errs, fmt.Errorf(format, args...)) }

	if cfg.Info == nil || cfg.Info.Version != bsffactory.BsfExpectedConfigVersion {
		add("info.version must be %q", bsffactory.BsfExpectedConfigVersion)
	}
	if cfg.Logger == nil {
		add("logger is required")
	} else if !validLogLevels[cfg.Logger.Level] {
		add("logger.level %q is invalid", cfg.Logger.Level)
	}

	c := cfg.Configuration
	if c == nil {
		add("configuration is required")
		return errors.Join(errs...)
	}
	if s := c.Sbi; s == nil {
		add("configuration.sbi is required")
	} else {
		if !isHTTPScheme(s.Scheme) {
			add("configuration.sbi.scheme %q must be http or https", s.Scheme)
		}
		if s.RegisterIPv4 == "" && s.RegisterIPv6 == "" {
			add("configuration.sbi.registerIPv4 or registerIPv6 is required")
		}
		if s.BindingIPv4 == "" && s.BindingIPv6 == "" {
			add("configuration.sbi.bindingIPv4 or bindingIPv6 is required")
		}
		if s.Port < 0 || s.Port > 65535 {
			add("configuration.sbi.port %d must be in 0..65535", s.Port)
		}
	}
	if !isAbsoluteURL(c.NrfUri) {
		add("configuration.nrfUri %q must be an absolute URL", c.NrfUri)
	}
	if m := c.MongoDB; m == nil || m.Name == "" || m.Url == "" {
		add("configuration.mongodb.name and url are required")
	} else if !isMongoURL(m.Url) {
		add("configuration.mongodb.url %q must start with mongodb:// or mongodb+srv://", m.Url)
	}
	return errors.Join(errs...)
}
