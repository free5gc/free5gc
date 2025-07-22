package util

import (
	"fmt"

	"github.com/go-playground/validator/v10"

	"github.com/free5gc/nssf/internal/logger"
	"github.com/free5gc/openapi/models"
)

func BindErrorInvalidParamsMessages(err error) []models.InvalidParam {
	var invalidParams []models.InvalidParam
	if errs, ok := err.(validator.ValidationErrors); ok {
		for _, e := range errs {
			ip := models.InvalidParam{
				Param: e.Field(),
			}

			switch e.Tag() {
			case "required":
				ip.Reason = fmt.Sprintf("The `%s` field is required.", e.Field())
			case "oneof":
				ip.Reason = fmt.Sprintf("The `%s` field must be one of '%s'.", e.Field(), e.Param())
			case "required_with":
				ip.Reason = fmt.Sprintf("The `%s` field is required when `%s` is present.", e.Field(), e.Param())
			case "required_without":
				ip.Reason = fmt.Sprintf("The `%s` field is required when `%s` is not present.", e.Field(), e.Param())
			case "uuid":
				ip.Reason = fmt.Sprintf("The `%s` field must be a valid UUID.", e.Field())
			default:
				ip.Reason = fmt.Sprintf("Failed on the `%s` tag.", e.Tag())
			}

			invalidParams = append(invalidParams, ip)
		}
	} else {
		logger.NsselLog.Errorf("Unknown error type: %+v", err)
	}

	return invalidParams
}
