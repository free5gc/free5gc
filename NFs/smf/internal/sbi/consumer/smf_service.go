package consumer

import (
	"context"
	"sync"

	"github.com/free5gc/openapi"
	"github.com/free5gc/openapi/models"
	"github.com/free5gc/openapi/smf/PDUSession"
	"github.com/free5gc/smf/internal/logger"
)

type nsmfService struct {
	consumer *Consumer

	PDUSessionMu sync.RWMutex

	PDUSessionClients map[string]*PDUSession.APIClient
}

func (s *nsmfService) getPDUSessionClient(uri string) *PDUSession.APIClient {
	if uri == "" {
		return nil
	}
	s.PDUSessionMu.RLock()
	client, ok := s.PDUSessionClients[uri]
	if ok {
		s.PDUSessionMu.RUnlock()
		return client
	}

	configuration := PDUSession.NewConfiguration()
	configuration.SetBasePath(uri)
	client = PDUSession.NewAPIClient(configuration)

	s.PDUSessionMu.RUnlock()
	s.PDUSessionMu.Lock()
	defer s.PDUSessionMu.Unlock()
	s.PDUSessionClients[uri] = client
	return client
}

func (s *nsmfService) SendSMContextStatusNotification(uri string) (*models.ProblemDetails, error) {
	if uri != "" {
		request := &PDUSession.PostSmContextsSmContextStatusNotificationPostRequest{
			SmfPduSessionSmContextStatusNotification: &models.SmfPduSessionSmContextStatusNotification{
				StatusInfo: &models.StatusInfo{
					ResourceStatus: models.ResourceStatus_RELEASED,
				},
			},
		}

		client := s.getPDUSessionClient(uri)

		logger.CtxLog.Infoln("[SMF] Send SMContext Status Notification")
		_, localErr := client.SMContextsCollectionApi.
			PostSmContextsSmContextStatusNotificationPost(context.Background(), uri, request)

		switch err := localErr.(type) {
		case openapi.GenericOpenAPIError:
			switch errModel := err.Model().(type) {
			case PDUSession.PostSmContextsSmContextStatusNotificationPostError:
				return &errModel.ProblemDetails, nil
			case error:
				return openapi.ProblemDetailsSystemFailure(errModel.Error()), nil
			default:
				return nil, openapi.ReportError("openapi error")
			}
		case error:
			return openapi.ProblemDetailsSystemFailure(err.Error()), nil
		case nil:
			logger.PduSessLog.Tracef("Send SMContextStatus Notification Success")
			return nil, nil
		default:
			logger.PduSessLog.Warnf("Send SMContextStatus Notification Unknown Error: %+v", err)
			return nil, openapi.ReportError("server no response")
		}
	}
	return nil, nil
}
