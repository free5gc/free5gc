package gmm

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"free5gc/lib/fsm"
	"free5gc/lib/nas"
	"free5gc/lib/nas/nasMessage"
	"free5gc/lib/openapi/models"
	"free5gc/src/amf/amf_context"
	"free5gc/src/amf/gmm/gmm_event"
	"free5gc/src/amf/gmm/gmm_handler"
	"free5gc/src/amf/logger"
)

var GmmLog *logrus.Entry

func init() {
	GmmLog = logger.GmmLog
}

func DeRegistered_3gpp(sm *fsm.FSM, event fsm.Event, args fsm.Args) error {
	return register_event_3gpp(sm, event, args)
}
func Registered_3gpp(sm *fsm.FSM, event fsm.Event, args fsm.Args) error {
	return register_event_3gpp(sm, event, args)
}

func register_event_3gpp(sm *fsm.FSM, event fsm.Event, args fsm.Args) error {
	var amfUe *amf_context.AmfUe
	var procedureCode int64
	switch event {
	case fsm.EVENT_ENTRY:
		return nil
	case gmm_event.EVENT_GMM_MESSAGE:
		amfUe = args[gmm_event.AMF_UE].(*amf_context.AmfUe)
		procedureCode = args[gmm_event.PROCEDURE_CODE].(int64)
		gmmMessage := args[gmm_event.GMM_MESSAGE].(*nas.GmmMessage)
		switch gmmMessage.GetMessageType() {
		case nas.MsgTypeULNASTransport:
			return gmm_handler.HandleULNASTransport(amfUe, models.AccessType__3_GPP_ACCESS, gmmMessage.ULNASTransport)
		case nas.MsgTypeRegistrationRequest:
			if err := gmm_handler.HandleRegistrationRequest(amfUe, models.AccessType__3_GPP_ACCESS, gmmMessage.RegistrationRequest); err != nil {
				return err
			}
		case nas.MsgTypeIdentityResponse:
			if err := gmm_handler.HandleIdentityResponse(amfUe, gmmMessage.IdentityResponse); err != nil {
				return err
			}
		case nas.MsgTypeConfigurationUpdateComplete:
			if err := gmm_handler.HandleConfigurationUpdateComplete(amfUe, gmmMessage.ConfigurationUpdateComplete); err != nil {
				return err
			}
		case nas.MsgTypeServiceRequest:
			if err := gmm_handler.HandleServiceRequest(amfUe, models.AccessType__3_GPP_ACCESS, procedureCode, gmmMessage.ServiceRequest); err != nil {
				return err
			}
		case nas.MsgTypeDeregistrationRequestUEOriginatingDeregistration:
			return gmm_handler.HandleDeregistrationRequest(amfUe, models.AccessType__3_GPP_ACCESS, gmmMessage.DeregistrationRequestUEOriginatingDeregistration)
		case nas.MsgTypeDeregistrationAcceptUETerminatedDeregistration:
			return gmm_handler.HandleDeregistrationAccept(amfUe, models.AccessType__3_GPP_ACCESS, gmmMessage.DeregistrationAcceptUETerminatedDeregistration)
		case nas.MsgTypeStatus5GMM:
			if err := gmm_handler.HandleStatus5GMM(amfUe, models.AccessType__3_GPP_ACCESS, gmmMessage.Status5GMM); err != nil {
				return err
			}
		default:
			GmmLog.Errorf("Unknown GmmMessage[%d]\n", gmmMessage.GetMessageType())
		}
	default:
		return fmt.Errorf("Unknown Event[%s]\n", event)
	}

	switch amfUe.RegistrationType5GS {
	case nasMessage.RegistrationType5GSInitialRegistration:
		return gmm_handler.HandleInitialRegistration(amfUe, models.AccessType__3_GPP_ACCESS)
	case nasMessage.RegistrationType5GSMobilityRegistrationUpdating:
		fallthrough
	case nasMessage.RegistrationType5GSPeriodicRegistrationUpdating:
		return gmm_handler.HandleMobilityAndPeriodicRegistrationUpdating(amfUe, models.AccessType__3_GPP_ACCESS, procedureCode)
	}
	return nil
}

func Authentication_3gpp(sm *fsm.FSM, event fsm.Event, args fsm.Args) error {
	switch event {
	case fsm.EVENT_ENTRY:
	case gmm_event.EVENT_GMM_MESSAGE:
		amfUe := args[gmm_event.AMF_UE].(*amf_context.AmfUe)
		gmmMessage := args[gmm_event.GMM_MESSAGE].(*nas.GmmMessage)
		switch gmmMessage.GetMessageType() {
		case nas.MsgTypeAuthenticationResponse:
			return gmm_handler.HandleAuthenticationResponse(amfUe, models.AccessType__3_GPP_ACCESS, gmmMessage.AuthenticationResponse)
		case nas.MsgTypeAuthenticationFailure:
			return gmm_handler.HandleAuthenticationFailure(amfUe, models.AccessType__3_GPP_ACCESS, gmmMessage.AuthenticationFailure)
		case nas.MsgTypeStatus5GMM:
			return gmm_handler.HandleStatus5GMM(amfUe, models.AccessType__3_GPP_ACCESS, gmmMessage.Status5GMM)
		}
	default:
		GmmLog.Errorf("Unknown Event[%s]\n", event)
	}
	return nil
}

