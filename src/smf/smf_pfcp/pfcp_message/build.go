package pfcp_message

import (
	"net"

	"free5gc/lib/pfcp"
	"free5gc/lib/pfcp/pfcpType"
	"free5gc/src/smf/smf_context"
	"free5gc/src/smf/smf_pfcp/pfcp_udp"
)

func BuildPfcpAssociationSetupRequest() (pfcp.PFCPAssociationSetupRequest, error) {
	msg := pfcp.PFCPAssociationSetupRequest{}

	msg.NodeID = &smf_context.SMF_Self().CPNodeID

	msg.RecoveryTimeStamp = &pfcpType.RecoveryTimeStamp{
		RecoveryTimeStamp: pfcp_udp.ServerStartTime,
	}

	msg.CPFunctionFeatures = &pfcpType.CPFunctionFeatures{
		SupportedFeatures: 0,
	}

	return msg, nil
}

func BuildPfcpAssociationSetupResponse(cause pfcpType.Cause) (pfcp.PFCPAssociationSetupResponse, error) {
	msg := pfcp.PFCPAssociationSetupResponse{}

	msg.NodeID = &smf_context.SMF_Self().CPNodeID

	msg.Cause = &cause

	msg.RecoveryTimeStamp = &pfcpType.RecoveryTimeStamp{
		RecoveryTimeStamp: pfcp_udp.ServerStartTime,
	}

	msg.CPFunctionFeatures = &pfcpType.CPFunctionFeatures{
		SupportedFeatures: 0,
	}

	return msg, nil
}

func BuildPfcpAssociationReleaseRequest() (pfcp.PFCPAssociationReleaseRequest, error) {
	msg := pfcp.PFCPAssociationReleaseRequest{}

	msg.NodeID = &smf_context.SMF_Self().CPNodeID

	return msg, nil
}

func BuildPfcpAssociationReleaseResponse(cause pfcpType.Cause) (pfcp.PFCPAssociationReleaseResponse, error) {
	msg := pfcp.PFCPAssociationReleaseResponse{}

	msg.NodeID = &smf_context.SMF_Self().CPNodeID

	msg.Cause = &cause

	return msg, nil
}

func pdrToCreatePDR(pdr *smf_context.PDR) *pfcp.CreatePDR {
	createPDR := new(pfcp.CreatePDR)

	createPDR.PDRID = new(pfcpType.PacketDetectionRuleID)
	createPDR.PDRID.RuleId = pdr.PDRID

	createPDR.Precedence = new(pfcpType.Precedence)
	createPDR.Precedence.PrecedenceValue = pdr.Precedence

	createPDR.PDI = &pfcp.PDI{
		SourceInterface: &pdr.PDI.SourceInterface,
		LocalFTEID:      &pdr.PDI.LocalFTeid,
		NetworkInstance: &pdr.PDI.NetworkInstance,
		UEIPAddress:     pdr.PDI.UEIPAddress,
	}

	createPDR.OuterHeaderRemoval = pdr.OuterHeaderRemoval

	createPDR.FARID = &pfcpType.FARID{
		FarIdValue: pdr.FAR.FARID,
	}

	return createPDR
}

func farToCreateFAR(far *smf_context.FAR) *pfcp.CreateFAR {
	createFAR := new(pfcp.CreateFAR)

	createFAR.FARID = new(pfcpType.FARID)
	createFAR.FARID.FarIdValue = far.FARID

	createFAR.ApplyAction = new(pfcpType.ApplyAction)
	createFAR.ApplyAction.Forw = true

	if far.BAR != nil {
		createFAR.BARID = new(pfcpType.BARID)
		createFAR.BARID.BarIdValue = far.BAR.BARID
	}

	createFAR.ForwardingParameters = new(pfcp.ForwardingParametersIEInFAR)
	createFAR.ForwardingParameters.DestinationInterface = &far.ForwardingParameters.DestinationInterface
	createFAR.ForwardingParameters.NetworkInstance = &far.ForwardingParameters.NetworkInstance
	createFAR.ForwardingParameters.OuterHeaderCreation = far.ForwardingParameters.OuterHeaderCreation

	return createFAR
}

func barToCreateBAR(bar *smf_context.BAR) *pfcp.CreateBAR {

	createBAR := new(pfcp.CreateBAR)

	createBAR.BARID = new(pfcpType.BARID)
	createBAR.BARID.BarIdValue = bar.BARID

	createBAR.DownlinkDataNotificationDelay = new(pfcpType.DownlinkDataNotificationDelay)

	//createBAR.SuggestedBufferingPacketsCount = new(pfcpType.SuggestedBufferingPacketsCount)

	return createBAR
}

