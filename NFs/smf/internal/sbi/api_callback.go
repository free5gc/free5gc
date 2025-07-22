package sbi

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/free5gc/openapi"
	"github.com/free5gc/openapi/models"
	"github.com/free5gc/smf/internal/logger"
)

func (s *Server) getCallbackRoutes() []Route {
	return []Route{
		{
			Name:    "SmPolicyUpdateNotification",
			Method:  http.MethodPost,
			Pattern: "/sm-policies/:smContextRef/update",
			APIFunc: s.HTTPSmPolicyUpdateNotification,
		},
		{
			Name:    "SmPolicyControlTerminationRequestNotification",
			Method:  http.MethodPost,
			Pattern: "/sm-policies/:smContextRef/terminate",
			APIFunc: s.SmPolicyControlTerminationRequestNotification,
		},
		{
			Name:    "ChargingNotification",
			Method:  http.MethodPost,
			Pattern: "/:notifyUri",
			APIFunc: s.HTTPChargingNotification,
		},
	}
}

// SubscriptionsPost -
func (s *Server) HTTPSmPolicyUpdateNotification(c *gin.Context) {
	var request models.SmPolicyNotification

	reqBody, err := c.GetRawData()
	if err != nil {
		logger.PduSessLog.Errorln("GetRawData failed")
	}

	err = openapi.Deserialize(&request, reqBody, c.ContentType())
	if err != nil {
		logger.PduSessLog.Errorln("Deserialize request failed")
	}

	smContextRef := c.Params.ByName("smContextRef")
	s.Processor().HandleSMPolicyUpdateNotify(c, request, smContextRef)
}

func (s *Server) SmPolicyControlTerminationRequestNotification(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{})
}

func (s *Server) HTTPChargingNotification(c *gin.Context) {
	var req models.ChargingNotifyRequest

	requestBody, err := c.GetRawData()
	if err != nil {
		logger.PduSessLog.Errorln("GetRawData failed")
	}

	err = openapi.Deserialize(&req, requestBody, APPLICATION_JSON)
	if err != nil {
		logger.PduSessLog.Errorln("Deserialize request failed")
	}

	smContextRef := strings.Split(c.Params.ByName("notifyUri"), "_")[1]

	s.Processor().HandleChargingNotification(c, req, smContextRef)
}
