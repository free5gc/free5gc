/*
 * BSF PCF UE Bindings Processor
 */

package processor

import (
	"net/http"

	"github.com/gin-gonic/gin"

	bsfContext "github.com/free5gc/bsf/internal/context"
	"github.com/free5gc/bsf/internal/logger"
	"github.com/free5gc/bsf/internal/metrics/business"
	"github.com/free5gc/bsf/internal/util"
	"github.com/free5gc/openapi/models"
)

// CreatePCFforUEBinding handles POST /pcf-ue-bindings
func CreatePCFforUEBinding(c *gin.Context) {
	logger.ProcLog.Infof("Handle CreatePCFforUEBinding")

	var request models.PcfForUeBinding
	if err := c.ShouldBindJSON(&request); err != nil {
		problemDetail := models.ProblemDetails{
			Status: http.StatusBadRequest,
			Cause:  "MALFORMED_REQUEST",
		}
		c.JSON(http.StatusBadRequest, problemDetail)
		return
	}

	// Convert to internal representation
	binding := &bsfContext.PcfForUeBinding{
		Supi:                request.Supi,
		Gpsi:                util.StringToPtr(request.Gpsi),
		PcfForUeFqdn:        util.StringToPtr(request.PcfForUeFqdn),
		PcfForUeIpEndPoints: request.PcfForUeIpEndPoints,
		PcfId:               util.StringToPtr(request.PcfId),
		PcfSetId:            util.StringToPtr(request.PcfSetId),
		BindLevel:           (*models.BindingLevel)(&request.BindLevel),
		SuppFeat:            util.StringToPtr(request.SuppFeat),
	}

	// Create new binding
	bindingId := bsfContext.BsfSelf.CreatePcfForUeBinding(binding)

	// Update metrics
	business.IncrPCFBindingGauge(business.PCF_UE_BINDING_TYPE_VALUE)
	business.IncrPCFBindingEventCounter(business.PCF_UE_BINDING_TYPE_VALUE, business.BINDING_EVENT_CREATE_VALUE, business.RESULT_SUCCESS_VALUE)

	// Convert back to response format
	response := models.PcfForUeBinding{
		Supi:                binding.Supi,
		Gpsi:                util.PtrToString(binding.Gpsi),
		PcfForUeFqdn:        util.PtrToString(binding.PcfForUeFqdn),
		PcfForUeIpEndPoints: binding.PcfForUeIpEndPoints,
		PcfId:               util.PtrToString(binding.PcfId),
		PcfSetId:            util.PtrToString(binding.PcfSetId),
		BindLevel:           (*binding.BindLevel),
		SuppFeat:            util.PtrToString(binding.SuppFeat),
	}

	locationHeader := "/nbsf-management/v1/pcf-ue-bindings/" + bindingId
	c.Header("Location", locationHeader)
	c.JSON(http.StatusCreated, response)
}

// GetPCFForUeBindings handles GET /pcf-ue-bindings
func GetPCFForUeBindings(c *gin.Context) {
	logger.ProcLog.Infof("Handle GetPCFForUeBindings")

	// Extract query parameters
	supi := c.Query("supi")
	gpsi := c.Query("gpsi")

	// Query bindings
	bindings := bsfContext.BsfSelf.QueryPcfForUeBindings(supi, gpsi)

	// Update metrics
	if len(bindings) > 0 {
		business.IncrPCFBindingEventCounter(business.PCF_UE_BINDING_TYPE_VALUE, business.BINDING_EVENT_QUERY_VALUE, business.RESULT_SUCCESS_VALUE)
	} else {
		business.IncrPCFBindingEventCounter(business.PCF_UE_BINDING_TYPE_VALUE, business.BINDING_EVENT_QUERY_VALUE, business.RESULT_FAILURE_VALUE)
	}

	if len(bindings) == 0 {
		c.JSON(http.StatusOK, []models.PcfForUeBinding{})
		return
	}

	// Convert to response format
	var response []models.PcfForUeBinding
	for _, binding := range bindings {
		response = append(response, models.PcfForUeBinding{
			Supi:                binding.Supi,
			Gpsi:                util.PtrToString(binding.Gpsi),
			PcfForUeFqdn:        util.PtrToString(binding.PcfForUeFqdn),
			PcfForUeIpEndPoints: binding.PcfForUeIpEndPoints,
			PcfId:               util.PtrToString(binding.PcfId),
			PcfSetId:            util.PtrToString(binding.PcfSetId),
			BindLevel:           (*binding.BindLevel),
			SuppFeat:            util.PtrToString(binding.SuppFeat),
		})
	}

	c.JSON(http.StatusOK, response)
}

