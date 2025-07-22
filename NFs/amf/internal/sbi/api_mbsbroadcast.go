package sbi

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) getMbsBroadcastRoutes() []Route {
	return []Route{
		{
			Method:  http.MethodGet,
			Pattern: "/",
			APIFunc: func(c *gin.Context) {
				c.String(http.StatusOK, "Hello World!")
			},
		},
		{
			Name:    "ContextCreate",
			Method:  http.MethodPost,
			Pattern: "/mbs-contexts",
			APIFunc: s.HTTPContextCreate,
		},
		{
			Name:    "ContextUpdate",
			Method:  http.MethodPost,
			Pattern: "/mbs-contexts/:mbsContextRef/update",
			APIFunc: s.HTTPContextUpdate,
		},
		{
			Name:    "ContextReleas",
			Method:  http.MethodDelete,
			Pattern: "/mbs-contexts/:mbsContextRef",
			APIFunc: s.HTTPContextRelease,
		},
	}
}

func (s *Server) HTTPContextCreate(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{})
}

func (s *Server) HTTPContextUpdate(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{})
}

func (s *Server) HTTPContextRelease(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{})
}
