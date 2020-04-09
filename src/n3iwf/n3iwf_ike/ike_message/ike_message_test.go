package ike_message

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	Crand "crypto/rand"
	"encoding/binary"
	"encoding/hex"
	"io"
	Mrand "math/rand"
	"net"
	"testing"
)

var conn net.Conn

func init() {
	conn, _ = net.Dial("udp", "127.0.0.1:500")
}

// TestEncodeDecode tests the Encode() and Decode() function using the data
// build manually.
// First, build each payload with correct value, then the IKE message for
// IKE_SA_INIT type.
// Second, encode/decode the IKE message using Encode/Decode function, and then
// re-encode the decoded message again.
// Third, send the encoded data to the UDP connection for verification with Wireshark.
// Compare the dataFirstEncode and dataSecondEncode and return the result.
func TestEncodeDecode(t *testing.T) {
	testPacket := &IKEMessage{}

	// random an SPI
	src := Mrand.NewSource(63579)
	localRand := Mrand.New(src)
	ispi := localRand.Uint64()

	testPacket.InitiatorSPI = ispi
	testPacket.Version = 0x20
	testPacket.ExchangeType = 34 // IKE_SA_INIT
	testPacket.Flags = 16        // flagI is set
	testPacket.MessageID = 0     // for IKE_SA_INIT

	testSA := &SecurityAssociation{}

	testProposal1 := &Proposal{}
	testProposal1.ProposalNumber = 1 // first
	testProposal1.ProtocolID = 1     // IKE

	testtransform1 := &Transform{}
	testtransform1.TransformType = 1 // ENCR
	testtransform1.TransformID = 12  // ENCR_AES_CBC
	testtransform1.AttributePresent = true
	testtransform1.AttributeFormat = 1
	testtransform1.AttributeType = 14
	testtransform1.AttributeValue = 128

	testProposal1.EncryptionAlgorithm = append(testProposal1.EncryptionAlgorithm, testtransform1)

	testtransform2 := &Transform{}
	testtransform2.TransformType = 1 // ENCR
	testtransform2.TransformID = 12  // ENCR_AES_CBC
	testtransform2.AttributePresent = true
	testtransform2.AttributeFormat = 1
	testtransform2.AttributeType = 14
	testtransform2.AttributeValue = 192

	testProposal1.EncryptionAlgorithm = append(testProposal1.EncryptionAlgorithm, testtransform2)

	testtransform3 := &Transform{}
	testtransform3.TransformType = 3 // INTEG
	testtransform3.TransformID = 5   // AUTH_AES_XCBC_96
	testtransform3.AttributePresent = false

	testProposal1.IntegrityAlgorithm = append(testProposal1.IntegrityAlgorithm, testtransform3)

	testtransform4 := &Transform{}
	testtransform4.TransformType = 3 // INTEG
	testtransform4.TransformID = 2   // AUTH_HMAC_SHA1_96
	testtransform4.AttributePresent = false

	testProposal1.IntegrityAlgorithm = append(testProposal1.IntegrityAlgorithm, testtransform4)

	testSA.Proposals = append(testSA.Proposals, testProposal1)

	testProposal2 := &Proposal{}
	testProposal2.ProposalNumber = 2 // second
	testProposal2.ProtocolID = 1     // IKE

	testtransform1 = &Transform{}
	testtransform1.TransformType = 1 // ENCR
	testtransform1.TransformID = 12  // ENCR_AES_CBC
	testtransform1.AttributePresent = true
	testtransform1.AttributeFormat = 1
	testtransform1.AttributeType = 14
	testtransform1.AttributeValue = 128

	testProposal2.EncryptionAlgorithm = append(testProposal2.EncryptionAlgorithm, testtransform1)

	testtransform2 = &Transform{}
	testtransform2.TransformType = 1 // ENCR
	testtransform2.TransformID = 12  // ENCR_AES_CBC
	testtransform2.AttributePresent = true
	testtransform2.AttributeFormat = 1
	testtransform2.AttributeType = 14
	testtransform2.AttributeValue = 192

	testProposal2.EncryptionAlgorithm = append(testProposal2.EncryptionAlgorithm, testtransform2)

	testtransform3 = &Transform{}
	testtransform3.TransformType = 3 // INTEG
	testtransform3.TransformID = 1   // AUTH_HMAC_MD5_96
	testtransform3.AttributePresent = false

	testProposal2.IntegrityAlgorithm = append(testProposal2.IntegrityAlgorithm, testtransform3)

	testtransform4 = &Transform{}
	testtransform4.TransformType = 3 // INTEG
	testtransform4.TransformID = 2   // AUTH_HMAC_SHA1_96
	testtransform4.AttributePresent = false

	testProposal2.IntegrityAlgorithm = append(testProposal2.IntegrityAlgorithm, testtransform4)

	testSA.Proposals = append(testSA.Proposals, testProposal2)

	testPacket.IKEPayload = append(testPacket.IKEPayload, testSA)

	testKE := &KeyExchange{}

	testKE.DiffieHellmanGroup = 1
	for i := 0; i < 8; i++ {
		partKeyExchangeData := make([]byte, 8)
		binary.BigEndian.PutUint64(partKeyExchangeData, 7482105748278537214)
		testKE.KeyExchangeData = append(testKE.KeyExchangeData, partKeyExchangeData...)
	}

	testPacket.IKEPayload = append(testPacket.IKEPayload, testKE)

	testIDr := &IdentificationResponder{}

	testIDr.IDType = 3
	for i := 0; i < 8; i++ {
		partIdentification := make([]byte, 8)
		binary.BigEndian.PutUint64(partIdentification, 4378215321473912643)
		testIDr.IDData = append(testIDr.IDData, partIdentification...)
	}

	testPacket.IKEPayload = append(testPacket.IKEPayload, testIDr)

	testCert := &Certificate{}

	testCert.CertificateEncoding = 1
	for i := 0; i < 8; i++ {
		partCertificate := make([]byte, 8)
		binary.BigEndian.PutUint64(partCertificate, 4378217432157543265)
		testCert.CertificateData = append(testCert.CertificateData, partCertificate...)
	}

	testPacket.IKEPayload = append(testPacket.IKEPayload, testCert)

	testCertReq := &CertificateRequest{}

	testCertReq.CertificateEncoding = 1
	for i := 0; i < 8; i++ {
		partCertificateRquest := make([]byte, 8)
		binary.BigEndian.PutUint64(partCertificateRquest, 7438274381754372584)
		testCertReq.CertificationAuthority = append(testCertReq.CertificationAuthority, partCertificateRquest...)
	}

	testPacket.IKEPayload = append(testPacket.IKEPayload, testCertReq)

	testAuth := &Authentication{}

	testAuth.AuthenticationMethod = 1
	for i := 0; i < 8; i++ {
		partAuthentication := make([]byte, 8)
		binary.BigEndian.PutUint64(partAuthentication, 4632714362816473824)
		testAuth.AuthenticationData = append(testAuth.AuthenticationData, partAuthentication...)
	}

	testPacket.IKEPayload = append(testPacket.IKEPayload, testAuth)

	testNonce := &Nonce{}

	for i := 0; i < 8; i++ {
		partNonce := make([]byte, 8)
		binary.BigEndian.PutUint64(partNonce, 8984327463782167381)
		testNonce.NonceData = append(testNonce.NonceData, partNonce...)
	}

	testPacket.IKEPayload = append(testPacket.IKEPayload, testNonce)

	testNotification := &Notification{}

	testNotification.ProtocolID = 1
	testNotification.NotifyMessageType = 2

	for i := 0; i < 5; i++ {
		partSPI := make([]byte, 8)
		binary.BigEndian.PutUint64(partSPI, 4372847328749832794)
		testNotification.SPI = append(testNotification.SPI, partSPI...)
	}

	for i := 0; i < 19; i++ {
		partNotification := make([]byte, 8)
		binary.BigEndian.PutUint64(partNotification, 9721437148392747354)
		testNotification.NotificationData = append(testNotification.NotificationData, partNotification...)
	}

	testPacket.IKEPayload = append(testPacket.IKEPayload, testNotification)

	testDelete := &Delete{}

	testDelete.ProtocolID = 1
	testDelete.SPISize = 9
	testDelete.NumberOfSPI = 4

	for i := 0; i < 36; i++ {
		testDelete.SPIs = append(testDelete.SPIs, 87)
	}

	testPacket.IKEPayload = append(testPacket.IKEPayload, testDelete)

	testVendor := &VendorID{}

	for i := 0; i < 5; i++ {
		partVendorData := make([]byte, 8)
		binary.BigEndian.PutUint64(partVendorData, 5421487329873941748)
		testVendor.VendorIDData = append(testVendor.VendorIDData, partVendorData...)
	}

	testPacket.IKEPayload = append(testPacket.IKEPayload, testVendor)

	testTSi := &TrafficSelectorResponder{}

	testIndividualTS := &IndividualTrafficSelector{}

	testIndividualTS.TSType = 7
	testIndividualTS.IPProtocolID = 6
	testIndividualTS.StartPort = 1989
	testIndividualTS.EndPort = 2020

	testIndividualTS.StartAddress = []byte{192, 168, 0, 15}
	testIndividualTS.EndAddress = []byte{192, 168, 0, 192}

	testTSi.TrafficSelectors = append(testTSi.TrafficSelectors, testIndividualTS)

	testIndividualTS = &IndividualTrafficSelector{}

	testIndividualTS.TSType = 8
	testIndividualTS.IPProtocolID = 6
	testIndividualTS.StartPort = 2010
	testIndividualTS.EndPort = 2050

	testIndividualTS.StartAddress = net.ParseIP("2001:db8::68")
	testIndividualTS.EndAddress = net.ParseIP("2001:db8::72")

	testTSi.TrafficSelectors = append(testTSi.TrafficSelectors, testIndividualTS)

	testPacket.IKEPayload = append(testPacket.IKEPayload, testTSi)

	testCP := new(Configuration)

	testCP.ConfigurationType = 1

	testIndividualConfigurationAttribute := new(IndividualConfigurationAttribute)

	testIndividualConfigurationAttribute.Type = 1
	testIndividualConfigurationAttribute.Value = []byte{10, 1, 14, 1}

	testCP.ConfigurationAttribute = append(testCP.ConfigurationAttribute, testIndividualConfigurationAttribute)

	testPacket.IKEPayload = append(testPacket.IKEPayload, testCP)

	testEAP := new(EAP)

	testEAP.Code = 1
	testEAP.Identifier = 123

	testEAPExpanded := new(EAPExpanded)

	testEAPExpanded.VendorID = 26838
	testEAPExpanded.VendorType = 1
	testEAPExpanded.VendorData = []byte{9, 4, 8, 7}

	testEAPNotification := new(EAPNotification)

	rawstr := "I'm tired"
	testEAPNotification.NotificationData = []byte(rawstr)

	testEAP.EAPTypeData = append(testEAP.EAPTypeData, testEAPNotification)

	testPacket.IKEPayload = append(testPacket.IKEPayload, testEAP)

	testSK := new(Encrypted)

	testSK.NextPayload = TypeSA

	ikePayload := []IKEPayloadType{
		testSA,
		testAuth,
	}

	ikePayloadDataForSK, retErr := EncodePayload(ikePayload)
	if retErr != nil {
		t.Fatalf("EncodePayload failed: %+v", retErr)
	}

	// aes 128 key
	key, retErr := hex.DecodeString("6368616e676520746869732070617373")
	if retErr != nil {
		t.Fatalf("HexDecoding failed: %+v", retErr)
	}
	block, retErr := aes.NewCipher(key)
	if retErr != nil {
		t.Fatalf("AES NewCipher failed: %+v", retErr)
	}

	// padding plaintext
	padNum := len(ikePayloadDataForSK) % aes.BlockSize
	for i := 0; i < (aes.BlockSize - padNum); i++ {
		ikePayloadDataForSK = append(ikePayloadDataForSK, byte(padNum))
	}

	// ciphertext
	cipherText := make([]byte, aes.BlockSize+len(ikePayloadDataForSK))
	iv := cipherText[:aes.BlockSize]
	if _, err := io.ReadFull(Crand.Reader, iv); err != nil {
		t.Fatalf("IO ReadFull failed: %+v", err)
	}

	// CBC mode
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(cipherText[aes.BlockSize:], ikePayloadDataForSK)

	testSK.EncryptedData = cipherText

	testPacket.IKEPayload = append(testPacket.IKEPayload, testSK)

	var dataFirstEncode, dataSecondEncode []byte
	var err error
	var decodedPacket *IKEMessage

	if dataFirstEncode, err = Encode(testPacket); err != nil {
		t.Fatalf("Encode failed: %+v", err)
	}

	t.Logf("%+v", dataFirstEncode)

	if decodedPacket, err = Decode(dataFirstEncode); err != nil {
		t.Fatalf("Decode failed: %+v", err)
	}

	if dataSecondEncode, err = Encode(decodedPacket); err != nil {
		t.Fatalf("Encode failed: %+v", err)
	}

	t.Logf("Original IKE Message: %+v", dataFirstEncode)
	t.Logf("Result IKE Message:   %+v", dataSecondEncode)

	_, err = conn.Write(dataFirstEncode)
	if err != nil {
		t.Fatalf("Error: %+v", err)
	}

	if !bytes.Equal(dataFirstEncode, dataSecondEncode) {
		t.FailNow()
	}

}

