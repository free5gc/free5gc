package test

import (
	"encoding/binary"
	"fmt"
	"math/big"
	"net"
	"strconv"
	"test/consumerTestdata/UDM/TestGenAuthData"
	"test/nasTestpacket"
	"testing"
	"time"

	"github.com/pkg/errors"

	"github.com/go-ping/ping"
	"github.com/stretchr/testify/assert"
	"github.com/vishvananda/netlink"
	"golang.org/x/sys/unix"

	"github.com/davecgh/go-spew/spew"

	"github.com/free5gc/ike"
	"github.com/free5gc/ike/message"
	ike_security "github.com/free5gc/ike/security"
	"github.com/free5gc/ngap/ngapType"

	ike_message "github.com/free5gc/ike/message"
	"github.com/free5gc/ike/security/dh"
	"github.com/free5gc/ike/security/encr"
	"github.com/free5gc/ike/security/integ"
	"github.com/free5gc/ike/security/prf"

	"github.com/free5gc/nas"
	"github.com/free5gc/nas/nasMessage"
	"github.com/free5gc/nas/nasType"
	nasSecurity "github.com/free5gc/nas/security"
	"github.com/free5gc/ngap"
	"github.com/free5gc/openapi/models"
	"github.com/free5gc/util/ueauth"
)

var (
	n3iwfInfo_IPSecIfaceAddr     = "192.168.127.1"
	n3ueInfo_IPSecIfaceAddr      = "192.168.127.2"
	n3ueInfo_SmPolicy_SNSSAI_SST = "1"
	n3ueInfo_SmPolicy_SNSSAI_SD  = "fedcba"
	n3ueInfo_IPSecIfaceName      = "veth3"
	n3ueInfo_XfrmiName           = "ipsec"
	n3ueInfo_XfrmiId             = uint32(1)
	n3ueInfo_GreIfaceName        = "gretun"
	ueInnerAddr                  = new(net.IPNet)
)

type N3IWFUe struct {
	N3IWFIkeUe
	N3IWFRanUe
}

type N3IWFIkeUe struct {
	/* UE identity */
	IPSecInnerIP     net.IP
	IPSecInnerIPAddr *net.IPAddr // Used to send UP packets to UE

	/* IKE Security Association */
	N3IWFIKESecurityAssociation   *IKESecurityAssociation
	N3IWFChildSecurityAssociation map[uint32]*ChildSecurityAssociation // inbound SPI as key

	/* Temporary Mapping of two SPIs */
	// Exchange Message ID(including a SPI) and ChildSA(including a SPI)
	// Mapping of Message ID of exchange in IKE and Child SA when creating new child SA
	TemporaryExchangeMsgIDChildSAMapping map[uint32]*ChildSecurityAssociation // Message ID as a key

	/* Security */
	Kn3iwf []uint8 // 32 bytes (256 bits), value is from NGAP IE "Security Key"

	// Length of PDU Session List
	PduSessionListLen int
}

type N3IWFRanUe struct {
	/* UE identity */
	RanUeNgapId  int64
	AmfUeNgapId  int64
	IPAddrv4     string
	IPAddrv6     string
	PortNumber   int32
	MaskedIMEISV *ngapType.MaskedIMEISV // TS 38.413 9.3.1.54
	Guti         string

	// UE send CREATE_CHILD_SA response
	TemporaryCachedNASMessage []byte

	/* NAS TCP Connection Established */
	IsNASTCPConnEstablished         bool
	IsNASTCPConnEstablishedComplete bool

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
	RRCEstablishmentCause            int16
	PduSessionReleaseList            ngapType.PDUSessionResourceReleasedListRelRes
}

type IKESecurityAssociation struct {
	*ike_security.IKESAKey
	// SPI
	RemoteSPI uint64
	LocalSPI  uint64

	// Message ID
	InitiatorMessageID uint32
	ResponderMessageID uint32

	// Authentication data
	ResponderSignedOctets []byte
	InitiatorSignedOctets []byte

	// Used for key generating
	ConcatenatedNonce []byte

	// State for IKE_AUTH
	State uint8

	// Temporary data stored for the use in later exchange
	IKEAuthResponseSA *message.SecurityAssociation
}

type ChildSecurityAssociation struct {
	*ike_security.ChildSAKey

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

	// Encapsulate
	EnableEncapsulate bool
	N3IWFPort         int
	NATPort           int
}

type PDUQoSInfo struct {
	pduSessionID    uint8
	qfiList         []uint8
	isDefault       bool
	isDSCPSpecified bool
	DSCP            uint8
}

func generateSPI(n3ue *N3IWFUe) ([]byte, error) {
	var spi uint32
	spiByte := make([]byte, 4)
	for {
		randomBigInt, err := ike_security.GenerateRandomNumber()
		if err != nil {
			return nil, errors.Wrapf(err, "GenerateSPI()")
		}
		randomUint64 := randomBigInt.Uint64()
		if _, ok := n3ue.N3IWFIkeUe.N3IWFChildSecurityAssociation[uint32(randomUint64)]; !ok {
			spi = uint32(randomUint64)
			binary.BigEndian.PutUint32(spiByte, spi)
			break
		}
	}
	return spiByte, nil
}

func setupIPsecXfrmi(xfrmIfaceName, parentIfaceName string, xfrmIfaceId uint32, xfrmIfaceAddr *net.IPNet) (netlink.Link, error) {
	var (
		xfrmi, parent netlink.Link
		err           error
	)

	if parent, err = netlink.LinkByName(parentIfaceName); err != nil {
		return nil, err
	}

	link := &netlink.Xfrmi{
		LinkAttrs: netlink.LinkAttrs{
			MTU:         1478,
			Name:        xfrmIfaceName,
			ParentIndex: parent.Attrs().Index,
		},
		Ifid: xfrmIfaceId,
	}

	// ip link add
	if err := netlink.LinkAdd(link); err != nil {
		return nil, err
	}

	if xfrmi, err = netlink.LinkByName(xfrmIfaceName); err != nil {
		return nil, err
	}

	// ip addr add
	linkIPSecAddr := &netlink.Addr{
		IPNet: xfrmIfaceAddr,
	}

	if err := netlink.AddrAdd(xfrmi, linkIPSecAddr); err != nil {
		return nil, err
	}

	// ip link set ... up
	if err := netlink.LinkSetUp(xfrmi); err != nil {
		return nil, err
	}

	return xfrmi, nil
}

func setupGreTunnel(greIfaceName, parentIfaceName string, ueTunnelAddr, n3iwfTunnelAddr, pduAddr net.IP, qoSInfo *PDUQoSInfo, t *testing.T) (netlink.Link, error) {
	var (
		parent      netlink.Link
		greKeyField uint32
		err         error
	)

	if qoSInfo != nil {
		greKeyField |= (uint32(qoSInfo.qfiList[0]) & 0x3F) << 24
	}

	if parent, err = netlink.LinkByName(parentIfaceName); err != nil {
		return nil, err
	}

	// New GRE tunnel interface
	newGRETunnel := &netlink.Gretun{
		LinkAttrs: netlink.LinkAttrs{
			Name: greIfaceName,
			MTU:  1438, // remain for endpoint IP header(most 40 bytes if IPv6) and ESP header (22 bytes)
		},
		Link:   uint32(parent.Attrs().Index), // PHYS_DEV in iproute2; IFLA_GRE_LINK in linux kernel
		Local:  ueTunnelAddr,
		Remote: n3iwfTunnelAddr,
		IKey:   greKeyField,
		OKey:   greKeyField,
	}

	t.Logf("GRE Key Field: 0x%x", greKeyField)

	if err := netlink.LinkAdd(newGRETunnel); err != nil {
		return nil, err
	}

	// Get link info
	linkGRE, err := netlink.LinkByName(greIfaceName)
	if err != nil {
		return nil, fmt.Errorf("No link named %s", greIfaceName)
	}

	linkGREAddr := &netlink.Addr{
		IPNet: &net.IPNet{
			IP:   pduAddr,
			Mask: net.IPv4Mask(255, 255, 255, 255),
		},
	}

	if err := netlink.AddrAdd(linkGRE, linkGREAddr); err != nil {
		return nil, err
	}

	// Set GRE interface up
	if err := netlink.LinkSetUp(linkGRE); err != nil {
		return nil, err
	}

	return linkGRE, nil
}

func getAuthSubscription() (authSubs models.AuthenticationSubscription) {
	authSubs.EncPermanentKey = TestGenAuthData.MilenageTestSet19.K
	authSubs.EncOpcKey = TestGenAuthData.MilenageTestSet19.OPC
	authSubs.AuthenticationManagementField = "8000"

	authSubs.SequenceNumber = &models.SequenceNumber{
		Sqn: TestGenAuthData.MilenageTestSet19.SQN,
	}
	authSubs.AuthenticationMethod = models.AuthMethod__5_G_AKA
	return
}

func setupUDPSocket() (*net.UDPConn, error) {
	bindAddr := n3ueInfo_IPSecIfaceAddr + ":500"
	udpAddr, err := net.ResolveUDPAddr("udp", bindAddr)
	if err != nil {
		return nil, fmt.Errorf("Resolve UDP address failed: %+v", err)
	}
	udpListener, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		return nil, fmt.Errorf("Resolve UDP address failed: %+v", err)
	}
	return udpListener, nil
}

// [TS 24502] 9.3.2.2.2 EAP-Response/5G-NAS message
// Define EAP-Response/5G-NAS message and AN-Parameters Format.

// [TS 24501] 8.2.6.1.1  REGISTRATION REQUEST message content
// For dealing with EAP-5G start, return EAP-5G response including
// "AN-Parameters and NASPDU of Registration Request"

