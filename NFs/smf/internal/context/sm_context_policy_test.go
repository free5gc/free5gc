package context_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/free5gc/openapi/models"
	smf_context "github.com/free5gc/smf/internal/context"
	"github.com/free5gc/smf/pkg/factory"
)

var userPlaneConfig = factory.UserPlaneInformation{
	UPNodes: map[string]*factory.UPNode{
		"GNodeB": {
			Type: "AN",
		},
		"UPF1": {
			Type:   "UPF",
			NodeID: "10.4.0.11",
			Addr:   "10.4.0.11",
			SNssaiInfos: []*factory.SnssaiUpfInfoItem{
				{
					SNssai: &models.Snssai{
						Sst: 1,
						Sd:  "010203",
					},
					DnnUpfInfoList: []*factory.DnnUpfInfoItem{
						{
							Dnn:      "internet",
							DnaiList: []string{"mec"},
						},
					},
				},
			},
			InterfaceUpfInfoList: []*factory.InterfaceUpfInfoItem{
				{
					InterfaceType: "N3",
					Endpoints: []string{
						"10.3.0.11",
					},
					NetworkInstances: []string{"internet"},
				},
				{
					InterfaceType: "N9",
					Endpoints: []string{
						"10.3.0.11",
					},
					NetworkInstances: []string{"internet"},
				},
			},
		},
		"UPF2": {
			Type:   "UPF",
			NodeID: "10.4.0.12",
			Addr:   "10.4.0.12",
			SNssaiInfos: []*factory.SnssaiUpfInfoItem{
				{
					SNssai: &models.Snssai{
						Sst: 1,
						Sd:  "010203",
					},
					DnnUpfInfoList: []*factory.DnnUpfInfoItem{
						{
							Dnn: "internet",
							Pools: []*factory.UEIPPool{
								{Cidr: "10.60.0.0/16"},
							},
						},
					},
				},
			},
			InterfaceUpfInfoList: []*factory.InterfaceUpfInfoItem{
				{
					InterfaceType: "N9",
					Endpoints: []string{
						"10.3.0.12",
					},
					NetworkInstances: []string{"internet"},
				},
			},
		},
	},
	Links: []*factory.UPLink{
		{
			A: "GNodeB",
			B: "UPF1",
		},
		{
			A: "UPF1",
			B: "UPF2",
		},
	},
}

var testConfig = factory.Config{
	Info: &factory.Info{
		Version:     "1.0.0",
		Description: "SMF procdeure test configuration",
	},
	Configuration: &factory.Configuration{
		Sbi: &factory.Sbi{
			Scheme:       "http",
			RegisterIPv4: "127.0.0.1",
			BindingIPv4:  "127.0.0.1",
			Port:         8000,
		},
		UserPlaneInformation: userPlaneConfig,
	},
}

func initConfig() {
	smf_context.InitSmfContext(&testConfig)
	factory.SmfConfig = &testConfig
}

