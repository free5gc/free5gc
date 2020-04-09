package smf_context_test

import (
	"fmt"
	"free5gc/lib/pfcp/pfcpType"
	"free5gc/src/smf/smf_context"
	"testing"
)

var upf *smf_context.UPF
var pdrQueue []*smf_context.PDR
var farQueue []*smf_context.FAR
var barQueue []*smf_context.BAR

func init() {
	NodeID := new(pfcpType.NodeID)
	upf = smf_context.NewUPF(NodeID)
	pdrQueue = make([]*smf_context.PDR, 0)
	farQueue = make([]*smf_context.FAR, 0)
	barQueue = make([]*smf_context.BAR, 0)

	for i := 0; i < 6; i++ {
		pdr, _ := upf.AddPDR()
		far, _ := upf.AddFAR()
		bar, _ := upf.AddBAR()

		pdrQueue = append(pdrQueue, pdr)
		farQueue = append(farQueue, far)
		barQueue = append(barQueue, bar)
	}

}

func TestRemovePDR(t *testing.T) {
	var exist bool

	pdr := pdrQueue[0]
	_ = upf.RemovePDR(pdr)
	exist = upf.CheckPDRIDExist(1)
	assertEqual(exist, false)

	pdr = pdrQueue[3]
	_ = upf.RemovePDR(pdr)
	exist = upf.CheckPDRIDExist(4)
	assertEqual(exist, false)

	pdr = pdrQueue[5]
	_ = upf.RemovePDR(pdr)
	exist = upf.CheckPDRIDExist(6)
	assertEqual(exist, false)

	upf.AddPDR()

	exist = upf.CheckPDRIDExist(1)
	assertEqual(exist, true)
}

func TestRemoveFAR(t *testing.T) {
	var exist bool

	far := farQueue[0]
	_ = upf.RemoveFAR(far)

	exist = upf.CheckFARIDExist(2)
	assertEqual(exist, false)

	far = farQueue[3]
	_ = upf.RemoveFAR(far)

	exist = upf.CheckFARIDExist(8)
	assertEqual(exist, false)
	far = farQueue[5]
	_ = upf.RemoveFAR(far)

	exist = upf.CheckFARIDExist(12)
	assertEqual(exist, false)

	upf.AddFAR()

	exist = upf.CheckFARIDExist(2)
	assertEqual(exist, true)
}

func TestRemoveBAR(t *testing.T) {
	var exist bool

	bar := barQueue[0]
	_ = upf.RemoveBAR(bar)
	exist = upf.CheckBARIDExist(1)
	assertEqual(exist, false)

	bar = barQueue[3]
	_ = upf.RemoveBAR(bar)
	exist = upf.CheckBARIDExist(4)
	assertEqual(exist, false)

	bar = barQueue[5]
	_ = upf.RemoveBAR(bar)
	exist = upf.CheckBARIDExist(6)
	assertEqual(exist, false)

	bar, _ = upf.AddBAR()
	bar, _ = upf.AddBAR()
	exist = upf.CheckBARIDExist(1)
	assertEqual(exist, true)
	exist = upf.CheckBARIDExist(4)
	assertEqual(exist, true)
}

func assertEqual(a, b bool) {
	if a != b {
		panic(fmt.Sprintln("Not Equal: ", a, " ", b))
	}
}
