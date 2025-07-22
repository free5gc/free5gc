package sbi

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/free5gc/openapi"
	"github.com/free5gc/openapi/models"
	"github.com/free5gc/pcf/internal/logger"
	"github.com/free5gc/pcf/internal/util"
)

func (s *Server) HTTPDeleteIndividualAMPolicyAssociation(c *gin.Context) {
	polAssoId, _ := c.Params.Get("polAssoId")

	if polAssoId == "" {
		problemDetails := &models.ProblemDetails{
			Title:  util.ERROR_INITIAL_PARAMETERS,
			Status: http.StatusBadRequest,
		}
		c.JSON(http.StatusBadRequest, problemDetails)
		return
	}

	s.Processor().HandleDeletePoliciesPolAssoId(c, polAssoId)
}

func (s *Server) HTTPReadIndividualAMPolicyAssociation(c *gin.Context) {
	polAssoId, _ := c.Params.Get("polAssoId")

	if polAssoId == "" {
		problemDetails := &models.ProblemDetails{
			Title:  util.ERROR_INITIAL_PARAMETERS,
			Status: http.StatusBadRequest,
		}
		c.JSON(http.StatusBadRequest, problemDetails)
		return
	}

	s.Processor().HandleGetPoliciesPolAssoId(c, polAssoId)
}

func (s *Server) HTTPReportObservedEventTriggersForIndividualAMPolicyAssociation(c *gin.Context) {
	var policyAssociationUpdateRequest models.PcfAmPolicyControlPolicyAssociationUpdateRequest

	requestBody, err := c.GetRawData()
	if err != nil {
		problemDetail := models.ProblemDetails{
			Title:  "System failure",
			Status: http.StatusInternalServerError,
			Detail: err.Error(),
			Cause:  "SYSTEM_FAILURE",
		}
		logger.AmPolicyLog.Errorf("Get Request Body error: %+v", err)
		c.JSON(http.StatusInternalServerError, problemDetail)
		return
	}

	err = openapi.Deserialize(&policyAssociationUpdateRequest, requestBody, "application/json")
	if err != nil {
		problemDetail := "[Request Body] " + err.Error()
		rsp := models.ProblemDetails{
			Title:  "Malformed request syntax",
			Status: http.StatusBadRequest,
			Detail: problemDetail,
		}
		logger.AmPolicyLog.Errorln(problemDetail)
		c.JSON(http.StatusBadRequest, rsp)
		return
	}

	polAssoId, _ := c.Params.Get("polAssoId")

	if polAssoId == "" {
		problemDetails := &models.ProblemDetails{
			Title:  util.ERROR_INITIAL_PARAMETERS,
			Status: http.StatusBadRequest,
		}
		c.JSON(http.StatusBadRequest, problemDetails)
		return
	}

	s.Processor().HandleUpdatePostPoliciesPolAssoId(c, polAssoId, policyAssociationUpdateRequest)
}

func (s *Server) HTTPCreateIndividualAMPolicyAssociation(c *gin.Context) {
	var policyAssociationRequest models.PcfAmPolicyControlPolicyAssociationRequest
	requestBody, err := c.GetRawData()
	if err != nil {
		problemDetail := models.ProblemDetails{
			Title:  "System failure",
			Status: http.StatusInternalServerError,
			Detail: err.Error(),
			Cause:  "SYSTEM_FAILURE",
		}
		logger.AmPolicyLog.Errorf("Get Request Body error: %+v", err)
		c.JSON(http.StatusInternalServerError, problemDetail)
		return
	}

	err = openapi.Deserialize(&policyAssociationRequest, requestBody, "application/json")
	if err != nil {
		problemDetail := "[Request Body] " + err.Error()
		rsp := models.ProblemDetails{
			Title:  "Malformed request syntax",
			Status: http.StatusBadRequest,
			Detail: problemDetail,
		}
		logger.AmPolicyLog.Errorln(problemDetail)
		c.JSON(http.StatusBadRequest, rsp)
		return
	}

	if policyAssociationRequest.Supi == "" || policyAssociationRequest.NotificationUri == "" {
		rsp := util.GetProblemDetail("Missing Mandatory IE", util.ERROR_REQUEST_PARAMETERS)
		logger.AmPolicyLog.Errorln(rsp.Detail)
		c.JSON(int(rsp.Status), rsp)
		return
	}

	polAssoId, _ := c.Params.Get("polAssoId")

	s.Processor().HandlePostPolicies(c, polAssoId, policyAssociationRequest)
}

func (s *Server) getAmPolicyRoutes() []Route {
	return []Route{
		{
			Name:    "ReadIndividualAMPolicyAssociation",
			Method:  http.MethodGet,
			Pattern: "/policies/:polAssoId",
			APIFunc: s.HTTPReadIndividualAMPolicyAssociation,
		},
		{
			Name:    "DeleteIndividualAMPolicyAssociation",
			Method:  http.MethodDelete,
			Pattern: "/policies/:polAssoId",
			APIFunc: s.HTTPDeleteIndividualAMPolicyAssociation,
		},
		{
			Name:    "ReportObservedEventTriggersForIndividualAMPolicyAssociation",
			Method:  http.MethodPost,
			Pattern: "/policies/:polAssoId/update",
			APIFunc: s.HTTPReportObservedEventTriggersForIndividualAMPolicyAssociation,
		},
		{
			Name:    "CreateIndividualAMPolicyAssociation",
			Method:  http.MethodPost,
			Pattern: "/policies",
			APIFunc: s.HTTPCreateIndividualAMPolicyAssociation,
		},
	}
}
