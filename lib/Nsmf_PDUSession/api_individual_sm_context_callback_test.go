package Nsmf_PDUSession_test

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"free5gc/lib/CommonConsumerTestData/AMF/TestAmf"
	"free5gc/lib/Nsmf_PDUSession"
	"free5gc/lib/http2_util"
	"free5gc/lib/openapi/models"
	"free5gc/src/amf/HttpCallback"
	"free5gc/src/amf/amf_context"
	"free5gc/src/amf/amf_handler"
	"free5gc/src/amf/amf_util"
	"net/http"
	"strconv"
	"testing"
	"time"
)

func init() {
	TestAmf.AmfInit()

	router := gin.Default()

	Namf_Callback.AddService(router)
	go amf_handler.Handle()

	addr := fmt.Sprintf("%s:%d", TestAmf.TestAmf.HttpIPv4Address, TestAmf.TestAmf.HttpIpv4Port)
	server, err := http2_util.NewServer(addr, amf_util.AmfLogPath, router)
	if err == nil && server != nil {
		go server.ListenAndServeTLS(amf_util.AmfPemPath, amf_util.AmfKeyPath)
	}
}

func TestSMContextCreate(t *testing.T) {
	time.Sleep(100 * time.Millisecond)
	TestAmf.UeAttach(models.AccessType__3_GPP_ACCESS)
	ue := TestAmf.TestAmf.UePool["imsi-2089300007487"]

	pduSession := models.PduSessionContext{
		PduSessionId: 10,
		Dnn:          "nctu.edu.tw",
		SNssai: &models.Snssai{
			Sst: 1,
			Sd:  "020304",
		},
	}

	uri := TestAmf.TestAmf.GetIPv4Uri() + "/namf-callback/v1/smContextStatus/" + ue.Guti + "/" + strconv.Itoa(int(pduSession.PduSessionId))

	var smContext amf_context.SmContext
	smContext.PduSessionContext = &pduSession
	ue.SmContextList[pduSession.PduSessionId] = &smContext

	request := models.SmContextStatusNotification{}
	request.StatusInfo = &models.StatusInfo{
		ResourceStatus: models.ResourceStatus_RELEASED,
	}
	configuration := Nsmf_PDUSession.NewConfiguration()
	client := Nsmf_PDUSession.NewAPIClient(configuration)
	httpResponse, err := client.IndividualSMContextNotificationApi.SMContextNotification(context.Background(), uri, request)
	assert.Nil(t, err)
	if assert.NotNil(t, httpResponse) {
		assert.Equal(t, http.StatusNoContent, httpResponse.StatusCode)
		assert.Equal(t, 0, len(ue.SmContextList))
	}

}