func buildEAP5GANParameters() []byte {
	var anParameters []byte

	// [TS 24.502] 9.3.2.2.2.3
	// AN-parameter value field in GUAMI, PLMN ID and NSSAI is coded as value part
	// Therefore, IEI of AN-parameter is not needed to be included.

	// anParameter = AN-parameter Type | AN-parameter Length | Value part of IE

	// Build GUAMI
	anParameter := make([]byte, 2)
	guami := make([]byte, 6)
	guami[0] = 0x02
	guami[1] = 0xf8
	guami[2] = 0x39
	guami[3] = 0xca
	guami[4] = 0xfe
	guami[5] = 0x0
	anParameter[0] = ike_message.ANParametersTypeGUAMI
	anParameter[1] = byte(len(guami))
	anParameter = append(anParameter, guami...)

	anParameters = append(anParameters, anParameter...)

	// Build Establishment Cause
	anParameter = make([]byte, 2)
	establishmentCause := make([]byte, 1)
	establishmentCause[0] = ike_message.EstablishmentCauseMO_Signaling
	anParameter[0] = ike_message.ANParametersTypeEstablishmentCause
	anParameter[1] = byte(len(establishmentCause))
	anParameter = append(anParameter, establishmentCause...)

	anParameters = append(anParameters, anParameter...)

	// Build PLMN ID
	anParameter = make([]byte, 2)
	plmnID := make([]byte, 3)
	plmnID[0] = 0x02
	plmnID[1] = 0xf8
	plmnID[2] = 0x39
	anParameter[0] = ike_message.ANParametersTypeSelectedPLMNID
	anParameter[1] = byte(len(plmnID))
	anParameter = append(anParameter, plmnID...)

	anParameters = append(anParameters, anParameter...)

	// Build NSSAI
	anParameter = make([]byte, 2)
	var nssai []byte
	// s-nssai = s-nssai length(1 byte) | SST(1 byte) | SD(3 bytes)
	snssai := make([]byte, 5)
	snssai[0] = 4
	snssai[1] = 1
	snssai[2] = 0x01
	snssai[3] = 0x02
	snssai[4] = 0x03
	nssai = append(nssai, snssai...)
	snssai = make([]byte, 5)
	snssai[0] = 4
	snssai[1] = 1
	snssai[2] = 0x11
	snssai[3] = 0x22
	snssai[4] = 0x33
	nssai = append(nssai, snssai...)
	anParameter[0] = ike_message.ANParametersTypeRequestedNSSAI
	anParameter[1] = byte(len(nssai))
	anParameter = append(anParameter, nssai...)

	anParameters = append(anParameters, anParameter...)

	return anParameters
}

func parseIPAddressInformationToChildSecurityAssociation(
	childSecurityAssociation *ChildSecurityAssociation,
	trafficSelectorLocal *ike_message.IndividualTrafficSelector,
	trafficSelectorRemote *ike_message.IndividualTrafficSelector) error {

	if childSecurityAssociation == nil {
		return errors.New("childSecurityAssociation is nil")
	}

	childSecurityAssociation.PeerPublicIPAddr = net.ParseIP(n3iwfInfo_IPSecIfaceAddr)
	childSecurityAssociation.LocalPublicIPAddr = net.ParseIP(n3ueInfo_IPSecIfaceAddr)

	childSecurityAssociation.TrafficSelectorLocal = net.IPNet{
		IP:   trafficSelectorLocal.StartAddress,
		Mask: []byte{255, 255, 255, 255},
	}

	childSecurityAssociation.TrafficSelectorRemote = net.IPNet{
		IP:   trafficSelectorRemote.StartAddress,
		Mask: []byte{255, 255, 255, 255},
	}
	return nil
}

func parse5GQoSInfoNotify(n *ike_message.Notification) (info *PDUQoSInfo, err error) {
	info = new(PDUQoSInfo)
	var offset int = 0
	data := n.NotificationData
	dataLen := int(data[0])
	info.pduSessionID = data[1]
	qfiListLen := int(data[2])
	offset += (3 + qfiListLen)

	if offset > dataLen {
		return nil, errors.New("parse5GQoSInfoNotify err: Length and content of 5G-QoS-Info-Notify mismatch")
	}

	info.qfiList = make([]byte, qfiListLen)
	copy(info.qfiList, data[3:3+qfiListLen])

	info.isDefault = (data[offset] & ike_message.NotifyType5G_QOS_INFOBitDCSICheck) > 0
	info.isDSCPSpecified = (data[offset] & ike_message.NotifyType5G_QOS_INFOBitDSCPICheck) > 0

	return
}

type XFRMEncryptionAlgorithmType uint16

func (xfrmEncryptionAlgorithmType XFRMEncryptionAlgorithmType) String() string {
	switch xfrmEncryptionAlgorithmType {
	case ike_message.ENCR_DES:
		return "cbc(des)"
	case ike_message.ENCR_3DES:
		return "cbc(des3_ede)"
	case ike_message.ENCR_CAST:
		return "cbc(cast5)"
	case ike_message.ENCR_BLOWFISH:
		return "cbc(blowfish)"
	case ike_message.ENCR_NULL:
		return "ecb(cipher_null)"
	case ike_message.ENCR_AES_CBC:
		return "cbc(aes)"
	case ike_message.ENCR_AES_CTR:
		return "rfc3686(ctr(aes))"
	default:
		return ""
	}
}

type XFRMIntegrityAlgorithmType uint16

func (xfrmIntegrityAlgorithmType XFRMIntegrityAlgorithmType) String() string {
	switch xfrmIntegrityAlgorithmType {
	case ike_message.AUTH_HMAC_MD5_96:
		return "hmac(md5)"
	case ike_message.AUTH_HMAC_SHA1_96:
		return "hmac(sha1)"
	case ike_message.AUTH_AES_XCBC_96:
		return "xcbc(aes)"
	default:
		return ""
	}
}

func applyXFRMRule(ue_is_initiator bool, ifId uint32, childSecurityAssociation *ChildSecurityAssociation) error {
	// Build XFRM information data structure for incoming traffic.

	// Mark
	// mark := &netlink.XfrmMark{
	// 	Value: ifMark, // n3ueInfo.XfrmMark,
	// }

	// Direction: N3IWF -> UE
	// State
	var xfrmEncryptionAlgorithm, xfrmIntegrityAlgorithm *netlink.XfrmStateAlgo
	if ue_is_initiator {
		xfrmEncryptionAlgorithm = &netlink.XfrmStateAlgo{
			Name: XFRMEncryptionAlgorithmType(childSecurityAssociation.EncrKInfo.TransformID()).String(),
			Key:  childSecurityAssociation.ResponderToInitiatorEncryptionKey,
		}
		if childSecurityAssociation.IntegKInfo != nil {
			xfrmIntegrityAlgorithm = &netlink.XfrmStateAlgo{
				Name: XFRMIntegrityAlgorithmType(childSecurityAssociation.IntegKInfo.TransformID()).String(),
				Key:  childSecurityAssociation.ResponderToInitiatorIntegrityKey,
			}
		}
	} else {
		xfrmEncryptionAlgorithm = &netlink.XfrmStateAlgo{
			Name: XFRMEncryptionAlgorithmType(childSecurityAssociation.EncrKInfo.TransformID()).String(),
			Key:  childSecurityAssociation.InitiatorToResponderEncryptionKey,
		}
		if childSecurityAssociation.IntegKInfo != nil {
			xfrmIntegrityAlgorithm = &netlink.XfrmStateAlgo{
				Name: XFRMIntegrityAlgorithmType(childSecurityAssociation.IntegKInfo.TransformID()).String(),
				Key:  childSecurityAssociation.InitiatorToResponderIntegrityKey,
			}
		}
	}

	xfrmState := new(netlink.XfrmState)

	xfrmState.Src = childSecurityAssociation.PeerPublicIPAddr
	xfrmState.Dst = childSecurityAssociation.LocalPublicIPAddr
	xfrmState.Proto = netlink.XFRM_PROTO_ESP
	xfrmState.Mode = netlink.XFRM_MODE_TUNNEL
	xfrmState.Spi = int(childSecurityAssociation.InboundSPI)
	xfrmState.Ifid = int(ifId)
	xfrmState.Auth = xfrmIntegrityAlgorithm
	xfrmState.Crypt = xfrmEncryptionAlgorithm
	xfrmState.ESN = childSecurityAssociation.EsnInfo.GetNeedESN()

	// Commit xfrm state to netlink
	var err error
	if err = netlink.XfrmStateAdd(xfrmState); err != nil {
		return fmt.Errorf("Set XFRM state rule failed: %+v", err)
	}

	// Policy
	xfrmPolicyTemplate := netlink.XfrmPolicyTmpl{
		Src:   xfrmState.Src,
		Dst:   xfrmState.Dst,
		Proto: xfrmState.Proto,
		Mode:  xfrmState.Mode,
		Spi:   xfrmState.Spi,
	}

	xfrmPolicy := new(netlink.XfrmPolicy)

	if childSecurityAssociation.SelectedIPProtocol == 0 {
		return errors.New("Protocol == 0")
	}

	xfrmPolicy.Src = &childSecurityAssociation.TrafficSelectorRemote
	xfrmPolicy.Dst = &childSecurityAssociation.TrafficSelectorLocal
	xfrmPolicy.Proto = netlink.Proto(childSecurityAssociation.SelectedIPProtocol)
	xfrmPolicy.Dir = netlink.XFRM_DIR_IN
	xfrmPolicy.Ifid = int(ifId)
	xfrmPolicy.Tmpls = []netlink.XfrmPolicyTmpl{
		xfrmPolicyTemplate,
	}

	// Commit xfrm policy to netlink
	if err = netlink.XfrmPolicyAdd(xfrmPolicy); err != nil {
		return fmt.Errorf("Set XFRM policy rule failed: %+v", err)
	}

	// Direction: UE -> N3IWF
	// State
	if ue_is_initiator {
		xfrmEncryptionAlgorithm.Key = childSecurityAssociation.InitiatorToResponderEncryptionKey
		if childSecurityAssociation.IntegKInfo != nil {
			xfrmIntegrityAlgorithm.Key = childSecurityAssociation.InitiatorToResponderIntegrityKey
		}
	} else {
		xfrmEncryptionAlgorithm.Key = childSecurityAssociation.ResponderToInitiatorEncryptionKey
		if childSecurityAssociation.IntegKInfo != nil {
			xfrmIntegrityAlgorithm.Key = childSecurityAssociation.ResponderToInitiatorIntegrityKey
		}
	}

	xfrmState.Src, xfrmState.Dst = xfrmState.Dst, xfrmState.Src
	xfrmState.Spi = int(childSecurityAssociation.OutboundSPI)

	// Commit xfrm state to netlink
	if err = netlink.XfrmStateAdd(xfrmState); err != nil {
		return fmt.Errorf("Set XFRM state rule failed: %+v", err)
	}

	// Policy
	xfrmPolicyTemplate.Src, xfrmPolicyTemplate.Dst = xfrmPolicyTemplate.Dst, xfrmPolicyTemplate.Src
	xfrmPolicyTemplate.Spi = int(childSecurityAssociation.OutboundSPI)

	xfrmPolicy.Src, xfrmPolicy.Dst = xfrmPolicy.Dst, xfrmPolicy.Src
	xfrmPolicy.Dir = netlink.XFRM_DIR_OUT
	xfrmPolicy.Tmpls = []netlink.XfrmPolicyTmpl{
		xfrmPolicyTemplate,
	}

	// Commit xfrm policy to netlink
	if err = netlink.XfrmPolicyAdd(xfrmPolicy); err != nil {
		return fmt.Errorf("Set XFRM policy rule failed: %+v", err)
	}
	return nil
}

