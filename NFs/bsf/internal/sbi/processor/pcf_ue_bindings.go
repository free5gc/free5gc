/*
 * BSF PCF UE Bindings Processor
 */

package processor

import (
	"net/http"

	"github.com/gin-gonic/gin"

	bsfContext "github.com/free5gc/bsf/internal/context"
	"github.com/free5gc/bsf/internal/logger"
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
		Gpsi:                stringToPtr(request.Gpsi),
		PcfForUeFqdn:        stringToPtr(request.PcfForUeFqdn),
		PcfForUeIpEndPoints: request.PcfForUeIpEndPoints,
		PcfId:               stringToPtr(request.PcfId),
		PcfSetId:            stringToPtr(request.PcfSetId),
		BindLevel:           (*models.BindingLevel)(&request.BindLevel),
		SuppFeat:            stringToPtr(request.SuppFeat),
	}

	// Create new binding
	bindingId := bsfContext.BsfSelf.CreatePcfForUeBinding(binding)

	// Convert back to response format
	response := models.PcfForUeBinding{
		Supi:                binding.Supi,
		Gpsi:                ptrToString(binding.Gpsi),
		PcfForUeFqdn:        ptrToString(binding.PcfForUeFqdn),
		PcfForUeIpEndPoints: binding.PcfForUeIpEndPoints,
		PcfId:               ptrToString(binding.PcfId),
		PcfSetId:            ptrToString(binding.PcfSetId),
		BindLevel:           (*binding.BindLevel),
		SuppFeat:            ptrToString(binding.SuppFeat),
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

	if len(bindings) == 0 {
		c.JSON(http.StatusOK, []models.PcfForUeBinding{})
		return
	}

	// Convert to response format
	var response []models.PcfForUeBinding
	for _, binding := range bindings {
		response = append(response, models.PcfForUeBinding{
			Supi:                binding.Supi,
			Gpsi:                ptrToString(binding.Gpsi),
			PcfForUeFqdn:        ptrToString(binding.PcfForUeFqdn),
			PcfForUeIpEndPoints: binding.PcfForUeIpEndPoints,
			PcfId:               ptrToString(binding.PcfId),
			PcfSetId:            ptrToString(binding.PcfSetId),
			BindLevel:           (*binding.BindLevel),
			SuppFeat:            ptrToString(binding.SuppFeat),
		})
	}

	c.JSON(http.StatusOK, response)
}

// DeleteIndPCFforUEBinding handles DELETE /pcf-ue-bindings/{bindingId}
func DeleteIndPCFforUEBinding(c *gin.Context) {
	logger.ProcLog.Infof("Handle DeleteIndPCFforUEBinding")

	bindingId := c.Param("bindingId")

	if bsfContext.BsfSelf.DeletePcfForUeBinding(bindingId) {
		c.Status(http.StatusNoContent)
	} else {
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
		problemDetail := models.ProblemDetails{
			Status: http.StatusNotFound,
			Cause:  "RESOURCE_NOT_FOUND",
		}
		c.JSON(http.StatusNotFound, problemDetail)
		return
	}

	// Apply patch
	if patchRequest.PcfForUeFqdn != "" {
		binding.PcfForUeFqdn = stringToPtr(patchRequest.PcfForUeFqdn)
	}
	if patchRequest.PcfForUeIpEndPoints != nil {
		binding.PcfForUeIpEndPoints = patchRequest.PcfForUeIpEndPoints
	}
	if patchRequest.PcfId != "" {
		binding.PcfId = stringToPtr(patchRequest.PcfId)
	}

	// Update binding
	bsfContext.BsfSelf.UpdatePcfForUeBinding(bindingId, binding)

	// Return updated binding
	response := models.PcfForUeBinding{
		Supi:                binding.Supi,
		Gpsi:                ptrToString(binding.Gpsi),
		PcfForUeFqdn:        ptrToString(binding.PcfForUeFqdn),
		PcfForUeIpEndPoints: binding.PcfForUeIpEndPoints,
		PcfId:               ptrToString(binding.PcfId),
		PcfSetId:            ptrToString(binding.PcfSetId),
		BindLevel:           (*binding.BindLevel),
		SuppFeat:            ptrToString(binding.SuppFeat),
	}

	c.JSON(http.StatusOK, response)
}
