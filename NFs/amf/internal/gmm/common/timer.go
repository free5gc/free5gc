package common

import (
	"github.com/free5gc/amf/internal/context"
	"github.com/free5gc/openapi/models"
)

func StopAll5GSMMTimers(ue *context.AmfUe) {
	if ue.T3513 != nil {
		ue.T3513.Stop()
		ue.T3513 = nil // clear the timer
		if ue.OnGoing(models.AccessType__3_GPP_ACCESS).Procedure == context.OnGoingProcedurePaging {
			ue.SetOnGoing(models.AccessType__3_GPP_ACCESS, &context.OnGoing{
				Procedure: context.OnGoingProcedureNothing,
			})
		}
	}
	if ue.T3522 != nil {
		ue.T3522.Stop()
		ue.T3522 = nil // clear the timer
	}
	if ue.T3550 != nil {
		ue.T3550.Stop()
		ue.T3550 = nil // clear the timer
	}
	if ue.T3560 != nil {
		ue.T3560.Stop()
		ue.T3560 = nil // clear the timer
	}
	if ue.T3565 != nil {
		ue.T3565.Stop()
		ue.T3565 = nil // clear the timer
	}
	if ue.T3570 != nil {
		ue.T3570.Stop()
		ue.T3570 = nil // clear the timer
	}
	if ue.T3555 != nil {
		ue.T3555.Stop()
		ue.T3555 = nil // clear the timer
	}
}
