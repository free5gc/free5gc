package sbi

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/free5gc/openapi"
	"github.com/free5gc/openapi/models"
	Nudr_DataRepository "github.com/free5gc/openapi/udr/DataRepository"
	"github.com/free5gc/udm/internal/logger"
)

func (s *Server) getUEContextManagementRoutes() []Route {
	return []Route{
		{
			"Index",
			http.MethodGet,
			"/",
			s.HandleIndex,
		},

		{
			"Get3GppRegistration",
			http.MethodGet,
			"/:ueId/registrations/amf-3gpp-access",
			s.HandleGetAmf3gppAccess,
		},

		{
			"GetNon3GppRegistration",
			http.MethodGet,
			"/:ueId/registrations/amf-non-3gpp-access",
			s.HandleGetAmfNon3gppAccess,
		},

		{
			"Call3GppRegistration",
			http.MethodPut,
			"/:ueId/registrations/amf-3gpp-access",
			s.HandleRegistrationAmf3gppAccess,
		},

		{
			"Non3GppRegistration",
			http.MethodPut,
			"/:ueId/registrations/amf-non-3gpp-access",
			s.HandleRegistrationAmfNon3gppAccess,
		},

		{
			"Update3GppRegistration",
			http.MethodPatch,
			"/:ueId/registrations/amf-3gpp-access",
			s.HandleUpdateAmf3gppAccess,
		},

		{
			"UpdateNon3GppRegistration",
			http.MethodPatch,
			"/:ueId/registrations/amf-non-3gpp-access",
			s.HandleUpdateAmfNon3gppAccess,
		},

		{
			"SmfDeregistration",
			http.MethodDelete,
			"/:ueId/registrations/smf-registrations/:pduSessionId",
			s.HandleDeregistrationSmfRegistrations,
		},

		{
			"Registration",
			http.MethodPut,
			"/:ueId/registrations/smf-registrations/:pduSessionId",
			s.HandleRegistrationSmfRegistrations,
		},

		{
			"Get3GppSmsfRegistration",
			http.MethodGet,
			"/:ueId/registrations/smsf-3gpp-access",
			s.HandleGetSmsf3gppAccess,
		},

		{
			"Call3GppSmsfDeregistration",
			http.MethodDelete,
			"/:ueId/registrations/smsf-3gpp-access",
			s.HandleDeregistrationSmsf3gppAccess,
		},

		{
			"Non3GppSmsfDeregistration",
			http.MethodDelete,
			"/:ueId/registrations/smsf-non-3gpp-access",
			s.HandleDeregistrationSmsfNon3gppAccess,
		},

		{
			"GetNon3GppSmsfRegistration",
			http.MethodGet,
			"/:ueId/registrations/smsf-non-3gpp-access",
			s.HandleGetSmsfNon3gppAccess,
		},

		{
			"Call3GppSmsfRegistration",
			http.MethodPut,
			"/:ueId/registrations/smsf-3gpp-access",
			s.HandleUpdateSMSFReg3GPP,
		},

		{
			"Non3GppSmsfRegistration",
			http.MethodPut,
			"/:ueId/registrations/smsf-non-3gpp-access",
			s.HandleRegistrationSmsfNon3gppAccess,
		},

		{
			"DeregAMF",
			http.MethodPost,
			"/:ueId/registrations/amf-3gpp-access/dereg-amf",
			s.HandleDeregAMF,
		},

		{
			"GetIpSmGwRegistration",
			http.MethodGet,
			"/:ueId/registrations/ip-sm-gw",
			s.HandleGetIpSmGwRegistration,
		},

		{
			"GetLocationInfo",
			http.MethodGet,
			"/:ueId/registrations/location",
			s.HandleGetLocationInfo,
		},

		{
			"GetNwdafRegistration",
			http.MethodGet,
			"/:ueId/registrations/nwdaf-registrations",
			s.HandleGetNwdafRegistration,
		},

		{
			"GetRegistrations",
			http.MethodGet,
			"/:ueId/registrations",
			s.HandleGetRegistrations,
		},

		{
			"GetSmfRegistration",
			http.MethodGet,
			"/:ueId/registrations/smf-registrations",
			s.HandleGetSmfRegistration,
		},

		{
			"IpSmGwDeregistration",
			http.MethodDelete,
			"/:ueId/registrations/ip-sm-gw",
			s.HandleIpSmGwDeregistration,
		},

		{
			"IpSmGwRegistration",
			http.MethodPut,
			"/:ueId/registrations/ip-sm-gw",
			s.HandleIpSmGwRegistration,
		},

		{
			"NwdafDeregistration",
			http.MethodDelete,
			"/:ueId/registrations/nwdaf-registrations/:nwdafRegistrationId",
			s.HandleNwdafDeregistration,
		},

		{
			"NwdafRegistration",
			http.MethodPut,
			"/:ueId/registrations/nwdaf-registrations/:nwdafRegistrationId",
			s.HandleNwdafRegistration,
		},

		{
			"PeiUpdate",
			http.MethodPost,
			"/:ueId/registrations/amf-3gpp-access/pei-update",
			s.HandlePeiUpdate,
		},

		{
			"RetrieveSmfRegistration",
			http.MethodGet,
			"/:ueId/registrations/smf-registrations/:pduSessionId",
			s.HandleRetrieveSmfRegistration,
		},

		{
			"SendRoutingInfoSm",
			http.MethodPost,
			"/:ueId/registrations/send-routing-info-sm",
			s.HandleSendRoutingInfoSm,
		},

		{
			"TriggerPCSCFRestoration",
			http.MethodPost,
			"/restore-pcscf",
			s.HandleTriggerPCSCFRestoration,
		},

		{
			"UpdateNwdafRegistration",
			http.MethodPatch,
			"/:ueId/registrations/nwdaf-registrations/:nwdafRegistrationId",
			s.HandleUpdateNwdafRegistration,
		},

		{
			"UpdateRoamingInformation",
			http.MethodPost,
			"/:ueId/registrations/amf-3gpp-access/roaming-info-update",
			s.HandleUpdateRoamingInformation,
		},

		{
			"UpdateSmfRegistration",
			http.MethodPatch,
			"/:ueId/registrations/smf-registrations/:pduSessionId",
			s.HandleUpdateSmfRegistration,
		},
	}
}

