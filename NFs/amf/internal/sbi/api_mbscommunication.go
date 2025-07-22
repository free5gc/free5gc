package sbi

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) getMbsCommunicationRoutes() []Route {
	return []Route{
		{
			Method:  http.MethodGet,
			Pattern: "/",
			APIFunc: func(c *gin.Context) {
				c.String(http.StatusOK, "Hello World!")
			},
		},
		{
			Name:    "N2MessageTransfer",
			Method:  http.MethodPost,
			Pattern: "/n2-messages/transfer",
			APIFunc: s.HTTPN2MessageTransfer,
		},
	}
}

func (s *Server) HTTPN2MessageTransfer(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{})
}
