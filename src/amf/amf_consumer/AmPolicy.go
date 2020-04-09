package amf_consumer

import (
	"context"
	"free5gc/lib/Npcf_AMPolicy"
	"free5gc/lib/openapi/common"
	"free5gc/lib/openapi/models"
	"free5gc/src/amf/amf_context"
	"free5gc/src/amf/logger"
	"regexp"
)

func AMPolicyControlCreate(ue *amf_context.AmfUe, anType models.AccessType) (problemDetails *models.ProblemDetails, err error) {

	configuration := Npcf_AMPolicy.NewConfiguration()
	configuration.SetBasePath(ue.PcfUri)
	client := Npcf_AMPolicy.NewAPIClient(configuration)

	amfSelf := amf_context.AMF_Self()

	policyAssociationRequest := models.PolicyAssociationRequest{
		NotificationUri: amfSelf.GetIPv4Uri() + "/namf-callback/v1/am-policy/",
		Supi:            ue.Supi,
		Pei:             ue.Pei,
		Gpsi:            ue.Gpsi,
		AccessType:      anType,
		ServingPlmn: &models.NetworkId{
			Mcc: ue.PlmnId.Mcc,
			Mnc: ue.PlmnId.Mnc,
		},
		Guami: &amfSelf.ServedGuamiList[0],
	}

	if ue.AccessAndMobilitySubscriptionData != nil {
		policyAssociationRequest.Rfsp = ue.AccessAndMobilitySubscriptionData.RfspIndex
	}

	res, httpResp, localErr := client.DefaultApi.PoliciesPost(context.Background(), policyAssociationRequest)
	if localErr == nil {
		locationHeader := httpResp.Header.Get("Location")
		logger.ConsumerLog.Debugf("location header: %+v", locationHeader)
		ue.AmPolicyUri = locationHeader

		re := regexp.MustCompile("/policies/.*")
		match := re.FindStringSubmatch(locationHeader)

		ue.PolicyAssociationId = match[0][10:]
		ue.AmPolicyAssociation = &res

		if res.Triggers != nil {
			for _, trigger := range res.Triggers {
				if trigger == models.RequestTrigger_LOC_CH {
					ue.RequestTriggerLocationChange = true
				}
				if trigger == models.RequestTrigger_PRA_CH {
					// TODO: Presence Reporting Area handling (TS 23.503 6.1.2.5, TS 23.501 5.6.11)
				}
			}
		}

		logger.ConsumerLog.Debugf("UE AM Policy Association ID: %s", ue.PolicyAssociationId)
		logger.ConsumerLog.Debugf("AmPolicyAssociation: %+v", ue.AmPolicyAssociation)
	} else if httpResp != nil {
		if httpResp.Status != localErr.Error() {
			err = localErr
			return
		}
		problem := localErr.(common.GenericOpenAPIError).Model().(models.ProblemDetails)
		problemDetails = &problem
	} else {
		err = common.ReportError("server no response")
	}
	return
}

func AMPolicyControlUpdate(ue *amf_context.AmfUe, updateRequest models.PolicyAssociationUpdateRequest) (problemDetails *models.ProblemDetails, err error) {
	configuration := Npcf_AMPolicy.NewConfiguration()
	configuration.SetBasePath(ue.PcfUri)
	client := Npcf_AMPolicy.NewAPIClient(configuration)

	res, httpResp, localErr := client.DefaultApi.PoliciesPolAssoIdUpdatePost(context.Background(), ue.PolicyAssociationId, updateRequest)
	if localErr == nil {
		if res.ServAreaRes != nil {
			ue.AmPolicyAssociation.ServAreaRes = res.ServAreaRes
		}
		if res.Rfsp != 0 {
			ue.AmPolicyAssociation.Rfsp = res.Rfsp
		}
		ue.AmPolicyAssociation.Triggers = res.Triggers
		ue.RequestTriggerLocationChange = false
		for _, trigger := range res.Triggers {
			if trigger == models.RequestTrigger_LOC_CH {
				ue.RequestTriggerLocationChange = true
			}
			if trigger == models.RequestTrigger_PRA_CH {
				// TODO: Presence Reporting Area handling (TS 23.503 6.1.2.5, TS 23.501 5.6.11)
			}
		}
		return
	} else if httpResp != nil {
		if httpResp.Status != localErr.Error() {
			err = localErr
			return
		}
		problem := localErr.(common.GenericOpenAPIError).Model().(models.ProblemDetails)
		problemDetails = &problem
	} else {
		err = common.ReportError("server no response")
	}
	return
}

func AMPolicyControlDelete(ue *amf_context.AmfUe) (problemDetails *models.ProblemDetails, err error) {

	configuration := Npcf_AMPolicy.NewConfiguration()
	configuration.SetBasePath(ue.PcfUri)
	client := Npcf_AMPolicy.NewAPIClient(configuration)

	httpResp, localErr := client.DefaultApi.PoliciesPolAssoIdDelete(context.Background(), ue.PolicyAssociationId)
	if localErr == nil {
		ue.RemoveAmPolicyAssociation()
	} else if httpResp != nil {
		if httpResp.Status != localErr.Error() {
			err = localErr
			return
		}
		problem := localErr.(common.GenericOpenAPIError).Model().(models.ProblemDetails)
		problemDetails = &problem
	} else {
		err = common.ReportError("server no response")
	}

	return
}
