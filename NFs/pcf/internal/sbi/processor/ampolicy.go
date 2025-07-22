package processor

import (
	"fmt"
	"net/http"
	"reflect"

	"github.com/gin-gonic/gin"
	"github.com/mohae/deepcopy"

	"github.com/free5gc/openapi"
	"github.com/free5gc/openapi/models"
	pcf_context "github.com/free5gc/pcf/internal/context"
	"github.com/free5gc/pcf/internal/logger"
	"github.com/free5gc/pcf/internal/util"
)

func (p *Processor) HandleDeletePoliciesPolAssoId(
	c *gin.Context,
	polAssoId string,
) {
	logger.AmPolicyLog.Infof("Handle AM Policy Association Delete")

	ue := p.Context().PCFUeFindByPolicyId(polAssoId)
	if ue == nil || ue.AMPolicyData[polAssoId] == nil {
		problemDetails := util.GetProblemDetail("polAssoId not found  in PCF", util.CONTEXT_NOT_FOUND)
		c.JSON(int(problemDetails.Status), problemDetails)
	}

	delete(ue.AMPolicyData, polAssoId)
	c.JSON(http.StatusNoContent, nil)
}

// PoliciesPolAssoIdGet -
func (p *Processor) HandleGetPoliciesPolAssoId(
	c *gin.Context,
	polAssoId string,
) {
	logger.AmPolicyLog.Infof("Handle AM Policy Association Get")

	ue := p.Context().PCFUeFindByPolicyId(polAssoId)
	if ue == nil || ue.AMPolicyData[polAssoId] == nil {
		problemDetails := util.GetProblemDetail("polAssoId not found  in PCF", util.CONTEXT_NOT_FOUND)
		c.JSON(int(problemDetails.Status), problemDetails)
		return
	}
	amPolicyData := ue.AMPolicyData[polAssoId]
	rsp := models.PcfAmPolicyControlPolicyAssociation{
		SuppFeat: amPolicyData.SuppFeat,
	}
	if amPolicyData.Rfsp != 0 {
		rsp.Rfsp = amPolicyData.Rfsp
	}
	if amPolicyData.ServAreaRes != nil {
		rsp.ServAreaRes = amPolicyData.ServAreaRes
	}
	if amPolicyData.Triggers != nil {
		rsp.Triggers = amPolicyData.Triggers
		for _, trigger := range amPolicyData.Triggers {
			if trigger == models.PcfAmPolicyControlRequestTrigger_PRA_CH {
				rsp.Pras = amPolicyData.Pras
				break
			}
		}
	}
	c.JSON(http.StatusOK, rsp)
}

func (p *Processor) HandleUpdatePostPoliciesPolAssoId(
	c *gin.Context,
	polAssoId string,
	policyAssociationUpdateRequest models.PcfAmPolicyControlPolicyAssociationUpdateRequest,
) {
	logger.AmPolicyLog.Infof("Handle AM Policy Association Update")

	response, problemDetails := p.UpdatePostPoliciesPolAssoIdProcedure(polAssoId, policyAssociationUpdateRequest)
	if response != nil {
		c.JSON(http.StatusOK, response)
	} else if problemDetails != nil {
		c.JSON(int(problemDetails.Status), problemDetails)
	}

	problemDetails = &models.ProblemDetails{
		Status: http.StatusForbidden,
		Cause:  "UNSPECIFIED",
	}
	c.JSON(int(problemDetails.Status), problemDetails)
}

