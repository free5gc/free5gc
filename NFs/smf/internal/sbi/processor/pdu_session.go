package processor

import (
	"encoding/hex"
	"errors"
	"fmt"
	"net"
	"net/http"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/free5gc/nas"
	"github.com/free5gc/nas/nasMessage"
	"github.com/free5gc/openapi"
	"github.com/free5gc/openapi/models"
	"github.com/free5gc/openapi/udm/SubscriberDataManagement"
	"github.com/free5gc/pfcp/pfcpType"
	smf_context "github.com/free5gc/smf/internal/context"
	"github.com/free5gc/smf/internal/logger"
	smf_errors "github.com/free5gc/smf/pkg/errors"
	"github.com/free5gc/smf/pkg/factory"
)

func (p *Processor) HandlePDUSessionSMContextCreate(
	c *gin.Context,
	request models.PostSmContextsRequest,
	isDone <-chan struct{},
) {
	// GSM State
	// PDU Session Establishment Accept/Reject
	var response models.PostSmContextsResponse201
	response.JsonData = new(models.SmfPduSessionSmContextCreatedData)
	logger.PduSessLog.Infoln("In HandlePDUSessionSMContextCreate")

	// Check has PDU Session Establishment Request
	m := nas.NewMessage()
	if err := m.GsmMessageDecode(&request.BinaryDataN1SmMessage); err != nil ||
		m.GsmHeader.GetMessageType() != nas.MsgTypePDUSessionEstablishmentRequest {
		logger.PduSessLog.Warnln("GsmMessageDecode Error: ", err)
		postSmContextsError := models.PostSmContextsError{
			JsonData: &models.SmContextCreateError{
				Error: &smf_errors.N1SmError,
			},
		}
		c.JSON(http.StatusForbidden, postSmContextsError)
		return
	}

	createData := request.JsonData
	// Check duplicate SM Context
	if dup_smCtx := smf_context.GetSMContextById(createData.Supi, createData.PduSessionId); dup_smCtx != nil {
		p.HandlePDUSessionSMContextLocalRelease(dup_smCtx, createData)
	}

	smContext := smf_context.NewSMContext(createData.Supi, createData.PduSessionId)
	smContext.SetState(smf_context.ActivePending)
	smContext.SmfPduSessionSmContextCreateData = createData
	smContext.SmStatusNotifyUri = createData.SmContextStatusUri

	smContext.SMLock.Lock()
	needUnlock := true
	defer func() {
		if needUnlock {
			smContext.SMLock.Unlock()
		}
	}()

	upi := smf_context.GetUserPlaneInformation()
	upi.Mu.RLock()
	defer upi.Mu.RUnlock()

	// DNN Information from config
	smContext.DNNInfo = smf_context.RetrieveDnnInformation(smContext.SNssai, smContext.Dnn)
	if smContext.DNNInfo == nil {
		logger.PduSessLog.Errorf("S-NSSAI[sst: %d, sd: %s] DNN[%s] not matched DNN Config",
			smContext.SNssai.Sst, smContext.SNssai.Sd, smContext.Dnn)
	}
	smContext.Log.Debugf("S-NSSAI[sst: %d, sd: %s] DNN[%s]",
		smContext.SNssai.Sst, smContext.SNssai.Sd, smContext.Dnn)

	// Query UDM
	if problemDetails, err := p.Consumer().SendNFDiscoveryUDM(); err != nil {
		smContext.Log.Warnf("Send NF Discovery Serving UDM Error[%v]", err)
	} else if problemDetails != nil {
		smContext.Log.Warnf("Send NF Discovery Serving UDM Problem[%+v]", problemDetails)
	} else {
		smContext.Log.Infoln("Send NF Discovery Serving UDM Successfully")
	}

	smPlmnID := createData.Guami.PlmnId

	smDataParams := &SubscriberDataManagement.GetSmDataRequest{
		Dnn: &createData.Dnn,
		PlmnId: &models.PlmnId{
			Mcc: smPlmnID.Mcc,
			Mnc: smPlmnID.Mnc,
		},
		SingleNssai: smContext.SNssai,
	}

	ctx, _, oauthErr := smf_context.GetSelf().GetTokenCtx(models.ServiceName_NUDM_SDM, models.NrfNfManagementNfType_UDM)
	if oauthErr != nil {
		smContext.Log.Errorf("Get Token Context Error[%v]", oauthErr)
		return
	}

	if sessSubData, err := p.Consumer().
		GetSmData(ctx, smContext.Supi, smDataParams); err != nil {
		smContext.Log.Errorln("Get SessionManagementSubscriptionData error:", err)
	} else {
		if len(sessSubData) > 0 {
			smContext.DnnConfiguration = sessSubData[0].DnnConfigurations[smContext.Dnn]
			// UP Security info present in session management subscription data
			if smContext.DnnConfiguration.UpSecurity != nil {
				smContext.UpSecurity = smContext.DnnConfiguration.UpSecurity
			}
		} else {
			smContext.Log.Errorln("SessionManagementSubscriptionData from UDM is nil")
		}
	}

	doSubscribe := false
	defer func() {
		if doSubscribe {
			if !p.Context().Ues.UeExists(smContext.Supi) {
				if problemDetails, err := p.Consumer().
					Subscribe(ctx, smContext, smPlmnID); problemDetails != nil {
					smContext.Log.Errorln("SDM Subscription Failed Problem:", problemDetails)
				} else if err != nil {
					smContext.Log.Errorln("SDM Subscription Error:", err)
				}
			} else {
				p.Context().Ues.IncrementPduSessionCount(smContext.Supi)
			}
		}
	}()

	establishmentRequest := m.PDUSessionEstablishmentRequest
	if err := HandlePDUSessionEstablishmentRequest(smContext, establishmentRequest); err != nil {
		smContext.Log.Errorf("PDU Session Establishment fail by %s", err)
		gsmError := &GSMError{}
		if errors.As(err, &gsmError) {
			p.makeEstRejectResAndReleaseSMContext(c, smContext,
				gsmError.GSMCause,
				&smf_errors.N1SmError)
			return
		}
		p.makeEstRejectResAndReleaseSMContext(c, smContext,
			nasMessage.Cause5GSMRequestRejectedUnspecified,
			&smf_errors.N1SmError)
		return
	}

	// Discover and new Namf_Comm client for use later
	if problemDetails, err := p.Consumer().SendNFDiscoveryServingAMF(smContext); err != nil {
		smContext.Log.Warnf("Send NF Discovery Serving AMF Error[%v]", err)
	} else if problemDetails != nil {
		smContext.Log.Warnf("Send NF Discovery Serving AMF Problem[%+v]", problemDetails)
	} else {
		smContext.Log.Traceln("Send NF Discovery Serving AMF successfully")
	}

	for _, service := range smContext.AMFProfile.NfServices {
		if service.ServiceName == models.ServiceName_NAMF_COMM {
			smContext.CommunicationClientApiPrefix = service.ApiPrefix
		}
	}

	if err := smContext.AllocUeIP(); err != nil {
		smContext.SetState(smf_context.InActive)
		smContext.Log.Errorf("PDUSessionSMContextCreate err: %v", err)
		p.makeEstRejectResAndReleaseSMContext(c, smContext,
			nasMessage.Cause5GSMInsufficientResourcesForSpecificSliceAndDNN,
			&smf_errors.InsufficientResourceSliceDnn)
		return
	}

	if err := p.Consumer().PCFSelection(smContext); err != nil {
		smContext.Log.Errorln("pcf selection error:", err)
	}

	smPolicyID, smPolicyDecision, err := p.Consumer().SendSMPolicyAssociationCreate(smContext)
	if err != nil {
		if openapiError, ok := err.(openapi.GenericOpenAPIError); ok {
			problemDetails := openapiError.Model().(models.ProblemDetails)
			smContext.Log.Errorln("setup sm policy association failed:", err, problemDetails)
			smContext.SetState(smf_context.InActive)
			if problemDetails.Cause == "USER_UNKNOWN" {
				p.makeEstRejectResAndReleaseSMContext(c, smContext,
					nasMessage.Cause5GSMRequestRejectedUnspecified,
					&smf_errors.SubscriptionDenied)
				return
			}
		}
		p.makeEstRejectResAndReleaseSMContext(c, smContext,
			nasMessage.Cause5GSMNetworkFailure,
			&smf_errors.NetworkFailure)
		return
	}
	smContext.SMPolicyID = smPolicyID

	// PDU　session create is a charging event
	logger.PduSessLog.Infof("CHF Selection for SMContext SUPI[%s] PDUSessionID[%d]\n",
		smContext.Supi, smContext.PDUSessionID)
	if err = p.Consumer().CHFSelection(smContext); err != nil {
		logger.PduSessLog.Errorln("chf selection error:", err)
	} else {
		p.CreateChargingSession(smContext)
	}

	// Update SessionRule from decision
	if err = smContext.ApplySessionRules(smPolicyDecision); err != nil {
		smContext.Log.Errorf("PDUSessionSMContextCreate err: %v", err)
		p.makeEstRejectResAndReleaseSMContext(c, smContext,
			nasMessage.Cause5GSMRequestRejectedUnspecified,
			&smf_errors.SubscriptionDenied)
		return
	}

	// If PCF prepares default Pcc Rule, SMF do not need to create defaultDataPath.
	if err = smContext.ApplyPccRules(smPolicyDecision); err != nil {
		smContext.Log.Errorf("apply sm policy decision error: %+v", err)
	}

	// SelectDefaultDataPath() will create a default data path if default data path is not found.
	if err = smContext.SelectDefaultDataPath(); err != nil {
		smContext.SetState(smf_context.InActive)
		smContext.Log.Errorf("PDUSessionSMContextCreate err: %v", err)
		p.makeEstRejectResAndReleaseSMContext(c, smContext,
			nasMessage.Cause5GSMInsufficientResourcesForSpecificSliceAndDNN,
			&smf_errors.InsufficientResourceSliceDnn)
		return
	}

	// SelectULCLDataPaths() will create other paths if ULCL is enabled.
	if err = smContext.SelectULCLDataPaths(); err != nil {
		smContext.SetState(smf_context.InActive)
		smContext.Log.Errorf("PDUSessionSMContextCreate err: %v", err)
		p.makeEstRejectResAndReleaseSMContext(c, smContext,
			nasMessage.Cause5GSMInsufficientResourcesForSpecificSliceAndDNN,
			&smf_errors.InsufficientResourceSliceDnn)
		return
	}

	// generate goroutine to handle PFCP and
	// reply PDUSessionSMContextCreate rsp immediately
	needUnlock = false
	go func() {
		defer smContext.SMLock.Unlock()

		smContext.SendUpPathChgNotification("EARLY", SendUpPathChgEventExposureNotification)

		handler := func(smContext *smf_context.SMContext, success bool) {
			p.EstHandler(isDone, smContext, success)
		}

		ActivateUPFSession(smContext, handler)

		smContext.SendUpPathChgNotification("LATE", SendUpPathChgEventExposureNotification)

		smContext.PostRemoveDataPath()
	}()

	doSubscribe = true
	response.JsonData = smContext.BuildCreatedData()

	// default location value will only be used in test environment
	// in real environment, location value will be formatted as a full URI
	location := smContext.Ref // this is the default location value
	if c.Request != nil {
		protocol := "http"
		if c.Request.TLS != nil {
			protocol += "s"
		}
		smContextRefParts := strings.Split(smContext.Ref, ":")
		if len(smContextRefParts) > 2 {
			location = fmt.Sprintf("%s://%s%s/%s",
				protocol,
				c.Request.Host,
				strings.TrimSuffix(c.Request.URL.Path, "/"),
				smContextRefParts[2])
		} else {
			logger.PduSessLog.Errorln("smContext.Ref(uuid) format is incorrect")
		}
	}
	c.Header("Location", location)
	c.JSON(http.StatusCreated, response)
}

