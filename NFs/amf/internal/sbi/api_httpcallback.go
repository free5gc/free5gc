package sbi

import (
	"net/http"

	"github.com/gin-gonic/gin"

	amf_context "github.com/free5gc/amf/internal/context"
	"github.com/free5gc/amf/internal/logger"
	"github.com/free5gc/openapi"
	"github.com/free5gc/openapi/models"
)

func (s *Server) getHttpCallBackRoutes() []Route {
	return []Route{
		{
			Method:  http.MethodGet,
			Pattern: "/",
			APIFunc: func(c *gin.Context) {
				c.String(http.StatusOK, "Hello World!")
			},
		},
		{
			Name:    "AmPolicyControlUpdateNotifyUpdate",
			Method:  http.MethodPost,
			Pattern: "/am-policy/:polAssoId/update",
			APIFunc: s.HTTPAmPolicyControlUpdateNotifyUpdate,
		},
		{
			Name:    "AmPolicyControlUpdateNotifyTerminate",
			Method:  http.MethodPost,
			Pattern: "/am-policy/:polAssoId/terminate",
			APIFunc: s.HTTPAmPolicyControlUpdateNotifyTerminate,
		},
		{
			Name:    "SmContextStatusNotify",
			Method:  http.MethodPost,
			Pattern: "/smContextStatus/:supi/:pduSessionId",
			APIFunc: s.HTTPSmContextStatusNotify,
		},
		{
			Name:    "N1MessageNotify",
			Method:  http.MethodPost,
			Pattern: "/n1-message-notify",
			APIFunc: s.HTTPN1MessageNotify,
		},
		{
			Name:    "HandleDeregistrationNotification",
			Method:  http.MethodPost,
			Pattern: "/deregistration/:ueid",
			APIFunc: s.HTTPHandleDeregistrationNotification,
		},
	}
}

func (s *Server) HTTPAmPolicyControlUpdateNotifyUpdate(c *gin.Context) {
	var policyUpdate models.PcfAmPolicyControlPolicyUpdate

	requestBody, err := c.GetRawData()
	if err != nil {
		logger.CallbackLog.Errorf("Get Request Body error: %+v", err)
		problemDetail := models.ProblemDetails{
			Title:  "System failure",
			Status: http.StatusInternalServerError,
			Detail: err.Error(),
			Cause:  "SYSTEM_FAILURE",
		}
		c.JSON(http.StatusInternalServerError, problemDetail)
		return
	}

	err = openapi.Deserialize(&policyUpdate, requestBody, "application/json")
	if err != nil {
		problemDetail := reqbody + err.Error()
		rsp := models.ProblemDetails{
			Title:  "Malformed request syntax",
			Status: http.StatusBadRequest,
			Detail: problemDetail,
		}
		logger.CallbackLog.Errorln(problemDetail)
		c.JSON(http.StatusBadRequest, rsp)
		return
	}
	s.Processor().HandleAmPolicyControlUpdateNotifyUpdate(c, policyUpdate)
}

func (s *Server) HTTPAmPolicyControlUpdateNotifyTerminate(c *gin.Context) {
	var terminationNotification models.PcfAmPolicyControlTerminationNotification

	requestBody, err := c.GetRawData()
	if err != nil {
		logger.CallbackLog.Errorf("Get Request Body error: %+v", err)
		problemDetail := models.ProblemDetails{
			Title:  "System failure",
			Status: http.StatusInternalServerError,
			Detail: err.Error(),
			Cause:  "SYSTEM_FAILURE",
		}
		c.JSON(http.StatusInternalServerError, problemDetail)
		return
	}

	err = openapi.Deserialize(&terminationNotification, requestBody, "application/json")
	if err != nil {
		problemDetail := reqbody + err.Error()
		rsp := models.ProblemDetails{
			Title:  "Malformed request syntax",
			Status: http.StatusBadRequest,
			Detail: problemDetail,
		}
		logger.CallbackLog.Errorln(problemDetail)
		c.JSON(http.StatusBadRequest, rsp)
		return
	}
	s.Processor().HandleAmPolicyControlUpdateNotifyTerminate(c, terminationNotification)
}

func (s *Server) HTTPN1MessageNotify(c *gin.Context) {
	var n1MessageNotify models.N1MessageNotifyRequest
	n1MessageNotify.JsonData = new(models.N1MessageNotification)

	err := c.ShouldBindWith(&n1MessageNotify, openapi.MultipartRelatedBinding{})
	if err != nil {
		problemDetail := reqbody + err.Error()
		rsp := models.ProblemDetails{
			Title:  "Malformed request syntax",
			Status: http.StatusBadRequest,
			Detail: problemDetail,
		}
		logger.CallbackLog.Errorln(problemDetail)
		c.JSON(http.StatusBadRequest, rsp)
		return
	}
	s.Processor().HandleN1MessageNotify(c, n1MessageNotify)
}

func (s *Server) HTTPSmContextStatusNotify(c *gin.Context) {
	var smContextStatusNotification models.SmfPduSessionSmContextStatusNotification

	requestBody, err := c.GetRawData()
	if err != nil {
		logger.CallbackLog.Errorf("Get Request Body error: %+v", err)
		problemDetail := models.ProblemDetails{
			Title:  "System failure",
			Status: http.StatusInternalServerError,
			Detail: err.Error(),
			Cause:  "SYSTEM_FAILURE",
		}
		c.JSON(http.StatusInternalServerError, problemDetail)
		return
	}

	err = openapi.Deserialize(&smContextStatusNotification, requestBody, "application/json")
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

	s.Processor().HandleSmContextStatusNotify(c, smContextStatusNotification)
}

