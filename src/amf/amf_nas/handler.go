package amf_nas

import (
	"free5gc/lib/nas"
	"free5gc/src/amf/amf_context"
	"free5gc/src/amf/amf_nas/nas_security"
	"free5gc/src/amf/logger"
)

func HandleNAS(ue *amf_context.RanUe, procedureCode int64, nasPdu []byte) {

	if ue == nil {
		logger.NasLog.Error("RanUe is nil")
		return
	}

	if nasPdu == nil {
		logger.NasLog.Error("nasPdu is nil")
		return
	}

	var msg *nas.Message

	if ue.AmfUe != nil {
		var err error
		msg, err = nas_security.Decode(ue.AmfUe, nas.GetSecurityHeaderType(nasPdu)&0x0f, nasPdu)
		if err != nil {
			logger.NasLog.Error(err.Error())
			return
		}
	} else {
		msg = new(nas.Message)
		err := msg.PlainNasDecode(&nasPdu)
		if err != nil {
			logger.NasLog.Error(err.Error())
			return
		}
	}

	if err := Dispatch(ue.AmfUe, ue.Ran.AnType, procedureCode, msg); err != nil {
		logger.NgapLog.Errorf("Send to Nas Error: %v", err)
	}
}
