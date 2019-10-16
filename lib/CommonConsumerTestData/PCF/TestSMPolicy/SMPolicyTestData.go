package TestSMPolicy

import (
	"time"

	"free5gc/lib/openapi/models"
)

func CreateTestData() models.SmPolicyContextData {
	t := time.Date(2000, 2, 1, 12, 30, 0, 0, time.UTC)
	smReqData := models.SmPolicyContextData{
		AccNetChId: &models.AccNetChId{
			AccNetChaIdValue: 0,
			RefPccRuleIds:    []string{"A", "B", "C"},
			SessionChScope:   true,
		},
		ChargEntityAddr: &models.AccNetChargingAddress{
			AnChargIpv4Addr: "198.51.100.1",
			AnChargIpv6Addr: "2001:db8:85a3::8a2e:370:7334",
		},
		Gpsi:                    "string1",
		Supi:                    "string1",
		InterGrpIds:             []string{"A", "B", "C"},
		PduSessionId:            123,
		Chargingcharacteristics: "string",
		Dnn:                     "string",
		NotificationUri:         "https://localhost:8081/nsmf/NotificationUri",
		AccessType:              "3GPP_ACCESS",
		PduSessionType:          "string",
		ServingNetwork: &models.NetworkId{
			Mnc: "string",
			Mcc: "string",
		},
		UserLocationInfo: &models.UserLocation{
			EutraLocation: &models.EutraLocation{
				Tai: &models.Tai{
					PlmnId: &models.PlmnId{
						Mcc: "string",
						Mnc: "string",
					},
					Tac: "string",
				},
				Ecgi: &models.Ecgi{
					PlmnId: &models.PlmnId{
						Mcc: "string",
						Mnc: "string",
					},
					EutraCellId: "string",
				},
				AgeOfLocationInformation: 0,
				UeLocationTimestamp:      &t,
				GeographicalInformation:  "string",
				GeodeticInformation:      "string",
				GlobalNgenbId: &models.GlobalRanNodeId{
					PlmnId: &models.PlmnId{
						Mcc: "string",
						Mnc: "string",
					},
					N3IwfId: "string",
					GNbId: &models.GNbId{
						BitLength: 0,
						GNBValue:  "string",
					},
					NgeNbId: "string",
				},
			},
			NrLocation: &models.NrLocation{
				Tai: &models.Tai{
					PlmnId: &models.PlmnId{
						Mcc: "string",
						Mnc: "string",
					},
					Tac: "string",
				},
				Ncgi: &models.Ncgi{
					PlmnId: &models.PlmnId{
						Mcc: "string",
						Mnc: "string",
					},
					NrCellId: "string",
				},
				AgeOfLocationInformation: 0,
				UeLocationTimestamp:      &t,
				GeographicalInformation:  "string",
				GeodeticInformation:      "string",
				GlobalGnbId: &models.GlobalRanNodeId{
					PlmnId: &models.PlmnId{
						Mcc: "string",
						Mnc: "string",
					},
					N3IwfId: "string",
					GNbId: &models.GNbId{
						BitLength: 0,
						GNBValue:  "string",
					},
					NgeNbId: "string",
				},
			},
			N3gaLocation: &models.N3gaLocation{
				N3gppTai: &models.Tai{
					PlmnId: &models.PlmnId{
						Mcc: "string",
						Mnc: "string",
					},
					Tac: "string",
				},
				N3IwfId:    "string",
				UeIpv4Addr: "198.51.100.1",
				UeIpv6Addr: "2001:db8:85a3::8a2e:370:7334",
				PortNumber: 0,
			},
		},
		UeTimeZone:        "string",
		Pei:               "string",
		Ipv4Address:       "198.51.100.1",
		Ipv6AddressPrefix: "2001:db8:abcd:12::0/64",
		IpDomain:          "string",
		SubsSessAmbr: &models.Ambr{
			Uplink:   "string",
			Downlink: "string",
		},
		SubsDefQos: &models.SubscribedDefaultQos{
			Var5qi: 0,
			Arp: &models.Arp{
				PriorityLevel: 0,
			},
			PriorityLevel: 0,
		},
		NumOfPackFilter:        0,
		Online:                 true,
		Offline:                true,
		Var3gppPsDataOffStatus: true,
		RefQosIndication:       true,
		TraceReq: &models.TraceData{
			TraceRef:                 "string",
			NeTypeList:               "string",
			EventList:                "string",
			CollectionEntityIpv4Addr: "198.51.100.1",
			CollectionEntityIpv6Addr: "2001:db8:85a3::8a2e:370:7334",
			InterfaceList:            "string",
		},
		SliceInfo: &models.Snssai{
			Sst: 1,
			Sd:  "string",
		},
		SuppFeat: "string",
	}
	return smReqData
}

func UpdateTestData() models.SmPolicyUpdateContextData {
	t := time.Date(2000, 2, 1, 12, 30, 0, 0, time.UTC)
	smUpData := models.SmPolicyUpdateContextData{
		RepPolicyCtrlReqTriggers: []models.PolicyControlRequestTrigger{"PLMN_CH", "string", "C"},
		AccNetChIds: []models.AccNetChId{
			{
				AccNetChaIdValue: 0,
				RefPccRuleIds:    []string{"string"},
				SessionChScope:   true,
			},
		},
		AccessType: "3GPP_ACCESS",
		ServingNetwork: &models.NetworkId{
			Mnc: "string",
			Mcc: "string",
		},
		UserLocationInfo: &models.UserLocation{
			EutraLocation: &models.EutraLocation{
				Tai: &models.Tai{
					PlmnId: &models.PlmnId{
						Mcc: "string",
						Mnc: "string",
					},
					Tac: "string",
				},
				Ecgi: &models.Ecgi{
					PlmnId: &models.PlmnId{
						Mcc: "string",
						Mnc: "string",
					},
					EutraCellId: "string",
				},
				AgeOfLocationInformation: 0,
				UeLocationTimestamp:      &t,
				GeographicalInformation:  "string",
				GeodeticInformation:      "string",
				GlobalNgenbId: &models.GlobalRanNodeId{
					PlmnId: &models.PlmnId{
						Mcc: "string",
						Mnc: "string",
					},
					N3IwfId: "string",
					GNbId: &models.GNbId{
						BitLength: 0,
						GNBValue:  "string",
					},
					NgeNbId: "string",
				},
			},
		},
	}
	return smUpData
}
func DeldateTestData() models.SmPolicyDeleteData {
	smDelData := models.SmPolicyDeleteData{
		UserLocationInfo: &models.UserLocation{
			EutraLocation: &models.EutraLocation{
				Tai: &models.Tai{
					PlmnId: &models.PlmnId{
						Mcc: "string",
						Mnc: "string",
					},
				},
			},
		},
	}
	return smDelData
}

func CreateTestData12345() models.SmPolicyContextData {
	smReqDataTest12345 := models.SmPolicyContextData{
		PduSessionId:    12345,
		Dnn:             "string",
		NotificationUri: "string",
		PduSessionType:  "string",
		SliceInfo: &models.Snssai{
			Sst: 1,
			Sd:  "string",
		},
		Supi: "string2",
		Gpsi: "string2",
		ChargEntityAddr: &models.AccNetChargingAddress{
			AnChargIpv4Addr: "198.51.100.1",
			AnChargIpv6Addr: "2001:db8:85a3::8a2e:370:7334",
		},
		Ipv4Address:       "198.51.100.1",
		Ipv6AddressPrefix: "2001:db8:abcd:12::0/64",
	}
	return smReqDataTest12345
}
