package processor

import (
	"net/http"
	"testing"

	"github.com/free5gc/openapi/models"
	"github.com/free5gc/openapi/models_nef"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"gopkg.in/h2non/gock.v1"
)

var (
	tiSub1ForAf1 = models_nef.TrafficInfluSub{
		AfServiceId: "Service1",
		AfAppId:     "App1",
		Dnn:         "internet",
		Snssai: &models.Snssai{
			Sst: 1,
			Sd:  "010203",
		},
		AnyUeInd: true,
		TrafficFilters: []models.FlowInfo{
			{
				FlowId: 1,
				FlowDescriptions: []string{
					"permit out ip from 192.168.0.21 to 10.60.0.0/16",
				},
			},
		},
		TrafficRoutes: []models.RouteToLocation{
			{
				Dnai: "mec",
				RouteInfo: &models.RouteInformation{
					Ipv4Addr:   "10.60.0.1",
					PortNumber: 0,
				},
			},
		},
	}

	tiSub2ForAf1 = models_nef.TrafficInfluSub{
		AfServiceId: "Service2",
		AfAppId:     "App2",
		Dnn:         "internet",
		Snssai: &models.Snssai{
			Sst: 1,
			Sd:  "010203",
		},
		AnyUeInd: true,
		TrafficFilters: []models.FlowInfo{
			{
				FlowId: 1,
				FlowDescriptions: []string{
					"permit out ip from 192.168.0.22 to 10.60.0.0/16",
				},
			},
		},
		TrafficRoutes: []models.RouteToLocation{
			{
				Dnai: "mec",
				RouteInfo: &models.RouteInformation{
					Ipv4Addr:   "10.60.0.1",
					PortNumber: 0,
				},
			},
		},
	}

	tiSub3ForAf1 = models_nef.TrafficInfluSub{
		AfServiceId: "Service3",
		AfAppId:     "App3",
		Dnn:         "internet",
		Snssai: &models.Snssai{
			Sst: 1,
			Sd:  "010203",
		},
		Ipv4Addr: "10.60.0.10",
		TrafficFilters: []models.FlowInfo{
			{
				FlowId: 1,
				FlowDescriptions: []string{
					"permit out ip from 192.168.0.23 to 10.60.0.10",
				},
			},
		},
		TrafficRoutes: []models.RouteToLocation{
			{
				Dnai: "mec",
				RouteInfo: &models.RouteInformation{
					Ipv4Addr:   "10.60.0.1",
					PortNumber: 0,
				},
			},
		},
	}

	tiSub4ForAf1 = models_nef.TrafficInfluSub{
		AfServiceId: "Service4",
		Dnn:         "internet",
		Snssai: &models.Snssai{
			Sst: 1,
			Sd:  "010203",
		},
		Ipv4Addr: "10.60.0.10",
	}

	tiSub5ForAf1 = models_nef.TrafficInfluSub{
		AfServiceId: "Service5",
		AfAppId:     "App5",
		TrafficFilters: []models.FlowInfo{
			{
				FlowId: 1,
				FlowDescriptions: []string{
					"permit out ip from 192.168.0.23 to 10.60.0.10",
				},
			},
		},
		TrafficRoutes: []models.RouteToLocation{
			{
				Dnai: "mec",
				RouteInfo: &models.RouteInformation{
					Ipv4Addr:   "10.60.0.1",
					PortNumber: 0,
				},
			},
		},
	}

	tiSubPatch1ForAf1 = models_nef.TrafficInfluSubPatch{
		TrafficFilters: []models.FlowInfo{
			{
				FlowId: 1,
				FlowDescriptions: []string{
					"permit out ip from 192.168.0.25 to 10.60.0.10",
				},
			},
		},
		TrafficRoutes: []models.RouteToLocation{
			{
				Dnai: "mec5",
			},
		},
	}
)

