package context

import (
	"fmt"
	"math"
	"net"

	"github.com/pkg/errors"
	"github.com/vishvananda/netlink"

	ike_message "github.com/free5gc/ike/message"
	ike_security "github.com/free5gc/ike/security"
)

const (
	AmfUeNgapIdUnspecified int64 = 0xffffffffff
)

type N3IWFIkeUe struct {
	N3iwfCtx *N3IWFContext

	// UE identity
	IPSecInnerIP     net.IP
	IPSecInnerIPAddr *net.IPAddr // Used to send UP packets to UE

	// IKE Security Association
	N3IWFIKESecurityAssociation   *IKESecurityAssociation
	N3IWFChildSecurityAssociation map[uint32]*ChildSecurityAssociation // inbound SPI as key

	// Temporary Mapping of two SPIs
	// Exchange Message ID(including a SPI) and ChildSA(including a SPI)
	// Mapping of Message ID of exchange in IKE and Child SA when creating new child SA
	TemporaryExchangeMsgIDChildSAMapping map[uint32]*ChildSecurityAssociation // Message ID as a key

	// Security
	Kn3iwf []uint8 // 32 bytes (256 bits), value is from NGAP IE "Security Key"

	// NAS IKE Connection
	IKEConnection *UDPSocketInfo

	// Length of PDU Session List
	PduSessionListLen int
}

type IkeMsgTemporaryData struct {
	SecurityAssociation      *ike_message.SecurityAssociation
	TrafficSelectorInitiator *ike_message.TrafficSelectorInitiator
	TrafficSelectorResponder *ike_message.TrafficSelectorResponder
}

type IKESecurityAssociation struct {
	*ike_security.IKESAKey
	// SPI
	RemoteSPI uint64
	LocalSPI  uint64

	// Message ID
	InitiatorMessageID uint32
	ResponderMessageID uint32

	// Used for key generating
	ConcatenatedNonce []byte

	// State for IKE_AUTH
	State uint8

	// Temporary data stored for the use in later exchange
	InitiatorID              *ike_message.IdentificationInitiator
	InitiatorCertificate     *ike_message.Certificate
	IKEAuthResponseSA        *ike_message.SecurityAssociation
	TrafficSelectorInitiator *ike_message.TrafficSelectorInitiator
	TrafficSelectorResponder *ike_message.TrafficSelectorResponder
	LastEAPIdentifier        uint8

	// UDP Connection
	IKEConnection *UDPSocketInfo

	// Authentication data
	ResponderSignedOctets []byte
	InitiatorSignedOctets []byte

	// NAT detection
	UeBehindNAT    bool // If true, N3IWF should enable NAT traversal and
	N3iwfBehindNAT bool // TODO: If true, N3IWF should send UDP keepalive periodically

	// IKE UE context
	IkeUE *N3IWFIkeUe

	// Temporary store the receive ike message
	TemporaryIkeMsg *IkeMsgTemporaryData

	DPDReqRetransTimer *Timer // The time from sending the DPD request to receiving the response
	CurrentRetryTimes  int32  // Accumulate the number of times the DPD response wasn't received
	IKESAClosedCh      chan struct{}
	IsUseDPD           bool
}

func (ikeSA *IKESecurityAssociation) String() string {
	return "====== IKE Security Association Info =====" +
		"\nInitiator's SPI: " + fmt.Sprintf("%016x", ikeSA.RemoteSPI) +
		"\nResponder's SPI: " + fmt.Sprintf("%016x", ikeSA.LocalSPI) +
		"\nIKESAKey: " + ikeSA.IKESAKey.String()
}

// Temporary State Data Args
const (
	ArgsUEUDPConn string = "UE UDP Socket Info"
)

type ChildSecurityAssociation struct {
	// SPI
	InboundSPI  uint32 // N3IWF Specify
	OutboundSPI uint32 // Non-3GPP UE Specify

	// Associated XFRM interface
	XfrmIface netlink.Link

	XfrmStateList  []netlink.XfrmState
	XfrmPolicyList []netlink.XfrmPolicy

	// IP address
	PeerPublicIPAddr  net.IP
	LocalPublicIPAddr net.IP

	// Traffic selector
	SelectedIPProtocol    uint8
	TrafficSelectorLocal  net.IPNet
	TrafficSelectorRemote net.IPNet

	// Security
	*ike_security.ChildSAKey

	// Encapsulate
	EnableEncapsulate bool
	N3IWFPort         int
	NATPort           int

	// PDU Session IDs associated with this child SA
	PDUSessionIds []int64

	// IKE UE context
	IkeUE *N3IWFIkeUe

	LocalIsInitiator bool
}

