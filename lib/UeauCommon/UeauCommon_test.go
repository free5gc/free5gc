//go:binary-only-package

package UeauCommon

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"testing"
)

type TestKDF struct {
	NetworkName string
	SQNxorAK    string
	CK          string
	IK          string
	FC          string
	DerivedKey  string
}

const (
	SUCCESS_CASE = "success"
	FAILURE_CASE = "failure"
)

var TestKDFTable = make(map[string]*TestKDF)

func init() {}

func TestGetKDFValue(t *testing.T) {}
