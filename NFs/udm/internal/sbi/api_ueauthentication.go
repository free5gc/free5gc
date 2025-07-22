package sbi

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/free5gc/openapi"
	"github.com/free5gc/openapi/models"
	"github.com/free5gc/udm/internal/logger"
)

func (s *Server) getUEAuthenticationRoutes() []Route {
	return []Route{
		{
			"Index",
			http.MethodGet,
			"/",
			s.HandleIndex,
		},
	}
}

// ConfirmAuth - Create a new confirmation event
func (s *Server) HandleConfirmAuth(c *gin.Context) {
	var authEvent models.AuthEvent
	requestBody, err := c.GetRawData()
	if err != nil {
		problemDetail := models.ProblemDetails{
			Title:  "System failure",
			Status: http.StatusInternalServerError,
			Detail: err.Error(),
			Cause:  "SYSTEM_FAILURE",
		}
		logger.UeauLog.Errorf("Get Request Body error: %+v", err)
		c.JSON(http.StatusInternalServerError, problemDetail)
		return
	}

	err = openapi.Deserialize(&authEvent, requestBody, "application/json")
	if err != nil {
		problemDetail := "[Request Body] " + err.Error()
		rsp := models.ProblemDetails{
			Title:  "Malformed request syntax",
			Status: http.StatusBadRequest,
			Detail: problemDetail,
		}
		logger.UeauLog.Errorln(problemDetail)
		c.JSON(http.StatusBadRequest, rsp)
		return
	}

	supi := c.Params.ByName("supi")

	logger.UeauLog.Infoln("Handle ConfirmAuthDataRequest")

	s.Processor().ConfirmAuthDataProcedure(c, authEvent, supi)
}

// GenerateAuthData - Generate authentication data for the UE
func (s *Server) HandleGenerateAuthData(c *gin.Context) {
	var authInfoReq models.AuthenticationInfoRequest

	requestBody, err := c.GetRawData()
	if err != nil {
		problemDetail := models.ProblemDetails{
			Title:  "System failure",
			Status: http.StatusInternalServerError,
			Detail: err.Error(),
			Cause:  "SYSTEM_FAILURE",
		}
		logger.UeauLog.Errorf("Get Request Body error: %+v", err)
		c.JSON(http.StatusInternalServerError, problemDetail)
		return
	}

	err = openapi.Deserialize(&authInfoReq, requestBody, "application/json")
	if err != nil {
		problemDetail := "[Request Body] " + err.Error()
		rsp := models.ProblemDetails{
			Title:  "Malformed request syntax",
			Status: http.StatusBadRequest,
			Detail: problemDetail,
		}
		logger.UeauLog.Errorln(problemDetail)
		c.JSON(http.StatusBadRequest, rsp)
		return
	}

	logger.UeauLog.Infoln("Handle GenerateAuthDataRequest")

	supiOrSuci := c.Param("supiOrSuci")

	s.Processor().GenerateAuthDataProcedure(c, authInfoReq, supiOrSuci)
}

func (s *Server) HandleDeleteAuth(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{})
}

func (s *Server) HandleGenerateAv(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{})
}

func (s *Server) HandleGenerateGbaAv(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{})
}

func (s *Server) HandleGenerateProseAV(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{})
}

func (s *Server) HandleGetRgAuthData(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{})
}

func (s *Server) UEAUTwoLayerPathHandlerFunc(c *gin.Context) {
	twoLayer := c.Param("twoLayer")

	// for "/:supi/auth-events"
	if twoLayer == "auth-events" && http.MethodPost == c.Request.Method {
		s.HandleConfirmAuth(c)
		return
	}

	// for "/:supiOrSuci/security-information-rg"
	if twoLayer == "security-information-rg" && http.MethodGet == c.Request.Method {
		var tmpParams gin.Params
		tmpParams = append(tmpParams, gin.Param{Key: "supiOrSuci", Value: c.Param("supi")})
		c.Params = tmpParams
		s.HandleGetRgAuthData(c)
		return
	}

	c.String(http.StatusNotFound, "404 page not found")
}

func (s *Server) UEAUThreeLayerPathHandlerFunc(c *gin.Context) {
	twoLayer := c.Param("twoLayer")

	// for "/:supi/auth-events/:authEventId"
	if twoLayer == "auth-events" && http.MethodPut == c.Request.Method {
		s.HandleDeleteAuth(c)
		return
	}

	// for "/:supi/gba-security-information/generate-av"
	if twoLayer == "gba-security-information" && http.MethodPost == c.Request.Method {
		s.HandleGenerateGbaAv(c)
		return
	}

	// for "/:supiOrSuci/prose-security-information/generate-av"
	if twoLayer == "prose-security-information" && http.MethodPost == c.Request.Method {
		var tmpParams gin.Params
		tmpParams = append(tmpParams, gin.Param{Key: "supiOrSuci", Value: c.Param("supi")})
		c.Params = tmpParams
		s.HandleGenerateProseAV(c)
		return
	}

	// for "/:supiOrSuci/security-information/generate-auth-data"
	if twoLayer == "security-information" && http.MethodPost == c.Request.Method {
		var tmpParams gin.Params
		tmpParams = append(tmpParams, gin.Param{Key: "supiOrSuci", Value: c.Param("supi")})
		c.Params = tmpParams
		s.HandleGenerateAuthData(c)
		return
	}

	c.String(http.StatusNotFound, "404 page not found")
}
