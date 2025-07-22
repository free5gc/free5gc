package processor

import (
	"net/http"

	"github.com/free5gc/nef/internal/logger"
	"github.com/free5gc/nef/pkg/factory"
	"github.com/free5gc/openapi"
	"github.com/free5gc/openapi/models"
	"github.com/free5gc/openapi/models_nef"
	"github.com/google/uuid"
)

func (p *Processor) GetTrafficInfluenceSubscription(
	afID string,
) *HandlerResponse {
	logger.TrafInfluLog.Infof("GetTrafficInfluenceSubscription - afID[%s]", afID)

	af := p.Context().GetAf(afID)
	if af == nil {
		pd := openapi.ProblemDetailsDataNotFound("AF is not found")
		return &HandlerResponse{http.StatusNotFound, nil, pd}
	}

	af.Mu.RLock()
	defer af.Mu.RUnlock()

	var tiSubs []models_nef.TrafficInfluSub
	for _, sub := range af.Subs {
		if sub.TiSub == nil {
			continue
		}
		tiSubs = append(tiSubs, *sub.TiSub)
	}
	return &HandlerResponse{http.StatusOK, nil, &tiSubs}
}

func (p *Processor) PostTrafficInfluenceSubscription(
	afID string,
	tiSub *models_nef.TrafficInfluSub,
) *HandlerResponse {
	logger.TrafInfluLog.Infof("PostTrafficInfluenceSubscription - afID[%s]", afID)

	rsp := validateTrafficInfluenceData(tiSub)
	if rsp != nil {
		return rsp
	}

	nefCtx := p.Context()
	af := nefCtx.GetAf(afID)
	if af == nil {
		af = nefCtx.NewAf(afID)
		if af == nil {
			pd := openapi.ProblemDetailsSystemFailure("No resource can be allocated")
			return &HandlerResponse{int(pd.Status), nil, pd}
		}
	}

	af.Mu.Lock()
	defer af.Mu.Unlock()

	correID := nefCtx.NewCorreID()
	afSub := af.NewSub(correID, tiSub)
	if afSub == nil {
		pd := openapi.ProblemDetailsSystemFailure("No resource can be allocated")
		return &HandlerResponse{int(pd.Status), nil, pd}
	}

	if len(tiSub.Gpsi) > 0 || len(tiSub.Ipv4Addr) > 0 || len(tiSub.Ipv6Addr) > 0 {
		// Single UE, sent to PCF
		asc := p.convertTrafficInfluSubToAppSessionContext(tiSub, afSub.NotifCorreID)
		rspStatus, rspBody, appSessID := p.Consumer().PostAppSessions(asc)
		if rspStatus != http.StatusCreated {
			return &HandlerResponse{rspStatus, nil, rspBody}
		}
		afSub.AppSessID = appSessID
	} else if len(tiSub.ExternalGroupId) > 0 || tiSub.AnyUeInd {
		// Group or any UE, sent to UDR
		afSub.InfluID = uuid.New().String()
		tiData := p.convertTrafficInfluSubToTrafficInfluData(tiSub, afSub.NotifCorreID)
		rspStatus, rspBody := p.Consumer().AppDataInfluenceDataPut(afSub.InfluID, tiData)
		if rspStatus != http.StatusOK &&
			rspStatus != http.StatusCreated &&
			rspStatus != http.StatusNoContent {
			return &HandlerResponse{rspStatus, nil, rspBody}
		}
	} else {
		// Invalid case. Return Error
		pd := openapi.ProblemDetailsMalformedReqSyntax("Not individual UE case, nor group case")
		return &HandlerResponse{int(pd.Status), nil, pd}
	}

	af.Subs[afSub.SubID] = afSub
	af.Log.Infoln("Subscription is added")

	nefCtx.AddAf(af)

	// Create Location URI
	tiSub.Self = p.genTrafficInfluSubURI(afID, afSub.SubID)
	headers := map[string][]string{
		"Location": {tiSub.Self},
	}
	return &HandlerResponse{http.StatusCreated, headers, tiSub}
}

