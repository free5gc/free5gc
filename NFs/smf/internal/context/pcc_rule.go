package context

import (
	"fmt"
	"net"
	"strconv"

	"github.com/pkg/errors"

	"github.com/free5gc/nas/nasType"
	"github.com/free5gc/openapi/models"
	"github.com/free5gc/smf/internal/logger"
	"github.com/free5gc/smf/pkg/factory"
	"github.com/free5gc/util/flowdesc"
)

// PCCRule - Policy and Charging Rule
type PCCRule struct {
	*models.PccRule
	QFI uint8
	// related Data
	Datapath *DataPath
}

// NewPCCRule - create PCC rule from OpenAPI models
func NewPCCRule(mPcc *models.PccRule) *PCCRule {
	if mPcc == nil {
		return nil
	}

	return &PCCRule{
		PccRule: mPcc,
	}
}

func (r *PCCRule) FlowDescription() string {
	if len(r.FlowInfos) > 0 {
		// now 1 pcc rule only maps to 1 FlowInfo
		return r.FlowInfos[0].FlowDescription
	}
	return ""
}

func (r *PCCRule) RefChgDataID() string {
	if len(r.RefChgData) > 0 {
		// now 1 pcc rule only maps to 1 Charging data
		return r.RefChgData[0]
	}
	return ""
}

func (r *PCCRule) RefQosDataID() string {
	if len(r.RefQosData) > 0 {
		// now 1 pcc rule only maps to 1 QoS data
		return r.RefQosData[0]
	}
	return ""
}

func (r *PCCRule) SetQFI(qfi uint8) {
	r.QFI = qfi
}

func (r *PCCRule) RefTcDataID() string {
	if len(r.RefTcData) > 0 {
		// now 1 pcc rule only maps to 1 Traffic Control data
		return r.RefTcData[0]
	}
	return ""
}

func (r *PCCRule) IdentifyChargingLevel() (ChargingLevel, error) {
	dlIPFilterRule, err := flowdesc.Decode(r.FlowDescription())
	if err != nil {
		return 0, err
	}
	// For the PCC rule that are applicable for all datapath,
	// it's charging level will be PDU-based
	if dlIPFilterRule.Src == "any" && dlIPFilterRule.Dst == "assigned" {
		return PduSessionCharging, nil
	} else {
		// For the PCC rule that is applicable for all datapath for a datapath,
		// it's charging level will be flow-based
		return FlowCharging, nil
	}
}

func (r *PCCRule) UpdateDataPathFlowDescription(dlFlowDesc string) error {
	if r.Datapath == nil {
		return fmt.Errorf("pcc[%s]: no data path", r.PccRuleId)
	}

	if dlFlowDesc == "" {
		return fmt.Errorf("pcc[%s]: no flow description", r.PccRuleId)
	}

	ulFlowDesc := dlFlowDesc
	r.Datapath.UpdateFlowDescription(ulFlowDesc, dlFlowDesc) // UL, DL flow description should be same
	return nil
}

func (r *PCCRule) AddDataPathForwardingParameters(c *SMContext,
	tgtRoute *models.RouteToLocation,
) {
	if tgtRoute == nil {
		return
	}

	if r.Datapath == nil {
		logger.CtxLog.Warnf("AddDataPathForwardingParameters pcc[%s]: no data path", r.PccRuleId)
		return
	}

	var routeProf factory.RouteProfile
	routeProfExist := false
	// specify N6 routing information
	if tgtRoute.RouteProfId != "" {
		routeProf, routeProfExist = factory.UERoutingConfig.RouteProf[factory.RouteProfID(tgtRoute.RouteProfId)]
		if !routeProfExist {
			logger.CtxLog.Warnf("Route Profile ID [%s] is not support", tgtRoute.RouteProfId)
			return
		}
	}
	if c.Tunnel.DataPathPool.GetDefaultPath() == nil {
		logger.CtxLog.Infoln("No Default Data Path")
	} else {
		r.Datapath.AddForwardingParameters(routeProf.ForwardingPolicyID,
			c.Tunnel.DataPathPool.GetDefaultPath().FirstDPNode.GetUpLinkPDR().PDI.LocalFTeid.Teid)
	}
}

