package sbi

import (
	"net/http"

	"github.com/free5gc/openapi/models"
	"github.com/gin-gonic/gin"
)

func (s *Server) getPFDManagementEndpoints() []Endpoint {
	return []Endpoint{
		{
			Method:  http.MethodGet,
			Pattern: "/:scsAsID/transactions",
			APIFunc: s.apiGetPFDManagementTransactions,
		},
		{
			Method:  http.MethodPost,
			Pattern: "/:scsAsID/transactions",
			APIFunc: s.apiPostPFDManagementTransactions,
		},
		{
			Method:  http.MethodDelete,
			Pattern: "/:scsAsID/transactions",
			APIFunc: s.apiDeletePFDManagementTransactions,
		},
		{
			Method:  http.MethodGet,
			Pattern: "/:scsAsID/transactions/:transID",
			APIFunc: s.apiGetIndividualPFDManagementTransaction,
		},
		{
			Method:  http.MethodPut,
			Pattern: "/:scsAsID/transactions/:transID",
			APIFunc: s.apiPutIndividualPFDManagementTransaction,
		},
		{
			Method:  http.MethodDelete,
			Pattern: "/:scsAsID/transactions/:transID",
			APIFunc: s.apiDeleteIndividualPFDManagementTransaction,
		},
		{
			Method:  http.MethodGet,
			Pattern: "/:scsAsID/transactions/:transID/applications/:appID",
			APIFunc: s.apiGetIndividualApplicationPFDManagement,
		},
		{
			Method:  http.MethodDelete,
			Pattern: "/:scsAsID/transactions/:transID/applications/:appID",
			APIFunc: s.apiDeleteIndividualApplicationPFDManagement,
		},
		{
			Method:  http.MethodPut,
			Pattern: "/:scsAsID/transactions/:transID/applications/:appID",
			APIFunc: s.apiPutIndividualApplicationPFDManagement,
		},
		{
			Method:  http.MethodPatch,
			Pattern: "/:scsAsID/transactions/:transID/applications/:appID",
			APIFunc: s.apiPatchIndividualApplicationPFDManagement,
		},
	}
}

func (s *Server) apiGetPFDManagementTransactions(gc *gin.Context) {
	hdlRsp := s.Processor().GetPFDManagementTransactions(
		gc.Param("scsAsID"))

	s.buildAndSendHttpResponse(gc, hdlRsp, false)
}

func (s *Server) apiPostPFDManagementTransactions(gc *gin.Context) {
	contentType, err := checkContentTypeIsJSON(gc)
	if err != nil {
		return
	}

	var pfdMng models.PfdManagement
	if err := s.deserializeData(gc, &pfdMng, contentType); err != nil {
		return
	}

	hdlRsp := s.Processor().PostPFDManagementTransactions(
		gc.Param("scsAsID"), &pfdMng)

	s.buildAndSendHttpResponse(gc, hdlRsp, false)
}

func (s *Server) apiDeletePFDManagementTransactions(gc *gin.Context) {
	hdlRsp := s.Processor().DeletePFDManagementTransactions(
		gc.Param("scsAsID"))

	s.buildAndSendHttpResponse(gc, hdlRsp, false)
}

func (s *Server) apiGetIndividualPFDManagementTransaction(gc *gin.Context) {
	hdlRsp := s.Processor().GetIndividualPFDManagementTransaction(
		gc.Param("scsAsID"), gc.Param("transID"))

	s.buildAndSendHttpResponse(gc, hdlRsp, false)
}

func (s *Server) apiPutIndividualPFDManagementTransaction(gc *gin.Context) {
	contentType, err := checkContentTypeIsJSON(gc)
	if err != nil {
		return
	}

	var pfdMng models.PfdManagement
	if err := s.deserializeData(gc, &pfdMng, contentType); err != nil {
		return
	}

	hdlRsp := s.Processor().PutIndividualPFDManagementTransaction(
		gc.Param("scsAsID"), gc.Param("transID"), &pfdMng)

	s.buildAndSendHttpResponse(gc, hdlRsp, false)
}

func (s *Server) apiDeleteIndividualPFDManagementTransaction(gc *gin.Context) {
	hdlRsp := s.Processor().DeleteIndividualPFDManagementTransaction(
		gc.Param("scsAsID"), gc.Param("transID"))

	s.buildAndSendHttpResponse(gc, hdlRsp, false)
}

func (s *Server) apiGetIndividualApplicationPFDManagement(gc *gin.Context) {
	hdlRsp := s.Processor().GetIndividualApplicationPFDManagement(
		gc.Param("scsAsID"), gc.Param("transID"), gc.Param("appID"))

	s.buildAndSendHttpResponse(gc, hdlRsp, false)
}

func (s *Server) apiDeleteIndividualApplicationPFDManagement(gc *gin.Context) {
	hdlRsp := s.Processor().DeleteIndividualApplicationPFDManagement(
		gc.Param("scsAsID"), gc.Param("transID"), gc.Param("appID"))

	s.buildAndSendHttpResponse(gc, hdlRsp, false)
}

func (s *Server) apiPutIndividualApplicationPFDManagement(gc *gin.Context) {
	contentType, err := checkContentTypeIsJSON(gc)
	if err != nil {
		return
	}

	var pfdData models.PfdData
	if err := s.deserializeData(gc, &pfdData, contentType); err != nil {
		return
	}

	hdlRsp := s.Processor().PutIndividualApplicationPFDManagement(
		gc.Param("scsAsID"), gc.Param("transID"), gc.Param("appID"), &pfdData)

	s.buildAndSendHttpResponse(gc, hdlRsp, false)
}

func (s *Server) apiPatchIndividualApplicationPFDManagement(gc *gin.Context) {
	contentType, err := checkContentTypeIsJSON(gc)
	if err != nil {
		return
	}

	var pfdData models.PfdData
	if err := s.deserializeData(gc, &pfdData, contentType); err != nil {
		return
	}

	hdlRsp := s.Processor().PatchIndividualApplicationPFDManagement(
		gc.Param("scsAsID"), gc.Param("transID"), gc.Param("appID"), &pfdData)

	s.buildAndSendHttpResponse(gc, hdlRsp, false)
}