func (p *Processor) GetIndividualTrafficInfluenceSubscription(
	afID, subID string,
) *HandlerResponse {
	logger.TrafInfluLog.Infof("GetIndividualTrafficInfluenceSubscription - afID[%s], subID[%s]", afID, subID)

	af := p.Context().GetAf(afID)
	if af == nil {
		pd := openapi.ProblemDetailsDataNotFound("AF is not found")
		return &HandlerResponse{http.StatusNotFound, nil, pd}
	}

	af.Mu.RLock()
	defer af.Mu.RUnlock()

	afSub, ok := af.Subs[subID]
	if !ok {
		pd := openapi.ProblemDetailsDataNotFound("Subscription is not found")
		return &HandlerResponse{http.StatusNotFound, nil, pd}
	}

	return &HandlerResponse{http.StatusOK, nil, afSub.TiSub}
}

func (p *Processor) PutIndividualTrafficInfluenceSubscription(
	afID, subID string,
	tiSub *models_nef.TrafficInfluSub,
) *HandlerResponse {
	logger.TrafInfluLog.Infof("PutIndividualTrafficInfluenceSubscription - afID[%s], subID[%s]", afID, subID)

	rsp := validateTrafficInfluenceData(tiSub)
	if rsp != nil {
		return rsp
	}

	af := p.Context().GetAf(afID)
	if af == nil {
		pd := openapi.ProblemDetailsDataNotFound("AF is not found")
		return &HandlerResponse{http.StatusNotFound, nil, pd}
	}

	af.Mu.Lock()
	defer af.Mu.Unlock()

	afSub, ok := af.Subs[subID]
	if !ok {
		pd := openapi.ProblemDetailsDataNotFound("Subscription is not found")
		return &HandlerResponse{http.StatusNotFound, nil, pd}
	}

	afSub.TiSub = tiSub
	if afSub.AppSessID != "" {
		asc := p.convertTrafficInfluSubToAppSessionContext(tiSub, afSub.NotifCorreID)
		rspStatus, rspBody, appSessID := p.Consumer().PostAppSessions(asc)
		if rspStatus != http.StatusCreated {
			return &HandlerResponse{rspStatus, nil, rspBody}
		}
		afSub.AppSessID = appSessID
	} else if afSub.InfluID != "" {
		tiData := p.convertTrafficInfluSubToTrafficInfluData(tiSub, afSub.NotifCorreID)
		rspStatus, rspBody := p.Consumer().AppDataInfluenceDataPut(afSub.InfluID, tiData)
		if rspStatus != http.StatusOK &&
			rspStatus != http.StatusCreated &&
			rspStatus != http.StatusNoContent {
			return &HandlerResponse{rspStatus, nil, rspBody}
		}
	} else {
		pd := openapi.ProblemDetailsDataNotFound("No AppSessID or InfluID")
		return &HandlerResponse{int(pd.Status), nil, pd}
	}

	return &HandlerResponse{http.StatusOK, nil, afSub.TiSub}
}

func (p *Processor) PatchIndividualTrafficInfluenceSubscription(
	afID, subID string,
	tiSubPatch *models_nef.TrafficInfluSubPatch,
) *HandlerResponse {
	logger.TrafInfluLog.Infof("PatchIndividualTrafficInfluenceSubscription - afID[%s], subID[%s]", afID, subID)

	af := p.Context().GetAf(afID)
	if af == nil {
		pd := openapi.ProblemDetailsDataNotFound("AF is not found")
		return &HandlerResponse{http.StatusNotFound, nil, pd}
	}

	af.Mu.Lock()
	defer af.Mu.Unlock()

	afSub, ok := af.Subs[subID]
	if !ok {
		pd := openapi.ProblemDetailsDataNotFound("Subscription is not found")
		return &HandlerResponse{http.StatusNotFound, nil, pd}
	}

	if afSub.AppSessID != "" {
		ascUpdateData := p.convertTrafficInfluSubPatchToAppSessionContextUpdateData(tiSubPatch)
		rspStatus, rspBody := p.Consumer().PatchAppSession(afSub.AppSessID, ascUpdateData)
		if rspStatus != http.StatusOK &&
			rspStatus != http.StatusNoContent {
			return &HandlerResponse{rspStatus, nil, rspBody}
		}
	} else if afSub.InfluID != "" {
		tiDataPatch := p.convertTrafficInfluSubPatchToTrafficInfluDataPatch(tiSubPatch)
		rspStatus, rspBody := p.Consumer().AppDataInfluenceDataPatch(afSub.InfluID, tiDataPatch)
		if rspStatus != http.StatusOK &&
			rspStatus != http.StatusNoContent {
			return &HandlerResponse{rspStatus, nil, rspBody}
		}
	} else {
		pd := openapi.ProblemDetailsDataNotFound("No AppSessID or InfluID")
		return &HandlerResponse{int(pd.Status), nil, pd}
	}

	afSub.PatchTiSubData(tiSubPatch)
	return &HandlerResponse{http.StatusOK, nil, afSub.TiSub}
}

