package test_test

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	"test"
	"test/consumerTestdata/UDM/TestGenAuthData"
	"test/nasTestpacket"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/free5gc/nas"
	"github.com/free5gc/nas/nasMessage"
	"github.com/free5gc/nas/nasType"
	"github.com/free5gc/nas/security"
	"github.com/free5gc/ngap"
	"github.com/free5gc/ngap/ngapType"
	"github.com/free5gc/openapi/models"
)

func TestULCLAndMultiUPF(t *testing.T) {
	testULCLSessionBase(t, 4, 6)
}

func testULCLSessionBase(t *testing.T, ueCount int, upfNum int) {
	ranN2Ipv4Addr := "127.0.0.1"
	amfN2Ipv4Addr := "127.0.0.18"
	ranN3Ipv4Addr := "10.200.200.1"

	var n int
	var sendMsg []byte
	var recvMsg = make([]byte, 2048)

	// RAN connect to AMF
	conn, err := test.ConnectToAmf(amfN2Ipv4Addr, ranN2Ipv4Addr, 38412, 9487)
	assert.Nil(t, err)

	// send NGSetupRequest Msg
	sendMsg, err = test.GetNGSetupRequest([]byte("\x00\x01\x02"), 24, "free5gc")
	assert.Nil(t, err)
	_, err = conn.Write(sendMsg)
	assert.Nil(t, err)

	// receive NGSetupResponse Msg
	n, err = conn.Read(recvMsg)
	assert.Nil(t, err)
	ngapPdu, err := ngap.Decoder(recvMsg[:n])
	assert.Nil(t, err)
	assert.True(t, ngapPdu.Present == ngapType.NGAPPDUPresentSuccessfulOutcome && ngapPdu.SuccessfulOutcome.ProcedureCode.Value == ngapType.ProcedureCodeNGSetup, "No NGSetupResponse received.")

	ueList := []*test.RanUeContext{}
	mobileIdentity5GSList := map[string]nasType.MobileIdentity5GS{}

	servingPlmnId := "20893"
	sNssai := models.Snssai{
		Sst: 1,
		Sd:  "010203",
	}

	for i := 0; i < ueCount; i++ {
		// New UE
		imsi_e := fmt.Sprintf("%04d", i)
		imsi_all := "imsi-208930000" + imsi_e
		fmt.Println(imsi_all)
		ue := test.NewRanUeContext(imsi_all, int64(i+1), security.AlgCiphering128NEA0, security.AlgIntegrity128NIA2,
			models.AccessType__3_GPP_ACCESS)
		ue.AmfUeNgapId = int64(i + 1)
		ue.AuthenticationSubs = test.GetAuthSubscription(TestGenAuthData.MilenageTestSet19.K,
			TestGenAuthData.MilenageTestSet19.OPC,
			TestGenAuthData.MilenageTestSet19.OP)
		// insert UE data to MongoDB
		test.InsertAuthSubscriptionToMongoDB(ue.Supi, ue.AuthenticationSubs)
		getData := test.GetAuthSubscriptionFromMongoDB(ue.Supi)
		assert.NotNil(t, getData)
		{
			amData := test.GetAccessAndMobilitySubscriptionData()
			test.InsertAccessAndMobilitySubscriptionDataToMongoDB(ue.Supi, amData, servingPlmnId)
			getData := test.GetAccessAndMobilitySubscriptionDataFromMongoDB(ue.Supi, servingPlmnId)
			assert.NotNil(t, getData)
		}
		{
			smfSelData := test.GetSmfSelectionSubscriptionData()
			test.InsertSmfSelectionSubscriptionDataToMongoDB(ue.Supi, smfSelData, servingPlmnId)
			getData := test.GetSmfSelectionSubscriptionDataFromMongoDB(ue.Supi, servingPlmnId)
			assert.NotNil(t, getData)
		}
		{
			smSelData := test.GetSessionManagementSubscriptionData()
			test.InsertSessionManagementSubscriptionDataToMongoDB(ue.Supi, servingPlmnId, smSelData)
			getData := test.GetSessionManagementDataFromMongoDB(ue.Supi, servingPlmnId)
			assert.NotNil(t, getData)
		}
		{
			amPolicyData := test.GetAmPolicyData()
			test.InsertAmPolicyDataToMongoDB(ue.Supi, amPolicyData)
			getData := test.GetAmPolicyDataFromMongoDB(ue.Supi)
			assert.NotNil(t, getData)
		}
		{
			smPolicyData := test.GetSmPolicyData()
			test.InsertSmPolicyDataToMongoDB(ue.Supi, smPolicyData)
			getData := test.GetSmPolicyDataFromMongoDB(ue.Supi)
			assert.NotNil(t, getData)
		}

		// send InitialUeMessage(Registration Request)
		//i%100
		i_e2 := 16*((i%100)%10) + ((i % 100) / 10)
		//i/100
		i_e4 := 16*((i/100)%10) + ((i / 100) / 10)
		mobileIdentity5GS := nasType.MobileIdentity5GS{
			Len:    12, // suci
			Buffer: []uint8{0x01, 0x02, 0xf8, 0x39, 0xf0, 0xff, 0x00, 0x00, 0x00, 0x00, uint8(i_e4), uint8(i_e2)},
		}
		mobileIdentity5GSList[ue.Supi] = mobileIdentity5GS

		ueSecurityCapability := ue.GetUESecurityCapability()
		registrationRequest := nasTestpacket.GetRegistrationRequest(
			nasMessage.RegistrationType5GSInitialRegistration, mobileIdentity5GS, nil, ueSecurityCapability, nil, nil, nil)
		sendMsg, err = test.GetInitialUEMessage(ue.RanUeNgapId, registrationRequest, "")
		assert.Nil(t, err)
		_, err = conn.Write(sendMsg)
		assert.Nil(t, err)

		// receive NAS Authentication Request Msg
		n, err = conn.Read(recvMsg)
		assert.Nil(t, err)
		ngapPdu, err = ngap.Decoder(recvMsg[:n])
		assert.Nil(t, err)
		assert.True(t, ngapPdu.Present == ngapType.NGAPPDUPresentInitiatingMessage, "No NGAP Initiating Message received.")

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
		_, err = conn.Write(sendMsg)
		assert.Nil(t, err)

		// receive NAS Security Mode Command Msg
		n, err = conn.Read(recvMsg)
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
		_, err = conn.Write(sendMsg)
		assert.Nil(t, err)

		// receive ngap Initial Context Setup Request Msg
		n, err = conn.Read(recvMsg)
		assert.Nil(t, err)
		ngapPdu, err = ngap.Decoder(recvMsg[:n])
		assert.Nil(t, err)
		assert.True(t, ngapPdu.Present == ngapType.NGAPPDUPresentInitiatingMessage &&
			ngapPdu.InitiatingMessage.ProcedureCode.Value == ngapType.ProcedureCodeInitialContextSetup,
			"No InitialContextSetup received.")

		// send ngap Initial Context Setup Response Msg
		sendMsg, err = test.GetInitialContextSetupResponse(ue.AmfUeNgapId, ue.RanUeNgapId)
		assert.Nil(t, err)
		_, err = conn.Write(sendMsg)
		assert.Nil(t, err)

		// send NAS Registration Complete Msg
		pdu = nasTestpacket.GetRegistrationComplete(nil)
		pdu, err = test.EncodeNasPduWithSecurity(ue, pdu, nas.SecurityHeaderTypeIntegrityProtectedAndCiphered, true, false)
		assert.Nil(t, err)
		sendMsg, err = test.GetUplinkNASTransport(ue.AmfUeNgapId, ue.RanUeNgapId, pdu)
		assert.Nil(t, err)
		_, err = conn.Write(sendMsg)
		assert.Nil(t, err)

		time.Sleep(100 * time.Millisecond)
		// send GetPduSessionEstablishmentRequest Msg

		pdu = nasTestpacket.GetUlNasTransport_PduSessionEstablishmentRequest(10, nasMessage.ULNASTransportRequestTypeInitialRequest, "internet", &sNssai)
		pdu, err = test.EncodeNasPduWithSecurity(ue, pdu, nas.SecurityHeaderTypeIntegrityProtectedAndCiphered, true, false)
		assert.Nil(t, err)
		sendMsg, err = test.GetUplinkNASTransport(ue.AmfUeNgapId, ue.RanUeNgapId, pdu)
		assert.Nil(t, err)
		_, err = conn.Write(sendMsg)
		assert.Nil(t, err)

		// receive 12. NGAP-PDU Session Resource Setup Request(DL nas transport((NAS msg-PDU session setup Accept)))
		n, err = conn.Read(recvMsg)
		assert.Nil(t, err)
		ngapPdu, err = ngap.Decoder(recvMsg[:n])
		assert.Nil(t, err)
		assert.True(t, ngapPdu.Present == ngapType.NGAPPDUPresentInitiatingMessage &&
			ngapPdu.InitiatingMessage.ProcedureCode.Value == ngapType.ProcedureCodePDUSessionResourceSetup,
			"No PDUSessionResourceSetup received.")

		// send 14. NGAP-PDU Session Resource Setup Response
		sendMsg, err = test.GetPDUSessionResourceSetupResponse(10, ue.AmfUeNgapId, ue.RanUeNgapId, ranN3Ipv4Addr)
		assert.Nil(t, err)
		_, err = conn.Write(sendMsg)
		assert.Nil(t, err)
		ueList = append(ueList, ue)

		// check PDR and FAR start(For ULCL)
		dir, _ := os.Getwd()
		cmdPath := dir + "/../go-gtp5gnl/bin/"
		gtp5gTunnelCmdPath := filepath.Clean(cmdPath)

		for ns_num := 1; ns_num < upfNum+1; ns_num++ {
			ns_name := fmt.Sprintf("UPFns0%d", ns_num)
			fmt.Println("---- List PDR ---")
			cmd := exec.Command("sudo", "ip", "netns", "exec", ns_name, "bash", "-c",
				gtp5gTunnelCmdPath+"/gtp5g-tunnel list pdr")
			out, err := cmd.Output()
			if err != nil {
				fmt.Println(err)
			}
			fmt.Printf("%s\n", out)

			fmt.Println("---- List FAR ---")
			cmd = exec.Command("sudo", "ip", "netns", "exec", ns_name, "bash", "-c",
				gtp5gTunnelCmdPath+"/gtp5g-tunnel list far")
			out, err = cmd.Output()
			if err != nil {
				fmt.Println(err)
			}
			fmt.Printf("%s\n", out)

			// wait 1s
			time.Sleep(1 * time.Second)
		}
		// check PDR and FAR end

		// wait 1s
		time.Sleep(1 * time.Second)
	}
	for _, ue := range ueList {
		// Send Pdu Session Establishment Release Request
		pdu := nasTestpacket.GetUlNasTransport_PduSessionReleaseRequest(10)
		pdu, err = test.EncodeNasPduWithSecurity(ue, pdu, nas.SecurityHeaderTypeIntegrityProtectedAndCiphered, true, false)
		assert.Nil(t, err)
		sendMsg, err = test.GetUplinkNASTransport(ue.AmfUeNgapId, ue.RanUeNgapId, pdu)
		assert.Nil(t, err)
		_, err = conn.Write(sendMsg)
		assert.Nil(t, err)

		time.Sleep(1000 * time.Millisecond)

		// receive NGAP-PDU Session Resource Release Request
		n, err = conn.Read(recvMsg)
		require.Nil(t, err)
		ngapPdu, err = ngap.Decoder(recvMsg[:n])
		require.Nil(t, err)
		require.True(t, ngapPdu.Present == ngapType.NGAPPDUPresentInitiatingMessage &&
			ngapPdu.InitiatingMessage.ProcedureCode.Value == ngapType.ProcedureCodePDUSessionResourceRelease,
			"No PDUSessionResourceRelease received.")

		// send N2 Resource Release Ack(PDUSession Resource Release Response)
		sendMsg, err = test.GetPDUSessionResourceReleaseResponse(ue.AmfUeNgapId, ue.RanUeNgapId)
		assert.Nil(t, err)
		_, err = conn.Write(sendMsg)
		assert.Nil(t, err)

		// wait 10 ms
		time.Sleep(1000 * time.Millisecond)

		//send N1 PDU Session Release Ack PDU session release complete
		pdu = nasTestpacket.GetUlNasTransport_PduSessionReleaseComplete(10, nasMessage.ULNASTransportRequestTypeExistingPduSession, "internet", &sNssai)
		pdu, err = test.EncodeNasPduWithSecurity(ue, pdu, nas.SecurityHeaderTypeIntegrityProtectedAndCiphered, true, false)
		assert.Nil(t, err)
		sendMsg, err = test.GetUplinkNASTransport(ue.AmfUeNgapId, ue.RanUeNgapId, pdu)
		assert.Nil(t, err)
		_, err = conn.Write(sendMsg)
		assert.Nil(t, err)

		// wait result
		time.Sleep(1 * time.Second)

		// send NAS Deregistration Request (UE Originating)
		pdu = nasTestpacket.GetDeregistrationRequest(nasMessage.AccessType3GPP, 0, 0x04, mobileIdentity5GSList[ue.Supi])
		pdu, err = test.EncodeNasPduWithSecurity(ue, pdu, nas.SecurityHeaderTypeIntegrityProtectedAndCiphered, true, false)
		require.Nil(t, err)
		sendMsg, err = test.GetUplinkNASTransport(ue.AmfUeNgapId, ue.RanUeNgapId, pdu)
		require.Nil(t, err)
		_, err = conn.Write(sendMsg)
		require.Nil(t, err)

		time.Sleep(500 * time.Millisecond)

		// receive Deregistration Accept
		n, err = conn.Read(recvMsg)
		require.Nil(t, err)
		ngapPdu, err = ngap.Decoder(recvMsg[:n])
		require.Nil(t, err)
		require.True(t, ngapPdu.Present == ngapType.NGAPPDUPresentInitiatingMessage &&
			ngapPdu.InitiatingMessage.ProcedureCode.Value == ngapType.ProcedureCodeDownlinkNASTransport,
			"No DownlinkNASTransport received.")
		nasPdu := test.GetNasPdu(ue, ngapPdu.InitiatingMessage.Value.DownlinkNASTransport)
		require.NotNil(t, nasPdu, "NAS PDU is nil")
		require.NotNil(t, nasPdu.GmmMessage, "GMM message is nil")
		require.Equal(t, nasPdu.GmmHeader.GetMessageType(), nas.MsgTypeDeregistrationAcceptUEOriginatingDeregistration,
			"Received wrong GMM message. Expected Deregistration Accept.")

		// receive ngap UE Context Release Command
		n, err = conn.Read(recvMsg)
		require.Nil(t, err)
		ngapPdu, err = ngap.Decoder(recvMsg[:n])
		require.Nil(t, err)
		require.True(t, ngapPdu.Present == ngapType.NGAPPDUPresentInitiatingMessage &&
			ngapPdu.InitiatingMessage.ProcedureCode.Value == ngapType.ProcedureCodeUEContextRelease,
			"No UEContextReleaseCommand received.")

		// send ngap UE Context Release Complete
		sendMsg, err = test.GetUEContextReleaseComplete(ue.AmfUeNgapId, ue.RanUeNgapId, nil)
		require.Nil(t, err)
		_, err = conn.Write(sendMsg)
		require.Nil(t, err)

		time.Sleep(100 * time.Millisecond)
	}

	for _, ue := range ueList {
		// delete test data
		test.DelAuthSubscriptionToMongoDB(ue.Supi)
		test.DelAccessAndMobilitySubscriptionDataFromMongoDB(ue.Supi, servingPlmnId)
		test.DelSmfSelectionSubscriptionDataFromMongoDB(ue.Supi, servingPlmnId)
	}

	// close Connection
	conn.Close()

	// terminate all NF
	NfTerminate()
}
