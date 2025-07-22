package notifier

import (
	"context"
	"errors"
	"runtime/debug"
	"strconv"
	"sync"

	"github.com/free5gc/nef/internal/logger"
	"github.com/free5gc/openapi/Nnef_PFDmanagement"
	"github.com/free5gc/openapi/models"
)

type PfdChangeNotifier struct {
	clientPfdManagement *Nnef_PFDmanagement.APIClient
	mu                  sync.RWMutex

	numPfdSubID   uint64
	appIdToSubIDs map[string]map[string]bool
	subIdToURI    map[string]string
}

type PfdNotifyContext struct {
	notifier             *PfdChangeNotifier
	appIdToNotification  map[string]models.PfdChangeNotification
	subIdToChangedAppIDs map[string][]string
}

func NewPfdChangeNotifier() (*PfdChangeNotifier, error) {
	return &PfdChangeNotifier{
		appIdToSubIDs: make(map[string]map[string]bool),
		subIdToURI:    make(map[string]string),
	}, nil
}

func (n *PfdChangeNotifier) initPfdManagementApiClient() {
	n.mu.Lock()
	defer n.mu.Unlock()

	if n.clientPfdManagement != nil {
		return
	}

	config := Nnef_PFDmanagement.NewConfiguration()
	n.clientPfdManagement = Nnef_PFDmanagement.NewAPIClient(config)
}

func (n *PfdChangeNotifier) AddPfdSub(pfdSub *models.PfdSubscription) string {
	n.initPfdManagementApiClient()

	n.mu.Lock()
	defer n.mu.Unlock()

	n.numPfdSubID++
	subID := strconv.FormatUint(n.numPfdSubID, 10)
	n.subIdToURI[subID] = pfdSub.NotifyUri
	// TODO: If pfdSub.ApplicationIds is empty, it may means monitoring all appIDs
	for _, appID := range pfdSub.ApplicationIds {
		if _, exist := n.appIdToSubIDs[appID]; !exist {
			n.appIdToSubIDs[appID] = make(map[string]bool)
		}
		n.appIdToSubIDs[appID][subID] = true
	}

	return subID
}

func (n *PfdChangeNotifier) DeletePfdSub(subID string) error {
	n.mu.Lock()
	defer n.mu.Unlock()

	if _, exist := n.subIdToURI[subID]; !exist {
		return errors.New("subscription not found")
	}
	delete(n.subIdToURI, subID)
	for _, subIDs := range n.appIdToSubIDs {
		delete(subIDs, subID)
	}
	return nil
}

func (n *PfdChangeNotifier) getSubIDs(appID string) []string {
	n.mu.RLock()
	defer n.mu.RUnlock()

	subIDs := make([]string, 0, len(n.appIdToSubIDs[appID]))
	for subID := range n.appIdToSubIDs[appID] {
		subIDs = append(subIDs, subID)
	}
	return subIDs
}

func (n *PfdChangeNotifier) getSubURI(subID string) string {
	n.mu.RLock()
	defer n.mu.RUnlock()
	return n.subIdToURI[subID]
}

func (n *PfdChangeNotifier) NewPfdNotifyContext() *PfdNotifyContext {
	return &PfdNotifyContext{
		notifier:             n,
		appIdToNotification:  make(map[string]models.PfdChangeNotification),
		subIdToChangedAppIDs: make(map[string][]string),
	}
}

func (nc *PfdNotifyContext) AddNotification(appID string, notif *models.PfdChangeNotification) {
	nc.appIdToNotification[appID] = *notif
	for _, subID := range nc.notifier.getSubIDs(appID) {
		nc.subIdToChangedAppIDs[subID] = append(nc.subIdToChangedAppIDs[subID], appID)
	}
}

func (nc *PfdNotifyContext) FlushNotifications() {
	for subID, appIDs := range nc.subIdToChangedAppIDs {
		pfdChangeNotifications := make([]models.PfdChangeNotification, 0, len(appIDs))
		for _, appID := range appIDs {
			pfdChangeNotifications = append(pfdChangeNotifications, nc.appIdToNotification[appID])
		}

		go func(id string) {
			defer func() {
				if p := recover(); p != nil {
					// Print stack for panic to log. Fatalf() will let program exit.
					logger.PFDManageLog.Fatalf("panic: %v\n%s", p, string(debug.Stack()))
				}
			}()

			_, _, err := nc.notifier.clientPfdManagement.NotificationApi.NotificationPost(
				context.TODO(), nc.notifier.getSubURI(id), pfdChangeNotifications)
			if err != nil {
				logger.PFDManageLog.Fatal(err)
			}
		}(subID)
		// TODO: Handle the response of notification properly
	}
}