func (p *Processor) HandlePDUSessionSMContextUpdate(
	c *gin.Context,
	body models.UpdateSmContextRequest,
	smContextRef string,
) {
	// GSM State
	// PDU Session Modification Reject(Cause Value == 43 || Cause Value != 43)/Complete
	// PDU Session Release Command/Complete
	var buf []byte
	var n2Buf []byte
	var err error

	smContext := smf_context.GetSMContextByRef(smContextRef)

	upi := smf_context.GetUserPlaneInformation()
	upi.Mu.RLock()
	defer upi.Mu.RUnlock()

	if smContext == nil {
		logger.PduSessLog.Warnf("SMContext[%s] is not found", smContextRef)

		updateSmContextError := models.UpdateSmContextResponse400{
			JsonData: &models.SmContextUpdateError{
				UpCnxState: models.UpCnxState_DEACTIVATED,
				Error: &models.SmfPduSessionExtProblemDetails{
					Type:   "Resource Not Found",
					Title:  "SMContext Ref is not found",
					Status: http.StatusNotFound,
				},
			},
		}
		c.JSON(http.StatusNotFound, updateSmContextError)
		return
	}

	smContext.SMLock.Lock()
	defer smContext.SMLock.Unlock()

	var sendPFCPModification bool
	var pfcpResponseStatus smf_context.PFCPSessionResponseStatus
	var response models.UpdateSmContextResponse200
	response.JsonData = new(models.SmContextUpdatedData)

	smContextUpdateData := body.JsonData

	if body.BinaryDataN1SmMessage != nil {
		m := nas.NewMessage()
		err = m.GsmMessageDecode(&body.BinaryDataN1SmMessage)
		smContext.Log.Tracef("N1 Message: %s", hex.EncodeToString(body.BinaryDataN1SmMessage))
		if err != nil {
			smContext.Log.Errorf("N1 Message parse failed: %v", err)
			updateSmContextError := models.UpdateSmContextResponse400{
				JsonData: &models.SmContextUpdateError{
					Error: &smf_errors.N1SmError,
				},
			} // Depends on the reason why N4 fail
			c.JSON(http.StatusForbidden, updateSmContextError)
			return
		}

		switch m.GsmHeader.GetMessageType() {
		case nas.MsgTypePDUSessionReleaseRequest:
			smContext.CheckState(smf_context.Active)
			// Wait till the state becomes Active again
			// TODO: implement sleep wait in concurrent architecture

			HandlePDUSessionReleaseRequest(smContext, m.PDUSessionReleaseRequest)
			if smContext.SelectedUPF != nil && smContext.PDUAddress != nil {
				smContext.Log.Infof("Release IP[%s]", smContext.PDUAddress)
				upi.ReleaseUEIP(smContext.SelectedUPF, smContext.PDUAddress, smContext.UseStaticIP)
				smContext.PDUAddress = nil
				// keep SelectedUPF until PDU Session Release is completed
			}

			// remove SM Policy Association
			if smContext.SMPolicyID != "" {
				if err = p.Consumer().SendSMPolicyAssociationTermination(smContext); err != nil {
					smContext.Log.Errorf("SM Policy Termination failed: %s", err)
				} else {
					smContext.SMPolicyID = ""
				}
			}

			if smf_context.GetSelf().Ues.UeExists(smContext.Supi) {
				problemDetails, clientErr := p.Consumer().UnSubscribe(smContext)
				if problemDetails != nil {
					logger.PduSessLog.Errorf("SDM UnSubscription Failed Problem[%+v]", problemDetails)
				} else if clientErr != nil {
					logger.PduSessLog.Errorf("SDM UnSubscription Error[%+v]", err)
				}
			}

			if smContext.UeCmRegistered {
				problemDetails, errUeCmDeregistration := p.Consumer().UeCmDeregistration(smContext)
				if problemDetails != nil {
					if problemDetails.Cause != CONTEXT_NOT_FOUND {
						logger.PduSessLog.Errorf("UECM_DeRegistration Failed Problem[%+v]", problemDetails)
					}
				} else if errUeCmDeregistration != nil {
					logger.PduSessLog.Errorf("UECM_DeRegistration Error[%+v]", errUeCmDeregistration)
				} else {
					logger.PduSessLog.Traceln("UECM_DeRegistration successful")
				}
			}

			cause := nasMessage.Cause5GSMRegularDeactivation
			if m.PDUSessionReleaseRequest.Cause5GSM != nil {
				cause = m.PDUSessionReleaseRequest.Cause5GSM.GetCauseValue()
			}

			if buf, err = smf_context.
				BuildGSMPDUSessionReleaseCommand(smContext, cause, true); err != nil {
				smContext.Log.Errorf("Build GSM PDUSessionReleaseCommand failed: %+v", err)
			} else {
				response.BinaryDataN1SmMessage = buf
				response.JsonData.N1SmMsg = &models.RefToBinaryData{ContentId: "PDUSessionReleaseCommand"}
				p.sendGSMPDUSessionReleaseCommand(smContext, buf)
			}

			// Only send N2 PDU Session Resource Release when UP connection is active
			if smContext.UpCnxState == models.UpCnxState_ACTIVATED {
				if buf, err = smf_context.
					BuildPDUSessionResourceReleaseCommandTransfer(smContext); err != nil {
					smContext.Log.Errorf("Build PDUSessionResourceReleaseCommandTransfer failed: %+v", err)
				} else {
					response.JsonData.N2SmInfoType = models.N2SmInfoType_PDU_RES_REL_CMD
					response.BinaryDataN2SmInformation = buf
					response.JsonData.N2SmInfo = &models.RefToBinaryData{ContentId: "PDUResourceReleaseCommand"}
				}
			}

			smContext.SetState(smf_context.PFCPModification)

			pfcpResponseStatus = releaseSession(smContext)
		case nas.MsgTypePDUSessionReleaseComplete:
			smContext.CheckState(smf_context.InActivePending)
			// Wait till the state becomes Active again
			// TODO: implement sleep wait in concurrent architecture

			smContext.SetState(smf_context.InActive)
			response.JsonData.UpCnxState = models.UpCnxState_DEACTIVATED
			smContext.StopT3592()

			// If CN tunnel resource is released, should
			if smContext.Tunnel.ANInformation.IPAddress == nil {
				p.RemoveSMContextFromAllNF(smContext, true)
			}
		case nas.MsgTypePDUSessionModificationRequest:
			if rsp, errHandleReq := p.
				HandlePDUSessionModificationRequest(smContext, m.PDUSessionModificationRequest); errHandleReq != nil {
				if buf, err = smf_context.BuildGSMPDUSessionModificationReject(smContext); err != nil {
					smContext.Log.Errorf("build GSM PDUSessionModificationReject failed: %+v", err)
				} else {
					response.BinaryDataN1SmMessage = buf
				}
			} else {
				if buf, err = rsp.PlainNasEncode(); err != nil {
					smContext.Log.Errorf("build GSM PDUSessionModificationCommand failed: %+v", err)
				} else {
					response.BinaryDataN1SmMessage = buf
					p.sendGSMPDUSessionModificationCommand(smContext, buf)
				}
			}

			if buf, err = smf_context.BuildPDUSessionResourceModifyRequestTransfer(smContext); err != nil {
				smContext.Log.Errorf("build N2 BuildPDUSessionResourceModifyRequestTransfer failed: %v", err)
			} else {
				response.BinaryDataN2SmInformation = buf
				response.JsonData.N2SmInfo = &models.RefToBinaryData{ContentId: "PDU_RES_MOD"}
				response.JsonData.N2SmInfoType = models.N2SmInfoType_PDU_RES_MOD_REQ
			}

			response.JsonData.N1SmMsg = &models.RefToBinaryData{ContentId: "PDUSessionModificationReject"}
			c.Render(http.StatusOK, openapi.MultipartRelatedRender{Data: response})
			return
		case nas.MsgTypePDUSessionModificationComplete:
			smContext.StopT3591()
		case nas.MsgTypePDUSessionModificationReject:
			smContext.StopT3591()
		}
	}

	tunnel := smContext.Tunnel
	dcTunnel := smContext.DCTunnel
	pdrList := []*smf_context.PDR{}
	farList := []*smf_context.FAR{}
	barList := []*smf_context.BAR{}
	qerList := []*smf_context.QER{}
	urrList := []*smf_context.URR{}

	switch smContextUpdateData.UpCnxState {
	case models.UpCnxState_ACTIVATING:
		smContext.CheckState(smf_context.Active)
		// Wait till the state becomes Active again
		// TODO: implement sleep wait in concurrent architecture

		smContext.SetState(smf_context.ModificationPending)
		response.JsonData.N2SmInfo = &models.RefToBinaryData{ContentId: "PDUSessionResourceSetupRequestTransfer"}
		response.JsonData.UpCnxState = models.UpCnxState_ACTIVATING
		response.JsonData.N2SmInfoType = models.N2SmInfoType_PDU_RES_SETUP_REQ

		n2Buf, err = smf_context.BuildPDUSessionResourceSetupRequestTransfer(smContext)
		if err != nil {
			logger.PduSessLog.Errorf("Build PDUSession Resource Setup Request Transfer Error(%s)", err.Error())
		} else {
			response.JsonData.N2SmInfoType = models.N2SmInfoType_PDU_RES_SETUP_REQ
			response.BinaryDataN2SmInformation = n2Buf
			response.JsonData.N2SmInfo = &models.RefToBinaryData{ContentId: "PDUSessionResourceSetupRequestTransfer"}
		}
		smContext.UpCnxState = models.UpCnxState_ACTIVATING
	case models.UpCnxState_DEACTIVATED:
		smContext.CheckState(smf_context.Active)
		// Wait till the state becomes Active again
		// TODO: implement sleep wait in concurrent architecture

		// If the PDU session has been released, skip sending PFCP Session Modification Request
		if smContext.CheckState(smf_context.InActivePending) {
			logger.CtxLog.Infof("Skip sending PFCP Session Modification Request of PDUSessionID:%d of SUPI:%s",
				smContext.PDUSessionID, smContext.Supi)
			response.JsonData.UpCnxState = models.UpCnxState_DEACTIVATED
			c.Render(http.StatusOK, openapi.MultipartRelatedRender{Data: response})
			return
		}
		smContext.SetState(smf_context.ModificationPending)
		response.JsonData.UpCnxState = models.UpCnxState_DEACTIVATED
		smContext.UpCnxState = body.JsonData.UpCnxState
		// UE location change is a charging event
		// TODO: This is not tested yet
		if !reflect.DeepEqual(smContext.UeLocation, body.JsonData.UeLocation) {
			// All rating group related to this Pdu session should send charging request
			for _, dataPath := range tunnel.DataPathPool {
				if dataPath.Activated {
					for curDataPathNode := dataPath.FirstDPNode; curDataPathNode != nil; curDataPathNode = curDataPathNode.Next() {
						if curDataPathNode.IsANUPF() {
							urrList = append(urrList, curDataPathNode.UpLinkTunnel.PDR.URR...)
							QueryReport(smContext, curDataPathNode.UPF, urrList, models.ChfConvergedChargingTriggerType_USER_LOCATION_CHANGE)
						}
					}
				}
			}

			p.ReportUsageAndUpdateQuota(smContext)
		}

		smContext.UeLocation = body.JsonData.UeLocation

		// Set FAR and An, N3 Release Info
		// TODO: Deactivate all datapath in ANUPF
		farList = []*smf_context.FAR{}
		for _, dataPath := range smContext.Tunnel.DataPathPool {
			ANUPF := dataPath.FirstDPNode
			DLPDR := ANUPF.DownLinkTunnel.PDR
			if DLPDR == nil {
				smContext.Log.Warnf("Access network resource is released")
			} else {
				DLPDR.FAR.State = smf_context.RULE_UPDATE
				DLPDR.FAR.ApplyAction.Forw = false
				DLPDR.FAR.ApplyAction.Buff = true
				DLPDR.FAR.ApplyAction.Nocp = true
				farList = append(farList, DLPDR.FAR)
				sendPFCPModification = true
				smContext.SetState(smf_context.PFCPModification)
			}
		}
	}

	switch smContextUpdateData.N2SmInfoType {
	case models.N2SmInfoType_PDU_RES_SETUP_RSP:
		smContext.CheckState(smf_context.Active)
		// Wait till the state becomes Active again
		// TODO: implement sleep wait in concurrent architecture

		smContext.SetState(smf_context.ModificationPending)
		pdrList = []*smf_context.PDR{}
		farList = []*smf_context.FAR{}

		for _, dataPath := range tunnel.DataPathPool {
			if dataPath.Activated {
				ANUPF := dataPath.FirstDPNode
				DLPDR := ANUPF.DownLinkTunnel.PDR

				DLPDR.FAR.ApplyAction = pfcpType.ApplyAction{
					Buff: false,
					Drop: false,
					Dupl: false,
					Forw: true,
					Nocp: false,
				}
				DLPDR.FAR.ForwardingParameters = &smf_context.ForwardingParameters{
					DestinationInterface: pfcpType.DestinationInterface{
						InterfaceValue: pfcpType.DestinationInterfaceAccess,
					},
					NetworkInstance: &pfcpType.NetworkInstance{
						NetworkInstance: smContext.Dnn,
						FQDNEncoding:    factory.SmfConfig.Configuration.NwInstFqdnEncoding,
					},
				}

				DLPDR.State = smf_context.RULE_UPDATE
				DLPDR.FAR.State = smf_context.RULE_UPDATE

				pdrList = append(pdrList, DLPDR)
				farList = append(farList, DLPDR.FAR)
			}
		}

		if err = smf_context.
			HandlePDUSessionResourceSetupResponseTransfer(body.BinaryDataN2SmInformation, smContext); err != nil {
			smContext.Log.Errorf("Handle PDUSessionResourceSetupResponseTransfer failed: %+v", err)
		} else if smContext.NrdcIndicator {
			for _, pdr := range pdrList {
				// Remove all PDRs except the default PDR
				if pdr.Precedence != 255 {
					pdr.State = smf_context.RULE_REMOVE
				}
			}
			if err = smContext.ApplyDcPccRulesOnDcTunnel(); err != nil {
				smContext.Log.Errorf("ApplyDcPccRulesOnDcTunnel failed: %+v", err)
			}
			for _, dataPath := range dcTunnel.DataPathPool {
				if dataPath.Activated {
					ANUPF := dataPath.FirstDPNode
					ULPDR := ANUPF.UpLinkTunnel.PDR
					DLPDR := ANUPF.DownLinkTunnel.PDR

					ULPDR.FAR.ApplyAction = pfcpType.ApplyAction{
						Buff: false,
						Drop: false,
						Dupl: false,
						Forw: true,
						Nocp: false,
					}
					DLPDR.FAR.ApplyAction = pfcpType.ApplyAction{
						Buff: false,
						Drop: false,
						Dupl: false,
						Forw: true,
						Nocp: false,
					}

					// ULPDR.FAR.ForwardingParameters = &smf_context.ForwardingParameters{
					// 	DestinationInterface: pfcpType.DestinationInterface{
					// 		InterfaceValue: pfcpType.DestinationInterfaceCore,
					// 	},
					// 	NetworkInstance: &pfcpType.NetworkInstance{
					// 		NetworkInstance: smContext.Dnn,
					// 		FQDNEncoding:    factory.SmfConfig.Configuration.NwInstFqdnEncoding,
					// 	},
					// }
					// DLPDR.FAR.ForwardingParameters = &smf_context.ForwardingParameters{
					// 	DestinationInterface: pfcpType.DestinationInterface{
					// 		InterfaceValue: pfcpType.DestinationInterfaceAccess,
					// 	},
					// 	NetworkInstance: &pfcpType.NetworkInstance{
					// 		NetworkInstance: smContext.Dnn,
					// 		FQDNEncoding:    factory.SmfConfig.Configuration.NwInstFqdnEncoding,
					// 	},
					// }

					DLPDR.State = smf_context.RULE_INITIAL
					DLPDR.FAR.State = smf_context.RULE_INITIAL
					ULPDR.State = smf_context.RULE_INITIAL
					ULPDR.FAR.State = smf_context.RULE_INITIAL

					pdrList = append(pdrList, DLPDR)
					farList = append(farList, DLPDR.FAR)

					pdrList = append(pdrList, ULPDR)
					farList = append(farList, ULPDR.FAR)
				}
			}
		}
		sendPFCPModification = true
		smContext.SetState(smf_context.PFCPModification)
	case models.N2SmInfoType_PDU_RES_SETUP_FAIL:
		if err = smf_context.
			HandlePDUSessionResourceSetupUnsuccessfulTransfer(body.BinaryDataN2SmInformation, smContext); err != nil {
			smContext.Log.Errorf("Handle PDUSessionResourceSetupResponseTransfer failed: %+v", err)
		}
	case models.N2SmInfoType_PDU_RES_MOD_RSP:
		if err = smf_context.
			HandlePDUSessionResourceModifyResponseTransfer(body.BinaryDataN2SmInformation, smContext); err != nil {
			smContext.Log.Errorf("Handle PDUSessionResourceModifyResponseTransfer failed: %+v", err)
		}
	case models.N2SmInfoType_PDU_RES_REL_RSP:
		// remove an tunnel info
		smContext.Log.Infoln("Handle N2 PDU Resource Release Response")
		smContext.Tunnel.ANInformation = struct {
			IPAddress net.IP
			TEID      uint32
		}{nil, 0}

		if smContext.PDUSessionRelease_DUE_TO_DUP_PDU_ID {
			smContext.CheckState(smf_context.InActivePending)
			// Wait till the state becomes Active again
			// TODO: implement sleep wait in concurrent architecture
			smContext.Log.Infoln("Release_DUE_TO_DUP_PDU_ID: Send Update SmContext Response")
			response.JsonData.UpCnxState = models.UpCnxState_DEACTIVATED
			// If NAS layer is inActive, the context should be remove
			if smContext.CheckState(smf_context.InActive) {
				p.RemoveSMContextFromAllNF(smContext, true)
			}
		} else if smContext.CheckState(smf_context.InActive) { // normal case
			// Wait till the state becomes Active again
			// TODO: implement sleep wait in concurrent architecture

			// If N1 PDU Session Release Complete is received, smContext state is InActive.
			// Remove SMContext when receiving N2 PDU Resource Release Response.
			// Use go routine to send Notification to prevent blocking the handling process
			p.RemoveSMContextFromAllNF(smContext, true)
		}
	case models.N2SmInfoType_PATH_SWITCH_REQ:
		smContext.Log.Traceln("Handle Path Switch Request")
		smContext.CheckState(smf_context.Active)
		// Wait till the state becomes Active again
		// TODO: implement sleep wait in concurrent architecture

		smContext.SetState(smf_context.ModificationPending)

		if err = smf_context.HandlePathSwitchRequestTransfer(body.BinaryDataN2SmInformation, smContext); err != nil {
			smContext.Log.Errorf("Handle PathSwitchRequestTransfer: %+v", err)
		}

		if n2Buf, err = smf_context.BuildPathSwitchRequestAcknowledgeTransfer(smContext); err != nil {
			smContext.Log.Errorf("Build Path Switch Transfer Error(%+v)", err)
		} else {
			response.JsonData.N2SmInfoType = models.N2SmInfoType_PATH_SWITCH_REQ_ACK
			response.BinaryDataN2SmInformation = n2Buf
			response.JsonData.N2SmInfo = &models.RefToBinaryData{
				ContentId: "PATH_SWITCH_REQ_ACK",
			}
		}

		for _, dataPath := range tunnel.DataPathPool {
			if dataPath.Activated {
				ANUPF := dataPath.FirstDPNode
				DLPDR := ANUPF.DownLinkTunnel.PDR

				pdrList = append(pdrList, DLPDR)
				farList = append(farList, DLPDR.FAR)
			}
		}

		sendPFCPModification = true
		smContext.SetState(smf_context.PFCPModification)
	case models.N2SmInfoType_PATH_SWITCH_SETUP_FAIL:
		smContext.CheckState(smf_context.Active)
		// Wait till the state becomes Active again
		// TODO: implement sleep wait in concurrent architecture

		smContext.SetState(smf_context.ModificationPending)
		err = smf_context.HandlePathSwitchRequestSetupFailedTransfer(
			body.BinaryDataN2SmInformation, smContext)
		if err != nil {
			smContext.Log.Errorf("HandlePathSwitchRequestSetupFailedTransfer failed: %v", err)
		}
	case models.N2SmInfoType_HANDOVER_REQUIRED:
		smContext.CheckState(smf_context.Active)
		// Wait till the state becomes Active again
		// TODO: implement sleep wait in concurrent architecture
		smContext.SetState(smf_context.ModificationPending)
		response.JsonData.N2SmInfo = &models.RefToBinaryData{ContentId: "Handover"}
	}

	switch smContextUpdateData.HoState {
	case models.HoState_PREPARING:
		smContext.Log.Traceln("In HoState_PREPARING")
		smContext.CheckState(smf_context.Active)
		// Wait till the state becomes Active again
		// TODO: implement sleep wait in concurrent architecture

		smContext.SetState(smf_context.ModificationPending)
		smContext.HoState = models.HoState_PREPARING
		err = smf_context.HandleHandoverRequiredTransfer(
			body.BinaryDataN2SmInformation, smContext)
		if err != nil {
			smContext.Log.Errorf("Handle HandoverRequiredTransfer failed: %+v", err)
		}
		response.JsonData.N2SmInfoType = models.N2SmInfoType_PDU_RES_SETUP_REQ

		if n2Buf, err = smf_context.BuildPDUSessionResourceSetupRequestTransfer(smContext); err != nil {
			smContext.Log.Errorf("Build PDUSession Resource Setup Request Transfer Error(%s)", err.Error())
		} else {
			response.BinaryDataN2SmInformation = n2Buf
			response.JsonData.N2SmInfoType = models.N2SmInfoType_PDU_RES_SETUP_REQ
			response.JsonData.N2SmInfo = &models.RefToBinaryData{
				ContentId: "PDU_RES_SETUP_REQ",
			}
		}
		response.JsonData.HoState = models.HoState_PREPARING
	case models.HoState_PREPARED:
		smContext.Log.Traceln("In HoState_PREPARED")
		smContext.CheckState(smf_context.Active)
		// Wait till the state becomes Active again
		// TODO: implement sleep wait in concurrent architecture

		smContext.SetState(smf_context.ModificationPending)
		smContext.HoState = models.HoState_PREPARED
		response.JsonData.HoState = models.HoState_PREPARED
		err = smf_context.HandleHandoverRequestAcknowledgeTransfer(
			body.BinaryDataN2SmInformation, smContext)
		if err != nil {
			smContext.Log.Errorf("Handle HandoverRequestAcknowledgeTransfer failed: %+v", err)
		}

		// request UPF establish indirect forwarding path for DL
		if smContext.DLForwardingType == smf_context.IndirectForwarding {
			ANUPF := smContext.IndirectForwardingTunnel.FirstDPNode
			IndirectForwardingPDR := smContext.IndirectForwardingTunnel.FirstDPNode.UpLinkTunnel.PDR

			pdrList = append(pdrList, IndirectForwardingPDR)
			farList = append(farList, IndirectForwardingPDR.FAR)

			// release indirect forwading path
			if err = ANUPF.UPF.RemovePDR(IndirectForwardingPDR); err != nil {
				logger.PduSessLog.Errorln("release indirect path: ", err)
			}

			sendPFCPModification = true
			smContext.SetState(smf_context.PFCPModification)
		}

		if n2Buf, err = smf_context.BuildHandoverCommandTransfer(smContext); err != nil {
			smContext.Log.Errorf("Build HandoverCommandTransfer failed: %v", err)
		} else {
			response.BinaryDataN2SmInformation = n2Buf
			response.JsonData.N2SmInfoType = models.N2SmInfoType_HANDOVER_CMD
			response.JsonData.N2SmInfo = &models.RefToBinaryData{
				ContentId: "HANDOVER_CMD",
			}
		}
		response.JsonData.HoState = models.HoState_PREPARING
	case models.HoState_COMPLETED:
		smContext.Log.Traceln("In HoState_COMPLETED")
		smContext.CheckState(smf_context.Active)
		// Wait till the state becomes Active again
		// TODO: implement sleep wait in concurrent architecture

		for _, dataPath := range tunnel.DataPathPool {
			if dataPath.Activated {
				ANUPF := dataPath.FirstDPNode
				DLPDR := ANUPF.DownLinkTunnel.PDR

				pdrList = append(pdrList, DLPDR)
				farList = append(farList, DLPDR.FAR)
			}
		}

		// remove indirect forwarding path
		if smContext.DLForwardingType == smf_context.IndirectForwarding {
			indirectForwardingPDR := smContext.IndirectForwardingTunnel.FirstDPNode.GetUpLinkPDR()
			indirectForwardingPDR.State = smf_context.RULE_REMOVE
			indirectForwardingPDR.FAR.State = smf_context.RULE_REMOVE
			pdrList = append(pdrList, indirectForwardingPDR)
			farList = append(farList, indirectForwardingPDR.FAR)
		}

		sendPFCPModification = true
		smContext.SetState(smf_context.PFCPModification)
		smContext.HoState = models.HoState_COMPLETED
		response.JsonData.HoState = models.HoState_COMPLETED
	}

	if smContextUpdateData.Cause == models.SmfPduSessionCause_REL_DUE_TO_DUPLICATE_SESSION_ID {
		// * release PDU Session Here
		smContext.Log.Infoln("[SMF] Cause_REL_DUE_TO_DUPLICATE_SESSION_ID")
		if smContext.CheckState(smf_context.Active) {
			// Wait till the state becomes Active again
			// TODO: implement sleep wait in concurrent architecture
			logger.PduSessLog.Warnf("SMContext[%s-%02d] should be Active, but actual %s",
				smContext.Supi, smContext.PDUSessionID, smContext.State().String())
		}

		smContext.PDUSessionRelease_DUE_TO_DUP_PDU_ID = true

		switch smContext.State() {
		case smf_context.ActivePending, smf_context.ModificationPending, smf_context.Active:
			if buf, err = smf_context.BuildPDUSessionResourceReleaseCommandTransfer(smContext); err != nil {
				smContext.Log.Errorf("Build PDUSessionResourceReleaseCommandTransfer failed: %v", err)
			} else {
				response.BinaryDataN2SmInformation = buf
				response.JsonData.N2SmInfoType = models.N2SmInfoType_PDU_RES_REL_CMD
				response.JsonData.N2SmInfo = &models.RefToBinaryData{
					ContentId: "PDUResourceReleaseCommand",
				}
			}

			pfcpResponseStatus = releaseSession(smContext)
		default:
			smContext.Log.Infof("Not needs to send pfcp release")
		}

		smContext.Log.Infoln("[SMF] Cause_REL_DUE_TO_DUPLICATE_SESSION_ID")
	}

	// Check FSM and take corresponding action
	switch smContext.State() {
	case smf_context.PFCPModification:
		smContext.Log.Traceln("In case PFCPModification")

		if sendPFCPModification {
			pfcpResponseStatus = p.updateAnUpfPfcpSession(smContext, pdrList, farList, barList, qerList, urrList)
		}

		switch pfcpResponseStatus {
		case smf_context.SessionUpdateSuccess:
			smContext.Log.Traceln("In case SessionUpdateSuccess")
			smContext.SetState(smf_context.Active)
			c.Render(http.StatusOK, openapi.MultipartRelatedRender{Data: response})
		case smf_context.SessionUpdateFailed:
			smContext.Log.Traceln("In case SessionUpdateFailed")
			smContext.SetState(smf_context.Active)
			// It is just a template
			updateSmContextError := models.UpdateSmContextResponse400{
				JsonData: &models.SmContextUpdateError{
					Error: &smf_errors.N1SmError,
				},
			} // Depends on the reason why N4 fail
			c.JSON(http.StatusForbidden, updateSmContextError)

		case smf_context.SessionReleaseSuccess:
			p.ReleaseChargingSession(smContext)

			smContext.Log.Traceln("In case SessionReleaseSuccess")
			smContext.SetState(smf_context.InActivePending)
			c.Render(http.StatusOK, openapi.MultipartRelatedRender{Data: response})

		case smf_context.SessionReleaseFailed:
			// Update SmContext Request(N1 PDU Session Release Request)
			// Send PDU Session Release Reject
			smContext.Log.Traceln("In case SessionReleaseFailed")
			problemDetail := models.SmfPduSessionExtProblemDetails{
				Status: http.StatusInternalServerError,
				Cause:  "SYSTEM_FAILURE",
			}
			smContext.SetState(smf_context.Active)
			errResponse := models.UpdateSmContextResponse400{
				JsonData: &models.SmContextUpdateError{
					Error: &problemDetail,
				},
			}
			if smContextUpdateData.Cause != models.SmfPduSessionCause_REL_DUE_TO_DUPLICATE_SESSION_ID {
				if buf, err = smf_context.BuildGSMPDUSessionReleaseReject(smContext); err != nil {
					logger.PduSessLog.Errorf("build GSM PDUSessionReleaseReject failed: %+v", err)
				} else {
					errResponse.BinaryDataN1SmMessage = buf
					errResponse.JsonData.N1SmMsg = &models.RefToBinaryData{ContentId: "PDUSessionReleaseReject"}
				}
			}
			c.JSON(int(problemDetail.Status), errResponse)
		}
		smContext.PostRemoveDataPath()
	case smf_context.ModificationPending:
		smContext.Log.Traceln("In case ModificationPending")
		smContext.SetState(smf_context.Active)
		c.Render(http.StatusOK, openapi.MultipartRelatedRender{Data: response})
	case smf_context.InActive, smf_context.InActivePending:
		smContext.Log.Traceln("In case InActive, InActivePending")
		c.Render(http.StatusOK, openapi.MultipartRelatedRender{Data: response})
	default:
		c.Render(http.StatusOK, openapi.MultipartRelatedRender{Data: response})
	}

	if smContext.PDUSessionRelease_DUE_TO_DUP_PDU_ID {
		// Note:
		// We don't want to launch timer to wait for N2SmInfoType_PDU_RES_REL_RSP.
		// So, local release smCtx and notify AMF after sending PDUSessionResourceReleaseCommand
		p.RemoveSMContextFromAllNF(smContext, true)
	}
}