func sendPduSessionEstablishmentRequest(
	pduSessionId uint8,
	ue *RanUeContext,
	n3Info *N3IWFUe,
	ikeSA *IKESecurityAssociation,
	ikeConn *net.UDPConn,
	nasConn *net.TCPConn,
	t *testing.T) ([]netlink.Link, error) {

	var ifaces []netlink.Link

	// Build S-NSSA
	sst, err := strconv.ParseInt(n3ueInfo_SmPolicy_SNSSAI_SST, 16, 0)

	if err != nil {
		return ifaces, fmt.Errorf("Parse SST Fail:%+v", err)
	}

	sNssai := models.Snssai{
		Sst: int32(sst),
		Sd:  n3ueInfo_SmPolicy_SNSSAI_SD,
	}

	// PDU session establishment request
	// TS 24.501 9.11.3.47.1 Request type
	pdu := nasTestpacket.GetUlNasTransport_PduSessionEstablishmentRequest(pduSessionId, nasMessage.ULNASTransportRequestTypeInitialRequest, "internet", &sNssai)
	pdu, err = EncodeNasPduInEnvelopeWithSecurity(ue, pdu, nas.SecurityHeaderTypeIntegrityProtectedAndCiphered, true, false)
	if err != nil {
		return ifaces, fmt.Errorf("Encode NAS PDU In Envelope Fail:%+v", err)
	}
	if _, err = nasConn.Write(pdu); err != nil {
		return ifaces, fmt.Errorf("Send NAS Message Fail:%+v", err)
	}

	buffer := make([]byte, 65535)

	t.Logf("Waiting for N3IWF reply from IKE")

	// Receive N3IWF reply
	n, _, err := ikeConn.ReadFromUDP(buffer)
	if err != nil {
		return ifaces, fmt.Errorf("Read IKE Message Fail:%+v", err)
	}

	ikeMessage := new(ike_message.IKEMessage)
	ikeMessage.Payloads.Reset()
	ikeMessage, err = ike.DecodeDecrypt(buffer[:n], nil,
		ikeSA.IKESAKey, ike_message.Role_Initiator)
	if err != nil {
		t.Fatalf("Decode IKE meesage: %v", err)
	}

	var qoSInfo *PDUQoSInfo

	var responseSecurityAssociation *ike_message.SecurityAssociation
	var responseTrafficSelectorInitiator *ike_message.TrafficSelectorInitiator
	var responseTrafficSelectorResponder *ike_message.TrafficSelectorResponder
	var outboundSPI uint32
	var upIPAddr net.IP
	for _, ikePayload := range ikeMessage.Payloads {
		switch ikePayload.Type() {
		case ike_message.TypeSA:
			responseSecurityAssociation = ikePayload.(*ike_message.SecurityAssociation)
			outboundSPI = binary.BigEndian.Uint32(responseSecurityAssociation.Proposals[0].SPI)
		case ike_message.TypeTSi:
			responseTrafficSelectorInitiator = ikePayload.(*ike_message.TrafficSelectorInitiator)
		case ike_message.TypeTSr:
			responseTrafficSelectorResponder = ikePayload.(*ike_message.TrafficSelectorResponder)
		case ike_message.TypeN:
			notification := ikePayload.(*ike_message.Notification)
			if notification.NotifyMessageType == ike_message.Vendor3GPPNotifyType5G_QOS_INFO {
				t.Logf("Received Qos Flow settings")
				if info, err := parse5GQoSInfoNotify(notification); err == nil {
					qoSInfo = info
					t.Logf("NotificationData:%+v", notification.NotificationData)
					if qoSInfo.isDSCPSpecified {
						t.Logf("DSCP is specified but test not support")
					}
				} else {
					t.Logf("%+v", err)
				}
			}
			if notification.NotifyMessageType == ike_message.Vendor3GPPNotifyTypeUP_IP4_ADDRESS {
				upIPAddr = notification.NotificationData[:4]
				t.Logf("UP IP Address: %+v\n", upIPAddr)
			}
		case ike_message.TypeNiNr:
			responseNonce := ikePayload.(*ike_message.Nonce)
			ikeSA.ConcatenatedNonce = responseNonce.NonceData
		}
	}

	// IKE CREATE_CHILD_SA response
	ikeMessage.Payloads.Reset()
	n3Info.N3IWFIkeUe.N3IWFIKESecurityAssociation.ResponderMessageID = ikeMessage.MessageID

	var ikePayload ike_message.IKEPayloadContainer
	ikePayload.Reset()

	// SA
	inboundSPI, err := generateSPI(n3Info)
	if err != nil {
		t.Fatal(err)
	}
	responseSecurityAssociation.Proposals[0].SPI = inboundSPI
	ikePayload = append(ikePayload, responseSecurityAssociation)

	// TSi
	ikePayload = append(ikePayload, responseTrafficSelectorInitiator)

	// TSr
	ikePayload = append(ikePayload, responseTrafficSelectorResponder)

	// Nonce
	localNonceBigInt, err := ike_security.GenerateRandomNumber()
	if err != nil {
		t.Fatalf("Generate local nonce: %v", err)
	}
	localNonce := localNonceBigInt.Bytes()
	ikeSA.ConcatenatedNonce = append(ikeSA.ConcatenatedNonce, localNonce...)
	ikePayload.BuildNonce(localNonce)

	ikeMessage = ike_message.NewMessage(
		ikeSA.LocalSPI,
		ikeSA.RemoteSPI,
		ike_message.CREATE_CHILD_SA,
		true, true,
		ikeSA.InitiatorMessageID,
		ikePayload,
	)

	ikeMessageData, err := ike.EncodeEncrypt(ikeMessage, ikeSA.IKESAKey,
		ike_message.Role_Initiator)
	if err != nil {
		t.Fatalf("EncodeEncrypt IKE message failed: %+v", err)
	}

	n3iwfUDPAddr, err := net.ResolveUDPAddr("udp", n3iwfInfo_IPSecIfaceAddr+":500")

	if err != nil {
		return ifaces, fmt.Errorf("Resolve N3IWF IPSec IP Addr Fail:%+v", err)
	}

	_, err = ikeConn.WriteToUDP(ikeMessageData, n3iwfUDPAddr)
	if err != nil {
		t.Errorf("Write IKE maessage fail: %+v", err)
		return ifaces, err
	}

	n3Info.N3IWFIkeUe.CreateHalfChildSA(n3Info.N3IWFIkeUe.N3IWFIKESecurityAssociation.ResponderMessageID, binary.BigEndian.Uint32(inboundSPI), int64(pduSessionId))
	childSecurityAssociationContextUserPlane, err := n3Info.N3IWFIkeUe.CompleteChildSA(
		n3Info.N3IWFIkeUe.N3IWFIKESecurityAssociation.ResponderMessageID, outboundSPI, responseSecurityAssociation)
	if err != nil {
		return ifaces, fmt.Errorf("Create child security association context failed: %+v", err)
	}

	err = parseIPAddressInformationToChildSecurityAssociation(
		childSecurityAssociationContextUserPlane,
		responseTrafficSelectorResponder.TrafficSelectors[0],
		responseTrafficSelectorInitiator.TrafficSelectors[0])

	if err != nil {
		return ifaces, fmt.Errorf("Parse IP address to child security association failed: %+v", err)
	}
	// Select GRE traffic
	childSecurityAssociationContextUserPlane.SelectedIPProtocol = unix.IPPROTO_GRE

	if err := childSecurityAssociationContextUserPlane.GenerateKeyForChildSA(ikeSA.IKESAKey,
		ikeSA.ConcatenatedNonce); err != nil {
		return ifaces, fmt.Errorf("Generate key for child SA failed: %+v", err)
	}

	// ====== Inbound ======
	t.Logf("====== IPSec/Child SA for 3GPP UP Inbound =====")
	t.Logf("[UE:%+v] <- [N3IWF:%+v]",
		childSecurityAssociationContextUserPlane.LocalPublicIPAddr, childSecurityAssociationContextUserPlane.PeerPublicIPAddr)
	t.Logf("IPSec SPI: 0x%016x", childSecurityAssociationContextUserPlane.InboundSPI)
	t.Logf("IPSec Encryption Algorithm: %d", childSecurityAssociationContextUserPlane.EncrKInfo.TransformID())
	t.Logf("IPSec Encryption Key: 0x%x", childSecurityAssociationContextUserPlane.InitiatorToResponderEncryptionKey)
	t.Logf("IPSec Integrity  Algorithm: %d", childSecurityAssociationContextUserPlane.IntegKInfo.TransformID())
	t.Logf("IPSec Integrity  Key: 0x%x", childSecurityAssociationContextUserPlane.InitiatorToResponderIntegrityKey)
	// ====== Outbound ======
	t.Logf("====== IPSec/Child SA for 3GPP UP Outbound =====")
	t.Logf("[UE:%+v] -> [N3IWF:%+v]",
		childSecurityAssociationContextUserPlane.LocalPublicIPAddr, childSecurityAssociationContextUserPlane.PeerPublicIPAddr)
	t.Logf("IPSec SPI: 0x%016x", childSecurityAssociationContextUserPlane.OutboundSPI)
	t.Logf("IPSec Encryption Algorithm: %d", childSecurityAssociationContextUserPlane.EncrKInfo.TransformID())
	t.Logf("IPSec Encryption Key: 0x%x", childSecurityAssociationContextUserPlane.ResponderToInitiatorEncryptionKey)
	t.Logf("IPSec Integrity  Algorithm: %d", childSecurityAssociationContextUserPlane.IntegKInfo.TransformID())
	t.Logf("IPSec Integrity  Key: 0x%x", childSecurityAssociationContextUserPlane.ResponderToInitiatorIntegrityKey)
	t.Logf("State function: encr: %d, auth: %d", childSecurityAssociationContextUserPlane.EncrKInfo.TransformID(),
		childSecurityAssociationContextUserPlane.IntegKInfo.TransformID())

	// Aplly XFRM rules
	n3ueInfo_XfrmiId++
	err = applyXFRMRule(false, n3ueInfo_XfrmiId, childSecurityAssociationContextUserPlane)

	if err != nil {
		t.Errorf("Applying XFRM rules failed: %+v", err)
		return ifaces, err
	}

	var linkIPSec netlink.Link

	// Setup interface for ipsec
	newXfrmiName := fmt.Sprintf("%s-%d", n3ueInfo_XfrmiName, n3ueInfo_XfrmiId)
	if linkIPSec, err = setupIPsecXfrmi(newXfrmiName, n3ueInfo_IPSecIfaceName, n3ueInfo_XfrmiId, ueInnerAddr); err != nil {
		return ifaces, fmt.Errorf("Setup XFRMi interface %s fail: %+v", newXfrmiName, err)
	}

	ifaces = append(ifaces, linkIPSec)

	t.Logf("Setup XFRM interface %s successfully", newXfrmiName)

	var pduAddr net.IP

	// Read NAS from N3IWF
	if n, err := nasConn.Read(buffer); err != nil {
		return ifaces, fmt.Errorf("Read NAS Message Fail:%+v", err)
	} else {
		nasMsg, err := DecodePDUSessionEstablishmentAccept(ue, n, buffer)
		if err != nil {
			t.Errorf("DecodePDUSessionEstablishmentAccept Fail: %+v", err)
		}
		spew.Config.Indent = "\t"
		nasStr := spew.Sdump(nasMsg)
		t.Log("Dump DecodePDUSessionEstablishmentAccept:\n", nasStr)

		pduAddr, err = GetPDUAddress(nasMsg.GsmMessage.PDUSessionEstablishmentAccept)
		if err != nil {
			t.Errorf("GetPDUAddress Fail: %+v", err)
		}

		t.Logf("PDU Address: %s", pduAddr.String())
	}

	var linkGRE netlink.Link

	newGREName := fmt.Sprintf("%s-id-%d", n3ueInfo_GreIfaceName, n3ueInfo_XfrmiId)

	if linkGRE, err = setupGreTunnel(newGREName, newXfrmiName, ueInnerAddr.IP, upIPAddr, pduAddr, qoSInfo, t); err != nil {
		return ifaces, fmt.Errorf("Setup GRE tunnel %s Fail %+v", newGREName, err)
	}

	ifaces = append(ifaces, linkGRE)

	return ifaces, nil
}

