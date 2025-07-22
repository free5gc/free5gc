package sbi

import (
	"net/http"

	"github.com/free5gc/openapi/models"
	"github.com/gin-gonic/gin"
)

func (s *Server) getPFDFEndpoints() []Endpoint {
	return []Endpoint{
		{
			Method:  http.MethodGet,
			Pattern: "/applications",
			APIFunc: s.apiGetApplicationsPFD,
		},
		{
			Method:  http.MethodGet,
			Pattern: "/applications/:appID",
			APIFunc: s.apiGetIndividualApplicationPFD,
		},
		{
			Method:  http.MethodPost,
			Pattern: "/subscriptions",
			APIFunc: s.apiPostPFDSubscriptions,
		},
		{
			Method:  http.MethodDelete,
			Pattern: "/subscriptions/:subID",
			APIFunc: s.apiDeleteIndividualPFDSubscription,
		},
	}
}

func (s *Server) apiGetApplicationsPFD(gc *gin.Context) {
	// TODO: support URI query parameters: supported-features
	hdlRsp := s.Processor().GetApplicationsPFD(gc.QueryArray("application-ids"))

	s.buildAndSendHttpResponse(gc, hdlRsp, false)
}

func (s *Server) apiGetIndividualApplicationPFD(gc *gin.Context) {
	// TODO: support URI query parameters: supported-features
	hdlRsp := s.Processor().GetIndividualApplicationPFD(gc.Param("appID"))

	s.buildAndSendHttpResponse(gc, hdlRsp, false)
}

func (s *Server) apiPostPFDSubscriptions(gc *gin.Context) {
	contentType, err := checkContentTypeIsJSON(gc)
	if err != nil {
		return
	}

	var pfdSubsc models.PfdSubscription
	if err := s.deserializeData(gc, &pfdSubsc, contentType); err != nil {
		return
	}

	hdlRsp := s.Processor().PostPFDSubscriptions(&pfdSubsc)

	s.buildAndSendHttpResponse(gc, hdlRsp, false)
}

func (s *Server) apiDeleteIndividualPFDSubscription(gc *gin.Context) {
	hdlRsp := s.Processor().DeleteIndividualPFDSubscription(gc.Param("subID"))

	s.buildAndSendHttpResponse(gc, hdlRsp, false)
}