func (p *Processor) HandlePDUSessionSMContextRelease(
	c *gin.Context,
	body models.ReleaseSmContextRequest,
	smContextRef string,
) {
	logger.PduSessLog.Infoln("In HandlePDUSessionSMContextRelease")
	smContext := smf_context.GetSMContextByRef(smContextRef)

	if smContext == nil {
		logger.PduSessLog.Warnf("SMContext[%s] is not found", smContextRef)

		updateSmContextError := &models.UpdateSmContextResponse400{
			JsonData: &models.SmContextUpdateError{
				UpCnxState: models.UpCnxState_DEACTIVATED,
				Error: &models.SmfPduSessionExtProblemDetails{
					Type:   "Resource Not Found",
					Title:  "SMContext Ref is not found",
					Status: http.StatusNotFound,
				},
			},
		}
		c.JSON(http.StatusNotFound, updateSmContextError)
		return
	}

	smContext.SMLock.Lock()
	defer smContext.SMLock.Unlock()

	smContext.StopT3591()
	smContext.StopT3592()

	// remove SM Policy Association
	if smContext.SMPolicyID != "" {
		if err := p.Consumer().SendSMPolicyAssociationTermination(smContext); err != nil {
			smContext.Log.Errorf("SM Policy Termination failed: %s", err)
		} else {
			smContext.SMPolicyID = ""
		}
	}

	if p.Context().Ues.UeExists(smContext.Supi) {
		problemDetails, err := p.Consumer().UnSubscribe(smContext)
		if problemDetails != nil {
			logger.PduSessLog.Errorf("SDM UnSubscription Failed Problem[%+v]", problemDetails)
		} else if err != nil {
			logger.PduSessLog.Errorf("SDM UnSubscription Error[%+v]", err)
		}
	}

	if smContext.UeCmRegistered {
		problemDetails, err := p.Consumer().UeCmDeregistration(smContext)
		if problemDetails != nil {
			if problemDetails.Cause != CONTEXT_NOT_FOUND {
				logger.PduSessLog.Errorf("UECM_DeRegistration Failed Problem[%+v]", problemDetails)
			}
		} else if err != nil {
			logger.PduSessLog.Errorf("UECM_DeRegistration Error[%+v]", err)
		} else {
			logger.PduSessLog.Traceln("UECM_DeRegistration successful")
		}
	}

	if !smContext.CheckState(smf_context.InActive) {
		smContext.SetState(smf_context.PFCPModification)
	}
	pfcpResponseStatus := releaseSession(smContext)

	switch pfcpResponseStatus {
	case smf_context.SessionReleaseSuccess:
		p.ReleaseChargingSession(smContext)

		smContext.Log.Traceln("In case SessionReleaseSuccess")
		smContext.SetState(smf_context.InActive)
		c.Status(http.StatusNoContent)

	case smf_context.SessionReleaseFailed:
		// Update SmContext Request(N1 PDU Session Release Request)
		// Send PDU Session Release Reject
		smContext.Log.Traceln("In case SessionReleaseFailed")
		problemDetail := models.SmfPduSessionExtProblemDetails{
			Status: http.StatusInternalServerError,
			Cause:  "SYSTEM_FAILURE",
		}
		smContext.SetState(smf_context.Active)
		errResponse := models.UpdateSmContextResponse400{
			JsonData: &models.SmContextUpdateError{
				Error: &problemDetail,
			},
		}
		if buf, err := smf_context.BuildGSMPDUSessionReleaseReject(smContext); err != nil {
			smContext.Log.Errorf("Build GSM PDUSessionReleaseReject failed: %+v", err)
		} else {
			errResponse.BinaryDataN1SmMessage = buf
			errResponse.JsonData.N1SmMsg = &models.RefToBinaryData{ContentId: "PDUSessionReleaseReject"}
		}

		c.JSON(int(problemDetail.Status), errResponse)

	default:
		smContext.Log.Warnf("The state shouldn't be [%s]\n", pfcpResponseStatus)

		smContext.Log.Traceln("In case Unknown")
		problemDetail := models.SmfPduSessionExtProblemDetails{
			Status: http.StatusInternalServerError,
			Cause:  "SYSTEM_FAILURE",
		}
		smContext.SetState(smf_context.Active)
		errResponse := models.UpdateSmContextResponse400{
			JsonData: &models.SmContextUpdateError{
				Error: &problemDetail,
			},
		}
		if buf, err := smf_context.BuildGSMPDUSessionReleaseReject(smContext); err != nil {
			smContext.Log.Errorf("Build GSM PDUSessionReleaseReject failed: %+v", err)
		} else {
			errResponse.BinaryDataN1SmMessage = buf
			errResponse.JsonData.N1SmMsg = &models.RefToBinaryData{ContentId: "PDUSessionReleaseReject"}
		}

		c.JSON(int(problemDetail.Status), errResponse)
	}

	p.RemoveSMContextFromAllNF(smContext, false)
}

