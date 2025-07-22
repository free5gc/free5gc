package processor

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/free5gc/openapi/models"
	udm_context "github.com/free5gc/udm/internal/context"
	"github.com/free5gc/udm/internal/logger"
)

// EE service
func (p *Processor) CreateEeSubscriptionProcedure(c *gin.Context, ueIdentity string,
	eesubscription models.UdmEeEeSubscription,
) {
	udmSelf := p.Context()
	logger.EeLog.Debugf("udIdentity: %s", ueIdentity)
	switch {
	// GPSI (MSISDN identifier) represents a single UE
	case strings.HasPrefix(ueIdentity, "msisdn-"):
		fallthrough
	// GPSI (External identifier) represents a single UE
	case strings.HasPrefix(ueIdentity, "extid-"):
		if ue, ok := udmSelf.UdmUeFindByGpsi(ueIdentity); ok {
			id, err := udmSelf.EeSubscriptionIDGenerator.Allocate()
			if err != nil {
				problemDetails := &models.ProblemDetails{
					Status: http.StatusInternalServerError,
					Cause:  "UNSPECIFIED_NF_FAILURE",
				}
				c.JSON(int(problemDetails.Status), problemDetails)
				return
			}
			subscriptionID := strconv.Itoa(int(id))
			ue.EeSubscriptions[subscriptionID] = &eesubscription
			createdEeSubscription := &models.UdmEeCreatedEeSubscription{
				EeSubscription: &eesubscription,
			}
			c.JSON(http.StatusCreated, createdEeSubscription)
		} else {
			problemDetails := &models.ProblemDetails{
				Status: http.StatusNotFound,
				Cause:  "USER_NOT_FOUND",
			}
			c.JSON(int(problemDetails.Status), problemDetails)
		}
	// external groupID represents a group of UEs
	case strings.HasPrefix(ueIdentity, "extgroupid-"):
		id, err := udmSelf.EeSubscriptionIDGenerator.Allocate()
		if err != nil {
			problemDetails := &models.ProblemDetails{
				Status: http.StatusInternalServerError,
				Cause:  "UNSPECIFIED_NF_FAILURE",
			}
			c.JSON(int(problemDetails.Status), problemDetails)
			return
		}
		subscriptionID := strconv.Itoa(int(id))
		createdEeSubscription := &models.UdmEeCreatedEeSubscription{
			EeSubscription: &eesubscription,
		}

		udmSelf.UdmUePool.Range(func(key, value interface{}) bool {
			ue := value.(*udm_context.UdmUeContext)
			if ue.ExternalGroupID == ueIdentity {
				ue.EeSubscriptions[subscriptionID] = &eesubscription
			}
			return true
		})
		c.JSON(http.StatusCreated, createdEeSubscription)
	// represents any UEs
	case ueIdentity == "anyUE":
		id, err := udmSelf.EeSubscriptionIDGenerator.Allocate()
		if err != nil {
			problemDetails := &models.ProblemDetails{
				Status: http.StatusInternalServerError,
				Cause:  "UNSPECIFIED_NF_FAILURE",
			}
			c.JSON(int(problemDetails.Status), problemDetails)
			return
		}
		subscriptionID := strconv.Itoa(int(id))
		createdEeSubscription := &models.UdmEeCreatedEeSubscription{
			EeSubscription: &eesubscription,
		}
		udmSelf.UdmUePool.Range(func(key, value interface{}) bool {
			ue := value.(*udm_context.UdmUeContext)
			ue.EeSubscriptions[subscriptionID] = &eesubscription
			return true
		})
		c.JSON(http.StatusCreated, createdEeSubscription)
	default:
		problemDetails := &models.ProblemDetails{
			Status: http.StatusBadRequest,
			Cause:  "MANDATORY_IE_INCORRECT",
			InvalidParams: []models.InvalidParam{
				{
					Param:  "ueIdentity",
					Reason: "incorrect format",
				},
			},
		}
		c.JSON(int(problemDetails.Status), problemDetails)
	}
}

