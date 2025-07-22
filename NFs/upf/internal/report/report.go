package report

import (
	"encoding/binary"
	"time"

	"github.com/pkg/errors"
	"github.com/wmnsk/go-pfcp/ie"
)

type ReportType int

// 29244-ga0 8.2.21 Report Type
const (
	DLDR ReportType = iota + 1
	USAR
	ERIR
	UPIR
	TMIR
	SESR
	UISR
)

func (t ReportType) String() string {
	str := []string{"", "DLDR", "USAR", "ERIR", "UPIR", "TMIR", "SESR", "UISR"}
	return str[t]
}

type Report interface {
	Type() ReportType
}

type DLDReport struct {
	PDRID  uint16
	Action uint16
	BufPkt []byte
}

func (r DLDReport) Type() ReportType {
	return DLDR
}

type MeasureMethod struct {
	DURAT bool
	VOLUM bool
	EVENT bool
}

type MeasureInformation struct {
	MBQE  bool
	INAM  bool
	RADI  bool
	ISTM  bool
	MNOP  bool
	SSPOC bool
	ASPOC bool
	CIAM  bool
}

type USAReport struct {
	URRID        uint32
	URSEQN       uint32
	USARTrigger  UsageReportTrigger
	VolumMeasure VolumeMeasure
	DuratMeasure DurationMeasure
	QueryUrrRef  uint32
	StartTime    time.Time
	EndTime      time.Time
}

func (r USAReport) Type() ReportType {
	return USAR
}

func (r USAReport) IEsWithinSessReportReq(
	method MeasureMethod, info MeasureInformation,
) []*ie.IE {
	ies := []*ie.IE{
		ie.NewURRID(r.URRID),
		ie.NewURSEQN(r.URSEQN),
		r.USARTrigger.IE(),
	}
	if !r.USARTrigger.START() && !r.USARTrigger.STOPT() && !r.USARTrigger.MACAR() {
		// These IEs shall be present, except if the Usage Report
		// Trigger indicates 'Start of Traffic', 'Stop of Traffic' or 'MAC
		// Addresses Reporting'.
		ies = append(ies, ie.NewStartTime(r.StartTime), ie.NewEndTime(r.EndTime))
	}
	if method.VOLUM {
		r.VolumMeasure.SetFlags(info.MNOP)
		ies = append(ies, r.VolumMeasure.IE())
	}
	if method.DURAT {
		ies = append(ies, r.DuratMeasure.IE())
	}
	return ies
}

func (r USAReport) IEsWithinSessModRsp(
	method MeasureMethod, info MeasureInformation,
) []*ie.IE {
	ies := []*ie.IE{
		ie.NewURRID(r.URRID),
		ie.NewURSEQN(r.URSEQN),
		r.USARTrigger.IE(),
	}
	if !r.USARTrigger.START() && !r.USARTrigger.STOPT() && !r.USARTrigger.MACAR() {
		// These IEs shall be present, except if the Usage Report
		// Trigger indicates 'Start of Traffic', 'Stop of Traffic' or 'MAC
		// Addresses Reporting'.
		ies = append(ies, ie.NewStartTime(r.StartTime), ie.NewEndTime(r.EndTime))
	}
	if method.VOLUM {
		r.VolumMeasure.SetFlags(info.MNOP)
		ies = append(ies, r.VolumMeasure.IE())
	}
	if method.DURAT {
		ies = append(ies, r.DuratMeasure.IE())
	}
	return ies
}

func (r USAReport) IEsWithinSessDelRsp(
	method MeasureMethod, info MeasureInformation,
) []*ie.IE {
	ies := []*ie.IE{
		ie.NewURRID(r.URRID),
		ie.NewURSEQN(r.URSEQN),
		r.USARTrigger.IE(),
	}
	if !r.USARTrigger.START() && !r.USARTrigger.STOPT() && !r.USARTrigger.MACAR() {
		// These IEs shall be present, except if the Usage Report
		// Trigger indicates 'Start of Traffic', 'Stop of Traffic' or 'MAC
		// Addresses Reporting'.
		ies = append(ies, ie.NewStartTime(r.StartTime), ie.NewEndTime(r.EndTime))
	}
	if method.VOLUM {
		r.VolumMeasure.SetFlags(info.MNOP)
		ies = append(ies, r.VolumMeasure.IE())
	}
	if method.DURAT {
		ies = append(ies, r.DuratMeasure.IE())
	}
	return ies
}

