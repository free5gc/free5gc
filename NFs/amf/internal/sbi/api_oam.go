package sbi

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) getOAMRoutes() []Route {
	return []Route{
		{
			Method:  http.MethodGet,
			Pattern: "/",
			APIFunc: func(c *gin.Context) {
				c.String(http.StatusOK, "Hello World!")
			},
		},
		{
			Name:    "RegisteredUEContext",
			Method:  http.MethodGet,
			Pattern: "/registered-ue-context",
			APIFunc: s.HTTPRegisteredUEContext,
		},
		{
			Name:    "RegisteredUEContext",
			Method:  http.MethodGet,
			Pattern: "/registered-ue-context/:supi",
			APIFunc: s.HTTPRegisteredUEContext,
		},
	}
}

func (s *Server) setCorsHeader(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
	c.Writer.Header().Set("Access-Control-Allow-Headers",
		"Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, "+
			"Authorization, accept, origin, Cache-Control, X-Requested-With")
	c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, PATCH, DELETE")
}

func (s *Server) HTTPRegisteredUEContext(c *gin.Context) {
	s.setCorsHeader(c)
	s.Processor().HandleOAMRegisteredUEContext(c)
}
