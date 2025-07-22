package consumer

import (
	"fmt"
	"net"
	"regexp"
	"strconv"
	"strings"
	"sync"

	"github.com/pkg/errors"

	"github.com/free5gc/nas/nasConvert"
	"github.com/free5gc/nas/nasType"
	"github.com/free5gc/openapi/models"
	"github.com/free5gc/openapi/pcf/SMPolicyControl"
	smf_context "github.com/free5gc/smf/internal/context"
	"github.com/free5gc/util/flowdesc"
)

type npcfService struct {
	consumer *Consumer

	SMPolicyControlMu sync.RWMutex

	SMPolicyControlClients map[string]*SMPolicyControl.APIClient
}

func (s *npcfService) getSMPolicyControlClient(uri string) *SMPolicyControl.APIClient {
	if uri == "" {
		return nil
	}
	s.SMPolicyControlMu.RLock()
	client, ok := s.SMPolicyControlClients[uri]
	if ok {
		s.SMPolicyControlMu.RUnlock()
		return client
	}

	configuration := SMPolicyControl.NewConfiguration()
	configuration.SetBasePath(uri)
	client = SMPolicyControl.NewAPIClient(configuration)

	s.SMPolicyControlMu.RUnlock()
	s.SMPolicyControlMu.Lock()
	defer s.SMPolicyControlMu.Unlock()
	s.SMPolicyControlClients[uri] = client
	return client
}

// SendSMPolicyAssociationCreate create the session management association to the PCF
func (s *npcfService) SendSMPolicyAssociationCreate(smContext *smf_context.SMContext) (
	string, *models.SmPolicyDecision, error,
) {
	var client *SMPolicyControl.APIClient

	// Create SMPolicyControl Client for this SM Context
	for _, service := range smContext.SelectedPCFProfile.NfServices {
		if service.ServiceName == models.ServiceName_NPCF_SMPOLICYCONTROL {
			client = s.getSMPolicyControlClient(service.ApiPrefix)
		}
	}

	if client == nil {
		return "", nil, errors.Errorf("smContext not selected PCF")
	}

	smPolicyData := models.SmPolicyContextData{}

	smPolicyData.Supi = smContext.Supi
	smPolicyData.PduSessionId = smContext.PDUSessionID
	smPolicyData.NotificationUri = fmt.Sprintf("%s://%s:%d/nsmf-callback/sm-policies/%s",
		smf_context.GetSelf().URIScheme,
		smf_context.GetSelf().RegisterIPv4,
		smf_context.GetSelf().SBIPort,
		smContext.Ref,
	)
	smPolicyData.Dnn = smContext.Dnn
	smPolicyData.PduSessionType = nasConvert.PDUSessionTypeToModels(smContext.SelectedPDUSessionType)
	smPolicyData.AccessType = smContext.AnType
	smPolicyData.RatType = smContext.RatType
	smPolicyData.Ipv4Address = smContext.PDUAddress.To4().String()
	smPolicyData.SubsSessAmbr = smContext.DnnConfiguration.SessionAmbr
	smPolicyData.SubsDefQos = smContext.DnnConfiguration.Var5gQosProfile
	smPolicyData.SliceInfo = smContext.SNssai
	smPolicyData.ServingNetwork = &models.PlmnIdNid{
		Mcc: smContext.ServingNetwork.Mcc,
		Mnc: smContext.ServingNetwork.Mnc,
	}
	smPolicyData.SuppFeat = "F"

	ctx, _, err := smf_context.GetSelf().
		GetTokenCtx(models.ServiceName_NPCF_SMPOLICYCONTROL, models.NrfNfManagementNfType_PCF)
	if err != nil {
		return "", nil, err
	}

	var smPolicyID string
	var smPolicyDecision *models.SmPolicyDecision
	request := &SMPolicyControl.CreateSMPolicyRequest{
		SmPolicyContextData: &smPolicyData,
	}

	smPolicyDecisionFromPCF, err := client.SMPoliciesCollectionApi.CreateSMPolicy(ctx, request)
	if err != nil || smPolicyDecisionFromPCF == nil {
		return "", nil, err
	}

	smPolicyDecision = &smPolicyDecisionFromPCF.SmPolicyDecision
	loc := smPolicyDecisionFromPCF.Location
	if smPolicyID = s.extractSMPolicyIDFromLocation(loc); len(smPolicyID) == 0 {
		return "", nil, fmt.Errorf("SMPolicy ID parse failed")
	}
	return smPolicyID, smPolicyDecision, nil
}