// Reporting Triggers IE bits definition
const (
	RPT_TRIG_PERIO = 1 << iota
	RPT_TRIG_VOLTH
	RPT_TRIG_TIMTH
	RPT_TRIG_QUHTI
	RPT_TRIG_START
	RPT_TRIG_STOPT
	RPT_TRIG_DROTH
	RPT_TRIG_LIUSA
	RPT_TRIG_VOLQU
	RPT_TRIG_TIMQU
	RPT_TRIG_ENVCL
	RPT_TRIG_MACAR
	RPT_TRIG_EVETH
	RPT_TRIG_EVEQU
	RPT_TRIG_IPMJL
	RPT_TRIG_QUVTI
	RPT_TRIG_REEMR
	RPT_TRIG_UPINT
)

type ReportingTrigger struct {
	Flags uint32
}

func (r *ReportingTrigger) Unmarshal(b []byte) error {
	if len(b) < 2 {
		return errors.Errorf("ReportingTrigger Unmarshal: less than 2 bytes")
	}
	// slice len might be 2 or 3; enlarge slice to 4 bytes at least
	v := make([]byte, len(b)+2)
	copy(v, b)
	r.Flags = binary.LittleEndian.Uint32(v)
	return nil
}

func (r *ReportingTrigger) IE() *ie.IE {
	b := make([]byte, 4)
	binary.LittleEndian.PutUint32(b, r.Flags)
	return ie.NewReportingTriggers(b[:3]...)
}

func (r *ReportingTrigger) PERIO() bool {
	return r.Flags&RPT_TRIG_PERIO != 0
}

func (r *ReportingTrigger) VOLTH() bool {
	return r.Flags&RPT_TRIG_VOLTH != 0
}

func (r *ReportingTrigger) TIMTH() bool {
	return r.Flags&RPT_TRIG_TIMTH != 0
}

func (r *ReportingTrigger) QUHTI() bool {
	return r.Flags&RPT_TRIG_QUHTI != 0
}

func (r *ReportingTrigger) START() bool {
	return r.Flags&RPT_TRIG_START != 0
}

func (r *ReportingTrigger) STOPT() bool {
	return r.Flags&RPT_TRIG_STOPT != 0
}

func (r *ReportingTrigger) DROTH() bool {
	return r.Flags&RPT_TRIG_DROTH != 0
}

func (r *ReportingTrigger) LIUSA() bool {
	return r.Flags&RPT_TRIG_LIUSA != 0
}

func (r *ReportingTrigger) VOLQU() bool {
	return r.Flags&RPT_TRIG_VOLQU != 0
}

func (r *ReportingTrigger) TIMQU() bool {
	return r.Flags&RPT_TRIG_TIMQU != 0
}

func (r *ReportingTrigger) ENVCL() bool {
	return r.Flags&RPT_TRIG_ENVCL != 0
}

func (r *ReportingTrigger) MACAR() bool {
	return r.Flags&RPT_TRIG_MACAR != 0
}

func (r *ReportingTrigger) EVETH() bool {
	return r.Flags&RPT_TRIG_EVETH != 0
}

func (r *ReportingTrigger) EVEQU() bool {
	return r.Flags&RPT_TRIG_EVEQU != 0
}

func (r *ReportingTrigger) IPMJL() bool {
	return r.Flags&RPT_TRIG_IPMJL != 0
}

func (r *ReportingTrigger) QUVTI() bool {
	return r.Flags&RPT_TRIG_QUVTI != 0
}

func (r *ReportingTrigger) REEMR() bool {
	return r.Flags&RPT_TRIG_REEMR != 0
}

func (r *ReportingTrigger) UPINT() bool {
	return r.Flags&RPT_TRIG_UPINT != 0
}

// Usage Report Trigger IE bits definition
const (
	USAR_TRIG_PERIO = 1 << iota
	USAR_TRIG_VOLTH
	USAR_TRIG_TIMTH
	USAR_TRIG_QUHTI
	USAR_TRIG_START
	USAR_TRIG_STOPT
	USAR_TRIG_DROTH
	USAR_TRIG_IMMER
	USAR_TRIG_VOLQU
	USAR_TRIG_TIMQU
	USAR_TRIG_LIUSA
	USAR_TRIG_TERMR
	USAR_TRIG_MONIT
	USAR_TRIG_ENVCL
	USAR_TRIG_MACAR
	USAR_TRIG_EVETH
	USAR_TRIG_EVEQU
	USAR_TRIG_TEBUR
	USAR_TRIG_IPMJL
	USAR_TRIG_QUVTI
	USAR_TRIG_EMRRE
	USAR_TRIG_UPINT
)

