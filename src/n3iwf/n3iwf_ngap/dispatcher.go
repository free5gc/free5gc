package n3iwf_ngap

import (
	"github.com/sirupsen/logrus"
	"free5gc/lib/ngap"
	"free5gc/lib/ngap/ngapType"
	"free5gc/src/n3iwf/logger"
	"free5gc/src/n3iwf/n3iwf_context"
	"free5gc/src/n3iwf/n3iwf_ngap/ngap_handler"
)

var Ngaplog *logrus.Entry

func init() {
	Ngaplog = logger.NgapLog
}

func Dispatch(sctpAddr string, msg []byte) {
	pdu, err := ngap.Decoder(msg)
	if err != nil {
		Ngaplog.Errorf("NGAP decode error: %+v\n", err)
		return
	}
	amf := n3iwf_context.N3IWFSelf().AMFPool[sctpAddr]

	switch pdu.Present {
	case ngapType.NGAPPDUPresentInitiatingMessage:
		initiatingMessage := pdu.InitiatingMessage
		if initiatingMessage == nil {
			Ngaplog.Errorln("Initiating Message is nil")
			return
		}

		switch initiatingMessage.ProcedureCode.Value {
		case ngapType.ProcedureCodeNGReset:
			ngap_handler.HandleNGReset(amf, pdu)
		case ngapType.ProcedureCodeInitialContextSetup:
			ngap_handler.HandleInitialContextSetupRequest(amf, pdu)
		case ngapType.ProcedureCodeUEContextModification:
			ngap_handler.HandleUEContextModificationRequest(amf, pdu)
		case ngapType.ProcedureCodeUEContextRelease:
			ngap_handler.HandleUEContextReleaseCommand(amf, pdu)
		case ngapType.ProcedureCodeDownlinkNASTransport:
			ngap_handler.HandleDownlinkNASTransport(amf, pdu)
		case ngapType.ProcedureCodePDUSessionResourceSetup:
			ngap_handler.HandlePDUSessionResourceSetupRequest(amf, pdu)
		case ngapType.ProcedureCodePDUSessionResourceModify:
			ngap_handler.HandlePDUSessionResourceModifyRequest(amf, pdu)
		case ngapType.ProcedureCodePDUSessionResourceRelease:
			ngap_handler.HandlePDUSessionResourceReleaseCommand(amf, pdu)
		case ngapType.ProcedureCodeErrorIndication:
			ngap_handler.HandleErrorIndication(amf, pdu)
		case ngapType.ProcedureCodeUERadioCapabilityCheck:
			ngap_handler.HandleUERadioCapabilityCheckRequest(amf, pdu)
		case ngapType.ProcedureCodeAMFConfigurationUpdate:
			ngap_handler.HandleAMFConfigurationUpdate(amf, pdu)
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
			ngap_handler.HandleOverloadStart(amf, pdu)
		case ngapType.ProcedureCodeOverloadStop:
			ngap_handler.HandleOverloadStop(amf, pdu)
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
			ngap_handler.HandleNGSetupResponse(sctpAddr, pdu)
		case ngapType.ProcedureCodeNGReset:
			ngap_handler.HandleNGResetAcknowledge(amf, pdu)
		case ngapType.ProcedureCodePDUSessionResourceModifyIndication:
			ngap_handler.HandlePDUSessionResourceModifyConfirm(amf, pdu)
		case ngapType.ProcedureCodeRANConfigurationUpdate:
			ngap_handler.HandleRANConfigurationUpdateAcknowledge(amf, pdu)
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
			ngap_handler.HandleNGSetupFailure(sctpAddr, pdu)
		case ngapType.ProcedureCodeRANConfigurationUpdate:
			ngap_handler.HandleRANConfigurationUpdateFailure(amf, pdu)
		default:
			Ngaplog.Warnf("Not implemented NGAP message(unsuccessfulOutcome), procedureCode:%d]\n", unsuccessfulOutcome.ProcedureCode.Value)
		}
	}
}
