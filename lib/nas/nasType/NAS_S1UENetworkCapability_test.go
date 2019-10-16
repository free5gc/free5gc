//go:binary-only-package

package nasType_test

import (
	"free5gc/lib/nas/nasMessage"
	"free5gc/lib/nas/nasType"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNasTypeNewS1UENetworkCapability(t *testing.T) {}

var nasTypeServiceS1UENetworkCapabilityTable = []NasTypeIeiData{
	{nasMessage.RegistrationRequestS1UENetworkCapabilityType, nasMessage.RegistrationRequestS1UENetworkCapabilityType},
}

func TestNasTypeS1UENetworkCapabilityGetSetIei(t *testing.T) {}

var nasTypeServiceS1UENetworkCapabilityLenTable = []NasTypeLenuint8Data{
	{2, 2},
}

func TestNasTypeS1UENetworkCapabilityGetSetLen(t *testing.T) {}

type nasTypeS1UENetworkCapabilityEEA0 struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeS1UENetworkCapabilityEEA0Table = []nasTypeS1UENetworkCapabilityEEA0{
	{7, 0x01, 0x01},
}

func TestNasTypeS1UENetworkCapabilityGetSetEEA0(t *testing.T) {}

type nasTypeS1UENetworkCapabilityEEA1_128 struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeS1UENetworkCapabilityEEA1_128Table = []nasTypeS1UENetworkCapabilityEEA1_128{
	{7, 0x01, 0x01},
}

func TestNasTypeS1UENetworkCapabilityGetSetEEA1_128(t *testing.T) {}

type nasTypeS1UENetworkCapabilityEEA2_128 struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeS1UENetworkCapabilityEEA2_128Table = []nasTypeS1UENetworkCapabilityEEA2_128{
	{7, 0x01, 0x01},
}

func TestNasTypeS1UENetworkCapabilityGetSetEEA2_128(t *testing.T) {}

type nasTypeS1UENetworkCapabilityEEA3_128 struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeS1UENetworkCapabilityEEA3_128Table = []nasTypeS1UENetworkCapabilityEEA3_128{
	{7, 0x01, 0x01},
}

func TestNasTypeS1UENetworkCapabilityGetSetEEA3_128(t *testing.T) {}

type nasTypeS1UENetworkCapabilityEEA4 struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeS1UENetworkCapabilityEEA4Table = []nasTypeS1UENetworkCapabilityEEA4{
	{7, 0x01, 0x01},
}

func TestNasTypeS1UENetworkCapabilityGetSetEEA4(t *testing.T) {}

type nasTypeS1UENetworkCapabilityEEA5 struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeS1UENetworkCapabilityEEA5Table = []nasTypeS1UENetworkCapabilityEEA5{
	{7, 0x01, 0x01},
}

func TestNasTypeS1UENetworkCapabilityGetSetEEA5(t *testing.T) {}

type nasTypeS1UENetworkCapabilityEEA6 struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeS1UENetworkCapabilityEEA6Table = []nasTypeS1UENetworkCapabilityEEA6{
	{7, 0x01, 0x01},
}

func TestNasTypeS1UENetworkCapabilityGetSetEEA6(t *testing.T) {}

type nasTypeS1UENetworkCapabilityEEA7 struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeS1UENetworkCapabilityEEA7Table = []nasTypeS1UENetworkCapabilityEEA7{
	{7, 0x01, 0x01},
}

func TestNasTypeS1UENetworkCapabilityGetSetEEA7(t *testing.T) {}

type nasTypeS1UENetworkCapabilityEIA0 struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeS1UENetworkCapabilityEIA0Table = []nasTypeS1UENetworkCapabilityEIA0{
	{7, 0x01, 0x01},
}

func TestNasTypeS1UENetworkCapabilityGetSetEIA0(t *testing.T) {}

type nasTypeS1UENetworkCapabilityEIA1_128 struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeS1UENetworkCapabilityEIA1_128Table = []nasTypeS1UENetworkCapabilityEIA1_128{
	{7, 0x01, 0x01},
}

func TestNasTypeS1UENetworkCapabilityGetSetEIA1_128(t *testing.T) {}

type nasTypeS1UENetworkCapabilityEIA2_128 struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeS1UENetworkCapabilityEIA2_128Table = []nasTypeS1UENetworkCapabilityEIA2_128{
	{7, 0x01, 0x01},
}

func TestNasTypeS1UENetworkCapabilityGetSetEIA2_128(t *testing.T) {}

type nasTypeS1UENetworkCapabilityEIA3_128 struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeS1UENetworkCapabilityEIA3_128Table = []nasTypeS1UENetworkCapabilityEIA3_128{
	{7, 0x01, 0x01},
}

func TestNasTypeS1UENetworkCapabilityGetSetEIA3_128(t *testing.T) {}

type nasTypeS1UENetworkCapabilityEIA4 struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeS1UENetworkCapabilityEIA4Table = []nasTypeS1UENetworkCapabilityEIA4{
	{7, 0x01, 0x01},
}

func TestNasTypeS1UENetworkCapabilityGetSetEIA4(t *testing.T) {}

type nasTypeS1UENetworkCapabilityEIA5 struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeS1UENetworkCapabilityEIA5Table = []nasTypeS1UENetworkCapabilityEIA4{
	{7, 0x01, 0x01},
}

func TestNasTypeS1UENetworkCapabilityGetSetEIA5(t *testing.T) {}

type nasTypeS1UENetworkCapabilityEIA6 struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeS1UENetworkCapabilityEIA6Table = []nasTypeS1UENetworkCapabilityEIA6{
	{7, 0x01, 0x01},
}

func TestNasTypeS1UENetworkCapabilityGetSetEIA6(t *testing.T) {}

type nasTypeS1UENetworkCapabilityEIA7 struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeS1UENetworkCapabilityEIA7Table = []nasTypeS1UENetworkCapabilityEIA7{
	{7, 0x01, 0x01},
}

func TestNasTypeS1UENetworkCapabilityGetSetEIA7(t *testing.T) {}

type nasTypeS1UENetworkCapabilityUEA0 struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeS1UENetworkCapabilityUEA0Table = []nasTypeS1UENetworkCapabilityUEA0{
	{7, 0x01, 0x01},
}

func TestNasTypeS1UENetworkCapabilityGetSetUEA0(t *testing.T) {}

type nasTypeS1UENetworkCapabilityUEA1 struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeS1UENetworkCapabilityUEA1Table = []nasTypeS1UENetworkCapabilityUEA1{
	{7, 0x01, 0x01},
}

func TestNasTypeS1UENetworkCapabilityGetSetUEA1(t *testing.T) {}

type nasTypeS1UENetworkCapabilityUEA2 struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeS1UENetworkCapabilityUEA2Table = []nasTypeS1UENetworkCapabilityUEA2{
	{7, 0x01, 0x01},
}

func TestNasTypeS1UENetworkCapabilityGetSetUEA2(t *testing.T) {}

type nasTypeS1UENetworkCapabilityUEA3 struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeS1UENetworkCapabilityUEA3Table = []nasTypeS1UENetworkCapabilityUEA3{
	{7, 0x01, 0x01},
}

func TestNasTypeS1UENetworkCapabilityGetSetUEA3(t *testing.T) {}

type nasTypeS1UENetworkCapabilityUEA4 struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeS1UENetworkCapabilityUEA4Table = []nasTypeS1UENetworkCapabilityUEA4{
	{7, 0x01, 0x01},
}

func TestNasTypeS1UENetworkCapabilityGetSetUEA4(t *testing.T) {}

type nasTypeS1UENetworkCapabilityUEA5 struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeS1UENetworkCapabilityUEA5Table = []nasTypeS1UENetworkCapabilityUEA4{
	{7, 0x01, 0x01},
}

func TestNasTypeS1UENetworkCapabilityGetSetUEA5(t *testing.T) {}

type nasTypeS1UENetworkCapabilityUEA6 struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeS1UENetworkCapabilityUEA6Table = []nasTypeS1UENetworkCapabilityUEA6{
	{7, 0x01, 0x01},
}

func TestNasTypeS1UENetworkCapabilityGetSetUEA6(t *testing.T) {}

type nasTypeS1UENetworkCapabilityUEA7 struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeS1UENetworkCapabilityUEA7Table = []nasTypeS1UENetworkCapabilityUEA7{
	{7, 0x01, 0x01},
}

func TestNasTypeS1UENetworkCapabilityGetSetUEA7(t *testing.T) {}

type nasTypeS1UENetworkCapabilityUCS2 struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeS1UENetworkCapabilityUCS2Table = []nasTypeS1UENetworkCapabilityUCS2{
	{7, 0x01, 0x01},
}

func TestNasTypeS1UENetworkCapabilityGetSetUCS2(t *testing.T) {}

type nasTypeS1UENetworkCapabilityUIA1 struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeS1UENetworkCapabilityUIA1Table = []nasTypeS1UENetworkCapabilityUIA1{
	{7, 0x01, 0x01},
}

func TestNasTypeS1UENetworkCapabilityGetSetUIA1(t *testing.T) {}

type nasTypeS1UENetworkCapabilityUIA2 struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeS1UENetworkCapabilityUIA2Table = []nasTypeS1UENetworkCapabilityUIA2{
	{7, 0x01, 0x01},
}

func TestNasTypeS1UENetworkCapabilityGetSetUIA2(t *testing.T) {}

type nasTypeS1UENetworkCapabilityUIA3 struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeS1UENetworkCapabilityUIA3Table = []nasTypeS1UENetworkCapabilityUIA3{
	{7, 0x01, 0x01},
}

func TestNasTypeS1UENetworkCapabilityGetSetUIA3(t *testing.T) {}

type nasTypeS1UENetworkCapabilityUIA4 struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeS1UENetworkCapabilityUIA4Table = []nasTypeS1UENetworkCapabilityUIA4{
	{7, 0x01, 0x01},
}

func TestNasTypeS1UENetworkCapabilityGetSetUIA4(t *testing.T) {}

type nasTypeS1UENetworkCapabilityUIA5 struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeS1UENetworkCapabilityUIA5Table = []nasTypeS1UENetworkCapabilityUIA4{
	{7, 0x01, 0x01},
}

func TestNasTypeS1UENetworkCapabilityGetSetUIA5(t *testing.T) {}

type nasTypeS1UENetworkCapabilityUIA6 struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeS1UENetworkCapabilityUIA6Table = []nasTypeS1UENetworkCapabilityUIA6{
	{7, 0x01, 0x01},
}

func TestNasTypeS1UENetworkCapabilityGetSetUIA6(t *testing.T) {}

type nasTypeS1UENetworkCapabilityUIA7 struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeS1UENetworkCapabilityUIA7Table = []nasTypeS1UENetworkCapabilityUIA7{
	{7, 0x01, 0x01},
}

func TestNasTypeS1UENetworkCapabilityGetSetUIA7(t *testing.T) {}

type nasTypeS1UENetworkCapabilityProSedd struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeS1UENetworkCapabilityProSeddTable = []nasTypeS1UENetworkCapabilityProSedd{
	{7, 0x01, 0x01},
}

func TestNasTypeS1UENetworkCapabilityGetSetProSedd(t *testing.T) {}

type nasTypeS1UENetworkCapabilityProSe struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeS1UENetworkCapabilityProSeTable = []nasTypeS1UENetworkCapabilityProSe{
	{7, 0x01, 0x01},
}

func TestNasTypeS1UENetworkCapabilityGetSetProSe(t *testing.T) {}

type nasTypeS1UENetworkCapabilityH245ASH struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeS1UENetworkCapabilityH245ASHTable = []nasTypeS1UENetworkCapabilityH245ASH{
	{7, 0x01, 0x01},
}

func TestNasTypeS1UENetworkCapabilityGetSetH245ASH(t *testing.T) {}

type nasTypeS1UENetworkCapabilityACCCSFB struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeS1UENetworkCapabilityACCCSFBTable = []nasTypeS1UENetworkCapabilityACCCSFB{
	{7, 0x01, 0x01},
}

func TestNasTypeS1UENetworkCapabilityGetSetACCCSFB(t *testing.T) {}

type nasTypeS1UENetworkCapabilityLPP struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeS1UENetworkCapabilityLPPTable = []nasTypeS1UENetworkCapabilityACCCSFB{
	{7, 0x01, 0x01},
}

func TestNasTypeS1UENetworkCapabilityGetSetLPP(t *testing.T) {}

type nasTypeS1UENetworkCapabilityLCS struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeS1UENetworkCapabilityLCSTable = []nasTypeS1UENetworkCapabilityLCS{
	{7, 0x01, 0x01},
}

func TestNasTypeS1UENetworkCapabilityGetSetLCS(t *testing.T) {}

type nasTypeS1UENetworkCapabilityxSRVCC struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeS1UENetworkCapabilityxSRVCCTable = []nasTypeS1UENetworkCapabilityxSRVCC{
	{7, 0x01, 0x01},
}

func TestNasTypeS1UENetworkCapabilityGetSetxSRVCC(t *testing.T) {}

type nasTypeS1UENetworkCapabilityNF struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeS1UENetworkCapabilityNFTable = []nasTypeS1UENetworkCapabilityNF{
	{7, 0x01, 0x01},
}

func TestNasTypeS1UENetworkCapabilityGetSetNF(t *testing.T) {}

type nasTypeS1UENetworkCapabilityEPCO struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeS1UENetworkCapabilityEPCOTable = []nasTypeS1UENetworkCapabilityEPCO{
	{7, 0x01, 0x01},
}

func TestNasTypeS1UENetworkCapabilityGetSetEPCO(t *testing.T) {}

type nasTypeS1UENetworkCapabilityHCCPCIOT struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeS1UENetworkCapabilityHCCPCIOTTable = []nasTypeS1UENetworkCapabilityHCCPCIOT{
	{7, 0x01, 0x01},
}

func TestNasTypeS1UENetworkCapabilityGetSetHCCPCIOT(t *testing.T) {}

type nasTypeS1UENetworkCapabilityERwoPDN struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeS1UENetworkCapabilityERwoPDNTable = []nasTypeS1UENetworkCapabilityERwoPDN{
	{7, 0x01, 0x01},
}

func TestNasTypeS1UENetworkCapabilityGetSetERwoPDN(t *testing.T) {}

type nasTypeS1UENetworkCapabilityS1UData struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeS1UENetworkCapabilityS1UDataTable = []nasTypeS1UENetworkCapabilityERwoPDN{
	{7, 0x01, 0x01},
}

func TestNasTypeS1UENetworkCapabilityGetSetS1UData(t *testing.T) {}

type nasTypeS1UENetworkCapabilityUPCIot struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeS1UENetworkCapabilityUPCIotTable = []nasTypeS1UENetworkCapabilityUPCIot{
	{7, 0x01, 0x01},
}

func TestNasTypeS1UENetworkCapabilityGetSetUPCIot(t *testing.T) {}

type nasTypeS1UENetworkCapabilityCPCIot struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeS1UENetworkCapabilityCPCIotTable = []nasTypeS1UENetworkCapabilityCPCIot{
	{7, 0x01, 0x01},
}

func TestNasTypeS1UENetworkCapabilityGetSetCPCIot(t *testing.T) {}

type nasTypeS1UENetworkCapabilityProserelay struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeS1UENetworkCapabilityProserelayTable = []nasTypeS1UENetworkCapabilityProserelay{
	{7, 0x01, 0x01},
}

func TestNasTypeS1UENetworkCapabilityGetSetProserelay(t *testing.T) {}

type nasTypeS1UENetworkCapabilityProSedc struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeS1UENetworkCapabilityProSedcTable = []nasTypeS1UENetworkCapabilityProSedc{
	{7, 0x01, 0x01},
}

func TestNasTypeS1UENetworkCapabilityGetSetProSedc(t *testing.T) {}

type nasTypeS1UENetworkCapabilityBearer15 struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeS1UENetworkCapabilityBearer15Table = []nasTypeS1UENetworkCapabilityBearer15{
	{7, 0x01, 0x01},
}

func TestNasTypeS1UENetworkCapabilityGetSetBearer15(t *testing.T) {}

type nasTypeS1UENetworkCapabilitySGC struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeS1UENetworkCapabilitySGCTable = []nasTypeS1UENetworkCapabilitySGC{
	{7, 0x01, 0x01},
}

func TestNasTypeS1UENetworkCapabilityGetSetSGC(t *testing.T) {}

type nasTypeS1UENetworkCapabilityN1mode struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeS1UENetworkCapabilityN1modeTable = []nasTypeS1UENetworkCapabilitySGC{
	{7, 0x01, 0x01},
}

func TestNasTypeS1UENetworkCapabilityGetSetN1mode(t *testing.T) {}

type nasTypeS1UENetworkCapabilityDCNR struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeS1UENetworkCapabilityDCNRTable = []nasTypeS1UENetworkCapabilityDCNR{
	{7, 0x01, 0x01},
}

func TestNasTypeS1UENetworkCapabilityGetSetDCNR(t *testing.T) {}

type nasTypeS1UENetworkCapabilityCPbackoff struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeS1UENetworkCapabilityCPbackoffTable = []nasTypeS1UENetworkCapabilityCPbackoff{
	{7, 0x01, 0x01},
}

func TestNasTypeS1UENetworkCapabilityGetSetCPbackoff(t *testing.T) {}

type nasTypeS1UENetworkCapabilityRestrictEC struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeS1UENetworkCapabilityRestrictECTable = []nasTypeS1UENetworkCapabilityRestrictEC{
	{7, 0x01, 0x01},
}

func TestNasTypeS1UENetworkCapabilityGetSetRestrictEC(t *testing.T) {}

type nasTypeS1UENetworkCapabilityV2XPC5 struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeS1UENetworkCapabilityV2XPC5Table = []nasTypeS1UENetworkCapabilityV2XPC5{
	{7, 0x01, 0x01},
}

func TestNasTypeS1UENetworkCapabilityGetSetV2XPC5(t *testing.T) {}

type nasTypeS1UENetworkCapabilityMulitpeDRB struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeS1UENetworkCapabilityMulitpeDRBTable = []nasTypeS1UENetworkCapabilityMulitpeDRB{
	{7, 0x01, 0x01},
}

func TestNasTypeS1UENetworkCapabilityGetSetMulitpeDRB(t *testing.T) {}

type nasTypeS1UENetworkCapabilitySpare struct {
	inLen uint8
	in    []uint8
	out   []uint8
}

var nasTypeS1UENetworkCapabilitySpareTable = []nasTypeS1UENetworkCapabilitySpare{
	{9, []uint8{0x01, 0x01}, []uint8{0x01, 0x01}},
}

func TestNasTypeS1UENetworkCapabilityGetSetSpare(t *testing.T) {}

type testS1UENetworkCapabilityDataTemplate struct {
	in  nasType.S1UENetworkCapability
	out nasType.S1UENetworkCapability
}

var S1UENetworkCapabilityTestData = []nasType.S1UENetworkCapability{
	{nasMessage.RegistrationRequestS1UENetworkCapabilityType, 9, []uint8{}},
}

var S1UENetworkCapabilityExpectedData = []nasType.S1UENetworkCapability{
	{nasMessage.RegistrationRequestS1UENetworkCapabilityType, 9, []uint8{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x01, 0x01}},
}

var S1UENetworkCapabilityTable = []testS1UENetworkCapabilityDataTemplate{
	{S1UENetworkCapabilityTestData[0], S1UENetworkCapabilityExpectedData[0]},
}

func TestNasTypeS1UENetworkCapability(t *testing.T) {}
