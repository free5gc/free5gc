package BDTPolicy_test

import (
	"context"
	"flag"
	TestBDTPolicy "free5gc/lib/CommonConsumerTestData/PCF/TestBDTPolicy"
	Npcf_BDTPolicy "free5gc/lib/Npcf_BDTPolicy"
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
func TestCreateBDTPolicy(t *testing.T) {
	pcfInit()

	configuration := Npcf_BDTPolicy.NewConfiguration()
	configuration.SetBasePath(pcf_util.PCF_BASIC_PATH + pcf_context.BdtUri)
	client := Npcf_BDTPolicy.NewAPIClient(configuration)

	// get test data
	bdtReqData1, _, _, _ := TestBDTPolicy.GetCreateTestData()

	// test create service
	_, httpRsp, err := client.BDTPoliciesCollectionApi.CreateBDTPolicy(context.Background(), bdtReqData1)
	assert.True(t, err == nil)
	assert.True(t, httpRsp != nil)
	assert.Equal(t, "201 Created", httpRsp.Status)

}
