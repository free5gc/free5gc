package test

import (
	"encoding/binary"
	"encoding/hex"
	"errors"
	"gofree5gc/lib/CommonConsumerTestData/UDM/TestGenAuthData"
	"gofree5gc/lib/nas"
	"gofree5gc/lib/nas/nasMessage"
	"gofree5gc/lib/nas/nasTestpacket"
	"gofree5gc/lib/nas/nasType"
	"gofree5gc/lib/openapi/models"
	"gofree5gc/src/n3iwf/n3iwf_context"
	"gofree5gc/src/n3iwf/n3iwf_ike/ike_handler"
	"gofree5gc/src/n3iwf/n3iwf_ike/ike_message"
	"hash"
	"math/big"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vishvananda/netlink"
	"golang.org/x/sys/unix"
)

func createIKEChildSecurityAssociation(chosenSecurityAssociation *ike_message.SecurityAssociation) (*n3iwf_context.ChildSecurityAssociation, error) {
	childSecurityAssociation := new(n3iwf_context.ChildSecurityAssociation)

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

	return childSecurityAssociation, nil
}

func getAuthSubscription() (authSubs models.AuthenticationSubscription) {
	authSubs.PermanentKey = &models.PermanentKey{
		PermanentKeyValue: TestGenAuthData.MilenageTestSet19.K,
	}
	authSubs.Opc = &models.Opc{
		OpcValue: TestGenAuthData.MilenageTestSet19.OPC,
	}
	authSubs.Milenage = &models.Milenage{
		Op: &models.Op{
			OpValue: TestGenAuthData.MilenageTestSet19.OP,
		},
	}
	authSubs.AuthenticationManagementField = "8000"

	authSubs.SequenceNumber = TestGenAuthData.MilenageTestSet19.SQN
	authSubs.AuthenticationMethod = models.AuthMethod__5_G_AKA
	return
}

func setupUDPSocket(t *testing.T) *net.UDPConn {
	bindAddr := "192.168.127.2:500"
	udpAddr, err := net.ResolveUDPAddr("udp", bindAddr)
	if err != nil {
		t.Fatal("Resolve UDP address failed")
	}
	udpListener, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		t.Fatal("Listen UDP socket failed")
	}
	return udpListener
}

func concatenateNonceAndSPI(nonce []byte, SPI_initiator uint64, SPI_responder uint64) []byte {
	spi := make([]byte, 8)

	binary.BigEndian.PutUint64(spi, SPI_initiator)
	newSlice := append(nonce, spi...)
	binary.BigEndian.PutUint64(spi, SPI_responder)
	newSlice = append(newSlice, spi...)

	return newSlice
}

