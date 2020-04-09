package smf_producer

import (
	"free5gc/lib/pfcp/pfcpType"
	"free5gc/lib/pfcp/pfcpUdp"
	"free5gc/src/smf/logger"
	"free5gc/src/smf/smf_context"
	"free5gc/src/smf/smf_pfcp/pfcp_message"
	"net"
)

func SetUpUplinkUserPlane(root *smf_context.DataPathNode, smContext *smf_context.SMContext) {

	visited := make(map[*smf_context.DataPathNode]bool)
	AllocateUpLinkPDRandTEID(root, smContext, visited)

	for node, _ := range visited {
		visited[node] = false
	}

	// SendUplinkPFCPRule(root, smContext, visited)
}

func SetUpDownLinkUserPlane(root *smf_context.DataPathNode, smContext *smf_context.SMContext) {

	visited := make(map[*smf_context.DataPathNode]bool)
	AllocateDownLinkPDR(root, smContext, visited)

	for node, _ := range visited {
		visited[node] = false
	}

	AllocateDownLinkTEID(root, smContext, visited)

	for node, _ := range visited {
		visited[node] = false
	}

	SendDownLinkPFCPRule(root, smContext, visited)
}

func AllocateUpLinkPDRandTEID(node *smf_context.DataPathNode, smContext *smf_context.SMContext, visited map[*smf_context.DataPathNode]bool) {

	if !visited[node] {
		visited[node] = true
	}

	var err error
	upLink := node.DataPathToAN

	teid, err := node.UPF.GenerateTEID()

	if err != nil {
		logger.PduSessLog.Error(err)
		return
	}

	upLink.UpLinkPDR, err = node.UPF.AddPDR()
	if err != nil {
		logger.PduSessLog.Error(err)
		return
	}
	upLink.UpLinkPDR.PDRID = upLink.UpLinkPDR.PDRID - 2
	upLink.UpLinkPDR.FAR.FARID = upLink.UpLinkPDR.FAR.FARID - 2

	upLink.UpLinkPDR.Precedence = 32
	upLink.UpLinkPDR.PDI = smf_context.PDI{
		SourceInterface: pfcpType.SourceInterface{
			//Todo:
			//Have to change source interface for different upf
			InterfaceValue: pfcpType.SourceInterfaceAccess,
		},
		LocalFTeid: &pfcpType.FTEID{
			V4:          true,
			Teid:        teid - 2,
			Ipv4Address: node.UPF.UPIPInfo.Ipv4Address,
		},
		NetworkInstance: []byte(smContext.Dnn),
		UEIPAddress: &pfcpType.UEIPAddress{
			V4:          true,
			Ipv4Address: smContext.PDUAddress.To4(),
		},
	}
	upLink.UpLinkPDR.OuterHeaderRemoval = new(pfcpType.OuterHeaderRemoval)
	upLink.UpLinkPDR.OuterHeaderRemoval.OuterHeaderRemovalDescription = pfcpType.OuterHeaderRemovalGtpUUdpIpv4
	upLink.UpLinkPDR.State = smf_context.RULE_INITIAL

	upLink.UpLinkPDR.FAR.ApplyAction.Forw = true
	upLink.UpLinkPDR.FAR.State = smf_context.RULE_INITIAL
	upLink.UpLinkPDR.FAR.ForwardingParameters = &smf_context.ForwardingParameters{
		DestinationInterface: pfcpType.DestinationInterface{
			InterfaceValue: pfcpType.DestinationInterfaceCore,
		},
		NetworkInstance: []byte(smContext.Dnn),
	}

	parent := node.GetParent()
	if parent != nil {

		parentUpLinkFAR := parent.GetUpLinkFAR()

		parentUpLinkFAR.ForwardingParameters.OuterHeaderCreation = new(pfcpType.OuterHeaderCreation)
		parentUpLinkFAR.ForwardingParameters.OuterHeaderCreation.OuterHeaderCreationDescription = pfcpType.OuterHeaderCreationGtpUUdpIpv4
		parentUpLinkFAR.ForwardingParameters.OuterHeaderCreation.Teid = uint32(teid - 2)
		parentUpLinkFAR.ForwardingParameters.OuterHeaderCreation.Ipv4Address = node.UPF.UPIPInfo.Ipv4Address
	}

	for _, upf_link := range node.DataPathToDN {

		child := upf_link.To
		if !visited[child] && child.InUse {
			AllocateUpLinkPDRandTEID(child, smContext, visited)
		}

	}

}