func SecurityMode_3gpp(sm *fsm.FSM, event fsm.Event, args fsm.Args) error {
	switch event {
	case fsm.EVENT_ENTRY:
	case gmm_event.EVENT_GMM_MESSAGE:
		amfUe := args[gmm_event.AMF_UE].(*amf_context.AmfUe)
		procedureCode := args[gmm_event.PROCEDURE_CODE].(int64)
		gmmMessage := args[gmm_event.GMM_MESSAGE].(*nas.GmmMessage)
		switch gmmMessage.GetMessageType() {
		case nas.MsgTypeSecurityModeComplete:
			return gmm_handler.HandleSecurityModeComplete(amfUe, models.AccessType__3_GPP_ACCESS, procedureCode, gmmMessage.SecurityModeComplete)
		case nas.MsgTypeSecurityModeReject:
			return gmm_handler.HandleSecurityModeReject(amfUe, models.AccessType__3_GPP_ACCESS, gmmMessage.SecurityModeReject)
		case nas.MsgTypeStatus5GMM:
			return gmm_handler.HandleStatus5GMM(amfUe, models.AccessType__3_GPP_ACCESS, gmmMessage.Status5GMM)
		}
	default:
		GmmLog.Errorf("Unknown Event[%s]\n", event)
	}
	return nil
}

func InitialContextSetup_3gpp(sm *fsm.FSM, event fsm.Event, args fsm.Args) error {
	switch event {
	case fsm.EVENT_ENTRY:
	case gmm_event.EVENT_GMM_MESSAGE:
		amfUe := args[gmm_event.AMF_UE].(*amf_context.AmfUe)
		gmmMessage := args[gmm_event.GMM_MESSAGE].(*nas.GmmMessage)
		switch gmmMessage.GetMessageType() {
		case nas.MsgTypeRegistrationComplete:
			return gmm_handler.HandleRegistrationComplete(amfUe, models.AccessType__3_GPP_ACCESS, gmmMessage.RegistrationComplete)
		case nas.MsgTypeStatus5GMM:
			return gmm_handler.HandleStatus5GMM(amfUe, models.AccessType__3_GPP_ACCESS, gmmMessage.Status5GMM)
		}
	default:
		GmmLog.Errorf("Unknown Event[%s]\n", event)
	}
	return nil
}

func DeRegistered_non_3gpp(sm *fsm.FSM, event fsm.Event, args fsm.Args) error {
	return register_event_non_3gpp(sm, event, args)
}
func Registered_non_3gpp(sm *fsm.FSM, event fsm.Event, args fsm.Args) error {
	return register_event_non_3gpp(sm, event, args)
}

func register_event_non_3gpp(sm *fsm.FSM, event fsm.Event, args fsm.Args) error {
	var amfUe *amf_context.AmfUe
	var procedureCode int64
	switch event {
	case fsm.EVENT_ENTRY:
		return nil
	case gmm_event.EVENT_GMM_MESSAGE:
		amfUe = args[gmm_event.AMF_UE].(*amf_context.AmfUe)
		gmmMessage := args[gmm_event.GMM_MESSAGE].(*nas.GmmMessage)
		procedureCode = args[gmm_event.PROCEDURE_CODE].(int64)
		switch gmmMessage.GetMessageType() {
		case nas.MsgTypeULNASTransport:
			return gmm_handler.HandleULNASTransport(amfUe, models.AccessType_NON_3_GPP_ACCESS, gmmMessage.ULNASTransport)
		case nas.MsgTypeRegistrationRequest:
			if err := gmm_handler.HandleRegistrationRequest(amfUe, models.AccessType_NON_3_GPP_ACCESS, gmmMessage.RegistrationRequest); err != nil {
				return nil
			}
		case nas.MsgTypeIdentityResponse:
			if err := gmm_handler.HandleIdentityResponse(amfUe, gmmMessage.IdentityResponse); err != nil {
				return err
			}
		case nas.MsgTypeNotificationResponse:
			if err := gmm_handler.HandleNotificationResponse(amfUe, gmmMessage.NotificationResponse); err != nil {
				return err
			}
		case nas.MsgTypeConfigurationUpdateComplete:
			if err := gmm_handler.HandleConfigurationUpdateComplete(amfUe, gmmMessage.ConfigurationUpdateComplete); err != nil {
				return err
			}
		case nas.MsgTypeServiceRequest:
			if err := gmm_handler.HandleServiceRequest(amfUe, models.AccessType_NON_3_GPP_ACCESS, procedureCode, gmmMessage.ServiceRequest); err != nil {
				return err
			}
		case nas.MsgTypeDeregistrationRequestUEOriginatingDeregistration:
			return gmm_handler.HandleDeregistrationRequest(amfUe, models.AccessType_NON_3_GPP_ACCESS, gmmMessage.DeregistrationRequestUEOriginatingDeregistration)
		case nas.MsgTypeDeregistrationAcceptUETerminatedDeregistration:
			return gmm_handler.HandleDeregistrationAccept(amfUe, models.AccessType_NON_3_GPP_ACCESS, gmmMessage.DeregistrationAcceptUETerminatedDeregistration)
		case nas.MsgTypeStatus5GMM:
			if err := gmm_handler.HandleStatus5GMM(amfUe, models.AccessType_NON_3_GPP_ACCESS, gmmMessage.Status5GMM); err != nil {
				return err
			}
		default:
			GmmLog.Errorf("Unknown GmmMessage[%d]\n", gmmMessage.GetMessageType())
		}
	default:
		GmmLog.Errorf("Unknown Event[%s]\n", event)
	}

	switch amfUe.RegistrationType5GS {
	case nasMessage.RegistrationType5GSInitialRegistration:
		return gmm_handler.HandleInitialRegistration(amfUe, models.AccessType_NON_3_GPP_ACCESS)
	case nasMessage.RegistrationType5GSMobilityRegistrationUpdating:
		fallthrough
	case nasMessage.RegistrationType5GSPeriodicRegistrationUpdating:
		return gmm_handler.HandleMobilityAndPeriodicRegistrationUpdating(amfUe, models.AccessType_NON_3_GPP_ACCESS, procedureCode)
	}
	return nil
}

