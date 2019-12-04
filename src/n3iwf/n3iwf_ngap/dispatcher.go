package n3iwf_ngap

import (
	"github.com/sirupsen/logrus"
	"free5gc/lib/ngap"
	"free5gc/lib/ngap/ngapType"
	"free5gc/src/n3iwf/logger"
	"free5gc/src/n3iwf/n3iwf_ngap/ngap_handler"
)

var Ngaplog *logrus.Entry

func init() {
	Ngaplog = logger.NgapLog
}

func Dispatch(sessionID string, msg []byte) {
	pdu, err := ngap.Decoder(msg)
	if err != nil {
		Ngaplog.Errorf("NGAP decode error: %+v\n", err)
		return
	}

	switch pdu.Present {
	case ngapType.NGAPPDUPresentInitiatingMessage:
		initiatingMessage := pdu.InitiatingMessage
		if initiatingMessage == nil {
			Ngaplog.Errorln("Initiating Message is nil")
			return
		}

		switch initiatingMessage.ProcedureCode.Value {
		case ngapType.ProcedureCodeNGReset:
			ngap_handler.HandleNGReset(pdu)
		case ngapType.ProcedureCodeInitialContextSetup:
			ngap_handler.HandleInitialContextSetupRequest(sessionID, pdu)
		case ngapType.ProcedureCodeUEContextModification:
			ngap_handler.HandleUEContextModificationRequest(sessionID, pdu)
		case ngapType.ProcedureCodeUEContextRelease:
			ngap_handler.HandleUEContextReleaseCommand(sessionID, pdu)
		case ngapType.ProcedureCodeDownlinkNASTransport:
			ngap_handler.HandleDownlinkNASTransport(pdu)
		case ngapType.ProcedureCodePDUSessionResourceSetup:
			ngap_handler.HandlePDUSessionResourceSetupRequest(pdu)
		case ngapType.ProcedureCodePDUSessionResourceModify:
			ngap_handler.HandlePDUSessionResourceModifyRequest(pdu)
		case ngapType.ProcedureCodePDUSessionResourceRelease:
			ngap_handler.HandlePDUSessionResourceReleaseCommand(pdu)
		case ngapType.ProcedureCodeErrorIndication:
			ngap_handler.HandleErrorIndication(pdu)
		case ngapType.ProcedureCodeUERadioCapabilityCheck:
			ngap_handler.HandleUERadioCapabilityCheckRequest(pdu)
		case ngapType.ProcedureCodeAMFConfigurationUpdate:
			ngap_handler.HandleAMFConfigurationUpdate(pdu)
		case ngapType.ProcedureCodeDownlinkRANConfigurationTransfer:
			ngap_handler.HandleDownlinkRANConfigurationTransfer(pdu)
		case ngapType.ProcedureCodeDownlinkRANStatusTransfer:
			ngap_handler.HandleDownlinkRANStatusTransfer(pdu)
		case ngapType.ProcedureCodeAMFStatusIndication:
			ngap_handler.HandleAMFStatusIndication(pdu)
		case ngapType.ProcedureCodeLocationReportingControl:
			ngap_handler.HandleLocationReportingControl(pdu)
		case ngapType.ProcedureCodeUETNLABindingRelease:
			ngap_handler.HandleUETNLAReleaseRequest(pdu)
		case ngapType.ProcedureCodeOverloadStart:
			ngap_handler.HandleOverloadStart(pdu)
		case ngapType.ProcedureCodeOverloadStop:
			ngap_handler.HandleOverloadStop(pdu)
		default:
			Ngaplog.Warnf("Not implemented NGAP message(initiatingMessage), procedureCode:%d]\n", initiatingMessage.ProcedureCode.Value)
		}
	case ngapType.NGAPPDUPresentSuccessfulOutcome:
		successfulOutcome := pdu.SuccessfulOutcome
		if successfulOutcome == nil {
			Ngaplog.Errorln("Successful Outcome is nil")
			return
		}

		switch successfulOutcome.ProcedureCode.Value {
		case ngapType.ProcedureCodeNGSetup:
			ngap_handler.HandleNGSetupResponse(sessionID, pdu)
		case ngapType.ProcedureCodeNGReset:
			ngap_handler.HandleNGResetAcknowledge(pdu)
		case ngapType.ProcedureCodePDUSessionResourceModifyIndication:
			ngap_handler.HandlePDUSessionResourceModifyConfirm(pdu)
		case ngapType.ProcedureCodeRANConfigurationUpdate:
			ngap_handler.HandleRANConfigurationUpdateAcknowledge(pdu)
		default:
			Ngaplog.Warnf("Not implemented NGAP message(successfulOutcome), procedureCode:%d]\n", successfulOutcome.ProcedureCode.Value)
		}
	case ngapType.NGAPPDUPresentUnsuccessfulOutcome:
		unsuccessfulOutcome := pdu.UnsuccessfulOutcome
		if unsuccessfulOutcome == nil {
			Ngaplog.Errorln("Unsuccessful Outcome is nil")
			return
		}

		switch unsuccessfulOutcome.ProcedureCode.Value {
		case ngapType.ProcedureCodeNGSetup:
			ngap_handler.HandleNGSetupFailure(sessionID, pdu)
		case ngapType.ProcedureCodeRANConfigurationUpdate:
			ngap_handler.HandleRANConfigurationUpdateFailure(pdu)
		default:
			Ngaplog.Warnf("Not implemented NGAP message(unsuccessfulOutcome), procedureCode:%d]\n", unsuccessfulOutcome.ProcedureCode.Value)
		}
	}
}
