package sbi

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) getMTRoutes() []Route {
	return []Route{
		{
			"Index",
			http.MethodGet,
			"/",
			s.HandleIndex,
		},

		{
			"ProvideLocationInfo",
			http.MethodPost,
			"/:supi/loc-info/provide-loc-info",
			s.HandleProvideLocationInfo,
		},

		{
			"QueryUeInfo",
			http.MethodGet,
			"/:supi",
			s.HandleQueryUeInfo,
		},
	}
}

func (s *Server) HandleProvideLocationInfo(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{})
}

func (s *Server) HandleQueryUeInfo(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{})
}
