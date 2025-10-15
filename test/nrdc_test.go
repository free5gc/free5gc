package test_test

import (
	"encoding/hex"
	"fmt"
	"net"
	"test"
	"test/consumerTestdata/UDM/TestGenAuthData"
	"test/nasTestpacket"
	"testing"
	"time"

	"github.com/calee0219/fatal"
	"github.com/free5gc/aper"
	"github.com/free5gc/nas"
	"github.com/free5gc/nas/nasMessage"
	"github.com/free5gc/nas/nasType"
	"github.com/free5gc/nas/security"
	"github.com/free5gc/ngap"
	"github.com/free5gc/ngap/ngapConvert"
	"github.com/free5gc/ngap/ngapType"
	"github.com/free5gc/openapi/models"
	"github.com/free5gc/sctp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
)

const (
	amfN2Addr  = "127.0.0.1"
	mranN2Addr = "127.0.0.1"
	sranN2Addr = "127.0.0.2"
	upfN3Addr  = "10.200.200.102"
	mranN3Addr = "10.200.200.1"
	sranN3Addr = "10.200.200.2"

	googleDNS    = "9.9.10.10"
	cloudFareDNS = "9.9.9.9"

	ueIp = "10.60.0.1"
	loIp = "10.60.0.101"

	amfPort    = 38412
	mranN2Port = 9487
	sranN2Port = 9488
	mupfN3Port = 2152
	supfN3Port = 2152
	mranN3Port = 2152
	sranN3Port = 2152

	servingPlmnId = "20893"

	mranULTeid = "00000002"
	sranULTeid = "00000003"
	mranDLTeid = "\x00\x00\x00\x01"
	sranDLTeid = "\x00\x00\x00\x02"

	ENABLE_DC_AT_PDU_SESSION_ESTABLISHMENT    = true
	UN_ENABLE_DC_AT_PDU_SESSION_ESTABLISHMENT = false

	ENABLE_DC_AT_PDU_SESSION_MODIFY_INDICATION  = true
	DISABLE_DC_AT_PDU_SESSION_MODIFY_INDICATION = false

	EXPECTED_ERROR    = true
	EXPECTED_NO_ERROR = false
)

func connectRANsToAMF(t *testing.T) (*sctp.SCTPConn, *sctp.SCTPConn) {
	// Master RAN connect to AMF
	MranConn, err := test.ConnectToAmf(amfN2Addr, mranN2Addr, amfPort, mranN2Port)
	if err != nil {
		t.Logf("Master RAN connect to AMF failed: %v", err)
		return nil, nil
	}
	assert.NotNil(t, MranConn)

	// Second RAN connect to AMF
	SranConn, err := test.ConnectToAmf(amfN2Addr, sranN2Addr, amfPort, sranN2Port)
	if err != nil {
		t.Logf("Secondary RAN connect to AMF failed: %v", err)
		if MranConn != nil {
			MranConn.Close()
		}
		return nil, nil
	}
	assert.NotNil(t, SranConn)

	return MranConn, SranConn
}

func connectRANsToUPF(t *testing.T) (*net.UDPConn, *net.UDPConn) {
	// Master RAN connect to UPF
	MupfConn, err := test.ConnectToUpf(mranN3Addr, upfN3Addr, mupfN3Port, mranN3Port)
	if err != nil {
		t.Errorf("Master RAN connect to UPF failed: %v", err)
		return nil, nil
	}
	assert.NotNil(t, MupfConn)

	// Second RAN connect to UPF
	SupfConn, err := test.ConnectToUpf(sranN3Addr, upfN3Addr, supfN3Port, sranN3Port)
	if err != nil {
		t.Errorf("Secondary RAN connect to UPF failed: %v", err)
		if MupfConn != nil {
			MupfConn.Close()
		}
		return nil, nil
	}
	assert.NotNil(t, SupfConn)

	return MupfConn, SupfConn
}

func nGsSetup(t *testing.T, MranConn *sctp.SCTPConn, SranConn *sctp.SCTPConn) {
	var n int
	var recvMsg = make([]byte, 2048)

	// send Master RAN NGSetupRequest Msg
	sendMsg, err := test.GetNGSetupRequest([]byte("\x00\x01\x02"), 24, "MasterRAN")
	assert.Nil(t, err)
	_, err = MranConn.Write(sendMsg)
	assert.Nil(t, err)

	// receive Master RAN NGSetupResponse Msg
	n, err = MranConn.Read(recvMsg)
	assert.Nil(t, err)
	ngapPdu, err := ngap.Decoder(recvMsg[:n])
	assert.Nil(t, err)
	assert.True(t, ngapPdu.Present == ngapType.NGAPPDUPresentSuccessfulOutcome && ngapPdu.SuccessfulOutcome.ProcedureCode.Value == ngapType.ProcedureCodeNGSetup, "No NGSetupResponse received.")

	// send Second RAN NGSetupRequest Msg
	sendMsg, err = test.GetNGSetupRequest([]byte("\x00\x03\x04"), 24, "SecondRAN")
	assert.Nil(t, err)
	_, err = SranConn.Write(sendMsg)
	assert.Nil(t, err)

	// receive Second RAN NGSetupResponse Msg
	n, err = SranConn.Read(recvMsg)
	assert.Nil(t, err)
	ngapPdu, err = ngap.Decoder(recvMsg[:n])
	assert.Nil(t, err)
	assert.True(t, ngapPdu.Present == ngapType.NGAPPDUPresentSuccessfulOutcome && ngapPdu.SuccessfulOutcome.ProcedureCode.Value == ngapType.ProcedureCodeNGSetup, "No NGSetupResponse received.")
}