func TestApplySessionRules(t *testing.T) {
	initConfig()

	testCases := []struct {
		name              string
		decision          *models.SmPolicyDecision
		noErr             bool
		expectedSessRules map[string]*smf_context.SessionRule
	}{
		{
			name:  "nil decision",
			noErr: false,
		},
		{
			name: "Install first session rule",
			decision: &models.SmPolicyDecision{
				SessRules: map[string]*models.SessionRule{
					"SessRuleId-1": {
						AuthSessAmbr: &models.Ambr{
							Uplink:   "1000 Kbps",
							Downlink: "1000 Kbps",
						},
						AuthDefQos: &models.AuthorizedDefaultQos{
							Var5qi: 9,
							Arp: &models.Arp{
								PriorityLevel: 8,
							},
							PriorityLevel: 8,
						},
						SessRuleId: "SessRuleId-1",
					},
				},
			},
			expectedSessRules: map[string]*smf_context.SessionRule{
				"SessRuleId-1": {
					SessionRule: &models.SessionRule{
						AuthSessAmbr: &models.Ambr{
							Uplink:   "1000 Kbps",
							Downlink: "1000 Kbps",
						},
						AuthDefQos: &models.AuthorizedDefaultQos{
							Var5qi: 9,
							Arp: &models.Arp{
								PriorityLevel: 8,
							},
							PriorityLevel: 8,
						},
						SessRuleId: "SessRuleId-1",
					},
					DefQosQFI: 1,
				},
			},
			noErr: true,
		},
		{
			name: "Install second session rule",
			decision: &models.SmPolicyDecision{
				SessRules: map[string]*models.SessionRule{
					"SessRuleId-2": {
						AuthSessAmbr: &models.Ambr{
							Uplink:   "1000 Kbps",
							Downlink: "1000 Kbps",
						},
						AuthDefQos: &models.AuthorizedDefaultQos{
							Var5qi: 9,
							Arp: &models.Arp{
								PriorityLevel: 8,
							},
							PriorityLevel: 8,
						},
						SessRuleId: "SessRuleId-2",
					},
				},
			},
			expectedSessRules: map[string]*smf_context.SessionRule{
				"SessRuleId-1": {
					SessionRule: &models.SessionRule{
						AuthSessAmbr: &models.Ambr{
							Uplink:   "1000 Kbps",
							Downlink: "1000 Kbps",
						},
						AuthDefQos: &models.AuthorizedDefaultQos{
							Var5qi: 9,
							Arp: &models.Arp{
								PriorityLevel: 8,
							},
							PriorityLevel: 8,
						},
						SessRuleId: "SessRuleId-1",
					},
					DefQosQFI: 1,
				},
				"SessRuleId-2": {
					SessionRule: &models.SessionRule{
						AuthSessAmbr: &models.Ambr{
							Uplink:   "1000 Kbps",
							Downlink: "1000 Kbps",
						},
						AuthDefQos: &models.AuthorizedDefaultQos{
							Var5qi: 9,
							Arp: &models.Arp{
								PriorityLevel: 8,
							},
							PriorityLevel: 8,
						},
						SessRuleId: "SessRuleId-2",
					},
					DefQosQFI: 1,
				},
			},
			noErr: true,
		},
		{
			name: "Modify first session rule",
			decision: &models.SmPolicyDecision{
				SessRules: map[string]*models.SessionRule{
					"SessRuleId-1": {
						AuthSessAmbr: &models.Ambr{
							Uplink:   "5000 Kbps",
							Downlink: "5000 Kbps",
						},
						AuthDefQos: &models.AuthorizedDefaultQos{
							Var5qi: 9,
							Arp: &models.Arp{
								PriorityLevel: 8,
							},
							PriorityLevel: 8,
						},
						SessRuleId: "SessRuleId-1",
					},
				},
			},
			expectedSessRules: map[string]*smf_context.SessionRule{
				"SessRuleId-1": {
					SessionRule: &models.SessionRule{
						AuthSessAmbr: &models.Ambr{
							Uplink:   "5000 Kbps",
							Downlink: "5000 Kbps",
						},
						AuthDefQos: &models.AuthorizedDefaultQos{
							Var5qi: 9,
							Arp: &models.Arp{
								PriorityLevel: 8,
							},
							PriorityLevel: 8,
						},
						SessRuleId: "SessRuleId-1",
					},
					DefQosQFI: 1,
				},
				"SessRuleId-2": {
					SessionRule: &models.SessionRule{
						AuthSessAmbr: &models.Ambr{
							Uplink:   "1000 Kbps",
							Downlink: "1000 Kbps",
						},
						AuthDefQos: &models.AuthorizedDefaultQos{
							Var5qi: 9,
							Arp: &models.Arp{
								PriorityLevel: 8,
							},
							PriorityLevel: 8,
						},
						SessRuleId: "SessRuleId-2",
					},
					DefQosQFI: 1,
				},
			},
			noErr: true,
		},
		{
			name: "Delete first session rule",
			decision: &models.SmPolicyDecision{
				SessRules: map[string]*models.SessionRule{
					"SessRuleId-1": nil,
				},
			},
			expectedSessRules: map[string]*smf_context.SessionRule{
				"SessRuleId-2": {
					SessionRule: &models.SessionRule{
						AuthSessAmbr: &models.Ambr{
							Uplink:   "1000 Kbps",
							Downlink: "1000 Kbps",
						},
						AuthDefQos: &models.AuthorizedDefaultQos{
							Var5qi: 9,
							Arp: &models.Arp{
								PriorityLevel: 8,
							},
							PriorityLevel: 8,
						},
						SessRuleId: "SessRuleId-2",
					},
					DefQosQFI: 1,
				},
			},
			noErr: true,
		},
		{
			name: "Delete second session rule",
			decision: &models.SmPolicyDecision{
				SessRules: map[string]*models.SessionRule{
					"SessRuleId-2": nil,
				},
			},
			expectedSessRules: nil,
			noErr:             false,
		},
	}

	smctx := smf_context.NewSMContext("imsi-208930000000001", 10)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := smctx.ApplySessionRules(tc.decision)
			if tc.noErr {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
			}
			if tc.expectedSessRules != nil {
				require.Equal(t, tc.expectedSessRules, smctx.SessionRules)
			}
		})
	}
}

