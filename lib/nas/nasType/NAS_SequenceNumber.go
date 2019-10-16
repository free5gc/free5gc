//go:binary-only-package

package nasType

// SequenceNumber 9.10
// SQN Row, sBit, len = [0, 0], 8 , 8
type SequenceNumber struct {
	Octet uint8
}

func NewSequenceNumber() (sequenceNumber *SequenceNumber) {}

// SequenceNumber 9.10
// SQN Row, sBit, len = [0, 0], 8 , 8
func (a *SequenceNumber) GetSQN() (sQN uint8) {}

// SequenceNumber 9.10
// SQN Row, sBit, len = [0, 0], 8 , 8
func (a *SequenceNumber) SetSQN(sQN uint8) {}
