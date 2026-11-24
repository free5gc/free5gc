// Package adapter applies the global configuration to individual NF configs.
//
// Each NF has a typed Adapter that knows exactly which fields of that NF's
// factory.Config correspond to which global fields. Adapters are bound to a
// concrete config value via Bind, producing a type-erased Unit that the
// generator can iterate over uniformly.
package adapter

import (
	"fmt"

	"gopkg.in/yaml.v2"

	"github.com/free5gc/config/cfg-generator/global"
)

// NFName identifies an NF and is also its key in the integrated input file.
type NFName string

const (
	AMF   NFName = "amf"
	AUSF  NFName = "ausf"
	BSF   NFName = "bsf"
	CHF   NFName = "chf"
	N3IWF NFName = "n3iwf"
	NEF   NFName = "nef"
	NRF   NFName = "nrf"
	NSSF  NFName = "nssf"
	PCF   NFName = "pcf"
	SCP   NFName = "scp"
	SMF   NFName = "smf"
	TNGF  NFName = "tngf"
	UDM   NFName = "udm"
	UDR   NFName = "udr"
	UPF   NFName = "upf"
)

// Adapter maps the global configuration onto one NF's config type T.
type Adapter[T any] interface {
	NF() NFName
	// FileName is the output file name, e.g. "amfcfg.yaml".
	FileName() string
	// Apply writes the provided global values into cfg and records every
	// field it actually changes. It must not touch fields the global block
	// does not provide.
	Apply(g *global.Global, cfg *T, r *Recorder) error
	// Validate checks that cfg is acceptable to the NF.
	Validate(cfg *T) error
}

// Unit is an Adapter bound to a concrete config, with the type erased.
type Unit struct {
	NF       NFName
	FileName string
	// Cfg is the bound config, e.g. *amffactory.Config. It is what gets
	// serialized to the output file.
	Cfg      any
	apply    func(*global.Global, *Recorder) error
	validate func() error
}

// Apply runs the bound adapter's Apply.
func (u *Unit) Apply(g *global.Global, r *Recorder) error { return u.apply(g, r) }

// Validate runs the bound adapter's Validate.
func (u *Unit) Validate() error { return u.validate() }

// Bind binds an adapter to cfg. It returns nil when cfg is nil, which means
// the NF is absent from the input and no file is generated for it.
func Bind[T any](a Adapter[T], cfg *T) *Unit {
	if cfg == nil {
		return nil
	}
	return &Unit{
		NF:       a.NF(),
		FileName: a.FileName(),
		Cfg:      cfg,
		apply:    func(g *global.Global, r *Recorder) error { return a.Apply(g, cfg, r) },
		validate: func() error { return validateCopy(a, cfg) },
	}
}

// validateCopy validates a copy of cfg decoded from its YAML serialization.
//
// NF Validate methods have side effects: AMF and SMF fill an empty
// nfInstanceId with a random UUID. Validating cfg itself would leak that
// UUID into the output file, freezing an ID that the NF is meant to generate
// at every start. Decoding the serialized bytes also means we validate
// exactly what the NF will read from disk.
func validateCopy[T any](a Adapter[T], cfg *T) error {
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("serialize for validation: %w", err)
	}
	cp := new(T)
	if err = yaml.Unmarshal(data, cp); err != nil {
		return fmt.Errorf("decode for validation: %w", err)
	}
	return a.Validate(cp)
}

// Change is one field modification made by an adapter.
type Change struct {
	NF   NFName
	Path string
	Old  string
	New  string
}

func (c Change) String() string {
	return fmt.Sprintf("%s: %s: %q -> %q", c.NF, c.Path, c.Old, c.New)
}

// Warning is a generic, free-form note about an NF config. Warnings never
// fail a run. A global value that an NF cannot use is a NotApplied instead.
type Warning struct {
	NF  NFName
	Msg string
}

func (w Warning) String() string { return fmt.Sprintf("%s: %s", w.NF, w.Msg) }

// NotApplied reports a global value that cannot reach an NF because the NF
// hard-codes the corresponding setting. The run still succeeds unless the
// CLI is started with --strict.
type NotApplied struct {
	NF NFName
	// Item names what could not be applied, e.g. "NRF registration PLMN".
	Item string
	// HardCoded is the value the NF uses regardless of its config.
	HardCoded string
	// Requested is the global value that was not applied.
	Requested string
}

func (n NotApplied) String() string {
	return fmt.Sprintf("%s: %s stays %s (hard-coded); global %s not applied",
		n.NF, n.Item, n.HardCoded, n.Requested)
}

// Recorder collects the changes adapters make. Tests use it to assert that
// exactly the expected paths were modified; the CLI prints it as a summary.
type Recorder struct {
	Changes    []Change
	Warnings   []Warning
	NotApplied []NotApplied
}

// Warn records a warning for nf.
func (r *Recorder) Warn(nf NFName, format string, args ...any) {
	r.Warnings = append(r.Warnings, Warning{NF: nf, Msg: fmt.Sprintf(format, args...)})
}

// NotApply records that the global value requested for item cannot reach nf
// because nf hard-codes hardCoded.
func (r *Recorder) NotApply(nf NFName, item, hardCoded, requested string) {
	r.NotApplied = append(r.NotApplied, NotApplied{
		NF: nf, Item: item, HardCoded: hardCoded, Requested: requested,
	})
}

// Set assigns v to *dst and records the change if the value differs.
// Unchanged assignments are not recorded.
func Set[V comparable](r *Recorder, nf NFName, path string, dst *V, v V) {
	if *dst == v {
		return
	}
	r.Changes = append(r.Changes, Change{
		NF:   nf,
		Path: path,
		Old:  fmt.Sprint(*dst),
		New:  fmt.Sprint(v),
	})
	*dst = v
}
