package sbi

import (
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"github.com/free5gc/openapi/models"
	"github.com/free5gc/pcf/internal/util"
)

const (
	CorsConfigMaxAge = 86400
)

func (s *Server) setCorsHeader(c *gin.Context) {
	// TODO: 1. turn these values into configurable variables
	// TODO: 2. use the official cors middleware
	s.router.Use(cors.New(cors.Config{
		AllowMethods: []string{"GET", "POST", "OPTIONS", "PUT", "PATCH", "DELETE"},
		AllowHeaders: []string{
			"Origin", "Content-Length", "Content-Type", "User-Agent",
			"Referrer", "Host", "Token", "X-Requested-With",
		},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowAllOrigins:  true,
		MaxAge:           CorsConfigMaxAge,
	}))

	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
	c.Writer.Header().Set(
		"Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, "+
			"X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
	c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, PATCH, DELETE")
}

func (s *Server) HTTPOAMGetAmPolicy(c *gin.Context) {
	s.setCorsHeader(c)

	supi := c.Params.ByName("supi")
	if supi == "" {
		problemDetails := &models.ProblemDetails{
			Title:  util.ERROR_INITIAL_PARAMETERS,
			Status: http.StatusBadRequest,
		}
		c.JSON(http.StatusBadRequest, problemDetails)
		return
	}
	s.Processor().HandleOAMGetAmPolicyRequest(c, supi)
}

func (s *Server) getOamRoutes() []Route {
	return []Route{
		{
			Method:  http.MethodGet,
			Pattern: "/am-policy/:supi",
			APIFunc: s.HTTPOAMGetAmPolicy,
		},
	}
}
