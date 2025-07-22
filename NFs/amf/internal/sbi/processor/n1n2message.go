package processor

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/free5gc/amf/internal/context"
	gmm_message "github.com/free5gc/amf/internal/gmm/message"
	"github.com/free5gc/amf/internal/logger"
	ngap_message "github.com/free5gc/amf/internal/ngap/message"
	callback "github.com/free5gc/amf/internal/sbi/processor/notifier"
	"github.com/free5gc/aper"
	"github.com/free5gc/nas/nasMessage"
	"github.com/free5gc/ngap/ngapType"
	"github.com/free5gc/openapi/models"
)

// TS23502 4.2.3.3, 4.2.4.3, 4.3.2.2, 4.3.2.3, 4.3.3.2, 4.3.7
func (p *Processor) HandleN1N2MessageTransferRequest(c *gin.Context,
	n1n2MessageTransferRequest models.N1N2MessageTransferRequest,
) {
	logger.ProducerLog.Infof("Handle N1N2 Message Transfer Request")

	ueContextID := c.Param("ueContextId")
	reqUri := c.Param("reqUri")

	n1n2MessageTransferRspData, locationHeader, problemDetails, transferErr := p.N1N2MessageTransferProcedure(
		ueContextID, reqUri, n1n2MessageTransferRequest)

	if problemDetails != nil {
		c.JSON(int(problemDetails.Status), problemDetails)
		return
	} else if transferErr != nil {
		c.JSON(int(transferErr.Error.Status), transferErr)
		return
	} else if n1n2MessageTransferRspData != nil {
		switch n1n2MessageTransferRspData.Cause {
		case models.N1N2MessageTransferCause_N1_MSG_NOT_TRANSFERRED:
			fallthrough
		case models.N1N2MessageTransferCause_N1_N2_TRANSFER_INITIATED:
			c.JSON(http.StatusOK, n1n2MessageTransferRspData)
			return
		case models.N1N2MessageTransferCause_ATTEMPTING_TO_REACH_UE:
			headers := http.Header{
				"Location": {locationHeader},
			}
			c.JSON(http.StatusAccepted, gin.H{"headers": headers, "data": n1n2MessageTransferRspData})
			return
		}
	}

	problemDetails = &models.ProblemDetails{
		Status: http.StatusForbidden,
		Cause:  "UNSPECIFIED",
	}
	c.JSON(http.StatusForbidden, problemDetails)
}

