package sbi

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/free5gc/openapi"
	"github.com/free5gc/openapi/models"
	"github.com/free5gc/pcf/internal/logger"
	"github.com/free5gc/pcf/internal/util"
)

func (s *Server) getHttpCallBackRoutes() []Route {
	return []Route{
		{
			Method:  http.MethodPost,
			Pattern: "/nudr-notify/policy-data/:supi",
			APIFunc: s.HTTPUdrPolicyDataChangeNotify,
		},
		{
			Method:  http.MethodPost,
			Pattern: "/nudr-notify/influence-data/:supi/:pduSessionId",
			APIFunc: s.HTTPUdrInfluenceDataUpdateNotify,
		},
		{
			Method:  http.MethodPost,
			Pattern: "/amfstatus",
			APIFunc: s.HTTPAmfStatusChangeNotify,
		},
	}
}

// amf_status_change
func (s *Server) HTTPAmfStatusChangeNotify(c *gin.Context) {
	var amfStatusChangeNotification models.AmfStatusChangeNotification

	requestBody, err := c.GetRawData()
	if err != nil {
		problemDetail := models.ProblemDetails{
			Title:  "System failure",
			Status: http.StatusInternalServerError,
			Detail: err.Error(),
			Cause:  "SYSTEM_FAILURE",
		}
		logger.CallbackLog.Errorf("Get Request Body error: %+v", err)
		c.JSON(http.StatusInternalServerError, problemDetail)
		return
	}

	err = openapi.Deserialize(&amfStatusChangeNotification, requestBody, "application/json")
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

	s.Processor().HandleAmfStatusChangeNotify(c, amfStatusChangeNotification)
}

// sm_policy_notify
// Nudr-Notify-smpolicy
func (s *Server) HTTPUdrPolicyDataChangeNotify(c *gin.Context) {
	var policyDataChangeNotification models.PolicyDataChangeNotification

	requestBody, err := c.GetRawData()
	if err != nil {
		problemDetail := models.ProblemDetails{
			Title:  "System failure",
			Status: http.StatusInternalServerError,
			Detail: err.Error(),
			Cause:  "SYSTEM_FAILURE",
		}
		logger.CallbackLog.Errorf("Get Request Body error: %+v", err)
		c.JSON(http.StatusInternalServerError, problemDetail)
		return
	}

	err = openapi.Deserialize(&policyDataChangeNotification, requestBody, "application/json")
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

	supi := c.Params.ByName("supi")
	if supi == "" {
		problemDetails := &models.ProblemDetails{
			Title:  util.ERROR_INITIAL_PARAMETERS,
			Status: http.StatusBadRequest,
		}
		c.JSON(http.StatusBadRequest, problemDetails)
		return
	}
	s.Processor().HandlePolicyDataChangeNotify(c, supi, policyDataChangeNotification)
}

// Influence Data Update Notification
func (s *Server) HTTPUdrInfluenceDataUpdateNotify(c *gin.Context) {
	var trafficInfluDataNotif []models.TrafficInfluDataNotif

	requestBody, err := c.GetRawData()
	if err != nil {
		problemDetail := models.ProblemDetails{
			Title:  "System failure",
			Status: http.StatusInternalServerError,
			Detail: err.Error(),
			Cause:  "SYSTEM_FAILURE",
		}
		logger.CallbackLog.Errorf("Get Request Body error: %+v", err)
		c.JSON(http.StatusInternalServerError, problemDetail)
		return
	}

	err = openapi.Deserialize(&trafficInfluDataNotif, requestBody, "application/json")
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

	supi := c.Params.ByName("supi")
	pduSessionId := c.Params.ByName("pduSessionId")
	if supi == "" || pduSessionId == "" {
		problemDetails := &models.ProblemDetails{
			Title:  util.ERROR_INITIAL_PARAMETERS,
			Status: http.StatusBadRequest,
		}
		c.JSON(http.StatusBadRequest, problemDetails)
		return
	}
	s.Processor().HandleInfluenceDataUpdateNotify(c, supi, pduSessionId, trafficInfluDataNotif)
}