func TestGetTrafficInfluenceSubscription(t *testing.T) {
	testCases := []struct {
		description      string
		afID             string
		expectedResponse *HandlerResponse
	}{
		{
			description: "TC1: AfID found, should return all TrafficInfluSub",
			afID:        "af1",
			expectedResponse: &HandlerResponse{
				Status: http.StatusOK,
				Body:   &[]models_nef.TrafficInfluSub{tiSub1ForAf1, tiSub2ForAf1},
			},
		},
		{
			description: "TC2: AfID not found, should return ProblemDetails",
			afID:        "af3",
			expectedResponse: &HandlerResponse{
				Status: http.StatusNotFound,
				Body: &models.ProblemDetails{
					Status: http.StatusNotFound,
					Title:  "Data not found",
					Detail: "AF is not found",
				},
			},
		},
	}

	nefCtx := nefApp.Context()
	af1 := nefCtx.NewAf("af1")
	af1.Mu.Lock()
	correID1 := nefCtx.NewCorreID()
	afSub1 := af1.NewSub(correID1, &tiSub1ForAf1)
	af1.Subs[afSub1.SubID] = afSub1

	correID2 := nefCtx.NewCorreID()
	afSub2 := af1.NewSub(correID2, &tiSub2ForAf1)
	af1.Subs[afSub2.SubID] = afSub2
	nefCtx.AddAf(af1)
	af1.Mu.Unlock()

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			rsp := nefApp.Processor().GetTrafficInfluenceSubscription(tc.afID)
			require.Equal(t, tc.expectedResponse.Status, rsp.Status)
			require.Equal(t, tc.expectedResponse.Headers, rsp.Headers)
			if trafficInfluSub, ok := tc.expectedResponse.Body.(*[]models_nef.TrafficInfluSub); ok {
				require.ElementsMatch(t, *trafficInfluSub, *rsp.Body.(*[]models_nef.TrafficInfluSub))
			} else {
				require.Equal(t, tc.expectedResponse.Body, rsp.Body)
			}
		})
	}

	nefCtx.DeleteAf(af1.AfID)
	nefCtx.ResetCorreID()
}

func TestGetGetIndividualTrafficInfluenceSubscription(t *testing.T) {
	testCases := []struct {
		description      string
		afID             string
		subID            string
		expectedResponse *HandlerResponse
	}{
		{
			description: "TC1: AfID & SubID found, should return the PfdDataforApp",
			afID:        "af1",
			subID:       "1",
			expectedResponse: &HandlerResponse{
				Status: http.StatusOK,
				Body:   &tiSub1ForAf1,
			},
		},
		{
			description: "TC2: AfID found but SubID not found , should return the ProblemDetails",
			afID:        "af1",
			subID:       "2",
			expectedResponse: &HandlerResponse{
				Status: http.StatusNotFound,
				Body: &models.ProblemDetails{
					Status: http.StatusNotFound,
					Title:  "Data not found",
					Detail: "Subscription is not found",
				},
			},
		},
		{
			description: "TC3: AF ID not found, should return ProblemDetails",
			afID:        "af3",
			subID:       "3",
			expectedResponse: &HandlerResponse{
				Status: http.StatusNotFound,
				Body: &models.ProblemDetails{
					Status: http.StatusNotFound,
					Title:  "Data not found",
					Detail: "AF is not found",
				},
			},
		},
	}

	nefCtx := nefApp.Context()
	af1 := nefCtx.NewAf("af1")
	af1.Mu.Lock()
	correID1 := nefCtx.NewCorreID()
	afSub1 := af1.NewSub(correID1, &tiSub1ForAf1)
	af1.Subs[afSub1.SubID] = afSub1
	nefCtx.AddAf(af1)
	af1.Mu.Unlock()

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			rsp := nefApp.Processor().GetIndividualTrafficInfluenceSubscription(tc.afID, tc.subID)
			require.Equal(t, tc.expectedResponse, rsp)
		})
	}

	nefCtx.DeleteAf(af1.AfID)
	nefCtx.ResetCorreID()
}

