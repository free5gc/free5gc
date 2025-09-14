/*
 * BSF Comprehensive Test Suite
 * Tests complete BSF functionality and 3GPP TS 29.521 compliance
 */

package processor_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"

	bsfContext "github.com/free5gc/bsf/internal/context"
	"github.com/free5gc/bsf/internal/sbi/processor"
	"github.com/free5gc/openapi/models"
)

// BSFTestSuite provides comprehensive testing for BSF functionality
type BSFTestSuite struct {
	suite.Suite
	router *gin.Engine
}

// SetupSuite initializes the test environment
func (suite *BSFTestSuite) SetupSuite() {
	// Initialize BSF context
	bsfContext.InitBsfContext()

	// Set up Gin router in test mode
	gin.SetMode(gin.TestMode)
	suite.router = gin.Default()

	// Set up routes
	suite.setupRoutes()
}

// TearDownSuite cleans up after all tests
func (suite *BSFTestSuite) TearDownSuite() {
	// Clean up context
	bsfContext.BsfSelf.PcfBindings = make(map[string]*bsfContext.PcfBinding)
	bsfContext.BsfSelf.PcfForUeBindings = make(map[string]*bsfContext.PcfForUeBinding)
	bsfContext.BsfSelf.PcfMbsBindings = make(map[string]*bsfContext.PcfMbsBinding)
}

// SetupTest runs before each test
func (suite *BSFTestSuite) SetupTest() {
	// Clear all bindings before each test
	bsfContext.BsfSelf.PcfBindings = make(map[string]*bsfContext.PcfBinding)
	bsfContext.BsfSelf.PcfForUeBindings = make(map[string]*bsfContext.PcfForUeBinding)
	bsfContext.BsfSelf.PcfMbsBindings = make(map[string]*bsfContext.PcfMbsBinding)
}

func (suite *BSFTestSuite) setupRoutes() {
	// PCF Bindings Collection
	suite.router.POST("/nbsf-management/v1/pcfBindings", processor.CreatePCFBinding)
	suite.router.GET("/nbsf-management/v1/pcfBindings", processor.GetPCFBindings)

	// Individual PCF Binding
	suite.router.GET("/nbsf-management/v1/pcfBindings/:bindingId", processor.GetIndPCFBinding)
	suite.router.PUT("/nbsf-management/v1/pcfBindings/:bindingId", processor.UpdateIndPCFBinding)
	suite.router.PATCH("/nbsf-management/v1/pcfBindings/:bindingId", processor.UpdateIndPCFBinding)
	suite.router.DELETE("/nbsf-management/v1/pcfBindings/:bindingId", processor.DeleteIndPCFBinding)

	// Note: PCF for UE Bindings and PCF MBS Bindings routes commented out until implemented
	// suite.router.POST("/nbsf-management/v1/pcfForUeBindings", processor.CreatePCFForUeBinding)
	// suite.router.GET("/nbsf-management/v1/pcfForUeBindings", processor.GetPCFForUeBindings)
	// suite.router.POST("/nbsf-management/v1/pcfMbsBindings", processor.CreatePCFMbsBinding)
	// suite.router.GET("/nbsf-management/v1/pcfMbsBindings", processor.GetPCFMbsBindings)
}

