package context

import (
	"encoding/binary"
	"errors"
	"fmt"

	"github.com/free5gc/aper"
	"github.com/free5gc/ngap/ngapType"
	"github.com/free5gc/openapi/models"
	"github.com/free5gc/pfcp/pfcpType"
	"github.com/free5gc/smf/internal/logger"
)

func strNgapCause(cause *ngapType.Cause) string {
	ret := ""
	switch cause.Present {
	case ngapType.CausePresentRadioNetwork:
		ret = fmt.Sprintf("Cause by RadioNetwork[%d]",
			cause.RadioNetwork.Value)
	case ngapType.CausePresentTransport:
		ret = fmt.Sprintf("Cause by Transport[%d]",
			cause.Transport.Value)
	case ngapType.CausePresentNas:
		ret = fmt.Sprintf("Cause by NAS[%d]",
			cause.Nas.Value)
	case ngapType.CausePresentProtocol:
		ret = fmt.Sprintf("Cause by Protocol[%d]",
			cause.Protocol.Value)
	case ngapType.CausePresentMisc:
		ret = fmt.Sprintf("Cause by Protocol[%d]",
			cause.Misc.Value)
	case ngapType.CausePresentChoiceExtensions:
		ret = fmt.Sprintf("Cause by Protocol[%v]",
			cause.ChoiceExtensions)
	default:
		ret = "Cause [unspecific]"
	}

	return ret
}

func HandlePDUSessionResourceSetupResponseTransfer(b []byte, ctx *SMContext) error {
	resourceSetupResponseTransfer := ngapType.PDUSessionResourceSetupResponseTransfer{}

	err := aper.UnmarshalWithParams(b, &resourceSetupResponseTransfer, "valueExt")
	if err != nil {
		return err
	}

	QosFlowPerTNLInformation := resourceSetupResponseTransfer.DLQosFlowPerTNLInformation
	var DCQosFlowPerTNLInformationItem ngapType.QosFlowPerTNLInformationItem
	DCQosFlowPerTNLInformation := resourceSetupResponseTransfer.AdditionalDLQosFlowPerTNLInformation
	if DCQosFlowPerTNLInformation != nil && len(DCQosFlowPerTNLInformation.List) > 0 {
		ctx.NrdcIndicator = true
		DCQosFlowPerTNLInformationItem = DCQosFlowPerTNLInformation.List[0]
	}

	if QosFlowPerTNLInformation.UPTransportLayerInformation.Present !=
		ngapType.UPTransportLayerInformationPresentGTPTunnel {
		return errors.New("resourceSetupResponseTransfer.QosFlowPerTNLInformation.UPTransportLayerInformation.Present")
	}
	if ctx.NrdcIndicator && DCQosFlowPerTNLInformationItem.QosFlowPerTNLInformation.UPTransportLayerInformation.Present !=
		ngapType.UPTransportLayerInformationPresentGTPTunnel {
		return errors.New(
			"resourceSetupResponseTransfer.AdditionalQosFlowPerTNLInformation." +
				"QosFlowPerTNLInformation.UPTransportLayerInformation.Present")
	}

	GTPTunnel := QosFlowPerTNLInformation.UPTransportLayerInformation.GTPTunnel
	DCGTPTunnel := &ngapType.GTPTunnel{}
	if ctx.NrdcIndicator {
		DCGTPTunnel = DCQosFlowPerTNLInformationItem.QosFlowPerTNLInformation.UPTransportLayerInformation.GTPTunnel
	}

	ctx.Tunnel.UpdateANInformation(
		GTPTunnel.TransportLayerAddress.Value.Bytes,
		binary.BigEndian.Uint32(GTPTunnel.GTPTEID.Value))
	if ctx.NrdcIndicator {
		ctx.DCTunnel.UpdateANInformation(
			DCGTPTunnel.TransportLayerAddress.Value.Bytes,
			binary.BigEndian.Uint32(DCGTPTunnel.GTPTEID.Value))
	}

	ctx.UpCnxState = models.UpCnxState_ACTIVATED
	for _, qos := range ctx.AdditonalQosFlows {
		qos.State = QoSFlowSet
	}
	return nil
}

