package PDUSession_test

import (
	"context"
	"github.com/stretchr/testify/assert"
	"free5gc/lib/CommonConsumerTestData/SMF/TestPDUSession"
	"free5gc/lib/Nsmf_PDUSession"
	"free5gc/lib/openapi/models"
	"free5gc/src/smf/PDUSession"
	"free5gc/src/smf/smf_handler"
	"testing"
)

func TestUpdateSmContext(t *testing.T) {
	go smf_handler.Handle()
	go PDUSession.DummyServer()
	configuration := Nsmf_PDUSession.NewConfiguration()
	configuration.SetBasePath("https://127.0.0.10:29502")
	client := Nsmf_PDUSession.NewAPIClient(configuration)
	var request models.UpdateSmContextRequest

	table := TestPDUSession.ConsumerSMFPDUSessionUpdateContextTable["ACTIVATING"]
	request.JsonData = table.JsonData
	request.BinaryDataN1SmMessage = table.BinaryDataN1SmMessage

	_, httpRsp, _ := client.IndividualSMContextApi.UpdateSmContext(context.Background(), "123", request)
	assert.True(t, httpRsp != nil)
	assert.Equal(t, "404 Not Found", httpRsp.Status)

}