// There are 4 possible return value for this function:
//   - n1n2MessageTransferRspData: if AMF handle N1N2MessageTransfer Request successfully.
//   - locationHeader: if response status code is 202, then it will return a non-empty string location header for
//     response
//   - problemDetails: if AMF reject the request due to application error, e.g. UE context not found.
//   - TransferErr: if AMF reject the request due to procedure error, e.g. UE has an ongoing procedure.
//
// see TS 29.518 6.1.3.5.3.1 for more details.
func (p *Processor) N1N2MessageTransferProcedure(ueContextID string, reqUri string,
	n1n2MessageTransferRequest models.N1N2MessageTransferRequest) (
	n1n2MessageTransferRspData *models.N1N2MessageTransferRspData,
	locationHeader string, problemDetails *models.ProblemDetails,
	transferErr *models.N1N2MessageTransferError,
) {
	var (
		requestData = n1n2MessageTransferRequest.JsonData
		n2Info      = n1n2MessageTransferRequest.BinaryDataN2Information
		n1Msg       = n1n2MessageTransferRequest.BinaryDataN1Message

		ue        *context.AmfUe
		ok        bool
		smContext *context.SmContext
		n1MsgType uint8
		anType    = models.AccessType__3_GPP_ACCESS
	)

	amfSelf := context.GetSelf()

	if ue, ok = amfSelf.AmfUeFindByUeContextID(ueContextID); !ok {
		logger.CtxLog.Warnf("AmfUe Context[%s] not found", ueContextID)
		problemDetails = &models.ProblemDetails{
			Status: http.StatusNotFound,
			Cause:  "CONTEXT_NOT_FOUND",
		}
		return nil, "", problemDetails, nil
	}

	ue.Lock.Lock()
	defer ue.Lock.Unlock()

	if requestData.N1MessageContainer != nil {
		switch requestData.N1MessageContainer.N1MessageClass {
		case models.N1MessageClass_SM:
			ue.ProducerLog.Debugf("Receive N1 SM Message (PDU Session ID: %d)", requestData.PduSessionId)
			n1MsgType = nasMessage.PayloadContainerTypeN1SMInfo
			if smContext, ok = ue.SmContextFindByPDUSessionID(requestData.PduSessionId); !ok {
				ue.ProducerLog.Errorf("SmContext[PDU Session ID:%d] not found", requestData.PduSessionId)
				problemDetails = &models.ProblemDetails{
					Status: http.StatusNotFound,
					Cause:  "CONTEXT_NOT_FOUND",
				}
				return nil, "", problemDetails, nil
			} else {
				anType = smContext.AccessType()
			}
		case models.N1MessageClass_SMS:
			n1MsgType = nasMessage.PayloadContainerTypeSMS
		case models.N1MessageClass_LPP:
			n1MsgType = nasMessage.PayloadContainerTypeLPP
		case models.N1MessageClass_UPDP:
			n1MsgType = nasMessage.PayloadContainerTypeUEPolicy
		default:
		}
	}

	if requestData.N2InfoContainer != nil {
		switch requestData.N2InfoContainer.N2InformationClass {
		case models.N2InformationClass_SM:
			ue.ProducerLog.Debugf("Receive N2 SM Message (PDU Session ID: %d)", requestData.PduSessionId)
			if smContext == nil {
				if smContext, ok = ue.SmContextFindByPDUSessionID(requestData.PduSessionId); !ok {
					ue.ProducerLog.Errorf("SmContext[PDU Session ID:%d] not found", requestData.PduSessionId)
					problemDetails = &models.ProblemDetails{
						Status: http.StatusNotFound,
						Cause:  "CONTEXT_NOT_FOUND",
					}
					return nil, "", problemDetails, nil
				} else {
					anType = smContext.AccessType()
				}
			}
		default:
			ue.ProducerLog.Warnf("N2 Information type [%s] is not supported", requestData.N2InfoContainer.N2InformationClass)
			problemDetails = &models.ProblemDetails{
				Status: http.StatusNotImplemented,
				Cause:  "NOT_IMPLEMENTED",
			}
			return nil, "", problemDetails, nil
		}
	}

	onGoing := ue.OnGoing(anType)
	// 4xx response cases
	// TODO: Error Status 307, 403 in TS29.518 Table 6.1.3.5.3.1-3
	switch onGoing.Procedure {
	case context.OnGoingProcedurePaging:
		if requestData.Ppi == 0 || (onGoing.Ppi != 0 && onGoing.Ppi <= requestData.Ppi) {
			transferErr = new(models.N1N2MessageTransferError)
			transferErr.Error = &models.ProblemDetails{
				Status: http.StatusConflict,
				Cause:  "HIGHER_PRIORITY_REQUEST_ONGOING",
			}
			return nil, "", nil, transferErr
		}
		ue.StopT3513()
		callback.SendN1N2TransferFailureNotification(ue, models.N1N2MessageTransferCause_UE_NOT_RESPONDING)
	case context.OnGoingProcedureRegistration:
		transferErr = new(models.N1N2MessageTransferError)
		transferErr.Error = &models.ProblemDetails{
			Status: http.StatusConflict,
			Cause:  "TEMPORARY_REJECT_REGISTRATION_ONGOING",
		}
		return nil, "", nil, transferErr
	case context.OnGoingProcedureN2Handover:
		transferErr = new(models.N1N2MessageTransferError)
		transferErr.Error = &models.ProblemDetails{
			Status: http.StatusConflict,
			Cause:  "TEMPORARY_REJECT_HANDOVER_ONGOING",
		}
		return nil, "", nil, transferErr
	}

	// UE is CM-Connected
	if ue.CmConnect(anType) {
		var (
			nasPdu []byte
			err    error
		)
		if n1Msg != nil {
			nasPdu, err = gmm_message.
				BuildDLNASTransport(ue, anType, n1MsgType, n1Msg, uint8(requestData.PduSessionId), nil, nil, 0)
			if err != nil {
				ue.ProducerLog.Errorf("Build DL NAS Transport error: %+v", err)
				problemDetails = &models.ProblemDetails{
					Title:  "System failure",
					Status: http.StatusInternalServerError,
					Detail: err.Error(),
					Cause:  "SYSTEM_FAILURE",
				}
				return nil, "", problemDetails, nil
			}
			if n2Info == nil {
				ue.ProducerLog.Debug("Forward N1 Message to UE")
				ngap_message.SendDownlinkNasTransport(ue.RanUe[anType], nasPdu, nil)
				n1n2MessageTransferRspData = new(models.N1N2MessageTransferRspData)
				n1n2MessageTransferRspData.Cause = models.N1N2MessageTransferCause_N1_N2_TRANSFER_INITIATED
				return n1n2MessageTransferRspData, "", nil, nil
			}
		}

		// TODO: only support transfer N2 SM information now
		if n2Info != nil {
			smInfo := requestData.N2InfoContainer.SmInfo
			switch smInfo.N2InfoContent.NgapIeType {
			case models.AmfCommunicationNgapIeType_PDU_RES_SETUP_REQ:
				ue.ProducerLog.Debugln("AMF Transfer NGAP PDU Session Resource Setup Request from SMF")
				if ue.RanUe[anType].InitialContextSetup {
					list := ngapType.PDUSessionResourceSetupListSUReq{}
					ngap_message.AppendPDUSessionResourceSetupListSUReq(&list, smInfo.PduSessionId, *smInfo.SNssai, nasPdu, n2Info)
					ngap_message.SendPDUSessionResourceSetupRequest(ue.RanUe[anType], nil, &list)
				} else {
					list := ngapType.PDUSessionResourceSetupListCxtReq{}
					ngap_message.AppendPDUSessionResourceSetupListCxtReq(&list, smInfo.PduSessionId, *smInfo.SNssai, nasPdu, n2Info)
					ngap_message.SendInitialContextSetupRequest(ue, anType, nil, &list, nil, nil, nil)
					ue.RanUe[anType].InitialContextSetup = true
				}
				n1n2MessageTransferRspData = new(models.N1N2MessageTransferRspData)
				n1n2MessageTransferRspData.Cause = models.N1N2MessageTransferCause_N1_N2_TRANSFER_INITIATED
				return n1n2MessageTransferRspData, "", nil, nil
			case models.AmfCommunicationNgapIeType_PDU_RES_MOD_REQ:
				ue.ProducerLog.Debugln("AMF Transfer NGAP PDU Session Resource Modify Request from SMF")
				list := ngapType.PDUSessionResourceModifyListModReq{}
				ngap_message.AppendPDUSessionResourceModifyListModReq(&list, smInfo.PduSessionId, nasPdu, n2Info)
				ngap_message.SendPDUSessionResourceModifyRequest(ue.RanUe[anType], list)
				n1n2MessageTransferRspData = new(models.N1N2MessageTransferRspData)
				n1n2MessageTransferRspData.Cause = models.N1N2MessageTransferCause_N1_N2_TRANSFER_INITIATED
				return n1n2MessageTransferRspData, "", nil, nil
			case models.AmfCommunicationNgapIeType_PDU_RES_REL_CMD:
				ue.ProducerLog.Debugln("AMF Transfer NGAP PDU Session Resource Release Command from SMF")
				list := ngapType.PDUSessionResourceToReleaseListRelCmd{}
				ngap_message.AppendPDUSessionResourceToReleaseListRelCmd(&list, smInfo.PduSessionId, n2Info)
				ngap_message.SendPDUSessionResourceReleaseCommand(ue.RanUe[anType], nasPdu, list)
				n1n2MessageTransferRspData = new(models.N1N2MessageTransferRspData)
				n1n2MessageTransferRspData.Cause = models.N1N2MessageTransferCause_N1_N2_TRANSFER_INITIATED
				return n1n2MessageTransferRspData, "", nil, nil
			default:
				ue.ProducerLog.Errorf("NGAP IE Type[%s] is not supported for SmInfo", smInfo.N2InfoContent.NgapIeType)
				problemDetails = &models.ProblemDetails{
					Status: http.StatusForbidden,
					Cause:  "UNSPECIFIED",
				}
				return nil, "", problemDetails, nil
			}
		}
	}

	// UE is CM-IDLE

	// 409: transfer a N2 PDU Session Resource Release Command to a 5G-AN and if the UE is in CM-IDLE
	if n2Info != nil &&
		requestData.N2InfoContainer.SmInfo.N2InfoContent.NgapIeType == models.AmfCommunicationNgapIeType_PDU_RES_REL_CMD {
		transferErr = new(models.N1N2MessageTransferError)
		transferErr.Error = &models.ProblemDetails{
			Status: http.StatusConflict,
			Cause:  "UE_IN_CM_IDLE_STATE",
		}
		return nil, "", nil, transferErr
	}
	// 504: the UE in MICO mode or the UE is only registered over Non-3GPP access and its state is CM-IDLE
	if !ue.State[models.AccessType__3_GPP_ACCESS].Is(context.Registered) {
		transferErr = new(models.N1N2MessageTransferError)
		transferErr.Error = &models.ProblemDetails{
			Status: http.StatusGatewayTimeout,
			Cause:  "UE_NOT_REACHABLE",
		}
		return nil, "", nil, transferErr
	}

	n1n2MessageTransferRspData = new(models.N1N2MessageTransferRspData)

	var pagingPriority *ngapType.PagingPriority

	var n1n2MessageID int64
	if n1n2MessageIDTmp, err := ue.N1N2MessageIDGenerator.Allocate(); err != nil {
		ue.ProducerLog.Errorf("Allocate n1n2MessageID error: %+v", err)
		problemDetails = &models.ProblemDetails{
			Status: http.StatusInternalServerError,
			Cause:  "SYSTEM_FAILURE",
			Detail: err.Error(),
		}
		return n1n2MessageTransferRspData, locationHeader, problemDetails, transferErr
	} else {
		n1n2MessageID = n1n2MessageIDTmp
	}
	locationHeader = context.GetSelf().GetIPv4Uri() + reqUri + "/" + strconv.Itoa(int(n1n2MessageID))

	// Case A (UE is CM-IDLE in 3GPP access and the associated access type is 3GPP access)
	// in subclause 5.2.2.3.1.2 of TS29518
	if anType == models.AccessType__3_GPP_ACCESS {
		if requestData.SkipInd && n2Info == nil {
			n1n2MessageTransferRspData.Cause = models.N1N2MessageTransferCause_N1_MSG_NOT_TRANSFERRED
		} else {
			n1n2MessageTransferRspData.Cause = models.N1N2MessageTransferCause_ATTEMPTING_TO_REACH_UE
			message := context.N1N2Message{
				Request:     n1n2MessageTransferRequest,
				Status:      n1n2MessageTransferRspData.Cause,
				ResourceUri: locationHeader,
			}
			ue.N1N2Message = &message
			ue.SetOnGoing(anType, &context.OnGoing{
				Procedure: context.OnGoingProcedurePaging,
				Ppi:       requestData.Ppi,
			})

			if onGoing.Ppi != 0 {
				pagingPriority = new(ngapType.PagingPriority)
				pagingPriority.Value = aper.Enumerated(onGoing.Ppi)
			}
			pkg, err := ngap_message.BuildPaging(ue, pagingPriority, false)
			if err != nil {
				logger.NgapLog.Errorf("Build Paging failed : %s", err.Error())
				return n1n2MessageTransferRspData, locationHeader, problemDetails, transferErr
			}
			ngap_message.SendPaging(ue, pkg)
		}
		// TODO: WAITING_FOR_ASYNCHRONOUS_TRANSFER
		return n1n2MessageTransferRspData, locationHeader, nil, nil
	} else {
		// Case B (UE is CM-IDLE in Non-3GPP access but CM-CONNECTED in 3GPP access and the associated
		// access type is Non-3GPP access)in subclause 5.2.2.3.1.2 of TS29518
		if ue.CmConnect(models.AccessType__3_GPP_ACCESS) {
			if n2Info == nil {
				n1n2MessageTransferRspData.Cause = models.N1N2MessageTransferCause_N1_N2_TRANSFER_INITIATED
				gmm_message.SendDLNASTransport(ue.RanUe[models.AccessType__3_GPP_ACCESS],
					nasMessage.PayloadContainerTypeN1SMInfo, n1Msg, requestData.PduSessionId, 0, nil, 0)
			} else {
				n1n2MessageTransferRspData.Cause = models.N1N2MessageTransferCause_ATTEMPTING_TO_REACH_UE
				message := context.N1N2Message{
					Request:     n1n2MessageTransferRequest,
					Status:      n1n2MessageTransferRspData.Cause,
					ResourceUri: locationHeader,
				}
				ue.N1N2Message = &message
				nasMsg, err := gmm_message.BuildNotification(ue, models.AccessType_NON_3_GPP_ACCESS)
				if err != nil {
					logger.GmmLog.Errorf("Build Notification failed : %s", err.Error())
					return n1n2MessageTransferRspData, locationHeader, problemDetails, transferErr
				}
				gmm_message.SendNotification(ue.RanUe[models.AccessType__3_GPP_ACCESS], nasMsg)
			}
			return n1n2MessageTransferRspData, locationHeader, nil, nil
		} else {
			// Case C ( UE is CM-IDLE in both Non-3GPP access and 3GPP access and the associated access ype is Non-3GPP access)
			// in subclause 5.2.2.3.1.2 of TS29518
			n1n2MessageTransferRspData.Cause = models.N1N2MessageTransferCause_ATTEMPTING_TO_REACH_UE
			message := context.N1N2Message{
				Request:     n1n2MessageTransferRequest,
				Status:      n1n2MessageTransferRspData.Cause,
				ResourceUri: locationHeader,
			}
			ue.N1N2Message = &message

			ue.SetOnGoing(anType, &context.OnGoing{
				Procedure: context.OnGoingProcedurePaging,
				Ppi:       requestData.Ppi,
			})
			if onGoing.Ppi != 0 {
				pagingPriority = new(ngapType.PagingPriority)
				pagingPriority.Value = aper.Enumerated(onGoing.Ppi)
			}
			pkg, err := ngap_message.BuildPaging(ue, pagingPriority, true)
			if err != nil {
				logger.NgapLog.Errorf("Build Paging failed : %s", err.Error())
			}
			ngap_message.SendPaging(ue, pkg)
			return n1n2MessageTransferRspData, locationHeader, nil, nil
		}
	}
}

