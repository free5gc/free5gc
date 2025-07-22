package consumer

import (
	"regexp"
	"sync"

	amf_context "github.com/free5gc/amf/internal/context"
	"github.com/free5gc/amf/internal/logger"
	"github.com/free5gc/amf/pkg/factory"
	"github.com/free5gc/openapi"
	"github.com/free5gc/openapi/models"
	Npcf_AMPolicy "github.com/free5gc/openapi/pcf/AMPolicyControl"
)

type npcfService struct {
	consumer *Consumer

	AMPolicyMu sync.RWMutex

	AMPolicyClients map[string]*Npcf_AMPolicy.APIClient
}

func (s *npcfService) getAMPolicyClient(uri string) *Npcf_AMPolicy.APIClient {
	if uri == "" {
		return nil
	}
	s.AMPolicyMu.RLock()
	client, ok := s.AMPolicyClients[uri]
	if ok {
		s.AMPolicyMu.RUnlock()
		return client
	}

	configuration := Npcf_AMPolicy.NewConfiguration()
	configuration.SetBasePath(uri)
	client = Npcf_AMPolicy.NewAPIClient(configuration)

	s.AMPolicyMu.RUnlock()
	s.AMPolicyMu.Lock()
	defer s.AMPolicyMu.Unlock()
	s.AMPolicyClients[uri] = client
	return client
}

func (s *npcfService) AMPolicyControlCreate(
	ue *amf_context.AmfUe, anType models.AccessType,
) (*models.ProblemDetails, error) {
	client := s.getAMPolicyClient(ue.PcfUri)
	if client == nil {
		return nil, openapi.ReportError("pcf not found")
	}
	amfSelf := amf_context.GetSelf()
	ctx, _, err := amf_context.GetSelf().GetTokenCtx(models.ServiceName_NPCF_AM_POLICY_CONTROL,
		models.NrfNfManagementNfType_PCF)
	if err != nil {
		return nil, err
	}

	policyAssociationRequest := models.PcfAmPolicyControlPolicyAssociationRequest{
		NotificationUri: amfSelf.GetIPv4Uri() + factory.AmfCallbackResUriPrefix + "/am-policy/",
		Supi:            ue.Supi,
		Pei:             ue.Pei,
		Gpsi:            ue.Gpsi,
		AccessType:      anType,
		ServingPlmn: &models.PlmnIdNid{
			Mcc: ue.PlmnId.Mcc,
			Mnc: ue.PlmnId.Mnc,
		},
		Guami: &amfSelf.ServedGuamiList[0],
	}
	var policyAssociationreq Npcf_AMPolicy.CreateIndividualAMPolicyAssociationRequest

	policyAssociationreq.SetPcfAmPolicyControlPolicyAssociationRequest(policyAssociationRequest)

	if ue.AccessAndMobilitySubscriptionData != nil {
		policyAssociationRequest.Rfsp = ue.AccessAndMobilitySubscriptionData.RfspIndex
	}

	res, localErr := client.AMPolicyAssociationsCollectionApi.
		CreateIndividualAMPolicyAssociation(ctx, &policyAssociationreq)
	if localErr == nil {
		locationHeader := res.Location
		logger.ConsumerLog.Debugf("location header: %+v", locationHeader)
		ue.AmPolicyUri = locationHeader

		re := regexp.MustCompile("/policies/.*")
		match := re.FindStringSubmatch(locationHeader)

		ue.PolicyAssociationId = match[0][10:]
		ue.AmPolicyAssociation = &res.PcfAmPolicyControlPolicyAssociation

		if res.PcfAmPolicyControlPolicyAssociation.Triggers != nil {
			for _, trigger := range res.PcfAmPolicyControlPolicyAssociation.Triggers {
				if trigger == models.PcfAmPolicyControlRequestTrigger_LOC_CH {
					ue.RequestTriggerLocationChange = true
				}
				// if trigger == models.RequestTrigger_PRA_CH {
				// TODO: Presence Reporting Area handling (TS 23.503 6.1.2.5, TS 23.501 5.6.11)
				// }
			}
		}

		logger.ConsumerLog.Debugf("UE AM Policy Association ID: %s", ue.PolicyAssociationId)
		logger.ConsumerLog.Debugf("AmPolicyAssociation: %+v", ue.AmPolicyAssociation)
	} else {
		switch apiErr := localErr.(type) {
		// API error
		case openapi.GenericOpenAPIError:
			switch errorModel := apiErr.Model().(type) {
			case Npcf_AMPolicy.CreateIndividualAMPolicyAssociationError:
				return &errorModel.ProblemDetails, nil
			case error:
				return openapi.ProblemDetailsSystemFailure(errorModel.Error()), nil
			default:
				return nil, openapi.ReportError("openapi error")
			}
		case error:
			return openapi.ProblemDetailsSystemFailure(apiErr.Error()), nil
		default:
			return nil, openapi.ReportError("openapi error")
		}
	}
	return nil, nil
}

