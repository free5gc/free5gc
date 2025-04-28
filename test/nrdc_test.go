package test_test

import (
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
)

const (
	amfN2Addr  = "127.0.0.1"
	mranN2Addr = "127.0.0.1"
	sranN2Addr = "127.0.0.1"
	upfN3Addr  = "10.200.200.102"
	mranN3Addr = "10.200.200.1"
	sranN3Addr = "10.200.200.2"

	amfPort    = 38412
	mranN2Port = 9487
	sranN2Port = 9488
	mupfN3Port = 2152
	supfN3Port = 2153
	mranN3Port = 2152
	sranN3Port = 2153

	servingPlmnId = "20893"

	mranTeid = "\x00\x00\x00\x01"
	sranTeid = "\x00\x00\x00\x02"
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
		t.Logf("Slave RAN connect to AMF failed: %v", err)
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
		t.Errorf("Slave RAN connect to UPF failed: %v", err)
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
	ue := test.NewRanUeContext("imsi-2089300007487", 1, security.AlgCiphering128NEA0, security.AlgIntegrity128NIA2,
		models.AccessType__3_GPP_ACCESS)
	ue.AmfUeNgapId = 1
	ue.AuthenticationSubs = test.GetAuthSubscription(TestGenAuthData.MilenageTestSet19.K,
		TestGenAuthData.MilenageTestSet19.OPC,
		TestGenAuthData.MilenageTestSet19.OP)

	// insert UE data to MongoDB
	test.DelUeFromMongoDB(t, ue, servingPlmnId)
	test.InsertUeToMongoDB(t, ue, servingPlmnId)

	// send InitialUeMessage(Registration Request)(imsi-2089300007487)
	mobileIdentity5GS := nasType.MobileIdentity5GS{
		Len:    12, // suci
		Buffer: []uint8{0x01, 0x02, 0xf8, 0x39, 0xf0, 0xff, 0x00, 0x00, 0x00, 0x00, 0x47, 0x78},
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

func pduSessionEstablishment(t *testing.T, ue *test.RanUeContext, MranConn *sctp.SCTPConn) {
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
		upTransportLayerInformation.GTPTunnel.GTPTEID.Value = aper.OctetString(mranTeid)
		upTransportLayerInformation.GTPTunnel.TransportLayerAddress = ngapConvert.IPAddressToNgap(mranN3Addr, "")

		// Associated QoS Flow List in QoS Flow per TNL Information
		associatedQosFlowList := &qosFlowPerTNLInformation.AssociatedQosFlowList

		associatedQosFlowItem := ngapType.AssociatedQosFlowItem{}
		associatedQosFlowItem.QosFlowIdentifier.Value = 1
		associatedQosFlowList.List = append(associatedQosFlowList.List, associatedQosFlowItem)

		// DC QoS Flow per TNL Information
		DCQosFlowPerTNLInformationItem := ngapType.QosFlowPerTNLInformationItem{}
		DCQosFlowPerTNLInformationItem.QosFlowPerTNLInformation.UPTransportLayerInformation.Present = ngapType.UPTransportLayerInformationPresentGTPTunnel

		// DC Transport Layer Information in QoS Flow per TNL Information
		DCUpTransportLayerInformation := &DCQosFlowPerTNLInformationItem.QosFlowPerTNLInformation.UPTransportLayerInformation
		DCUpTransportLayerInformation.Present = ngapType.UPTransportLayerInformationPresentGTPTunnel
		DCUpTransportLayerInformation.GTPTunnel = new(ngapType.GTPTunnel)
		DCUpTransportLayerInformation.GTPTunnel.GTPTEID.Value = aper.OctetString(sranTeid)
		DCUpTransportLayerInformation.GTPTunnel.TransportLayerAddress = ngapConvert.IPAddressToNgap(sranN3Addr, "")

		// DC Associated QoS Flow List in QoS Flow per TNL Information
		DCAssociatedQosFlowList := &DCQosFlowPerTNLInformationItem.QosFlowPerTNLInformation.AssociatedQosFlowList
		DCAssociatedQosFlowItem := ngapType.AssociatedQosFlowItem{}
		DCAssociatedQosFlowItem.QosFlowIdentifier.Value = 1
		DCAssociatedQosFlowList.List = append(DCAssociatedQosFlowList.List, DCAssociatedQosFlowItem)

		// Additional DL QoS Flow per TNL Information
		data.AdditionalDLQosFlowPerTNLInformation = new(ngapType.QosFlowPerTNLInformationList)
		data.AdditionalDLQosFlowPerTNLInformation.List = append(data.AdditionalDLQosFlowPerTNLInformation.List, DCQosFlowPerTNLInformationItem)

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

func TestDCRegistration(t *testing.T) {
	// RANs connect to AMF
	MranConn, SranConn := connectRANsToAMF(t)
	if MranConn == nil || SranConn == nil {
		t.Fatal("Failed to connect to AMF")
		return
	}
	defer MranConn.Close()
	defer SranConn.Close()
	t.Log("Master RAN and Slave RAN connect to AMF successfully")

	// RANs connect to UPF
	MupfConn, SupfConn := connectRANsToUPF(t)
	if MupfConn == nil || SupfConn == nil {
		t.Fatal("Failed to connect to UPF")
		return
	}
	defer MupfConn.Close()
	defer SupfConn.Close()
	t.Log("Master RAN and Slave RAN connect to UPF successfully")

	// NGSetup
	nGsSetup(t, MranConn, SranConn)
	t.Log("Master RAN and Slave RAN NGSetup successfully")

	// New UE and initial registration(NAS/NGAP)
	ue := newUEAndInitialRegistration(t, MranConn)
	defer test.DelUeFromMongoDB(t, ue, servingPlmnId)
	t.Log("New UE and initial registration(NAS/NGAP) successfully")

	// PDU Session Establishment
	pduSessionEstablishment(t, ue, MranConn)
	t.Log("PDU Session Establishment successfully")
}
