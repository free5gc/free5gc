package AMPolicy_test

import (
	"context"
	"flag"
	"free5gc/lib/CommonConsumerTestData/PCF/TestAMPolicy"
	"free5gc/lib/Npcf_AMPolicy"
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
func TestCreateAMPolicy(t *testing.T) {
	pcfInit()

	configuration := Npcf_AMPolicy.NewConfiguration()
	configuration.SetBasePath(pcf_util.PCF_BASIC_PATH + pcf_context.AmpolicyUri)
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
