package n3iwf_context

import (
	"fmt"
	"free5gc/lib/ngap/ngapType"
)

const (
	AmfUeNgapIdUnspecified int64 = 1099511627776
)

type N3IWFUe struct {
	/* UE identity*/
	RanUeNgapId  int64
	AmfUeNgapId  int64
	IPAddrv4     string
	IPAddrv6     string
	PortNumber   int32
	MaskedIMEISV *ngapType.MaskedIMEISV // TS 38.413 9.3.1.54

	/* PDU Session */
	PduSessionList map[int64]*PDUSession // pduSessionId as key

	/* Security */
	Kn3iwf               []uint8                          // 32 bytes (256 bits), value is from NGAP IE "Security Key"
	SecurityCapabilities *ngapType.UESecurityCapabilities // TS 38.413 9.3.1.86

	/* Others */
	Guami                            *ngapType.GUAMI
	IndexToRfsp                      int64
	Ambr                             *ngapType.UEAggregateMaximumBitRate
	AllowedNssai                     *ngapType.AllowedNSSAI
	RadioCapability                  *ngapType.UERadioCapability
	CoreNetworkAssistanceInformation *ngapType.CoreNetworkAssistanceInformation // TS 38.413 9.3.1.15
}

type PDUSession struct {
	Id              int64 // PDU Session ID
	Type            ngapType.PDUSessionType
	Snssai          ngapType.SNSSAI
	GTPEndpointIPv4 string
	GTPEndpointIPv6 string
	TEID            string
	QosFlows        map[int64]*QosFlow // QosFlowIdentifier as key
}

type QosFlow struct {
	Identifier int64
	Parameters ngapType.QosFlowLevelQosParameters
}

func (ue *N3IWFUe) init() {
	ue.PduSessionList = make(map[int64]*PDUSession)
}

func (ue *N3IWFUe) Remove() {
	n3iwfSelf := N3IWFSelf()
	delete(n3iwfSelf.UePool, ue.RanUeNgapId)
}

func (ue *N3IWFUe) FindPDUSession(pduSessionID int64) *PDUSession {
	if _, exists := ue.PduSessionList[pduSessionID]; exists {
		return ue.PduSessionList[pduSessionID]
	}
	return nil
}

func (ue *N3IWFUe) CreatePDUSession(pduSessionID int64, snssai ngapType.SNSSAI) (*PDUSession, error) {
	if _, exists := ue.PduSessionList[pduSessionID]; exists {
		return nil, fmt.Errorf("PDU Session[ID:%d] is already exists", pduSessionID)
	}
	pduSession := PDUSession{}
	pduSession.Id = pduSessionID
	pduSession.Snssai = snssai
	pduSession.QosFlows = make(map[int64]*QosFlow)
	ue.PduSessionList[pduSessionID] = &pduSession
	return &pduSession, nil
}
