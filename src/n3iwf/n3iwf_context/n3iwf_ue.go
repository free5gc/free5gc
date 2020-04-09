package n3iwf_context

import (
	"encoding/binary"
	"errors"
	"fmt"
	"free5gc/lib/ngap/ngapType"
	"free5gc/src/n3iwf/n3iwf_handler/n3iwf_message"
	"free5gc/src/n3iwf/n3iwf_ike/ike_message"
	"net"

	gtpv1 "github.com/wmnsk/go-gtp/v1"
)

const (
	AmfUeNgapIdUnspecified int64 = 0xffffffffff
)

type N3IWFUe struct {
	/* UE identity*/
	RanUeNgapId           int64
	AmfUeNgapId           int64
	IPAddrv4              string
	IPAddrv6              string
	PortNumber            int32
	MaskedIMEISV          *ngapType.MaskedIMEISV // TS 38.413 9.3.1.54
	Guti                  string
	RRCEstablishmentCause int16
	IPSecInnerIP          string

	/* Relative Context */
	AMF *N3IWFAMF

	/* PDU Session */
	PduSessionList map[int64]*PDUSession // pduSessionId as key

	/* PDU Session Setup Temporary Data */
	TemporaryPDUSessionSetupData *PDUSessionSetupTemporaryData

	/* Security */
	Kn3iwf               []uint8                          // 32 bytes (256 bits), value is from NGAP IE "Security Key"
	SecurityCapabilities *ngapType.UESecurityCapabilities // TS 38.413 9.3.1.86

	/* IKE Security Association */
	N3IWFIKESecurityAssociation   *IKESecurityAssociation
	N3IWFChildSecurityAssociation *ChildSecurityAssociation

	/* NAS IKE Connection */
	UDPSendInfoGroup *n3iwf_message.UDPSendInfoGroup
	/* NAS TCP Connection */
	TCPConnection net.Conn

	/* Others */
	Guami                            *ngapType.GUAMI
	IndexToRfsp                      int64
	Ambr                             *ngapType.UEAggregateMaximumBitRate
	AllowedNssai                     *ngapType.AllowedNSSAI
	RadioCapability                  *ngapType.UERadioCapability                // TODO: This is for RRC, can be deleted
	CoreNetworkAssistanceInformation *ngapType.CoreNetworkAssistanceInformation // TS 38.413 9.3.1.15
	IMSVoiceSupported                int32
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
	GTPConnection                    *GTPConnectionInfo
	QFIList                          []uint8
	QosFlows                         map[int64]*QosFlow // QosFlowIdentifier as key
}

type PDUSessionSetupTemporaryData struct {
	// Slice of unactivated PDU session
	UnactivatedPDUSession []int64 // PDUSessionID as content
	// NGAPProcedureCode is used to identify which type of
	// response shall be used
	NGAPProcedureCode ngapType.ProcedureCode
	// PDU session setup list response
	SetupListCxtRes  *ngapType.PDUSessionResourceSetupListCxtRes
	FailedListCxtRes *ngapType.PDUSessionResourceFailedToSetupListCxtRes
	SetupListSURes   *ngapType.PDUSessionResourceSetupListSURes
	FailedListSURes  *ngapType.PDUSessionResourceFailedToSetupListSURes
}

type QosFlow struct {
	Identifier int64
	Parameters ngapType.QosFlowLevelQosParameters
}

type GTPConnectionInfo struct {
	UPFIPAddr           string
	UPFUDPAddr          net.Addr
	IncomingTEID        uint32
	OutgoingTEID        uint32
	UserPlaneConnection *gtpv1.UPlaneConn
}

type IKESecurityAssociation struct {
	// SPI
	RemoteSPI uint64
	LocalSPI  uint64

	// Message ID
	MessageID uint32

	// Transforms for IKE SA
	EncryptionAlgorithm    *ike_message.Transform
	PseudorandomFunction   *ike_message.Transform
	IntegrityAlgorithm     *ike_message.Transform
	DiffieHellmanGroup     *ike_message.Transform
	ExpandedSequenceNumber *ike_message.Transform

	// Used for key generating
	ConcatenatedNonce      []byte
	DiffieHellmanSharedKey []byte

	// Keys
	SK_d  []byte // used for child SA key deriving
	SK_ai []byte // used by initiator for integrity checking
	SK_ar []byte // used by responder for integrity checking
	SK_ei []byte // used by initiator for encrypting
	SK_er []byte // used by responder for encrypting
	SK_pi []byte // used by initiator for IKE authentication
	SK_pr []byte // used by responder for IKE authentication

	// State for IKE_AUTH
	State uint8

	// Temporary data stored for the use in later exchange
	InitiatorID              *ike_message.IdentificationInitiator
	InitiatorCertificate     *ike_message.Certificate
	IKEAuthResponseSA        *ike_message.SecurityAssociation
	TrafficSelectorInitiator *ike_message.TrafficSelectorInitiator
	TrafficSelectorResponder *ike_message.TrafficSelectorResponder
	LastEAPIdentifier        uint8

	// Authentication data
	LocalUnsignedAuthentication  []byte
	RemoteUnsignedAuthentication []byte

	// UE context
	ThisUE *N3IWFUe
}

