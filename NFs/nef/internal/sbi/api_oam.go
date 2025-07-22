package sbi

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) getOamEndpoints() []Endpoint {
	return []Endpoint{
		{
			Method:  http.MethodGet,
			Pattern: "/",
			APIFunc: s.apiGetOamIndex,
		},
	}
}

func (s *Server) apiGetOamIndex(gc *gin.Context) {
	hdlRsp := s.Processor().GetOamIndex()
	s.buildAndSendHttpResponse(gc, hdlRsp, false)
}
