package util

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/free5gc/openapi/models"
)

func GinProblemJson(c *gin.Context, pd *models.ProblemDetails) {
	c.JSON(int(pd.Status), pd)
	c.Writer.Header().Set("Content-Type", "application/problem+json")
}

func EmptyUeIdProblemJson(c *gin.Context) {
	problemDetail := &models.ProblemDetails{
		Title:  MALFORMED_REQUEST,
		Status: http.StatusBadRequest,
		Detail: "ueId is required",
	}
	GinProblemJson(c, problemDetail)
}
