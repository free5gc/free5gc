package util

import (
	"net/http"

	"github.com/gin-gonic/gin"

	ausf_context "github.com/free5gc/ausf/internal/context"
	"github.com/free5gc/ausf/internal/logger"
	"github.com/free5gc/openapi/models"
)

type RouterAuthorizationCheck struct {
	serviceName models.ServiceName
}

func NewRouterAuthorizationCheck(serviceName models.ServiceName) *RouterAuthorizationCheck {
	return &RouterAuthorizationCheck{
		serviceName: serviceName,
	}
}

func (rac *RouterAuthorizationCheck) Check(c *gin.Context, ausfContext ausf_context.NFContext) {
	token := c.Request.Header.Get("Authorization")
	err := ausfContext.AuthorizationCheck(token, rac.serviceName)
	if err != nil {
		logger.UtilLog.Debugf("RouterAuthorizationCheck: Check Unauthorized: %s", err.Error())
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		c.Abort()
		return
	}

	logger.UtilLog.Debugf("RouterAuthorizationCheck: Check Authorized")
}
