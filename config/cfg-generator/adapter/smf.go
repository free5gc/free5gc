package adapter

import (
	"errors"
	"fmt"
	"sort"

	"github.com/free5gc/openapi/models"

	"github.com/free5gc/config/cfg-generator/global"
	smffactory "github.com/free5gc/smf/pkg/factory"
)

// SMFAdapter applies global fields to smfcfg.yaml (docs/global-mapping.md §1).
//
//   - plmn:       plmnList[*]
//   - snssaiList: snssaiInfos[*].sNssai, then
//     userplaneInformation.upNodes.<node>.sNssaiUpfInfos[*].sNssai  (value mapping)
//   - dnnList:    snssaiInfos[*].dnnInfos[*].dnn, then
//     upNodes.<node>.sNssaiUpfInfos[*].dnnUpfInfoList[*].dnn          (value mapping);
//     upNodes.<node>.interfaces[*].networkInstances[*] follow the same mapping
//   - nrf, sbiScheme, logger, timeouts: see applyCommon
//
// upNodes is a map; nodes are always visited in sorted name order so that
// value collection and the change log are deterministic.
type SMFAdapter struct{}

var _ Adapter[smffactory.Config] = SMFAdapter{}

func (SMFAdapter) NF() NFName       { return SMF }
func (SMFAdapter) FileName() string { return "smfcfg.yaml" }

func (SMFAdapter) Apply(g *global.Global, cfg *smffactory.Config, r *Recorder) error {
	if g == nil {
		return nil
	}
	c := cfg.Configuration
	if c == nil {
		return errors.New("configuration block is missing")
	}
	nodes := sortedNodeNames(c.UserPlaneInformation.UPNodes)

	snssaiMap, err := buildValueMap(smfSnssais(c, nodes), g.Snssais, "snssaiList",
		"smf snssaiInfos/userplaneInformation")
	if err != nil {
		return err
	}
	dnnMap, err := buildValueMap(smfDnns(c, nodes), g.Dnns, "dnnList",
		"smf snssaiInfos/userplaneInformation")
	if err != nil {
		return err
	}

	if p := g.Plmn; p != nil {
		for i := range c.PLMNList {
			base := fmt.Sprintf("plmnList[%d]", i)
			Set(r, SMF, base+".mcc", &c.PLMNList[i].Mcc, p.Mcc)
			Set(r, SMF, base+".mnc", &c.PLMNList[i].Mnc, p.Mnc)
		}
	}

	for i, info := range c.SNssaiInfo {
		if info == nil {
			continue
		}
		base := fmt.Sprintf("snssaiInfos[%d]", i)
		mapSnssai(r, SMF, base+".sNssai", info.SNssai, snssaiMap)
		for j, dnnInfo := range info.DnnInfos {
			if dnnInfo != nil {
				mapString(r, SMF, fmt.Sprintf("%s.dnnInfos[%d].dnn", base, j), &dnnInfo.Dnn, dnnMap)
			}
		}
	}

	for _, name := range nodes {
		node := c.UserPlaneInformation.UPNodes[name]
		base := "userplaneInformation.upNodes." + name
		for i, s := range node.SNssaiInfos {
			if s == nil {
				continue
			}
			sBase := fmt.Sprintf("%s.sNssaiUpfInfos[%d]", base, i)
			mapSnssai(r, SMF, sBase+".sNssai", s.SNssai, snssaiMap)
			for j, d := range s.DnnUpfInfoList {
				if d != nil {
					mapString(r, SMF, fmt.Sprintf("%s.dnnUpfInfoList[%d].dnn", sBase, j), &d.Dnn, dnnMap)
				}
			}
		}
		for i, itf := range node.InterfaceUpfInfoList {
			if itf == nil {
				continue
			}
			for j := range itf.NetworkInstances {
				mapString(r, SMF, fmt.Sprintf("%s.interfaces[%d].networkInstances[%d]", base, i, j),
					&itf.NetworkInstances[j], dnnMap)
			}
		}
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
		cfg.Logger = &smffactory.Logger{}
	}
	if cfg.Logger != nil {
		f.LoggerEnable = &cfg.Logger.Enable
		f.LoggerLevel = &cfg.Logger.Level
		f.LoggerReportCaller = &cfg.Logger.ReportCaller
	}
	applyCommon(g, SMF, f, r)
	return nil
}

func sortedNodeNames(nodes map[string]*smffactory.UPNode) []string {
	names := make([]string, 0, len(nodes))
	for name, node := range nodes {
		if node != nil {
			names = append(names, name)
		}
	}
	sort.Strings(names)
	return names
}

// smfSnssais returns the distinct S-NSSAIs: snssaiInfos first, then UPF
// sNssaiUpfInfos in sorted node order.
func smfSnssais(c *smffactory.Configuration, nodes []string) []models.Snssai {
	var d distinct[models.Snssai]
	for _, info := range c.SNssaiInfo {
		if info != nil && info.SNssai != nil {
			d.add(*info.SNssai)
		}
	}
	for _, name := range nodes {
		for _, s := range c.UserPlaneInformation.UPNodes[name].SNssaiInfos {
			if s != nil && s.SNssai != nil {
				d.add(*s.SNssai)
			}
		}
	}
	return d.order
}

// smfDnns returns the distinct DNNs from snssaiInfos and UPF dnnUpfInfoList.
// networkInstances are not counted: they usually equal a DNN and are mapped
// only when they do.
func smfDnns(c *smffactory.Configuration, nodes []string) []string {
	var d distinct[string]
	for _, info := range c.SNssaiInfo {
		if info == nil {
			continue
		}
		for _, dnnInfo := range info.DnnInfos {
			if dnnInfo != nil {
				d.add(dnnInfo.Dnn)
			}
		}
	}
	for _, name := range nodes {
		for _, s := range c.UserPlaneInformation.UPNodes[name].SNssaiInfos {
			if s == nil {
				continue
			}
			for _, dnn := range s.DnnUpfInfoList {
				if dnn != nil {
					d.add(dnn.Dnn)
				}
			}
		}
	}
	return d.order
}

func (SMFAdapter) Validate(cfg *smffactory.Config) error {
	if _, err := cfg.Validate(); err != nil {
		return err
	}
	return nil
}
