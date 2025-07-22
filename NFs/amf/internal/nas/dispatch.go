package nas

import (
	"errors"
	"fmt"

	"github.com/free5gc/amf/internal/context"
	"github.com/free5gc/amf/internal/gmm"
	"github.com/free5gc/amf/internal/logger"
	"github.com/free5gc/nas"
	"github.com/free5gc/openapi/models"
	"github.com/free5gc/util/fsm"
)

func Dispatch(ue *context.AmfUe, accessType models.AccessType, procedureCode int64, msg *nas.Message) error {
	if msg.GmmMessage == nil {
		return errors.New("gmm Message is nil")
	}

	if msg.GsmMessage != nil {
		return errors.New("GSM Message should include in GMM Message")
	}

	if ue.State[accessType] == nil {
		return fmt.Errorf("UE State is empty (accessType=%q). Can't send GSM Message", accessType)
	}

	return gmm.GmmFSM.SendEvent(ue.State[accessType], gmm.GmmMessageEvent, fsm.ArgsType{
		gmm.ArgAmfUe:         ue,
		gmm.ArgAccessType:    accessType,
		gmm.ArgNASMessage:    msg.GmmMessage,
		gmm.ArgProcedureCode: procedureCode,
	}, logger.GmmLog)
}
