/*
 * BSF Management Processor
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

// Helper functions for type conversion
func stringToPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func ptrToString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

// CreatePCFBinding handles POST /pcfBindings
func CreatePCFBinding(c *gin.Context) {
	logger.ProcLog.Infof("Handle CreatePCFBinding")

	var request models.PcfBinding
	if err := c.ShouldBindJSON(&request); err != nil {
		problemDetail := models.ProblemDetails{
			Status: http.StatusBadRequest,
			Cause:  "MALFORMED_REQUEST",
		}
		c.JSON(http.StatusBadRequest, problemDetail)
		return
	}

	// Validate required fields per 3GPP TS 29.521
	if request.Dnn == "" {
		problemDetail := models.ProblemDetails{
			Status: http.StatusBadRequest,
			Cause:  "MANDATORY_IE_MISSING",
			Detail: "dnn is required",
		}
		c.JSON(http.StatusBadRequest, problemDetail)
		return
	}

	if request.Snssai == nil {
		problemDetail := models.ProblemDetails{
			Status: http.StatusBadRequest,
			Cause:  "MANDATORY_IE_MISSING",
			Detail: "snssai is required",
		}
		c.JSON(http.StatusBadRequest, problemDetail)
		return
	}

	// Convert to internal representation
	binding := &bsfContext.PcfBinding{
		Supi:               stringToPtr(request.Supi),
		Gpsi:               stringToPtr(request.Gpsi),
		Ipv4Addr:           stringToPtr(request.Ipv4Addr),
		Ipv6Prefix:         stringToPtr(request.Ipv6Prefix),
		AddIpv6Prefixes:    request.AddIpv6Prefixes,
		IpDomain:           stringToPtr(request.IpDomain),
		MacAddr48:          stringToPtr(request.MacAddr48),
		AddMacAddrs:        request.AddMacAddrs,
		Dnn:                request.Dnn,
		PcfFqdn:            stringToPtr(request.PcfFqdn),
		PcfIpEndPoints:     request.PcfIpEndPoints,
		PcfDiamHost:        stringToPtr(request.PcfDiamHost),
		PcfDiamRealm:       stringToPtr(request.PcfDiamRealm),
		PcfSmFqdn:          stringToPtr(request.PcfSmFqdn),
		PcfSmIpEndPoints:   request.PcfSmIpEndPoints,
		Snssai:             request.Snssai,
		SuppFeat:           stringToPtr(request.SuppFeat),
		PcfId:              stringToPtr(request.PcfId),
		PcfSetId:           stringToPtr(request.PcfSetId),
		ParaCom:            request.ParaCom,
		BindLevel:          (*models.BindingLevel)(&request.BindLevel),
		Ipv4FrameRouteList: request.Ipv4FrameRouteList,
		Ipv6FrameRouteList: request.Ipv6FrameRouteList,
	}

	if request.RecoveryTime != nil {
		binding.RecoveryTime = request.RecoveryTime
	}

	// Check for existing binding based on combination
	existingBindings := bsfContext.BsfSelf.QueryPcfBindings(
		request.Supi,
		request.Gpsi,
		request.Dnn,
		request.Ipv4Addr,
		request.Ipv6Prefix,
		request.MacAddr48,
		request.IpDomain,
		request.Snssai,
	)

	if len(existingBindings) > 0 {
		// Return 403 with existing PCF binding info per 3GPP TS 29.521
		problemDetail := models.ProblemDetails{
			Status: http.StatusForbidden,
			Cause:  "BINDING_ALREADY_EXISTS",
			Detail: "Existing PCF binding information stored in BSF for the indicated combination",
		}
		c.JSON(http.StatusForbidden, problemDetail)
		return
	}

	// Create new binding
	bindingId := bsfContext.BsfSelf.CreatePcfBinding(binding)

	// Convert back to response format
	response := models.PcfBinding{
		Supi:               ptrToString(binding.Supi),
		Gpsi:               ptrToString(binding.Gpsi),
		Ipv4Addr:           ptrToString(binding.Ipv4Addr),
		Ipv6Prefix:         ptrToString(binding.Ipv6Prefix),
		AddIpv6Prefixes:    binding.AddIpv6Prefixes,
		IpDomain:           ptrToString(binding.IpDomain),
		MacAddr48:          ptrToString(binding.MacAddr48),
		AddMacAddrs:        binding.AddMacAddrs,
		Dnn:                binding.Dnn,
		PcfFqdn:            ptrToString(binding.PcfFqdn),
		PcfIpEndPoints:     binding.PcfIpEndPoints,
		PcfDiamHost:        ptrToString(binding.PcfDiamHost),
		PcfDiamRealm:       ptrToString(binding.PcfDiamRealm),
		PcfSmFqdn:          ptrToString(binding.PcfSmFqdn),
		PcfSmIpEndPoints:   binding.PcfSmIpEndPoints,
		Snssai:             binding.Snssai,
		SuppFeat:           ptrToString(binding.SuppFeat),
		PcfId:              ptrToString(binding.PcfId),
		PcfSetId:           ptrToString(binding.PcfSetId),
		ParaCom:            binding.ParaCom,
		BindLevel:          (*binding.BindLevel),
		Ipv4FrameRouteList: binding.Ipv4FrameRouteList,
		Ipv6FrameRouteList: binding.Ipv6FrameRouteList,
	}

	if binding.RecoveryTime != nil {
		response.RecoveryTime = binding.RecoveryTime
	}

	locationHeader := "/nbsf-management/v1/pcfBindings/" + bindingId
	c.Header("Location", locationHeader)
	c.JSON(http.StatusCreated, response)
}

// GetPCFBindings handles GET /pcfBindings
func GetPCFBindings(c *gin.Context) {
	logger.ProcLog.Infof("Handle GetPCFBindings")

	// Extract query parameters
	ipv4Addr := c.Query("ipv4Addr")
	ipv6Prefix := c.Query("ipv6Prefix")
	macAddr48 := c.Query("macAddr48")
	dnn := c.Query("dnn")
	supi := c.Query("supi")
	gpsi := c.Query("gpsi")
	ipDomain := c.Query("ipDomain")

	var snssai *models.Snssai
	if snssaiParam := c.Query("snssai"); snssaiParam != "" {
		if err := json.Unmarshal([]byte(snssaiParam), &snssai); err != nil {
			problemDetail := models.ProblemDetails{
				Status: http.StatusBadRequest,
				Cause:  "MALFORMED_REQUEST",
			}
			c.JSON(http.StatusBadRequest, problemDetail)
			return
		}
	}

	// Query bindings
	bindings := bsfContext.BsfSelf.QueryPcfBindings(supi, gpsi, dnn, ipv4Addr, ipv6Prefix, macAddr48, ipDomain, snssai)

	if len(bindings) == 0 {
		c.Status(http.StatusNoContent)
		return
	}

	// Convert first match to response format
	binding := bindings[0]
	response := models.PcfBinding{
		Supi:               ptrToString(binding.Supi),
		Gpsi:               ptrToString(binding.Gpsi),
		Ipv4Addr:           ptrToString(binding.Ipv4Addr),
		Ipv6Prefix:         ptrToString(binding.Ipv6Prefix),
		AddIpv6Prefixes:    binding.AddIpv6Prefixes,
		IpDomain:           ptrToString(binding.IpDomain),
		MacAddr48:          ptrToString(binding.MacAddr48),
		AddMacAddrs:        binding.AddMacAddrs,
		Dnn:                binding.Dnn,
		PcfFqdn:            ptrToString(binding.PcfFqdn),
		PcfIpEndPoints:     binding.PcfIpEndPoints,
		PcfDiamHost:        ptrToString(binding.PcfDiamHost),
		PcfDiamRealm:       ptrToString(binding.PcfDiamRealm),
		PcfSmFqdn:          ptrToString(binding.PcfSmFqdn),
		PcfSmIpEndPoints:   binding.PcfSmIpEndPoints,
		Snssai:             binding.Snssai,
		SuppFeat:           ptrToString(binding.SuppFeat),
		PcfId:              ptrToString(binding.PcfId),
		PcfSetId:           ptrToString(binding.PcfSetId),
		ParaCom:            binding.ParaCom,
		BindLevel:          (*binding.BindLevel),
		Ipv4FrameRouteList: binding.Ipv4FrameRouteList,
		Ipv6FrameRouteList: binding.Ipv6FrameRouteList,
	}

	if binding.RecoveryTime != nil {
		response.RecoveryTime = binding.RecoveryTime
	}

	c.JSON(http.StatusOK, response)
}

// DeleteIndPCFBinding handles DELETE /pcfBindings/{bindingId}
func DeleteIndPCFBinding(c *gin.Context) {
	logger.ProcLog.Infof("Handle DeleteIndPCFBinding")

	bindingId := c.Param("bindingId")

	if bsfContext.BsfSelf.DeletePcfBinding(bindingId) {
		c.Status(http.StatusNoContent)
	} else {
		problemDetail := models.ProblemDetails{
			Status: http.StatusNotFound,
			Cause:  "RESOURCE_NOT_FOUND",
		}
		c.JSON(http.StatusNotFound, problemDetail)
	}
}

// UpdateIndPCFBinding handles PATCH /pcfBindings/{bindingId}
func UpdateIndPCFBinding(c *gin.Context) {
	logger.ProcLog.Infof("Handle UpdateIndPCFBinding")

	bindingId := c.Param("bindingId")

	var patchRequest models.PcfBindingPatch
	if err := c.ShouldBindJSON(&patchRequest); err != nil {
		problemDetail := models.ProblemDetails{
			Status: http.StatusBadRequest,
			Cause:  "MALFORMED_REQUEST",
		}
		c.JSON(http.StatusBadRequest, problemDetail)
		return
	}

	binding, exists := bsfContext.BsfSelf.GetPcfBinding(bindingId)
	if !exists {
		problemDetail := models.ProblemDetails{
			Status: http.StatusNotFound,
			Cause:  "RESOURCE_NOT_FOUND",
		}
		c.JSON(http.StatusNotFound, problemDetail)
		return
	}

	// Apply patch
	if patchRequest.Ipv4Addr != "" {
		binding.Ipv4Addr = stringToPtr(patchRequest.Ipv4Addr)
	}
	if patchRequest.IpDomain != "" {
		binding.IpDomain = stringToPtr(patchRequest.IpDomain)
	}
	if patchRequest.Ipv6Prefix != "" {
		binding.Ipv6Prefix = stringToPtr(patchRequest.Ipv6Prefix)
	}
	if patchRequest.AddIpv6Prefixes != nil {
		binding.AddIpv6Prefixes = patchRequest.AddIpv6Prefixes
	}
	if patchRequest.MacAddr48 != "" {
		binding.MacAddr48 = stringToPtr(patchRequest.MacAddr48)
	}
	if patchRequest.AddMacAddrs != nil {
		binding.AddMacAddrs = patchRequest.AddMacAddrs
	}
	if patchRequest.PcfId != "" {
		binding.PcfId = stringToPtr(patchRequest.PcfId)
	}
	if patchRequest.PcfFqdn != "" {
		binding.PcfFqdn = stringToPtr(patchRequest.PcfFqdn)
	}
	if patchRequest.PcfIpEndPoints != nil {
		binding.PcfIpEndPoints = patchRequest.PcfIpEndPoints
	}
	if patchRequest.PcfDiamHost != "" {
		binding.PcfDiamHost = stringToPtr(patchRequest.PcfDiamHost)
	}
	if patchRequest.PcfDiamRealm != "" {
		binding.PcfDiamRealm = stringToPtr(patchRequest.PcfDiamRealm)
	}

	// Update binding
	bsfContext.BsfSelf.UpdatePcfBinding(bindingId, binding)

	// Return updated binding
	response := models.PcfBinding{
		Supi:               ptrToString(binding.Supi),
		Gpsi:               ptrToString(binding.Gpsi),
		Ipv4Addr:           ptrToString(binding.Ipv4Addr),
		Ipv6Prefix:         ptrToString(binding.Ipv6Prefix),
		AddIpv6Prefixes:    binding.AddIpv6Prefixes,
		IpDomain:           ptrToString(binding.IpDomain),
		MacAddr48:          ptrToString(binding.MacAddr48),
		AddMacAddrs:        binding.AddMacAddrs,
		Dnn:                binding.Dnn,
		PcfFqdn:            ptrToString(binding.PcfFqdn),
		PcfIpEndPoints:     binding.PcfIpEndPoints,
		PcfDiamHost:        ptrToString(binding.PcfDiamHost),
		PcfDiamRealm:       ptrToString(binding.PcfDiamRealm),
		PcfSmFqdn:          ptrToString(binding.PcfSmFqdn),
		PcfSmIpEndPoints:   binding.PcfSmIpEndPoints,
		Snssai:             binding.Snssai,
		SuppFeat:           ptrToString(binding.SuppFeat),
		PcfId:              ptrToString(binding.PcfId),
		PcfSetId:           ptrToString(binding.PcfSetId),
		ParaCom:            binding.ParaCom,
		BindLevel:          (*binding.BindLevel),
		Ipv4FrameRouteList: binding.Ipv4FrameRouteList,
		Ipv6FrameRouteList: binding.Ipv6FrameRouteList,
	}

	if binding.RecoveryTime != nil {
		response.RecoveryTime = binding.RecoveryTime
	}

	c.JSON(http.StatusOK, response)
}

// Helper function to safely get string value from pointer
func getStringValue(ptr *string) string {
	if ptr != nil {
		return *ptr
	}
	return ""
}
