/*
 * NRF NSSAI Availability
 *
 * NRF NSSAI Availability Service
 */

package Nnrf_NFManagement

import (
	"context"
	"encoding/json"
	"fmt"
	"free5gc/lib/http2_util"
	. "free5gc/lib/openapi/models"
	"free5gc/lib/path_util"
	"free5gc/src/nssf/test"
	"log"
	"net/http"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

var testingNotification = test.TestingNssaiavailability{
	ConfigFile:             test.ConfigFileFromArgs,
	NfNssaiAvailabilityUri: "https://localhost:29510/notification",
}

func generateNotificationRequest() NotificationData {
	const jsonRequest = `
        {
            "event": "NF_REGISTERED",
            "nfInstanceUri" : "127.0.0.1:123456"
        }
    `

	var n NotificationData
	if err := json.NewDecoder(strings.NewReader(jsonRequest)).Decode(&n); err != nil {
		fmt.Printf("Decode error: %v", err)
	}

	return n
}

// Path of HTTP2 key and log file
var (
	NRF_LOG_PATH = path_util.Gofree5gcPath("free5gc/src/nrf/Management/sslkeylog.log")
	NRF_PEM_PATH = path_util.Gofree5gcPath("free5gc/support/TLS/nrf.pem")
	NRF_KEY_PATH = path_util.Gofree5gcPath("free5gc/support/TLS/nrf.key")
)

func TestNotificationPost(t *testing.T) {
	var (
		requestBody string
	)

	// Create a server to accept testing requests
	router := gin.Default()
	router.POST("/notification", func(c *gin.Context) {
		/*buf, err := c.GetRawData()
		if err != nil {
			t.Errorf(err.Error())
		}
		// Remove NL line feed, new line character
		//requestBody = string(buf[:len(buf)-1])*/
		var ND NotificationData

		if err := c.ShouldBindJSON(&ND); err != nil {
			log.Panic(err.Error())
		}
		fmt.Println(ND)
		c.JSON(http.StatusNoContent, gin.H{})
	})

	srv, err := http2_util.NewServer(":29510", NRF_LOG_PATH, router)
	if err != nil {
		t.Fatal(err)
	}

	go func() {
		err := srv.ListenAndServeTLS(NRF_PEM_PATH, NRF_KEY_PATH)
		if err != nil && err != http.ErrServerClosed {
			t.Fatal(err)
		}
	}()

	configuration := NewConfiguration()
	configuration.SetBasePathNoGroup(testingNotification.NfNssaiAvailabilityUri)
	apiClient := NewAPIClient(configuration)

	subtests := []struct {
		name                string
		generateRequestBody func() NotificationData
		expectRequestBody   string
	}{
		{
			name:                "Notify",
			generateRequestBody: generateNotificationRequest,
			expectRequestBody:   ``,
		},
	}

	for _, subtest := range subtests {
		t.Run(subtest.name, func(t *testing.T) {
			var (
				n    NotificationData
				resp *http.Response
			)

			// Start to generate and send notification request after channel is closed
			if subtest.generateRequestBody != nil {
				n = subtest.generateRequestBody()
			}

			resp, err := apiClient.NotificationApi.NotificationPost(context.Background(), n)

			if err != nil {
				t.Errorf(err.Error())
			}

			if resp.StatusCode != http.StatusNoContent {
				t.Errorf("Incorrect status code: expected %d, got %d", http.StatusNoContent, resp.StatusCode)
			}

			if requestBody != subtest.expectRequestBody {
				t.Errorf("Incorrect request body:\nexpected\n%s\n, got\n%s", subtest.expectRequestBody, requestBody)
			}

			err = srv.Shutdown(context.Background())
			if err != nil {
				t.Fatal(err)
			}
		})
	}
}
