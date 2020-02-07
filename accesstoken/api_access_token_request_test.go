package accesstoken_test

import (
	"context"
	"crypto/tls"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/antihax/optional"

	"github.com/free5gc/MongoDBLibrary"
	"github.com/free5gc/nrf/accesstoken"
	"github.com/free5gc/nrf/factory"
	"github.com/free5gc/nrf/logger"
	"github.com/free5gc/nrf/util"
	"github.com/free5gc/openapi/Nnrf_AccessToken"
	"github.com/free5gc/openapi/models"
)

func TestAccessTokenRequest(t *testing.T) {
	// run accesstoken Server Routine
	go func() {
		kl, _ := os.OpenFile(util.NrfLogPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o600)
		router := accesstoken.NewRouter()

		server := http.Server{
			Addr: factory.NRF_DEFAULT_IPV4 + ":" + factory.NRF_DEFAULT_PORT,
			TLSConfig: &tls.Config{
				KeyLogWriter: kl,
			},

			Handler: router,
		}
		_ = server.ListenAndServeTLS(util.NrfPemPath, util.NrfKeyPath)
	}()
	time.Sleep(time.Duration(2) * time.Second)

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
		logger.AppLog.Errorln(err)
	}
	if res != nil {
		if status := res.StatusCode; status != http.StatusOK {
			logger.AppLog.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusOK)
		}
	}

	t.Logf("%+v", rep)
}
