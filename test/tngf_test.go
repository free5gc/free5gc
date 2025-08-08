package test

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/rand"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"hash"
	"math/big"
	"net"
	"strconv"
	"testing"
	"time"

	"test/nasTestpacket"

	"github.com/davecgh/go-spew/spew"
	"github.com/go-ping/ping"
	"github.com/stretchr/testify/assert"
	"github.com/vishvananda/netlink"
	"golang.org/x/sys/unix"

	ike_security "github.com/free5gc/ike/security"
	"github.com/free5gc/nas"
	"github.com/free5gc/nas/nasMessage"
	"github.com/free5gc/nas/nasType"
	nasSecurity "github.com/free5gc/nas/security"
	"github.com/free5gc/openapi/models"
	"github.com/free5gc/tngf/pkg/context"
	"github.com/free5gc/tngf/pkg/ike/handler"
	"github.com/free5gc/tngf/pkg/ike/message"
	"github.com/free5gc/tngf/pkg/ike/xfrm"
	radiusHandler "github.com/free5gc/tngf/pkg/radius/handler"
	radiusMessage "github.com/free5gc/tngf/pkg/radius/message"
	"github.com/free5gc/util/ueauth"
)

var (
	tngfInfo_IPSecIfaceAddr        = "192.168.127.1"
	tngfueInfo_IPSecIfaceAddr      = "192.168.127.2"
	tngfueInfo_SmPolicy_SNSSAI_SST = "1"
	tngfueInfo_SmPolicy_SNSSAI_SD  = "fedcba"
	tngfueInfo_IPSecIfaceName      = "veth3"
	tngfueInfo_XfrmiName           = "ipsec"
	tngfueInfo_XfrmiId             = uint32(1)
	tngfueInfo_GreIfaceName        = "gretun"
	tngfueInnerAddr                = new(net.IPNet)
)

func concatenateNonceAndSPI(nonce []byte, spi_initiator uint64, spi_responder uint64) []byte {
	var newSlice []byte
	spi := make([]byte, 8)

	newSlice = append(newSlice, nonce...)
	binary.BigEndian.PutUint64(spi, spi_initiator)
	newSlice = append(newSlice, spi...)
	binary.BigEndian.PutUint64(spi, spi_responder)
	newSlice = append(newSlice, spi...)

	return newSlice
}

func tngfGenerateSPI(tngfue *context.TNGFUe) []byte {
	var spi uint32
	spiByte := make([]byte, 4)
	for {
		randomUint64 := handler.GenerateRandomNumber().Uint64()
		if _, ok := tngfue.TNGFChildSecurityAssociation[uint32(randomUint64)]; !ok {
			spi = uint32(randomUint64)
			binary.BigEndian.PutUint32(spiByte, spi)
			break
		}
	}
	return spiByte
}