type ChildSecurityAssociation struct {
	// SPI
	SPI uint32

	// IP address
	PeerPublicIPAddr  net.IP
	LocalPublicIPAddr net.IP

	// Traffic selector
	SelectedIPProtocol    uint8
	TrafficSelectorLocal  net.IPNet
	TrafficSelectorRemote net.IPNet

	// Security
	EncryptionAlgorithm               uint16
	InitiatorToResponderEncryptionKey []byte
	ResponderToInitiatorEncryptionKey []byte
	IntegrityAlgorithm                uint16
	InitiatorToResponderIntegrityKey  []byte
	ResponderToInitiatorIntegrityKey  []byte
	ESN                               bool

	// UE context
	ThisUE *N3IWFUe
}

func (ue *N3IWFUe) init() {
	ue.PduSessionList = make(map[int64]*PDUSession)
}

func (ue *N3IWFUe) Remove() {
	n3iwfSelf := N3IWFSelf()
	ue.DetachAMF()
	delete(n3iwfSelf.UePool, ue.RanUeNgapId)
}

func (ue *N3IWFUe) FindPDUSession(pduSessionID int64) *PDUSession {
	if pduSession, ok := ue.PduSessionList[pduSessionID]; ok {
		return pduSession
	} else {
		return nil
	}
}

func (ue *N3IWFUe) CreatePDUSession(pduSessionID int64, snssai ngapType.SNSSAI) (*PDUSession, error) {
	if _, exists := ue.PduSessionList[pduSessionID]; exists {
		return nil, fmt.Errorf("PDU Session[ID:%d] is already exists", pduSessionID)
	}
	pduSession := &PDUSession{
		Id:       pduSessionID,
		Snssai:   snssai,
		QosFlows: make(map[int64]*QosFlow),
	}
	ue.PduSessionList[pduSessionID] = pduSession
	return pduSession, nil
}

func (ue *N3IWFUe) CreateIKEChildSecurityAssociation(chosenSecurityAssociation *ike_message.SecurityAssociation) (*ChildSecurityAssociation, error) {
	childSecurityAssociation := new(ChildSecurityAssociation)

	if chosenSecurityAssociation == nil {
		return nil, errors.New("chosenSecurityAssociation is nil")
	}

	if len(chosenSecurityAssociation.Proposals) == 0 {
		return nil, errors.New("No proposal")
	}

	childSecurityAssociation.SPI = binary.BigEndian.Uint32(chosenSecurityAssociation.Proposals[0].SPI)

	if len(chosenSecurityAssociation.Proposals[0].EncryptionAlgorithm) != 0 {
		childSecurityAssociation.EncryptionAlgorithm = chosenSecurityAssociation.Proposals[0].EncryptionAlgorithm[0].TransformID
	}
	if len(chosenSecurityAssociation.Proposals[0].IntegrityAlgorithm) != 0 {
		childSecurityAssociation.IntegrityAlgorithm = chosenSecurityAssociation.Proposals[0].IntegrityAlgorithm[0].TransformID
	}
	if len(chosenSecurityAssociation.Proposals[0].ExtendedSequenceNumbers) != 0 {
		if chosenSecurityAssociation.Proposals[0].ExtendedSequenceNumbers[0].TransformID == 0 {
			childSecurityAssociation.ESN = false
		} else {
			childSecurityAssociation.ESN = true
		}
	}

	// Link UE context
	childSecurityAssociation.ThisUE = ue
	// Record to N3IWF context
	n3iwfContext.ChildSA[childSecurityAssociation.SPI] = childSecurityAssociation

	ue.N3IWFChildSecurityAssociation = childSecurityAssociation

	return childSecurityAssociation, nil
}

func (ue *N3IWFUe) AttachAMF(sctpAddr string) error {
	amf, err := N3IWFSelf().FindAMFBySCTPAddr(sctpAddr)
	if err != nil {
		return err
	}
	amf.N3iwfUeList[ue.RanUeNgapId] = ue
	ue.AMF = amf
	return nil
}
func (ue *N3IWFUe) DetachAMF() {
	amf := ue.AMF
	if amf == nil {
		return
	}
	delete(amf.N3iwfUeList, ue.RanUeNgapId)
}
