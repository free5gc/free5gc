/*
 * BSF Context Management
 */

package context

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/free5gc/bsf/internal/logger"
	"github.com/free5gc/bsf/pkg/factory"
	"github.com/free5gc/openapi/models"
)

var (
	bsfContext = BSFContext{}
	BsfSelf    = &bsfContext
)

func init() {
	BsfSelf.Name = "bsf"
	BsfSelf.NfId = uuid.New().String()
	BsfSelf.PcfBindings = make(map[string]*PcfBinding)
	BsfSelf.PcfForUeBindings = make(map[string]*PcfForUeBinding)
	BsfSelf.PcfMbsBindings = make(map[string]*PcfMbsBinding)
	BsfSelf.Subscriptions = make(map[string]*BsfSubscription)
	BsfSelf.mutex = sync.RWMutex{}
}

type BSFContext struct {
	mutex        sync.RWMutex
	NfId         string
	Name         string
	UriScheme    string
	RegisterIPv4 string
	SBIPort      int
	BindingIPv4  string
	NrfUri       string

	// MongoDB
	MongoDBName string
	MongoDBUrl  string

	// BSF Business Logic
	PcfBindings      map[string]*PcfBinding      // bindingId -> PcfBinding
	PcfForUeBindings map[string]*PcfForUeBinding // bindingId -> PcfForUeBinding
	PcfMbsBindings   map[string]*PcfMbsBinding   // bindingId -> PcfMbsBinding
	Subscriptions    map[string]*BsfSubscription // subId -> BsfSubscription
}

type PcfBinding struct {
	BindingId          string
	Supi               *string
	Gpsi               *string
	Ipv4Addr           *string
	Ipv6Prefix         *string
	AddIpv6Prefixes    []string
	IpDomain           *string
	MacAddr48          *string
	AddMacAddrs        []string
	Dnn                string
	PcfFqdn            *string
	PcfIpEndPoints     []models.IpEndPoint
	PcfDiamHost        *string
	PcfDiamRealm       *string
	PcfSmFqdn          *string
	PcfSmIpEndPoints   []models.IpEndPoint
	Snssai             *models.Snssai
	SuppFeat           *string
	PcfId              *string
	PcfSetId           *string
	RecoveryTime       *time.Time
	ParaCom            *models.ParameterCombination
	BindLevel          *models.BindingLevel
	Ipv4FrameRouteList []string
	Ipv6FrameRouteList []string
}

type PcfForUeBinding struct {
	BindingId           string
	Supi                string
	Gpsi                *string
	PcfForUeFqdn        *string
	PcfForUeIpEndPoints []models.IpEndPoint
	PcfId               *string
	PcfSetId            *string
	BindLevel           *models.BindingLevel
	SuppFeat            *string
}

type PcfMbsBinding struct {
	BindingId      string
	MbsSessionId   *models.MbsSessionId
	PcfFqdn        *string
	PcfIpEndPoints []models.IpEndPoint
	PcfId          *string
	PcfSetId       *string
	BindLevel      *models.BindingLevel
	RecoveryTime   *time.Time
	SuppFeat       *string
}

type BsfSubscription struct {
	SubId             string
	Events            []models.BsfEvent
	NotifUri          string
	NotifCorreId      string
	Supi              string
	Gpsi              *string
	SnssaiDnnPairs    *models.SnssaiDnnPair
	AddSnssaiDnnPairs []models.SnssaiDnnPair
	SuppFeat          *string
}

func (c *BSFContext) GetSelf() *BSFContext {
	return &bsfContext
}

func InitBsfContext() {
	config := factory.BsfConfig
	if config == nil {
		logger.CtxLog.Error("Config is nil")
		return
	}

	logger.CtxLog.Infof("bsfconfig Info: Version[%s] Description[%s]", config.Info.Version, config.Info.Description)

	configuration := config.Configuration
	if configuration.Sbi == nil {
		logger.CtxLog.Errorln("Configuration needs \"sbi\" value")
		return
	} else {
		BsfSelf.RegisterIPv4 = configuration.Sbi.RegisterIPv4
		BsfSelf.SBIPort = configuration.Sbi.Port
		BsfSelf.BindingIPv4 = configuration.Sbi.BindingIPv4
		BsfSelf.UriScheme = configuration.Sbi.Scheme
		if configuration.Sbi.Tls != nil {
			logger.CtxLog.Infoln("TLS enabled")
		}
	}

	if configuration.NrfUri != "" {
		BsfSelf.NrfUri = configuration.NrfUri
	} else {
		logger.CtxLog.Warn("NRF Uri is empty! BSF will not register to NRF")
	}

	if configuration.MongoDB == nil {
		logger.CtxLog.Warn("MongoDB is nil")
	} else {
		BsfSelf.MongoDBName = configuration.MongoDB.Name
		BsfSelf.MongoDBUrl = configuration.MongoDB.Url
	}
}

