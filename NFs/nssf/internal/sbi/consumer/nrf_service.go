/*
 * NSSF Consumer
 *
 * Network Function Management
 */

package consumer

import (
	"context"
	"fmt"
	"strings"
	"time"

	nssf_context "github.com/free5gc/nssf/internal/context"
	"github.com/free5gc/nssf/internal/logger"
	"github.com/free5gc/openapi"
	"github.com/free5gc/openapi/models"
	"github.com/free5gc/openapi/nrf/NFManagement"
)

type NrfService struct {
	nrfNfMgmtClient *NFManagement.APIClient
	// NOTE: No mutex needed. One connection at a time.
}

func (ns *NrfService) buildNFProfile(context *nssf_context.NSSFContext) (
	profile models.NrfNfManagementNfProfile, err error,
) {
	profile.NfInstanceId = context.NfId
	profile.NfType = models.NrfNfManagementNfType_NSSF
	profile.NfStatus = models.NrfNfManagementNfStatus_REGISTERED
	profile.PlmnList = context.SupportedPlmnList
	profile.Ipv4Addresses = []string{context.RegisterIPv4}
	var services []models.NrfNfManagementNfService
	for _, nfService := range context.NfService {
		services = append(services, nfService)
	}
	if len(services) > 0 {
		profile.NfServices = services
	}
	return
}

func (ns *NrfService) SendRegisterNFInstance(ctx context.Context, nssfCtx *nssf_context.NSSFContext) (
	resourceNrfUri string, retrieveNfInstanceId string, err error,
) {
	nfInstanceId := nssfCtx.NfId
	profile, err := ns.buildNFProfile(nssfCtx)
	if err != nil {
		return "", "", fmt.Errorf("failed to build nrf profile: %s", err.Error())
	}
	apiClient := ns.nrfNfMgmtClient

	var res *NFManagement.RegisterNFInstanceResponse
	var nf models.NrfNfManagementNfProfile
	finish := false
	for !finish {
		select {
		case <-ctx.Done():
			return "", "", fmt.Errorf("context done")

		default:
			req := &NFManagement.RegisterNFInstanceRequest{
				NfInstanceID:             &nfInstanceId,
				NrfNfManagementNfProfile: &profile,
			}

			res, err = apiClient.NFInstanceIDDocumentApi.RegisterNFInstance(ctx, req)
			if err != nil || res == nil {
				// TODO : add log
				logger.ConsumerLog.Errorf("NSSF register to NRF Error[%s]", err.Error())
				const retryInterval = 2 * time.Second
				time.Sleep(retryInterval)
				continue
			}

			resourceUri := res.Location
			resourceNrfUri, _, _ = strings.Cut(resourceUri, "/nnrf-nfm/")
			retrieveNfInstanceId = resourceUri[strings.LastIndex(resourceUri, "/")+1:]
			nf = res.NrfNfManagementNfProfile

			oauth2 := false
			if nf.CustomInfo != nil {
				v, ok := nf.CustomInfo["oauth2"].(bool)
				if ok {
					oauth2 = v
					logger.MainLog.Infoln("OAuth2 setting receive from NRF:", oauth2)
				}
			}
			nssf_context.GetSelf().OAuth2Required = oauth2
			if oauth2 && nssf_context.GetSelf().NrfCertPem == "" {
				logger.CfgLog.Error("OAuth2 enable but no nrfCertPem provided in config.")
			}
			finish = true
		}
	}
	return resourceNrfUri, retrieveNfInstanceId, err
}

func (ns *NrfService) SendDeregisterNFInstance(nfInstanceId string) (*models.ProblemDetails, error) {
	logger.ConsumerLog.Infof("Send Deregister NFInstance [%s]", nfInstanceId)

	var err error

	ctx, pd, err := nssf_context.GetSelf().GetTokenCtx(models.ServiceName_NNRF_NFM, models.NrfNfManagementNfType_NRF)
	if err != nil {
		return pd, err
	}

	client := ns.nrfNfMgmtClient

	req := &NFManagement.DeregisterNFInstanceRequest{
		NfInstanceID: &nfInstanceId,
	}

	_, err = client.NFInstanceIDDocumentApi.DeregisterNFInstance(ctx, req)
	if err != nil {
		if apiErr, ok := err.(openapi.GenericOpenAPIError); ok {
			// API error
			if deregError, ok2 := apiErr.Model().(NFManagement.DeregisterNFInstanceError); ok2 {
				return &deregError.ProblemDetails, err
			}
			return nil, err
		}

		// Golang error
		return nil, err
	}

	return nil, nil
}
