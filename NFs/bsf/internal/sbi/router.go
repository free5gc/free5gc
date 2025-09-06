/*
 * BSF SBI Router
 */

package sbi

import (
	"github.com/gin-gonic/gin"

	"github.com/free5gc/bsf/internal/logger"
)

// AddService initializes the BSF SBI service with proper routing
// This function maintains backward compatibility with existing BSF initialization
func AddService(engine *gin.Engine) {
	// Apply BSF Management routes to the provided engine
	managementGroup := engine.Group("/nbsf-management/v1")
	managementRoutes := getManagementRoutes()
	applyRoutes(managementGroup, managementRoutes)

	logger.SbiLog.Infof("BSF SBI server initialized")
}