func HandlePDUSessionResourceModifyResponseTransfer(b []byte, ctx *SMContext) error {
	resourceModifyResponseTransfer := ngapType.PDUSessionResourceModifyResponseTransfer{}

	err := aper.UnmarshalWithParams(b, &resourceModifyResponseTransfer, "valueExt")
	if err != nil {
		return err
	}

	if DLInfo := resourceModifyResponseTransfer.DLNGUUPTNLInformation; DLInfo != nil {
		GTPTunnel := DLInfo.GTPTunnel

		ctx.Tunnel.UpdateANInformation(
			GTPTunnel.TransportLayerAddress.Value.Bytes,
			binary.BigEndian.Uint32(GTPTunnel.GTPTEID.Value))
	}

	if qosInfoList := resourceModifyResponseTransfer.QosFlowAddOrModifyResponseList; qosInfoList != nil {
		for _, item := range qosInfoList.List {
			qfi := uint8(item.QosFlowIdentifier.Value)
			ctx.AdditonalQosFlows[qfi].State = QoSFlowSet
		}
	}

	if qosFailedInfoList := resourceModifyResponseTransfer.QosFlowFailedToAddOrModifyList; qosFailedInfoList != nil {
		for _, item := range qosFailedInfoList.List {
			qfi := uint8(item.QosFlowIdentifier.Value)
			logger.PduSessLog.Warnf("PDU Session Resource Modify QFI[%d] %s",
				qfi, strNgapCause(&item.Cause))

			ctx.AdditonalQosFlows[qfi].State = QoSFlowUnset
		}
	}

	return nil
}

func HandlePDUSessionResourceSetupUnsuccessfulTransfer(b []byte, ctx *SMContext) error {
	resourceSetupUnsuccessfulTransfer := ngapType.PDUSessionResourceSetupUnsuccessfulTransfer{}

	err := aper.UnmarshalWithParams(b, &resourceSetupUnsuccessfulTransfer, "valueExt")
	if err != nil {
		return err
	}

	switch resourceSetupUnsuccessfulTransfer.Cause.Present {
	case ngapType.CausePresentRadioNetwork:
		logger.PduSessLog.Warnf("PDU Session Resource Setup Unsuccessful by RadioNetwork[%d]",
			resourceSetupUnsuccessfulTransfer.Cause.RadioNetwork.Value)
	case ngapType.CausePresentTransport:
		logger.PduSessLog.Warnf("PDU Session Resource Setup Unsuccessful by Transport[%d]",
			resourceSetupUnsuccessfulTransfer.Cause.Transport.Value)
	case ngapType.CausePresentNas:
		logger.PduSessLog.Warnf("PDU Session Resource Setup Unsuccessful by NAS[%d]",
			resourceSetupUnsuccessfulTransfer.Cause.Nas.Value)
	case ngapType.CausePresentProtocol:
		logger.PduSessLog.Warnf("PDU Session Resource Setup Unsuccessful by Protocol[%d]",
			resourceSetupUnsuccessfulTransfer.Cause.Protocol.Value)
	case ngapType.CausePresentMisc:
		logger.PduSessLog.Warnf("PDU Session Resource Setup Unsuccessful by Protocol[%d]",
			resourceSetupUnsuccessfulTransfer.Cause.Misc.Value)
	case ngapType.CausePresentChoiceExtensions:
		logger.PduSessLog.Warnf("PDU Session Resource Setup Unsuccessful by Protocol[%v]",
			resourceSetupUnsuccessfulTransfer.Cause.ChoiceExtensions)
	}

	ctx.UpCnxState = models.UpCnxState_ACTIVATING

	return nil
}

