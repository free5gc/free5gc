package smf_producer

import (
	"free5gc/lib/CommonConsumerTestData/SMF/TestPDUSession"
	"free5gc/lib/Namf_Communication"
	"free5gc/lib/Nsmf_PDUSession"
	"free5gc/lib/http_wrapper"
	"free5gc/lib/nas"
	"free5gc/lib/openapi/models"
	"free5gc/lib/pfcp/pfcpType"
	"free5gc/lib/pfcp/pfcpUdp"
	"free5gc/src/smf/logger"
	"free5gc/src/smf/smf_consumer"
	"free5gc/src/smf/smf_context"
	"free5gc/src/smf/smf_handler/smf_message"
	"free5gc/src/smf/smf_pfcp/pfcp_message"
	"net"
	"net/http"
)

func HandlePDUSessionSMContextCreate(rspChan chan smf_message.HandlerResponseMessage, request models.PostSmContextsRequest) {
	var response models.PostSmContextsResponse
	response.JsonData = new(models.SmContextCreatedData)

	createData := request.JsonData
	smContext := smf_context.NewSMContext(createData.Supi, createData.PduSessionId)

	smContext.PDUAddress = smf_context.AllocUEIP()

	if request.BinaryDataN1SmMessage == nil {
		rspChan <- smf_message.HandlerResponseMessage{
			HTTPResponse: &http_wrapper.Response{
				Header: nil,
				Status: http.StatusForbidden,
				Body: models.PostSmContextsErrorResponse{
					JsonData: &models.SmContextCreateError{
						Error: &Nsmf_PDUSession.N1SmError,
					},
				},
			},
		}
		return
	}

	selectedUPF := smf_context.SelectUPFByDnn(createData.Dnn)
	if selectedUPF == nil {
		logger.PduSessLog.Errorf("UPF for serve DNN[%s] not found\n", createData.Dnn)
		rspChan <- smf_message.HandlerResponseMessage{
			HTTPResponse: &http_wrapper.Response{
				Header: nil,
				Status: http.StatusForbidden,
				Body: models.PostSmContextsErrorResponse{
					JsonData: &models.SmContextCreateError{
						Error:   &Nsmf_PDUSession.DnnNotSupported,
						N1SmMsg: &models.RefToBinaryData{ContentId: "N1Msg"},
					},
				},
			},
		}
	}

	m := nas.NewMessage()
	err := m.GsmMessageDecode(&request.BinaryDataN1SmMessage)
	if err != nil || m.GsmHeader.GetMessageType() != nas.MsgTypePDUSessionEstablishmentRequest {
		rspChan <- smf_message.HandlerResponseMessage{
			HTTPResponse: &http_wrapper.Response{
				Header: nil,
				Status: http.StatusForbidden,
				Body: models.PostSmContextsErrorResponse{
					JsonData: &models.SmContextCreateError{
						Error: &Nsmf_PDUSession.N1SmError,
					},
				},
			},
		}
		return
	}

	establishmentRequest := m.PDUSessionEstablishmentRequest

	smContext.PDUSessionID = int32(establishmentRequest.PDUSessionID.Octet)
	smContext.SetCreateData(createData)
	response.JsonData = smContext.BuildCreatedData()
	rspChan <- smf_message.HandlerResponseMessage{HTTPResponse: &http_wrapper.Response{
		Header: http.Header{
			"Location": {smContext.Ref},
		},
		Status: http.StatusCreated,
		Body:   response,
	}}

	// TODO: UECM registration

	smContext.Tunnel = new(smf_context.UPTunnel)
	smContext.Tunnel.Node = selectedUPF
	tunnel := smContext.Tunnel
	// Establish UP

	tunnel.ULTEID = tunnel.Node.GenerateTEID()

	tunnel.ULPDR = smContext.Tunnel.Node.AddPDR()
	tunnel.ULPDR.Precedence = 32
	tunnel.ULPDR.PDI = smf_context.PDI{
		SourceInterface: pfcpType.SourceInterface{
			InterfaceValue: pfcpType.SourceInterfaceAccess,
		},
		LocalFTeid: pfcpType.FTEID{
			V4:          true,
			Teid:        tunnel.ULTEID,
			Ipv4Address: tunnel.Node.UPIPInfo.Ipv4Address,
		},
		NetworkInstance: []byte(smContext.Dnn),
		UEIPAddress: &pfcpType.UEIPAddress{
			V4:          true,
			Ipv4Address: smContext.PDUAddress.To4(),
		},
	}
	tunnel.ULPDR.OuterHeaderRemoval = new(pfcpType.OuterHeaderRemoval)
	tunnel.ULPDR.OuterHeaderRemoval.OuterHeaderRemovalDescription = pfcpType.OuterHeaderRemovalGtpUUdpIpv4

	tunnel.ULPDR.FAR.ApplyAction.Forw = true
	tunnel.ULPDR.FAR.ForwardingParameters = &smf_context.ForwardingParameters{
		DestinationInterface: pfcpType.DestinationInterface{
			InterfaceValue: pfcpType.DestinationInterfaceCore,
		},
		NetworkInstance: []byte(smContext.Dnn),
	}

	// TODO: PCF Selection

	addr := net.UDPAddr{
		IP:   smContext.Tunnel.Node.NodeID.NodeIdValue,
		Port: pfcpUdp.PFCP_PORT,
	}
	pfcp_message.SendPfcpSessionEstablishmentRequest(&addr, smContext)

	smf_consumer.SendNFDiscoveryServingAMF(smContext)

	// Workaround AMF Profile
	// smContext.AMFProfile = models.NfProfile{
	// 	NfServices: &[]models.NfService{
	// 		{
	// 			ServiceName: models.ServiceName_NAMF_COMM,
	// 			ApiPrefix:   "https://127.0.0.1:29518",
	// 		},
	// 	},
	// }

	for _, service := range *smContext.AMFProfile.NfServices {
		if service.ServiceName == models.ServiceName_NAMF_COMM {
			communicationConf := Namf_Communication.NewConfiguration()
			communicationConf.SetBasePath(service.ApiPrefix)
			smContext.CommunicationClient = Namf_Communication.NewAPIClient(communicationConf)
		}
	}
}

