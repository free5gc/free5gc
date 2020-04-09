package smf_producer

import (
	"fmt"
	"free5gc/lib/flowdesc"
	"free5gc/lib/pfcp/pfcpType"
	"free5gc/lib/pfcp/pfcpUdp"
	"free5gc/src/smf/logger"
	"free5gc/src/smf/smf_context"
	"free5gc/src/smf/smf_pfcp/pfcp_message"
	"net"
)

func AddPDUSessionAnchorAndULCL(smContext *smf_context.SMContext) {

	bpManager := smContext.BPManager
	upfRoot := smContext.Tunnel.ULCLRoot
	//select PSA2
	bpManager.SelectPSA2()
	err := upfRoot.EnableUserPlanePath(bpManager.PSA2Path)
	if err != nil {
		logger.PduSessLog.Errorln(err)
		return
	}
	//select an upf as ULCL
	err = bpManager.FindULCL(smContext)
	if err != nil {
		logger.PduSessLog.Errorln(err)
		return
	}

	//Establish PSA2
	EstablishPSA2(smContext)
	//Establish ULCL
	EstablishULCL(smContext)

	//updatePSA1 downlink
	//UpdatePSA1DownLink(smContext)
	//updatePSA2 downlink
	UpdatePSA2DownLink(smContext)
	//update AN for new CN Info

}

func EstablishPSA2(smContext *smf_context.SMContext) {

	//upfRoot := smContext.Tunnel.UpfRoot
	bpMGR := smContext.BPManager
	psa2_path := bpMGR.PSA2Path

	curDataPathNode := bpMGR.ULCLDataPathNode
	upperBound := len(psa2_path) - 1

	if bpMGR.ULCLState == smf_context.IsOnlyULCL {
		for idx := bpMGR.ULCLIdx; idx <= upperBound; idx++ {

			if idx == bpMGR.ULCLIdx {

				nextUPFID := psa2_path[idx+1].UPF.GetUPFID()
				curDataPathNode = curDataPathNode.DataPathToDN[nextUPFID].To
			} else {

				SetUPPSA2Path(smContext, psa2_path[idx:], curDataPathNode)
				break
			}
		}
	}

	logger.PduSessLog.Traceln("End of EstablishPSA2")

	return
}

