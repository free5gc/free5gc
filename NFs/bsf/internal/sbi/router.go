/*
 * BSF SBI Router
 */

package sbi

import (
	"github.com/gin-gonic/gin"

	"github.com/free5gc/bsf/internal/logger"
	"github.com/free5gc/bsf/internal/sbi/processor"
)

func AddService(engine *gin.Engine) {
	managementGroup := engine.Group("/nbsf-management/v1")

	// PCF Bindings Collection
	managementGroup.POST("/pcfBindings", processor.CreatePCFBinding)
	managementGroup.GET("/pcfBindings", processor.GetPCFBindings)

	// Individual PCF Binding Document
	managementGroup.DELETE("/pcfBindings/:bindingId", processor.DeleteIndPCFBinding)
	managementGroup.PATCH("/pcfBindings/:bindingId", processor.UpdateIndPCFBinding)

	// PCF UE Bindings Collection
	managementGroup.POST("/pcf-ue-bindings", processor.CreatePCFforUEBinding)
	managementGroup.GET("/pcf-ue-bindings", processor.GetPCFForUeBindings)

	// Individual PCF UE Binding Document
	managementGroup.DELETE("/pcf-ue-bindings/:bindingId", processor.DeleteIndPCFforUEBinding)
	managementGroup.PATCH("/pcf-ue-bindings/:bindingId", processor.UpdateIndPCFforUEBinding)

	// PCF MBS Bindings Collection
	managementGroup.POST("/pcf-mbs-bindings", processor.CreatePCFMbsBinding)
	managementGroup.GET("/pcf-mbs-bindings", processor.GetPCFMbsBinding)

	// Individual PCF MBS Binding Document
	managementGroup.PATCH("/pcf-mbs-bindings/:bindingId", processor.ModifyIndPCFMbsBinding)
	managementGroup.DELETE("/pcf-mbs-bindings/:bindingId", processor.DeleteIndPCFMbsBinding)

	// Subscriptions Collection
	managementGroup.POST("/subscriptions", processor.CreateIndividualSubcription)

	// Individual Subscription Document
	managementGroup.PUT("/subscriptions/:subId", processor.ReplaceIndividualSubcription)
	managementGroup.DELETE("/subscriptions/:subId", processor.DeleteIndividualSubcription)

	logger.SbiLog.Infof("BSF SBI server initialized")
}
