package util

import (
	"net/http"

	"github.com/free5gc/openapi/models"
)

const (
	UdmDefaultKeyLogPath = "./log/udmsslkey.log"
	UdmDefaultPemPath    = "./config/TLS/udm.pem"
	UdmDefaultKeyPath    = "./config/TLS/udm.key"
	UdmDefaultConfigPath = "./config/udmcfg.yaml"
)

func ProblemDetailsSystemFailure(detail string) *models.ProblemDetails {
	return &models.ProblemDetails{
		Title:  "System failure",
		Status: http.StatusInternalServerError,
		Detail: detail,
		Cause:  "SYSTEM_FAILURE",
	}
}