type UsageReportTrigger struct {
	Flags uint32
}

func (t *UsageReportTrigger) SetReportingTrigger(r uint32) {
	switch r {
	case RPT_TRIG_PERIO:
		t.Flags |= USAR_TRIG_PERIO
	case RPT_TRIG_VOLTH:
		t.Flags |= USAR_TRIG_VOLTH
	case RPT_TRIG_TIMTH:
		t.Flags |= USAR_TRIG_TIMTH
	case RPT_TRIG_QUHTI:
		t.Flags |= USAR_TRIG_QUHTI
	case RPT_TRIG_START:
		t.Flags |= USAR_TRIG_START
	case RPT_TRIG_STOPT:
		t.Flags |= USAR_TRIG_STOPT
	case RPT_TRIG_DROTH:
		t.Flags |= USAR_TRIG_DROTH
	case RPT_TRIG_LIUSA:
		t.Flags |= USAR_TRIG_LIUSA
	case RPT_TRIG_VOLQU:
		t.Flags |= USAR_TRIG_VOLQU
	case RPT_TRIG_TIMQU:
		t.Flags |= USAR_TRIG_TIMQU
	case RPT_TRIG_ENVCL:
		t.Flags |= USAR_TRIG_ENVCL
	case RPT_TRIG_MACAR:
		t.Flags |= USAR_TRIG_MACAR
	case RPT_TRIG_EVETH:
		t.Flags |= USAR_TRIG_EVETH
	case RPT_TRIG_EVEQU:
		t.Flags |= USAR_TRIG_EVEQU
	case RPT_TRIG_IPMJL:
		t.Flags |= USAR_TRIG_IPMJL
	case RPT_TRIG_QUVTI:
		t.Flags |= USAR_TRIG_QUVTI
	case RPT_TRIG_UPINT:
		t.Flags |= USAR_TRIG_UPINT
	}
}

func (t *UsageReportTrigger) IE() *ie.IE {
	b := make([]byte, 4)
	binary.LittleEndian.PutUint32(b, t.Flags)
	return ie.NewUsageReportTrigger(b[:3]...)
}

func (t *UsageReportTrigger) PERIO() bool {
	return t.Flags&USAR_TRIG_PERIO != 0
}

func (t *UsageReportTrigger) VOLTH() bool {
	return t.Flags&USAR_TRIG_VOLTH != 0
}

func (t *UsageReportTrigger) TIMTH() bool {
	return t.Flags&USAR_TRIG_TIMTH != 0
}

func (t *UsageReportTrigger) QUHTI() bool {
	return t.Flags&USAR_TRIG_QUHTI != 0
}

func (t *UsageReportTrigger) START() bool {
	return t.Flags&USAR_TRIG_START != 0
}

func (t *UsageReportTrigger) STOPT() bool {
	return t.Flags&USAR_TRIG_STOPT != 0
}

func (t *UsageReportTrigger) DROTH() bool {
	return t.Flags&USAR_TRIG_DROTH != 0
}

func (t *UsageReportTrigger) IMMER() bool {
	return t.Flags&USAR_TRIG_IMMER != 0
}

func (t *UsageReportTrigger) VOLQU() bool {
	return t.Flags&USAR_TRIG_VOLQU != 0
}

func (t *UsageReportTrigger) TIMQU() bool {
	return t.Flags&USAR_TRIG_TIMQU != 0
}

func (t *UsageReportTrigger) LIUSA() bool {
	return t.Flags&USAR_TRIG_LIUSA != 0
}

func (t *UsageReportTrigger) TERMR() bool {
	return t.Flags&USAR_TRIG_TERMR != 0
}

func (t *UsageReportTrigger) MONIT() bool {
	return t.Flags&USAR_TRIG_MONIT != 0
}

func (t *UsageReportTrigger) ENVCL() bool {
	return t.Flags&USAR_TRIG_ENVCL != 0
}

func (t *UsageReportTrigger) MACAR() bool {
	return t.Flags&USAR_TRIG_MACAR != 0
}

func (t *UsageReportTrigger) EVETH() bool {
	return t.Flags&USAR_TRIG_EVETH != 0
}

func (t *UsageReportTrigger) EVEQU() bool {
	return t.Flags&USAR_TRIG_EVEQU != 0
}

func (t *UsageReportTrigger) TEBUR() bool {
	return t.Flags&USAR_TRIG_TEBUR != 0
}

