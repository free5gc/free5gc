package util

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/free5gc/nrf/pkg/factory"
	"github.com/free5gc/openapi/oauth"
)

// This function would check the OAuth2 token
func AuthorizationCheck(c *gin.Context, serviceName string) error {
	if factory.NrfConfig.GetOAuth() {
		oauth_err := oauth.VerifyOAuth(c.Request.Header.Get("Authorization"), serviceName,
			factory.NrfConfig.GetNrfCertPemPath())

		if oauth_err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": oauth_err.Error()})
			return oauth_err
		}
	}
	return nil
}
