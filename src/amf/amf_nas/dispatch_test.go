package amf_nas_test

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"flag"
	"fmt"
	"github.com/urfave/cli"
	"free5gc/lib/CommonConsumerTestData/AMF/TestAmf"
	"free5gc/lib/CommonConsumerTestData/AMF/TestComm"
	"free5gc/lib/CommonConsumerTestData/AUSF/TestUEAuth"
	"free5gc/lib/CommonConsumerTestData/UDM/TestGenAuthData"
	"free5gc/lib/MongoDBLibrary"
	"free5gc/lib/Namf_Communication"
	"free5gc/lib/http2_util"
	"free5gc/lib/nas/nasMessage"
	"free5gc/lib/nas/nasTestpacket"
	"free5gc/lib/nas/nasType"
	"free5gc/lib/ngap/ngapType"
	"free5gc/lib/openapi/common"
	"free5gc/lib/openapi/models"
	"free5gc/lib/path_util"
	"free5gc/src/amf/Communication"
	"free5gc/src/amf/amf_consumer"
	"free5gc/src/amf/amf_handler"
	"free5gc/src/amf/amf_nas"
	"free5gc/src/amf/gmm/gmm_state"
	"free5gc/src/amf/logger"
	Nausf_UEAU "free5gc/src/ausf/UEAuthentication"
	"free5gc/src/ausf/ausf_context"
	"free5gc/src/ausf/ausf_handler"
	"free5gc/src/ausf/ausf_producer"
	"free5gc/src/nrf/nrf_handler"
	"free5gc/src/nrf/nrf_service"
	"free5gc/src/smf/smf_service"
	Nudm_UEAU "free5gc/src/udm/UEAuthentication"
	"free5gc/src/udm/udm_handler"
	"log"
	"net/http"
	"testing"
	"time"

	"github.com/bronze1man/radius"
	"github.com/gin-gonic/gin"
	"github.com/ishidawataru/sctp"
	"github.com/mohae/deepcopy"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
)

var flags flag.FlagSet
var c = cli.NewContext(nil, &flags, nil)

func addNon3gppRan() {
	TestAmf.Conn2, _ = sctp.DialSCTP("sctp", TestAmf.Laddr2, TestAmf.ServerAddr)
	time.Sleep(10 * time.Millisecond)
	ran := TestAmf.TestAmf.AmfRanPool[TestAmf.Laddr2.String()]
	ran.AnType = models.AccessType_NON_3_GPP_ACCESS
	ranUe := ran.NewRanUe()
	ue := TestAmf.TestAmf.UePool["imsi-2089300007487"]
	ue.AttachRanUe(ranUe)
}

func ausfInit() {
	go func() { // ausf server
		router := Nausf_UEAU.NewRouter()

		ausfLogPath := path_util.Gofree5gcPath("free5gc/ausfsslkey.log")
		ausfPemPath := path_util.Gofree5gcPath("free5gc/support/TLS/ausf.pem")
		ausfKeyPath := path_util.Gofree5gcPath("free5gc/support/TLS/ausf.key")

		server, err := http2_util.NewServer(":29509", ausfLogPath, router)
		if err == nil && server != nil {
			logger.InitLog.Infoln(server.ListenAndServeTLS(ausfPemPath, ausfKeyPath))
		}
	}()

	go ausf_handler.Handle()
}

func udmInit() {
	go func() { // fake udm server
		router := Nudm_UEAU.NewRouter()

		udmLogPath := path_util.Gofree5gcPath("free5gc/udmsslkey.log")
		udmPemPath := path_util.Gofree5gcPath("free5gc/support/TLS/udm.pem")
		udmKeyPath := path_util.Gofree5gcPath("free5gc/support/TLS/udm.key")

		server, err := http2_util.NewServer(":29503", udmLogPath, router)
		if err == nil && server != nil {
			logger.InitLog.Infoln(server.ListenAndServeTLS(udmPemPath, udmKeyPath))
		}
	}()

	go udm_handler.Handle()
}

func udrInit() {
	go func() { // fake udr server
		router := gin.Default()

		router.GET("/nudr-dr/v1/subscription-data/:ueId/authentication-data/authentication-subscription", func(c *gin.Context) {
			ueId := c.Param("ueId")
			fmt.Println("ueId: ", ueId)
			var authSubs models.AuthenticationSubscription
			var pk models.PermanentKey
			var opc models.Opc
			var var_milenage models.Milenage
			var op models.Op

			pk.PermanentKeyValue = TestGenAuthData.MilenageTestSet19.K
			opc.OpcValue = TestGenAuthData.MilenageTestSet19.OPC
			op.OpValue = TestGenAuthData.MilenageTestSet19.OP
			var_milenage.Op = &op

			authSubs.PermanentKey = &pk
			authSubs.Opc = &opc
			authSubs.Milenage = &var_milenage
			authSubs.SequenceNumber = TestGenAuthData.MilenageTestSet19.SQN
			authSubs.AuthenticationMethod = models.AuthMethod__5_G_AKA
			// authSubs.AuthenticationMethod = models.AuthMethod_EAP_AKA_PRIME

			c.JSON(http.StatusOK, authSubs)
		})

		router.PUT("/nudr-dr/v1/subscription-data/:ueId/authentication-data/authentication-status", func(c *gin.Context) {
			ueId := c.Param("ueId")
			fmt.Println("===================================")
			fmt.Println("ueId: ", ueId)
			c.JSON(http.StatusNoContent, gin.H{})
		})

		udrLogPath := path_util.Gofree5gcPath("free5gc/udrsslkey.log")
		udrPemPath := path_util.Gofree5gcPath("free5gc/support/TLS/udr.pem")
		udrKeyPath := path_util.Gofree5gcPath("free5gc/support/TLS/udr.key")

		server, err := http2_util.NewServer(":29504", udrLogPath, router)
		if err == nil && server != nil {
			logger.InitLog.Infoln(server.ListenAndServeTLS(udrPemPath, udrKeyPath))
		}
	}()
}

