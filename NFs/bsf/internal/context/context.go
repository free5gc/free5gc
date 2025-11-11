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
	"go.mongodb.org/mongo-driver/bson"
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

// MongoDB collection names
const (
	PCF_BINDINGS_COLLECTION = "pcfBindings"
	// TODO: Implement persistence for UE-specific PCF binding management
	// PCF_FOR_UE_BINDINGS_COLLECTION = "pcfForUeBindings"

	// TODO: Implement persistence for MBS (Multicast/Broadcast Service) binding management
	// PCF_MBS_BINDINGS_COLLECTION    = "pcfMbsBindings"

	// TODO: Implement persistence  for subscription management for BSF events
	// SUBSCRIPTIONS_COLLECTION       = "subscriptions"
)

func init() {
	BsfSelf.Name = "bsf"
	BsfSelf.NfId = uuid.New().String()
	BsfSelf.PcfBindings = make(map[string]*PcfBinding)
	BsfSelf.PcfForUeBindings = make(map[string]*PcfForUeBinding)
	BsfSelf.PcfMbsBindings = make(map[string]*PcfMbsBinding)
	BsfSelf.Subscriptions = make(map[string]*BsfSubscription)
	BsfSelf.mutex = sync.RWMutex{}

	// Initialize lifecycle management defaults
	BsfSelf.DefaultBindingTTL = 24 * time.Hour // 24 hours default TTL
	BsfSelf.CleanupInterval = 10 * time.Minute // Cleanup every 10 minutes
	BsfSelf.MaxInactiveTime = 1 * time.Hour    // Delete if inactive for 1 hour
	BsfSelf.ShutdownChannel = make(chan bool, 1)
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
	MongoDBName   string
	MongoDBUrl    string
	MongoClient   *mongo.Client
	MongoDatabase *mongo.Database

	// BSF Business Logic
	PcfBindings      map[string]*PcfBinding      // bindingId -> PcfBinding
	PcfForUeBindings map[string]*PcfForUeBinding // bindingId -> PcfForUeBinding
	PcfMbsBindings   map[string]*PcfMbsBinding   // bindingId -> PcfMbsBinding
	Subscriptions    map[string]*BsfSubscription // subId -> BsfSubscription

	// Lifecycle management
	DefaultBindingTTL time.Duration // Default TTL for new bindings
	CleanupInterval   time.Duration // How often to run cleanup
	MaxInactiveTime   time.Duration // Max time without access before cleanup
	CleanupTicker     *time.Ticker  // Ticker for periodic cleanup
	ShutdownChannel   chan bool     // Channel for graceful shutdown
}

type PcfBinding struct {
	BindingId          string                       `bson:"_id,omitempty"`
	Supi               *string                      `bson:"supi,omitempty"`
	Gpsi               *string                      `bson:"gpsi,omitempty"`
	Ipv4Addr           *string                      `bson:"ipv4_addr,omitempty"`
	Ipv6Prefix         *string                      `bson:"ipv6_prefix,omitempty"`
	AddIpv6Prefixes    []string                     `bson:"add_ipv6_prefixes,omitempty"`
	IpDomain           *string                      `bson:"ip_domain,omitempty"`
	MacAddr48          *string                      `bson:"mac_addr48,omitempty"`
	AddMacAddrs        []string                     `bson:"add_mac_addrs,omitempty"`
	Dnn                string                       `bson:"dnn"`
	PcfFqdn            *string                      `bson:"pcf_fqdn,omitempty"`
	PcfIpEndPoints     []models.IpEndPoint          `bson:"pcf_ip_endpoints,omitempty"`
	PcfDiamHost        *string                      `bson:"pcf_diam_host,omitempty"`
	PcfDiamRealm       *string                      `bson:"pcf_diam_realm,omitempty"`
	PcfSmFqdn          *string                      `bson:"pcf_sm_fqdn,omitempty"`
	PcfSmIpEndPoints   []models.IpEndPoint          `bson:"pcf_sm_ip_endpoints,omitempty"`
	Snssai             *models.Snssai               `bson:"snssai,omitempty"`
	SuppFeat           *string                      `bson:"supp_feat,omitempty"`
	PcfId              *string                      `bson:"pcf_id,omitempty"`
	PcfSetId           *string                      `bson:"pcf_set_id,omitempty"`
	RecoveryTime       *time.Time                   `bson:"recovery_time,omitempty"`
	ParaCom            *models.ParameterCombination `bson:"para_com,omitempty"`
	BindLevel          *models.BindingLevel         `bson:"bind_level,omitempty"`
	Ipv4FrameRouteList []string                     `bson:"ipv4_frame_route_list,omitempty"`
	Ipv6FrameRouteList []string                     `bson:"ipv6_frame_route_list,omitempty"`

	// Lifecycle management fields
	CreatedTime    time.Time  `bson:"created_time"`
	LastAccessTime time.Time  `bson:"last_access_time"`
	ExpiryTime     *time.Time `bson:"expiry_time,omitempty"`
	IsActive       bool       `bson:"is_active"`
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
func (c *BSFContext) ConnectMongoDB(ctx context.Context) error {
	if c.MongoDBUrl == "" {
		return fmt.Errorf("MongoDB URL is empty")
	}

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(c.MongoDBUrl))
	if err != nil {
		return fmt.Errorf("failed to connect to MongoDB: %+v", err)
	}

	// Test the connection
	err = client.Ping(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to ping MongoDB: %+v", err)
	}

	// Store client and database references
	c.MongoClient = client
	c.MongoDatabase = client.Database(c.MongoDBName)

	logger.CtxLog.Infof("Connected to MongoDB at %s, database: %s", c.MongoDBUrl, c.MongoDBName)
	return nil
}

