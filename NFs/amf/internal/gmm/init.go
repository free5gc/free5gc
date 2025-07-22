package gmm

import (
	"github.com/free5gc/amf/internal/context"
	"github.com/free5gc/amf/internal/logger"
	"github.com/free5gc/util/fsm"
)

const (
	GmmMessageEvent           fsm.EventType = "Gmm Message"
	StartAuthEvent            fsm.EventType = "Start Authentication"
	AuthSuccessEvent          fsm.EventType = "Authentication Success"
	AuthErrorEvent            fsm.EventType = "Authentication Error"
	AuthRestartEvent          fsm.EventType = "Authentication Restart"
	AuthFailEvent             fsm.EventType = "Authentication Fail"
	SecurityModeSuccessEvent  fsm.EventType = "SecurityMode Success"
	SecurityModeFailEvent     fsm.EventType = "SecurityMode Fail"
	ContextSetupSuccessEvent  fsm.EventType = "ContextSetup Success"
	ContextSetupFailEvent     fsm.EventType = "ContextSetup Fail"
	InitDeregistrationEvent   fsm.EventType = "Initialize Deregistration"
	DeregistrationAcceptEvent fsm.EventType = "Deregistration Accept"
)

const (
	ArgAmfUe               string = "AMF Ue"
	ArgNASMessage          string = "NAS Message"
	ArgProcedureCode       string = "Procedure Code"
	ArgAccessType          string = "Access Type"
	ArgEAPSuccess          string = "EAP Success"
	ArgEAPMessage          string = "EAP Message"
	Arg3GPPDeregistered    string = "3GPP Deregistered"
	ArgNon3GPPDeregistered string = "Non3GPP Deregistered"
)

var transitions = fsm.Transitions{
	{Event: GmmMessageEvent, From: context.Deregistered, To: context.Deregistered},
	{Event: GmmMessageEvent, From: context.Authentication, To: context.Authentication},
	{Event: GmmMessageEvent, From: context.SecurityMode, To: context.SecurityMode},
	{Event: GmmMessageEvent, From: context.ContextSetup, To: context.ContextSetup},
	{Event: GmmMessageEvent, From: context.Registered, To: context.Registered},
	{Event: StartAuthEvent, From: context.Deregistered, To: context.Authentication},
	{Event: StartAuthEvent, From: context.Registered, To: context.Authentication},
	{Event: AuthRestartEvent, From: context.Authentication, To: context.Authentication},
	{Event: AuthSuccessEvent, From: context.Authentication, To: context.SecurityMode},
	{Event: AuthFailEvent, From: context.Authentication, To: context.Deregistered},
	{Event: AuthErrorEvent, From: context.Authentication, To: context.Deregistered},
	{Event: SecurityModeSuccessEvent, From: context.SecurityMode, To: context.ContextSetup},
	{Event: SecurityModeFailEvent, From: context.SecurityMode, To: context.Deregistered},
	{Event: ContextSetupSuccessEvent, From: context.ContextSetup, To: context.Registered},
	{Event: ContextSetupFailEvent, From: context.ContextSetup, To: context.Deregistered},
	{Event: InitDeregistrationEvent, From: context.Registered, To: context.DeregistrationInitiated},
	{Event: DeregistrationAcceptEvent, From: context.DeregistrationInitiated, To: context.Deregistered},
}

var callbacks = fsm.Callbacks{
	context.Deregistered:            DeRegistered,
	context.Authentication:          Authentication,
	context.SecurityMode:            SecurityMode,
	context.ContextSetup:            ContextSetup,
	context.Registered:              Registered,
	context.DeregistrationInitiated: DeregisteredInitiated,
}

var GmmFSM *fsm.FSM

func init() {
	if f, err := fsm.NewFSM(transitions, callbacks); err != nil {
		logger.GmmLog.Errorf("Initialize Gmm FSM Error: %+v", err)
	} else {
		GmmFSM = f
	}
}
