/*
 * NSSF NSSAI Availability
 *
 * NSSF NSSAI Availability Service
 */

package NSSAIAvailability

import (
	"context"
	"net/http"
	"testing"

	"free5gc/lib/Nnssf_NSSAIAvailability"
	"free5gc/lib/http2_util"
	"free5gc/src/nssf/factory"
	"free5gc/src/nssf/nssf_handler"
	"free5gc/src/nssf/test"
	"free5gc/src/nssf/util"
)

var testingNssaiavailabilityUnsubscribeApi = test.TestingNssaiavailability{
	ConfigFile: test.ConfigFileFromArgs,
}

func TestNSSAIAvailabilityUnsubscribe(t *testing.T) {
	factory.InitConfigFactory(testingNssaiavailabilityUnsubscribeApi.ConfigFile)

	router := NewRouter()
	srv, err := http2_util.NewServer(":29531", util.NSSF_LOG_PATH, router)
	if err != nil {
		t.Fatal(err)
	}

	go nssf_handler.Handle()

	go func() {
		err := srv.ListenAndServeTLS(util.NSSF_PEM_PATH, util.NSSF_KEY_PATH)
		if err != nil && err != http.ErrServerClosed {
			t.Fatal(err)
		}
	}()

	configuration := Nnssf_NSSAIAvailability.NewConfiguration()
	configuration.SetBasePath("https://localhost:29531")
	apiClient := Nnssf_NSSAIAvailability.NewAPIClient(configuration)

	subtests := []struct {
		name           string
		subscriptionId string
		expectStatus   int
	}{
		{
			name:           "Delete",
			subscriptionId: "3",
			expectStatus:   http.StatusNoContent,
		},
	}

	for _, subtest := range subtests {
		t.Run(subtest.name, func(t *testing.T) {
			var (
				resp *http.Response
			)

			resp, err := apiClient.SubscriptionIDDocumentApi.NSSAIAvailabilityUnsubscribe(context.Background(), subtest.subscriptionId)

			if err != nil {
				t.Errorf(err.Error())
			}

			if resp.StatusCode != subtest.expectStatus {
				t.Errorf("Incorrect status code: expected %d, got %d", subtest.expectStatus, resp.StatusCode)
			}
		})
	}

	err = srv.Shutdown(context.Background())
	if err != nil {
		t.Fatal(err)
	}
}
