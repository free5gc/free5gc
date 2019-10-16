//go:binary-only-package

package nasType_test

import (
	"free5gc/lib/nas/nasMessage"
	"free5gc/lib/nas/nasType"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNasTypeNewNetworkFeatureSupport5GS(t *testing.T) {}

var nasTypeNetworkFeatureSupport5GSRegistrationAcceptNetworkFeatureSupport5GSTypeTable = []NasTypeIeiData{
	{nasMessage.RegistrationAcceptNetworkFeatureSupport5GSType, nasMessage.RegistrationAcceptNetworkFeatureSupport5GSType},
}

func TestNasTypeNetworkFeatureSupport5GSGetSetIei(t *testing.T) {}

var nasTypeNetworkFeatureSupport5GSLenTable = []NasTypeLenuint8Data{
	{2, 2},
}

func TestNasTypeNetworkFeatureSupport5GSGetSetLen(t *testing.T) {}

type nasTypeNetworkFeatureSupport5GSMPSIData struct {
	in  uint8
	out uint8
}

var nasTypeNetworkFeatureSupport5GSMPSITable = []nasTypeNetworkFeatureSupport5GSMPSIData{
	{0x01, 0x01},
}

func TestNasTypeNetworkFeatureSupport5GSGetSetMPSI(t *testing.T) {}

type nasTypeNetworkFeatureSupport5GSIWKN26Data struct {
	in  uint8
	out uint8
}

var nasTypeNetworkFeatureSupport5GSIWKN26Table = []nasTypeNetworkFeatureSupport5GSIWKN26Data{
	{0x01, 0x01},
}

func TestNasTypeNetworkFeatureSupport5GSGetSetIWKN26(t *testing.T) {}

type nasTypeNetworkFeatureSupport5GSEMFData struct {
	in  uint8
	out uint8
}

var nasTypeNetworkFeatureSupport5GSEMFTable = []nasTypeNetworkFeatureSupport5GSEMFData{
	{0x03, 0x03},
}

func TestNasTypeNetworkFeatureSupport5GSGetSetEMF(t *testing.T) {}

type nasTypeNetworkFeatureSupport5GSEMCData struct {
	in  uint8
	out uint8
}

var nasTypeNetworkFeatureSupport5GSEMCTable = []nasTypeNetworkFeatureSupport5GSEMCData{
	{0x03, 0x03},
}

func TestNasTypeNetworkFeatureSupport5GSGetSetEMC(t *testing.T) {}

type nasTypeNetworkFeatureSupport5GSIMSVoPSN3GPPData struct {
	in  uint8
	out uint8
}

var nasTypeNetworkFeatureSupport5GSIMSVoPSN3GPPTable = []nasTypeNetworkFeatureSupport5GSIMSVoPSN3GPPData{
	{0x01, 0x01},
}

func TestNasTypeNetworkFeatureSupport5GSGetSetIMSVoPSN3GPP(t *testing.T) {}

type nasTypeNetworkFeatureSupport5GSIMSVoPS3GPPData struct {
	in  uint8
	out uint8
}

var nasTypeNetworkFeatureSupport5GSIMSVoPS3GPPTable = []nasTypeNetworkFeatureSupport5GSIMSVoPS3GPPData{
	{0x01, 0x01},
}

func TestNasTypeNetworkFeatureSupport5GSGetSetIMSVoPS3GPP(t *testing.T) {}

type nasTypeNetworkFeatureSupport5GSMCSIData struct {
	in  uint8
	out uint8
}

var nasTypeNetworkFeatureSupport5GSMCSITable = []nasTypeNetworkFeatureSupport5GSMCSIData{
	{0x01, 0x01},
}

func TestNasTypeNetworkFeatureSupport5GSGetSetMCSI(t *testing.T) {}

type nasTypeNetworkFeatureSupport5GSEMCNData struct {
	in  uint8
	out uint8
}

var nasTypeNetworkFeatureSupport5GSEMCNTable = []nasTypeNetworkFeatureSupport5GSEMCNData{
	{0x01, 0x01},
}

func TestNasTypeNetworkFeatureSupport5GSGetSetEMCN(t *testing.T) {}

type nasTypeNetworkFeatureSupport5GSSpareData struct {
	in  uint8
	out uint8
}

var nasTypeNetworkFeatureSupport5GSSpareTable = []nasTypeNetworkFeatureSupport5GSSpareData{
	{0x00, 0x00},
}

func TestNasTypeNetworkFeatureSupport5GSGetSetSpare(t *testing.T) {}

type testNetworkFeatureSupport5GSDataTemplate struct {
	inIei           uint8
	inLen           uint8
	inMPSI          uint8
	inIWKN26        uint8
	inEMF           uint8
	inEMC           uint8
	inIMSVoPSN3GPP  uint8
	inIMSVoPS3GPP   uint8
	inMCSI          uint8
	inEMCN          uint8
	inSpare         uint8
	outIei          uint8
	outLen          uint8
	outMPSI         uint8
	outIWKN26       uint8
	outEMF          uint8
	outEMC          uint8
	outIMSVoPSN3GPP uint8
	outIMSVoPS3GPP  uint8
	outMCSI         uint8
	outEMCN         uint8
	outSpare        uint8
}

var testNetworkFeatureSupport5GSTestTable = []testNetworkFeatureSupport5GSDataTemplate{
	{nasMessage.RegistrationAcceptNetworkFeatureSupport5GSType, 2, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01,
		nasMessage.RegistrationAcceptNetworkFeatureSupport5GSType, 2, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01},
}

func TestNasTypeNetworkFeatureSupport5GS(t *testing.T) {}