// MongoDB related functions
func (c *BSFContext) ConnectMongoDB() error {
	if c.MongoDBUrl == "" {
		return fmt.Errorf("MongoDB URL is empty")
	}

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(c.MongoDBUrl))
	if err != nil {
		return fmt.Errorf("failed to connect to MongoDB: %+v", err)
	}

	// Test the connection
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		return fmt.Errorf("failed to ping MongoDB: %+v", err)
	}

	logger.CtxLog.Infof("Connected to MongoDB at %s", c.MongoDBUrl)
	return nil
}

// PCF Binding Management
func (c *BSFContext) CreatePcfBinding(binding *PcfBinding) string {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	bindingId := uuid.New().String()
	binding.BindingId = bindingId
	c.PcfBindings[bindingId] = binding

	logger.CtxLog.Debugf("Created PCF binding with ID: %s", bindingId)

	return bindingId
}

func (c *BSFContext) GetPcfBinding(bindingId string) (*PcfBinding, bool) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	binding, exists := c.PcfBindings[bindingId]
	return binding, exists
}

func (c *BSFContext) UpdatePcfBinding(bindingId string, binding *PcfBinding) bool {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if _, exists := c.PcfBindings[bindingId]; exists {
		binding.BindingId = bindingId
		c.PcfBindings[bindingId] = binding
		logger.CtxLog.Debugf("Updated PCF binding with ID: %s", bindingId)
		return true
	}
	return false
}

func (c *BSFContext) DeletePcfBinding(bindingId string) bool {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if _, exists := c.PcfBindings[bindingId]; exists {
		delete(c.PcfBindings, bindingId)
		logger.CtxLog.Debugf("Deleted PCF binding with ID: %s", bindingId)

		return true
	}
	return false
}

// PCF UE Binding Management
func (c *BSFContext) CreatePcfForUeBinding(binding *PcfForUeBinding) string {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	bindingId := uuid.New().String()
	binding.BindingId = bindingId
	c.PcfForUeBindings[bindingId] = binding

	logger.CtxLog.Debugf("Created PCF UE binding with ID: %s", bindingId)
	return bindingId
}

func (c *BSFContext) GetPcfForUeBinding(bindingId string) (*PcfForUeBinding, bool) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	binding, exists := c.PcfForUeBindings[bindingId]
	return binding, exists
}

func (c *BSFContext) UpdatePcfForUeBinding(bindingId string, binding *PcfForUeBinding) bool {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if _, exists := c.PcfForUeBindings[bindingId]; exists {
		binding.BindingId = bindingId
		c.PcfForUeBindings[bindingId] = binding
		logger.CtxLog.Debugf("Updated PCF UE binding with ID: %s", bindingId)
		return true
	}
	return false
}

func (c *BSFContext) DeletePcfForUeBinding(bindingId string) bool {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if _, exists := c.PcfForUeBindings[bindingId]; exists {
		delete(c.PcfForUeBindings, bindingId)
		logger.CtxLog.Debugf("Deleted PCF UE binding with ID: %s", bindingId)
		return true
	}
	return false
}

// PCF MBS Binding Management
func (c *BSFContext) CreatePcfMbsBinding(binding *PcfMbsBinding) string {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	bindingId := uuid.New().String()
	binding.BindingId = bindingId
	c.PcfMbsBindings[bindingId] = binding

	logger.CtxLog.Debugf("Created PCF MBS binding with ID: %s", bindingId)
	return bindingId
}

func (c *BSFContext) GetPcfMbsBinding(bindingId string) (*PcfMbsBinding, bool) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	binding, exists := c.PcfMbsBindings[bindingId]
	return binding, exists
}

func (c *BSFContext) UpdatePcfMbsBinding(bindingId string, binding *PcfMbsBinding) bool {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if _, exists := c.PcfMbsBindings[bindingId]; exists {
		binding.BindingId = bindingId
		c.PcfMbsBindings[bindingId] = binding
		logger.CtxLog.Debugf("Updated PCF MBS binding with ID: %s", bindingId)
		return true
	}
	return false
}

func (c *BSFContext) DeletePcfMbsBinding(bindingId string) bool {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if _, exists := c.PcfMbsBindings[bindingId]; exists {
		delete(c.PcfMbsBindings, bindingId)
		logger.CtxLog.Debugf("Deleted PCF MBS binding with ID: %s", bindingId)
		return true
	}
	return false
}

// Subscription Management
func (c *BSFContext) CreateSubscription(sub *BsfSubscription) string {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	subId := uuid.New().String()
	sub.SubId = subId
	c.Subscriptions[subId] = sub

	logger.CtxLog.Debugf("Created subscription with ID: %s", subId)
	return subId
}

func (c *BSFContext) GetSubscription(subId string) (*BsfSubscription, bool) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	sub, exists := c.Subscriptions[subId]
	return sub, exists
}

func (c *BSFContext) UpdateSubscription(subId string, sub *BsfSubscription) bool {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if _, exists := c.Subscriptions[subId]; exists {
		sub.SubId = subId
		c.Subscriptions[subId] = sub
		logger.CtxLog.Debugf("Updated subscription with ID: %s", subId)
		return true
	}
	return false
}

