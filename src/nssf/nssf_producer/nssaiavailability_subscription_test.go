/*
 * NSSF NSSAI Availability
 *
 * NSSF NSSAI Availability Service
 */

package nssf_producer

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strings"
	"testing"
	"time"

	"gopkg.in/yaml.v2"

	. "free5gc/lib/openapi/models"
	"free5gc/src/nssf/factory"
	"free5gc/src/nssf/test"
)

var testingSubscription = test.TestingNssaiavailability{
	ConfigFile:     test.ConfigFileFromArgs,
	SubscriptionId: "3",
}

func checkSubscriptionExist(subscriptionId string) bool {
	for _, subscription := range factory.NssfConfig.Subscriptions {
		if subscription.SubscriptionId == subscriptionId {
			return true
		}
	}
	return false
}

func generateSubscriptionRequest() NssfEventSubscriptionCreateData {
	const jsonRequest = `
        {
            "nfNssaiAvailabilityUri": "http://free5gc-amf2.nctu.me:29518/namf-nssaiavailability/v1/nssai-availability/notify",
            "taiList": [
                {
                    "plmnId": {
                        "mcc": "466",
                        "mnc": "92"
                    },
                    "tac": "33456"
                },
                {
                    "plmnId": {
                        "mcc": "466",
                        "mnc": "92"
                    },
                    "tac": "33458"
                }
            ],
            "event": "SNSSAI_STATUS_CHANGE_REPORT",
            "expiry": "2019-06-24T16:35:31+08:00"
        }
    `

	var n NssfEventSubscriptionCreateData
	if err := json.NewDecoder(strings.NewReader(jsonRequest)).Decode(&n); err != nil {
		fmt.Printf("Decode error: %v", err)
	}

	return n
}

func TestSubscriptionTemplate(t *testing.T) {
	t.Skip()

	// Tests may have different configuration files
	factory.InitConfigFactory(testingSubscription.ConfigFile)

	d, _ := yaml.Marshal(*factory.NssfConfig.Info)
	t.Logf("%s", string(d))
}

func TestSubscriptionPost(t *testing.T) {
	factory.InitConfigFactory(testingSubscription.ConfigFile)

	subtests := []struct {
		name                          string
		generateRequestBody           func() NssfEventSubscriptionCreateData
		expectStatus                  int
		expectSubscriptionCreatedData *NssfEventSubscriptionCreatedData
		expectProblemDetails          *ProblemDetails
	}{
		{
			name:                "Subscribe",
			generateRequestBody: generateSubscriptionRequest,
			expectStatus:        http.StatusCreated,
			expectSubscriptionCreatedData: &NssfEventSubscriptionCreatedData{
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
				n      NssfEventSubscriptionCreateData
				status int
				c      NssfEventSubscriptionCreatedData
				d      ProblemDetails
			)

			if subtest.generateRequestBody != nil {
				n = subtest.generateRequestBody()
			}

			status = subscriptionPost(n, &c, &d)

			if status == http.StatusCreated {
				if reflect.DeepEqual(c, *subtest.expectSubscriptionCreatedData) == false {
					e, _ := json.Marshal(*subtest.expectSubscriptionCreatedData)
					r, _ := json.Marshal(c)
					t.Errorf("Incorrect NSSF event subscription created data:\nexpected\n%s\n, got\n%s", string(e), string(r))
				}
			} else {
				if reflect.DeepEqual(d, *subtest.expectProblemDetails) == false {
					e, _ := json.Marshal(*subtest.expectProblemDetails)
					r, _ := json.Marshal(d)
					t.Errorf("Incorrect problem details:\nexpected\n%s\n, got\n%s", string(e), string(r))
				}
			}
		})
	}
}

func TestSubscriptionDelete(t *testing.T) {
	factory.InitConfigFactory(testingSubscription.ConfigFile)

	subtests := []struct {
		name                 string
		expectStatus         int
		expectProblemDetails *ProblemDetails
	}{
		{
			name:         "Unsubscribe",
			expectStatus: http.StatusNoContent,
		},
	}

	for _, subtest := range subtests {
		t.Run(subtest.name, func(t *testing.T) {
			var (
				status int
				d      ProblemDetails
			)

			status = subscriptionDelete(testingSubscription.SubscriptionId, &d)

			if status == http.StatusNoContent {
				if checkSubscriptionExist(testingSubscription.SubscriptionId) == true {
					t.Errorf("Subscription ID '%s' in configuration should be deleted, but still exists", testingSubscription.SubscriptionId)
				}
			} else {
				if reflect.DeepEqual(d, *subtest.expectProblemDetails) == false {
					e, _ := json.Marshal(*subtest.expectProblemDetails)
					r, _ := json.Marshal(d)
					t.Errorf("Incorrect problem details:\nexpected\n%s\n, got\n%s", string(e), string(r))
				}
			}
		})
	}
}