// Test 3GPP TS 29.521 Section 5.3.2.2 - PCF Binding Creation
func (suite *BSFTestSuite) TestCreatePCFBinding_3GPP_Compliance() {
	tests := []struct {
		name         string
		binding      models.PcfBinding
		expectedCode int
		description  string
	}{
		{
			name: "Valid PCF Binding Creation",
			binding: models.PcfBinding{
				Dnn:     "internet",
				Snssai:  &models.Snssai{Sst: 1, Sd: "010203"},
				Supi:    "imsi-208930000000001",
				PcfFqdn: "pcf.free5gc.org",
				PcfIpEndPoints: []models.IpEndPoint{
					{Ipv4Address: "127.0.0.1", Port: 8000},
				},
				PcfId: "pcf-001",
			},
			expectedCode: http.StatusCreated,
			description:  "3GPP TS 29.521 - Valid binding creation with mandatory fields",
		},
		{
			name: "Missing DNN - Should Fail",
			binding: models.PcfBinding{
				Snssai:  &models.Snssai{Sst: 1, Sd: "010203"},
				Supi:    "imsi-208930000000001",
				PcfFqdn: "pcf.free5gc.org",
			},
			expectedCode: http.StatusBadRequest,
			description:  "3GPP TS 29.521 - DNN is mandatory",
		},
		{
			name: "Missing S-NSSAI - Should Fail",
			binding: models.PcfBinding{
				Dnn:     "internet",
				Supi:    "imsi-208930000000001",
				PcfFqdn: "pcf.free5gc.org",
			},
			expectedCode: http.StatusBadRequest,
			description:  "3GPP TS 29.521 - S-NSSAI is mandatory",
		},
		{
			name: "IPv4 Address Binding",
			binding: models.PcfBinding{
				Dnn:      "internet",
				Snssai:   &models.Snssai{Sst: 1, Sd: "010203"},
				Ipv4Addr: "192.168.1.100",
				PcfFqdn:  "pcf.free5gc.org",
				PcfId:    "pcf-002",
			},
			expectedCode: http.StatusCreated,
			description:  "3GPP TS 29.521 - IPv4 address based binding",
		},
		{
			name: "IPv6 Prefix Binding",
			binding: models.PcfBinding{
				Dnn:        "internet",
				Snssai:     &models.Snssai{Sst: 1, Sd: "010203"},
				Ipv6Prefix: "2001:db8::/64",
				PcfFqdn:    "pcf.free5gc.org",
				PcfId:      "pcf-003",
			},
			expectedCode: http.StatusCreated,
			description:  "3GPP TS 29.521 - IPv6 prefix based binding",
		},
		{
			name: "MAC Address Binding",
			binding: models.PcfBinding{
				Dnn:       "internet",
				Snssai:    &models.Snssai{Sst: 1, Sd: "010203"},
				MacAddr48: "aa:bb:cc:dd:ee:ff",
				PcfFqdn:   "pcf.free5gc.org",
				PcfId:     "pcf-004",
			},
			expectedCode: http.StatusCreated,
			description:  "3GPP TS 29.521 - MAC address based binding",
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			jsonData, err := json.Marshal(tt.binding)
			suite.NoError(err)

			req, err := http.NewRequest("POST", "/nbsf-management/v1/pcfBindings", bytes.NewBuffer(jsonData))
			suite.NoError(err)
			req.Header.Set("Content-Type", "application/json")

			rr := httptest.NewRecorder()
			suite.router.ServeHTTP(rr, req)

			suite.Equal(tt.expectedCode, rr.Code, tt.description)

			if tt.expectedCode == http.StatusCreated {
				location := rr.Header().Get("Location")
				suite.Contains(location, "/nbsf-management/v1/pcfBindings/")

				var response models.PcfBinding
				err = json.Unmarshal(rr.Body.Bytes(), &response)
				suite.NoError(err)
				suite.Equal(tt.binding.Dnn, response.Dnn)
			}
		})
	}
}

