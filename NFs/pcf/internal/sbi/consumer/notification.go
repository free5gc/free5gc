package consumer

import (
	"reflect"
	"sync"

	"github.com/free5gc/openapi/models"
	"github.com/free5gc/openapi/pcf/AMPolicyControl"
	pcf_context "github.com/free5gc/pcf/internal/context"
	"github.com/free5gc/pcf/internal/logger"
	"github.com/free5gc/pcf/internal/util"
)

type npcfService struct {
	consumer *Consumer

	nfAMPolicyControlMu sync.RWMutex

	nfAMPolicyControlClient map[string]*AMPolicyControl.APIClient
}

func (s *npcfService) getAMPolicyControl(uri string) *AMPolicyControl.APIClient {
	if uri == "" {
		return nil
	}
	s.nfAMPolicyControlMu.RLock()
	client, ok := s.nfAMPolicyControlClient[uri]
	if ok {
		defer s.nfAMPolicyControlMu.RUnlock()
		return client
	}

	configuration := AMPolicyControl.NewConfiguration()
	configuration.SetBasePath(uri)
	client = AMPolicyControl.NewAPIClient(configuration)

	s.nfAMPolicyControlMu.RUnlock()
	s.nfAMPolicyControlMu.Lock()
	defer s.nfAMPolicyControlMu.Unlock()
	s.nfAMPolicyControlClient[uri] = client
	return client
}

// Send AM Policy Update to AMF if policy has changed
func (s *npcfService) SendAMPolicyUpdateNotification(ue *pcf_context.UeContext,
	PolId string, request models.PcfAmPolicyControlPolicyUpdate,
) {
	if ue == nil {
		logger.ConsumerLog.Warnln("Policy Update Notification Error[Ue is nil]")
		return
	}
	amPolicyData := ue.AMPolicyData[PolId]
	if amPolicyData == nil {
		logger.ConsumerLog.Warnf("Policy Update Notification Error[Can't find polAssoId[%s] in UE(%s)]", PolId, ue.Supi)
		return
	}

	uri := amPolicyData.NotificationUri

	if uri == "" {
		logger.ConsumerLog.Warnln("NPcf client can't find call back uri")
		return
	}

	if reflect.DeepEqual(request, models.PcfAmPolicyControlPolicyUpdate{}) {
		logger.ConsumerLog.Warnln("SendAMPolicyUpdateNotification request is nil")
		return
	}

	ctx, problemDetails, err := s.consumer.Context().GetTokenCtx(models.ServiceName_NPCF_AM_POLICY_CONTROL,
		models.NrfNfManagementNfType_PCF)
	if err != nil {
		logger.ConsumerLog.Warnf("Policy Update Notification Error[%s]", err.Error())
		return
	} else if problemDetails != nil {
		logger.ConsumerLog.Warnf("Policy Update Notification Fault[%s]", problemDetails.Detail)
		return
	}

	client := s.getAMPolicyControl(uri)
	param := AMPolicyControl.CreateIndividualAMPolicyAssociationPolicyUpdateNotificationPostRequest{
		PcfAmPolicyControlPolicyUpdate: &request,
	}
	rsp, err := client.AMPolicyAssociationsCollectionApi.
		CreateIndividualAMPolicyAssociationPolicyUpdateNotificationPost(
			ctx, uri, &param)
	if err != nil {
		logger.ConsumerLog.Warnf("SendAMPolicyUpdateNotification function in consumer Error[%s]",
			err.Error())
		return
	} else if rsp == nil {
		logger.ConsumerLog.Warnf("SendAMPolicyUpdateNotification function in consumer Failed[Response is nil]")
		return
	}
}

func (s *npcfService) SendAMPolicyAssociationPolicyAssocitionTerminationRequestNotification(
	uri string, request *models.PcfAmPolicyControlTerminationNotification,
) (
	problemDetails *models.ProblemDetails, err error,
) {
	if uri == "" {
		problemDetail := util.GetProblemDetail("NPcf client can't find call back uri",
			"SendAMPolicyAssociationPolicyAssocitionTerminationRequestNotification Can't find URI")
		return &problemDetail, nil
	}

	if request == nil {
		problemDetail := util.GetProblemDetail(
			"SendAMPolicyAssociationPolicyAssocitionTerminationRequestNotification request is nil",
			"SendAMPolicyAssociationPolicyAssocitionTerminationRequestNotification request is nil")
		return &problemDetail, nil
	}

	ctx, problemDetails, err := s.consumer.Context().GetTokenCtx(
		models.ServiceName_NPCF_AM_POLICY_CONTROL,
		models.NrfNfManagementNfType_PCF)
	if err != nil {
		return nil, err
	} else if problemDetails != nil {
		return problemDetails, nil
	}

	client := s.getAMPolicyControl(uri)
	param := AMPolicyControl.
		CreateIndividualAMPolicyAssociationPolicyAssocitionTerminationRequestNotificationPostRequest{
		PcfAmPolicyControlTerminationNotification: request,
	}
	_, err = client.AMPolicyAssociationsCollectionApi.
		CreateIndividualAMPolicyAssociationPolicyAssocitionTerminationRequestNotificationPost(
			ctx, uri, &param)
	if err != nil {
		logger.AmPolicyLog.Warnf(
			"CreateIndividualAMPolicyAssociationPolicyAssocitionTerminationRequestNotificationPost Error[%s]",
			err.Error())
		return nil, err
	}
	return nil, nil
}