func TestApplyPccRules(t *testing.T) {
	testCases := []struct {
		name             string
		decision         *models.SmPolicyDecision
		noErr            bool
		expectedPCCRules map[string]*smf_context.PCCRule
		expectedQosDatas map[string]*models.QosData
		expectedTcDatas  map[string]*smf_context.TrafficControlData
	}{
		{
			name:  "nil decision",
			noErr: false,
		},
		{
			name: "Install first pcc rule",
			decision: &models.SmPolicyDecision{
				PccRules: map[string]*models.PccRule{
					"PccRuleId-1": {
						FlowInfos: []models.FlowInformation{
							{
								FlowDescription: "permit out ip from 192.168.0.21 to 10.60.0.0/16",
							},
						},
						PccRuleId:  "PccRuleId-1",
						Precedence: 23,
						RefQosData: []string{"QosId-1"},
						RefTcData:  []string{"TcId-1"},
					},
				},
				QosDecs: map[string]*models.QosData{
					"QosId-1": {
						QosId: "QosId-1",
					},
				},
				TraffContDecs: map[string]*models.TrafficControlData{
					"TcId-1": {
						TcId: "TcId-1",
						RouteToLocs: []*models.RouteToLocation{
							{
								Dnai: "mec",
							},
						},
					},
				},
			},
			expectedPCCRules: map[string]*smf_context.PCCRule{
				"PccRuleId-1": {
					PccRule: &models.PccRule{
						FlowInfos: []models.FlowInformation{
							{
								FlowDescription: "permit out ip from 192.168.0.21 to 10.60.0.0/16",
							},
						},
						PccRuleId:  "PccRuleId-1",
						Precedence: 23,
						RefQosData: []string{"QosId-1"},
						RefTcData:  []string{"TcId-1"},
					},
				},
			},
			expectedQosDatas: map[string]*models.QosData{
				"QosId-1": {
					QosId: "QosId-1",
				},
			},
			expectedTcDatas: map[string]*smf_context.TrafficControlData{
				"TcId-1": {
					TrafficControlData: &models.TrafficControlData{
						TcId: "TcId-1",
						RouteToLocs: []*models.RouteToLocation{
							{
								Dnai: "mec",
							},
						},
					},
				},
			},
			noErr: true,
		},
		{
			name: "Install second pcc rule",
			decision: &models.SmPolicyDecision{
				PccRules: map[string]*models.PccRule{
					"PccRuleId-2": {
						FlowInfos: []models.FlowInformation{
							{
								FlowDescription: "permit out ip from 192.168.0.31 to 10.60.0.0/16",
							},
						},
						PccRuleId:  "PccRuleId-2",
						Precedence: 23,
						RefQosData: []string{"QosId-2"},
						RefTcData:  []string{"TcId-1"},
					},
				},
				QosDecs: map[string]*models.QosData{
					"QosId-2": {
						QosId: "QosId-2",
					},
				},
			},
			expectedPCCRules: map[string]*smf_context.PCCRule{
				"PccRuleId-1": {
					PccRule: &models.PccRule{
						FlowInfos: []models.FlowInformation{
							{
								FlowDescription: "permit out ip from 192.168.0.21 to 10.60.0.0/16",
							},
						},
						PccRuleId:  "PccRuleId-1",
						Precedence: 23,
						RefQosData: []string{"QosId-1"},
						RefTcData:  []string{"TcId-1"},
					},
				},
				"PccRuleId-2": {
					PccRule: &models.PccRule{
						FlowInfos: []models.FlowInformation{
							{
								FlowDescription: "permit out ip from 192.168.0.31 to 10.60.0.0/16",
							},
						},
						PccRuleId:  "PccRuleId-2",
						Precedence: 23,
						RefQosData: []string{"QosId-2"},
						RefTcData:  []string{"TcId-1"},
					},
				},
			},
			expectedQosDatas: map[string]*models.QosData{
				"QosId-1": {
					QosId: "QosId-1",
				},
				"QosId-2": {
					QosId: "QosId-2",
				},
			},
			expectedTcDatas: map[string]*smf_context.TrafficControlData{
				"TcId-1": {
					TrafficControlData: &models.TrafficControlData{
						TcId: "TcId-1",
						RouteToLocs: []*models.RouteToLocation{
							{
								Dnai: "mec",
							},
						},
					},
				},
			},
			noErr: true,
		},
		{
			name: "modify first pcc rule",
			decision: &models.SmPolicyDecision{
				PccRules: map[string]*models.PccRule{
					"PccRuleId-1": {
						FlowInfos: []models.FlowInformation{
							{
								FlowDescription: "permit out ip from 192.168.0.22 to 10.60.0.0/16",
							},
						},
						PccRuleId:  "PccRuleId-1",
						Precedence: 23,
						RefQosData: []string{"QosId-3"},
						RefTcData:  []string{"TcId-1"},
					},
				},
				QosDecs: map[string]*models.QosData{
					"QosId-3": {
						QosId: "QosId-3",
					},
				},
			},
			expectedPCCRules: map[string]*smf_context.PCCRule{
				"PccRuleId-1": {
					PccRule: &models.PccRule{
						FlowInfos: []models.FlowInformation{
							{
								FlowDescription: "permit out ip from 192.168.0.22 to 10.60.0.0/16",
							},
						},
						PccRuleId:  "PccRuleId-1",
						Precedence: 23,
						RefQosData: []string{"QosId-3"},
						RefTcData:  []string{"TcId-1"},
					},
				},
				"PccRuleId-2": {
					PccRule: &models.PccRule{
						FlowInfos: []models.FlowInformation{
							{
								FlowDescription: "permit out ip from 192.168.0.31 to 10.60.0.0/16",
							},
						},
						PccRuleId:  "PccRuleId-2",
						Precedence: 23,
						RefQosData: []string{"QosId-2"},
						RefTcData:  []string{"TcId-1"},
					},
				},
			},
			expectedQosDatas: map[string]*models.QosData{
				"QosId-2": {
					QosId: "QosId-2",
				},
				"QosId-3": {
					QosId: "QosId-3",
				},
			},
			expectedTcDatas: map[string]*smf_context.TrafficControlData{
				"TcId-1": {
					TrafficControlData: &models.TrafficControlData{
						TcId: "TcId-1",
						RouteToLocs: []*models.RouteToLocation{
							{
								Dnai: "mec",
							},
						},
					},
				},
			},
			noErr: true,
		},
		{
			name: "delete second pcc rule",
			decision: &models.SmPolicyDecision{
				PccRules: map[string]*models.PccRule{
					"PccRuleId-2": nil,
				},
			},
			expectedPCCRules: map[string]*smf_context.PCCRule{
				"PccRuleId-1": {
					PccRule: &models.PccRule{
						FlowInfos: []models.FlowInformation{
							{
								FlowDescription: "permit out ip from 192.168.0.22 to 10.60.0.0/16",
							},
						},
						PccRuleId:  "PccRuleId-1",
						Precedence: 23,
						RefQosData: []string{"QosId-3"},
						RefTcData:  []string{"TcId-1"},
					},
				},
			},
			expectedQosDatas: map[string]*models.QosData{
				"QosId-3": {
					QosId: "QosId-3",
				},
			},
			expectedTcDatas: map[string]*smf_context.TrafficControlData{
				"TcId-1": {
					TrafficControlData: &models.TrafficControlData{
						TcId: "TcId-1",
						RouteToLocs: []*models.RouteToLocation{
							{
								Dnai: "mec",
							},
						},
					},
				},
			},
			noErr: true,
		},
		{
			name: "delete first pcc rule",
			decision: &models.SmPolicyDecision{
				PccRules: map[string]*models.PccRule{
					"PccRuleId-1": nil,
				},
			},
			expectedPCCRules: map[string]*smf_context.PCCRule{},
			expectedQosDatas: map[string]*models.QosData{},
			expectedTcDatas:  map[string]*smf_context.TrafficControlData{},
			noErr:            true,
		},
	}

	smfContext := smf_context.GetSelf()
	smfContext.UserPlaneInformation = smf_context.NewUserPlaneInformation(&userPlaneConfig)
	for _, n := range smfContext.UserPlaneInformation.UPFs {
		n.UPF.AssociationContext = context.Background()
	}

	smctx := smf_context.NewSMContext("imsi-208930000000002", 10)

	smctx.SMLock.Lock()
	defer smctx.SMLock.Unlock()

	smctx.SmfPduSessionSmContextCreateData = &models.SmfPduSessionSmContextCreateData{
		Supi:         "imsi-208930000000002",
		Pei:          "imeisv-1110000000000000",
		Gpsi:         "msisdn-0900000000",
		PduSessionId: 10,
		Dnn:          "internet",
		SNssai: &models.Snssai{
			Sst: 1,
			Sd:  "010203",
		},
		ServingNfId: "c8d0ee65-f466-48aa-a42f-235ec771cb52",
		Guami: &models.Guami{
			PlmnId: &models.PlmnIdNid{
				Mcc: "208",
				Mnc: "93",
			},
			AmfId: "cafe00",
		},
		AnType: "3GPP_ACCESS",
		ServingNetwork: &models.PlmnIdNid{
			Mcc: "208",
			Mnc: "93",
		},
	}
	smctx.SelectedPDUSessionType = 1
	smctx.SessionRules["SessRuleId-1"] = &smf_context.SessionRule{
		SessionRule: &models.SessionRule{
			AuthSessAmbr: &models.Ambr{
				Uplink:   "1000 Kbps",
				Downlink: "1000 Kbps",
			},
			AuthDefQos: &models.AuthorizedDefaultQos{
				Var5qi: 9,
				Arp: &models.Arp{
					PriorityLevel: 8,
				},
				PriorityLevel: 8,
			},
			SessRuleId: "SessRuleId-1",
		},
	}
	smctx.SelectedSessionRuleID = "SessRuleId-1"
	smctx.DnnConfiguration = models.DnnConfiguration{
		SessionAmbr: &models.Ambr{
			Uplink:   "1 Gbps",
			Downlink: "1 Gbps",
		},
	}
	err := smctx.AllocUeIP()
	require.NoError(t, err)
	err = smctx.SelectDefaultDataPath()
	require.NoError(t, err)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err = smctx.ApplyPccRules(tc.decision)
			if tc.noErr {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
			}
			for id, expPcc := range tc.expectedPCCRules {
				require.Equal(t, expPcc.PccRule, smctx.PCCRules[id].PccRule)
			}
			if tc.expectedQosDatas != nil {
				require.Equal(t, tc.expectedQosDatas, smctx.QosDatas)
			}
			if tc.expectedTcDatas != nil {
				require.Equal(t, tc.expectedTcDatas, smctx.TrafficControlDatas)
			}
		})
	}
}