func (p *Processor) UpdatePostPoliciesPolAssoIdProcedure(polAssoId string,
	policyAssociationUpdateRequest models.PcfAmPolicyControlPolicyAssociationUpdateRequest,
) (*models.PcfAmPolicyControlPolicyUpdate, *models.ProblemDetails) {
	ue := p.Context().PCFUeFindByPolicyId(polAssoId)
	if ue == nil || ue.AMPolicyData[polAssoId] == nil {
		problemDetails := util.GetProblemDetail("polAssoId not found  in PCF", util.CONTEXT_NOT_FOUND)
		return nil, &problemDetails
	}

	amPolicyData := ue.AMPolicyData[polAssoId]
	var response models.PcfAmPolicyControlPolicyUpdate
	if policyAssociationUpdateRequest.NotificationUri != "" {
		amPolicyData.NotificationUri = policyAssociationUpdateRequest.NotificationUri
	}
	if policyAssociationUpdateRequest.AltNotifIpv4Addrs != nil {
		amPolicyData.AltNotifIpv4Addrs = policyAssociationUpdateRequest.AltNotifIpv4Addrs
	}
	if policyAssociationUpdateRequest.AltNotifIpv6Addrs != nil {
		amPolicyData.AltNotifIpv6Addrs = policyAssociationUpdateRequest.AltNotifIpv6Addrs
	}
	for _, trigger := range policyAssociationUpdateRequest.Triggers {
		// TODO: Modify the value according to policies
		switch trigger {
		case models.PcfAmPolicyControlRequestTrigger_LOC_CH:
			// TODO: report to AF subscriber
			if policyAssociationUpdateRequest.UserLoc == nil {
				problemDetail := util.GetProblemDetail("UserLoc are nli", util.ERROR_REQUEST_PARAMETERS)
				logger.AmPolicyLog.Warnln(
					"UserLoc doesn't exist in Policy Association Requset Update while Triggers include LOC_CH")
				return nil, &problemDetail
			}
			amPolicyData.UserLoc = policyAssociationUpdateRequest.UserLoc
			logger.AmPolicyLog.Infof("Ue[%s] UserLocation %+v", ue.Supi, amPolicyData.UserLoc)
		case models.PcfAmPolicyControlRequestTrigger_PRA_CH:
			if policyAssociationUpdateRequest.PraStatuses == nil {
				problemDetail := util.GetProblemDetail("PraStatuses are nli", util.ERROR_REQUEST_PARAMETERS)
				logger.AmPolicyLog.Warnln("PraStatuses doesn't exist in Policy Association",
					"Requset Update while Triggers include PRA_CH")
				return nil, &problemDetail
			}
			for praId, praInfo := range policyAssociationUpdateRequest.PraStatuses {
				// TODO: report to AF subscriber
				logger.AmPolicyLog.Infof("Policy Association Presence Id[%s] change state to %s", praId, praInfo.PresenceState)
			}
		case models.PcfAmPolicyControlRequestTrigger_SERV_AREA_CH:
			if policyAssociationUpdateRequest.ServAreaRes == nil {
				problemDetail := util.GetProblemDetail("ServAreaRes are nli", util.ERROR_REQUEST_PARAMETERS)
				logger.AmPolicyLog.Warnln("ServAreaRes doesn't exist in Policy Association",
					"Requset Update while Triggers include SERV_AREA_CH")
				return nil, &problemDetail
			} else {
				amPolicyData.ServAreaRes = policyAssociationUpdateRequest.ServAreaRes
				response.ServAreaRes = policyAssociationUpdateRequest.ServAreaRes
			}
		case models.PcfAmPolicyControlRequestTrigger_RFSP_CH:
			if policyAssociationUpdateRequest.Rfsp == 0 {
				problemDetail := util.GetProblemDetail("Rfsp are nli", util.ERROR_REQUEST_PARAMETERS)
				logger.AmPolicyLog.Warnln("Rfsp doesn't exist in Policy Association Requset Update while Triggers include RFSP_CH")
				return nil, &problemDetail
			} else {
				amPolicyData.Rfsp = policyAssociationUpdateRequest.Rfsp
				response.Rfsp = policyAssociationUpdateRequest.Rfsp
			}
		}
	}
	// TODO: handle TraceReq
	// TODO: Change Request Trigger Policies if needed
	response.Triggers = amPolicyData.Triggers
	// TODO: Change Policies if needed
	// rsp.Pras
	return &response, nil
}

