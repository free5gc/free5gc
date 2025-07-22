package consumer

import (
	"strconv"
	"strings"
	"sync"

	"github.com/free5gc/openapi"
	"github.com/free5gc/openapi/models"
	"github.com/free5gc/openapi/udr/DataRepository"
	pcf_context "github.com/free5gc/pcf/internal/context"
	"github.com/free5gc/pcf/internal/logger"
	"github.com/free5gc/pcf/internal/util"
)

type nudrService struct {
	consumer *Consumer

	nfDataSubMu sync.RWMutex

	nfDataSubClients map[string]*DataRepository.APIClient
}

func (s *nudrService) getDataSubscription(uri string) *DataRepository.APIClient {
	if uri == "" {
		return nil
	}
	s.nfDataSubMu.RLock()
	client, ok := s.nfDataSubClients[uri]
	if ok {
		defer s.nfDataSubMu.RUnlock()
		return client
	}

	configuration := DataRepository.NewConfiguration()
	configuration.SetBasePath(uri)
	client = DataRepository.NewAPIClient(configuration)

	s.nfDataSubMu.RUnlock()
	s.nfDataSubMu.Lock()
	defer s.nfDataSubMu.Unlock()
	s.nfDataSubClients[uri] = client
	return client
}

func (s *nudrService) GetSessionManagementPolicyData(uri string,
	ueId string, sliceInfo *models.Snssai, dnn string) (
	resp *DataRepository.ReadSessionManagementPolicyDataResponse,
	problemDetails *models.ProblemDetails, err error,
) {
	if uri == "" {
		problemDetail := util.GetProblemDetail("Can't find any UDR which supported to this PCF",
			"GetSessionManagementPolicyData Can't find UDR URI")
		return nil, &problemDetail, nil
	}

	if ueId == "" {
		problemDetail := util.GetProblemDetail("Can't find any UDR which supported to this PCF",
			"GetSessionManagementPolicyData Can't find UE ID")
		return nil, &problemDetail, nil
	}

	if sliceInfo == nil {
		problemDetails := util.GetProblemDetail("Can't find any UDR which supported to this PCF",
			"GetSessionManagementPolicyData Can't find Slice Info")
		return nil, &problemDetails, nil
	}

	if dnn == "" {
		problemDetails := util.GetProblemDetail("Can't find any UDR which supported to this PCF",
			"GetSessionManagementPolicyData Can't find DNN")
		return nil, &problemDetails, nil
	}

	ctx, pd, err := s.consumer.Context().GetTokenCtx(models.ServiceName_NUDR_DR, models.NrfNfManagementNfType_UDR)
	if err != nil {
		return nil, nil, err
	} else if pd != nil {
		return nil, pd, nil
	}

	client := s.getDataSubscription(uri)
	param := DataRepository.ReadSessionManagementPolicyDataRequest{
		UeId:   &ueId,
		Snssai: sliceInfo,
		Dnn:    &dnn,
	}
	resp, localErr := client.SessionManagementPolicyDataDocumentApi.ReadSessionManagementPolicyData(ctx, &param)
	if localErr == nil {
		return resp, nil, nil
	}
	if genericErr, ok := localErr.(openapi.GenericOpenAPIError); ok {
		if problemDetails, ok := genericErr.Model().(models.ProblemDetails); ok {
			return nil, &problemDetails, nil
		}

		logger.ConsumerLog.Errorf("openapi error: %+v", localErr)
		return nil, nil, localErr
	}

	return nil, nil, localErr
}

func (s *nudrService) CreateBdtData(uri string, bdtData *models.BdtData) (
	problemDetails *models.ProblemDetails, err error,
) {
	if uri == "" {
		problemDetail := util.GetProblemDetail("Can't find any UDR which supported to this PCF",
			"CreateBdtData Can't find UDR URI")
		return &problemDetail, nil
	}

	ctx, pd, err := s.consumer.Context().GetTokenCtx(models.ServiceName_NUDR_DR, models.NrfNfManagementNfType_UDR)
	if err != nil {
		return nil, err
	} else if pd != nil {
		return pd, nil
	}

	client := s.getDataSubscription(uri)
	param := DataRepository.CreateIndividualBdtDataRequest{
		BdtData: bdtData,
	}
	_, localErr := client.IndividualBdtDataDocumentApi.CreateIndividualBdtData(ctx, &param)
	if localErr == nil {
		return nil, nil
	}
	if genericErr, ok := localErr.(openapi.GenericOpenAPIError); ok {
		if problemDetails, ok := genericErr.Model().(models.ProblemDetails); ok {
			return &problemDetails, nil
		}

		logger.ConsumerLog.Errorf("openapi error: %+v", localErr)
		return nil, localErr
	}

	return nil, localErr
}

