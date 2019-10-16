package amf_util

import (
	"free5gc/lib/ngap/ngapType"
	"free5gc/lib/openapi/models"
	"free5gc/lib/timer"
	"free5gc/src/amf/amf_context"
	"free5gc/src/amf/amf_handler/amf_message"
	"free5gc/src/amf/logger"
)

func StartT3513(ue *amf_context.AmfUe) {
	if ue == nil {
		logger.UtilLog.Error("AmfUe is nil")
		return
	}
	msg := amf_message.HandlerMessage{
		Event: amf_message.EventGMMT3513,
		Value: ue,
	}
	ue.T3513 = timer.StartTimer(amf_context.TimeT3513, func(msg interface{}) {
		amf_message.SendMessage(msg.(amf_message.HandlerMessage))
	}, msg)
}
func ClearT3513(ue *amf_context.AmfUe) {
	if ue == nil {
		logger.UtilLog.Error("AmfUe is nil")
		return
	}
	if ue.T3513 != nil {
		ue.T3513.Stop()
		ue.T3513 = nil
	}
	ue.PagingRetryTimes = 0
	ue.LastPagingPkg = nil
	ue.OnGoing[models.AccessType__3_GPP_ACCESS].Ppi = 0
	ue.OnGoing[models.AccessType__3_GPP_ACCESS].Procedure = amf_context.OnGoingProcedureNothing
}

func StartT3565(ue *amf_context.RanUe) {
	if ue == nil {
		logger.UtilLog.Error("RanUe is nil")
		return
	}
	if ue.AmfUe == nil {
		logger.UtilLog.Error("AmfUe is nil")
		return
	}
	msg := amf_message.HandlerMessage{
		Event: amf_message.EventGMMT3565,
		Value: ue,
	}
	ue.AmfUe.T3565 = timer.StartTimer(amf_context.TimeT3565, func(msg interface{}) {
		amf_message.SendMessage(msg.(amf_message.HandlerMessage))
	}, msg)
}

func ClearT3565(ue *amf_context.AmfUe) {
	if ue == nil {
		logger.UtilLog.Error("AmfUe is nil")
		return
	}
	if ue.T3565 != nil {
		ue.T3565.Stop()
		ue.T3565 = nil
	}
	ue.NotificationRetryTimes = 0
	ue.LastNotificationPkg = nil
}

func StartT3560(ue *amf_context.RanUe, event amf_message.Event, eapSuccess *bool, eapMessage *string) {
	if ue == nil {
		logger.UtilLog.Error("RanUe is nil")
		return
	}
	if ue.AmfUe == nil {
		logger.UtilLog.Error("AmfUe is nil")
		return
	}

	var msg amf_message.HandlerMessage
	switch event {
	case amf_message.EventGMMT3560ForAuthenticationRequest:
		msg = amf_message.HandlerMessage{
			Event: amf_message.EventGMMT3560ForAuthenticationRequest,
			Value: ue,
		}
	case amf_message.EventGMMT3560ForSecurityModeCommand:
		msg = amf_message.HandlerMessage{
			Event: amf_message.EventGMMT3560ForSecurityModeCommand,
			Value: amf_message.EventGMMT3560ValueForSecurityCommand{
				RanUe:      ue,
				EapSuccess: *eapSuccess,
				EapMessage: *eapMessage,
			},
		}
	}
	ue.AmfUe.T3560 = timer.StartTimer(amf_context.TimeT3560, func(msg interface{}) {
		amf_message.SendMessage(msg.(amf_message.HandlerMessage))
	}, msg)
}

func ClearT3560(ue *amf_context.AmfUe) {
	if ue == nil {
		logger.UtilLog.Error("AmfUe is nil")
		return
	}

	if ue.T3560 != nil {
		ue.T3560.Stop()
		ue.T3560 = nil
	}
	ue.T3560RetryTimes = 0
}

func StartT3550(
	ue *amf_context.AmfUe,
	accessType models.AccessType,
	pDUSessionStatus *[16]bool,
	reactivationResult *[16]bool,
	errPduSessionId, errCause []uint8,
	pduSessionResourceSetupList *ngapType.PDUSessionResourceSetupListCxtReq) {

	if ue == nil {
		logger.UtilLog.Error("AmfUe is nil")
		return
	}
	msg := amf_message.HandlerMessage{
		Event: amf_message.EventGMMT3550,
		Value: amf_message.EventGMMT3550Value{
			AmfUe:                       ue,
			AccessType:                  accessType,
			PDUSessionStatus:            pDUSessionStatus,
			ReactivationResult:          reactivationResult,
			ErrPduSessionId:             errPduSessionId,
			ErrCause:                    errCause,
			PduSessionResourceSetupList: pduSessionResourceSetupList,
		},
	}
	ue.T3550 = timer.StartTimer(amf_context.TimeT3550, func(msg interface{}) {
		amf_message.SendMessage(msg.(amf_message.HandlerMessage))
	}, msg)
}

func ClearT3550(ue *amf_context.AmfUe) {
	if ue == nil {
		logger.UtilLog.Error("AmfUe is nil")
		return
	}

	if ue.T3550 != nil {
		ue.T3550.Stop()
		ue.T3550 = nil
	}
	ue.T3550RetryTimes = 0
}

func StartT3522(ue *amf_context.RanUe, accessType *uint8, reRegistrationRequired *bool, cause5GMM *uint8) {
	if ue == nil {
		logger.UtilLog.Error("RanUe is nil")
		return
	}
	if ue.AmfUe == nil {
		logger.UtilLog.Error("AmfUe is nil")
		return
	}
	msg := amf_message.HandlerMessage{
		Event: amf_message.EventGMMT3522,
		Value: amf_message.EventGMMT3522Value{
			RanUe:                  ue,
			AccessType:             *accessType,
			ReRegistrationRequired: *reRegistrationRequired,
			Cause5GMM:              *cause5GMM,
		},
	}
	ue.AmfUe.T3522 = timer.StartTimer(amf_context.TimeT3522, func(msg interface{}) {
		amf_message.SendMessage(msg.(amf_message.HandlerMessage))
	}, msg)
}
func ClearT3522(ue *amf_context.AmfUe) {
	if ue == nil {
		logger.UtilLog.Error("AmfUe is nil")
		return
	}

	if ue.T3522 != nil {
		ue.T3522.Stop()
		ue.T3522 = nil
	}
	ue.T3522RetryTimes = 0

}
