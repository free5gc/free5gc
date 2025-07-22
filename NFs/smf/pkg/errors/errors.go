package errors

import (
	"net/http"

	"github.com/free5gc/openapi/models"
)

var (
	N1SmError = models.SmfPduSessionExtProblemDetails{
		Title:  "Invalid N1 Message",
		Status: http.StatusForbidden,
		Detail: "N1 Message Error",
		Cause:  "N1_SM_ERROR",
	}
	N2SmError = models.SmfPduSessionExtProblemDetails{
		Title:  "Invalid N2 Message",
		Status: http.StatusForbidden,
		Detail: "N2 Message Error",
		Cause:  "N2_SM_ERROR",
	}
	DnnDeniedError = models.SmfPduSessionExtProblemDetails{
		Title:         "DNN Denied",
		Status:        http.StatusForbidden,
		Detail:        "The subscriber does not have the necessary subscription to access the DNN",
		Cause:         "DNN_DENIED",
		InvalidParams: nil,
	}
	DnnNotSupported = models.SmfPduSessionExtProblemDetails{
		Title:         "DNN Not Supported",
		Status:        http.StatusForbidden,
		Detail:        "The DNN is not supported by the SMF.",
		Cause:         "DNN_NOT_SUPPORTED",
		InvalidParams: nil,
	}
	InsufficientResourceSliceDnn = models.SmfPduSessionExtProblemDetails{
		Title:         "DNN Resource insufficient",
		Status:        http.StatusInternalServerError,
		Detail:        "The request cannot be provided due to insufficient resources for the specific slice and DNN.",
		Cause:         "INSUFFICIENT_RESOURCES_SLICE_DNN",
		InvalidParams: nil,
	}
	SubscriptionDenied = models.SmfPduSessionExtProblemDetails{
		Title:  "Subscription Denied",
		Status: http.StatusForbidden,
		Detail: "This indicates an error, other than those listed in this table, " +
			"due to lack of necessary subscription to serve the UE request.",
		Cause:         "SUBSCRIPTION_DENIED",
		InvalidParams: nil,
	}
	NetworkFailure = models.SmfPduSessionExtProblemDetails{
		Title:         "Network failure",
		Status:        http.StatusGatewayTimeout,
		Detail:        "The request is rejected due to a network problem.",
		Cause:         "NETWORK_FAILURE",
		InvalidParams: nil,
	}
	SmContextStateMismatchActive = models.SmfPduSessionExtProblemDetails{
		Title:  "SMContext state mismatch",
		Status: http.StatusForbidden,
		Detail: "The SMContext State should be Active State.",
	}
	SmContextStateMismatchInActive = models.SmfPduSessionExtProblemDetails{
		Title:  "SMContext state mismatch",
		Status: http.StatusForbidden,
		Detail: "The SMContext State should be InActive State.",
	}
)