func (s *npcfService) AMPolicyControlUpdate(
	ue *amf_context.AmfUe, updateRequest models.PcfAmPolicyControlPolicyAssociationUpdateRequest,
) (problemDetails *models.ProblemDetails, err error) {
	client := s.getAMPolicyClient(ue.PcfUri)
	if client == nil {
		return nil, openapi.ReportError("pcf not found")
	}

	ctx, _, err := amf_context.GetSelf().GetTokenCtx(models.ServiceName_NPCF_AM_POLICY_CONTROL,
		models.NrfNfManagementNfType_PCF)
	if err != nil {
		return nil, err
	}

	var policyUpdateReq Npcf_AMPolicy.ReportObservedEventTriggersForIndividualAMPolicyAssociationRequest

	policyUpdateReq.SetPolAssoId(ue.PolicyAssociationId)
	policyUpdateReq.SetPcfAmPolicyControlPolicyAssociationUpdateRequest(updateRequest)

	res, localErr := client.IndividualAMPolicyAssociationDocumentApi.
		ReportObservedEventTriggersForIndividualAMPolicyAssociation(ctx, &policyUpdateReq)
	if localErr == nil {
		if res.PcfAmPolicyControlPolicyUpdate.ServAreaRes != nil {
			ue.AmPolicyAssociation.ServAreaRes = res.PcfAmPolicyControlPolicyUpdate.ServAreaRes
		}
		if res.PcfAmPolicyControlPolicyUpdate.Rfsp != 0 {
			ue.AmPolicyAssociation.Rfsp = res.PcfAmPolicyControlPolicyUpdate.Rfsp
		}
		ue.AmPolicyAssociation.Triggers = res.PcfAmPolicyControlPolicyUpdate.Triggers
		ue.RequestTriggerLocationChange = false
		for _, trigger := range res.PcfAmPolicyControlPolicyUpdate.Triggers {
			if trigger == models.PcfAmPolicyControlRequestTrigger_LOC_CH {
				ue.RequestTriggerLocationChange = true
			}
			// if trigger == models.RequestTrigger_PRA_CH {
			// TODO: Presence Reporting Area handling (TS 23.503 6.1.2.5, TS 23.501 5.6.11)
			// }
		}
	} else {
		switch apiErr := localErr.(type) {
		// API error
		case openapi.GenericOpenAPIError:
			switch errorModel := apiErr.Model().(type) {
			case Npcf_AMPolicy.ReportObservedEventTriggersForIndividualAMPolicyAssociationError:
				return &errorModel.ProblemDetails, nil
			case error:
				return openapi.ProblemDetailsSystemFailure(errorModel.Error()), nil
			default:
				err = openapi.ReportError("openapi error")
			}
		case error:
			return openapi.ProblemDetailsSystemFailure(apiErr.Error()), nil
		default:
			err = openapi.ReportError("server no response")
		}
	}
	return nil, err
}

func (s *npcfService) AMPolicyControlDelete(ue *amf_context.AmfUe) (problemDetails *models.ProblemDetails, err error) {
	client := s.getAMPolicyClient(ue.PcfUri)
	if client == nil {
		return nil, openapi.ReportError("pcf not found")
	}

	ctx, _, ctxErr := amf_context.GetSelf().GetTokenCtx(models.ServiceName_NPCF_AM_POLICY_CONTROL,
		models.NrfNfManagementNfType_PCF)
	if ctxErr != nil {
		return nil, ctxErr
	}

	var deleteReq Npcf_AMPolicy.DeleteIndividualAMPolicyAssociationRequest
	deleteReq.SetPolAssoId(ue.PolicyAssociationId)

	_, err = client.IndividualAMPolicyAssociationDocumentApi.DeleteIndividualAMPolicyAssociation(ctx, &deleteReq)
	if err == nil {
		ue.RemoveAmPolicyAssociation()
	} else {
		switch apiErr := err.(type) {
		// API error
		case openapi.GenericOpenAPIError:
			switch errorModel := apiErr.Model().(type) {
			case Npcf_AMPolicy.DeleteIndividualAMPolicyAssociationError:
				return &errorModel.ProblemDetails, nil
			case error:
				return openapi.ProblemDetailsSystemFailure(errorModel.Error()), nil
			default:
				err = openapi.ReportError("openapi error")
			}
		case error:
			return openapi.ProblemDetailsSystemFailure(apiErr.Error()), nil
		default:
			err = openapi.ReportError("server no response")
		}
	}
	return nil, err
}
