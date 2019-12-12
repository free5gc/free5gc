package PolicyAuthorization_test

import (
	"context"
	"flag"
	"free5gc/lib/CommonConsumerTestData/PCF/TestPolicyAuthorization"
	"free5gc/lib/CommonConsumerTestData/PCF/TestSMPolicy"
	"free5gc/lib/Npcf_PolicyAuthorization"
	"free5gc/lib/Npcf_SMPolicyControl"
	"free5gc/src/pcf/pcf_context"
	"free5gc/src/pcf/pcf_service"
	"free5gc/src/pcf/pcf_util"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli"
)

func pcfInit() {
	flags := flag.FlagSet{}
	c := cli.NewContext(nil, &flags, nil)
	pcf := &pcf_service.PCF{}
	pcf.Initialize(c)
	go pcf.Start()
	time.Sleep(100 * time.Millisecond)
}

func TestApplicationSessionsCollection(t *testing.T) {
	pcfInit()
	configuration := Npcf_PolicyAuthorization.NewConfiguration()
	configuration.SetBasePath(pcf_util.PCF_BASIC_PATH + pcf_context.PolicyAuthorizationUri)
	client := Npcf_PolicyAuthorization.NewAPIClient(configuration)

	//Get Test Data
	PostAppSessions201Data := TestPolicyAuthorization.GetPostAppSessions201Data()
	PostAppSessions403Data := TestPolicyAuthorization.GetPostAppSessions403Data()

	//Create UE
	// pcf_context.NewPCFUe("123")
	pcf_context.PCF_Self().NewPCFUe("string1")
	smReqData := TestSMPolicy.CreateTestData()
	configuration2 := Npcf_SMPolicyControl.NewConfiguration()
	configuration2.SetBasePath(pcf_util.PCF_BASIC_PATH + pcf_context.SmUri)
	client2 := Npcf_SMPolicyControl.NewAPIClient(configuration2)
	_, httpRsp, err := client2.DefaultApi.SmPoliciesPost(context.Background(), smReqData)

	//Test PostAppSessions 201
	_, httpRsp, err = client.ApplicationSessionsCollectionApi.PostAppSessions(context.Background(), PostAppSessions201Data)
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
