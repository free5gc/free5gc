package processor

import (
	"fmt"
	"net/http"

	"github.com/free5gc/nef/internal/logger"
	"github.com/free5gc/nef/pkg/factory"
	"github.com/free5gc/openapi"
	"github.com/free5gc/openapi/models"
)

func (p *Processor) GetApplicationsPFD(appIDs []string) *HandlerResponse {
	logger.PFDFLog.Infof("GetApplicationsPFD - appIDs: %v", appIDs)

	// TODO: Support SupportedFeatures
	rspCode, rspBody := p.Consumer().AppDataPfdsGet(appIDs)

	return &HandlerResponse{rspCode, nil, rspBody}
}

func (p *Processor) GetIndividualApplicationPFD(appID string) *HandlerResponse {
	logger.PFDFLog.Infof("GetIndividualApplicationPFD - appID[%s]", appID)

	// TODO: Support SupportedFeatures
	rspCode, rspBody := p.Consumer().AppDataPfdsAppIdGet(appID)

	return &HandlerResponse{rspCode, nil, rspBody}
}

func (p *Processor) PostPFDSubscriptions(pfdSubsc *models.PfdSubscription) *HandlerResponse {
	logger.PFDFLog.Infof("PostPFDSubscriptions - appIDs: %v", pfdSubsc.ApplicationIds)

	// TODO: Support SupportedFeatures
	if len(pfdSubsc.NotifyUri) == 0 {
		pd := openapi.ProblemDetailsDataNotFound("Absent of Notify URI")
		return &HandlerResponse{int(pd.Status), nil, pd}
	}

	subID := p.Notifier().PfdChangeNotifier.AddPfdSub(pfdSubsc)
	hdrs := make(map[string][]string)
	addLocationheader(hdrs, p.genPfdSubscriptionURI(subID))

	return &HandlerResponse{http.StatusCreated, hdrs, pfdSubsc}
}

func (p *Processor) DeleteIndividualPFDSubscription(subID string) *HandlerResponse {
	logger.PFDFLog.Infof("DeleteIndividualPFDSubscription - subID[%s]", subID)

	if err := p.Notifier().PfdChangeNotifier.DeletePfdSub(subID); err != nil {
		pd := openapi.ProblemDetailsDataNotFound(err.Error())
		return &HandlerResponse{int(pd.Status), nil, pd}
	}

	return &HandlerResponse{http.StatusNoContent, nil, nil}
}

func (p *Processor) genPfdSubscriptionURI(subID string) string {
	// E.g. "https://localhost:29505/nnef-pfdmanagement/v1/subscriptions/{subscriptionId}
	return fmt.Sprintf("%s/subscriptions/%s", p.Config().ServiceUri(factory.ServiceNefPfd), subID)
}