func (p *Processor) HandleN1N2MessageTransferStatusRequest(c *gin.Context) {
	logger.CommLog.Info("Handle N1N2Message Transfer Status Request")

	ueContextID := c.Param("ueContextId")
	reqUri := c.Param("reqUri")

	status, problemDetails := p.N1N2MessageTransferStatusProcedure(ueContextID, reqUri)
	if problemDetails != nil {
		c.JSON(int(problemDetails.Status), problemDetails)
	} else {
		c.JSON(http.StatusOK, status)
	}
}

func (p *Processor) N1N2MessageTransferStatusProcedure(ueContextID string,
	reqUri string) (models.N1N2MessageTransferCause,
	*models.ProblemDetails,
) {
	amfSelf := context.GetSelf()

	ue, ok := amfSelf.AmfUeFindByUeContextID(ueContextID)
	if !ok {
		logger.CtxLog.Warnf("AmfUe Context[%s] not found", ueContextID)
		problemDetails := &models.ProblemDetails{
			Status: http.StatusNotFound,
			Cause:  "CONTEXT_NOT_FOUND",
		}
		return "", problemDetails
	}

	ue.Lock.Lock()
	defer ue.Lock.Unlock()

	resourceUri := amfSelf.GetIPv4Uri() + reqUri
	n1n2Message := ue.N1N2Message
	if n1n2Message == nil || n1n2Message.ResourceUri != resourceUri {
		problemDetails := &models.ProblemDetails{
			Status: http.StatusNotFound,
			Cause:  "CONTEXT_NOT_FOUND",
		}
		return "", problemDetails
	}

	return n1n2Message.Status, nil
}

