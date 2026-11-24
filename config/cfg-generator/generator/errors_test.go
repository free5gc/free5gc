package generator_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/free5gc/config/cfg-generator/generator"
)

func TestStageString(t *testing.T) {
	assert.Equal(t, "input", generator.StageInput.String())
	assert.Equal(t, "global", generator.StageGlobal.String())
	assert.Equal(t, "adapter", generator.StageAdapter.String())
	assert.Equal(t, "validate", generator.StageValidate.String())
	assert.Equal(t, "publish", generator.StagePublish.String())
	assert.Equal(t, "strict", generator.StageStrict.String())
	assert.Equal(t, "unknown", generator.Stage(99).String())
}

func TestExitCode(t *testing.T) {
	assert.Equal(t, 0, generator.ExitCode(nil))
	assert.Equal(t, 1, generator.ExitCode(errors.New("plain error")))
	assert.Equal(t, 3, generator.ExitCode(&generator.StageError{Stage: generator.StageAdapter}))
}

func TestStageErrorFormatAndUnwrap(t *testing.T) {
	e1, e2 := errors.New("first"), errors.New("second")
	se := &generator.StageError{Stage: generator.StageGlobal, Errs: []error{e1, e2}}

	assert.Equal(t, "[global] first\n[global] second", se.Error())
	assert.ErrorIs(t, se, e1)
	assert.ErrorIs(t, se, e2)
}
