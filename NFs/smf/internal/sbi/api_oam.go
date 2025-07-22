package sbi

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) getOAMRoutes() []Route {
	return []Route{
		{
			Name:    "Index",
			Method:  http.MethodGet,
			Pattern: "/",
			APIFunc: func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"status": "Service Available"})
			},
		},
		{
			Name:    "Get UE PDU Session Info",
			Method:  http.MethodGet,
			Pattern: "/ue-pdu-session-info/:smContextRef",
			APIFunc: s.HTTPGetUEPDUSessionInfo,
		},
		{
			Name:    "Get SMF Userplane Information",
			Method:  http.MethodGet,
			Pattern: "/user-plane-info/",
			APIFunc: s.HTTPGetSMFUserPlaneInfo,
		},
	}
}

func (s *Server) HTTPGetUEPDUSessionInfo(c *gin.Context) {
	smContextRef := c.Params.ByName("smContextRef")

	s.Processor().HandleOAMGetUEPDUSessionInfo(c, smContextRef)
}

func (s *Server) HTTPGetSMFUserPlaneInfo(c *gin.Context) {
	s.Processor().HandleGetSMFUserPlaneInfo(c)
}
