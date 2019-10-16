/*
 * NSSF Utility
 */

package util

import (
	"encoding/json"
	"reflect"
	"testing"

	"gopkg.in/yaml.v2"

	. "free5gc/lib/openapi/models"
	"free5gc/src/nssf/factory"
	"free5gc/src/nssf/test"
)

var testingUtil = test.TestingUtil{
	ConfigFile: test.ConfigFileFromArgs,
}

func TestPluginTemplate(t *testing.T) {
	t.Skip()

	factory.InitConfigFactory(testingUtil.ConfigFile)

	d, _ := yaml.Marshal(*factory.NssfConfig.Info)
	t.Logf("%s", string(d))
}

func TestAddAmfInformation(t *testing.T) {
	factory.InitConfigFactory(testingUtil.ConfigFile)

	subtests := []struct {
		name                             string
		tai                              *Tai
		authorizedNetworkSliceInfo       *AuthorizedNetworkSliceInfo
		expectAuthorizedNetworkSliceInfo *AuthorizedNetworkSliceInfo
	}{
		{
			name: "Add Candidate AMF List from AMF List",
			tai: &Tai{
				PlmnId: &PlmnId{
					Mcc: "466",
					Mnc: "92",
				},
				Tac: "33456",
			},
			authorizedNetworkSliceInfo: &AuthorizedNetworkSliceInfo{
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
									Sd:  "2",
								},
								NsiInformationList: []NsiInformation{
									{
										NrfId: "http://free5gc-nrf-12-1.nctu.me:29510/nnrf-nfm/v1/nf-instances",
										NsiId: "12",
									},
									{
										NrfId: "http://free5gc-nrf-12-2.nctu.me:29510/nnrf-nfm/v1/nf-instances",
										NsiId: "12",
									},
								},
							},
						},
						AccessType: func() AccessType { a := AccessType__3_GPP_ACCESS; return a }(),
					},
				},
			},
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
									Sd:  "2",
								},
								NsiInformationList: []NsiInformation{
									{
										NrfId: "http://free5gc-nrf-12-1.nctu.me:29510/nnrf-nfm/v1/nf-instances",
										NsiId: "12",
									},
									{
										NrfId: "http://free5gc-nrf-12-2.nctu.me:29510/nnrf-nfm/v1/nf-instances",
										NsiId: "12",
									},
								},
							},
						},
						AccessType: func() AccessType { a := AccessType__3_GPP_ACCESS; return a }(),
					},
				},
				CandidateAmfList: []string{
					"469de254-2fe5-4ca0-8381-af3f500af77c",
					"b9e6e2cb-5ce8-4cb6-9173-a266dd9a2f0c",
				},
			},
		},
		{
			name: "Add Candidate AMF List from AMF Set",
			tai: &Tai{
				PlmnId: &PlmnId{
					Mcc: "466",
					Mnc: "92",
				},
				Tac: "33456",
			},
			authorizedNetworkSliceInfo: &AuthorizedNetworkSliceInfo{
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
							{
								AllowedSnssai: &Snssai{
									Sst: 1,
									Sd:  "2",
								},
								NsiInformationList: []NsiInformation{
									{
										NrfId: "http://free5gc-nrf-12-1.nctu.me:29510/nnrf-nfm/v1/nf-instances",
										NsiId: "12",
									},
									{
										NrfId: "http://free5gc-nrf-12-2.nctu.me:29510/nnrf-nfm/v1/nf-instances",
										NsiId: "12",
									},
								},
							},
						},
						AccessType: func() AccessType { a := AccessType__3_GPP_ACCESS; return a }(),
					},
				},
			},
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
							{
								AllowedSnssai: &Snssai{
									Sst: 1,
									Sd:  "2",
								},
								NsiInformationList: []NsiInformation{
									{
										NrfId: "http://free5gc-nrf-12-1.nctu.me:29510/nnrf-nfm/v1/nf-instances",
										NsiId: "12",
									},
									{
										NrfId: "http://free5gc-nrf-12-2.nctu.me:29510/nnrf-nfm/v1/nf-instances",
										NsiId: "12",
									},
								},
							},
						},
						AccessType: func() AccessType { a := AccessType__3_GPP_ACCESS; return a }(),
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
			name: "Add Target AMF Set",
			tai: &Tai{
				PlmnId: &PlmnId{
					Mcc: "466",
					Mnc: "92",
				},
				Tac: "33456",
			},
			authorizedNetworkSliceInfo: &AuthorizedNetworkSliceInfo{
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
										NrfId: "http://free5gc-nrf.nctu.me:8081/nnrf-nfm/v1/nf-instances",
										NsiId: "1",
									},
								},
							},
							{
								AllowedSnssai: &Snssai{
									Sst: 1,
									Sd:  "3",
								},
								NsiInformationList: []NsiInformation{
									{
										NrfId: "http://free5gc-nrf.nctu.me:8084/nnrf-nfm/v1/nf-instances",
									},
								},
							},
						},
						AccessType: func() AccessType { a := AccessType__3_GPP_ACCESS; return a }(),
					},
				},
			},
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
										NrfId: "http://free5gc-nrf.nctu.me:8081/nnrf-nfm/v1/nf-instances",
										NsiId: "1",
									},
								},
							},
							{
								AllowedSnssai: &Snssai{
									Sst: 1,
									Sd:  "3",
								},
								NsiInformationList: []NsiInformation{
									{
										NrfId: "http://free5gc-nrf.nctu.me:8084/nnrf-nfm/v1/nf-instances",
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
	}

	for _, subtest := range subtests {
		t.Run(subtest.name, func(t *testing.T) {
			a := *subtest.authorizedNetworkSliceInfo

			AddAmfInformation(*subtest.tai, &a)

			if reflect.DeepEqual(a, *subtest.expectAuthorizedNetworkSliceInfo) != true {
				e, _ := json.Marshal(*subtest.expectAuthorizedNetworkSliceInfo)
				r, _ := json.Marshal(a)
				t.Errorf("Incorrect authorized network slice info:\nexpected\n%s\n, got\n%s", string(e), string(r))
			}
		})
	}
}
