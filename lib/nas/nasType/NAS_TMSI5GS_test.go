//go:binary-only-package

package nasType_test

import (
	"free5gc/lib/nas/nasType"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNasTypeNewTMSI5GS(t *testing.T) {}

var nasTypeServiceRequestTMSI5GSTable = []NasTypeIeiData{
	{0x01, 0x01},
}

func TestNasTypeTMSI5GSGetSetIei(t *testing.T) {}

type nasTypeTMSI5GSLen struct {
	in  uint16
	out uint16
}

var nasTypeServiceRequestTMSI5GSLenTable = []nasTypeTMSI5GSLen{
	{2, 2},
}

func TestNasTypeTMSI5GSGetSetLen(t *testing.T) {}

type nasTypeTMSI5GSSpare struct {
	inLen uint16
	in    uint8
	out   uint8
}

var nasTypeTMSI5GSSpareTable = []nasTypeTMSI5GSSpare{
	{2, 0x01, 0x01},
}

func TestNasTypeTMSI5GSGetSetSpare(t *testing.T) {}

type nasTypeTMSI5GSTypeOfIdentity struct {
	inLen uint16
	in    uint8
	out   uint8
}

var nasTypeTMSI5GSTypeOfIdentityTable = []nasTypeTMSI5GSTypeOfIdentity{
	{3, 0x01, 0x01},
}

func TestNasTypeTMSI5GSGetSetTypeOfIdentity(t *testing.T) {}

type nasTypeTMSI5GSAMFSetID struct {
	inLen uint16
	in    uint16
	out   uint16
}

var nasTypeTMSI5GSAMFSetIDTable = []nasTypeTMSI5GSAMFSetID{
	{2, 0x01, 0x01},
}

func TestNasTypeTMSI5GSGetSetAMFSetID(t *testing.T) {}

type nasTypeTMSI5GSAMFPointer struct {
	inLen uint16
	in    uint8
	out   uint8
}

var nasTypeTMSI5GSAMFPointerTable = []nasTypeTMSI5GSAMFPointer{
	{3, 0x01, 0x01},
}

func TestNasTypeTMSI5GSGetSetAMFPointer(t *testing.T) {}

type nasTypeTMSI5GSTMSI5G struct {
	inLen uint16
	in    [4]uint8
	out   [4]uint8
}

var nasTypeTMSI5GSTMSI5GTable = []nasTypeTMSI5GSTMSI5G{
	{4, [4]uint8{0x01, 0x01, 0x01, 0x01}, [4]uint8{0x01, 0x01, 0x01, 0x01}},
}

func TestNasTypeTMSI5GSGetSetTMSI5G(t *testing.T) {}

type testTMSI5GTypeOfIdentityataTemplate struct {
	in  nasType.TMSI5GS
	out nasType.TMSI5GS
}

var TMSI5GSTestData = []nasType.TMSI5GS{
	{0x01, 7, [7]uint8{}},
}

var TMSI5GSExpectedData = []nasType.TMSI5GS{
	{0x01, 7, [7]uint8{0x09, 0x00, 0x41, 0x01, 0x01, 0x01, 0x01}},
}

var TMSI5GSTable = []testTMSI5GTypeOfIdentityataTemplate{
	{TMSI5GSTestData[0], TMSI5GSExpectedData[0]},
}

func TestNasTypeTMSI5GS(t *testing.T) {}
