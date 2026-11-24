package test_test

import (
	"encoding/json"
	"net/http"
	"net/url"
	"testing"

	"test"
	"test/consumerTestdata/UDM/TestGenAuthData"
	"test/nasTestpacket"

	nasIE "github.com/free5gc/nas/ie"
	nasMessage "github.com/free5gc/nas/message"
	ngapMessage "github.com/free5gc/ngap/message"
	"github.com/free5gc/openapi"
	"github.com/free5gc/openapi/models"
	"github.com/stretchr/testify/require"
)

const (
	scpBaseURL = "http://127.0.0.6:8000"
	nrfBaseURL = "http://127.0.0.10:8000"
)

// TestSCP verifies that NF service discovery routes a complete 5G-AKA
// authentication through SCP without making the consumer NFs SCP-aware.
func TestSCP(t *testing.T) {
	ue := test.NewRanUeContext(
		"imsi-208930000007487",
		1,
		nasMessage.AlgCiphering128NEA0,
		nasMessage.AlgIntegrity128NIA2,
		models.AccessType_3_GPP_ACCESS,
	)
	ue.AmfUeNgapId = 1
	ue.AuthenticationSubs = test.GetAuthSubscription(
		TestGenAuthData.MilenageTestSet19.K,
		TestGenAuthData.MilenageTestSet19.OPC,
		TestGenAuthData.MilenageTestSet19.OP,
	)
	servingPlmnID := "20893"

	defer func() {
		test.DelUeFromMongoDB(t, ue, servingPlmnID)
		NfTerminate()
	}()

	test.DelUeFromMongoDB(t, ue, servingPlmnID)
	test.InsertUeToMongoDB(t, ue, servingPlmnID)

	t.Run("NRF advertises producer services through SCP", func(t *testing.T) {
		requireDiscoveredServiceURI(t,
			models.Nrf_NFMgmt_NFType_AMF,
			models.Nrf_NFMgmt_NFType_AUSF,
			models.Nrf_NFMgmt_ServiceName_NAUSF_AUTH,
			scpBaseURL,
		)
		requireDiscoveredServiceURI(t,
			models.Nrf_NFMgmt_NFType_AUSF,
			models.Nrf_NFMgmt_NFType_UDM,
			models.Nrf_NFMgmt_ServiceName_NUDM_UEAU,
			scpBaseURL,
		)
		requireDiscoveredServiceURI(t,
			models.Nrf_NFMgmt_NFType_UDM,
			models.Nrf_NFMgmt_NFType_UDR,
			models.Nrf_NFMgmt_ServiceName_NUDR_DR,
			scpBaseURL,
		)
	})

	t.Run("complete 5G-AKA authentication through SCP", func(t *testing.T) {
		conn, err := test.ConnectToAmf("127.0.0.1", "127.0.0.1", 38412, 9487)
		require.NoError(t, err)
		defer conn.Close()

		sendMsg, err := test.GetNGSetupRequest([]byte("\x00\x01\x02"), 24, "free5gc", "", "", "")
		require.NoError(t, err)
		_, err = conn.Write(sendMsg)
		require.NoError(t, err)

		recvMsg := make([]byte, 2048)
		n, err := conn.Read(recvMsg)
		require.NoError(t, err)
		ngapPdu, err := ngapMessage.Parse(recvMsg[:n])
		require.NoError(t, err)
		require.Equal(t, ngapMessage.MessageTypeSuccessfulOutcome, ngapPdu.MessageType())
		require.Equal(t, ngapMessage.ProcedureCodeNGSetup, ngapPdu.ProcedureCode())

		mobileIdentity5GS := nasTestpacket.MobileIdentity5GS(
			[]byte{0x01, 0x02, 0xf8, 0x39, 0xf0, 0xff, 0x00, 0x00, 0x00, 0x00, 0x00, 0x47, 0x78},
		)
		ueSecurityCapability := ue.GetUESecurityCapability()
		registrationRequest := nasTestpacket.GetRegistrationRequest(
			nasIE.RegType_InitialReg,
			mobileIdentity5GS,
			nil,
			ueSecurityCapability,
			nil,
			nil,
			nil,
		)
		sendMsg, err = test.GetInitialUEMessage(ue.RanUeNgapId, registrationRequest, "")
		require.NoError(t, err)
		_, err = conn.Write(sendMsg)
		require.NoError(t, err)

		n, err = conn.Read(recvMsg)
		require.NoError(t, err)
		ngapPdu, err = ngapMessage.Parse(recvMsg[:n])
		require.NoError(t, err)
		require.Equal(t, ngapMessage.MessageTypeInitiatingMessage, ngapPdu.MessageType())
		require.Equal(t, ngapMessage.ProcedureCodeDownlinkNASTransport, ngapPdu.ProcedureCode())

		nasPdu := test.GetNasPdu(ue, ngapPdu.(*ngapMessage.DownlinkNASTransport))
		require.NotNil(t, nasPdu)
		require.Equal(t, nasMessage.MsgTypeAuthReq, nasPdu.MsgType(),
			"expected NAS Authentication Request")
		rand := nasPdu.(*nasMessage.AuthReq).AuthParamRAND5GAuthChlg.Rand
		resStar := ue.DeriveRESstarAndSetKey(
			ue.AuthenticationSubs,
			rand[:],
			"5G:mnc093.mcc208.3gppnetwork.org",
		)

		pdu := nasTestpacket.GetAuthenticationResponse(resStar, "")
		sendMsg, err = test.GetUplinkNASTransport(ue.AmfUeNgapId, ue.RanUeNgapId, pdu)
		require.NoError(t, err)
		_, err = conn.Write(sendMsg)
		require.NoError(t, err)

		// AMF sends Security Mode Command only after AUSF accepts RES* via
		// the 5G-AKA confirmation endpoint, proving authentication succeeded.
		n, err = conn.Read(recvMsg)
		require.NoError(t, err)
		ngapPdu, err = ngapMessage.Parse(recvMsg[:n])
		require.NoError(t, err)
		require.Equal(t, ngapMessage.MessageTypeInitiatingMessage, ngapPdu.MessageType())
		require.Equal(t, ngapMessage.ProcedureCodeDownlinkNASTransport, ngapPdu.ProcedureCode())

		nasPdu = test.GetNasPdu(ue, ngapPdu.(*ngapMessage.DownlinkNASTransport))
		require.NotNil(t, nasPdu)
		require.Equal(t, nasMessage.MsgTypeSecModeCmd, nasPdu.MsgType(),
			"authentication did not complete successfully")
	})
}

func requireDiscoveredServiceURI(
	t *testing.T,
	requester models.Nrf_NFMgmt_NFType,
	target models.Nrf_NFMgmt_NFType,
	service models.Nrf_NFMgmt_ServiceName,
	wantURI string,
) {
	t.Helper()

	query := url.Values{
		"requester-nf-type": []string{string(requester)},
		"target-nf-type":    []string{string(target)},
		"service-names":     []string{string(service)},
	}
	resp, err := http.Get(nrfBaseURL + "/nnrf-disc/v1/nf-instances?" + query.Encode())
	require.NoError(t, err)
	defer resp.Body.Close()
	require.Equal(t, http.StatusOK, resp.StatusCode)

	var result models.Nrf_NFDisc_SearchResult
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&result))
	require.NotEmpty(t, result.NfInstances)
	_, gotURI, err := openapi.GetServiceNfProfileAndUri(result.NfInstances, service)
	require.NoError(t, err)
	require.Equal(t, wantURI, gotURI)
}