var smPolicyRegexp = regexp.MustCompile(`http[s]?\://.*/npcf-smpolicycontrol/v\d+/sm-policies/(.*)`)

func (s *npcfService) extractSMPolicyIDFromLocation(location string) string {
	match := smPolicyRegexp.FindStringSubmatch(location)
	if len(match) > 1 {
		return match[1]
	}
	// not match submatch
	return ""
}

func (s *npcfService) SendSMPolicyAssociationUpdateByUERequestModification(
	smContext *smf_context.SMContext,
	qosRules nasType.QoSRules, qosFlowDescs nasType.QoSFlowDescs,
) (*models.SmPolicyDecision, error) {
	updateSMPolicy := models.SmPolicyUpdateContextData{}

	updateSMPolicy.RepPolicyCtrlReqTriggers = []models.PolicyControlRequestTrigger{
		models.PolicyControlRequestTrigger_RES_MO_RE,
	}

	// UE SHOULD only create ONE QoS Flow in a request (TS 24.501 6.4.2.2)
	if len(qosRules) == 0 {
		return nil, errors.New("QoS Rule not found")
	}
	if len(qosFlowDescs) == 0 {
		return nil, errors.New("QoS Flow Description not found")
	}

	rule := qosRules[0]
	flowDesc := qosFlowDescs[0]

	var ruleOp models.RuleOperation
	switch rule.Operation {
	case nasType.OperationCodeCreateNewQoSRule:
		ruleOp = models.RuleOperation_CREATE_PCC_RULE
	case nasType.OperationCodeDeleteExistingQoSRule:
		ruleOp = models.RuleOperation_DELETE_PCC_RULE
	case nasType.OperationCodeModifyExistingQoSRuleAndAddPacketFilters:
		ruleOp = models.RuleOperation_MODIFY_PCC_RULE_AND_ADD_PACKET_FILTERS
	case nasType.OperationCodeModifyExistingQoSRuleAndDeletePacketFilters:
		ruleOp = models.RuleOperation_MODIFY_PCC_RULE_AND_DELETE_PACKET_FILTERS
	case nasType.OperationCodeModifyExistingQoSRuleAndReplaceAllPacketFilters:
		ruleOp = models.RuleOperation_MODIFY_PCC_RULE_AND_REPLACE_PACKET_FILTERS
	case nasType.OperationCodeModifyExistingQoSRuleWithoutModifyingPacketFilters:
		ruleOp = models.RuleOperation_MODIFY_PCC_RULE_WITHOUT_MODIFY_PACKET_FILTERS
	default:
		return nil, errors.New("QoS Rule Operation Unknown")
	}

	ueInitResReq := &models.UeInitiatedResourceRequest{}
	ueInitResReq.RuleOp = ruleOp
	ueInitResReq.Precedence = int32(rule.Precedence)
	ueInitResReq.ReqQos = new(models.RequestedQos)

	for _, parameter := range flowDesc.Parameters {
		switch parameter.Identifier() {
		case nasType.ParameterIdentifier5QI:
			para5Qi := parameter.(*nasType.QoSFlow5QI)
			ueInitResReq.ReqQos.Var5qi = int32(para5Qi.FiveQI)
		case nasType.ParameterIdentifierGFBRUplink:
			paraGFBRUplink := parameter.(*nasType.QoSFlowGFBRUplink)
			ueInitResReq.ReqQos.GbrUl = s.nasBitRateToString(paraGFBRUplink.Value, paraGFBRUplink.Unit)
		case nasType.ParameterIdentifierGFBRDownlink:
			paraGFBRDownlink := parameter.(*nasType.QoSFlowGFBRDownlink)
			ueInitResReq.ReqQos.GbrDl = s.nasBitRateToString(paraGFBRDownlink.Value, paraGFBRDownlink.Unit)
		}
	}

	updateSMPolicy.UeInitResReq = ueInitResReq

	for _, pf := range rule.PacketFilterList {
		if PackFiltInfo, err := s.buildPktFilterInfo(pf); err != nil {
			smContext.Log.Warning("Build PackFiltInfo failed", err)
			continue
		} else {
			updateSMPolicy.UeInitResReq.PackFiltInfo = append(updateSMPolicy.UeInitResReq.PackFiltInfo, *PackFiltInfo)
		}
	}

	ctx, _, err := smf_context.GetSelf().
		GetTokenCtx(models.ServiceName_NPCF_SMPOLICYCONTROL, models.NrfNfManagementNfType_PCF)
	if err != nil {
		return nil, err
	}

	var client *SMPolicyControl.APIClient

	// Create SMPolicyControl Client for this SM Context
	for _, service := range smContext.SelectedPCFProfile.NfServices {
		if service.ServiceName == models.ServiceName_NPCF_SMPOLICYCONTROL {
			client = s.getSMPolicyControlClient(service.ApiPrefix)
		}
	}

	var smPolicyDecision *models.SmPolicyDecision
	request := &SMPolicyControl.UpdateSMPolicyRequest{
		SmPolicyId:                &smContext.SMPolicyID,
		SmPolicyUpdateContextData: &updateSMPolicy,
	}

	smPolicyDecisionFromPCF, err := client.IndividualSMPolicyDocumentApi.UpdateSMPolicy(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("update sm policy [%s] association failed: %s", smContext.SMPolicyID, err)
	}
	smPolicyDecision = &smPolicyDecisionFromPCF.SmPolicyDecision
	return smPolicyDecision, nil
}

