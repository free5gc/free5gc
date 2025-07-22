package forwarder

import (
	"bytes"
	"net"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/khirono/go-nl"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/wmnsk/go-pfcp/ie"

	"github.com/free5gc/go-gtp5gnl"
	"github.com/free5gc/go-upf/internal/report"
	"github.com/free5gc/go-upf/pkg/factory"
)

func Test_convertSlice(t *testing.T) {
	t.Run("convert slices", func(t *testing.T) {
		b := convertSlice([][]uint16{{1}, {2, 4}})
		want := []byte{0x01, 0x00, 0x01, 0x00, 0x04, 0x00, 0x02, 0x00}
		if !bytes.Equal(b, want) {
			t.Errorf("want %x; but got %x\n", want, b)
		}
	})
}

type testHandler struct{}

var testSessRpts map[uint64]*report.SessReport // key: SEID

func (h *testHandler) NotifySessReport(sessRpt report.SessReport) {
	testSessRpts[sessRpt.SEID] = &sessRpt
}

func (h *testHandler) PopBufPkt(lSeid uint64, pdrid uint16) ([]byte, bool) {
	return nil, true
}

func TestGtp5g_CreateRules(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping testing in short mode")
	}

	var wg sync.WaitGroup
	g, err := OpenGtp5g(&wg, ":"+strconv.Itoa(factory.UpfGtpDefaultPort), 1400)
	if err != nil {
		t.Fatal(err)
	}
	defer g.Close()

	testSessRpts = make(map[uint64]*report.SessReport)
	g.HandleReport(&testHandler{})

	lSeid := uint64(1)
	t.Run("create rules", func(t *testing.T) {
		far := ie.NewCreateFAR(
			ie.NewFARID(2),
			ie.NewApplyAction(0x2),
			ie.NewForwardingParameters(
				ie.NewDestinationInterface(ie.DstInterfaceSGiLANN6LAN),
				ie.NewNetworkInstance("internet"),
			),
		)

		err = g.CreateFAR(lSeid, far)
		if err != nil {
			t.Fatal(err)
		}

		far = ie.NewCreateFAR(
			ie.NewFARID(4),
			ie.NewApplyAction(0x2),
		)

		err = g.CreateFAR(lSeid, far)
		if err != nil {
			t.Fatal(err)
		}

		qer := ie.NewCreateQER(
			ie.NewQERID(1),
			ie.NewGateStatus(ie.GateStatusOpen, ie.GateStatusOpen),
			ie.NewMBR(200000, 100000),
			ie.NewQFI(10),
		)

		err = g.CreateQER(lSeid, qer)
		if err != nil {
			t.Fatal(err)
		}

		rptTrig := report.ReportingTrigger{
			Flags: report.RPT_TRIG_PERIO,
		}

		urr := ie.NewCreateURR(
			ie.NewURRID(1),
			ie.NewMeasurementPeriod(1*time.Second),
			ie.NewMeasurementMethod(0, 1, 0),
			rptTrig.IE(),
			ie.NewMeasurementInformation(4),
		)
		err = g.CreateURR(lSeid, urr)
		if err != nil {
			t.Fatal(err)
		}

		rptTrig.Flags = report.RPT_TRIG_VOLTH | report.RPT_TRIG_VOLQU
		urr = ie.NewCreateURR(
			ie.NewURRID(2),
			ie.NewMeasurementMethod(0, 1, 0),
			rptTrig.IE(),
			ie.NewMeasurementInformation(4),
			ie.NewVolumeThreshold(7, 10000, 20000, 30000),
			ie.NewVolumeQuota(7, 40000, 50000, 60000),
		)
		err = g.CreateURR(lSeid, urr)
		if err != nil {
			t.Fatal(err)
		}

		pdr := ie.NewCreatePDR(
			ie.NewPDRID(1),
			ie.NewPrecedence(255),
			ie.NewPDI(
				ie.NewSourceInterface(ie.SrcInterfaceAccess),
				ie.NewFTEID(
					0x01,
					1,
					net.ParseIP("30.30.30.2"),
					nil,
					0,
				),
				ie.NewNetworkInstance(""),
				ie.NewUEIPAddress(
					0x02,
					"60.60.0.1",
					"",
					0,
					0,
				),
			),
			ie.NewOuterHeaderRemoval(0, 0),
			ie.NewFARID(2),
			ie.NewQERID(1),
			ie.NewURRID(1),
			ie.NewURRID(2),
		)

		err = g.CreatePDR(lSeid, pdr)
		if err != nil {
			t.Fatal(err)
		}

		pdr = ie.NewCreatePDR(
			ie.NewPDRID(3),
			ie.NewPrecedence(255),
			ie.NewPDI(
				ie.NewSourceInterface(ie.SrcInterfaceCore),
				ie.NewNetworkInstance("internet"),
				ie.NewUEIPAddress(
					0x02,
					"60.60.0.1",
					"",
					0,
					0,
				),
			),
			ie.NewFARID(4),
			ie.NewQERID(1),
			ie.NewURRID(1),
		)

		err = g.CreatePDR(lSeid, pdr)
		if err != nil {
			t.Fatal(err)
		}

		time.Sleep(1100 * time.Millisecond)

		require.Contains(t, testSessRpts, lSeid)
		require.Equal(t, len(testSessRpts[lSeid].Reports), 1)
		require.Equal(t, testSessRpts[lSeid].Reports[0].(report.USAReport).URRID, uint32(1))
	})

	t.Run("update rules", func(t *testing.T) {
		rpt := report.ReportingTrigger{
			Flags: report.RPT_TRIG_PERIO,
		}

		urr := ie.NewUpdateURR(
			ie.NewURRID(1),
			ie.NewMeasurementPeriod(2*time.Second),
			rpt.IE(),
		)
		rs, err := g.UpdateURR(lSeid, urr)
		if err != nil {
			t.Fatal(err)
		}

		// TODO: should apply PERIO updateURR and receive final report from old URR
		require.Nil(t, rs)
		// require.NotNil(t, r)
		// require.Equal(t, r.URRID, uint32(1))

		far := ie.NewUpdateFAR(
			ie.NewFARID(4),
			ie.NewApplyAction(0x2),
			ie.NewUpdateForwardingParameters(
				ie.NewDestinationInterface(ie.DstInterfaceAccess),
				ie.NewNetworkInstance("internet"),
				ie.NewOuterHeaderCreation(
					0x0100,
					1,
					"30.30.30.1",
					"",
					0,
					0,
					0,
				),
			),
		)

		err = g.UpdateFAR(lSeid, far)
		if err != nil {
			t.Fatal(err)
		}

		pdr := ie.NewUpdatePDR(
			ie.NewPDRID(3),
			ie.NewPrecedence(255),
			ie.NewPDI(
				ie.NewSourceInterface(ie.SrcInterfaceCore),
				ie.NewNetworkInstance("internet"),
				ie.NewUEIPAddress(
					0x02,
					"60.60.0.1",
					"",
					0,
					0,
				),
			),
			ie.NewFARID(4),
		)

		err = g.UpdatePDR(lSeid, pdr)
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("remove rules", func(t *testing.T) {
		urr := ie.NewRemoveURR(
			ie.NewURRID(1),
		)

		rs, err1 := g.RemoveURR(lSeid, urr)
		if err1 != nil {
			t.Fatal(err1)
		}
		g.log.Infof("Receive final report from URR(%d), rpts: %+v", rs[0].URRID, rs)

		require.NotNil(t, rs)
		g.log.Infof("Receive final report from URR(%d)", rs[0].URRID)

		urr = ie.NewRemoveURR(
			ie.NewURRID(2),
		)
		rs, err1 = g.RemoveURR(lSeid, urr)
		if err1 != nil {
			t.Fatal(err1)
		}

		g.log.Infof("Receive final report from URR(%d), rpts: %+v", rs[0].URRID, rs)

		require.NotNil(t, rs)
		g.log.Infof("Receive final reports from URR(%d)", rs[0].URRID)
	})
}

func TestNewFlowDesc(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping testing in short mode")
	}

	var wg sync.WaitGroup
	g, err := OpenGtp5g(&wg, ":"+strconv.Itoa(factory.UpfGtpDefaultPort), 1400)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		g.Close()
		wg.Wait()
	}()

	cases := []struct {
		name       string
		s          string
		swapSrcDst bool
		attrs      nl.AttrList
		err        error
	}{
		{
			name:       "permit out any to assigned",
			s:          "permit out ip from any to assigned",
			swapSrcDst: false,
			attrs: nl.AttrList{
				nl.Attr{
					Type:  gtp5gnl.FLOW_DESCRIPTION_ACTION,
					Value: nl.AttrU8(gtp5gnl.SDF_FILTER_PERMIT),
				},
				nl.Attr{
					Type:  gtp5gnl.FLOW_DESCRIPTION_DIRECTION,
					Value: nl.AttrU8(gtp5gnl.SDF_FILTER_OUT),
				},
			},
			err: nil,
		},
		{
			name:       "network addr (UL)",
			s:          "permit out ip from 10.20.30.40/24 to 50.60.70.80/16",
			swapSrcDst: false,
			attrs: nl.AttrList{
				nl.Attr{
					Type:  gtp5gnl.FLOW_DESCRIPTION_ACTION,
					Value: nl.AttrU8(gtp5gnl.SDF_FILTER_PERMIT),
				},
				nl.Attr{
					Type:  gtp5gnl.FLOW_DESCRIPTION_DIRECTION,
					Value: nl.AttrU8(gtp5gnl.SDF_FILTER_OUT),
				},
				nl.Attr{
					Type:  gtp5gnl.FLOW_DESCRIPTION_SRC_IPV4,
					Value: nl.AttrBytes(net.IPv4(10, 20, 30, 0).To4()),
				},
				nl.Attr{
					Type:  gtp5gnl.FLOW_DESCRIPTION_DEST_IPV4,
					Value: nl.AttrBytes(net.IPv4(50, 60, 0, 0).To4()),
				},
			},
			err: nil,
		},
		{
			name:       "network addr (DL)",
			s:          "permit out ip from 10.20.30.40/24 to 50.60.70.80/16",
			swapSrcDst: true,
			attrs: nl.AttrList{
				nl.Attr{
					Type:  gtp5gnl.FLOW_DESCRIPTION_ACTION,
					Value: nl.AttrU8(gtp5gnl.SDF_FILTER_PERMIT),
				},
				nl.Attr{
					Type:  gtp5gnl.FLOW_DESCRIPTION_DIRECTION,
					Value: nl.AttrU8(gtp5gnl.SDF_FILTER_OUT),
				},
				nl.Attr{
					Type:  gtp5gnl.FLOW_DESCRIPTION_SRC_IPV4,
					Value: nl.AttrBytes(net.IPv4(50, 60, 0, 0).To4()),
				},
				nl.Attr{
					Type:  gtp5gnl.FLOW_DESCRIPTION_DEST_IPV4,
					Value: nl.AttrBytes(net.IPv4(10, 20, 30, 0).To4()),
				},
			},
			err: nil,
		},
		{
			name:       "source port (DL)",
			s:          "permit out ip from 10.20.30.40/24 345,789-792,1023-1026 to 50.60.70.80/16 456-458,1088,1089",
			swapSrcDst: false,
			attrs: nl.AttrList{
				nl.Attr{
					Type:  gtp5gnl.FLOW_DESCRIPTION_ACTION,
					Value: nl.AttrU8(gtp5gnl.SDF_FILTER_PERMIT),
				},
				nl.Attr{
					Type:  gtp5gnl.FLOW_DESCRIPTION_DIRECTION,
					Value: nl.AttrU8(gtp5gnl.SDF_FILTER_OUT),
				},
				nl.Attr{
					Type:  gtp5gnl.FLOW_DESCRIPTION_SRC_IPV4,
					Value: nl.AttrBytes(net.IPv4(10, 20, 30, 0).To4()),
				},
				nl.Attr{
					Type:  gtp5gnl.FLOW_DESCRIPTION_DEST_IPV4,
					Value: nl.AttrBytes(net.IPv4(50, 60, 0, 0).To4()),
				},
				nl.Attr{
					Type: gtp5gnl.FLOW_DESCRIPTION_SRC_PORT,
					Value: nl.AttrBytes(convertSlice([][]uint16{
						{345},
						{789, 792},
						{1023, 1026},
					})),
				},
				nl.Attr{
					Type: gtp5gnl.FLOW_DESCRIPTION_DEST_PORT,
					Value: nl.AttrBytes(convertSlice([][]uint16{
						{456, 458},
						{1088},
						{1089},
					})),
				},
			},
			err: nil,
		},
		{
			name:       "source port (UL)",
			s:          "permit out ip from 10.20.30.40/24 345,789-792,1023-1026 to 50.60.70.80/16 456-458,1088,1089",
			swapSrcDst: true,
			attrs: nl.AttrList{
				nl.Attr{
					Type:  gtp5gnl.FLOW_DESCRIPTION_ACTION,
					Value: nl.AttrU8(gtp5gnl.SDF_FILTER_PERMIT),
				},
				nl.Attr{
					Type:  gtp5gnl.FLOW_DESCRIPTION_DIRECTION,
					Value: nl.AttrU8(gtp5gnl.SDF_FILTER_OUT),
				},
				nl.Attr{
					Type:  gtp5gnl.FLOW_DESCRIPTION_SRC_IPV4,
					Value: nl.AttrBytes(net.IPv4(50, 60, 0, 0).To4()),
				},
				nl.Attr{
					Type:  gtp5gnl.FLOW_DESCRIPTION_DEST_IPV4,
					Value: nl.AttrBytes(net.IPv4(10, 20, 30, 0).To4()),
				},
				nl.Attr{
					Type: gtp5gnl.FLOW_DESCRIPTION_SRC_PORT,
					Value: nl.AttrBytes(convertSlice([][]uint16{
						{456, 458},
						{1088},
						{1089},
					})),
				},
				nl.Attr{
					Type: gtp5gnl.FLOW_DESCRIPTION_DEST_PORT,
					Value: nl.AttrBytes(convertSlice([][]uint16{
						{345},
						{789, 792},
						{1023, 1026},
					})),
				},
			},
			err: nil,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			attrs, err := g.newFlowDesc(tt.s, tt.swapSrcDst)
			if tt.err == nil {
				if err != nil {
					t.Fatal(err)
				}
				assert.Subset(t, attrs, tt.attrs)
			} else if err != tt.err {
				t.Errorf("wantErr %v; but got %v", tt.err, err)
			}
		})
	}
}

// TODO
// Test on newSdfFilter()