// GetAmfNon3gppAccess - retrieve the AMF registration for non-3GPP access information
func (s *Server) HandleGetAmfNon3gppAccess(c *gin.Context) {
	logger.UecmLog.Infoln("Handle GetAmfNon3gppAccessRequest")

	ueId := c.Param("ueId")
	supportedFeatures := c.Query("supported-features")
	var queryAmfContextNon3gppRequest Nudr_DataRepository.QueryAmfContextNon3gppRequest
	queryAmfContextNon3gppRequest.SupportedFeatures = &supportedFeatures
	queryAmfContextNon3gppRequest.UeId = &ueId
	s.Processor().GetAmfNon3gppAccessProcedure(c, queryAmfContextNon3gppRequest, ueId)
}

// Register - register as AMF for non-3GPP access
func (s *Server) HandleRegistrationAmfNon3gppAccess(c *gin.Context) {
	var amfNon3GppAccessRegistration models.AmfNon3GppAccessRegistration

	requestBody, err := c.GetRawData()
	if err != nil {
		problemDetail := models.ProblemDetails{
			Title:  "System failure",
			Status: http.StatusInternalServerError,
			Detail: err.Error(),
			Cause:  "SYSTEM_FAILURE",
		}
		logger.UecmLog.Errorf("Get Request Body error: %+v", err)
		c.JSON(http.StatusInternalServerError, problemDetail)
		return
	}

	err = openapi.Deserialize(&amfNon3GppAccessRegistration, requestBody, "application/json")
	if err != nil {
		problemDetail := "[Request Body] " + err.Error()
		rsp := models.ProblemDetails{
			Title:  "Malformed request syntax",
			Status: http.StatusBadRequest,
			Detail: problemDetail,
		}
		logger.UecmLog.Errorln(problemDetail)
		c.JSON(http.StatusBadRequest, rsp)
		return
	}

	logger.UecmLog.Infof("Handle RegisterAmfNon3gppAccessRequest")

	ueID := c.Param("ueId")

	s.Processor().RegisterAmfNon3gppAccessProcedure(c, amfNon3GppAccessRegistration, ueID)
}

