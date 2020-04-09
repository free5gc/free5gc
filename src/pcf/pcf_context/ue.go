package pcf_context

import (
	"fmt"
	"free5gc/lib/openapi/models"
	"math"
	"reflect"
	"strconv"
	"strings"
)

// key is supi
type UeContext struct {
	// Ue Context
	Supi                      string
	Gpsi                      string
	Pei                       string
	GroupIds                  []string
	PolAssociationIDGenerator uint32
	AMPolicyData              map[string]*UeAMPolicyData // use PolAssoId(ue.Supi-numPolId) as key

	// Udr Ref
	UdrUri string
	// SMPolicy
	SmPolicyData map[string]*UeSmPolicyData // use smPolicyId(ue.Supi-pduSessionId) as key
	// App Session Related
	AppSessionIdGenerator uint64
	// PolicyAuth
	AfRoutReq *models.AfRoutingRequirement
	AspId     string
	// Policy Decision
	AppSessionIdStore           *AppSessionIdStore
	PolicyDataSubscriptionStore *models.PolicyDataSubscription
	PolicyDataChangeStore       *models.PolicyDataChangeNotification
}

type UeAMPolicyData struct {
	PolAssoId         string
	AccessType        models.AccessType
	NotificationUri   string
	ServingPlmn       *models.NetworkId
	AltNotifIpv4Addrs []string
	AltNotifIpv6Addrs []string
	// TODO: AMF Status Change
	AmfStatusUri string
	Guami        *models.Guami
	ServiveName  string
	// TraceReq *TraceData
	// Policy Association
	Triggers    []models.RequestTrigger
	ServAreaRes *models.ServiceAreaRestriction
	Rfsp        int32
	UserLoc     *models.UserLocation
	TimeZone    string
	SuppFeat    string
	// about AF request
	Pras map[string]models.PresenceInfo
	// related to UDR Subscription Data
	AmPolicyData *models.AmPolicyData // Svbscription Data
	// Corresponding UE
	PcfUe *UeContext
	// AMF status change subscription
	AmfStatusChangeSubscription *AMFStatusSubscriptionData
}

type UeSmPolicyData struct {
	// PduSessionId    int32
	// DNN             string
	// NotificationUri string
	// Snssai                 models.Snssai
	// PduSessionType         models.PduSessionType
	// IPAddress              models.IpAddress
	// IPDomain               string
	// Var3gppPsDataOffStatus bool
	// SmfId                  string
	// TraceReq *TraceData
	// RecoveryTime     *time.Time
	PackFiltIdGenarator int32
	PccRuleIdGenarator  int32
	ChargingIdGenarator int32
	// FlowMapsToPackFiltIds  map[string][]string // use Flow Description(in TS 29214) as key map to pcc rule ids
	PackFiltMapToPccRuleId map[string]string // use PackFiltId as Key
	// Related to GBR
	RemainGbrUL *float64
	RemainGbrDL *float64
	// related to UDR Subscription Data
	SmPolicyData *models.SmPolicyData // Svbscription Data
	// related to Policy
	PolicyContext  *models.SmPolicyContextData
	PolicyDecision *models.SmPolicyDecision
	// related to AppSession
	AppSessions map[string]bool // related appSessionId
	// Corresponding UE
	PcfUe *UeContext
}

// NewUeAMPolicyData returns created UeAMPolicyData data and insert this data to Ue.AMPolicyData with assolId as key
func (ue *UeContext) NewUeAMPolicyData(assolId string, req models.PolicyAssociationRequest) *UeAMPolicyData {
	ue.Gpsi = req.Gpsi
	ue.Pei = req.Pei
	ue.GroupIds = req.GroupIds
	ue.AMPolicyData[assolId] = &UeAMPolicyData{
		PolAssoId:         assolId,
		ServAreaRes:       req.ServAreaRes,
		AltNotifIpv4Addrs: req.AltNotifIpv4Addrs,
		AltNotifIpv6Addrs: req.AltNotifIpv6Addrs,
		AccessType:        req.AccessType,
		NotificationUri:   req.NotificationUri,
		ServingPlmn:       req.ServingPlmn,
		TimeZone:          req.TimeZone,
		Rfsp:              req.Rfsp,
		Guami:             req.Guami,
		UserLoc:           req.UserLoc,
		ServiveName:       req.ServiveName,
		PcfUe:             ue,
	}
	ue.AMPolicyData[assolId].Pras = make(map[string]models.PresenceInfo)
	return ue.AMPolicyData[assolId]
}