func newUEAndInitialRegistration(t *testing.T, MranConn *sctp.SCTPConn) *test.RanUeContext {
	var n int
	var sendMsg []byte
	var recvMsg = make([]byte, 2048)
	var err error

	// New UE
	ue := test.NewRanUeContext("imsi-208930000007487", 1, security.AlgCiphering128NEA0, security.AlgIntegrity128NIA2,
		models.AccessType__3_GPP_ACCESS)
	ue.AmfUeNgapId = 1
	ue.AuthenticationSubs = test.GetAuthSubscription(TestGenAuthData.MilenageTestSet19.K,
		TestGenAuthData.MilenageTestSet19.OPC,
		TestGenAuthData.MilenageTestSet19.OP)

	// insert UE data to MongoDB
	test.DelUeFromMongoDB(t, ue, servingPlmnId)
	test.InsertUeToMongoDB(t, ue, servingPlmnId)

	// send InitialUeMessage(Registration Request)(imsi-208930000007487)
	mobileIdentity5GS := nasType.MobileIdentity5GS{
		Len:    13, // suci
		Buffer: []uint8{0x01, 0x02, 0xf8, 0x39, 0xf0, 0xff, 0x00, 0x00, 0x00, 0x00, 0x00, 0x47, 0x78},
	}

	ueSecurityCapability := ue.GetUESecurityCapability()
	registrationRequest := nasTestpacket.GetRegistrationRequest(
		nasMessage.RegistrationType5GSInitialRegistration, mobileIdentity5GS, nil, ueSecurityCapability, nil, nil, nil)
	sendMsg, err = test.GetInitialUEMessage(ue.RanUeNgapId, registrationRequest, "")
	assert.Nil(t, err)
	_, err = MranConn.Write(sendMsg)
	assert.Nil(t, err)

	// receive NAS Authentication Request Msg
	n, err = MranConn.Read(recvMsg)
	assert.Nil(t, err)
	ngapPdu, err := ngap.Decoder(recvMsg[:n])
	assert.Nil(t, err)
	assert.True(t, ngapPdu.Present == ngapType.NGAPPDUPresentInitiatingMessage && ngapPdu.InitiatingMessage.ProcedureCode.Value == ngapType.ProcedureCodeDownlinkNASTransport, "No NAS Authentication Request received.")

	// Calculate for RES*
	nasPdu := test.GetNasPdu(ue, ngapPdu.InitiatingMessage.Value.DownlinkNASTransport)
	require.NotNil(t, nasPdu)
	require.NotNil(t, nasPdu.GmmMessage, "GMM message is nil")
	require.Equal(t, nasPdu.GmmHeader.GetMessageType(), nas.MsgTypeAuthenticationRequest,
		"Received wrong GMM message. Expected Authentication Request.")
	rand := nasPdu.AuthenticationRequest.GetRANDValue()
	resStat := ue.DeriveRESstarAndSetKey(ue.AuthenticationSubs, rand[:], "5G:mnc093.mcc208.3gppnetwork.org")

	// send NAS Authentication Response
	pdu := nasTestpacket.GetAuthenticationResponse(resStat, "")
	sendMsg, err = test.GetUplinkNASTransport(ue.AmfUeNgapId, ue.RanUeNgapId, pdu)
	assert.Nil(t, err)
	_, err = MranConn.Write(sendMsg)
	assert.Nil(t, err)

	// receive NAS Security Mode Command Msg
	n, err = MranConn.Read(recvMsg)
	assert.Nil(t, err)
	ngapPdu, err = ngap.Decoder(recvMsg[:n])
	assert.Nil(t, err)
	assert.NotNil(t, ngapPdu)
	nasPdu = test.GetNasPdu(ue, ngapPdu.InitiatingMessage.Value.DownlinkNASTransport)
	require.NotNil(t, nasPdu)
	require.NotNil(t, nasPdu.GmmMessage, "GMM message is nil")
	require.Equal(t, nasPdu.GmmHeader.GetMessageType(), nas.MsgTypeSecurityModeCommand,
		"Received wrong GMM message. Expected Security Mode Command.")

	// send NAS Security Mode Complete Msg
	registrationRequestWith5GMM := nasTestpacket.GetRegistrationRequest(nasMessage.RegistrationType5GSInitialRegistration,
		mobileIdentity5GS, nil, ueSecurityCapability, ue.Get5GMMCapability(), nil, nil)
	pdu = nasTestpacket.GetSecurityModeComplete(registrationRequestWith5GMM)
	pdu, err = test.EncodeNasPduWithSecurity(ue, pdu, nas.SecurityHeaderTypeIntegrityProtectedAndCipheredWithNew5gNasSecurityContext, true, true)
	assert.Nil(t, err)
	sendMsg, err = test.GetUplinkNASTransport(ue.AmfUeNgapId, ue.RanUeNgapId, pdu)
	assert.Nil(t, err)
	_, err = MranConn.Write(sendMsg)
	assert.Nil(t, err)

	// receive ngap Initial Context Setup Request Msg
	n, err = MranConn.Read(recvMsg)
	assert.Nil(t, err)
	ngapPdu, err = ngap.Decoder(recvMsg[:n])
	assert.Nil(t, err)
	assert.True(t, ngapPdu.Present == ngapType.NGAPPDUPresentInitiatingMessage &&
		ngapPdu.InitiatingMessage.ProcedureCode.Value == ngapType.ProcedureCodeInitialContextSetup,
		"No InitialContextSetup received.")

	// send ngap Initial Context Setup Response Msg
	sendMsg, err = test.GetInitialContextSetupResponse(ue.AmfUeNgapId, ue.RanUeNgapId)
	assert.Nil(t, err)
	_, err = MranConn.Write(sendMsg)
	assert.Nil(t, err)

	// send NAS Registration Complete Msg
	pdu = nasTestpacket.GetRegistrationComplete(nil)
	pdu, err = test.EncodeNasPduWithSecurity(ue, pdu, nas.SecurityHeaderTypeIntegrityProtectedAndCiphered, true, false)
	assert.Nil(t, err)
	sendMsg, err = test.GetUplinkNASTransport(ue.AmfUeNgapId, ue.RanUeNgapId, pdu)
	assert.Nil(t, err)
	_, err = MranConn.Write(sendMsg)
	assert.Nil(t, err)

	// receive UE Configuration Update Command Msg
	recvUeConfigUpdateCmd(t, recvMsg, MranConn)

	time.Sleep(100 * time.Millisecond)

	return ue
}