func (t *UsageReportTrigger) IPMJL() bool {
	return t.Flags&USAR_TRIG_IPMJL != 0
}

func (t *UsageReportTrigger) QUVTI() bool {
	return t.Flags&USAR_TRIG_QUVTI != 0
}

func (t *UsageReportTrigger) EMRRE() bool {
	return t.Flags&USAR_TRIG_EMRRE != 0
}

func (t *UsageReportTrigger) UPINT() bool {
	return t.Flags&USAR_TRIG_UPINT != 0
}

// Volume Measurement IE Flag bits definition
const (
	TOVOL uint8 = 1 << iota
	ULVOL
	DLVOL
	TONOP
	ULNOP
	DLNOP
)

type VolumeMeasure struct {
	Flags          uint8
	TotalVolume    uint64
	UplinkVolume   uint64
	DownlinkVolume uint64
	TotalPktNum    uint64
	UplinkPktNum   uint64
	DownlinkPktNum uint64
}

func (m *VolumeMeasure) SetFlags(mnop bool) {
	m.Flags |= (TOVOL | ULVOL | DLVOL)
	if mnop {
		m.Flags |= (TONOP | ULNOP | DLNOP)
	}
}

func (m *VolumeMeasure) IE() *ie.IE {
	return ie.NewVolumeMeasurement(
		m.Flags,
		m.TotalVolume,
		m.UplinkVolume,
		m.DownlinkVolume,
		m.TotalPktNum,
		m.UplinkPktNum,
		m.DownlinkPktNum,
	)
}

type DurationMeasure struct {
	DurationValue uint64
}

func (m *DurationMeasure) IE() *ie.IE {
	return ie.NewDurationMeasurement(time.Duration(m.DurationValue))
}

// Apply Action IE bits definition
const (
	APPLY_ACT_DROP = 1 << iota
	APPLY_ACT_FORW
	APPLY_ACT_BUFF
	APPLY_ACT_NOCP
	APPLY_ACT_DUPL
	APPLY_ACT_IPMA
	APPLY_ACT_IPMD
	APPLY_ACT_DFRT
	APPLY_ACT_EDRT
	APPLY_ACT_BDPN
	APPLY_ACT_DDPN
	APPLY_ACT_FSSM
	APPLY_ACT_MBSU
)

type ApplyAction struct {
	Flags uint16
}

func (a *ApplyAction) Unmarshal(b []byte) error {
	if len(b) < 1 {
		return errors.Errorf("ApplyAction Unmarshal: less than 1 bytes")
	}

	// slice len might be 1 or 2; enlarge slice to 2 bytes at least
	v := make([]byte, max(2, len(b)))
	copy(v, b)
	a.Flags = binary.LittleEndian.Uint16(v)
	return nil
}

func (a *ApplyAction) DROP() bool {
	return a.Flags&APPLY_ACT_DROP != 0
}

func (a *ApplyAction) FORW() bool {
	return a.Flags&APPLY_ACT_FORW != 0
}

func (a *ApplyAction) BUFF() bool {
	return a.Flags&APPLY_ACT_BUFF != 0
}

func (a *ApplyAction) NOCP() bool {
	return a.Flags&APPLY_ACT_NOCP != 0
}

func (a *ApplyAction) DUPL() bool {
	return a.Flags&APPLY_ACT_DUPL != 0
}

func (a *ApplyAction) IPMA() bool {
	return a.Flags&APPLY_ACT_IPMA != 0
}

func (a *ApplyAction) IPMD() bool {
	return a.Flags&APPLY_ACT_IPMD != 0
}

func (a *ApplyAction) DFRT() bool {
	return a.Flags&APPLY_ACT_DFRT != 0
}

func (a *ApplyAction) EDRT() bool {
	return a.Flags&APPLY_ACT_EDRT != 0
}

func (a *ApplyAction) BDPN() bool {
	return a.Flags&APPLY_ACT_BDPN != 0
}

func (a *ApplyAction) DDPN() bool {
	return a.Flags&APPLY_ACT_DDPN != 0
}

func (a *ApplyAction) FSSM() bool {
	return a.Flags&APPLY_ACT_FSSM != 0
}

func (a *ApplyAction) MBSU() bool {
	return a.Flags&APPLY_ACT_MBSU != 0
}

type SessReport struct {
	SEID    uint64
	Reports []Report
}

type BufInfo struct {
	SEID  uint64
	PDRID uint16
}