func (childSA *ChildSecurityAssociation) String(xfrmiId uint32) string {
	var inboundEncryptionKey, inboundIntegrityKey, outboundEncryptionKey, outboundIntegrityKey []byte

	if childSA.LocalIsInitiator {
		inboundEncryptionKey = childSA.ResponderToInitiatorEncryptionKey
		inboundIntegrityKey = childSA.ResponderToInitiatorIntegrityKey
		outboundEncryptionKey = childSA.InitiatorToResponderEncryptionKey
		outboundIntegrityKey = childSA.InitiatorToResponderIntegrityKey
	} else {
		inboundEncryptionKey = childSA.InitiatorToResponderEncryptionKey
		inboundIntegrityKey = childSA.InitiatorToResponderIntegrityKey
		outboundEncryptionKey = childSA.ResponderToInitiatorEncryptionKey
		outboundIntegrityKey = childSA.ResponderToInitiatorIntegrityKey
	}

	return fmt.Sprintf("====== IPSec/Child SA Info ======"+
		"\n====== Inbound ======"+
		"\nXFRM interface if_id: %d"+
		"\nIPSec Inbound  SPI: 0x%08x"+
		"\n[UE:%+v] -> [N3IWF:%+v]"+
		"\nIPSec Encryption Algorithm: %d"+
		"\nIPSec Encryption Key: 0x%x"+
		"\nIPSec Integrity  Algorithm: %d"+
		"\nIPSec Integrity  Key: 0x%x"+
		"\n====== IPSec/Child SA Info ======"+
		"\n====== Outbound ======"+
		"\nXFRM interface if_id: %d"+
		"\nIPSec Outbound  SPI: 0x%08x"+
		"\n[N3IWF:%+v] -> [UE:%+v]"+
		"\nIPSec Encryption Algorithm: %d"+
		"\nIPSec Encryption Key: 0x%x"+
		"\nIPSec Integrity  Algorithm: %d"+
		"\nIPSec Integrity  Key: 0x%x",
		xfrmiId,
		childSA.InboundSPI,
		childSA.PeerPublicIPAddr,
		childSA.LocalPublicIPAddr,
		childSA.EncrKInfo.TransformID(),
		inboundEncryptionKey,
		childSA.IntegKInfo.TransformID(),
		inboundIntegrityKey,
		xfrmiId,
		childSA.OutboundSPI,
		childSA.LocalPublicIPAddr,
		childSA.PeerPublicIPAddr,
		childSA.EncrKInfo.TransformID(),
		outboundEncryptionKey,
		childSA.IntegKInfo.TransformID(),
		outboundIntegrityKey,
	)
}

type UDPSocketInfo struct {
	Conn      *net.UDPConn
	N3IWFAddr *net.UDPAddr
	UEAddr    *net.UDPAddr
}

func (ikeUe *N3IWFIkeUe) init() {
	ikeUe.N3IWFChildSecurityAssociation = make(map[uint32]*ChildSecurityAssociation)
	ikeUe.TemporaryExchangeMsgIDChildSAMapping = make(map[uint32]*ChildSecurityAssociation)
}

func (ikeUe *N3IWFIkeUe) Remove() error {
	if ikeUe.N3IWFIKESecurityAssociation.IsUseDPD {
		ikeUe.N3IWFIKESecurityAssociation.IKESAClosedCh <- struct{}{}
	}

	// remove from IKE UE context
	n3iwfCtx := ikeUe.N3iwfCtx
	n3iwfCtx.DeleteIKESecurityAssociation(ikeUe.N3IWFIKESecurityAssociation.LocalSPI)
	n3iwfCtx.DeleteInternalUEIPAddr(ikeUe.IPSecInnerIP.String())

	err := n3iwfCtx.IPSecInnerIPPool.Release(net.ParseIP(ikeUe.IPSecInnerIP.String()).To4())
	if err != nil {
		return errors.Wrapf(err, "N3IWFIkeUe Remove()")
	}

	for _, childSA := range ikeUe.N3IWFChildSecurityAssociation {
		if err := ikeUe.DeleteChildSA(childSA); err != nil {
			return err
		}
	}
	n3iwfCtx.DeleteIKEUe(ikeUe.N3IWFIKESecurityAssociation.LocalSPI)

	return nil
}