func pduSessionEstablishment(t *testing.T, ue *test.RanUeContext, MranConn *sctp.SCTPConn, enableDC bool) {
	var n int
	var sendMsg []byte
	var recvMsg = make([]byte, 2048)
	var err error

	buildPDUSessionResourceSetupResponseTransferWithDC := func() ngapType.PDUSessionResourceSetupResponseTransfer {
		var data ngapType.PDUSessionResourceSetupResponseTransfer
		// QoS Flow per TNL Information
		qosFlowPerTNLInformation := &data.DLQosFlowPerTNLInformation
		qosFlowPerTNLInformation.UPTransportLayerInformation.Present = ngapType.UPTransportLayerInformationPresentGTPTunnel

		// UP Transport Layer Information in QoS Flow per TNL Information
		upTransportLayerInformation := &qosFlowPerTNLInformation.UPTransportLayerInformation
		upTransportLayerInformation.Present = ngapType.UPTransportLayerInformationPresentGTPTunnel
		upTransportLayerInformation.GTPTunnel = new(ngapType.GTPTunnel)
		upTransportLayerInformation.GTPTunnel.GTPTEID.Value = aper.OctetString(mranDLTeid)
		upTransportLayerInformation.GTPTunnel.TransportLayerAddress = ngapConvert.IPAddressToNgap(mranN3Addr, "")

		// Associated QoS Flow List in QoS Flow per TNL Information
		associatedQosFlowList := &qosFlowPerTNLInformation.AssociatedQosFlowList

		associatedQosFlowItem := ngapType.AssociatedQosFlowItem{}
		associatedQosFlowItem.QosFlowIdentifier.Value = 1
		associatedQosFlowList.List = append(associatedQosFlowList.List, associatedQosFlowItem)

		if enableDC {
			// DC QoS Flow per TNL Information
			DCQosFlowPerTNLInformationItem := ngapType.QosFlowPerTNLInformationItem{}
			DCQosFlowPerTNLInformationItem.QosFlowPerTNLInformation.UPTransportLayerInformation.Present = ngapType.UPTransportLayerInformationPresentGTPTunnel

			// DC Transport Layer Information in QoS Flow per TNL Information
			DCUpTransportLayerInformation := &DCQosFlowPerTNLInformationItem.QosFlowPerTNLInformation.UPTransportLayerInformation
			DCUpTransportLayerInformation.Present = ngapType.UPTransportLayerInformationPresentGTPTunnel
			DCUpTransportLayerInformation.GTPTunnel = new(ngapType.GTPTunnel)
			DCUpTransportLayerInformation.GTPTunnel.GTPTEID.Value = aper.OctetString(sranDLTeid)
			DCUpTransportLayerInformation.GTPTunnel.TransportLayerAddress = ngapConvert.IPAddressToNgap(sranN3Addr, "")

			// DC Associated QoS Flow List in QoS Flow per TNL Information
			DCAssociatedQosFlowList := &DCQosFlowPerTNLInformationItem.QosFlowPerTNLInformation.AssociatedQosFlowList
			DCAssociatedQosFlowItem := ngapType.AssociatedQosFlowItem{}
			DCAssociatedQosFlowItem.QosFlowIdentifier.Value = 1
			DCAssociatedQosFlowList.List = append(DCAssociatedQosFlowList.List, DCAssociatedQosFlowItem)

			// Additional DL QoS Flow per TNL Information
			data.AdditionalDLQosFlowPerTNLInformation = new(ngapType.QosFlowPerTNLInformationList)
			data.AdditionalDLQosFlowPerTNLInformation.List = append(data.AdditionalDLQosFlowPerTNLInformation.List, DCQosFlowPerTNLInformationItem)
		}

		return data
	}

	getPDUSessionResourceSetupResponseTransferWithDC := func() []byte {
		data := buildPDUSessionResourceSetupResponseTransferWithDC()
		encodeData, err := aper.MarshalWithParams(data, "valueExt")
		if err != nil {
			fatal.Fatalf("aper MarshalWithParams error in GetPDUSessionResourceSetupResponseTransfer: %+v", err)
		}
		return encodeData
	}

	buildPDUSessionResourceSetupResponseForRegistrationTestWithDC := func(pduSessionId int64, amfUeNgapId int64, ranUeNgapId int64) (pdu ngapType.NGAPPDU) {
		pdu.Present = ngapType.NGAPPDUPresentSuccessfulOutcome
		pdu.SuccessfulOutcome = new(ngapType.SuccessfulOutcome)

		successfulOutcome := pdu.SuccessfulOutcome
		successfulOutcome.ProcedureCode.Value = ngapType.ProcedureCodePDUSessionResourceSetup
		successfulOutcome.Criticality.Value = ngapType.CriticalityPresentReject

		successfulOutcome.Value.Present = ngapType.SuccessfulOutcomePresentPDUSessionResourceSetupResponse
		successfulOutcome.Value.PDUSessionResourceSetupResponse = new(ngapType.PDUSessionResourceSetupResponse)

		pDUSessionResourceSetupResponse := successfulOutcome.Value.PDUSessionResourceSetupResponse
		pDUSessionResourceSetupResponseIEs := &pDUSessionResourceSetupResponse.ProtocolIEs

		// AMF UE NGAP ID
		ie := ngapType.PDUSessionResourceSetupResponseIEs{}
		ie.Id.Value = ngapType.ProtocolIEIDAMFUENGAPID
		ie.Criticality.Value = ngapType.CriticalityPresentIgnore
		ie.Value.Present = ngapType.PDUSessionResourceSetupResponseIEsPresentAMFUENGAPID
		ie.Value.AMFUENGAPID = new(ngapType.AMFUENGAPID)

		aMFUENGAPID := ie.Value.AMFUENGAPID
		aMFUENGAPID.Value = amfUeNgapId

		pDUSessionResourceSetupResponseIEs.List = append(pDUSessionResourceSetupResponseIEs.List, ie)

		// RAN UE NGAP ID
		ie = ngapType.PDUSessionResourceSetupResponseIEs{}
		ie.Id.Value = ngapType.ProtocolIEIDRANUENGAPID
		ie.Criticality.Value = ngapType.CriticalityPresentIgnore
		ie.Value.Present = ngapType.PDUSessionResourceSetupResponseIEsPresentRANUENGAPID
		ie.Value.RANUENGAPID = new(ngapType.RANUENGAPID)

		rANUENGAPID := ie.Value.RANUENGAPID
		rANUENGAPID.Value = ranUeNgapId

		pDUSessionResourceSetupResponseIEs.List = append(pDUSessionResourceSetupResponseIEs.List, ie)

		// PDU Session Resource Setup Response List
		ie = ngapType.PDUSessionResourceSetupResponseIEs{}
		ie.Id.Value = ngapType.ProtocolIEIDPDUSessionResourceSetupListSURes
		ie.Criticality.Value = ngapType.CriticalityPresentIgnore
		ie.Value.Present = ngapType.PDUSessionResourceSetupResponseIEsPresentPDUSessionResourceSetupListSURes
		ie.Value.PDUSessionResourceSetupListSURes = new(ngapType.PDUSessionResourceSetupListSURes)

		pDUSessionResourceSetupListSURes := ie.Value.PDUSessionResourceSetupListSURes

		// PDU Session Resource Setup Response Item in PDU Session Resource Setup Response List
		pDUSessionResourceSetupItemSURes := ngapType.PDUSessionResourceSetupItemSURes{}
		pDUSessionResourceSetupItemSURes.PDUSessionID.Value = pduSessionId

		pDUSessionResourceSetupItemSURes.PDUSessionResourceSetupResponseTransfer =
			getPDUSessionResourceSetupResponseTransferWithDC()

		pDUSessionResourceSetupListSURes.List = append(pDUSessionResourceSetupListSURes.List, pDUSessionResourceSetupItemSURes)

		pDUSessionResourceSetupResponseIEs.List = append(pDUSessionResourceSetupResponseIEs.List, ie)

		return pdu
	}

	getPDUSessionResourceSetupResponseWithDC := func(pduSessionId int64, amfUeNgapId int64, ranUeNgapId int64) ([]byte, error) {
		message := buildPDUSessionResourceSetupResponseForRegistrationTestWithDC(pduSessionId, amfUeNgapId, ranUeNgapId)
		return ngap.Encoder(message)
	}

	// send GetPduSessionEstablishmentRequest Msg
	sNssai := models.Snssai{
		Sst: 1,
		Sd:  "fedcba",
	}
	pdu := nasTestpacket.GetUlNasTransport_PduSessionEstablishmentRequest(10, nasMessage.ULNASTransportRequestTypeInitialRequest, "internet", &sNssai)
	pdu, err = test.EncodeNasPduWithSecurity(ue, pdu, nas.SecurityHeaderTypeIntegrityProtectedAndCiphered, true, false)
	assert.Nil(t, err)
	sendMsg, err = test.GetUplinkNASTransport(ue.AmfUeNgapId, ue.RanUeNgapId, pdu)
	assert.Nil(t, err)
	_, err = MranConn.Write(sendMsg)
	assert.Nil(t, err)

	// receive ngap PDU Session Resource Setup Request Msg
	n, err = MranConn.Read(recvMsg)
	assert.Nil(t, err)
	ngapPdu, err := ngap.Decoder(recvMsg[:n])
	assert.Nil(t, err)
	assert.True(t, ngapPdu.Present == ngapType.NGAPPDUPresentInitiatingMessage &&
		ngapPdu.InitiatingMessage.ProcedureCode.Value == ngapType.ProcedureCodePDUSessionResourceSetup,
		"No PDU Session Resource Setup Request received.")

	// send ngap PDU Session Resource Setup Response Msg
	sendMsg, err = getPDUSessionResourceSetupResponseWithDC(10, ue.AmfUeNgapId, ue.RanUeNgapId)
	assert.Nil(t, err)
	_, err = MranConn.Write(sendMsg)
	assert.Nil(t, err)

	time.Sleep(1 * time.Second)
}

