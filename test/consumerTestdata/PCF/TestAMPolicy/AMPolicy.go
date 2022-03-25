package TestAMPolicy

import (
	"time"

	"github.com/free5gc/openapi/models"
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
	timeNow := time.Now()
	//d := time.Date(2019, 7, 5, 12, 30, 0, 0, time.UTC)
	amCreateReqData := models.PolicyAssociationRequest{
		NotificationUri: "http://127.0.0.1:29518/namf-callback/v1/am-policy/imsi-2089300007487-1",
		Supi:            "imsi-2089300007487",
		SuppFeat:        "1",
		Pei:             "123456789123456",
		RatType:         models.RatType_NR,
		Rfsp:            123,
		AccessType:      models.AccessType__3_GPP_ACCESS,
		ServAreaRes: &models.ServiceAreaRestriction{
			RestrictionType: "ALLOWED_AREAS",
			Areas: []models.Area{
				{
					Tacs: []string{
						"000001",
					},
					AreaCodes: "+886",
				},
			},
			MaxNumOfTAs: 999,
		},
		UserLoc: &models.UserLocation{
			NrLocation: &models.NrLocation{
				Tai: &models.Tai{
					PlmnId: &models.PlmnId{
						Mcc: "208",
						Mnc: "93",
					},
					Tac: "000001",
				},
				Ncgi: &models.Ncgi{
					PlmnId: &models.PlmnId{
						Mcc: "208",
						Mnc: "93",
					},
					NrCellId: "000000001",
				},
				AgeOfLocationInformation: 123,
				//UeLocationTimestamp:      "2019-04-15T09:47:35.505Z",
				UeLocationTimestamp: &timeNow,
				GlobalGnbId: &models.GlobalRanNodeId{
					PlmnId: &models.PlmnId{
						Mcc: "208",
						Mnc: "93",
					},
					GNbId: &models.GNbId{
						BitLength: 24,
						GNBValue:  "000001",
					},
				},
			},
		},
		ServingPlmn: &models.NetworkId{
			Mcc: "208",
			Mnc: "93",
		},
		GroupIds: []string{
			"groupids",
		},
		Guami: &models.Guami{
			PlmnId: &models.PlmnId{
				Mcc: "208",
				Mnc: "93",
			},
			AmfId: "cafe00",
		},
		TimeZone: "+08:00+1h",
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
		Triggers: []models.RequestTrigger{
			"RFSP_CH",
			"SERV_AREA_CH",
		},
		ServAreaRes: &models.ServiceAreaRestriction{
			RestrictionType: "ALLOWED_AREAS",
			Areas: []models.Area{
				{
					Tacs: []string{
						"000002",
					},
					AreaCodes: "+886",
				},
			},
			MaxNumOfTAs: 123,
		},
		Rfsp: 100,
		TraceReq: &models.TraceData{
			TraceDepth:               "string",
			TraceRef:                 "string",
			NeTypeList:               "string",
			EventList:                "string",
			CollectionEntityIpv4Addr: "198.51.100.1",
			CollectionEntityIpv6Addr: "2001:db8:85a3::8a2e:370:7334",
			InterfaceList:            "string",
		},
		NotificationUri:   "http://127.0.0.1:29518/namf-callback/v1/am-policy-backbup/imsi-2089300007487-1",
		AltNotifIpv4Addrs: []string{"http://127.0.0.1:29518/namf-callback/v1/am-policy/imsi-2089300007487-1"},
	}
	return amUpdateReqData
}

//-------------------------------------------------------------------------------------------------

//-------------------------------------------------------------------------------------------------
//Fail Test (Create part)
func GetamCreatefailnotifyURIData() models.PolicyAssociationRequest {
	//d := time.Date(2019, 7, 5, 12, 30, 0, 0, time.UTC)
	amCreatefailnotifyURIData := GetAMreqdata()
	amCreatefailnotifyURIData.NotificationUri = ""
	return amCreatefailnotifyURIData
}
func GetamCreatefailsupiData() models.PolicyAssociationRequest {
	//d := time.Date(2019, 7, 5, 12, 30, 0, 0, time.UTC)
	amCreatefailnotifysupiData := GetAMreqdata()
	amCreatefailnotifysupiData.Supi = "dadfasdfasd"
	return amCreatefailnotifysupiData
}
