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

	"gopkg.in/yaml.v2"

	. "free5gc/lib/openapi/models"
	"free5gc/src/nssf/factory"
	. "free5gc/src/nssf/plugin"
	"free5gc/src/nssf/test"
)

var testingNssaiavailabilityStore = test.TestingNssaiavailability{
	ConfigFile: test.ConfigFileFromArgs,
	NfId:       "469de254-2fe5-4ca0-8381-af3f500af77c",
}

func checkAmfExist(nfId string) bool {
	for _, amfConfig := range factory.NssfConfig.Configuration.AmfList {
		if amfConfig.NfId == nfId {
			return true
		}
	}
	return false
}

func generatePatchAddRequest() PatchDocument {
	const jsonRequest = `
        [
            {
                "op": "add",
                "path": "/0/supportedSnssaiList/-",
                "value": {
                    "sst": 1,
                    "sd": "1"
                }
            }
        ]
    `

	var p PatchDocument
	if err := json.NewDecoder(strings.NewReader(jsonRequest)).Decode(&p); err != nil {
		fmt.Printf("Decode error: %v", err)
	}

	return p
}

func generatePatchCopyRequest() PatchDocument {
	const jsonRequest = `
        [
            {
                "op": "copy",
                "path": "/1/supportedSnssaiList/0",
                "from": "/0/supportedSnssaiList/0"
            }
        ]
    `

	var p PatchDocument
	if err := json.NewDecoder(strings.NewReader(jsonRequest)).Decode(&p); err != nil {
		fmt.Printf("Decode error: %v", err)
	}

	return p
}

func generatePatchMoveRequest() PatchDocument {
	const jsonRequest = `
        [
            {
                "op": "move",
                "path": "/0/supportedSnssaiList/1",
                "from": "/0/supportedSnssaiList/3"
            }
        ]
    `

	var p PatchDocument
	if err := json.NewDecoder(strings.NewReader(jsonRequest)).Decode(&p); err != nil {
		fmt.Printf("Decode error: %v", err)
	}

	return p
}

func generatePatchRemoveRequest() PatchDocument {
	const jsonRequest = `
        [
            {
                "op": "remove",
                "path": "/0/supportedSnssaiList/2"
            }
        ]
    `

	var p PatchDocument
	if err := json.NewDecoder(strings.NewReader(jsonRequest)).Decode(&p); err != nil {
		fmt.Printf("Decode error: %v", err)
	}

	return p
}

func generatePatchReplaceRequest() PatchDocument {
	const jsonRequest = `
        [
            {
                "op": "replace",
                "path": "/1/supportedSnssaiList/2",
                "value": {
                    "sst": 2
                }
            }
        ]
    `

	var p PatchDocument
	if err := json.NewDecoder(strings.NewReader(jsonRequest)).Decode(&p); err != nil {
		fmt.Printf("Decode error: %v", err)
	}

	return p
}

func generatePatchTestRequest() PatchDocument {
	const jsonRequest = `
        [
            {
                "op": "test",
                "path": "/1/supportedSnssaiList/1",
                "value": {
                    "sst": 1,
                    "sd": "1"
                }
            }
        ]
    `

	var p PatchDocument
	if err := json.NewDecoder(strings.NewReader(jsonRequest)).Decode(&p); err != nil {
		fmt.Printf("Decode error: %v", err)
	}

	return p
}

func generatePutRequest() NssaiAvailabilityInfo {
	const jsonRequest = `
        {
            "supportedNssaiAvailabilityData": [
                {
                    "tai": {
                        "plmnId": {
                            "mcc": "466",
                            "mnc": "92"
                        },
                        "tac": "33456"
                    },
                    "supportedSnssaiList": [
                        {
                            "sst": 1
                        },
                        {
                            "sst": 1,
                            "sd": "1"
                        },
                        {
                            "sst": 1,
                            "sd": "2"
                        }
                    ]
                },
                {
                    "tai": {
                        "plmnId": {
                            "mcc": "466",
                            "mnc": "92"
                        },
                        "tac": "33458"
                    },
                    "supportedSnssaiList": [
                        {
                            "sst": 1
                        },
                        {
                            "sst": 1,
                            "sd": "1"
                        },
                        {
                            "sst": 1,
                            "sd": "3"
                        }
                    ]
                }
            ],
            "supportedFeatures": ""
        }
    `

	var n NssaiAvailabilityInfo
	if err := json.NewDecoder(strings.NewReader(jsonRequest)).Decode(&n); err != nil {
		fmt.Printf("Decode error: %v", err)
	}

	return n
}

func TestNssaiavailabilityTemplate(t *testing.T) {
	t.Skip()

	// Tests may have different configuration files
	factory.InitConfigFactory(testingNssaiavailabilityStore.ConfigFile)

	d, _ := yaml.Marshal(*factory.NssfConfig.Info)
	t.Logf("%s", string(d))
}

