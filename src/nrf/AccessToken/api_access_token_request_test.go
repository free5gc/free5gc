package AccessToken_test

import (
	"context"
	"crypto/tls"
	"github.com/antihax/optional"
	"free5gc/lib/MongoDBLibrary"
	"free5gc/lib/Nnrf_AccessToken"
	"free5gc/lib/openapi/models"
	"free5gc/src/nrf/AccessToken"
	"free5gc/src/nrf/logger"
	"free5gc/src/nrf/nrf_handler"
	"free5gc/src/nrf/nrf_util"
	"net/http"
	"os"
	"testing"
	"time"
	//"github.com/stretchr/testify/assert"
)

func TestAccessTokenRequest(t *testing.T) {
	// run AccessToken Server Routine
	go func() {
		kl, _ := os.OpenFile(nrf_util.NrfLogPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
		router := AccessToken.NewRouter()

		server := http.Server{
			Addr: "127.0.0.1:29510",
			TLSConfig: &tls.Config{
				KeyLogWriter: kl,
			},

			Handler: router,
		}
		_ = server.ListenAndServeTLS(nrf_util.NrfPemPath, nrf_util.NrfKeyPath)

	}()
	time.Sleep(time.Duration(2) * time.Second)

	go nrf_handler.Handle()

	// connect to mongoDB
	MongoDBLibrary.SetMongoDB("free5gc", "mongodb://140.113.214.205:30030")

	// Set client and set url
	configuration := Nnrf_AccessToken.NewConfiguration()
	configuration.SetBasePath("https://127.0.0.1:29510")
	client := Nnrf_AccessToken.NewAPIClient(configuration)

	// Set test data (with expected data)
	grantType := "client_credentials"
	nfInstanceId := "0" // nfInstanceId of service consumer
	scope := "nnrf-nfm"
	localVarOptionals := Nnrf_AccessToken.AccessTokenRequestParamOpts{
		NfType:             optional.NewInterface(models.NfType_NRF),                     // nfType of service consumer
		TargetNfType:       optional.NewInterface(models.NfType_NRF),                     // nfType of service producer
		TargetNfInstanceId: optional.NewInterface("2"),                                   // nfInstanceId of service producer
		RequesterPlmn:      optional.NewInterface("{\"mcc\": \"111\",\"mnc\": \"111\"}"), // plmn of service consumer
		TargetPlmn:         optional.NewInterface("{\"mcc\": \"111\",\"mnc\": \"111\"}"), // plmn of service producer
	}

	// Check test data (Use RESTful GET)
	rep, res, err := client.AccessTokenRequestApi.AccessTokenRequest(context.TODO(), grantType, nfInstanceId, scope, &localVarOptionals)
	if err != nil {
		logger.AppLog.Panic(err)
	}
	if res != nil {
		if status := res.StatusCode; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusOK)
		}
	}
	t.Logf("%+v", rep)
}