func AllocateDownLinkPDR(node *smf_context.DataPathNode, smContext *smf_context.SMContext, visited map[*smf_context.DataPathNode]bool) {
	var err error
	var teid uint32

	if !visited[node] {
		visited[node] = true
	}

	for _, downLink := range node.DataPathToDN {

		downLink.DownLinkPDR, err = node.UPF.AddPDR()

		if err != nil {
			logger.PduSessLog.Error(err)
		}

		teid, err = node.UPF.GenerateTEID()

		if err != nil {
			logger.PduSessLog.Error(err)
		}

		downLink.DownLinkPDR.Precedence = 32
		downLink.DownLinkPDR.PDI = smf_context.PDI{
			SourceInterface: pfcpType.SourceInterface{
				//Todo:
				//Have to change source interface for different upf
				InterfaceValue: pfcpType.SourceInterfaceAccess,
			},
			LocalFTeid: &pfcpType.FTEID{
				V4:          true,
				Teid:        teid,
				Ipv4Address: node.UPF.UPIPInfo.Ipv4Address,
			},
			NetworkInstance: []byte(smContext.Dnn),
			UEIPAddress: &pfcpType.UEIPAddress{
				V4:          true,
				Ipv4Address: smContext.PDUAddress.To4(),
			},
		}

		downLink.DownLinkPDR.OuterHeaderRemoval = new(pfcpType.OuterHeaderRemoval)
		downLink.DownLinkPDR.OuterHeaderRemoval.OuterHeaderRemovalDescription = pfcpType.OuterHeaderRemovalGtpUUdpIpv4
		downLink.DownLinkPDR.State = smf_context.RULE_INITIAL

		downLink.DownLinkPDR.FAR.ApplyAction.Forw = true
		downLink.DownLinkPDR.FAR.State = smf_context.RULE_INITIAL
		downLink.DownLinkPDR.FAR.ForwardingParameters = &smf_context.ForwardingParameters{
			DestinationInterface: pfcpType.DestinationInterface{
				InterfaceValue: pfcpType.DestinationInterfaceCore,
			},
			NetworkInstance: []byte(smContext.Dnn),
		}

	}

	for _, upf_link := range node.DataPathToDN {

		child := upf_link.To
		if !visited[child] && child.InUse {
			AllocateDownLinkPDR(child, smContext, visited)
		}

	}

	if node.IsAnchorUPF() {

		downLink := node.DLDataPathLinkForPSA
		downLink.DownLinkPDR, err = node.UPF.AddPDR()
		if err != nil {
			logger.PduSessLog.Error(err)
		}
		downLink.DownLinkPDR.Precedence = 32
		downLink.DownLinkPDR.PDI = smf_context.PDI{
			SourceInterface: pfcpType.SourceInterface{
				//Todo:
				//Have to change source interface for different upf
				InterfaceValue: pfcpType.SourceInterfaceAccess,
			},
			LocalFTeid: &pfcpType.FTEID{
				V4:          true,
				Teid:        0,
				Ipv4Address: node.UPF.UPIPInfo.Ipv4Address,
			},
			NetworkInstance: []byte(smContext.Dnn),
			UEIPAddress: &pfcpType.UEIPAddress{
				V4:          true,
				Ipv4Address: smContext.PDUAddress.To4(),
			},
		}

		downLink.DownLinkPDR.OuterHeaderRemoval = new(pfcpType.OuterHeaderRemoval)
		downLink.DownLinkPDR.OuterHeaderRemoval.OuterHeaderRemovalDescription = pfcpType.OuterHeaderRemovalGtpUUdpIpv4
		downLink.DownLinkPDR.State = smf_context.RULE_INITIAL

		downLink.DownLinkPDR.FAR.ApplyAction.Forw = true
		downLink.DownLinkPDR.FAR.State = smf_context.RULE_INITIAL
		downLink.DownLinkPDR.FAR.ForwardingParameters = &smf_context.ForwardingParameters{
			DestinationInterface: pfcpType.DestinationInterface{
				InterfaceValue: pfcpType.DestinationInterfaceCore,
			},
			NetworkInstance: []byte(smContext.Dnn),
		}
	}
}

