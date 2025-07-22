package sbi

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/free5gc/amf/internal/logger"
	"github.com/free5gc/openapi"
	"github.com/free5gc/openapi/models"
)

func Index(c *gin.Context) {
	c.String(http.StatusOK, "Hello World!")
}

func (s *Server) getCommunicationRoutes() []Route {
	return []Route{
		{
			Method:  http.MethodGet,
			Pattern: "/",
			APIFunc: func(c *gin.Context) {
				c.String(http.StatusOK, "Hello World!")
			},
		},
		{
			Name:    "AMFStatusChangeSubscribeModfy",
			Method:  http.MethodPut,
			Pattern: "/subscriptions/:subscriptionId",
			APIFunc: s.HTTPAMFStatusChangeSubscribeModify,
		},
		{
			Name:    "AMFStatusChangeUnSubscribe",
			Method:  http.MethodDelete,
			Pattern: "/subscriptions/:subscriptionId",
			APIFunc: s.HTTPAMFStatusChangeUnSubscribe,
		},
		{
			Name:    "CreateUEContext",
			Method:  http.MethodPut,
			Pattern: "/ue-contexts/:ueContextId",
			APIFunc: s.HTTPCreateUEContext,
		},
		{
			Name:    "EBIAssignment",
			Method:  http.MethodPost,
			Pattern: "/ue-contexts/:ueContextId/assign-ebi",
			APIFunc: s.HTTPEBIAssignment,
		},
		{
			Name:    "RegistrationStatusUpdate",
			Method:  http.MethodPost,
			Pattern: "/ue-contexts/:ueContextId/transfer-update",
			APIFunc: s.HTTPRegistrationStatusUpdate,
		},
		{
			Name:    "ReleaseUEContext",
			Method:  http.MethodPost,
			Pattern: "/ue-contexts/:ueContextId/release",
			APIFunc: s.HTTPReleaseUEContext,
		},
		{
			Name:    "UEContextTransfer",
			Method:  http.MethodPost,
			Pattern: "/ue-contexts/:ueContextId/transfer",
			APIFunc: s.HTTPUEContextTransfer,
		},
		{
			Name:    "RelocateUEContext",
			Method:  http.MethodPost,
			Pattern: "/ue-contexts/:ueContextId/relocate",
			APIFunc: s.HTTPRelocateUEContext,
		},
		{
			Name:    "CancelRelocateUEContext",
			Method:  http.MethodPost,
			Pattern: "/ue-contexts/:ueContextId/cancel-relocate",
			APIFunc: s.HTTPCancelRelocateUEContext,
		},
		{
			Name:    "N1N2MessageUnSubscribe",
			Method:  http.MethodDelete,
			Pattern: "/ue-contexts/:ueContextId/n1-n2-messages/subscriptions/:subscriptionId",
			APIFunc: s.HTTPN1N2MessageUnSubscribe,
		},
		{
			Name:    "N1N2MessageTransfer",
			Method:  http.MethodPost,
			Pattern: "/ue-contexts/:ueContextId/n1-n2-messages",
			APIFunc: s.HTTPN1N2MessageTransfer,
		},
		{
			Name:    "N1N2MessageTransferStatus",
			Method:  http.MethodGet,
			Pattern: "/ue-contexts/:ueContextId/n1-n2-messages/:n1N2MessageId",
			APIFunc: s.HTTPN1N2MessageTransferStatus,
		},
		{
			Name:    "N1N2MessageSubscribe",
			Method:  http.MethodPost,
			Pattern: "/ue-contexts/:ueContextId/n1-n2-messages/subscriptions",
			APIFunc: s.HTTPN1N2MessageSubscribe,
		},
		{
			Name:    "NonUeN2InfoUnSubscribe",
			Method:  http.MethodDelete,
			Pattern: "/non-ue-n2-messages/subscriptions/:n2NotifySubscriptionId",
			APIFunc: s.HTTPNonUeN2InfoUnSubscribe,
		},
		{
			Name:    "NonUeN2MessageTransfer",
			Method:  http.MethodPost,
			Pattern: "/non-ue-n2-messages/transfer",
			APIFunc: s.HTTPNonUeN2MessageTransfer,
		},
		{
			Name:    "NonUeN2InfoSubscribe",
			Method:  http.MethodPost,
			Pattern: "/non-ue-n2-messages/subscriptions",
			APIFunc: s.HTTPNonUeN2InfoSubscribe,
		},
		{
			Name:    "AMFStatusChangeSubscribe",
			Method:  http.MethodPost,
			Pattern: "/subscriptions",
			APIFunc: s.HTTPAMFStatusChangeSubscribe,
		},
	}
}