func nrfInit() {
	nrf := &nrf_service.NRF{}

	nrf.Initialize(c)
	go nrf.Start()
	time.Sleep(10 * time.Millisecond)

}

func smfInit() {
	smf := &smf_service.SMF{}

	smf.Initialize(c)
	go smf.Start()
	time.Sleep(10 * time.Millisecond)
}

func init() {

	go amf_handler.Handle()
	time.Sleep(10 * time.Millisecond)
	go nrf_handler.Handle()

	time.Sleep(10 * time.Millisecond)
	TestAmf.SctpSever()
	time.Sleep(10 * time.Millisecond)

	udmInit()
	udrInit()
	ausfInit()
	nrfInit()
	smfInit()
}

func sendN1N2Transfer(client *Namf_Communication.APIClient, supi string, request models.N1N2MessageTransferRequest) {
	n1N2MessageTransferResponse, httpResponse, err := client.N1N2MessageCollectionDocumentApi.N1N2MessageTransfer(context.Background(), supi, request)
	if err != nil {
		if httpResponse == nil {
			log.Panic(err)
		} else if err.Error() != httpResponse.Status {
			log.Panic(err)
		} else if httpResponse.StatusCode == 504 || httpResponse.StatusCode == 409 {
			transferError := err.(common.GenericOpenAPIError).Model().(models.N1N2MessageTransferError)
			TestAmf.Config.Dump(transferError)
		} else {
			probelmDetail := err.(common.GenericOpenAPIError).Model().(models.ProblemDetails)
			TestAmf.Config.Dump(probelmDetail)
		}
	} else {
		TestAmf.Config.Dump(n1N2MessageTransferResponse)
	}

}

func TestULNASTransportPDUSessionEstablishemnt(t *testing.T) {

	time.Sleep(200 * time.Millisecond)
	MongoDBLibrary.RestfulAPIDeleteMany("NfProfile", bson.M{})

	// smf register to nrf
	uuid, profile := TestAmf.BuildSmfNfProfile()
	uri, err := amf_consumer.SendRegisterNFInstance("https://localhost:29510", uuid, profile)
	if err != nil {
		t.Error(err.Error())
	} else {
		TestAmf.Config.Dump(uri)
	}

	TestAmf.AmfInit()
	TestAmf.SctpConnectToServer(models.AccessType__3_GPP_ACCESS)
	time.Sleep(10 * time.Millisecond)
	ue := TestAmf.TestAmf.UePool["imsi-2089300007487"]
	sNssai := models.Snssai{
		Sst: 1,
		Sd:  "010203",
	}
	ue.PlmnId.Mcc = "208"
	ue.PlmnId.Mnc = "93"
	// InitialRequest with known dnn and Snssai (success)
	nasPdu := nasTestpacket.GetUlNasTransport_PduSessionEstablishmentRequest(10, nasMessage.ULNASTransportRequestTypeInitialRequest, "internet", &sNssai)

	amf_nas.HandleNAS(ue.RanUe[models.AccessType__3_GPP_ACCESS], ngapType.ProcedureCodeUplinkNASTransport, nasPdu)

	// InitialRequest with exist pdusession id (Duplicatd pdu session id)(success)
	nasPdu = nasTestpacket.GetUlNasTransport_PduSessionEstablishmentRequest(10, nasMessage.ULNASTransportRequestTypeInitialRequest, "internet", &sNssai)

	amf_nas.HandleNAS(ue.RanUe[models.AccessType__3_GPP_ACCESS], ngapType.ProcedureCodeUplinkNASTransport, nasPdu)

	// InitialRequest with unknown dnn and Snssai (success)
	nasPdu = nasTestpacket.GetUlNasTransport_PduSessionEstablishmentRequest(11, nasMessage.ULNASTransportRequestTypeInitialRequest, "", nil)

	amf_nas.HandleNAS(ue.RanUe[models.AccessType__3_GPP_ACCESS], ngapType.ProcedureCodeUplinkNASTransport, nasPdu)

	// InitialRequest with unknown dnn (failed)
	nasPdu = nasTestpacket.GetUlNasTransport_PduSessionEstablishmentRequest(12, nasMessage.ULNASTransportRequestTypeInitialRequest, "nctu.edu.tw", &sNssai)

	amf_nas.HandleNAS(ue.RanUe[models.AccessType__3_GPP_ACCESS], ngapType.ProcedureCodeUplinkNASTransport, nasPdu)

	// handover with allow snssai is empty (failed)
	nasPdu = nasTestpacket.GetUlNasTransport_PduSessionEstablishmentRequest(11, nasMessage.ULNASTransportRequestTypeExistingPduSession, "", nil)

	amf_nas.HandleNAS(ue.RanUe[models.AccessType__3_GPP_ACCESS], ngapType.ProcedureCodeUplinkNASTransport, nasPdu)

	// handover with allow snssai is not empty (success)
	ue.AllowedNssai[models.AccessType_NON_3_GPP_ACCESS] = deepcopy.Copy(ue.AllowedNssai[models.AccessType__3_GPP_ACCESS]).([]models.Snssai)
	nasPdu = nasTestpacket.GetUlNasTransport_PduSessionEstablishmentRequest(11, nasMessage.ULNASTransportRequestTypeExistingPduSession, "", nil)

	amf_nas.HandleNAS(ue.RanUe[models.AccessType__3_GPP_ACCESS], ngapType.ProcedureCodeUplinkNASTransport, nasPdu)

	TestAmf.Conn.Close()
	time.Sleep(10 * time.Millisecond)

}

