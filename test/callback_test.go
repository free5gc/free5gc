package test_test

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"sync/atomic"
	"test"
	"testing"
	"time"

	"github.com/free5gc/openapi/models"
	"github.com/free5gc/openapi/oauth"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const testAFInstanceID = "88924b10-60b6-455b-9d41-356c3ee72e1f"

type oauthTokenResponse struct {
	AccessToken string `json:"access_token"`
	Error       string `json:"error"`
}

func prepareNFInstanceCertificate(
	t *testing.T, clientNfType models.Nrf_NFMgmt_NFType, clientNfID string,
) {
	t.Helper()

	certDir, ok := test.OAuthCertificateDirectory()
	require.True(t, ok, "OAuth certificate directory is not configured")
	rootCert, err := oauth.ParseCertFromPEM("../cert/root.pem")
	require.NoError(t, err)
	rootPrivateKey, err := oauth.ParsePrivateKeyFromPEM("../cert/root.key")
	require.NoError(t, err)
	nfPrivateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)

	certPath := oauth.GetNFCertPath(certDir, string(clientNfType), clientNfID)
	_, err = oauth.GenerateCertificate(
		string(clientNfType), clientNfID, certPath,
		&nfPrivateKey.PublicKey, rootCert, rootPrivateKey,
	)
	require.NoError(t, err)
}

func registerTestNF(
	t *testing.T,
	clientNfID string,
	clientNfType models.Nrf_NFMgmt_NFType,
	serviceName models.Nrf_NFMgmt_ServiceName,
) {
	t.Helper()

	nfProfile := models.Nrf_NFMgmt_NFProfile{
		NfInstanceId: clientNfID,
		NfType:       clientNfType,
		NfStatus:     models.Nrf_NFMgmt_NFStatus_REGISTERED,
		NfServices: []models.Nrf_NFMgmt_NFService{{
			ServiceInstanceId: "1",
			ServiceName:       serviceName,
			Versions: []models.Nrf_NFMgmt_NFServiceVersion{{
				ApiVersionInUri: "v1",
				ApiFullVersion:  "1.0.0",
			}},
			Scheme:          models.UriScheme_HTTP,
			NfServiceStatus: models.Nrf_NFMgmt_NFServiceStatus_REGISTERED,
		}},
	}
	b, err := json.Marshal(nfProfile)
	require.NoError(t, err)
	req, err := http.NewRequest(http.MethodPut,
		"http://127.0.0.10:8000/nnrf-nfm/v1/nf-instances/"+clientNfID, bytes.NewReader(b))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()
	require.Contains(t, []int{http.StatusOK, http.StatusCreated}, resp.StatusCode)
}

func requestOAuthToken(
	t *testing.T,
	clientNfID string,
	clientNfType models.Nrf_NFMgmt_NFType,
	targetNfType models.Nrf_NFMgmt_NFType,
	targetNfInstanceID string,
	scope models.Nrf_NFMgmt_ServiceName,
) (int, oauthTokenResponse) {
	t.Helper()

	data := url.Values{}
	data.Set("grant_type", "client_credentials")
	data.Set("nfInstanceId", clientNfID)
	data.Set("nfType", string(clientNfType))
	data.Set("targetNfType", string(targetNfType))
	if targetNfInstanceID != "" {
		data.Set("targetNfInstanceId", targetNfInstanceID)
	}
	data.Set("scope", string(scope))

	req, err := http.NewRequest(http.MethodPost,
		"http://127.0.0.10:8000/oauth2/token", strings.NewReader(data.Encode()))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	var result oauthTokenResponse
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&result))
	return resp.StatusCode, result
}

func requireOAuthToken(
	t *testing.T,
	clientNfID string,
	clientNfType models.Nrf_NFMgmt_NFType,
	targetNfType models.Nrf_NFMgmt_NFType,
	targetNfInstanceID string,
	scope models.Nrf_NFMgmt_ServiceName,
) string {
	t.Helper()

	status, result := requestOAuthToken(
		t, clientNfID, clientNfType, targetNfType, targetNfInstanceID, scope,
	)
	require.Equal(t, http.StatusOK, status, "OAuth error: %s", result.Error)
	require.NotEmpty(t, result.AccessToken)
	return result.AccessToken
}

func nfInstanceID(t *testing.T, nfType models.Nrf_NFMgmt_NFType) string {
	t.Helper()
	id, ok := test.GetNFInstanceID(nfType)
	require.True(t, ok, "test NF %s is not configured", nfType)
	require.NotEmpty(t, id)
	return id
}