// Create AM Policy
func (p *Processor) HandlePostPolicies(
	c *gin.Context,
	polAssoId string,
	policyAssociationRequest models.PcfAmPolicyControlPolicyAssociationRequest,
) {
	logger.AmPolicyLog.Infof("Handle AM Policy Create Request")

	response, locationHeader, problemDetails := p.PostPoliciesProcedure(polAssoId, policyAssociationRequest)
	if response != nil {
		// TODO: set gin header
		c.Header("Location", locationHeader)
		c.JSON(http.StatusCreated, response)
		return
	} else if problemDetails != nil {
		c.JSON(int(problemDetails.Status), problemDetails)
		return
	}

	problemDetails = &models.ProblemDetails{
		Status: http.StatusForbidden,
		Cause:  "UNSPECIFIED",
	}
	c.JSON(int(problemDetails.Status), problemDetails)
}

func (p *Processor) PostPoliciesProcedure(polAssoId string,
	policyAssociationRequest models.PcfAmPolicyControlPolicyAssociationRequest,
) (*models.PcfAmPolicyControlPolicyAssociation, string, *models.ProblemDetails) {
	var response models.PcfAmPolicyControlPolicyAssociation
	pcfSelf := p.Context()
	var ue *pcf_context.UeContext
	if val, ok := pcfSelf.UePool.Load(policyAssociationRequest.Supi); ok {
		ue = val.(*pcf_context.UeContext)
	}
	if ue == nil {
		if newUe, err := pcfSelf.NewPCFUe(policyAssociationRequest.Supi); err != nil {
			// supi format dose not match "imsi-..."
			problemDetail := util.GetProblemDetail("Supi Format Error", util.ERROR_REQUEST_PARAMETERS)
			logger.AmPolicyLog.Errorln(err.Error())
			return nil, "", &problemDetail
		} else {
			ue = newUe
		}
	}
	udrUri := p.getUdrUri(ue)
	if udrUri == "" {
		// Can't find any UDR support this Ue
		pcfSelf.UePool.Delete(ue.Supi)
		problemDetail := util.GetProblemDetail("Ue is not supported in PCF", util.USER_UNKNOWN)
		logger.AmPolicyLog.Errorf("Ue[%s] is not supported in PCF", ue.Supi)
		return nil, "", &problemDetail
	}
	ue.UdrUri = udrUri

	response.Request = deepcopy.Copy(&policyAssociationRequest).(*models.PcfAmPolicyControlPolicyAssociationRequest)
	assolId := fmt.Sprintf("%s-%d", ue.Supi, ue.PolAssociationIDGenerator)
	amPolicy := ue.AMPolicyData[assolId]

	if amPolicy == nil || amPolicy.AmPolicyData == nil {
		rsp, problemDetail, err := p.Consumer().GetAccessAndMobilityPolicyData(ue)
		if err != nil {
			logger.AmPolicyLog.Errorf("Get Access And Mobility Policy Data Error Problem[%+v]", err)
			problemDetail := util.GetProblemDetail("Get Access And Mobility Policy Data error[%+v]", err.Error())
			return nil, "", &problemDetail
		} else if problemDetail != nil {
			logger.AmPolicyLog.Errorf("Get Access And Mobility Policy Data Failed Problem[%+v]", problemDetail)
			problemDetail := util.GetProblemDetail("Get Access And Mobility Policy Data Failed[%+v]", problemDetail.Cause)
			return nil, "", &problemDetail
		}

		if amPolicy == nil {
			amPolicy = ue.NewUeAMPolicyData(assolId, policyAssociationRequest)
		}
		amPolicy.AmPolicyData = rsp
	}

	// TODO: according to PCF Policy to determine ServAreaRes, Rfsp, SuppFeat
	// amPolicy.ServAreaRes =
	// amPolicy.Rfsp =
	var requestSuppFeat openapi.SupportedFeature
	if suppFeat, err := openapi.NewSupportedFeature(policyAssociationRequest.SuppFeat); err != nil {
		logger.AmPolicyLog.Warnln(err)
	} else {
		requestSuppFeat = suppFeat
	}
	amPolicy.SuppFeat = pcfSelf.PcfSuppFeats[models.
		ServiceName_NPCF_AM_POLICY_CONTROL].NegotiateWith(
		requestSuppFeat).String()
	if amPolicy.Rfsp != 0 {
		response.Rfsp = amPolicy.Rfsp
	}
	response.SuppFeat = amPolicy.SuppFeat
	// TODO: add Reports
	// rsp.Triggers
	// rsp.Pras
	ue.PolAssociationIDGenerator++
	// Create location header for update, delete, get
	locationHeader := util.GetResourceUri(models.ServiceName_NPCF_AM_POLICY_CONTROL, assolId)
	logger.AmPolicyLog.Tracef("AMPolicy association Id[%s] Create", assolId)

	// if consumer is AMF then subscribe this AMF Status
	if policyAssociationRequest.Guami != nil {
		// if policyAssociationRequest.Guami has been subscribed, then no need to subscribe again
		needSubscribe := true
		pcfSelf.AMFStatusSubsData.Range(func(key, value interface{}) bool {
			data := value.(pcf_context.AMFStatusSubscriptionData)
			for _, guami := range data.GuamiList {
				if reflect.DeepEqual(guami, *policyAssociationRequest.Guami) {
					needSubscribe = false
					break
				}
			}
			// if no need to subscribe => stop iteration
			return needSubscribe
		})

		if needSubscribe {
			logger.AmPolicyLog.Debugf("Subscribe AMF status change[GUAMI: %+v]", *policyAssociationRequest.Guami)
			amfUri := p.Consumer().SendNFInstancesAMF(pcfSelf.NrfUri,
				*policyAssociationRequest.Guami, models.ServiceName_NAMF_COMM)
			if amfUri != "" {
				problemDetails, err := p.Consumer().AmfStatusChangeSubscribe(amfUri,
					[]models.Guami{*policyAssociationRequest.Guami})
				if err != nil {
					logger.AmPolicyLog.Errorf("Subscribe AMF status change error[%+v]", err)
				} else if problemDetails != nil {
					logger.AmPolicyLog.Errorf("Subscribe AMF status change failed[%+v]", problemDetails)
				} else {
					amPolicy.Guami = policyAssociationRequest.Guami
				}
			}
		} else {
			logger.AmPolicyLog.Debugf("AMF status[GUAMI: %+v] has been subscribed", *policyAssociationRequest.Guami)
		}
	}
	return &response, locationHeader, nil
}

