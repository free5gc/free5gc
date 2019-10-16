package nssf_message

type Event int

const (
	NSSelectionGet Event = iota + 1
	NSSAIAvailabilityPut
	NSSAIAvailabilityPatch
	NSSAIAvailabilityDelete
	NSSAIAvailabilityPost
	NSSAIAvailabilityUnsubscribe
)
