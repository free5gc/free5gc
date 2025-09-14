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
	"github.com/free5gc/bsf/internal/metrics/business"
	"github.com/free5gc/bsf/internal/util"
	"github.com/free5gc/openapi/models"
)

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
	} // Validate required fields per 3GPP TS 29.521
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
		Supi:               util.StringToPtr(request.Supi),
		Gpsi:               util.StringToPtr(request.Gpsi),
		Ipv4Addr:           util.StringToPtr(request.Ipv4Addr),
		Ipv6Prefix:         util.StringToPtr(request.Ipv6Prefix),
		AddIpv6Prefixes:    request.AddIpv6Prefixes,
		IpDomain:           util.StringToPtr(request.IpDomain),
		MacAddr48:          util.StringToPtr(request.MacAddr48),
		AddMacAddrs:        request.AddMacAddrs,
		Dnn:                request.Dnn,
		PcfFqdn:            util.StringToPtr(request.PcfFqdn),
		PcfIpEndPoints:     request.PcfIpEndPoints,
		PcfDiamHost:        util.StringToPtr(request.PcfDiamHost),
		PcfDiamRealm:       util.StringToPtr(request.PcfDiamRealm),
		PcfSmFqdn:          util.StringToPtr(request.PcfSmFqdn),
		PcfSmIpEndPoints:   request.PcfSmIpEndPoints,
		Snssai:             request.Snssai,
		SuppFeat:           util.StringToPtr(request.SuppFeat),
		PcfId:              util.StringToPtr(request.PcfId),
		PcfSetId:           util.StringToPtr(request.PcfSetId),
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

	// Update metrics
	business.IncrPCFBindingGauge(business.PCF_BINDING_TYPE_VALUE)
	business.IncrPCFBindingEventCounter(business.PCF_BINDING_TYPE_VALUE, business.BINDING_EVENT_CREATE_VALUE, business.RESULT_SUCCESS_VALUE)

	// Convert back to response format
	response := models.PcfBinding{
		Supi:               util.PtrToString(binding.Supi),
		Gpsi:               util.PtrToString(binding.Gpsi),
		Ipv4Addr:           util.PtrToString(binding.Ipv4Addr),
		Ipv6Prefix:         util.PtrToString(binding.Ipv6Prefix),
		AddIpv6Prefixes:    binding.AddIpv6Prefixes,
		IpDomain:           util.PtrToString(binding.IpDomain),
		MacAddr48:          util.PtrToString(binding.MacAddr48),
		AddMacAddrs:        binding.AddMacAddrs,
		Dnn:                binding.Dnn,
		PcfFqdn:            util.PtrToString(binding.PcfFqdn),
		PcfIpEndPoints:     binding.PcfIpEndPoints,
		PcfDiamHost:        util.PtrToString(binding.PcfDiamHost),
		PcfDiamRealm:       util.PtrToString(binding.PcfDiamRealm),
		PcfSmFqdn:          util.PtrToString(binding.PcfSmFqdn),
		PcfSmIpEndPoints:   binding.PcfSmIpEndPoints,
		Snssai:             binding.Snssai,
		SuppFeat:           util.PtrToString(binding.SuppFeat),
		PcfId:              util.PtrToString(binding.PcfId),
		PcfSetId:           util.PtrToString(binding.PcfSetId),
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

	// Validate conditional parameters as per 3GPP TS 29.521 Table 5.3.2.3.2-1
	// NOTE 1: One and only one of query parameter ipv4Addr, ipv6Prefix or macAddr48 shall be present
	ueAddressParams := []string{ipv4Addr, ipv6Prefix, macAddr48}
	nonEmptyParams := 0
	for _, param := range ueAddressParams {
		if param != "" {
			nonEmptyParams++
		}
	}

	if nonEmptyParams > 1 {
		problemDetail := models.ProblemDetails{
			Status: http.StatusBadRequest,
			Cause:  "MALFORMED_REQUEST",
			Detail: "Only one of ipv4Addr, ipv6Prefix, or macAddr48 parameters shall be present (3GPP TS 29.521 NOTE 1)",
		}
		c.JSON(http.StatusBadRequest, problemDetail)
		return
	}

	// NOTE 2: The query parameter ipDomain, if applicable, shall be present with query parameter ipv4Addr
	if ipDomain != "" && ipv4Addr == "" {
		problemDetail := models.ProblemDetails{
			Status: http.StatusBadRequest,
			Cause:  "MALFORMED_REQUEST",
			Detail: "ipDomain parameter shall be present only with ipv4Addr parameter (3GPP TS 29.521 NOTE 2)",
		}
		c.JSON(http.StatusBadRequest, problemDetail)
		return
	}

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

	// Update metrics
	if len(bindings) > 0 {
		business.IncrPCFBindingEventCounter(business.PCF_BINDING_TYPE_VALUE, business.BINDING_EVENT_QUERY_VALUE, business.RESULT_SUCCESS_VALUE)
	} else {
		business.IncrPCFBindingEventCounter(business.PCF_BINDING_TYPE_VALUE, business.BINDING_EVENT_QUERY_VALUE, business.RESULT_FAILURE_VALUE)
	}

	if len(bindings) == 0 {
		c.Status(http.StatusNoContent)
		return
	}

	// Convert first match to response format
	binding := bindings[0]
	response := models.PcfBinding{
		Supi:               util.PtrToString(binding.Supi),
		Gpsi:               util.PtrToString(binding.Gpsi),
		Ipv4Addr:           util.PtrToString(binding.Ipv4Addr),
		Ipv6Prefix:         util.PtrToString(binding.Ipv6Prefix),
		AddIpv6Prefixes:    binding.AddIpv6Prefixes,
		IpDomain:           util.PtrToString(binding.IpDomain),
		MacAddr48:          util.PtrToString(binding.MacAddr48),
		AddMacAddrs:        binding.AddMacAddrs,
		Dnn:                binding.Dnn,
		PcfFqdn:            util.PtrToString(binding.PcfFqdn),
		PcfIpEndPoints:     binding.PcfIpEndPoints,
		PcfDiamHost:        util.PtrToString(binding.PcfDiamHost),
		PcfDiamRealm:       util.PtrToString(binding.PcfDiamRealm),
		PcfSmFqdn:          util.PtrToString(binding.PcfSmFqdn),
		PcfSmIpEndPoints:   binding.PcfSmIpEndPoints,
		Snssai:             binding.Snssai,
		SuppFeat:           util.PtrToString(binding.SuppFeat),
		PcfId:              util.PtrToString(binding.PcfId),
		PcfSetId:           util.PtrToString(binding.PcfSetId),
		ParaCom:            binding.ParaCom,
		BindLevel:          (*binding.BindLevel),
		Ipv4FrameRouteList: binding.Ipv4FrameRouteList,
		Ipv6FrameRouteList: binding.Ipv6FrameRouteList,
	}

	// Add binding ID for reference (enhancement for query operations)
	c.Header("X-BSF-Binding-ID", binding.BindingId)

	if binding.RecoveryTime != nil {
		response.RecoveryTime = binding.RecoveryTime
	}

	c.JSON(http.StatusOK, response)
}