func HandlePDUSessionSMContextUpdate(rspChan chan smf_message.HandlerResponseMessage, smContextRef string, body models.UpdateSmContextRequest) (seqNum uint32, resBody models.UpdateSmContextResponse) {
	smContext := smf_context.GetSMContext(smContextRef)

	if smContext == nil {
		rspChan <- smf_message.HandlerResponseMessage{HTTPResponse: &http_wrapper.Response{
			Header: nil,
			Status: http.StatusNotFound,
			Body: models.UpdateSmContextErrorResponse{
				JsonData: &models.SmContextUpdateError{
					UpCnxState: models.UpCnxState_DEACTIVATED,
					Error: &models.ProblemDetails{
						Type:   "Resource Not Found",
						Title:  "SMContext Ref is not found",
						Status: http.StatusNotFound,
					},
				},
			},
		}}
		return
	}

	var response models.UpdateSmContextResponse
	response.JsonData = new(models.SmContextUpdatedData)

	smContextUpdateData := body.JsonData

	if body.BinaryDataN1SmMessage != nil {
		m := nas.NewMessage()
		err := m.GsmMessageDecode(&body.BinaryDataN1SmMessage)
		if err != nil {
			logger.PduSessLog.Error(err)
			return
		}
		switch m.GsmHeader.GetMessageType() {
		case nas.MsgTypePDUSessionReleaseRequest:
			smContext.HandlePDUSessionReleaseRequest(m.PDUSessionReleaseRequest)
			buf, _ := smf_context.BuildGSMPDUSessionReleaseCommand(smContext)
			response.BinaryDataN1SmMessage = buf
			response.JsonData.N1SmMsg = &models.RefToBinaryData{ContentId: "PDUSessionReleaseCommand"}

			response.JsonData.N2SmInfo = &models.RefToBinaryData{ContentId: "PDUResourceReleaseCommand"}
			response.JsonData.N2SmInfoType = models.N2SmInfoType_PDU_RES_REL_CMD

			buf, err := smf_context.BuildPDUSessionResourceReleaseCommandTransfer(smContext)
			response.BinaryDataN2SmInformation = buf
			if err != nil {
				logger.PduSessLog.Error(err)
			}
		}

	}

	tunnel := smContext.Tunnel
	pdr_list := []*smf_context.PDR{}
	far_list := []*smf_context.FAR{}
	bar_list := []*smf_context.BAR{}

	switch smContextUpdateData.UpCnxState {
	case models.UpCnxState_ACTIVATING:
		response.JsonData.N2SmInfo = &models.RefToBinaryData{ContentId: "PDUSessionResourceSetupRequestTransfer"}
		response.JsonData.UpCnxState = TestPDUSession.ACTIVATING
		response.JsonData.N2SmInfoType = models.N2SmInfoType_PDU_RES_SETUP_REQ

		n2Buf, err := smf_context.BuildPDUSessionResourceSetupRequestTransfer(smContext)
		if err != nil {
			logger.PduSessLog.Errorf("Build PDUSession Resource Setup Request Transfer Error(%s)", err.Error())
		}
		response.BinaryDataN2SmInformation = n2Buf
		response.JsonData.N2SmInfoType = models.N2SmInfoType_PDU_RES_SETUP_REQ
	case models.UpCnxState_DEACTIVATED:
		response.JsonData.UpCnxState = models.UpCnxState_DEACTIVATED
		smContext.UpCnxState = body.JsonData.UpCnxState
		smContext.UeLocation = body.JsonData.UeLocation
		// TODO: Deactivate N2 downlink tunnel
		//Set FAR and An, N3 Release Info
		if tunnel.DLPDR == nil {
			logger.PduSessLog.Errorf("Release Error")
		} else {
			tunnel.DLPDR.FAR.State = smf_context.RULE_UPDATE
			tunnel.DLPDR.FAR.ApplyAction.Forw = false
			tunnel.DLPDR.FAR.ApplyAction.Buff = true
			tunnel.DLPDR.FAR.ApplyAction.Nocp = true

			if tunnel.DLPDR.FAR.BAR == nil {
				tunnel.DLPDR.FAR.BAR = smContext.Tunnel.Node.AddBAR()
				bar_list = []*smf_context.BAR{tunnel.DLPDR.FAR.BAR}
			}

		}

		far_list = []*smf_context.FAR{tunnel.DLPDR.FAR}
	}

	var err error

	switch smContextUpdateData.N2SmInfoType {
	case models.N2SmInfoType_PDU_RES_SETUP_RSP:
		if tunnel.DLPDR == nil {
			tunnel.DLPDR = smContext.Tunnel.Node.AddPDR()
		} else {
			tunnel.DLPDR.State = smf_context.RULE_UPDATE
		}
		tunnel.DLPDR.Precedence = 32
		tunnel.DLPDR.PDI = smf_context.PDI{
			SourceInterface: pfcpType.SourceInterface{
				InterfaceValue: pfcpType.SourceInterfaceSgiLanN6Lan,
			},
			LocalFTeid: pfcpType.FTEID{
				V4:          true,
				Teid:        tunnel.ULTEID,
				Ipv4Address: tunnel.Node.UPIPInfo.Ipv4Address,
			},
			NetworkInstance: []byte(smContext.Dnn),
			UEIPAddress: &pfcpType.UEIPAddress{
				V4:          true,
				Ipv4Address: smContext.PDUAddress.To4(),
			},
		}

		tunnel.DLPDR.FAR.ApplyAction.Forw = true
		tunnel.DLPDR.FAR.ApplyAction.Nocp = false
		tunnel.DLPDR.FAR.ApplyAction.Drop = false
		tunnel.DLPDR.FAR.ApplyAction.Buff = false
		tunnel.DLPDR.FAR.ApplyAction.Dupl = false
		tunnel.DLPDR.FAR.ForwardingParameters = &smf_context.ForwardingParameters{
			DestinationInterface: pfcpType.DestinationInterface{
				InterfaceValue: pfcpType.DestinationInterfaceAccess,
			},
			NetworkInstance: []byte(smContext.Dnn),
		}
		err = smf_context.HandlePDUSessionResourceSetupResponseTransfer(body.BinaryDataN2SmInformation, smContext)

		pdr_list = []*smf_context.PDR{tunnel.DLPDR}
		far_list = []*smf_context.FAR{tunnel.DLPDR.FAR}

	case models.N2SmInfoType_PATH_SWITCH_REQ:
		err = smf_context.HandlePathSwitchRequestTransfer(body.BinaryDataN2SmInformation, smContext)
		n2Buf, err := smf_context.BuildPathSwitchRequestAcknowledgeTransfer(smContext)
		if err != nil {
			logger.PduSessLog.Errorf("Build Path Switch Transfer Error(%s)", err.Error())
		}

		response.BinaryDataN2SmInformation = n2Buf
		response.JsonData.N2SmInfoType = models.N2SmInfoType_PATH_SWITCH_REQ_ACK
		response.JsonData.N2SmInfo = &models.RefToBinaryData{
			ContentId: "PATH_SWITCH_REQ_ACK",
		}

		pdr_list = []*smf_context.PDR{tunnel.DLPDR}
		far_list = []*smf_context.FAR{tunnel.DLPDR.FAR}

	case models.N2SmInfoType_PATH_SWITCH_SETUP_FAIL:
		err = smf_context.HandlePathSwitchRequestSetupFailedTransfer(body.BinaryDataN2SmInformation, smContext)
	case models.N2SmInfoType_HANDOVER_REQUIRED:
		response.JsonData.N2SmInfo = &models.RefToBinaryData{ContentId: "Handover"}
	}

	switch smContextUpdateData.HoState {
	case models.HoState_PREPARING:
		smContext.HoState = models.HoState_PREPARING
		err = smf_context.HandleHandoverRequiredTransfer(body.BinaryDataN2SmInformation, smContext)
		response.JsonData.N2SmInfoType = models.N2SmInfoType_PDU_RES_SETUP_REQ

		n2Buf, err := smf_context.BuildPDUSessionResourceSetupRequestTransfer(smContext)
		if err != nil {
			logger.PduSessLog.Errorf("Build PDUSession Resource Setup Request Transfer Error(%s)", err.Error())
		}
		response.BinaryDataN2SmInformation = n2Buf
		response.JsonData.N2SmInfoType = models.N2SmInfoType_PDU_RES_SETUP_REQ
		response.JsonData.N2SmInfo = &models.RefToBinaryData{
			ContentId: "PDU_RES_SETUP_REQ",
		}
		response.JsonData.HoState = models.HoState_PREPARING
	case models.HoState_PREPARED:
		smContext.HoState = models.HoState_PREPARED
		response.JsonData.HoState = models.HoState_PREPARED
		err = smf_context.HandleHandoverRequestAcknowledgeTransfer(body.BinaryDataN2SmInformation, smContext)
		n2Buf, err := smf_context.BuildHandoverCommandTransfer(smContext)
		if err != nil {
			logger.PduSessLog.Errorf("Build PDUSession Resource Setup Request Transfer Error(%s)", err.Error())
		}
		response.BinaryDataN2SmInformation = n2Buf
		response.JsonData.N2SmInfoType = models.N2SmInfoType_HANDOVER_CMD
		response.JsonData.N2SmInfo = &models.RefToBinaryData{
			ContentId: "HANDOVER_CMD",
		}
		response.JsonData.HoState = models.HoState_PREPARING
	case models.HoState_COMPLETED:
		smContext.HoState = models.HoState_COMPLETED
		response.JsonData.HoState = models.HoState_COMPLETED
	}

	if err != nil {
		logger.PduSessLog.Error(err)
	}

	addr := net.UDPAddr{
		IP:   smContext.Tunnel.Node.NodeID.NodeIdValue,
		Port: pfcpUdp.PFCP_PORT,
	}

	seqNum = pfcp_message.SendPfcpSessionModificationRequest(&addr, smContext, pdr_list, far_list, bar_list)

	return seqNum, response
}

func HandlePDUSessionSMContextRelease(rspChan chan smf_message.HandlerResponseMessage, smContextRef string, body models.ReleaseSmContextRequest) {
	smContext := smf_context.GetSMContext(smContextRef)

	smf_context.RemoveSMContext(smContext.Ref)

	addr := net.UDPAddr{
		IP:   smContext.Tunnel.Node.NodeID.NodeIdValue,
		Port: pfcpUdp.PFCP_PORT,
	}

	pfcp_message.SendPfcpSessionDeletionRequest(&addr, smContext)

	rspChan <- smf_message.HandlerResponseMessage{HTTPResponse: &http_wrapper.Response{
		Header: nil,
		Status: http.StatusNoContent,
		Body:   nil,
	}}
}
