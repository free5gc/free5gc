package adapter

import (
	"fmt"

	"github.com/free5gc/openapi/models"

	"github.com/free5gc/config/cfg-generator/global"
	n3iwffactory "github.com/free5gc/n3iwf/pkg/factory"
)

// N3IWFAdapter applies global fields to n3iwfcfg.yaml (docs/global-mapping.md §1).
//
//   - plmn: n3iwfInformation.globalN3IWFID.plmnID,
//     n3iwfInformation.supportedTAList[*].broadcastPlmnList[*].plmnID
//   - tac:  n3iwfInformation.supportedTAList[*].tac
//   - snssaiList: ...broadcastPlmnList[*].taiSliceSupportList[*].snssai (value mapping)
//   - logger, timeouts (idleTimeout only): see applyCommon
//
// The n3iwf factory package only builds on Linux.
type N3IWFAdapter struct{}

var _ Adapter[n3iwffactory.Config] = N3IWFAdapter{}

func (N3IWFAdapter) NF() NFName       { return N3IWF }
func (N3IWFAdapter) FileName() string { return "n3iwfcfg.yaml" }

func (N3IWFAdapter) Apply(g *global.Global, cfg *n3iwffactory.Config, r *Recorder) error {
	if g == nil {
		return nil
	}
	c := cfg.Configuration
	if c == nil {
		return errMissingConfiguration
	}

	if info := c.N3IWFInfo; info != nil {
		var found distinct[models.Snssai]
		for _, ta := range info.SupportedTAList {
			for _, bp := range ta.BroadcastPLMNList {
				for _, s := range bp.TAISliceSupportList {
					found.add(models.Snssai{Sst: s.SNSSAI.SST, Sd: s.SNSSAI.SD})
				}
			}
		}
		snssaiMap, err := buildValueMap(found.order, g.Snssais, "snssaiList", "n3iwf supportedTAList")
		if err != nil {
			return err
		}

		const base = "n3iwfInformation"
		if p := g.Plmn; p != nil && info.GlobalN3IWFID != nil {
			if info.GlobalN3IWFID.PLMNID == nil {
				info.GlobalN3IWFID.PLMNID = &n3iwffactory.PLMNID{}
			}
			id := info.GlobalN3IWFID.PLMNID
			Set(r, N3IWF, base+".globalN3IWFID.plmnID.mcc", &id.Mcc, p.Mcc)
			Set(r, N3IWF, base+".globalN3IWFID.plmnID.mnc", &id.Mnc, p.Mnc)
		}
		for i := range info.SupportedTAList {
			ta := &info.SupportedTAList[i]
			taPath := fmt.Sprintf("%s.supportedTAList[%d]", base, i)
			if g.Tac != nil {
				Set(r, N3IWF, taPath+".tac", &ta.TAC, *g.Tac)
			}
			for j := range ta.BroadcastPLMNList {
				bp := &ta.BroadcastPLMNList[j]
				bpPath := fmt.Sprintf("%s.broadcastPlmnList[%d]", taPath, j)
				if p := g.Plmn; p != nil {
					if bp.PLMNID == nil {
						bp.PLMNID = &n3iwffactory.PLMNID{}
					}
					Set(r, N3IWF, bpPath+".plmnID.mcc", &bp.PLMNID.Mcc, p.Mcc)
					Set(r, N3IWF, bpPath+".plmnID.mnc", &bp.PLMNID.Mnc, p.Mnc)
				}
				for k := range bp.TAISliceSupportList {
					s := &bp.TAISliceSupportList[k].SNSSAI
					to, ok := snssaiMap[models.Snssai{Sst: s.SST, Sd: s.SD}]
					if !ok {
						continue
					}
					sPath := fmt.Sprintf("%s.taiSliceSupportList[%d].snssai", bpPath, k)
					Set(r, N3IWF, sPath+".sst", &s.SST, to.Sst)
					Set(r, N3IWF, sPath+".sd", &s.SD, to.Sd)
				}
			}
		}
	}

	f := commonFields{IdleTimeout: &c.IdleTimeout}
	if cfg.Logger == nil && g.Logger != nil {
		cfg.Logger = &n3iwffactory.Logger{}
	}
	if cfg.Logger != nil {
		f.LoggerEnable, f.LoggerLevel, f.LoggerReportCaller =
			&cfg.Logger.Enable, &cfg.Logger.Level, &cfg.Logger.ReportCaller
	}
	applyCommon(g, N3IWF, f, r)
	return nil
}

func (N3IWFAdapter) Validate(cfg *n3iwffactory.Config) error {
	return cfg.Validate()
}
