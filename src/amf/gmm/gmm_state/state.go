package gmm_state

import (
	"free5gc/lib/fsm"
)

const (
	DE_REGISTERED         fsm.State = "De-Registered"
	REGISTERED            fsm.State = "Registered"
	AUTHENTICATION        fsm.State = "Authentication"
	SECURITY_MODE         fsm.State = "Security Mode"
	INITIAL_CONTEXT_SETUP fsm.State = "Initial Context Setup"
	EXCEPTION             fsm.State = "Exception"
)
