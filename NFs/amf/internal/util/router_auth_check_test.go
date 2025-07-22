package util_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"

	"github.com/free5gc/amf/internal/util"
	"github.com/free5gc/openapi"
	"github.com/free5gc/openapi/models"
)

const (
	Valid   = "valid"
	Invalid = "invalid"
)

type mockAMFContext struct{}

func newMockAMFContext() *mockAMFContext {
	return &mockAMFContext{}
}

func (m *mockAMFContext) AuthorizationCheck(token string, serviceName models.ServiceName) error {
	if token == Valid {
		return nil
	}

	return errors.New("invalid token")
}

func TestRouterAuthorizationCheck_Check(t *testing.T) {
	// Mock gin.Context
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	var err error
	c.Request, err = http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Errorf("error on http request: %+v", err)
	}

	type Args struct {
		token string
	}
	type Want struct {
		statusCode int
	}

	tests := []struct {
		name string
		args Args
		want Want
	}{
		{
			name: "Valid Token",
			args: Args{
				token: Valid,
			},
			want: Want{
				statusCode: http.StatusOK,
			},
		},
		{
			name: "Invalid Token",
			args: Args{
				token: Invalid,
			},
			want: Want{
				statusCode: http.StatusUnauthorized,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w = httptest.NewRecorder()
			c, _ = gin.CreateTestContext(w)
			c.Request, err = http.NewRequest("GET", "/", nil)
			if err != nil {
				t.Errorf("error on http request: %+v", err)
			}
			c.Request.Header.Set("Authorization", tt.args.token)

			var testService models.ServiceName = "testService"

			rac := util.NewRouterAuthorizationCheck(testService)
			rac.Check(c, newMockAMFContext())
			if w.Code != tt.want.statusCode {
				t.Errorf("StatusCode should be %d, but got %d", tt.want.statusCode, w.Code)
			}
		})
	}
}

// Test for smContextStatusNotify
type mockProcessor struct {
	called bool
}

func (m *mockProcessor) HandleSmContextStatusNotify(c *gin.Context,
	notif models.SmfPduSessionSmContextStatusNotification,
) {
	m.called = true
	c.JSON(http.StatusNoContent, nil)
}

type mockServer struct {
	processor *mockProcessor
}

func (s *mockServer) Processor() *mockProcessor {
	return s.processor
}

func (s *mockServer) HTTPSmContextStatusNotify(c *gin.Context) {
	var smContextStatusNotification models.SmfPduSessionSmContextStatusNotification

	requestBody, err := c.GetRawData()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	err = openapi.Deserialize(&smContextStatusNotification, requestBody, "application/json")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	s.Processor().HandleSmContextStatusNotify(c, smContextStatusNotification)
}

func TestHTTPSmContextStatusNotify(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	processor := &mockProcessor{}
	server := &mockServer{processor: processor}

	router.POST("/sm-context-status-notify", server.HTTPSmContextStatusNotify)

	jsonBody := `{
		"pduSessionId": "1",
		"statusInfo": {
			"cause": "RELEASED_DUE_TO_5GSM_CAUSE"
		}
	}`

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	var err error
	c.Request, err = http.NewRequest("POST", "/sm-context-status-notify", strings.NewReader(jsonBody))
	if err != nil {
		t.Errorf("error on http request: %+v", err)
	}
	c.Request.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, c.Request)

	if w.Code != http.StatusNoContent {
		t.Errorf("StatusCode should be %d, but got %d", http.StatusNoContent, w.Code)
	}
	if !processor.called {
		t.Errorf("Expected HandleSmContextStatusNotify to be called")
	}
}