// DisconnectMongoDB gracefully disconnects from MongoDB
func (c *BSFContext) DisconnectMongoDB() error {
	if c.MongoClient != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := c.MongoClient.Disconnect(ctx); err != nil {
			return fmt.Errorf("failed to disconnect from MongoDB: %+v", err)
		}

		c.MongoClient = nil
		c.MongoDatabase = nil
		logger.CtxLog.Info("Disconnected from MongoDB")
	}
	return nil
}

// PCF Binding Management
func (c *BSFContext) CreatePcfBinding(binding *PcfBinding) string {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	bindingId := uuid.New().String()
	binding.BindingId = bindingId

	// Set lifecycle fields
	now := time.Now()
	binding.CreatedTime = now
	binding.LastAccessTime = now
	if binding.ExpiryTime == nil {
		expiryTime := now.Add(c.DefaultBindingTTL)
		binding.ExpiryTime = &expiryTime
	}
	binding.IsActive = true

	// Store in memory for fast access
	c.PcfBindings[bindingId] = binding

	// Persist to MongoDB
	if c.MongoDatabase != nil {
		collection := c.MongoDatabase.Collection(PCF_BINDINGS_COLLECTION)
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		logger.CtxLog.Debugf("Attempting to insert PCF binding to MongoDB: %s", bindingId)
		result, err := collection.InsertOne(ctx, binding)
		if err != nil {
			logger.CtxLog.Errorf("Failed to insert PCF binding to MongoDB: %+v", err)
		} else {
			logger.CtxLog.Infof("PCF binding persisted to MongoDB with ID: %s, InsertedID: %v", bindingId, result.InsertedID)
		}
	} else {
		logger.CtxLog.Warnf("MongoDB database is nil, cannot persist PCF binding: %s", bindingId)
	}

	logger.CtxLog.Debugf("Created PCF binding with ID: %s", bindingId)
	return bindingId
}

func (c *BSFContext) GetPcfBinding(bindingId string) (*PcfBinding, bool) {
	c.mutex.RLock()

	// First check in-memory cache
	binding, exists := c.PcfBindings[bindingId]
	if exists {
		// Update last access time
		binding.LastAccessTime = time.Now()
		c.mutex.RUnlock()

		// Update in MongoDB asynchronously
		go c.updateLastAccessTimeInMongoDB(bindingId, binding.LastAccessTime)

		return binding, true
	}
	c.mutex.RUnlock()

	// If not in memory, try MongoDB
	if c.MongoDatabase != nil {
		collection := c.MongoDatabase.Collection(PCF_BINDINGS_COLLECTION)
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		var dbBinding PcfBinding
		err := collection.FindOne(ctx, bson.M{"_id": bindingId}).Decode(&dbBinding)
		if err == nil {
			// Update last access time
			dbBinding.LastAccessTime = time.Now()

			// Load into memory cache for future access
			c.mutex.Lock()
			c.PcfBindings[bindingId] = &dbBinding
			c.mutex.Unlock()

			// Update in MongoDB asynchronously
			go c.updateLastAccessTimeInMongoDB(bindingId, dbBinding.LastAccessTime)

			logger.CtxLog.Debugf("PCF binding loaded from MongoDB: %s", bindingId)
			return &dbBinding, true
		} else if err != mongo.ErrNoDocuments {
			logger.CtxLog.Errorf("Failed to query PCF binding from MongoDB: %+v", err)
		}
	}

	return nil, false
}