func icmpTest(t *testing.T, upfConn *net.UDPConn, ulTeid, destIp string, expectedError bool) {
	gtpHdr, err := hex.DecodeString(fmt.Sprintf("32ff0034%s00000000", ulTeid))
	assert.Nil(t, err)
	icmpData, err := hex.DecodeString("8c870d0000000000101112131415161718191a1b1c1d1e1f202122232425262728292a2b2c2d2e2f3031323334353637")
	assert.Nil(t, err)

	ipv4hdr := ipv4.Header{
		Version:  4,
		Len:      20,
		Protocol: 1,
		Flags:    0,
		TotalLen: 48,
		TTL:      64,
		Src:      net.ParseIP(ueIp).To4(),
		Dst:      net.ParseIP(destIp).To4(),
		ID:       1,
	}
	checksum := test.CalculateIpv4HeaderChecksum(&ipv4hdr)
	ipv4hdr.Checksum = int(checksum)

	v4HdrBuf, err := ipv4hdr.Marshal()
	assert.Nil(t, err)
	tt := append(gtpHdr, v4HdrBuf...)

	m := icmp.Message{
		Type: ipv4.ICMPTypeEcho, Code: 0,
		Body: &icmp.Echo{
			ID: 12394, Seq: 1,
			Data: icmpData,
		},
	}
	b, err := m.Marshal(nil)
	assert.Nil(t, err)
	b[2] = 0xaf
	b[3] = 0x88
	_, err = upfConn.Write(append(tt, b...))
	assert.Nil(t, err)

	recvMsg := make([]byte, 2048)
	err = upfConn.SetReadDeadline(time.Now().Add(2 * time.Second))
	assert.Nil(t, err)
	n, err := upfConn.Read(recvMsg)
	if expectedError {
		assert.NotNil(t, err)
	} else {
		assert.Nil(t, err)
		assert.Equal(t, 64, n)
	}
	err = upfConn.SetReadDeadline(time.Time{})
	assert.Nil(t, err)
}

