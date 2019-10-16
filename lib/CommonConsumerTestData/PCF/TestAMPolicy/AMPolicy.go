package TestAMPolicy

import (
	"free5gc/lib/openapi/models"
	"time"
)

//Data type definition TS 29.571//
/*
Supi : IMSI 46611123456789
SuppFeat : ()
Pei : IMEI 123456789123456 (length 15)
RatType : NR ,WLAN ,EUTRA,VIRTUAL
Rfsp : 1~256
AccessType : "3GPP_ACCESS" , "NON_3GPP_ACCESS"
ServAreaRes :
	RestrictionType : "ALLOWED_AREAS","NOT_ALLOWED_AREAS"
	Tacs : "0000" & "FFFF" /reserve
	AreaCodes : "+886"
	MaxNumOfTAs :  The "maxNumOfTAs" attribute value cannot be lower than the number of TAIs included in the "tacs" attribute.
UserLoc :
	EutraLocation :
		Tai : Tracking  Area  Identity
			PlmnId :
				Mcc : Mobile Country Code (466 Taiwan)
				Mnc : Mobile Network Code (11 CHT)
			Tac : Tracking area code
		Ecgi : E-UTRAN Cell Global Identifier
			PlmnId :
				Mcc : Mobile Country Code (466 Taiwan)
				Mnc : Mobile Network Code (11 CHT)
			EutraCellId : E-UTRA Cell Identity
		AgeOfLocationInformation : 0
		GeographicalInformation : Allowed characters are 0-9 and A-F, See TS23.032 7.3.2
		GeodeticInformation : Allowed characters are 0-9 and A-F, See ITU-T Recommendation Q.763 3.88.2
		GlobalNgenbId : Global identity of the ng-eNodeB in which the UE is currently located, See TS38.413 9.3.1.8
			PlmnId :
				Mnc : Mobile Country Code (466 Taiwan)
				Mcc : Mobile Network Code (11 CHT)
			N3IwfId : This IE shall be included if the RAN node belongs to non 3GPP access
			GNbId :
				BitLength : 22 to 32, TS38.413
				GNBValue : Examples:
							A 30 bit value "382A3F47" indicates a gNB ID with value 0x382A3F47
							A 22 bit value "2A3F47" indicates a gNB ID with value 0x2A3F47
			NgeNbId : Examples:
						" SMacroNGeNB-34B89" indicates a Short Macro NG-eNB ID with value 0x34B89.
*/
//Success Test
func GetAMreqdata() models.PolicyAssociationRequest {
	//d := time.Date(2019, 7, 5, 12, 30, 0, 0, time.UTC)
	amCreateReqData := models.PolicyAssociationRequest{
		NotificationUri: "abc@gmail.com",
		Supi:            "46611123456789",
		SuppFeat:        "1",
		Pei:             "123456789123456",
		RatType:         "NR",
		Rfsp:            123,
		AccessType:      "3GPP_ACCESS",
		ServAreaRes: models.ServiceAreaRestriction{
			RestrictionType: "ALLOWED_AREAS",
			Areas: []models.Area{
				{
					Tacs: []string{
						"0000",
					},
					AreaCodes: "+886",
				},
			},
			MaxNumOfTAs: 999,
		},
		UserLoc: models.UserLocation{
			EutraLocation: &models.EutraLocation{
				Tai: &models.Tai{
					PlmnId: &models.PlmnId{
						Mcc: "466",
						Mnc: "11",
					},
					Tac: "0000",
				},
				Ecgi: &models.Ecgi{
					PlmnId: &models.PlmnId{
						Mcc: "466",
						Mnc: "11",
					},
					EutraCellId: "string",
				},
				AgeOfLocationInformation: 0,
				//UeLocationTimestamp:      d,
				// UeLocationTimestamp: &time.Time{
				// 	,
				// },
				GeographicalInformation: "0A9F",
				GeodeticInformation:     "0A9F",
				GlobalNgenbId: &models.GlobalRanNodeId{
					PlmnId: &models.PlmnId{
						Mcc: "466",
						Mnc: "11",
					},
					N3IwfId: "string",
					GNbId: &models.GNbId{
						BitLength: 123,
						GNBValue:  "string",
					},
					NgeNbId: "string",
				},
			},
			NrLocation: &models.NrLocation{
				Tai: &models.Tai{
					PlmnId: &models.PlmnId{
						Mcc: "466",
						Mnc: "11",
					},
					Tac: "string",
				},
				Ncgi: &models.Ncgi{
					PlmnId: &models.PlmnId{
						Mcc: "466",
						Mnc: "11",
					},
					NrCellId: "string",
				},
				AgeOfLocationInformation: 123,
				//UeLocationTimestamp:      "2019-04-15T09:47:35.505Z",
				UeLocationTimestamp:     &time.Time{},
				GeographicalInformation: "string",
				GeodeticInformation:     "string",
				GlobalGnbId: &models.GlobalRanNodeId{
					PlmnId: &models.PlmnId{
						Mcc: "466",
						Mnc: "11",
					},
					N3IwfId: "string",
					GNbId: &models.GNbId{
						BitLength: 123,
						GNBValue:  "string",
					},
					NgeNbId: "string",
				},
			},
			N3gaLocation: &models.N3gaLocation{
				N3gppTai: &models.Tai{
					PlmnId: &models.PlmnId{
						Mcc: "466",
						Mnc: "11",
					},
					Tac: "string",
				},
				N3IwfId:    "string",
				UeIpv4Addr: "198.51.100.1",
				UeIpv6Addr: "2001:db8:85a3::8a2e:370:7334",
				PortNumber: 123,
			},
		},
		ServingPlmn: models.NetworkId{
			Mcc: "466",
			Mnc: "11",
		},
		GroupIds: []string{
			"groupids",
		},
		Guami: models.Guami{
			PlmnId: &models.PlmnId{
				Mcc: "466",
				Mnc: "11",
			},
			AmfId: "string",
		},
		TimeZone: "string",
		TraceReq: &models.TraceData{
			TraceDepth:               "string",
			TraceRef:                 "string",
			NeTypeList:               "string",
			EventList:                "string",
			CollectionEntityIpv4Addr: "198.51.100.1",
			CollectionEntityIpv6Addr: "2001:db8:85a3::8a2e:370:7334",
			InterfaceList:            "string",
		},
	}
	return amCreateReqData
}
func GetAMUpdateReqData() models.PolicyAssociationUpdateRequest {
	amUpdateReqData := models.PolicyAssociationUpdateRequest{
		AllowedSnssais: []models.Snssai{
			{
				Sst: 123,
				Sd:  "string",
			},
		},
		Triggers: []models.RequestTrigger{
			"LOC_CH",
			"PRA_CH",
			"RSFP_CH",
			"SERV_AREA_CH",
		},
		ServAreaRes: &models.ServiceAreaRestriction{
			RestrictionType: "string",
			Areas: []models.Area{
				{
					Tacs: []string{
						"tacs",
					},
					AreaCodes: "123",
				},
			},
			MaxNumOfTAs: 123,
		},
		Rfsp: 123,
		TraceReq: &models.TraceData{
			TraceDepth:               "string",
			TraceRef:                 "string",
			NeTypeList:               "string",
			EventList:                "string",
			CollectionEntityIpv4Addr: "198.51.100.1",
			CollectionEntityIpv6Addr: "2001:db8:85a3::8a2e:370:7334",
			InterfaceList:            "string",
		},
		UserLoc: &models.UserLocation{
			EutraLocation: &models.EutraLocation{
				Tai: &models.Tai{
					PlmnId: &models.PlmnId{
						Mcc: "466",
						Mnc: "11",
					},
					Tac: "123",
				},
				Ecgi: &models.Ecgi{
					PlmnId: &models.PlmnId{
						Mcc: "466",
						Mnc: "11",
					},
					EutraCellId: "string",
				},
				AgeOfLocationInformation: 123,
				//UeLocationTimestamp:      d,
				// UeLocationTimestamp: &time.Time{
				// 	,
				// },
				GeographicalInformation: "string",
				GeodeticInformation:     "string",
				GlobalNgenbId: &models.GlobalRanNodeId{
					PlmnId: &models.PlmnId{
						Mcc: "466",
						Mnc: "11",
					},
					N3IwfId: "string",
					GNbId: &models.GNbId{
						BitLength: 123,
						GNBValue:  "string",
					},
					NgeNbId: "string",
				},
			},
			NrLocation: &models.NrLocation{
				Tai: &models.Tai{
					PlmnId: &models.PlmnId{
						Mcc: "466",
						Mnc: "11",
					},
					Tac: "string",
				},
				Ncgi: &models.Ncgi{
					PlmnId: &models.PlmnId{
						Mcc: "466",
						Mnc: "11",
					},
					NrCellId: "string",
				},
				AgeOfLocationInformation: 123,
				//UeLocationTimestamp:      "2019-04-15T09:47:35.505Z",
				UeLocationTimestamp:     &time.Time{},
				GeographicalInformation: "string",
				GeodeticInformation:     "string",
				GlobalGnbId: &models.GlobalRanNodeId{
					PlmnId: &models.PlmnId{
						Mcc: "466",
						Mnc: "11",
					},
					N3IwfId: "string",
					GNbId: &models.GNbId{
						BitLength: 123,
						GNBValue:  "string",
					},
					NgeNbId: "string",
				},
			},
			N3gaLocation: &models.N3gaLocation{
				N3gppTai: &models.Tai{
					PlmnId: &models.PlmnId{
						Mcc: "466",
						Mnc: "11",
					},
					Tac: "string",
				},
				N3IwfId:    "string",
				UeIpv4Addr: "198.51.100.1",
				UeIpv6Addr: "2001:db8:85a3::8a2e:370:7334",
				PortNumber: 123,
			},
		},
		PraStatuses: map[string]models.PresenceInfo{
			"Pras1": {PraId: "11111",
				PresenceState: "IN_AREA",
				TrackingAreaList: []models.Tai{
					{Tac: "string",
						PlmnId: &models.PlmnId{
							Mcc: "string",
							Mnc: "string",
						},
					},
				},
				EcgiList: []models.Ecgi{
					{EutraCellId: "string",
						PlmnId: &models.PlmnId{
							Mcc: "string",
							Mnc: "string",
						},
					},
				},
				NcgiList: []models.Ncgi{
					{NrCellId: "string",
						PlmnId: &models.PlmnId{
							Mcc: "string",
							Mnc: "string",
						},
					},
				},
				GlobalRanNodeIdList: []models.GlobalRanNodeId{
					{
						PlmnId: &models.PlmnId{
							Mcc: "string",
							Mnc: "string",
						},
						N3IwfId: "string",
						GNbId: &models.GNbId{
							BitLength: 123,
							GNBValue:  "string",
						},
						NgeNbId: "string",
					},
				},
			},
		},
		NotificationUri:   "string",
		AltNotifIpv4Addrs: []string{"ipv4addr"},
		AltNotifIpv6Addrs: []string{"ipv6addr"},
	}
	return amUpdateReqData
}

