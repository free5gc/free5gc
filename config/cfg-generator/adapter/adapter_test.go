package adapter_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/free5gc/config/cfg-generator/adapter"
	"github.com/free5gc/config/cfg-generator/global"
)

func TestSetRecordsOnlyActualChanges(t *testing.T) {
	r := &adapter.Recorder{}
	v := "208"

	adapter.Set(r, adapter.AMF, "x.mcc", &v, "208")
	assert.Empty(t, r.Changes, "same value must not be recorded")
	assert.Equal(t, "208", v)

	adapter.Set(r, adapter.AMF, "x.mcc", &v, "466")
	require.Len(t, r.Changes, 1)
	assert.Equal(t, adapter.Change{NF: adapter.AMF, Path: "x.mcc", Old: "208", New: "466"}, r.Changes[0])
	assert.Equal(t, "466", v)
	assert.Equal(t, `amf: x.mcc: "208" -> "466"`, r.Changes[0].String())
}

type fakeCfg struct {
	Value string `yaml:"value"`
	ID    string `yaml:"id,omitempty"`
}

type fakeAdapter struct {
	validateErr error
	seen        *fakeCfg // what Validate received
}

func (fakeAdapter) NF() adapter.NFName { return "fake" }
func (fakeAdapter) FileName() string   { return "fakecfg.yaml" }
func (fakeAdapter) Apply(_ *global.Global, cfg *fakeCfg, r *adapter.Recorder) error {
	adapter.Set(r, "fake", "value", &cfg.Value, "applied")
	return nil
}

// Validate mimics AMF/SMF: it fills an empty ID as a side effect.
func (a fakeAdapter) Validate(cfg *fakeCfg) error {
	if cfg.ID == "" {
		cfg.ID = "random-uuid"
	}
	if a.seen != nil {
		*a.seen = *cfg
	}
	return a.validateErr
}

func TestBind(t *testing.T) {
	t.Run("nil config means NF absent", func(t *testing.T) {
		assert.Nil(t, adapter.Bind[fakeCfg](fakeAdapter{}, nil))
	})

	t.Run("bound unit forwards to adapter", func(t *testing.T) {
		cfg := &fakeCfg{Value: "orig"}
		wantErr := errors.New("invalid")
		u := adapter.Bind[fakeCfg](fakeAdapter{validateErr: wantErr}, cfg)
		require.NotNil(t, u)

		assert.Equal(t, adapter.NFName("fake"), u.NF)
		assert.Equal(t, "fakecfg.yaml", u.FileName)
		assert.Same(t, cfg, u.Cfg)

		r := &adapter.Recorder{}
		require.NoError(t, u.Apply(&global.Global{}, r))
		assert.Equal(t, "applied", cfg.Value)
		assert.Len(t, r.Changes, 1)

		assert.ErrorIs(t, u.Validate(), wantErr)
	})

	t.Run("validate side effects do not reach the bound config", func(t *testing.T) {
		cfg := &fakeCfg{Value: "orig"}
		seen := &fakeCfg{}
		u := adapter.Bind[fakeCfg](fakeAdapter{seen: seen}, cfg)

		require.NoError(t, u.Validate())
		assert.Equal(t, "orig", seen.Value, "Validate sees the same values")
		assert.Equal(t, "random-uuid", seen.ID, "Validate ran its side effect on the copy")
		assert.Empty(t, cfg.ID, "the bound config, which gets written, is untouched")
	})
}
