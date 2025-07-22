package gmm

import (
	"fmt"
	"testing"

	"github.com/free5gc/util/fsm"
)

func TestGmmFSM(t *testing.T) {
	if err := fsm.ExportDot(GmmFSM, "gmm"); err != nil {
		fmt.Printf("fsm export data return error: %+v", err)
	}
}
