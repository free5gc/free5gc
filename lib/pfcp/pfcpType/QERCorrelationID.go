//go:binary-only-package

package pfcpType

import (
	"encoding/binary"
	"fmt"
)

type QERCorrelationID struct {
	QerCorrelationIdValue uint32
}

func (q *QERCorrelationID) MarshalBinary() (data []byte, err error) {}

func (q *QERCorrelationID) UnmarshalBinary(data []byte) error {}
