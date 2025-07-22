package processor

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"github.com/free5gc/openapi/models"
	"github.com/stretchr/testify/require"
	"gopkg.in/h2non/gock.v1"
)

func TestGetApplicationsPFD(t *testing.T) {
	initUDRDrGetPfdDatasStub()
	defer gock.Off()

	testCases := []struct {
		description      string
		appIDs           []string
		expectedResponse *HandlerResponse
	}{
		{
			description: "TC1: All App IDs found, should return all PfdDataforApp",
			appIDs:      []string{"app1", "app2"},
			expectedResponse: &HandlerResponse{
				Status: http.StatusOK,
				Body:   &[]models.PfdDataForApp{pfdDataForApp1, pfdDataForApp2},
			},
		},
		{
			description: "TC2: All App ID not found, should return ProblemDetails",
			appIDs:      []string{"app3"},
			expectedResponse: &HandlerResponse{
				Status: http.StatusNotFound,
				Body:   &models.ProblemDetails{Status: http.StatusNotFound},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			rsp := nefApp.Processor().GetApplicationsPFD(tc.appIDs)
			require.Equal(t, tc.expectedResponse, rsp)
		})
	}
}

func TestGetIndividualApplicationPFD(t *testing.T) {
	initUDRDrGetPfdDataStub()
	defer gock.Off()

	testCases := []struct {
		description      string
		appID            string
		expectedResponse *HandlerResponse
	}{
		{
			description: "TC1: App ID found, should return the PfdDataforApp",
			appID:       "app1",
			expectedResponse: &HandlerResponse{
				Status: http.StatusOK,
				Body:   &pfdDataForApp1,
			},
		},
		{
			description: "TC2: App ID not found, should return ProblemDetails",
			appID:       "app3",
			expectedResponse: &HandlerResponse{
				Status: http.StatusNotFound,
				Body:   &models.ProblemDetails{Status: http.StatusNotFound},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			rsp := nefApp.Processor().GetIndividualApplicationPFD(tc.appID)
			require.Equal(t, tc.expectedResponse, rsp)
		})
	}
}

func TestPostPFDSubscriptions(t *testing.T) {
	pfdSubsc := &models.PfdSubscription{
		ApplicationIds: []string{"app1", "app2"},
		NotifyUri:      "http://pfdSub1URI/notify",
	}

	testCases := []struct {
		description      string
		subscription     *models.PfdSubscription
		expectedResponse *HandlerResponse
	}{
		{
			description:  "TC1: Successful subscription, should return PfdSubscription",
			subscription: pfdSubsc,
			expectedResponse: &HandlerResponse{
				Status: http.StatusCreated,
				Headers: map[string][]string{
					"Location": {nefApp.Processor().genPfdSubscriptionURI("1")},
				},
				Body: pfdSubsc,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			rsp := nefApp.Processor().PostPFDSubscriptions(tc.subscription)
			require.Equal(t, tc.expectedResponse, rsp)
		})
	}
}

func TestDeleteIndividualPFDSubscription(t *testing.T) {
	testCases := []struct {
		description      string
		subscriptionID   string
		expectedResponse *HandlerResponse
	}{
		{
			description:    "TC1: Successful unsubscription",
			subscriptionID: "1",
			expectedResponse: &HandlerResponse{
				Status: http.StatusNoContent,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			rsp := nefApp.Processor().DeleteIndividualPFDSubscription(tc.subscriptionID)
			require.Equal(t, tc.expectedResponse, rsp)
		})
	}
}

var (
	// `notifChan` are used in `TestPostPfdChangeReports()` to pass the notification requests intercepted by gock.
	notifChan   = make(chan *http.Request)
	pfdContent1 = models.PfdContent{
		PfdId: "pfd1",
		FlowDescriptions: []string{
			"permit in ip from 10.68.28.39 80 to any",
			"permit out ip from any to 10.68.28.39 80",
		},
	}
)

func TestPostPfdChangeReports(t *testing.T) {
	// Note: Because TestPostPFDSubscriptions() already used subscription ID 1, the ID will start from 2 here.
	initUDRDrPutPfdDataStub(http.StatusOK)
	initUDRDrDeletePfdDataStub()
	initNEFNotificationStub("http://pfdSub2URI")
	initNEFNotificationStub("http://pfdSub3URI")
	defer gock.Off()
	gock.Observe(func(request *http.Request, mock gock.Mock) {
		if strings.Contains(request.URL.String(), "pfdSub") {
			notifChan <- request
		}
	})

	af := nefApp.Context().NewAf("af1")
	nefApp.Context().AddAf(af)
	defer nefApp.Context().DeleteAf("af1")

	af.Mu.Lock()
	afPfdTr := af.NewPfdTrans()
	af.PfdTrans[afPfdTr.TransID] = afPfdTr
	afPfdTr.AddExtAppID("app1")
	afPfdTr.AddExtAppID("app2")
	af.Mu.Unlock()

	subsID1 := nefApp.Notifier().PfdChangeNotifier.AddPfdSub(&models.PfdSubscription{
		ApplicationIds: []string{"app1"},
		NotifyUri:      "http://pfdSub2URI",
	})
	subsID2 := nefApp.Notifier().PfdChangeNotifier.AddPfdSub(&models.PfdSubscription{
		ApplicationIds: []string{"app1", "app2"},
		NotifyUri:      "http://pfdSub3URI",
	})
	defer func() {
		if err := nefApp.Notifier().PfdChangeNotifier.DeletePfdSub(subsID1); err != nil {
			t.Fatal(err)
		}
		if err := nefApp.Notifier().PfdChangeNotifier.DeletePfdSub(subsID2); err != nil {
			t.Fatal(err)
		}
	}()

	testCases := []struct {
		description           string
		triggerFunc           func()
		expectedNotifications map[string][]models.PfdChangeNotification
	}{
		{
			description: "Update app1, should send notification for subscription 2 and 3",
			triggerFunc: func() {
				nefApp.Processor().PutIndividualApplicationPFDManagement("af1", "1", "app1", &models.PfdData{
					ExternalAppId: "app1",
					Pfds: map[string]models.Pfd{
						"pfd1": pfd1,
					},
				})
			},
			expectedNotifications: map[string][]models.PfdChangeNotification{
				"http://pfdSub2URI/notify": {
					{
						ApplicationId: "app1",
						Pfds: []models.PfdContent{
							pfdContent1,
						},
					},
				},
				"http://pfdSub3URI/notify": {
					{
						ApplicationId: "app1",
						Pfds: []models.PfdContent{
							pfdContent1,
						},
					},
				},
			},
		},
		{
			description: "Delete app2, should send notification for subscription 3",
			triggerFunc: func() {
				nefApp.Processor().DeleteIndividualApplicationPFDManagement("af1", "1", "app2")
			},
			expectedNotifications: map[string][]models.PfdChangeNotification{
				"http://pfdSub3URI/notify": {
					{
						ApplicationId: "app2",
						RemovalFlag:   true,
					},
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			tc.triggerFunc()
			for i := 0; i < len(tc.expectedNotifications); i++ {
				r := <-notifChan

				var getNotifications []models.PfdChangeNotification
				if err := json.NewDecoder(r.Body).Decode(&getNotifications); err != nil {
					t.Fatal(err)
				}
				require.Equal(t, tc.expectedNotifications[r.URL.String()], getNotifications)
			}
		})
	}
}

func initNEFNotificationStub(notifyURI string) {
	gock.New(notifyURI).
		Post("/notify").
		Persist().
		Reply(http.StatusNoContent)
}
