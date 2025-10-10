package context_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/free5gc/nrf/internal/context"
	"github.com/free5gc/nrf/pkg/factory"
	"github.com/free5gc/openapi/models"
)

// Test 1.1.1: 驗證 Nrf_NfInstanceID 不為空
func TestNrfInstanceIDNotEmpty(t *testing.T) {
	// Given: NRF 初始化配置
	setupTestConfig()

	// When: 執行 InitNrfContext()
	err := context.InitNrfContext()
	require.NoError(t, err, "InitNrfContext should not return error")

	// Then: Nrf_NfInstanceID 不為空
	nrfCtx := context.GetSelf()
	require.NotEmpty(t, nrfCtx.Nrf_NfInstanceID,
		"Nrf_NfInstanceID should not be empty after initialization")
}

// Test 1.1.2: 驗證兩個 Instance ID 字段一致
func TestNrfInstanceIDConsistency(t *testing.T) {
	// Given: NRF 已初始化
	setupTestConfig()
	err := context.InitNrfContext()
	require.NoError(t, err)

	// When: 檢查 Nrf_NfInstanceID 和 NrfNfProfile.NfInstanceId
	nrfCtx := context.GetSelf()

	// Then: 兩者應該相同
	require.Equal(t, nrfCtx.NrfNfProfile.NfInstanceId, nrfCtx.Nrf_NfInstanceID,
		"Nrf_NfInstanceID should match NrfNfProfile.NfInstanceId")
}

// Test 1.1.3: 驗證 UUID 格式正確
func TestNrfInstanceIDValidUUID(t *testing.T) {
	// Given: NRF 已初始化
	setupTestConfig()
	err := context.InitNrfContext()
	require.NoError(t, err)

	// When: 解析 Nrf_NfInstanceID
	nrfCtx := context.GetSelf()
	instanceID := nrfCtx.Nrf_NfInstanceID

	// Then: 應該是有效的 UUID v4 格式
	require.NotEmpty(t, instanceID)
	// UUID v4 格式：xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx (36 chars with dashes)
	require.Len(t, instanceID, 36, "UUID should be 36 characters long")
	require.Contains(t, instanceID, "-", "UUID should contain dashes")
}

// setupTestConfig 設置測試用的 NRF 配置
func setupTestConfig() {
	factory.NrfConfig = &factory.Config{
		Info: &factory.Info{
			Version:     "1.0.2",
			Description: "NRF test configuration",
		},
		Configuration: &factory.Configuration{
			MongoDBName: "test_free5gc",
			MongoDBUrl:  "mongodb://127.0.0.1:27017",
			Sbi: &factory.Sbi{
				Scheme:       "http",
				RegisterIPv4: "127.0.0.10",
				BindingIPv4:  "127.0.0.10",
				Port:         8000,
				Cert: &factory.Cert{
					Pem: "../../../cert/nrf.pem",
					Key: "../../../cert/nrf.key",
				},
				RootCert: &factory.Cert{
					Pem: "../../../cert/root.pem",
					Key: "../../../cert/root.key",
				},
				OAuth: false, // 先不啟用 OAuth 以簡化測試
			},
			DefaultPlmnId: models.PlmnId{
				Mcc: "208",
				Mnc: "93",
			},
			ServiceNameList: []string{
				"nnrf-nfm",
				"nnrf-disc",
			},
		},
		Logger: &factory.Logger{
			Enable:       true,
			Level:        "info",
			ReportCaller: false,
		},
	}
}
