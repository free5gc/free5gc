package Namf_Callback

import (
	"free5gc/lib/http_wrapper"
	"free5gc/lib/openapi/models"
	"free5gc/src/amf/amf_handler/amf_message"
	"free5gc/src/amf/logger"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SmContextStatusNotify(c *gin.Context) {

	var request models.SmContextStatusNotification

	err := c.ShouldBindJSON(&request)
	if err != nil {
		problemDetail := "[Request Body] " + err.Error()
		rsp := models.ProblemDetails{
			Title:  "Malformed request syntax",
			Status: http.StatusBadRequest,
			Detail: problemDetail,
		}
		logger.CallbackLog.Errorln(problemDetail)
		c.JSON(http.StatusBadRequest, rsp)
		return
	}

	req := http_wrapper.NewRequest(c.Request, request)
	req.Params["guti"] = c.Params.ByName("guti")
	req.Params["pduSessionId"] = c.Params.ByName("pduSessionId")

	handlerMsg := amf_message.NewHandlerMessage(amf_message.EventSmContextStatusNotify, req)
	amf_message.SendMessage(handlerMsg)

	rsp := <-handlerMsg.ResponseChan

	HTTPResponse := rsp.HTTPResponse

	c.JSON(HTTPResponse.Status, HTTPResponse.Body)

}
