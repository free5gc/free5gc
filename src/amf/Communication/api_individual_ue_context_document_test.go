package Communication_test

import (
	"context"
	"github.com/stretchr/testify/assert"
	"free5gc/lib/CommonConsumerTestData/AMF/TestAmf"
	"free5gc/lib/CommonConsumerTestData/AMF/TestComm"
	Namf_Communication_Client "free5gc/lib/Namf_Communication"
	"free5gc/lib/http2_util"
	"free5gc/lib/openapi/common"
	"free5gc/lib/openapi/models"
	Namf_Communication_Server "free5gc/src/amf/Communication"
	"free5gc/src/amf/amf_context"
	"free5gc/src/amf/amf_handler"
	"free5gc/src/amf/gmm"
	"log"
	"testing"
	"time"
)

func sendCreateUEContextRequestAndPrintResult(client *Namf_Communication_Client.APIClient, supi string, request models.CreateUeContextRequest) {
	ueContextInfo, httpResponse, err := client.IndividualUeContextDocumentApi.CreateUEContext(context.Background(), supi, request)
	if err != nil {
		if httpResponse == nil {
			log.Println(err)
		} else if err.Error() != httpResponse.Status {
			log.Println(err)
		} else {
			var ueContextCreateError models.UeContextCreateError
			ueContextCreateError = err.(common.GenericOpenAPIError).Model().(models.UeContextCreateError)
			TestAmf.Config.Dump(ueContextCreateError.Error)
		}
	} else {
		TestAmf.Config.Dump(ueContextInfo)
	}
}

func sendReleaseContextRequestAndPrintResult(client *Namf_Communication_Client.APIClient, supi string, request models.UeContextRelease) {
	httpResponse, err := client.IndividualUeContextDocumentApi.ReleaseUEContext(context.Background(), supi, request)
	if err != nil {
		if httpResponse == nil {
			log.Println(err)
		} else if err.Error() != httpResponse.Status {
			log.Println(err)
		}
	} else {

	}
}

func sendContextTransferRequestAndPrintResult(client *Namf_Communication_Client.APIClient, supi string, request models.UeContextTransferRequest) {
	ueContextTransferResponse, httpResponse, err := client.IndividualUeContextDocumentApi.UEContextTransfer(context.Background(), supi, request)
	if err != nil {
		if httpResponse == nil {
			log.Println(err)
		} else if err.Error() != httpResponse.Status {
			log.Println(err)
		} else {

		}
	} else {
		TestAmf.Config.Dump(ueContextTransferResponse)
	}
}

func sendEBIAssignmentRequestAndPrintResult(client *Namf_Communication_Client.APIClient, supi string, request models.AssignEbiData) {
	assignedEbiData, httpResponse, err := client.IndividualUeContextDocumentApi.EBIAssignment(context.Background(), supi, request)
	if err != nil {
		if httpResponse == nil {
			log.Println(err)
		} else if err.Error() != httpResponse.Status {
			log.Println(err)
		} else {

		}
	} else {
		TestAmf.Config.Dump(assignedEbiData)
	}
}

func sendRegistrationStatusUpdateRequestAndPrintResult(client *Namf_Communication_Client.APIClient, supi string, request models.UeRegStatusUpdateReqData) {
	ueRegStatusUpdateRspData, httpResponse, err := client.IndividualUeContextDocumentApi.RegistrationStatusUpdate(context.Background(), supi, request)
	if err != nil {
		if httpResponse == nil {
			log.Println(err)
		} else if err.Error() != httpResponse.Status {
			log.Println(err)
		} else {

		}
	} else {
		TestAmf.Config.Dump(ueRegStatusUpdateRspData)
	}
}

