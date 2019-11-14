package n3iwf_message

type Event int

const (
	EventN1UDPMessage Event = iota
	EventN1TUNMessage
	EventNGAPMessage
	EventGTPMessage
)