func (c *BSFContext) UpdatePcfBinding(bindingId string, binding *PcfBinding) bool {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	// Check if binding exists (either in memory or MongoDB)
	_, exists := c.GetPcfBindingUnsafe(bindingId)
	if !exists {
		return false
	}

	binding.BindingId = bindingId
	c.PcfBindings[bindingId] = binding

	// Update in MongoDB
	if c.MongoDatabase != nil {
		collection := c.MongoDatabase.Collection(PCF_BINDINGS_COLLECTION)
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		filter := bson.M{"_id": bindingId}
		_, err := collection.ReplaceOne(ctx, filter, binding)
		if err != nil {
			logger.CtxLog.Errorf("Failed to update PCF binding in MongoDB: %+v", err)
		} else {
			logger.CtxLog.Debugf("PCF binding updated in MongoDB: %s", bindingId)
		}
	}

	logger.CtxLog.Debugf("Updated PCF binding with ID: %s", bindingId)
	return true
}

// Helper function for internal use without additional locking
func (c *BSFContext) GetPcfBindingUnsafe(bindingId string) (*PcfBinding, bool) {
	// Check memory first
	if binding, exists := c.PcfBindings[bindingId]; exists {
		return binding, true
	}

	// Check MongoDB
	if c.MongoDatabase != nil {
		collection := c.MongoDatabase.Collection(PCF_BINDINGS_COLLECTION)
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		var dbBinding PcfBinding
		err := collection.FindOne(ctx, bson.M{"_id": bindingId}).Decode(&dbBinding)
		if err == nil {
			return &dbBinding, true
		}
	}

	return nil, false
}

func (c *BSFContext) DeletePcfBinding(bindingId string) bool {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	// Check if binding exists (either in memory or MongoDB)
	_, exists := c.GetPcfBindingUnsafe(bindingId)
	if !exists {
		return false
	}

	// Delete from memory
	delete(c.PcfBindings, bindingId)

	// Delete from MongoDB
	if c.MongoDatabase != nil {
		collection := c.MongoDatabase.Collection(PCF_BINDINGS_COLLECTION)
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		filter := bson.M{"_id": bindingId}
		logger.CtxLog.Debugf("Attempting to delete PCF binding from MongoDB: %s", bindingId)
		result, err := collection.DeleteOne(ctx, filter)
		if err != nil {
			logger.CtxLog.Errorf("Failed to delete PCF binding from MongoDB: %+v", err)
		} else if result.DeletedCount > 0 {
			logger.CtxLog.Infof("PCF binding deleted from MongoDB: %s, DeletedCount: %d", bindingId, result.DeletedCount)
		} else {
			logger.CtxLog.Warnf("PCF binding not found in MongoDB for deletion: %s", bindingId)
		}
	} else {
		logger.CtxLog.Warnf("MongoDB database is nil, cannot delete PCF binding: %s", bindingId)
	}

	logger.CtxLog.Debugf("Deleted PCF binding with ID: %s", bindingId)
	return true
}

