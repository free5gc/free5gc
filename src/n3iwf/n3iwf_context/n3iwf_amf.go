package n3iwf_context

import (
	"free5gc/lib/ngap/ngapType"
)

type N3IWFAMF struct {
	AMFName             ngapType.AMFName
	ServedGUAMIList     ngapType.ServedGUAMIList
	RelativeAMFCapacity ngapType.RelativeAMFCapacity
	PLMNSupportList     ngapType.PLMNSupportList
}
