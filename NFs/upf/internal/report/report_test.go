package report_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/free5gc/go-upf/internal/report"
)

func TestApplyAction0(t *testing.T) {
	var act report.ApplyAction
	e := act.Unmarshal([]byte{})
	assert.Error(t, e)
}

func TestApplyAction1(t *testing.T) {
	var act report.ApplyAction
	e := act.Unmarshal([]byte{0x02})
	assert.NoError(t, e)
	assert.Equal(t, uint16(0x0002), act.Flags)
	assert.False(t, act.DROP())
	assert.True(t, act.FORW())
	assert.False(t, act.BUFF())
	assert.False(t, act.NOCP())
	assert.False(t, act.DUPL())
	assert.False(t, act.IPMA())
	assert.False(t, act.IPMD())
	assert.False(t, act.DFRT())
	assert.False(t, act.EDRT())
	assert.False(t, act.BDPN())
	assert.False(t, act.DDPN())
	assert.False(t, act.FSSM())
	assert.False(t, act.MBSU())
}

func TestApplyAction2(t *testing.T) {
	var act report.ApplyAction
	e := act.Unmarshal([]byte{0x0C, 0x00})
	assert.NoError(t, e)
	assert.Equal(t, uint16(0x000C), act.Flags)
	assert.False(t, act.DROP())
	assert.False(t, act.FORW())
	assert.True(t, act.BUFF())
	assert.True(t, act.NOCP())
	assert.False(t, act.DUPL())
	assert.False(t, act.IPMA())
	assert.False(t, act.IPMD())
	assert.False(t, act.DFRT())
	assert.False(t, act.EDRT())
	assert.False(t, act.BDPN())
	assert.False(t, act.DDPN())
	assert.False(t, act.FSSM())
	assert.False(t, act.MBSU())
}
