package consumer

import (
	"sync"

	amf_context "github.com/free5gc/amf/internal/context"
	"github.com/free5gc/openapi"
	"github.com/free5gc/openapi/models"
	Nnssf_NSSelection "github.com/free5gc/openapi/nssf/NSSelection"
)

type nssfService struct {
	consumer *Consumer

	NSSelectionMu sync.RWMutex

	NSSelectionClients map[string]*Nnssf_NSSelection.APIClient
}

func (s *nssfService) getNSSelectionClient(uri string) *Nnssf_NSSelection.APIClient {
	if uri == "" {
		return nil
	}
	s.NSSelectionMu.RLock()
	client, ok := s.NSSelectionClients[uri]
	if ok {
		s.NSSelectionMu.RUnlock()
		return client
	}

	configuration := Nnssf_NSSelection.NewConfiguration()
	configuration.SetBasePath(uri)
	client = Nnssf_NSSelection.NewAPIClient(configuration)

	s.NSSelectionMu.RUnlock()
	s.NSSelectionMu.Lock()
	defer s.NSSelectionMu.Unlock()
	s.NSSelectionClients[uri] = client
	return client
}

func (s *nssfService) NSSelectionGetForRegistration(ue *amf_context.AmfUe, requestedNssai []models.MappingOfSnssai) (
	*models.ProblemDetails, error,
) {
	client := s.getNSSelectionClient(ue.NssfUri)
	if client == nil {
		return nil, openapi.ReportError("nssf not found")
	}

	amfSelf := amf_context.GetSelf()
	ctx, _, err := amf_context.GetSelf().GetTokenCtx(models.ServiceName_NNSSF_NSSELECTION,
		models.NrfNfManagementNfType_NSSF)
	if err != nil {
		return nil, err
	}
	sliceInfo := models.SliceInfoForRegistration{
		SubscribedNssai: ue.SubscribedNssai,
	}

	for _, snssai := range requestedNssai {
		sliceInfo.RequestedNssai = append(sliceInfo.RequestedNssai, *snssai.ServingSnssai)
		if snssai.HomeSnssai != nil {
			sliceInfo.MappingOfNssai = append(sliceInfo.MappingOfNssai, snssai)
		}
	}

	var paramOpt Nnssf_NSSelection.NSSelectionGetRequest

	testNfType := models.NrfNfManagementNfType_AMF

	paramOpt = Nnssf_NSSelection.NSSelectionGetRequest{
		NfType:                          &testNfType,
		NfId:                            &amfSelf.NfId,
		SliceInfoRequestForRegistration: &sliceInfo,
		Tai:                             &ue.Tai, // TS 29.531 R15.3 6.1.3.2.3.1
	}

	res, localErr := client.NetworkSliceInformationDocumentApi.NSSelectionGet(ctx,
		&paramOpt)
	if localErr == nil {
		ue.NetworkSliceInfo = &res.AuthorizedNetworkSliceInfo
		for _, allowedNssai := range res.AuthorizedNetworkSliceInfo.AllowedNssaiList {
			ue.AllowedNssai[allowedNssai.AccessType] = allowedNssai.AllowedSnssaiList
		}
		ue.ConfiguredNssai = res.AuthorizedNetworkSliceInfo.ConfiguredNssai
	} else {
		switch apiErr := err.(type) {
		// API error
		case openapi.GenericOpenAPIError:
			switch errModel := apiErr.Model().(type) {
			case Nnssf_NSSelection.NSSelectionGetError:
				return &errModel.ProblemDetails, localErr
			case error:
				return openapi.ProblemDetailsSystemFailure(errModel.Error()), nil
			default:
				return nil, openapi.ReportError("openapi error")
			}
		case error:
			return openapi.ProblemDetailsSystemFailure(apiErr.Error()), nil
		default:
			return nil, openapi.ReportError("openapi error")
		}
	}

	return nil, nil
}

func (s *nssfService) NSSelectionGetForPduSession(ue *amf_context.AmfUe, snssai models.Snssai) (
	*models.AuthorizedNetworkSliceInfo, *models.ProblemDetails, error,
) {
	client := s.getNSSelectionClient(ue.NssfUri)
	if client == nil {
		return nil, nil, openapi.ReportError("nssf not found")
	}

	amfSelf := amf_context.GetSelf()
	sliceInfoForPduSession := models.SliceInfoForPduSession{
		SNssai:            &snssai,
		RoamingIndication: models.RoamingIndication_NON_ROAMING, // not support roaming
	}

	testNfType := models.NrfNfManagementNfType_AMF

	paramOpt := Nnssf_NSSelection.NSSelectionGetRequest{
		NfType:                        &testNfType,
		NfId:                          &amfSelf.NfId,
		SliceInfoRequestForPduSession: &sliceInfoForPduSession,
		Tai:                           &ue.Tai, // TS 29.531 R15.3 6.1.3.2.3.1
	}

	ctx, _, err := amf_context.GetSelf().GetTokenCtx(models.ServiceName_NNSSF_NSSELECTION,
		models.NrfNfManagementNfType_NSSF)
	if err != nil {
		return nil, nil, err
	}
	res, localErr := client.NetworkSliceInformationDocumentApi.NSSelectionGet(ctx, &paramOpt)

	if localErr == nil {
		return &res.AuthorizedNetworkSliceInfo, nil, nil
	} else {
		switch apiErr := localErr.(type) {
		// API error
		case openapi.GenericOpenAPIError:
			switch errModel := apiErr.Model().(type) {
			case Nnssf_NSSelection.NSSelectionGetError:
				return nil, &errModel.ProblemDetails, localErr
			case error:
				return nil, openapi.ProblemDetailsSystemFailure(errModel.Error()), nil
			default:
				return nil, nil, openapi.ReportError("openapi error")
			}
		case error:
			return nil, openapi.ProblemDetailsSystemFailure(apiErr.Error()), nil
		default:
			return nil, nil, openapi.ReportError("server no response")
		}
	}
}