func HandlePathSwitchRequestTransfer(b []byte, ctx *SMContext) error {
	pathSwitchRequestTransfer := ngapType.PathSwitchRequestTransfer{}

	if err := aper.UnmarshalWithParams(b, &pathSwitchRequestTransfer, "valueExt"); err != nil {
		return err
	}

	if pathSwitchRequestTransfer.DLNGUUPTNLInformation.Present != ngapType.UPTransportLayerInformationPresentGTPTunnel {
		return errors.New("pathSwitchRequestTransfer.DLNGUUPTNLInformation.Present")
	}

	GTPTunnel := pathSwitchRequestTransfer.DLNGUUPTNLInformation.GTPTunnel

	ctx.Tunnel.UpdateANInformation(
		GTPTunnel.TransportLayerAddress.Value.Bytes,
		binary.BigEndian.Uint32(GTPTunnel.GTPTEID.Value))

	ctx.UpSecurityFromPathSwitchRequestSameAsLocalStored = true

	// Verify whether UP security in PathSwitchRequest same as SMF locally stored or not TS 33.501 6.6.1
	if ctx.UpSecurity != nil && pathSwitchRequestTransfer.UserPlaneSecurityInformation != nil {
		rcvSecurityIndication := pathSwitchRequestTransfer.UserPlaneSecurityInformation.SecurityIndication
		rcvUpSecurity := new(models.UpSecurity)
		switch rcvSecurityIndication.IntegrityProtectionIndication.Value {
		case ngapType.IntegrityProtectionIndicationPresentRequired:
			rcvUpSecurity.UpIntegr = models.UpIntegrity_REQUIRED
		case ngapType.IntegrityProtectionIndicationPresentPreferred:
			rcvUpSecurity.UpIntegr = models.UpIntegrity_PREFERRED
		case ngapType.IntegrityProtectionIndicationPresentNotNeeded:
			rcvUpSecurity.UpIntegr = models.UpIntegrity_NOT_NEEDED
		}
		switch rcvSecurityIndication.ConfidentialityProtectionIndication.Value {
		case ngapType.ConfidentialityProtectionIndicationPresentRequired:
			rcvUpSecurity.UpConfid = models.UpConfidentiality_REQUIRED
		case ngapType.ConfidentialityProtectionIndicationPresentPreferred:
			rcvUpSecurity.UpConfid = models.UpConfidentiality_PREFERRED
		case ngapType.ConfidentialityProtectionIndicationPresentNotNeeded:
			rcvUpSecurity.UpConfid = models.UpConfidentiality_NOT_NEEDED
		}

		if rcvUpSecurity.UpIntegr != ctx.UpSecurity.UpIntegr ||
			rcvUpSecurity.UpConfid != ctx.UpSecurity.UpConfid {
			ctx.UpSecurityFromPathSwitchRequestSameAsLocalStored = false

			// SMF shall support logging capabilities for this mismatch event TS 33.501 6.6.1
			logger.PduSessLog.Warnf("Received UP security policy mismatch from SMF locally stored")
		}
	}

	return nil
}

func HandlePathSwitchRequestSetupFailedTransfer(b []byte, ctx *SMContext) error {
	pathSwitchRequestSetupFailedTransfer := ngapType.PathSwitchRequestSetupFailedTransfer{}

	err := aper.UnmarshalWithParams(b, &pathSwitchRequestSetupFailedTransfer, "valueExt")
	if err != nil {
		return err
	}

	// TODO: finish handler
	return nil
}

func HandleHandoverRequiredTransfer(b []byte, ctx *SMContext) error {
	handoverRequiredTransfer := ngapType.HandoverRequiredTransfer{}

	err := aper.UnmarshalWithParams(b, &handoverRequiredTransfer, "valueExt")

	directForwardingPath := handoverRequiredTransfer.DirectForwardingPathAvailability
	if directForwardingPath != nil {
		logger.PduSessLog.Infoln("Direct Forwarding Path Available")
		ctx.DLForwardingType = DirectForwarding
	} else {
		logger.PduSessLog.Infoln("Direct Forwarding Path Unavailable")
		ctx.DLForwardingType = IndirectForwarding
	}

	if err != nil {
		return err
	}
	return nil
}

