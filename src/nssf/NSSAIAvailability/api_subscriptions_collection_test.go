/*
 * NSSF NSSAI Availability
 *
 * NSSF NSSAI Availability Service
 */

package NSSAIAvailability

import (
	"context"
	"encoding/json"
	"net/http"
	"reflect"
	"testing"
	"time"

	"free5gc/lib/Nnssf_NSSAIAvailability"
	"free5gc/lib/http2_util"
	. "free5gc/lib/openapi/models"
	"free5gc/src/nssf/factory"
	"free5gc/src/nssf/nssf_handler"
	"free5gc/src/nssf/test"
	"free5gc/src/nssf/util"
)

var testingNssaiavailabilitySubscribeApi = test.TestingNssaiavailability{
	ConfigFile: test.ConfigFileFromArgs,
}

func TestNSSAIAvailabilityPost(t *testing.T) {
	factory.InitConfigFactory(testingNssaiavailabilitySubscribeApi.ConfigFile)

	router := NewRouter()
	srv, err := http2_util.NewServer(":29531", util.NSSF_LOG_PATH, router)
	if err != nil {
		t.Fatal(err)
	}

	go nssf_handler.Handle()

	go func() {
		err := srv.ListenAndServeTLS(util.NSSF_PEM_PATH, util.NSSF_KEY_PATH)
		if err != nil && err != http.ErrServerClosed {
			t.Fatal(err)
		}
	}()

	configuration := Nnssf_NSSAIAvailability.NewConfiguration()
	configuration.SetBasePath("https://localhost:29531")
	apiClient := Nnssf_NSSAIAvailability.NewAPIClient(configuration)

	subtests := []struct {
		name                                   string
		nssfEventSubscriptionCreateData        *NssfEventSubscriptionCreateData
		expectStatus                           int
		expectNssfEventSubscriptionCreatedData *NssfEventSubscriptionCreatedData
	}{
		{
			name: "Post",
			nssfEventSubscriptionCreateData: &NssfEventSubscriptionCreateData{
				NfNssaiAvailabilityUri: "http://free5gc-amf2.nctu.me:29518/namf-nssaiavailability/v1/nssai-availability/notify",
				TaiList: []Tai{
					{
						PlmnId: &PlmnId{
							Mcc: "466",
							Mnc: "92",
						},
						Tac: "33456",
					},
					{
						PlmnId: &PlmnId{
							Mcc: "466",
							Mnc: "92",
						},
						Tac: "33458",
					},
				},
				Event:  func() NssfEventType { n := NssfEventType_SNSSAI_STATUS_CHANGE_REPORT; return n }(),
				Expiry: func() *time.Time { t, _ := time.Parse(time.RFC3339, "2019-06-24T16:35:31+08:00"); return &t }(),
			},
			expectStatus: http.StatusCreated,
			expectNssfEventSubscriptionCreatedData: &NssfEventSubscriptionCreatedData{
				SubscriptionId: "2",
				Expiry:         func() *time.Time { t, _ := time.Parse(time.RFC3339, "2019-06-24T16:35:31+08:00"); return &t }(),
				AuthorizedNssaiAvailabilityData: []AuthorizedNssaiAvailabilityData{
					{
						Tai: &Tai{
							PlmnId: &PlmnId{
								Mcc: "466",
								Mnc: "92",
							},
							Tac: "33456",
						},
						SupportedSnssaiList: []Snssai{
							{
								Sst: 1,
							},
							{
								Sst: 1,
								Sd:  "1",
							},
							{
								Sst: 1,
								Sd:  "2",
							},
							{
								Sst: 2,
							},
						},
					},
					{
						Tai: &Tai{
							PlmnId: &PlmnId{
								Mcc: "466",
								Mnc: "92",
							},
							Tac: "33458",
						},
						SupportedSnssaiList: []Snssai{
							{
								Sst: 1,
							},
							{
								Sst: 1,
								Sd:  "1",
							},
							{
								Sst: 1,
								Sd:  "3",
							},
							{
								Sst: 2,
							},
						},
						RestrictedSnssaiList: []RestrictedSnssai{
							{
								HomePlmnId: &PlmnId{
									Mcc: "310",
									Mnc: "560",
								},
								SNssaiList: []Snssai{
									{
										Sst: 1,
										Sd:  "3",
									},
								},
							},
						},
					},
				},
			},
		},
	}

	for _, subtest := range subtests {
		t.Run(subtest.name, func(t *testing.T) {
			var (
				resp *http.Response
				n    NssfEventSubscriptionCreatedData
			)

			n, resp, err := apiClient.SubscriptionsCollectionApi.NSSAIAvailabilityPost(context.Background(), *subtest.nssfEventSubscriptionCreateData)

			if err != nil {
				t.Errorf(err.Error())
			}

			if resp.StatusCode != subtest.expectStatus {
				t.Errorf("Incorrect status code: expected %d, got %d", subtest.expectStatus, resp.StatusCode)
			}

			if reflect.DeepEqual(n, *subtest.expectNssfEventSubscriptionCreatedData) == false {
				e, _ := json.Marshal(*subtest.expectNssfEventSubscriptionCreatedData)
				r, _ := json.Marshal(n)
				t.Errorf("Incorrect NSSF event subscription created data:\nexpected\n%s\n, got\n%s", string(e), string(r))
			}
		})
	}

	err = srv.Shutdown(context.Background())
	if err != nil {
		t.Fatal(err)
	}
}