func (s *npcfService) nasBitRateToString(value uint16, unit nasType.QoSFlowBitRateUnit) string {
	var base int
	var unitStr string
	switch unit {
	case nasType.QoSFlowBitRateUnit1Kbps:
		base = 1
		unitStr = "Kbps"
	case nasType.QoSFlowBitRateUnit4Kbps:
		base = 4
		unitStr = "Kbps"
	case nasType.QoSFlowBitRateUnit16Kbps:
		base = 16
		unitStr = "Kbps"
	case nasType.QoSFlowBitRateUnit64Kbps:
		base = 64
		unitStr = "Kbps"
	case nasType.QoSFlowBitRateUnit256Kbps:
		base = 256
		unitStr = "Kbps"
	case nasType.QoSFlowBitRateUnit1Mbps:
		base = 1
		unitStr = "Mbps"
	case nasType.QoSFlowBitRateUnit4Mbps:
		base = 4
		unitStr = "Mbps"
	case nasType.QoSFlowBitRateUnit16Mbps:
		base = 16
		unitStr = "Mbps"
	case nasType.QoSFlowBitRateUnit64Mbps:
		base = 64
		unitStr = "Mbps"
	case nasType.QoSFlowBitRateUnit256Mbps:
		base = 256
		unitStr = "Mbps"
	case nasType.QoSFlowBitRateUnit1Gbps:
		base = 1
		unitStr = "Gbps"
	case nasType.QoSFlowBitRateUnit4Gbps:
		base = 4
		unitStr = "Gbps"
	case nasType.QoSFlowBitRateUnit16Gbps:
		base = 16
		unitStr = "Gbps"
	case nasType.QoSFlowBitRateUnit64Gbps:
		base = 64
		unitStr = "Gbps"
	case nasType.QoSFlowBitRateUnit256Gbps:
		base = 256
		unitStr = "Gbps"
	case nasType.QoSFlowBitRateUnit1Tbps:
		base = 1
		unitStr = "Tbps"
	case nasType.QoSFlowBitRateUnit4Tbps:
		base = 4
		unitStr = "Tbps"
	case nasType.QoSFlowBitRateUnit16Tbps:
		base = 16
		unitStr = "Tbps"
	case nasType.QoSFlowBitRateUnit64Tbps:
		base = 64
		unitStr = "Tbps"
	case nasType.QoSFlowBitRateUnit256Tbps:
		base = 256
		unitStr = "Tbps"
	case nasType.QoSFlowBitRateUnit1Pbps:
		base = 1
		unitStr = "Pbps"
	case nasType.QoSFlowBitRateUnit4Pbps:
		base = 4
		unitStr = "Pbps"
	case nasType.QoSFlowBitRateUnit16Pbps:
		base = 16
		unitStr = "Pbps"
	case nasType.QoSFlowBitRateUnit64Pbps:
		base = 64
		unitStr = "Pbps"
	case nasType.QoSFlowBitRateUnit256Pbps:
		base = 256
		unitStr = "Pbps"
	default:
		base = 1
		unitStr = "Kbps"
	}

	return fmt.Sprintf("%d %s", base*int(value), unitStr)
}

