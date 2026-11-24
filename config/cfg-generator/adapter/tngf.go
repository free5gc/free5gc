package adapter

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/free5gc/openapi/models"

	"github.com/free5gc/config/cfg-generator/global"
	tngffactory "github.com/free5gc/tngf/pkg/factory"
)

// TNGFAdapter applies global fields to tngfcfg.yaml (docs/global-mapping.md §1).
// TNGF uses capitalised YAML keys.
//
//   - plmn: TNGFInformation.GlobalTNGFID.PLMNID,
//     TNGFInformation.SupportedTAList[*].BroadcastPLMNList[*].PLMNID
//   - tac:  TNGFInformation.SupportedTAList[*].TAC
//   - snssaiList: ...BroadcastPLMNList[*].TAISliceSupportList[*].SNSSAI (value
//     mapping). TNGF stores SST as a single hex digit, so only SST 0..15 can be
//     represented; a larger global SST is an error.
//   - logger, timeouts (idleTimeout only): see applyCommon
//
// The tngf factory package only builds on Linux.
type TNGFAdapter struct{}

var _ Adapter[tngffactory.Config] = TNGFAdapter{}

func (TNGFAdapter) NF() NFName       { return TNGF }
func (TNGFAdapter) FileName() string { return "tngfcfg.yaml" }

// tngfSst parses TNGF's single-hex-digit SST.
func tngfSst(s string) (int32, error) {
	v, err := strconv.ParseInt(s, 16, 32)
	if err != nil || len(s) != 1 {
		return 0, fmt.Errorf("invalid TNGF SST %q: must be a single hex digit", s)
	}
	return int32(v), nil
}

// formatTngfSst formats sst as TNGF's single uppercase hex digit.
func formatTngfSst(sst int32) (string, error) {
	if sst < 0 || sst > 15 {
		return "", fmt.Errorf("global SST %d cannot be represented by TNGF, which stores SST as one hex digit (0-15)", sst)
	}
	return strings.ToUpper(strconv.FormatInt(int64(sst), 16)), nil
}

func (TNGFAdapter) Apply(g *global.Global, cfg *tngffactory.Config, r *Recorder) error {
	if g == nil {
		return nil
	}
	c := cfg.Configuration
	if c == nil {
		return errMissingConfiguration
	}
	info := &c.TNGFInfo

	// Collect and validate everything before changing anything.
	var found distinct[models.Snssai]
	for _, ta := range info.SupportedTAList {
		for _, bp := range ta.BroadcastPLMNList {
			for _, s := range bp.TAISliceSupportList {
				sst, err := tngfSst(s.SNSSAI.SST)
				if err != nil {
					return err
				}
				found.add(models.Snssai{Sst: sst, Sd: s.SNSSAI.SD})
			}
		}
	}
	snssaiMap, err := buildValueMap(found.order, g.Snssais, "snssaiList", "tngf SupportedTAList")
	if err != nil {
		return err
	}
	sstText := map[int32]string{}
	for _, to := range snssaiMap {
		if sstText[to.Sst], err = formatTngfSst(to.Sst); err != nil {
			return err
		}
	}

	const base = "TNGFInformation"
	if p := g.Plmn; p != nil {
		id := &info.GlobalTNGFID.PLMNID
		Set(r, TNGF, base+".GlobalTNGFID.PLMNID.MCC", &id.Mcc, p.Mcc)
		Set(r, TNGF, base+".GlobalTNGFID.PLMNID.MNC", &id.Mnc, p.Mnc)
	}
	for i := range info.SupportedTAList {
		ta := &info.SupportedTAList[i]
		taPath := fmt.Sprintf("%s.SupportedTAList[%d]", base, i)
		if g.Tac != nil {
			Set(r, TNGF, taPath+".TAC", &ta.TAC, *g.Tac)
		}
		for j := range ta.BroadcastPLMNList {
			bp := &ta.BroadcastPLMNList[j]
			bpPath := fmt.Sprintf("%s.BroadcastPLMNList[%d]", taPath, j)
			if p := g.Plmn; p != nil {
				Set(r, TNGF, bpPath+".PLMNID.MCC", &bp.PLMNID.Mcc, p.Mcc)
				Set(r, TNGF, bpPath+".PLMNID.MNC", &bp.PLMNID.Mnc, p.Mnc)
			}
			for k := range bp.TAISliceSupportList {
				s := &bp.TAISliceSupportList[k].SNSSAI
				sst, _ := tngfSst(s.SST) // validated above
				to, ok := snssaiMap[models.Snssai{Sst: sst, Sd: s.SD}]
				if !ok {
					continue
				}
				sPath := fmt.Sprintf("%s.TAISliceSupportList[%d].SNSSAI", bpPath, k)
				if to.Sst != sst {
					Set(r, TNGF, sPath+".SST", &s.SST, sstText[to.Sst])
				}
				Set(r, TNGF, sPath+".SD", &s.SD, to.Sd)
			}
		}
	}

	f := commonFields{IdleTimeout: &c.IdleTimeout}
	if cfg.Logger == nil && g.Logger != nil {
		cfg.Logger = &tngffactory.Logger{}
	}
	if cfg.Logger != nil {
		f.LoggerEnable, f.LoggerLevel, f.LoggerReportCaller =
			&cfg.Logger.Enable, &cfg.Logger.Level, &cfg.Logger.ReportCaller
	}
	applyCommon(g, TNGF, f, r)
	return nil
}

// Validate runs TNGF's own Validate and checks the version, because TNGF's
// CheckConfigVersion has no caller in the TNGF repository
// (docs/nf-config-defaults.md).
func (TNGFAdapter) Validate(cfg *tngffactory.Config) error {
	var errs []error
	if _, err := cfg.Validate(); err != nil {
		errs = append(errs, err)
	}
	if cfg.Info == nil || cfg.Info.Version != tngffactory.TngfExpectedConfigVersion {
		errs = append(errs, fmt.Errorf("info.version must be %q", tngffactory.TngfExpectedConfigVersion))
	}
	return errors.Join(errs...)
}
