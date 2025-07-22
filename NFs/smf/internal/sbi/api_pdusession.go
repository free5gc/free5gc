package sbi

import (
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/free5gc/openapi"
	"github.com/free5gc/openapi/models"
	"github.com/free5gc/smf/internal/logger"
)

func (s *Server) getPDUSessionRoutes() []Route {
	return []Route{
		{
			Name:    "Index",
			Method:  http.MethodGet,
			Pattern: "/",
			APIFunc: func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"status": "Service Available"})
			},
		},
		{
			Name:    "PostSmContexts",
			Method:  http.MethodPost,
			Pattern: "/sm-contexts",
			APIFunc: s.HTTPPostSmContexts,
		},
		{
			Name:    "UpdateSmContext",
			Method:  http.MethodPost,
			Pattern: "/sm-contexts/:smContextRef/modify",
			APIFunc: s.HTTPUpdateSmContext,
		},
		{
			Name:    "RetrieveSmContext",
			Method:  http.MethodPost,
			Pattern: "/sm-contexts/:smContextRef/retrieve",
			APIFunc: s.HTTPRetrieveSmContext,
		},
		{
			Name:    "ReleaseSmContext",
			Method:  http.MethodPost,
			Pattern: "/sm-contexts/:smContextRef/release",
			APIFunc: s.HTTPReleaseSmContext,
		},
		{
			Name:    "SendMoData",
			Method:  http.MethodPost,
			Pattern: "/sm-contexts/:smContextRef/send-mo-data",
			APIFunc: s.HTTPSendMoData,
		},
		{
			Name:    "PostPduSessions",
			Method:  http.MethodPatch,
			Pattern: "/pdu-sessions",
			APIFunc: s.HTTPPostPduSessions,
		},
		{
			Name:    "UpdatePduSession",
			Method:  http.MethodPost,
			Pattern: "/pdu-sessions/:pduSessionRef/modify",
			APIFunc: s.HTTPUpdatePduSession,
		},
		{
			Name:    "ReleasePduSession",
			Method:  http.MethodPost,
			Pattern: "/pdu-sessions/:pduSessionRef/release",
			APIFunc: s.HTTPReleasePduSession,
		},
		{
			Name:    "RetrievePduSession",
			Method:  http.MethodPost,
			Pattern: "/pdu-sessions/:pduSessionRef/retrieve",
			APIFunc: s.HTTPRetrievePduSession,
		},
		{
			Name:    "TransferMoData",
			Method:  http.MethodPost,
			Pattern: "/pdu-sessions/:pduSessionRef/transfer-mo-data",
			APIFunc: s.HTTPTransferMoData,
		},
	}
}

// HTTPPostSmContexts - Create SM Context
func (s *Server) HTTPPostSmContexts(c *gin.Context) {
	logger.PduSessLog.Info("Receive Create SM Context Request")
	var request models.PostSmContextsRequest

	request.JsonData = new(models.SmfPduSessionSmContextCreateData)

	contentType := strings.Split(c.GetHeader("Content-Type"), ";")
	var err error
	switch contentType[0] {
	case APPLICATION_JSON:
		err = c.ShouldBindJSON(request.JsonData)
	case MULTIPART_RELATED:
		err = c.ShouldBindWith(&request, openapi.MultipartRelatedBinding{})
	}

	if err != nil {
		problemDetail := "[Request Body] " + err.Error()
		logger.PduSessLog.Errorln(problemDetail)
		c.JSON(http.StatusBadRequest, openapi.ProblemDetailsMalformedReqSyntax(problemDetail))
		return
	}

	isDone := c.Done()
	s.Processor().HandlePDUSessionSMContextCreate(c, request, isDone)
}

// HTTPUpdateSmContext - Update SM Context
func (s *Server) HTTPUpdateSmContext(c *gin.Context) {
	logger.PduSessLog.Info("Receive Update SM Context Request")
	var request models.UpdateSmContextRequest
	request.JsonData = new(models.SmfPduSessionSmContextUpdateData)

	contentType := strings.Split(c.GetHeader("Content-Type"), ";")
	var err error
	switch contentType[0] {
	case APPLICATION_JSON:
		err = c.ShouldBindJSON(request.JsonData)
	case MULTIPART_RELATED:
		err = c.ShouldBindWith(&request, openapi.MultipartRelatedBinding{})
	}
	if err != nil {
		log.Print(err)
		return
	}

	smContextRef := c.Params.ByName("smContextRef")
	s.Processor().HandlePDUSessionSMContextUpdate(c, request, smContextRef)
}

// HTTPRetrieveSmContext - Retrieve SM Context
func (s *Server) HTTPRetrieveSmContext(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{})
}

// HTTPReleaseSmContext - Release SM Context
func (s *Server) HTTPReleaseSmContext(c *gin.Context) {
	logger.PduSessLog.Info("Receive Release SM Context Request")
	var request models.ReleaseSmContextRequest
	request.JsonData = new(models.SmfPduSessionSmContextReleaseData)

	contentType := strings.Split(c.GetHeader("Content-Type"), ";")
	var err error
	switch contentType[0] {
	case APPLICATION_JSON:
		err = c.ShouldBindJSON(request.JsonData)
	case MULTIPART_RELATED:
		err = c.ShouldBindWith(&request, openapi.MultipartRelatedBinding{})
	}
	if err != nil {
		log.Print(err)
		return
	}

	smContextRef := c.Params.ByName("smContextRef")
	s.Processor().HandlePDUSessionSMContextRelease(c, request, smContextRef)
}

func (s *Server) HTTPSendMoData(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{})
}

// HTTPPostPduSessions - Create
func (s *Server) HTTPPostPduSessions(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{})
}

// HTTPUpdatePduSession - Update (initiated by V-SMF)
func (s *Server) HTTPUpdatePduSession(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{})
}

// HTTPReleasePduSession - Release
func (s *Server) HTTPReleasePduSession(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{})
}

func (s *Server) HTTPRetrievePduSession(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{})
}

func (s *Server) HTTPTransferMoData(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{})
}