func (p *Processor) DeleteIndividualTrafficInfluenceSubscription(
	afID, subID string,
) *HandlerResponse {
	logger.TrafInfluLog.Infof("DeleteIndividualTrafficInfluenceSubscription - afID[%s], subID[%s]", afID, subID)

	af := p.Context().GetAf(afID)
	if af == nil {
		pd := openapi.ProblemDetailsDataNotFound("AF is not found")
		return &HandlerResponse{http.StatusNotFound, nil, pd}
	}

	af.Mu.Lock()
	defer af.Mu.Unlock()

	sub, ok := af.Subs[subID]
	if !ok {
		pd := openapi.ProblemDetailsDataNotFound("Subscription is not found")
		return &HandlerResponse{http.StatusNotFound, nil, pd}
	}

	if sub.AppSessID != "" {
		rspStatus, rspBody := p.Consumer().DeleteAppSession(sub.AppSessID)
		if rspStatus != http.StatusOK &&
			rspStatus != http.StatusNoContent {
			return &HandlerResponse{rspStatus, nil, rspBody}
		}
	} else {
		rspStatus, rspBody := p.Consumer().AppDataInfluenceDataDelete(sub.InfluID)
		if rspStatus != http.StatusOK &&
			rspStatus != http.StatusNoContent {
			return &HandlerResponse{rspStatus, nil, rspBody}
		}
	}
	delete(af.Subs, subID)
	return &HandlerResponse{http.StatusNoContent, nil, nil}
}

func validateTrafficInfluenceData(
	tiSub *models_nef.TrafficInfluSub,
) *HandlerResponse {
	// TS29.522: One of "afAppId", "trafficFilters" or "ethTrafficFilters" shall be included.
	if tiSub.AfAppId == "" &&
		len(tiSub.TrafficFilters) == 0 &&
		len(tiSub.EthTrafficFilters) == 0 {
		pd := openapi.
			ProblemDetailsMalformedReqSyntax(
				"Missing one of afAppId, trafficFilters or ethTrafficFilters")
		return &HandlerResponse{int(pd.Status), nil, pd}
	}

	// TS29.522: One of individual UE identifier
	// (i.e. "gpsi", “macAddr”, "ipv4Addr" or "ipv6Addr"),
	// External Group Identifier (i.e. "externalGroupId") or
	// any UE indication "anyUeInd" shall be included.
	if tiSub.Gpsi == "" &&
		tiSub.Ipv4Addr == "" &&
		tiSub.Ipv6Addr == "" &&
		tiSub.ExternalGroupId == "" &&
		!tiSub.AnyUeInd {
		pd := openapi.
			ProblemDetailsMalformedReqSyntax(
				"Missing one of Gpsi, Ipv4Addr, Ipv6Addr, ExternalGroupId, AnyUeInd")
		return &HandlerResponse{int(pd.Status), nil, pd}
	}
	return nil
}

func (p *Processor) genTrafficInfluSubURI(
	afID, subscriptionId string,
) string {
	// E.g. https://localhost:29505/3gpp-traffic-Influence/v1/{afId}/subscriptions/{subscriptionId}
	return p.Config().ServiceUri(factory.ServiceTraffInflu) + "/" + afID + "/subscriptions/" + subscriptionId
}

