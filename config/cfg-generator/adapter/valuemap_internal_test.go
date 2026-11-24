package adapter

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/free5gc/openapi/models"
)

func TestDistinctKeepsFirstAppearanceOrder(t *testing.T) {
	var d distinct[string]
	for _, v := range []string{"b", "a", "b", "c", "a"} {
		d.add(v)
	}
	assert.Equal(t, []string{"b", "a", "c"}, d.order)

	var empty distinct[string]
	assert.Nil(t, empty.order)
}

func TestBuildValueMap(t *testing.T) {
	t.Run("nil global means not provided", func(t *testing.T) {
		m, err := buildValueMap([]string{"a"}, nil, "dnnList", "nf x")
		require.NoError(t, err)
		assert.Nil(t, m)
	})

	t.Run("maps by position", func(t *testing.T) {
		m, err := buildValueMap([]string{"a", "b"}, []string{"x", "y"}, "dnnList", "nf x")
		require.NoError(t, err)
		assert.Equal(t, map[string]string{"a": "x", "b": "y"}, m)
	})

	t.Run("identity and swap are allowed", func(t *testing.T) {
		m, err := buildValueMap([]string{"a", "b"}, []string{"b", "a"}, "dnnList", "nf x")
		require.NoError(t, err)
		assert.Equal(t, map[string]string{"a": "b", "b": "a"}, m)
	})

	t.Run("count mismatch", func(t *testing.T) {
		_, err := buildValueMap([]string{"a", "b"}, []string{"x"}, "dnnList", "nf x")
		require.Error(t, err)
		assert.Equal(t, "global dnnList has 1 entries but nf x has 2 distinct value(s) [a b]", err.Error())
	})

	t.Run("NF has no values but global has some", func(t *testing.T) {
		_, err := buildValueMap(nil, []string{"x"}, "dnnList", "nf x")
		assert.Error(t, err)
	})
}

func TestMapSnssai(t *testing.T) {
	m := map[models.Snssai]models.Snssai{
		{Sst: 1, Sd: "010203"}: {Sst: 1, Sd: "000001"}, // sd only
		{Sst: 1, Sd: "112233"}: {Sst: 2, Sd: "112233"}, // sst only
	}

	r := &Recorder{}
	a := &models.Snssai{Sst: 1, Sd: "010203"}
	mapSnssai(r, AMF, "p0", a, m)
	b := &models.Snssai{Sst: 1, Sd: "112233"}
	mapSnssai(r, AMF, "p1", b, m)
	unmapped := &models.Snssai{Sst: 9, Sd: "ffffff"}
	mapSnssai(r, AMF, "p2", unmapped, m)
	mapSnssai(r, AMF, "p3", nil, m)
	mapSnssai(r, AMF, "p4", &models.Snssai{Sst: 1, Sd: "010203"}, nil) // field not provided

	assert.Equal(t, models.Snssai{Sst: 1, Sd: "000001"}, *a)
	assert.Equal(t, models.Snssai{Sst: 2, Sd: "112233"}, *b)
	assert.Equal(t, models.Snssai{Sst: 9, Sd: "ffffff"}, *unmapped, "unmapped value kept")
	assert.Equal(t, []Change{
		{NF: AMF, Path: "p0.sd", Old: "010203", New: "000001"},
		{NF: AMF, Path: "p1.sst", Old: "1", New: "2"},
	}, r.Changes, "only fields that actually change are recorded")
}

func TestMapString(t *testing.T) {
	r := &Recorder{}
	m := map[string]string{"internet": "ims"}
	a, b := "internet", "other"
	mapString(r, SMF, "a", &a, m)
	mapString(r, SMF, "b", &b, m)
	c := "internet"
	mapString(r, SMF, "c", &c, nil)

	assert.Equal(t, "ims", a)
	assert.Equal(t, "other", b)
	assert.Equal(t, "internet", c)
	assert.Equal(t, []Change{{NF: SMF, Path: "a", Old: "internet", New: "ims"}}, r.Changes)
}

func TestSetDuration(t *testing.T) {
	t.Run("unset to set", func(t *testing.T) {
		r := &Recorder{}
		var d *time.Duration
		setDuration(r, AMF, "idleTimeout", &d, 60*time.Second)
		require.NotNil(t, d)
		assert.Equal(t, 60*time.Second, *d)
		assert.Equal(t, []Change{{NF: AMF, Path: "idleTimeout", Old: "<unset>", New: "1m0s"}}, r.Changes)
	})

	t.Run("same value not recorded and pointer kept", func(t *testing.T) {
		r := &Recorder{}
		v := 15 * time.Second
		d := &v
		setDuration(r, AMF, "pingTimeout", &d, 15*time.Second)
		assert.Same(t, &v, d)
		assert.Empty(t, r.Changes)
	})

	t.Run("change replaces pointer, never writes through it", func(t *testing.T) {
		r := &Recorder{}
		shared := 15 * time.Second
		d := &shared
		setDuration(r, AMF, "pingTimeout", &d, 30*time.Second)
		assert.Equal(t, 30*time.Second, *d)
		assert.Equal(t, 15*time.Second, shared, "the old pointee is not modified")
		assert.Equal(t, []Change{{NF: AMF, Path: "pingTimeout", Old: "15s", New: "30s"}}, r.Changes)
	})
}