func TestULNASTransportPDUSessionModification(t *testing.T) {
	TestAmf.AmfInit()
	TestAmf.SctpConnectToServer(models.AccessType__3_GPP_ACCESS)

	time.Sleep(10 * time.Millisecond)
	ue := TestAmf.TestAmf.UePool["imsi-2089300007487"]
	sNssai := models.Snssai{
		Sst: 1,
		Sd:  "010203",
	}
	// Pdu session Establishment InitialRequest with known dnn and Snssai (success)
	nasPdu := nasTestpacket.GetUlNasTransport_PduSessionEstablishmentRequest(10, nasMessage.ULNASTransportRequestTypeInitialRequest, "internet", &sNssai)

	amf_nas.HandleNAS(ue.RanUe[models.AccessType__3_GPP_ACCESS], ngapType.ProcedureCodeUplinkNASTransport, nasPdu)

	// Pdu session Modification Request(success)
	nasPdu = nasTestpacket.GetUlNasTransport_PduSessionCommonData(10, nasTestpacket.PDUSesModiReq)

	amf_nas.HandleNAS(ue.RanUe[models.AccessType__3_GPP_ACCESS], ngapType.ProcedureCodeUplinkNASTransport, nasPdu)

	// Pdu session Modification Complete (success)
	nasPdu = nasTestpacket.GetUlNasTransport_PduSessionCommonData(10, nasTestpacket.PDUSesModiCmp)

	amf_nas.HandleNAS(ue.RanUe[models.AccessType__3_GPP_ACCESS], ngapType.ProcedureCodeUplinkNASTransport, nasPdu)

	// Pdu session Modification Command Reject (success)
	nasPdu = nasTestpacket.GetUlNasTransport_PduSessionCommonData(10, nasTestpacket.PDUSesModiCmdRej)

	amf_nas.HandleNAS(ue.RanUe[models.AccessType__3_GPP_ACCESS], ngapType.ProcedureCodeUplinkNASTransport, nasPdu)
	TestAmf.Conn.Close()

	time.Sleep(10 * time.Millisecond)
}

func TestULNASTransportPDUSessionRelease(t *testing.T) {

	TestAmf.AmfInit()
	TestAmf.SctpConnectToServer(models.AccessType__3_GPP_ACCESS)

	time.Sleep(10 * time.Millisecond)
	ue := TestAmf.TestAmf.UePool["imsi-2089300007487"]
	sNssai := models.Snssai{
		Sst: 1,
		Sd:  "010203",
	}
	// Pdu session Establishment InitialRequest with known dnn and Snssai (success)
	nasPdu := nasTestpacket.GetUlNasTransport_PduSessionEstablishmentRequest(10, nasMessage.ULNASTransportRequestTypeInitialRequest, "internet", &sNssai)

	amf_nas.HandleNAS(ue.RanUe[models.AccessType__3_GPP_ACCESS], ngapType.ProcedureCodeUplinkNASTransport, nasPdu)

	// Pdu session Release Request(success)
	nasPdu = nasTestpacket.GetUlNasTransport_PduSessionCommonData(10, nasTestpacket.PDUSesRelReq)

	amf_nas.HandleNAS(ue.RanUe[models.AccessType__3_GPP_ACCESS], ngapType.ProcedureCodeUplinkNASTransport, nasPdu)

	// Pdu session Release Complete (success)
	nasPdu = nasTestpacket.GetUlNasTransport_PduSessionCommonData(10, nasTestpacket.PDUSesRelCmp)

	amf_nas.HandleNAS(ue.RanUe[models.AccessType__3_GPP_ACCESS], ngapType.ProcedureCodeUplinkNASTransport, nasPdu)

	// Pdu session Release Command Reject (success)
	nasPdu = nasTestpacket.GetUlNasTransport_PduSessionCommonData(10, nasTestpacket.PDUSesRelRej)

	amf_nas.HandleNAS(ue.RanUe[models.AccessType__3_GPP_ACCESS], ngapType.ProcedureCodeUplinkNASTransport, nasPdu)
	TestAmf.Conn.Close()

	time.Sleep(10 * time.Millisecond)
}

func TestULNASTransportPDUSessionAuthentication(t *testing.T) {

	TestAmf.AmfInit()
	TestAmf.SctpConnectToServer(models.AccessType__3_GPP_ACCESS)

	time.Sleep(10 * time.Millisecond)
	ue := TestAmf.TestAmf.UePool["imsi-2089300007487"]
	sNssai := models.Snssai{
		Sst: 1,
		Sd:  "010203",
	}
	// Pdu session Establishment InitialRequest with known dnn and Snssai (success)
	nasPdu := nasTestpacket.GetUlNasTransport_PduSessionEstablishmentRequest(10, nasMessage.ULNASTransportRequestTypeInitialRequest, "internet", &sNssai)

	amf_nas.HandleNAS(ue.RanUe[models.AccessType__3_GPP_ACCESS], ngapType.ProcedureCodeUplinkNASTransport, nasPdu)

	// Pdu session Release Complete (success)
	nasPdu = nasTestpacket.GetUlNasTransport_PduSessionCommonData(10, nasTestpacket.PDUSesAuthCmp)

	amf_nas.HandleNAS(ue.RanUe[models.AccessType__3_GPP_ACCESS], ngapType.ProcedureCodeUplinkNASTransport, nasPdu)

	TestAmf.Conn.Close()

	time.Sleep(10 * time.Millisecond)
}

