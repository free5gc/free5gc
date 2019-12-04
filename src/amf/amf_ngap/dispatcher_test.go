package amf_ngap_test

import (
	"encoding/hex"
	"flag"
	"github.com/gin-gonic/gin"
	"github.com/ishidawataru/sctp"
	"github.com/mohae/deepcopy"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"free5gc/lib/CommonConsumerTestData/AMF/TestAmf"
	"free5gc/lib/CommonConsumerTestData/SMF/TestPDUSession"
	"free5gc/lib/CommonConsumerTestData/UDM/TestGenAuthData"
	"free5gc/lib/aper"
	"free5gc/lib/http2_util"
	"free5gc/lib/nas/nasMessage"
	"free5gc/lib/nas/nasTestpacket"
	"free5gc/lib/nas/nasType"
	"free5gc/lib/ngap"
	"free5gc/lib/openapi/models"
	"free5gc/lib/path_util"
	"free5gc/src/amf/amf_consumer"
	"free5gc/src/amf/amf_context"
	"free5gc/src/amf/amf_handler"
	"free5gc/src/amf/amf_ngap"
	"free5gc/src/amf/amf_ngap/ngap_message"
	"free5gc/src/amf/logger"
	Nausf_UEAU "free5gc/src/ausf/UEAuthentication"
	"free5gc/src/ausf/ausf_context"
	"free5gc/src/ausf/ausf_handler"
	"free5gc/src/nrf/nrf_service"
	"free5gc/src/smf/smf_service"
	"free5gc/src/test/ngapTestpacket"
	Nudm_UEAU "free5gc/src/udm/UEAuthentication"
	"free5gc/src/udm/udm_handler"
	"net/http"
	"testing"
	"time"
)