func pduSessionModifyIndication(t *testing.T, ue *test.RanUeContext, MranConn *sctp.SCTPConn, enableDC bool) {
	var n int
	var sendMsg []byte
	var recvMsg = make([]byte, 2048)
	var err error

	buildPDUSessionResourceModifyIndicationTransferWithDC := func() ngapType.PDUSessionResourceModifyIndicationTransfer {
		var data ngapType.PDUSessionResourceModifyIndicationTransfer
		// QoS Flow per TNL Information
		qosFlowPerTNLInformation := &data.DLQosFlowPerTNLInformation
		qosFlowPerTNLInformation.UPTransportLayerInformation.Present = ngapType.UPTransportLayerInformationPresentGTPTunnel

		// UP Transport Layer Information in QoS Flow per TNL Information
		upTransportLayerInformation := &qosFlowPerTNLInformation.UPTransportLayerInformation
		upTransportLayerInformation.Present = ngapType.UPTransportLayerInformationPresentGTPTunnel
		upTransportLayerInformation.GTPTunnel = new(ngapType.GTPTunnel)
		upTransportLayerInformation.GTPTunnel.GTPTEID.Value = aper.OctetString(mranDLTeid)
		upTransportLayerInformation.GTPTunnel.TransportLayerAddress = ngapConvert.IPAddressToNgap(mranN3Addr, "")

		// Associated QoS Flow List in QoS Flow per TNL Information
		associatedQosFlowList := &qosFlowPerTNLInformation.AssociatedQosFlowList

		associatedQosFlowItem := ngapType.AssociatedQosFlowItem{}
		associatedQosFlowItem.QosFlowIdentifier.Value = 1
		associatedQosFlowList.List = append(associatedQosFlowList.List, associatedQosFlowItem)

		if enableDC {
			// DC QoS Flow per TNL Information
			DCQosFlowPerTNLInformationItem := ngapType.QosFlowPerTNLInformationItem{}
			DCQosFlowPerTNLInformationItem.QosFlowPerTNLInformation.UPTransportLayerInformation.Present = ngapType.UPTransportLayerInformationPresentGTPTunnel

			// DC Transport Layer Information in QoS Flow per TNL Information
			DCUpTransportLayerInformation := &DCQosFlowPerTNLInformationItem.QosFlowPerTNLInformation.UPTransportLayerInformation
			DCUpTransportLayerInformation.Present = ngapType.UPTransportLayerInformationPresentGTPTunnel
			DCUpTransportLayerInformation.GTPTunnel = new(ngapType.GTPTunnel)
			DCUpTransportLayerInformation.GTPTunnel.GTPTEID.Value = aper.OctetString(sranDLTeid)
			DCUpTransportLayerInformation.GTPTunnel.TransportLayerAddress = ngapConvert.IPAddressToNgap(sranN3Addr, "")

			// DC Associated QoS Flow List in QoS Flow per TNL Information
			DCAssociatedQosFlowList := &DCQosFlowPerTNLInformationItem.QosFlowPerTNLInformation.AssociatedQosFlowList
			DCAssociatedQosFlowItem := ngapType.AssociatedQosFlowItem{}
			DCAssociatedQosFlowItem.QosFlowIdentifier.Value = 1
			DCAssociatedQosFlowList.List = append(DCAssociatedQosFlowList.List, DCAssociatedQosFlowItem)

			// Additional DL QoS Flow per TNL Information
			data.AdditionalDLQosFlowPerTNLInformation = new(ngapType.QosFlowPerTNLInformationList)
			data.AdditionalDLQosFlowPerTNLInformation.List = append(data.AdditionalDLQosFlowPerTNLInformation.List, DCQosFlowPerTNLInformationItem)
		}

		return data
	}

	getPDUSessionResourceModifyIndicationTransferWithDC := func() []byte {
		data := buildPDUSessionResourceModifyIndicationTransferWithDC()
		encodeData, err := aper.MarshalWithParams(data, "valueExt")
		if err != nil {
			fatal.Fatalf("aper MarshalWithParams error in GetPDUSessionResourceModifyIndicationTransfer: %+v", err)
		}
		return encodeData
	}

	buildPDUSessionResourceModifyIndication := func(pduSessionId int64, amfUeNgapId int64, ranUeNgapId int64) (pdu ngapType.NGAPPDU) {
		pdu.Present = ngapType.NGAPPDUPresentInitiatingMessage
		pdu.InitiatingMessage = new(ngapType.InitiatingMessage)

		initiatingMessage := pdu.InitiatingMessage
		initiatingMessage.ProcedureCode.Value = ngapType.ProcedureCodePDUSessionResourceModifyIndication
		initiatingMessage.Criticality.Value = ngapType.CriticalityPresentReject

		initiatingMessage.Value.Present = ngapType.InitiatingMessagePresentPDUSessionResourceModifyIndication
		initiatingMessage.Value.PDUSessionResourceModifyIndication = new(ngapType.PDUSessionResourceModifyIndication)

		indication := initiatingMessage.Value.PDUSessionResourceModifyIndication
		indicationIEs := &indication.ProtocolIEs

		// AMF UE NGAP ID
		ie := ngapType.PDUSessionResourceModifyIndicationIEs{}
		ie.Id.Value = ngapType.ProtocolIEIDAMFUENGAPID
		ie.Criticality.Value = ngapType.CriticalityPresentIgnore
		ie.Value.Present = ngapType.PDUSessionResourceModifyIndicationIEsPresentAMFUENGAPID
		ie.Value.AMFUENGAPID = new(ngapType.AMFUENGAPID)
		ie.Value.AMFUENGAPID.Value = amfUeNgapId
		indicationIEs.List = append(indicationIEs.List, ie)

		// RAN UE NGAP ID
		ie = ngapType.PDUSessionResourceModifyIndicationIEs{}
		ie.Id.Value = ngapType.ProtocolIEIDRANUENGAPID
		ie.Criticality.Value = ngapType.CriticalityPresentIgnore
		ie.Value.Present = ngapType.PDUSessionResourceModifyIndicationIEsPresentRANUENGAPID
		ie.Value.RANUENGAPID = new(ngapType.RANUENGAPID)
		ie.Value.RANUENGAPID.Value = ranUeNgapId
		indicationIEs.List = append(indicationIEs.List, ie)

		// PDU Session Resource Modify List
		ie = ngapType.PDUSessionResourceModifyIndicationIEs{}
		ie.Id.Value = ngapType.ProtocolIEIDPDUSessionResourceModifyListModInd
		ie.Criticality.Value = ngapType.CriticalityPresentReject
		ie.Value.Present = ngapType.PDUSessionResourceModifyIndicationIEsPresentPDUSessionResourceModifyListModInd
		ie.Value.PDUSessionResourceModifyListModInd = new(ngapType.PDUSessionResourceModifyListModInd)

		modifyItem := ngapType.PDUSessionResourceModifyItemModInd{}
		modifyItem.PDUSessionID.Value = pduSessionId
		modifyItem.PDUSessionResourceModifyIndicationTransfer = getPDUSessionResourceModifyIndicationTransferWithDC()

		ie.Value.PDUSessionResourceModifyListModInd.List = append(
			ie.Value.PDUSessionResourceModifyListModInd.List, modifyItem)

		indicationIEs.List = append(indicationIEs.List, ie)

		return pdu
	}

	getPDUSessionResourceModifyIndication := func(pduSessionId int64, amfUeNgapId int64, ranUeNgapId int64) ([]byte, error) {
		message := buildPDUSessionResourceModifyIndication(pduSessionId, amfUeNgapId, ranUeNgapId)
		return ngap.Encoder(message)
	}

	// send pdu session resource modify indication
	sendMsg, err = getPDUSessionResourceModifyIndication(10, ue.AmfUeNgapId, ue.RanUeNgapId)
	assert.Nil(t, err)
	_, err = MranConn.Write(sendMsg)
	assert.Nil(t, err)

	// receive pdu session resource modify confirm
	n, err = MranConn.Read(recvMsg)
	assert.Nil(t, err)
	ngapPdu, err := ngap.Decoder(recvMsg[:n])
	assert.Nil(t, err)
	assert.True(t, ngapPdu.Present == ngapType.NGAPPDUPresentSuccessfulOutcome &&
		ngapPdu.SuccessfulOutcome.ProcedureCode.Value == ngapType.ProcedureCodePDUSessionResourceModifyIndication,
		"No PDU Session Resource Modify Confirm received.")

	// check successful outcome
	for _, ie := range ngapPdu.SuccessfulOutcome.Value.PDUSessionResourceModifyConfirm.ProtocolIEs.List {
		switch ie.Id.Value {
		case ngapType.ProtocolIEIDAMFUENGAPID:
		case ngapType.ProtocolIEIDRANUENGAPID:
		case ngapType.ProtocolIEIDPDUSessionResourceModifyListModCfm:
			t.Log("PDU session modify indication request successful")
		case ngapType.ProtocolIEIDPDUSessionResourceFailedToModifyListModCfm:
			t.Fatalf("PDU session modify indication request failed")
		}
	}

	time.Sleep(1 * time.Second)
}

