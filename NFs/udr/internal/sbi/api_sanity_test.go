package sbi

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"

	"github.com/free5gc/openapi/models"
	db "github.com/free5gc/udr/internal/database"
	"github.com/free5gc/udr/internal/logger"
	"github.com/free5gc/udr/internal/sbi/processor"
	"github.com/free5gc/udr/pkg/factory"
	util_logger "github.com/free5gc/util/logger"
	"github.com/free5gc/util/mongoapi"
)

type testdata struct {
	influId string
	supi    string
}

func setupHttpServer(t *testing.T) *gin.Engine {
	router := util_logger.NewGinWithLogrus(logger.GinLog)
	dataRepositoryGroup := router.Group(factory.UdrDrResUriPrefix)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	udr := NewMockUDR(ctrl)
	factory.UdrConfig = &factory.Config{
		Configuration: &factory.Configuration{
			DbConnectorType: "mongodb",
			Mongodb:         &factory.Mongodb{},
			Sbi: &factory.Sbi{
				BindingIPv4: "127.0.0.1",
				Port:        8000,
			},
		},
	}
	udr.EXPECT().
		Config().
		Return(factory.UdrConfig).
		AnyTimes()

	processor := processor.NewProcessor(udr)
	udr.EXPECT().Processor().Return(processor).AnyTimes()

	s := NewServer(udr, "")
	dataRepositoryRoutes := s.getDataRepositoryRoutes()
	AddService(dataRepositoryGroup, dataRepositoryRoutes)
	return router
}

func setupMongoDB(t *testing.T) {
	err := mongoapi.SetMongoDB("test5gc", "mongodb://localhost:27017")
	require.Nil(t, err)
	err = mongoapi.Drop(db.APPDATA_INFLUDATA_DB_COLLECTION_NAME)
	require.Nil(t, err)
	err = mongoapi.Drop(db.APPDATA_INFLUDATA_SUBSC_DB_COLLECTION_NAME)
	require.Nil(t, err)
	err = mongoapi.Drop(db.APPDATA_PFD_DB_COLLECTION_NAME)
	require.Nil(t, err)
}

func getUri(t *testing.T, baseUri, extUri string) *httptest.ResponseRecorder {
	server := setupHttpServer(t)
	reqUri := baseUri + extUri
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, reqUri, nil)
	require.Nil(t, err)
	rsp := httptest.NewRecorder()
	server.ServeHTTP(rsp, req)
	return rsp
}

func getInfluData(supi string) *models.TrafficInfluData {
	return &models.TrafficInfluData{
		Dnn: "internet",
		Snssai: &models.Snssai{
			Sst: 1, Sd: "010203",
		},
		Supi: supi,
		TrafficFilters: []models.FlowInfo{{
			FlowId:           1,
			FlowDescriptions: []string{"permit out ip from 60.60.0.1 8080 to any"},
		}},
		TrafficRoutes: []*models.RouteToLocation{{
			Dnai: "edge1", RouteProfId: "1",
		}, {
			Dnai: "edge2", RouteProfId: "2",
		}},
		NwAreaInfo: &models.NetworkAreaInfo{
			Tais: []models.Tai{{
				PlmnId: &models.PlmnId{
					Mcc: "208", Mnc: "93",
				},
				Tac: "1",
			}},
		},
	}
}

func postPutInfluData(t *testing.T, method string, baseUri, extUri string, influData *models.TrafficInfluData) (
	*httptest.ResponseRecorder, []byte,
) {
	server := setupHttpServer(t)
	reqUri := baseUri + extUri
	bjson, err := json.Marshal(influData)
	require.Nil(t, err)
	reqBody := bytes.NewReader(bjson)
	req, err := http.NewRequestWithContext(context.Background(), method, reqUri, reqBody)
	require.Nil(t, err)
	rsp := httptest.NewRecorder()
	server.ServeHTTP(rsp, req)
	return rsp, bjson
}

func postInfluData(t *testing.T, baseUri, extUri string, influData *models.TrafficInfluData) (
	*httptest.ResponseRecorder, []byte,
) {
	return postPutInfluData(t, http.MethodPost, baseUri, extUri, influData)
}

