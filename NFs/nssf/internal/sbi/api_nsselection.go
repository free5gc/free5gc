package sbi

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/free5gc/nssf/internal/logger"
	"github.com/free5gc/nssf/internal/sbi/processor"
	"github.com/free5gc/nssf/internal/util"
	"github.com/free5gc/openapi/models"
)

func (s *Server) getNsSelectionRoutes() []Route {
	return []Route{
		{
			"Health Check",
			http.MethodGet,
			"/",
			func(ctx *gin.Context) {
				ctx.JSON(http.StatusOK, gin.H{"status": "Service Available"})
			},
		},

		{
			"NSSelectionGet",
			http.MethodGet,
			"/network-slice-information",
			s.NetworkSliceInformationGet,
		},
	}
}

func (s *Server) NetworkSliceInformationGet(c *gin.Context) {
	logger.NsselLog.Infof("Handle NSSelectionGet")

	var query processor.NetworkSliceInformationGetQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		logger.NsselLog.Errorf("BindQuery failed: %+v", err)
		problemDetail := &models.ProblemDetails{
			Title:         "Malformed Request",
			Status:        http.StatusBadRequest,
			Detail:        err.Error(),
			Instance:      "",
			InvalidParams: util.BindErrorInvalidParamsMessages(err),
		}
		util.GinProblemJson(c, problemDetail)
		return
	}

	s.Processor().NSSelectionSliceInformationGet(c, query)
}