func (r *PCCRule) AddDataPathForwardingParametersOnDcTunnel(c *SMContext,
	tgtRoute *models.RouteToLocation,
) {
	if tgtRoute == nil {
		return
	}

	if r.Datapath == nil {
		logger.CtxLog.Warnf("AddDataPathForwardingParametersOnDcTunnel pcc[%s]: no data path", r.PccRuleId)
		return
	}

	var routeProf factory.RouteProfile
	routeProfExist := false
	// specify N6 routing information
	if tgtRoute.RouteProfId != "" {
		routeProf, routeProfExist = factory.UERoutingConfig.RouteProf[factory.RouteProfID(tgtRoute.RouteProfId)]
		if !routeProfExist {
			logger.CtxLog.Warnf("Route Profile ID [%s] is not support on DCTunnel", tgtRoute.RouteProfId)
			return
		}
	}

	if c.DCTunnel.DataPathPool.GetDefaultPath() == nil {
		logger.CtxLog.Infoln("No Default Data Path")
	} else {
		r.Datapath.AddForwardingParameters(routeProf.ForwardingPolicyID,
			c.DCTunnel.DataPathPool.GetDefaultPath().FirstDPNode.GetUpLinkPDR().PDI.LocalFTeid.Teid)
	}
}

func (r *PCCRule) BuildNasQoSRule(smCtx *SMContext,
	opCode nasType.QoSRuleOperationCode,
) (*nasType.QoSRule, error) {
	rule := nasType.QoSRule{}
	rule.Operation = nasType.OperationCodeCreateNewQoSRule
	rule.Precedence = uint8(r.Precedence)
	pfList := make(nasType.PacketFilterList, 0)
	for _, flowInfo := range r.FlowInfos {
		if pfs, err := BuildNASPacketFiltersFromFlowInformation(&flowInfo, smCtx); err != nil {
			logger.CtxLog.Warnf("BuildNasQoSRule: Build packet filter fail: %s\n", err)
			continue
		} else {
			pfList = append(pfList, pfs...)
		}
	}
	rule.PacketFilterList = pfList
	rule.QFI = r.QFI

	return &rule, nil
}

func createNasPacketFilter(
	pfInfo *models.FlowInformation,
	smCtx *SMContext,
	ipFilterRule *flowdesc.IPFilterRule,
	srcP *flowdesc.PortRange,
	dstP *flowdesc.PortRange,
) (*nasType.PacketFilter, error) {
	pf := new(nasType.PacketFilter)

	pfId, errAllocate := smCtx.PacketFilterIDGenerator.Allocate()
	if errAllocate != nil {
		return nil, errAllocate
	}
	pf.Identifier = uint8(pfId)
	smCtx.PacketFilterIDToNASPFID[pfInfo.PackFiltId] = uint8(pfId)

	switch pfInfo.FlowDirection {
	case models.FlowDirection_DOWNLINK:
		pf.Direction = nasType.PacketFilterDirectionDownlink
	case models.FlowDirection_UPLINK:
		pf.Direction = nasType.PacketFilterDirectionUplink
	case models.FlowDirection_BIDIRECTIONAL:
		pf.Direction = nasType.PacketFilterDirectionBidirectional
	}

	pfComponents := make(nasType.PacketFilterComponentList, 0)
	if pfInfo.FlowLabel != "" {
		if label, parseErr := strconv.ParseInt(pfInfo.FlowLabel, 16, 32); parseErr != nil {
			return nil, fmt.Errorf("parse flow label fail: %s", parseErr)
		} else {
			pfComponents = append(pfComponents, &nasType.PacketFilterFlowLabel{
				Label: uint32(label),
			})
		}
	}

	if pfInfo.Spi != "" {
		if spi, parseErr := strconv.ParseInt(pfInfo.Spi, 16, 32); parseErr != nil {
			return nil, fmt.Errorf("parse security parameter index fail: %s", parseErr)
		} else {
			pfComponents = append(pfComponents, &nasType.PacketFilterSecurityParameterIndex{
				Index: uint32(spi),
			})
		}
	}

	if pfInfo.TosTrafficClass != "" {
		if tos, parseErr := strconv.ParseInt(pfInfo.TosTrafficClass, 16, 32); parseErr != nil {
			return nil, fmt.Errorf("parse security parameter index fail: %s", parseErr)
		} else {
			pfComponents = append(pfComponents, &nasType.PacketFilterServiceClass{
				Class: uint8(tos >> 8),
				Mask:  uint8(tos & 0x00FF),
			})
		}
	}

	if ipFilterRule.Dst != "assigned" {
		_, ipNet, errParseCIDR := net.ParseCIDR(ipFilterRule.Dst)
		if errParseCIDR != nil {
			return nil, fmt.Errorf("parse IP fail: %s", errParseCIDR)
		}
		pfComponents = append(pfComponents, &nasType.PacketFilterIPv4LocalAddress{
			Address: ipNet.IP.To4(),
			Mask:    ipNet.Mask,
		})
	}
	if dstP != nil {
		if dstP.Start != dstP.End {
			pfComponents = append(pfComponents, &nasType.PacketFilterLocalPortRange{
				LowLimit:  dstP.Start,
				HighLimit: dstP.End,
			})
		} else if dstP.Start != 0 && dstP.End != 0 {
			pfComponents = append(pfComponents, &nasType.PacketFilterSingleLocalPort{
				Value: dstP.Start,
			})
		}
	}

	if ipFilterRule.Src != "any" {
		_, ipNet, errParseCIDR := net.ParseCIDR(ipFilterRule.Src)
		if errParseCIDR != nil {
			return nil, fmt.Errorf("parse IP fail: %s", errParseCIDR)
		}
		pfComponents = append(pfComponents, &nasType.PacketFilterIPv4RemoteAddress{
			Address: ipNet.IP.To4(),
			Mask:    ipNet.Mask,
		})
	}
	if srcP != nil {
		if srcP.Start != srcP.End {
			pfComponents = append(pfComponents, &nasType.PacketFilterRemotePortRange{
				LowLimit:  srcP.Start,
				HighLimit: srcP.End,
			})
		} else if srcP.Start != 0 && srcP.End != 0 {
			pfComponents = append(pfComponents, &nasType.PacketFilterSingleRemotePort{
				Value: srcP.Start,
			})
		}
	}

	if ipFilterRule.Proto != flowdesc.ProtocolNumberAny {
		pfComponents = append(pfComponents, &nasType.PacketFilterProtocolIdentifier{
			Value: ipFilterRule.Proto,
		})
	}

	if len(pfComponents) == 0 {
		pfComponents = append(pfComponents, &nasType.PacketFilterMatchAll{})
	}

	pf.Components = pfComponents
	return pf, nil
}