// AMFStatusChangeSubscribeModify - Namf_Communication AMF Status Change Subscribe Modify service Operation
func (s *Server) HTTPAMFStatusChangeSubscribeModify(c *gin.Context) {
	var subscriptionData models.AmfCommunicationSubscriptionData

	requestBody, err := c.GetRawData()
	if err != nil {
		logger.CommLog.Errorf("Get Request Body error: %+v", err)
		problemDetail := models.ProblemDetails{
			Title:  "System failure",
			Status: http.StatusInternalServerError,
			Detail: err.Error(),
			Cause:  "SYSTEM_FAILURE",
		}
		c.JSON(http.StatusInternalServerError, problemDetail)
		return
	}

	err = openapi.Deserialize(&subscriptionData, requestBody, applicationjson)
	if err != nil {
		problemDetail := reqbody + err.Error()
		rsp := models.ProblemDetails{
			Title:  "Malformed request syntax",
			Status: http.StatusBadRequest,
			Detail: problemDetail,
		}
		logger.CommLog.Errorln(problemDetail)
		c.JSON(http.StatusBadRequest, rsp)
		return
	}
	s.Processor().HandleAMFStatusChangeSubscribeModify(c, subscriptionData)
}

// AMFStatusChangeUnSubscribe - Namf_Communication AMF Status Change UnSubscribe service Operation
func (s *Server) HTTPAMFStatusChangeUnSubscribe(c *gin.Context) {
	s.Processor().HandleAMFStatusChangeUnSubscribeRequest(c)
}

func (s *Server) HTTPCreateUEContext(c *gin.Context) {
	var createUeContextRequest models.CreateUeContextRequest
	createUeContextRequest.JsonData = new(models.UeContextCreateData)

	requestBody, err := c.GetRawData()
	if err != nil {
		logger.CommLog.Errorf("Get Request Body error: %+v", err)
		problemDetail := models.ProblemDetails{
			Title:  "System failure",
			Status: http.StatusInternalServerError,
			Detail: err.Error(),
			Cause:  "SYSTEM_FAILURE",
		}
		c.JSON(http.StatusInternalServerError, problemDetail)
		return
	}

	contentType := c.GetHeader("Content-Type")
	str := strings.Split(contentType, ";")
	switch str[0] {
	case applicationjson:
		err = openapi.Deserialize(createUeContextRequest.JsonData, requestBody, contentType)
	case multipartrelate:
		err = openapi.Deserialize(&createUeContextRequest, requestBody, contentType)
	default:
		err = fmt.Errorf("wrong content type")
	}

	if err != nil {
		problemDetail := reqbody + err.Error()
		rsp := models.ProblemDetails{
			Title:  "Malformed request syntax",
			Status: http.StatusBadRequest,
			Detail: problemDetail,
		}
		logger.CommLog.Errorln(problemDetail)
		c.JSON(http.StatusBadRequest, rsp)
		return
	}
	s.Processor().HandleCreateUEContextRequest(c, createUeContextRequest)
}

// EBIAssignment - Namf_Communication EBI Assignment service Operation
func (s *Server) HTTPEBIAssignment(c *gin.Context) {
	var assignEbiData models.AssignEbiData

	requestBody, err := c.GetRawData()
	if err != nil {
		problemDetail := models.ProblemDetails{
			Title:  "System failure",
			Status: http.StatusInternalServerError,
			Detail: err.Error(),
			Cause:  "SYSTEM_FAILURE",
		}
		logger.CommLog.Errorf("Get Request Body error: %+v", err)
		c.JSON(http.StatusInternalServerError, problemDetail)
		return
	}

	err = openapi.Deserialize(&assignEbiData, requestBody, applicationjson)
	if err != nil {
		problemDetail := reqbody + err.Error()
		rsp := models.ProblemDetails{
			Title:  "Malformed request syntax",
			Status: http.StatusBadRequest,
			Detail: problemDetail,
		}
		logger.CommLog.Errorln(problemDetail)
		c.JSON(http.StatusBadRequest, rsp)
		return
	}
	s.Processor().HandleAssignEbiDataRequest(c, assignEbiData)
}

