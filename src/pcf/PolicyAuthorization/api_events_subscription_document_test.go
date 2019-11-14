package PolicyAuthorization_test

import (
	"context"
	TestPolicyAuthorization "free5gc/lib/CommonConsumerTestData/PCF/TestPolicyAuthorization"
	"free5gc/lib/Npcf_PolicyAuthorization"
	"free5gc/src/pcf/pcf_context"
	"free5gc/src/pcf/pcf_util"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEventsSubscription(t *testing.T) {
	pcfInit()
	configuration := Npcf_PolicyAuthorization.NewConfiguration()
	configuration.SetBasePath(pcf_util.PCF_BASIC_PATH + pcf_context.PolicyAuthorizationUri)
	client := Npcf_PolicyAuthorization.NewAPIClient(configuration)

	//Get Test Data
	PostAppSessions201Data := TestPolicyAuthorization.GetPostAppSessions201Data()
	UpdateEventsSubsc201Data := TestPolicyAuthorization.GetUpdateEventsSubsc201Data()
	UpdateEventsSubsc200Data := TestPolicyAuthorization.GetUpdateEventsSubsc200Data()
	UpdateEventsSubsc403Data := TestPolicyAuthorization.GetUpdateEventsSubsc403Data()
	pcf_context.PCF_Self().NewPCFUe("123")
	//Test PostAppSessions Success
	_, httpRsp, err := client.ApplicationSessionsCollectionApi.PostAppSessions(context.Background(), PostAppSessions201Data)
	assert.True(t, err == nil)
	assert.True(t, httpRsp != nil)
	assert.Equal(t, "201 Created", httpRsp.Status)
	var appSessionId = httpRsp.Header.Get("Location")

	//Test UpdateEventsSubsc 201
	_, httpRsp, err = client.EventsSubscriptionDocumentApi.UpdateEventsSubsc(context.Background(), appSessionId, UpdateEventsSubsc201Data)
	assert.True(t, err == nil)
	assert.True(t, httpRsp != nil)
	assert.Equal(t, "201 Created", httpRsp.Status)

	//Test UpdateEventsSubsc 200
	_, httpRsp, err = client.EventsSubscriptionDocumentApi.UpdateEventsSubsc(context.Background(), appSessionId, UpdateEventsSubsc200Data)
	assert.True(t, err == nil)
	assert.True(t, httpRsp != nil)
	assert.Equal(t, "200 OK", httpRsp.Status)

	//Test UpdateEventsSubsc 403
	_, httpRsp, err = client.EventsSubscriptionDocumentApi.UpdateEventsSubsc(context.Background(), appSessionId, UpdateEventsSubsc403Data)
	assert.True(t, err != nil)
	assert.True(t, httpRsp != nil)
	assert.Equal(t, "403 Forbidden", httpRsp.Status)

	//Test UpdateEventsSubsc 404
	_, httpRsp, err = client.EventsSubscriptionDocumentApi.UpdateEventsSubsc(context.Background(), "456", UpdateEventsSubsc201Data)
	assert.True(t, err != nil)
	assert.True(t, httpRsp != nil)
	assert.Equal(t, "404 Not Found", httpRsp.Status)

	// Test DeleteEventsSubsc 204
	httpRsp, err = client.EventsSubscriptionDocumentApi.DeleteEventsSubsc(context.Background(), appSessionId)
	assert.True(t, err == nil)
	assert.True(t, httpRsp != nil)
	assert.Equal(t, "204 No Content", httpRsp.Status)

	// Test DeleteEventsSubsc 404
	httpRsp, err = client.EventsSubscriptionDocumentApi.DeleteEventsSubsc(context.Background(), "456")
	assert.True(t, err != nil)
	assert.True(t, httpRsp != nil)
	assert.Equal(t, "404 Not Found", httpRsp.Status)

	var deleteAppSessionParamOpts Npcf_PolicyAuthorization.DeleteAppSessionParamOpts
	_, httpRsp, err = client.IndividualApplicationSessionContextDocumentApi.DeleteAppSession(context.Background(), appSessionId, &deleteAppSessionParamOpts)
}
