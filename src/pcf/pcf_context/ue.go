package pcf_context

import (
	"fmt"
	"free5gc/lib/Nudr_DataRepository"
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
	SmPolicyData map[string]*UeSmPolicyData // use PolAssoId(ue.Supi-pduSessionId) as key
	// PolicyAuth
	AfRoutReq *models.AfRoutingRequirement
	AspId     string
	// Policy Decision
	AppSessionIdStore           *AppSessionIdStore
	SmPolicyControlStore        *models.SmPolicyControl
	PolicyDataSubscriptionStore *models.PolicyDataSubscription
	BdtPolicyTimeout            bool
	BdtPolicyStore              *models.BdtPolicyData
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
	PackFiltIdGenarator    int32
	PccRuleIdGenarator     int32
	PackFiltMapToPccRuleId map[string]string // use PackFiltId as Key
	RemainGbrUL            *float64
	RemainGbrDL            *float64
	SmPolicyData           *models.SmPolicyData
	PolicyContext          *models.SmPolicyContextData
	PolicyDecision         *models.SmPolicyDecision
	// PccRulePool        map[string]models.PccRule // use PccRuleId as key
	// RefToAmPolicy          *UeAMPolicyData
}

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
		SuppFeat:          req.SuppFeat,
	}
	ue.AMPolicyData[assolId].Pras = make(map[string]models.PresenceInfo)
	return ue.AMPolicyData[assolId]
}

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
	// data.RefToAmPolicy = amData
	data.PccRuleIdGenarator = 1
	ue.SmPolicyData[key] = &data
	return &data
}

func (ue *UeContext) GetNudrClient() *Nudr_DataRepository.APIClient {
	if ue.UdrUri == "" {
		return nil
	}
	configuration := Nudr_DataRepository.NewConfiguration()
	BasePath := ue.UdrUri
	configuration.SetBasePath(BasePath)
	client := Nudr_DataRepository.NewAPIClient(configuration)
	return client
}

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

func (policy *UeSmPolicyData) IncreaseRemainGBR(qosId string) (origUl, origDl *float64) {
	decision := policy.PolicyDecision
	if decision == nil {
		return
	}
	if qos, exist := decision.QosDecs[qosId]; exist {
		if qos.Var5qi <= 4 {
			// Add GBR
			if policy.RemainGbrDL != nil && qos.GbrDl != "" {
				bitRate, err := ConvertBitRateToKbps(qos.GbrDl)
				if err == nil {
					origDl = new(float64)
					*origDl = *policy.RemainGbrDL
					*policy.RemainGbrDL += bitRate
				}
			}
			if policy.RemainGbrUL != nil {
				bitRate, err := ConvertBitRateToKbps(qos.GbrUl)
				if err == nil {
					origUl = new(float64)
					*origUl = *policy.RemainGbrUL
					*policy.RemainGbrUL += bitRate
				}
			}
		}
	}
	return
}

func (policy *UeSmPolicyData) DecreaseRemainGBR(req *models.RequestedQos, remainGbrDL, remainGbrUL *float64) (gbrDl, gbrUl string, err error) {
	if req == nil {
		return "", "", nil
	}
	if req.Var5qi <= 4 {
		if req.GbrDl != "" {
			gbrDL, err1 := ConvertBitRateToKbps(req.GbrDl)
			if err1 != nil {
				err = err1
				return
			} else if remainGbrDL != nil {
				if *remainGbrDL < gbrDL {
					err = fmt.Errorf("Request DL GBR exceed Dnn Aggregate DL GBR of UE")
					return
				}
				*remainGbrDL -= gbrDL
				gbrDl = req.GbrDl
			}
		}
		if req.GbrUl != "" {
			gbrUL, err1 := ConvertBitRateToKbps(req.GbrUl)
			if err1 != nil {
				err = err1
				return
			} else if remainGbrUL != nil {
				if *remainGbrUL < gbrUL {
					err = fmt.Errorf("Request DL GBR exceed Dnn Aggregate DL GBR of UE")
					return
				}
				*remainGbrUL -= gbrUL
				gbrUl = req.GbrUl
			}
		}
	}
	return
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

func ConvertBitRateToKbps(bitRate string) (mBitRate float64, err error) {
	index := strings.Index(bitRate, "bps")
	if index == -1 || index == 0 {
		err := fmt.Errorf("bitRate format error")
		return 0, err
	}
	tmp := bitRate[:index]
	lastIndex := len(tmp) - 1
	times := 0.0
	switch tmp[lastIndex] {
	case 'T':
		times = 30.0
	case 'G':
		times = 20.0
	case 'M':
		times = 10.0
	case 'K':
		times = 0.0
	default:
		lastIndex = lastIndex + 1
		times = -10.0
	}
	tmp = tmp[:lastIndex]
	mBitRate, err = strconv.ParseFloat(tmp, 64)
	if err == nil {
		mBitRate = mBitRate * math.Pow(2, times)
	} else {
		mBitRate = 0.0
	}
	return
}
