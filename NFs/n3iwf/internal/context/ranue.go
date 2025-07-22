package context

import (
	"fmt"
	"net"

	"github.com/free5gc/ngap/ngapType"
)

type UeCtxRelState bool

const (
	// NGAP has already received UE Context release command
	UeCtxRelStateNone    UeCtxRelState = false
	UeCtxRelStateOngoing UeCtxRelState = true
)

type PduSessResRelState bool

const (
	// NGAP has not received Pdu Session resouces release request
	PduSessResRelStateNone    PduSessResRelState = false
	PduSessResRelStateOngoing PduSessResRelState = true
)

type RanUe interface {
	// Get Attributes
	GetUserLocationInformation() *ngapType.UserLocationInformation
	GetSharedCtx() *RanUeSharedCtx

	// User Plane Traffic
	// ForwardDL(gtpQoSMsg.QoSTPDUPacket)
	// ForwardUL()

	// Others
	CreatePDUSession(int64, ngapType.SNSSAI) (*PDUSession, error)
	DeletePDUSession(int64)
	FindPDUSession(int64) *PDUSession
	Remove() error
}

type RanUeSharedCtx struct {
	// UE identity
	RanUeNgapId  int64
	AmfUeNgapId  int64
	IPAddrv4     string
	IPAddrv6     string
	PortNumber   int32
	MaskedIMEISV *ngapType.MaskedIMEISV // TS 38.413 9.3.1.54
	Guti         string

	// Relative Context
	N3iwfCtx *N3IWFContext
	AMF      *N3IWFAMF

	// Security
	SecurityCapabilities *ngapType.UESecurityCapabilities // TS 38.413 9.3.1.86

	// PDU Session
	PduSessionList map[int64]*PDUSession // pduSessionId as key

	// PDU Session Setup Temporary Data
	TemporaryPDUSessionSetupData *PDUSessionSetupTemporaryData

	// Others
	Guami                            *ngapType.GUAMI
	IndexToRfsp                      int64
	Ambr                             *ngapType.UEAggregateMaximumBitRate
	AllowedNssai                     *ngapType.AllowedNSSAI
	RadioCapability                  *ngapType.UERadioCapability                // TODO: This is for RRC, can be deleted
	CoreNetworkAssistanceInformation *ngapType.CoreNetworkAssistanceInformation // TS 38.413 9.3.1.15
	IMSVoiceSupported                int32
	RRCEstablishmentCause            int16
	PduSessionReleaseList            ngapType.PDUSessionResourceReleasedListRelRes
	UeCtxRelState                    UeCtxRelState
	PduSessResRelState               PduSessResRelState
}

type PDUSession struct {
	Id                               int64 // PDU Session ID
	Type                             *ngapType.PDUSessionType
	Ambr                             *ngapType.PDUSessionAggregateMaximumBitRate
	Snssai                           ngapType.SNSSAI
	NetworkInstance                  *ngapType.NetworkInstance
	SecurityCipher                   bool
	SecurityIntegrity                bool
	MaximumIntegrityDataRateUplink   *ngapType.MaximumIntegrityProtectedDataRate
	MaximumIntegrityDataRateDownlink *ngapType.MaximumIntegrityProtectedDataRate
	GTPConnInfo                      *GTPConnectionInfo
	QFIList                          []uint8
	QosFlows                         map[int64]*QosFlow // QosFlowIdentifier as key
}

type QosFlow struct {
	Identifier int64
	Parameters ngapType.QosFlowLevelQosParameters
}

type GTPConnectionInfo struct {
	UPFIPAddr    string
	UPFUDPAddr   net.Addr
	IncomingTEID uint32
	OutgoingTEID uint32
}

type PDUSessionSetupTemporaryData struct {
	// Slice of unactivated PDU session
	UnactivatedPDUSession []*PDUSession // PDUSession as content
	// NGAPProcedureCode is used to identify which type of
	// response shall be used
	NGAPProcedureCode ngapType.ProcedureCode
	// PDU session setup list response
	SetupListCxtRes  *ngapType.PDUSessionResourceSetupListCxtRes
	FailedListCxtRes *ngapType.PDUSessionResourceFailedToSetupListCxtRes
	SetupListSURes   *ngapType.PDUSessionResourceSetupListSURes
	FailedListSURes  *ngapType.PDUSessionResourceFailedToSetupListSURes
	// List of Error for failed setup PDUSessionID
	FailedErrStr []EvtError // Error string as content
	// Current Index of UnactivatedPDUSession
	Index int
}

func (ranUe *RanUeSharedCtx) GetSharedCtx() *RanUeSharedCtx {
	return ranUe
}

func (ranUe *RanUeSharedCtx) FindPDUSession(pduSessionID int64) *PDUSession {
	if pduSession, ok := ranUe.PduSessionList[pduSessionID]; ok {
		return pduSession
	} else {
		return nil
	}
}

func (ranUe *RanUeSharedCtx) CreatePDUSession(pduSessionID int64, snssai ngapType.SNSSAI) (*PDUSession, error) {
	if _, exists := ranUe.PduSessionList[pduSessionID]; exists {
		return nil, fmt.Errorf("PDU Session[ID:%d] is already exists", pduSessionID)
	}
	pduSession := &PDUSession{
		Id:       pduSessionID,
		Snssai:   snssai,
		QosFlows: make(map[int64]*QosFlow),
	}
	ranUe.PduSessionList[pduSessionID] = pduSession
	return pduSession, nil
}

func (ranUe *RanUeSharedCtx) DeletePDUSession(pduSessionId int64) {
	delete(ranUe.PduSessionList, pduSessionId)
}
