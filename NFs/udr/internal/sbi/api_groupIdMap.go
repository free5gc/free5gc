package sbi

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) getGroupIdMap() []Route {
	return []Route{
		{
			Name:        "Index",
			Method:      "GET",
			Pattern:     "/",
			HandlerFunc: Index,
		},
		{
			Name:        "GetNfGroupIDs",
			Method:      "GET",
			Pattern:     "/nf-group-ids",
			HandlerFunc: s.HTTPGetNfGroupIDs,
		},
	}
}

// GetNfGroupIDs - Retrieves NF-Group IDs for provided Subscriber and NF types
func (s *Server) HTTPGetNfGroupIDs(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{})
}
