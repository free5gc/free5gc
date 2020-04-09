package amf_consumer_test

import (
	"github.com/stretchr/testify/assert"
	"free5gc/lib/CommonConsumerTestData/AMF/TestAmf"
	"free5gc/lib/CommonConsumerTestData/AMF/TestComm"
	"free5gc/lib/http2_util"
	"free5gc/lib/nas/nasMessage"
	"free5gc/lib/nas/nasTestpacket"
	"free5gc/lib/nas/nasType"
	"free5gc/lib/ngap/ngapType"
	"free5gc/lib/openapi/models"
	"free5gc/src/amf/Communication"
	"free5gc/src/amf/amf_consumer"
	"free5gc/src/amf/amf_context"
	"free5gc/src/amf/amf_handler"
	"free5gc/src/amf/gmm"
	"testing"
	"time"
)

func TestCreateUEContextRequest(t *testing.T) {
	if len(TestAmf.TestAmf.UePool) == 0 {
		go func() {
			router := Communication.NewRouter()
			server, err := http2_util.NewServer(":29518", TestAmf.AmfLogPath, router)
			if err == nil && server != nil {
				err = server.ListenAndServeTLS(TestAmf.AmfPemPath, TestAmf.AmfKeyPath)
			}
			assert.True(t, err == nil)
		}()

		go amf_handler.Handle()
		TestAmf.AmfInit()
		time.Sleep(100 * time.Millisecond)
	}

	/* init ue info*/
	ue := TestAmf.TestAmf.UePool["imsi-2089300007487"]

	ueContextCreateData := TestComm.ConsumerAMFCreateUEContextRequsetTable[TestComm.CreateUEContext201]
	ue.TargetAmfUri = "https://localhost:29518"
	ueContextCreatedData, problemDetails, err := amf_consumer.CreateUEContextRequest(ue, *ueContextCreateData.JsonData)
	if err != nil {
		t.Error(err)
	} else if problemDetails != nil {
		t.Errorf("Create Ue Context Request Failed: %+v", problemDetails)
	} else {
		TestAmf.Config.Dump(ueContextCreatedData)
	}
}

func TestReleaseUEContextRequest(t *testing.T) {
	if len(TestAmf.TestAmf.UePool) == 0 {
		TestCreateUEContextRequest(t)
	}

	/* init ue info*/
	self := amf_context.AMF_Self()
	supi := "imsi-0010202"
	ue := self.NewAmfUe(supi)
	if err := gmm.InitAmfUeSm(ue); err != nil {
		t.Errorf("InitAmfUeSm error: %v", err)
	}
	ue.Supi = "imsi-111222"

	ue = TestAmf.TestAmf.UePool["imsi-2089300007487"]
	ngapCause := models.NgApCause{
		Group: int32(ngapType.CausePresentProtocol),
		Value: int32(ngapType.CauseProtocolPresentUnspecified),
	}
	ue.TargetAmfUri = "https://localhost:29518"
	problemDetails, err := amf_consumer.ReleaseUEContextRequest(ue, ngapCause)
	if err != nil {
		t.Error(err)
	} else if problemDetails != nil {
		t.Errorf("Release Ue Context Request Failed: %+v", problemDetails)
	}
}

func TestUEContextTransferRequest(t *testing.T) {
	if len(TestAmf.TestAmf.UePool) == 0 {
		TestCreateUEContextRequest(t)
	}

	/* init ue info*/
	self := amf_context.AMF_Self()
	supi := "imsi-0010202"
	ue := self.NewAmfUe(supi)
	if err := gmm.InitAmfUeSm(ue); err != nil {
		t.Errorf("InitAmfUeSm error: %v", err.Error())
	}
	mobileIdentity5GS := nasType.MobileIdentity5GS{
		Len:    12, // suci
		Buffer: []uint8{0x01, 0x02, 0xf8, 0x39, 0xf0, 0xff, 0x00, 0x00, 0x00, 0x00, 0x47, 0x78},
	}
	registrationRequest := nasTestpacket.GetRegistrationRequestWith5GMM(nasMessage.RegistrationType5GSInitialRegistration, mobileIdentity5GS, nil, nil)
	ue.RegistrationRequest = nasMessage.NewRegistrationRequest(0)
	ue.RegistrationRequest.DecodeRegistrationRequest(&registrationRequest)
	ue.TargetAmfUri = "https://localhost:29518"
	ueContextTransferRspData, problemDetails, err := amf_consumer.UEContextTransferRequest(ue, models.AccessType__3_GPP_ACCESS, models.TransferReason_INIT_REG)
	if err != nil {
		t.Error(err)
	} else if problemDetails != nil {
		t.Errorf("Ue Context Transfer Request Failed: %+v", problemDetails)
	} else {
		TestAmf.Config.Dump(ueContextTransferRspData)
	}
}

func TestRegistrationStatusUpdate(t *testing.T) {
	if len(TestAmf.TestAmf.UePool) == 0 {
		TestCreateUEContextRequest(t)
	}

	ue := TestAmf.TestAmf.UePool["imsi-2089300007487"]
	ueRegStatusUpdateReqData := TestComm.ConsumerRegistrationStatusUpdateTable[TestComm.RegistrationStatusUpdate200]
	ue.TargetAmfUri = "https://localhost:29518"
	regStatusTransferComplete, problemDetails, err := amf_consumer.RegistrationStatusUpdate(ue, ueRegStatusUpdateReqData)
	if err != nil {
		t.Error(err)
	} else if problemDetails != nil {
		t.Errorf("Registration Status Update Failed: %+v", problemDetails)
	} else {
		TestAmf.Config.Dump(regStatusTransferComplete)
	}
}
