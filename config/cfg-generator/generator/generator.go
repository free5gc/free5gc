// Package generator runs the config generation pipeline:
// Load -> Global Validate -> Adapter Apply -> NF Validate -> Marshal -> Publish.
package generator

import (
	"fmt"
	"time"

	"github.com/free5gc/config/cfg-generator/adapter"
	"github.com/free5gc/config/cfg-generator/global"
)

// Options configures one generator run.
type Options struct {
	InputPath string
	OutDir    string
	// DryRun stops after validation: changes are reported, nothing is written.
	DryRun bool
	// Strict fails the run with StageStrict (exit code 6) before anything is
	// written when a global value could not be applied to some NF.
	Strict bool
	// Now is used for backup file names. Defaults to time.Now.
	Now func() time.Time
}

// Result describes a run.
type Result struct {
	Changes    []adapter.Change
	Warnings   []adapter.Warning
	NotApplied []adapter.NotApplied
	Files      []string // output file names, in publish order
}

// Run executes the pipeline. Each stage collects all of its errors before
// returning; a failed stage stops the pipeline. Nothing is written unless
// every stage up to publish succeeds.
//
// A StageStrict error is the only error returned together with a non-nil
// Result, so that the caller can report what was not applied.
func Run(opts Options) (*Result, error) {
	if opts.Now == nil {
		opts.Now = time.Now
	}

	in, err := Load(opts.InputPath)
	if err != nil {
		return nil, err
	}

	if errs := global.Validate(in.Global); len(errs) > 0 {
		return nil, stageErr(StageGlobal, errs...)
	}

	units := in.Units()
	res := &Result{}
	rec := &adapter.Recorder{}

	var errs []error
	for _, u := range units {
		if err = u.Apply(in.Global, rec); err != nil {
			errs = append(errs, fmt.Errorf("%s: %w", u.NF, err))
		}
	}
	if len(errs) > 0 {
		return nil, stageErr(StageAdapter, errs...)
	}
	res.Changes = rec.Changes
	res.Warnings = rec.Warnings
	res.NotApplied = rec.NotApplied

	for _, u := range units {
		if err = u.Validate(); err != nil {
			errs = append(errs, fmt.Errorf("%s: %w", u.NF, err))
		}
	}
	if len(errs) > 0 {
		return nil, stageErr(StageValidate, errs...)
	}

	files := make([]File, 0, len(units))
	for _, u := range units {
		data, mErr := Marshal(u.Cfg)
		if mErr != nil {
			errs = append(errs, fmt.Errorf("%s: %w", u.NF, mErr))
			continue
		}
		files = append(files, File{Name: u.FileName, Data: data})
		res.Files = append(res.Files, u.FileName)
	}
	if len(errs) > 0 {
		return nil, stageErr(StagePublish, errs...)
	}

	if opts.Strict && len(res.NotApplied) > 0 {
		return res, stageErr(StageStrict, fmt.Errorf(
			"%d global value(s) not applied; nothing written (run without --strict to accept)",
			len(res.NotApplied)))
	}

	if opts.DryRun {
		return res, nil
	}
	if err = Publish(opts.OutDir, files, opts.Now()); err != nil {
		return nil, stageErr(StagePublish, err)
	}
	return res, nil
}