func TestIdentityResponse(t *testing.T) {

	TestAmf.AmfInit()
	TestAmf.SctpConnectToServer(models.AccessType__3_GPP_ACCESS)

	time.Sleep(10 * time.Millisecond)
	ue := TestAmf.TestAmf.UePool["imsi-2089300007487"]

	// imei
	mobilityIdentity := nasType.MobileIdentity{
		Len:    8,
		Buffer: []uint8{0x9b, 0x32, 0x01, 0x46, 0x11, 0x22, 0x33, 0x77},
	}
	nasPdu := nasTestpacket.GetIdentityResponse(mobilityIdentity)

	amf_nas.HandleNAS(ue.RanUe[models.AccessType__3_GPP_ACCESS], ngapType.ProcedureCodeUplinkNASTransport, nasPdu)

	// imeisv
	mobilityIdentity = nasType.MobileIdentity{
		Len:    9,
		Buffer: []uint8{0x95, 0x32, 0x01, 0x46, 0x11, 0x22, 0x33, 0x77, 0xf8},
	}
	nasPdu = nasTestpacket.GetIdentityResponse(mobilityIdentity)

	amf_nas.HandleNAS(ue.RanUe[models.AccessType__3_GPP_ACCESS], ngapType.ProcedureCodeUplinkNASTransport, nasPdu)
	TestAmf.Conn.Close()
}

func TestNotificationResponse(t *testing.T) {

	TestAmf.AmfInit()
	TestAmf.UeAttach(models.AccessType_NON_3_GPP_ACCESS)

	time.Sleep(10 * time.Millisecond)
	ue := TestAmf.TestAmf.UePool["imsi-2089300007487"]
	sNssai := models.Snssai{
		Sst: 1,
		Sd:  "010203",
	}
	// InitialRequest with known dnn and Snssai (success)
	nasPdu := nasTestpacket.GetUlNasTransport_PduSessionEstablishmentRequest(10, nasMessage.ULNASTransportRequestTypeInitialRequest, "internet", &sNssai)

	amf_nas.HandleNAS(ue.RanUe[models.AccessType_NON_3_GPP_ACCESS], ngapType.ProcedureCodeUplinkNASTransport, nasPdu)

	// Notification Request
	nasPdu = nasTestpacket.GetNotificationResponse([]uint8{0x00, 0x40})

	amf_nas.HandleNAS(ue.RanUe[models.AccessType_NON_3_GPP_ACCESS], ngapType.ProcedureCodeUplinkNASTransport, nasPdu)

	time.Sleep(10 * time.Millisecond)
	TestAmf.Conn2.Close()
}
func TestUeConfiguratioinUpdateComplete(t *testing.T) {

	TestAmf.AmfInit()
	TestAmf.SctpConnectToServer(models.AccessType__3_GPP_ACCESS)

	time.Sleep(10 * time.Millisecond)
	ue := TestAmf.TestAmf.UePool["imsi-2089300007487"]

	nasPdu := nasTestpacket.GetConfigurationUpdateComplete()

	amf_nas.HandleNAS(ue.RanUe[models.AccessType__3_GPP_ACCESS], ngapType.ProcedureCodeUplinkNASTransport, nasPdu)
	TestAmf.Conn.Close()
}

