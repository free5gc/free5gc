package processor

import (
	"fmt"

	"github.com/free5gc/nas"
	"github.com/free5gc/nas/nasConvert"
	"github.com/free5gc/nas/nasMessage"
	"github.com/free5gc/nas/nasType"
	"github.com/free5gc/openapi/models"
	smf_context "github.com/free5gc/smf/internal/context"
	"github.com/free5gc/smf/internal/logger"
)

type GSMError struct {
	GSMCause uint8
}

var _ error = &GSMError{}

func (e *GSMError) Error() string {
	return fmt.Sprintf("gsm error cause[%d]", e.GSMCause)
}

func HandlePDUSessionEstablishmentRequest(
	smCtx *smf_context.SMContext, req *nasMessage.PDUSessionEstablishmentRequest,
) error {
	// Retrieve PDUSessionID
	smCtx.PDUSessionID = int32(req.PDUSessionID.GetPDUSessionID())
	logger.GsmLog.Infoln("In HandlePDUSessionEstablishmentRequest")

	// Retrieve PTI (Procedure transaction identity)
	smCtx.Pti = req.GetPTI()

	// Retrieve MaxIntegrityProtectedDataRate of UE for UP Security
	switch req.GetMaximumDataRatePerUEForUserPlaneIntegrityProtectionForUpLink() {
	case 0x00:
		smCtx.MaximumDataRatePerUEForUserPlaneIntegrityProtectionForUpLink = models.
			MaxIntegrityProtectedDataRate__64_KBPS
	case 0xff:
		smCtx.MaximumDataRatePerUEForUserPlaneIntegrityProtectionForUpLink = models.
			MaxIntegrityProtectedDataRate_MAX_UE_RATE
	}
	switch req.GetMaximumDataRatePerUEForUserPlaneIntegrityProtectionForDownLink() {
	case 0x00:
		smCtx.MaximumDataRatePerUEForUserPlaneIntegrityProtectionForDownLink = models.
			MaxIntegrityProtectedDataRate__64_KBPS
	case 0xff:
		smCtx.MaximumDataRatePerUEForUserPlaneIntegrityProtectionForDownLink = models.
			MaxIntegrityProtectedDataRate_MAX_UE_RATE
	}
	// Handle PDUSessionType
	if req.PDUSessionType != nil {
		requestedPDUSessionType := req.PDUSessionType.GetPDUSessionTypeValue()
		if err := smCtx.IsAllowedPDUSessionType(requestedPDUSessionType); err != nil {
			logger.CtxLog.Errorf("%s", err)
			return &GSMError{
				GSMCause: nasMessage.Cause5GSMPDUSessionTypeIPv4OnlyAllowed,
			}
		}
	} else {
		// Set to default supported PDU Session Type
		switch smf_context.GetSelf().SupportedPDUSessionType {
		case "IPv4":
			smCtx.SelectedPDUSessionType = nasMessage.PDUSessionTypeIPv4
		case "IPv6":
			smCtx.SelectedPDUSessionType = nasMessage.PDUSessionTypeIPv6
		case "IPv4v6":
			smCtx.SelectedPDUSessionType = nasMessage.PDUSessionTypeIPv4IPv6
		case "Ethernet":
			smCtx.SelectedPDUSessionType = nasMessage.PDUSessionTypeEthernet
		default:
			smCtx.SelectedPDUSessionType = nasMessage.PDUSessionTypeIPv4
		}
	}

	if req.ExtendedProtocolConfigurationOptions != nil {
		EPCOContents := req.ExtendedProtocolConfigurationOptions.GetExtendedProtocolConfigurationOptionsContents()
		protocolConfigurationOptions := nasConvert.NewProtocolConfigurationOptions()
		unmarshalErr := protocolConfigurationOptions.UnMarshal(EPCOContents)
		if unmarshalErr != nil {
			logger.GsmLog.Errorf("Parsing PCO failed: %s", unmarshalErr)
		}
		logger.GsmLog.Infoln("Protocol Configuration Options")
		logger.GsmLog.Infoln(protocolConfigurationOptions)

		for _, container := range protocolConfigurationOptions.ProtocolOrContainerList {
			logger.GsmLog.Traceln("Container ID: ", container.ProtocolOrContainerID)
			logger.GsmLog.Traceln("Container Length: ", container.LengthOfContents)
			switch container.ProtocolOrContainerID {
			case nasMessage.PCSCFIPv6AddressRequestUL:
				logger.GsmLog.Infoln("Didn't Implement container type PCSCFIPv6AddressRequestUL")
			case nasMessage.IMCNSubsystemSignalingFlagUL:
				logger.GsmLog.Infoln("Didn't Implement container type IMCNSubsystemSignalingFlagUL")
			case nasMessage.DNSServerIPv6AddressRequestUL:
				smCtx.ProtocolConfigurationOptions.DNSIPv6Request = true
			case nasMessage.NotSupportedUL:
				logger.GsmLog.Infoln("Didn't Implement container type NotSupportedUL")
			case nasMessage.MSSupportOfNetworkRequestedBearerControlIndicatorUL:
				logger.GsmLog.Infoln("Didn't Implement container type MSSupportOfNetworkRequestedBearerControlIndicatorUL")
			case nasMessage.DSMIPv6HomeAgentAddressRequestUL:
				logger.GsmLog.Infoln("Didn't Implement container type DSMIPv6HomeAgentAddressRequestUL")
			case nasMessage.DSMIPv6HomeNetworkPrefixRequestUL:
				logger.GsmLog.Infoln("Didn't Implement container type DSMIPv6HomeNetworkPrefixRequestUL")
			case nasMessage.DSMIPv6IPv4HomeAgentAddressRequestUL:
				logger.GsmLog.Infoln("Didn't Implement container type DSMIPv6IPv4HomeAgentAddressRequestUL")
			case nasMessage.IPAddressAllocationViaNASSignallingUL:
				logger.GsmLog.Infoln("Didn't Implement container type IPAddressAllocationViaNASSignallingUL")
			case nasMessage.IPv4AddressAllocationViaDHCPv4UL:
				logger.GsmLog.Infoln("Didn't Implement container type IPv4AddressAllocationViaDHCPv4UL")
			case nasMessage.PCSCFIPv4AddressRequestUL:
				smCtx.ProtocolConfigurationOptions.PCSCFIPv4Request = true
			case nasMessage.DNSServerIPv4AddressRequestUL:
				smCtx.ProtocolConfigurationOptions.DNSIPv4Request = true
			case nasMessage.MSISDNRequestUL:
				logger.GsmLog.Infoln("Didn't Implement container type MSISDNRequestUL")
			case nasMessage.IFOMSupportRequestUL:
				logger.GsmLog.Infoln("Didn't Implement container type IFOMSupportRequestUL")
			case nasMessage.IPv4LinkMTURequestUL:
				smCtx.ProtocolConfigurationOptions.IPv4LinkMTURequest = true
			case nasMessage.MSSupportOfLocalAddressInTFTIndicatorUL:
				logger.GsmLog.Infoln("Didn't Implement container type MSSupportOfLocalAddressInTFTIndicatorUL")
			case nasMessage.PCSCFReSelectionSupportUL:
				logger.GsmLog.Infoln("Didn't Implement container type PCSCFReSelectionSupportUL")
			case nasMessage.NBIFOMRequestIndicatorUL:
				logger.GsmLog.Infoln("Didn't Implement container type NBIFOMRequestIndicatorUL")
			case nasMessage.NBIFOMModeUL:
				logger.GsmLog.Infoln("Didn't Implement container type NBIFOMModeUL")
			case nasMessage.NonIPLinkMTURequestUL:
				logger.GsmLog.Infoln("Didn't Implement container type NonIPLinkMTURequestUL")
			case nasMessage.APNRateControlSupportIndicatorUL:
				logger.GsmLog.Infoln("Didn't Implement container type APNRateControlSupportIndicatorUL")
			case nasMessage.UEStatus3GPPPSDataOffUL:
				logger.GsmLog.Infoln("Didn't Implement container type UEStatus3GPPPSDataOffUL")
			case nasMessage.ReliableDataServiceRequestIndicatorUL:
				logger.GsmLog.Infoln("Didn't Implement container type ReliableDataServiceRequestIndicatorUL")
			case nasMessage.AdditionalAPNRateControlForExceptionDataSupportIndicatorUL:
				logger.GsmLog.Infoln(
					"Didn't Implement container type AdditionalAPNRateControlForExceptionDataSupportIndicatorUL",
				)
			case nasMessage.PDUSessionIDUL:
				logger.GsmLog.Infoln("Didn't Implement container type PDUSessionIDUL")
			case nasMessage.EthernetFramePayloadMTURequestUL:
				logger.GsmLog.Infoln("Didn't Implement container type EthernetFramePayloadMTURequestUL")
			case nasMessage.UnstructuredLinkMTURequestUL:
				logger.GsmLog.Infoln("Didn't Implement container type UnstructuredLinkMTURequestUL")
			case nasMessage.I5GSMCauseValueUL:
				logger.GsmLog.Infoln("Didn't Implement container type 5GSMCauseValueUL")
			case nasMessage.QoSRulesWithTheLengthOfTwoOctetsSupportIndicatorUL:
				logger.GsmLog.Infoln("Didn't Implement container type QoSRulesWithTheLengthOfTwoOctetsSupportIndicatorUL")
			case nasMessage.QoSFlowDescriptionsWithTheLengthOfTwoOctetsSupportIndicatorUL:
				logger.GsmLog.Infoln(
					"Didn't Implement container type QoSFlowDescriptionsWithTheLengthOfTwoOctetsSupportIndicatorUL",
				)
			case nasMessage.LinkControlProtocolUL:
				logger.GsmLog.Infoln("Didn't Implement container type LinkControlProtocolUL")
			case nasMessage.PushAccessControlProtocolUL:
				logger.GsmLog.Infoln("Didn't Implement container type PushAccessControlProtocolUL")
			case nasMessage.ChallengeHandshakeAuthenticationProtocolUL:
				logger.GsmLog.Infoln("Didn't Implement container type ChallengeHandshakeAuthenticationProtocolUL")
			case nasMessage.InternetProtocolControlProtocolUL:
				logger.GsmLog.Infoln("Didn't Implement container type InternetProtocolControlProtocolUL")
			default:
				logger.GsmLog.Infof("Unknown Container ID [%d]", container.ProtocolOrContainerID)
			}
		}
	}
	return nil
}

