package perio

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/free5gc/go-upf/internal/report"
)

type testHandler struct{}

func NewTestHandler() *testHandler {
	return &testHandler{}
}

var testSessRpts map[uint64]*report.SessReport // key: SEID

func (h *testHandler) NotifySessReport(sessRpt report.SessReport) {
	testSessRpts[sessRpt.SEID] = &sessRpt
}

func (h *testHandler) PopBufPkt(lSeid uint64, pdrid uint16) ([]byte, bool) {
	return nil, true
}

func testGetUSAReport(lSeidUrridsMap map[uint64][]uint32) (map[uint64][]report.USAReport, error) {
	sessUsars := make(map[uint64][]report.USAReport)

	v := report.VolumeMeasure{
		UplinkVolume:   10,
		DownlinkVolume: 20,
		TotalVolume:    30,
	}

	for lSeid, urrids := range lSeidUrridsMap {
		for _, urrid := range urrids {
			sessUsars[lSeid] = append(sessUsars[lSeid], report.USAReport{
				URRID:        urrid,
				USARTrigger:  report.UsageReportTrigger{Flags: report.USAR_TRIG_PERIO},
				VolumMeasure: v,
			})
		}
	}

	return sessUsars, nil
}

func TestServer(t *testing.T) {
	var wg sync.WaitGroup
	s, err := OpenServer(&wg)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		s.Close()
		wg.Wait()
	}()

	testSessRpts = make(map[uint64]*report.SessReport)
	h := NewTestHandler()
	s.Handle(h, testGetUSAReport)

	// 1. Add 3 PERIO URRs
	s.AddPeriodReportTimer(
		1, // lSeid
		1, // urrid
		1*time.Second)
	s.AddPeriodReportTimer(
		1, // lSeid
		2, // urrid
		1*time.Second)
	s.AddPeriodReportTimer(
		2, // lSeid
		1, // urrid
		2*time.Second)

	time.Sleep(2100 * time.Millisecond)

	expectedSessRpts := map[uint64]*report.SessReport{
		1: {
			SEID: 1,
			Reports: []report.Report{
				report.USAReport{
					URRID: 1,
					USARTrigger: report.UsageReportTrigger{
						Flags: report.USAR_TRIG_PERIO,
					},
					VolumMeasure: report.VolumeMeasure{
						TotalVolume:    30,
						UplinkVolume:   10,
						DownlinkVolume: 20,
					},
				},
				report.USAReport{
					URRID: 2,
					USARTrigger: report.UsageReportTrigger{
						Flags: report.USAR_TRIG_PERIO,
					},
					VolumMeasure: report.VolumeMeasure{
						TotalVolume:    30,
						UplinkVolume:   10,
						DownlinkVolume: 20,
					},
				},
			},
		},
		2: {
			SEID: 2,
			Reports: []report.Report{
				report.USAReport{
					URRID: 1,
					USARTrigger: report.UsageReportTrigger{
						Flags: report.USAR_TRIG_PERIO,
					},
					VolumMeasure: report.VolumeMeasure{
						TotalVolume:    30,
						UplinkVolume:   10,
						DownlinkVolume: 20,
					},
				},
			},
		},
	}

	// Check the reports
	require.Contains(t, testSessRpts, uint64(1))
	require.Contains(t, testSessRpts, uint64(2))
	require.ElementsMatch(t, testSessRpts[1].Reports, expectedSessRpts[1].Reports)
	require.ElementsMatch(t, testSessRpts[2].Reports, expectedSessRpts[2].Reports)

	testSessRpts = make(map[uint64]*report.SessReport)
	expectedSessRpts2 := map[uint64]*report.SessReport{
		1: {
			SEID:    1,
			Reports: expectedSessRpts[1].Reports[1:2],
		},
	}

	// 2. Delete 2 PERIO URRs
	s.DelPeriodReportTimer(
		1, // lSeid
		1, // urrid
	)
	s.DelPeriodReportTimer(
		2, // lSeid
		1, // urrid
	)

	time.Sleep(1100 * time.Millisecond)

	// Check the reports
	require.Contains(t, testSessRpts, uint64(1))
	require.ElementsMatch(t, testSessRpts[1].Reports, expectedSessRpts2[1].Reports)

	// 3. Make sure SEID(2) PERIO timer not launched
	testSessRpts = make(map[uint64]*report.SessReport)
	time.Sleep(1100 * time.Millisecond)

	require.Contains(t, testSessRpts, uint64(1))
	require.NotContains(t, testSessRpts, uint64(2))
	require.ElementsMatch(t, testSessRpts[1].Reports, expectedSessRpts2[1].Reports)
}