func TestServiceRequest(t *testing.T) {
	go func() {
		router := Communication.NewRouter()
		server, err := http2_util.NewServer(":29518", TestAmf.AmfLogPath, router)
		if err == nil && server != nil {
			err = server.ListenAndServeTLS(TestAmf.AmfPemPath, TestAmf.AmfKeyPath)
			assert.True(t, err == nil, err.Error())
		}
	}()
	TestAmf.AmfInit()
	TestAmf.SctpConnectToServer(models.AccessType__3_GPP_ACCESS)
	ue := TestAmf.TestAmf.UePool["imsi-2089300007487"]
	ue.AccessAndMobilitySubscriptionData = &models.AccessAndMobilitySubscriptionData{
		SubscribedUeAmbr: &models.AmbrRm{
			Uplink:   "500",
			Downlink: "500",
		},
	}
	addNon3gppRan()

	err := ue.Sm[models.AccessType_NON_3_GPP_ACCESS].Transfer(gmm_state.REGISTERED, nil)
	assert.True(t, err == nil)
	time.Sleep(10 * time.Millisecond)

	err = ue.Sm[models.AccessType__3_GPP_ACCESS].Transfer(gmm_state.REGISTERED, nil)
	assert.True(t, err == nil)
	ue.NCC = 5
	ue.NH, _ = hex.DecodeString("ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff")
	ue.SecurityContextAvailable = true
	sNssai := models.Snssai{
		Sst: 1,
		Sd:  "010203",
	}
	// Pdu session Establishment InitialRequest with known dnn and Snssai (success)
	nasPdu := nasTestpacket.GetUlNasTransport_PduSessionEstablishmentRequest(10, nasMessage.ULNASTransportRequestTypeInitialRequest, "internet", &sNssai)

	amf_nas.HandleNAS(ue.RanUe[models.AccessType__3_GPP_ACCESS], ngapType.ProcedureCodeUplinkNASTransport, nasPdu)

	// Pdu session Establishment InitialRequest with known dnn and Snssai (success)
	nasPdu = nasTestpacket.GetUlNasTransport_PduSessionEstablishmentRequest(11, nasMessage.ULNASTransportRequestTypeInitialRequest, "internet", &sNssai)

	amf_nas.HandleNAS(ue.RanUe[models.AccessType_NON_3_GPP_ACCESS], ngapType.ProcedureCodeUplinkNASTransport, nasPdu)

	// Trigger By Ue, Data (Service Accept)
	nasPdu = nasTestpacket.GetServiceRequest(nasMessage.ServiceTypeData)

	amf_nas.HandleNAS(ue.RanUe[models.AccessType__3_GPP_ACCESS], ngapType.ProcedureCodeUplinkNASTransport, nasPdu)

	// Trigger By Ue, Data (Service Accept (initial context setup))
	nasPdu = nasTestpacket.GetServiceRequest(nasMessage.ServiceTypeData)

	amf_nas.HandleNAS(ue.RanUe[models.AccessType__3_GPP_ACCESS], ngapType.ProcedureCodeInitialUEMessage, nasPdu)

	// n1n2Transfer

	ranUe := ue.RanUe[models.AccessType_NON_3_GPP_ACCESS]
	err = ranUe.Remove()
	assert.True(t, err == nil)
	ue.DetachRanUe(models.AccessType_NON_3_GPP_ACCESS)

	ue.RegistrationArea[models.AccessType__3_GPP_ACCESS] = []models.Tai{
		{
			PlmnId: &models.PlmnId{
				Mcc: "208",
				Mnc: "93",
			},
			Tac: "000001",
		},
	}
	configuration := Namf_Communication.NewConfiguration()
	configuration.SetBasePath("https://localhost:29518")
	client := Namf_Communication.NewAPIClient(configuration)

	var n1N2MessageTransferRequest models.N1N2MessageTransferRequest
	n1N2MessageTransferRequest.BinaryDataN2Information = []byte{0x00}
	n1N2MessageTransferRequest.JsonData = TestComm.ConsumerAMFN1N2MessageTransferRequsetTable[TestComm.PDU_SETUP_REQ_11]
	sendN1N2Transfer(client, ue.Supi, n1N2MessageTransferRequest)

	time.Sleep(200 * time.Millisecond)
	// Trigger By Network (Service Accept)
	nasPdu = nasTestpacket.GetServiceRequest(nasMessage.ServiceTypeMobileTerminatedServices)

	amf_nas.HandleNAS(ue.RanUe[models.AccessType__3_GPP_ACCESS], ngapType.ProcedureCodeUplinkNASTransport, nasPdu)
	time.Sleep(100 * time.Millisecond)

	// Trigger By Ue (Service Reject)
	ue.MacFailed = true
	nasPdu = nasTestpacket.GetServiceRequest(nasMessage.ServiceTypeData)

	amf_nas.HandleNAS(ue.RanUe[models.AccessType__3_GPP_ACCESS], ngapType.ProcedureCodeUplinkNASTransport, nasPdu)

	time.Sleep(100 * time.Millisecond)
	TestAmf.Conn.Close()
	TestAmf.Conn2.Close()
}

func TestAuthenticationResponse5gAka(t *testing.T) {

	TestAmf.AmfInit()
	TestAmf.SctpConnectToServer(models.AccessType__3_GPP_ACCESS)

	time.Sleep(10 * time.Millisecond)
	ue := TestAmf.TestAmf.UePool["imsi-2089300007487"]
	err := ue.Sm[models.AccessType__3_GPP_ACCESS].Transfer(gmm_state.AUTHENTICATION, nil)
	assert.True(t, err == nil)

	rand := "0123456789abcdef0123456789abcdef"
	p0 := rand
	p1 := TestUEAuth.TestUe5gAuthTable[TestUEAuth.SUCCESS_CASE].XresStar
	concat, _ := hex.DecodeString(p0 + p1)
	hXResStarBytes := sha256.Sum256(concat)
	hXResStar := hex.EncodeToString(hXResStarBytes[:])

	ue.AuthenticationCtx = &models.UeAuthenticationCtx{
		AuthType: models.AuthType__5_G_AKA,
		Var5gAuthData: map[string]interface{}{
			"rand":      rand,
			"hxresStar": hXResStar,
		},
		Links: map[string]models.LinksValueSchema{
			"link": {
				Href: "https://localhost:29509/nausf-auth/v1/ue-authentications/imsi-2089300007487/5g-aka-confirmation",
			},
		},
	}

	t.Run("Success_Case", func(t *testing.T) {
		ausf_context.Init()
		ausfUe := ausf_context.NewAusfUeContext("imsi-2089300007487")
		ausfUe.ServingNetworkName = "5G:mnc093.mcc208.3gppnetwork.org"
		ausfUe.XresStar = TestUEAuth.TestUe5gAuthTable[TestUEAuth.SUCCESS_CASE].XresStar
		ausfUe.AuthStatus = models.AuthResult_ONGOING
		ausf_context.AddAusfUeContextToPool(ausfUe)

		resStar, _ := hex.DecodeString(TestUEAuth.TestUe5gAuthTable[TestUEAuth.SUCCESS_CASE].ResStar)
		nasPdu := nasTestpacket.GetAuthenticationResponse(resStar, "")
		amf_nas.HandleNAS(ue.RanUe[models.AccessType__3_GPP_ACCESS], ngapType.ProcedureCodeUplinkNASTransport, nasPdu)
	})

	t.Run("Failure_Case", func(t *testing.T) {
		ausf_context.Init()
		ausfUe := ausf_context.NewAusfUeContext("imsi-2089300007487")
		ausfUe.ServingNetworkName = "5G:mnc093.mcc208.3gppnetwork.org"
		ausfUe.XresStar = TestUEAuth.TestUe5gAuthTable[TestUEAuth.FAILURE_CASE].XresStar
		ausfUe.AuthStatus = models.AuthResult_ONGOING
		ausf_context.AddAusfUeContextToPool(ausfUe)

		err = ue.Sm[models.AccessType__3_GPP_ACCESS].Transfer(gmm_state.AUTHENTICATION, nil)
		assert.True(t, err == nil)
		resStar, _ := hex.DecodeString(TestUEAuth.TestUe5gAuthTable[TestUEAuth.FAILURE_CASE].ResStar)
		nasPdu := nasTestpacket.GetAuthenticationResponse(resStar, "")
		amf_nas.HandleNAS(ue.RanUe[models.AccessType__3_GPP_ACCESS], ngapType.ProcedureCodeUplinkNASTransport, nasPdu)
	})
	TestAmf.Conn.Close()

}