func tngfSetupIPsecXfrmi(xfrmIfaceName, parentIfaceName string, xfrmIfaceId uint32, xfrmIfaceAddr *net.IPNet) (netlink.Link, error) {
	var (
		xfrmi, parent netlink.Link
		err           error
	)

	if parent, err = netlink.LinkByName(parentIfaceName); err != nil {
		return nil, err
	}

	if oldLink, err := netlink.LinkByName(xfrmIfaceName); err == nil {
		fmt.Println("Deleting old XFRM interface...")
		_ = netlink.LinkDel(oldLink)
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
	fmt.Printf("Expected XFRM Interface Addr: %s\n", xfrmIfaceAddr.String())
	return xfrmi, nil
}

func setupRadiusSocket() (*net.UDPConn, error) {
	bindAddr := tngfueInfo_IPSecIfaceAddr + ":48744"
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

func tngfGenerateKeyForIKESA(ikeSecurityAssociation *context.IKESecurityAssociation) error {
	// Transforms
	transformPseudorandomFunction := ikeSecurityAssociation.PseudorandomFunction

	// Get key length of SK_d, SK_ai, SK_ar, SK_ei, SK_er, SK_pi, SK_pr
	var length_SK_d, length_SK_ai, length_SK_ar, length_SK_ei, length_SK_er, length_SK_pi, length_SK_pr, totalKeyLength int
	var ok bool

	length_SK_d = 20
	length_SK_ai = 20
	length_SK_ar = length_SK_ai
	length_SK_ei = 0
	length_SK_er = length_SK_ei
	length_SK_pi, length_SK_pr = length_SK_d, length_SK_d
	totalKeyLength = length_SK_d + length_SK_ai + length_SK_ar + length_SK_ei + length_SK_er + length_SK_pi + length_SK_pr

	// Generate IKE SA key as defined in RFC7296 Section 1.3 and Section 1.4
	var pseudorandomFunction hash.Hash

	if pseudorandomFunction, ok = handler.NewPseudorandomFunction(ikeSecurityAssociation.ConcatenatedNonce, transformPseudorandomFunction.TransformID); !ok {
		return errors.New("New pseudorandom function failed")
	}

	if _, err := pseudorandomFunction.Write(ikeSecurityAssociation.DiffieHellmanSharedKey); err != nil {
		return errors.New("Pseudorandom function write failed")
	}

	SKEYSEED := pseudorandomFunction.Sum(nil)

	seed := concatenateNonceAndSPI(ikeSecurityAssociation.ConcatenatedNonce, ikeSecurityAssociation.LocalSPI, ikeSecurityAssociation.RemoteSPI)

	var keyStream, generatedKeyBlock []byte
	var index byte
	for index = 1; len(keyStream) < totalKeyLength; index++ {
		if pseudorandomFunction, ok = handler.NewPseudorandomFunction(SKEYSEED, transformPseudorandomFunction.TransformID); !ok {
			return errors.New("New pseudorandom function failed")
		}
		if _, err := pseudorandomFunction.Write(append(append(generatedKeyBlock, seed...), index)); err != nil {
			return errors.New("Pseudorandom function write failed")
		}
		generatedKeyBlock = pseudorandomFunction.Sum(nil)
		keyStream = append(keyStream, generatedKeyBlock...)
	}

	// Assign keys into context
	ikeSecurityAssociation.SK_d = keyStream[:length_SK_d]
	keyStream = keyStream[length_SK_d:]
	ikeSecurityAssociation.SK_ai = keyStream[:length_SK_ai]
	keyStream = keyStream[length_SK_ai:]
	ikeSecurityAssociation.SK_ar = keyStream[:length_SK_ar]
	keyStream = keyStream[length_SK_ar:]
	ikeSecurityAssociation.SK_ei = keyStream[:length_SK_ei]
	keyStream = keyStream[length_SK_ei:]
	ikeSecurityAssociation.SK_er = keyStream[:length_SK_er]
	keyStream = keyStream[length_SK_er:]
	ikeSecurityAssociation.SK_pi = keyStream[:length_SK_pi]
	keyStream = keyStream[length_SK_pi:]
	ikeSecurityAssociation.SK_pr = keyStream[:length_SK_pr]
	keyStream = keyStream[length_SK_pr:]

	return nil
}

func tngfDecryptProcedure(ikeSecurityAssociation *context.IKESecurityAssociation, ikeMessage *message.IKEMessage, encryptedPayload *message.Encrypted) (message.IKEPayloadContainer, error) {
	// Load needed information
	transformIntegrityAlgorithm := ikeSecurityAssociation.IntegrityAlgorithm
	transformEncryptionAlgorithm := ikeSecurityAssociation.EncryptionAlgorithm
	checksumLength := 12 // HMAC_SHA1_96

	// Checksum
	checksum := encryptedPayload.EncryptedData[len(encryptedPayload.EncryptedData)-checksumLength:]

	ikeMessageData, err := ikeMessage.Encode()
	if err != nil {
		return nil, errors.New("Encoding IKE message failed")
	}

	ok, err := handler.VerifyIKEChecksum(ikeSecurityAssociation.SK_ar, ikeMessageData[:len(ikeMessageData)-checksumLength], checksum, transformIntegrityAlgorithm.TransformID)
	if err != nil {
		return nil, errors.New("Error verify checksum")
	}
	if !ok {
		return nil, errors.New("Checksum failed, drop.")
	}

	// Decrypt
	encryptedData := encryptedPayload.EncryptedData[:len(encryptedPayload.EncryptedData)-checksumLength]
	plainText, err := handler.DecryptMessage(ikeSecurityAssociation.SK_er, encryptedData, transformEncryptionAlgorithm.TransformID)
	if err != nil {
		return nil, errors.New("Error decrypting message")
	}

	var decryptedIKEPayload message.IKEPayloadContainer
	err = decryptedIKEPayload.Decode(encryptedPayload.NextPayload, plainText)
	if err != nil {
		return nil, errors.New("Decoding decrypted payload failed")
	}

	return decryptedIKEPayload, nil

}

func tngfEncryptProcedure(ikeSecurityAssociation *context.IKESecurityAssociation, ikePayload message.IKEPayloadContainer, responseIKEMessage *message.IKEMessage) error {
	// Load needed information
	transformIntegrityAlgorithm := ikeSecurityAssociation.IntegrityAlgorithm
	transformEncryptionAlgorithm := ikeSecurityAssociation.EncryptionAlgorithm
	checksumLength := 12 // HMAC_SHA1_96

	// Encrypting
	notificationPayloadData, err := ikePayload.Encode()
	if err != nil {
		return errors.New("Encoding IKE payload failed.")
	}

	encryptedData, err := handler.EncryptMessage(ikeSecurityAssociation.SK_ei, notificationPayloadData, transformEncryptionAlgorithm.TransformID)
	if err != nil {
		return errors.New("Error encrypting message")
	}

	encryptedData = append(encryptedData, make([]byte, checksumLength)...)
	sk := responseIKEMessage.Payloads.BuildEncrypted(ikePayload[0].Type(), encryptedData)

	// Calculate checksum
	responseIKEMessageData, err := responseIKEMessage.Encode()
	if err != nil {
		return errors.New("Encoding IKE message error")
	}
	checksumOfMessage, err := handler.CalculateChecksum(ikeSecurityAssociation.SK_ai, responseIKEMessageData[:len(responseIKEMessageData)-checksumLength], transformIntegrityAlgorithm.TransformID)
	if err != nil {
		return errors.New("Error calculating checksum")
	}
	checksumField := sk.EncryptedData[len(sk.EncryptedData)-checksumLength:]
	copy(checksumField, checksumOfMessage)

	return nil

}

// [TS 24502] 9.3.2.2.2 EAP-Response/5G-NAS message
// Define EAP-Response/5G-NAS message and AN-Parameters Format.

// [TS 24501] 8.2.6.1.1  REGISTRATION REQUEST message content
// For dealing with EAP-5G start, return EAP-5G response including
// "AN-Parameters and NASPDU of Registration Request"

func tngfBuildEAP5GANParameters(mobileIdentity5GS nasType.MobileIdentity5GS) []byte {
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
	anParameter[0] = radiusMessage.ANParametersTypeGUAMI
	anParameter[1] = byte(len(guami))
	anParameter = append(anParameter, guami...)

	anParameters = append(anParameters, anParameter...)

	// Build Establishment Cause
	anParameter = make([]byte, 2)
	establishmentCause := make([]byte, 1)
	establishmentCause[0] = radiusMessage.EstablishmentCauseMO_Signaling
	anParameter[0] = radiusMessage.ANParametersTypeEstablishmentCause
	anParameter[1] = byte(len(establishmentCause))
	anParameter = append(anParameter, establishmentCause...)

	anParameters = append(anParameters, anParameter...)

	// Build PLMN ID
	anParameter = make([]byte, 2)
	plmnID := make([]byte, 3)
	plmnID[0] = 0x02
	plmnID[1] = 0xf8
	plmnID[2] = 0x39
	anParameter[0] = radiusMessage.ANParametersTypeSelectedPLMNID
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
	anParameter[0] = radiusMessage.ANParametersTypeRequestedNSSAI
	anParameter[1] = byte(len(nssai))
	anParameter = append(anParameter, nssai...)

	anParameters = append(anParameters, anParameter...)

	// Build UE ID
	anParameter = make([]byte, 3)
	anParameter[0] = radiusMessage.ANParametersTypeUEIdentity
	anParameter[1] = byte(16)
	anParameter[2] = mobileIdentity5GS.GetIei()
	anParameterLength := make([]byte, 2)
	binary.BigEndian.PutUint16(anParameterLength, mobileIdentity5GS.GetLen())
	anParameter = append(anParameter, anParameterLength...)
	anParameter = append(anParameter, mobileIdentity5GS.Buffer...)

	anParameters = append(anParameters, anParameter...)

	return anParameters
}

func tngfParseIPAddressInformationToChildSecurityAssociation(
	childSecurityAssociation *context.ChildSecurityAssociation,
	trafficSelectorLocal *message.IndividualTrafficSelector,
	trafficSelectorRemote *message.IndividualTrafficSelector) error {

	if childSecurityAssociation == nil {
		return errors.New("childSecurityAssociation is nil")
	}

	childSecurityAssociation.PeerPublicIPAddr = net.ParseIP(tngfInfo_IPSecIfaceAddr)
	childSecurityAssociation.LocalPublicIPAddr = net.ParseIP(tngfueInfo_IPSecIfaceAddr)

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

func tngfParse5GQoSInfoNotify(n *message.Notification) (info *PDUQoSInfo, err error) {
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

	info.isDefault = (data[offset] & message.NotifyType5G_QOS_INFOBitDCSICheck) > 0
	info.isDSCPSpecified = (data[offset] & message.NotifyType5G_QOS_INFOBitDSCPICheck) > 0

	return
}

func tngfApplyXFRMRule(ue_is_initiator bool, ifId uint32, childSecurityAssociation *context.ChildSecurityAssociation) error {
	// Build XFRM information data structure for incoming traffic.

	// Mark
	// mark := &netlink.XfrmMark{
	// 	Value: ifMark, // tngfueInfo.XfrmMark,
	// }

	// Direction: TNGF -> UE
	// State
	var xfrmEncryptionAlgorithm, xfrmIntegrityAlgorithm *netlink.XfrmStateAlgo
	if ue_is_initiator {
		xfrmEncryptionAlgorithm = &netlink.XfrmStateAlgo{
			Name: "ecb(cipher_null)",
			Key:  nil,
		}
		if childSecurityAssociation.IntegrityAlgorithm != 0 {
			xfrmIntegrityAlgorithm = &netlink.XfrmStateAlgo{
				Name: xfrm.XFRMIntegrityAlgorithmType(childSecurityAssociation.IntegrityAlgorithm).String(),
				Key:  childSecurityAssociation.ResponderToInitiatorIntegrityKey,
			}
		}
	} else {
		xfrmEncryptionAlgorithm = &netlink.XfrmStateAlgo{
			Name: "ecb(cipher_null)",
			Key:  nil,
		}
		if childSecurityAssociation.IntegrityAlgorithm != 0 {
			xfrmIntegrityAlgorithm = &netlink.XfrmStateAlgo{
				Name: xfrm.XFRMIntegrityAlgorithmType(childSecurityAssociation.IntegrityAlgorithm).String(),
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
	xfrmState.ESN = childSecurityAssociation.ESN

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

	// Direction: UE -> TNGF
	// State
	if ue_is_initiator {
		xfrmEncryptionAlgorithm.Key = nil
		if childSecurityAssociation.IntegrityAlgorithm != 0 {
			xfrmIntegrityAlgorithm.Key = childSecurityAssociation.InitiatorToResponderIntegrityKey
		}
	} else {
		xfrmEncryptionAlgorithm.Key = nil
		if childSecurityAssociation.IntegrityAlgorithm != 0 {
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

func tngfSendPduSessionEstablishmentRequest(
	pduSessionId uint8,
	ue *RanUeContext,
	tngfInfo *context.TNGFUe,
	ikeSA *context.IKESecurityAssociation,
	ikeConn *net.UDPConn,
	nasConn *net.TCPConn,
	t *testing.T) ([]netlink.Link, error) {

	var ifaces []netlink.Link

	// Build S-NSSA
	sst, err := strconv.ParseInt(tngfueInfo_SmPolicy_SNSSAI_SST, 16, 0)

	if err != nil {
		return ifaces, fmt.Errorf("Parse SST Fail:%+v", err)
	}

	sNssai := models.Snssai{
		Sst: int32(sst),
		Sd:  tngfueInfo_SmPolicy_SNSSAI_SD,
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

	t.Logf("Waiting for TNGF reply from IKE")

	// Receive TNGF reply
	n, _, err := ikeConn.ReadFromUDP(buffer)
	if err != nil {
		return ifaces, fmt.Errorf("Read IKE Message Fail:%+v", err)
	}

	ikeMessage := new(message.IKEMessage)
	ikeMessage.Payloads.Reset()
	err = ikeMessage.Decode(buffer[:n])
	if err != nil {
		return ifaces, fmt.Errorf("Decode IKE Message Fail:%+v", err)
	}
	t.Logf("IKE message exchange type: %d", ikeMessage.ExchangeType)
	t.Logf("IKE message ID: %d", ikeMessage.MessageID)

	encryptedPayload, ok := ikeMessage.Payloads[0].(*message.Encrypted)
	if !ok {
		return ifaces, errors.New("Received pakcet is not an encrypted payload")
	}
	decryptedIKEPayload, err := tngfDecryptProcedure(ikeSA, ikeMessage, encryptedPayload)
	if err != nil {
		return ifaces, fmt.Errorf("Decrypt IKE Message Fail:%+v", err)
	}

	var qoSInfo *PDUQoSInfo

	var responseSecurityAssociation *message.SecurityAssociation
	var responseTrafficSelectorInitiator *message.TrafficSelectorInitiator
	var responseTrafficSelectorResponder *message.TrafficSelectorResponder
	var outboundSPI uint32
	var upIPAddr net.IP
	for _, ikePayload := range decryptedIKEPayload {
		switch ikePayload.Type() {
		case message.TypeSA:
			responseSecurityAssociation = ikePayload.(*message.SecurityAssociation)
			outboundSPI = binary.BigEndian.Uint32(responseSecurityAssociation.Proposals[0].SPI)
		case message.TypeTSi:
			responseTrafficSelectorInitiator = ikePayload.(*message.TrafficSelectorInitiator)
		case message.TypeTSr:
			responseTrafficSelectorResponder = ikePayload.(*message.TrafficSelectorResponder)
		case message.TypeN:
			notification := ikePayload.(*message.Notification)
			if notification.NotifyMessageType == message.Vendor3GPPNotifyType5G_QOS_INFO {
				t.Logf("Received Qos Flow settings")
				if info, err := tngfParse5GQoSInfoNotify(notification); err == nil {
					qoSInfo = info
					t.Logf("NotificationData:%+v", notification.NotificationData)
					if qoSInfo.isDSCPSpecified {
						t.Logf("DSCP is specified but test not support")
					}
				} else {
					t.Logf("%+v", err)
				}
			}
			if notification.NotifyMessageType == message.Vendor3GPPNotifyTypeUP_IP4_ADDRESS {
				upIPAddr = notification.NotificationData[:4]
				t.Logf("UP IP Address: %+v\n", upIPAddr)
			}
		case message.TypeNiNr:
			responseNonce := ikePayload.(*message.Nonce)
			ikeSA.ConcatenatedNonce = responseNonce.NonceData
		}
	}

	// IKE CREATE_CHILD_SA response
	ikeMessage.Payloads.Reset()
	tngfInfo.TNGFIKESecurityAssociation.ResponderMessageID = ikeMessage.MessageID
	ikeMessage.BuildIKEHeader(ikeMessage.InitiatorSPI, ikeMessage.ResponderSPI,
		message.CREATE_CHILD_SA, message.ResponseBitCheck|message.InitiatorBitCheck,
		tngfInfo.TNGFIKESecurityAssociation.ResponderMessageID)

	var ikePayload message.IKEPayloadContainer
	ikePayload.Reset()

	// SA
	inboundSPI := tngfGenerateSPI(tngfInfo)
	responseSecurityAssociation.Proposals[0].SPI = inboundSPI
	ikePayload = append(ikePayload, responseSecurityAssociation)

	// TSi
	ikePayload = append(ikePayload, responseTrafficSelectorInitiator)

	// TSr
	ikePayload = append(ikePayload, responseTrafficSelectorResponder)

	// Nonce
	localNonce := handler.GenerateRandomNumber().Bytes()
	ikeSA.ConcatenatedNonce = append(ikeSA.ConcatenatedNonce, localNonce...)
	ikePayload.BuildNonce(localNonce)

	if err := tngfEncryptProcedure(ikeSA, ikePayload, ikeMessage); err != nil {
		t.Errorf("Encrypt IKE message failed: %+v", err)
		return ifaces, err
	}

	// Send to TNGF
	ikeMessageData, err := ikeMessage.Encode()
	if err != nil {
		return ifaces, fmt.Errorf("Encode IKE Message Fail:%+v", err)
	}

	tngfUDPAddr, err := net.ResolveUDPAddr("udp", tngfInfo_IPSecIfaceAddr+":500")

	if err != nil {
		return ifaces, fmt.Errorf("Resolve TNGF IPSec IP Addr Fail:%+v", err)
	}

	_, err = ikeConn.WriteToUDP(ikeMessageData, tngfUDPAddr)
	if err != nil {
		t.Errorf("Write IKE maessage fail: %+v", err)
		return ifaces, err
	}

	tngfInfo.CreateHalfChildSA(tngfInfo.TNGFIKESecurityAssociation.ResponderMessageID, binary.BigEndian.Uint32(inboundSPI), int64(pduSessionId))
	childSecurityAssociationContextUserPlane, err := tngfInfo.CompleteChildSA(
		tngfInfo.TNGFIKESecurityAssociation.ResponderMessageID, outboundSPI, responseSecurityAssociation)
	if err != nil {
		return ifaces, fmt.Errorf("Create child security association context failed: %+v", err)
	}

	err = tngfParseIPAddressInformationToChildSecurityAssociation(
		childSecurityAssociationContextUserPlane,
		responseTrafficSelectorResponder.TrafficSelectors[0],
		responseTrafficSelectorInitiator.TrafficSelectors[0])

	if err != nil {
		return ifaces, fmt.Errorf("Parse IP address to child security association failed: %+v", err)
	}
	// Select GRE traffic
	childSecurityAssociationContextUserPlane.SelectedIPProtocol = unix.IPPROTO_GRE

	if err := handler.GenerateKeyForChildSA(ikeSA, childSecurityAssociationContextUserPlane); err != nil {
		return ifaces, fmt.Errorf("Generate key for child SA failed: %+v", err)
	}

	// ====== Inbound ======
	t.Logf("====== IPSec/Child SA for 3GPP UP Inbound =====")
	t.Logf("[UE:%+v] <- [TNGF:%+v]",
		childSecurityAssociationContextUserPlane.LocalPublicIPAddr, childSecurityAssociationContextUserPlane.PeerPublicIPAddr)
	t.Logf("IPSec SPI: 0x%016x", childSecurityAssociationContextUserPlane.InboundSPI)
	t.Logf("IPSec Encryption Algorithm: %d", childSecurityAssociationContextUserPlane.EncryptionAlgorithm)
	t.Logf("IPSec Encryption Key: 0x%x", childSecurityAssociationContextUserPlane.InitiatorToResponderEncryptionKey)
	t.Logf("IPSec Integrity  Algorithm: %d", childSecurityAssociationContextUserPlane.IntegrityAlgorithm)
	t.Logf("IPSec Integrity  Key: 0x%x", childSecurityAssociationContextUserPlane.InitiatorToResponderIntegrityKey)
	// ====== Outbound ======
	t.Logf("====== IPSec/Child SA for 3GPP UP Outbound =====")
	t.Logf("[UE:%+v] -> [TNGF:%+v]",
		childSecurityAssociationContextUserPlane.LocalPublicIPAddr, childSecurityAssociationContextUserPlane.PeerPublicIPAddr)
	t.Logf("IPSec SPI: 0x%016x", childSecurityAssociationContextUserPlane.OutboundSPI)
	t.Logf("IPSec Encryption Algorithm: %d", childSecurityAssociationContextUserPlane.EncryptionAlgorithm)
	t.Logf("IPSec Encryption Key: 0x%x", childSecurityAssociationContextUserPlane.ResponderToInitiatorEncryptionKey)
	t.Logf("IPSec Integrity  Algorithm: %d", childSecurityAssociationContextUserPlane.IntegrityAlgorithm)
	t.Logf("IPSec Integrity  Key: 0x%x", childSecurityAssociationContextUserPlane.ResponderToInitiatorIntegrityKey)
	t.Logf("State function: encr: %d, auth: %d", childSecurityAssociationContextUserPlane.EncryptionAlgorithm, childSecurityAssociationContextUserPlane.IntegrityAlgorithm)

	// Aplly XFRM rules
	tngfueInfo_XfrmiId++
	err = tngfApplyXFRMRule(false, tngfueInfo_XfrmiId, childSecurityAssociationContextUserPlane)

	if err != nil {
		t.Errorf("Applying XFRM rules failed: %+v", err)
		return ifaces, err
	}

	var linkIPSec netlink.Link

	// Setup interface for ipsec
	newXfrmiName := fmt.Sprintf("%s-%d", tngfueInfo_XfrmiName, tngfueInfo_XfrmiId)
	if linkIPSec, err = setupIPsecXfrmi(newXfrmiName, tngfueInfo_IPSecIfaceName, tngfueInfo_XfrmiId, tngfueInnerAddr); err != nil {
		return ifaces, fmt.Errorf("Setup XFRMi interface %s fail: %+v", newXfrmiName, err)
	}

	ifaces = append(ifaces, linkIPSec)

	t.Logf("Setup XFRM interface %s successfully", newXfrmiName)

	var pduAddr net.IP

	// Read NAS from TNGF
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

	newGREName := fmt.Sprintf("%s-id-%d", tngfueInfo_GreIfaceName, tngfueInfo_XfrmiId)

	if linkGRE, err = setupGreTunnel(newGREName, newXfrmiName, tngfueInnerAddr.IP, upIPAddr, pduAddr, qoSInfo, t); err != nil {
		return ifaces, fmt.Errorf("Setup GRE tunnel %s Fail %+v", newGREName, err)
	}

	ifaces = append(ifaces, linkGRE)

	return ifaces, nil
}

// create EAP Identity and append to Radius payload
func BuildEAPIdentity(container *radiusMessage.RadiusPayloadContainer, identifier uint8, identityData []byte) {
	eap := new(radiusMessage.EAP)
	eap.Code = radiusMessage.EAPCodeResponse
	eap.Identifier = identifier
	eapIdentity := new(radiusMessage.EAPIdentity)
	eapIdentity.IdentityData = identityData
	eap.EAPTypeData = append(eap.EAPTypeData, eapIdentity)
	eapPayload, err := eap.Marshal()
	if err != nil {
		return
	}
	payload := new(radiusMessage.RadiusPayload)
	payload.Type = radiusMessage.TypeEAPMessage
	payload.Val = eapPayload

	*container = append(*container, *payload)
}

func BuildEAP5GNAS(container *radiusMessage.RadiusPayloadContainer, identifier uint8, vendorData []byte) {
	eap := new(radiusMessage.EAP)
	eap.Code = radiusMessage.EAPCodeResponse
	eap.Identifier = identifier
	eap.EAPTypeData.BuildEAPExpanded(radiusMessage.VendorID3GPP, radiusMessage.VendorTypeEAP5G, vendorData)
	eapPayload, err := eap.Marshal()
	if err != nil {
		return
	}

	payload := new(radiusMessage.RadiusPayload)
	payload.Type = radiusMessage.TypeEAPMessage
	payload.Val = eapPayload

	*container = append(*container, *payload)
}

func BuildEAP5GNotification(container *radiusMessage.RadiusPayloadContainer, identifier uint8) {
	eap := new(radiusMessage.EAP)
	eap.Code = radiusMessage.EAPCodeResponse
	eap.Identifier = identifier
	vendorData := make([]byte, 2)
	vendorData[0] = radiusMessage.EAP5GType5GNotification
	eap.EAPTypeData.BuildEAPExpanded(radiusMessage.VendorID3GPP, radiusMessage.VendorTypeEAP5G, vendorData)
	eapPayload, err := eap.Marshal()
	if err != nil {
		return
	}

	payload := new(radiusMessage.RadiusPayload)
	payload.Type = radiusMessage.TypeEAPMessage
	payload.Val = eapPayload

	*container = append(*container, *payload)
}

func UEencode(radiusMessage *radiusMessage.RadiusMessage) ([]byte, error) {

	radiusMessageData := make([]byte, 4)

	radiusMessageData[0] = radiusMessage.Code
	radiusMessageData[1] = radiusMessage.PktID
	radiusMessageData = append(radiusMessageData, radiusMessage.Auth...)

	radiusMessagePayloadData, err := radiusMessage.Payloads.Encode()
	if err != nil {
		return nil, fmt.Errorf("Encode(): EncodePayload failed: %+v", err)
	}
	radiusMessageData = append(radiusMessageData, radiusMessagePayloadData...)
	binary.BigEndian.PutUint16(radiusMessageData[2:4], uint16(len(radiusMessageData)))

	return radiusMessageData, nil
}

func GetMessageAuthenticator(message *radiusMessage.RadiusMessage) []byte {
	radius_secret := []byte("free5gctngf")
	radiusMessageData := make([]byte, 4)

	radiusMessageData[0] = message.Code
	radiusMessageData[1] = message.PktID
	radiusMessageData = append(radiusMessageData, message.Auth...)

	radiusMessagePayloadData, err := message.Payloads.Encode()
	if err != nil {
		return nil
	}
	radiusMessageData = append(radiusMessageData, radiusMessagePayloadData...)
	binary.BigEndian.PutUint16(radiusMessageData[2:4], uint16(len(radiusMessageData)))

	hmacFun := hmac.New(md5.New, radius_secret) // radius_secret is same as cfg's radius_secret
	hmacFun.Write(radiusMessageData)
	return hmacFun.Sum(nil)
}

func TestTngfUE(t *testing.T) {
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
	tngfue := context.TNGFSelf().NewTngfUe()
	tngfue.PduSessionList = make(map[int64]*context.PDUSession)
	tngfue.TNGFChildSecurityAssociation = make(map[uint32]*context.ChildSecurityAssociation)
	tngfue.TemporaryExchangeMsgIDChildSAMapping = make(map[uint32]*context.ChildSecurityAssociation)

	tngfRadiusUDPAddr, err := net.ResolveUDPAddr("udp", tngfInfo_IPSecIfaceAddr+":1812")
	if err != nil {
		t.Fatalf("Resolve UDP address %s fail: %+v", tngfInfo_IPSecIfaceAddr+":1812", err)
	}
	tngfUDPAddr, err := net.ResolveUDPAddr("udp", tngfInfo_IPSecIfaceAddr+":500")
	if err != nil {
		t.Fatalf("Resolve UDP address %s fail: %+v", tngfInfo_IPSecIfaceAddr+":500", err)
	}

	udpConnection, err := setupUDPSocket()
	if err != nil {
		t.Fatalf("Setup UDP socket Fail: %v", err)
	}
	radiusConnection, err := setupRadiusSocket()
	if err != nil {
		t.Fatalf("Setup Radius socket Fail: %+v", err)
	}

	// calling station payload
	callingStationPayload := new(radiusMessage.RadiusPayload)
	callingStationPayload.Type = radiusMessage.TypeCallingStationId
	callingStationPayload.Length = uint8(19)
	callingStationPayload.Val = []byte("C4-85-08-77-A7-D1")
	// called station payload
	calledStationPayload := new(radiusMessage.RadiusPayload)
	calledStationPayload.Type = radiusMessage.TypeCalledStationId
	calledStationPayload.Length = uint8(30)
	calledStationPayload.Val = []byte("D4-6E-0E-65-AC-A2:free5gc-ap")
	// UE user name payload
	ueUserNamePayload := new(radiusMessage.RadiusPayload)
	ueUserNamePayload.Type = radiusMessage.TypeUserName
	ueUserNamePayload.Length = uint8(8)
	ueUserNamePayload.Val = []byte("tngfue")

	var pkt []byte

	// Step3: AAA message, send to tngf
	// create a new radius message
	ueRadiusMessage := new(radiusMessage.RadiusMessage)
	radiusAuthenticator := make([]byte, 16)
	rand.Read(radiusAuthenticator) // request authenticator is random
	if err != nil {
		fmt.Printf("Failed to decode hex string: %v\n", err)
		return
	}

	ueRadiusMessage.BuildRadiusHeader(radiusMessage.AccessRequest, 0x05, radiusAuthenticator)
	// create Radius payload
	ueRadiusPayload := new(radiusMessage.RadiusPayloadContainer)
	*ueRadiusPayload = append(*ueRadiusPayload, *ueUserNamePayload, *calledStationPayload, *callingStationPayload)

	// create EAP message (Identity) payload
	identifier, err := radiusHandler.GenerateRandomUint8()
	if err != nil {
		t.Errorf("Random number failed: %+v", err)
		return
	}
	BuildEAPIdentity(ueRadiusPayload, identifier, []byte("tngfue"))

	// create Authenticator payload
	authPayload := new(radiusMessage.RadiusPayload)
	authPayload.Type = radiusMessage.TypeMessageAuthenticator
	authPayload.Length = uint8(18)
	authPayload.Val = make([]byte, 16)

	ueRadiusMessage.Payloads = *ueRadiusPayload
	ueRadiusMessage.Payloads = append(ueRadiusMessage.Payloads, *authPayload)
	authPayload.Val = GetMessageAuthenticator(ueRadiusMessage)
	*ueRadiusPayload = append(*ueRadiusPayload, *authPayload)
	ueRadiusMessage.Payloads = *ueRadiusPayload

	pkt, err = UEencode(ueRadiusMessage)

	if err != nil {
		t.Fatalf("Radius Message Encoding error: %+v", err)
	}
	// send to tngf
	if _, err := radiusConnection.WriteToUDP(pkt, tngfRadiusUDPAddr); err != nil {
		t.Fatalf("Write Radius maessage fail: %+v", err)
	}

	// Step 4: receive TNGF reply
	buffer := make([]byte, 65535)
	n, _, err := radiusConnection.ReadFromUDP(buffer)
	if err != nil {
		t.Fatalf("Read Radius message failed: %+v", err)
	}

	// Step 5: 5GNAS
	ueRadiusMessage = new(radiusMessage.RadiusMessage)
	radiusAuthenticator, err = hex.DecodeString("ea408c3a615fc82899bb8f2fa2e374e9")
	if err != nil {
		fmt.Printf("Failed to decode hex string: %v\n", err)
		return
	}

	ueRadiusMessage.BuildRadiusHeader(radiusMessage.AccessRequest, 0x06, radiusAuthenticator)
	// create Radius payload
	ueRadiusPayload = new(radiusMessage.RadiusPayloadContainer)
	*ueRadiusPayload = append(*ueRadiusPayload, *ueUserNamePayload, *calledStationPayload, *callingStationPayload)

	// create EAP message (Expanded) payload
	identifier, err = radiusHandler.GenerateRandomUint8()
	if err != nil {
		t.Errorf("Random number failed: %+v", err)
		return
	}
	// EAP-5G vendor type data
	eapVendorTypeData := make([]byte, 2)
	eapVendorTypeData[0] = message.EAP5GType5GNAS
	// AN Parameters
	anParameters := tngfBuildEAP5GANParameters(mobileIdentity5GS)
	anParametersLength := make([]byte, 2)
	binary.BigEndian.PutUint16(anParametersLength, uint16(len(anParameters)))
	eapVendorTypeData = append(eapVendorTypeData, anParametersLength...)
	eapVendorTypeData = append(eapVendorTypeData, anParameters...)

	// NAS-PDU (Registration Request)
	ueSecurityCapability := ue.GetUESecurityCapability()
	registrationRequest := nasTestpacket.GetRegistrationRequest(nasMessage.RegistrationType5GSInitialRegistration,
		mobileIdentity5GS, nil, ueSecurityCapability, nil, nil, nil)

	nasLength := make([]byte, 2)
	binary.BigEndian.PutUint16(nasLength, uint16(len(registrationRequest)))
	eapVendorTypeData = append(eapVendorTypeData, nasLength...)
	eapVendorTypeData = append(eapVendorTypeData, registrationRequest...)

	BuildEAP5GNAS(ueRadiusPayload, identifier, eapVendorTypeData)

	ueRadiusMessage.Payloads = *ueRadiusPayload
	pkt, err = ueRadiusMessage.Encode()
	if err != nil {
		t.Fatalf("Radius Message Encoding error: %+v", err)
	}
	// Send to tngf
	if _, err := radiusConnection.WriteToUDP(pkt, tngfRadiusUDPAddr); err != nil {
		t.Fatalf("Write Radius maessage fail: %+v", err)
	}

	// Step 6: Receive TNGF reply
	buffer = make([]byte, 65535)
	n, _, err = radiusConnection.ReadFromUDP(buffer)
	if err != nil {
		t.Fatalf("Read Radius message failed: %+v", err)
	}

	err = ueRadiusMessage.Decode(buffer[:n])
	if err != nil {
		t.Fatalf("Decode Radius message failed: %+v", err)
	}
	var eapMessage []byte

	for _, radiusPayload := range ueRadiusMessage.Payloads {
		switch radiusPayload.Type {
		case radiusMessage.TypeEAPMessage:
			eapMessage = radiusPayload.Val
		}
	}
	eap := new(radiusMessage.EAP)
	err = eap.Unmarshal(eapMessage)
	if eap.Code != radiusMessage.EAPCodeRequest {
		t.Fatalf("[EAP] Received an EAP payload with code other than request. Drop the payload.")
	}

	eapTypeData := eap.EAPTypeData[0]
	var eapExpanded *radiusMessage.EAPExpanded

	var decodedNAS *nas.Message

	eapExpanded = eapTypeData.(*radiusMessage.EAPExpanded)

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

	// Send Authentication

	ueRadiusMessage = new(radiusMessage.RadiusMessage)
	radiusAuthenticator, err = hex.DecodeString("ea408c3a615fc82899bb8f2fa2e374e9")

	if err != nil {
		fmt.Printf("Failed to decode hex string: %v\n", err)
		return
	}
	ueRadiusMessage.BuildRadiusHeader(radiusMessage.AccessRequest, 0x07, radiusAuthenticator)
	ueRadiusPayload = new(radiusMessage.RadiusPayloadContainer)
	*ueRadiusPayload = append(*ueRadiusPayload, *ueUserNamePayload, *calledStationPayload, *callingStationPayload)
	// create EAP message (Expanded) payload
	identifier, err = radiusHandler.GenerateRandomUint8()
	if err != nil {
		t.Errorf("Random number failed: %+v", err)
		return
	}
	// EAP-5G vendor type data
	eapVendorTypeData = make([]byte, 2)
	eapVendorTypeData[0] = message.EAP5GType5GNAS

	// AN Parameters
	eapVendorTypeData = append(eapVendorTypeData, anParametersLength...)
	eapVendorTypeData = append(eapVendorTypeData, anParameters...)

	authenticationResponse := nasTestpacket.GetAuthenticationResponse(resStat, "")
	nasLength = make([]byte, 2)
	binary.BigEndian.PutUint16(nasLength, uint16(len(authenticationResponse)))
	eapVendorTypeData = append(eapVendorTypeData, nasLength...)
	eapVendorTypeData = append(eapVendorTypeData, authenticationResponse...)

	BuildEAP5GNAS(ueRadiusPayload, identifier, eapVendorTypeData)

	ueRadiusMessage.Payloads = *ueRadiusPayload
	pkt, err = ueRadiusMessage.Encode()
	if err != nil {
		t.Fatalf("Radius Message Encoding error: %+v", err)
	}
	// Send to tngf
	if _, err := radiusConnection.WriteToUDP(pkt, tngfRadiusUDPAddr); err != nil {
		t.Fatalf("Write Radius maessage fail: %+v", err)
	}

	// Step 9b: Receive TNGF reply
	buffer = make([]byte, 65535)
	n, _, err = radiusConnection.ReadFromUDP(buffer)
	if err != nil {
		t.Fatalf("Read Radius message failed: %+v", err)
	}

	err = ueRadiusMessage.Decode(buffer[:n])
	if err != nil {
		t.Fatalf("Decode Radius message failed: %+v", err)
	}

	// Step 9c:
	ueRadiusMessage = new(radiusMessage.RadiusMessage)
	radiusAuthenticator, err = hex.DecodeString("ea408c3a615fc82899bb8f2fa2e374e9")
	if err != nil {
		fmt.Printf("Failed to decode hex string: %v\n", err)
		return
	}

	ueRadiusMessage.BuildRadiusHeader(radiusMessage.AccessRequest, 0x08, radiusAuthenticator)
	// create Radius payload
	ueRadiusPayload = new(radiusMessage.RadiusPayloadContainer)
	*ueRadiusPayload = append(*ueRadiusPayload, *ueUserNamePayload, *calledStationPayload, *callingStationPayload)

	// create EAP message (Expanded) payload
	identifier, err = radiusHandler.GenerateRandomUint8()
	if err != nil {
		t.Errorf("Random number failed: %+v", err)
		return
	}
	// EAP-5G vendor type data
	eapVendorTypeData = make([]byte, 2)
	eapVendorTypeData[0] = message.EAP5GType5GNAS

	// AN Parameters
	anParameters = tngfBuildEAP5GANParameters(mobileIdentity5GS)
	anParametersLength = make([]byte, 2)
	binary.BigEndian.PutUint16(anParametersLength, uint16(len(anParameters)))
	eapVendorTypeData = append(eapVendorTypeData, anParametersLength...)
	eapVendorTypeData = append(eapVendorTypeData, anParameters...)

	// NAS-PDU (SMC Complete)
	registrationRequestWith5GMM := nasTestpacket.GetRegistrationRequest(nasMessage.RegistrationType5GSInitialRegistration,
		mobileIdentity5GS, nil, ueSecurityCapability, ue.Get5GMMCapability(), nil, nil)
	smcComplete := nasTestpacket.GetSecurityModeComplete(registrationRequestWith5GMM)
	smcComplete, err = EncodeNasPduWithSecurity(ue, smcComplete, nas.SecurityHeaderTypeIntegrityProtectedAndCipheredWithNew5gNasSecurityContext, true, true)
	assert.Nil(t, err)
	nasLength = make([]byte, 2)
	binary.BigEndian.PutUint16(nasLength, uint16(len(smcComplete)))
	eapVendorTypeData = append(eapVendorTypeData, nasLength...)
	eapVendorTypeData = append(eapVendorTypeData, smcComplete...)

	BuildEAP5GNAS(ueRadiusPayload, identifier, eapVendorTypeData)

	ueRadiusMessage.Payloads = *ueRadiusPayload
	pkt, err = ueRadiusMessage.Encode()
	if err != nil {
		t.Fatalf("Radius Message Encoding error: %+v", err)
	}
	// Send to tngf
	if _, err := radiusConnection.WriteToUDP(pkt, tngfRadiusUDPAddr); err != nil {
		t.Fatalf("Write Radius maessage fail: %+v", err)
	}

	// Step 10b: Receive TNGF reply
	buffer = make([]byte, 65535)
	n, _, err = radiusConnection.ReadFromUDP(buffer)
	if err != nil {
		t.Fatalf("Read Radius message failed: %+v", err)
	}

	err = ueRadiusMessage.Decode(buffer[:n])
	if err != nil {
		t.Fatalf("Decode Radius message failed: %+v", err)
	}

	// 10c: EAP-Res/5G-Notification
	ueRadiusMessage = new(radiusMessage.RadiusMessage)
	radiusAuthenticator, err = hex.DecodeString("ea408c3a615fc82899bb8f2fa2e374e9")
	if err != nil {
		fmt.Printf("Failed to decode hex string: %v\n", err)
		return
	}

	ueRadiusMessage.BuildRadiusHeader(radiusMessage.AccessRequest, 0x09, radiusAuthenticator)
	// create Radius payload
	ueRadiusPayload = new(radiusMessage.RadiusPayloadContainer)
	*ueRadiusPayload = append(*ueRadiusPayload, *ueUserNamePayload, *calledStationPayload, *callingStationPayload)

	// create EAP message (Expanded) payload
	identifier, err = radiusHandler.GenerateRandomUint8()
	if err != nil {
		t.Errorf("Random number failed: %+v", err)
		return
	}
	BuildEAP5GNotification(ueRadiusPayload, identifier)

	ueRadiusMessage.Payloads = *ueRadiusPayload
	pkt, err = ueRadiusMessage.Encode()
	if err != nil {
		t.Fatalf("Radius Message Encoding error: %+v", err)
	}
	// Send to tngf
	if _, err := radiusConnection.WriteToUDP(pkt, tngfRadiusUDPAddr); err != nil {
		t.Fatalf("Write Radius maessage fail: %+v", err)
	}

	// 10e: EAP-Success
	// Receive TNGF reply
	buffer = make([]byte, 65535)
	n, _, err = radiusConnection.ReadFromUDP(buffer)
	if err != nil {
		t.Fatalf("Read Radius message failed: %+v", err)
	}

	err = ueRadiusMessage.Decode(buffer[:n])
	if err != nil {
		t.Fatalf("Decode Radius message failed: %+v", err)
	}

	// time.Sleep(10000 * time.Millisecond)
	// IKE_SA_INIT
	ikeInitiatorSPI := uint64(123123)
	ikeMessage := new(message.IKEMessage)
	ikeMessage.BuildIKEHeader(ikeInitiatorSPI, 0, message.IKE_SA_INIT, message.InitiatorBitCheck, 0)

	// Security Association
	securityAssociation := ikeMessage.Payloads.BuildSecurityAssociation()
	// Proposal 1
	proposal := securityAssociation.Proposals.BuildProposal(1, message.TypeIKE, nil)
	// ENCR
	proposal.EncryptionAlgorithm.BuildTransform(message.TypeEncryptionAlgorithm, message.ENCR_NULL, nil, nil, nil)
	// INTEG
	proposal.IntegrityAlgorithm.BuildTransform(message.TypeIntegrityAlgorithm, message.AUTH_HMAC_SHA1_96, nil, nil, nil)
	// PRF
	proposal.PseudorandomFunction.BuildTransform(message.TypePseudorandomFunction, message.PRF_HMAC_SHA1, nil, nil, nil)
	// DH
	proposal.DiffieHellmanGroup.BuildTransform(message.TypeDiffieHellmanGroup, message.DH_2048_BIT_MODP, nil, nil, nil)

	// Key exchange data
	generator := new(big.Int).SetUint64(handler.Group14Generator)
	factor, ok := new(big.Int).SetString(handler.Group14PrimeString, 16)
	if !ok {
		t.Fatalf("Generate key exchange data failed")
	}
	secert := handler.GenerateRandomNumber()
	localPublicKeyExchangeValue := new(big.Int).Exp(generator, secert, factor).Bytes()
	prependZero := make([]byte, len(factor.Bytes())-len(localPublicKeyExchangeValue))
	localPublicKeyExchangeValue = append(prependZero, localPublicKeyExchangeValue...)
	ikeMessage.Payloads.BUildKeyExchange(message.DH_2048_BIT_MODP, localPublicKeyExchangeValue)

	// Nonce
	localNonce := handler.GenerateRandomNumber().Bytes()
	ikeMessage.Payloads.BuildNonce(localNonce)

	// Send to TNGF
	ikeMessageData, err := ikeMessage.Encode()
	if err != nil {
		t.Fatalf("Encode IKE Message fail: %+v", err)
	}
	if _, err := udpConnection.WriteToUDP(ikeMessageData, tngfUDPAddr); err != nil {
		t.Fatalf("Write IKE maessage fail: %+v", err)
	}
	realMessage1, _ := ikeMessage.Encode()
	ikeSecurityAssociation := &context.IKESecurityAssociation{
		ResponderSignedOctets: realMessage1,
	}

	// Receive TNGF reply
	buffer = make([]byte, 65535)
	n, _, err = udpConnection.ReadFromUDP(buffer)
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
		case message.TypeSA:
			t.Log("Get SA payload")
		case message.TypeKE:
			remotePublicKeyExchangeValue := ikePayload.(*message.KeyExchange).KeyExchangeData
			var i int = 0
			for {
				if remotePublicKeyExchangeValue[i] != 0 {
					break
				}
			}
			remotePublicKeyExchangeValue = remotePublicKeyExchangeValue[i:]
			remotePublicKeyExchangeValueBig := new(big.Int).SetBytes(remotePublicKeyExchangeValue)
			sharedKeyExchangeData = new(big.Int).Exp(remotePublicKeyExchangeValueBig, secert, factor).Bytes()
		case message.TypeNiNr:
			remoteNonce = ikePayload.(*message.Nonce).NonceData
		}
	}

	ikeSecurityAssociation = &context.IKESecurityAssociation{
		LocalSPI:               ikeInitiatorSPI,
		RemoteSPI:              ikeMessage.ResponderSPI,
		InitiatorMessageID:     0,
		ResponderMessageID:     0,
		EncryptionAlgorithm:    proposal.EncryptionAlgorithm[0],
		IntegrityAlgorithm:     proposal.IntegrityAlgorithm[0],
		PseudorandomFunction:   proposal.PseudorandomFunction[0],
		DiffieHellmanGroup:     proposal.DiffieHellmanGroup[0],
		ConcatenatedNonce:      append(localNonce, remoteNonce...),
		DiffieHellmanSharedKey: sharedKeyExchangeData,
		ResponderSignedOctets:  append(ikeSecurityAssociation.ResponderSignedOctets, remoteNonce...),
	}

	if err := tngfGenerateKeyForIKESA(ikeSecurityAssociation); err != nil {
		t.Fatalf("Generate key for IKE SA failed: %+v", err)
	}

	tngfue.TNGFIKESecurityAssociation = ikeSecurityAssociation

	// IKE_AUTH (negociate IKE_CHILD_SA)
	ikeMessage.Payloads.Reset()
	tngfue.TNGFIKESecurityAssociation.InitiatorMessageID++
	ikeMessage.BuildIKEHeader(
		tngfue.TNGFIKESecurityAssociation.LocalSPI, tngfue.TNGFIKESecurityAssociation.RemoteSPI,
		message.IKE_AUTH, message.InitiatorBitCheck, tngfue.TNGFIKESecurityAssociation.InitiatorMessageID)

	var ikePayload message.IKEPayloadContainer

	// Identification
	ikePayload.BuildIdentificationInitiator(message.ID_KEY_ID, mobileIdentity5GS.GetMobileIdentity5GSContents())

	// Security Association
	securityAssociation = ikePayload.BuildSecurityAssociation()
	// Proposal 1
	inboundSPI := tngfGenerateSPI(tngfue)
	proposal = securityAssociation.Proposals.BuildProposal(1, message.TypeESP, inboundSPI)
	// ENCR (use null encryption for ESP)
	proposal.EncryptionAlgorithm.BuildTransform(message.TypeEncryptionAlgorithm, message.ENCR_NULL, nil, nil, nil)
	// INTEG
	proposal.IntegrityAlgorithm.BuildTransform(message.TypeIntegrityAlgorithm, message.AUTH_HMAC_SHA1_96, nil, nil, nil)
	// ESN
	proposal.ExtendedSequenceNumbers.BuildTransform(message.TypeExtendedSequenceNumbers, message.ESN_NO, nil, nil, nil)

	// Traffic Selector
	tsi := ikePayload.BuildTrafficSelectorInitiator()
	tsi.TrafficSelectors.BuildIndividualTrafficSelector(message.TS_IPV4_ADDR_RANGE, 0, 0, 65535, []byte{0, 0, 0, 0}, []byte{255, 255, 255, 255})
	tsr := ikePayload.BuildTrafficSelectorResponder()
	tsr.TrafficSelectors.BuildIndividualTrafficSelector(message.TS_IPV4_ADDR_RANGE, 0, 0, 65535, []byte{0, 0, 0, 0}, []byte{255, 255, 255, 255})

	// Authentication
	// Derive Ktngf
	P0 := make([]byte, 4)
	binary.BigEndian.PutUint32(P0, ue.ULCount.Get()-1)
	L0 := ueauth.KDFLen(P0)
	P1 := []byte{nasSecurity.AccessTypeNon3GPP}
	L1 := ueauth.KDFLen(P1)
	Ktngf, err := ueauth.GetKDFValue(ue.Kamf, ueauth.FC_FOR_KGNB_KN3IWF_DERIVATION, P0, L0, P1, L1)
	if err != nil {
		t.Fatalf("Get Ktngf error : %+v", err)
	}

	fmt.Println("ktngf: ", hex.Dump(Ktngf))

	pseudorandomFunction, random_ok := handler.NewPseudorandomFunction(ikeSecurityAssociation.SK_pi,
		ikeSecurityAssociation.PseudorandomFunction.TransformID)
	if !random_ok {
		t.Fatal("Get an unsupported pseudorandom funcion. This may imply an unsupported transform is chosen.")
		return
	}
	var idPayload message.IKEPayloadContainer
	idPayload.BuildIdentificationInitiator(message.ID_KEY_ID, mobileIdentity5GS.GetMobileIdentity5GSContents())
	idPayloadData, err := idPayload.Encode()
	if err != nil {
		t.Fatalf("Encode IKE payload failed : %+v", err)
	}
	if _, err := pseudorandomFunction.Write(idPayloadData[4:]); err != nil {
		t.Fatalf("Pseudorandom function write error: %+v", err)
	}
	ikeSecurityAssociation.ResponderSignedOctets = append(ikeSecurityAssociation.ResponderSignedOctets,
		pseudorandomFunction.Sum(nil)...)

	p0 := []byte{0x01}
	Ktipsec, err := ueauth.GetKDFValue(Ktngf, ueauth.FC_FOR_KTIPSEC_KTNAP_DERIVATION, p0, ueauth.KDFLen(p0))
	if err != nil {
		t.Fatal("UE authentication get KDF value Fatal.")
		return
	}
	fmt.Println("ktipsec: ", hex.Dump(Ktipsec))
	pseudorandomFunction, ok = handler.NewPseudorandomFunction(Ktipsec, ikeSecurityAssociation.PseudorandomFunction.TransformID)
	if !ok {
		t.Fatal("Get an unsupported pseudorandom funcion. This may imply an unsupported transform is chosen.")
		return
	}
	if _, random_err := pseudorandomFunction.Write([]byte("Key Pad for IKEv2")); random_err != nil {
		t.Fatalf("Pseudorandom function write error: %+v", random_err)
		return
	}

	secret := pseudorandomFunction.Sum(nil)
	t.Log("Using key to authentication:", hex.Dump(secret))
	pseudorandomFunction, ok = handler.NewPseudorandomFunction(secret, ikeSecurityAssociation.PseudorandomFunction.TransformID)
	if !ok {
		t.Fatal("Get an unsupported pseudorandom funcion. This may imply an unsupported transform is chosen.")
		return
	}

	pseudorandomFunction.Reset()
	t.Log("InitoatorSignedOctets: ", hex.Dump(ikeSecurityAssociation.ResponderSignedOctets))
	if _, random_err := pseudorandomFunction.Write(ikeSecurityAssociation.ResponderSignedOctets); random_err != nil {
		t.Fatalf("Pseudorandom function write error: %+v", random_err)
		return
	}

	ikePayload.BuildAuthentication(message.SharedKeyMesageIntegrityCode, pseudorandomFunction.Sum(nil))

	// Configuration Request
	configurationRequest := ikePayload.BuildConfiguration(message.CFG_REQUEST)
	configurationRequest.ConfigurationAttribute.BuildConfigurationAttribute(message.INTERNAL_IP4_ADDRESS, nil)

	if err := tngfEncryptProcedure(ikeSecurityAssociation, ikePayload, ikeMessage); err != nil {
		t.Fatalf("Encrypting IKE message failed: %+v", err)
	}

	// Send to TNGF
	ikeMessageData, err = ikeMessage.Encode()
	if err != nil {
		t.Fatalf("Encode IKE message failed: %+v", err)
	}
	if _, err := udpConnection.WriteToUDP(ikeMessageData, tngfUDPAddr); err != nil {
		t.Fatalf("Write IKE message failed: %+v", err)
	}

	tngfue.CreateHalfChildSA(tngfue.TNGFIKESecurityAssociation.InitiatorMessageID, binary.BigEndian.Uint32(inboundSPI), -1)

	// Receive TNGF reply
	n, _, err = udpConnection.ReadFromUDP(buffer)
	if err != nil {
		t.Fatalf("Read IKE message failed: %+v", err)
	}
	ikeMessage.Payloads.Reset()
	err = ikeMessage.Decode(buffer[:n])
	if err != nil {
		t.Fatalf("Decode IKE message failed: %+v", err)
	}

	encryptedPayload, ok := ikeMessage.Payloads[0].(*message.Encrypted)
	if !ok {
		t.Fatalf("Received payload is not an encrypted payload")
	}

	decryptedIKEPayload, err := tngfDecryptProcedure(ikeSecurityAssociation, ikeMessage, encryptedPayload)
	if err != nil {
		t.Fatalf("Decrypt IKE message failed: %+v", err)
	}

	// AUTH, SAr2, TSi, Tsr, N(NAS_IP_ADDRESS), N(NAS_TCP_PORT)
	var responseSecurityAssociation *message.SecurityAssociation
	var responseTrafficSelectorInitiator *message.TrafficSelectorInitiator
	var responseTrafficSelectorResponder *message.TrafficSelectorResponder
	var responseConfiguration *message.Configuration
	tngfNASAddr := new(net.TCPAddr)

	for _, ikePayload := range decryptedIKEPayload {
		switch ikePayload.Type() {
		case message.TypeAUTH:
			t.Log("Get Authentication from TNGF")
		case message.TypeSA:
			responseSecurityAssociation = ikePayload.(*message.SecurityAssociation)
			tngfue.TNGFIKESecurityAssociation.IKEAuthResponseSA = responseSecurityAssociation
		case message.TypeTSi:
			responseTrafficSelectorInitiator = ikePayload.(*message.TrafficSelectorInitiator)
		case message.TypeTSr:
			responseTrafficSelectorResponder = ikePayload.(*message.TrafficSelectorResponder)
		case message.TypeN:
			notification := ikePayload.(*message.Notification)
			if notification.NotifyMessageType == message.Vendor3GPPNotifyTypeNAS_IP4_ADDRESS {
				tngfNASAddr.IP = net.IPv4(notification.NotificationData[0], notification.NotificationData[1], notification.NotificationData[2], notification.NotificationData[3])
			}
			if notification.NotifyMessageType == message.Vendor3GPPNotifyTypeNAS_TCP_PORT {
				tngfNASAddr.Port = int(binary.BigEndian.Uint16(notification.NotificationData))
			}
		case message.TypeCP:
			responseConfiguration = ikePayload.(*message.Configuration)
			if responseConfiguration.ConfigurationType == message.CFG_REPLY {
				for _, configAttr := range responseConfiguration.ConfigurationAttribute {
					if configAttr.Type == message.INTERNAL_IP4_ADDRESS {
						ueInnerAddr.IP = configAttr.Value
					}
					if configAttr.Type == message.INTERNAL_IP4_NETMASK {
						ueInnerAddr.Mask = configAttr.Value
					}
				}
			}
		}
	}

	OutboundSPI := binary.BigEndian.Uint32(tngfue.TNGFIKESecurityAssociation.IKEAuthResponseSA.Proposals[0].SPI)
	childSecurityAssociationContext, err := tngfue.CompleteChildSA(
		0x01, OutboundSPI, tngfue.TNGFIKESecurityAssociation.IKEAuthResponseSA)
	if err != nil {
		t.Fatalf("Create child security association context failed: %+v", err)
	}
	err = tngfParseIPAddressInformationToChildSecurityAssociation(childSecurityAssociationContext,
		responseTrafficSelectorInitiator.TrafficSelectors[0],
		responseTrafficSelectorResponder.TrafficSelectors[0])

	if err != nil {
		t.Fatalf("Parse IP address to child security association failed: %+v", err)
	}
	// Select TCP traffic
	childSecurityAssociationContext.SelectedIPProtocol = unix.IPPROTO_TCP

	if err := handler.GenerateKeyForChildSA(ikeSecurityAssociation, childSecurityAssociationContext); err != nil {
		t.Fatalf("Generate key for child SA failed: %+v", err)
	}

	var linkIPSec netlink.Link

	// Setup interface for ipsec
	newXfrmiName := fmt.Sprintf("%s-default", tngfueInfo_XfrmiName)
	if linkIPSec, err = tngfSetupIPsecXfrmi(newXfrmiName, tngfueInfo_IPSecIfaceName, tngfueInfo_XfrmiId, ueInnerAddr); err != nil {
		t.Fatalf("Setup XFRM interface %s fail: %+v", newXfrmiName, err)
	}

	defer func() {
		if err := netlink.LinkDel(linkIPSec); err != nil {
			t.Fatalf("Delete XFRM interface %s fail: %+v", newXfrmiName, err)
		} else {
			t.Logf("Delete XFRM interface: %s", newXfrmiName)
		}
	}()

	// Apply XFRM rules
	if err = tngfApplyXFRMRule(true, tngfueInfo_XfrmiId, childSecurityAssociationContext); err != nil {
		t.Fatalf("Applying XFRM rules failed: %+v", err)
	}

	defer func() {
		_ = netlink.XfrmPolicyFlush()
		_ = netlink.XfrmStateFlush(netlink.XFRM_PROTO_IPSEC_ANY)
	}()

	// Waiting for downlink data to be stored in the cache, then establishing a TCP connection.
	time.Sleep(1 * time.Second)

	localTCPAddr := &net.TCPAddr{
		IP: ueInnerAddr.IP,
	}
	tcpConnWithTNGF, err := net.DialTCP("tcp", localTCPAddr, tngfNASAddr)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("Establish TCP connection with TNGF successfully!!")

	// For timeout
	tcpConnWithTNGF.SetReadDeadline(time.Now().Add(20 * time.Second))

	nasEnv := make([]byte, 65535)

	n, err = tcpConnWithTNGF.Read(nasEnv)
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
	t.Logf("Get Registration Accept Message:\n %+v", nasStr)

	// send NAS Registration Complete Msg
	pdu := nasTestpacket.GetRegistrationComplete(nil)
	pdu, err = EncodeNasPduInEnvelopeWithSecurity(ue, pdu, nas.SecurityHeaderTypeIntegrityProtectedAndCiphered, true, false)
	if err != nil {
		t.Fatal(err)
		return
	}
	_, err = tcpConnWithTNGF.Write(pdu)
	if err != nil {
		t.Fatal(err)
		return
	}

	// Do not read the response, just wait for the AMF to finish the registration process
	time.Sleep(1 * time.Second)

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
	_, err = tcpConnWithTNGF.Write(pdu)
	if err != nil {
		t.Fatal(err)
		return
	}

	// Receive TNGF reply
	n, _, err = udpConnection.ReadFromUDP(buffer)
	if err != nil {
		t.Fatalf("Read IKE Message fail: %+v", err)
	}
	ikeMessage.Payloads.Reset()
	err = ikeMessage.Decode(buffer[:n])
	if err != nil {
		t.Fatalf("Decode IKE Message fail: %+v", err)
	}
	t.Logf("IKE message exchange type: %d", ikeMessage.ExchangeType)
	t.Logf("IKE message ID: %d", ikeMessage.MessageID)
	encryptedPayload, ok = ikeMessage.Payloads[0].(*message.Encrypted)
	if !ok {
		t.Fatal("Received pakcet is not an encrypted payload")
		return
	}
	decryptedIKEPayload, err = tngfDecryptProcedure(ikeSecurityAssociation, ikeMessage, encryptedPayload)
	if err != nil {
		t.Fatal(err)
		return
	}

	var QoSInfo *PDUQoSInfo

	var upIPAddr net.IP
	for _, ikePayload := range decryptedIKEPayload {
		switch ikePayload.Type() {
		case message.TypeSA:
			responseSecurityAssociation = ikePayload.(*message.SecurityAssociation)
			OutboundSPI = binary.BigEndian.Uint32(responseSecurityAssociation.Proposals[0].SPI)
		case message.TypeTSi:
			responseTrafficSelectorInitiator = ikePayload.(*message.TrafficSelectorInitiator)
		case message.TypeTSr:
			responseTrafficSelectorResponder = ikePayload.(*message.TrafficSelectorResponder)
		case message.TypeN:
			notification := ikePayload.(*message.Notification)
			if notification.NotifyMessageType == message.Vendor3GPPNotifyType5G_QOS_INFO {
				t.Log("Received Qos Flow settings")
				if info, err := tngfParse5GQoSInfoNotify(notification); err == nil {
					QoSInfo = info
					t.Logf("NotificationData:%+v", notification.NotificationData)
					if QoSInfo.isDSCPSpecified {
						t.Logf("DSCP is specified but test not support")
					}
				} else {
					t.Logf("%+v", err)
				}
			}
			if notification.NotifyMessageType == message.Vendor3GPPNotifyTypeUP_IP4_ADDRESS {
				upIPAddr = notification.NotificationData[:4]
				t.Logf("UP IP Address: %+v\n", upIPAddr)
			}
		case message.TypeNiNr:
			responseNonce := ikePayload.(*message.Nonce)
			ikeSecurityAssociation.ConcatenatedNonce = responseNonce.NonceData
		}
	}

	// IKE CREATE_CHILD_SA response
	ikeMessage.Payloads.Reset()
	tngfue.TNGFIKESecurityAssociation.ResponderMessageID = ikeMessage.MessageID
	ikeMessage.BuildIKEHeader(ikeMessage.InitiatorSPI, ikeMessage.ResponderSPI, message.CREATE_CHILD_SA,
		message.ResponseBitCheck|message.InitiatorBitCheck, tngfue.TNGFIKESecurityAssociation.ResponderMessageID)

	ikePayload.Reset()

	// SA
	inboundSPI = tngfGenerateSPI(tngfue)
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
	localNonce = localNonceBigInt.Bytes()
	ikeSecurityAssociation.ConcatenatedNonce = append(ikeSecurityAssociation.ConcatenatedNonce, localNonce...)
	ikePayload.BuildNonce(localNonce)

	if err := tngfEncryptProcedure(ikeSecurityAssociation, ikePayload, ikeMessage); err != nil {
		t.Fatalf("Encrypt IKE message failed: %+v", err)
	}

	// Send to TNGF
	ikeMessageData, err = ikeMessage.Encode()
	if err != nil {
		t.Fatalf("Encode IKE Message fail: %+v", err)
	}
	_, err = udpConnection.WriteToUDP(ikeMessageData, tngfUDPAddr)
	if err != nil {
		t.Fatalf("Write IKE message failed: %+v", err)
	}

	tngfue.CreateHalfChildSA(tngfue.TNGFIKESecurityAssociation.ResponderMessageID, binary.BigEndian.Uint32(inboundSPI), -1)
	childSecurityAssociationContextUserPlane, err := tngfue.CompleteChildSA(
		tngfue.TNGFIKESecurityAssociation.ResponderMessageID, OutboundSPI, responseSecurityAssociation)

	if err != nil {
		t.Fatalf("Create child security association context failed: %+v", err)
	}
	err = tngfParseIPAddressInformationToChildSecurityAssociation(childSecurityAssociationContextUserPlane, responseTrafficSelectorResponder.TrafficSelectors[0], responseTrafficSelectorInitiator.TrafficSelectors[0])
	if err != nil {
		t.Fatalf("Parse IP address to child security association failed: %+v", err)
	}
	// Select GRE traffic
	childSecurityAssociationContextUserPlane.SelectedIPProtocol = unix.IPPROTO_GRE

	if err := handler.GenerateKeyForChildSA(ikeSecurityAssociation, childSecurityAssociationContextUserPlane); err != nil {
		t.Fatalf("Generate key for child SA failed: %+v", err)
	}

	// Aplly XFRM rules
	if err = tngfApplyXFRMRule(false, tngfueInfo_XfrmiId, childSecurityAssociationContextUserPlane); err != nil {
		t.Fatalf("Applying XFRM rules failed: %+v", err)
	}

	var pduAddress net.IP

	// Read NAS from TNGF
	if n, err := tcpConnWithTNGF.Read(buffer); err != nil {
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

	newGREName := fmt.Sprintf("%s-id-%d", tngfueInfo_GreIfaceName, tngfueInfo_XfrmiId)

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
