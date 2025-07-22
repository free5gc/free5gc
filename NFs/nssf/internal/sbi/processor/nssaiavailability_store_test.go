package processor_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"go.uber.org/mock/gomock"

	"github.com/free5gc/nssf/internal/sbi/processor"
	"github.com/free5gc/nssf/internal/util"
	"github.com/free5gc/nssf/pkg/app"
	"github.com/free5gc/nssf/pkg/factory"
	"github.com/free5gc/openapi/models"
)

func setup() {
	// Set the default values for the factory.NssfConfig
	factory.NssfConfig = &factory.Config{
		Configuration: &factory.Configuration{},
	}
}

func TestMain(m *testing.M) {
	// Run the tests
	setup()
	m.Run()
}

func TestNfInstanceDelete(t *testing.T) {
	mockNssfApp := app.NewMockNssfApp(gomock.NewController(t))
	processor := processor.NewProcessor(mockNssfApp)
	httpRecorder := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(httpRecorder)

	// Create a sample AMF list
	amfList := []factory.AmfConfig{
		{
			NfId: "nf1",
		},
		{
			NfId: "nf2",
		},
		{
			NfId: "nf3",
		},
	}

	// Set the sample AMF list in the factory.NssfConfig.Configuration
	factory.NssfConfig.Configuration.AmfList = amfList

	// Test case 1: Delete an existing NF instance
	nfIdToDelete := "nf2"
	processor.NssaiAvailabilityNfInstanceDelete(c, nfIdToDelete)
	if c.Writer.Status() != http.StatusNoContent {
		t.Errorf("Expected status code %d, got: %d", http.StatusNoContent, httpRecorder.Code)
	}

	// Verify that the NF instance is deleted from the AMF list
	for _, amfConfig := range factory.NssfConfig.Configuration.AmfList {
		if amfConfig.NfId == nfIdToDelete {
			t.Errorf("Expected NF instance '%s' to be deleted, but it still exists", nfIdToDelete)
		}
	}

	// Test case 2: Delete a non-existing NF instance
	nfIdToDelete = "nf4"
	expectedDetail := fmt.Sprintf("AMF ID '%s' does not exist", nfIdToDelete)
	processor.NssaiAvailabilityNfInstanceDelete(c, nfIdToDelete)
	if httpRecorder.Code != http.StatusNotFound {
		t.Errorf("Expected status code %d, got: %d", http.StatusNotFound, httpRecorder.Code)
	}

	var problemDetails models.ProblemDetails
	if err := json.Unmarshal(httpRecorder.Body.Bytes(), &problemDetails); err != nil {
		t.Errorf("Error unmarshalling response body: %v", err)
	}
	if problemDetails.Title != util.UNSUPPORTED_RESOURCE {
		t.Errorf("Expected problemDetails.Title to be '%s', got: '%s'", util.UNSUPPORTED_RESOURCE, problemDetails.Title)
	}
	if problemDetails.Status != http.StatusNotFound {
		t.Errorf("Expected problemDetails.Status to be %d, got: %d", http.StatusNotFound, problemDetails.Status)
	}
	if problemDetails.Detail != expectedDetail {
		t.Errorf("Expected problemDetails.Detail to be '%s', got: '%s'", expectedDetail, problemDetails.Detail)
	}
}