func (s *npcfService) StringToNasBitRate(str string) (uint16, nasType.QoSFlowBitRateUnit, error) {
	strSegment := strings.Split(str, " ")

	var unit nasType.QoSFlowBitRateUnit
	switch strSegment[1] {
	case "Kbps":
		unit = nasType.QoSFlowBitRateUnit1Kbps
	case "Mbps":
		unit = nasType.QoSFlowBitRateUnit1Mbps
	case "Gbps":
		unit = nasType.QoSFlowBitRateUnit1Gbps
	case "Tbps":
		unit = nasType.QoSFlowBitRateUnit1Tbps
	case "Pbps":
		unit = nasType.QoSFlowBitRateUnit1Pbps
	default:
		unit = nasType.QoSFlowBitRateUnit1Kbps
	}

	if value, err := strconv.Atoi(strSegment[0]); err != nil {
		return 0, 0, err
	} else {
		return uint16(value), unit, err
	}
}

func (s *npcfService) buildPktFilterInfo(pf nasType.PacketFilter) (*models.PacketFilterInfo, error) {
	pfInfo := &models.PacketFilterInfo{}

	switch pf.Direction {
	case nasType.PacketFilterDirectionDownlink:
		pfInfo.FlowDirection = models.FlowDirection_DOWNLINK
	case nasType.PacketFilterDirectionUplink:
		pfInfo.FlowDirection = models.FlowDirection_UPLINK
	case nasType.PacketFilterDirectionBidirectional:
		pfInfo.FlowDirection = models.FlowDirection_BIDIRECTIONAL
	default:
		pfInfo.FlowDirection = models.FlowDirection_UNSPECIFIED
	}

	const ProtocolNumberAny = 0xfc
	packetFilter := &flowdesc.IPFilterRule{
		Action: "permit",
		Dir:    "out",
		Proto:  ProtocolNumberAny,
	}

	for _, component := range pf.Components {
		switch component.Type() {
		case nasType.PacketFilterComponentTypeIPv4RemoteAddress:
			ipv4Remote := component.(*nasType.PacketFilterIPv4RemoteAddress)
			remoteIPnet := net.IPNet{
				IP:   ipv4Remote.Address,
				Mask: ipv4Remote.Mask,
			}
			packetFilter.Src = remoteIPnet.String()
		case nasType.PacketFilterComponentTypeIPv4LocalAddress:
			ipv4Local := component.(*nasType.PacketFilterIPv4LocalAddress)
			localIPnet := net.IPNet{
				IP:   ipv4Local.Address,
				Mask: ipv4Local.Mask,
			}
			packetFilter.Dst = localIPnet.String()
		case nasType.PacketFilterComponentTypeProtocolIdentifierOrNextHeader:
			protoNumber := component.(*nasType.PacketFilterProtocolIdentifier)
			packetFilter.Proto = protoNumber.Value

		case nasType.PacketFilterComponentTypeSingleLocalPort:
			localPort := component.(*nasType.PacketFilterSingleLocalPort)
			packetFilter.DstPorts = append(packetFilter.DstPorts, flowdesc.PortRange{
				Start: localPort.Value,
				End:   localPort.Value,
			})
		case nasType.PacketFilterComponentTypeLocalPortRange:
			localPortRange := component.(*nasType.PacketFilterLocalPortRange)
			packetFilter.DstPorts = append(packetFilter.DstPorts, flowdesc.PortRange{
				Start: localPortRange.LowLimit,
				End:   localPortRange.HighLimit,
			})
		case nasType.PacketFilterComponentTypeSingleRemotePort:
			remotePort := component.(*nasType.PacketFilterSingleRemotePort)
			packetFilter.SrcPorts = append(packetFilter.SrcPorts, flowdesc.PortRange{
				Start: remotePort.Value,
				End:   remotePort.Value,
			})
		case nasType.PacketFilterComponentTypeRemotePortRange:
			remotePortRange := component.(*nasType.PacketFilterRemotePortRange)
			packetFilter.SrcPorts = append(packetFilter.SrcPorts, flowdesc.PortRange{
				Start: remotePortRange.LowLimit,
				End:   remotePortRange.HighLimit,
			})
		case nasType.PacketFilterComponentTypeSecurityParameterIndex:
			securityParameter := component.(*nasType.PacketFilterSecurityParameterIndex)
			pfInfo.Spi = fmt.Sprintf("%04x", securityParameter.Index)
		case nasType.PacketFilterComponentTypeTypeOfServiceOrTrafficClass:
			serviceClass := component.(*nasType.PacketFilterServiceClass)
			pfInfo.TosTrafficClass = fmt.Sprintf("%x%x", serviceClass.Class, serviceClass.Mask)
		case nasType.PacketFilterComponentTypeFlowLabel:
			flowLabel := component.(*nasType.PacketFilterFlowLabel)
			pfInfo.FlowLabel = fmt.Sprintf("%03x", flowLabel.Label)
		}
	}

	if desc, err := flowdesc.Encode(packetFilter); err != nil {
		return nil, err
	} else {
		pfInfo.PackFiltCont = desc
	}
	// according TS 29.212 IPFilterRule cannot use [options]
	return pfInfo, nil
}

func (s *npcfService) SendSMPolicyAssociationTermination(smContext *smf_context.SMContext) error {
	var client *SMPolicyControl.APIClient

	// Create SMPolicyControl Client for this SM Context
	for _, service := range smContext.SelectedPCFProfile.NfServices {
		if service.ServiceName == models.ServiceName_NPCF_SMPOLICYCONTROL {
			client = s.getSMPolicyControlClient(service.ApiPrefix)
		}
	}

	if client == nil {
		return errors.Errorf("smContext not selected PCF")
	}

	ctx, _, err := smf_context.GetSelf().
		GetTokenCtx(models.ServiceName_NPCF_SMPOLICYCONTROL, models.NrfNfManagementNfType_PCF)
	if err != nil {
		return err
	}

	request := &SMPolicyControl.DeleteSMPolicyRequest{
		SmPolicyId:         &smContext.SMPolicyID,
		SmPolicyDeleteData: &models.SmPolicyDeleteData{},
	}

	_, err = client.IndividualSMPolicyDocumentApi.DeleteSMPolicy(ctx, request)
	if err != nil {
		return fmt.Errorf("SM Policy termination failed: %v", err)
	}
	return nil
}