func (ikeUe *N3IWFIkeUe) DeleteChildSAXfrm(childSA *ChildSecurityAssociation) error {
	n3iwfCtx := ikeUe.N3iwfCtx
	iface := childSA.XfrmIface

	// Delete child SA xfrmState
	for idx := range childSA.XfrmStateList {
		xfrmState := childSA.XfrmStateList[idx]
		if err := netlink.XfrmStateDel(&xfrmState); err != nil {
			return errors.Wrapf(err, "Delete xfrmstate")
		}
	}
	// Delete child SA xfrmPolicy
	for idx := range childSA.XfrmPolicyList {
		xfrmPolicy := childSA.XfrmPolicyList[idx]
		if err := netlink.XfrmPolicyDel(&xfrmPolicy); err != nil {
			return errors.Wrapf(err, "Delete xfrmPolicy")
		}
	}

	if iface == nil || iface.Attrs().Name == "xfrmi-default" {
	} else if err := netlink.LinkDel(iface); err != nil {
		return errors.Wrapf(err, "Delete interface[%s]", iface.Attrs().Name)
	} else {
		ifId := childSA.XfrmStateList[0].Ifid
		if ifId < 0 || ifId > math.MaxUint32 {
			return errors.Errorf("DeleteChildSAXfrm Ifid has out of uint32 range value: %d", ifId)
		}
		n3iwfCtx.XfrmIfaces.Delete(uint32(ifId))
	}

	childSA.XfrmStateList = nil
	childSA.XfrmPolicyList = nil

	return nil
}

func (ikeUe *N3IWFIkeUe) DeleteChildSA(childSA *ChildSecurityAssociation) error {
	if err := ikeUe.DeleteChildSAXfrm(childSA); err != nil {
		return err
	}

	delete(ikeUe.N3IWFChildSecurityAssociation, childSA.InboundSPI)

	return nil
}

// When N3IWF send CREATE_CHILD_SA request to N3UE, the inbound SPI of childSA will be only stored first until
// receive response and call CompleteChildSAWithProposal to fill the all data of childSA
func (ikeUe *N3IWFIkeUe) CreateHalfChildSA(msgID, inboundSPI uint32, pduSessionID int64) {
	childSA := new(ChildSecurityAssociation)
	childSA.InboundSPI = inboundSPI
	childSA.PDUSessionIds = append(childSA.PDUSessionIds, pduSessionID)
	// Link UE context
	childSA.IkeUE = ikeUe
	// Map Exchange Message ID and Child SA data until get paired response
	ikeUe.TemporaryExchangeMsgIDChildSAMapping[msgID] = childSA
}

func (ikeUe *N3IWFIkeUe) CompleteChildSA(msgID uint32, outboundSPI uint32,
	chosenSecurityAssociation *ike_message.SecurityAssociation,
) (*ChildSecurityAssociation, error) {
	childSA, ok := ikeUe.TemporaryExchangeMsgIDChildSAMapping[msgID]

	if !ok {
		return nil, fmt.Errorf("there's not a half child SA created by the exchange with message ID %d", msgID)
	}

	// Remove mapping of exchange msg ID and child SA
	delete(ikeUe.TemporaryExchangeMsgIDChildSAMapping, msgID)

	if chosenSecurityAssociation == nil {
		return nil, errors.New("chosenSecurityAssociation is nil")
	}

	if len(chosenSecurityAssociation.Proposals) == 0 {
		return nil, errors.New("no proposal")
	}

	childSA.OutboundSPI = outboundSPI

	var err error
	childSA.ChildSAKey, err = ike_security.NewChildSAKeyByProposal(chosenSecurityAssociation.Proposals[0])
	if err != nil {
		return nil, errors.Wrapf(err, "CompleteChildSA")
	}

	// Record to UE context with inbound SPI as key
	ikeUe.N3IWFChildSecurityAssociation[childSA.InboundSPI] = childSA
	// Record to N3IWF context with inbound SPI as key
	ikeUe.N3iwfCtx.ChildSA.Store(childSA.InboundSPI, childSA)

	return childSA, nil
}
