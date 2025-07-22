package sbi

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) getSorprotectionRoutes() []Route {
	return []Route{
		{
			Name:    "Index",
			Method:  http.MethodGet,
			Pattern: "/",
			APIFunc: func(c *gin.Context) {
				c.String(http.StatusOK, "Hello free5GC!")
			},
		},
		{
			Name:    "SupiUeSorPost",
			Method:  http.MethodPost,
			Pattern: "/:supi/ue-sor/generate-sor-data",
			APIFunc: s.HTTPSupiUeSorPost,
		},
	}
}

func (s *Server) HTTPSupiUeSorPost(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{})
}
