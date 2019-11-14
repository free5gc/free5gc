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

func TestIndividualApplicationSessionContext(t *testing.T) {
	pcfInit()
	configuration := Npcf_PolicyAuthorization.NewConfiguration()
	configuration.SetBasePath(pcf_util.PCF_BASIC_PATH + pcf_context.PolicyAuthorizationUri)
	client := Npcf_PolicyAuthorization.NewAPIClient(configuration)

	//Get Test Data
	PostAppSessions201Data := TestPolicyAuthorization.GetPostAppSessions201Data()
	DeleteAppSession204Data := TestPolicyAuthorization.GetDeleteAppSession204Data()
	ModAppSession200Data := TestPolicyAuthorization.GetModAppSession200Data()
	ModAppSession403Data := TestPolicyAuthorization.GetModAppSession403Data()
	pcf_context.PCF_Self().NewPCFUe("123")

	//Test PostAppSessions Success
	_, httpRsp, err := client.ApplicationSessionsCollectionApi.PostAppSessions(context.Background(), PostAppSessions201Data)
	assert.True(t, err == nil)
	assert.True(t, httpRsp != nil)
	assert.Equal(t, "201 Created", httpRsp.Status)
	var appSessionId = httpRsp.Header.Get("Location")

	// Test ModAppSession 200
	_, httpRsp, err = client.IndividualApplicationSessionContextDocumentApi.ModAppSession(context.Background(), appSessionId, ModAppSession200Data)
	assert.True(t, err == nil)
	assert.True(t, httpRsp != nil)
	assert.Equal(t, "200 OK", httpRsp.Status)

	// Test ModAppSession 403
	_, httpRsp, err = client.IndividualApplicationSessionContextDocumentApi.ModAppSession(context.Background(), appSessionId, ModAppSession403Data)
	assert.True(t, err != nil)
	assert.True(t, httpRsp != nil)
	assert.Equal(t, "403 Forbidden", httpRsp.Status)

	// Test ModAppSession 404
	_, httpRsp, err = client.IndividualApplicationSessionContextDocumentApi.ModAppSession(context.Background(), "456", ModAppSession200Data)
	assert.True(t, err != nil)
	assert.True(t, httpRsp != nil)
	assert.Equal(t, "404 Not Found", httpRsp.Status)

	// Test GetAppSession 200
	_, httpRsp, err = client.IndividualApplicationSessionContextDocumentApi.GetAppSession(context.Background(), appSessionId)
	assert.True(t, err == nil)
	assert.True(t, httpRsp != nil)
	assert.Equal(t, "200 OK", httpRsp.Status)

	// Test GetAppSession 404
	_, httpRsp, err = client.IndividualApplicationSessionContextDocumentApi.GetAppSession(context.Background(), "456")
	assert.True(t, err != nil)
	assert.True(t, httpRsp != nil)
	assert.Equal(t, "404 Not Found", httpRsp.Status)

	// Test DeleteAppSession 200
	var deleteAppSessionParamOpts Npcf_PolicyAuthorization.DeleteAppSessionParamOpts
	_, httpRsp, err = client.IndividualApplicationSessionContextDocumentApi.DeleteAppSession(context.Background(), appSessionId, &deleteAppSessionParamOpts)
	assert.True(t, err == nil)
	assert.True(t, httpRsp != nil)
	assert.Equal(t, "200 OK", httpRsp.Status)

	// Test DeleteAppSession 204
	_, httpRsp, err = client.ApplicationSessionsCollectionApi.PostAppSessions(context.Background(), DeleteAppSession204Data)
	appSessionId = httpRsp.Header.Get("Location")
	_, httpRsp, err = client.IndividualApplicationSessionContextDocumentApi.DeleteAppSession(context.Background(), appSessionId, nil)
	assert.True(t, err == nil)
	assert.True(t, httpRsp != nil)
	assert.Equal(t, "204 No Content", httpRsp.Status)
}
