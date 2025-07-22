package context

import (
	"math"

	"github.com/free5gc/util/idgenerator"
)

func NewTestContext(n3iwf n3iwf) (*N3IWFContext, error) {
	n := &N3IWFContext{
		n3iwf:                n3iwf,
		RANUENGAPIDGenerator: idgenerator.NewGenerator(0, math.MaxInt64),
		TEIDGenerator:        idgenerator.NewGenerator(1, math.MaxUint32),
	}
	return n, nil
}
