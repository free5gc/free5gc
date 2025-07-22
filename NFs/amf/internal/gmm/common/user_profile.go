package common

import (
	"github.com/free5gc/amf/internal/context"
	"github.com/free5gc/amf/internal/logger"
	ngap_message "github.com/free5gc/amf/internal/ngap/message"
	"github.com/free5gc/amf/internal/sbi/consumer"
	"github.com/free5gc/ngap/ngapType"
	"github.com/free5gc/openapi/models"
)

func RemoveAmfUe(ue *context.AmfUe, notifyNF bool) {
	if notifyNF {
		// notify SMF to release all sessions
		ue.SmContextList.Range(func(key, value interface{}) bool {
			smContext := value.(*context.SmContext)

			problemDetail, err := consumer.GetConsumer().SendReleaseSmContextRequest(ue, smContext, nil, "", nil)
			if problemDetail != nil {
				ue.GmmLog.Errorf("Release SmContext Failed Problem[%+v]", problemDetail)
			} else if err != nil {
				ue.GmmLog.Errorf("Release SmContext Error[%v]", err.Error())
			}
			return true
		})

		// notify PCF to terminate AmPolicy association
		if ue.AmPolicyAssociation != nil {
			problemDetails, err := consumer.GetConsumer().AMPolicyControlDelete(ue)
			if problemDetails != nil {
				ue.GmmLog.Errorf("AM Policy Control Delete Failed Problem[%+v]", problemDetails)
			} else if err != nil {
				ue.GmmLog.Errorf("AM Policy Control Delete Error[%v]", err.Error())
			}
		}
	}

	PurgeAmfUeSubscriberData(ue)
	ue.Remove()
}

func PurgeAmfUeSubscriberData(ue *context.AmfUe) {
	if ue.RanUe[models.AccessType__3_GPP_ACCESS] != nil {
		err := PurgeSubscriberData(ue, models.AccessType__3_GPP_ACCESS)
		if err != nil {
			logger.GmmLog.Errorf("Purge subscriber data Error[%v]", err.Error())
		}
	}
	if ue.RanUe[models.AccessType_NON_3_GPP_ACCESS] != nil {
		err := PurgeSubscriberData(ue, models.AccessType_NON_3_GPP_ACCESS)
		if err != nil {
			logger.GmmLog.Errorf("Purge subscriber data Error[%v]", err.Error())
		}
	}
}

func AttachRanUeToAmfUeAndReleaseOldIfAny(amfUe *context.AmfUe, ranUe *context.RanUe) {
	if oldRanUe := amfUe.RanUe[ranUe.Ran.AnType]; oldRanUe != nil {
		oldRanUe.Log.Infof("Implicit Deregistration - RanUeNgapID[%d]", oldRanUe.RanUeNgapId)
		oldRanUe.DetachAmfUe()
		if amfUe.T3550 != nil {
			amfUe.State[ranUe.Ran.AnType].Set(context.Registered)
		}
		StopAll5GSMMTimers(amfUe)
		causeGroup := ngapType.CausePresentRadioNetwork
		causeValue := ngapType.CauseRadioNetworkPresentReleaseDueToNgranGeneratedReason
		ngap_message.SendUEContextReleaseCommand(oldRanUe, context.UeContextReleaseUeContext, causeGroup, causeValue)
	}
	amfUe.AttachRanUe(ranUe)
}

func AttachRanUeToAmfUeAndReleaseOldHandover(amfUe *context.AmfUe, sourceRanUe, targetRanUe *context.RanUe) {
	logger.GmmLog.Debugln("In AttachRanUeToAmfUeAndReleaseOldHandover")

	if sourceRanUe != nil {
		sourceRanUe.DetachAmfUe()
		if amfUe.T3550 != nil {
			amfUe.State[targetRanUe.Ran.AnType].Set(context.Registered)
		}
		StopAll5GSMMTimers(amfUe)
		causeGroup := ngapType.CausePresentRadioNetwork
		causeValue := ngapType.CauseRadioNetworkPresentSuccessfulHandover
		ngap_message.SendUEContextReleaseCommand(sourceRanUe, context.UeContextReleaseHandover, causeGroup, causeValue)
	} else {
		// This function will be call only by N2 Handover, so we can assume sourceRanUe will not be nil
		logger.GmmLog.Errorln("AttachRanUeToAmfUeAndReleaseOldHandover() is called but sourceRanUe is nil")
	}
	amfUe.AttachRanUe(targetRanUe)
}

func ClearHoldingRanUe(ranUe *context.RanUe) {
	if ranUe != nil {
		ranUe.DetachAmfUe()
		ranUe.Log.Infof("Clear Holding RanUE")
		causeGroup := ngapType.CausePresentRadioNetwork
		causeValue := ngapType.CauseRadioNetworkPresentReleaseDueToNgranGeneratedReason
		ngap_message.SendUEContextReleaseCommand(ranUe, context.UeContextReleaseUeContext, causeGroup, causeValue)
	} else {
		logger.GmmLog.Warnf("RanUE is nil")
	}
}

func PurgeSubscriberData(ue *context.AmfUe, accessType models.AccessType) error {
	logger.GmmLog.Debugln("PurgeSubscriberData")

	if !ue.ContextValid {
		return nil
	}
	// Purge of subscriber data in AMF described in TS 23.502 4.5.3
	if ue.SdmSubscriptionId != "" {
		problemDetails, err := consumer.GetConsumer().SDMUnsubscribe(ue)
		if problemDetails != nil {
			logger.GmmLog.Errorf("SDM Unubscribe Failed Problem[%+v]", problemDetails)
		} else if err != nil {
			logger.GmmLog.Errorf("SDM Unubscribe Error[%+v]", err)
		}
		ue.SdmSubscriptionId = ""
	}

	if ue.UeCmRegistered[accessType] {
		problemDetails, err := consumer.GetConsumer().UeCmDeregistration(ue, accessType)
		if problemDetails != nil {
			logger.GmmLog.Errorf("UECM Deregistration Failed Problem[%+v]", problemDetails)
		} else if err != nil {
			logger.GmmLog.Errorf("UECM Deregistration Error[%+v]", err)
		}
		ue.UeCmRegistered[accessType] = false
	}
	return nil
}