func TestPostTrafficInfluenceSubscription(t *testing.T) {
	initNRFDiscPCFStub()
	initUDRDrPutTiDataStub(http.StatusNoContent)
	initPCFPaPostAppSessionsStub(http.StatusCreated)
	defer gock.Off()

	rspTiSub1 := tiSub1ForAf1
	rspTiSub1.Self = nefApp.Processor().genTrafficInfluSubURI("af1", "1")

	rspTiSub2 := tiSub3ForAf1
	rspTiSub2.Self = nefApp.Processor().genTrafficInfluSubURI("af1", "2")

	testCases := []struct {
		description      string
		afID             string
		tiSub            *models_nef.TrafficInfluSub
		expectedResponse *HandlerResponse
	}{
		{
			description: "TC1: Successful AnyUE subscription, should put tiData to UDR",
			afID:        "af1",
			tiSub:       &tiSub1ForAf1,
			expectedResponse: &HandlerResponse{
				Status: http.StatusCreated,
				Headers: map[string][]string{
					"Location": {rspTiSub1.Self},
				},
				Body: &rspTiSub1,
			},
		},
		{
			description: "TC2: Successful UEIPv4 subscription, should post AppSession to PCF",
			afID:        "af1",
			tiSub:       &tiSub3ForAf1,
			expectedResponse: &HandlerResponse{
				Status: http.StatusCreated,
				Headers: map[string][]string{
					"Location": {rspTiSub2.Self},
				},
				Body: &rspTiSub2,
			},
		},
		{
			description: "TC3: Missing one of afAppId, trafficFilters or ethTrafficFilters",
			afID:        "af1",
			tiSub:       &tiSub4ForAf1,
			expectedResponse: &HandlerResponse{
				Status: http.StatusBadRequest,
				Body: &models.ProblemDetails{
					Status: http.StatusBadRequest,
					Title:  "Malformed request syntax",
					Detail: "Missing one of afAppId, trafficFilters or ethTrafficFilters",
				},
			},
		},
		{
			description: "TC4: Missing one of Gpsi, Ipv4Addr, Ipv6Addr, ExternalGroupId, AnyUeInd",
			afID:        "af1",
			tiSub:       &tiSub5ForAf1,
			expectedResponse: &HandlerResponse{
				Status: http.StatusBadRequest,
				Body: &models.ProblemDetails{
					Status: http.StatusBadRequest,
					Title:  "Malformed request syntax",
					Detail: "Missing one of Gpsi, Ipv4Addr, Ipv6Addr, ExternalGroupId, AnyUeInd",
				},
			},
		},
	}

	nefCtx := nefApp.Context()
	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			rsp := nefApp.Processor().PostTrafficInfluenceSubscription(tc.afID, tc.tiSub)
			require.Equal(t, tc.expectedResponse, rsp)
		})
	}
	nefCtx.DeleteAf("af1")
	nefCtx.ResetCorreID()
}

func TestDeleteIndividualTrafficInfluenceSubscription(t *testing.T) {
	initNRFDiscPCFStub()
	initUDRDrDeleteTiDataStub(http.StatusNoContent)
	initPCFPaDeleteAppSessionsStub(http.StatusNoContent)
	defer gock.Off()

	testCases := []struct {
		description      string
		afID             string
		subID            string
		expectedResponse *HandlerResponse
	}{
		{
			description: "TC1: Successful delete TI subscription to UDR",
			afID:        "af1",
			subID:       "1",
			expectedResponse: &HandlerResponse{
				Status: http.StatusNoContent,
			},
		},
		{
			description: "TC1: Successful delete TI subscription to PCF",
			afID:        "af1",
			subID:       "2",
			expectedResponse: &HandlerResponse{
				Status: http.StatusNoContent,
			},
		},
		{
			description: "TC3: Delete non-existed TI subscription",
			afID:        "af1",
			subID:       "3",
			expectedResponse: &HandlerResponse{
				Status: http.StatusNotFound,
				Body: &models.ProblemDetails{
					Status: http.StatusNotFound,
					Title:  "Data not found",
					Detail: "Subscription is not found",
				},
			},
		},
	}

	nefCtx := nefApp.Context()
	af1 := nefCtx.NewAf("af1")
	af1.Mu.Lock()
	correID1 := nefCtx.NewCorreID()
	afSub1 := af1.NewSub(correID1, &tiSub1ForAf1)
	afSub1.InfluID = uuid.New().String()
	af1.Subs[afSub1.SubID] = afSub1

	correID2 := nefCtx.NewCorreID()
	afSub2 := af1.NewSub(correID2, &tiSub3ForAf1)
	af1.Subs[afSub2.SubID] = afSub2
	afSub2.AppSessID = "12345"
	nefCtx.AddAf(af1)
	af1.Mu.Unlock()

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			rsp := nefApp.Processor().DeleteIndividualTrafficInfluenceSubscription(tc.afID, tc.subID)
			require.Equal(t, tc.expectedResponse, rsp)
		})
	}
	nefCtx.DeleteAf(af1.AfID)
	nefCtx.ResetCorreID()
}