func putInfluData(t *testing.T, baseUri, extUri string, influData *models.TrafficInfluData) (
	*httptest.ResponseRecorder, []byte,
) {
	return postPutInfluData(t, http.MethodPut, baseUri, extUri, influData)
}

func delUri(t *testing.T, baseUri, extUri string) *httptest.ResponseRecorder {
	server := setupHttpServer(t)
	reqUri := baseUri + extUri
	req, err := http.NewRequestWithContext(context.Background(), http.MethodDelete, reqUri, nil)
	require.Nil(t, err)
	rsp := httptest.NewRecorder()
	server.ServeHTTP(rsp, req)
	return rsp
}

func TestUDR_Root(t *testing.T) {
	server := setupHttpServer(t)
	reqUri := factory.UdrDrResUriPrefix + "/"

	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, reqUri, nil)
	require.Nil(t, err)
	rsp := httptest.NewRecorder()
	server.ServeHTTP(rsp, req)

	t.Run("UDR Root", func(t *testing.T) {
		require.Equal(t, http.StatusNotImplemented, rsp.Code)
		require.Equal(t, "Hello World!", rsp.Body.String())
	})
}

func TestUDR_GetSubs2Notify_GetBeforeCreateingOne(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping testing in short mode")
	}

	server := setupHttpServer(t)
	reqUri := factory.UdrDrResUriPrefix + "/application-data/influenceData/subs-to-notify?dnn=internet"

	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, reqUri, nil)
	require.Nil(t, err)
	rsp := httptest.NewRecorder()
	server.ServeHTTP(rsp, req)

	t.Run("UDR subs-to-notify Get before Create, dnn==internet", func(t *testing.T) {
		require.Equal(t, http.StatusOK, rsp.Code)
		require.Equal(t, "null", rsp.Body.String())
	})
}

func TestUDR_GetSubs2Notify_CreateThenGet(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping testing in short mode")
	}

	server := setupHttpServer(t)
	baseUri := factory.UdrDrResUriPrefix + "/application-data/influenceData/subs-to-notify"
	reqUri := baseUri

	test := models.TrafficInfluSub{
		Dnns: []string{"internet", "outernet"},
		Snssais: []models.Snssai{{
			Sst: 1, Sd: "010203",
		}, {
			Sst: 1, Sd: "112233",
		}},
		NotificationUri: "http://127.0.0.1/notifyMePlease",
	}
	bjson, err := json.Marshal(test)
	require.Nil(t, err)
	reqBody := bytes.NewReader(bjson)

	test.NotificationUri = ""
	bjson_bad, err := json.Marshal(test)
	require.Nil(t, err)
	reqBody_bad := bytes.NewReader(bjson_bad)

	// Create one - w/o the mandatory NotificationUri:
	req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, reqUri, reqBody_bad)
	require.Nil(t, err)
	rsp := httptest.NewRecorder()
	server.ServeHTTP(rsp, req)
	t.Run("UDR subs-to-notify CreateThenGet - Create w/o mandatory notificationUri", func(t *testing.T) {
		require.Equal(t, http.StatusBadRequest, rsp.Code)
	})

	// Create one - normal
	req, err = http.NewRequestWithContext(context.Background(), http.MethodPost, reqUri, reqBody)
	require.Nil(t, err)
	rsp = httptest.NewRecorder()
	server.ServeHTTP(rsp, req)
	// Linter complains not closing response body even tried to close,
	// ==> remove the comments manually to test if location header is there.
	// location := rsp.Result().Header.Get("Location")
	// err = rsp.Result().Body.Close()
	// require.Nil(t, err)
	// require.NotNil(t, location)
	// t.Log("location:", location)
	t.Run("UDR subs-to-notify CreateThenGet - Create normal case", func(t *testing.T) {
		require.Equal(t, http.StatusCreated, rsp.Code)
		require.Equal(t, string(bjson), rsp.Body.String())
		// require.True(t, strings.Contains(location, baseUri+"/"))
		// require.True(t, strings.HasPrefix(location, udr_context.GetSelf().GetIPv4Uri()+baseUri+"/"))
	})

	// Get success
	rsp = getUri(t, baseUri, "?dnn=internet")
	t.Run("UDR subs-to-notify CreateThenGet - get", func(t *testing.T) {
		require.Equal(t, http.StatusOK, rsp.Code)
		require.Equal(t, "["+string(bjson)+"]", rsp.Body.String())
	})

	// Get without a filter
	rsp = getUri(t, baseUri, "")
	t.Run("UDR subs-to-notify CreateThenGet - get w/o a filter", func(t *testing.T) {
		require.Equal(t, http.StatusBadRequest, rsp.Code)
	})

	// Get a non-exist DNN
	rsp = getUri(t, baseUri, "?dnn=ThisIsABadDNN")
	t.Run("UDR subs-to-notify CreateThenGet - get bad DNN", func(t *testing.T) {
		require.Equal(t, http.StatusOK, rsp.Code)
		require.Equal(t, "null", rsp.Body.String())
	})
}