func (p *Processor) HandlePDUSessionSMContextLocalRelease(
	smContext *smf_context.SMContext, createData *models.SmfPduSessionSmContextCreateData,
) {
	smContext.SMLock.Lock()
	defer smContext.SMLock.Unlock()

	// remove SM Policy Association
	if smContext.SMPolicyID != "" {
		if err := p.Consumer().SendSMPolicyAssociationTermination(smContext); err != nil {
			logger.PduSessLog.Errorf("SM Policy Termination failed: %s", err)
		} else {
			smContext.SMPolicyID = ""
		}
	}

	if p.Context().Ues.UeExists(smContext.Supi) {
		problemDetails, err := p.Consumer().UnSubscribe(smContext)
		if problemDetails != nil {
			logger.PduSessLog.Errorf("SDM UnSubscription Failed Problem[%+v]", problemDetails)
		} else if err != nil {
			logger.PduSessLog.Errorf("SDM UnSubscription Error[%+v]", err)
		}
	}

	if smContext.UeCmRegistered {
		problemDetails, err := p.Consumer().UeCmDeregistration(smContext)
		if problemDetails != nil {
			if problemDetails.Cause != CONTEXT_NOT_FOUND {
				logger.PduSessLog.Errorf("UECM_DeRegistration Failed Problem[%+v]", problemDetails)
			}
		} else if err != nil {
			logger.PduSessLog.Errorf("UECM_DeRegistration Error[%+v]", err)
		} else {
			logger.PduSessLog.Traceln("UECM_DeRegistration successful")
		}
	}

	smContext.SetState(smf_context.PFCPModification)

	pfcpResponseStatus := releaseSession(smContext)

	switch pfcpResponseStatus {
	case smf_context.SessionReleaseSuccess:
		p.ReleaseChargingSession(smContext)

		logger.CtxLog.Traceln("In case SessionReleaseSuccess")
		smContext.SetState(smf_context.InActivePending)
		if createData.SmContextStatusUri != smContext.SmStatusNotifyUri {
			problemDetails, err := p.Consumer().SendSMContextStatusNotification(smContext.SmStatusNotifyUri)
			if problemDetails != nil || err != nil {
				if problemDetails != nil {
					logger.PduSessLog.Warnf("Send SMContext Status Notification Problem[%+v]", problemDetails)
				}

				if err != nil {
					logger.PduSessLog.Warnf("Send SMContext Status Notification Error[%v]", err)
				}
			} else {
				logger.PduSessLog.Traceln("Send SMContext Status Notification successfully")
			}
		}
		p.RemoveSMContextFromAllNF(smContext, false)

	case smf_context.SessionReleaseFailed:
		logger.CtxLog.Traceln("In case SessionReleaseFailed")
		smContext.SetState(smf_context.Active)

	default:
		smContext.Log.Warnf("The state shouldn't be [%s]", pfcpResponseStatus)
		logger.CtxLog.Traceln("In case Unknown")
		smContext.SetState(smf_context.Active)
	}
}

