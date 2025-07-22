package processor

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/free5gc/amf/internal/context"
	"github.com/free5gc/amf/internal/logger"
	"github.com/free5gc/openapi/models"
)

func (p *Processor) HandleProvideLocationInfoRequest(c *gin.Context, requestLocInfo models.RequestLocInfo) {
	logger.ProducerLog.Info("Handle Provide Location Info Request")

	ueContextID := c.Param("ueContextId")

	provideLocInfo, problemDetails := p.ProvideLocationInfoProcedure(requestLocInfo, ueContextID)
	if problemDetails != nil {
		c.JSON(int(problemDetails.Status), problemDetails)
	} else {
		c.JSON(http.StatusOK, provideLocInfo)
	}
}

func (p *Processor) ProvideLocationInfoProcedure(requestLocInfo models.RequestLocInfo, ueContextID string) (
	*models.ProvideLocInfo, *models.ProblemDetails,
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

	anType := ue.GetAnType()
	if anType == "" {
		problemDetails := &models.ProblemDetails{
			Status: http.StatusNotFound,
			Cause:  "CONTEXT_NOT_FOUND",
		}
		return nil, problemDetails
	}

	provideLocInfo := new(models.ProvideLocInfo)

	ranUe := ue.RanUe[anType]
	if requestLocInfo.Req5gsLoc || requestLocInfo.ReqCurrentLoc {
		provideLocInfo.CurrentLoc = true
		provideLocInfo.Location = &ue.Location
	}

	if requestLocInfo.ReqRatType {
		provideLocInfo.RatType = ue.RatType
	}

	if requestLocInfo.ReqTimeZone {
		provideLocInfo.Timezone = ue.TimeZone
	}

	if requestLocInfo.SupportedFeatures != "" {
		provideLocInfo.SupportedFeatures = ranUe.SupportedFeatures
	}
	return provideLocInfo, nil
}
