/*
 * NSSF NS Selection
 *
 * NSSF Network Slice Selection Service
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
	"free5gc/src/nssf/util"
)

var testingNsselectionForRegistration = test.TestingNsselection{
	ConfigFile: test.ConfigFileFromArgs,
	GenerateNonRoamingQueryParameter: func() NsselectionQueryParameter {
		const jsonQuery = `
            {
                "nf-type": "AMF",
                "nf-id": "469de254-2fe5-4ca0-8381-af3f500af77c",
                "slice-info-request-for-registration": {
                    "subscribedNssai": [
                        {
                            "subscribedSnssai": {
                                "sst": 1
                            },
                            "defaultIndication": false
                        },
                        {
                            "subscribedSnssai": {
                                "sst": 1,
                                "sd": "1"
                            },
                            "defaultIndication": true
                        }
                    ],
                    "allowedNssaiCurrentAccess": {
                        "allowedSnssaiList": [
                            {
                                "allowedSnssai": {
                                    "sst": 1,
                                    "sd": "1"
                                },
                                "nsiInformationList": [
                                    {
                                        "nrfId": "http://free5gc-nrf.nctu.me:29510/nnrf-nfm/v1/nf-instances",
                                        "nsiId": "1"
                                    }
                                ]
                            }
                        ],
                        "accessType": "3GPP_ACCESS"
                    },
                    "requestedNssai": [
                        {
                            "sst": 1,
                            "sd": "1"
                        },
                        {
                            "sst": 1,
                            "sd": "2"
                        },
                        {
                            "sst": 1,
                            "sd": "3"
                        }
                    ],
                    "defaultConfiguredSnssaiInd": false
                },
                "tai": {
                    "plmnId": {
                        "mcc": "466",
                        "mnc": "92"
                    },
                    "tac": "33456"
                }
            }
    `

		var p NsselectionQueryParameter
		if err := json.NewDecoder(strings.NewReader(jsonQuery)).Decode(&p); err != nil {
			fmt.Printf("Decode error: %v", err)
		}

		return p
	},
	GenerateRoamingQueryParameter: func() NsselectionQueryParameter {
		const jsonQuery = `
            {
                "nf-type": "AMF",
                "nf-id": "469de254-2fe5-4ca0-8381-af3f500af77c",
                "slice-info-request-for-registration": {
                    "subscribedNssai": [
                        {
                            "subscribedSnssai": {
                                "sst": 1
                            },
                            "defaultIndication": false
                        },
                        {
                            "subscribedSnssai": {
                                "sst": 1,
                                "sd": "1"
                            },
                            "defaultIndication": true
                        }
                    ],
                    "allowedNssaiCurrentAccess": {
                        "allowedSnssaiList": [
                            {
                                "allowedSnssai": {
                                    "sst": 1,
                                    "sd": "1"
                                },
                                "nsiInformationList": [
                                    {
                                        "nrfId": "http://free5gc-nrf.nctu.me:29510/nnrf-nfm/v1/nf-instances",
                                        "nsiId": "1"
                                    }
                                ],
                                "mappedHomeSnssai": {
                                    "sst": 1,
                                    "sd": "1"
                                }
                            }
                        ],
                        "accessType": "3GPP_ACCESS"
                    },
                    "requestedNssai": [
                        {
                            "sst": 1,
                            "sd": "1"
                        },
                        {
                            "sst": 1,
                            "sd": "2"
                        },
                        {
                            "sst": 1,
                            "sd": "3"
                        }
                    ],
                    "defaultConfiguredSnssaiInd": false,
                    "mappingofNssai": [
                        {
                            "servingSnssai": {
                                "sst": 1,
                                "sd": "1"
                            },
                            "homeSnssai": {
                                "sst": 1,
                                "sd": "1"
                            }
                        },
                        {
                            "servingSnssai":{
                                "sst": 1,
                                "sd": "2"
                            },
                            "homeSnssai": {
                                "sst": 1,
                                "sd": "3"
                            }
                        },
                        {
                            "servingSnssai":{
                                "sst": 1,
                                "sd": "3"
                            },
                            "homeSnssai": {
                                "sst": 1,
                                "sd": "4"
                            }
                        }
                    ],
                    "requestMapping": false
                },
                "home-plmn-id": {
                    "mcc": "440",
                    "mnc": "10"
                },
                "tai": {
                    "plmnId": {
                        "mcc": "466",
                        "mnc": "92"
                    },
                    "tac": "33456"
                }
            }
        `

		var p NsselectionQueryParameter
		if err := json.NewDecoder(strings.NewReader(jsonQuery)).Decode(&p); err != nil {
			fmt.Printf("Decode error: %v", err)
		}

		return p
	},
}

func setUnsupportedHomePlmnId(p *NsselectionQueryParameter) {
	*p.HomePlmnId = PlmnId{
		Mcc: "123",
		Mnc: "456",
	}
}

func setUnsupportedTai(p *NsselectionQueryParameter) {
	*p.Tai = Tai{
		PlmnId: &PlmnId{
			Mcc: "466",
			Mnc: "92",
		},
		Tac: "12345",
	}
}

func setUnsupportedSnssai(p *NsselectionQueryParameter) {
	p.SliceInfoRequestForRegistration.RequestedNssai = append(p.SliceInfoRequestForRegistration.RequestedNssai, Snssai{
		Sst: 9,
		Sd:  "9",
	})
}

func setDefaultConfiguredSnssai(p *NsselectionQueryParameter) {
	p.SliceInfoRequestForRegistration.RequestedNssai = []Snssai{
		{
			Sst: 1,
		},
		{
			Sst: 2,
		},
	}

	p.SliceInfoRequestForRegistration.DefaultConfiguredSnssaiInd = true
}

func setAllRequestedSnssaiValid(p *NsselectionQueryParameter) {
	p.SliceInfoRequestForRegistration.RequestedNssai = []Snssai{
		{
			Sst: 1,
		},
		{
			Sst: 1,
			Sd:  "1",
		},
	}
}

func setRequestMapping(p *NsselectionQueryParameter) {
	p.SliceInfoRequestForRegistration.SNssaiForMapping = []Snssai{
		{
			Sst: 1,
			Sd:  "3",
		},
	}

	p.SliceInfoRequestForRegistration.RequestMapping = true
}

func removeRequestedNssai(p *NsselectionQueryParameter) {
	p.SliceInfoRequestForRegistration.RequestedNssai = []Snssai{}
}

func TestNsselectionForRegistrationTemplate(t *testing.T) {
	t.Skip()

	// Tests may have different configuration files
	factory.InitConfigFactory(testingNsselectionForRegistration.ConfigFile)

	d, _ := yaml.Marshal(*factory.NssfConfig.Info)
	t.Logf("%s", string(d))
}

func TestNsselectionForRegistrationGeneral(t *testing.T) {
	factory.InitConfigFactory(testingNsselectionForRegistration.ConfigFile)

	subtests := []struct {
		name                             string
		modifyQueryParameter             func(*NsselectionQueryParameter)
		expectStatus                     int
		expectAuthorizedNetworkSliceInfo *AuthorizedNetworkSliceInfo
		expectProblemDetails             *ProblemDetails
	}{
		{
			name:                 "Unsupported TA",
			modifyQueryParameter: setUnsupportedTai,
			expectStatus:         http.StatusOK,
			expectAuthorizedNetworkSliceInfo: &AuthorizedNetworkSliceInfo{
				RejectedNssaiInTa: []Snssai{
					{
						Sst: 1,
						Sd:  "1",
					},
					{
						Sst: 1,
						Sd:  "2",
					},
					{
						Sst: 1,
						Sd:  "3",
					},
				},
			},
		},
		{
			name:                 "Unsupported S-NSSAI",
			modifyQueryParameter: setUnsupportedSnssai,
			expectStatus:         http.StatusForbidden,
			expectProblemDetails: &ProblemDetails{
				Title:  util.UNSUPPORTED_RESOURCE,
				Status: http.StatusForbidden,
				Detail: "S-NSSAI in Requested NSSAI is not supported in PLMN",
				Cause:  "SNSSAI_NOT_SUPPORTED",
			},
		},
		{
			name:                 "Default Configured S-NSSAI",
			modifyQueryParameter: setDefaultConfiguredSnssai,
			expectStatus:         http.StatusOK,
			expectAuthorizedNetworkSliceInfo: &AuthorizedNetworkSliceInfo{
				AllowedNssaiList: []AllowedNssai{
					{
						AllowedSnssaiList: []AllowedSnssai{
							{
								AllowedSnssai: &Snssai{
									Sst: 1,
								},
								NsiInformationList: []NsiInformation{
									{
										NrfId: "http://free5gc-nrf-10.nctu.me:29510/nnrf-nfm/v1/nf-instances",
										NsiId: "10",
									},
								},
							},
						},
						AccessType: func() AccessType { a := AccessType__3_GPP_ACCESS; return a }(),
					},
				},
				ConfiguredNssai: []ConfiguredSnssai{
					{
						ConfiguredSnssai: &Snssai{
							Sst: 1,
						},
					},
				},
				RejectedNssaiInPlmn: []Snssai{
					{
						Sst: 2,
					},
				},
			},
		},
	}

	for _, subtest := range subtests {
		t.Run(subtest.name, func(t *testing.T) {
			var (
				status int
				a      AuthorizedNetworkSliceInfo
				d      ProblemDetails
			)

			p := testingNsselectionForRegistration.GenerateNonRoamingQueryParameter()

			if subtest.modifyQueryParameter != nil {
				subtest.modifyQueryParameter(&p)
			}

			status = nsselectionForRegistration(p, &a, &d)

			if status != subtest.expectStatus {
				t.Errorf("Incorrect status code: expected %d, got %d", subtest.expectStatus, status)
			}

			if status == http.StatusOK {
				if reflect.DeepEqual(a, *subtest.expectAuthorizedNetworkSliceInfo) == false {
					e, _ := json.Marshal(*subtest.expectAuthorizedNetworkSliceInfo)
					r, _ := json.Marshal(a)
					t.Errorf("Incorrect authorized network slice info:\nexpected\n%s\n, got\n%s", string(e), string(r))
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

func TestNsselectionForRegistrationNonRoaming(t *testing.T) {
	factory.InitConfigFactory(testingNsselectionForRegistration.ConfigFile)

	subtests := []struct {
		name                             string
		modifyQueryParameter             func(*NsselectionQueryParameter)
		expectStatus                     int
		expectAuthorizedNetworkSliceInfo *AuthorizedNetworkSliceInfo
		expectProblemDetails             *ProblemDetails
	}{
		{
			name:                 "Requested with All Allowed",
			modifyQueryParameter: setAllRequestedSnssaiValid,
			expectStatus:         http.StatusOK,
			expectAuthorizedNetworkSliceInfo: &AuthorizedNetworkSliceInfo{
				AllowedNssaiList: []AllowedNssai{
					{
						AllowedSnssaiList: []AllowedSnssai{
							{
								AllowedSnssai: &Snssai{
									Sst: 1,
								},
								NsiInformationList: []NsiInformation{
									{
										NrfId: "http://free5gc-nrf-10.nctu.me:29510/nnrf-nfm/v1/nf-instances",
										NsiId: "10",
									},
								},
							},
							{
								AllowedSnssai: &Snssai{
									Sst: 1,
									Sd:  "1",
								},
								NsiInformationList: []NsiInformation{
									{
										NrfId: "http://free5gc-nrf-11.nctu.me:29510/nnrf-nfm/v1/nf-instances",
										NsiId: "11",
									},
								},
							},
						},
						AccessType: func() AccessType { a := AccessType__3_GPP_ACCESS; return a }(),
					},
				},
				TargetAmfSet: "2",
				NrfAmfSet:    "http://free5gc-nrf.nctu.me:8084/nnrf-nfm/v1/nf-instances",
			},
		},
		{
			name:         "Requested with Two Rejected in PLMN and TA and One Allowed",
			expectStatus: http.StatusOK,
			expectAuthorizedNetworkSliceInfo: &AuthorizedNetworkSliceInfo{
				AllowedNssaiList: []AllowedNssai{
					{
						AllowedSnssaiList: []AllowedSnssai{
							{
								AllowedSnssai: &Snssai{
									Sst: 1,
									Sd:  "1",
								},
								NsiInformationList: []NsiInformation{
									{
										NrfId: "http://free5gc-nrf-11.nctu.me:29510/nnrf-nfm/v1/nf-instances",
										NsiId: "11",
									},
								},
							},
						},
						AccessType: func() AccessType { a := AccessType__3_GPP_ACCESS; return a }(),
					},
				},
				ConfiguredNssai: []ConfiguredSnssai{
					{
						ConfiguredSnssai: &Snssai{
							Sst: 1,
						},
					},
					{
						ConfiguredSnssai: &Snssai{
							Sst: 1,
							Sd:  "1",
						},
					},
				},
				CandidateAmfList: []string{
					"ffa2e8d7-3275-49c7-8631-6af1df1d9d26",
					"0e8831c3-6286-4689-ab27-1e2161e15cb1",
					"a1fba9ba-2e39-4e22-9c74-f749da571d0d",
				},
				RejectedNssaiInPlmn: []Snssai{
					{
						Sst: 1,
						Sd:  "2",
					},
				},
				RejectedNssaiInTa: []Snssai{
					{
						Sst: 1,
						Sd:  "3",
					},
				},
			},
		},
		{
			name:                 "No Requested and One Subscribed Marked as Default",
			modifyQueryParameter: removeRequestedNssai,
			expectStatus:         http.StatusOK,
			expectAuthorizedNetworkSliceInfo: &AuthorizedNetworkSliceInfo{
				AllowedNssaiList: []AllowedNssai{
					{
						AllowedSnssaiList: []AllowedSnssai{
							{
								AllowedSnssai: &Snssai{
									Sst: 1,
									Sd:  "1",
								},
								NsiInformationList: []NsiInformation{
									{
										NrfId: "http://free5gc-nrf-11.nctu.me:29510/nnrf-nfm/v1/nf-instances",
										NsiId: "11",
									},
								},
							},
						},
						AccessType: func() AccessType { a := AccessType__3_GPP_ACCESS; return a }(),
					},
				},
				ConfiguredNssai: []ConfiguredSnssai{
					{
						ConfiguredSnssai: &Snssai{
							Sst: 1,
						},
					},
					{
						ConfiguredSnssai: &Snssai{
							Sst: 1,
							Sd:  "1",
						},
					},
				},
				CandidateAmfList: []string{
					"ffa2e8d7-3275-49c7-8631-6af1df1d9d26",
					"0e8831c3-6286-4689-ab27-1e2161e15cb1",
					"a1fba9ba-2e39-4e22-9c74-f749da571d0d",
				},
			},
		},
	}

	for _, subtest := range subtests {
		t.Run(subtest.name, func(t *testing.T) {
			var (
				status int
				a      AuthorizedNetworkSliceInfo
				d      ProblemDetails
			)

			p := testingNsselectionForRegistration.GenerateNonRoamingQueryParameter()

			if subtest.modifyQueryParameter != nil {
				subtest.modifyQueryParameter(&p)
			}

			status = nsselectionForRegistration(p, &a, &d)

			if status != subtest.expectStatus {
				t.Errorf("Incorrect status code: expected %d, got %d", subtest.expectStatus, status)
			}

			if status == http.StatusOK {
				if reflect.DeepEqual(a, *subtest.expectAuthorizedNetworkSliceInfo) == false {
					e, _ := json.Marshal(*subtest.expectAuthorizedNetworkSliceInfo)
					r, _ := json.Marshal(a)
					t.Errorf("Incorrect authorized network slice info:\nexpected\n%s\n, got\n%s", string(e), string(r))
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

func TestNsselectionForRegistrationRoaming(t *testing.T) {
	factory.InitConfigFactory(testingNsselectionForRegistration.ConfigFile)

	subtests := []struct {
		name                             string
		modifyQueryParameter             func(*NsselectionQueryParameter)
		expectStatus                     int
		expectAuthorizedNetworkSliceInfo *AuthorizedNetworkSliceInfo
		expectProblemDetails             *ProblemDetails
	}{
		{
			name:                 "Unsupported HPLMN",
			modifyQueryParameter: setUnsupportedHomePlmnId,
			expectStatus:         http.StatusOK,
			expectAuthorizedNetworkSliceInfo: &AuthorizedNetworkSliceInfo{
				RejectedNssaiInPlmn: []Snssai{
					{
						Sst: 1,
						Sd:  "1",
					},
					{
						Sst: 1,
						Sd:  "2",
					},
					{
						Sst: 1,
						Sd:  "3",
					},
				},
			},
		},
		{
			name:                 "Requested with All Allowed",
			modifyQueryParameter: setAllRequestedSnssaiValid,
			expectStatus:         http.StatusOK,
			expectAuthorizedNetworkSliceInfo: &AuthorizedNetworkSliceInfo{
				AllowedNssaiList: []AllowedNssai{
					{
						AllowedSnssaiList: []AllowedSnssai{
							{
								AllowedSnssai: &Snssai{
									Sst: 1,
								},
								NsiInformationList: []NsiInformation{
									{
										NrfId: "http://free5gc-nrf-10.nctu.me:29510/nnrf-nfm/v1/nf-instances",
										NsiId: "10",
									},
								},
							},
							{
								AllowedSnssai: &Snssai{
									Sst: 1,
									Sd:  "1",
								},
								NsiInformationList: []NsiInformation{
									{
										NrfId: "http://free5gc-nrf-11.nctu.me:29510/nnrf-nfm/v1/nf-instances",
										NsiId: "11",
									},
								},
								MappedHomeSnssai: &Snssai{
									Sst: 1,
									Sd:  "1",
								},
							},
						},
						AccessType: func() AccessType { a := AccessType__3_GPP_ACCESS; return a }(),
					},
				},
				TargetAmfSet: "2",
				NrfAmfSet:    "http://free5gc-nrf.nctu.me:8084/nnrf-nfm/v1/nf-instances",
			},
		},
		{
			name:         "Requested with Two Rejected in PLMN and TA and One Allowed",
			expectStatus: http.StatusOK,
			expectAuthorizedNetworkSliceInfo: &AuthorizedNetworkSliceInfo{
				AllowedNssaiList: []AllowedNssai{
					{
						AllowedSnssaiList: []AllowedSnssai{
							{
								AllowedSnssai: &Snssai{
									Sst: 1,
									Sd:  "1",
								},
								NsiInformationList: []NsiInformation{
									{
										NrfId: "http://free5gc-nrf-11.nctu.me:29510/nnrf-nfm/v1/nf-instances",
										NsiId: "11",
									},
								},
								MappedHomeSnssai: &Snssai{
									Sst: 1,
									Sd:  "1",
								},
							},
						},
						AccessType: func() AccessType { a := AccessType__3_GPP_ACCESS; return a }(),
					},
				},
				ConfiguredNssai: []ConfiguredSnssai{
					{
						ConfiguredSnssai: &Snssai{
							Sst: 1,
						},
					},
					{
						ConfiguredSnssai: &Snssai{
							Sst: 1,
							Sd:  "1",
						},
						MappedHomeSnssai: &Snssai{
							Sst: 1,
							Sd:  "1",
						},
					},
				},
				CandidateAmfList: []string{
					"ffa2e8d7-3275-49c7-8631-6af1df1d9d26",
					"0e8831c3-6286-4689-ab27-1e2161e15cb1",
					"a1fba9ba-2e39-4e22-9c74-f749da571d0d",
				},
				RejectedNssaiInPlmn: []Snssai{
					{
						Sst: 1,
						Sd:  "2",
					},
				},
				RejectedNssaiInTa: []Snssai{
					{
						Sst: 1,
						Sd:  "3",
					},
				},
			},
		},
		{
			name:                 "No Requested and One Subscribed Marked as Default",
			modifyQueryParameter: removeRequestedNssai,
			expectStatus:         http.StatusOK,
			expectAuthorizedNetworkSliceInfo: &AuthorizedNetworkSliceInfo{
				AllowedNssaiList: []AllowedNssai{
					{
						AllowedSnssaiList: []AllowedSnssai{
							{
								AllowedSnssai: &Snssai{
									Sst: 1,
									Sd:  "1",
								},
								NsiInformationList: []NsiInformation{
									{
										NrfId: "http://free5gc-nrf-11.nctu.me:29510/nnrf-nfm/v1/nf-instances",
										NsiId: "11",
									},
								},
								MappedHomeSnssai: &Snssai{
									Sst: 1,
									Sd:  "1",
								},
							},
						},
						AccessType: func() AccessType { a := AccessType__3_GPP_ACCESS; return a }(),
					},
				},
				ConfiguredNssai: []ConfiguredSnssai{
					{
						ConfiguredSnssai: &Snssai{
							Sst: 1,
						},
					},
					{
						ConfiguredSnssai: &Snssai{
							Sst: 1,
							Sd:  "1",
						},
						MappedHomeSnssai: &Snssai{
							Sst: 1,
							Sd:  "1",
						},
					},
				},
				CandidateAmfList: []string{
					"ffa2e8d7-3275-49c7-8631-6af1df1d9d26",
					"0e8831c3-6286-4689-ab27-1e2161e15cb1",
					"a1fba9ba-2e39-4e22-9c74-f749da571d0d",
				},
			},
		},
		{
			name:                 "Request Mapping",
			modifyQueryParameter: setRequestMapping,
			expectStatus:         http.StatusOK,
			expectAuthorizedNetworkSliceInfo: &AuthorizedNetworkSliceInfo{
				AllowedNssaiList: []AllowedNssai{
					{
						AllowedSnssaiList: []AllowedSnssai{
							{
								AllowedSnssai: &Snssai{
									Sst: 1,
									Sd:  "1",
								},
								MappedHomeSnssai: &Snssai{
									Sst: 1,
									Sd:  "1",
								},
							},
							{
								AllowedSnssai: &Snssai{
									Sst: 1,
									Sd:  "2",
								},
								MappedHomeSnssai: &Snssai{
									Sst: 1,
									Sd:  "3",
								},
							},
						},
						AccessType: func() AccessType { a := AccessType__3_GPP_ACCESS; return a }(),
					},
				},
			},
		},
	}

	for _, subtest := range subtests {
		t.Run(subtest.name, func(t *testing.T) {
			var (
				status int
				a      AuthorizedNetworkSliceInfo
				d      ProblemDetails
			)

			p := testingNsselectionForRegistration.GenerateRoamingQueryParameter()

			if subtest.modifyQueryParameter != nil {
				subtest.modifyQueryParameter(&p)
			}

			status = nsselectionForRegistration(p, &a, &d)

			if status != subtest.expectStatus {
				t.Errorf("Incorrect status code: expected %d, got %d", subtest.expectStatus, status)
			}

			if status == http.StatusOK {
				if reflect.DeepEqual(a, *subtest.expectAuthorizedNetworkSliceInfo) == false {
					e, _ := json.Marshal(*subtest.expectAuthorizedNetworkSliceInfo)
					r, _ := json.Marshal(a)
					t.Errorf("Incorrect authorized network slice info:\nexpected\n%s\n, got\n%s", string(e), string(r))
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