func buildEapPkt(testCase string) (eapMsgStr string) {
	var eapPkt radius.EapPacket

	eapPkt.Code = radius.EapCode(radius.EapCodeResponse)
	eapPkt.Type = radius.EapType(50) // accroding to RFC5448 6.1
	eapPkt.Identifier = 0x01
	atRes, _ := ausf_producer.EapEncodeAttribute("AT_RES", TestUEAuth.TestUeEapAuthTable[testCase].Res)
	atMAC, _ := ausf_producer.EapEncodeAttribute("AT_MAC", "")

	dataArrayBeforeMAC := atRes + atMAC
	eapPkt.Data = []byte(dataArrayBeforeMAC)
	encodedPktBeforeMAC := eapPkt.Encode()

	MACvalue := ausf_producer.CalculateAtMAC([]byte(TestUEAuth.TestUeEapAuthTable[testCase].K_aut), encodedPktBeforeMAC)

	atMacNum := fmt.Sprintf("%02x", ausf_context.AT_MAC_ATTRIBUTE)
	atMACfirstRow, _ := hex.DecodeString(atMacNum + "05" + "0000")
	wholeAtMAC := append(atMACfirstRow, MACvalue...)

	atMAC = string(wholeAtMAC)
	dataArrayAfterMAC := atRes + atMAC

	eapPkt.Data = []byte(dataArrayAfterMAC)
	encodedPktAfterMAC := eapPkt.Encode()

	eapMsg := nasType.EAPMessage{
		Len:    uint16(len(encodedPktAfterMAC)),
		Buffer: encodedPktAfterMAC,
	}

	eapMsgStr = base64.StdEncoding.EncodeToString(eapMsg.GetEAPMessage())
	return
}

func TestAuthenticationResponseEapAkaPrime(t *testing.T) {

	TestAmf.AmfInit()
	TestAmf.SctpConnectToServer(models.AccessType__3_GPP_ACCESS)

	time.Sleep(10 * time.Millisecond)
	ue := TestAmf.TestAmf.UePool["imsi-2089300007487"]
	err := ue.Sm[models.AccessType__3_GPP_ACCESS].Transfer(gmm_state.AUTHENTICATION, nil)
	assert.True(t, err == nil)

	ue.AuthenticationCtx = &models.UeAuthenticationCtx{
		AuthType: models.AuthType_EAP_AKA_PRIME,
		Links: map[string]models.LinksValueSchema{
			"link": {
				Href: "https://localhost:29509/nausf-auth/v1/ue-authentications/imsi-2089300007487/eap-session",
			},
		},
	}

	t.Run("Success_Case", func(t *testing.T) {
		ausf_context.Init()
		ausfUe := ausf_context.NewAusfUeContext("imsi-2089300007487")
		ausfUe.ServingNetworkName = "5G:mnc093.mcc208.3gppnetwork.org"
		ausfUe.XRES = TestUEAuth.TestUeEapAuthTable[TestUEAuth.SUCCESS_CASE].Xres
		ausfUe.K_aut = TestUEAuth.TestUeEapAuthTable[TestUEAuth.SUCCESS_CASE].K_aut
		ausfUe.AuthStatus = models.AuthResult_ONGOING
		ausf_context.AddAusfUeContextToPool(ausfUe)

		eapMsg := buildEapPkt(TestUEAuth.SUCCESS_CASE)
		nasPdu := nasTestpacket.GetAuthenticationResponse(nil, eapMsg)
		amf_nas.HandleNAS(ue.RanUe[models.AccessType__3_GPP_ACCESS], ngapType.ProcedureCodeUplinkNASTransport, nasPdu)
	})

	t.Run("Failure_Case", func(t *testing.T) {
		ausf_context.Init()
		ausfUe := ausf_context.NewAusfUeContext("imsi-2089300007487")
		ausfUe.ServingNetworkName = "5G:mnc093.mcc208.3gppnetwork.org"
		ausfUe.XRES = TestUEAuth.TestUeEapAuthTable[TestUEAuth.FAILURE_CASE].Xres
		ausfUe.K_aut = TestUEAuth.TestUeEapAuthTable[TestUEAuth.FAILURE_CASE].K_aut
		ausfUe.AuthStatus = models.AuthResult_ONGOING
		ausf_context.AddAusfUeContextToPool(ausfUe)

		err = ue.Sm[models.AccessType__3_GPP_ACCESS].Transfer(gmm_state.AUTHENTICATION, nil)
		assert.True(t, err == nil)
		eapMsg := buildEapPkt(TestUEAuth.FAILURE_CASE)
		nasPdu := nasTestpacket.GetAuthenticationResponse(nil, eapMsg)
		amf_nas.HandleNAS(ue.RanUe[models.AccessType__3_GPP_ACCESS], ngapType.ProcedureCodeUplinkNASTransport, nasPdu)
	})
	TestAmf.Conn.Close()
}

