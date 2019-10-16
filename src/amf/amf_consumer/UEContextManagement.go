package amf_consumer

import (
	"context"
	"free5gc/lib/Nudm_UEContextManagement"
	"free5gc/lib/openapi/common"
	"free5gc/lib/openapi/models"
	"free5gc/src/amf/amf_context"
)

func UeCmRegistration(ue *amf_context.AmfUe, accessType models.AccessType, initialRegistrationInd bool) (problemDetails *models.ProblemDetails, err error) {

	configuration := Nudm_UEContextManagement.NewConfiguration()
	configuration.SetBasePath(ue.NudmUECMUri)
	client := Nudm_UEContextManagement.NewAPIClient(configuration)

	amfSelf := amf_context.AMF_Self()

	switch accessType {
	case models.AccessType__3_GPP_ACCESS:
		registrationData := models.Amf3GppAccessRegistration{
			AmfInstanceId:          amfSelf.NfId,
			InitialRegistrationInd: initialRegistrationInd,
			Guami:                  &amfSelf.ServedGuamiList[0],
			RatType:                ue.RatType,
			ImsVoPs:                models.ImsVoPs_HOMOGENEOUS_NON_SUPPORT, // TODO: not support Homogenous Support of IMS Voice over PS Sessions this stage
		}

		_, httpResp, localErr := client.AMFRegistrationFor3GPPAccessApi.Registration(context.Background(), ue.Supi, registrationData)
		if localErr == nil {
			return
		} else if httpResp != nil {
			if httpResp.Status != localErr.Error() {
				err = localErr
				return
			}
			problem := localErr.(common.GenericOpenAPIError).Model().(models.ProblemDetails)
			problemDetails = &problem
		} else {
			err = common.ReportError("server no response")
			return
		}
	case models.AccessType_NON_3_GPP_ACCESS:
		registrationData := models.AmfNon3GppAccessRegistration{
			AmfInstanceId: amfSelf.NfId,
			Guami:         &amfSelf.ServedGuamiList[0],
			RatType:       ue.RatType,
		}

		_, httpResp, localErr := client.AMFRegistrationForNon3GPPAccessApi.Register(context.Background(), ue.Supi, registrationData)
		if localErr == nil {
			return
		} else if httpResp != nil {
			if httpResp.Status != localErr.Error() {
				err = localErr
				return
			}
			problem := localErr.(common.GenericOpenAPIError).Model().(models.ProblemDetails)
			problemDetails = &problem
		} else {
			err = common.ReportError("server no response")
			return
		}
	}
	return
}
