package consumer

import (
	"net/http"
	"strings"
	"sync"

	"github.com/antihax/optional"
	"github.com/free5gc/nef/internal/logger"
	"github.com/free5gc/openapi/Npcf_PolicyAuthorization"
	"github.com/free5gc/openapi/models"
)

type npcfService struct {
	consumer *Consumer

	mu      sync.RWMutex
	clients map[string]*Npcf_PolicyAuthorization.APIClient
}

func (s *npcfService) getClient(uri string) *Npcf_PolicyAuthorization.APIClient {
	s.mu.RLock()
	if client, ok := s.clients[uri]; ok {
		defer s.mu.RUnlock()
		return client
	} else {
		configuration := Npcf_PolicyAuthorization.NewConfiguration()
		configuration.SetBasePath(uri)
		cli := Npcf_PolicyAuthorization.NewAPIClient(configuration)

		s.mu.RUnlock()
		s.mu.Lock()
		defer s.mu.Unlock()
		s.clients[uri] = cli
		return cli
	}
}

func (s *npcfService) getPcfPolicyAuthUri() (string, error) {
	uri := s.consumer.Context().PcfPaUri()
	if uri == "" {
		_, sUri, err := s.consumer.SearchNFInstances(s.consumer.Config().NrfUri(),
			models.ServiceName_NPCF_POLICYAUTHORIZATION, nil)
		if err == nil {
			s.consumer.Context().SetPcfPaUri(sUri)
		}
		return sUri, err
	}
	return uri, nil
}

func (s *npcfService) GetAppSession(appSessionId string) (int, interface{}) {
	var (
		err     error
		rspCode int
		rspBody interface{}
		result  models.AppSessionContext
		rsp     *http.Response
	)

	uri, err := s.getPcfPolicyAuthUri()
	if err != nil {
		return rspCode, rspBody
	}
	client := s.getClient(uri)

	ctx, _, err := s.consumer.Context().GetTokenCtx(models.ServiceName_NPCF_POLICYAUTHORIZATION, models.NfType_PCF)
	if err != nil {
		return rspCode, rspBody
	}

	result, rsp, err = client.IndividualApplicationSessionContextDocumentApi.
		GetAppSession(ctx, appSessionId)

	if rsp != nil {
		defer func() {
			if rsp.Request.Response != nil {
				rsp_err := rsp.Request.Response.Body.Close()
				if rsp_err != nil {
					logger.ConsumerLog.Errorf("ResponseBody can't be close: %+v", err)
				}
			}
		}()

		rspCode = rsp.StatusCode
		if rsp.StatusCode == http.StatusOK {
			rspBody = &result
		} else if err != nil {
			rspCode, rspBody = handleAPIServiceResponseError(rsp, err)
		}
	} else {
		// API Service Internal Error or Server No Response
		rspCode, rspBody = handleAPIServiceNoResponse(err)
	}

	return rspCode, rspBody
}

func (s *npcfService) PostAppSessions(asc *models.AppSessionContext) (int, interface{}, string) {
	var (
		err       error
		rspCode   int
		rspBody   interface{}
		appSessID string
		result    models.AppSessionContext
		rsp       *http.Response
	)

	uri, err := s.getPcfPolicyAuthUri()
	if err != nil {
		return rspCode, rspBody, appSessID
	}
	client := s.getClient(uri)

	ctx, _, err := s.consumer.Context().GetTokenCtx(models.ServiceName_NPCF_POLICYAUTHORIZATION, models.NfType_PCF)
	if err != nil {
		return rspCode, rspBody, appSessID
	}

	result, rsp, err = client.ApplicationSessionsCollectionApi.PostAppSessions(ctx, *asc)
	if rsp != nil {
		defer func() {
			if rsp.Request.Response != nil {
				rsp_err := rsp.Request.Response.Body.Close()
				if rsp_err != nil {
					logger.ConsumerLog.Errorf("ResponseBody can't be close: %+v", err)
				}
			}
		}()

		rspCode = rsp.StatusCode
		if rsp.StatusCode == http.StatusCreated {
			logger.ConsumerLog.Debugf("PostAppSessions RspData: %+v", result)
			rspBody = &result
			appSessID = getAppSessIDFromRspLocationHeader(rsp)
		} else if err != nil {
			rspCode, rspBody = handleAPIServiceResponseError(rsp, err)
		}
	} else {
		// API Service Internal Error or Server No Response
		rspCode, rspBody = handleAPIServiceNoResponse(err)
	}

	return rspCode, rspBody, appSessID
}

