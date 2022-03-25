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
	SUCCESS_SERVING_NETWORK_NAME = "5G:mnc216.mcc415.3gppnetwork.org"
	FAILURE_SERVING_NETWORK_NAME = "abc"

	// [EphemeralKey_profileA]
	// privKey25519: f5717b2b096949597a37470a9b9fa74130f758e339efc2eebe49db25b079b1d4
	// pubKey25519: 445c8f4bbcfb488523c71b5bd90598ddb80449f55c1261dcc42224b779b9e54c
	// Additional info for this specific testset:
	// shared key: 8f068802410819ef5a7169fdae57ed593bf15b0f50cc5b95b7d2eadb4b8e6c49

	// [EphemeralKey_profileB]
	// Private key: da62c847f56ca586982f84e924aac2a440c9990a656e334fc3d93c33832ba8f1
	// X of Public key: b59c228fa869ad0c57b5f9ad575ffcd3af0994098afd5d52d25c32beb96d19f1
	// Y of Public key: 5300799a2e29449305ec491bdc5d5c4f009de3916b2ef7d005d673651cf6c0e4
	// Additional info for this specific testset:
	// shared key: cc47048554b97744cc0e5856493952f3b3d0b49b529ffaffac300dea66e69fd6

	TEST_SUCI_NULL_SCHEME = "suci-0-274-012-0001-0-01-00012080f6"
	TEST_SUCI_PROFILE_A   = "suci-0-274-012-0001-1-01-445c8f4bbcfb488523c71b5bd90598ddb80449f55c1261dcc42224b779b9e54cfbb72af9893fea0a0d7fef3739"
	TEST_SUCI_PROFILE_B   = "suci-0-274-012-0001-2-01-02b59c228fa869ad0c57b5f9ad575ffcd3af0994098afd5d52d25c32beb96d19f1291c615cfe9f8d480ae336f9a9"
	TEST_SUPI_ANS         = "imsi-27401200012080f6"
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
