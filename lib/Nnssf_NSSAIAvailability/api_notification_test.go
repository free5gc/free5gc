/*
 * NSSF NSSAI Availability
 *
 * NSSF NSSAI Availability Service
 */

package Nnssf_NSSAIAvailability

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"

	"free5gc/lib/http2_util"
	. "free5gc/lib/openapi/models"
	"free5gc/src/nssf/test"
	"free5gc/src/nssf/util"
)

var testingNotification = test.TestingNssaiavailability{
	ConfigFile:             test.ConfigFileFromArgs,
	NfNssaiAvailabilityUri: "https://localhost:29531/notification",
}

func generateNotificationRequest() NssfEventNotification {
	const jsonRequest = `
        {
            "subscriptionId": "1",
            "authorizedNssaiAvailabilityData": [
                {
                    "tai": {
                        "plmnId": {
                            "mcc": "466",
                            "mnc": "92"
                        },
                        "tac": "33456"
                    },
                    "supportedSnssaiList": [
                        {
                            "sst": 1
                        },
                        {
                            "sst": 1,
                            "sd": "1"
                        },
                        {
                            "sst": 1,
                            "sd": "2"
                        },
                        {
                            "sst": 2
                        }
                    ]
                }
            ]
        }
    `

	var n NssfEventNotification
	if err := json.NewDecoder(strings.NewReader(jsonRequest)).Decode(&n); err != nil {
		fmt.Printf("Decode error: %v", err)
	}

	return n
}

func TestNotificationPost(t *testing.T) {
	var (
		requestBody string
	)

	// Create a server to accept testing requests
	router := gin.Default()
	router.POST("/notification", func(c *gin.Context) {
		buf, err := c.GetRawData()
		if err != nil {
			t.Errorf(err.Error())
		}
		// Remove NL line feed, new line character
		requestBody = string(buf[:len(buf)-1])

		c.JSON(http.StatusNoContent, gin.H{})
	})

	srv, err := http2_util.NewServer(":29531", util.NSSF_LOG_PATH, router)
	if err != nil {
		t.Fatal(err)
	}

	go func() {
		err := srv.ListenAndServeTLS(util.NSSF_PEM_PATH, util.NSSF_KEY_PATH)
		if err != nil && err != http.ErrServerClosed {
			t.Fatal(err)
		}
	}()

	configuration := NewConfiguration()
	configuration.SetBasePathNoGroup(testingNotification.NfNssaiAvailabilityUri)
	apiClient := NewAPIClient(configuration)

	subtests := []struct {
		name                string
		generateRequestBody func() NssfEventNotification
		expectRequestBody   string
	}{
		{
			name:                "Notify",
			generateRequestBody: generateNotificationRequest,
			expectRequestBody: `{"subscriptionId":"1","authorizedNssaiAvailabilityData":[{` +
				`"tai":{"plmnId":{"mcc":"466","mnc":"92"},"tac":"33456"},` +
				`"supportedSnssaiList":[{"sst":1},{"sst":1,"sd":"1"},{"sst":1,"sd":"2"},{"sst":2}]}]}`,
		},
	}

	for _, subtest := range subtests {
		t.Run(subtest.name, func(t *testing.T) {
			var (
				n    NssfEventNotification
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
		})
	}

	err = srv.Shutdown(context.Background())
	if err != nil {
		t.Fatal(err)
	}
}