func generateKeyForIKESA(ikeSecurityAssociation *n3iwf_context.IKESecurityAssociation) error {
	// Transforms
	transformPseudorandomFunction := ikeSecurityAssociation.PseudorandomFunction

	// Get key length of SK_d, SK_ai, SK_ar, SK_ei, SK_er, SK_pi, SK_pr
	var length_SK_d, length_SK_ai, length_SK_ar, length_SK_ei, length_SK_er, length_SK_pi, length_SK_pr, totalKeyLength int
	var ok bool

	length_SK_d = 20
	length_SK_ai = 20
	length_SK_ar = length_SK_ai
	length_SK_ei = 32
	length_SK_er = length_SK_ei
	length_SK_pi, length_SK_pr = length_SK_d, length_SK_d
	totalKeyLength = length_SK_d + length_SK_ai + length_SK_ar + length_SK_ei + length_SK_er + length_SK_pi + length_SK_pr

	// Generate IKE SA key as defined in RFC7296 Section 1.3 and Section 1.4
	var pseudorandomFunction hash.Hash

	if pseudorandomFunction, ok = ike_handler.NewPseudorandomFunction(ikeSecurityAssociation.ConcatenatedNonce, transformPseudorandomFunction.TransformID); !ok {
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
		if pseudorandomFunction, ok = ike_handler.NewPseudorandomFunction(SKEYSEED, transformPseudorandomFunction.TransformID); !ok {
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

func generateKeyForChildSA(ikeSecurityAssociation *n3iwf_context.IKESecurityAssociation, childSecurityAssociation *n3iwf_context.ChildSecurityAssociation) error {
	// Transforms
	transformPseudorandomFunction := ikeSecurityAssociation.PseudorandomFunction
	var transformIntegrityAlgorithmForIPSec *ike_message.Transform
	if len(ikeSecurityAssociation.IKEAuthResponseSA.Proposals[0].IntegrityAlgorithm) != 0 {
		transformIntegrityAlgorithmForIPSec = ikeSecurityAssociation.IKEAuthResponseSA.Proposals[0].IntegrityAlgorithm[0]
	}

	// Get key length for encryption and integrity key for IPSec
	var lengthEncryptionKeyIPSec, lengthIntegrityKeyIPSec, totalKeyLength int
	var ok bool

	lengthEncryptionKeyIPSec = 32
	if transformIntegrityAlgorithmForIPSec != nil {
		lengthIntegrityKeyIPSec = 20
	}
	totalKeyLength = lengthEncryptionKeyIPSec + lengthIntegrityKeyIPSec
	totalKeyLength = totalKeyLength * 2

	// Generate key for child security association as specified in RFC 7296 section 2.17
	seed := ikeSecurityAssociation.ConcatenatedNonce
	var pseudorandomFunction hash.Hash

	var keyStream, generatedKeyBlock []byte
	var index byte
	for index = 1; len(keyStream) < totalKeyLength; index++ {
		if pseudorandomFunction, ok = ike_handler.NewPseudorandomFunction(ikeSecurityAssociation.SK_d, transformPseudorandomFunction.TransformID); !ok {
			return errors.New("New pseudorandom function failed")
		}
		if _, err := pseudorandomFunction.Write(append(append(generatedKeyBlock, seed...), index)); err != nil {
			return errors.New("Pseudorandom function write failed")
		}
		generatedKeyBlock = pseudorandomFunction.Sum(nil)
		keyStream = append(keyStream, generatedKeyBlock...)
	}

	childSecurityAssociation.IncomingEncryptionKey = append(childSecurityAssociation.IncomingEncryptionKey, keyStream[:lengthEncryptionKeyIPSec]...)
	keyStream = keyStream[lengthEncryptionKeyIPSec:]
	childSecurityAssociation.IncomingIntegrityKey = append(childSecurityAssociation.IncomingIntegrityKey, keyStream[:lengthIntegrityKeyIPSec]...)
	keyStream = keyStream[lengthIntegrityKeyIPSec:]
	childSecurityAssociation.OutgoingEncryptionKey = append(childSecurityAssociation.OutgoingEncryptionKey, keyStream[:lengthEncryptionKeyIPSec]...)
	keyStream = keyStream[lengthEncryptionKeyIPSec:]
	childSecurityAssociation.OutgoingIntegrityKey = append(childSecurityAssociation.OutgoingIntegrityKey, keyStream[:lengthIntegrityKeyIPSec]...)

	return nil

}

func decryptProcedure(ikeSecurityAssociation *n3iwf_context.IKESecurityAssociation, message *ike_message.IKEMessage, encryptedPayload *ike_message.Encrypted) ([]ike_message.IKEPayloadType, error) {
	// Load needed information
	transformIntegrityAlgorithm := ikeSecurityAssociation.IntegrityAlgorithm
	transformEncryptionAlgorithm := ikeSecurityAssociation.EncryptionAlgorithm
	checksumLength := 12 // HMAC_SHA1_96

	// Checksum
	checksum := encryptedPayload.EncryptedData[len(encryptedPayload.EncryptedData)-checksumLength:]

	ikeMessageData, err := ike_message.Encode(message)
	if err != nil {
		return nil, errors.New("Encoding IKE message failed")
	}

	ok, err := ike_handler.VerifyIKEChecksum(ikeSecurityAssociation.SK_ar, ikeMessageData[:len(ikeMessageData)-checksumLength], checksum, transformIntegrityAlgorithm.TransformID)
	if err != nil {
		return nil, errors.New("Error verify checksum")
	}
	if !ok {
		return nil, errors.New("Checksum failed, drop.")
	}

	// Decrypt
	encryptedData := encryptedPayload.EncryptedData[:len(encryptedPayload.EncryptedData)-checksumLength]
	plainText, err := ike_handler.DecryptMessage(ikeSecurityAssociation.SK_er, encryptedData, transformEncryptionAlgorithm.TransformID)
	if err != nil {
		return nil, errors.New("Error decrypting message")
	}

	decryptedIKEPayload, err := ike_message.DecodePayload(encryptedPayload.NextPayload, plainText)
	if err != nil {
		return nil, errors.New("Decoding decrypted payload failed")
	}

	return decryptedIKEPayload, nil

}

func encryptProcedure(ikeSecurityAssociation *n3iwf_context.IKESecurityAssociation, ikePayload []ike_message.IKEPayloadType, responseIKEMessage *ike_message.IKEMessage) error {
	// Load needed information
	transformIntegrityAlgorithm := ikeSecurityAssociation.IntegrityAlgorithm
	transformEncryptionAlgorithm := ikeSecurityAssociation.EncryptionAlgorithm
	checksumLength := 12 // HMAC_SHA1_96

	// Encrypting
	notificationPayloadData, err := ike_message.EncodePayload(ikePayload)
	if err != nil {
		return errors.New("Encoding IKE payload failed.")
	}

	encryptedData, err := ike_handler.EncryptMessage(ikeSecurityAssociation.SK_ei, notificationPayloadData, transformEncryptionAlgorithm.TransformID)
	if err != nil {
		return errors.New("Error encrypting message")
	}

	encryptedData = append(encryptedData, make([]byte, checksumLength)...)
	responseEncryptedPayload := ike_message.BuildEncrypted(ikePayload[0].Type(), encryptedData)

	responseIKEMessage.IKEPayload = append(responseIKEMessage.IKEPayload, responseEncryptedPayload)

	// Calculate checksum
	responseIKEMessageData, err := ike_message.Encode(responseIKEMessage)
	if err != nil {
		return errors.New("Encoding IKE message error")
	}
	checksumOfMessage, err := ike_handler.CalculateChecksum(ikeSecurityAssociation.SK_ai, responseIKEMessageData[:len(responseIKEMessageData)-checksumLength], transformIntegrityAlgorithm.TransformID)
	if err != nil {
		return errors.New("Error calculating checksum")
	}
	checksumField := responseEncryptedPayload.EncryptedData[len(responseEncryptedPayload.EncryptedData)-checksumLength:]
	copy(checksumField, checksumOfMessage)

	return nil

}

func buildEAP5GANParameters() []byte {
	var anParameters []byte

	// Build GUAMI
	anParameter := make([]byte, 2)
	guami := make([]byte, 7)
	guami[1] = 0x02
	guami[2] = 0xf8
	guami[3] = 0x39
	guami[4] = 0xca
	guami[5] = 0xfe
	guami[6] = 0x0
	anParameter[0] = ike_message.ANParametersTypeGUAMI
	anParameter[1] = byte(len(guami))
	anParameter = append(anParameter, guami...)

	anParameters = append(anParameters, anParameter...)

	// Build Establishment Cause
	anParameter = make([]byte, 2)
	establishmentCause := make([]byte, 2)
	establishmentCause[1] = ike_message.EstablishmentCauseMO_Data
	anParameter[0] = ike_message.ANParametersTypeEstablishmentCause
	anParameter[1] = byte(len(establishmentCause))
	anParameter = append(anParameter, establishmentCause...)

	anParameters = append(anParameters, anParameter...)

	// Build PLMN ID
	anParameter = make([]byte, 2)
	plmnID := make([]byte, 5)
	plmnID[1] = 3
	plmnID[2] = 0x02
	plmnID[3] = 0xf8
	plmnID[4] = 0x39
	anParameter[0] = ike_message.ANParametersTypeSelectedPLMNID
	anParameter[1] = byte(len(plmnID))
	anParameter = append(anParameter, plmnID...)

	anParameters = append(anParameters, anParameter...)

	// Build NSSAI
	anParameter = make([]byte, 2)
	nssai := make([]byte, 2)
	snssai := make([]byte, 6)
	snssai[1] = 4
	snssai[2] = 1
	snssai[3] = 0x01
	snssai[4] = 0x02
	snssai[5] = 0x03
	nssai = append(nssai, snssai...)
	snssai = make([]byte, 6)
	snssai[1] = 4
	snssai[2] = 1
	snssai[3] = 0x11
	snssai[4] = 0x22
	snssai[5] = 0x33
	nssai = append(nssai, snssai...)
	nssai[1] = 12
	anParameter[0] = ike_message.ANParametersTypeRequestedNSSAI
	anParameter[1] = byte(len(nssai))
	anParameter = append(anParameter, nssai...)

	anParameters = append(anParameters, anParameter...)

	return anParameters
}

func parseIPAddressInformationToChildSecurityAssociation(
	childSecurityAssociation *n3iwf_context.ChildSecurityAssociation,
	n3iwfPublicIPAddr net.IP,
	trafficSelectorInitiator *ike_message.TrafficSelectorInitiator,
	trafficSelectorResponder *ike_message.TrafficSelectorResponder) error {

	if childSecurityAssociation == nil {
		return errors.New("childSecurityAssociation is nil")
	}

	childSecurityAssociation.PeerPublicIPAddr = n3iwfPublicIPAddr
	childSecurityAssociation.LocalPublicIPAddr = net.ParseIP("192.168.127.2")

	childSecurityAssociation.TrafficSelectorInitiator = net.IPNet{
		IP:   trafficSelectorInitiator.TrafficSelectors[0].StartAddress,
		Mask: []byte{255, 255, 255, 255},
	}

	childSecurityAssociation.TrafficSelectorResponder = net.IPNet{
		IP:   trafficSelectorResponder.TrafficSelectors[0].StartAddress,
		Mask: []byte{255, 255, 255, 255},
	}

	return nil
}

func applyXFRMRule(childSecurityAssociation *n3iwf_context.ChildSecurityAssociation) error {
	// Build XFRM information data structure for incoming traffic.

	// Mark
	mark := &netlink.XfrmMark{
		Value: 5,
	}

	// Direction: {private_network} -> this_server
	// State
	xfrmEncryptionAlgorithm := &netlink.XfrmStateAlgo{
		Name: ike_handler.XFRMEncryptionAlgorithmType(childSecurityAssociation.EncryptionAlgorithm).String(),
		Key:  childSecurityAssociation.IncomingEncryptionKey,
	}
	xfrmIntegrityAlgorithm := &netlink.XfrmStateAlgo{
		Name: ike_handler.XFRMIntegrityAlgorithmType(childSecurityAssociation.IntegrityAlgorithm).String(),
		Key:  childSecurityAssociation.IncomingIntegrityKey,
	}

	xfrmState := new(netlink.XfrmState)

	xfrmState.Src = childSecurityAssociation.LocalPublicIPAddr
	xfrmState.Dst = childSecurityAssociation.PeerPublicIPAddr
	xfrmState.Proto = netlink.XFRM_PROTO_ESP
	xfrmState.Mode = netlink.XFRM_MODE_TUNNEL
	xfrmState.Spi = int(childSecurityAssociation.SPI)
	xfrmState.Mark = mark
	xfrmState.Auth = xfrmIntegrityAlgorithm
	xfrmState.Crypt = xfrmEncryptionAlgorithm
	xfrmState.ESN = childSecurityAssociation.ESN

	// Commit xfrm state to netlink
	var err error
	if err = netlink.XfrmStateAdd(xfrmState); err != nil {
		return errors.New("Set XFRM state rule failed")
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

	xfrmPolicy.Src = &childSecurityAssociation.TrafficSelectorInitiator
	xfrmPolicy.Dst = &childSecurityAssociation.TrafficSelectorResponder
	xfrmPolicy.Proto = netlink.Proto(childSecurityAssociation.SelectedIPProtocol)
	xfrmPolicy.Dir = netlink.XFRM_DIR_OUT
	xfrmPolicy.Mark = mark
	xfrmPolicy.Tmpls = []netlink.XfrmPolicyTmpl{
		xfrmPolicyTemplate,
	}

	// Commit xfrm policy to netlink
	if err = netlink.XfrmPolicyAdd(xfrmPolicy); err != nil {
		return errors.New("Set XFRM policy rule failed")
	}

	// Direction: this_server -> {private_network}
	// State
	xfrmEncryptionAlgorithm.Key = childSecurityAssociation.OutgoingEncryptionKey
	xfrmIntegrityAlgorithm.Key = childSecurityAssociation.OutgoingIntegrityKey

	xfrmState.Src, xfrmState.Dst = xfrmState.Dst, xfrmState.Src

	// Commit xfrm state to netlink
	if err = netlink.XfrmStateAdd(xfrmState); err != nil {
		return errors.New("Set XFRM state rule failed")
	}

	// Policy
	xfrmPolicyTemplate.Src, xfrmPolicyTemplate.Dst = xfrmPolicyTemplate.Dst, xfrmPolicyTemplate.Src

	xfrmPolicy.Src, xfrmPolicy.Dst = xfrmPolicy.Dst, xfrmPolicy.Src
	xfrmPolicy.Dir = netlink.XFRM_DIR_IN
	xfrmPolicy.Tmpls = []netlink.XfrmPolicyTmpl{
		xfrmPolicyTemplate,
	}

	// Commit xfrm policy to netlink
	if err = netlink.XfrmPolicyAdd(xfrmPolicy); err != nil {
		return errors.New("Set XFRM policy rule failed")
	}

	return nil
}

func TestNon3GPPUE(t *testing.T) {
	// New UE
	ue := NewRanUeContext("imsi-2089300007487", 1, ALG_CIPHERING_128_NEA2, ALG_INTEGRITY_128_NIA2)
	ue.AmfUeNgapId = 1
	ue.AuthenticationSubs = getAuthSubscription()
	mobileIdentity5GS := nasType.MobileIdentity5GS{
		Len:    12, // suci
		Buffer: []uint8{0x01, 0x02, 0xf8, 0x39, 0xf0, 0xff, 0x00, 0x00, 0x00, 0x00, 0x47, 0x78},
	}

	n3iwfUDPAddr, err := net.ResolveUDPAddr("udp", "192.168.127.1:500")
	if err != nil {
		t.Fatal(err)
	}
	udpConnection := setupUDPSocket(t)

	// IKE_SA_INIT
	ikeMessage := ike_message.BuildIKEHeader(123123, 0, ike_message.IKE_SA_INIT, ike_message.InitiatorBitCheck, 0)

	// Security Association
	proposal := ike_message.BuildProposal(1, ike_message.TypeIKE, nil)
	var attributeType uint16 = ike_message.AttributeTypeKeyLength
	var keyLength uint16 = 256
	encryptTransform := ike_message.BuildTransform(ike_message.TypeEncryptionAlgorithm, ike_message.ENCR_AES_CBC, &attributeType, &keyLength, nil)
	ike_message.AppendTransformToProposal(proposal, encryptTransform)
	integrityTransform := ike_message.BuildTransform(ike_message.TypeIntegrityAlgorithm, ike_message.AUTH_HMAC_SHA1_96, nil, nil, nil)
	ike_message.AppendTransformToProposal(proposal, integrityTransform)
	pseudorandomFunctionTransform := ike_message.BuildTransform(ike_message.TypePseudorandomFunction, ike_message.PRF_HMAC_SHA1, nil, nil, nil)
	ike_message.AppendTransformToProposal(proposal, pseudorandomFunctionTransform)
	diffiehellmanTransform := ike_message.BuildTransform(ike_message.TypeDiffieHellmanGroup, ike_message.DH_2048_BIT_MODP, nil, nil, nil)
	ike_message.AppendTransformToProposal(proposal, diffiehellmanTransform)
	securityAssociation := ike_message.BuildSecurityAssociation([]*ike_message.Proposal{proposal})
	ikeMessage.IKEPayload = append(ikeMessage.IKEPayload, securityAssociation)

	// Key exchange data
	generator := new(big.Int).SetUint64(ike_handler.Group14Generator)
	factor, ok := new(big.Int).SetString(ike_handler.Group14PrimeString, 16)
	if !ok {
		t.Fatal("Generate key exchange datd failed")
	}
	secert := ike_handler.GenerateRandomNumber()
	localPublicKeyExchangeValue := new(big.Int).Exp(generator, secert, factor).Bytes()
	prependZero := make([]byte, len(factor.Bytes())-len(localPublicKeyExchangeValue))
	localPublicKeyExchangeValue = append(prependZero, localPublicKeyExchangeValue...)
	keyExchangeData := ike_message.BUildKeyExchange(ike_message.DH_2048_BIT_MODP, localPublicKeyExchangeValue)
	ikeMessage.IKEPayload = append(ikeMessage.IKEPayload, keyExchangeData)

	// Nonce
	localNonce := ike_handler.GenerateRandomNumber().Bytes()
	nonce := ike_message.BuildNonce(localNonce)
	ikeMessage.IKEPayload = append(ikeMessage.IKEPayload, nonce)

	// Send to N3IWF
	ikeMessageData, err := ike_message.Encode(ikeMessage)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := udpConnection.WriteToUDP(ikeMessageData, n3iwfUDPAddr); err != nil {
		t.Fatal(err)
	}

	// Receive N3IWF reply
	buffer := make([]byte, 65535)
	n, _, err := udpConnection.ReadFromUDP(buffer)
	if err != nil {
		t.Fatal(err)
	}
	ikeMessage, err = ike_message.Decode(buffer[:n])
	if err != nil {
		t.Fatal(err)
	}

	var sharedKeyExchangeData []byte
	var remoteNonce []byte

	for _, ikePayload := range ikeMessage.IKEPayload {
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

	ikeSecurityAssociation := &n3iwf_context.IKESecurityAssociation{
		LocalSPI:               123123,
		RemoteSPI:              ikeMessage.ResponderSPI,
		EncryptionAlgorithm:    encryptTransform,
		IntegrityAlgorithm:     integrityTransform,
		PseudorandomFunction:   pseudorandomFunctionTransform,
		DiffieHellmanGroup:     diffiehellmanTransform,
		ConcatenatedNonce:      append(localNonce, remoteNonce...),
		DiffieHellmanSharedKey: sharedKeyExchangeData,
	}

	t.Logf("Concatenated Nonce:\n%s\n", hex.Dump(ikeSecurityAssociation.ConcatenatedNonce))
	t.Logf("Diffie-Hellman Value:\n%s\n", hex.Dump(ikeSecurityAssociation.DiffieHellmanSharedKey))

	if err := generateKeyForIKESA(ikeSecurityAssociation); err != nil {
		t.Fatalf("Generate key for IKE SA failed: %+v", err)
	}

	// IKE_AUTH
	ikeMessage = ike_message.BuildIKEHeader(123123, ikeSecurityAssociation.RemoteSPI, ike_message.IKE_AUTH, ike_message.InitiatorBitCheck, 1)

	var ikePayload []ike_message.IKEPayloadType

	// Identification
	identification := ike_message.BuildIdentificationInitiator(ike_message.ID_FQDN, []byte("UE"))
	ikePayload = append(ikePayload, identification)

	// Security Association
	proposal = ike_message.BuildProposal(1, ike_message.TypeESP, []byte{0, 0, 0, 1})
	encryptTransform = ike_message.BuildTransform(ike_message.TypeEncryptionAlgorithm, ike_message.ENCR_AES_CBC, &attributeType, &keyLength, nil)
	ike_message.AppendTransformToProposal(proposal, encryptTransform)
	integrityTransform = ike_message.BuildTransform(ike_message.TypeIntegrityAlgorithm, ike_message.AUTH_HMAC_SHA1_96, nil, nil, nil)
	ike_message.AppendTransformToProposal(proposal, integrityTransform)
	extendedSequenceNumbersTransform := ike_message.BuildTransform(ike_message.TypeExtendedSequenceNumbers, ike_message.ESN_NO, nil, nil, nil)
	ike_message.AppendTransformToProposal(proposal, extendedSequenceNumbersTransform)
	securityAssociation = ike_message.BuildSecurityAssociation([]*ike_message.Proposal{proposal})
	ikePayload = append(ikePayload, securityAssociation)

	// Traffic Selector
	inidividualTrafficSelector := ike_message.BuildIndividualTrafficSelector(ike_message.TS_IPV4_ADDR_RANGE, 0, 0, 65535, []byte{0, 0, 0, 0}, []byte{255, 255, 255, 255})
	trafficSelectorInitiator := ike_message.BuildTrafficSelectorInitiator([]*ike_message.IndividualTrafficSelector{inidividualTrafficSelector})
	ikePayload = append(ikePayload, trafficSelectorInitiator)
	trafficSelectorResponder := ike_message.BuildTrafficSelectorResponder([]*ike_message.IndividualTrafficSelector{inidividualTrafficSelector})
	ikePayload = append(ikePayload, trafficSelectorResponder)

	if err := encryptProcedure(ikeSecurityAssociation, ikePayload, ikeMessage); err != nil {
		t.Fatalf("Encrypting IKE message failed: %+v", err)
	}

	// Send to N3IWF
	ikeMessageData, err = ike_message.Encode(ikeMessage)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := udpConnection.WriteToUDP(ikeMessageData, n3iwfUDPAddr); err != nil {
		t.Fatal(err)
	}

	// Receive N3IWF reply
	n, _, err = udpConnection.ReadFromUDP(buffer)
	if err != nil {
		t.Fatal(err)
	}
	ikeMessage, err = ike_message.Decode(buffer[:n])
	if err != nil {
		t.Fatal(err)
	}

	encryptedPayload, ok := ikeMessage.IKEPayload[0].(*ike_message.Encrypted)
	if !ok {
		t.Fatal("Received payload is not an encrypted payload")
	}

	decryptedIKEPayload, err := decryptProcedure(ikeSecurityAssociation, ikeMessage, encryptedPayload)
	if err != nil {
		t.Fatalf("[IKE] Decrypt IKE message failed: %+v", err)
	}

	var eapIdentifier uint8

	for _, ikePayload := range decryptedIKEPayload {
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
	ikeMessage = ike_message.BuildIKEHeader(123123, ikeSecurityAssociation.RemoteSPI, ike_message.IKE_AUTH, ike_message.InitiatorBitCheck, 2)

	ikePayload = []ike_message.IKEPayloadType{}

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
	registrationRequest := nasTestpacket.GetRegistrationRequestWith5GMM(nasMessage.RegistrationType5GSInitialRegistration, mobileIdentity5GS, nil, nil)
	nasLength := make([]byte, 2)
	binary.BigEndian.PutUint16(nasLength, uint16(len(registrationRequest)))
	eapVendorTypeData = append(eapVendorTypeData, nasLength...)
	eapVendorTypeData = append(eapVendorTypeData, registrationRequest...)

	eapExpanded := ike_message.BuildEAPExpanded(ike_message.VendorID3GPP, ike_message.VendorTypeEAP5G, eapVendorTypeData)
	eap := ike_message.BuildEAP(ike_message.EAPCodeResponse, eapIdentifier, eapExpanded)

	ikePayload = append(ikePayload, eap)

	if err := encryptProcedure(ikeSecurityAssociation, ikePayload, ikeMessage); err != nil {
		t.Fatal(err)
	}

	// Send to N3IWF
	ikeMessageData, err = ike_message.Encode(ikeMessage)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := udpConnection.WriteToUDP(ikeMessageData, n3iwfUDPAddr); err != nil {
		t.Fatal(err)
	}

	// Receive N3IWF reply
	n, _, err = udpConnection.ReadFromUDP(buffer)
	if err != nil {
		t.Fatal(err)
	}
	ikeMessage, err = ike_message.Decode(buffer[:n])
	if err != nil {
		t.Fatal(err)
	}
	encryptedPayload, ok = ikeMessage.IKEPayload[0].(*ike_message.Encrypted)
	if !ok {
		t.Fatal("Received payload is not an encrypted payload")
	}
	decryptedIKEPayload, err = decryptProcedure(ikeSecurityAssociation, ikeMessage, encryptedPayload)
	if err != nil {
		t.Fatalf("[IKE] Decrypt IKE message failed: %+v", err)
	}

	var eapReq *ike_message.EAP

	eapReq, ok = decryptedIKEPayload[0].(*ike_message.EAP)
	if !ok {
		t.Fatal("Received packet is not an EAP payload")
	}

	var decodedNAS *nas.Message

	eapExpanded, ok = eapReq.EAPTypeData[0].(*ike_message.EAPExpanded)
	if !ok {
		t.Fatal("The EAP data is not an EAP expended.")
	}

	// Decode NAS - Authentication Request
	nasData := eapExpanded.VendorData[4:]
	decodedNAS = new(nas.Message)
	if err := decodedNAS.PlainNasDecode(&nasData); err != nil {
		t.Fatal(err)
	}

	// Calculate for RES*
	assert.NotNil(t, decodedNAS)
	rand := decodedNAS.AuthenticationRequest.GetRANDValue()
	resStat := ue.DeriveRESstarAndSetKey(ue.AuthenticationSubs, rand[:], "5G:mnc093.mcc208.3gppnetwork.org")

	// send NAS Authentication Response
	pdu := nasTestpacket.GetAuthenticationResponse(resStat, "")

	// IKE_AUTH - EAP exchange
	ikeMessage = ike_message.BuildIKEHeader(123123, ikeSecurityAssociation.RemoteSPI, ike_message.IKE_AUTH, ike_message.InitiatorBitCheck, 3)

	ikePayload = []ike_message.IKEPayloadType{}

	// EAP-5G vendor type data
	eapVendorTypeData = make([]byte, 4)
	eapVendorTypeData[0] = ike_message.EAP5GType5GNAS

	// NAS - Authentication Response
	nasLength = make([]byte, 2)
	binary.BigEndian.PutUint16(nasLength, uint16(len(pdu)))
	eapVendorTypeData = append(eapVendorTypeData, nasLength...)
	eapVendorTypeData = append(eapVendorTypeData, pdu...)

	eapExpanded = ike_message.BuildEAPExpanded(ike_message.VendorID3GPP, ike_message.VendorTypeEAP5G, eapVendorTypeData)
	eap = ike_message.BuildEAP(ike_message.EAPCodeResponse, eapReq.Identifier, eapExpanded)

	ikePayload = append(ikePayload, eap)

	err = encryptProcedure(ikeSecurityAssociation, ikePayload, ikeMessage)
	if err != nil {
		t.Fatal(err)
	}

	// Send to N3IWF
	ikeMessageData, err = ike_message.Encode(ikeMessage)
	if err != nil {
		t.Fatal(err)
	}
	_, err = udpConnection.WriteToUDP(ikeMessageData, n3iwfUDPAddr)
	if err != nil {
		t.Fatal(err)
	}

	// Receive N3IWF reply
	n, _, err = udpConnection.ReadFromUDP(buffer)
	if err != nil {
		t.Fatal(err)
	}
	ikeMessage, err = ike_message.Decode(buffer[:n])
	if err != nil {
		t.Fatal(err)
	}
	encryptedPayload, ok = ikeMessage.IKEPayload[0].(*ike_message.Encrypted)
	if !ok {
		t.Fatal("Received pakcet is not and encrypted payload")
	}
	decryptedIKEPayload, err = decryptProcedure(ikeSecurityAssociation, ikeMessage, encryptedPayload)
	if err != nil {
		t.Fatal(err)
	}
	eapReq, ok = decryptedIKEPayload[0].(*ike_message.EAP)
	if !ok {
		t.Fatal("Received packet is not an EAP payload")
	}
	eapExpanded, ok = eapReq.EAPTypeData[0].(*ike_message.EAPExpanded)
	if !ok {
		t.Fatal("Received packet is not an EAP expended payload")
	}

	nasData = eapExpanded.VendorData[4:]

	// Send NAS Security Mode Complete Msg
	pdu = nasTestpacket.GetSecurityModeComplete(registrationRequest)
	pdu, err = EncodeNasPduWithSecurity(ue, pdu)
	assert.Nil(t, err)

	// IKE_AUTH - EAP exchange
	ikeMessage = ike_message.BuildIKEHeader(123123, ikeSecurityAssociation.RemoteSPI, ike_message.IKE_AUTH, ike_message.InitiatorBitCheck, 4)

	ikePayload = []ike_message.IKEPayloadType{}

	// EAP-5G vendor type data
	eapVendorTypeData = make([]byte, 4)
	eapVendorTypeData[0] = ike_message.EAP5GType5GNAS

	// NAS - Authentication Response
	nasLength = make([]byte, 2)
	binary.BigEndian.PutUint16(nasLength, uint16(len(pdu)))
	eapVendorTypeData = append(eapVendorTypeData, nasLength...)
	eapVendorTypeData = append(eapVendorTypeData, pdu...)

	eapExpanded = ike_message.BuildEAPExpanded(ike_message.VendorID3GPP, ike_message.VendorTypeEAP5G, eapVendorTypeData)
	eap = ike_message.BuildEAP(ike_message.EAPCodeResponse, eapReq.Identifier, eapExpanded)

	ikePayload = append(ikePayload, eap)

	err = encryptProcedure(ikeSecurityAssociation, ikePayload, ikeMessage)
	if err != nil {
		t.Fatal(err)
	}

	// Send to N3IWF
	ikeMessageData, err = ike_message.Encode(ikeMessage)
	if err != nil {
		t.Fatal(err)
	}
	_, err = udpConnection.WriteToUDP(ikeMessageData, n3iwfUDPAddr)
	if err != nil {
		t.Fatal(err)
	}

	// Receive N3IWF reply
	n, _, err = udpConnection.ReadFromUDP(buffer)
	if err != nil {
		t.Fatal(err)
	}
	ikeMessage, err = ike_message.Decode(buffer[:n])
	if err != nil {
		t.Fatal(err)
	}
	encryptedPayload, ok = ikeMessage.IKEPayload[0].(*ike_message.Encrypted)
	if !ok {
		t.Fatal("Received pakcet is not and encrypted payload")
	}
	decryptedIKEPayload, err = decryptProcedure(ikeSecurityAssociation, ikeMessage, encryptedPayload)
	if err != nil {
		t.Fatal(err)
	}
	eapReq, ok = decryptedIKEPayload[0].(*ike_message.EAP)
	if !ok {
		t.Fatal("Received packet is not an EAP payload")
	}
	if eapReq.Code != ike_message.EAPCodeSuccess {
		t.Fatal("Not Success")
	}

	// IKE_AUTH - Authentication
	ikeMessage = ike_message.BuildIKEHeader(123123, ikeSecurityAssociation.RemoteSPI, ike_message.IKE_AUTH, ike_message.InitiatorBitCheck, 5)

	ikePayload = []ike_message.IKEPayloadType{}

	// Authentication
	auth := ike_message.BuildAuthentication(ike_message.SharedKeyMesageIntegrityCode, []byte{1, 2, 3})
	ikePayload = append(ikePayload, auth)

	// Configuration Request
	configurationAttribute := ike_message.BuildConfigurationAttribute(ike_message.INTERNAL_IP4_ADDRESS, nil)
	configurationRequest := ike_message.BuildConfiguration(ike_message.CFG_REQUEST, []*ike_message.IndividualConfigurationAttribute{configurationAttribute})
	ikePayload = append(ikePayload, configurationRequest)

	err = encryptProcedure(ikeSecurityAssociation, ikePayload, ikeMessage)
	if err != nil {
		t.Fatal(err)
	}

	// Send to N3IWF
	ikeMessageData, err = ike_message.Encode(ikeMessage)
	if err != nil {
		t.Fatal(err)
	}
	_, err = udpConnection.WriteToUDP(ikeMessageData, n3iwfUDPAddr)
	if err != nil {
		t.Fatal(err)
	}

	// Receive N3IWF reply
	n, _, err = udpConnection.ReadFromUDP(buffer)
	if err != nil {
		t.Fatal(err)
	}
	ikeMessage, err = ike_message.Decode(buffer[:n])
	if err != nil {
		t.Fatal(err)
	}
	encryptedPayload, ok = ikeMessage.IKEPayload[0].(*ike_message.Encrypted)
	if !ok {
		t.Fatal("Received pakcet is not and encrypted payload")
	}
	decryptedIKEPayload, err = decryptProcedure(ikeSecurityAssociation, ikeMessage, encryptedPayload)
	if err != nil {
		t.Fatal(err)
	}

	// AUTH, SAr2, TSi, Tsr, N(NAS_IP_ADDRESS), N(NAS_TCP_PORT)
	var responseSecurityAssociation *ike_message.SecurityAssociation
	var responseTrafficSelectorInitiator *ike_message.TrafficSelectorInitiator
	var responseTrafficSelectorResponder *ike_message.TrafficSelectorResponder
	var responseConfiguration *ike_message.Configuration
	n3iwfNASAddr := new(net.TCPAddr)
	ueAddr := new(net.IPNet)

	for _, ikePayload := range decryptedIKEPayload {
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
						ueAddr.IP = configAttr.Value
					}
					if configAttr.Type == ike_message.INTERNAL_IP4_NETMASK {
						ueAddr.Mask = configAttr.Value
					}
				}
			}
		}
	}

	childSecurityAssociationContext, err := createIKEChildSecurityAssociation(ikeSecurityAssociation.IKEAuthResponseSA)
	if err != nil {
		t.Fatalf("[IKE] Create child security association context failed: %+v", err)
		return
	}
	err = parseIPAddressInformationToChildSecurityAssociation(childSecurityAssociationContext, net.ParseIP("192.168.127.1"), responseTrafficSelectorInitiator, responseTrafficSelectorResponder)
	if err != nil {
		t.Fatalf("[IKE] Parse IP address to child security association failed: %+v", err)
		return
	}
	// Select TCP traffic
	childSecurityAssociationContext.SelectedIPProtocol = unix.IPPROTO_TCP

	if err := generateKeyForChildSA(ikeSecurityAssociation, childSecurityAssociationContext); err != nil {
		t.Fatalf("[IKE] Generate key for child SA failed: %+v", err)
		return
	}

	// Aplly XFRM rules
	if err = applyXFRMRule(childSecurityAssociationContext); err != nil {
		t.Fatalf("[IKE] Applying XFRM rules failed: %+v", err)
		return
	}

	// Get link ipsec0
	links, err := netlink.LinkList()
	if err != nil {
		t.Fatal(err)
	}

	var linkIPSec netlink.Link
	for _, link := range links {
		if link.Attrs() != nil {
			if link.Attrs().Name == "ipsec0" {
				linkIPSec = link
				break
			}
		}
	}
	if linkIPSec == nil {
		t.Fatal("No link named ipsec0")
	}

	linkIPSecAddr := &netlink.Addr{
		IPNet: ueAddr,
	}

	if err := netlink.AddrAdd(linkIPSec, linkIPSecAddr); err != nil {
		t.Fatalf("Set ipsec0 addr failed: %v", err)
	}

	defer func() {
		_ = netlink.AddrDel(linkIPSec, linkIPSecAddr)
		_ = netlink.XfrmPolicyFlush()
		_ = netlink.XfrmStateFlush(netlink.XFRM_PROTO_IPSEC_ANY)
	}()

	localTCPAddr := &net.TCPAddr{
		IP: ueAddr.IP,
	}
	tcpConnWithN3IWF, err := net.DialTCP("tcp", localTCPAddr, n3iwfNASAddr)
	if err != nil {
		t.Fatal(err)
	}

	_, err = tcpConnWithN3IWF.Write([]byte("Dial"))
	if err != nil {
		t.Fatal(err)
	}
}
