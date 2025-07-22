package context

import (
	"github.com/free5gc/openapi/models"
)

// TrafficControlData - Traffic control data defines how traffic data flows
// associated with a rule are treated (e.g. blocked, redirected).
type TrafficControlData struct {
	*models.TrafficControlData
}

// NewTrafficControlData - create the traffic control data from OpenAPI model
func NewTrafficControlData(model *models.TrafficControlData) *TrafficControlData {
	if model == nil {
		return nil
	}

	return &TrafficControlData{
		TrafficControlData: model,
	}
}
