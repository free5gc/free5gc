package adapter

import (
	"time"

	"github.com/free5gc/config/cfg-generator/global"
)

// commonFields points at the fields that most NFs share under the same YAML
// names but with NF-specific Go types. A nil pointer means the NF has no such
// field and the corresponding global value is skipped for it.
type commonFields struct {
	LoggerEnable       *bool
	LoggerLevel        *string
	LoggerReportCaller *bool

	IdleTimeout     **time.Duration // configuration.idleTimeout
	ReadIdleTimeout **time.Duration // configuration.readIdleTimeout
	PingTimeout     **time.Duration // configuration.pingTimeout

	SbiScheme *string // configuration.sbi.scheme
	NrfUri    *string // configuration.nrfUri

	MongoName *string // configuration.mongodb.name (NRF: MongoDBName)
	MongoUrl  *string // configuration.mongodb.url  (NRF: MongoDBUrl)
	MongoPath string  // path prefix for the change log, e.g. "mongodb"
}

// applyCommon applies logger, timeouts, sbiScheme, nrf and mongodb
// (docs/global-mapping.md §1.5-§1.9).
func applyCommon(g *global.Global, nf NFName, f commonFields, r *Recorder) {
	if l := g.Logger; l != nil {
		if l.Enable != nil && f.LoggerEnable != nil {
			Set(r, nf, "logger.enable", f.LoggerEnable, *l.Enable)
		}
		if l.Level != nil && f.LoggerLevel != nil {
			Set(r, nf, "logger.level", f.LoggerLevel, *l.Level)
		}
		if l.ReportCaller != nil && f.LoggerReportCaller != nil {
			Set(r, nf, "logger.reportCaller", f.LoggerReportCaller, *l.ReportCaller)
		}
	}

	if t := g.Timeouts; t != nil {
		if t.Idle != nil && f.IdleTimeout != nil {
			setDuration(r, nf, "idleTimeout", f.IdleTimeout, *t.Idle)
		}
		if t.ReadIdle != nil && f.ReadIdleTimeout != nil {
			setDuration(r, nf, "readIdleTimeout", f.ReadIdleTimeout, *t.ReadIdle)
		}
		if t.Ping != nil && f.PingTimeout != nil {
			setDuration(r, nf, "pingTimeout", f.PingTimeout, *t.Ping)
		}
	}

	if g.SbiScheme != nil && f.SbiScheme != nil {
		Set(r, nf, "sbi.scheme", f.SbiScheme, *g.SbiScheme)
	}
	if g.Nrf != nil && f.NrfUri != nil {
		Set(r, nf, "nrfUri", f.NrfUri, g.Nrf.URI())
	}
	if m := g.MongoDB; m != nil {
		name, url := f.MongoPath+".name", f.MongoPath+".url"
		if f.MongoPath == "" { // NRF keeps both fields at the top level
			name, url = "MongoDBName", "MongoDBUrl"
		}
		if m.Name != nil && f.MongoName != nil {
			Set(r, nf, name, f.MongoName, *m.Name)
		}
		if m.Url != nil && f.MongoUrl != nil {
			Set(r, nf, url, f.MongoUrl, *m.Url)
		}
	}
}
