package adapter

import (
	"fmt"

	"github.com/asaskevich/govalidator"

	"github.com/free5gc/config/cfg-generator/global"
	upffactory "github.com/free5gc/go-upf/pkg/factory"
)

// UPFAdapter applies global fields to upfcfg.yaml (docs/global-mapping.md §1).
// UPF's config has no info/configuration wrapper; paths are top-level.
//
//   - dnnList: dnnList[*].dnn (value mapping). UPF itself only uses the cidr;
//     the name is kept consistent with SMF.
//   - logger:  logger.*
//
// UPF has no SBI, NRF, timeouts or database fields.
type UPFAdapter struct{}

var _ Adapter[upffactory.Config] = UPFAdapter{}

func (UPFAdapter) NF() NFName       { return UPF }
func (UPFAdapter) FileName() string { return "upfcfg.yaml" }

func (UPFAdapter) Apply(g *global.Global, cfg *upffactory.Config, r *Recorder) error {
	if g == nil {
		return nil
	}

	var found distinct[string]
	for _, d := range cfg.DnnList {
		found.add(d.Dnn)
	}
	dnnMap, err := buildValueMap(found.order, g.Dnns, "dnnList", "upf dnnList")
	if err != nil {
		return err
	}
	for i := range cfg.DnnList {
		mapString(r, UPF, fmt.Sprintf("dnnList[%d].dnn", i), &cfg.DnnList[i].Dnn, dnnMap)
	}

	var f commonFields
	if cfg.Logger == nil && g.Logger != nil {
		cfg.Logger = &upffactory.Logger{}
	}
	if cfg.Logger != nil {
		f.LoggerEnable, f.LoggerLevel, f.LoggerReportCaller =
			&cfg.Logger.Enable, &cfg.Logger.Level, &cfg.Logger.ReportCaller
	}
	applyCommon(g, UPF, f, r)
	return nil
}

// Validate mirrors upf/pkg/factory.ReadConfig: UPF has no Validate method and
// validates with govalidator plus a custom "cidr" tag inside ReadConfig. The
// NodeID DNS resolution done there is skipped; it depends on the host.
func (UPFAdapter) Validate(cfg *upffactory.Config) error {
	govalidator.TagMap["cidr"] = govalidator.Validator(govalidator.IsCIDR)
	_, err := govalidator.ValidateStruct(cfg)
	return err
}