func callCallback(t *testing.T, method, callbackURL string, body []byte, token string) int {
	t.Helper()
	req, err := http.NewRequest(method, callbackURL, bytes.NewReader(body))
	require.NoError(t, err)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()
	return resp.StatusCode
}

func TestOAuth2Callback(t *testing.T) {
	t.Log("[TestOAuth2Callback] Running in STRICT OAuth2 mode.")

	clientNfID := testAFInstanceID
	prepareNFInstanceCertificate(t, models.Nrf_NFMgmt_NFType_AF, clientNfID)
	registerTestNF(t, clientNfID, models.Nrf_NFMgmt_NFType_AF, models.Nrf_NFMgmt_ServiceName("nnef-callback"))
	t.Logf("[TestOAuth2Callback] Using fixed AF Instance ID: %s", clientNfID)

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
	reqSub.Header.Set("Authorization", "Bearer "+requireOAuthToken(
		t, clientNfID, models.Nrf_NFMgmt_NFType_AF,
		models.Nrf_NFMgmt_NFType_NEF, nfInstanceID(t, models.Nrf_NFMgmt_NFType_NEF),
		models.Nrf_NFMgmt_ServiceName("3gpp-traffic-influence"),
	))

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
	reqNotif.Header.Set("Authorization", "Bearer "+requireOAuthToken(
		t, nfInstanceID(t, models.Nrf_NFMgmt_NFType_SMF), models.Nrf_NFMgmt_NFType_SMF,
		models.Nrf_NFMgmt_NFType_NEF, nfInstanceID(t, models.Nrf_NFMgmt_NFType_NEF),
		models.Nrf_NFMgmt_ServiceName("nnef-callback"),
	))

	notifResp, err := http.DefaultClient.Do(reqNotif)
	require.NoError(t, err)
	defer notifResp.Body.Close()

	assert.Equal(t, http.StatusNoContent, notifResp.StatusCode)

	time.Sleep(200 * time.Millisecond)
	assert.Equal(t, int64(1), atomic.LoadInt64(&afCallCount))

	loc := subResp.Header.Get("Location")
	if loc != "" {
		delReq, _ := http.NewRequest(http.MethodDelete, loc, nil)
		delReq.Header.Set("Authorization", "Bearer "+requireOAuthToken(
			t, clientNfID, models.Nrf_NFMgmt_NFType_AF,
			models.Nrf_NFMgmt_NFType_NEF, nfInstanceID(t, models.Nrf_NFMgmt_NFType_NEF),
			models.Nrf_NFMgmt_ServiceName("3gpp-traffic-influence"),
		))
		delResp, err := http.DefaultClient.Do(delReq)
		if err == nil {
			defer delResp.Body.Close()
			assert.Equal(t, http.StatusNoContent, delResp.StatusCode)
		}
	}

	callbackCases := []struct {
		name           string
		consumerType   models.Nrf_NFMgmt_NFType
		targetType     models.Nrf_NFMgmt_NFType
		scope          models.Nrf_NFMgmt_ServiceName
		method         string
		callbackURL    string
		body           []byte
		expectedStatus int
		typeScoped     bool
	}{
		{
			name: "PCF to AMF", consumerType: models.Nrf_NFMgmt_NFType_PCF,
			targetType: models.Nrf_NFMgmt_NFType_AMF, scope: models.Nrf_NFMgmt_ServiceName("namf-callback"),
			method: http.MethodGet, callbackURL: "http://127.0.0.18:8000/namf-callback/v1/",
			expectedStatus: http.StatusOK,
		},
		{
			name: "UDM to AMF", consumerType: models.Nrf_NFMgmt_NFType_UDM,
			targetType: models.Nrf_NFMgmt_NFType_AMF, scope: models.Nrf_NFMgmt_ServiceName("namf-callback"),
			method: http.MethodGet, callbackURL: "http://127.0.0.18:8000/namf-callback/v1/",
			expectedStatus: http.StatusOK,
		},
		{
			name: "AMF to PCF", consumerType: models.Nrf_NFMgmt_NFType_AMF,
			targetType: models.Nrf_NFMgmt_NFType_PCF, scope: models.Nrf_NFMgmt_ServiceName("npcf-callback"),
			method: http.MethodPost, callbackURL: "http://127.0.0.7:8000/npcf-callback/v1/amfstatus",
			body: []byte(`{}`), expectedStatus: http.StatusNoContent,
		},
		{
			name: "UDR to PCF", consumerType: models.Nrf_NFMgmt_NFType_UDR,
			targetType: models.Nrf_NFMgmt_NFType_PCF, scope: models.Nrf_NFMgmt_ServiceName("npcf-callback"),
			method:      http.MethodPost,
			callbackURL: "http://127.0.0.7:8000/npcf-callback/v1/nudr-notify/policy-data/imsi-test",
			body:        []byte(`{}`), expectedStatus: http.StatusNotImplemented,
		},
		{
			name: "UDR to UDM", consumerType: models.Nrf_NFMgmt_NFType_UDR,
			targetType: models.Nrf_NFMgmt_NFType_UDM, scope: models.Nrf_NFMgmt_ServiceName_NUDM_SDM,
			method:         http.MethodPost,
			callbackURL:    "http://127.0.0.3:8000/imsi-208930000000001/sdm-subscriptions",
			body:           []byte(`{"notifyItems":[{"resourceId":"resource"}]}`),
			expectedStatus: http.StatusNoContent,
		},
		{
			name: "SMF to AMF", consumerType: models.Nrf_NFMgmt_NFType_SMF,
			targetType: models.Nrf_NFMgmt_NFType_AMF, scope: models.Nrf_NFMgmt_ServiceName("namf-callback"),
			method: http.MethodGet, callbackURL: "http://127.0.0.18:8000/namf-callback/v1/",
			expectedStatus: http.StatusOK, typeScoped: true,
		},
		{
			name: "PCF to SMF", consumerType: models.Nrf_NFMgmt_NFType_PCF,
			targetType: models.Nrf_NFMgmt_NFType_SMF, scope: models.Nrf_NFMgmt_ServiceName("nsmf-callback"),
			method:      http.MethodPost,
			callbackURL: "http://127.0.0.2:8000/nsmf-callback/v1/sm-policies/ref/terminate",
			body:        []byte(`{}`), expectedStatus: http.StatusNotImplemented, typeScoped: true,
		},
	}

	for _, tc := range callbackCases {
		t.Run(tc.name, func(t *testing.T) {
			targetInstanceID := nfInstanceID(t, tc.targetType)
			if tc.typeScoped {
				targetInstanceID = ""
			}
			token := requireOAuthToken(
				t, nfInstanceID(t, tc.consumerType), tc.consumerType,
				tc.targetType, targetInstanceID, tc.scope,
			)
			status := callCallback(t, tc.method, tc.callbackURL, tc.body, token)
			require.Equal(t, tc.expectedStatus, status)
		})
	}

	rejectedRequests := []struct {
		name         string
		consumerID   string
		consumerType models.Nrf_NFMgmt_NFType
		targetType   models.Nrf_NFMgmt_NFType
		scope        models.Nrf_NFMgmt_ServiceName
	}{
		{
			name: "AF cannot call AMF callback", consumerID: clientNfID,
			consumerType: models.Nrf_NFMgmt_NFType_AF,
			targetType:   models.Nrf_NFMgmt_NFType_AMF, scope: models.Nrf_NFMgmt_ServiceName("namf-callback"),
		},
		{
			name:         "AUSF cannot call AMF callback",
			consumerType: models.Nrf_NFMgmt_NFType_AUSF,
			targetType:   models.Nrf_NFMgmt_NFType_AMF, scope: models.Nrf_NFMgmt_ServiceName("namf-callback"),
		},
		{
			name:         "callback scope with wrong target",
			consumerType: models.Nrf_NFMgmt_NFType_AMF,
			targetType:   models.Nrf_NFMgmt_NFType_AMF, scope: models.Nrf_NFMgmt_ServiceName("npcf-callback"),
		},
		{
			name:         "wrong callback scope",
			consumerType: models.Nrf_NFMgmt_NFType_UDR,
			targetType:   models.Nrf_NFMgmt_NFType_PCF, scope: models.Nrf_NFMgmt_ServiceName("namf-callback"),
		},
	}

	for _, tc := range rejectedRequests {
		t.Run(tc.name, func(t *testing.T) {
			consumerID := tc.consumerID
			if consumerID == "" {
				consumerID = nfInstanceID(t, tc.consumerType)
			}
			status, result := requestOAuthToken(
				t, consumerID, tc.consumerType,
				tc.targetType, nfInstanceID(t, tc.targetType), tc.scope,
			)
			require.Equal(t, http.StatusBadRequest, status)
			require.Equal(t, "invalid_scope", result.Error)
			require.Empty(t, result.AccessToken)
		})
	}

	t.Log("[TestOAuth2Callback] PASS - E2E Callback and Service Registration verified.")
}