func (s *nudrService) CreateBdtPolicyContext(uri string, req *DataRepository.ReadBdtDataRequest) (
	resp *DataRepository.ReadBdtDataResponse, problemDetails *models.ProblemDetails, err error,
) {
	if uri == "" {
		problemDetail := util.GetProblemDetail("Can't find any UDR which supported to this PCF",
			"CreateBdtData Can't find UDR URI")
		return nil, &problemDetail, nil
	}
	ctx, pd, err := s.consumer.Context().GetTokenCtx(models.ServiceName_NUDR_DR, models.NrfNfManagementNfType_UDR)
	if err != nil {
		return nil, nil, err
	} else if pd != nil {
		return nil, pd, nil
	}

	client := s.getDataSubscription(uri)
	resp, localErr := client.BdtDataStoreApi.ReadBdtData(ctx, req)
	if localErr == nil {
		return resp, nil, nil
	}
	if genericErr, ok := localErr.(openapi.GenericOpenAPIError); ok {
		if problemDetails, ok := genericErr.Model().(models.ProblemDetails); ok {
			return nil, &problemDetails, nil
		}

		logger.ConsumerLog.Errorf("openapi error: %+v", localErr)
		return nil, nil, localErr
	}

	return nil, nil, localErr
}

func (s *nudrService) GetBdtData(uri string, bdtRefId string) (
	resp *DataRepository.ReadIndividualBdtDataResponse, problemDetails *models.ProblemDetails, err error,
) {
	if uri == "" {
		problemDetail := util.GetProblemDetail("Can't find any UDR which supported to this PCF",
			"GetBdtData Can't find UDR URI")
		return nil, &problemDetail, nil
	}

	if bdtRefId == "" {
		problemDetail := util.GetProblemDetail("Can't find any UDR which supported to this PCF",
			"GetBdtData Can't find BdtRefId")
		return nil, &problemDetail, nil
	}

	ctx, pd, err := s.consumer.Context().GetTokenCtx(models.ServiceName_NUDR_DR, models.NrfNfManagementNfType_UDR)
	if err != nil {
		return nil, nil, err
	} else if pd != nil {
		return nil, pd, nil
	}

	readBdtDataReq := DataRepository.ReadIndividualBdtDataRequest{
		BdtReferenceId: &bdtRefId,
	}

	client := s.getDataSubscription(uri)
	resp, localErr := client.IndividualBdtDataDocumentApi.ReadIndividualBdtData(ctx, &readBdtDataReq)
	if localErr == nil {
		return resp, nil, nil
	}

	if genericErr, ok := localErr.(openapi.GenericOpenAPIError); ok {
		if problemDetails, ok := genericErr.Model().(models.ProblemDetails); ok {
			return nil, &problemDetails, nil
		}

		logger.ConsumerLog.Errorf("openapi error: %+v", localErr)
		return nil, nil, localErr
	}

	return nil, nil, localErr
}

func (s *nudrService) GetAccessAndMobilityPolicyData(ue *pcf_context.UeContext) (
	amPolicyData *models.AmPolicyData,
	problemDetails *models.ProblemDetails, err error,
) {
	if ue.Supi == "" {
		problemDetail := util.GetProblemDetail("Can't find corresponding SUPI with UE", util.USER_UNKNOWN)
		logger.ConsumerLog.Warn("Can't find corresponding SUPI with UE")
		return nil, &problemDetail, nil
	}

	if ue.UdrUri == "" {
		problemDetail := util.GetProblemDetail("Can't find corresponding UDR with UE", util.USER_UNKNOWN)
		logger.ConsumerLog.Warnf("Can't find corresponding UDR with UE[%s]", ue.Supi)
		return nil, &problemDetail, nil
	}

	ctx, pd, err := s.consumer.Context().GetTokenCtx(models.ServiceName_NUDR_DR, models.NrfNfManagementNfType_UDR)
	if err != nil {
		return nil, nil, err
	} else if pd != nil {
		return nil, pd, nil
	}

	client := s.getDataSubscription(ue.UdrUri)
	param := DataRepository.ReadAccessAndMobilityPolicyDataRequest{
		UeId: &ue.Supi,
	}
	resp, localErr := client.AccessAndMobilityPolicyDataDocumentApi.ReadAccessAndMobilityPolicyData(ctx, &param)
	if localErr == nil {
		return &resp.AmPolicyData, nil, nil
	}
	if genericErr, ok := localErr.(openapi.GenericOpenAPIError); ok {
		if problemDetails, ok := genericErr.Model().(models.ProblemDetails); ok {
			return nil, &problemDetails, nil
		}

		logger.ConsumerLog.Errorf("openapi error: %+v", localErr)
		return nil, nil, localErr
	}

	return nil, nil, localErr
}

