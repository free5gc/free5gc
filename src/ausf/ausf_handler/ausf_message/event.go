package ausf_message

type Event int

const (
	EventUeAuthPost Event = iota
	EventAuth5gAkaComfirm
	EventEapAuthComfirm
)