func HandleHandoverRequestAcknowledgeTransfer(b []byte, ctx *SMContext) error {
	handoverRequestAcknowledgeTransfer := ngapType.HandoverRequestAcknowledgeTransfer{}

	err := aper.UnmarshalWithParams(b, &handoverRequestAcknowledgeTransfer, "valueExt")
	if err != nil {
		return err
	}

	DLNGUUPGTPTunnel := handoverRequestAcknowledgeTransfer.DLNGUUPTNLInformation.GTPTunnel

	ctx.Tunnel.UpdateANInformation(
		DLNGUUPGTPTunnel.TransportLayerAddress.Value.Bytes,
		binary.BigEndian.Uint32(DLNGUUPGTPTunnel.GTPTEID.Value))

	DLForwardingInfo := handoverRequestAcknowledgeTransfer.DLForwardingUPTNLInformation

	if DLForwardingInfo == nil {
		ctx.DLForwardingType = NoForwarding
		logger.PduSessLog.Warnf("Handle HandoverRequestAcknowledgeTransfer warned: %+v", "DL Forwarding Info not provision")
		return nil
	}

	switch ctx.DLForwardingType {
	case IndirectForwarding:
		DLForwardingGTPTunnel := DLForwardingInfo.GTPTunnel

		ctx.IndirectForwardingTunnel = NewDataPath()
		ctx.IndirectForwardingTunnel.FirstDPNode = NewDataPathNode()
		ctx.IndirectForwardingTunnel.FirstDPNode.UPF = ctx.Tunnel.DataPathPool.GetDefaultPath().FirstDPNode.UPF
		ctx.IndirectForwardingTunnel.FirstDPNode.UpLinkTunnel = &GTPTunnel{}

		ANUPF := ctx.IndirectForwardingTunnel.FirstDPNode.UPF

		var indirectFowardingPDR *PDR

		if pdr, errAddPDR := ANUPF.AddPDR(); errAddPDR != nil {
			return errAddPDR
		} else {
			indirectFowardingPDR = pdr
		}

		originPDR := ctx.Tunnel.DataPathPool.GetDefaultPath().FirstDPNode.UpLinkTunnel.PDR

		if teid, errGenerateTEID := GenerateTEID(); errGenerateTEID != nil {
			return errGenerateTEID
		} else {
			ctx.IndirectForwardingTunnel.FirstDPNode.UpLinkTunnel.TEID = teid
			ctx.IndirectForwardingTunnel.FirstDPNode.UpLinkTunnel.PDR = indirectFowardingPDR
			indirectFowardingPDR.PDI.LocalFTeid = &pfcpType.FTEID{
				V4:          originPDR.PDI.LocalFTeid.V4,
				Teid:        ctx.IndirectForwardingTunnel.FirstDPNode.UpLinkTunnel.TEID,
				Ipv4Address: originPDR.PDI.LocalFTeid.Ipv4Address,
			}
			indirectFowardingPDR.OuterHeaderRemoval = &pfcpType.OuterHeaderRemoval{
				OuterHeaderRemovalDescription: pfcpType.OuterHeaderRemovalGtpUUdpIpv4,
			}

			indirectFowardingPDR.FAR.ApplyAction = pfcpType.ApplyAction{
				Forw: true,
			}
			indirectFowardingPDR.FAR.ForwardingParameters = &ForwardingParameters{
				DestinationInterface: pfcpType.DestinationInterface{
					InterfaceValue: pfcpType.DestinationInterfaceAccess,
				},
				OuterHeaderCreation: &pfcpType.OuterHeaderCreation{
					OuterHeaderCreationDescription: pfcpType.OuterHeaderCreationGtpUUdpIpv4,
					Teid:                           binary.BigEndian.Uint32(DLForwardingGTPTunnel.GTPTEID.Value),
					Ipv4Address:                    DLForwardingGTPTunnel.TransportLayerAddress.Value.Bytes,
				},
			}
		}
	case DirectForwarding:
		ctx.DLDirectForwardingTunnel = DLForwardingInfo
	}

	return nil
}