func (s *Server) HTTPHandleDeregistrationNotification(c *gin.Context) {
	// TS 23.502 - 4.2.2.2.2 - step 14d
	logger.CallbackLog.Traceln("Handle Deregistration Notification")

	var deregData models.DeregistrationData

	requestBody, err := c.GetRawData()
	if err != nil {
		logger.CallbackLog.Errorf("Get Request Body error: %+v", err)
		problemDetails := models.ProblemDetails{
			Title:  "System failure",
			Status: http.StatusInternalServerError,
			Detail: err.Error(),
			Cause:  "SYSTEM_FAILURE",
		}
		c.JSON(http.StatusInternalServerError, problemDetails)
		return
	}

	err = openapi.Deserialize(&deregData, requestBody, "application/json")
	if err != nil {
		problemDetails := models.ProblemDetails{
			Title:  "Malformed request syntax",
			Status: http.StatusBadRequest,
			Detail: reqbody + err.Error(),
		}
		logger.CallbackLog.Errorln(problemDetails.Detail)
		c.JSON(http.StatusBadRequest, problemDetails)
		return
	}

	ueid := c.Param("ueid")
	ue, ok := amf_context.GetSelf().AmfUeFindByUeContextID(ueid)
	if !ok {
		logger.CallbackLog.Errorf("AmfUe Context[%s] not found", ueid)
		problemDetails := models.ProblemDetails{
			Status: http.StatusNotFound,
			Cause:  "CONTEXT_NOT_FOUND",
		}
		c.JSON(http.StatusNotFound, problemDetails)
		return
	}

	problemDetails, err := s.DeregistrationNotificationProcedure(ue, deregData)
	if problemDetails != nil {
		ue.GmmLog.Errorf("Deregistration Notification Procedure Failed Problem[%+v]", problemDetails)
	} else if err != nil {
		ue.GmmLog.Errorf("Deregistration Notification Procedure Error[%v]", err.Error())
	}
	// TS 23.503 - 5.3.2.3.2 UDM initiated NF Deregistration
	// The AMF acknowledges the Nudm_UECM_DeRegistrationNotification to the UDM.
	c.JSON(http.StatusNoContent, nil)
}

// TS 23.502 - 4.2.2.3.3 Network-initiated Deregistration
// The AMF can initiate this procedure for either explicit (e.g. by O&M intervention) or
// implicit (e.g. expiring of Implicit Deregistration timer)
func (s *Server) DeregistrationNotificationProcedure(ue *amf_context.AmfUe, deregData models.DeregistrationData) (
	problemDetails *models.ProblemDetails, err error,
) {
	// The AMF does not send the Deregistration Request message to the UE for Implicit Deregistration.
	if deregData.DeregReason == models.DeregistrationReason_UE_INITIAL_REGISTRATION {
		// TS 23.502 - 4.2.2.2.2 General Registration
		// Invokes the Nsmf_PDUSession_ReleaseSMContext for the corresponding access type
		ue.SmContextList.Range(func(key, value interface{}) bool {
			smContext := value.(*amf_context.SmContext)
			if smContext.AccessType() == deregData.AccessType {
				problemDetails, err = s.Consumer().SendReleaseSmContextRequest(ue, smContext, nil, "", nil)
				if problemDetails != nil {
					ue.GmmLog.Errorf("Release SmContext Failed Problem[%+v]", problemDetails)
				} else if err != nil {
					ue.GmmLog.Errorf("Release SmContext Error[%v]", err.Error())
				}
			}
			return true
		})
	}
	// TS 23.502 - 4.2.2.2.2 General Registration - 14e
	// TODO: (R16) If old AMF does not have UE context for another access type (i.e. non-3GPP access),
	// the Old AMF unsubscribes with the UDM for subscription data using Nudm_SDM_unsubscribe
	if ue.SdmSubscriptionId != "" {
		problemDetails, err = s.Consumer().SDMUnsubscribe(ue)
		if problemDetails != nil {
			logger.GmmLog.Errorf("SDM Unubscribe Failed Problem[%+v]", problemDetails)
		} else if err != nil {
			logger.GmmLog.Errorf("SDM Unubscribe Error[%+v]", err)
		}
		ue.SdmSubscriptionId = ""
	}

	// TS 23.502 - 4.2.2.2.2 General Registration - 20 AMF-Initiated Policy Association Termination
	// For UE_INITIAL_REGISTRATION and SUBSCRIPTION_WITHDRAW, do AMF-Initiated Policy Association Termination directly.
	if ue.PolicyAssociationId != "" {
		// TODO: For REGISTRATION_AREA_CHANGE, old AMF performs an AMF-initiated Policy Association Termination
		// procedure if the old AMF has established an AM Policy Association and a UE Policy Association with the PCF(s)
		// and the old AMF did not transfer the PCF ID(s) to the new AMF. (Ref: TS 23.502 - 4.2.2.2.2)
		// Currently, old AMF will transfer the PCF ID but new AMF will not utilize the PCF ID
		problemDetailsAMPolicyControlDelete, errAMPolicyControlDelete := s.Consumer().AMPolicyControlDelete(ue)
		if problemDetailsAMPolicyControlDelete != nil {
			logger.GmmLog.Errorf("Delete AM policy Failed Problem[%+v]", problemDetailsAMPolicyControlDelete)
		} else if err != nil {
			logger.GmmLog.Errorf("Delete AM policy Error[%+v]", errAMPolicyControlDelete)
		}
	}

	// The old AMF should clean the UE context
	// TODO: (R16) Only remove the target access UE context
	ue.Remove()

	return nil, nil
}
