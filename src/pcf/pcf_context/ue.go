package pcf_context

import (
	"free5gc/lib/openapi/models"
)

// key is supi

type PCFUeContext struct {
	Supi   string
	Gpsi   string
	UeIpv4 string
	UeIpv6 string
	Udruri string

	// AMPolicy triggers
	UserLoc            models.UserLocation
	ServAreaRes        models.ServiceAreaRestriction
	Rfsp               int32
	Pras               map[string]models.PresenceInfoRm
	allowedSnssais     []models.Snssai
	AMFNotificationuri string
	// SMPolicy
	SmPolicyData map[string]*PCFUeSmPolicyData // use DNN as key
	// PolicyAuth
	AfRoutReq *models.AfRoutingRequirement
	AspId     string
	// Policy Decision
	AppSessionIdStore           *AppSessionIdStore
	PolAssociationIDStore       *PolAssociationIDStore
	SmPolicyControlStore        *models.SmPolicyControl
	PolicyDataSubscriptionStore *models.PolicyDataSubscription
	BdtPolicyTimeout            bool
	BdtPolicyStore              *models.BdtPolicyData
	PolicyDataChangeStore       *models.PolicyDataChangeNotification
}

type PCFUeSmPolicyData struct {
	PduSessionId       int32
	DNN                string
	SMFNotificationUri string
	SMSubscriptionData *models.SessionManagementSubscriptionData
	PccRuleIdGenarator int
	PolicyDecision     *models.SmPolicyDecision
	// PccRulePool        map[string]models.PccRule // use PccRuleId as key
}

func (ue *PCFUeContext) NewPCFUeSmPolicyData(pduSessionId int32, dnn string, smData *models.SessionManagementSubscriptionData) *PCFUeSmPolicyData {
	if pduSessionId == 0 || dnn == "" || smData == nil {
		return nil
	}
	data := PCFUeSmPolicyData{}
	data.DNN = dnn
	data.PduSessionId = pduSessionId
	data.SMSubscriptionData = smData
	data.PccRuleIdGenarator = 1
	ue.SmPolicyData[dnn] = &data
	return &data
}

// polAssoidTemp -
type PolAssociationIDStore struct {
	PolAssoId                     string
	PolAssoidTemp                 models.PolicyAssociation
	PolAssoidUpdateTemp           models.PolicyUpdate
	PolAssoidSubcCatsTemp         models.AmPolicyData
	PolAssoidDataSubscriptionTemp models.PolicyDataSubscription
	PolAssoidSubscriptiondataTemp models.SubscriptionData
}

//var polAssociationContextStore []PolAssociationIDStore

// AppSessionIdStore -
type AppSessionIdStore struct {
	AppSessionId      string
	AppSessionContext models.AppSessionContext
}

var AppSessionContextStore []AppSessionIdStore

// BdtPolicyData_store -
var BdtPolicyData_store []models.BdtPolicyData
var CreateFailBdtDateStore []models.BdtData
