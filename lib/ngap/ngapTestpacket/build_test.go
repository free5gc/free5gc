//go:binary-only-package

package ngapTestpacket_test

import (
	"bytes"
	"encoding/hex"
	"free5gc/lib/aper"
	"free5gc/lib/nas"
	"free5gc/lib/nas/nasMessage"
	"free5gc/lib/nas/nasTestpacket"
	"free5gc/lib/nas/nasType"
	"free5gc/lib/ngap"
	"free5gc/lib/ngap/logger"
	"free5gc/lib/ngap/ngapTestpacket"
	"free5gc/lib/ngap/ngapType"
	"reflect"
	"testing"
)

type testEncodeData struct {
	out []byte
	in  ngapType.NGAPPDU
}

type testDecodeData struct {
	in  []byte
	out ngapType.NGAPPDU
}

var ngapTestEncodeData = []testEncodeData{}
var ngapTestDecodeData = []testDecodeData{}

var hexString = []string{
	"00150035000004001B00080002F83910454647005240090300667265653547430066001000000000010002F839000010080102030015400140",
}
var pduList = []ngapType.NGAPPDU{
	ngapTestpacket.BuildNGSetupRequest(),
}

func init() {}
func TestNgapEncode(t *testing.T) {}

func TestNgapDecode(t *testing.T) {}

func TestBuildNGSetupRequest(t *testing.T) {}

func TestBuildInitialUEMessage(t *testing.T) {}

func TestBuildErrorIndication(t *testing.T) {}

func TestBuildUEContextReleaseRequest(t *testing.T) {}

func TestBuildUEContextReleaseComplete(t *testing.T) {}

func TestBuildUEContextModificationResponse(t *testing.T) {}

func TestBuildNGReset(t *testing.T) {}

func TestBuildNGResetAcknowledge(t *testing.T) {}

func TestBuildUplinkNasTransport(t *testing.T) {}

func TestBuildInitialContextSetupResponse(t *testing.T) {}

func TestBuildInitialContextSetupFailure(t *testing.T) {}

func TestBuildPathSwitchRequest(t *testing.T) {}

func TestBuildHandoverRequestAcknowledge(t *testing.T) {}

func TestBuildHandoverFailure(t *testing.T) {}

func TestBuildPDUSessionResourceReleaseResponse(t *testing.T) {}

func TestBuildAMFConfigurationUpdateFailure(t *testing.T) {}

func TestBuildUERadioCapabilityCheckResponse(t *testing.T) {}

func TestBuildHandoverCancel(t *testing.T) {}

func TestBuildLocationReportingFailureIndication(t *testing.T) {}

func TestBuildPDUSessionResourceSetupResponse(t *testing.T) {}

func TestBuildPDUSessionResourceModifyResponse(t *testing.T) {}

func TestBuildPDUSessionResourceNotify(t *testing.T) {}

func TestBuildPDUSessionResourceModifyIndication(t *testing.T) {}

func TestBuildUEContextModificationFailure(t *testing.T) {}

func TestBuildRRCInactiveTransitionReport(t *testing.T) {}

func TestBuildHandoverNotify(t *testing.T) {}

func TestBuildUplinkRanStatusTransfer(t *testing.T) {}

func TestBuildNasNonDeliveryIndication(t *testing.T) {}

func TestBuildRanConfigurationUpdate(t *testing.T) {}

func TestBuildAMFStatusIndication(t *testing.T) {}

func TestBuildUplinkRanConfigurationTransfer(t *testing.T) {}

func TestBuildUplinkUEAssociatedNRPPATransport(t *testing.T) {}

func TestBuildUplinkNonUEAssociatedNRPPATransport(t *testing.T) {}

func TestBuildLocationReport(t *testing.T) {}

func TestBuildUETNLABindingReleaseRequest(t *testing.T) {}

func TestBuildUERadioCapabilityInfoIndication(t *testing.T) {}

func TestBuildAMFConfigurationUpdateAcknowledge(t *testing.T) {}

func TestBuildHandoverRequired(t *testing.T) {}

func TestCellTrafficTrace(t *testing.T) {}

func TestBuildPDUSessionResourceReleaseResponseForReleaseTest(t *testing.T) {}
