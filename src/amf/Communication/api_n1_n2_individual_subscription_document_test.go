package Communication_test

import (
	"context"
	"free5gc/lib/CommonConsumerTestData/AMF/TestAmf"
	Namf_Communication_Client "free5gc/lib/Namf_Communication"
	"free5gc/lib/openapi/common"
	"free5gc/lib/openapi/models"
	"free5gc/src/amf/gmm"
	"log"
	"testing"
)

func sendN1N2MessageUnSubscribeRequestAndPrintResult(client *Namf_Communication_Client.APIClient, supi string, subscriptionID string) {
	httpResponse, err := client.N1N2IndividualSubscriptionDocumentApi.N1N2MessageUnSubscribe(context.Background(), supi, subscriptionID)
	if err != nil {
		if httpResponse == nil {
			log.Panic(err)
		} else if err.Error() != httpResponse.Status {
			log.Panic(err)
		} else {
			var probelmDetail models.ProblemDetails
			probelmDetail = err.(common.GenericOpenAPIError).Model().(models.ProblemDetails)
			TestAmf.Config.Dump(probelmDetail)
		}
	} else {

	}
}

func TestN1N2MessageUnSubscribe(t *testing.T) {
	if len(TestAmf.TestAmf.UePool) == 0 {
		TestN1N2MessageSubscribe(t)
	}
	configuration := Namf_Communication_Client.NewConfiguration()
	configuration.SetBasePath("https://localhost:29518")
	client := Namf_Communication_Client.NewAPIClient(configuration)

	/* init ue info*/
	ue, err := TestAmf.TestAmf.UePool["imsi-2089300007487"]
	if err == false {
		// ue imsi-2089300007487 does not in ue pool
		supi := "imsi-2089300007487"
		ue = TestAmf.TestAmf.NewAmfUe(supi)
		if err := gmm.InitAmfUeSm(ue); err != nil {
			t.Errorf("InitAmfUeSm error: %v", err.Error())
		}
	}

	sendN1N2MessageUnSubscribeRequestAndPrintResult(client, ue.Supi, "0")
}
