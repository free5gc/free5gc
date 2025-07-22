package sbi

import (
	"net/http"

	"github.com/free5gc/openapi/models_nef"
	"github.com/gin-gonic/gin"
)

func (s *Server) getTrafficInfluenceEndpoints() []Endpoint {
	return []Endpoint{
		{
			Method:  http.MethodGet,
			Pattern: "/:afID/subscriptions",
			APIFunc: s.apiGetTrafficInfluenceSubscription,
		},
		{
			Method:  http.MethodPost,
			Pattern: "/:afID/subscriptions",
			APIFunc: s.apiPostTrafficInfluenceSubscription,
		},
		{
			Method:  http.MethodGet,
			Pattern: "/:afID/subscriptions/:subID",
			APIFunc: s.apiGetIndividualTrafficInfluenceSubscription,
		},
		{
			Method:  http.MethodPut,
			Pattern: "/:afID/subscriptions/:subID",
			APIFunc: s.apiPutIndividualTrafficInfluenceSubscription,
		},
		{
			Method:  http.MethodPatch,
			Pattern: "/:afID/subscriptions/:subID",
			APIFunc: s.apiPatchIndividualTrafficInfluenceSubscription,
		},
		{
			Method:  http.MethodDelete,
			Pattern: "/:afID/subscriptions/:subID",
			APIFunc: s.apiDeleteIndividualTrafficInfluenceSubscription,
		},
	}
}

func (s *Server) apiGetTrafficInfluenceSubscription(gc *gin.Context) {
	hdlRsp := s.Processor().GetTrafficInfluenceSubscription(
		gc.Param("afID"))

	s.buildAndSendHttpResponse(gc, hdlRsp, false)
}

func (s *Server) apiPostTrafficInfluenceSubscription(gc *gin.Context) {
	contentType, err := checkContentTypeIsJSON(gc)
	if err != nil {
		return
	}

	var tiSub models_nef.TrafficInfluSub
	if err := s.deserializeData(gc, &tiSub, contentType); err != nil {
		return
	}

	hdlRsp := s.Processor().PostTrafficInfluenceSubscription(
		gc.Param("afID"), &tiSub)

	s.buildAndSendHttpResponse(gc, hdlRsp, false)
}

func (s *Server) apiGetIndividualTrafficInfluenceSubscription(gc *gin.Context) {
	hdlRsp := s.Processor().GetIndividualTrafficInfluenceSubscription(
		gc.Param("afID"), gc.Param("subID"))

	s.buildAndSendHttpResponse(gc, hdlRsp, false)
}

func (s *Server) apiPutIndividualTrafficInfluenceSubscription(gc *gin.Context) {
	contentType, err := checkContentTypeIsJSON(gc)
	if err != nil {
		return
	}

	var tiSub models_nef.TrafficInfluSub
	if err := s.deserializeData(gc, &tiSub, contentType); err != nil {
		return
	}

	hdlRsp := s.Processor().PutIndividualTrafficInfluenceSubscription(
		gc.Param("afID"), gc.Param("subID"), &tiSub)

	s.buildAndSendHttpResponse(gc, hdlRsp, false)
}

func (s *Server) apiPatchIndividualTrafficInfluenceSubscription(gc *gin.Context) {
	contentType, err := checkContentTypeIsJSON(gc)
	if err != nil {
		return
	}

	var tiSubPatch models_nef.TrafficInfluSubPatch
	if err := s.deserializeData(gc, &tiSubPatch, contentType); err != nil {
		return
	}

	hdlRsp := s.Processor().PatchIndividualTrafficInfluenceSubscription(
		gc.Param("afID"), gc.Param("subID"), &tiSubPatch)

	s.buildAndSendHttpResponse(gc, hdlRsp, false)
}

func (s *Server) apiDeleteIndividualTrafficInfluenceSubscription(gc *gin.Context) {
	hdlRsp := s.Processor().DeleteIndividualTrafficInfluenceSubscription(
		gc.Param("afID"), gc.Param("subID"))

	s.buildAndSendHttpResponse(gc, hdlRsp, false)
}
