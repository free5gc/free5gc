package PolicyAuthorization_test

import (
	"context"
	TestPolicyAuthorization "free5gc/lib/CommonConsumerTestData/PCF/TestPolicyAuthorization"
	"free5gc/lib/Npcf_PolicyAuthorization"
	"free5gc/lib/http2_util"
	"free5gc/lib/path_util"
	"free5gc/src/pcf/PolicyAuthorization"
	"free5gc/src/pcf/pcf_context"
	"free5gc/src/pcf/pcf_handler"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestApplicationSessionsCollection(t *testing.T) {
	go func() {
		pcfrouter := PolicyAuthorization.NewRouter()
		pcfserver, err := http2_util.NewServer(":29507", path_util.Gofree5gcPath("free5gc/pcfsslkey.log"), pcfrouter)
		if err == nil && pcfserver != nil {
			err := pcfserver.ListenAndServeTLS(path_util.Gofree5gcPath("free5gc/support/TLS/pcf.pem"), path_util.Gofree5gcPath("free5gc/support/TLS/pcf.key"))
			assert.True(t, err == nil)
		}
	}()
	go pcf_handler.Handle()
	configuration := Npcf_PolicyAuthorization.NewConfiguration()
	configuration.SetBasePath("https://localhost:29507/npcf-policyauthorization/v1")
	client := Npcf_PolicyAuthorization.NewAPIClient(configuration)

	//Get Test Data
	PostAppSessions201Data := TestPolicyAuthorization.GetPostAppSessions201Data()
	PostAppSessions403Data := TestPolicyAuthorization.GetPostAppSessions403Data()

	//Create UE
	pcf_context.NewPCFUe("123")

	//Test PostAppSessions 201
	_, httpRsp, err := client.ApplicationSessionsCollectionApi.PostAppSessions(context.Background(), PostAppSessions201Data)
	assert.True(t, err == nil)
	assert.True(t, httpRsp != nil)
	assert.Equal(t, "201 Created", httpRsp.Status)
	var appSessionId = httpRsp.Header.Get("Location")

	//Test PostAppSessions 303
	_, httpRsp, err = client.ApplicationSessionsCollectionApi.PostAppSessions(context.Background(), PostAppSessions201Data)
	assert.True(t, err == nil)
	assert.True(t, httpRsp != nil)
	assert.Equal(t, "200 OK", httpRsp.Status)

	//Test PostAppSessions 403
	_, httpRsp, err = client.ApplicationSessionsCollectionApi.PostAppSessions(context.Background(), PostAppSessions403Data)
	assert.True(t, err != nil)
	assert.True(t, httpRsp != nil)
	assert.Equal(t, "403 Forbidden", httpRsp.Status)

	var deleteAppSessionParamOpts Npcf_PolicyAuthorization.DeleteAppSessionParamOpts
	_, httpRsp, err = client.IndividualApplicationSessionContextDocumentApi.DeleteAppSession(context.Background(), appSessionId, &deleteAppSessionParamOpts)
}
