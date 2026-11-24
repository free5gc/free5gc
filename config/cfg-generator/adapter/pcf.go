package adapter

import (
	"github.com/free5gc/config/cfg-generator/global"
	pcffactory "github.com/free5gc/pcf/pkg/factory"
)

// PCFAdapter applies global fields to pcfcfg.yaml (docs/global-mapping.md §1).
//
//   - dnnList: PCF has no DNN field; it registers PcfInfo.DnnList hard-coded
//     as [free5gc internet] (pcf/internal/sbi/consumer/nrf_service.go). A
//     global DNN outside that list is recorded as NotApplied.
//   - nrf, sbiScheme, mongodb, logger, timeouts: see applyCommon
type PCFAdapter struct{}

var _ Adapter[pcffactory.Config] = PCFAdapter{}

// pcfRegisteredDnns is the DNN list PCF registers at the NRF regardless of
// its config.
var pcfRegisteredDnns = map[string]bool{"free5gc": true, "internet": true}

func (PCFAdapter) NF() NFName       { return PCF }
func (PCFAdapter) FileName() string { return "pcfcfg.yaml" }

func (PCFAdapter) Apply(g *global.Global, cfg *pcffactory.Config, r *Recorder) error {
	if g == nil {
		return nil
	}
	c := cfg.Configuration
	if c == nil {
		return errMissingConfiguration
	}

	for _, dnn := range g.Dnns {
		if !pcfRegisteredDnns[dnn] {
			r.NotApply(PCF, "NRF registration DNN list", "[free5gc internet]", dnn)
		}
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
		c.Mongodb = &pcffactory.Mongodb{}
	}
	if c.Mongodb != nil {
		f.MongoName, f.MongoUrl = &c.Mongodb.Name, &c.Mongodb.Url
	}
	if cfg.Logger == nil && g.Logger != nil {
		cfg.Logger = &pcffactory.Logger{}
	}
	if cfg.Logger != nil {
		f.LoggerEnable, f.LoggerLevel, f.LoggerReportCaller =
			&cfg.Logger.Enable, &cfg.Logger.Level, &cfg.Logger.ReportCaller
	}
	applyCommon(g, PCF, f, r)
	return nil
}

func (PCFAdapter) Validate(cfg *pcffactory.Config) error {
	_, err := cfg.Validate()
	return err
}
