package smf_producer

import (
	"context"
	"fmt"
	"free5gc/lib/Namf_Communication"
	"free5gc/lib/Nsmf_PDUSession"
	"free5gc/lib/Nudm_SubscriberDataManagement"
	"free5gc/lib/http_wrapper"
	"free5gc/lib/nas"
	"free5gc/lib/nas/nasConvert"
	"free5gc/lib/openapi"
	"free5gc/lib/openapi/common"
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

	"github.com/antihax/optional"
)

func HandlePDUSessionSMContextCreate(rspChan chan smf_message.HandlerResponseMessage, request models.PostSmContextsRequest) {
	var err error
	var response models.PostSmContextsResponse
	response.JsonData = new(models.SmContextCreatedData)

	// Check has PDU Session Establishment Request
	m := nas.NewMessage()
	err = m.GsmMessageDecode(&request.BinaryDataN1SmMessage)
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

	createData := request.JsonData
	smContext := smf_context.NewSMContext(createData.Supi, createData.PduSessionId)
	smContext.SetCreateData(createData)
	smContext.SmStatusNotifyUri = createData.SmContextStatusUri

	// Query UDM
	smf_consumer.SendNFDiscoveryUDM()

	smPlmnID := createData.Guami.PlmnId

	smDataParams := &Nudm_SubscriberDataManagement.GetSmDataParamOpts{
		Dnn:         optional.NewString(createData.Dnn),
		PlmnId:      optional.NewInterface(smPlmnID.Mcc + smPlmnID.Mnc),
		SingleNssai: optional.NewInterface(openapi.MarshToJsonString(smContext.Snssai)),
	}

	SubscriberDataManagementClient := smf_context.SMF_Self().SubscriberDataManagementClient

	sessSubData, _, err := SubscriberDataManagementClient.SessionManagementSubscriptionDataRetrievalApi.GetSmData(context.Background(), smContext.Supi, smDataParams)

	if err != nil {
		logger.PduSessLog.Errorln("Get SessionManagementSubscriptionData error:", err)
	}

	if sessSubData != nil && len(sessSubData) > 0 {
		smContext.DnnConfiguration = sessSubData[0].DnnConfigurations[smContext.Dnn]
	} else {
		logger.PduSessLog.Errorln("SessionManagementSubscriptionData from UDM is nil")
	}

	establishmentRequest := m.PDUSessionEstablishmentRequest
	smContext.HandlePDUSessionEstablishmentRequest(establishmentRequest)

	logger.PduSessLog.Infof("PCF Selection for SMContext SUPI[%s] PDUSessionID[%d]\n", smContext.Supi, smContext.PDUSessionID)
	err = smContext.PCFSelection()

	if err != nil {
		logger.PduSessLog.Errorln("pcf selection error:", err)
	}

	smPolicyData := models.SmPolicyContextData{}

	smPolicyData.Supi = smContext.Supi
	smPolicyData.PduSessionId = smContext.PDUSessionID
	smPolicyData.NotificationUri = fmt.Sprintf("https://%s:%d/", smf_context.SMF_Self().HTTPAddress, smf_context.SMF_Self().HTTPPort)
	smPolicyData.Dnn = smContext.Dnn
	smPolicyData.PduSessionType = nasConvert.PDUSessionTypeToModels(smContext.SelectedPDUSessionType)
	smPolicyData.AccessType = smContext.AnType
	smPolicyData.RatType = smContext.RatType
	smPolicyData.Ipv4Address = smContext.PDUAddress.To4().String()
	smPolicyData.SubsSessAmbr = smContext.DnnConfiguration.SessionAmbr
	smPolicyData.SubsDefQos = smContext.DnnConfiguration.Var5gQosProfile
	smPolicyData.SliceInfo = smContext.Snssai
	smPolicyData.ServingNetwork = &models.NetworkId{
		Mcc: smContext.ServingNetwork.Mcc,
		Mnc: smContext.ServingNetwork.Mnc,
	}

	smPolicyDecision, _, err := smContext.SMPolicyClient.DefaultApi.SmPoliciesPost(context.Background(), smPolicyData)

	if err != nil {
		openapiError := err.(common.GenericOpenAPIError)
		problemDetails := openapiError.Model().(models.ProblemDetails)
		logger.PduSessLog.Errorln("setup sm policy association failed:", err, problemDetails)
	}

	for _, sessRule := range smPolicyDecision.SessRules {
		smContext.SessionRule = sessRule
		break
	}

	smContext.Tunnel = new(smf_context.UPTunnel)

	var dataPathRoot *smf_context.DataPathNode
	defaultUPPath := smf_context.GetUserPlaneInformation().GetDefaultUserPlanePathByDNN(createData.Dnn)
	dataPathRoot = smf_context.GenerateDataPath(defaultUPPath, smContext)
	smContext.Tunnel.UpfRoot = dataPathRoot
	if smf_context.SMF_Self().ULCLSupport && smf_context.CheckUEHasPreConfig(createData.Supi) {
		logger.PduSessLog.Infof("SUPI[%s] has pre-config route", createData.Supi)
		ueRoutingGraph := smf_context.GetUERoutingGraph(createData.Supi)
		dataPathRoot = ueRoutingGraph.GetGraphRoot()
		smContext.Tunnel.ULCLRoot = dataPathRoot
		psaPath := smf_context.GetUserPlaneInformation().GetDefaultUserPlanePathByDNN(createData.Dnn)

		err := dataPathRoot.EnableUserPlanePath(psaPath)
		if err != nil {
			logger.PduSessLog.Error(err)
			return
		}

		smContext.BPManager = smf_context.NewBPManager(createData.Supi)
		smContext.BPManager.SetPSAStatus(psaPath)
		smContext.BPManager.PSA1Path = psaPath
		SetUpUplinkUserPlane(smContext.Tunnel.ULCLRoot, smContext)
	}

	if dataPathRoot == nil {
		logger.PduSessLog.Errorf("Path for serve DNN[%s] not found\n", createData.Dnn)
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

	response.JsonData = smContext.BuildCreatedData()
	rspChan <- smf_message.HandlerResponseMessage{HTTPResponse: &http_wrapper.Response{
		Header: http.Header{
			"Location": {smContext.Ref},
		},
		Status: http.StatusCreated,
		Body:   response,
	}}

	// TODO: UECM registration

	SendPFCPRule(smContext, smContext.Tunnel.UpfRoot)

	smf_consumer.SendNFDiscoveryServingAMF(smContext)

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

			curDataPathNode := smContext.Tunnel.UpfRoot

			for curDataPathNode != nil {
				seqNum = pfcp_message.SendPfcpSessionDeletionRequest(curDataPathNode.UPF.PFCPAddr(), smContext)
				curDataPathNode = curDataPathNode.DownLinkTunnel.SrcEndPoint
			}

			return seqNum, response
		case nas.MsgTypePDUSessionReleaseComplete:
			// Send Release Notify to AMF
			smf_consumer.SendSMContextStatusNotification(smContext.SmStatusNotifyUri)
			smf_context.RemoveSMContext(smContext.Ref)
			response.JsonData.UpCnxState = models.UpCnxState_DEACTIVATED

			return seqNum, response
		}

	}

	tunnel := smContext.Tunnel
	pdrList := []*smf_context.PDR{}
	farList := []*smf_context.FAR{}
	barList := []*smf_context.BAR{}

	switch smContextUpdateData.UpCnxState {
	case models.UpCnxState_ACTIVATING:
		response.JsonData.N2SmInfo = &models.RefToBinaryData{ContentId: "PDUSessionResourceSetupRequestTransfer"}
		response.JsonData.UpCnxState = models.UpCnxState_ACTIVATING
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
		DLPDR := tunnel.UpfRoot.DownLinkTunnel.MatchedPDR
		if DLPDR == nil {
			logger.PduSessLog.Errorf("Release Error")
		} else {
			DLPDR.FAR.State = smf_context.RULE_UPDATE
			DLPDR.FAR.ApplyAction.Forw = false
			DLPDR.FAR.ApplyAction.Buff = true
			DLPDR.FAR.ApplyAction.Nocp = true
		}

		farList = []*smf_context.FAR{DLPDR.FAR}
	}

	var err error

	switch smContextUpdateData.N2SmInfoType {
	case models.N2SmInfoType_PDU_RES_SETUP_RSP:
		DLPDR := tunnel.UpfRoot.DownLinkTunnel.MatchedPDR

		// TODO: Setup Uplink Routing

		DLPDR.Precedence = 32
		DLPDR.PDI = smf_context.PDI{
			SourceInterface: pfcpType.SourceInterface{
				InterfaceValue: pfcpType.SourceInterfaceSgiLanN6Lan,
			},
			NetworkInstance: []byte(smContext.Dnn),
			// TODO: Should Uncomment this after FR5GC-1029 is solved
			// UEIPAddress: &pfcpType.UEIPAddress{
			// 	V4:          true,
			// 	Ipv4Address: smContext.PDUAddress.To4(),
			// },
		}

		DLPDR.FAR.ApplyAction = pfcpType.ApplyAction{Buff: false, Drop: false, Dupl: false, Forw: true, Nocp: false}
		DLPDR.FAR.ForwardingParameters = &smf_context.ForwardingParameters{
			DestinationInterface: pfcpType.DestinationInterface{
				InterfaceValue: pfcpType.DestinationInterfaceAccess,
			},
			NetworkInstance: []byte(smContext.Dnn),
		}
		err = smf_context.HandlePDUSessionResourceSetupResponseTransfer(body.BinaryDataN2SmInformation, smContext)

		DLPDR.State = smf_context.RULE_UPDATE
		DLPDR.FAR.State = smf_context.RULE_UPDATE
		pdrList = []*smf_context.PDR{DLPDR}
		farList = []*smf_context.FAR{DLPDR.FAR}

	case models.N2SmInfoType_PATH_SWITCH_REQ:
		DLPDR := tunnel.UpfRoot.DownLinkTunnel.MatchedPDR
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

		pdrList = []*smf_context.PDR{DLPDR}
		farList = []*smf_context.FAR{DLPDR.FAR}

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
		IP:   smContext.Tunnel.UpfRoot.UPF.NodeID.NodeIdValue,
		Port: pfcpUdp.PFCP_PORT,
	}

	seqNum = pfcp_message.SendPfcpSessionModificationRequest(&addr, smContext, pdrList, farList, barList)

	return seqNum, response
}