func TestPatchIndividualTrafficInfluenceSubscription(t *testing.T) {
	initNRFDiscPCFStub()
	initUDRDrPatchTiDataStub(http.StatusNoContent)
	initPCFPaPatchAppSessionsStub(http.StatusNoContent)
	defer gock.Off()

	rspTiSub1 := tiSub1ForAf1
	rspTiSub1.TrafficFilters = tiSubPatch1ForAf1.TrafficFilters
	rspTiSub1.TrafficRoutes = tiSubPatch1ForAf1.TrafficRoutes

	rspTiSub2 := tiSub3ForAf1
	rspTiSub2.TrafficFilters = tiSubPatch1ForAf1.TrafficFilters
	rspTiSub2.TrafficRoutes = tiSubPatch1ForAf1.TrafficRoutes

	testCases := []struct {
		description      string
		afID             string
		subID            string
		tiSubPatch       *models_nef.TrafficInfluSubPatch
		expectedResponse *HandlerResponse
	}{
		{
			description: "TC1: Successful patch TI subscription to UDR",
			afID:        "af1",
			subID:       "1",
			tiSubPatch:  &tiSubPatch1ForAf1,
			expectedResponse: &HandlerResponse{
				Status: http.StatusOK,
				Body:   &rspTiSub1,
			},
		},
		{
			description: "TC2: Successful patch TI subscription to PCF",
			afID:        "af1",
			subID:       "2",
			tiSubPatch:  &tiSubPatch1ForAf1,
			expectedResponse: &HandlerResponse{
				Status: http.StatusOK,
				Body:   &rspTiSub2,
			},
		},
		{
			description: "TC3: Patch non-existed TI subscription",
			afID:        "af1",
			subID:       "3",
			tiSubPatch:  &tiSubPatch1ForAf1,
			expectedResponse: &HandlerResponse{
				Status: http.StatusNotFound,
				Body: &models.ProblemDetails{
					Status: http.StatusNotFound,
					Title:  "Data not found",
					Detail: "Subscription is not found",
				},
			},
		},
	}

	nefCtx := nefApp.Context()
	af1 := nefCtx.NewAf("af1")
	af1.Mu.Lock()
	correID1 := nefCtx.NewCorreID()
	afSub1 := af1.NewSub(correID1, &tiSub1ForAf1)
	afSub1.InfluID = uuid.New().String()
	af1.Subs[afSub1.SubID] = afSub1

	correID2 := nefCtx.NewCorreID()
	afSub2 := af1.NewSub(correID2, &tiSub3ForAf1)
	af1.Subs[afSub2.SubID] = afSub2
	afSub2.AppSessID = "12345"
	nefCtx.AddAf(af1)
	af1.Mu.Unlock()

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			rsp := nefApp.Processor().PatchIndividualTrafficInfluenceSubscription(
				tc.afID, tc.subID, tc.tiSubPatch)
			require.Equal(t, tc.expectedResponse, rsp)
		})
	}
	nefCtx.DeleteAf(af1.AfID)
	nefCtx.ResetCorreID()
}