// LoadPcfBindingsFromMongoDB loads all PCF bindings from MongoDB into memory cache
func (c *BSFContext) LoadPcfBindingsFromMongoDB() error {
	if c.MongoDatabase == nil {
		return fmt.Errorf("MongoDB database not initialized")
	}

	collection := c.MongoDatabase.Collection(PCF_BINDINGS_COLLECTION)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return fmt.Errorf("failed to query PCF bindings from MongoDB: %+v", err)
	}
	defer func() {
		if closeErr := cursor.Close(ctx); closeErr != nil {
			logger.CtxLog.Errorf("Failed to close cursor: %+v", closeErr)
		}
	}()

	c.mutex.Lock()
	defer c.mutex.Unlock()

	count := 0
	for cursor.Next(ctx) {
		var binding PcfBinding
		if decodeErr := cursor.Decode(&binding); decodeErr != nil {
			logger.CtxLog.Errorf("Failed to decode PCF binding: %+v", decodeErr)
			continue
		}
		c.PcfBindings[binding.BindingId] = &binding
		count++
	}

	if cursorErr := cursor.Err(); cursorErr != nil {
		return fmt.Errorf("cursor error while loading PCF bindings: %+v", cursorErr)
	}

	logger.CtxLog.Infof("Loaded %d PCF bindings from MongoDB", count)
	return nil
}

// updateLastAccessTimeInMongoDB updates the last access time in MongoDB
func (c *BSFContext) updateLastAccessTimeInMongoDB(bindingId string, lastAccessTime time.Time) {
	if c.MongoDatabase == nil {
		return
	}

	collection := c.MongoDatabase.Collection(PCF_BINDINGS_COLLECTION)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"_id": bindingId}
	update := bson.M{"$set": bson.M{"last_access_time": lastAccessTime}}

	_, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		logger.CtxLog.Errorf("Failed to update last access time in MongoDB for binding %s: %+v", bindingId, err)
	}
}

// StartCleanupRoutine starts the periodic cleanup of expired and inactive bindings
func (c *BSFContext) StartCleanupRoutine() {
	logger.CtxLog.Infof("Starting PCF binding cleanup routine (interval: %v)", c.CleanupInterval)

	c.CleanupTicker = time.NewTicker(c.CleanupInterval)

	go func() {
		for {
			select {
			case <-c.CleanupTicker.C:
				c.CleanupExpiredBindings()
			case <-c.ShutdownChannel:
				logger.CtxLog.Info("Stopping PCF binding cleanup routine")
				c.CleanupTicker.Stop()
				return
			}
		}
	}()
}

// StopCleanupRoutine stops the cleanup routine
func (c *BSFContext) StopCleanupRoutine() {
	if c.CleanupTicker != nil {
		select {
		case c.ShutdownChannel <- true:
		default:
		}
	}
}

// CleanupExpiredBindings removes expired and inactive bindings
func (c *BSFContext) CleanupExpiredBindings() {
	logger.CtxLog.Debug("Running PCF binding cleanup")

	now := time.Now()
	expiredBindings := []string{}
	inactiveBindings := []string{}

	c.mutex.RLock()
	for bindingId, binding := range c.PcfBindings {
		// Check expiry time
		if binding.ExpiryTime != nil && now.After(*binding.ExpiryTime) {
			expiredBindings = append(expiredBindings, bindingId)
			continue
		}

		// Check last access time for inactive bindings
		if now.Sub(binding.LastAccessTime) > c.MaxInactiveTime {
			inactiveBindings = append(inactiveBindings, bindingId)
		}
	}
	c.mutex.RUnlock()

	// Delete expired bindings
	for _, bindingId := range expiredBindings {
		logger.CtxLog.Infof("Deleting expired PCF binding: %s", bindingId)
		c.DeletePcfBinding(bindingId)
	}

	// Delete inactive bindings
	for _, bindingId := range inactiveBindings {
		logger.CtxLog.Infof("Deleting inactive PCF binding: %s (last access: %v ago)",
			bindingId, now.Sub(c.PcfBindings[bindingId].LastAccessTime))
		c.DeletePcfBinding(bindingId)
	}

	if len(expiredBindings) > 0 || len(inactiveBindings) > 0 {
		logger.CtxLog.Infof("Cleanup completed: %d expired, %d inactive bindings removed",
			len(expiredBindings), len(inactiveBindings))
	}
}