func HandlePDUSessionReleaseRequest(
	smCtx *smf_context.SMContext, req *nasMessage.PDUSessionReleaseRequest,
) {
	logger.GsmLog.Infof("Handle Pdu Session Release Request")

	// Retrieve PTI (Procedure transaction identity)
	smCtx.Pti = req.GetPTI()
}

func (p *Processor) HandlePDUSessionModificationRequest(
	smCtx *smf_context.SMContext, req *nasMessage.PDUSessionModificationRequest,
) (*nas.Message, error) {
	logger.GsmLog.Infof("Handle Pdu Session Modification Request")

	// Retrieve PTI (Procedure transaction identity)
	smCtx.Pti = req.GetPTI()

	rsp := nas.NewMessage()
	rsp.GsmMessage = nas.NewGsmMessage()
	rsp.GsmHeader.SetMessageType(nas.MsgTypePDUSessionModificationCommand)
	rsp.GsmHeader.SetExtendedProtocolDiscriminator(nasMessage.Epd5GSSessionManagementMessage)
	rsp.PDUSessionModificationCommand = nasMessage.NewPDUSessionModificationCommand(0x00)
	pDUSessionModificationCommand := rsp.PDUSessionModificationCommand
	pDUSessionModificationCommand.SetExtendedProtocolDiscriminator(nasMessage.Epd5GSSessionManagementMessage)
	pDUSessionModificationCommand.SetPDUSessionID(uint8(smCtx.PDUSessionID))
	pDUSessionModificationCommand.SetPTI(smCtx.Pti)
	pDUSessionModificationCommand.SetMessageType(nas.MsgTypePDUSessionModificationCommand)

	reqQoSRules := nasType.QoSRules{}
	reqQoSFlowDescs := nasType.QoSFlowDescs{}

	if req.RequestedQosRules != nil {
		qosRuleBytes := req.GetQoSRules()
		if err := reqQoSRules.UnmarshalBinary(qosRuleBytes); err != nil {
			smCtx.Log.Warning("QoS rule parse failed:", err)
		}
	}

	if req.RequestedQosFlowDescriptions != nil {
		qosFlowDescsBytes := req.GetQoSFlowDescriptions()
		if err := reqQoSFlowDescs.UnmarshalBinary(qosFlowDescsBytes); err != nil {
			smCtx.Log.Warning("QoS flow descriptions parse failed:", err)
		}
	}

	smPolicyDecision, err_ := p.Consumer().SendSMPolicyAssociationUpdateByUERequestModification(
		smCtx, reqQoSRules, reqQoSFlowDescs)
	if err_ != nil {
		return nil, fmt.Errorf("sm policy update failed: %s", err_)
	}

	// Update SessionRule from decision
	if errApplySessionRules := smCtx.ApplySessionRules(smPolicyDecision); errApplySessionRules != nil {
		return nil, fmt.Errorf("PDUSessionSMContextCreate err: %v", errApplySessionRules)
	}

	if errApplyPccRules := smCtx.ApplyPccRules(smPolicyDecision); errApplyPccRules != nil {
		smCtx.Log.Errorf("apply sm policy decision error: %+v", errApplyPccRules)
	}

	authQoSRules := nasType.QoSRules{}
	authQoSFlowDesc := reqQoSFlowDescs

	for id := range smPolicyDecision.PccRules {
		// get op code from request
		opCode := reqQoSRules[0].Operation
		// build nas Qos Rule
		pccRule := smCtx.PCCRules[id]
		rule, err := pccRule.BuildNasQoSRule(smCtx, opCode)
		if err != nil {
			return nil, err
		}

		authQoSRules = append(authQoSRules, *rule)
	}

	if len(authQoSRules) > 0 {
		if buf, err := authQoSRules.MarshalBinary(); err != nil {
			return nil, err
		} else {
			pDUSessionModificationCommand.AuthorizedQosRules = nasType.NewAuthorizedQosRules(0x7A)
			pDUSessionModificationCommand.AuthorizedQosRules.SetLen(uint16(len(buf)))
			pDUSessionModificationCommand.AuthorizedQosRules.SetQosRule(buf)
		}
	}

	if len(authQoSFlowDesc) > 0 {
		if buf, err := authQoSFlowDesc.MarshalBinary(); err != nil {
			return nil, err
		} else {
			pDUSessionModificationCommand.AuthorizedQosFlowDescriptions = nasType.NewAuthorizedQosFlowDescriptions(0x79)
			pDUSessionModificationCommand.AuthorizedQosFlowDescriptions.SetLen(uint16(len(buf)))
			pDUSessionModificationCommand.AuthorizedQosFlowDescriptions.SetQoSFlowDescriptions(buf)
		}
	}

	return rsp, nil
}