func TestPutIndividualTrafficInfluenceSubscription(t *testing.T) {
	initNRFDiscPCFStub()
	initUDRDrPutTiDataStub(http.StatusNoContent)
	initPCFPaPostAppSessionsStub(http.StatusCreated)
	defer gock.Off()

	testCases := []struct {
		description      string
		afID             string
		subID            string
		tiSub            *models_nef.TrafficInfluSub
		expectedResponse *HandlerResponse
	}{
		{
			description: "TC1: Successful put TI subscription to UDR",
			afID:        "af1",
			subID:       "1",
			tiSub:       &tiSub2ForAf1,
			expectedResponse: &HandlerResponse{
				Status: http.StatusOK,
				Body:   &tiSub2ForAf1,
			},
		},
		{
			description: "TC2: Successful put TI subscription to PCF",
			afID:        "af1",
			subID:       "2",
			tiSub:       &tiSub2ForAf1,
			expectedResponse: &HandlerResponse{
				Status: http.StatusOK,
				Body:   &tiSub2ForAf1,
			},
		},
		{
			description: "TC3: Put non-existed TI subscription",
			afID:        "af1",
			subID:       "3",
			tiSub:       &tiSub2ForAf1,
			expectedResponse: &HandlerResponse{
				Status: http.StatusNotFound,
				Body: &models.ProblemDetails{
					Status: http.StatusNotFound,
					Title:  "Data not found",
					Detail: "Subscription is not found",
				},
			},
		},
		{
			description: "TC4: Missing one of afAppId, trafficFilters or ethTrafficFilters",
			afID:        "af1",
			subID:       "4",
			tiSub:       &tiSub4ForAf1,
			expectedResponse: &HandlerResponse{
				Status: http.StatusBadRequest,
				Body: &models.ProblemDetails{
					Status: http.StatusBadRequest,
					Title:  "Malformed request syntax",
					Detail: "Missing one of afAppId, trafficFilters or ethTrafficFilters",
				},
			},
		},
		{
			description: "TC5: Missing one of Gpsi, Ipv4Addr, Ipv6Addr, ExternalGroupId, AnyUeInd",
			afID:        "af1",
			subID:       "5",
			tiSub:       &tiSub5ForAf1,
			expectedResponse: &HandlerResponse{
				Status: http.StatusBadRequest,
				Body: &models.ProblemDetails{
					Status: http.StatusBadRequest,
					Title:  "Malformed request syntax",
					Detail: "Missing one of Gpsi, Ipv4Addr, Ipv6Addr, ExternalGroupId, AnyUeInd",
				},
			},
		},
	}

	nefCtx := nefApp.Context()
	af1 := nefCtx.NewAf("af1")
	af1.Mu.Lock()
	correID1 := nefCtx.NewCorreID()
	afSub1 := af1.NewSub(correID1, &tiSub1ForAf1)
	afSub1.InfluID = uuid.New().String()
	af1.Subs[afSub1.SubID] = afSub1

	correID2 := nefCtx.NewCorreID()
	afSub2 := af1.NewSub(correID2, &tiSub3ForAf1)
	af1.Subs[afSub2.SubID] = afSub2
	afSub2.AppSessID = "12345"
	nefCtx.AddAf(af1)
	af1.Mu.Unlock()

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			rsp := nefApp.Processor().PutIndividualTrafficInfluenceSubscription(
				tc.afID, tc.subID, tc.tiSub)
			require.Equal(t, tc.expectedResponse, rsp)
		})
	}
	nefCtx.DeleteAf(af1.AfID)
	nefCtx.ResetCorreID()
}

func initUDRDrPutTiDataStub(statusCode int) {
	gock.New("http://127.0.0.4:8000/nudr-dr/v1").
		Put("/application-data/influenceData/.*").
		Persist().
		Reply(statusCode)
}

func initUDRDrPatchTiDataStub(statusCode int) {
	gock.New("http://127.0.0.4:8000/nudr-dr/v1").
		Patch("/application-data/influenceData/.*").
		Persist().
		Reply(statusCode)
}

func initUDRDrDeleteTiDataStub(statusCode int) {
	gock.New("http://127.0.0.4:8000/nudr-dr/v1").
		Delete("/application-data/influenceData/.*").
		Persist().
		Reply(statusCode)
}

func initPCFPaPostAppSessionsStub(statusCode int) {
	asc3ForAf1 := &models.AppSessionContext{
		AscReqData: &models.AppSessionContextReqData{
			AfAppId: tiSub3ForAf1.AfAppId,
			AfRoutReq: &models.AfRoutingRequirement{
				AppReloc:    tiSub3ForAf1.AppReloInd,
				RouteToLocs: tiSub3ForAf1.TrafficRoutes,
				TempVals:    tiSub3ForAf1.TempValidities,
			},
			UeIpv4:    tiSub3ForAf1.Ipv4Addr,
			UeIpv6:    tiSub3ForAf1.Ipv6Addr,
			UeMac:     tiSub3ForAf1.MacAddr,
			NotifUri:  tiSub3ForAf1.NotificationDestination,
			SuppFeat:  tiSub3ForAf1.SuppFeat,
			Dnn:       tiSub3ForAf1.Dnn,
			SliceInfo: tiSub3ForAf1.Snssai,
			// Supi: ,
		},
	}

	gock.New("http://127.0.0.7:8000/npcf-policyauthorization/v1").
		Post("/app-sessions").
		Persist().
		Reply(statusCode).
		SetHeader("Location", "http://127.0.0.7:8000/npcf-policyauthorization/v1/app-sessions/12345").
		JSON(asc3ForAf1)
}

func initPCFPaPatchAppSessionsStub(statusCode int) {
	gock.New("http://127.0.0.7:8000/npcf-policyauthorization/v1").
		Patch("/app-sessions/12345").
		Persist().
		Reply(statusCode)
}

func initPCFPaDeleteAppSessionsStub(statusCode int) {
	gock.New("http://127.0.0.7:8000/npcf-policyauthorization/v1").
		Post("/app-sessions/12345/delete").
		Persist().
		Reply(statusCode)
}
