package sbi

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/free5gc/amf/internal/logger"
	"github.com/free5gc/openapi"
	"github.com/free5gc/openapi/models"
)

func (s *Server) getEventexposureRoutes() []Route {
	return []Route{
		{
			Method:  http.MethodGet,
			Pattern: "/",
			APIFunc: func(c *gin.Context) {
				c.String(http.StatusOK, "Hello World!")
			},
		},
		{
			Name:    "DeleteSubscription",
			Method:  http.MethodDelete,
			Pattern: "/subscriptions/:subscriptionId",
			APIFunc: s.HTTPDeleteSubscription,
		},
		{
			Name:    "ModifySubscription",
			Method:  http.MethodPatch,
			Pattern: "/subscriptions/:subscriptionId",
			APIFunc: s.HTTPModifySubscription,
		},
		{
			Name:    "CreateSubscription",
			Method:  http.MethodPost,
			Pattern: "/subscriptions",
			APIFunc: s.HTTPCreateSubscription,
		},
	}
}

// DeleteSubscription - Namf_EventExposure Unsubscribe service Operation
func (s *Server) HTTPDeleteSubscription(c *gin.Context) {
	s.Processor().HandleDeleteAMFEventSubscription(c)
}

// ModifySubscription - Namf_EventExposure Subscribe Modify service Operation
func (s *Server) HTTPModifySubscription(c *gin.Context) {
	var modifySubscriptionRequest models.ModifySubscriptionRequest

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

	err = openapi.Deserialize(&modifySubscriptionRequest, requestBody, "application/json")
	if err != nil {
		problemDetail := reqbody + err.Error()
		rsp := models.ProblemDetails{
			Title:  "Malformed request syntax",
			Status: http.StatusBadRequest,
			Detail: problemDetail,
		}
		logger.EeLog.Errorln(problemDetail)
		c.JSON(http.StatusBadRequest, rsp)
		return
	}
	s.Processor().HandleModifyAMFEventSubscription(c, modifySubscriptionRequest)
}

func (s *Server) HTTPCreateSubscription(c *gin.Context) {
	var createEventSubscription models.AmfCreateEventSubscription

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

	err = openapi.Deserialize(&createEventSubscription, requestBody, "application/json")
	if err != nil {
		problemDetail := reqbody + err.Error()
		rsp := models.ProblemDetails{
			Title:  "Malformed request syntax",
			Status: http.StatusBadRequest,
			Detail: problemDetail,
		}
		logger.EeLog.Errorln(problemDetail)
		c.JSON(http.StatusBadRequest, rsp)
		return
	}
	s.Processor().HandleCreateAMFEventSubscription(c, createEventSubscription)
}
