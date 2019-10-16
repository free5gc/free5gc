package Nsmf_PDUSession

import (
	"context"
	"fmt"
	"free5gc/lib/openapi/models"
	"log"
	"testing"
)

func TestSMContextsCollectionAPI(t *testing.T) {
	configuration := NewConfiguration()
	configuration.SetBasePath("http://localhost:8080")

	client := NewAPIClient(configuration)

	var smContextsRequest models.PostSmContextsRequest
	smContextsRequest.JsonData = &models.SmContextCreateData{
		N1SmMsg: &models.RefToBinaryData{
			ContentId: "1000",
		},
	}
	tmp := []byte("123")
	smContextsRequest.BinaryDataN1SmMessage = tmp

	postSmContextReponse, httpResponse, err := client.SMContextsCollectionApi.PostSmContexts(context.Background(), smContextsRequest)
	if err != nil {
		log.Fatal(err.Error())
	} else {
		fmt.Println(httpResponse.Status)
		fmt.Println(httpResponse.Proto)
		fmt.Println("Content-Type: ", httpResponse.Header.Get("Content-Type"))
		fmt.Printf("\n  *jsonData: %v\n  *N2SmInformation:\n\tcontentID: %s\n\tData: 0x%0x\n",
			*postSmContextReponse.JsonData,
			postSmContextReponse.JsonData.N2SmInfo.ContentId,
			postSmContextReponse.BinaryDataN2SmInformation)
	}

}