func TestNon3GPPUE(t *testing.T) {
	// New UE
	ue := NewRanUeContext("imsi-208930000007487", 1, nasSecurity.AlgCiphering128NEA0, nasSecurity.AlgIntegrity128NIA2,
		models.AccessType_NON_3_GPP_ACCESS)
	ue.AmfUeNgapId = 1
	ue.AuthenticationSubs = getAuthSubscription()
	mobileIdentity5GS := nasType.MobileIdentity5GS{
		Len:    13, // suci
		Buffer: []uint8{0x01, 0x02, 0xf8, 0x39, 0xf0, 0xff, 0x00, 0x00, 0x00, 0x00, 0x00, 0x47, 0x78},
	}

	// Used to save IPsec/IKE related data
	n3ue := new(N3IWFUe)
	n3ue.N3IWFIkeUe.N3IWFChildSecurityAssociation = make(map[uint32]*ChildSecurityAssociation)
	n3ue.N3IWFIkeUe.TemporaryExchangeMsgIDChildSAMapping = make(map[uint32]*ChildSecurityAssociation)

	n3iwfUDPAddr, err := net.ResolveUDPAddr("udp", n3iwfInfo_IPSecIfaceAddr+":500")
	if err != nil {
		t.Fatalf("Resolve UDP address %s fail: %+v", n3iwfInfo_IPSecIfaceAddr+":500", err)
	}
	udpConnection, err := setupUDPSocket()

	if err != nil {
		t.Fatalf("Setup UDP socket Fail: %+v", err)
	}

	// IKE_SA_INIT
	ikeInitiatorSPI := uint64(123123)
	payload := new(ike_message.IKEPayloadContainer)

	// Security Association
	securityAssociation := payload.BuildSecurityAssociation()
	// Proposal 1
	proposal := securityAssociation.Proposals.BuildProposal(1, ike_message.TypeIKE, nil)
	// ENCR
	var attributeType uint16 = ike_message.AttributeTypeKeyLength
	var keyLength uint16 = 256
	proposal.EncryptionAlgorithm.BuildTransform(ike_message.TypeEncryptionAlgorithm, ike_message.ENCR_AES_CBC, &attributeType, &keyLength, nil)
	// INTEG
	proposal.IntegrityAlgorithm.BuildTransform(ike_message.TypeIntegrityAlgorithm, ike_message.AUTH_HMAC_SHA1_96, nil, nil, nil)
	// PRF
	proposal.PseudorandomFunction.BuildTransform(ike_message.TypePseudorandomFunction, ike_message.PRF_HMAC_SHA1, nil, nil, nil)
	// DH
	proposal.DiffieHellmanGroup.BuildTransform(ike_message.TypeDiffieHellmanGroup, ike_message.DH_2048_BIT_MODP, nil, nil, nil)

	// Key exchange data
	generator := new(big.Int).SetUint64(dh.Group14Generator)
	factor, ok := new(big.Int).SetString(dh.Group14PrimeString, 16)
	if !ok {
		t.Fatalf("Generate key exchange data failed")
	}
	secert, err := ike_security.GenerateRandomNumber()
	if err != nil {
		t.Fatalf("Generate secert: %v", err)
	}
	localPublicKeyExchangeValue := new(big.Int).Exp(generator, secert, factor).Bytes()
	prependZero := make([]byte, len(factor.Bytes())-len(localPublicKeyExchangeValue))
	localPublicKeyExchangeValue = append(prependZero, localPublicKeyExchangeValue...)
	payload.BUildKeyExchange(ike_message.DH_2048_BIT_MODP, localPublicKeyExchangeValue)

	// Nonce
	localNonceBigInt, err := ike_security.GenerateRandomNumber()
	if err != nil {
		t.Fatalf("Generate localNonce : %v", err)
	}
	localNonce := localNonceBigInt.Bytes()
	payload.BuildNonce(localNonce)

	ikeMessage := ike_message.NewMessage(ikeInitiatorSPI, 0, ike_message.IKE_SA_INIT,
		false, true, 0, *payload)
	// Send to N3IWF
	ikeMessageData, err := ike.EncodeEncrypt(ikeMessage, nil, ike_message.Role_Initiator)
	if err != nil {
		t.Fatalf("Encode IKE Message fail: %+v", err)
	}
	if _, err := udpConnection.WriteToUDP(ikeMessageData, n3iwfUDPAddr); err != nil {
		t.Fatalf("Write IKE maessage fail: %+v", err)
	}
	realMessage1, _ := ikeMessage.Encode()
	ikeSecurityAssociation := &IKESecurityAssociation{
		ResponderSignedOctets: realMessage1,
	}

	// Receive N3IWF reply
	buffer := make([]byte, 65535)
	n, _, err := udpConnection.ReadFromUDP(buffer)
	if err != nil {
		t.Fatalf("Read IKE Message fail: %+v", err)
	}
	ikeMessage.Payloads.Reset()
	err = ikeMessage.Decode(buffer[:n])
	if err != nil {
		t.Fatalf("Decode IKE Message fail: %+v", err)
	}

	var sharedKeyExchangeData []byte
	var remoteNonce []byte

	for _, ikePayload := range ikeMessage.Payloads {
		switch ikePayload.Type() {
		case ike_message.TypeSA:
			t.Log("Get SA payload")
		case ike_message.TypeKE:
			remotePublicKeyExchangeValue := ikePayload.(*ike_message.KeyExchange).KeyExchangeData
			var i int = 0
			for {
				if remotePublicKeyExchangeValue[i] != 0 {
					break
				}
			}
			remotePublicKeyExchangeValue = remotePublicKeyExchangeValue[i:]
			remotePublicKeyExchangeValueBig := new(big.Int).SetBytes(remotePublicKeyExchangeValue)
			sharedKeyExchangeData = new(big.Int).Exp(remotePublicKeyExchangeValueBig, secert, factor).Bytes()
		case ike_message.TypeNiNr:
			remoteNonce = ikePayload.(*ike_message.Nonce).NonceData
		}
	}

	ikeSecurityAssociation = &IKESecurityAssociation{
		LocalSPI:           ikeInitiatorSPI,
		RemoteSPI:          ikeMessage.ResponderSPI,
		InitiatorMessageID: 0,
		ResponderMessageID: 0,
		IKESAKey: &ike_security.IKESAKey{
			EncrInfo:  encr.DecodeTransform(proposal.EncryptionAlgorithm[0]),
			IntegInfo: integ.DecodeTransform(proposal.IntegrityAlgorithm[0]),
			PrfInfo:   prf.DecodeTransform(proposal.PseudorandomFunction[0]),
			DhInfo:    dh.DecodeTransform(proposal.DiffieHellmanGroup[0]),
		},
		ConcatenatedNonce:     append(localNonce, remoteNonce...),
		ResponderSignedOctets: append(ikeSecurityAssociation.ResponderSignedOctets, remoteNonce...),
	}

	err = ikeSecurityAssociation.IKESAKey.GenerateKeyForIKESA(ikeSecurityAssociation.ConcatenatedNonce,
		sharedKeyExchangeData, ikeSecurityAssociation.LocalSPI, ikeSecurityAssociation.RemoteSPI)
	if err != nil {
		t.Fatalf("Generate key for IKE SA failed: %+v", err)
	}

	n3ue.N3IWFIkeUe.N3IWFIKESecurityAssociation = ikeSecurityAssociation

	// IKE_AUTH
	ikeMessage.Payloads.Reset()
	ikeSecurityAssociation.InitiatorMessageID++

	var ikePayload ike_message.IKEPayloadContainer

	// Identification
	ikePayload.BuildIdentificationInitiator(ike_message.ID_KEY_ID, []byte("UE"))

	// Security Association
	securityAssociation = ikePayload.BuildSecurityAssociation()
	// Proposal 1
	inboundSPI, err := generateSPI(n3ue)
	if err != nil {
		t.Fatal(err)
	}
	proposal = securityAssociation.Proposals.BuildProposal(1, ike_message.TypeESP, inboundSPI)
	// ENCR
	proposal.EncryptionAlgorithm.BuildTransform(ike_message.TypeEncryptionAlgorithm, ike_message.ENCR_AES_CBC, &attributeType, &keyLength, nil)
	// INTEG
	proposal.IntegrityAlgorithm.BuildTransform(ike_message.TypeIntegrityAlgorithm, ike_message.AUTH_HMAC_SHA1_96, nil, nil, nil)
	// ESN
	proposal.ExtendedSequenceNumbers.BuildTransform(ike_message.TypeExtendedSequenceNumbers, ike_message.ESN_DISABLE, nil, nil, nil)

	// Traffic Selector
	tsi := ikePayload.BuildTrafficSelectorInitiator()
	tsi.TrafficSelectors.BuildIndividualTrafficSelector(ike_message.TS_IPV4_ADDR_RANGE, 0, 0, 65535, []byte{0, 0, 0, 0}, []byte{255, 255, 255, 255})
	tsr := ikePayload.BuildTrafficSelectorResponder()
	tsr.TrafficSelectors.BuildIndividualTrafficSelector(ike_message.TS_IPV4_ADDR_RANGE, 0, 0, 65535, []byte{0, 0, 0, 0}, []byte{255, 255, 255, 255})

	ikeMessage = ike_message.NewMessage(
		ikeSecurityAssociation.LocalSPI,
		ikeSecurityAssociation.RemoteSPI,
		ike_message.IKE_AUTH, false, true,
		ikeSecurityAssociation.InitiatorMessageID,
		ikePayload,
	)

	// Send to N3IWF
	ikeMessageData, err = ike.EncodeEncrypt(ikeMessage, ikeSecurityAssociation.IKESAKey,
		ike_message.Role_Initiator)
	if err != nil {
		t.Fatalf("EncodeEncrypt IKE message failed: %+v", err)
	}
	if _, err := udpConnection.WriteToUDP(ikeMessageData, n3iwfUDPAddr); err != nil {
		t.Fatalf("Write IKE message failed: %+v", err)
	}

	n3ue.N3IWFIkeUe.CreateHalfChildSA(ikeSecurityAssociation.InitiatorMessageID,
		binary.BigEndian.Uint32(inboundSPI), -1)

	// Receive N3IWF reply
	n, _, err = udpConnection.ReadFromUDP(buffer)
	if err != nil {
		t.Fatalf("Read IKE message failed: %+v", err)
	}
	ikeMessage.Payloads.Reset()

	ikeMessage, err = ike.DecodeDecrypt(buffer[:n], nil,
		ikeSecurityAssociation.IKESAKey, ike_message.Role_Initiator)
	if err != nil {
		t.Fatalf("Decode IKE meesage: %v", err)
	}

	var eapIdentifier uint8

	for _, ikePayload := range ikeMessage.Payloads {
		switch ikePayload.Type() {
		case ike_message.TypeIDr:
			t.Log("Get IDr")
		case ike_message.TypeAUTH:
			t.Log("Get AUTH")
		case ike_message.TypeCERT:
			t.Log("Get CERT")
		case ike_message.TypeEAP:
			eapIdentifier = ikePayload.(*ike_message.EAP).Identifier
			t.Log("Get EAP")
		}
	}

	// IKE_AUTH - EAP exchange
	ikeMessage.Payloads.Reset()
	ikeSecurityAssociation.InitiatorMessageID++

	ikePayload.Reset()

	// EAP-5G vendor type data
	eapVendorTypeData := make([]byte, 2)
	eapVendorTypeData[0] = ike_message.EAP5GType5GNAS

	// AN Parameters
	anParameters := buildEAP5GANParameters()
	anParametersLength := make([]byte, 2)
	binary.BigEndian.PutUint16(anParametersLength, uint16(len(anParameters)))
	eapVendorTypeData = append(eapVendorTypeData, anParametersLength...)
	eapVendorTypeData = append(eapVendorTypeData, anParameters...)

	// NAS
	ueSecurityCapability := ue.GetUESecurityCapability()
	registrationRequest := nasTestpacket.GetRegistrationRequest(nasMessage.RegistrationType5GSInitialRegistration,
		mobileIdentity5GS, nil, ueSecurityCapability, nil, nil, nil)

	nasLength := make([]byte, 2)
	binary.BigEndian.PutUint16(nasLength, uint16(len(registrationRequest)))
	eapVendorTypeData = append(eapVendorTypeData, nasLength...)
	eapVendorTypeData = append(eapVendorTypeData, registrationRequest...)

	eap := ikePayload.BuildEAP(ike_message.EAPCodeResponse, eapIdentifier)
	eap.EAPTypeData.BuildEAPExpanded(ike_message.VendorID3GPP, ike_message.VendorTypeEAP5G, eapVendorTypeData)

	ikeMessage = ike_message.NewMessage(
		ikeSecurityAssociation.LocalSPI,
		ikeSecurityAssociation.RemoteSPI,
		ike_message.IKE_AUTH,
		false, true,
		ikeSecurityAssociation.InitiatorMessageID,
		ikePayload,
	)

	// Send to N3IWF
	ikeMessageData, err = ike.EncodeEncrypt(ikeMessage, ikeSecurityAssociation.IKESAKey,
		ike_message.Role_Initiator)
	if err != nil {
		t.Fatalf("EncodeEncrypt IKE message failed: %+v", err)
	}
	if _, err := udpConnection.WriteToUDP(ikeMessageData, n3iwfUDPAddr); err != nil {
		t.Fatalf("Write IKE message failed: %+v", err)
	}

	// Receive N3IWF reply
	n, _, err = udpConnection.ReadFromUDP(buffer)
	if err != nil {
		t.Fatalf("Read IKE message failed: %+v", err)
	}

	ikeMessage.Payloads.Reset()

	ikeMessage, err = ike.DecodeDecrypt(buffer[:n], nil,
		ikeSecurityAssociation.IKESAKey, ike_message.Role_Initiator)
	if err != nil {
		t.Fatalf("Decode IKE meesage: %v", err)
	}

	var eapReq *ike_message.EAP
	var eapExpanded *ike_message.EAPExpanded

	eapReq, ok = ikeMessage.Payloads[0].(*ike_message.EAP)
	if !ok {
		t.Fatalf("Received packet is not an EAP payload")
	}

	var decodedNAS *nas.Message

	eapExpanded, ok = eapReq.EAPTypeData[0].(*ike_message.EAPExpanded)
	if !ok {
		t.Fatalf("The EAP data is not an EAP expended.")
	}

	// Decode NAS - Authentication Request
	nasData := eapExpanded.VendorData[4:]
	decodedNAS = new(nas.Message)
	if err := decodedNAS.PlainNasDecode(&nasData); err != nil {
		t.Fatalf("Decode plain NAS fail: %+v", err)
	}

	// Calculate for RES*
	assert.NotNil(t, decodedNAS)
	rand := decodedNAS.AuthenticationRequest.GetRANDValue()
	resStat := ue.DeriveRESstarAndSetKey(ue.AuthenticationSubs, rand[:], "5G:mnc093.mcc208.3gppnetwork.org")

	// send NAS Authentication Response
	pdu := nasTestpacket.GetAuthenticationResponse(resStat, "")

	// IKE_AUTH - EAP exchange
	ikeMessage.Payloads.Reset()
	ikeSecurityAssociation.InitiatorMessageID++

	ikePayload.Reset()

	// EAP-5G vendor type data
	eapVendorTypeData = make([]byte, 4)
	eapVendorTypeData[0] = ike_message.EAP5GType5GNAS

	// NAS - Authentication Response
	nasLength = make([]byte, 2)
	binary.BigEndian.PutUint16(nasLength, uint16(len(pdu)))
	eapVendorTypeData = append(eapVendorTypeData, nasLength...)
	eapVendorTypeData = append(eapVendorTypeData, pdu...)

	eap = ikePayload.BuildEAP(ike_message.EAPCodeResponse, eapReq.Identifier)
	eap.EAPTypeData.BuildEAPExpanded(ike_message.VendorID3GPP, ike_message.VendorTypeEAP5G, eapVendorTypeData)

	ikeMessage = ike_message.NewMessage(
		ikeSecurityAssociation.LocalSPI,
		ikeSecurityAssociation.RemoteSPI,
		ike_message.IKE_AUTH,
		false, true,
		ikeSecurityAssociation.InitiatorMessageID,
		ikePayload,
	)
	// Send to N3IWF
	ikeMessageData, err = ike.EncodeEncrypt(ikeMessage, ikeSecurityAssociation.IKESAKey,
		ike_message.Role_Initiator)
	if err != nil {
		t.Fatalf("EncodeEncrypt IKE message failed: %+v", err)
	}
	_, err = udpConnection.WriteToUDP(ikeMessageData, n3iwfUDPAddr)
	if err != nil {
		t.Fatalf("Write IKE message failed: %+v", err)
	}

	// Receive N3IWF reply
	n, _, err = udpConnection.ReadFromUDP(buffer)
	if err != nil {
		t.Fatalf("Read IKE Message fail: %+v", err)
	}
	ikeMessage.Payloads.Reset()
	ikeMessage, err = ike.DecodeDecrypt(buffer[:n], nil,
		ikeSecurityAssociation.IKESAKey, ike_message.Role_Initiator)
	if err != nil {
		t.Fatalf("Decode IKE meesage: %v", err)
	}

	eapReq, ok = ikeMessage.Payloads[0].(*ike_message.EAP)
	if !ok {
		t.Fatal("Received packet is not an EAP payload")
		return
	}
	eapExpanded, ok = eapReq.EAPTypeData[0].(*ike_message.EAPExpanded)
	if !ok {
		t.Fatal("Received packet is not an EAP expended payload")
		return
	}

	nasData = eapExpanded.VendorData[4:]

	// Send NAS Security Mode Complete Msg
	registrationRequestWith5GMM := nasTestpacket.GetRegistrationRequest(nasMessage.RegistrationType5GSInitialRegistration,
		mobileIdentity5GS, nil, ueSecurityCapability, ue.Get5GMMCapability(), nil, nil)
	pdu = nasTestpacket.GetSecurityModeComplete(registrationRequestWith5GMM)
	pdu, err = EncodeNasPduWithSecurity(ue, pdu, nas.SecurityHeaderTypeIntegrityProtectedAndCipheredWithNew5gNasSecurityContext, true, true)
	assert.Nil(t, err)

	// IKE_AUTH - EAP exchange
	ikeMessage.Payloads.Reset()
	ikeSecurityAssociation.InitiatorMessageID++

	ikePayload.Reset()

	// EAP-5G vendor type data
	eapVendorTypeData = make([]byte, 4)
	eapVendorTypeData[0] = ike_message.EAP5GType5GNAS

	// NAS - Authentication Response
	nasLength = make([]byte, 2)
	binary.BigEndian.PutUint16(nasLength, uint16(len(pdu)))
	eapVendorTypeData = append(eapVendorTypeData, nasLength...)
	eapVendorTypeData = append(eapVendorTypeData, pdu...)

	eap = ikePayload.BuildEAP(ike_message.EAPCodeResponse, eapReq.Identifier)
	eap.EAPTypeData.BuildEAPExpanded(ike_message.VendorID3GPP, ike_message.VendorTypeEAP5G, eapVendorTypeData)

	ikeMessage = ike_message.NewMessage(
		ikeSecurityAssociation.LocalSPI,
		ikeSecurityAssociation.RemoteSPI,
		ike_message.IKE_AUTH,
		false, true,
		ikeSecurityAssociation.InitiatorMessageID,
		ikePayload,
	)

	// Send to N3IWF
	ikeMessageData, err = ike.EncodeEncrypt(ikeMessage, ikeSecurityAssociation.IKESAKey,
		ike_message.Role_Initiator)
	if err != nil {
		t.Fatalf("EncodeEncrypt IKE message failed: %+v", err)
	}
	_, err = udpConnection.WriteToUDP(ikeMessageData, n3iwfUDPAddr)
	if err != nil {
		t.Fatalf("Write IKE message failed: %+v", err)
	}

	// Receive N3IWF reply
	n, _, err = udpConnection.ReadFromUDP(buffer)
	if err != nil {
		t.Fatalf("Read IKE Message fail: %+v", err)
		return
	}
	ikeMessage.Payloads.Reset()

	ikeMessage, err = ike.DecodeDecrypt(buffer[:n], nil,
		ikeSecurityAssociation.IKESAKey, ike_message.Role_Initiator)
	if err != nil {
		t.Fatalf("Decode IKE meesage: %v", err)
	}

	eapReq, ok = ikeMessage.Payloads[0].(*ike_message.EAP)
	if !ok {
		t.Fatal("Received packet is not an EAP payload")
	}
	if eapReq.Code != ike_message.EAPCodeSuccess {
		t.Fatal("Not Success")
	}

	// IKE_AUTH - Authentication
	ikeMessage.Payloads.Reset()
	ikeSecurityAssociation.InitiatorMessageID++

	ikePayload.Reset()

	// Authentication
	// Derive Kn3iwf
	P0 := make([]byte, 4)
	binary.BigEndian.PutUint32(P0, ue.ULCount.Get()-1)
	L0 := ueauth.KDFLen(P0)
	P1 := []byte{nasSecurity.AccessTypeNon3GPP}
	L1 := ueauth.KDFLen(P1)

	Kn3iwf, err := ueauth.GetKDFValue(ue.Kamf, ueauth.FC_FOR_KGNB_KN3IWF_DERIVATION, P0, L0, P1, L1)
	if err != nil {
		t.Fatalf("Get Kn3iwf error : %+v", err)
	}

	var idPayload ike_message.IKEPayloadContainer
	idPayload.BuildIdentificationInitiator(ike_message.ID_KEY_ID, []byte("UE"))
	idPayloadData, err := idPayload.Encode()
	if err != nil {
		t.Fatalf("Encode IKE payload failed : %+v", err)
	}
	if _, err = ikeSecurityAssociation.Prf_i.Write(idPayloadData[4:]); err != nil {
		t.Fatalf("Pseudorandom function write error: %+v", err)
	}
	ikeSecurityAssociation.ResponderSignedOctets = append(
		ikeSecurityAssociation.ResponderSignedOctets,
		ikeSecurityAssociation.Prf_i.Sum(nil)...)

	pseudorandomFunction := ikeSecurityAssociation.PrfInfo.Init(Kn3iwf)
	if _, err = pseudorandomFunction.Write([]byte("Key Pad for IKEv2")); err != nil {
		t.Fatalf("Pseudorandom function write error: %+v", err)
	}
	secret := pseudorandomFunction.Sum(nil)
	pseudorandomFunction = ikeSecurityAssociation.PrfInfo.Init(secret)
	pseudorandomFunction.Reset()
	if _, err = pseudorandomFunction.Write(ikeSecurityAssociation.ResponderSignedOctets); err != nil {
		t.Fatalf("Pseudorandom function write error: %+v", err)
	}

	ikePayload.BuildAuthentication(ike_message.SharedKeyMesageIntegrityCode, pseudorandomFunction.Sum(nil))

	// Configuration Request
	configurationRequest := ikePayload.BuildConfiguration(ike_message.CFG_REQUEST)
	configurationRequest.ConfigurationAttribute.BuildConfigurationAttribute(ike_message.INTERNAL_IP4_ADDRESS, nil)

	ikeMessage = ike_message.NewMessage(
		ikeSecurityAssociation.LocalSPI,
		ikeSecurityAssociation.RemoteSPI,
		ike_message.IKE_AUTH,
		false, true,
		ikeSecurityAssociation.InitiatorMessageID,
		ikePayload,
	)

	ikeMessageData, err = ike.EncodeEncrypt(ikeMessage, ikeSecurityAssociation.IKESAKey,
		ike_message.Role_Initiator)
	if err != nil {
		t.Fatalf("EncodeEncrypt IKE message failed: %+v", err)
	}
	_, err = udpConnection.WriteToUDP(ikeMessageData, n3iwfUDPAddr)
	if err != nil {
		t.Fatalf("Write IKE message failed: %+v", err)
	}

	// Receive N3IWF reply
	n, _, err = udpConnection.ReadFromUDP(buffer)
	if err != nil {
		t.Fatalf("Read IKE Message fail: %+v", err)
	}

	ikeMessage, err = ike.DecodeDecrypt(buffer[:n], nil,
		ikeSecurityAssociation.IKESAKey, ike_message.Role_Initiator)
	if err != nil {
		t.Fatalf("Decode IKE meesage: %v", err)
	}

	// AUTH, SAr2, TSi, Tsr, N(NAS_IP_ADDRESS), N(NAS_TCP_PORT)
	var responseSecurityAssociation *ike_message.SecurityAssociation
	var responseTrafficSelectorInitiator *ike_message.TrafficSelectorInitiator
	var responseTrafficSelectorResponder *ike_message.TrafficSelectorResponder
	var responseConfiguration *ike_message.Configuration
	n3iwfNASAddr := new(net.TCPAddr)

	for _, ikePayload := range ikeMessage.Payloads {
		switch ikePayload.Type() {
		case ike_message.TypeAUTH:
			t.Log("Get Authentication from N3IWF")
		case ike_message.TypeSA:
			responseSecurityAssociation = ikePayload.(*ike_message.SecurityAssociation)
			ikeSecurityAssociation.IKEAuthResponseSA = responseSecurityAssociation
		case ike_message.TypeTSi:
			responseTrafficSelectorInitiator = ikePayload.(*ike_message.TrafficSelectorInitiator)
		case ike_message.TypeTSr:
			responseTrafficSelectorResponder = ikePayload.(*ike_message.TrafficSelectorResponder)
		case ike_message.TypeN:
			notification := ikePayload.(*ike_message.Notification)
			if notification.NotifyMessageType == ike_message.Vendor3GPPNotifyTypeNAS_IP4_ADDRESS {
				n3iwfNASAddr.IP = net.IPv4(notification.NotificationData[0], notification.NotificationData[1], notification.NotificationData[2], notification.NotificationData[3])
			}
			if notification.NotifyMessageType == ike_message.Vendor3GPPNotifyTypeNAS_TCP_PORT {
				n3iwfNASAddr.Port = int(binary.BigEndian.Uint16(notification.NotificationData))
			}
		case ike_message.TypeCP:
			responseConfiguration = ikePayload.(*ike_message.Configuration)
			if responseConfiguration.ConfigurationType == ike_message.CFG_REPLY {
				for _, configAttr := range responseConfiguration.ConfigurationAttribute {
					if configAttr.Type == ike_message.INTERNAL_IP4_ADDRESS {
						ueInnerAddr.IP = configAttr.Value
					}
					if configAttr.Type == ike_message.INTERNAL_IP4_NETMASK {
						ueInnerAddr.Mask = configAttr.Value
					}
				}
			}
		}
	}

	OutboundSPI := binary.BigEndian.Uint32(ikeSecurityAssociation.IKEAuthResponseSA.Proposals[0].SPI)
	childSecurityAssociationContext, err := n3ue.N3IWFIkeUe.CompleteChildSA(
		0x01, OutboundSPI, ikeSecurityAssociation.IKEAuthResponseSA)
	if err != nil {
		t.Fatalf("Create child security association context failed: %+v", err)
	}
	err = parseIPAddressInformationToChildSecurityAssociation(childSecurityAssociationContext,
		responseTrafficSelectorInitiator.TrafficSelectors[0],
		responseTrafficSelectorResponder.TrafficSelectors[0])

	if err != nil {
		t.Fatalf("Parse IP address to child security association failed: %+v", err)
	}
	// Select TCP traffic
	childSecurityAssociationContext.SelectedIPProtocol = unix.IPPROTO_TCP

	if err := childSecurityAssociationContext.GenerateKeyForChildSA(ikeSecurityAssociation.IKESAKey,
		ikeSecurityAssociation.ConcatenatedNonce); err != nil {
		t.Fatalf("Generate key for child SA failed: %+v", err)
	}

	var linkIPSec netlink.Link

	// Setup interface for ipsec
	newXfrmiName := fmt.Sprintf("%s-default", n3ueInfo_XfrmiName)
	if linkIPSec, err = setupIPsecXfrmi(newXfrmiName, n3ueInfo_IPSecIfaceName, n3ueInfo_XfrmiId, ueInnerAddr); err != nil {
		t.Fatalf("Setup XFRM interface %s fail: %+v", newXfrmiName, err)
	}

	defer func() {
		if err := netlink.LinkDel(linkIPSec); err != nil {
			t.Fatalf("Delete XFRM interface %s fail: %+v", newXfrmiName, err)
		} else {
			t.Logf("Delete XFRM interface: %s", newXfrmiName)
		}
	}()

	// Aplly XFRM rules
	if err = applyXFRMRule(true, n3ueInfo_XfrmiId, childSecurityAssociationContext); err != nil {
		t.Fatalf("Applying XFRM rules failed: %+v", err)
	}

	defer func() {
		_ = netlink.XfrmPolicyFlush()
		_ = netlink.XfrmStateFlush(netlink.XFRM_PROTO_IPSEC_ANY)
	}()

	localTCPAddr := &net.TCPAddr{
		IP: ueInnerAddr.IP,
	}
	tcpConnWithN3IWF, err := net.DialTCP("tcp", localTCPAddr, n3iwfNASAddr)
	if err != nil {
		t.Fatal(err)
	}

	nasEnv := make([]byte, 65535)

	n, err = tcpConnWithN3IWF.Read(nasEnv)
	if err != nil {
		t.Fatal(err)
		return
	}

	nasEnv, n, err = DecapNasPduFromEnvelope(nasEnv[:n])
	if err != nil {
		t.Fatal(err)
	}
	nasMsg, err := NASDecode(ue, nas.SecurityHeaderTypeIntegrityProtectedAndCiphered, nasEnv[:n])
	if err != nil {
		t.Fatalf("NAS Decode Fail: %+v", err)
	}

	spew.Config.Indent = "\t"
	nasStr := spew.Sdump(nasMsg)
	t.Logf("Get NAS Security Mode Command Message:\n %+v", nasStr)

	// send NAS Registration Complete Msg
	pdu = nasTestpacket.GetRegistrationComplete(nil)
	pdu, err = EncodeNasPduInEnvelopeWithSecurity(ue, pdu, nas.SecurityHeaderTypeIntegrityProtectedAndCiphered, true, false)
	if err != nil {
		t.Fatal(err)
		return
	}
	_, err = tcpConnWithN3IWF.Write(pdu)
	if err != nil {
		t.Fatal(err)
		return
	}

	time.Sleep(500 * time.Millisecond)

	// UE request PDU session setup
	sNssai := models.Snssai{
		Sst: 1,
		Sd:  "fedcba",
	}

	var pduSessionId uint8 = 1

	pdu = nasTestpacket.GetUlNasTransport_PduSessionEstablishmentRequest(pduSessionId, nasMessage.ULNASTransportRequestTypeInitialRequest, "internet", &sNssai)
	pdu, err = EncodeNasPduInEnvelopeWithSecurity(ue, pdu, nas.SecurityHeaderTypeIntegrityProtectedAndCiphered, true, false)
	if err != nil {
		t.Fatal(err)
		return
	}
	_, err = tcpConnWithN3IWF.Write(pdu)
	if err != nil {
		t.Fatal(err)
		return
	}

	// Receive N3IWF reply
	n, _, err = udpConnection.ReadFromUDP(buffer)
	if err != nil {
		t.Fatalf("Read IKE Message fail: %+v", err)
	}
	ikeMessage.Payloads.Reset()
	ikeMessage, err = ike.DecodeDecrypt(buffer[:n], nil,
		ikeSecurityAssociation.IKESAKey, ike_message.Role_Initiator)
	if err != nil {
		t.Fatalf("Decode IKE meesage: %v", err)
	}

	var QoSInfo *PDUQoSInfo

	var upIPAddr net.IP
	for _, ikePayload := range ikeMessage.Payloads {
		switch ikePayload.Type() {
		case ike_message.TypeSA:
			responseSecurityAssociation = ikePayload.(*ike_message.SecurityAssociation)
			OutboundSPI = binary.BigEndian.Uint32(responseSecurityAssociation.Proposals[0].SPI)
		case ike_message.TypeTSi:
			responseTrafficSelectorInitiator = ikePayload.(*ike_message.TrafficSelectorInitiator)
		case ike_message.TypeTSr:
			responseTrafficSelectorResponder = ikePayload.(*ike_message.TrafficSelectorResponder)
		case ike_message.TypeN:
			notification := ikePayload.(*ike_message.Notification)
			if notification.NotifyMessageType == ike_message.Vendor3GPPNotifyType5G_QOS_INFO {
				t.Log("Received Qos Flow settings")
				if info, err := parse5GQoSInfoNotify(notification); err == nil {
					QoSInfo = info
					t.Logf("NotificationData:%+v", notification.NotificationData)
					if QoSInfo.isDSCPSpecified {
						t.Logf("DSCP is specified but test not support")
					}
				} else {
					t.Logf("%+v", err)
				}
			}
			if notification.NotifyMessageType == ike_message.Vendor3GPPNotifyTypeUP_IP4_ADDRESS {
				upIPAddr = notification.NotificationData[:4]
				t.Logf("UP IP Address: %+v\n", upIPAddr)
			}
		case ike_message.TypeNiNr:
			responseNonce := ikePayload.(*ike_message.Nonce)
			ikeSecurityAssociation.ConcatenatedNonce = responseNonce.NonceData
		}
	}

	// IKE CREATE_CHILD_SA response
	ikeMessage.Payloads.Reset()
	ikeSecurityAssociation.ResponderMessageID = ikeMessage.MessageID

	ikePayload.Reset()

	// SA
	inboundSPI, err = generateSPI(n3ue)
	if err != nil {
		t.Fatal(err)
	}
	responseSecurityAssociation.Proposals[0].SPI = inboundSPI
	ikePayload = append(ikePayload, responseSecurityAssociation)

	// TSi
	ikePayload = append(ikePayload, responseTrafficSelectorInitiator)

	// TSr
	ikePayload = append(ikePayload, responseTrafficSelectorResponder)

	// Nonce
	localNonceBigInt, err = ike_security.GenerateRandomNumber()
	if err != nil {
		t.Fatalf("Generate local nonce: %v", err)
	}
	localNonce = localNonceBigInt.Bytes()
	ikeSecurityAssociation.ConcatenatedNonce = append(ikeSecurityAssociation.ConcatenatedNonce, localNonce...)
	ikePayload.BuildNonce(localNonce)

	ikeMessage = ike_message.NewMessage(
		ikeSecurityAssociation.LocalSPI,
		ikeSecurityAssociation.RemoteSPI,
		ike_message.CREATE_CHILD_SA,
		true, true,
		ikeSecurityAssociation.InitiatorMessageID,
		ikePayload,
	)

	ikeMessageData, err = ike.EncodeEncrypt(ikeMessage, ikeSecurityAssociation.IKESAKey,
		ike_message.Role_Initiator)
	if err != nil {
		t.Fatalf("EncodeEncrypt IKE message failed: %+v", err)
	}
	_, err = udpConnection.WriteToUDP(ikeMessageData, n3iwfUDPAddr)
	if err != nil {
		t.Fatalf("Write IKE message failed: %+v", err)
	}

	n3ue.N3IWFIkeUe.CreateHalfChildSA(ikeSecurityAssociation.ResponderMessageID,
		binary.BigEndian.Uint32(inboundSPI), -1)
	childSecurityAssociationContextUserPlane, err := n3ue.N3IWFIkeUe.CompleteChildSA(
		ikeSecurityAssociation.ResponderMessageID, OutboundSPI, responseSecurityAssociation)

	if err != nil {
		t.Fatalf("Create child security association context failed: %+v", err)
	}
	err = parseIPAddressInformationToChildSecurityAssociation(childSecurityAssociationContextUserPlane, responseTrafficSelectorResponder.TrafficSelectors[0], responseTrafficSelectorInitiator.TrafficSelectors[0])
	if err != nil {
		t.Fatalf("Parse IP address to child security association failed: %+v", err)
	}
	// Select GRE traffic
	childSecurityAssociationContextUserPlane.SelectedIPProtocol = unix.IPPROTO_GRE

	if err := childSecurityAssociationContextUserPlane.GenerateKeyForChildSA(ikeSecurityAssociation.IKESAKey,
		ikeSecurityAssociation.ConcatenatedNonce); err != nil {
		t.Fatalf("Generate key for child SA failed: %+v", err)
	}

	// Aplly XFRM rules
	if err = applyXFRMRule(false, n3ueInfo_XfrmiId, childSecurityAssociationContextUserPlane); err != nil {
		t.Fatalf("Applying XFRM rules failed: %+v", err)
	}

	// TODO
	// We don't check any of message in UeConfigUpdate Message
	if n, err := tcpConnWithN3IWF.Read(buffer); err != nil {
		t.Fatalf("No UeConfigUpdate Message: %+v", err)
		_, err := ngap.Decoder(buffer[2:n])
		if err != nil {
			t.Fatalf("UeConfigUpdate Decode Error: %+v", err)
		}
	}

	var pduAddress net.IP

	// Read NAS from N3IWF
	if n, err := tcpConnWithN3IWF.Read(buffer); err != nil {
		t.Fatalf("Read NAS Message Fail:%+v", err)
	} else {
		nasMsg, err := DecodePDUSessionEstablishmentAccept(ue, n, buffer)
		if err != nil {
			t.Fatalf("DecodePDUSessionEstablishmentAccept Fail: %+v", err)
		}

		spew.Config.Indent = "\t"
		nasStr := spew.Sdump(nasMsg)
		t.Log("Dump DecodePDUSessionEstablishmentAccept:\n", nasStr)
		pduAddress, err = GetPDUAddress(nasMsg.GsmMessage.PDUSessionEstablishmentAccept)
		if err != nil {
			t.Fatalf("GetPDUAddress Fail: %+v", err)
		}

		t.Logf("PDU Address: %s", pduAddress.String())
	}

	var linkGRE netlink.Link

	newGREName := fmt.Sprintf("%s-id-%d", n3ueInfo_GreIfaceName, n3ueInfo_XfrmiId)

	if linkGRE, err = setupGreTunnel(newGREName, newXfrmiName, ueInnerAddr.IP, upIPAddr, pduAddress, QoSInfo, t); err != nil {
		t.Fatalf("Setup GRE tunnel %s Fail %+v", newGREName, err)
	}

	defer func() {
		_ = netlink.LinkDel(linkGRE)
		t.Logf("Delete interface: %s", linkGRE.Attrs().Name)
	}()

	// Add route
	upRoute := &netlink.Route{
		LinkIndex: linkGRE.Attrs().Index,
		Dst: &net.IPNet{
			IP:   net.IPv4zero,
			Mask: net.IPv4Mask(0, 0, 0, 0),
		},
	}
	if err := netlink.RouteAdd(upRoute); err != nil {
		t.Fatal(err)
	}

	for i := 1; i <= 3; i++ {
		var (
			ifaces []netlink.Link
			err    error
		)
		t.Logf("%d times PDU Session Est Request Start", i+1)
		if ifaces, err = sendPduSessionEstablishmentRequest(pduSessionId+uint8(i), ue, n3ue, ikeSecurityAssociation, udpConnection, tcpConnWithN3IWF, t); err != nil {
			t.Fatalf("Session Est Request Fail: %+v", err)
		} else {
			t.Logf("Create %d interfaces", len(ifaces))
		}

		defer func() {
			for _, iface := range ifaces {
				if err := netlink.LinkDel(iface); err != nil {
					t.Fatalf("Delete interface %s fail: %+v", iface.Attrs().Name, err)
				} else {
					t.Logf("Delete interface: %s", iface.Attrs().Name)
				}
			}
		}()
	}

	// Ping remote
	pinger, err := ping.NewPinger("10.60.0.101")
	if err != nil {
		t.Fatal(err)
		return
	}

	// Run with root
	pinger.SetPrivileged(true)

	pinger.OnRecv = func(pkt *ping.Packet) {
		t.Logf("%d bytes from %s: icmp_seq=%d time=%v\n",
			pkt.Nbytes, pkt.IPAddr, pkt.Seq, pkt.Rtt)
	}
	pinger.OnFinish = func(stats *ping.Statistics) {
		t.Logf("\n--- %s ping statistics ---\n", stats.Addr)
		t.Logf("%d packets transmitted, %d packets received, %v%% packet loss\n",
			stats.PacketsSent, stats.PacketsRecv, stats.PacketLoss)
		t.Logf("round-trip min/avg/max/stddev = %v/%v/%v/%v\n",
			stats.MinRtt, stats.AvgRtt, stats.MaxRtt, stats.StdDevRtt)
	}

	pinger.Count = 5
	pinger.Timeout = 10 * time.Second
	pinger.Source = "10.60.0.1"

	time.Sleep(3 * time.Second)

	pinger.Run()

	time.Sleep(1 * time.Second)

	stats := pinger.Statistics()
	if stats.PacketsSent != stats.PacketsRecv {
		t.Fatal("Ping Failed")
		return
	}
}

