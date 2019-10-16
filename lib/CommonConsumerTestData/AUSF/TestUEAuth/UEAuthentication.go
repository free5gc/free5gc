package TestUEAuth

type ResStarConfirmData struct {
	// these two should be the same
	XresStar string // stored in AUSF
	ResStar  string // calculated by UE
}
type ResConfirmData struct {
	// these two should be the same
	Xres  string // stored in AUSF
	Res   string // calculated by UE
	K_aut string
}

const (
	SUCCESS_CASE                 = "success"
	FAILURE_CASE                 = "failure"
	SUPI                         = "1111222233334444"
	SUCCESS_SERVING_NETWORK_NAME = "5G:mnc216.mcc415.3gppnetwork.org"
	FAILURE_SERVING_NETWORK_NAME = "abc"
)

var TestUe5gAuthTable = make(map[string]*ResStarConfirmData)
var TestUeEapAuthTable = make(map[string]*ResConfirmData)

func init() {
	// for 5G AKA
	TestUe5gAuthTable[SUCCESS_CASE] = &ResStarConfirmData{
		XresStar: "02f8bfc4c22a2e8da31a0da1ae1d4bd6",
		ResStar:  "02f8bfc4c22a2e8da31a0da1ae1d4bd6",
	}

	TestUe5gAuthTable[FAILURE_CASE] = &ResStarConfirmData{
		XresStar: "02f8bfc4c22a2e8da31a0da1ae1d4bd6",
		ResStar:  "dddddddddddddddddddddddddddddddd",
	}

	// for EAP-AKA', K_aut for SUPI_1 only
	TestUeEapAuthTable[SUCCESS_CASE] = &ResConfirmData{
		Xres:  "28d7b0f2a2ec3de5",
		Res:   "28d7b0f2a2ec3de5",
		K_aut: "8ec31ab631cd48b2c2698e8f6279a62bf0034c66d3efa173231f4e4b8aa015a0",
	}

	TestUeEapAuthTable[FAILURE_CASE] = &ResConfirmData{
		Xres:  "28d7b0f2a2ec3de5",
		Res:   "dddddddddddddddd",
		K_aut: "8ec31ab631cd48b2c2698e8f6279a62bf0034c66d3efa173231f4e4b8aa015a0",
	}
}