func (c *BSFContext) DeleteSubscription(subId string) bool {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if _, exists := c.Subscriptions[subId]; exists {
		delete(c.Subscriptions, subId)
		logger.CtxLog.Debugf("Deleted subscription with ID: %s", subId)
		return true
	}
	return false
}

// Query functions for PCF bindings based on parameters
func (c *BSFContext) QueryPcfBindings(supi, gpsi, dnn, ipv4Addr, ipv6Prefix, macAddr48, ipDomain string, snssai *models.Snssai) []*PcfBinding {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	var result []*PcfBinding

	for _, binding := range c.PcfBindings {
		match := true

		if supi != "" && (binding.Supi == nil || *binding.Supi != supi) {
			match = false
		}
		if gpsi != "" && (binding.Gpsi == nil || *binding.Gpsi != gpsi) {
			match = false
		}
		if dnn != "" && binding.Dnn != dnn {
			match = false
		}
		if ipv4Addr != "" && (binding.Ipv4Addr == nil || *binding.Ipv4Addr != ipv4Addr) {
			match = false
		}
		if ipv6Prefix != "" && (binding.Ipv6Prefix == nil || *binding.Ipv6Prefix != ipv6Prefix) {
			match = false
		}
		if macAddr48 != "" && (binding.MacAddr48 == nil || *binding.MacAddr48 != macAddr48) {
			match = false
		}
		if ipDomain != "" && (binding.IpDomain == nil || *binding.IpDomain != ipDomain) {
			match = false
		}
		if snssai != nil && !snssaiEquals(binding.Snssai, snssai) {
			match = false
		}

		if match {
			result = append(result, binding)
		}
	}

	return result
}

func (c *BSFContext) QueryPcfForUeBindings(supi, gpsi string) []*PcfForUeBinding {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	var result []*PcfForUeBinding

	for _, binding := range c.PcfForUeBindings {
		match := true

		if supi != "" && binding.Supi != supi {
			match = false
		}
		if gpsi != "" && (binding.Gpsi == nil || *binding.Gpsi != gpsi) {
			match = false
		}

		if match {
			result = append(result, binding)
		}
	}

	return result
}

func (c *BSFContext) QueryPcfMbsBindings(mbsSessionId *models.MbsSessionId) []*PcfMbsBinding {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	var result []*PcfMbsBinding

	for _, binding := range c.PcfMbsBindings {
		if mbsSessionId != nil && binding.MbsSessionId != nil {
			// Compare MBS Session IDs (implementation depends on MbsSessionId structure)
			if mbsSessionIdEquals(binding.MbsSessionId, mbsSessionId) {
				result = append(result, binding)
			}
		}
	}

	return result
}

// Helper functions
func snssaiEquals(snssai1, snssai2 *models.Snssai) bool {
	if snssai1 == nil && snssai2 == nil {
		return true
	}
	if snssai1 == nil || snssai2 == nil {
		return false
	}
	return snssai1.Sst == snssai2.Sst && snssai1.Sd == snssai2.Sd
}

func mbsSessionIdEquals(id1, id2 *models.MbsSessionId) bool {
	if id1 == nil && id2 == nil {
		return true
	}
	if id1 == nil || id2 == nil {
		return false
	}
	// This would need to be implemented based on the actual MbsSessionId structure
	// For now, assuming it has a string representation
	return fmt.Sprintf("%v", id1) == fmt.Sprintf("%v", id2)
}

// Get BSF NF Profile
func (c *BSFContext) GetBsfProfile() models.NrfNfManagementNfProfile {
	nfProfile := models.NrfNfManagementNfProfile{
		NfInstanceId:  c.NfId,
		NfType:        models.NrfNfManagementNfType_BSF,
		NfStatus:      models.NrfNfManagementNfStatus_REGISTERED,
		PlmnList:      []models.PlmnId{{Mcc: "208", Mnc: "93"}}, // Default PLMN
		Ipv4Addresses: []string{c.RegisterIPv4},
		NfServices:    []models.NrfNfManagementNfService{},
	}

	bsfMgmtService := models.NrfNfManagementNfService{
		ServiceInstanceId: uuid.New().String(),
		ServiceName:       models.ServiceName_NBSF_MANAGEMENT,
		Versions: []models.NfServiceVersion{{
			ApiVersionInUri: "v1",
			ApiFullVersion:  "1.5.0",
		}},
		Scheme:          models.UriScheme(c.UriScheme),
		NfServiceStatus: models.NfServiceStatus_REGISTERED,
		IpEndPoints: []models.IpEndPoint{{
			Ipv4Address: c.RegisterIPv4,
			Port:        int32(c.SBIPort),
		}},
		ApiPrefix: fmt.Sprintf("%s://%s:%d", c.UriScheme, c.RegisterIPv4, c.SBIPort),
	}

	nfProfile.NfServices = append(nfProfile.NfServices, bsfMgmtService)

	return nfProfile
}