func udrInit() {
	go func() { // fake udr server
		router := gin.Default()

		router.GET("/nudr-dr/v1/subscription-data/:ueId/authentication-data/authentication-subscription", func(c *gin.Context) {
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

		udrLogPath := path_util.Gofree5gcPath("free5gc/udrsslkey.log")
		udrPemPath := path_util.Gofree5gcPath("free5gc/support/TLS/udr.pem")
		udrKeyPath := path_util.Gofree5gcPath("free5gc/support/TLS/udr.key")

		server, err := http2_util.NewServer(":29504", udrLogPath, router)
		if err == nil && server != nil {
			logger.InitLog.Infoln(server.ListenAndServeTLS(udrPemPath, udrKeyPath))
		}
	}()
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
func smfInit() {
	flags := flag.FlagSet{}
	c := cli.NewContext(nil, &flags, nil)
	smf := &smf_service.SMF{}
	smf.Initialize(c)
	go smf.Start()
	time.Sleep(100 * time.Millisecond)
}
func nrfInit() {
	flags := flag.FlagSet{}
	c := cli.NewContext(nil, &flags, nil)
	nrf := &nrf_service.NRF{}
	nrf.Initialize(c)
	go nrf.Start()
	time.Sleep(100 * time.Millisecond)
}
func init() {
	logger.SetLogLevel(logrus.TraceLevel)
	logger.SetReportCaller(false)

	go amf_handler.Handle()
	nrfInit()
	smfInit()
	udrInit()
	udmInit()
	ausfInit()
	time.Sleep(100 * time.Millisecond)

	TestAmf.SctpSever()

}

// func sctpConnectToServer(ran *amf_context.AmfRan) {
// 	ipStr := "127.0.0.1"
// 	ips := []net.IPAddr{}
// 	if ip, err := net.ResolveIPAddr("ip", ipStr); err != nil {
// 		amf_ngap.Ngaplog.Errorf("Error resolving address '%s': %v", ipStr, err)
// 	} else {
// 		ips = append(ips, *ip)
// 	}
// 	addr1 := &sctp.SCTPAddr{
// 		IPAddrs: ips,
// 		Port:    38412,
// 	}
// 	amf_ngap.Ngaplog.Printf("raw TestAmf.Laddr: %+v\n", addr1.ToRawSockAddrBuf())

// 	var laddr *sctp.SCTPAddr
// 	conn, err := sctp.DialSCTP("sctp", laddr, addr1)

// 	if err != nil {
// 		amf_ngap.Ngaplog.Errorf("failed to dial: %v\n", err)
// 	}
// 	amf_ngap.Ngaplog.Printf("Dail LocalAddr: %s; RemoteAddr: %s", conn.LocalAddr(), conn.RemoteAddr())
// 	time.Sleep(time.Millisecond)

// 	ran.Conn = conn
// }

func TestDispatchNGSetupRequest(t *testing.T) {

	TestAmf.AmfInit()
	TestAmf.SctpConnectToServer(models.AccessType__3_GPP_ACCESS)
	time.Sleep(100 * time.Millisecond)

	message := ngapTestpacket.BuildNGSetupRequest()
	msg, err := ngap.Encoder(message)
	if err != nil {
		amf_ngap.Ngaplog.Errorln(err)
	}
	amf_ngap.Ngaplog.Warnln(TestAmf.Laddr.String())
	amf_ngap.Dispatch(TestAmf.Laddr.String(), msg)

	time.Sleep(10 * time.Millisecond)
	TestAmf.Conn.Close()
	// TestAmf.Config.Dump(ran)
}

func TestDispatchUplinkNasTransport(t *testing.T) {

	TestAmf.AmfInit()
	TestAmf.SctpConnectToServer(models.AccessType__3_GPP_ACCESS)
	time.Sleep(100 * time.Millisecond)

	ue := TestAmf.TestAmf.UePool["imsi-2089300007487"]

	ue.RanUe[models.AccessType__3_GPP_ACCESS].AmfUeNgapId = 1
	ue.RanUe[models.AccessType__3_GPP_ACCESS].RanUeNgapId = 2

	nasPdu := nasTestpacket.GetUlNasTransport_PduSessionEstablishmentRequest(1, nasMessage.ULNASTransportRequestTypeInitialRequest, "internet", &ue.AllowedNssai[models.AccessType__3_GPP_ACCESS][0])
	// amf_ngap.Ngaplog.Tracef("nas: %0x", nasPdu.Bytes())
	message := ngapTestpacket.BuildUplinkNasTransport(1, 2, nasPdu)
	msg, err := ngap.Encoder(message)
	if err != nil {
		amf_ngap.Ngaplog.Errorln(err)
	}
	amf_ngap.Dispatch(TestAmf.Laddr.String(), msg)
	TestAmf.Conn.Close()
}

func TestDispatchNGReset(t *testing.T) {

	TestAmf.AmfInit()
	TestAmf.SctpConnectToServer(models.AccessType__3_GPP_ACCESS)
	time.Sleep(100 * time.Millisecond)

	message := ngapTestpacket.BuildNGReset()
	msg, err := ngap.Encoder(message)
	if err != nil {
		amf_ngap.Ngaplog.Errorln(err)
	}
	amf_ngap.Dispatch(TestAmf.Laddr.String(), msg)
	time.Sleep(10 * time.Millisecond)
	TestAmf.Conn.Close()
}

func TestDispatchNGResetAcknowledge(t *testing.T) {

	TestAmf.AmfInit()
	TestAmf.UeAttach(models.AccessType__3_GPP_ACCESS)

	message := ngapTestpacket.BuildNGResetAcknowledge()
	msg, err := ngap.Encoder(message)
	if err != nil {
		amf_ngap.Ngaplog.Errorln(err)
	}
	amf_ngap.Dispatch(TestAmf.Laddr.String(), msg)
}

func TestDispatchUEContextReleaseComplete(t *testing.T) {

	TestAmf.AmfInit()
	TestAmf.UeAttach(models.AccessType__3_GPP_ACCESS)

	message := ngapTestpacket.BuildUEContextReleaseComplete(1, 2, nil)
	msg, err := ngap.Encoder(message)
	if err != nil {
		amf_ngap.Ngaplog.Errorln(err)
	}
	amf_ngap.Dispatch(TestAmf.Laddr.String(), msg)
}

func TestDispatchPDUSessionResourceReleaseResponse(t *testing.T) {

	TestAmf.AmfInit()
	TestAmf.UeAttach(models.AccessType__3_GPP_ACCESS)

	message := ngapTestpacket.BuildPDUSessionResourceReleaseResponse()
	msg, err := ngap.Encoder(message)
	if err != nil {
		amf_ngap.Ngaplog.Errorln(err)
	}
	amf_ngap.Dispatch(TestAmf.Laddr.String(), msg)
}

func TestDispatchUERadioCapabilityCheckResponse(t *testing.T) {

	TestAmf.AmfInit()
	TestAmf.UeAttach(models.AccessType__3_GPP_ACCESS)

	message := ngapTestpacket.BuildUERadioCapabilityCheckResponse()
	msg, err := ngap.Encoder(message)
	if err != nil {
		amf_ngap.Ngaplog.Errorln(err)
	}
	amf_ngap.Dispatch(TestAmf.Laddr.String(), msg)
}

func TestDispatchHandoverCancel(t *testing.T) {

	TestAmf.AmfInit()
	TestAmf.SctpConnectToServer(models.AccessType__3_GPP_ACCESS)

	ue := TestAmf.TestAmf.UePool["imsi-2089300007487"]
	targetUe := TestAmf.TestAmf.AmfRanPool[TestAmf.Laddr.String()].NewRanUe()
	ue.RanUe[models.AccessType__3_GPP_ACCESS].TargetUe = targetUe
	targetUe.AmfUe = ue

	message := ngapTestpacket.BuildHandoverCancel()
	msg, err := ngap.Encoder(message)
	if err != nil {
		amf_ngap.Ngaplog.Errorln(err)
	}
	amf_ngap.Dispatch(TestAmf.Laddr.String(), msg)
	time.Sleep(20 * time.Millisecond)
	TestAmf.Conn.Close()
}

func TestDispatchLocationReportingFailureIndication(t *testing.T) {

	TestAmf.AmfInit()
	TestAmf.UeAttach(models.AccessType__3_GPP_ACCESS)

	message := ngapTestpacket.BuildLocationReportingFailureIndication()
	msg, err := ngap.Encoder(message)
	if err != nil {
		amf_ngap.Ngaplog.Errorln(err)
	}
	amf_ngap.Dispatch(TestAmf.Laddr.String(), msg)
}

func TestDispatchInitialUEMessage(t *testing.T) {

	ausf_context.Init()
	TestAmf.AmfInit()
	TestAmf.SctpConnectToServer(models.AccessType__3_GPP_ACCESS)

	mobileIdentity5GS := nasType.MobileIdentity5GS{
		Len:    12, // suci
		Buffer: []uint8{0x01, 0x02, 0xf8, 0x39, 0xf0, 0x00, 0x00, 0x00, 0x00, 0x00, 0x47, 0x78},
	}

	requestedNSSAI := &nasType.RequestedNSSAI{
		Iei:    nasMessage.RegistrationRequestRequestedNSSAIType,
		Len:    5,
		Buffer: []uint8{0x04, 0x01, 0x01, 0x02, 0x03},
	}
	nasPdu := nasTestpacket.GetRegistrationRequest(nasMessage.RegistrationType5GSInitialRegistration, mobileIdentity5GS, requestedNSSAI, nil)

	message := ngapTestpacket.BuildInitialUEMessage(1, nasPdu, "")
	msg, err := ngap.Encoder(message)
	if err != nil {
		amf_ngap.Ngaplog.Errorln(err)
	}

	/* Wireshark test */
	ran := TestAmf.TestAmf.AmfRanPool[TestAmf.Laddr.String()]
	ngap_message.SendToRan(ran, msg)
	/* Wireshark test */

	amf_ngap.Dispatch(TestAmf.Laddr.String(), msg)
	TestAmf.Conn.Close()
}

func TestDispatchInitialUEMessageWithGuti(t *testing.T) {

	ausf_context.Init()
	TestAmf.AmfInit()
	TestAmf.SctpConnectToServer(models.AccessType__3_GPP_ACCESS)

	ue := TestAmf.TestAmf.UePool["imsi-2089300007487"]
	ue.Guti = "20893cafe0000000001"
	TestAmf.TestAmf.GutiPool[ue.Guti] = ue

	mobileIdentity5GS := nasType.MobileIdentity5GS{
		Len:    11, // 5g-guti
		Buffer: []uint8{0x02, 0x02, 0xf8, 0x39, 0xca, 0xfe, 0x00, 0x00, 0x00, 0x00, 0x01},
	}

	nasPdu := nasTestpacket.GetRegistrationRequest(nasMessage.RegistrationType5GSMobilityRegistrationUpdating, mobileIdentity5GS, nil, nil)

	message := ngapTestpacket.BuildInitialUEMessage(1, nasPdu, "fe0000000001")
	msg, err := ngap.Encoder(message)
	if err != nil {
		amf_ngap.Ngaplog.Errorln(err)
	}

	amf_ngap.Dispatch(TestAmf.Laddr.String(), msg)
	TestAmf.Conn.Close()
}

func TestDispatchPDUSessionResourceSetupResponse(t *testing.T) {

	// create test smContext
	smContextCreate()

	ran := TestAmf.TestAmf.AmfRanPool[TestAmf.Laddr.String()]
	ran.RanUeList[0].RanUeNgapId = 123

	message := ngapTestpacket.BuildPDUSessionResourceSetupResponse(1, 123, "10.200.200.1")
	msg, err := ngap.Encoder(message)
	if err != nil {
		amf_ngap.Ngaplog.Errorln(err)
	}
	amf_ngap.Dispatch(TestAmf.Laddr.String(), msg)
	time.Sleep(20 * time.Millisecond)
	TestAmf.Conn.Close()
}

func TestDispatchPDUSessionResourceModifyResponse(t *testing.T) {

	// create test smContext
	smContextCreate()

	ran := TestAmf.TestAmf.AmfRanPool[TestAmf.Laddr.String()]
	ran.RanUeList[0].AmfUeNgapId = 1
	ran.RanUeList[0].RanUeNgapId = 2
	message := ngapTestpacket.BuildPDUSessionResourceModifyResponse(1, 2)
	msg, err := ngap.Encoder(message)
	if err != nil {
		amf_ngap.Ngaplog.Errorln(err)
		t.Error("Encode testpacket failed")
	}
	amf_ngap.Dispatch(TestAmf.Laddr.String(), msg)
	time.Sleep(20 * time.Millisecond)
	TestAmf.Conn.Close()
}

func TestDispatchPDUSessionResourceNotify(t *testing.T) {

	TestAmf.AmfInit()
	TestAmf.UeAttach(models.AccessType__3_GPP_ACCESS)
	// ran := TestAmf.TestAmf.AmfRanPool[TestAmf.Laddr.String()]
	message := ngapTestpacket.BuildPDUSessionResourceNotify()
	msg, err := ngap.Encoder(message)
	if err != nil {
		amf_ngap.Ngaplog.Errorln(err)
	}
	amf_ngap.Dispatch(TestAmf.Laddr.String(), msg)

	// TestAmf.Config.Dump(ran)
}

func TestDispatchPDUSessionResourceModifyIndication(t *testing.T) {

	smContextCreate()

	ran := TestAmf.TestAmf.AmfRanPool[TestAmf.Laddr.String()]
	ran.RanUeList[0].AmfUeNgapId = 1
	ran.RanUeList[0].RanUeNgapId = 2
	message := ngapTestpacket.BuildPDUSessionResourceModifyIndication(1, 2)
	msg, err := ngap.Encoder(message)
	if err != nil {
		amf_ngap.Ngaplog.Errorln(err)
	}
	amf_ngap.Dispatch(TestAmf.Laddr.String(), msg)
	time.Sleep(20 * time.Millisecond)
	TestAmf.Conn.Close()
}

func TestDispatchInitialContextSetupResponse(t *testing.T) {

	TestAmf.AmfInit()
	TestAmf.UeAttach(models.AccessType__3_GPP_ACCESS)

	ran := TestAmf.TestAmf.AmfRanPool[TestAmf.Laddr.String()]
	ran.RanUeList[0].AmfUeNgapId = 1
	ran.RanUeList[0].RanUeNgapId = 2
	message := ngapTestpacket.BuildInitialContextSetupResponse(1, 2, "10.200.200.1", nil)
	msg, err := ngap.Encoder(message)
	if err != nil {
		amf_ngap.Ngaplog.Errorln(err)
	}
	amf_ngap.Dispatch(TestAmf.Laddr.String(), msg)
}

func TestDispatchInitialContextSetupFailure(t *testing.T) {

	TestAmf.AmfInit()
	TestAmf.UeAttach(models.AccessType__3_GPP_ACCESS)

	ran := TestAmf.TestAmf.AmfRanPool[TestAmf.Laddr.String()]
	ran.RanUeList[0].AmfUeNgapId = 1
	ran.RanUeList[0].RanUeNgapId = 2
	message := ngapTestpacket.BuildInitialContextSetupFailure(1, 2)
	msg, err := ngap.Encoder(message)
	if err != nil {
		amf_ngap.Ngaplog.Errorln(err)
	}
	amf_ngap.Dispatch(TestAmf.Laddr.String(), msg)
}

func TestDispatchUEContextReleaseRequest(t *testing.T) {

	smContextCreate()
	time.Sleep(10 * time.Millisecond)
	ran := TestAmf.TestAmf.AmfRanPool[TestAmf.Laddr.String()]

	ranUe := ran.RanUeList[0]
	ranUe.RanUeNgapId = 123
	message := ngapTestpacket.BuildUEContextReleaseRequest(1, 123, nil)
	msg, err := ngap.Encoder(message)
	if err != nil {
		amf_ngap.Ngaplog.Errorln(err)
	}
	amf_ngap.Dispatch(TestAmf.Laddr.String(), msg)
	time.Sleep(10 * time.Millisecond)
	TestAmf.Conn.Close()
}

func TestDispatchUEContextModificationResponse(t *testing.T) {

	TestAmf.AmfInit()
	TestAmf.UeAttach(models.AccessType__3_GPP_ACCESS)

	ran := TestAmf.TestAmf.AmfRanPool[TestAmf.Laddr.String()]
	ranUe := ran.RanUeList[0]
	ranUe.RanUeNgapId = 123

	message := ngapTestpacket.BuildUEContextModificationResponse(1, 123)
	msg, err := ngap.Encoder(message)
	if err != nil {
		amf_ngap.Ngaplog.Errorln(err)
	}
	amf_ngap.Dispatch(TestAmf.Laddr.String(), msg)
}

func TestDispatchUEContextModificationFailure(t *testing.T) {

	TestAmf.AmfInit()
	TestAmf.UeAttach(models.AccessType__3_GPP_ACCESS)

	ran := TestAmf.TestAmf.AmfRanPool[TestAmf.Laddr.String()]
	ranUe := ran.RanUeList[0]
	ranUe.RanUeNgapId = 123

	message := ngapTestpacket.BuildUEContextModificationFailure(1, 123)
	msg, err := ngap.Encoder(message)
	if err != nil {
		amf_ngap.Ngaplog.Errorln(err)
	}
	amf_ngap.Dispatch(TestAmf.Laddr.String(), msg)
}

func TestDispatchRRCInactiveTransitionReport(t *testing.T) {

	TestAmf.AmfInit()
	TestAmf.UeAttach(models.AccessType__3_GPP_ACCESS)

	// ran := TestAmf.TestAmf.AmfRanPool[TestAmf.Laddr.String()]
	message := ngapTestpacket.BuildRRCInactiveTransitionReport()
	msg, err := ngap.Encoder(message)
	if err != nil {
		amf_ngap.Ngaplog.Errorln(err)
	}
	amf_ngap.Dispatch(TestAmf.Laddr.String(), msg)

	// TestAmf.Config.Dump(ran)
}

func TestDispatchHandoverNotify(t *testing.T) {

	TestAmf.AmfInit()
	TestAmf.UeAttach(models.AccessType__3_GPP_ACCESS)

	message := ngapTestpacket.BuildHandoverNotify(0, 0)
	msg, err := ngap.Encoder(message)
	if err != nil {
		amf_ngap.Ngaplog.Errorln(err)
	}
	amf_ngap.Dispatch(TestAmf.Laddr.String(), msg)
}

func TestDispatchPathSwitchRequest(t *testing.T) {

	smContextCreate()

	ue := TestAmf.TestAmf.UePool["imsi-2089300007487"]
	ue.RanUe[models.AccessType__3_GPP_ACCESS].AmfUeNgapId = 1
	ue.RanUe[models.AccessType__3_GPP_ACCESS].RanUeNgapId = 2
	message := ngapTestpacket.BuildPathSwitchRequest(1, 2)
	msg, err := ngap.Encoder(message)
	if err != nil {
		amf_ngap.Ngaplog.Errorln(err)
	}
	amf_ngap.Dispatch(TestAmf.Laddr.String(), msg)
	time.Sleep(20 * time.Millisecond)
	TestAmf.Conn.Close()
}

func TestDispatchHandoverRequestAcknowledge(t *testing.T) {

	smContextCreate()
	ue := TestAmf.TestAmf.UePool["imsi-2089300007487"]
	ue.RanUe[models.AccessType__3_GPP_ACCESS].AmfUeNgapId = 1
	ue.RanUe[models.AccessType__3_GPP_ACCESS].RanUeNgapId = 2

	targetUe := TestAmf.TestAmf.AmfRanPool[TestAmf.Laddr.String()].NewRanUe()
	targetUe.AmfUe = ue
	targetUe.SourceUe = ue.RanUe[models.AccessType__3_GPP_ACCESS]
	// sourceAmfUe := TestAmf.TestAmf.NewAmfUe("imsi-2089300007488")
	// sourceAmfUe.AttachRanUe(sourceUe)

	message := ngapTestpacket.BuildHandoverRequestAcknowledge(2, 2)
	msg, err := ngap.Encoder(message)
	if err != nil {
		amf_ngap.Ngaplog.Errorln(err)
	}
	amf_ngap.Dispatch(TestAmf.Laddr.String(), msg)
	time.Sleep(20 * time.Millisecond)
	TestAmf.Conn.Close()
}

func TestDispatchHandoverFailure(t *testing.T) {

	TestAmf.AmfInit()
	TestAmf.SctpConnectToServer(models.AccessType__3_GPP_ACCESS)
	time.Sleep(100 * time.Millisecond)

	ue := TestAmf.TestAmf.UePool["imsi-2089300007487"]
	ue.RanUe[models.AccessType__3_GPP_ACCESS].AmfUeNgapId = 1
	ue.RanUe[models.AccessType__3_GPP_ACCESS].RanUeNgapId = 2

	sourceUe := TestAmf.TestAmf.AmfRanPool[TestAmf.Laddr.String()].NewRanUe()
	ue.RanUe[models.AccessType__3_GPP_ACCESS].SourceUe = sourceUe
	sourceUe.AmfUe = ue

	message := ngapTestpacket.BuildHandoverFailure(1)
	msg, err := ngap.Encoder(message)
	if err != nil {
		amf_ngap.Ngaplog.Errorln(err)
	}
	amf_ngap.Dispatch(TestAmf.Laddr.String(), msg)
	time.Sleep(10 * time.Millisecond)
	TestAmf.Conn.Close()
}

func TestDispatchUplinkRanStatusTransfer(t *testing.T) {

	TestAmf.AmfInit()
	TestAmf.UeAttach(models.AccessType__3_GPP_ACCESS)

	ran := TestAmf.TestAmf.AmfRanPool[TestAmf.Laddr.String()]
	ran.RanUeList[0].RanUeNgapId = 123
	message := ngapTestpacket.BuildUplinkRanStatusTransfer(122, 123)
	msg, err := ngap.Encoder(message)
	if err != nil {
		amf_ngap.Ngaplog.Errorln(err)
	}
	amf_ngap.Dispatch(TestAmf.Laddr.String(), msg)

}

func TestDispatchNasNonDeliveryIndication(t *testing.T) {

	TestAmf.AmfInit()
	TestAmf.UeAttach(models.AccessType__3_GPP_ACCESS)

	ran := TestAmf.TestAmf.AmfRanPool[TestAmf.Laddr.String()]
	ran.RanUeList[0].RanUeNgapId = 123
	message := ngapTestpacket.BuildNasNonDeliveryIndication(1, 123, aper.OctetString("\x01\x02\x03"))
	msg, err := ngap.Encoder(message)
	if err != nil {
		amf_ngap.Ngaplog.Errorln(err)
	}
	amf_ngap.Dispatch(TestAmf.Laddr.String(), msg)
}

func TestDispatchRanConfigurationUpdate(t *testing.T) {

	TestAmf.AmfInit()
	TestAmf.SctpConnectToServer(models.AccessType__3_GPP_ACCESS)
	time.Sleep(100 * time.Millisecond)

	ran := TestAmf.TestAmf.AmfRanPool[TestAmf.Laddr.String()]

	ran.RanUeList[0].RanUeNgapId = 123
	message := ngapTestpacket.BuildRanConfigurationUpdate()
	msg, err := ngap.Encoder(message)
	if err != nil {
		amf_ngap.Ngaplog.Errorln(err)
	}
	amf_ngap.Dispatch(TestAmf.Laddr.String(), msg)
	time.Sleep(10 * time.Millisecond)
	TestAmf.Conn.Close()
}

func TestDispatchUplinkRanConfigurationTransfer(t *testing.T) {

	TestAmf.AmfInit()
	TestAmf.SctpConnectToServer(models.AccessType__3_GPP_ACCESS)
	time.Sleep(100 * time.Millisecond)
	ran := TestAmf.TestAmf.AmfRanPool[TestAmf.Laddr.String()]

	ran.RanPresent = 1
	ran.RanId = new(models.GlobalRanNodeId)
	ran.RanId.GNbId = new(models.GNbId)
	ran.RanId.GNbId.GNBValue = "414240"

	// ran := TestAmf.TestAmf.AmfRanPool[TestAmf.Laddr.String()]
	message := ngapTestpacket.BuildUplinkRanConfigurationTransfer()
	msg, err := ngap.Encoder(message)
	if err != nil {
		amf_ngap.Ngaplog.Errorln(err)
	}
	amf_ngap.Dispatch(TestAmf.Laddr.String(), msg)
	time.Sleep(10 * time.Millisecond)
	TestAmf.Conn.Close()

	// TestAmf.Config.Dump(ran)
}

func TestDispatchUplinkUEAssociatedNRPPATransport(t *testing.T) {

	TestAmf.AmfInit()
	TestAmf.UeAttach(models.AccessType__3_GPP_ACCESS)

	// ran := TestAmf.TestAmf.AmfRanPool[TestAmf.Laddr.String()]
	message := ngapTestpacket.BuildUplinkUEAssociatedNRPPATransport()
	msg, err := ngap.Encoder(message)
	if err != nil {
		amf_ngap.Ngaplog.Errorln(err)
	}
	amf_ngap.Dispatch(TestAmf.Laddr.String(), msg)

	// TestAmf.Config.Dump(ran)
}

func TestDispatchUplinkNonUEAssociatedNRPPATransport(t *testing.T) {
	TestAmf.AmfInit()
	TestAmf.UeAttach(models.AccessType__3_GPP_ACCESS)

	ran := TestAmf.TestAmf.AmfRanPool[TestAmf.Laddr.String()]
	ran.RanUeList[0].RanUeNgapId = 123
	message := ngapTestpacket.BuildUplinkNonUEAssociatedNRPPATransport()
	msg, err := ngap.Encoder(message)
	if err != nil {
		amf_ngap.Ngaplog.Errorln(err)
	}
	amf_ngap.Dispatch(TestAmf.Laddr.String(), msg)

}

func TestDispatchLocationReport(t *testing.T) {

	TestAmf.AmfInit()
	TestAmf.UeAttach(models.AccessType__3_GPP_ACCESS)

	// ran := TestAmf.TestAmf.AmfRanPool[TestAmf.Laddr.String()]
	message := ngapTestpacket.BuildLocationReport()
	msg, err := ngap.Encoder(message)
	if err != nil {
		amf_ngap.Ngaplog.Errorln(err)
	}
	amf_ngap.Dispatch(TestAmf.Laddr.String(), msg)

	// TestAmf.Config.Dump(ran)
}

func TestDispatchUERadioCapabilityInfoIndication(t *testing.T) {

	TestAmf.AmfInit()
	TestAmf.UeAttach(models.AccessType__3_GPP_ACCESS)

	//ran := TestAmf.TestAmf.AmfRanPool[TestAmf.Laddr.String()]
	message := ngapTestpacket.BuildUERadioCapabilityInfoIndication()
	msg, err := ngap.Encoder(message)
	if err != nil {
		amf_ngap.Ngaplog.Errorln(err)
	}
	amf_ngap.Dispatch(TestAmf.Laddr.String(), msg)

	//TestAmf.Config.Dump(ran)
}

func TestDispatchAMFConfigurationUpdateFailure(t *testing.T) {

	TestAmf.AmfInit()
	TestAmf.UeAttach(models.AccessType__3_GPP_ACCESS)

	// ran := TestAmf.TestAmf.AmfRanPool[TestAmf.Laddr.String()]
	message := ngapTestpacket.BuildAMFConfigurationUpdateFailure()
	msg, err := ngap.Encoder(message)
	if err != nil {
		amf_ngap.Ngaplog.Errorln(err)
	}
	amf_ngap.Dispatch(TestAmf.Laddr.String(), msg)

	// TestAmf.Config.Dump(ran)
}

func TestDispatchAMFConfigurationUpdateAcknowledge(t *testing.T) {

	TestAmf.AmfInit()
	TestAmf.UeAttach(models.AccessType__3_GPP_ACCESS)

	// ran := TestAmf.TestAmf.AmfRanPool[TestAmf.Laddr.String()]
	message := ngapTestpacket.BuildAMFConfigurationUpdateAcknowledge()
	msg, err := ngap.Encoder(message)
	if err != nil {
		amf_ngap.Ngaplog.Errorln(err)
	}
	amf_ngap.Dispatch(TestAmf.Laddr.String(), msg)

	// TestAmf.Config.Dump(ran)
}

func TestDispatcherErrorIndication(t *testing.T) {

	TestAmf.AmfInit()
	TestAmf.UeAttach(models.AccessType__3_GPP_ACCESS)

	message := ngapTestpacket.BuildErrorIndication()
	msg, err := ngap.Encoder(message)
	if err != nil {
		amf_ngap.Ngaplog.Errorln(err)
	}

	amf_ngap.Dispatch(TestAmf.Laddr.String(), msg)
}

func TestDispatchHandoverRequired(t *testing.T) {

	smContextCreate()
	ue := TestAmf.TestAmf.UePool["imsi-2089300007487"]
	ue.SecurityContextAvailable = true
	ue.RanUe[models.AccessType__3_GPP_ACCESS].RanUeNgapId = 1

	ue.NCC = 5
	ue.NH, _ = hex.DecodeString("ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff")

	ue.SecurityCapabilities.NREncryptionAlgorithms = [2]byte{0xe0, 0x00}
	ue.SecurityCapabilities.NRIntegrityProtectionAlgorithms = [2]byte{0xe0, 0x00}
	ue.SecurityCapabilities.EUTRAEncryptionAlgorithms = [2]byte{0xe0, 0x00}
	ue.SecurityCapabilities.EUTRAIntegrityProtectionAlgorithms = [2]byte{0xe0, 0x00}
	// target ran
	conn, _ := sctp.DialSCTP("sctp", TestAmf.Laddr2, TestAmf.ServerAddr)
	time.Sleep(20 * time.Millisecond)

	ran := TestAmf.TestAmf.AmfRanPool[TestAmf.Laddr2.String()]
	ran.RanPresent = amf_context.RanPresentGNbId
	ran.RanId = new(models.GlobalRanNodeId)
	ran.RanId.GNbId = new(models.GNbId)
	ran.RanId.GNbId.GNBValue = "454647"

	message := ngapTestpacket.BuildHandoverRequired(1, 2, []byte{0x45, 0x46, 0x47}, []byte{0x01, 0x20})
	msg, err := ngap.Encoder(message)
	if err != nil {
		amf_ngap.Ngaplog.Errorln(err)
	}
	amf_ngap.Dispatch(TestAmf.Laddr.String(), msg)
	time.Sleep(20 * time.Millisecond)
	conn.Close()
	TestAmf.Conn.Close()

	//TestAmf.Config.Dump(ran)
}

func TestDispatchCellTrafficTrace(t *testing.T) {

	TestAmf.AmfInit()
	TestAmf.UeAttach(models.AccessType__3_GPP_ACCESS)

	ue := TestAmf.TestAmf.UePool["imsi-2089300007487"]
	ue.RanUe[models.AccessType__3_GPP_ACCESS].AmfUeNgapId = 1
	ue.RanUe[models.AccessType__3_GPP_ACCESS].RanUeNgapId = 2

	message := ngapTestpacket.BuildCellTrafficTrace(1, 2)
	msg, err := ngap.Encoder(message)
	if err != nil {
		amf_ngap.Ngaplog.Errorln(err)
	}
	amf_ngap.Dispatch(TestAmf.Laddr.String(), msg)
}

// Copy from TestSMContextCreate() in amf_consumer/sm_context_test.go
func smContextCreate() {

	TestAmf.AmfInit()
	TestAmf.SctpConnectToServer(models.AccessType__3_GPP_ACCESS)

	ue := TestAmf.TestAmf.UePool["imsi-2089300007487"]

	payload := TestPDUSession.GetEstablishmentRequestData(TestPDUSession.SERVICE_REQUEST)
	pduSession := models.PduSessionContext{
		PduSessionId: 10,
		Dnn:          "internet",
		SNssai: &models.Snssai{
			Sst: 1,
			Sd:  "020304",
		},
	}
	requestType := models.RequestType_INITIAL_REQUEST
	if anType := ue.GetAnType(); anType == "" {
		pduSession.AccessType = models.AccessType__3_GPP_ACCESS
	} else {
		pduSession.AccessType = anType
	}
	smContextCreateData := amf_consumer.BuildCreateSmContextRequest(ue, pduSession, requestType)
	// TODO: http://localhost:29502/ -> smfD smfUri which required from NRF
	smfUri := "https://localhost:29502"

	createPduSession(ue, &pduSession, smfUri, payload, smContextCreateData)

	pduSession2 := models.PduSessionContext{
		PduSessionId: 11,
		Dnn:          "internet",
		SNssai: &models.Snssai{
			Sst: 1,
			Sd:  "020304",
		},
	}
	requestType = models.RequestType_INITIAL_REQUEST
	if anType := ue.GetAnType(); anType == "" {
		pduSession2.AccessType = models.AccessType__3_GPP_ACCESS
	} else {
		pduSession2.AccessType = anType
	}
	smContextCreateData = amf_consumer.BuildCreateSmContextRequest(ue, pduSession2, requestType)

	createPduSession(ue, &pduSession2, smfUri, payload, smContextCreateData)

}

func createPduSession(ue *amf_context.AmfUe, pduSession *models.PduSessionContext, smfUri string, payload []byte, smContextCreateData models.SmContextCreateData) {

	response, smContextRef, _, _, err := amf_consumer.SendCreateSmContextRequest(ue, smfUri, payload, smContextCreateData)
	if response != nil {
		var smContext amf_context.SmContext
		pduSession.SmContextRef = smContextRef
		smContext.PduSessionContext = pduSession
		smContext.UserLocation = deepcopy.Copy(ue.Location).(models.UserLocation)
		smContext.SmfUri = smfUri
		// TODO: store SmfId
		// smContext.SmfId = ???
		ue.SmContextList[pduSession.PduSessionId] = &smContext
		// TODO: handle response(response N2SmInfo to RAN if exists)
	} else if err != nil {
		amf_ngap.Ngaplog.Errorf("[ERROR] " + err.Error())
	} else {
		// TODO: error handling
	}
}
