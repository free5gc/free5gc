package ngap

import (
	"testing"

	"github.com/stretchr/testify/require"

	n3iwf_context "github.com/free5gc/n3iwf/internal/context"
	"github.com/free5gc/n3iwf/internal/ike"
	"github.com/free5gc/n3iwf/pkg/factory"
)

func TestReleaseIkeUeAndRanUe(t *testing.T) {
	n3iwf, err := NewN3iwfTestApp(&factory.Config{})
	require.NoError(t, err)

	n3iwf.ngapServer, err = NewServer(n3iwf)
	require.NoError(t, err)

	n3iwf.ikeServer, err = ike.NewServer(n3iwf)
	require.NoError(t, err)

	n3iwfCtx := n3iwf.n3iwfCtx
	ranUe := &n3iwf_context.N3IWFRanUe{
		RanUeSharedCtx: n3iwf_context.RanUeSharedCtx{
			N3iwfCtx: n3iwfCtx,
		},
	}

	ranUeNgapId := int64(0x1234567890ABCDEF)
	spi := uint64(123)
	ranUe.RanUeNgapId = ranUeNgapId
	n3iwfCtx.RANUePool.Store(ranUeNgapId, ranUe)
	n3iwfCtx.NGAPIdToIKESPI.Store(ranUeNgapId, spi)
	n3iwfCtx.IKESPIToNGAPId.Store(spi, ranUeNgapId)

	stopCh := make(chan struct{})
	rcvIkeEvtCh := n3iwf.mockIkeEvtCh.GetRcvChan()

	go func() {
		for {
			select {
			case <-stopCh:
				return
			case rcvEvt := <-rcvIkeEvtCh:
				if rcvEvt.Type() != n3iwf_context.IKEDeleteRequest {
					t.Errorf("Receive Wrong Event")
				}
			}
		}
	}()

	err = n3iwf.ngapServer.releaseIkeUeAndRanUe(ranUe)
	require.NoError(t, err)

	_, ok := n3iwfCtx.RANUePool.Load(ranUeNgapId)
	if ok {
		t.Errorf("RanUe doesn't get remove")
	}

	stopCh <- struct{}{}
	n3iwf.ngapServer.Stop()
}
