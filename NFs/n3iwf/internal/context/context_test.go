package context_test

import (
	"context"
	"net"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"

	n3iwf_context "github.com/free5gc/n3iwf/internal/context"
	"github.com/free5gc/n3iwf/pkg/factory"
	"github.com/free5gc/util/ippool"
)

type n3iwfTestApp struct {
	cfg      *factory.Config
	n3iwfCtx *n3iwf_context.N3IWFContext
	ctx      context.Context
	cancel   context.CancelFunc
	wg       *sync.WaitGroup
}

func (a *n3iwfTestApp) Config() *factory.Config {
	return a.cfg
}

func (a *n3iwfTestApp) Context() *n3iwf_context.N3IWFContext {
	return a.n3iwfCtx
}

func (a *n3iwfTestApp) CancelContext() context.Context {
	return a.ctx
}

func NewN3iwfTestApp(cfg *factory.Config) (*n3iwfTestApp, error) {
	var err error
	ctx, cancel := context.WithCancel(context.Background())

	n3iwfApp := &n3iwfTestApp{
		cfg:    cfg,
		ctx:    ctx,
		cancel: cancel,
		wg:     &sync.WaitGroup{},
	}

	n3iwfApp.n3iwfCtx, err = n3iwf_context.NewTestContext(n3iwfApp)
	if err != nil {
		return nil, err
	}
	return n3iwfApp, err
}

func NewTestCfg() *factory.Config {
	return &factory.Config{
		Configuration: &factory.Configuration{
			IPSecGatewayAddr: "10.0.0.1",
			UEIPAddressRange: "10.0.0.0/24",
		},
	}
}

func TestNewInternalUEIPAddr(t *testing.T) {
	cfg := NewTestCfg()
	var app *n3iwfTestApp
	var err error
	var ip, invalidIP, invalidIP2 net.IP

	app, err = NewN3iwfTestApp(cfg)
	require.NoError(t, err)

	n3iwfCtx := app.n3iwfCtx

	invalidIP = net.ParseIP("10.0.0.0")
	invalidIP2 = net.ParseIP("10.0.0.255")
	n3iwfCtx.IPSecInnerIPPool, err = ippool.NewIPPool("10.0.0.0/24")
	require.NoError(t, err)

	for i := 1; i <= 253; i++ {
		ip, err = n3iwfCtx.NewIPsecInnerUEIP(&n3iwf_context.N3IWFIkeUe{})
		require.NoError(t, err)
		require.NotEqual(t, cfg.GetIPSecGatewayAddr(), ip.String())
		require.NotEqual(t, ip, invalidIP)
		require.NotEqual(t, ip, invalidIP2)
	}

	_, err = n3iwfCtx.NewIPsecInnerUEIP(&n3iwf_context.N3IWFIkeUe{})
	require.Error(t, err)

	n3iwfCtx.AllocatedUEIPAddress = sync.Map{}

	n3iwfCtx.IPSecInnerIPPool, err = ippool.NewIPPool("10.0.0.0/16")
	require.NoError(t, err)

	invalidIP2 = net.ParseIP("10.0.255.255")
	for i := 1; i <= 65533; i++ {
		ip, err = n3iwfCtx.NewIPsecInnerUEIP(&n3iwf_context.N3IWFIkeUe{})
		require.NoError(t, err)
		require.NotEqual(t, cfg.GetIPSecGatewayAddr(), ip.String())
		require.NotEqual(t, ip, invalidIP)
		require.NotEqual(t, ip, invalidIP2)
	}

	_, err = n3iwfCtx.NewIPsecInnerUEIP(&n3iwf_context.N3IWFIkeUe{})
	require.Error(t, err)
}