func releaseSession(smContext *smf_context.SMContext) smf_context.PFCPSessionResponseStatus {
	smContext.SetState(smf_context.PFCPModification)

	for _, res := range ReleaseTunnel(smContext) {
		if res.Status != smf_context.SessionReleaseSuccess {
			return res.Status
		}
	}
	if !smContext.NrdcIndicator {
		return smf_context.SessionReleaseSuccess
	}

	for _, res := range ReleaseDcTunnel(smContext) {
		if res.Status != smf_context.SessionReleaseSuccess {
			return res.Status
		}
	}
	return smf_context.SessionReleaseSuccess
}

func (p *Processor) makeEstRejectResAndReleaseSMContext(
	c *gin.Context,
	smContext *smf_context.SMContext,
	nasErrorCause uint8,
	sbiError *models.SmfPduSessionExtProblemDetails,
) {
	postSmContextsError := models.PostSmContextsError{
		JsonData: &models.SmContextCreateError{
			Error:   sbiError,
			N1SmMsg: &models.RefToBinaryData{ContentId: "n1SmMsg"},
		},
	}
	if buf, err := smf_context.
		BuildGSMPDUSessionEstablishmentReject(
			smContext,
			nasErrorCause); err != nil {
		logger.PduSessLog.Errorf("Build GSM PDUSessionEstablishmentReject failed: %+v", err)
	} else {
		postSmContextsError.BinaryDataN1SmMessage = buf
	}
	p.nasErrorResponse(c, int(sbiError.Status), postSmContextsError)
	p.RemoveSMContextFromAllNF(smContext, false)
}