// TestEncodeDecodeUsingPublicData tests the Encode() and Decode() function
// using the public data.
// Decode and encode the data, and compare the verifyData and the origin
// data and return the result.
func TestEncodeDecodeUsingPublicData(t *testing.T) {
	data := []byte{
		0x86, 0x43, 0x30, 0xac, 0x30, 0xe6, 0x56, 0x4d, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x21, 0x20, 0x22, 0x08, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0xc9, 0x22, 0x00, 0x00, 0x30, 0x00, 0x00, 0x00, 0x2c, 0x01, 0x01, 0x00, 0x04, 0x03, 0x00,
		0x00, 0x0c, 0x01, 0x00, 0x00, 0x0c, 0x80, 0x0e, 0x00, 0x80, 0x03, 0x00, 0x00, 0x08, 0x02, 0x00, 0x00, 0x02, 0x03, 0x00, 0x00,
		0x08, 0x03, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00, 0x08, 0x04, 0x00, 0x00, 0x02, 0x28, 0x00, 0x00, 0x88, 0x00, 0x02, 0x00, 0x00,
		0x03, 0xdc, 0xf5, 0x9a, 0x29, 0x05, 0x7b, 0x5a, 0x49, 0xbd, 0x55, 0x8c, 0x9b, 0x14, 0x7a, 0x11, 0x0e, 0xed, 0xff, 0xe5, 0xea,
		0x2d, 0x12, 0xc2, 0x1e, 0x5c, 0x7a, 0x5f, 0x5e, 0x9c, 0x99, 0xe3, 0xd1, 0xd3, 0x00, 0x24, 0x3c, 0x89, 0x73, 0x1e, 0x6c, 0x6d,
		0x63, 0x41, 0x7b, 0x33, 0xfa, 0xaf, 0x5a, 0xc7, 0x26, 0xe8, 0xb6, 0xf8, 0xc3, 0xb5, 0x2a, 0x14, 0xeb, 0xec, 0xd5, 0x6f, 0x1b,
		0xd9, 0x5b, 0x28, 0x32, 0x84, 0x9e, 0x26, 0xfc, 0x59, 0xee, 0xf1, 0x4e, 0x38, 0x5f, 0x55, 0xc2, 0x1b, 0xe8, 0xf6, 0xa3, 0xfb,
		0xc5, 0x55, 0xd7, 0x35, 0x92, 0x86, 0x24, 0x00, 0x62, 0x8b, 0xea, 0xce, 0x23, 0xf0, 0x47, 0xaf, 0xaa, 0xf8, 0x61, 0xe4, 0x5c,
		0x42, 0xba, 0x5c, 0xa1, 0x4a, 0x52, 0x6e, 0xd8, 0xe8, 0xf1, 0xb9, 0x74, 0xae, 0xe4, 0xd1, 0x9c, 0x9f, 0xa5, 0x9b, 0xf0, 0xd7,
		0xdb, 0x55, 0x2b, 0x00, 0x00, 0x44, 0x4c, 0xa7, 0xf3, 0x9b, 0xcd, 0x1d, 0xc2, 0x01, 0x79, 0xfa, 0xa2, 0xe4, 0x72, 0xe0, 0x61,
		0xc4, 0x45, 0x61, 0xe6, 0x49, 0x2d, 0xb3, 0x96, 0xae, 0xc9, 0x2c, 0xdb, 0x54, 0x21, 0xf4, 0x98, 0x4f, 0x72, 0xd2, 0x43, 0x78,
		0xab, 0x80, 0xe4, 0x6c, 0x01, 0x78, 0x6a, 0xc4, 0x64, 0x45, 0xbc, 0xa8, 0x1f, 0x56, 0xbc, 0xed, 0xf9, 0xb5, 0xd8, 0x21, 0x95,
		0x41, 0x71, 0xe9, 0x0e, 0xb4, 0x3c, 0x4e, 0x2b, 0x00, 0x00, 0x17, 0x43, 0x49, 0x53, 0x43, 0x4f, 0x2d, 0x44, 0x45, 0x4c, 0x45,
		0x54, 0x45, 0x2d, 0x52, 0x45, 0x41, 0x53, 0x4f, 0x4e, 0x2b, 0x00, 0x00, 0x3b, 0x43, 0x49, 0x53, 0x43, 0x4f, 0x28, 0x43, 0x4f,
		0x50, 0x59, 0x52, 0x49, 0x47, 0x48, 0x54, 0x29, 0x26, 0x43, 0x6f, 0x70, 0x79, 0x72, 0x69, 0x67, 0x68, 0x74, 0x20, 0x28, 0x63,
		0x29, 0x20, 0x32, 0x30, 0x30, 0x39, 0x20, 0x43, 0x69, 0x73, 0x63, 0x6f, 0x20, 0x53, 0x79, 0x73, 0x74, 0x65, 0x6d, 0x73, 0x2c,
		0x20, 0x49, 0x6e, 0x63, 0x2e, 0x29, 0x00, 0x00, 0x13, 0x43, 0x49, 0x53, 0x43, 0x4f, 0x2d, 0x47, 0x52, 0x45, 0x2d, 0x4d, 0x4f,
		0x44, 0x45, 0x02, 0x29, 0x00, 0x00, 0x1c, 0x01, 0x00, 0x40, 0x04, 0x7e, 0x57, 0x6c, 0xc0, 0x13, 0xd4, 0x05, 0x43, 0xa2, 0xe8,
		0x77, 0x7d, 0x00, 0x34, 0x68, 0xa5, 0xb1, 0x89, 0x0c, 0x58, 0x2b, 0x00, 0x00, 0x1c, 0x01, 0x00, 0x40, 0x05, 0x52, 0x64, 0x4d,
		0x87, 0xd4, 0x7c, 0x2d, 0x44, 0x23, 0xbd, 0x37, 0xe4, 0x48, 0xa9, 0xf5, 0x17, 0x01, 0x81, 0xcb, 0x8a, 0x00, 0x00, 0x00, 0x14,
		0x40, 0x48, 0xb7, 0xd5, 0x6e, 0xbc, 0xe8, 0x85, 0x25, 0xe7, 0xde, 0x7f, 0x00, 0xd6, 0xc2, 0xd3}

	ikePacket, err := Decode(data)
	if err != nil {
		t.Fatalf("Decode failed: %+v", err)
	}

	verifyData, err := Encode(ikePacket)
	if err != nil {
		t.Fatalf("Encode failed: %+v", err)
	}

	if !bytes.Equal(data, verifyData) {
		t.FailNow()
	}
}