func BuildNASPacketFiltersFromFlowInformation(pfInfo *models.FlowInformation,
	smCtx *SMContext,
) ([]nasType.PacketFilter, error) {
	var pfList []nasType.PacketFilter

	ipFilterRule := flowdesc.NewIPFilterRule()
	if pfInfo.FlowDescription != "" {
		var err error
		ipFilterRule, err = flowdesc.Decode(pfInfo.FlowDescription)
		if err != nil {
			return nil, fmt.Errorf("parse packet filter content fail: %s", err)
		}
	}

	// TS 24.501 9.11.4.13.4
	srcPLen := len(ipFilterRule.SrcPorts)
	dstPLen := len(ipFilterRule.DstPorts)
	switch {
	case srcPLen > 0 && dstPLen > 0:
		for _, srcP := range ipFilterRule.SrcPorts {
			for _, dstP := range ipFilterRule.DstPorts {
				pf, err := createNasPacketFilter(pfInfo, smCtx, ipFilterRule, &srcP, &dstP)
				if err != nil {
					return nil, errors.Wrap(err, "create packet filter fail")
				}
				pfList = append(pfList, *pf)
			}
		}
	case srcPLen == 0 && dstPLen > 0:
		for _, dstP := range ipFilterRule.DstPorts {
			pf, err := createNasPacketFilter(pfInfo, smCtx, ipFilterRule, nil, &dstP)
			if err != nil {
				return nil, errors.Wrap(err, "create packet filter fail")
			}
			pfList = append(pfList, *pf)
		}
	case srcPLen > 0 && dstPLen == 0:
		for _, srcP := range ipFilterRule.SrcPorts {
			pf, err := createNasPacketFilter(pfInfo, smCtx, ipFilterRule, &srcP, nil)
			if err != nil {
				return nil, errors.Wrap(err, "create packet filter fail")
			}
			pfList = append(pfList, *pf)
		}
	case srcPLen == 0 && dstPLen == 0:
		pf, err := createNasPacketFilter(pfInfo, smCtx, ipFilterRule, nil, nil)
		if err != nil {
			return nil, errors.Wrap(err, "create packet filter fail")
		}
		pfList = append(pfList, *pf)
	default:
		return nil, errors.Errorf("invalid srcPLen(%d) or dstPLen(%d)", srcPLen, dstPLen)
	}

	return pfList, nil
}