func TestDC(t *testing.T) {
	// RANs connect to AMF
	MranConn, SranConn := connectRANsToAMF(t)
	if MranConn == nil || SranConn == nil {
		t.Fatal("Failed to connect to AMF")
		return
	}
	defer MranConn.Close()
	defer SranConn.Close()
	t.Log("Master RAN and Secondary RAN connect to AMF successfully")

	// RANs connect to UPF
	MupfConn, SupfConn := connectRANsToUPF(t)
	if MupfConn == nil || SupfConn == nil {
		t.Fatal("Failed to connect to UPF")
		return
	}
	defer MupfConn.Close()
	defer SupfConn.Close()
	t.Log("Master RAN and Secondary RAN connect to UPF successfully")

	// NGSetup
	nGsSetup(t, MranConn, SranConn)
	t.Log("Master RAN and Secondary RAN NGSetup successfully")

	// New UE and initial registration(NAS/NGAP)
	ue := newUEAndInitialRegistration(t, MranConn)
	defer test.DelUeFromMongoDB(t, ue, servingPlmnId)
	t.Log("New UE and initial registration(NAS/NGAP) successfully")

	// PDU Session Establishment
	pduSessionEstablishment(t, ue, MranConn, ENABLE_DC_AT_PDU_SESSION_ESTABLISHMENT)
	t.Log("PDU Session Establishment successfully")

	// ping test via master RAN
	t.Run("ping test via master RAN", func(t *testing.T) {
		icmpTest(t, MupfConn, mranULTeid, googleDNS, EXPECTED_NO_ERROR)
		t.Log("ICMP test via master RAN successfully")
	})

	// ping test via Secondary RAN
	t.Run("ping test via Secondary RAN", func(t *testing.T) {
		icmpTest(t, SupfConn, sranULTeid, cloudFareDNS, EXPECTED_NO_ERROR)
		t.Log("ICMP test via Secondary RAN successfully")
	})

	NfTerminate()
}

