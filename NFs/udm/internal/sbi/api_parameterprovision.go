package sbi

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/free5gc/openapi"
	"github.com/free5gc/openapi/models"
	"github.com/free5gc/udm/internal/logger"
)

func (s *Server) getParameterProvisionRoutes() []Route {
	return []Route{
		{
			"Index",
			http.MethodGet,
			"/",
			s.HandleIndex,
		},

		{
			"Update",
			http.MethodPatch,
			"/:ueId/pp-data",
			s.HandleUpdate,
		},

		{
			"Create5GMBSGroup",
			http.MethodPut,
			"/mbs-group-membership/:extGroupId",
			s.HandleCreate5GMBSGroup,
		},

		{
			"Create5GVNGroup",
			http.MethodPut,
			"/5g-vn-groups/:extGroupId",
			s.HandleCreate5GVNGroup,
		},

		{
			"CreatePPDataEntry",
			http.MethodPut,
			"/:ueId/pp-data-store/:afInstanceId",
			s.HandleCreatePPDataEntry,
		},

		{
			"Delete5GMBSGroup",
			http.MethodDelete,
			"/mbs-group-membership/:extGroupId",
			s.HandleDelete5GMBSGroup,
		},

		{
			"Delete5GVNGroup",
			http.MethodDelete,
			"/5g-vn-groups/:extGroupId",
			s.HandleDelete5GVNGroup,
		},

		{
			"DeletePPDataEntry",
			http.MethodDelete,
			"/:ueId/pp-data-store/:afInstanceId",
			s.HandleDeletePPDataEntry,
		},

		{
			"Get5GMBSGroup",
			http.MethodGet,
			"/mbs-group-membership/:extGroupId",
			s.HandleGet5GMBSGroup,
		},

		{
			"Get5GVNGroup",
			http.MethodGet,
			"/5g-vn-groups/:extGroupId",
			s.HandleGet5GVNGroup,
		},

		{
			"GetPPDataEntry",
			http.MethodGet,
			"/:ueId/pp-data-store/:afInstanceId",
			s.HandleGetPPDataEntry,
		},

		{
			"Modify5GMBSGroup",
			http.MethodPatch,
			"/mbs-group-membership/:extGroupId",
			s.HandleModify5GMBSGroup,
		},

		{
			"Modify5GVNGroup",
			http.MethodPatch,
			"/5g-vn-groups/:extGroupId",
			s.HandleModify5GVNGroup,
		},
	}
}

func (s *Server) HandleUpdate(c *gin.Context) {
	var ppDataReq models.PpData

	// step 1: retrieve http request body
	requestBody, err := c.GetRawData()
	if err != nil {
		problemDetail := models.ProblemDetails{
			Title:  "System failure",
			Status: http.StatusInternalServerError,
			Detail: err.Error(),
			Cause:  "SYSTEM_FAILURE",
		}
		logger.PpLog.Errorf("Get Request Body error: %+v", err)
		c.JSON(http.StatusInternalServerError, problemDetail)
		return
	}

	// step 2: convert requestBody to openapi models
	err = openapi.Deserialize(&ppDataReq, requestBody, "application/json")
	if err != nil {
		problemDetail := "[Request Body] " + err.Error()
		rsp := models.ProblemDetails{
			Title:  "Malformed request syntax",
			Status: http.StatusBadRequest,
			Detail: problemDetail,
		}
		logger.PpLog.Errorln(problemDetail)
		c.JSON(http.StatusBadRequest, rsp)
		return
	}

	gpsi := c.Params.ByName("ueId")
	if gpsi == "" {
		problemDetails := &models.ProblemDetails{
			Status: http.StatusBadRequest,
			Cause:  "NO_GPSI",
		}
		c.JSON(int(problemDetails.Status), problemDetails)
		return
	}

	logger.PpLog.Infoln("Handle UpdateRequest")

	// step 3: handle the message
	s.Processor().UpdateProcedure(c, ppDataReq, gpsi)
}

func (s *Server) HandleCreate5GMBSGroup(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{})
}

func (s *Server) HandleCreate5GVNGroup(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{})
}

func (s *Server) HandleCreatePPDataEntry(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{})
}

func (s *Server) HandleDelete5GMBSGroup(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{})
}

func (s *Server) HandleDelete5GVNGroup(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{})
}

func (s *Server) HandleDeletePPDataEntry(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{})
}

func (s *Server) HandleGet5GMBSGroup(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{})
}

func (s *Server) HandleGet5GVNGroup(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{})
}

func (s *Server) HandleGetPPDataEntry(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{})
}

func (s *Server) HandleModify5GMBSGroup(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{})
}

func (s *Server) HandleModify5GVNGroup(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{})
}