func (s *nudrService) CreateInfluenceDataSubscription(ue *pcf_context.UeContext, request models.SmPolicyContextData) (
	subscriptionID string, problemDetails *models.ProblemDetails, err error,
) {
	if ue.UdrUri == "" {
		problemDetail := util.GetProblemDetail("Can't find corresponding UDR with UE", util.USER_UNKNOWN)
		logger.ConsumerLog.Warnf("Can't find corresponding UDR with UE[%s]", ue.Supi)
		return "", &problemDetail, nil
	}
	ctx, pd, err := s.consumer.Context().GetTokenCtx(models.ServiceName_NUDR_DR, models.NrfNfManagementNfType_UDR)
	if err != nil {
		return "", pd, err
	}
	client := s.getDataSubscription(ue.UdrUri)
	trafficInfluSub := s.buildTrafficInfluSub(request)
	individualInfluenceDataSubscReq := DataRepository.CreateIndividualInfluenceDataSubscriptionRequest{
		TrafficInfluSub: &trafficInfluSub,
	}
	httpResp, localErr := client.InfluenceDataSubscriptionsCollectionApi.
		CreateIndividualInfluenceDataSubscription(ctx, &individualInfluenceDataSubscReq)
	if localErr == nil {
		locationHeader := httpResp.Location
		subscriptionID = locationHeader[strings.LastIndex(locationHeader, "/")+1:]
		logger.ConsumerLog.Debugf("Influence Data Subscription ID: %s", subscriptionID)
		return subscriptionID, nil, nil
	}
	if genericErr, ok := localErr.(openapi.GenericOpenAPIError); ok {
		if problemDetails, ok := genericErr.Model().(models.ProblemDetails); ok {
			return "", &problemDetails, nil
		}

		logger.ConsumerLog.Errorf("openapi error: %+v", localErr)
		return "", nil, err
	}

	return "", nil, localErr
}

func (s *nudrService) buildTrafficInfluSub(request models.SmPolicyContextData) models.TrafficInfluSub {
	trafficInfluSub := models.TrafficInfluSub{
		Dnns:             []string{request.Dnn},
		Snssais:          []models.Snssai{*request.SliceInfo},
		InternalGroupIds: request.InterGrpIds,
		Supis:            []string{request.Supi},
		NotificationUri: s.consumer.Context().GetIPv4Uri() +
			pcf_context.InfluenceDataUpdateNotifyUri + "/" +
			request.Supi + "/" + strconv.Itoa(int(request.PduSessionId)),
		// TODO: support expiry time and resend subscription when expired
	}
	return trafficInfluSub
}

func (s *nudrService) RemoveInfluenceDataSubscription(ue *pcf_context.UeContext, subscriptionID string) (
	problemDetails *models.ProblemDetails, err error,
) {
	if ue.UdrUri == "" {
		problemDetail := util.GetProblemDetail("Can't find corresponding UDR with UE", util.USER_UNKNOWN)
		logger.ConsumerLog.Warnf("Can't find corresponding UDR with UE[%s]", ue.Supi)
		return &problemDetail, nil
	}
	ctx, pd, err := s.consumer.Context().GetTokenCtx(models.ServiceName_NUDR_DR, models.NrfNfManagementNfType_UDR)
	if err != nil {
		return pd, err
	}
	client := s.getDataSubscription(ue.UdrUri)
	deleteIndividualInfluenceDataSubscriptionReq := DataRepository.DeleteIndividualInfluenceDataSubscriptionRequest{
		SubscriptionId: &subscriptionID,
	}
	_, localErr := client.IndividualInfluenceDataSubscriptionDocumentApi.
		DeleteIndividualInfluenceDataSubscription(ctx, &deleteIndividualInfluenceDataSubscriptionReq)
	if localErr == nil {
		logger.ConsumerLog.Debugf("DataRepository Remove Influence Data Subscription Status With No Err")
	}
	if genericErr, ok := localErr.(openapi.GenericOpenAPIError); ok {
		if problemDetails, ok := genericErr.Model().(models.ProblemDetails); ok {
			return &problemDetails, nil
		}

		logger.ConsumerLog.Errorf("openapi error: %+v", localErr)
		return nil, localErr
	}

	return nil, localErr
}