func TestDynamicDC(t *testing.T) {
	// RANs connect to AMF
	MranConn, SranConn := connectRANsToAMF(t)
	if MranConn == nil || SranConn == nil {
		t.Fatal("Failed to connect to AMF")
		return
	}
	defer MranConn.Close()
	defer SranConn.Close()
	t.Log("Master RAN and Secondary RAN connect to AMF successfully")

	// RANs connect to UPF
	MupfConn, SupfConn := connectRANsToUPF(t)
	if MupfConn == nil || SupfConn == nil {
		t.Fatal("Failed to connect to UPF")
		return
	}
	defer MupfConn.Close()
	defer SupfConn.Close()
	t.Log("Master RAN and Secondary RAN connect to UPF successfully")

	// NGSetup
	nGsSetup(t, MranConn, SranConn)
	t.Log("Master RAN and Secondary RAN NGSetup successfully")

	// New UE and initial registration(NAS/NGAP)
	ue := newUEAndInitialRegistration(t, MranConn)
	defer test.DelUeFromMongoDB(t, ue, servingPlmnId)
	t.Log("New UE and initial registration(NAS/NGAP) successfully")

	// PDU Session Establishment
	pduSessionEstablishment(t, ue, MranConn, UN_ENABLE_DC_AT_PDU_SESSION_ESTABLISHMENT)
	t.Log("PDU Session Establishment successfully")

	// ICMP test before DC is enabled
	t.Run("ping test before DC is enabled", func(t *testing.T) {
		t.Run("ping test via master RAN", func(t *testing.T) {
			icmpTest(t, MupfConn, mranULTeid, googleDNS, EXPECTED_NO_ERROR)
			icmpTest(t, MupfConn, mranULTeid, cloudFareDNS, EXPECTED_NO_ERROR)
		})

		t.Run("ping test via secondary RAN", func(t *testing.T) {
			icmpTest(t, SupfConn, sranULTeid, cloudFareDNS, EXPECTED_ERROR)
		})
	})

	// PDU Session Modify Indication Enable DC
	pduSessionModifyIndication(t, ue, MranConn, ENABLE_DC_AT_PDU_SESSION_MODIFY_INDICATION)
	t.Log("PDU Session Modify Indication successfully")

	// ICMP test after DC is enabled
	t.Run("ping test after DC is enabled", func(t *testing.T) {
		t.Run("ping test via master RAN", func(t *testing.T) {
			icmpTest(t, MupfConn, mranULTeid, googleDNS, EXPECTED_NO_ERROR)
		})

		t.Run("ping test via secondary RAN", func(t *testing.T) {
			icmpTest(t, SupfConn, sranULTeid, cloudFareDNS, EXPECTED_NO_ERROR)
		})
	})

	// PDU Session Modify Indication Disable DC
	pduSessionModifyIndication(t, ue, MranConn, DISABLE_DC_AT_PDU_SESSION_MODIFY_INDICATION)
	t.Log("PDU Session Modify Indication successfully")

	// ICMP test after DC is disabled
	t.Run("ping test after DC is disabled", func(t *testing.T) {
		t.Run("ping test via master RAN", func(t *testing.T) {
			icmpTest(t, MupfConn, mranULTeid, googleDNS, EXPECTED_NO_ERROR)
			icmpTest(t, MupfConn, mranULTeid, cloudFareDNS, EXPECTED_NO_ERROR)
		})

		t.Run("ping test via secondary RAN", func(t *testing.T) {
			icmpTest(t, SupfConn, sranULTeid, cloudFareDNS, EXPECTED_ERROR)
		})
	})
}