func TestAuthenticationFailure(t *testing.T) {

	ausf_context.Init()
	ausfUe := ausf_context.NewAusfUeContext("imsi-2089300007487")
	ausfUe.ServingNetworkName = "5G:mnc093.mcc208.3gppnetwork.org"
	ausfUe.XresStar = TestUEAuth.TestUe5gAuthTable[TestUEAuth.SUCCESS_CASE].XresStar
	ausfUe.AuthStatus = models.AuthResult_ONGOING
	ausf_context.AddAusfUeContextToPool(ausfUe)

	TestAmf.AmfInit()
	TestAmf.SctpConnectToServer(models.AccessType__3_GPP_ACCESS)
	time.Sleep(10 * time.Millisecond)
	ue := TestAmf.TestAmf.UePool["imsi-2089300007487"]

	ue.AuthenticationCtx = &models.UeAuthenticationCtx{
		AuthType: models.AuthType__5_G_AKA,
	}

	t.Run("Cause5GMMMACFailure", func(t *testing.T) {
		err := ue.Sm[models.AccessType__3_GPP_ACCESS].Transfer(gmm_state.AUTHENTICATION, nil)
		assert.True(t, err == nil)
		nasPdu := nasTestpacket.GetAuthenticationFailure(nasMessage.Cause5GMMMACFailure, nil)
		amf_nas.HandleNAS(ue.RanUe[models.AccessType__3_GPP_ACCESS], ngapType.ProcedureCodeUplinkNASTransport, nasPdu)
	})

	t.Run("Cause5GMMngKSIAlreadyInUse", func(t *testing.T) {
		err := ue.Sm[models.AccessType__3_GPP_ACCESS].Transfer(gmm_state.AUTHENTICATION, nil)
		assert.True(t, err == nil)
		nasPdu := nasTestpacket.GetAuthenticationFailure(nasMessage.Cause5GMMngKSIAlreadyInUse, nil)
		amf_nas.HandleNAS(ue.RanUe[models.AccessType__3_GPP_ACCESS], ngapType.ProcedureCodeUplinkNASTransport, nasPdu)
	})

	t.Run("Cause5GMMSynchFailure", func(t *testing.T) {
		err := ue.Sm[models.AccessType__3_GPP_ACCESS].Transfer(gmm_state.AUTHENTICATION, nil)
		assert.True(t, err == nil)
		failureParam := []uint8{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x10, 0x11, 0x12, 0x13, 0x14}
		nasPdu := nasTestpacket.GetAuthenticationFailure(nasMessage.Cause5GMMSynchFailure, failureParam)
		amf_nas.HandleNAS(ue.RanUe[models.AccessType__3_GPP_ACCESS], ngapType.ProcedureCodeUplinkNASTransport, nasPdu)
	})
	TestAmf.Conn.Close()

}

func TestRegistrationComplete(t *testing.T) {

	TestAmf.AmfInit()
	TestAmf.SctpConnectToServer(models.AccessType__3_GPP_ACCESS)

	time.Sleep(10 * time.Millisecond)
	ue := TestAmf.TestAmf.UePool["imsi-2089300007487"]
	err := ue.Sm[models.AccessType__3_GPP_ACCESS].Transfer(gmm_state.INITIAL_CONTEXT_SETUP, nil)
	assert.True(t, err == nil)

	ue.RegistrationRequest = nasMessage.NewRegistrationRequest(0)
	nasPdu := nasTestpacket.GetRegistrationComplete(nil)
	amf_nas.HandleNAS(ue.RanUe[models.AccessType__3_GPP_ACCESS], ngapType.ProcedureCodeUplinkNASTransport, nasPdu)
	TestAmf.Conn.Close()

}

func TestSecurityModeComplete(t *testing.T) {

	TestAmf.AmfInit()
	TestAmf.UeAttach(models.AccessType__3_GPP_ACCESS)
	TestAmf.SctpConnectToServer(models.AccessType__3_GPP_ACCESS)

	time.Sleep(10 * time.Millisecond)
	ue := TestAmf.TestAmf.UePool["imsi-2089300007487"]

	ranUe := ue.RanUe[models.AccessType__3_GPP_ACCESS]
	ranUe.Tai = models.Tai{
		PlmnId: &models.PlmnId{
			Mcc: "208",
			Mnc: "93",
		},
		Tac: "000001",
	}

	err := ue.Sm[models.AccessType__3_GPP_ACCESS].Transfer(gmm_state.SECURITY_MODE, nil)
	assert.True(t, err == nil)
	nasPdu := nasTestpacket.GetSecurityModeComplete()
	amf_nas.HandleNAS(ue.RanUe[models.AccessType__3_GPP_ACCESS], ngapType.ProcedureCodeUplinkNASTransport, nasPdu)
	TestAmf.Conn.Close()

}