func setUESecurityCapability(ue *RanUeContext) (UESecurityCapability *nasType.UESecurityCapability) {
	UESecurityCapability = &nasType.UESecurityCapability{
		Iei:    nasMessage.RegistrationRequestUESecurityCapabilityType,
		Len:    8,
		Buffer: []uint8{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
	}
	switch ue.CipheringAlg {
	case nasSecurity.AlgCiphering128NEA0:
		UESecurityCapability.SetEA0_5G(1)
	case nasSecurity.AlgCiphering128NEA1:
		UESecurityCapability.SetEA1_128_5G(1)
	case nasSecurity.AlgCiphering128NEA2:
		UESecurityCapability.SetEA2_128_5G(1)
	case nasSecurity.AlgCiphering128NEA3:
		UESecurityCapability.SetEA3_128_5G(1)
	}

	switch ue.IntegrityAlg {
	case nasSecurity.AlgIntegrity128NIA0:
		UESecurityCapability.SetIA0_5G(1)
	case nasSecurity.AlgIntegrity128NIA1:
		UESecurityCapability.SetIA1_128_5G(1)
	case nasSecurity.AlgIntegrity128NIA2:
		UESecurityCapability.SetIA2_128_5G(1)
	case nasSecurity.AlgIntegrity128NIA3:
		UESecurityCapability.SetIA3_128_5G(1)
	}
	return
}

func (ikeUe *N3IWFIkeUe) CreateHalfChildSA(msgID, inboundSPI uint32, pduSessionID int64) {
	childSA := new(ChildSecurityAssociation)
	childSA.InboundSPI = inboundSPI
	// Map Exchange Message ID and Child SA data until get paired response
	ikeUe.TemporaryExchangeMsgIDChildSAMapping[msgID] = childSA
}

func (ikeUe *N3IWFIkeUe) CompleteChildSA(msgID uint32, outboundSPI uint32,
	chosenSecurityAssociation *message.SecurityAssociation,
) (*ChildSecurityAssociation, error) {
	childSA, ok := ikeUe.TemporaryExchangeMsgIDChildSAMapping[msgID]

	if !ok {
		return nil, errors.Errorf("CompleteChildSA(): There's not a half child SA created by the exchange with message ID %d.", msgID)
	}

	// Remove mapping of exchange msg ID and child SA
	delete(ikeUe.TemporaryExchangeMsgIDChildSAMapping, msgID)

	if chosenSecurityAssociation == nil {
		return nil, errors.Errorf("CompleteChildSA(): chosenSecurityAssociation is nil")
	}

	if len(chosenSecurityAssociation.Proposals) == 0 {
		return nil, errors.Errorf("CompleteChildSA(): No proposal")
	}

	childSA.OutboundSPI = outboundSPI

	var err error
	childSA.ChildSAKey, err = ike_security.NewChildSAKeyByProposal(chosenSecurityAssociation.Proposals[0])
	if err != nil {
		return nil, errors.Wrapf(err, "CompleteChildSA")
	}

	// Record to UE context with inbound SPI as key
	ikeUe.N3IWFChildSecurityAssociation[childSA.InboundSPI] = childSA

	return childSA, nil
}
