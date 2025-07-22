package util

import (
	"net/http"

	"github.com/free5gc/openapi/models"
)

// Title in ProblemDetails for UDR HTTP APIs
const (
	INVALID_REQUEST       = "Invalid request message framing"
	MALFORMED_REQUEST     = "Malformed request syntax"
	UNAUTHORIZED_CONSUMER = "Unauthorized NF service consumer"
	UNSUPPORTED_RESOURCE  = "Unsupported request resources"
)

func ProblemDetailsSystemFailure(detail string) *models.ProblemDetails {
	return &models.ProblemDetails{
		Title:  "System failure",
		Status: http.StatusInternalServerError,
		Detail: detail,
		Cause:  "SYSTEM_FAILURE",
	}
}

func ProblemDetailsMalformedReqSyntax(detail string) *models.ProblemDetails {
	return &models.ProblemDetails{
		Title:  "Malformed request syntax",
		Status: http.StatusBadRequest,
		Detail: detail,
	}
}

func ProblemDetailsNotFound(cause string) *models.ProblemDetails {
	title := ""
	switch cause {
	case "USER_NOT_FOUND":
		title = "User not found"
	case "SUBSCRIPTION_NOT_FOUND":
		title = "Subscription not found"
	case "AMFSUBSCRIPTION_NOT_FOUND":
		title = "AMF Subscription not found"
	default:
		title = "Data not found"
	}
	return &models.ProblemDetails{
		Title:  title,
		Status: http.StatusNotFound,
		Cause:  cause,
	}
}

func ProblemDetailsModifyNotAllowed(detail string) *models.ProblemDetails {
	return &models.ProblemDetails{
		Title:  "Modify not allowed",
		Status: http.StatusForbidden,
		Cause:  "MODIFY_NOT_ALLOWED",
		Detail: detail,
	}
}

func ProblemDetailsUpspecified(detail string) *models.ProblemDetails {
	return &models.ProblemDetails{
		Title:  "Unspecified",
		Status: http.StatusForbidden,
		Cause:  "UNSPECIFIED",
		Detail: detail,
	}
}