func TestSecurityModeReject(t *testing.T) {

	TestAmf.AmfInit()
	TestAmf.SctpConnectToServer(models.AccessType__3_GPP_ACCESS)

	time.Sleep(10 * time.Millisecond)
	ue := TestAmf.TestAmf.UePool["imsi-2089300007487"]

	err := ue.Sm[models.AccessType__3_GPP_ACCESS].Transfer(gmm_state.SECURITY_MODE, nil)
	assert.True(t, err == nil)
	nasPdu := nasTestpacket.GetSecurityModeReject(nasMessage.Cause5GMMSecurityModeRejectedUnspecified)
	amf_nas.HandleNAS(ue.RanUe[models.AccessType__3_GPP_ACCESS], ngapType.ProcedureCodeUplinkNASTransport, nasPdu)
	TestAmf.Conn.Close()

}

func TestDeregistrationRequestUEOriginating(t *testing.T) {

	TestAmf.AmfInit()
	TestAmf.SctpConnectToServer(models.AccessType__3_GPP_ACCESS)

	time.Sleep(10 * time.Millisecond)
	ue := TestAmf.TestAmf.UePool["imsi-2089300007487"]

	err := ue.Sm[models.AccessType__3_GPP_ACCESS].Transfer(gmm_state.REGISTERED, nil)
	assert.True(t, err == nil)

	// TODO: modify test data if needed
	mobileIdentity5GS := nasType.MobileIdentity5GS{
		Len:    11, // 5g-guti
		Buffer: []uint8{0x02, 0x02, 0xf8, 0x39, 0xca, 0xfe, 0x00, 0x00, 0x00, 0x00, 0x01},
	}
	nasPdu := nasTestpacket.GetDeregistrationRequest(nasMessage.AccessType3GPP, 1, 0x04, mobileIdentity5GS)
	amf_nas.HandleNAS(ue.RanUe[models.AccessType__3_GPP_ACCESS], ngapType.ProcedureCodeUplinkNASTransport, nasPdu)
	TestAmf.Conn.Close()

}

func TestDeregistrationAccpetUETerminated(t *testing.T) {

	TestAmf.AmfInit()
	TestAmf.SctpConnectToServer(models.AccessType__3_GPP_ACCESS)

	time.Sleep(10 * time.Millisecond)
	ue := TestAmf.TestAmf.UePool["imsi-2089300007487"]

	err := ue.Sm[models.AccessType__3_GPP_ACCESS].Transfer(gmm_state.REGISTERED, nil)
	assert.True(t, err == nil)
	nasPdu := nasTestpacket.GetDeregistrationAccept()
	amf_nas.HandleNAS(ue.RanUe[models.AccessType__3_GPP_ACCESS], ngapType.ProcedureCodeUplinkNASTransport, nasPdu)
	TestAmf.Conn.Close()
}

func TestStatus5GMM(t *testing.T) {

	TestAmf.AmfInit()
	TestAmf.SctpConnectToServer(models.AccessType__3_GPP_ACCESS)

	time.Sleep(10 * time.Millisecond)
	ue := TestAmf.TestAmf.UePool["imsi-2089300007487"]

	nasPdu := nasTestpacket.GetStatus5GMM(nasMessage.Cause5GMMIllegalUE)

	amf_nas.HandleNAS(ue.RanUe[models.AccessType__3_GPP_ACCESS], ngapType.ProcedureCodeUplinkNASTransport, nasPdu)
	TestAmf.Conn.Close()
}

func TestStatus5GSM(t *testing.T) {
	time.Sleep(200 * time.Millisecond)
	MongoDBLibrary.RestfulAPIDeleteMany("NfProfile", bson.M{})

	go func() {
		router := Communication.NewRouter()
		server, err := http2_util.NewServer(":29518", TestAmf.AmfLogPath, router)
		if err == nil && server != nil {
			err = server.ListenAndServeTLS(TestAmf.AmfPemPath, TestAmf.AmfKeyPath)
			assert.True(t, err == nil, err.Error())
		}
	}()

	// smf register to nrf
	uuid, profile := TestAmf.BuildSmfNfProfile()
	uri, err := amf_consumer.SendRegisterNFInstance("https://localhost:29510", uuid, profile)
	if err != nil {
		t.Error(err.Error())
	} else {
		TestAmf.Config.Dump(uri)
	}

	TestAmf.AmfInit()
	TestAmf.SctpConnectToServer(models.AccessType__3_GPP_ACCESS)
	time.Sleep(10 * time.Millisecond)
	ue := TestAmf.TestAmf.UePool["imsi-2089300007487"]
	sNssai := models.Snssai{
		Sst: 1,
		Sd:  "010203",
	}
	ue.PlmnId.Mcc = "208"
	ue.PlmnId.Mnc = "93"
	// InitialRequest with known dnn and Snssai (success)
	nasPdu := nasTestpacket.GetUlNasTransport_PduSessionEstablishmentRequest(10, nasMessage.ULNASTransportRequestTypeInitialRequest, "internet", &sNssai)
	amf_nas.HandleNAS(ue.RanUe[models.AccessType__3_GPP_ACCESS], ngapType.ProcedureCodeUplinkNASTransport, nasPdu)

	nasPdu = nasTestpacket.GetUlNasTransport_Status5GSM(10, nasMessage.Cause5GSMInvalidPDUSessionIdentity)

	amf_nas.HandleNAS(ue.RanUe[models.AccessType__3_GPP_ACCESS], ngapType.ProcedureCodeUplinkNASTransport, nasPdu)
	TestAmf.Conn.Close()
}
