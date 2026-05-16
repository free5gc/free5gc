package test_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func getOAuthToken(t *testing.T, targetNfType, scope string) string {
	t.Helper()
	nfID := "88924b10-60b6-455b-9d41-356c3ee72e1f"
	nfProfile := map[string]interface{}{
		"nfInstanceId": nfID,
		"nfType":       "AF",
		"nfStatus":     "REGISTERED",
		"nfServices": []map[string]interface{}{
			{
				"serviceInstanceId": "1",
				"serviceName":       "nnef-callback",
				"versions": []map[string]interface{}{
					{"apiVersionInUri": "v1", "apiFullVersion": "1.0.0"},
				},
				"scheme":          "http",
				"nfServiceStatus": "REGISTERED",
			},
		},
	}
	b, _ := json.Marshal(nfProfile)
	// Register AF in NRF
	req, _ := http.NewRequest(http.MethodPut,
		"http://127.0.0.10:8000/nnrf-nfm/v1/nf-instances/"+nfID, bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	data := url.Values{}
	data.Set("grant_type", "client_credentials")
	data.Set("nfInstanceId", nfID)
	data.Set("nfType", "AF")
	data.Set("targetNfType", targetNfType)
	data.Set("scope", scope)

	req2, _ := http.NewRequest(http.MethodPost,
		"http://127.0.0.10:8000/oauth2/token", strings.NewReader(data.Encode()))
	req2.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp2, err := http.DefaultClient.Do(req2)
	require.NoError(t, err)
	defer resp2.Body.Close()

	var res map[string]interface{}
	json.NewDecoder(resp2.Body).Decode(&res)

	if token, ok := res["access_token"].(string); ok {
		return token
	}

	t.Fatalf("Failed to get OAuth token for target %s (scope: %s). "+
		"Please check if the target NF registered this service in its YAML config. Error: %v",
		targetNfType, scope, res)
	return ""
}

func TestOAuth2Callback(t *testing.T) {
	t.Log("[TestOAuth2Callback] Running in STRICT OAuth2 mode.")

	var afCallCount int64
	afMock := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt64(&afCallCount, 1)
		t.Logf("[AF mock] notification received: method=%s path=%s", r.Method, r.URL.Path)
		w.WriteHeader(http.StatusNoContent)
	}))
	defer afMock.Close()

	afNotifURL := afMock.URL + "/callback/notify"
	afID := "af-e2e-callback"

	subBody := map[string]interface{}{
		"afAppId":                 "app_video_1",
		"notificationDestination": afNotifURL,
		"suppFeat":                "0",
		"dnn":                     "internet",
		"snssai": map[string]interface{}{
			"sst": 1,
			"sd":  "010203",
		},
		"anyUeInd": true,
		"trafficRoutes": []map[string]interface{}{
			{
				"ipv4Addr": "10.60.0.103",
			},
		},
	}
	subBodyJSON, _ := json.Marshal(subBody)

	subReqURL := "http://127.0.0.5:8000/3gpp-traffic-influence/v1/" + afID + "/subscriptions"
	reqSub, _ := http.NewRequest(http.MethodPost, subReqURL, bytes.NewReader(subBodyJSON))
	reqSub.Header.Set("Content-Type", "application/json")
	reqSub.Header.Set("Authorization", "Bearer "+getOAuthToken(t, "NEF", "3gpp-traffic-influence"))

	subResp, err := http.DefaultClient.Do(reqSub)
	require.NoError(t, err)
	defer subResp.Body.Close()
	require.Equal(t, http.StatusCreated, subResp.StatusCode)

	notifCorreID := "1"
	smfNotif := map[string]interface{}{
		"notifId":   notifCorreID,
		"notifType": "UP_PATH_CH",
		"eventNotifs": []map[string]interface{}{
			{"event": "UP_PATH_CH", "dnaiChgType": "EARLY"},
		},
	}
	notifBody, _ := json.Marshal(smfNotif)

	callbackURL := "http://127.0.0.5:8000/nnef-callback/v1/notification/smf"
	reqNotif, _ := http.NewRequest(http.MethodPost, callbackURL, bytes.NewReader(notifBody))
	reqNotif.Header.Set("Content-Type", "application/json")
	reqNotif.Header.Set("Authorization", "Bearer "+getOAuthToken(t, "NEF", "nnef-callback"))

	notifResp, err := http.DefaultClient.Do(reqNotif)
	require.NoError(t, err)
	defer notifResp.Body.Close()

	assert.Equal(t, http.StatusNoContent, notifResp.StatusCode)

	time.Sleep(200 * time.Millisecond)
	assert.Equal(t, int64(1), atomic.LoadInt64(&afCallCount))

	loc := subResp.Header.Get("Location")
	if loc != "" {
		delReq, _ := http.NewRequest(http.MethodDelete, loc, nil)
		delReq.Header.Set("Authorization", "Bearer "+getOAuthToken(t, "NEF", "3gpp-traffic-influence"))
		delResp, err := http.DefaultClient.Do(delReq)
		if err == nil {
			defer delResp.Body.Close()
			assert.Equal(t, http.StatusNoContent, delResp.StatusCode)
		}
	}

	nfCallbacks := map[string]string{
		"AMF": "namf-callback",
		"SMF": "nsmf-callback",
		"PCF": "npcf-callback",
	}

	for nfType, scope := range nfCallbacks {
		t.Logf("[Strict Check] Verifying OAuth2 registration for %s service: %s", nfType, scope)
		token := getOAuthToken(t, nfType, scope)
		require.NotEmpty(t, token)
		t.Logf("[Strict Check] Success: %s has registered %s and NRF issued a token.", nfType, scope)
	}

	t.Log("[TestOAuth2Callback] PASS - E2E Callback and Service Registration verified.")
}
