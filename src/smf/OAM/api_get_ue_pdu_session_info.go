package Nsmf_OAM

import (
	"github.com/gin-gonic/gin"
	"free5gc/lib/http_wrapper"
	"free5gc/src/smf/smf_handler/smf_message"
)

func GetUEPDUSessionInfo(c *gin.Context) {

	req := http_wrapper.NewRequest(c.Request, nil)
	req.Params["smContextRef"] = c.Params.ByName("smContextRef")

	handlerMsg := smf_message.NewHandlerMessage(smf_message.OAMGetUEPDUSessionInfo, req)
	smf_message.SendMessage(handlerMsg)

	rsp := <-handlerMsg.ResponseChan

	HTTPResponse := rsp.HTTPResponse

	c.JSON(HTTPResponse.Status, HTTPResponse.Body)
}