// Send AM Policy Update to AMF if policy has been terminated
func (p *Processor) SendAMPolicyTerminationRequestNotification(ue *pcf_context.UeContext,
	PolId string, request models.PcfAmPolicyControlTerminationNotification,
) {
	if ue == nil {
		logger.AmPolicyLog.Warnln("Policy Assocition Termination Request Notification Error[Ue is nil]")
		return
	}
	amPolicyData := ue.AMPolicyData[PolId]
	if amPolicyData == nil {
		logger.AmPolicyLog.Warnf(
			"Policy Assocition Termination Request Notification Error[Can't find polAssoId[%s] in UE(%s)]", PolId, ue.Supi)
		return
	}

	uri := amPolicyData.NotificationUri

	if uri != "" {
		pd, err := p.Consumer().SendAMPolicyAssociationPolicyAssocitionTerminationRequestNotification(uri, &request)
		if err != nil {
			logger.AmPolicyLog.Warnf("Policy Assocition Termination Request Notification Error[%s]", err.Error())
			return
		} else if pd != nil {
			logger.AmPolicyLog.Warnf("Policy Assocition Termination Request Notification Fault[%s]", pd.Detail)
			return
		}
		return
	}
}

// returns UDR Uri of Ue, if ue.UdrUri dose not exist, query NRF to get supported Udr Uri
func (p *Processor) getUdrUri(ue *pcf_context.UeContext) string {
	if ue.UdrUri != "" {
		return ue.UdrUri
	}
	return p.Consumer().SendNFInstancesUDR(p.Context().NrfUri, ue.Supi)
}