func (p *Processor) sendGSMPDUSessionReleaseCommand(smContext *smf_context.SMContext, nasPdu []byte) {
	n1n2Request := models.N1N2MessageTransferRequest{}
	n1n2Request.JsonData = &models.N1N2MessageTransferReqData{
		PduSessionId: smContext.PDUSessionID,
		N1MessageContainer: &models.N1MessageContainer{
			N1MessageClass:   "SM",
			N1MessageContent: &models.RefToBinaryData{ContentId: "GSM_NAS"},
		},
	}
	n1n2Request.BinaryDataN1Message = nasPdu
	if smContext.T3592 != nil {
		smContext.T3592.Stop()
		smContext.T3592 = nil
	}

	// Start T3592
	t3592 := factory.SmfConfig.Configuration.T3592
	if t3592.Enable {
		ctx, _, err := smf_context.GetSelf().GetTokenCtx(models.ServiceName_NAMF_COMM, models.NrfNfManagementNfType_AMF)
		if err != nil {
			smContext.Log.Warnf("Get namf-comm token failed: %+v", err)
			return
		}

		smContext.T3592 = smf_context.NewTimer(t3592.ExpireTime, t3592.MaxRetryTimes, func(expireTimes int32) {
			smContext.SMLock.Lock()
			rspData, errMsgTransfer := p.Consumer().
				N1N2MessageTransfer(ctx, smContext.Supi, n1n2Request, smContext.CommunicationClientApiPrefix)
			if errMsgTransfer != nil {
				logger.ConsumerLog.Warnf("N1N2MessageTransfer for GSMPDUSessionReleaseCommand failed: %+v", errMsgTransfer)
				return
			}

			if rspData.Cause == models.N1N2MessageTransferCause_N1_MSG_NOT_TRANSFERRED {
				smContext.Log.Warnf("%v", rspData.Cause)
			}

			smContext.SMLock.Unlock()
		}, func() {
			smContext.Log.Warn("T3592 Expires 3 times, abort notification procedure")
			smContext.SMLock.Lock()
			smContext.T3592 = nil
			p.SendReleaseNotification(smContext)
			smContext.SMLock.Unlock()
		})
	}
}

