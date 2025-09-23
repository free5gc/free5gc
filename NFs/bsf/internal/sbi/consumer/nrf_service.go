/*
 * BSF NRF Consumer
 */

package consumer

import (
	"context"
	"time"

	bsfContext "github.com/free5gc/bsf/internal/context"
	"github.com/free5gc/bsf/internal/logger"
	"github.com/free5gc/openapi"
	"github.com/free5gc/openapi/models"
	Nnrf_NFManagement "github.com/free5gc/openapi/nrf/NFManagement"
)

func BuildNFProfile(bsfContext *bsfContext.BSFContext) models.NrfNfManagementNfProfile {
	return bsfContext.GetBsfProfile()
}

func SendRegisterNFInstance(ctx context.Context) (*models.NrfNfManagementNfProfile, error) {
	// Set client and set url
	configuration := Nnrf_NFManagement.NewConfiguration()
	configuration.SetBasePath(bsfContext.BsfSelf.NrfUri)
	client := Nnrf_NFManagement.NewAPIClient(configuration)

	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}
		nfProfile := BuildNFProfile(bsfContext.BsfSelf)

		request := &Nnrf_NFManagement.RegisterNFInstanceRequest{
			NfInstanceID:             &bsfContext.BsfSelf.NfId,
			NrfNfManagementNfProfile: &nfProfile,
		}

		res, err := client.NFInstanceIDDocumentApi.RegisterNFInstance(context.TODO(), request)
		if err != nil {
			// Check if it's an OpenAPI error with more details
			if apiErr, ok := err.(*openapi.GenericOpenAPIError); ok {
				logger.ConsLog.Errorf("BSF register to NRF OpenAPI Error: %+v", apiErr.Error())
				logger.ConsLog.Errorf("BSF register to NRF Response Body: %s", string(apiErr.Body()))
				logger.ConsLog.Errorf("BSF register to NRF Response Model: %+v", apiErr.Model())
			} else {
				logger.ConsLog.Errorf("BSF register to NRF Error[%v]", err)
			}
			time.Sleep(2 * time.Second)
			continue
		}
		if res == nil {
			logger.ConsLog.Errorf("BSF register to NRF: received nil response")
			time.Sleep(2 * time.Second)
			continue
		}

		// Check if NFUpdate (no Location header) or NFRegister (has Location header)
		if res.Location == "" {
			// NFUpdate
			logger.ConsLog.Infof("BSF registration to NRF updated")
			return &res.NrfNfManagementNfProfile, nil
		} else {
			// NFRegister
			resourceUri := res.Location
			logger.ConsLog.Infof("BSF registration to NRF successful, resource: %s", resourceUri)
			return &res.NrfNfManagementNfProfile, nil
		}
	}
}

func SendDeregisterNFInstance() (*models.ProblemDetails, error) {
	// Set client and set url
	configuration := Nnrf_NFManagement.NewConfiguration()
	configuration.SetBasePath(bsfContext.BsfSelf.NrfUri)
	client := Nnrf_NFManagement.NewAPIClient(configuration)

	request := &Nnrf_NFManagement.DeregisterNFInstanceRequest{
		NfInstanceID: &bsfContext.BsfSelf.NfId,
	}

	_, err := client.NFInstanceIDDocumentApi.DeregisterNFInstance(context.TODO(), request)
	if err != nil {
		logger.ConsLog.Warnf("BSF deregistration from NRF failed[%v]", err)
		return nil, err
	}

	logger.ConsLog.Infof("BSF deregistration from NRF successful")
	return nil, nil
}