// RegistrationAmf3gppAccess - register as AMF for 3GPP access
func (s *Server) HandleRegistrationAmf3gppAccess(c *gin.Context) {
	var amf3GppAccessRegistration models.Amf3GppAccessRegistration
	requestBody, err := c.GetRawData()
	if err != nil {
		problemDetail := models.ProblemDetails{
			Title:  "System failure",
			Status: http.StatusInternalServerError,
			Detail: err.Error(),
			Cause:  "SYSTEM_FAILURE",
		}
		logger.UecmLog.Errorf("Get Request Body error: %+v", err)
		c.JSON(http.StatusInternalServerError, problemDetail)
		return
	}

	err = openapi.Deserialize(&amf3GppAccessRegistration, requestBody, "application/json")
	if err != nil {
		problemDetail := "[Request Body] " + err.Error()
		rsp := models.ProblemDetails{
			Title:  "Malformed request syntax",
			Status: http.StatusBadRequest,
			Detail: problemDetail,
		}
		logger.UecmLog.Errorln(problemDetail)
		c.JSON(http.StatusBadRequest, rsp)
		return
	}

	logger.UecmLog.Infof("Handle RegistrationAmf3gppAccess")

	ueID := c.Param("ueId")
	logger.UecmLog.Info("UEID: ", ueID)

	s.Processor().RegistrationAmf3gppAccessProcedure(c, amf3GppAccessRegistration, ueID)
}

// UpdateAmfNon3gppAccess - update a parameter in the AMF registration for non-3GPP access
func (s *Server) HandleUpdateAmfNon3gppAccess(c *gin.Context) {
	var amfNon3GppAccessRegistrationModification models.AmfNon3GppAccessRegistrationModification
	requestBody, err := c.GetRawData()
	if err != nil {
		problemDetail := models.ProblemDetails{
			Title:  "System failure",
			Status: http.StatusInternalServerError,
			Detail: err.Error(),
			Cause:  "SYSTEM_FAILURE",
		}
		logger.UecmLog.Errorf("Get Request Body error: %+v", err)
		c.JSON(http.StatusInternalServerError, problemDetail)
		return
	}

	err = openapi.Deserialize(&amfNon3GppAccessRegistrationModification, requestBody, "application/json")
	if err != nil {
		problemDetail := "[Request Body] " + err.Error()
		rsp := models.ProblemDetails{
			Title:  "Malformed request syntax",
			Status: http.StatusBadRequest,
			Detail: problemDetail,
		}
		logger.UecmLog.Errorln(problemDetail)
		c.JSON(http.StatusBadRequest, rsp)
		return
	}

	logger.UecmLog.Infof("Handle UpdateAmfNon3gppAccessRequest")

	ueID := c.Param("ueId")

	s.Processor().UpdateAmfNon3gppAccessProcedure(c, amfNon3GppAccessRegistrationModification, ueID)
}

// UpdateAmf3gppAccess - Update a parameter in the AMF registration for 3GPP access
func (s *Server) HandleUpdateAmf3gppAccess(c *gin.Context) {
	var amf3GppAccessRegistrationModification models.Amf3GppAccessRegistrationModification

	requestBody, err := c.GetRawData()
	if err != nil {
		problemDetail := models.ProblemDetails{
			Title:  "System failure",
			Status: http.StatusInternalServerError,
			Detail: err.Error(),
			Cause:  "SYSTEM_FAILURE",
		}
		logger.UecmLog.Errorf("Get Request Body error: %+v", err)
		c.JSON(http.StatusInternalServerError, problemDetail)
		return
	}

	err = openapi.Deserialize(&amf3GppAccessRegistrationModification, requestBody, "application/json")
	if err != nil {
		problemDetail := "[Request Body] " + err.Error()
		rsp := models.ProblemDetails{
			Title:  "Malformed request syntax",
			Status: http.StatusBadRequest,
			Detail: problemDetail,
		}
		logger.UecmLog.Errorln(problemDetail)
		c.JSON(http.StatusBadRequest, rsp)
		return
	}

	logger.UecmLog.Infof("Handle UpdateAmf3gppAccessRequest")

	ueID := c.Param("ueId")

	s.Processor().UpdateAmf3gppAccessProcedure(c, amf3GppAccessRegistrationModification, ueID)
}

// DeregistrationSmsfNon3gppAccess - delete SMSF registration for non 3GPP access
func (s *Server) HandleDeregistrationSmsfNon3gppAccess(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{})
}