func AllocateDownLinkTEID(node *smf_context.DataPathNode, smContext *smf_context.SMContext, visited map[*smf_context.DataPathNode]bool) {

	if !visited[node] {
		visited[node] = true
	}

	for _, downLink := range node.DataPathToDN {

		child := downLink.To
		allocatedDownLinkTEID := downLink.DownLinkPDR.PDI.LocalFTeid.Teid

		for _, child_downLink := range child.DataPathToDN {

			childDownLinkFAR := child_downLink.DownLinkPDR.FAR
			childDownLinkFAR.ForwardingParameters.OuterHeaderCreation = new(pfcpType.OuterHeaderCreation)
			childDownLinkFAR.ForwardingParameters.OuterHeaderCreation.OuterHeaderCreationDescription = pfcpType.OuterHeaderCreationGtpUUdpIpv4
			childDownLinkFAR.ForwardingParameters.OuterHeaderCreation.Teid = uint32(allocatedDownLinkTEID)
			childDownLinkFAR.ForwardingParameters.OuterHeaderCreation.Ipv4Address = node.UPF.UPIPInfo.Ipv4Address
		}

	}

	for _, upf_link := range node.DataPathToDN {

		child := upf_link.To
		if !visited[child] && child.InUse {
			AllocateDownLinkTEID(child, smContext, visited)
		}

	}
}

func SendUplinkPFCPRule(node *smf_context.DataPathNode, smContext *smf_context.SMContext, visited map[*smf_context.DataPathNode]bool) {

	if !visited[node] {
		visited[node] = true
	}

	addr := net.UDPAddr{
		IP:   node.UPF.NodeID.NodeIdValue,
		Port: pfcpUdp.PFCP_PORT,
	}

	upLink := node.DataPathToAN
	pdrList := []*smf_context.PDR{upLink.UpLinkPDR}
	farList := []*smf_context.FAR{upLink.UpLinkPDR.FAR}
	barList := []*smf_context.BAR{}

	pfcp_message.SendPfcpSessionEstablishmentRequestForULCL(&addr, smContext, pdrList, farList, barList)

	for _, upf_link := range node.DataPathToDN {

		child := upf_link.To
		if !visited[child] && child.InUse {
			SendUplinkPFCPRule(upf_link.To, smContext, visited)
		}
	}

}

func SendDownLinkPFCPRule(node *smf_context.DataPathNode, smContext *smf_context.SMContext, visited map[*smf_context.DataPathNode]bool) {

	if !visited[node] {
		visited[node] = true
	}

	addr := net.UDPAddr{
		IP:   node.UPF.NodeID.NodeIdValue,
		Port: pfcpUdp.PFCP_PORT,
	}

	for _, down_link := range node.DataPathToDN {

		pdrList := []*smf_context.PDR{down_link.DownLinkPDR}
		farList := []*smf_context.FAR{down_link.DownLinkPDR.FAR}
		barList := []*smf_context.BAR{}
		pfcp_message.SendPfcpSessionModificationRequest(&addr, smContext, pdrList, farList, barList)
	}

	if node.IsAnchorUPF() {

		down_link := node.DLDataPathLinkForPSA
		pdrList := []*smf_context.PDR{down_link.DownLinkPDR}
		farList := []*smf_context.FAR{down_link.DownLinkPDR.FAR}
		barList := []*smf_context.BAR{}
		pfcp_message.SendPfcpSessionModificationRequest(&addr, smContext, pdrList, farList, barList)
	}

	for _, upf_link := range node.DataPathToDN {

		child := upf_link.To
		if !visited[child] && child.InUse {
			SendDownLinkPFCPRule(upf_link.To, smContext, visited)
		}
	}

}

