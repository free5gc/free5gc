/*
 * BSF Unit Tests
 */

package processor_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	bsfContext "github.com/free5gc/bsf/internal/context"
	"github.com/free5gc/bsf/internal/sbi/processor"
	"github.com/free5gc/openapi/models"
)

func TestCreatePCFBinding(t *testing.T) {
	// Initialize BSF context
	bsfContext.InitBsfContext()

	// Set up Gin router
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.POST("/pcfBindings", processor.CreatePCFBinding)

	// Create test request
	pcfBinding := models.PcfBinding{
		Dnn:     "internet",
		Snssai:  &models.Snssai{Sst: 1, Sd: "010203"},
		Supi:    "imsi-208930000000001",
		PcfFqdn: "pcf.free5gc.org",
		PcfIpEndPoints: []models.IpEndPoint{
			{
				Ipv4Address: "127.0.0.1",
				Port:        8000,
			},
		},
	}

	jsonData, err := json.Marshal(pcfBinding)
	assert.NoError(t, err)

	// Create HTTP request
	req, err := http.NewRequestWithContext(context.Background(), "POST", "/pcfBindings", bytes.NewBuffer(jsonData))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	// Create response recorder
	rr := httptest.NewRecorder()

	// Execute request
	router.ServeHTTP(rr, req)

	// Assert response
	assert.Equal(t, http.StatusCreated, rr.Code)

	var response models.PcfBinding
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "internet", response.Dnn)
	assert.Equal(t, "imsi-208930000000001", response.Supi)

	// Check Location header contains binding ID
	location := rr.Header().Get("Location")
	assert.Contains(t, location, "/nbsf-management/v1/pcfBindings/")
}

func TestGetPCFBindings(t *testing.T) {
	// Initialize BSF context
	bsfContext.InitBsfContext()

	// Create a test binding first
	binding := &bsfContext.PcfBinding{
		Dnn:     "internet",
		Snssai:  &models.Snssai{Sst: 1, Sd: "010203"},
		Supi:    stringPtr("imsi-208930000000001"),
		PcfFqdn: stringPtr("pcf.free5gc.org"),
	}
	bindingId := bsfContext.BsfSelf.CreatePcfBinding(binding)

	// Set up Gin router
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.GET("/pcfBindings", processor.GetPCFBindings)

	// Create HTTP request with query parameters
	req, err := http.NewRequestWithContext(
		context.Background(),
		"GET",
		"/pcfBindings?supi=imsi-208930000000001&dnn=internet",
		nil,
	)
	assert.NoError(t, err)

	// Create response recorder
	rr := httptest.NewRecorder()

	// Execute request
	router.ServeHTTP(rr, req)

	// Assert response
	assert.Equal(t, http.StatusOK, rr.Code)

	var response models.PcfBinding
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "internet", response.Dnn)
	assert.Equal(t, "imsi-208930000000001", response.Supi)

	// Clean up
	bsfContext.BsfSelf.DeletePcfBinding(bindingId)
}

func TestDeleteIndPCFBinding(t *testing.T) {
	// Initialize BSF context
	bsfContext.InitBsfContext()

	// Create a test binding first
	binding := &bsfContext.PcfBinding{
		Dnn:     "internet",
		Snssai:  &models.Snssai{Sst: 1, Sd: "010203"},
		Supi:    stringPtr("imsi-208930000000001"),
		PcfFqdn: stringPtr("pcf.free5gc.org"),
	}
	bindingId := bsfContext.BsfSelf.CreatePcfBinding(binding)

	// Set up Gin router
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.DELETE("/pcfBindings/:bindingId", processor.DeleteIndPCFBinding)

	// Create HTTP request
	req, err := http.NewRequestWithContext(context.Background(), "DELETE", "/pcfBindings/"+bindingId, nil)
	assert.NoError(t, err)

	// Create response recorder
	rr := httptest.NewRecorder()

	// Execute request
	router.ServeHTTP(rr, req)

	// Assert response
	assert.Equal(t, http.StatusNoContent, rr.Code)

	// Verify binding was deleted
	_, exists := bsfContext.BsfSelf.PcfBindings[bindingId]
	assert.False(t, exists)
}

// Helper functions
func stringPtr(s string) *string {
	return &s
}