func pdrToUpdatePDR(pdr *smf_context.PDR) *pfcp.UpdatePDR {
	updatePDR := new(pfcp.UpdatePDR)

	updatePDR.PDRID = new(pfcpType.PacketDetectionRuleID)
	updatePDR.PDRID.RuleId = pdr.PDRID

	updatePDR.Precedence = new(pfcpType.Precedence)
	updatePDR.Precedence.PrecedenceValue = pdr.Precedence

	updatePDR.PDI = &pfcp.PDI{
		SourceInterface: &pdr.PDI.SourceInterface,
		LocalFTEID:      &pdr.PDI.LocalFTeid,
		NetworkInstance: &pdr.PDI.NetworkInstance,
		UEIPAddress:     pdr.PDI.UEIPAddress,
	}

	updatePDR.OuterHeaderRemoval = pdr.OuterHeaderRemoval

	updatePDR.FARID = &pfcpType.FARID{
		FarIdValue: pdr.FAR.FARID,
	}

	updatePDR.FARID = &pfcpType.FARID{
		FarIdValue: pdr.FAR.FARID,
	}

	return updatePDR
}

func farToUpdateFAR(far *smf_context.FAR) *pfcp.UpdateFAR {
	updateFAR := new(pfcp.UpdateFAR)

	updateFAR.FARID = new(pfcpType.FARID)
	updateFAR.FARID.FarIdValue = far.FARID

	if far.BAR != nil {
		updateFAR.BARID = new(pfcpType.BARID)
		updateFAR.BARID.BarIdValue = far.BAR.BARID
	}

	updateFAR.ApplyAction = new(pfcpType.ApplyAction)
	updateFAR.ApplyAction.Forw = far.ApplyAction.Forw
	updateFAR.ApplyAction.Buff = far.ApplyAction.Buff
	updateFAR.ApplyAction.Nocp = far.ApplyAction.Nocp
	updateFAR.ApplyAction.Dupl = far.ApplyAction.Dupl
	updateFAR.ApplyAction.Drop = far.ApplyAction.Drop

	updateFAR.UpdateForwardingParameters = new(pfcp.UpdateForwardingParametersIEInFAR)
	updateFAR.UpdateForwardingParameters.DestinationInterface = &far.ForwardingParameters.DestinationInterface
	updateFAR.UpdateForwardingParameters.NetworkInstance = &far.ForwardingParameters.NetworkInstance
	updateFAR.UpdateForwardingParameters.OuterHeaderCreation = far.ForwardingParameters.OuterHeaderCreation

	return updateFAR
}

// TODO: Replace dummy value in PFCP message
func BuildPfcpSessionEstablishmentRequest(smContext *smf_context.SMContext) (pfcp.PFCPSessionEstablishmentRequest, error) {
	msg := pfcp.PFCPSessionEstablishmentRequest{}

	msg.NodeID = &smf_context.SMF_Self().CPNodeID

	isv4 := smf_context.SMF_Self().CPNodeID.NodeIdType == 0
	msg.CPFSEID = &pfcpType.FSEID{
		V4:          isv4,
		V6:          !isv4,
		Seid:        smContext.LocalSEID,
		Ipv4Address: smf_context.SMF_Self().CPNodeID.NodeIdValue,
	}

	msg.CreatePDR = make([]*pfcp.CreatePDR, 0, 2)
	msg.CreateFAR = make([]*pfcp.CreateFAR, 0, 2)

	msg.CreatePDR = append(msg.CreatePDR, pdrToCreatePDR(smContext.Tunnel.ULPDR))
	msg.CreateFAR = append(msg.CreateFAR, farToCreateFAR(smContext.Tunnel.ULPDR.FAR))

	msg.PDNType = &pfcpType.PDNType{
		PdnType: pfcpType.PDNTypeIpv4,
	}

	return msg, nil
}

