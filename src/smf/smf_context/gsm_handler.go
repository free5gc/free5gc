package smf_context

import (
	"free5gc/lib/nas/nasConvert"
	"free5gc/lib/nas/nasMessage"
	"free5gc/src/smf/logger"
)

func (smContext *SMContext) HandlePDUSessionEstablishmentRequest(req *nasMessage.PDUSessionEstablishmentRequest) {
	// Retrieve PDUSessionID
	smContext.PDUSessionID = int32(req.PDUSessionID.GetPDUSessionID())

	// Handle PDUSessionType
	if req.PDUSessionType != nil {
		requestedPDUSessionType := req.PDUSessionType.GetPDUSessionTypeValue()
		if smContext.isAllowedPDUSessionType(requestedPDUSessionType) {
			smContext.SelectedPDUSessionType = requestedPDUSessionType
		} else {
			logger.CtxLog.Errorf("requested pdu session type [%s] is not in allowed type\n", nasConvert.PDUSessionTypeToModels(requestedPDUSessionType))
		}
	} else {
		// Default to IPv4
		// TODO: use Default PDU Session Type
		smContext.SelectedPDUSessionType = nasMessage.PDUSessionTypeIPv4
	}

	smContext.PDUAddress = AllocUEIP()
}

func (smContext *SMContext) HandlePDUSessionReleaseRequest(req *nasMessage.PDUSessionReleaseRequest) {
	logger.GsmLog.Infof("Handle Pdu Session Release Request")
}
