package TestGenAuthData

import (
	"github.com/free5gc/openapi/models"
)

type milenageTestSet struct {
	K      string
	RAND   string
	SQN    string
	AMF    string
	OP     string
	OPC    string
	F1     string
	F1star string
	F2     string
	F3     string
	F4     string
	F5     string
	F5star string
}

const (
	SUCCESS_CASE                 = "success"
	FAILURE_CASE                 = "failure"
	SUCCESS_SERVING_NETWORK_NAME = "free5gc"
	TESTSET_SERVING_NETWORK_NAME = "WLAN"
)

var TestGenAuthDataTable = make(map[string]*models.AuthenticationInfoRequest)
var MilenageTestSet19 milenageTestSet

func init() {
	TestGenAuthDataTable[SUCCESS_CASE] = &models.AuthenticationInfoRequest{
		ServingNetworkName: TESTSET_SERVING_NETWORK_NAME,
	}

	// TS 35.208 test set 19
	MilenageTestSet19 = milenageTestSet{
		K:      "5122250214c33e723a5dd523fc145fc0",
		RAND:   "81e92b6c0ee0e12ebceba8d92a99dfa5",
		SQN:    "16f3b3f70fc2",
		AMF:    "c3ab",
		OP:     "c9e8763286b5b9ffbdf56e1297d0887b",
		OPC:    "981d464c7c52eb6e5036234984ad0bcf",
		F1:     "2a5c23d15ee351d5",
		F1star: "62dae3853f3af9d2",
		F2:     "28d7b0f2a2ec3de5",
		F3:     "5349fbe098649f948f5d2e973a81c00f",
		F4:     "9744871ad32bf9bbd1dd5ce54e3e2e5a",
		F5:     "ada15aeb7bb8",
		F5star: "d461bc15475d",
	}
}