func HandlePDUSessionSMContextRelease(rspChan chan smf_message.HandlerResponseMessage, smContextRef string, body models.ReleaseSmContextRequest) (seqNum uint32) {
	smContext := smf_context.GetSMContext(smContextRef)

	// smf_context.RemoveSMContext(smContext.Ref)

	addr := net.UDPAddr{
		IP:   smContext.Tunnel.Node.NodeID.NodeIdValue,
		Port: pfcpUdp.PFCP_PORT,
	}

	seqNum = pfcp_message.SendPfcpSessionDeletionRequest(&addr, smContext)
	return seqNum

	// rspChan <- smf_message.HandlerResponseMessage{HTTPResponse: &http_wrapper.Response{
	// 	Header: nil,
	// 	Status: http.StatusNoContent,
	// 	Body:   nil,
	// }}
}

func SendPFCPRule(smContext *smf_context.SMContext, root *smf_context.DataPathNode) {
	curDataPathNode := root
	for curDataPathNode != nil {
		pdrList := make([]*smf_context.PDR, 0, 2)
		farList := make([]*smf_context.FAR, 0, 2)
		if !curDataPathNode.HaveSession {
			if curDataPathNode.UpLinkTunnel != nil && curDataPathNode.UpLinkTunnel.MatchedPDR != nil {
				pdrList = append(pdrList, curDataPathNode.UpLinkTunnel.MatchedPDR)
				farList = append(farList, curDataPathNode.UpLinkTunnel.MatchedPDR.FAR)
			}
			if curDataPathNode.DownLinkTunnel != nil && curDataPathNode.DownLinkTunnel.MatchedPDR != nil {
				pdrList = append(pdrList, curDataPathNode.DownLinkTunnel.MatchedPDR)
				farList = append(farList, curDataPathNode.DownLinkTunnel.MatchedPDR.FAR)
			}

			pfcp_message.SendPfcpSessionEstablishmentRequestForULCL(curDataPathNode.UPF.PFCPAddr(), smContext, pdrList, farList, nil)
			curDataPathNode.HaveSession = true
		} else {
			if curDataPathNode.UpLinkTunnel != nil && curDataPathNode.UpLinkTunnel.MatchedPDR != nil {
				pdrList = append(pdrList, curDataPathNode.UpLinkTunnel.MatchedPDR)
				farList = append(farList, curDataPathNode.UpLinkTunnel.MatchedPDR.FAR)
			}
			if curDataPathNode.DownLinkTunnel != nil && curDataPathNode.DownLinkTunnel.MatchedPDR != nil {
				pdrList = append(pdrList, curDataPathNode.DownLinkTunnel.MatchedPDR)
				farList = append(farList, curDataPathNode.DownLinkTunnel.MatchedPDR.FAR)
			}

			pfcp_message.SendPfcpSessionModificationRequest(curDataPathNode.UPF.PFCPAddr(), smContext, pdrList, farList, nil)
		}
		if curDataPathNode.DownLinkTunnel == nil {
			break
		}
		curDataPathNode = curDataPathNode.DownLinkTunnel.SrcEndPoint
	}
}