//-------------------------------------------------------------------------------------------------

//-------------------------------------------------------------------------------------------------
//Fail Test (Create part)
func GetamCreatefailnotifyURIData() models.PolicyAssociationRequest {
	//d := time.Date(2019, 7, 5, 12, 30, 0, 0, time.UTC)
	amCreatefailnotifyURIData := models.PolicyAssociationRequest{
		NotificationUri: "",
		Supi:            "46611123456789",
		SuppFeat:        "string",
		Pei:             "string",
		RatType:         "string",
		Rfsp:            123,
		AccessType:      "3GPP_ACCESS",
		ServAreaRes: models.ServiceAreaRestriction{
			RestrictionType: "string",
			Areas: []models.Area{
				{
					Tacs: []string{
						"tacs",
					},
					AreaCodes: "123",
				},
			},
			MaxNumOfTAs: 123,
		},
		UserLoc: models.UserLocation{
			EutraLocation: &models.EutraLocation{
				Tai: &models.Tai{
					PlmnId: &models.PlmnId{
						Mcc: "466",
						Mnc: "11",
					},
					Tac: "string",
				},
				Ecgi: &models.Ecgi{
					PlmnId: &models.PlmnId{
						Mcc: "466",
						Mnc: "11",
					},
					EutraCellId: "string",
				},
				AgeOfLocationInformation: 123,
				//UeLocationTimestamp:      d,
				// UeLocationTimestamp: &time.Time{
				// 	,
				// },
				GeographicalInformation: "string",
				GeodeticInformation:     "string",
				GlobalNgenbId: &models.GlobalRanNodeId{
					PlmnId: &models.PlmnId{
						Mcc: "466",
						Mnc: "11",
					},
					N3IwfId: "string",
					GNbId: &models.GNbId{
						BitLength: 123,
						GNBValue:  "string",
					},
					NgeNbId: "string",
				},
			},
			NrLocation: &models.NrLocation{
				Tai: &models.Tai{
					PlmnId: &models.PlmnId{
						Mcc: "466",
						Mnc: "11",
					},
					Tac: "string",
				},
				Ncgi: &models.Ncgi{
					PlmnId: &models.PlmnId{
						Mcc: "466",
						Mnc: "11",
					},
					NrCellId: "string",
				},
				AgeOfLocationInformation: 123,
				//UeLocationTimestamp:      "2019-04-15T09:47:35.505Z",
				UeLocationTimestamp:     &time.Time{},
				GeographicalInformation: "string",
				GeodeticInformation:     "string",
				GlobalGnbId: &models.GlobalRanNodeId{
					PlmnId: &models.PlmnId{
						Mcc: "466",
						Mnc: "11",
					},
					N3IwfId: "string",
					GNbId: &models.GNbId{
						BitLength: 123,
						GNBValue:  "string",
					},
					NgeNbId: "string",
				},
			},
			N3gaLocation: &models.N3gaLocation{
				N3gppTai: &models.Tai{
					PlmnId: &models.PlmnId{
						Mcc: "466",
						Mnc: "11",
					},
					Tac: "string",
				},
				N3IwfId:    "string",
				UeIpv4Addr: "198.51.100.1",
				UeIpv6Addr: "2001:db8:85a3::8a2e:370:7334",
				PortNumber: 123,
			},
		},
		ServingPlmn: models.NetworkId{
			Mcc: "466",
			Mnc: "11",
		},
		GroupIds: []string{
			"groupids",
		},
		Guami: models.Guami{
			PlmnId: &models.PlmnId{
				Mcc: "466",
				Mnc: "11",
			},
			AmfId: "string",
		},
		TimeZone: "string",
		TraceReq: &models.TraceData{
			TraceDepth:               "string",
			TraceRef:                 "string",
			NeTypeList:               "string",
			EventList:                "string",
			CollectionEntityIpv4Addr: "198.51.100.1",
			CollectionEntityIpv6Addr: "2001:db8:85a3::8a2e:370:7334",
			InterfaceList:            "string",
		},
	}
	return amCreatefailnotifyURIData
}
func GetamCreatefailsupiData() models.PolicyAssociationRequest {
	//d := time.Date(2019, 7, 5, 12, 30, 0, 0, time.UTC)
	amCreatefailsupiData := models.PolicyAssociationRequest{
		NotificationUri: "string",
		Supi:            "",
		SuppFeat:        "string",
		Pei:             "string",
		RatType:         "string",
		Rfsp:            123,
		AccessType:      "3GPP_ACCESS",
		ServAreaRes: models.ServiceAreaRestriction{
			RestrictionType: "string",
			Areas: []models.Area{
				{
					Tacs: []string{
						"tacs",
					},
					AreaCodes: "123",
				},
			},
			MaxNumOfTAs: 123,
		},
		UserLoc: models.UserLocation{
			EutraLocation: &models.EutraLocation{
				Tai: &models.Tai{
					PlmnId: &models.PlmnId{
						Mcc: "466",
						Mnc: "11",
					},
					Tac: "string",
				},
				Ecgi: &models.Ecgi{
					PlmnId: &models.PlmnId{
						Mcc: "466",
						Mnc: "11",
					},
					EutraCellId: "string",
				},
				AgeOfLocationInformation: 123,
				//UeLocationTimestamp:      d,
				// UeLocationTimestamp: &time.Time{
				// 	,
				// },
				GeographicalInformation: "string",
				GeodeticInformation:     "string",
				GlobalNgenbId: &models.GlobalRanNodeId{
					PlmnId: &models.PlmnId{
						Mcc: "466",
						Mnc: "11",
					},
					N3IwfId: "string",
					GNbId: &models.GNbId{
						BitLength: 123,
						GNBValue:  "string",
					},
					NgeNbId: "string",
				},
			},
			NrLocation: &models.NrLocation{
				Tai: &models.Tai{
					PlmnId: &models.PlmnId{
						Mcc: "466",
						Mnc: "11",
					},
					Tac: "string",
				},
				Ncgi: &models.Ncgi{
					PlmnId: &models.PlmnId{
						Mcc: "466",
						Mnc: "11",
					},
					NrCellId: "string",
				},
				AgeOfLocationInformation: 123,
				//UeLocationTimestamp:      "2019-04-15T09:47:35.505Z",
				UeLocationTimestamp:     &time.Time{},
				GeographicalInformation: "string",
				GeodeticInformation:     "string",
				GlobalGnbId: &models.GlobalRanNodeId{
					PlmnId: &models.PlmnId{
						Mcc: "466",
						Mnc: "11",
					},
					N3IwfId: "string",
					GNbId: &models.GNbId{
						BitLength: 123,
						GNBValue:  "string",
					},
					NgeNbId: "string",
				},
			},
			N3gaLocation: &models.N3gaLocation{
				N3gppTai: &models.Tai{
					PlmnId: &models.PlmnId{
						Mcc: "466",
						Mnc: "11",
					},
					Tac: "string",
				},
				N3IwfId:    "string",
				UeIpv4Addr: "198.51.100.1",
				UeIpv6Addr: "2001:db8:85a3::8a2e:370:7334",
				PortNumber: 123,
			},
		},
		ServingPlmn: models.NetworkId{
			Mcc: "466",
			Mnc: "11",
		},
		GroupIds: []string{
			"groupids",
		},
		Guami: models.Guami{
			PlmnId: &models.PlmnId{
				Mcc: "466",
				Mnc: "11",
			},
			AmfId: "string",
		},
		TimeZone: "string",
		TraceReq: &models.TraceData{
			TraceDepth:               "string",
			TraceRef:                 "string",
			NeTypeList:               "string",
			EventList:                "string",
			CollectionEntityIpv4Addr: "198.51.100.1",
			CollectionEntityIpv6Addr: "2001:db8:85a3::8a2e:370:7334",
			InterfaceList:            "string",
		},
	}
	return amCreatefailsupiData
}
func GetamCreatefailsuppfeatData() models.PolicyAssociationRequest {
	//d := time.Date(2019, 7, 5, 12, 30, 0, 0, time.UTC)
	amCreatefailsuppfeatData := models.PolicyAssociationRequest{
		NotificationUri: "string",
		Supi:            "46611123456789",
		SuppFeat:        "",
		Pei:             "string",
		RatType:         "string",
		Rfsp:            123,
		AccessType:      "3GPP_ACCESS",
		ServAreaRes: models.ServiceAreaRestriction{
			RestrictionType: "string",
			Areas: []models.Area{
				{
					Tacs: []string{
						"tacs",
					},
					AreaCodes: "123",
				},
			},
			MaxNumOfTAs: 123,
		},
		UserLoc: models.UserLocation{
			EutraLocation: &models.EutraLocation{
				Tai: &models.Tai{
					PlmnId: &models.PlmnId{
						Mcc: "466",
						Mnc: "11",
					},
					Tac: "string",
				},
				Ecgi: &models.Ecgi{
					PlmnId: &models.PlmnId{
						Mcc: "466",
						Mnc: "11",
					},
					EutraCellId: "string",
				},
				AgeOfLocationInformation: 123,
				//UeLocationTimestamp:      d,
				// UeLocationTimestamp: &time.Time{
				// 	,
				// },
				GeographicalInformation: "string",
				GeodeticInformation:     "string",
				GlobalNgenbId: &models.GlobalRanNodeId{
					PlmnId: &models.PlmnId{
						Mcc: "466",
						Mnc: "11",
					},
					N3IwfId: "string",
					GNbId: &models.GNbId{
						BitLength: 123,
						GNBValue:  "string",
					},
					NgeNbId: "string",
				},
			},
			NrLocation: &models.NrLocation{
				Tai: &models.Tai{
					PlmnId: &models.PlmnId{
						Mcc: "466",
						Mnc: "11",
					},
					Tac: "string",
				},
				Ncgi: &models.Ncgi{
					PlmnId: &models.PlmnId{
						Mcc: "466",
						Mnc: "11",
					},
					NrCellId: "string",
				},
				AgeOfLocationInformation: 123,
				//UeLocationTimestamp:      "2019-04-15T09:47:35.505Z",
				UeLocationTimestamp:     &time.Time{},
				GeographicalInformation: "string",
				GeodeticInformation:     "string",
				GlobalGnbId: &models.GlobalRanNodeId{
					PlmnId: &models.PlmnId{
						Mcc: "466",
						Mnc: "11",
					},
					N3IwfId: "string",
					GNbId: &models.GNbId{
						BitLength: 123,
						GNBValue:  "string",
					},
					NgeNbId: "string",
				},
			},
			N3gaLocation: &models.N3gaLocation{
				N3gppTai: &models.Tai{
					PlmnId: &models.PlmnId{
						Mcc: "466",
						Mnc: "11",
					},
					Tac: "string",
				},
				N3IwfId:    "string",
				UeIpv4Addr: "198.51.100.1",
				UeIpv6Addr: "2001:db8:85a3::8a2e:370:7334",
				PortNumber: 123,
			},
		},
		ServingPlmn: models.NetworkId{
			Mcc: "466",
			Mnc: "11",
		},
		GroupIds: []string{
			"groupids",
		},
		Guami: models.Guami{
			PlmnId: &models.PlmnId{
				Mcc: "466",
				Mnc: "11",
			},
			AmfId: "string",
		},
		TimeZone: "string",
		TraceReq: &models.TraceData{
			TraceDepth:               "string",
			TraceRef:                 "string",
			NeTypeList:               "string",
			EventList:                "string",
			CollectionEntityIpv4Addr: "198.51.100.1",
			CollectionEntityIpv6Addr: "2001:db8:85a3::8a2e:370:7334",
			InterfaceList:            "string",
		},
	}
	return amCreatefailsuppfeatData
}