func TestCreateUEContext(t *testing.T) {
	if len(TestAmf.TestAmf.UePool) == 0 {
		go func() {
			router := Namf_Communication_Server.NewRouter()
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
	configuration := Namf_Communication_Client.NewConfiguration()
	configuration.SetBasePath("https://localhost:29518")
	client := Namf_Communication_Client.NewAPIClient(configuration)

	/* init ue info*/
	ue := TestAmf.TestAmf.UePool["imsi-2089300007487"]

	ueContextCreateData := TestComm.ConsumerAMFCreateUEContextRequsetTable[TestComm.CreateUEContext403]
	sendCreateUEContextRequestAndPrintResult(client, ue.Supi, ueContextCreateData)

	ueContextCreateData = TestComm.ConsumerAMFCreateUEContextRequsetTable[TestComm.CreateUEContext201]
	sendCreateUEContextRequestAndPrintResult(client, ue.Supi, ueContextCreateData)
}

func TestReleaseUEContext(t *testing.T) {
	if len(TestAmf.TestAmf.UePool) == 0 {
		TestCreateUEContext(t)
	}

	configuration := Namf_Communication_Client.NewConfiguration()
	configuration.SetBasePath("https://localhost:29518")
	client := Namf_Communication_Client.NewAPIClient(configuration)

	/* init ue info*/
	self := amf_context.AMF_Self()
	supi := "imsi-0010202"
	ue := self.NewAmfUe(supi)
	if err := gmm.InitAmfUeSm(ue); err != nil {
		t.Errorf("InitAmfUeSm error: %v", err)
	}
	ue.Supi = "imsi-111222"
	ueContextReleaseData := TestComm.ConsumerAMFReleaseUEContextRequestTable[TestComm.UeContextRelease404]
	sendReleaseContextRequestAndPrintResult(client, ue.Supi, ueContextReleaseData)

	ue = TestAmf.TestAmf.UePool["imsi-2089300007487"]
	ueContextReleaseData = TestComm.ConsumerAMFReleaseUEContextRequestTable[TestComm.UeContextRelease201]
	sendReleaseContextRequestAndPrintResult(client, ue.Supi, ueContextReleaseData)

}

func TestUEContextTransfer(t *testing.T) {
	if len(TestAmf.TestAmf.UePool) == 0 {
		TestCreateUEContext(t)
	}

	configuration := Namf_Communication_Client.NewConfiguration()
	configuration.SetBasePath("https://localhost:29518")
	client := Namf_Communication_Client.NewAPIClient(configuration)

	/* init ue info*/
	self := amf_context.AMF_Self()
	supi := "imsi-0010202"
	ue := self.NewAmfUe(supi)
	if err := gmm.InitAmfUeSm(ue); err != nil {
		t.Errorf("InitAmfUeSm error: %v", err.Error())
	}
	ue.Supi = "imsi-111222"
	ueContextTransferData := TestComm.ConsumerAMFUEContextTransferRequestTable[TestComm.UeContextTransfer404]
	sendContextTransferRequestAndPrintResult(client, ue.Supi, ueContextTransferData)

	ue, err := TestAmf.TestAmf.UePool["imsi-2089300007487"]
	if err == false {
		// ue imsi-2089300007487 does not in ue pool
		supi := "imsi-2089300007487"
		ue = TestAmf.TestAmf.NewAmfUe(supi)
		if err := gmm.InitAmfUeSm(ue); err != nil {
			t.Errorf("InitAmfUeSm error: %v", err.Error())
		}
	}
	ueContextTransferData = TestComm.ConsumerAMFUEContextTransferRequestTable[TestComm.UeContextTransferINIT_REG200]
	sendContextTransferRequestAndPrintResult(client, ue.Supi, ueContextTransferData)
}

func TestEBIAssignment(t *testing.T) {
	if len(TestAmf.TestAmf.UePool) == 0 {
		TestCreateUEContext(t)
	}
	configuration := Namf_Communication_Client.NewConfiguration()
	configuration.SetBasePath("https://localhost:29518")
	client := Namf_Communication_Client.NewAPIClient(configuration)

	/* init ue info*/
	ue := TestAmf.TestAmf.UePool["imsi-2089300007487"]

	assignEbiData := TestComm.ConsumerAMFUEContextEBIAssignmentTable[TestComm.AssignEbiData403]
	sendEBIAssignmentRequestAndPrintResult(client, ue.Supi, assignEbiData)

	ue = TestAmf.TestAmf.UePool["imsi-2089300007487"]
	ue.SmContextList[10] = &amf_context.SmContext{
		SmfId:        "123",
		SmfUri:       "https://localhost:29503",
		PlmnId:       models.PlmnId{},
		UserLocation: models.UserLocation{},
		PduSessionContext: &models.PduSessionContext{
			PduSessionId:     10,
			SmContextRef:     "",
			SNssai:           nil,
			Dnn:              "",
			AccessType:       "",
			AllocatedEbiList: nil,
			HsmfId:           "",
			VsmfId:           "",
			NsInstance:       "",
		},
	}
	ebiArpMapping := models.EbiArpMapping{
		EpsBearerId: 10,
		Arp:         nil,
	}
	ue.SmContextList[10].PduSessionContext.AllocatedEbiList = append(ue.SmContextList[10].PduSessionContext.AllocatedEbiList, ebiArpMapping)
	assignEbiData = TestComm.ConsumerAMFUEContextEBIAssignmentTable[TestComm.AssignEbiData200]
	sendEBIAssignmentRequestAndPrintResult(client, ue.Supi, assignEbiData)
}

func TestRegistrationStatusUpdate(t *testing.T) {
	if len(TestAmf.TestAmf.UePool) == 0 {
		TestCreateUEContext(t)
	}
	configuration := Namf_Communication_Client.NewConfiguration()
	configuration.SetBasePath("https://localhost:29518")
	client := Namf_Communication_Client.NewAPIClient(configuration)

	ueRegStatusUpdateReqData := TestComm.ConsumerRegistrationStatusUpdateTable[TestComm.RegistrationStatusUpdate404]
	sendRegistrationStatusUpdateRequestAndPrintResult(client, "", ueRegStatusUpdateReqData)
	/* init ue info*/
	ue := TestAmf.TestAmf.UePool["imsi-2089300007487"]
	ueRegStatusUpdateReqData = TestComm.ConsumerRegistrationStatusUpdateTable[TestComm.RegistrationStatusUpdate200]
	sendRegistrationStatusUpdateRequestAndPrintResult(client, ue.Supi, ueRegStatusUpdateReqData)

}