// RegistrationStatusUpdate - Namf_Communication RegistrationStatusUpdate service Operation
func (s *Server) HTTPRegistrationStatusUpdate(c *gin.Context) {
	var ueRegStatusUpdateReqData models.UeRegStatusUpdateReqData

	requestBody, err := c.GetRawData()
	if err != nil {
		logger.CommLog.Errorf("Get Request Body error: %+v", err)
		problemDetail := models.ProblemDetails{
			Title:  "System failure",
			Status: http.StatusInternalServerError,
			Detail: err.Error(),
			Cause:  "SYSTEM_FAILURE",
		}
		c.JSON(http.StatusInternalServerError, problemDetail)
		return
	}

	err = openapi.Deserialize(&ueRegStatusUpdateReqData, requestBody, applicationjson)
	if err != nil {
		problemDetail := reqbody + err.Error()
		rsp := models.ProblemDetails{
			Title:  "Malformed request syntax",
			Status: http.StatusBadRequest,
			Detail: problemDetail,
		}
		logger.CommLog.Errorln(problemDetail)
		c.JSON(http.StatusBadRequest, rsp)
		return
	}
	s.Processor().HandleRegistrationStatusUpdateRequest(c, ueRegStatusUpdateReqData)
}

// ReleaseUEContext - Namf_Communication ReleaseUEContext service Operation
func (s *Server) HTTPReleaseUEContext(c *gin.Context) {
	var ueContextRelease models.UeContextRelease

	requestBody, err := c.GetRawData()
	if err != nil {
		logger.CommLog.Errorf("Get Request Body error: %+v", err)
		problemDetail := models.ProblemDetails{
			Title:  "System failure",
			Status: http.StatusInternalServerError,
			Detail: err.Error(),
			Cause:  "SYSTEM_FAILURE",
		}
		c.JSON(http.StatusInternalServerError, problemDetail)
		return
	}

	err = openapi.Deserialize(&ueContextRelease, requestBody, applicationjson)
	if err != nil {
		problemDetail := reqbody + err.Error()
		rsp := models.ProblemDetails{
			Title:  "Malformed request syntax",
			Status: http.StatusBadRequest,
			Detail: problemDetail,
		}
		logger.CommLog.Errorln(problemDetail)
		c.JSON(http.StatusBadRequest, rsp)
		return
	}
	s.Processor().HandleReleaseUEContextRequest(c, ueContextRelease)
}

// UEContextTransfer - Namf_Communication UEContextTransfer service Operation
func (s *Server) HTTPUEContextTransfer(c *gin.Context) {
	var ueContextTransferRequest models.UeContextTransferRequest
	ueContextTransferRequest.JsonData = new(models.UeContextTransferReqData)

	requestBody, err := c.GetRawData()
	if err != nil {
		logger.CommLog.Errorf("Get Request Body error: %+v", err)
		problemDetail := models.ProblemDetails{
			Title:  "System failure",
			Status: http.StatusInternalServerError,
			Detail: err.Error(),
			Cause:  "SYSTEM_FAILURE",
		}
		c.JSON(http.StatusInternalServerError, problemDetail)
		return
	}

	contentType := c.GetHeader("Content-Type")
	str := strings.Split(contentType, ";")
	switch str[0] {
	case applicationjson:
		err = openapi.Deserialize(ueContextTransferRequest.JsonData, requestBody, contentType)
	case multipartrelate:
		err = openapi.Deserialize(&ueContextTransferRequest, requestBody, contentType)
	}

	if err != nil {
		problemDetail := reqbody + err.Error()
		rsp := models.ProblemDetails{
			Title:  "Malformed request syntax",
			Status: http.StatusBadRequest,
			Detail: problemDetail,
		}
		logger.CommLog.Errorln(problemDetail)
		c.JSON(http.StatusBadRequest, rsp)
		return
	}
	s.Processor().HandleUEContextTransferRequest(c, ueContextTransferRequest)
}

func (s *Server) HTTPRelocateUEContext(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{})
}

func (s *Server) HTTPCancelRelocateUEContext(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{})
}

func (s *Server) HTTPN1N2MessageUnSubscribe(c *gin.Context) {
	s.Processor().HandleN1N2MessageUnSubscribeRequest(c)
}

