package cdrType

// Need to import "gofree5gc/lib/aper" if it uses "aper"

const (
	TriggerPresentNothing int = iota /* No components present */
	TriggerPresentSMFTrigger
)

type Trigger struct {
	Present    int         /* Choice Type */
	SMFTrigger *SMFTrigger `ber:"tagNum:0"`
}