// DeleteIndPCFforUEBinding handles DELETE /pcf-ue-bindings/{bindingId}
func DeleteIndPCFforUEBinding(c *gin.Context) {
	logger.ProcLog.Infof("Handle DeleteIndPCFforUEBinding")

	bindingId := c.Param("bindingId")

	if bsfContext.BsfSelf.DeletePcfForUeBinding(bindingId) {
		// Update metrics
		business.DecrPCFBindingGauge(business.PCF_UE_BINDING_TYPE_VALUE)
		business.IncrPCFBindingEventCounter(business.PCF_UE_BINDING_TYPE_VALUE, business.BINDING_EVENT_DELETE_VALUE, business.RESULT_SUCCESS_VALUE)
		c.Status(http.StatusNoContent)
	} else {
		business.IncrPCFBindingEventCounter(business.PCF_UE_BINDING_TYPE_VALUE, business.BINDING_EVENT_DELETE_VALUE, business.RESULT_FAILURE_VALUE)
		problemDetail := models.ProblemDetails{
			Status: http.StatusNotFound,
			Cause:  "RESOURCE_NOT_FOUND",
		}
		c.JSON(http.StatusNotFound, problemDetail)
	}
}

// UpdateIndPCFforUEBinding handles PATCH /pcf-ue-bindings/{bindingId}
func UpdateIndPCFforUEBinding(c *gin.Context) {
	logger.ProcLog.Infof("Handle UpdateIndPCFforUEBinding")

	bindingId := c.Param("bindingId")

	var patchRequest models.PcfForUeBindingPatch
	if err := c.ShouldBindJSON(&patchRequest); err != nil {
		problemDetail := models.ProblemDetails{
			Status: http.StatusBadRequest,
			Cause:  "MALFORMED_REQUEST",
		}
		c.JSON(http.StatusBadRequest, problemDetail)
		return
	}

	binding, exists := bsfContext.BsfSelf.GetPcfForUeBinding(bindingId)
	if !exists {
		business.IncrPCFBindingEventCounter(business.PCF_UE_BINDING_TYPE_VALUE, business.BINDING_EVENT_UPDATE_VALUE, business.RESULT_FAILURE_VALUE)
		problemDetail := models.ProblemDetails{
			Status: http.StatusNotFound,
			Cause:  "RESOURCE_NOT_FOUND",
		}
		c.JSON(http.StatusNotFound, problemDetail)
		return
	}

	// Apply patch
	if patchRequest.PcfForUeFqdn != "" {
		binding.PcfForUeFqdn = util.StringToPtr(patchRequest.PcfForUeFqdn)
	}
	if patchRequest.PcfForUeIpEndPoints != nil {
		binding.PcfForUeIpEndPoints = patchRequest.PcfForUeIpEndPoints
	}
	if patchRequest.PcfId != "" {
		binding.PcfId = util.StringToPtr(patchRequest.PcfId)
	}

	// Update binding
	bsfContext.BsfSelf.UpdatePcfForUeBinding(bindingId, binding)

	// Update metrics
	business.IncrPCFBindingEventCounter(business.PCF_UE_BINDING_TYPE_VALUE, business.BINDING_EVENT_UPDATE_VALUE, business.RESULT_SUCCESS_VALUE)

	// Return updated binding
	response := models.PcfForUeBinding{
		Supi:                binding.Supi,
		Gpsi:                util.PtrToString(binding.Gpsi),
		PcfForUeFqdn:        util.PtrToString(binding.PcfForUeFqdn),
		PcfForUeIpEndPoints: binding.PcfForUeIpEndPoints,
		PcfId:               util.PtrToString(binding.PcfId),
		PcfSetId:            util.PtrToString(binding.PcfSetId),
		BindLevel:           (*binding.BindLevel),
		SuppFeat:            util.PtrToString(binding.SuppFeat),
	}

	c.JSON(http.StatusOK, response)
}