// GetIndPCFBinding handles GET /pcfBindings/{bindingId}
func GetIndPCFBinding(c *gin.Context) {
	logger.ProcLog.Infof("Handle GetIndPCFBinding")

	bindingId := c.Param("bindingId")
	if bindingId == "" {
		problemDetail := models.ProblemDetails{
			Status: http.StatusBadRequest,
			Cause:  "MALFORMED_REQUEST",
			Detail: "bindingId is required",
		}
		c.JSON(http.StatusBadRequest, problemDetail)
		return
	}

	// Get binding by ID
	binding, exists := bsfContext.BsfSelf.GetPcfBinding(bindingId)
	if !exists {
		problemDetail := models.ProblemDetails{
			Status: http.StatusNotFound,
			Cause:  "RESOURCE_NOT_FOUND",
			Detail: "PCF binding not found",
		}
		c.JSON(http.StatusNotFound, problemDetail)
		return
	}

	// Convert to response format
	response := models.PcfBinding{
		Supi:               util.PtrToString(binding.Supi),
		Gpsi:               util.PtrToString(binding.Gpsi),
		Ipv4Addr:           util.PtrToString(binding.Ipv4Addr),
		Ipv6Prefix:         util.PtrToString(binding.Ipv6Prefix),
		AddIpv6Prefixes:    binding.AddIpv6Prefixes,
		IpDomain:           util.PtrToString(binding.IpDomain),
		MacAddr48:          util.PtrToString(binding.MacAddr48),
		AddMacAddrs:        binding.AddMacAddrs,
		Dnn:                binding.Dnn,
		PcfFqdn:            util.PtrToString(binding.PcfFqdn),
		PcfIpEndPoints:     binding.PcfIpEndPoints,
		PcfDiamHost:        util.PtrToString(binding.PcfDiamHost),
		PcfDiamRealm:       util.PtrToString(binding.PcfDiamRealm),
		PcfSmFqdn:          util.PtrToString(binding.PcfSmFqdn),
		PcfSmIpEndPoints:   binding.PcfSmIpEndPoints,
		Snssai:             binding.Snssai,
		SuppFeat:           util.PtrToString(binding.SuppFeat),
		PcfId:              util.PtrToString(binding.PcfId),
		PcfSetId:           util.PtrToString(binding.PcfSetId),
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
		// Update metrics
		business.DecrPCFBindingGauge(business.PCF_BINDING_TYPE_VALUE)
		business.IncrPCFBindingEventCounter(business.PCF_BINDING_TYPE_VALUE, business.BINDING_EVENT_DELETE_VALUE, business.RESULT_SUCCESS_VALUE)
		c.Status(http.StatusNoContent)
	} else {
		business.IncrPCFBindingEventCounter(business.PCF_BINDING_TYPE_VALUE, business.BINDING_EVENT_DELETE_VALUE, business.RESULT_FAILURE_VALUE)
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
		business.IncrPCFBindingEventCounter(business.PCF_BINDING_TYPE_VALUE, business.BINDING_EVENT_UPDATE_VALUE, business.RESULT_FAILURE_VALUE)
		problemDetail := models.ProblemDetails{
			Status: http.StatusNotFound,
			Cause:  "RESOURCE_NOT_FOUND",
		}
		c.JSON(http.StatusNotFound, problemDetail)
		return
	}

	// Apply patch
	if patchRequest.Ipv4Addr != "" {
		binding.Ipv4Addr = util.StringToPtr(patchRequest.Ipv4Addr)
	}
	if patchRequest.IpDomain != "" {
		binding.IpDomain = util.StringToPtr(patchRequest.IpDomain)
	}
	if patchRequest.Ipv6Prefix != "" {
		binding.Ipv6Prefix = util.StringToPtr(patchRequest.Ipv6Prefix)
	}
	if patchRequest.AddIpv6Prefixes != nil {
		binding.AddIpv6Prefixes = patchRequest.AddIpv6Prefixes
	}
	if patchRequest.MacAddr48 != "" {
		binding.MacAddr48 = util.StringToPtr(patchRequest.MacAddr48)
	}
	if patchRequest.AddMacAddrs != nil {
		binding.AddMacAddrs = patchRequest.AddMacAddrs
	}
	if patchRequest.PcfId != "" {
		binding.PcfId = util.StringToPtr(patchRequest.PcfId)
	}
	if patchRequest.PcfFqdn != "" {
		binding.PcfFqdn = util.StringToPtr(patchRequest.PcfFqdn)
	}
	if patchRequest.PcfIpEndPoints != nil {
		binding.PcfIpEndPoints = patchRequest.PcfIpEndPoints
	}
	if patchRequest.PcfDiamHost != "" {
		binding.PcfDiamHost = util.StringToPtr(patchRequest.PcfDiamHost)
	}
	if patchRequest.PcfDiamRealm != "" {
		binding.PcfDiamRealm = util.StringToPtr(patchRequest.PcfDiamRealm)
	}

	// Update binding
	bsfContext.BsfSelf.UpdatePcfBinding(bindingId, binding)

	// Update metrics
	business.IncrPCFBindingEventCounter(business.PCF_BINDING_TYPE_VALUE, business.BINDING_EVENT_UPDATE_VALUE, business.RESULT_SUCCESS_VALUE)

	// Return updated binding
	response := models.PcfBinding{
		Supi:               util.PtrToString(binding.Supi),
		Gpsi:               util.PtrToString(binding.Gpsi),
		Ipv4Addr:           util.PtrToString(binding.Ipv4Addr),
		Ipv6Prefix:         util.PtrToString(binding.Ipv6Prefix),
		AddIpv6Prefixes:    binding.AddIpv6Prefixes,
		IpDomain:           util.PtrToString(binding.IpDomain),
		MacAddr48:          util.PtrToString(binding.MacAddr48),
		AddMacAddrs:        binding.AddMacAddrs,
		Dnn:                binding.Dnn,
		PcfFqdn:            util.PtrToString(binding.PcfFqdn),
		PcfIpEndPoints:     binding.PcfIpEndPoints,
		PcfDiamHost:        util.PtrToString(binding.PcfDiamHost),
		PcfDiamRealm:       util.PtrToString(binding.PcfDiamRealm),
		PcfSmFqdn:          util.PtrToString(binding.PcfSmFqdn),
		PcfSmIpEndPoints:   binding.PcfSmIpEndPoints,
		Snssai:             binding.Snssai,
		SuppFeat:           util.PtrToString(binding.SuppFeat),
		PcfId:              util.PtrToString(binding.PcfId),
		PcfSetId:           util.PtrToString(binding.PcfSetId),
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