func Authentication_non_3gpp(sm *fsm.FSM, event fsm.Event, args fsm.Args) error {
	switch event {
	case fsm.EVENT_ENTRY:
	case gmm_event.EVENT_GMM_MESSAGE:
		amfUe := args[gmm_event.AMF_UE].(*amf_context.AmfUe)
		gmmMessage := args[gmm_event.GMM_MESSAGE].(*nas.GmmMessage)
		switch gmmMessage.GetMessageType() {
		case nas.MsgTypeAuthenticationResponse:
			return gmm_handler.HandleAuthenticationResponse(amfUe, models.AccessType_NON_3_GPP_ACCESS, gmmMessage.AuthenticationResponse)
		case nas.MsgTypeAuthenticationFailure:
			return gmm_handler.HandleAuthenticationFailure(amfUe, models.AccessType_NON_3_GPP_ACCESS, gmmMessage.AuthenticationFailure)
		case nas.MsgTypeStatus5GMM:
			return gmm_handler.HandleStatus5GMM(amfUe, models.AccessType_NON_3_GPP_ACCESS, gmmMessage.Status5GMM)
		}
	default:
		GmmLog.Errorf("Unknown Event[%s]\n", event)
	}
	return nil
}

func SecurityMode_non_3gpp(sm *fsm.FSM, event fsm.Event, args fsm.Args) error {
	switch event {
	case fsm.EVENT_ENTRY:
	case gmm_event.EVENT_GMM_MESSAGE:
		amfUe := args[gmm_event.AMF_UE].(*amf_context.AmfUe)
		procedureCode := args[gmm_event.PROCEDURE_CODE].(int64)
		gmmMessage := args[gmm_event.GMM_MESSAGE].(*nas.GmmMessage)
		switch gmmMessage.GetMessageType() {
		case nas.MsgTypeSecurityModeComplete:
			return gmm_handler.HandleSecurityModeComplete(amfUe, models.AccessType_NON_3_GPP_ACCESS, procedureCode, gmmMessage.SecurityModeComplete)
		case nas.MsgTypeSecurityModeReject:
			return gmm_handler.HandleSecurityModeReject(amfUe, models.AccessType_NON_3_GPP_ACCESS, gmmMessage.SecurityModeReject)
		case nas.MsgTypeStatus5GMM:
			return gmm_handler.HandleStatus5GMM(amfUe, models.AccessType_NON_3_GPP_ACCESS, gmmMessage.Status5GMM)
		}
	default:
		GmmLog.Errorf("Unknown Event[%s]\n", event)
	}
	return nil
}

func InitialContextSetup_non_3gpp(sm *fsm.FSM, event fsm.Event, args fsm.Args) error {
	switch event {
	case fsm.EVENT_ENTRY:
	case gmm_event.EVENT_GMM_MESSAGE:
		amfUe := args[gmm_event.AMF_UE].(*amf_context.AmfUe)
		gmmMessage := args[gmm_event.GMM_MESSAGE].(*nas.GmmMessage)
		switch gmmMessage.GetMessageType() {
		case nas.MsgTypeRegistrationComplete:
			return gmm_handler.HandleRegistrationComplete(amfUe, models.AccessType_NON_3_GPP_ACCESS, gmmMessage.RegistrationComplete)
		case nas.MsgTypeStatus5GMM:
			return gmm_handler.HandleStatus5GMM(amfUe, models.AccessType_NON_3_GPP_ACCESS, gmmMessage.Status5GMM)
		}
	default:
		GmmLog.Errorf("Unknown Event[%s]\n", event)
	}
	return nil
}

func Exception(sm *fsm.FSM, event fsm.Event, args fsm.Args) error {
	switch event {
	case fsm.EVENT_ENTRY:
	default:
		GmmLog.Errorf("Unknown Event[%s]\n", event)
	}
	return nil
}