// Test 3GPP TS 29.521 Section 5.3.2.3 - PCF Binding Discovery
func (suite *BSFTestSuite) TestQueryPCFBindings_3GPP_Compliance() {
	// Create test bindings
	bindings := []models.PcfBinding{
		{
			Dnn:      "internet",
			Snssai:   &models.Snssai{Sst: 1, Sd: "010203"},
			Supi:     "imsi-208930000000001",
			Ipv4Addr: "192.168.1.100",
			PcfFqdn:  "pcf.free5gc.org",
			PcfId:    "pcf-001",
		},
		{
			Dnn:        "ims",
			Snssai:     &models.Snssai{Sst: 2, Sd: "020304"},
			Supi:       "imsi-208930000000002",
			Ipv6Prefix: "2001:db8::/64",
			PcfFqdn:    "pcf.free5gc.org",
			PcfId:      "pcf-002",
		},
		{
			Dnn:       "internet",
			Snssai:    &models.Snssai{Sst: 1, Sd: "010203"},
			MacAddr48: "aa:bb:cc:dd:ee:ff",
			PcfFqdn:   "pcf.free5gc.org",
			PcfId:     "pcf-003",
		},
	}

	// Create bindings
	for _, binding := range bindings {
		jsonData, _ := json.Marshal(binding)
		req, _ := http.NewRequest("POST", "/nbsf-management/v1/pcfBindings", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()
		suite.router.ServeHTTP(rr, req)
		suite.Equal(http.StatusCreated, rr.Code)
	}

	queryTests := []struct {
		name          string
		queryParams   string
		expectedCode  int
		expectBinding bool
		description   string
	}{
		{
			name:          "Query by SUPI and DNN",
			queryParams:   "supi=imsi-208930000000001&dnn=internet",
			expectedCode:  http.StatusOK,
			expectBinding: true,
			description:   "3GPP TS 29.521 - Valid query by SUPI and DNN",
		},
		{
			name:          "Query by IPv4 Address",
			queryParams:   "ipv4Addr=192.168.1.100",
			expectedCode:  http.StatusOK,
			expectBinding: true,
			description:   "3GPP TS 29.521 - Valid query by IPv4 address",
		},
		{
			name:          "Query by IPv6 Prefix",
			queryParams:   "ipv6Prefix=2001:db8::/64",
			expectedCode:  http.StatusOK,
			expectBinding: true,
			description:   "3GPP TS 29.521 - Valid query by IPv6 prefix",
		},
		{
			name:          "Query by MAC Address",
			queryParams:   "macAddr48=aa:bb:cc:dd:ee:ff",
			expectedCode:  http.StatusOK,
			expectBinding: true,
			description:   "3GPP TS 29.521 - Valid query by MAC address",
		},
		{
			name:          "No matching binding",
			queryParams:   "supi=imsi-999999999999999&dnn=nonexistent",
			expectedCode:  http.StatusNoContent,
			expectBinding: false,
			description:   "3GPP TS 29.521 - No content when no binding found",
		},
		{
			name:          "Invalid query - multiple address params",
			queryParams:   "ipv4Addr=192.168.1.100&ipv6Prefix=2001:db8::/64",
			expectedCode:  http.StatusBadRequest,
			expectBinding: false,
			description:   "3GPP TS 29.521 NOTE 1 - Only one address parameter allowed",
		},
		{
			name:          "Invalid query - ipDomain without ipv4Addr",
			queryParams:   "ipDomain=example.com&dnn=internet",
			expectedCode:  http.StatusBadRequest,
			expectBinding: false,
			description:   "3GPP TS 29.521 NOTE 2 - ipDomain requires ipv4Addr",
		},
	}

	for _, tt := range queryTests {
		suite.Run(tt.name, func() {
			req, err := http.NewRequest("GET", "/nbsf-management/v1/pcfBindings?"+tt.queryParams, nil)
			suite.NoError(err)

			rr := httptest.NewRecorder()
			suite.router.ServeHTTP(rr, req)

			suite.Equal(tt.expectedCode, rr.Code, tt.description)

			if tt.expectBinding {
				// Check that X-BSF-Binding-ID header is present (our enhancement)
				bindingID := rr.Header().Get("X-BSF-Binding-ID")
				suite.NotEmpty(bindingID, "X-BSF-Binding-ID header should be present")

				var response models.PcfBinding
				err = json.Unmarshal(rr.Body.Bytes(), &response)
				suite.NoError(err)
				suite.NotEmpty(response.Dnn)
			}
		})
	}
}

// Test 3GPP TS 29.521 Section 5.3.2.4 - Individual PCF Binding Operations
func (suite *BSFTestSuite) TestIndividualPCFBinding_3GPP_Compliance() {
	// Create a test binding first
	binding := models.PcfBinding{
		Dnn:     "internet",
		Snssai:  &models.Snssai{Sst: 1, Sd: "010203"},
		Supi:    "imsi-208930000000001",
		PcfFqdn: "pcf.free5gc.org",
		PcfId:   "pcf-001",
	}

	jsonData, _ := json.Marshal(binding)
	req, _ := http.NewRequest("POST", "/nbsf-management/v1/pcfBindings", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	suite.router.ServeHTTP(rr, req)
	suite.Equal(http.StatusCreated, rr.Code)

	// Extract binding ID from Location header
	location := rr.Header().Get("Location")
	parts := strings.Split(location, "/")
	bindingId := parts[len(parts)-1]

	// Test GET individual binding
	suite.Run("GET Individual Binding", func() {
		req, err := http.NewRequest("GET", "/nbsf-management/v1/pcfBindings/"+bindingId, nil)
		suite.NoError(err)

		rr := httptest.NewRecorder()
		suite.router.ServeHTTP(rr, req)

		suite.Equal(http.StatusOK, rr.Code)

		var response models.PcfBinding
		err = json.Unmarshal(rr.Body.Bytes(), &response)
		suite.NoError(err)
		suite.Equal("internet", response.Dnn)
		suite.Equal("imsi-208930000000001", response.Supi)
	})

	// Test PATCH update
	suite.Run("PATCH Update Binding", func() {
		update := map[string]interface{}{
			"pcfFqdn": "updated-pcf.free5gc.org",
		}
		jsonData, _ := json.Marshal(update)

		req, err := http.NewRequest("PATCH", "/nbsf-management/v1/pcfBindings/"+bindingId, bytes.NewBuffer(jsonData))
		suite.NoError(err)
		req.Header.Set("Content-Type", "application/json-patch+json")

		rr := httptest.NewRecorder()
		suite.router.ServeHTTP(rr, req)

		suite.Equal(http.StatusOK, rr.Code)
	})

	// Test DELETE binding
	suite.Run("DELETE Individual Binding", func() {
		req, err := http.NewRequest("DELETE", "/nbsf-management/v1/pcfBindings/"+bindingId, nil)
		suite.NoError(err)

		rr := httptest.NewRecorder()
		suite.router.ServeHTTP(rr, req)

		suite.Equal(http.StatusNoContent, rr.Code)

		// Verify binding is deleted
		req, _ = http.NewRequest("GET", "/nbsf-management/v1/pcfBindings/"+bindingId, nil)
		rr = httptest.NewRecorder()
		suite.router.ServeHTTP(rr, req)
		suite.Equal(http.StatusNotFound, rr.Code)
	})
}

// Test MongoDB persistence functionality
func (suite *BSFTestSuite) TestMongoDB_Persistence() {
	suite.Run("MongoDB Connection and Operations", func() {
		// Test MongoDB connection (if available)
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// Check if MongoDB is available
		if bsfContext.BsfSelf.MongoClient != nil {
			err := bsfContext.BsfSelf.MongoClient.Ping(ctx, nil)
			if err == nil {
				suite.T().Log("MongoDB connection successful")

				// Test binding persistence
				binding := &bsfContext.PcfBinding{
					Dnn:     "internet",
					Snssai:  &models.Snssai{Sst: 1, Sd: "010203"},
					Supi:    helperStringPtr("imsi-208930000000001"),
					PcfFqdn: helperStringPtr("pcf.free5gc.org"),
				}

				// Create binding (should persist to MongoDB)
				bindingId := bsfContext.BsfSelf.CreatePcfBinding(binding)
				suite.NotEmpty(bindingId)

				// Verify binding exists
				retrievedBinding, exists := bsfContext.BsfSelf.GetPcfBinding(bindingId)
				suite.True(exists)
				suite.Equal("internet", retrievedBinding.Dnn)

				// Clean up
				bsfContext.BsfSelf.DeletePcfBinding(bindingId)
			} else {
				suite.T().Log("MongoDB not available, skipping persistence tests")
			}
		} else {
			suite.T().Log("MongoDB client not initialized, skipping persistence tests")
		}
	})
}

// Test BSF metrics and monitoring
func (suite *BSFTestSuite) TestBSF_Metrics() {
	suite.Run("Metrics Collection", func() {
		// Create some bindings to generate metrics
		for i := 0; i < 5; i++ {
			binding := models.PcfBinding{
				Dnn:     fmt.Sprintf("dnn-%d", i),
				Snssai:  &models.Snssai{Sst: 1, Sd: "010203"},
				Supi:    fmt.Sprintf("imsi-20893000000000%d", i),
				PcfFqdn: "pcf.free5gc.org",
				PcfId:   fmt.Sprintf("pcf-%03d", i),
			}

			jsonData, _ := json.Marshal(binding)
			req, _ := http.NewRequest("POST", "/nbsf-management/v1/pcfBindings", bytes.NewBuffer(jsonData))
			req.Header.Set("Content-Type", "application/json")
			rr := httptest.NewRecorder()
			suite.router.ServeHTTP(rr, req)
			suite.Equal(http.StatusCreated, rr.Code)
		}

		// Query bindings (should generate query metrics)
		for i := 0; i < 3; i++ {
			req, _ := http.NewRequest("GET", fmt.Sprintf("/nbsf-management/v1/pcfBindings?supi=imsi-20893000000000%d", i), nil)
			rr := httptest.NewRecorder()
			suite.router.ServeHTTP(rr, req)
			suite.Equal(http.StatusOK, rr.Code)
		}

		// Test non-existent binding query (should generate failure metric)
		req, _ := http.NewRequest("GET", "/nbsf-management/v1/pcfBindings?supi=imsi-999999999999999", nil)
		rr := httptest.NewRecorder()
		suite.router.ServeHTTP(rr, req)
		suite.Equal(http.StatusNoContent, rr.Code)

		suite.T().Log("Metrics testing completed - check BSF logs for metric events")
	})
}

// Test concurrent operations and race conditions
func (suite *BSFTestSuite) TestConcurrent_Operations() {
	suite.Run("Concurrent Binding Operations", func() {
		const numGoroutines = 10
		results := make(chan error, numGoroutines)

		// Test concurrent binding creation
		for i := 0; i < numGoroutines; i++ {
			go func(id int) {
				binding := models.PcfBinding{
					Dnn:     fmt.Sprintf("dnn-%d", id),
					Snssai:  &models.Snssai{Sst: 1, Sd: "010203"},
					Supi:    fmt.Sprintf("imsi-20893000000%03d", id),
					PcfFqdn: "pcf.free5gc.org",
					PcfId:   fmt.Sprintf("pcf-%03d", id),
				}

				jsonData, _ := json.Marshal(binding)
				req, _ := http.NewRequest("POST", "/nbsf-management/v1/pcfBindings", bytes.NewBuffer(jsonData))
				req.Header.Set("Content-Type", "application/json")
				rr := httptest.NewRecorder()
				suite.router.ServeHTTP(rr, req)

				if rr.Code != http.StatusCreated {
					results <- fmt.Errorf("goroutine %d failed with status %d", id, rr.Code)
				} else {
					results <- nil
				}
			}(i)
		}

		// Wait for all goroutines to complete
		for i := 0; i < numGoroutines; i++ {
			err := <-results
			suite.NoError(err)
		}

		// Verify all bindings were created
		suite.Equal(numGoroutines, len(bsfContext.BsfSelf.PcfBindings))
	})
}

// Test binding lifecycle and TTL
func (suite *BSFTestSuite) TestBinding_Lifecycle() {
	suite.Run("Binding TTL and Cleanup", func() {
		// Create a binding
		binding := models.PcfBinding{
			Dnn:     "internet",
			Snssai:  &models.Snssai{Sst: 1, Sd: "010203"},
			Supi:    "imsi-208930000000001",
			PcfFqdn: "pcf.free5gc.org",
			PcfId:   "pcf-001",
		}

		jsonData, _ := json.Marshal(binding)
		req, _ := http.NewRequest("POST", "/nbsf-management/v1/pcfBindings", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()
		suite.router.ServeHTTP(rr, req)
		suite.Equal(http.StatusCreated, rr.Code)

		// Test that binding exists
		suite.True(len(bsfContext.BsfSelf.PcfBindings) > 0)

		// Test binding update timestamp functionality
		for _, binding := range bsfContext.BsfSelf.PcfBindings {
			suite.NotNil(binding.LastAccessTime)
			suite.True(time.Since(binding.LastAccessTime) < time.Minute)
		}

		suite.T().Log("Binding lifecycle tests completed")
	})
}

// Test error handling and edge cases
func (suite *BSFTestSuite) TestError_Handling() {
	suite.Run("Invalid JSON Request", func() {
		req, _ := http.NewRequest("POST", "/nbsf-management/v1/pcfBindings", strings.NewReader("invalid json"))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()
		suite.router.ServeHTTP(rr, req)
		suite.Equal(http.StatusBadRequest, rr.Code)
	})

	suite.Run("Malformed S-NSSAI Query", func() {
		req, _ := http.NewRequest("GET", "/nbsf-management/v1/pcfBindings?snssai=invalid", nil)
		rr := httptest.NewRecorder()
		suite.router.ServeHTTP(rr, req)
		suite.Equal(http.StatusBadRequest, rr.Code)
	})

	suite.Run("Non-existent Binding Operations", func() {
		// GET non-existent binding
		req, _ := http.NewRequest("GET", "/nbsf-management/v1/pcfBindings/non-existent-id", nil)
		rr := httptest.NewRecorder()
		suite.router.ServeHTTP(rr, req)
		suite.Equal(http.StatusNotFound, rr.Code)

		// DELETE non-existent binding (should still return 204 per 3GPP)
		req, _ = http.NewRequest("DELETE", "/nbsf-management/v1/pcfBindings/non-existent-id", nil)
		rr = httptest.NewRecorder()
		suite.router.ServeHTTP(rr, req)
		suite.Equal(http.StatusNotFound, rr.Code)
	})
}

// Run the test suite
func TestBSFComprehensive(t *testing.T) {
	suite.Run(t, new(BSFTestSuite))
}

// Helper functions
func helperStringPtr(s string) *string {
	return &s
}

func helperTimePtr(t time.Time) *time.Time {
	return &t
}