func (s *Server) HTTPN1N2MessageTransfer(c *gin.Context) {
	var n1n2MessageTransferRequest models.N1N2MessageTransferRequest
	n1n2MessageTransferRequest.JsonData = new(models.N1N2MessageTransferReqData)

	requestBody, err := c.GetRawData()
	if err != nil {
		problemDetail := models.ProblemDetails{
			Title:  "System failure",
			Status: http.StatusInternalServerError,
			Detail: err.Error(),
			Cause:  "SYSTEM_FAILURE",
		}
		logger.CommLog.Errorf("Get Request Body error: %+v", err)
		c.JSON(http.StatusInternalServerError, problemDetail)
		return
	}

	contentType := c.GetHeader("Content-Type")
	str := strings.Split(contentType, ";")
	switch str[0] {
	case applicationjson:
		err = fmt.Errorf("N1 and N2 datas are both Empty in N1N2MessgeTransfer")
	case multipartrelate:
		err = openapi.Deserialize(&n1n2MessageTransferRequest, requestBody, contentType)
	default:
		err = fmt.Errorf("wrong content type")
	}

	if err != nil {
		problemDetail := reqbody + err.Error()
		rsp := models.ProblemDetails{
			Title:  "Malformed request syntax",
			Status: http.StatusBadRequest,
			Detail: problemDetail,
		}
		logger.CommLog.Errorln(problemDetail)
		c.JSON(http.StatusBadRequest, rsp)
		return
	}
	s.Processor().HandleN1N2MessageTransferRequest(c, n1n2MessageTransferRequest)
}

func (s *Server) HTTPN1N2MessageTransferStatus(c *gin.Context) {
	s.Processor().HandleN1N2MessageTransferStatusRequest(c)
}

func (s *Server) HTTPN1N2MessageSubscribe(c *gin.Context) {
	var ueN1N2InfoSubscriptionCreateData models.UeN1N2InfoSubscriptionCreateData

	requestBody, err := c.GetRawData()
	if err != nil {
		logger.CommLog.Errorf("Get Request Body error: %+v", err)
		problemDetail := models.ProblemDetails{
			Title:  "System failure",
			Status: http.StatusInternalServerError,
			Detail: err.Error(),
			Cause:  "SYSTEM_FAILURE",
		}
		c.JSON(http.StatusInternalServerError, problemDetail)
		return
	}

	err = openapi.Deserialize(&ueN1N2InfoSubscriptionCreateData, requestBody, applicationjson)
	if err != nil {
		problemDetail := reqbody + err.Error()
		rsp := models.ProblemDetails{
			Title:  "Malformed request syntax",
			Status: http.StatusBadRequest,
			Detail: problemDetail,
		}
		logger.CommLog.Errorln(problemDetail)
		c.JSON(http.StatusBadRequest, rsp)
		return
	}
	s.Processor().HandleN1N2MessageSubscribeRequest(c, ueN1N2InfoSubscriptionCreateData)
}

func (s *Server) HTTPNonUeN2InfoUnSubscribe(c *gin.Context) {
	logger.CommLog.Warnf("Handle Non Ue N2 Info UnSubscribe is not implemented.")
	c.JSON(http.StatusNotImplemented, gin.H{})
}

func (s *Server) HTTPNonUeN2MessageTransfer(c *gin.Context) {
	logger.CommLog.Warnf("Handle Non Ue N2 Message Transfer is not implemented.")
	c.JSON(http.StatusNotImplemented, gin.H{})
}

func (s *Server) HTTPNonUeN2InfoSubscribe(c *gin.Context) {
	logger.CommLog.Warnf("Handle Non Ue N2 Info Subscribe is not implemented.")
	c.JSON(http.StatusNotImplemented, gin.H{})
}

func (s *Server) HTTPAMFStatusChangeSubscribe(c *gin.Context) {
	var subscriptionData models.AmfCommunicationSubscriptionData

	requestBody, err := c.GetRawData()
	if err != nil {
		logger.CommLog.Errorf("Get Request Body error: %+v", err)
		problemDetail := models.ProblemDetails{
			Title:  "System failure",
			Status: http.StatusInternalServerError,
			Detail: err.Error(),
			Cause:  "SYSTEM_FAILURE",
		}
		c.JSON(http.StatusInternalServerError, problemDetail)
		return
	}

	err = openapi.Deserialize(&subscriptionData, requestBody, applicationjson)
	if err != nil {
		problemDetail := reqbody + err.Error()
		rsp := models.ProblemDetails{
			Title:  "Malformed request syntax",
			Status: http.StatusBadRequest,
			Detail: problemDetail,
		}
		logger.CommLog.Errorln(problemDetail)
		c.JSON(http.StatusBadRequest, rsp)
		return
	}
	s.Processor().HandleAMFStatusChangeSubscribeRequest(c, subscriptionData)
}
