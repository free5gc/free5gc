package generator

import (
	"errors"
	"strings"
)

// Stage identifies the pipeline stage that failed. Each stage maps to a
// distinct process exit code (docs/adapter-data-model.md §8).
type Stage int

const (
	StageInput    Stage = iota + 1 // 1: arguments, file read, YAML syntax
	StageGlobal                    // 2: global block validation
	StageAdapter                   // 3: applying global to NF configs
	StageValidate                  // 4: NF config validation
	StagePublish                   // 5: writing or publishing output
	// StageStrict (6): with Options.Strict, at least one global value could
	// not be applied (adapter.NotApplied). It is checked after NF validation
	// and before publishing, so nothing is written. Its code is 6 rather than
	// its position in the pipeline so that codes 1-5 keep their meaning.
	StageStrict
)

func (s Stage) String() string {
	switch s {
	case StageInput:
		return "input"
	case StageGlobal:
		return "global"
	case StageAdapter:
		return "adapter"
	case StageValidate:
		return "validate"
	case StagePublish:
		return "publish"
	case StageStrict:
		return "strict"
	default:
		return "unknown"
	}
}

// ExitCode returns the process exit code for the stage.
func (s Stage) ExitCode() int { return int(s) }

// StageError carries every error collected in one stage.
type StageError struct {
	Stage Stage
	Errs  []error
}

func (e *StageError) Error() string {
	lines := make([]string, 0, len(e.Errs))
	for _, err := range e.Errs {
		lines = append(lines, "["+e.Stage.String()+"] "+err.Error())
	}
	return strings.Join(lines, "\n")
}

func (e *StageError) Unwrap() []error { return e.Errs }

func stageErr(s Stage, errs ...error) *StageError {
	return &StageError{Stage: s, Errs: errs}
}

// ExitCode returns the exit code for err: 0 for nil, the stage code for a
// StageError, and 1 for anything else.
func ExitCode(err error) int {
	if err == nil {
		return 0
	}
	var se *StageError
	if errors.As(err, &se) {
		return se.Stage.ExitCode()
	}
	return StageInput.ExitCode()
}