// TODO: complete this procedure based on TS 29503 5.5
func (p *Processor) DeleteEeSubscriptionProcedure(c *gin.Context, ueIdentity string, subscriptionID string) {
	udmSelf := p.Context()

	switch {
	case strings.HasPrefix(ueIdentity, "msisdn-"):
		fallthrough
	case strings.HasPrefix(ueIdentity, "extid-"):
		if ue, ok := udmSelf.UdmUeFindByGpsi(ueIdentity); ok {
			delete(ue.EeSubscriptions, subscriptionID)
		}
	case strings.HasPrefix(ueIdentity, "extgroupid-"):
		udmSelf.UdmUePool.Range(func(key, value interface{}) bool {
			ue := value.(*udm_context.UdmUeContext)
			if ue.ExternalGroupID == ueIdentity {
				delete(ue.EeSubscriptions, subscriptionID)
			}
			return true
		})
	case ueIdentity == "anyUE":
		udmSelf.UdmUePool.Range(func(key, value interface{}) bool {
			ue := value.(*udm_context.UdmUeContext)
			delete(ue.EeSubscriptions, subscriptionID)
			return true
		})
	}
	if id, err := strconv.ParseInt(subscriptionID, 10, 64); err != nil {
		logger.EeLog.Warnf("subscriptionID covert type error: %+v", err)
	} else {
		udmSelf.EeSubscriptionIDGenerator.FreeID(id)
	}

	// only return 204 no content
	c.Status(http.StatusNoContent)
}

// TODO: complete this procedure based on TS 29503 5.5
func (p *Processor) UpdateEeSubscriptionProcedure(c *gin.Context, ueIdentity string, subscriptionID string,
	patchList []models.PatchItem,
) {
	udmSelf := p.Context()

	switch {
	case strings.HasPrefix(ueIdentity, "msisdn-"):
		fallthrough
	case strings.HasPrefix(ueIdentity, "extid-"):
		if ue, ok := udmSelf.UdmUeFindByGpsi(ueIdentity); ok {
			if _, ok := ue.EeSubscriptions[subscriptionID]; ok {
				for _, patchItem := range patchList {
					logger.EeLog.Debugf("patch item: %+v", patchItem)
					// TODO: patch the Eesubscription
				}
				c.Status(http.StatusNoContent)
			} else {
				problemDetails := &models.ProblemDetails{
					Status: http.StatusNotFound,
					Cause:  "SUBSCRIPTION_NOT_FOUND",
				}
				c.JSON(int(problemDetails.Status), problemDetails)
			}
		} else {
			problemDetails := &models.ProblemDetails{
				Status: http.StatusNotFound,
				Cause:  "SUBSCRIPTION_NOT_FOUND",
			}
			c.JSON(int(problemDetails.Status), problemDetails)
		}
	case strings.HasPrefix(ueIdentity, "extgroupid-"):
		udmSelf.UdmUePool.Range(func(key, value interface{}) bool {
			ue := value.(*udm_context.UdmUeContext)
			if ue.ExternalGroupID == ueIdentity {
				if _, ok := ue.EeSubscriptions[subscriptionID]; ok {
					for _, patchItem := range patchList {
						logger.EeLog.Debugf("patch item: %+v", patchItem)
						// TODO: patch the Eesubscription
					}
				}
			}
			return true
		})
		c.Status(http.StatusNoContent)
	case ueIdentity == "anyUE":
		udmSelf.UdmUePool.Range(func(key, value interface{}) bool {
			ue := value.(*udm_context.UdmUeContext)
			if _, ok := ue.EeSubscriptions[subscriptionID]; ok {
				for _, patchItem := range patchList {
					logger.EeLog.Debugf("patch item: %+v", patchItem)
					// TODO: patch the Eesubscription
				}
			}
			return true
		})
		c.Status(http.StatusNoContent)
	default:
		problemDetails := &models.ProblemDetails{
			Status: http.StatusBadRequest,
			Cause:  "MANDATORY_IE_INCORRECT",
			InvalidParams: []models.InvalidParam{
				{
					Param:  "ueIdentity",
					Reason: "incorrect format",
				},
			},
		}
		c.JSON(int(problemDetails.Status), problemDetails)
	}
}