func EstablishULCL(smContext *smf_context.SMContext) {

	logger.PduSessLog.Traceln("In EstablishULCL")

	bpMGR := smContext.BPManager
	ulcl := bpMGR.ULCLDataPathNode

	if ulcl.IsAnchorUPF() {
		return
	}

	if bpMGR.ULCLState == smf_context.IsOnlyULCL {

		psa1Path := bpMGR.PSA1Path
		psa2Path := bpMGR.PSA2Path
		var psa1NodeAfterUlcl *smf_context.DataPathNode
		var psa2NodeAfterUlcl *smf_context.DataPathNode
		var err error

		ulclIdx := bpMGR.ULCLIdx
		psa1NodeAfterUlcl = ulcl.DataPathToDN[psa1Path[ulclIdx+1].UPF.GetUPFID()].To
		psa2NodeAfterUlcl = ulcl.DataPathToDN[psa2Path[ulclIdx+1].UPF.GetUPFID()].To

		//Get the UPlinkPDR for PSA1
		var UpLinkForPSA1, UpLinkForPSA2 *smf_context.DataPathDownLink
		var DownLinkForPSA2 *smf_context.DataPathUpLink
		//Todo:
		//Put every uplink to BPUplink
		fmt.Println(ulcl.DataPathToAN.UpLinkPDR)
		upLinkIP := ulcl.DataPathToAN.UpLinkPDR.FAR.ForwardingParameters.OuterHeaderCreation.Ipv4Address.String()
		if upLinkIP != psa1NodeAfterUlcl.UPF.UPIPInfo.Ipv4Address.String() {
			UpLinkForPSA1 = ulcl.BPUpLinkPDRs[psa1NodeAfterUlcl.UPF.GetUPFID()]
		} else {
			UpLinkForPSA1 = ulcl.DataPathToAN
			UpLinkForPSA1.DestinationIP = ulcl.DataPathToDN[psa1NodeAfterUlcl.UPF.GetUPFID()].DestinationIP
			UpLinkForPSA1.DestinationPort = ulcl.DataPathToDN[psa1NodeAfterUlcl.UPF.GetUPFID()].DestinationPort
		}

		UpLinkForPSA2 = smf_context.NewDataPathDownLink()
		UpLinkForPSA2.To = UpLinkForPSA1.To
		UpLinkForPSA2.DestinationIP = ulcl.DataPathToDN[psa2NodeAfterUlcl.UPF.GetUPFID()].DestinationIP
		UpLinkForPSA2.DestinationPort = ulcl.DataPathToDN[psa2NodeAfterUlcl.UPF.GetUPFID()].DestinationPort

		UpLinkForPSA2.UpLinkPDR, err = ulcl.UPF.AddPDR()
		if err != nil {
			logger.PduSessLog.Error(err)
		}

		UpLinkForPSA2.UpLinkPDR.Precedence = 32
		UpLinkForPSA2.UpLinkPDR.PDI = smf_context.PDI{
			SourceInterface: pfcpType.SourceInterface{
				//Todo:
				//Have to change source interface for different upf
				InterfaceValue: pfcpType.SourceInterfaceAccess,
			},
			LocalFTeid: &pfcpType.FTEID{
				V4:          true,
				Teid:        UpLinkForPSA1.UpLinkPDR.PDI.LocalFTeid.Teid,
				Ipv4Address: ulcl.UPF.UPIPInfo.Ipv4Address,
			},
			NetworkInstance: []byte(smContext.Dnn),
			UEIPAddress: &pfcpType.UEIPAddress{
				V4:          true,
				Ipv4Address: smContext.PDUAddress.To4(),
			},
		}
		UpLinkForPSA2.UpLinkPDR.OuterHeaderRemoval = new(pfcpType.OuterHeaderRemoval)
		UpLinkForPSA2.UpLinkPDR.OuterHeaderRemoval.OuterHeaderRemovalDescription = pfcpType.OuterHeaderRemovalGtpUUdpIpv4
		UpLinkForPSA2.UpLinkPDR.State = smf_context.RULE_INITIAL

		UpLinkFARForPSA2 := UpLinkForPSA2.UpLinkPDR.FAR
		UpLinkFARForPSA2.ApplyAction.Forw = true
		UpLinkFARForPSA2.State = smf_context.RULE_INITIAL
		UpLinkFARForPSA2.ForwardingParameters = &smf_context.ForwardingParameters{
			DestinationInterface: pfcpType.DestinationInterface{
				InterfaceValue: pfcpType.DestinationInterfaceCore,
			},
			NetworkInstance: []byte(smContext.Dnn),
		}

		UpLinkFARForPSA2.ForwardingParameters.OuterHeaderCreation = new(pfcpType.OuterHeaderCreation)
		UpLinkFARForPSA2.ForwardingParameters.OuterHeaderCreation.OuterHeaderCreationDescription = pfcpType.OuterHeaderCreationGtpUUdpIpv4
		UpLinkFARForPSA2.ForwardingParameters.OuterHeaderCreation.Teid = psa2NodeAfterUlcl.GetUpLinkPDR().PDI.LocalFTeid.Teid
		UpLinkFARForPSA2.ForwardingParameters.OuterHeaderCreation.Ipv4Address = psa2NodeAfterUlcl.UPF.UPIPInfo.Ipv4Address

		UpLinkForPSA1.UpLinkPDR.State = smf_context.RULE_UPDATE
		// UpLinkFARForPSA1 := UpLinkForPSA1.UpLinkPDR.FAR
		// UpLinkFARForPSA1.State = smf_context.RULE_UPDATE
		// UpLinkFARForPSA1.ForwardingParameters.OuterHeaderCreation = new(pfcpType.OuterHeaderCreation)
		// UpLinkFARForPSA1.ForwardingParameters.OuterHeaderCreation.OuterHeaderCreationDescription = pfcpType.OuterHeaderCreationGtpUUdpIpv4
		// UpLinkFARForPSA1.ForwardingParameters.OuterHeaderCreation.Teid = psa1NodeAfterUlcl.GetUpLinkPDR().PDI.LocalFTeid.Teid
		// UpLinkFARForPSA1.ForwardingParameters.OuterHeaderCreation.Ipv4Address = psa1NodeAfterUlcl.UPF.UPIPInfo.Ipv4Address

		ulcl.BPUpLinkPDRs[psa2NodeAfterUlcl.UPF.GetUPFID()] = UpLinkForPSA2
		upLinks := []*smf_context.DataPathDownLink{UpLinkForPSA1, UpLinkForPSA2}

		for _, link := range upLinks {
			FlowDespcription := flowdesc.NewIPFilterRule()
			err = FlowDespcription.SetAction(true) //permit
			if err != nil {
				logger.PduSessLog.Errorf("Error occurs when setting flow despcription: %s\n", err)
			}
			err = FlowDespcription.SetDirection(true) //uplink
			if err != nil {
				logger.PduSessLog.Errorf("Error occurs when setting flow despcription: %s\n", err)
			}
			err = FlowDespcription.SetDestinationIp(link.DestinationIP)
			if err != nil {
				logger.PduSessLog.Errorf("Error occurs when setting flow despcription: %s\n", err)
			}
			err = FlowDespcription.SetDestinationPorts(link.DestinationPort)
			if err != nil {
				logger.PduSessLog.Errorf("Error occurs when setting flow despcription: %s\n", err)
			}
			err = FlowDespcription.SetSourceIp(smContext.PDUAddress.To4().String())
			if err != nil {
				logger.PduSessLog.Errorf("Error occurs when setting flow despcription: %s\n", err)
			}

			FlowDespcriptionStr, err := FlowDespcription.Encode()

			if err != nil {
				logger.PduSessLog.Errorf("Error occurs when encoding flow despcription: %s\n", err)
			}

			link.UpLinkPDR.PDI.SDFFilter = &pfcpType.SDFFilter{
				Bid:                     false,
				Fl:                      false,
				Spi:                     false,
				Ttc:                     false,
				Fd:                      true,
				LengthOfFlowDescription: uint16(len(FlowDespcriptionStr)),
				FlowDescription:         []byte(FlowDespcriptionStr),
			}

			link.UpLinkPDR.Precedence = 30
		}

		//DownLinkForPSA1 = ulcl.DataPathToDN[psa1NodeAfterUlcl.UPF.GetUPFID()]
		DownLinkForPSA2 = ulcl.DataPathToDN[psa2NodeAfterUlcl.UPF.GetUPFID()]

		DownLinkForPSA2.DownLinkPDR, err = ulcl.UPF.AddPDR()
		if err != nil {
			logger.PduSessLog.Error(err)
		}

		teid, err := ulcl.UPF.GenerateTEID()
		DownLinkForPSA2.DownLinkPDR.Precedence = 32
		DownLinkForPSA2.DownLinkPDR.PDI = smf_context.PDI{
			SourceInterface: pfcpType.SourceInterface{
				//Todo:
				//Have to change source interface for different upf
				InterfaceValue: pfcpType.SourceInterfaceAccess,
			},
			LocalFTeid: &pfcpType.FTEID{
				V4:          true,
				Teid:        teid,
				Ipv4Address: ulcl.UPF.UPIPInfo.Ipv4Address,
			},
			NetworkInstance: []byte(smContext.Dnn),
			UEIPAddress: &pfcpType.UEIPAddress{
				V4:          true,
				Ipv4Address: smContext.PDUAddress.To4(),
			},
		}
		DownLinkForPSA2.DownLinkPDR.OuterHeaderRemoval = new(pfcpType.OuterHeaderRemoval)
		DownLinkForPSA2.DownLinkPDR.OuterHeaderRemoval.OuterHeaderRemovalDescription = pfcpType.OuterHeaderRemovalGtpUUdpIpv4
		DownLinkForPSA2.DownLinkPDR.State = smf_context.RULE_INITIAL

		DownLinkFarForPSA2 := DownLinkForPSA2.DownLinkPDR.FAR
		DownLinkFarForPSA2.ApplyAction.Forw = true
		DownLinkFarForPSA2.State = smf_context.RULE_INITIAL
		DownLinkFarForPSA2.ForwardingParameters = &smf_context.ForwardingParameters{
			DestinationInterface: pfcpType.DestinationInterface{
				InterfaceValue: pfcpType.DestinationInterfaceCore,
			},
			NetworkInstance: []byte(smContext.Dnn),
		}

		//Todo:
		//Delete this after finishing new downlinking userplane
		//Todo:
		//Uncommemt after finishing new downlinking userplane
		DownLinkFarForPSA2.ForwardingParameters.OuterHeaderCreation = new(pfcpType.OuterHeaderCreation)
		DownLinkFarForPSA2.ForwardingParameters.OuterHeaderCreation.OuterHeaderCreationDescription = pfcpType.OuterHeaderCreationGtpUUdpIpv4
		fmt.Println("In EstablishULCL")
		fmt.Println(ulcl)
		//TODO: Change this workaround after release 3.0
		workAroundULCL := smContext.Tunnel.UpfRoot
		DownLinkFarForPSA2.ForwardingParameters.OuterHeaderCreation.Teid = workAroundULCL.DownLinkTunnel.MatchedPDR.FAR.ForwardingParameters.OuterHeaderCreation.Teid
		DownLinkFarForPSA2.ForwardingParameters.OuterHeaderCreation.Ipv4Address = workAroundULCL.DownLinkTunnel.MatchedPDR.FAR.ForwardingParameters.OuterHeaderCreation.Ipv4Address //DownLinkForPSA1.DownLinkPDR.FAR.ForwardingParameters.OuterHeaderCreation.Ipv4Address

		// addr := net.UDPAddr{
		// 	IP:   ulcl.Next[psa1NodeAfterUlcl.UPF.GetUPFID()].To.UPF.NodeID.NodeIdValue,
		// 	Port: pfcpUdp.PFCP_PORT,
		// }
		addr := net.UDPAddr{
			IP:   ulcl.UPF.NodeID.NodeIdValue,
			Port: pfcpUdp.PFCP_PORT,
		}
		pdrList := []*smf_context.PDR{UpLinkForPSA1.UpLinkPDR, UpLinkForPSA2.UpLinkPDR, DownLinkForPSA2.DownLinkPDR}
		farList := []*smf_context.FAR{UpLinkForPSA2.UpLinkPDR.FAR, DownLinkForPSA2.DownLinkPDR.FAR}
		barList := []*smf_context.BAR{}

		pfcp_message.SendPfcpSessionModificationRequest(&addr, smContext, pdrList, farList, barList)
		logger.PfcpLog.Info("[SMF] Establish ULCL msg has been send")
	}
}

