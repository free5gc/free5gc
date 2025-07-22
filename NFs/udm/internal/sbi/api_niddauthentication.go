package sbi

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) getNIDDAuthenticationRoutes() []Route {
	return []Route{
		{
			"Index",
			http.MethodGet,
			"/",
			s.HandleIndex,
		},

		{
			"AuthorizeNiddData",
			http.MethodPost,
			"/:ueIdentity/authorize",
			s.HandleAuthorizeNiddData,
		},
	}
}

func (s *Server) HandleAuthorizeNiddData(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{})
}