// DeregistrationSmsf3gppAccess - delete the SMSF registration for 3GPP access
func (s *Server) HandleDeregistrationSmsf3gppAccess(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{})
}

// GetSmsfNon3gppAccess - retrieve the SMSF registration for non-3GPP access information
func (s *Server) HandleGetSmsfNon3gppAccess(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{})
}

// RegistrationSmsfNon3gppAccess - register as SMSF for non-3GPP access
func (s *Server) HandleRegistrationSmsfNon3gppAccess(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{})
}

// UpdateSMSFReg3GPP - register as SMSF for 3GPP access
func (s *Server) HandleUpdateSMSFReg3GPP(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{})
}

// GetSmsf3gppAccess - retrieve the SMSF registration for 3GPP access information
func (s *Server) HandleGetSmsf3gppAccess(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{})
}

// DeregistrationSmfRegistrations - delete an SMF registration
func (s *Server) HandleDeregistrationSmfRegistrations(c *gin.Context) {
	logger.UecmLog.Infof("Handle DeregistrationSmfRegistrations")

	ueID := c.Params.ByName("ueId")
	pduSessionID := c.Params.ByName("pduSessionId")

	s.Processor().DeregistrationSmfRegistrationsProcedure(c, ueID, pduSessionID)
}

// RegistrationSmfRegistrations - register as SMF
func (s *Server) HandleRegistrationSmfRegistrations(c *gin.Context) {
	var smfRegistration models.SmfRegistration

	requestBody, err := c.GetRawData()
	if err != nil {
		problemDetail := models.ProblemDetails{
			Title:  "System failure",
			Status: http.StatusInternalServerError,
			Detail: err.Error(),
			Cause:  "SYSTEM_FAILURE",
		}
		logger.UecmLog.Errorf("Get Request Body error: %+v", err)
		c.JSON(http.StatusInternalServerError, problemDetail)
		return
	}

	err = openapi.Deserialize(&smfRegistration, requestBody, "application/json")
	if err != nil {
		problemDetail := "[Request Body] " + err.Error()
		rsp := models.ProblemDetails{
			Title:  "Malformed request syntax",
			Status: http.StatusBadRequest,
			Detail: problemDetail,
		}
		logger.UecmLog.Errorln(problemDetail)
		c.JSON(http.StatusBadRequest, rsp)
		return
	}

	logger.UecmLog.Infof("Handle RegistrationSmfRegistrations")

	ueID := c.Params.ByName("ueId")
	pduSessionID := c.Params.ByName("pduSessionId")

	s.Processor().RegistrationSmfRegistrationsProcedure(
		c,
		&smfRegistration,
		ueID,
		pduSessionID,
	)
}

// GetAmf3gppAccess - retrieve the AMF registration for 3GPP access information
func (s *Server) HandleGetAmf3gppAccess(c *gin.Context) {
	logger.UecmLog.Infof("Handle HandleGetAmf3gppAccessRequest")

	ueID := c.Param("ueId")
	supportedFeatures := c.Query("supported-features")

	s.Processor().GetAmf3gppAccessProcedure(c, ueID, supportedFeatures)
}

func (s *Server) HandleDeregAMF(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{})
}

func (s *Server) HandleGetIpSmGwRegistration(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{})
}

func (s *Server) HandleGetLocationInfo(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{})
}

func (s *Server) HandleGetNwdafRegistration(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{})
}

func (s *Server) HandleGetRegistrations(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{})
}

func (s *Server) HandleGetSmfRegistration(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{})
}

func (s *Server) HandleIpSmGwDeregistration(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{})
}

func (s *Server) HandleIpSmGwRegistration(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{})
}

func (s *Server) HandleNwdafDeregistration(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{})
}

func (s *Server) HandleNwdafRegistration(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{})
}

func (s *Server) HandlePeiUpdate(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{})
}

func (s *Server) HandleRetrieveSmfRegistration(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{})
}

func (s *Server) HandleSendRoutingInfoSm(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{})
}

func (s *Server) HandleTriggerPCSCFRestoration(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{})
}

func (s *Server) HandleUpdateNwdafRegistration(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{})
}

func (s *Server) HandleUpdateRoamingInformation(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{})
}

func (s *Server) HandleUpdateSmfRegistration(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{})
}
