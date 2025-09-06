package sbi

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/free5gc/bsf/internal/sbi/processor"
)

// getManagementRoutes returns all BSF Management API routes following 3GPP TS 29.521
func getManagementRoutes() Routes {
	return Routes{
		{
			Name:    "Index",
			Method:  http.MethodGet,
			Pattern: "/",
			HandlerFunc: func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"status": "BSF Management Service Available"})
			},
		},
		// PCF Bindings Collection
		{
			Name:        "CreatePCFBinding",
			Method:      http.MethodPost,
			Pattern:     "/pcfBindings",
			HandlerFunc: processor.CreatePCFBinding,
		},
		{
			Name:        "GetPCFBindings",
			Method:      http.MethodGet,
			Pattern:     "/pcfBindings",
			HandlerFunc: processor.GetPCFBindings,
		},
		// Individual PCF Binding Document
		{
			Name:        "GetIndPCFBinding",
			Method:      http.MethodGet,
			Pattern:     "/pcfBindings/:bindingId",
			HandlerFunc: processor.GetIndPCFBinding,
		},
		{
			Name:        "UpdateIndPCFBinding",
			Method:      http.MethodPatch,
			Pattern:     "/pcfBindings/:bindingId",
			HandlerFunc: processor.UpdateIndPCFBinding,
		},
		{
			Name:        "DeleteIndPCFBinding",
			Method:      http.MethodDelete,
			Pattern:     "/pcfBindings/:bindingId",
			HandlerFunc: processor.DeleteIndPCFBinding,
		},
		// PCF UE Bindings Collection
		{
			Name:        "CreatePCFforUEBinding",
			Method:      http.MethodPost,
			Pattern:     "/pcf-ue-bindings",
			HandlerFunc: processor.CreatePCFforUEBinding,
		},
		{
			Name:        "GetPCFForUeBindings",
			Method:      http.MethodGet,
			Pattern:     "/pcf-ue-bindings",
			HandlerFunc: processor.GetPCFForUeBindings,
		},
		// Individual PCF UE Binding Document
		{
			Name:        "UpdateIndPCFforUEBinding",
			Method:      http.MethodPatch,
			Pattern:     "/pcf-ue-bindings/:bindingId",
			HandlerFunc: processor.UpdateIndPCFforUEBinding,
		},
		{
			Name:        "DeleteIndPCFforUEBinding",
			Method:      http.MethodDelete,
			Pattern:     "/pcf-ue-bindings/:bindingId",
			HandlerFunc: processor.DeleteIndPCFforUEBinding,
		},
		// PCF MBS Bindings Collection
		{
			Name:        "CreatePCFMbsBinding",
			Method:      http.MethodPost,
			Pattern:     "/pcf-mbs-bindings",
			HandlerFunc: processor.CreatePCFMbsBinding,
		},
		{
			Name:        "GetPCFMbsBinding",
			Method:      http.MethodGet,
			Pattern:     "/pcf-mbs-bindings",
			HandlerFunc: processor.GetPCFMbsBinding,
		},
		// Individual PCF MBS Binding Document
		{
			Name:        "ModifyIndPCFMbsBinding",
			Method:      http.MethodPatch,
			Pattern:     "/pcf-mbs-bindings/:bindingId",
			HandlerFunc: processor.ModifyIndPCFMbsBinding,
		},
		{
			Name:        "DeleteIndPCFMbsBinding",
			Method:      http.MethodDelete,
			Pattern:     "/pcf-mbs-bindings/:bindingId",
			HandlerFunc: processor.DeleteIndPCFMbsBinding,
		},
		// Subscriptions Collection
		{
			Name:        "CreateIndividualSubcription",
			Method:      http.MethodPost,
			Pattern:     "/subscriptions",
			HandlerFunc: processor.CreateIndividualSubcription,
		},
		// Individual Subscription Document
		{
			Name:        "ReplaceIndividualSubcription",
			Method:      http.MethodPut,
			Pattern:     "/subscriptions/:subId",
			HandlerFunc: processor.ReplaceIndividualSubcription,
		},
		{
			Name:        "DeleteIndividualSubcription",
			Method:      http.MethodDelete,
			Pattern:     "/subscriptions/:subId",
			HandlerFunc: processor.DeleteIndividualSubcription,
		},
	}
}
