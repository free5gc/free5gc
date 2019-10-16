package smf_message

type Event int

const (
	PFCPMessage Event = iota
	PDUSessionSMContextCreate
	PDUSessionSMContextUpdate
	PDUSessionSMContextRelease
)