func (p *Processor) sendGSMPDUSessionModificationCommand(smContext *smf_context.SMContext, nasPdu []byte) {
	n1n2Request := models.N1N2MessageTransferRequest{}
	n1n2Request.JsonData = &models.N1N2MessageTransferReqData{
		PduSessionId: smContext.PDUSessionID,
		N1MessageContainer: &models.N1MessageContainer{
			N1MessageClass:   "SM",
			N1MessageContent: &models.RefToBinaryData{ContentId: "GSM_NAS"},
		},
	}
	n1n2Request.BinaryDataN1Message = nasPdu

	if smContext.T3591 != nil {
		smContext.T3591.Stop()
		smContext.T3591 = nil
	}

	// Start T3591
	t3591 := factory.SmfConfig.Configuration.T3591
	if t3591.Enable {
		ctx, _, err := smf_context.GetSelf().GetTokenCtx(models.ServiceName_NAMF_COMM, models.NrfNfManagementNfType_AMF)
		if err != nil {
			smContext.Log.Warnf("Get namf-comm token failed: %+v", err)
			return
		}

		smContext.T3591 = smf_context.NewTimer(t3591.ExpireTime, t3591.MaxRetryTimes, func(expireTimes int32) {
			smContext.SMLock.Lock()
			defer smContext.SMLock.Unlock()
			rspData, errMsgTransfer := p.Consumer().
				N1N2MessageTransfer(ctx, smContext.Supi, n1n2Request, smContext.CommunicationClientApiPrefix)
			if errMsgTransfer != nil {
				logger.ConsumerLog.Warnf("N1N2MessageTransfer for GSMPDUSessionModificationCommand failed: %+v", errMsgTransfer)
				return
			}

			if rspData.Cause == models.N1N2MessageTransferCause_N1_MSG_NOT_TRANSFERRED {
				smContext.Log.Warnf("%v", rspData.Cause)
			}
		}, func() {
			smContext.Log.Warn("T3591 Expires3 times, abort notification procedure")
			smContext.SMLock.Lock()
			defer smContext.SMLock.Unlock()
			smContext.T3591 = nil
		})
	}
}

func (p *Processor) nasErrorResponse(
	c *gin.Context,
	status int,
	errBody models.PostSmContextsError,
) {
	switch status {
	case http.StatusForbidden,
		http.StatusNotFound,
		http.StatusInternalServerError,
		http.StatusGatewayTimeout:
		logger.SBILog.Warnf("NAS Error Response: %v", errBody)
		if errBody.BinaryDataN1SmMessage != nil ||
			errBody.BinaryDataN2SmMessage != nil {
			rspBody, contentType, err := openapi.MultipartSerialize(errBody)
			if err != nil {
				logger.SBILog.Infof("MultipartSerialize error: %v", err)
				c.JSON(http.StatusInternalServerError, openapi.ProblemDetailsSystemFailure(err.Error()))
			} else {
				c.Data(status, contentType, rspBody)
			}
			return
		}
	}
	c.JSON(status, errBody)
}
