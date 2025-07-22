package processor

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/free5gc/amf/internal/context"
	"github.com/free5gc/amf/internal/logger"
	"github.com/free5gc/openapi/models"
)

func (p *Processor) HandleProvideDomainSelectionInfoRequest(c *gin.Context) {
	logger.MtLog.Info("Handle Provide Domain Selection Info Request")

	ueContextID := c.Param("ueContextId")
	infoClassQuery := c.Query("info-class")
	supportedFeaturesQuery := c.Query("supported-features")

	ueContextInfo, problemDetails := p.ProvideDomainSelectionInfoProcedure(ueContextID,
		infoClassQuery, supportedFeaturesQuery)
	if problemDetails != nil {
		c.JSON(int(problemDetails.Status), problemDetails)
	} else {
		c.JSON(http.StatusOK, ueContextInfo)
	}
}

func (p *Processor) ProvideDomainSelectionInfoProcedure(ueContextID string, infoClassQuery string,
	supportedFeaturesQuery string) (
	*models.UeContextInfo, *models.ProblemDetails,
) {
	amfSelf := context.GetSelf()

	ue, ok := amfSelf.AmfUeFindByUeContextID(ueContextID)
	if !ok {
		logger.CtxLog.Warnf("AmfUe Context[%s] not found", ueContextID)
		problemDetails := &models.ProblemDetails{
			Status: http.StatusNotFound,
			Cause:  "CONTEXT_NOT_FOUND",
		}
		return nil, problemDetails
	}

	ue.Lock.Lock()
	defer ue.Lock.Unlock()

	ueContextInfo := new(models.UeContextInfo)

	// TODO: Error Status 307, 403 in TS29.518 Table 6.3.3.3.3.1-3
	anType := ue.GetAnType()
	if anType != "" && infoClassQuery != "" {
		ranUe := ue.RanUe[anType]
		ueContextInfo.AccessType = anType
		ueContextInfo.LastActTime = ranUe.LastActTime
		ueContextInfo.RatType = ue.RatType
		ueContextInfo.SupportedFeatures = ranUe.SupportedFeatures
		ueContextInfo.SupportVoPS = ranUe.SupportVoPS
		ueContextInfo.SupportVoPSn3gpp = ranUe.SupportVoPSn3gpp
	}

	return ueContextInfo, nil
}
