package processor

import (
	"context"
	"fmt"
	"time"

	"github.com/free5gc/nas/nasMessage"
	"github.com/free5gc/openapi/models"
	"github.com/free5gc/pfcp"
	"github.com/free5gc/pfcp/pfcpType"
	smf_context "github.com/free5gc/smf/internal/context"
	"github.com/free5gc/smf/internal/logger"
	"github.com/free5gc/smf/internal/pfcp/message"
)

func (p *Processor) ToBeAssociatedWithUPF(smfPfcpContext context.Context, upf *smf_context.UPF) {
	var upfStr string
	if upf.NodeID.NodeIdType == pfcpType.NodeIdTypeFqdn {
		upfStr = fmt.Sprintf("[%s](%s)", upf.NodeID.FQDN, upf.NodeID.ResolveNodeIdToIp().String())
	} else {
		upfStr = fmt.Sprintf("[%s]", upf.NodeID.ResolveNodeIdToIp().String())
	}

	for {
		// check if SMF PFCP context (parent) was canceled
		// note: UPF AssociationContexts are children of smfPfcpContext
		select {
		case <-smfPfcpContext.Done():
			logger.MainLog.Infoln("Canceled SMF PFCP context")
			return
		default:
			ensureSetupPfcpAssociation(smfPfcpContext, upf, upfStr)
			if smf_context.GetSelf().PfcpHeartbeatInterval == 0 {
				return
			}
			keepHeartbeatTo(upf, upfStr)
			// returns when UPF heartbeat loss is detected or association is canceled

			p.releaseAllResourcesOfUPF(upf, upfStr)
		}
	}
}

func (p *Processor) ReleaseAllResourcesOfUPF(upf *smf_context.UPF) {
	var upfStr string
	if upf.NodeID.NodeIdType == pfcpType.NodeIdTypeFqdn {
		upfStr = fmt.Sprintf("[%s](%s)", upf.NodeID.FQDN, upf.NodeID.ResolveNodeIdToIp().String())
	} else {
		upfStr = fmt.Sprintf("[%s]", upf.NodeID.ResolveNodeIdToIp().String())
	}
	p.releaseAllResourcesOfUPF(upf, upfStr)
}

func ensureSetupPfcpAssociation(parentContext context.Context, upf *smf_context.UPF, upfStr string) {
	alertTime := time.Now()
	alertInterval := smf_context.GetSelf().AssocFailAlertInterval
	retryInterval := smf_context.GetSelf().AssocFailRetryInterval
	for {
		err := setupPfcpAssociation(upf, upfStr)
		if err == nil {
			// success
			// assign UPF an AssociationContext, with SMF PFCP Context as parent
			upf.AssociationContext, upf.CancelAssociation = context.WithCancel(parentContext)
			return
		}
		logger.MainLog.Warnf("Failed to setup an association with UPF[%s], error:%+v", upfStr, err)
		now := time.Now()
		logger.MainLog.Debugf("now %+v, alertTime %+v", now, alertTime)
		if now.After(alertTime.Add(alertInterval)) {
			logger.MainLog.Errorf("ALERT for UPF[%s]", upfStr)
			alertTime = now
		}
		logger.MainLog.Debugf("Wait %+v until next retry attempt", retryInterval)
		timer := time.After(retryInterval)
		select { // no default case, either case needs to be true to continue
		case <-parentContext.Done():
			logger.MainLog.Infoln("Canceled SMF PFCP context")
			return
		case <-timer:
			continue
		}
	}
}

func setupPfcpAssociation(upf *smf_context.UPF, upfStr string) error {
	logger.MainLog.Infof("Sending PFCP Association Request to UPF%s", upfStr)

	resMsg, err := message.SendPfcpAssociationSetupRequest(upf.NodeID)
	if err != nil {
		return err
	}

	rsp := resMsg.PfcpMessage.Body.(pfcp.PFCPAssociationSetupResponse)

	if rsp.Cause == nil || rsp.Cause.CauseValue != pfcpType.CauseRequestAccepted {
		return fmt.Errorf("received PFCP Association Setup Not Accepted Response from UPF%s", upfStr)
	}

	nodeID := rsp.NodeID
	if nodeID == nil {
		return fmt.Errorf("pfcp association needs NodeID")
	}

	logger.MainLog.Infof("Received PFCP Association Setup Accepted Response from UPF%s", upfStr)
	logger.MainLog.Infof("UPF(%s) setup association", upf.NodeID.ResolveNodeIdToIp().String())

	return nil
}

func keepHeartbeatTo(upf *smf_context.UPF, upfStr string) {
	for {
		err := doPfcpHeartbeat(upf, upfStr)
		if err != nil {
			logger.MainLog.Errorf("PFCP Heartbeat error: %v", err)
			return
		}

		timer := time.After(smf_context.GetSelf().PfcpHeartbeatInterval)
		select {
		case <-upf.AssociationContext.Done():
			logger.MainLog.Infof("Canceled association to UPF[%s]", upfStr)
			return
		case <-timer:
			continue
		}
	}
}