// returns UeSmPolicyData and insert related info to Ue with smPolId
func (ue *UeContext) NewUeSmPolicyData(key string, request models.SmPolicyContextData, smData *models.SmPolicyData) *UeSmPolicyData {
	if smData == nil {
		return nil
	}
	data := UeSmPolicyData{}
	data.PolicyContext = &request
	// data.DNN = request.Dnn
	// data.Snssai = *request.SliceInfo
	// data.PduSessionId = request.PduSessionId
	// data.PduSessionType = request.PduSessionType
	// switch request.PduSessionType {
	// case models.PduSessionType_IPV4:
	// 	data.IPAddress.Ipv4Addr = request.Ipv4Address
	// 	data.IPDomain = request.IpDomain
	// case models.PduSessionType_IPV6:
	// 	data.IPAddress.Ipv6Prefix = request.Ipv6AddressPrefix
	// case models.PduSessionType_IPV4_V6:
	// 	data.IPAddress.Ipv4Addr = request.Ipv4Address
	// 	data.IPAddress.Ipv6Prefix = request.Ipv6AddressPrefix
	// 	data.IPDomain = request.IpDomain
	// }
	// data.NotificationUri = request.NotificationUri
	// data.SmfId = request.SmfId
	// data.Var3gppPsDataOffStatus = request.Var3gppPsDataOffStatus
	data.SmPolicyData = smData
	data.PackFiltIdGenarator = 1
	data.PackFiltMapToPccRuleId = make(map[string]string)
	data.AppSessions = make(map[string]bool)
	// data.RefToAmPolicy = amData
	data.PccRuleIdGenarator = 1
	data.ChargingIdGenarator = 1
	data.PcfUe = ue
	ue.SmPolicyData[key] = &data
	return &data
}

// Remove Pcc rule wich PccRuleId in the policy
func (policy *UeSmPolicyData) RemovePccRule(pccRuleId string) error {
	decision := policy.PolicyDecision
	if decision == nil {
		return fmt.Errorf("Can't find the Policy Decision")
	}
	if rule, exist := decision.PccRules[pccRuleId]; exist {
		for _, info := range rule.FlowInfos {
			delete(policy.PackFiltMapToPccRuleId, info.PackFiltId)
		}
		for _, id := range rule.RefQosData {
			if decision.QosDecs != nil {
				policy.IncreaseRemainGBR(id)
				delete(decision.QosDecs, id)
				if len(decision.QosDecs) == 0 {
					decision.QosDecs = nil
				}
			} else {
				break
			}
		}
		if rule.RefCondData != "" {
			if decision.Conds != nil {
				delete(decision.Conds, rule.RefCondData)
				if len(decision.Conds) == 0 {
					decision.Conds = nil
				}
			}
		}
		for _, id := range rule.RefChgData {
			if decision.ChgDecs != nil {
				delete(decision.ChgDecs, id)
				if len(decision.ChgDecs) == 0 {
					decision.ChgDecs = nil
				}
			} else {
				break
			}
		}
		for _, id := range rule.RefTcData {
			if decision.TraffContDecs != nil {
				delete(decision.TraffContDecs, id)
				if len(decision.TraffContDecs) == 0 {
					decision.TraffContDecs = nil
				}
			} else {
				break
			}
		}
		for _, id := range rule.RefUmData {
			if decision.UmDecs != nil {
				delete(decision.UmDecs, id)
				if len(decision.UmDecs) == 0 {
					decision.UmDecs = nil
				}
			} else {
				break
			}
		}
		delete(decision.PccRules, pccRuleId)
	} else {
		return fmt.Errorf("Can't find the pccRuleId[%s] in Session[%d]", pccRuleId, policy.PolicyContext.PduSessionId)
	}
	return nil
}

