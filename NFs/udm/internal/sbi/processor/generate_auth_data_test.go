package processor

import (
	"io"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/h2non/gock"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/free5gc/openapi"
	"github.com/free5gc/openapi/models"
	udm_context "github.com/free5gc/udm/internal/context"
	"github.com/free5gc/udm/internal/sbi/consumer"
	mockapp "github.com/free5gc/udm/pkg/mockapp"
)

func TestGenerateAuthDataProcedure(t *testing.T) {
	defer gock.Off() // Flush pending mocks after test execution

	openapi.InterceptH2CClient()
	defer openapi.RestoreH2CClient()

	queryRes := models.AuthenticationSubscription{
		AuthenticationMethod:          models.AuthMethod__5_G_AKA,
		EncPermanentKey:               "8baf473f2f8fd09487cccbd7097c6862",
		ProtectionParameterId:         "8baf473f2f8fd09487cccbd7097c6862",
		SequenceNumber:                &models.SequenceNumber{Sqn: "000000000023"},
		AuthenticationManagementField: "8000",
		AlgorithmId:                   "128-EEA0",
		EncOpcKey:                     "8e27b6af0e692e750f32667a3b14605d",
		EncTopcKey:                    "8e27b6",
	}

	gock.New("http://127.0.0.4:8000/nudr-dr/v2").
		Get("/subscription-data/imsi-208930000000001/authentication-data/authentication-subscription").
		Reply(200).
		AddHeader("Content-Type", "application/json").
		JSON(queryRes)

	gock.New("http://127.0.0.4:8000").
		Patch("/nudr-dr/v2/subscription-data/imsi-208930000000001/authentication-data/authentication-subscription").
		Reply(204).
		JSON(map[string]string{})

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockApp := mockapp.NewMockApp(ctrl)
	testConsumer, err := consumer.NewConsumer(mockApp)
	require.NoError(t, err)
	testProcessor, err := NewProcessor(mockApp)
	require.NoError(t, err)
	udm_context.GetSelf().NrfUri = "http://127.0.0.10:8000"
	ue := new(udm_context.UdmUeContext)
	ue.Init()
	ue.Supi = "imsi-208930000000001"
	ue.UdrUri = "http://127.0.0.4:8000"
	udm_context.GetSelf().UdmUePool.Store("imsi-208930000000001", ue)

	mockApp.EXPECT().Consumer().Return(testConsumer).AnyTimes()
	mockApp.EXPECT().Context().Return(
		&udm_context.UDMContext{
			OAuth2Required: false,
			NrfUri:         "http://127.0.0.10:8000",
			NfId:           "1",
		},
	).AnyTimes()

	authInfoReq := models.AuthenticationInfoRequest{
		ServingNetworkName: "internet",
	}
	httpRecorder := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(httpRecorder)
	testProcessor.GenerateAuthDataProcedure(c, authInfoReq, "imsi-208930000000001")

	httpResp := httpRecorder.Result()
	if errClose := httpResp.Body.Close(); errClose != nil {
		t.Fatalf("Failed to close response body: %+v", errClose)
	}

	rawBytes, errReadAll := io.ReadAll(httpResp.Body)
	if errReadAll != nil {
		t.Fatalf("Failed to read response body: %+v", errReadAll)
	}

	var res models.UdmUeauAuthenticationInfoResult
	err = openapi.Deserialize(&res, rawBytes, httpResp.Header.Get("Content-Type"))
	if err != nil {
		t.Fatalf("Failed to deserialize response body: %+v", err)
	}

	expectResponse := models.UdmUeauAuthenticationInfoResult{
		AuthType: "5G_AKA",
		AuthenticationVector: &models.AuthenticationVector{
			AvType: "5G_HE_AKA",
			// Rand:     "6823f0dc9da02a61f6224d278a6f65b0",
			// Autn:     "1a7692538cf2800082662b9daa0396c5",
			// XresStar: "10e6e9e7d4ed291e7b2dc81e41f1da17",
			// Kausf:    "b0066a373555e6cd3efabfd14dfbc0bd91f5792445e2d80748daeb0a629d9b09",
		},
		Supi: "imsi-208930000000001",
	}

	require.Equal(t, 200, httpResp.StatusCode)

	// Since GenerateAuthData have randomness, only check fix value
	require.Equal(t, expectResponse.AuthType, res.AuthType)
	require.Equal(t, expectResponse.Supi, res.Supi)
	require.Equal(t, expectResponse.AuthenticationVector.AvType, res.AuthenticationVector.AvType)
}