func TestUDR_InfluData_GetBeforeCreateing(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping testing in short mode")
	}

	server := setupHttpServer(t)
	reqUri := factory.UdrDrResUriPrefix + "/application-data/influenceData"

	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, reqUri, nil)
	require.Nil(t, err)
	rsp := httptest.NewRecorder()
	server.ServeHTTP(rsp, req)

	t.Run("UDR influ-data Get before Create",
		func(t *testing.T) {
			require.Equal(t, http.StatusOK, rsp.Code)
			require.Equal(t, "[]", rsp.Body.String())
		})
}

func TestUDR_InfluData_CreateThenGet(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping testing in short mode")
	}

	// PUT, PATCH, DELETE
	setupMongoDB(t)
	server := setupHttpServer(t)
	baseUri := factory.UdrDrResUriPrefix + "/application-data/influenceData"
	td1 := testdata{"/influenceId0001", "imsi-208930000000001"}
	td2 := testdata{"/influenceId0002", "imsi-208930000000002"}

	// Create one - bad method (POST not allowed)
	influData := getInfluData(td1.supi)
	rsp, _ := postInfluData(t, baseUri, td1.influId, influData)
	t.Run("UDR influ-data CreateThenGet - Create one - bad method",
		func(t *testing.T) {
			require.Equal(t, http.StatusMethodNotAllowed, rsp.Code)
		})

	// Create one - normal
	influData = getInfluData(td1.supi)
	rsp, bjson := putInfluData(t, baseUri, td1.influId, influData)
	t.Run("UDR influ-data CreateThenGet - Create normal case",
		func(t *testing.T) {
			require.Equal(t, http.StatusCreated, rsp.Code)
			require.Equal(t, string(bjson), rsp.Body.String())
		})

	// Create one - update existing one with identical data
	influData = getInfluData(td1.supi)
	rsp, bjson = putInfluData(t, baseUri, td1.influId, influData)
	t.Run("UDR influ-data CreateThenGet - Create - update existing one-identical data",
		func(t *testing.T) {
			require.Equal(t, http.StatusOK, rsp.Code)
			require.Equal(t, string(bjson), rsp.Body.String())
		})

	// Create one - update existing one with some difference
	influData = getInfluData(td1.supi)
	influData.Snssai.Sst = 2
	rsp, bjson = putInfluData(t, baseUri, td1.influId, influData)
	t.Run("UDR influ-data CreateThenGet - Create - update existing one-with some difference",
		func(t *testing.T) {
			require.Equal(t, http.StatusOK, rsp.Code)
			require.Equal(t, string(bjson), rsp.Body.String())
		})

	// Note: NOT WORING
	// Patch - update existing one with some difference
	// influData = &models.TrafficInfluData{
	// 	Snssai: &models.Snssai{
	// 		Sst: 1, Sd: "995995",
	// 	}}
	// bjson, err = json.Marshal(influData)
	// require.Nil(t, err)
	// reqBody = bytes.NewReader(bjson)
	// req, err = http.NewRequestWithContext(context.Background(), http.MethodPatch, reqUri, reqBody)
	// require.Nil(t, err)
	// rsp = httptest.NewRecorder()
	// server.ServeHTTP(rsp, req)
	// t.Run("UDR influ-data CreateThenGet - Patch - update existing one-with some difference",
	// 	func(t *testing.T) {
	// 		require.Equal(t, http.StatusNoContent, rsp.Code)
	// 		require.Equal(t, string(bjson), rsp.Body.String())
	// 	})

	// Get success
	rsp = getUri(t, baseUri, "?dnns="+influData.Dnn)
	testRsp := []models.TrafficInfluData{}
	err := json.Unmarshal(rsp.Body.Bytes(), &testRsp)
	require.Nil(t, err)
	t.Run("UDR influ-data CreateThenGet - get",
		func(t *testing.T) {
			require.Equal(t, http.StatusOK, rsp.Code)
			require.Equal(t, influData.Dnn, testRsp[0].Dnn)
			require.Equal(t, influData.Snssai, testRsp[0].Snssai)
			// ResUri differs here
		})

	// Create with td2 - normal
	influData = getInfluData(td2.supi)
	rsp, bjson = putInfluData(t, baseUri, td2.influId, influData)
	t.Run("UDR influ-data CreateThenGet - Create normal case",
		func(t *testing.T) {
			require.Equal(t, http.StatusCreated, rsp.Code)
			require.Equal(t, string(bjson), rsp.Body.String())
		})

	// Get - 2 influencesIds
	rsp = getUri(t, baseUri, "?dnns="+influData.Dnn)
	err = json.Unmarshal(rsp.Body.Bytes(), &testRsp)
	t.Log(rsp.Body.String())
	require.Nil(t, err)
	t.Run("UDR influ-data CreateThenGet - get - 2 influData",
		func(t *testing.T) {
			require.Equal(t, http.StatusOK, rsp.Code)
			require.Equal(t, 2, len(testRsp))
		})

	// Get a non-exist Supi
	rsp = getUri(t, baseUri, "?Supi=BadSupi")
	err = json.Unmarshal(rsp.Body.Bytes(), &testRsp)
	require.Nil(t, err)
	t.Run("UDR influ-data CreateThenGet - Bad DNN",
		func(t *testing.T) {
			require.Equal(t, http.StatusOK, rsp.Code)
			// expect zero influData
			require.Equal(t, 0, len(testRsp))
		})

	// Delete td2
	reqUri := baseUri + td2.influId
	req, err := http.NewRequestWithContext(context.Background(), http.MethodDelete, reqUri, nil)
	require.Nil(t, err)
	rsp = httptest.NewRecorder()
	server.ServeHTTP(rsp, req)
	t.Run("UDR influ-data CreateThenGet - Delete td2",
		func(t *testing.T) {
			require.Equal(t, http.StatusNoContent, rsp.Code)
		})

	// Get - 1 influenceId left
	rsp = getUri(t, baseUri, "?dnns="+influData.Dnn)
	err = json.Unmarshal(rsp.Body.Bytes(), &testRsp)
	t.Log(rsp.Body.String())
	require.Nil(t, err)
	t.Run("UDR influ-data CreateThenGet - get - 1 influData left",
		func(t *testing.T) {
			require.Equal(t, http.StatusOK, rsp.Code)
			require.Equal(t, 1, len(testRsp))
		})

	// Delete td1
	rsp = delUri(t, baseUri, td1.influId)
	t.Run("UDR influ-data CreateThenGet - Delete td1",
		func(t *testing.T) {
			require.Equal(t, http.StatusNoContent, rsp.Code)
		})

	// Get without a filter - 0 influenceId
	rsp = getUri(t, baseUri, "?dnns="+influData.Dnn)
	err = json.Unmarshal(rsp.Body.Bytes(), &testRsp)
	t.Log(rsp.Body.String())
	require.Nil(t, err)
	t.Run("UDR influ-data CreateThenGet - get - all influData deleted",
		func(t *testing.T) {
			require.Equal(t, http.StatusOK, rsp.Code)
			require.Equal(t, 0, len(testRsp))
		})
}