func UpdatePSA2DownLink(smContext *smf_context.SMContext) {
	logger.PduSessLog.Traceln("In UpdatePSA2DownLink")

	bpMGR := smContext.BPManager
	ulcl := bpMGR.ULCLDataPathNode

	if bpMGR.ULCLState == smf_context.IsOnlyULCL {
		psa2Path := bpMGR.PSA2Path

		var psa2NodeAfterUlcl *smf_context.DataPathNode

		ulclIdx := bpMGR.ULCLIdx
		psa2NodeAfterUlcl = ulcl.DataPathToDN[psa2Path[ulclIdx+1].UPF.GetUPFID()].To
		farList := []*smf_context.FAR{}

		if psa2NodeAfterUlcl.IsAnchorUPF() {

			updateDownLinkFAR := psa2NodeAfterUlcl.DLDataPathLinkForPSA.DownLinkPDR.FAR
			updateDownLinkFAR.State = smf_context.RULE_UPDATE
			updateDownLinkFAR.ForwardingParameters.OuterHeaderCreation = new(pfcpType.OuterHeaderCreation)
			updateDownLinkFAR.ForwardingParameters.OuterHeaderCreation.OuterHeaderCreationDescription = pfcpType.OuterHeaderCreationGtpUUdpIpv4
			updateDownLinkFAR.ForwardingParameters.OuterHeaderCreation.Teid = ulcl.DataPathToDN[psa2NodeAfterUlcl.UPF.GetUPFID()].DownLinkPDR.PDI.LocalFTeid.Teid
			updateDownLinkFAR.ForwardingParameters.OuterHeaderCreation.Ipv4Address = ulcl.UPF.UPIPInfo.Ipv4Address

			farList = append(farList, updateDownLinkFAR)
		} else {

			for _, updateDownLink := range psa2NodeAfterUlcl.DataPathToDN {

				if updateDownLink.DownLinkPDR != nil {
					updateDownLinkFAR := updateDownLink.DownLinkPDR.FAR
					updateDownLinkFAR.State = smf_context.RULE_UPDATE
					updateDownLinkFAR.ForwardingParameters.OuterHeaderCreation = new(pfcpType.OuterHeaderCreation)
					updateDownLinkFAR.ForwardingParameters.OuterHeaderCreation.OuterHeaderCreationDescription = pfcpType.OuterHeaderCreationGtpUUdpIpv4
					updateDownLinkFAR.ForwardingParameters.OuterHeaderCreation.Teid = ulcl.DataPathToDN[psa2NodeAfterUlcl.UPF.GetUPFID()].DownLinkPDR.PDI.LocalFTeid.Teid
					updateDownLinkFAR.ForwardingParameters.OuterHeaderCreation.Ipv4Address = ulcl.UPF.UPIPInfo.Ipv4Address

					farList = append(farList, updateDownLinkFAR)
				}
			}
		}

		addr := net.UDPAddr{
			IP:   psa2NodeAfterUlcl.UPF.NodeID.NodeIdValue,
			Port: pfcpUdp.PFCP_PORT,
		}
		pdrList := []*smf_context.PDR{}
		barList := []*smf_context.BAR{}

		pfcp_message.SendPfcpSessionModificationRequest(&addr, smContext, pdrList, farList, barList)
		logger.PfcpLog.Info("[SMF] Update PSA2 downlink msg has been send")
	}
}
