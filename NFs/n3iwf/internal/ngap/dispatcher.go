package ngap

import (
	"runtime/debug"

	"github.com/free5gc/n3iwf/internal/logger"
	"github.com/free5gc/ngap"
	"github.com/free5gc/ngap/ngapType"
	"github.com/free5gc/sctp"
)

func (s *Server) NGAPDispatch(conn *sctp.SCTPConn, msg []byte) {
	ngapLog := logger.NgapLog

	defer func() {
		if p := recover(); p != nil {
			// Print stack for panic to log. Fatalf() will let program exit.
			ngapLog.Fatalf("panic: %v\n%s", p, string(debug.Stack()))
		}
	}()

	// AMF SCTP address
	sctpAddr := conn.RemoteAddr().String()
	// AMF context
	n3iwfCtx := s.Context()
	amf, _ := n3iwfCtx.AMFPoolLoad(sctpAddr)
	// Decode
	pdu, err := ngap.Decoder(msg)
	if err != nil {
		ngapLog.Errorf("NGAP decode error: %+v\n", err)
		return
	}

	switch pdu.Present {
	case ngapType.NGAPPDUPresentInitiatingMessage:
		initiatingMessage := pdu.InitiatingMessage
		if initiatingMessage == nil {
			ngapLog.Errorln("Initiating Message is nil")
			return
		}

		switch initiatingMessage.ProcedureCode.Value {
		case ngapType.ProcedureCodeNGReset:
			s.HandleNGReset(amf, pdu)
		case ngapType.ProcedureCodeInitialContextSetup:
			s.HandleInitialContextSetupRequest(amf, pdu)
		case ngapType.ProcedureCodeUEContextModification:
			s.HandleUEContextModificationRequest(amf, pdu)
		case ngapType.ProcedureCodeUEContextRelease:
			s.HandleUEContextReleaseCommand(amf, pdu)
		case ngapType.ProcedureCodeDownlinkNASTransport:
			s.HandleDownlinkNASTransport(amf, pdu)
		case ngapType.ProcedureCodePDUSessionResourceSetup:
			s.HandlePDUSessionResourceSetupRequest(amf, pdu)
		case ngapType.ProcedureCodePDUSessionResourceModify:
			s.HandlePDUSessionResourceModifyRequest(amf, pdu)
		case ngapType.ProcedureCodePDUSessionResourceRelease:
			s.HandlePDUSessionResourceReleaseCommand(amf, pdu)
		case ngapType.ProcedureCodeErrorIndication:
			s.HandleErrorIndication(amf, pdu)
		case ngapType.ProcedureCodeUERadioCapabilityCheck:
			s.HandleUERadioCapabilityCheckRequest(amf, pdu)
		case ngapType.ProcedureCodeAMFConfigurationUpdate:
			s.HandleAMFConfigurationUpdate(amf, pdu)
		case ngapType.ProcedureCodeDownlinkRANConfigurationTransfer:
			s.HandleDownlinkRANConfigurationTransfer(pdu)
		case ngapType.ProcedureCodeDownlinkRANStatusTransfer:
			s.HandleDownlinkRANStatusTransfer(pdu)
		case ngapType.ProcedureCodeAMFStatusIndication:
			s.HandleAMFStatusIndication(pdu)
		case ngapType.ProcedureCodeLocationReportingControl:
			s.HandleLocationReportingControl(pdu)
		case ngapType.ProcedureCodeUETNLABindingRelease:
			s.HandleUETNLAReleaseRequest(pdu)
		case ngapType.ProcedureCodeOverloadStart:
			s.HandleOverloadStart(amf, pdu)
		case ngapType.ProcedureCodeOverloadStop:
			s.HandleOverloadStop(amf, pdu)
		default:
			ngapLog.Warnf("Not implemented NGAP message(initiatingMessage), procedureCode:%d]\n",
				initiatingMessage.ProcedureCode.Value)
		}
	case ngapType.NGAPPDUPresentSuccessfulOutcome:
		successfulOutcome := pdu.SuccessfulOutcome
		if successfulOutcome == nil {
			ngapLog.Errorln("Successful Outcome is nil")
			return
		}

		switch successfulOutcome.ProcedureCode.Value {
		case ngapType.ProcedureCodeNGSetup:
			s.HandleNGSetupResponse(sctpAddr, conn, pdu)
		case ngapType.ProcedureCodeNGReset:
			s.HandleNGResetAcknowledge(amf, pdu)
		case ngapType.ProcedureCodePDUSessionResourceModifyIndication:
			s.HandlePDUSessionResourceModifyConfirm(amf, pdu)
		case ngapType.ProcedureCodeRANConfigurationUpdate:
			s.HandleRANConfigurationUpdateAcknowledge(amf, pdu)
		default:
			ngapLog.Warnf("Not implemented NGAP message(successfulOutcome), procedureCode:%d]\n",
				successfulOutcome.ProcedureCode.Value)
		}
	case ngapType.NGAPPDUPresentUnsuccessfulOutcome:
		unsuccessfulOutcome := pdu.UnsuccessfulOutcome
		if unsuccessfulOutcome == nil {
			ngapLog.Errorln("Unsuccessful Outcome is nil")
			return
		}

		switch unsuccessfulOutcome.ProcedureCode.Value {
		case ngapType.ProcedureCodeNGSetup:
			s.HandleNGSetupFailure(sctpAddr, conn, pdu)
		case ngapType.ProcedureCodeRANConfigurationUpdate:
			s.HandleRANConfigurationUpdateFailure(amf, pdu)
		default:
			ngapLog.Warnf("Not implemented NGAP message(unsuccessfulOutcome), procedureCode:%d]\n",
				unsuccessfulOutcome.ProcedureCode.Value)
		}
	}
}
