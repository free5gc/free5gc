package Communication_test

import (
	"context"
	"free5gc/lib/CommonConsumerTestData/AMF/TestAmf"
	"free5gc/lib/CommonConsumerTestData/AMF/TestComm"
	Namf_Communication_Client "free5gc/lib/Namf_Communication"
	"free5gc/lib/openapi/models"

	"log"
	"testing"
)

func sendAMFStatusUnSubscriptionRequestAndPrintResult(client *Namf_Communication_Client.APIClient, subscriptionId string) {
	httpResponse, err := client.IndividualSubscriptionDocumentApi.AMFStatusChangeUnSubscribe(context.Background(), subscriptionId)
	if err != nil {
		if httpResponse == nil {
			log.Println(err)
		} else if err.Error() != httpResponse.Status {
			log.Println(err)
		} else {

		}
	} else {

	}
}

func sendAMFStatusSubscriptionModfyRequestAndPrintResult(client *Namf_Communication_Client.APIClient, subscriptionID string, request models.SubscriptionData) {
	aMFStatusSubscription, httpResponse, err := client.IndividualSubscriptionDocumentApi.AMFStatusChangeSubscribeModfy(context.Background(), subscriptionID, request)
	if err != nil {
		if httpResponse == nil {
			log.Println(err)
		} else if err.Error() != httpResponse.Status {
			log.Println(err)
		} else {

		}
	} else {
		TestAmf.Config.Dump(aMFStatusSubscription)
	}
}

func TestAMFStatusChangeSubscribeModfy(t *testing.T) {
	if len(TestAmf.TestAmf.UePool) == 0 {
		TestAMFStatusChangeSubscribe(t)
	}
	configuration := Namf_Communication_Client.NewConfiguration()
	configuration.SetBasePath("https://localhost:29518")
	client := Namf_Communication_Client.NewAPIClient(configuration)

	subscriptionData := TestComm.ConsumerAMFStatusChangeSubscribeModfyTable[TestComm.AMFStatusSubscriptionModfy403]
	sendAMFStatusSubscriptionModfyRequestAndPrintResult(client, "0", subscriptionData)
	//
	subscriptionData = TestComm.ConsumerAMFStatusChangeSubscribeModfyTable[TestComm.AMFStatusSubscriptionModfy200]
	sendAMFStatusSubscriptionModfyRequestAndPrintResult(client, "1", subscriptionData)
}

func TestAMFStatusChangeUnSubscribe(t *testing.T) {
	if len(TestAmf.TestAmf.UePool) == 0 {
		TestAMFStatusChangeSubscribe(t)
	}
	configuration := Namf_Communication_Client.NewConfiguration()
	configuration.SetBasePath("https://localhost:29518")
	client := Namf_Communication_Client.NewAPIClient(configuration)

	subscriptionID := TestComm.ConsumerAMFStatusUnSubscriptionTable[TestComm.AMFStatusUnSubscription403]
	sendAMFStatusUnSubscriptionRequestAndPrintResult(client, subscriptionID)
	//
	subscriptionID = TestComm.ConsumerAMFStatusUnSubscriptionTable[TestComm.AMFStatusUnSubscription204]
	sendAMFStatusUnSubscriptionRequestAndPrintResult(client, subscriptionID)
}