func SetUPPSA2Path(smContext *smf_context.SMContext, psa2_path_after_ulcl []*smf_context.UPNode, start_node *smf_context.DataPathNode) {

	lowerBound := 0
	upperBound := len(psa2_path_after_ulcl) - 1
	curDataPathNode := start_node
	var downLink *smf_context.DataPathUpLink

	//Allocate upLink and downLink PDR
	logger.PduSessLog.Traceln("In SetUPPSA2Path")
	for i, node := range psa2_path_after_ulcl {

		logger.PduSessLog.Traceln("Node ", i, ": ", node.UPF.GetUPFIP())
	}
	for idx, _ := range psa2_path_after_ulcl {

		upLink := curDataPathNode.GetUpLink()

		teid, err := curDataPathNode.UPF.GenerateTEID()

		if err != nil {
			logger.PduSessLog.Error(err)
		}

		upLink.UpLinkPDR, err = curDataPathNode.UPF.AddPDR()
		if err != nil {
			logger.PduSessLog.Error(err)
		}

		upLink.UpLinkPDR.Precedence = 32
		upLink.UpLinkPDR.PDI = smf_context.PDI{
			SourceInterface: pfcpType.SourceInterface{
				//Todo:
				//Have to change source interface for different upf
				InterfaceValue: pfcpType.SourceInterfaceAccess,
			},
			LocalFTeid: &pfcpType.FTEID{
				V4:          true,
				Teid:        teid,
				Ipv4Address: curDataPathNode.UPF.UPIPInfo.Ipv4Address,
			},
			NetworkInstance: []byte(smContext.Dnn),
			UEIPAddress: &pfcpType.UEIPAddress{
				V4:          true,
				Ipv4Address: smContext.PDUAddress.To4(),
			},
		}
		upLink.UpLinkPDR.OuterHeaderRemoval = new(pfcpType.OuterHeaderRemoval)
		upLink.UpLinkPDR.OuterHeaderRemoval.OuterHeaderRemovalDescription = pfcpType.OuterHeaderRemovalGtpUUdpIpv4
		upLink.UpLinkPDR.State = smf_context.RULE_INITIAL

		upLink.UpLinkPDR.FAR.ApplyAction.Forw = true
		upLink.UpLinkPDR.FAR.State = smf_context.RULE_INITIAL
		upLink.UpLinkPDR.FAR.ForwardingParameters = &smf_context.ForwardingParameters{
			DestinationInterface: pfcpType.DestinationInterface{
				InterfaceValue: pfcpType.DestinationInterfaceCore,
			},
			NetworkInstance: []byte(smContext.Dnn),
		}

		if curDataPathNode.IsAnchorUPF() {

			downLink = curDataPathNode.DLDataPathLinkForPSA
		} else {
			nextUPFID := psa2_path_after_ulcl[idx+1].UPF.GetUPFID()
			downLink = curDataPathNode.DataPathToDN[nextUPFID]
		}

		downLink.DownLinkPDR, err = curDataPathNode.UPF.AddPDR()

		if err != nil {
			logger.PduSessLog.Error(err)
		}

		teid, err = curDataPathNode.UPF.GenerateTEID()

		if err != nil {
			logger.PduSessLog.Error(err)
		}

		downLink.DownLinkPDR.Precedence = 32
		downLink.DownLinkPDR.PDI = smf_context.PDI{
			SourceInterface: pfcpType.SourceInterface{
				//Todo:
				//Have to change source interface for different upf
				InterfaceValue: pfcpType.SourceInterfaceAccess,
			},
			LocalFTeid: &pfcpType.FTEID{
				V4:          true,
				Teid:        teid,
				Ipv4Address: curDataPathNode.UPF.UPIPInfo.Ipv4Address,
			},
			NetworkInstance: []byte(smContext.Dnn),
			UEIPAddress: &pfcpType.UEIPAddress{
				V4:          true,
				Ipv4Address: smContext.PDUAddress.To4(),
			},
		}

		if !curDataPathNode.IsAnchorUPF() {
			downLink.DownLinkPDR.OuterHeaderRemoval = new(pfcpType.OuterHeaderRemoval)
			downLink.DownLinkPDR.OuterHeaderRemoval.OuterHeaderRemovalDescription = pfcpType.OuterHeaderRemovalGtpUUdpIpv4
			downLink.DownLinkPDR.State = smf_context.RULE_INITIAL
		}

		downLink.DownLinkPDR.FAR.ApplyAction.Forw = true
		downLink.DownLinkPDR.FAR.State = smf_context.RULE_INITIAL
		downLink.DownLinkPDR.FAR.ForwardingParameters = &smf_context.ForwardingParameters{
			DestinationInterface: pfcpType.DestinationInterface{
				InterfaceValue: pfcpType.DestinationInterfaceCore,
			},
			NetworkInstance: []byte(smContext.Dnn),
		}

		if idx != upperBound {
			curDataPathNode = downLink.To
		}

	}

	curDataPathNode = start_node

	//Allocate upLink and downLink TEID
	for idx, _ := range psa2_path_after_ulcl {

		switch idx {
		case lowerBound:

			if !curDataPathNode.IsAnchorUPF() {
				nextUPFID := psa2_path_after_ulcl[idx+1].UPF.GetUPFID()
				downLink = curDataPathNode.DataPathToDN[nextUPFID]
				allocatedDownLinkTEID := downLink.DownLinkPDR.PDI.LocalFTeid.Teid
				child := downLink.To

				var childDownLinkFAR *smf_context.FAR

				if child.IsAnchorUPF() {

					childDownLinkFAR = child.DLDataPathLinkForPSA.DownLinkPDR.FAR
				} else {

					nextNextUPFID := psa2_path_after_ulcl[idx+2].UPF.GetUPFID()
					childDownLinkFAR = child.DataPathToDN[nextNextUPFID].DownLinkPDR.FAR
				}
				childDownLinkFAR.ForwardingParameters.OuterHeaderCreation = new(pfcpType.OuterHeaderCreation)
				childDownLinkFAR.ForwardingParameters.OuterHeaderCreation.OuterHeaderCreationDescription = pfcpType.OuterHeaderCreationGtpUUdpIpv4
				childDownLinkFAR.ForwardingParameters.OuterHeaderCreation.Teid = uint32(allocatedDownLinkTEID)
				childDownLinkFAR.ForwardingParameters.OuterHeaderCreation.Ipv4Address = curDataPathNode.UPF.UPIPInfo.Ipv4Address

			}

		case upperBound:

			parent := curDataPathNode.GetParent()
			if parent != nil {
				allocatedUPLinkTEID := curDataPathNode.DataPathToAN.UpLinkPDR.PDI.LocalFTeid.Teid
				parentUpLinkFAR := parent.GetUpLinkFAR()

				parentUpLinkFAR.ForwardingParameters.OuterHeaderCreation = new(pfcpType.OuterHeaderCreation)
				parentUpLinkFAR.ForwardingParameters.OuterHeaderCreation.OuterHeaderCreationDescription = pfcpType.OuterHeaderCreationGtpUUdpIpv4
				parentUpLinkFAR.ForwardingParameters.OuterHeaderCreation.Teid = uint32(allocatedUPLinkTEID)
				parentUpLinkFAR.ForwardingParameters.OuterHeaderCreation.Ipv4Address = curDataPathNode.UPF.UPIPInfo.Ipv4Address
			}
		default:

			nextUPFID := psa2_path_after_ulcl[idx+1].UPF.GetUPFID()
			downLink = curDataPathNode.DataPathToDN[nextUPFID]
			allocatedDownLinkTEID := downLink.DownLinkPDR.PDI.LocalFTeid.Teid
			child := downLink.To

			var childDownLinkFAR *smf_context.FAR

			if child.IsAnchorUPF() {

				childDownLinkFAR = child.DLDataPathLinkForPSA.DownLinkPDR.FAR
			} else {

				nextNextUPFID := psa2_path_after_ulcl[idx+2].UPF.GetUPFID()
				childDownLinkFAR = child.DataPathToDN[nextNextUPFID].DownLinkPDR.FAR
			}
			childDownLinkFAR.ForwardingParameters.OuterHeaderCreation = new(pfcpType.OuterHeaderCreation)
			childDownLinkFAR.ForwardingParameters.OuterHeaderCreation.OuterHeaderCreationDescription = pfcpType.OuterHeaderCreationGtpUUdpIpv4
			childDownLinkFAR.ForwardingParameters.OuterHeaderCreation.Teid = uint32(allocatedDownLinkTEID)
			childDownLinkFAR.ForwardingParameters.OuterHeaderCreation.Ipv4Address = curDataPathNode.UPF.UPIPInfo.Ipv4Address

			parent := curDataPathNode.GetParent()
			if parent != nil {
				allocatedUPLinkTEID := curDataPathNode.DataPathToAN.UpLinkPDR.PDI.LocalFTeid.Teid
				parentUpLinkFAR := parent.GetUpLinkFAR()

				parentUpLinkFAR.ForwardingParameters.OuterHeaderCreation = new(pfcpType.OuterHeaderCreation)
				parentUpLinkFAR.ForwardingParameters.OuterHeaderCreation.OuterHeaderCreationDescription = pfcpType.OuterHeaderCreationGtpUUdpIpv4
				parentUpLinkFAR.ForwardingParameters.OuterHeaderCreation.Teid = uint32(allocatedUPLinkTEID)
				parentUpLinkFAR.ForwardingParameters.OuterHeaderCreation.Ipv4Address = curDataPathNode.UPF.UPIPInfo.Ipv4Address
			}

		}

		if idx != upperBound {
			curDataPathNode = downLink.To
		}
	}

	curDataPathNode = start_node
	logger.PduSessLog.Traceln("Start Node is PSA: ", curDataPathNode.IsAnchorUPF())
	for idx, _ := range psa2_path_after_ulcl {

		addr := net.UDPAddr{
			IP:   curDataPathNode.UPF.NodeID.NodeIdValue,
			Port: pfcpUdp.PFCP_PORT,
		}

		logger.PduSessLog.Traceln("Send to upf addr: ", addr.String())

		upLink := curDataPathNode.DataPathToAN

		if curDataPathNode.IsAnchorUPF() {

			downLink = curDataPathNode.DLDataPathLinkForPSA
		} else {

			nextUPFID := psa2_path_after_ulcl[idx+1].UPF.GetUPFID()
			downLink = curDataPathNode.DataPathToDN[nextUPFID]
		}

		pdrList := []*smf_context.PDR{upLink.UpLinkPDR, downLink.DownLinkPDR}
		farList := []*smf_context.FAR{upLink.UpLinkPDR.FAR, downLink.DownLinkPDR.FAR}
		barList := []*smf_context.BAR{}

		pfcp_message.SendPfcpSessionEstablishmentRequestForULCL(&addr, smContext, pdrList, farList, barList)
	}

}