func TestDCHandover(t *testing.T) {
	// RANs connect to AMF
	MranConn, SranConn := connectRANsToAMF(t)
	if MranConn == nil || SranConn == nil {
		t.Fatal("Failed to connect to AMF")
		return
	}
	defer MranConn.Close()
	defer SranConn.Close()
	t.Log("Master RAN and Secondary RAN connect to AMF successfully")

	// RANs connect to UPF
	MupfConn, SupfConn := connectRANsToUPF(t)
	if MupfConn == nil || SupfConn == nil {
		t.Fatal("Failed to connect to UPF")
		return
	}
	defer MupfConn.Close()
	defer SupfConn.Close()
	t.Log("Master RAN and Secondary RAN connect to UPF successfully")

	// NGSetup
	nGsSetup(t, MranConn, SranConn)
	t.Log("Master RAN and Secondary RAN NGSetup successfully")

	// New UE and initial registration(NAS/NGAP)
	ue := newUEAndInitialRegistration(t, MranConn)
	defer test.DelUeFromMongoDB(t, ue, servingPlmnId)
	t.Log("New UE and initial registration(NAS/NGAP) successfully")

	// PDU Session Establishment
	pduSessionEstablishment(t, ue, MranConn, ENABLE_DC_AT_PDU_SESSION_ESTABLISHMENT)
	t.Log("PDU Session Establishment successfully")

	// ping test via master RAN
	t.Run("ping test via master RAN", func(t *testing.T) {
		icmpTest(t, MupfConn, mranULTeid, googleDNS, EXPECTED_NO_ERROR)
		t.Log("ICMP test via master RAN successfully")
	})

	// ping test via secondary RAN
	t.Run("ping test via secondary RAN", func(t *testing.T) {
		icmpTest(t, SupfConn, sranULTeid, cloudFareDNS, EXPECTED_NO_ERROR)
		t.Log("ICMP test via secondary RAN successfully")
	})
}