// Check if the afEvent exists in smPolicy
func (policy *UeSmPolicyData) CheckRelatedAfEvent(event models.AfEvent) (found bool) {
	for appSessionId := range policy.AppSessions {
		if appSession, exist := PCF_Self().AppSessionPool[appSessionId]; exist {
			for afEvent := range appSession.Events {
				if afEvent == event {
					return true
				}
			}
		}
	}
	return false
}

// Arrange Exist Event policy Sm policy about afevents and return if it changes or not and
func (policy *UeSmPolicyData) ArrangeExistEventSubscription() (changed bool) {
	triggers := []models.PolicyControlRequestTrigger{}
	for _, trigger := range policy.PolicyDecision.PolicyCtrlReqTriggers {
		var afEvent models.AfEvent
		switch trigger {
		case models.PolicyControlRequestTrigger_PLMN_CH: // PLMN Change
			afEvent = models.AfEvent_PLMN_CHG
		case models.PolicyControlRequestTrigger_QOS_NOTIF: // SMF notify PCF when receiving from RAN that QoS can/can't be guaranteed (subsclause 4.2.4.20 in TS29512) (always)
			afEvent = models.AfEvent_QOS_NOTIF
		case models.PolicyControlRequestTrigger_SUCC_RES_ALLO: // Successful resource allocation (subsclause 4.2.6.5.5, 4.2.4.14 in TS29512)
			afEvent = models.AfEvent_SUCCESSFUL_RESOURCES_ALLOCATION
		case models.PolicyControlRequestTrigger_AC_TY_CH: // Change of RatType
			afEvent = models.AfEvent_ACCESS_TYPE_CHANGE
		case models.PolicyControlRequestTrigger_US_RE: // UMC
			afEvent = models.AfEvent_USAGE_REPORT
		}
		if afEvent != "" && !policy.CheckRelatedAfEvent(afEvent) {
			changed = true
		} else {
			triggers = append(triggers, trigger)
		}
	}
	policy.PolicyDecision.PolicyCtrlReqTriggers = triggers
	return
}

// Increase remain GBR of this policy and returns original UL DL GBR for resume case
func (policy *UeSmPolicyData) IncreaseRemainGBR(qosId string) (origUl, origDl *float64) {
	decision := policy.PolicyDecision
	if decision == nil {
		return
	}
	if qos, exist := decision.QosDecs[qosId]; exist {
		if qos.Var5qi <= 4 {
			// Add GBR
			origUl = IncreaseRamainBitRate(policy.RemainGbrUL, qos.GbrUl)
			origDl = IncreaseRamainBitRate(policy.RemainGbrDL, qos.GbrDl)
		}
	}
	return
}

// Increase remain Bit Rate and returns original Bit Rate
func IncreaseRamainBitRate(remainBitRate *float64, reqBitRate string) (orig *float64) {
	if remainBitRate != nil && reqBitRate != "" {
		bitRate, err := ConvertBitRateToKbps(reqBitRate)
		if err == nil {
			orig = new(float64)
			*orig = *remainBitRate
			*remainBitRate += bitRate
		}
	}
	return
}

// Decrease remain GBR of this policy and returns UL DL GBR
func (policy *UeSmPolicyData) DecreaseRemainGBR(req *models.RequestedQos) (gbrDl, gbrUl string, err error) {
	if req == nil {
		return "", "", nil
	}
	if req.Var5qi <= 4 {
		err = DecreaseRamainBitRate(policy.RemainGbrDL, req.GbrDl)
		if err != nil {
			return
		}
		gbrDl = req.GbrDl
		err = DecreaseRamainBitRate(policy.RemainGbrUL, req.GbrUl)
		if err != nil {
			return
		}
		gbrUl = req.GbrUl
	}
	return
}

