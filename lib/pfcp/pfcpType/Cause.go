//go:binary-only-package

package pfcpType

import (
	"fmt"
)

// Acceptance in a response
const (
	CauseRequestAccepted uint8 = 1
)

// Rejection in a response
const (
	CauseRequestRejected uint8 = iota + 64
	CauseSessionContextNotFound
	CauseMandatoryIeMissing
	CauseConditionalIeMissing
	CauseInvalidLength
	CauseMandatoryIeIncorrect
	CauseInvalidForwardingPolicy
	CauseInvalidFTeidAllocationOption
	CauseNoEstablishedPfcpAssociation
	CauseRuleCreationModificationFailure
	CausePfcpEntityInCongestion
	CauseNoResourcesAvailable
	CauseServiceNotSupported
	CauseSystemFailure
)

type Cause struct {
	CauseValue uint8
}

func (c *Cause) MarshalBinary() (data []byte, err error) {}

func (c *Cause) UnmarshalBinary(data []byte) error {}