func doPfcpHeartbeat(upf *smf_context.UPF, upfStr string) error {
	if err := upf.IsAssociated(); err != nil {
		return fmt.Errorf("cancel heartbeat: %+v", err)
	}

	logger.MainLog.Debugf("Sending PFCP Heartbeat Request to UPF%s", upfStr)

	resMsg, err := message.SendPfcpHeartbeatRequest(upf)
	if err != nil {
		upf.CancelAssociation()
		upf.RecoveryTimeStamp = time.Time{}
		return fmt.Errorf("SendPfcpHeartbeatRequest error: %w", err)
	}

	rsp := resMsg.PfcpMessage.Body.(pfcp.HeartbeatResponse)
	if rsp.RecoveryTimeStamp == nil {
		logger.MainLog.Warnf("Received PFCP Heartbeat Response without timestamp from UPF%s", upfStr)
		return nil
	}

	logger.MainLog.Debugf("Received PFCP Heartbeat Response from UPF%s", upfStr)
	if upf.RecoveryTimeStamp.IsZero() {
		// first receive
		upf.RecoveryTimeStamp = rsp.RecoveryTimeStamp.RecoveryTimeStamp
	} else if upf.RecoveryTimeStamp.Before(rsp.RecoveryTimeStamp.RecoveryTimeStamp) {
		// received a newer recovery timestamp
		upf.CancelAssociation()
		upf.RecoveryTimeStamp = time.Time{}
		return fmt.Errorf("received PFCP Heartbeat Response RecoveryTimeStamp has been updated")
	}
	return nil
}

func (p *Processor) releaseAllResourcesOfUPF(upf *smf_context.UPF, upfStr string) {
	logger.MainLog.Infof("Release all resources of UPF %s", upfStr)

	upf.ProcEachSMContext(func(smContext *smf_context.SMContext) {
		smContext.SMLock.Lock()
		defer smContext.SMLock.Unlock()
		switch smContext.State() {
		case smf_context.Active, smf_context.ModificationPending, smf_context.PFCPModification:
			needToSendNotify, removeContext := p.requestAMFToReleasePDUResources(smContext)
			if needToSendNotify {
				p.SendReleaseNotification(smContext)
			}
			if removeContext {
				// Notification has already been sent, if it is needed
				p.RemoveSMContextFromAllNF(smContext, false)
			}
		}
	})
}

func (p *Processor) requestAMFToReleasePDUResources(
	smContext *smf_context.SMContext,
) (sendNotify bool, releaseContext bool) {
	n1n2Request := models.N1N2MessageTransferRequest{}
	// TS 23.502 4.3.4.2 3b. Send Namf_Communication_N1N2MessageTransfer Request, SMF->AMF
	n1n2Request.JsonData = &models.N1N2MessageTransferReqData{
		PduSessionId: smContext.PDUSessionID,
		SkipInd:      true,
	}
	cause := nasMessage.Cause5GSMNetworkFailure
	if buf, err := smf_context.BuildGSMPDUSessionReleaseCommand(smContext, cause, false); err != nil {
		logger.MainLog.Errorf("Build GSM PDUSessionReleaseCommand failed: %+v", err)
	} else {
		n1n2Request.BinaryDataN1Message = buf
		n1n2Request.JsonData.N1MessageContainer = &models.N1MessageContainer{
			N1MessageClass:   "SM",
			N1MessageContent: &models.RefToBinaryData{ContentId: "GSM_NAS"},
		}
	}
	if smContext.UpCnxState != models.UpCnxState_DEACTIVATED {
		if buf, err := smf_context.BuildPDUSessionResourceReleaseCommandTransfer(smContext); err != nil {
			logger.MainLog.Errorf("Build PDUSessionResourceReleaseCommandTransfer failed: %+v", err)
		} else {
			n1n2Request.BinaryDataN2Information = buf
			n1n2Request.JsonData.N2InfoContainer = &models.N2InfoContainer{
				N2InformationClass: models.N2InformationClass_SM,
				SmInfo: &models.N2SmInformation{
					PduSessionId: smContext.PDUSessionID,
					N2InfoContent: &models.N2InfoContent{
						NgapIeType: models.AmfCommunicationNgapIeType_PDU_RES_REL_CMD,
						NgapData: &models.RefToBinaryData{
							ContentId: "N2SmInformation",
						},
					},
					SNssai: smContext.SNssai,
				},
			}
		}
	}

	ctx, _, errToken := smf_context.GetSelf().GetTokenCtx(models.ServiceName_NAMF_COMM, models.NrfNfManagementNfType_AMF)
	if errToken != nil {
		return false, false
	}

	rspData, err := p.Consumer().
		N1N2MessageTransfer(ctx, smContext.Supi, n1n2Request, smContext.CommunicationClientApiPrefix)

	if err != nil || rspData == nil {
		logger.ConsumerLog.Warnf("N1N2MessageTransfer for RequestAMFToReleasePDUResources failed: %+v", err)
		// keep SM Context to avoid inconsistency with AMF
		smContext.SetState(smf_context.InActive)
	} else {
		switch rspData.Cause {
		case models.N1N2MessageTransferCause_N1_MSG_NOT_TRANSFERRED:
			// the PDU Session Release Command was not transferred to the UE since it is in CM-IDLE state.
			//   ref. step3b of "4.3.4.2 UE or network requested PDU Session Release for Non-Roaming and
			//        Roaming with Local Breakout" in TS23.502
			// it is needed to remove both AMF's and SMF's SM Contexts immediately
			smContext.SetState(smf_context.InActive)
			return true, true
		case models.N1N2MessageTransferCause_N1_N2_TRANSFER_INITIATED:
			// wait for N2 PDU Session Release Response
			smContext.SetState(smf_context.InActivePending)
		default:
			// other causes are unexpected.
			// keep SM Context to avoid inconsistency with AMF
			smContext.SetState(smf_context.InActive)
		}
	}
	return false, false
}