// Decrease remain Bit Rate
func DecreaseRamainBitRate(remainBitRate *float64, reqBitRate string) error {
	if reqBitRate != "" {
		bitRate, err := ConvertBitRateToKbps(reqBitRate)
		if err != nil {
			return err
		}
		if remainBitRate != nil {
			if *remainBitRate < bitRate {
				return fmt.Errorf("Request BitRate exceed Dnn Aggregate BitRate of UE")
			}
			*remainBitRate -= bitRate
		}
	}
	return nil
}

// Returns remin Bit rate string and decrease ir to zero
func DecreaseRamainBitRateToZero(remainBitRate *float64) string {
	if remainBitRate != nil {
		bitRate := ConvertBitRateToString(*remainBitRate)
		*remainBitRate = 0
		return bitRate
	}
	return ""
}

// returns AM Policy which AccessType and plmnId match
func (ue *UeContext) FindAMPolicy(anType models.AccessType, plmnId *models.NetworkId) *UeAMPolicyData {
	if ue == nil || plmnId == nil {
		return nil
	}
	for _, amPolicy := range ue.AMPolicyData {
		if amPolicy.AccessType == anType && reflect.DeepEqual(*amPolicy.ServingPlmn, *plmnId) {
			return amPolicy
		}
	}
	return nil
}

// Return App Session Id with format "ue.Supi-%d" which be allocated
func (ue *UeContext) AllocUeAppSessionId(context *PCFContext) string {
	appSessionId := fmt.Sprintf("%s-%d", ue.Supi, ue.AppSessionIdGenerator)
	_, exist := context.AppSessionPool[appSessionId]
	for exist {
		ue.AppSessionIdGenerator++
		appSessionId = fmt.Sprintf("%s-%d", ue.Supi, ue.AppSessionIdGenerator)
		_, exist = context.AppSessionPool[appSessionId]
	}
	ue.AppSessionIdGenerator++
	return appSessionId
}

// returns SM Policy by IPv4
func (ue *UeContext) SMPolicyFindByIpv4(v4 string) *UeSmPolicyData {
	for _, smPolicy := range ue.SmPolicyData {
		if smPolicy.PolicyContext.Ipv4Address == v4 {
			return smPolicy
		}
	}
	return nil
}

// returns SM Policy by IPv6
func (ue *UeContext) SMPolicyFindByIpv6(v6 string) *UeSmPolicyData {
	for _, smPolicy := range ue.SmPolicyData {
		if smPolicy.PolicyContext.Ipv6AddressPrefix == v6 {
			return smPolicy
		}
	}
	return nil
}

// AppSessionIdStore -
type AppSessionIdStore struct {
	AppSessionId      string
	AppSessionContext models.AppSessionContext
}

var AppSessionContextStore []AppSessionIdStore

// BdtPolicyData_store -
var BdtPolicyData_store []models.BdtPolicyData
var CreateFailBdtDateStore []models.BdtData

// Convert bitRate string to float64 with uint Kbps
func ConvertBitRateToKbps(bitRate string) (kBitRate float64, err error) {
	list := strings.Split(bitRate, " ")
	if len(list) != 2 {
		err := fmt.Errorf("bitRate format error")
		return 0, err
	}
	// parse exponential value with 2 as base
	exp := 0.0
	switch list[1] {
	case "Tbps":
		exp = 30.0
	case "Gbps":
		exp = 20.0
	case "Mbps":
		exp = 10.0
	case "Kbps":
		exp = 0.0
	case "bps":
		exp = -10.0
	default:
		err := fmt.Errorf("bitRate format error")
		return 0, err
	}
	// parse value from string to float64
	kBitRate, err = strconv.ParseFloat(list[0], 64)
	if err == nil {
		kBitRate = kBitRate * math.Pow(2, exp)
	} else {
		kBitRate = 0.0
	}
	return
}

// Convert bitRate from float64 to String
func ConvertBitRateToString(kBitRate float64) (bitRate string) {
	return fmt.Sprintf("%f Kbps", kBitRate)
}
