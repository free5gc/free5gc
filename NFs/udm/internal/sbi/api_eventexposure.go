package sbi

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/free5gc/openapi"
	"github.com/free5gc/openapi/models"
	"github.com/free5gc/udm/internal/logger"
)

func (s *Server) getEventExposureRoutes() []Route {
	return []Route{
		{
			"Index",
			http.MethodGet,
			"/",
			s.HandleIndex,
		},

		{
			"CreateEeSubscription",
			http.MethodPost,
			"/:ueIdentity/ee-subscriptions",
			s.HandleCreateEeSubscription,
		},

		{
			"DeleteEeSubscription",
			http.MethodDelete,
			"/:ueIdentity/ee-subscriptions/:subscriptionId",
			s.HandleDeleteEeSubscription,
		},

		{
			"UpdateEeSubscription",
			http.MethodPatch,
			"/:ueIdentity/ee-subscriptions/:subscriptionId",
			s.HandleUpdateEeSubscription,
		},
	}
}

// HTTPCreateEeSubscription - Subscribe
func (s *Server) HandleCreateEeSubscription(c *gin.Context) {
	var eesubscription models.UdmEeEeSubscription

	requestBody, err := c.GetRawData()
	if err != nil {
		logger.EeLog.Errorf("Get Request Body error: %+v", err)
		problemDetail := models.ProblemDetails{
			Title:  "System failure",
			Status: http.StatusInternalServerError,
			Detail: err.Error(),
			Cause:  "SYSTEM_FAILURE",
		}
		c.JSON(http.StatusInternalServerError, problemDetail)
		return
	}

	err = openapi.Deserialize(&eesubscription, requestBody, "application/json")
	if err != nil {
		problemDetail := "[Request Body] " + err.Error()
		rsp := models.ProblemDetails{
			Title:  "Malformed request syntax",
			Status: http.StatusBadRequest,
			Detail: problemDetail,
		}
		logger.EeLog.Errorln(problemDetail)
		c.JSON(http.StatusBadRequest, rsp)
		return
	}

	logger.EeLog.Infoln("Handle Create EE Subscription")

	ueIdentity := c.Params.ByName("ueIdentity")

	s.Processor().CreateEeSubscriptionProcedure(c, ueIdentity, eesubscription)
}

func (s *Server) HandleDeleteEeSubscription(c *gin.Context) {
	ueIdentity := c.Params.ByName("ueIdentity")
	subscriptionID := c.Params.ByName("subscriptionId")

	s.Processor().DeleteEeSubscriptionProcedure(c, ueIdentity, subscriptionID)
}

func (s *Server) HandleUpdateEeSubscription(c *gin.Context) {
	var patchList []models.PatchItem

	requestBody, err := c.GetRawData()
	if err != nil {
		logger.EeLog.Errorf("Get Request Body error: %+v", err)
		problemDetail := models.ProblemDetails{
			Title:  "System failure",
			Status: http.StatusInternalServerError,
			Detail: err.Error(),
			Cause:  "SYSTEM_FAILURE",
		}
		c.JSON(http.StatusInternalServerError, problemDetail)
		return
	}

	err = openapi.Deserialize(&patchList, requestBody, "application/json")
	if err != nil {
		problemDetail := "[Request Body] " + err.Error()
		rsp := models.ProblemDetails{
			Title:  "Malformed request syntax",
			Status: http.StatusBadRequest,
			Detail: problemDetail,
		}
		logger.EeLog.Errorln(problemDetail)
		c.JSON(http.StatusBadRequest, rsp)
		return
	}

	ueIdentity := c.Params.ByName("ueIdentity")
	subscriptionID := c.Params.ByName("subscriptionId")

	logger.EeLog.Infoln("Handle Update EE subscription")
	logger.EeLog.Warnln("Update EE Subscription is not implemented")

	s.Processor().UpdateEeSubscriptionProcedure(c, ueIdentity, subscriptionID, patchList)
}

func (s *Server) HandleIndex(c *gin.Context) {
	c.String(http.StatusOK, "Hello World!")
}