func TestNssaiavailabilityDelete(t *testing.T) {
	factory.InitConfigFactory(testingNssaiavailabilityStore.ConfigFile)

	subtests := []struct {
		name                 string
		expectStatus         int
		expectProblemDetails *ProblemDetails
	}{
		{
			name:         "Delete",
			expectStatus: http.StatusNoContent,
		},
	}

	for _, subtest := range subtests {
		t.Run(subtest.name, func(t *testing.T) {
			var (
				status int
				d      ProblemDetails
			)

			status = nssaiavailabilityDelete(testingNssaiavailabilityStore.NfId, &d)

			if status == http.StatusNoContent {
				if checkAmfExist(testingNssaiavailabilityStore.NfId) == true {
					t.Errorf("AMF ID '%s' in configuration should be deleted, but still exists", testingNssaiavailabilityStore.NfId)
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

func TestNssaiavailabilityPatch(t *testing.T) {
	factory.InitConfigFactory(testingNssaiavailabilityStore.ConfigFile)

	subtests := []struct {
		name                                  string
		generateRequestBody                   func() PatchDocument
		expectStatus                          int
		expectAuthorizedNssaiAvailabilityInfo *AuthorizedNssaiAvailabilityInfo
		expectProblemDetails                  *ProblemDetails
	}{
		{
			name:                "Add",
			generateRequestBody: generatePatchAddRequest,
			expectStatus:        http.StatusOK,
			expectAuthorizedNssaiAvailabilityInfo: &AuthorizedNssaiAvailabilityInfo{
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
								Sd:  "2",
							},
							{
								Sst: 2,
							},
							{
								Sst: 1,
								Sd:  "1",
							},
						},
					},
					{
						Tai: &Tai{
							PlmnId: &PlmnId{
								Mcc: "466",
								Mnc: "92",
							},
							Tac: "33457",
						},
						SupportedSnssaiList: []Snssai{
							{
								Sst: 1,
								Sd:  "1",
							},
							{
								Sst: 1,
								Sd:  "2",
							},
						},
					},
				},
			},
		},
		{
			name:                "Copy",
			generateRequestBody: generatePatchCopyRequest,
			expectStatus:        http.StatusOK,
			expectAuthorizedNssaiAvailabilityInfo: &AuthorizedNssaiAvailabilityInfo{
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
								Sd:  "2",
							},
							{
								Sst: 2,
							},
							{
								Sst: 1,
								Sd:  "1",
							},
						},
					},
					{
						Tai: &Tai{
							PlmnId: &PlmnId{
								Mcc: "466",
								Mnc: "92",
							},
							Tac: "33457",
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
						},
					},
				},
			},
		},
		{
			name:                "Move",
			generateRequestBody: generatePatchMoveRequest,
			expectStatus:        http.StatusOK,
			expectAuthorizedNssaiAvailabilityInfo: &AuthorizedNssaiAvailabilityInfo{
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
							Tac: "33457",
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
						},
					},
				},
			},
		},
		{
			name:                "Remove",
			generateRequestBody: generatePatchRemoveRequest,
			expectStatus:        http.StatusOK,
			expectAuthorizedNssaiAvailabilityInfo: &AuthorizedNssaiAvailabilityInfo{
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
							Tac: "33457",
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
						},
					},
				},
			},
		},
		{
			name:                "Replace",
			generateRequestBody: generatePatchReplaceRequest,
			expectStatus:        http.StatusOK,
			expectAuthorizedNssaiAvailabilityInfo: &AuthorizedNssaiAvailabilityInfo{
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
							Tac: "33457",
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
								Sst: 2,
							},
						},
					},
				},
			},
		},
		{
			name:                "Test",
			generateRequestBody: generatePatchTestRequest,
			expectStatus:        http.StatusOK,
			expectAuthorizedNssaiAvailabilityInfo: &AuthorizedNssaiAvailabilityInfo{
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
							Tac: "33457",
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
								Sst: 2,
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
				p      PatchDocument
				status int
				a      AuthorizedNssaiAvailabilityInfo
				d      ProblemDetails
			)

			if subtest.generateRequestBody != nil {
				p = subtest.generateRequestBody()
			}

			status = nssaiavailabilityPatch(testingNssaiavailabilityStore.NfId, p, &a, &d)

			if status == http.StatusOK {
				if reflect.DeepEqual(a, *subtest.expectAuthorizedNssaiAvailabilityInfo) == false {
					e, _ := json.Marshal(*subtest.expectAuthorizedNssaiAvailabilityInfo)
					r, _ := json.Marshal(a)
					t.Errorf("Incorrect authorized NSSAI availability info:\nexpected\n%s\n, got\n%s", string(e), string(r))
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

func TestNssaiavailabilityPut(t *testing.T) {
	factory.InitConfigFactory(testingNssaiavailabilityStore.ConfigFile)

	subtests := []struct {
		name                                  string
		generateRequestBody                   func() NssaiAvailabilityInfo
		expectStatus                          int
		expectAuthorizedNssaiAvailabilityInfo *AuthorizedNssaiAvailabilityInfo
		expectProblemDetails                  *ProblemDetails
	}{
		{
			name:                "Create and Replace",
			generateRequestBody: generatePutRequest,
			expectStatus:        http.StatusOK,
			expectAuthorizedNssaiAvailabilityInfo: &AuthorizedNssaiAvailabilityInfo{
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
				SupportedFeatures: "",
			},
		},
	}

	for _, subtest := range subtests {
		t.Run(subtest.name, func(t *testing.T) {
			var (
				n      NssaiAvailabilityInfo
				status int
				a      AuthorizedNssaiAvailabilityInfo
				d      ProblemDetails
			)

			if subtest.generateRequestBody != nil {
				n = subtest.generateRequestBody()
			}

			status = nssaiavailabilityPut(testingNssaiavailabilityStore.NfId, n, &a, &d)

			if status == http.StatusOK {
				if reflect.DeepEqual(a, *subtest.expectAuthorizedNssaiAvailabilityInfo) == false {
					e, _ := json.Marshal(*subtest.expectAuthorizedNssaiAvailabilityInfo)
					r, _ := json.Marshal(a)
					t.Errorf("Incorrect authorized NSSAI availability info:\nexpected\n%s\n, got\n%s", string(e), string(r))
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
