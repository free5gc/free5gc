package adapter

import (
	"errors"
	"fmt"

	"github.com/free5gc/openapi/models"

	amffactory "github.com/free5gc/amf/pkg/factory"
	"github.com/free5gc/config/cfg-generator/global"
)

// AMFAdapter applies global fields to amfcfg.yaml (docs/global-mapping.md §1).
//
//   - plmn:       servedGuamiList[*].plmnId (nid kept), supportTaiList[*].plmnId,
//     plmnSupportList[*].plmnId
//   - tac:        supportTaiList[*].tac
//   - snssaiList: plmnSupportList[*].snssaiList[*]                (value mapping)
//   - dnnList:    supportDnnList[*]                               (value mapping);
//     supportLadnList[*].dnn follows the same mapping
//   - nrf, sbiScheme, logger, timeouts: see applyCommon
type AMFAdapter struct{}

var _ Adapter[amffactory.Config] = AMFAdapter{}

func (AMFAdapter) NF() NFName       { return AMF }
func (AMFAdapter) FileName() string { return "amfcfg.yaml" }

func (a AMFAdapter) Apply(g *global.Global, cfg *amffactory.Config, r *Recorder) error {
	if g == nil {
		return nil
	}
	c := cfg.Configuration
	if c == nil {
		return errors.New("configuration block is missing")
	}

	// Build every value map first so a mismatch fails before any change.
	snssaiMap, err := buildValueMap(amfSnssais(c), g.Snssais, "snssaiList", "amf plmnSupportList[*].snssaiList")
	if err != nil {
		return err
	}
	dnnMap, err := buildValueMap(amfDnns(c), g.Dnns, "dnnList", "amf supportDnnList")
	if err != nil {
		return err
	}

	if p := g.Plmn; p != nil {
		amfApplyPlmn(c, p, r)
	}

	if g.Tac != nil {
		for i := range c.SupportTAIList {
			Set(r, AMF, fmt.Sprintf("supportTaiList[%d].tac", i), &c.SupportTAIList[i].Tac, *g.Tac)
		}
	}

	for i := range c.PlmnSupportList {
		for j := range c.PlmnSupportList[i].SNssaiList {
			mapSnssai(r, AMF, fmt.Sprintf("plmnSupportList[%d].snssaiList[%d]", i, j),
				&c.PlmnSupportList[i].SNssaiList[j], snssaiMap)
		}
	}

	for i := range c.SupportDnnList {
		mapString(r, AMF, fmt.Sprintf("supportDnnList[%d]", i), &c.SupportDnnList[i], dnnMap)
	}
	for i := range c.SupportLadnList {
		mapString(r, AMF, fmt.Sprintf("supportLadnList[%d].dnn", i), &c.SupportLadnList[i].Dnn, dnnMap)
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
	// Create the logger block only when the global block overrides it, so a
	// run without global.logger never changes the output.
	if cfg.Logger == nil && g.Logger != nil {
		cfg.Logger = &amffactory.Logger{}
	}
	if cfg.Logger != nil {
		f.LoggerEnable = &cfg.Logger.Enable
		f.LoggerLevel = &cfg.Logger.Level
		f.LoggerReportCaller = &cfg.Logger.ReportCaller
	}
	applyCommon(g, AMF, f, r)
	return nil
}

func amfApplyPlmn(c *amffactory.Configuration, p *models.PlmnId, r *Recorder) {
	for i := range c.ServedGumaiList {
		guami := &c.ServedGumaiList[i]
		if guami.PlmnId == nil {
			guami.PlmnId = &models.PlmnIdNid{}
		}
		base := fmt.Sprintf("servedGuamiList[%d].plmnId", i)
		Set(r, AMF, base+".mcc", &guami.PlmnId.Mcc, p.Mcc)
		Set(r, AMF, base+".mnc", &guami.PlmnId.Mnc, p.Mnc)
	}
	for i := range c.SupportTAIList {
		tai := &c.SupportTAIList[i]
		if tai.PlmnId == nil {
			tai.PlmnId = &models.PlmnId{}
		}
		base := fmt.Sprintf("supportTaiList[%d].plmnId", i)
		Set(r, AMF, base+".mcc", &tai.PlmnId.Mcc, p.Mcc)
		Set(r, AMF, base+".mnc", &tai.PlmnId.Mnc, p.Mnc)
	}
	for i := range c.PlmnSupportList {
		item := &c.PlmnSupportList[i]
		if item.PlmnId == nil {
			item.PlmnId = &models.PlmnId{}
		}
		base := fmt.Sprintf("plmnSupportList[%d].plmnId", i)
		Set(r, AMF, base+".mcc", &item.PlmnId.Mcc, p.Mcc)
		Set(r, AMF, base+".mnc", &item.PlmnId.Mnc, p.Mnc)
	}
}

// amfSnssais returns the distinct S-NSSAIs in path order.
func amfSnssais(c *amffactory.Configuration) []models.Snssai {
	var d distinct[models.Snssai]
	for _, item := range c.PlmnSupportList {
		for _, s := range item.SNssaiList {
			d.add(s)
		}
	}
	return d.order
}

// amfDnns returns the distinct DNNs of supportDnnList. LADN DNNs are not
// counted: they are optional and must be a subset of the supported DNNs.
func amfDnns(c *amffactory.Configuration) []string {
	var d distinct[string]
	for _, dnn := range c.SupportDnnList {
		d.add(dnn)
	}
	return d.order
}

func (AMFAdapter) Validate(cfg *amffactory.Config) error {
	if _, err := cfg.Validate(); err != nil {
		return err
	}
	return nil
}