// CleanupBySupi removes all bindings for a specific SUPI (when UE deregisters)
func (c *BSFContext) CleanupBySupi(supi string) {
	logger.CtxLog.Infof("Cleaning up PCF bindings for SUPI: %s", supi)

	bindingsToDelete := []string{}

	c.mutex.RLock()
	for bindingId, binding := range c.PcfBindings {
		if binding.Supi != nil && *binding.Supi == supi {
			bindingsToDelete = append(bindingsToDelete, bindingId)
		}
	}
	c.mutex.RUnlock()

	for _, bindingId := range bindingsToDelete {
		logger.CtxLog.Infof("Deleting PCF binding for SUPI %s: %s", supi, bindingId)
		c.DeletePcfBinding(bindingId)
	}

	logger.CtxLog.Infof("Cleaned up %d PCF bindings for SUPI: %s", len(bindingsToDelete), supi)
}

// CleanupByPcfId removes all bindings for a specific PCF (when PCF becomes unavailable)
func (c *BSFContext) CleanupByPcfId(pcfId string) {
	logger.CtxLog.Infof("Cleaning up PCF bindings for PCF ID: %s", pcfId)

	bindingsToDelete := []string{}

	c.mutex.RLock()
	for bindingId, binding := range c.PcfBindings {
		if binding.PcfId != nil && *binding.PcfId == pcfId {
			bindingsToDelete = append(bindingsToDelete, bindingId)
		}
	}
	c.mutex.RUnlock()

	for _, bindingId := range bindingsToDelete {
		logger.CtxLog.Infof("Deleting PCF binding for PCF ID %s: %s", pcfId, bindingId)
		c.DeletePcfBinding(bindingId)
	}

	logger.CtxLog.Infof("Cleaned up %d PCF bindings for PCF ID: %s", len(bindingsToDelete), pcfId)
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
func (c *BSFContext) QueryPcfBindings(
	supi, gpsi, dnn, ipv4Addr, ipv6Prefix, macAddr48, ipDomain string,
	snssai *models.Snssai,
) []*PcfBinding {
	c.mutex.RLock()
	var result []*PcfBinding

	// First search in-memory cache
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
	c.mutex.RUnlock()

	// If no results found in memory and MongoDB is available, search MongoDB
	if len(result) == 0 && c.MongoDatabase != nil {
		logger.CtxLog.Debugf("No PCF bindings found in memory cache, searching MongoDB for SUPI: %s, DNN: %s", supi, dnn)

		collection := c.MongoDatabase.Collection(PCF_BINDINGS_COLLECTION)
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// Build MongoDB query filter
		filter := bson.M{}
		if supi != "" {
			filter["supi"] = supi
		}
		if gpsi != "" {
			filter["gpsi"] = gpsi
		}
		if dnn != "" {
			filter["dnn"] = dnn
		}
		if ipv4Addr != "" {
			filter["ipv4_addr"] = ipv4Addr
		}
		if ipv6Prefix != "" {
			filter["ipv6_prefix"] = ipv6Prefix
		}
		if macAddr48 != "" {
			filter["mac_addr48"] = macAddr48
		}
		if ipDomain != "" {
			filter["ip_domain"] = ipDomain
		}
		if snssai != nil {
			filter["snssai.sst"] = snssai.Sst
			if snssai.Sd != "" {
				filter["snssai.sd"] = snssai.Sd
			}
		}

		cursor, err := collection.Find(ctx, filter)
		if err != nil {
			logger.CtxLog.Errorf("Failed to query PCF bindings from MongoDB: %+v", err)
			return result
		}
		defer func() {
			if closeErr := cursor.Close(ctx); closeErr != nil {
				logger.CtxLog.Errorf("Failed to close cursor: %+v", closeErr)
			}
		}()

		for cursor.Next(ctx) {
			var dbBinding PcfBinding
			if decodeErr := cursor.Decode(&dbBinding); decodeErr != nil {
				logger.CtxLog.Errorf("Failed to decode PCF binding from MongoDB: %+v", decodeErr)
				continue
			}

			// Load matching bindings into memory cache for future access
			c.mutex.Lock()
			c.PcfBindings[dbBinding.BindingId] = &dbBinding
			c.mutex.Unlock()

			result = append(result, &dbBinding)
			logger.CtxLog.Infof("Loaded PCF binding from MongoDB into cache: %s", dbBinding.BindingId)
		}

		if cursorErr := cursor.Err(); cursorErr != nil {
			logger.CtxLog.Errorf("Cursor error while querying PCF bindings: %+v", cursorErr)
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
