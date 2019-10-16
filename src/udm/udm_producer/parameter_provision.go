package udm_producer

import (
	"context"
	"free5gc/lib/openapi/common"
	"free5gc/lib/openapi/models"
	m "free5gc/lib/openapi/models"
	"free5gc/src/udm/udm_handler/udm_message"
	"net/http"
)

func HandleUpdate(httpChannel chan udm_message.HandlerResponseMessage, gpsi string, body m.PpData) {

	clientAPI := createUDMClientToUDR(gpsi, false)
	res, err := clientAPI.ProvisionedParameterDataDocumentApi.ModifyPpData(context.Background(), gpsi, nil)
	if err != nil {
		var problemDetails m.ProblemDetails
		problemDetails.Cause = err.(common.GenericOpenAPIError).Model().(models.ProblemDetails).Cause
		udm_message.SendHttpResponseMessage(httpChannel, nil, res.StatusCode, problemDetails)
		return
	}
	udm_message.SendHttpResponseMessage(httpChannel, nil, http.StatusNoContent, nil)
}