// TS 29.518 5.2.2.3.3
func (p *Processor) HandleN1N2MessageSubscribeRequest(c *gin.Context,
	ueN1N2InfoSubscriptionCreateData models.UeN1N2InfoSubscriptionCreateData,
) {
	ueContextID := c.Param("ueContextId")

	ueN1N2InfoSubscriptionCreatedData, problemDetails := p.
		N1N2MessageSubscribeProcedure(ueContextID, ueN1N2InfoSubscriptionCreateData)
	if problemDetails != nil {
		c.JSON(int(problemDetails.Status), problemDetails)
	} else {
		c.JSON(http.StatusCreated, ueN1N2InfoSubscriptionCreatedData)
	}
}

func (p *Processor) N1N2MessageSubscribeProcedure(ueContextID string,
	ueN1N2InfoSubscriptionCreateData models.UeN1N2InfoSubscriptionCreateData) (
	*models.UeN1N2InfoSubscriptionCreatedData, *models.ProblemDetails,
) {
	amfSelf := context.GetSelf()

	ue, ok := amfSelf.AmfUeFindByUeContextID(ueContextID)
	if !ok {
		logger.CtxLog.Warnf("AmfUe Context[%s] not found", ueContextID)
		problemDetails := &models.ProblemDetails{
			Status: http.StatusNotFound,
			Cause:  "CONTEXT_NOT_FOUND",
		}
		return nil, problemDetails
	}

	ue.Lock.Lock()
	defer ue.Lock.Unlock()

	ueN1N2InfoSubscriptionCreatedData := new(models.UeN1N2InfoSubscriptionCreatedData)

	if newSubscriptionID, err := ue.N1N2MessageSubscribeIDGenerator.Allocate(); err != nil {
		logger.CommLog.Errorf("Create subscriptionID Error: %+v", err)
		problemDetails := &models.ProblemDetails{
			Status: http.StatusInternalServerError,
			Cause:  "SYSTEM_FAILURE",
		}
		return nil, problemDetails
	} else {
		ueN1N2InfoSubscriptionCreatedData.N1n2NotifySubscriptionId = strconv.Itoa(int(newSubscriptionID))
		ue.N1N2MessageSubscription.Store(newSubscriptionID, ueN1N2InfoSubscriptionCreateData)
	}
	return ueN1N2InfoSubscriptionCreatedData, nil
}

func (p *Processor) HandleN1N2MessageUnSubscribeRequest(c *gin.Context) {
	logger.CommLog.Info("Handle N1N2Message Unsubscribe Request")

	ueContextID := c.Param("ueContextId")
	subscriptionID := c.Param("subscriptionId")

	problemDetails := p.N1N2MessageUnSubscribeProcedure(ueContextID, subscriptionID)
	if problemDetails != nil {
		c.JSON(int(problemDetails.Status), problemDetails)
	} else {
		c.Status(http.StatusNoContent)
	}
}

func (p *Processor) N1N2MessageUnSubscribeProcedure(ueContextID string, subscriptionID string) *models.ProblemDetails {
	amfSelf := context.GetSelf()

	ue, ok := amfSelf.AmfUeFindByUeContextID(ueContextID)
	if !ok {
		logger.CtxLog.Warnf("AmfUe Context[%s] not found", ueContextID)
		problemDetails := &models.ProblemDetails{
			Status: http.StatusNotFound,
			Cause:  "CONTEXT_NOT_FOUND",
		}
		return problemDetails
	}

	ue.Lock.Lock()
	defer ue.Lock.Unlock()

	ue.N1N2MessageSubscription.Delete(subscriptionID)
	return nil
}