func (p *Processor) genNotificationUri() string {
	return p.Config().ServiceUri(factory.ServiceNefCallback) + "/notification/smf"
}

func (p *Processor) convertTrafficInfluSubToAppSessionContext(
	tiSub *models_nef.TrafficInfluSub,
	notifCorreID string,
) *models.AppSessionContext {
	asc := &models.AppSessionContext{
		AscReqData: &models.AppSessionContextReqData{
			AfAppId: tiSub.AfAppId,
			AfRoutReq: &models.AfRoutingRequirement{
				AppReloc:    tiSub.AppReloInd,
				RouteToLocs: tiSub.TrafficRoutes,
				TempVals:    tiSub.TempValidities,
			},
			UeIpv4:    tiSub.Ipv4Addr,
			UeIpv6:    tiSub.Ipv6Addr,
			UeMac:     tiSub.MacAddr,
			NotifUri:  tiSub.NotificationDestination,
			SuppFeat:  tiSub.SuppFeat,
			Dnn:       tiSub.Dnn,
			SliceInfo: tiSub.Snssai,
			// Supi: ,
		},
	}

	if tiSub.DnaiChgType != "" {
		asc.AscReqData.AfRoutReq.UpPathChgSub = &models.UpPathChgEvent{
			DnaiChgType:     tiSub.DnaiChgType,
			NotificationUri: p.genNotificationUri(),
			NotifCorreId:    notifCorreID,
		}
	}
	return asc
}

func (p *Processor) convertTrafficInfluSubPatchToAppSessionContextUpdateData(
	tiSubPatch *models_nef.TrafficInfluSubPatch,
) *models.AppSessionContextUpdateData {
	ascUpdate := &models.AppSessionContextUpdateData{
		AfRoutReq: &models.AfRoutingRequirementRm{
			AppReloc:    tiSubPatch.AppReloInd,
			RouteToLocs: tiSubPatch.TrafficRoutes,
			TempVals:    tiSubPatch.TempValidities,
		},
	}
	return ascUpdate
}

func (p *Processor) convertTrafficInfluSubToTrafficInfluData(
	tiSub *models_nef.TrafficInfluSub,
	notifCorreID string,
) *models.TrafficInfluData {
	tiData := &models.TrafficInfluData{
		AfAppId:    tiSub.AfAppId,
		AppReloInd: tiSub.AppReloInd,
		// Supi: ,
		DnaiChgType:           tiSub.DnaiChgType,
		UpPathChgNotifUri:     p.genNotificationUri(),
		UpPathChgNotifCorreId: notifCorreID,
		Dnn:                   tiSub.Dnn,
		Snssai:                tiSub.Snssai,
		EthTrafficFilters:     tiSub.EthTrafficFilters,
		TrafficFilters:        tiSub.TrafficFilters,
		TrafficRoutes:         tiSub.TrafficRoutes,
		TraffCorreInd:         tiSub.TfcCorrInd,
		// ValidStartTime: ,
		// ValidEndTime: ,
		TempValidities:    tiSub.TempValidities,
		AfAckInd:          tiSub.AfAckInd,
		AddrPreserInd:     tiSub.AddrPreserInd,
		SupportedFeatures: tiSub.SuppFeat,
	}

	// TODO: handle ExternalGroupId
	if tiSub.AnyUeInd {
		tiData.InterGroupId = "AnyUE"
	}

	return tiData
}

func (p *Processor) convertTrafficInfluSubPatchToTrafficInfluDataPatch(
	tiSubPatch *models_nef.TrafficInfluSubPatch,
) *models.TrafficInfluDataPatch {
	tiDataPatch := &models.TrafficInfluDataPatch{
		AppReloInd:        tiSubPatch.AppReloInd,
		EthTrafficFilters: tiSubPatch.EthTrafficFilters,
		TrafficFilters:    tiSubPatch.TrafficFilters,
		TrafficRoutes:     tiSubPatch.TrafficRoutes,
	}
	return tiDataPatch
}
