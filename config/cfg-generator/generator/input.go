package generator

import (
	"fmt"
	"os"
	"sort"

	"gopkg.in/yaml.v2"

	amffactory "github.com/free5gc/amf/pkg/factory"
	ausffactory "github.com/free5gc/ausf/pkg/factory"
	bsffactory "github.com/free5gc/bsf/pkg/factory"
	chffactory "github.com/free5gc/chf/pkg/factory"
	"github.com/free5gc/config/cfg-generator/adapter"
	"github.com/free5gc/config/cfg-generator/global"
	upffactory "github.com/free5gc/go-upf/pkg/factory"
	n3iwffactory "github.com/free5gc/n3iwf/pkg/factory"
	neffactory "github.com/free5gc/nef/pkg/factory"
	nrffactory "github.com/free5gc/nrf/pkg/factory"
	nssffactory "github.com/free5gc/nssf/pkg/factory"
	pcffactory "github.com/free5gc/pcf/pkg/factory"
	scpfactory "github.com/free5gc/scp/pkg/factory"
	smffactory "github.com/free5gc/smf/pkg/factory"
	tngffactory "github.com/free5gc/tngf/pkg/factory"
	udmfactory "github.com/free5gc/udm/pkg/factory"
	udrfactory "github.com/free5gc/udr/pkg/factory"
)

// Input is the integrated input file: one YAML document whose top-level keys
// are `global` and one key per NF. NF blocks are decoded straight into each
// NF's own factory.Config, which keeps string fields such as sd "010203" and
// mnc "01" intact (a generic map would turn them into 4227 and 1).
//
// Every NF field is a pointer: an absent NF is nil and gets no output file.
// webui is not supported yet: the webconsole submodule is not checked out.
type Input struct {
	Global *global.Global       `yaml:"global"`
	UPF    *upffactory.Config   `yaml:"upf,omitempty"`
	NRF    *nrffactory.Config   `yaml:"nrf,omitempty"`
	AMF    *amffactory.Config   `yaml:"amf,omitempty"`
	SMF    *smffactory.Config   `yaml:"smf,omitempty"`
	UDR    *udrfactory.Config   `yaml:"udr,omitempty"`
	PCF    *pcffactory.Config   `yaml:"pcf,omitempty"`
	UDM    *udmfactory.Config   `yaml:"udm,omitempty"`
	NSSF   *nssffactory.Config  `yaml:"nssf,omitempty"`
	AUSF   *ausffactory.Config  `yaml:"ausf,omitempty"`
	CHF    *chffactory.Config   `yaml:"chf,omitempty"`
	NEF    *neffactory.Config   `yaml:"nef,omitempty"`
	SCP    *scpfactory.Config   `yaml:"scp,omitempty"`
	BSF    *bsffactory.Config   `yaml:"bsf,omitempty"`
	N3IWF  *n3iwffactory.Config `yaml:"n3iwf,omitempty"`
	TNGF   *tngffactory.Config  `yaml:"tngf,omitempty"`
}

// supportedNFs lists the top-level keys accepted besides `global`.
var supportedNFs = map[string]bool{
	string(adapter.UPF): true, string(adapter.NRF): true, string(adapter.AMF): true,
	string(adapter.SMF): true, string(adapter.UDR): true, string(adapter.PCF): true,
	string(adapter.UDM): true, string(adapter.NSSF): true, string(adapter.AUSF): true,
	string(adapter.CHF): true, string(adapter.NEF): true, string(adapter.SCP): true,
	string(adapter.BSF): true, string(adapter.N3IWF): true, string(adapter.TNGF): true,
}

// Load reads and decodes the integrated input file.
//
// It decodes twice:
//  1. strictly, to reject unknown keys inside `global` and unknown top-level
//     keys (yaml.v2 non-strict decoding would silently ignore a typo);
//  2. leniently with yaml.v2, exactly as the NFs' own loaders decode their
//     files, so NF blocks accept whatever the NFs accept.
//
// Step 2 deliberately does not use free-ran-ue/util.LoadFromYaml: on a decode
// error it returns before closing the file (util v0.2.1, yaml.go), which
// leaks the handle and, on Windows, keeps the file locked. It also pulls in
// gin, quic-go and the MongoDB driver as transitive dependencies. The decode
// semantics are identical (both are yaml.v2 non-strict).
func Load(path string) (*Input, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, stageErr(StageInput, err)
	}

	var probe struct {
		Global *global.Global         `yaml:"global"`
		Rest   map[string]interface{} `yaml:",inline"`
	}
	if err = yaml.UnmarshalStrict(data, &probe); err != nil {
		return nil, stageErr(StageInput, fmt.Errorf("%s: %w", path, err))
	}

	var unknown []error
	keys := make([]string, 0, len(probe.Rest))
	for k := range probe.Rest {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		if !supportedNFs[k] {
			unknown = append(unknown, fmt.Errorf("%s: unsupported top-level key %q", path, k))
		}
	}
	if len(unknown) > 0 {
		return nil, stageErr(StageInput, unknown...)
	}

	in := &Input{}
	if err = yaml.Unmarshal(data, in); err != nil {
		return nil, stageErr(StageInput, fmt.Errorf("%s: %w", path, err))
	}
	return in, nil
}

// Units binds every present NF block to its adapter, in run.sh start order.
func (in *Input) Units() []*adapter.Unit {
	var units []*adapter.Unit
	for _, u := range []*adapter.Unit{
		adapter.Bind[upffactory.Config](adapter.UPFAdapter{}, in.UPF),
		adapter.Bind[nrffactory.Config](adapter.NRFAdapter{}, in.NRF),
		adapter.Bind[amffactory.Config](adapter.AMFAdapter{}, in.AMF),
		adapter.Bind[smffactory.Config](adapter.SMFAdapter{}, in.SMF),
		adapter.Bind[udrfactory.Config](adapter.UDRAdapter{}, in.UDR),
		adapter.Bind[pcffactory.Config](adapter.PCFAdapter{}, in.PCF),
		adapter.Bind[udmfactory.Config](adapter.UDMAdapter{}, in.UDM),
		adapter.Bind[nssffactory.Config](adapter.NSSFAdapter{}, in.NSSF),
		adapter.Bind[ausffactory.Config](adapter.AUSFAdapter{}, in.AUSF),
		adapter.Bind[chffactory.Config](adapter.CHFAdapter{}, in.CHF),
		adapter.Bind[neffactory.Config](adapter.NEFAdapter{}, in.NEF),
		adapter.Bind[scpfactory.Config](adapter.SCPAdapter{}, in.SCP),
		adapter.Bind[bsffactory.Config](adapter.BSFAdapter{}, in.BSF),
		adapter.Bind[n3iwffactory.Config](adapter.N3IWFAdapter{}, in.N3IWF),
		adapter.Bind[tngffactory.Config](adapter.TNGFAdapter{}, in.TNGF),
	} {
		if u != nil {
			units = append(units, u)
		}
	}
	return units
}
