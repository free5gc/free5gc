/*
 * BSF PCF MBS Bindings Processor
 */

package processor

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"

	bsfContext "github.com/free5gc/bsf/internal/context"
	"github.com/free5gc/bsf/internal/logger"
	"github.com/free5gc/openapi/models"
)

// CreatePCFMbsBinding handles POST /pcf-mbs-bindings
func CreatePCFMbsBinding(c *gin.Context) {
	logger.ProcLog.Infof("Handle CreatePCFMbsBinding")

	var request models.PcfMbsBinding
	if err := c.ShouldBindJSON(&request); err != nil {
		problemDetail := models.ProblemDetails{
			Status: http.StatusBadRequest,
			Cause:  "MALFORMED_REQUEST",
		}
		c.JSON(http.StatusBadRequest, problemDetail)
		return
	}

	// Check for existing binding
	existingBindings := bsfContext.BsfSelf.QueryPcfMbsBindings(request.MbsSessionId)
	if len(existingBindings) > 0 {
		// Return existing binding information
		problemDetail := models.ProblemDetails{
			Status: http.StatusForbidden,
			Cause:  "BINDING_ALREADY_EXISTS",
		}
		c.JSON(http.StatusForbidden, problemDetail)
		return
	}

	// Convert to internal representation
	binding := &bsfContext.PcfMbsBinding{
		MbsSessionId:   request.MbsSessionId,
		PcfFqdn:        stringToPtr(request.PcfFqdn),
		PcfIpEndPoints: request.PcfIpEndPoints,
		PcfId:          stringToPtr(request.PcfId),
		PcfSetId:       stringToPtr(request.PcfSetId),
		BindLevel:      (*models.BindingLevel)(&request.BindLevel),
		SuppFeat:       stringToPtr(request.SuppFeat),
	}

	if request.RecoveryTime != nil {
		binding.RecoveryTime = request.RecoveryTime
	}

	// Create new binding
	bindingId := bsfContext.BsfSelf.CreatePcfMbsBinding(binding)

	// Convert back to response format
	response := models.PcfMbsBinding{
		MbsSessionId:   binding.MbsSessionId,
		PcfFqdn:        ptrToString(binding.PcfFqdn),
		PcfIpEndPoints: binding.PcfIpEndPoints,
		PcfId:          ptrToString(binding.PcfId),
		PcfSetId:       ptrToString(binding.PcfSetId),
		BindLevel:      (*binding.BindLevel),
		SuppFeat:       ptrToString(binding.SuppFeat),
	}

	if binding.RecoveryTime != nil {
		response.RecoveryTime = binding.RecoveryTime
	}

	locationHeader := "/nbsf-management/v1/pcf-mbs-bindings/" + bindingId
	c.Header("Location", locationHeader)
	c.JSON(http.StatusCreated, response)
}

// GetPCFMbsBinding handles GET /pcf-mbs-bindings
func GetPCFMbsBinding(c *gin.Context) {
	logger.ProcLog.Infof("Handle GetPCFMbsBinding")

	// Extract query parameters - MBS Session ID is required
	mbsSessionIdParam := c.Query("mbs-session-id")
	if mbsSessionIdParam == "" {
		problemDetail := models.ProblemDetails{
			Status: http.StatusBadRequest,
			Cause:  "MISSING_REQUIRED_PARAMETER",
		}
		c.JSON(http.StatusBadRequest, problemDetail)
		return
	}

	var mbsSessionId *models.MbsSessionId
	if err := json.Unmarshal([]byte(mbsSessionIdParam), &mbsSessionId); err != nil {
		problemDetail := models.ProblemDetails{
			Status: http.StatusBadRequest,
			Cause:  "MALFORMED_REQUEST",
		}
		c.JSON(http.StatusBadRequest, problemDetail)
		return
	}

	// Query bindings
	bindings := bsfContext.BsfSelf.QueryPcfMbsBindings(mbsSessionId)

	if len(bindings) == 0 {
		c.JSON(http.StatusOK, []models.PcfMbsBinding{})
		return
	}

	// Convert to response format
	var response []models.PcfMbsBinding
	for _, binding := range bindings {
		mbsBinding := models.PcfMbsBinding{
			MbsSessionId:   binding.MbsSessionId,
			PcfFqdn:        ptrToString(binding.PcfFqdn),
			PcfIpEndPoints: binding.PcfIpEndPoints,
			PcfId:          ptrToString(binding.PcfId),
			PcfSetId:       ptrToString(binding.PcfSetId),
			BindLevel:      (*binding.BindLevel),
			SuppFeat:       ptrToString(binding.SuppFeat),
		}

		if binding.RecoveryTime != nil {
			mbsBinding.RecoveryTime = binding.RecoveryTime
		}

		response = append(response, mbsBinding)
	}

	c.JSON(http.StatusOK, response)
}

// ModifyIndPCFMbsBinding handles PATCH /pcf-mbs-bindings/{bindingId}
func ModifyIndPCFMbsBinding(c *gin.Context) {
	logger.ProcLog.Infof("Handle ModifyIndPCFMbsBinding")

	bindingId := c.Param("bindingId")

	var patchRequest models.PcfMbsBindingPatch
	if err := c.ShouldBindJSON(&patchRequest); err != nil {
		problemDetail := models.ProblemDetails{
			Status: http.StatusBadRequest,
			Cause:  "MALFORMED_REQUEST",
		}
		c.JSON(http.StatusBadRequest, problemDetail)
		return
	}

	binding, exists := bsfContext.BsfSelf.GetPcfMbsBinding(bindingId)
	if !exists {
		problemDetail := models.ProblemDetails{
			Status: http.StatusNotFound,
			Cause:  "RESOURCE_NOT_FOUND",
		}
		c.JSON(http.StatusNotFound, problemDetail)
		return
	}

	// Apply patch
	if patchRequest.PcfFqdn != "" {
		binding.PcfFqdn = stringToPtr(patchRequest.PcfFqdn)
	}
	if patchRequest.PcfIpEndPoints != nil {
		binding.PcfIpEndPoints = patchRequest.PcfIpEndPoints
	}
	if patchRequest.PcfId != "" {
		binding.PcfId = stringToPtr(patchRequest.PcfId)
	}

	// Update binding
	bsfContext.BsfSelf.UpdatePcfMbsBinding(bindingId, binding)

	// Return updated binding
	response := models.PcfMbsBinding{
		MbsSessionId:   binding.MbsSessionId,
		PcfFqdn:        ptrToString(binding.PcfFqdn),
		PcfIpEndPoints: binding.PcfIpEndPoints,
		PcfId:          ptrToString(binding.PcfId),
		PcfSetId:       ptrToString(binding.PcfSetId),
		BindLevel:      (*binding.BindLevel),
		SuppFeat:       ptrToString(binding.SuppFeat),
	}

	if binding.RecoveryTime != nil {
		response.RecoveryTime = binding.RecoveryTime
	}

	c.JSON(http.StatusOK, response)
}

// DeleteIndPCFMbsBinding handles DELETE /pcf-mbs-bindings/{bindingId}
func DeleteIndPCFMbsBinding(c *gin.Context) {
	logger.ProcLog.Infof("Handle DeleteIndPCFMbsBinding")

	bindingId := c.Param("bindingId")

	if bsfContext.BsfSelf.DeletePcfMbsBinding(bindingId) {
		c.Status(http.StatusNoContent)
	} else {
		problemDetail := models.ProblemDetails{
			Status: http.StatusNotFound,
			Cause:  "RESOURCE_NOT_FOUND",
		}
		c.JSON(http.StatusNotFound, problemDetail)
	}
}
