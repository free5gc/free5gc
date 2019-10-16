package AMPolicy_test

import (
	"context"
	"free5gc/lib/CommonConsumerTestData/PCF/TestAMPolicy"
	"free5gc/lib/Npcf_AMPolicy"
	"free5gc/lib/http2_util"
	AMPolicy "free5gc/src/pcf/AMPolicy"
	"free5gc/src/pcf/pcf_handler"
	"free5gc/src/pcf/pcf_util"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateAMPolicy(t *testing.T) {
	go func() {
		pcfrouter := AMPolicy.NewRouter()
		pcfserver, err := http2_util.NewServer(":29507", pcf_util.PCF_LOG_PATH, pcfrouter)
		if err == nil && pcfserver != nil {
			err := pcfserver.ListenAndServeTLS(pcf_util.PCF_PEM_PATH, pcf_util.PCF_KEY_PATH)
			assert.True(t, err == nil)
		}
	}()
	go pcf_handler.Handle()
	configuration := Npcf_AMPolicy.NewConfiguration()
	configuration.SetBasePath("https://localhost:29507/npcf-am-policy-control/v1")
	client := Npcf_AMPolicy.NewAPIClient(configuration)

	//Test PostPolicies
	amCreateReqData := TestAMPolicy.GetAMreqdata()
	_, httpRsp, err := client.AMPolicyAssociationsCollectionApi.CreateIndividualAMPolicyAssociation(context.Background(), amCreateReqData)
	assert.True(t, err == nil)
	assert.True(t, httpRsp != nil)
	assert.Equal(t, "201 Created", httpRsp.Status)

	//Test GetPoliciesPolAssoId
	_, httpRsp, err = client.IndividualAMPolicyAssociationDocumentApi.ReadIndividualAMPolicyAssociation(context.Background(), "46611123456789")
	assert.True(t, err == nil)
	assert.True(t, httpRsp != nil)
	assert.Equal(t, "200 OK", httpRsp.Status)

	//Test Update
	amUpdateReqData := TestAMPolicy.GetAMUpdateReqData()
	_, httpRsp, err = client.IndividualAMPolicyAssociationDocumentApi.ReportObservedEventTriggersForIndividualAMPolicyAssociation(context.Background(), "46611123456789", amUpdateReqData)
	assert.True(t, err == nil)
	assert.True(t, httpRsp != nil)
	assert.Equal(t, "200 OK", httpRsp.Status)

	//Test PoliciesPolAssoIdDelete
	httpRsp, err = client.IndividualAMPolicyAssociationDocumentApi.DeleteIndividualAMPolicyAssociation(context.Background(), "46611123456789")
	assert.True(t, err == nil)
	assert.True(t, httpRsp != nil)
	assert.Equal(t, "204 No Content", httpRsp.Status)

	//---------------------------------------------------------------------------------------------------------
	//Fail Test (Create no notifyURI)
	amCreatefailnotifyURIData := TestAMPolicy.GetamCreatefailnotifyURIData()
	_, httpRsp, err = client.AMPolicyAssociationsCollectionApi.CreateIndividualAMPolicyAssociation(context.Background(), amCreatefailnotifyURIData)
	assert.True(t, err != nil)
	assert.True(t, httpRsp != nil)
	assert.Equal(t, "400 Bad Request", httpRsp.Status)

	//Fail Test (Create no supi)
	amCreatefailsupiData := TestAMPolicy.GetamCreatefailsupiData()
	_, httpRsp, err = client.AMPolicyAssociationsCollectionApi.CreateIndividualAMPolicyAssociation(context.Background(), amCreatefailsupiData)
	assert.True(t, err != nil)
	assert.True(t, httpRsp != nil)
	assert.Equal(t, "400 Bad Request", httpRsp.Status)

	//Fail Test (Create no suppfeat)
	amCreatefailsuppfeatData := TestAMPolicy.GetamCreatefailsuppfeatData()
	_, httpRsp, err = client.AMPolicyAssociationsCollectionApi.CreateIndividualAMPolicyAssociation(context.Background(), amCreatefailsuppfeatData)
	assert.True(t, err != nil)
	assert.True(t, httpRsp != nil)
	assert.Equal(t, "400 Bad Request", httpRsp.Status)
	//---------------------------------------------------------------------------------------------------

	//Test PostPolicies
	amCreateReqData = TestAMPolicy.GetAMreqdata()
	_, httpRsp, err = client.AMPolicyAssociationsCollectionApi.CreateIndividualAMPolicyAssociation(context.Background(), amCreateReqData)
	assert.True(t, err == nil)
	assert.True(t, httpRsp != nil)
	assert.Equal(t, "201 Created", httpRsp.Status)

	//Fail Test (GetID wrong ID)
	_, httpRsp, err = client.IndividualAMPolicyAssociationDocumentApi.ReadIndividualAMPolicyAssociation(context.Background(), "11111111")
	assert.True(t, err != nil)
	assert.True(t, httpRsp != nil)
	assert.Equal(t, "404 Not Found", httpRsp.Status)

	//Test PoliciesPolAssoIdDelete
	httpRsp, err = client.IndividualAMPolicyAssociationDocumentApi.DeleteIndividualAMPolicyAssociation(context.Background(), "46611123456789")
	assert.True(t, err == nil)
	assert.True(t, httpRsp != nil)
	assert.Equal(t, "204 No Content", httpRsp.Status)

	//Test Fail Update
	amCreateReqData = TestAMPolicy.GetAMreqdata()
	_, httpRsp, err = client.AMPolicyAssociationsCollectionApi.CreateIndividualAMPolicyAssociation(context.Background(), amCreateReqData)
	assert.True(t, err == nil)
	assert.True(t, httpRsp != nil)
	assert.Equal(t, "201 Created", httpRsp.Status)

	//Test Fail Update
	amUpdateReqData = TestAMPolicy.GetAMUpdateReqData()
	_, httpRsp, err = client.IndividualAMPolicyAssociationDocumentApi.ReportObservedEventTriggersForIndividualAMPolicyAssociation(context.Background(), "1111111", amUpdateReqData)
	assert.True(t, err != nil)
	assert.True(t, httpRsp != nil)
	assert.Equal(t, "404 Not Found", httpRsp.Status)

	//Test PoliciesPolAssoIdDelete
	httpRsp, err = client.IndividualAMPolicyAssociationDocumentApi.DeleteIndividualAMPolicyAssociation(context.Background(), "46611123456789")
	assert.True(t, err == nil)
	assert.True(t, httpRsp != nil)
	assert.Equal(t, "204 No Content", httpRsp.Status)

}