func (s *npcfService) PutAppSession(
	appSessionId string,
	ascUpdateData *models.AppSessionContextUpdateData,
	asc *models.AppSessionContext,
) (int, interface{}, string) {
	var (
		err       error
		rspCode   int
		rspBody   interface{}
		appSessID string
		result    models.AppSessionContext
		rsp       *http.Response
	)

	uri, err := s.getPcfPolicyAuthUri()
	if err != nil {
		return rspCode, rspBody, appSessID
	}
	client := s.getClient(uri)

	ctx, _, err := s.consumer.Context().GetTokenCtx(models.ServiceName_NPCF_POLICYAUTHORIZATION, models.NfType_PCF)
	if err != nil {
		return rspCode, rspBody, appSessID
	}

	appSessID = appSessionId
	result, rsp, err = client.IndividualApplicationSessionContextDocumentApi.
		GetAppSession(ctx, appSessionId)
	if rsp != nil {
		if rsp.Body != nil {
			if bodyCloseErr := rsp.Body.Close(); bodyCloseErr != nil {
				logger.ConsumerLog.Errorf("SearchNFInstances err: response body cannot close: %+v", bodyCloseErr)
			}
		}

		rspCode = rsp.StatusCode
		if rsp.StatusCode == http.StatusOK {
			// Patch
			result, rsp, err = client.IndividualApplicationSessionContextDocumentApi.ModAppSession(
				ctx, appSessionId, *ascUpdateData)
			if rsp != nil {
				defer func() {
					if rsp.Request.Response != nil {
						rsp_err := rsp.Request.Response.Body.Close()
						if rsp_err != nil {
							logger.ConsumerLog.Errorf("ResponseBody can't be close: %+v", err)
						}
					}
				}()

				rspCode = rsp.StatusCode
				if rsp.StatusCode == http.StatusOK {
					logger.ConsumerLog.Debugf("PatchAppSessions RspData: %+v", result)
					rspBody = &result
				} else if err != nil {
					rspCode, rspBody = handleAPIServiceResponseError(rsp, err)
				}
			} else {
				// API Service Internal Error or Server No Response
				rspCode, rspBody = handleAPIServiceNoResponse(err)
			}

			return rspCode, rspBody, appSessID
		}
		// TODO:
		// else if err != nil {
		// 	// Post
		// }
	} else {
		// API Service Internal Error or Server No Response
		rspCode, rspBody = handleAPIServiceNoResponse(err)
		return rspCode, rspBody, appSessID
	}

	return rspCode, rspBody, appSessID
}

func (s *npcfService) PatchAppSession(appSessionId string,
	ascUpdateData *models.AppSessionContextUpdateData,
) (int, interface{}) {
	var (
		err     error
		rspCode int
		rspBody interface{}
		result  models.AppSessionContext
		rsp     *http.Response
	)

	uri, err := s.getPcfPolicyAuthUri()
	if err != nil {
		return rspCode, rspBody
	}
	client := s.getClient(uri)

	ctx, _, err := s.consumer.Context().GetTokenCtx(models.ServiceName_NPCF_POLICYAUTHORIZATION, models.NfType_PCF)
	if err != nil {
		return rspCode, rspBody
	}

	result, rsp, err = client.IndividualApplicationSessionContextDocumentApi.ModAppSession(
		ctx, appSessionId, *ascUpdateData)
	if rsp != nil {
		defer func() {
			if rsp.Request.Response != nil {
				rsp_err := rsp.Request.Response.Body.Close()
				if rsp_err != nil {
					logger.ConsumerLog.Errorf("ResponseBody can't be close: %+v", err)
				}
			}
		}()

		rspCode = rsp.StatusCode
		if rsp.StatusCode == http.StatusOK {
			logger.ConsumerLog.Debugf("PatchAppSessions RspData: %+v", result)
			rspBody = &result
		} else if err != nil {
			rspCode, rspBody = handleAPIServiceResponseError(rsp, err)
		}
	} else {
		// API Service Internal Error or Server No Response
		rspCode, rspBody = handleAPIServiceNoResponse(err)
	}

	return rspCode, rspBody
}

func (s *npcfService) DeleteAppSession(appSessionId string) (int, interface{}) {
	var (
		err     error
		rspCode int
		rspBody interface{}
		result  models.AppSessionContext
		rsp     *http.Response
	)

	uri, err := s.getPcfPolicyAuthUri()
	if err != nil {
		return rspCode, rspBody
	}
	client := s.getClient(uri)

	param := &Npcf_PolicyAuthorization.DeleteAppSessionParamOpts{
		EventsSubscReqData: optional.NewInterface(models.EventsSubscReqData{}),
	}

	ctx, _, err := s.consumer.Context().GetTokenCtx(models.ServiceName_NPCF_POLICYAUTHORIZATION, models.NfType_PCF)
	if err != nil {
		return rspCode, rspBody
	}

	result, rsp, err = client.IndividualApplicationSessionContextDocumentApi.DeleteAppSession(
		ctx, appSessionId, param)
	if rsp != nil {
		defer func() {
			if rsp.Request.Response != nil {
				rsp_err := rsp.Request.Response.Body.Close()
				if rsp_err != nil {
					logger.ConsumerLog.Errorf("ResponseBody can't be close: %+v", err)
				}
			}
		}()

		rspCode = rsp.StatusCode
		if rsp.StatusCode == http.StatusOK {
			logger.ConsumerLog.Debugf("DeleteAppSessions RspData: %+v", result)
			rspBody = &result
		} else if err != nil {
			rspCode, rspBody = handleAPIServiceResponseError(rsp, err)
		}
	} else {
		// API Service Internal Error or Server No Response
		rspCode, rspBody = handleAPIServiceNoResponse(err)
	}

	return rspCode, rspBody
}

func getAppSessIDFromRspLocationHeader(rsp *http.Response) string {
	appSessID := ""
	loc := rsp.Header.Get("Location")
	if strings.Contains(loc, "http") {
		index := strings.LastIndex(loc, "/")
		appSessID = loc[index+1:]
	}
	logger.ConsumerLog.Infof("appSessID=%q", appSessID)
	return appSessID
}