func BuildPfcpSessionEstablishmentResponse() (pfcp.PFCPSessionEstablishmentResponse, error) {
	msg := pfcp.PFCPSessionEstablishmentResponse{}

	msg.NodeID = &smf_context.SMF_Self().CPNodeID

	msg.Cause = &pfcpType.Cause{
		CauseValue: pfcpType.CauseRequestAccepted,
	}

	msg.OffendingIE = &pfcpType.OffendingIE{
		TypeOfOffendingIe: 12345,
	}

	msg.UPFSEID = &pfcpType.FSEID{
		V4:          true,
		V6:          false, //;
		Seid:        123456789123456789,
		Ipv4Address: net.ParseIP("192.168.1.1").To4(),
	}

	msg.CreatedPDR = &pfcp.CreatedPDR{
		PDRID: &pfcpType.PacketDetectionRuleID{
			RuleId: 256,
		},
		LocalFTEID: &pfcpType.FTEID{
			Chid:        false,
			Ch:          false,
			V6:          false,
			V4:          true,
			Teid:        12345,
			Ipv4Address: net.ParseIP("192.168.1.1").To4(),
		},
	}

	return msg, nil
}

// TODO: Replace dummy value in PFCP message
func BuildPfcpSessionModificationRequest(smContext *smf_context.SMContext, pdr_list []*smf_context.PDR, far_list []*smf_context.FAR, bar_list []*smf_context.BAR) (pfcp.PFCPSessionModificationRequest, error) {
	msg := pfcp.PFCPSessionModificationRequest{}

	msg.UpdatePDR = make([]*pfcp.UpdatePDR, 0, 2)
	msg.UpdateFAR = make([]*pfcp.UpdateFAR, 0, 2)

	msg.CPFSEID = &pfcpType.FSEID{
		V4:          true,
		V6:          false,
		Seid:        smContext.LocalSEID,
		Ipv4Address: smf_context.SMF_Self().CPNodeID.NodeIdValue,
	}

	for _, pdr := range pdr_list {
		switch pdr.State {
		case smf_context.RULE_INITIAL:
			msg.CreatePDR = append(msg.CreatePDR, pdrToCreatePDR(pdr))
		case smf_context.RULE_UPDATE:
			msg.UpdatePDR = append(msg.UpdatePDR, pdrToUpdatePDR(pdr))
		}
	}

	for _, far := range far_list {
		switch far.State {
		case smf_context.RULE_INITIAL:
			msg.CreateFAR = append(msg.CreateFAR, farToCreateFAR(far))
		case smf_context.RULE_UPDATE:
			msg.UpdateFAR = append(msg.UpdateFAR, farToUpdateFAR(far))
		}
	}

	for _, bar := range bar_list {
		switch bar.State {
		case smf_context.RULE_INITIAL:
			msg.CreateBAR = append(msg.CreateBAR, barToCreateBAR(bar))
		}
	}

	return msg, nil
}

// TODO: Replace dummy value in PFCP message
func BuildPfcpSessionModificationResponse() (pfcp.PFCPSessionModificationResponse, error) {
	msg := pfcp.PFCPSessionModificationResponse{}

	msg.Cause = &pfcpType.Cause{
		CauseValue: pfcpType.CauseRequestAccepted,
	}

	msg.OffendingIE = &pfcpType.OffendingIE{
		TypeOfOffendingIe: 12345,
	}

	msg.CreatedPDR = &pfcp.CreatedPDR{
		PDRID: &pfcpType.PacketDetectionRuleID{
			RuleId: 256,
		},
		LocalFTEID: &pfcpType.FTEID{
			Chid:        false,
			Ch:          false,
			V6:          false,
			V4:          true,
			Teid:        12345,
			Ipv4Address: net.ParseIP("192.168.1.1").To4(),
		},
	}

	return msg, nil
}

func BuildPfcpSessionDeletionRequest() (pfcp.PFCPSessionDeletionRequest, error) {
	msg := pfcp.PFCPSessionDeletionRequest{}

	return msg, nil
}

// TODO: Replace dummy value in PFCP message
func BuildPfcpSessionDeletionResponse() (pfcp.PFCPSessionDeletionResponse, error) {
	msg := pfcp.PFCPSessionDeletionResponse{}

	msg.Cause = &pfcpType.Cause{
		CauseValue: pfcpType.CauseRequestAccepted,
	}

	msg.OffendingIE = &pfcpType.OffendingIE{
		TypeOfOffendingIe: 12345,
	}

	return msg, nil
}

func BuildPfcpSessionReportResponse(cause pfcpType.Cause) (pfcp.PFCPSessionReportResponse, error) {
	msg := pfcp.PFCPSessionReportResponse{}

	msg.Cause = &cause

	return msg, nil
}
