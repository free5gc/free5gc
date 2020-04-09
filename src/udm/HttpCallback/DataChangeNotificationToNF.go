package Nudm_Callback

import (
	"github.com/gin-gonic/gin"
	"free5gc/lib/http_wrapper"
	"free5gc/lib/openapi/models"
	"free5gc/src/udm/logger"
	"free5gc/src/udm/udm_handler"
	"free5gc/src/udm/udm_handler/udm_message"
	"net/http"
)

func DataChangeNotificationToNF(c *gin.Context) {

	var request models.DataChangeNotify

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
	}

	req := http_wrapper.NewRequest(c.Request, request)
	req.Params["supi"] = c.Params.ByName("supi")

	handleMsg := udm_message.NewHandlerMessage(udm_message.EventDataChangeNotificationToNF, req)
	udm_handler.SendMessage(handleMsg)

	rsp := <-handleMsg.ResponseChan

	HTTPResponse := rsp.HTTPResponse

	c.JSON(HTTPResponse.Status, HTTPResponse.Body)
}
