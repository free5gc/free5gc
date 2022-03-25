package ngapTestpacket_test

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"reflect"
	"testing"

	"test/nasTestpacket"
	"test/ngapTestpacket"

	"github.com/free5gc/aper"
	"github.com/free5gc/nas"
	"github.com/free5gc/nas/nasMessage"
	"github.com/free5gc/nas/nasType"
	"github.com/free5gc/ngap"
	"github.com/free5gc/ngap/ngapType"
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

func init() {
	for i := range hexString {
		rawData, _ := hex.DecodeString(hexString[i])
		encodeData := testEncodeData{
			in:  pduList[i],
			out: rawData,
		}
		decodeData := testDecodeData{
			out: pduList[i],
			in:  rawData,
		}
		ngapTestEncodeData = append(ngapTestEncodeData, encodeData)
		ngapTestDecodeData = append(ngapTestDecodeData, decodeData)
	}

}
func TestNgapEncode(t *testing.T) {
	fmt.Println("------Encode------")
	for i, test := range ngapTestEncodeData {
		fmt.Println("[TEST ", i+1, "]")
		// pdu, err := Decoder(test.in)
		// ngapTestTrace(2, fmt.Sprintf("	in : %0x", test.in))
		// ngapTestTrace(2, fmt.Sprintf("	out : \n%s", PrintResult(reflect.ValueOf(pdu).Elem(), 0)))
		output, err := ngap.Encoder(test.in)
		// fmt.Println("	out      : ", output)
		// fmt.Println("	expected : ", test.out)

		if err != nil {
			t.Error(err.Error())

		} else if reflect.DeepEqual(test.out, output) {
			fmt.Println("[PASS]")
			continue
		}
		fmt.Println("[FAIL]")
		t.Errorf("TEST %d is FAILED", i+1)
	}
}

func TestNgapDecode(t *testing.T) {
	fmt.Println("------Decode------")
	for i, test := range ngapTestDecodeData {
		fmt.Println("[ TEST ", i+1, "]")
		output, err := ngap.Decoder(test.in)
		// fmt.Println("	out : \n", ngap.PrintResult(reflect.ValueOf(output).Elem(), 0))
		// fmt.Println("	expected : \n", ngap.PrintResult(reflect.ValueOf(test.out), 0))

		if err != nil {
			t.Error(err.Error())
		} else if reflect.DeepEqual(test.out, *output) {
			fmt.Println("[PASS]")
			continue
		}
		fmt.Println("[FAIL]")
		t.Errorf("TEST %d is FAILED", i+1)
	}
}

func TestBuildNGSetupRequest(t *testing.T) {
	pdu := ngapTestpacket.BuildNGSetupRequest()
	encodeData, err := ngap.Encoder(pdu)
	if err != nil {
		t.Error(err.Error())
		t.Error("Encode is failed")
	} else {
		decodeData, err := ngap.Decoder(encodeData)
		if err != nil {
			t.Error(err.Error())
			t.Error("Decode is Failed")
		} else if reflect.DeepEqual(pdu, *decodeData) {
			fmt.Println("[PASS]")
		}
	}
}

func TestBuildInitialUEMessage(t *testing.T) {

	mobileIdentity5GS := nasType.MobileIdentity5GS{
		Len:    12, // suci
		Buffer: []uint8{0x01, 0x02, 0xf8, 0x39, 0xf0, 0x00, 0x00, 0x00, 0x00, 0x00, 0x47, 0x78},
	}
	nasPdu := nasTestpacket.GetRegistrationRequest(1, mobileIdentity5GS, nil, nil, nil, nil, nil)
	pdu := ngapTestpacket.BuildInitialUEMessage(123, nasPdu, "")
	encodeData, err := ngap.Encoder(pdu)
	if err != nil {
		t.Error(err.Error())
		t.Error("Encode is failed")
	} else {
		decodeData, err := ngap.Decoder(encodeData)
		if err != nil {
			t.Error(err.Error())
			t.Error("Decode is Failed")
		} else if reflect.DeepEqual(pdu, *decodeData) {
			fmt.Println("[PASS]")
		}
	}
}

func TestBuildErrorIndication(t *testing.T) {
	pdu := ngapTestpacket.BuildErrorIndication()
	encodeData, err := ngap.Encoder(pdu)
	if err != nil {
		t.Error(err.Error())
		t.Error("Encode is failed")
	} else {
		decodeData, err := ngap.Decoder(encodeData)
		if err != nil {
			t.Error(err.Error())
			t.Error("Decode is failed")
		} else if reflect.DeepEqual(pdu, *decodeData) {
			fmt.Println("[PASS]")
		}
	}
}

func TestBuildUEContextReleaseRequest(t *testing.T) {
	pdu := ngapTestpacket.BuildUEContextReleaseRequest(123, 456, nil)
	encodeData, err := ngap.Encoder(pdu)
	if err != nil {
		t.Error(err.Error())
		t.Error("Encode is failed")
	} else {
		decodeData, err := ngap.Decoder(encodeData)
		if err != nil {
			t.Error(err.Error())
			t.Error("Decode is Failed")
		} else if reflect.DeepEqual(pdu, *decodeData) {
			fmt.Println("[PASS]")
		}
	}
}

func TestBuildUEContextReleaseComplete(t *testing.T) {

	pdu := ngapTestpacket.BuildUEContextReleaseComplete(1, 2, nil)
	encodeData, err := ngap.Encoder(pdu)
	if err != nil {
		t.Error(err.Error())
		t.Error("Encode is Failed")
	} else {
		decodeData, err := ngap.Decoder(encodeData)
		if err != nil {
			t.Error(err.Error())
			t.Error("Decode is Failed")
		} else if reflect.DeepEqual(pdu, *decodeData) {
			fmt.Println("[PASS]")
		}
	}
}

func TestBuildUEContextModificationResponse(t *testing.T) {
	pdu := ngapTestpacket.BuildUEContextModificationResponse(123, 456)
	encodeData, err := ngap.Encoder(pdu)
	if err != nil {
		t.Error(err.Error())
		t.Error("Encode is failed")
	} else {
		decodeData, err := ngap.Decoder(encodeData)
		if err != nil {
			t.Error(err.Error())
			t.Error("Decode is Failed")
		} else if reflect.DeepEqual(pdu, *decodeData) {
			fmt.Println("[PASS]")
		}
	}
}

func TestBuildNGReset(t *testing.T) {
	pdu := ngapTestpacket.BuildNGReset(nil)
	encodeData, err := ngap.Encoder(pdu)
	if err != nil {
		t.Error(err.Error())
		t.Error("Encode is failed")
	} else {
		decodeData, err := ngap.Decoder(encodeData)
		if err != nil {
			t.Error(err.Error())
			t.Error("Decode is Failed")
		} else if reflect.DeepEqual(pdu, *decodeData) {
			fmt.Println("[PASS]")
		}
	}
}

func TestBuildNGResetAcknowledge(t *testing.T) {
	pdu := ngapTestpacket.BuildNGResetAcknowledge()
	encodeData, err := ngap.Encoder(pdu)
	if err != nil {
		t.Error(err.Error())
		t.Error("Encode is failed")
	} else {
		decodeData, err := ngap.Decoder(encodeData)
		if err != nil {
			t.Error(err.Error())
			t.Error("Decode is Failed")
		} else if reflect.DeepEqual(pdu, *decodeData) {
			fmt.Println("[PASS]")
		}
	}
}

func TestBuildUplinkNasTransport(t *testing.T) {

	pduSessionEstablishmentRequest := nasTestpacket.GetPduSessionEstablishmentRequest(1)

	m := nas.NewMessage()
	m.GmmMessage = nas.NewGmmMessage()
	m.GmmHeader.SetMessageType(nas.MsgTypeULNASTransport)

	ulNasTransport := nasMessage.NewULNASTransport(0)
	ulNasTransport.SpareHalfOctetAndSecurityHeaderType.SetSecurityHeaderType(nas.SecurityHeaderTypePlainNas)
	ulNasTransport.SetMessageType(nas.MsgTypeULNASTransport)
	ulNasTransport.ExtendedProtocolDiscriminator.SetExtendedProtocolDiscriminator(nasMessage.Epd5GSMobilityManagementMessage)
	ulNasTransport.PduSessionID2Value = new(nasType.PduSessionID2Value)
	ulNasTransport.PduSessionID2Value.SetIei(nasMessage.ULNASTransportPduSessionID2ValueType)
	ulNasTransport.PduSessionID2Value.SetPduSessionID2Value(1)
	ulNasTransport.RequestType = new(nasType.RequestType)
	ulNasTransport.RequestType.SetIei(nasMessage.ULNASTransportRequestTypeType)
	ulNasTransport.RequestType.SetRequestTypeValue(nasMessage.ULNASTransportRequestTypeInitialRequest)
	ulNasTransport.SpareHalfOctetAndPayloadContainerType.SetPayloadContainerType(nasMessage.PayloadContainerTypeN1SMInfo)
	ulNasTransport.PayloadContainer.SetLen(uint16(len(pduSessionEstablishmentRequest)))
	ulNasTransport.PayloadContainer.SetPayloadContainerContents(pduSessionEstablishmentRequest)

	m.GmmMessage.ULNASTransport = ulNasTransport
	nasPdu := bytes.Buffer{}
	if err := m.GmmMessageEncode(&nasPdu); err != nil {
		t.Error(err.Error())
		t.Error("GmmEncode is failed")
	}

	pdu := ngapTestpacket.BuildUplinkNasTransport(1, 2, nasPdu.Bytes())
	encodeData, err := ngap.Encoder(pdu)
	if err != nil {
		t.Error(err.Error())
		t.Error("Encode is failed")
	} else {
		decodeData, err := ngap.Decoder(encodeData)
		if err != nil {
			t.Error(err.Error())
			t.Error("Decode is Failed")
		} else if reflect.DeepEqual(pdu, *decodeData) {
			fmt.Println("[PASS]")
		}
	}
}

func TestBuildInitialContextSetupResponse(t *testing.T) {
	pdu := ngapTestpacket.BuildInitialContextSetupResponse(123, 456, "10.200.200.1", nil)
	encodeData, err := ngap.Encoder(pdu)
	if err != nil {
		t.Error(err.Error())
		t.Error("Encode is failed")
	} else {
		decodeData, err := ngap.Decoder(encodeData)
		if err != nil {
			t.Error(err.Error())
			t.Error("Decode is Failed")
		} else if reflect.DeepEqual(pdu, *decodeData) {
			fmt.Println("[PASS]")
		}
	}
}

func TestBuildInitialContextSetupFailure(t *testing.T) {
	pdu := ngapTestpacket.BuildInitialContextSetupFailure(123, 456)
	encodeData, err := ngap.Encoder(pdu)
	if err != nil {
		t.Error(err.Error())
		t.Error("Encode is failed")
	} else {
		decodeData, err := ngap.Decoder(encodeData)
		if err != nil {
			t.Error(err.Error())
			t.Error("Decode is Failed")
		} else if reflect.DeepEqual(pdu, *decodeData) {
			fmt.Println("[PASS]")
		}
	}
}

func TestBuildPathSwitchRequest(t *testing.T) {
	pdu := ngapTestpacket.BuildPathSwitchRequest(1, 2)
	encodeData, err := ngap.Encoder(pdu)
	if err != nil {
		t.Error(err.Error())
		t.Error("Encode is failed")
	} else {
		decodeData, err := ngap.Decoder(encodeData)
		if err != nil {
			t.Error(err.Error())
			t.Error("Decode is Failed")
		} else if reflect.DeepEqual(pdu, *decodeData) {
			fmt.Println("[PASS]")
		}
	}
}

func TestBuildHandoverRequestAcknowledge(t *testing.T) {
	pdu := ngapTestpacket.BuildHandoverRequestAcknowledge(1, 2)
	encodeData, err := ngap.Encoder(pdu)
	if err != nil {
		t.Error(err.Error())
		t.Error("Encode is failed")
	} else {
		decodeData, err := ngap.Decoder(encodeData)
		if err != nil {
			t.Error(err.Error())
			t.Error("Decode is failed")
		} else if reflect.DeepEqual(pdu, *decodeData) {
			fmt.Println("[PASS]")
		}
	}
}

func TestBuildHandoverFailure(t *testing.T) {
	pdu := ngapTestpacket.BuildHandoverFailure(1)
	encodeData, err := ngap.Encoder(pdu)
	if err != nil {
		t.Error(err.Error())
		t.Error("Encode is failed")
	} else {
		decodeData, err := ngap.Decoder(encodeData)
		if err != nil {
			t.Error(err.Error())
			t.Error("Decode is failed")
		} else if reflect.DeepEqual(pdu, *decodeData) {
			fmt.Println("[PASS]")
		}
	}
}

func TestBuildPDUSessionResourceReleaseResponse(t *testing.T) {

	pdu := ngapTestpacket.BuildPDUSessionResourceReleaseResponse()
	encodeData, err := ngap.Encoder(pdu)
	if err != nil {
		t.Error(err.Error())
		t.Error("Encode is failed")
	} else {
		decodeData, err := ngap.Decoder(encodeData)
		if err != nil {
			t.Error(err.Error())
			t.Error("Decode is Failed")
		} else if reflect.DeepEqual(pdu, *decodeData) {
			fmt.Println("[PASS]")
		}
	}
}

func TestBuildAMFConfigurationUpdateFailure(t *testing.T) {
	pdu := ngapTestpacket.BuildAMFConfigurationUpdateFailure()
	encodeData, err := ngap.Encoder(pdu)
	if err != nil {
		t.Error(err.Error())
		t.Error("Encode is failed")
	} else {
		decodeData, err := ngap.Decoder(encodeData)
		if err != nil {
			t.Error(err.Error())
			t.Error("Decode is Failed")
		} else if reflect.DeepEqual(pdu, *decodeData) {
			fmt.Println("[PASS]")
		}
	}
}

func TestBuildUERadioCapabilityCheckResponse(t *testing.T) {
	pdu := ngapTestpacket.BuildUERadioCapabilityCheckResponse()
	encodeData, err := ngap.Encoder(pdu)
	if err != nil {
		t.Error(err.Error())
		t.Error("Encode is failed")
	} else {
		decodeData, err := ngap.Decoder(encodeData)
		if err != nil {
			t.Error(err.Error())
			t.Error("Decode is Failed")
		} else if reflect.DeepEqual(pdu, *decodeData) {
			fmt.Println("[PASS]")
		}
	}
}

func TestBuildHandoverCancel(t *testing.T) {
	pdu := ngapTestpacket.BuildHandoverCancel()
	encodeData, err := ngap.Encoder(pdu)
	if err != nil {
		t.Error(err.Error())
		t.Error("Encode is failed")
	} else {
		decodeData, err := ngap.Decoder(encodeData)
		if err != nil {
			t.Error(err.Error())
			t.Error("Decode is Failed")
		} else if reflect.DeepEqual(pdu, *decodeData) {
			fmt.Println("[PASS]")
		}
	}
}

func TestBuildLocationReportingFailureIndication(t *testing.T) {
	pdu := ngapTestpacket.BuildLocationReportingFailureIndication()
	encodeData, err := ngap.Encoder(pdu)
	if err != nil {
		t.Error(err.Error())
		t.Error("Encode is failed")
	} else {
		decodeData, err := ngap.Decoder(encodeData)
		if err != nil {
			t.Error(err.Error())
			t.Error("Decode is Failed")
		} else if reflect.DeepEqual(pdu, *decodeData) {
			fmt.Println("[PASS]")
		}
	}
}

func TestBuildPDUSessionResourceSetupResponse(t *testing.T) {
	pdu := ngapTestpacket.BuildPDUSessionResourceSetupResponse(123, 456, "10.200.200.1")
	encodeData, err := ngap.Encoder(pdu)
	if err != nil {
		t.Error(err.Error())
		t.Error("Encode is failed")
	} else {
		decodeData, err := ngap.Decoder(encodeData)
		if err != nil {
			t.Error(err.Error())
			t.Error("Decode is Failed")
		} else if reflect.DeepEqual(pdu, *decodeData) {
			fmt.Println("[PASS]")
		}
	}
}

func TestBuildPDUSessionResourceModifyResponse(t *testing.T) {
	pdu := ngapTestpacket.BuildPDUSessionResourceModifyResponse(123, 456)
	encodeData, err := ngap.Encoder(pdu)
	if err != nil {
		t.Error(err.Error())
		t.Error("Encode is failed")
	} else {
		decodeData, err := ngap.Decoder(encodeData)
		if err != nil {
			t.Error(err.Error())
			t.Error("Decode is Failed")
		} else if reflect.DeepEqual(pdu, *decodeData) {
			fmt.Println("[PASS]")
		}
	}
}

func TestBuildPDUSessionResourceNotify(t *testing.T) {
	pdu := ngapTestpacket.BuildPDUSessionResourceNotify()
	encodeData, err := ngap.Encoder(pdu)
	if err != nil {
		t.Error(err.Error())
		t.Error("Encode is failed")
	} else {
		decodeData, err := ngap.Decoder(encodeData)
		if err != nil {
			t.Error(err.Error())
			t.Error("Decode is Failed")
		} else if reflect.DeepEqual(pdu, *decodeData) {
			fmt.Println("[PASS]")
		}
	}
}

func TestBuildPDUSessionResourceModifyIndication(t *testing.T) {
	pdu := ngapTestpacket.BuildPDUSessionResourceModifyIndication(123, 456)
	encodeData, err := ngap.Encoder(pdu)
	if err != nil {
		t.Error(err.Error())
		t.Error("Encode is failed")
	} else {
		decodeData, err := ngap.Decoder(encodeData)
		if err != nil {
			t.Error(err.Error())
			t.Error("Decode is Failed")
		} else if reflect.DeepEqual(pdu, *decodeData) {
			fmt.Println("[PASS]")
		}
	}
}

func TestBuildUEContextModificationFailure(t *testing.T) {
	pdu := ngapTestpacket.BuildUEContextModificationFailure(123, 456)
	encodeData, err := ngap.Encoder(pdu)
	if err != nil {
		t.Error(err.Error())
		t.Error("Encode is failed")
	} else {
		decodeData, err := ngap.Decoder(encodeData)
		if err != nil {
			t.Error(err.Error())
			t.Error("Decode is Failed")
		} else if reflect.DeepEqual(pdu, *decodeData) {
			fmt.Println("[PASS]")
		}
	}
}

func TestBuildRRCInactiveTransitionReport(t *testing.T) {
	pdu := ngapTestpacket.BuildRRCInactiveTransitionReport()
	encodeData, err := ngap.Encoder(pdu)
	if err != nil {
		t.Error(err.Error())
	} else {
		decodeData, err := ngap.Decoder(encodeData)
		if err != nil {
			t.Error(err.Error())
		} else if reflect.DeepEqual(pdu, *decodeData) {
			fmt.Println("[PASS]")
		}
	}
}

func TestBuildHandoverNotify(t *testing.T) {
	pdu := ngapTestpacket.BuildHandoverNotify(1, 2)
	encodeData, err := ngap.Encoder(pdu)
	if err != nil {
		t.Error(err.Error())
	} else {
		decodeData, err := ngap.Decoder(encodeData)
		if err != nil {
			t.Error(err.Error())
		} else if reflect.DeepEqual(pdu, *decodeData) {
			fmt.Println("[PASS]")
		}
	}
}

func TestBuildUplinkRanStatusTransfer(t *testing.T) {
	pdu := ngapTestpacket.BuildUplinkRanStatusTransfer(211, 321)
	encodeData, err := ngap.Encoder(pdu)
	if err != nil {
		t.Error(err.Error())
		t.Error("Encode is failed")
	} else {
		decodeData, err := ngap.Decoder(encodeData)
		if err != nil {
			t.Error(err.Error())
			t.Error("Decode is failed")
		} else if reflect.DeepEqual(pdu, *decodeData) {
			fmt.Println("[PASS]")
		}
	}
}

func TestBuildNasNonDeliveryIndication(t *testing.T) {

	pdu := ngapTestpacket.BuildNasNonDeliveryIndication(123, 456, aper.OctetString("\x01\x02\x03"))
	encodeData, err := ngap.Encoder(pdu)
	if err != nil {
		t.Error(err.Error())
		t.Error("Encode is failed")
	} else {
		decodeData, err := ngap.Decoder(encodeData)
		if err != nil {
			t.Error(err.Error())
			t.Error("Decode is Failed")
		} else if reflect.DeepEqual(pdu, *decodeData) {
			fmt.Println("[PASS]")
		}
	}
}

func TestBuildRanConfigurationUpdate(t *testing.T) {

	pdu := ngapTestpacket.BuildRanConfigurationUpdate()
	encodeData, err := ngap.Encoder(pdu)
	if err != nil {
		t.Error(err.Error())
		t.Error("Encode is failed")
	} else {
		decodeData, err := ngap.Decoder(encodeData)
		if err != nil {
			t.Error(err.Error())
			t.Error("Decode is Failed")
		} else if reflect.DeepEqual(pdu, *decodeData) {
			fmt.Println("[PASS]")
		}
	}
}

func TestBuildAMFStatusIndication(t *testing.T) {
}

func TestBuildUplinkRanConfigurationTransfer(t *testing.T) {
	pdu := ngapTestpacket.BuildUplinkRanConfigurationTransfer()
	encodeData, err := ngap.Encoder(pdu)
	if err != nil {
		t.Error(err.Error())
		t.Error("Encode is failed")
	} else {
		decodeData, err := ngap.Decoder(encodeData)
		if err != nil {
			t.Error(err.Error())
			t.Error("Decode is failed")
		} else if reflect.DeepEqual(pdu, *decodeData) {
			fmt.Println("[PASS]")
		}
	}
}

func TestBuildUplinkUEAssociatedNRPPATransport(t *testing.T) {

	pdu := ngapTestpacket.BuildUplinkUEAssociatedNRPPATransport()
	encodeData, err := ngap.Encoder(pdu)
	if err != nil {
		t.Error(err.Error())
	} else {
		decodeData, err := ngap.Decoder(encodeData)
		if err != nil {
			t.Error(err.Error())
		} else if reflect.DeepEqual(pdu, *decodeData) {
			fmt.Println("[PASS]")
		}
	}
}

func TestBuildUplinkNonUEAssociatedNRPPATransport(t *testing.T) {
	pdu := ngapTestpacket.BuildUplinkNonUEAssociatedNRPPATransport()
	encodeData, err := ngap.Encoder(pdu)
	if err != nil {
		t.Error(err.Error())
		t.Error("Encode is failed")
	} else {
		decodeData, err := ngap.Decoder(encodeData)
		if err != nil {
			t.Error(err.Error())
			t.Error("Decode is Failed")
		} else if reflect.DeepEqual(pdu, *decodeData) {
			fmt.Println("[PASS]")
		}
	}
}

func TestBuildLocationReport(t *testing.T) {

	pdu := ngapTestpacket.BuildLocationReport()
	encodeData, err := ngap.Encoder(pdu)
	if err != nil {
		t.Error(err.Error())
		t.Error("Encode is failed")
	} else {
		decodeData, err := ngap.Decoder(encodeData)
		if err != nil {
			t.Error(err.Error())
			t.Error("Decode is Failed")
		} else if reflect.DeepEqual(pdu, *decodeData) {
			fmt.Println("[PASS]")
		}
	}
}

func TestBuildUETNLABindingReleaseRequest(t *testing.T) {
}

func TestBuildUERadioCapabilityInfoIndication(t *testing.T) {
	pdu := ngapTestpacket.BuildUERadioCapabilityInfoIndication()
	encodeData, err := ngap.Encoder(pdu)
	if err != nil {
		t.Error(err.Error())
		t.Error("Encode is failed")
	} else {
		decodeData, err := ngap.Decoder(encodeData)
		if err != nil {
			t.Error(err.Error())
			t.Error("Decode is failed")
		} else if reflect.DeepEqual(pdu, *decodeData) {
			fmt.Println("[PASS]")
		}
	}
}

func TestBuildAMFConfigurationUpdateAcknowledge(t *testing.T) {
	pdu := ngapTestpacket.BuildAMFConfigurationUpdateAcknowledge()
	encodeData, err := ngap.Encoder(pdu)
	if err != nil {
		t.Error(err.Error())
		t.Error("Encode is failed")
	} else {
		decodeData, err := ngap.Decoder(encodeData)
		if err != nil {
			t.Error(err.Error())
			t.Error("Decode is failed")
		} else if reflect.DeepEqual(pdu, *decodeData) {
			fmt.Println("[PASS]")
		}
	}
}

func TestBuildHandoverRequired(t *testing.T) {
	pdu := ngapTestpacket.BuildHandoverRequired(1, 2, []byte{0x00, 0x01, 0x02}, []byte{0x01, 0x20})
	encodeData, err := ngap.Encoder(pdu)
	if err != nil {
		t.Error(err.Error())
		t.Error("Encode is failed")
	} else {
		decodeData, err := ngap.Decoder(encodeData)
		if err != nil {
			t.Error(err.Error())
			t.Error("Decode is failed")
		} else if reflect.DeepEqual(pdu, *decodeData) {
			fmt.Println("[PASS]")
		}
	}
}

func TestCellTrafficTrace(t *testing.T) {
	pdu := ngapTestpacket.BuildCellTrafficTrace(1, 2)
	encodeData, err := ngap.Encoder(pdu)
	if err != nil {
		t.Error(err.Error())
		t.Error("Encode is failed")
	} else {
		decodeData, err := ngap.Decoder(encodeData)
		if err != nil {
			t.Error(err.Error())
			t.Error("Decode is Failed")
		} else if reflect.DeepEqual(pdu, *decodeData) {
			fmt.Println("[PASS]")
		}
	}
}

func TestBuildPDUSessionResourceReleaseResponseForReleaseTest(t *testing.T) {
	pdu := ngapTestpacket.BuildPDUSessionResourceReleaseResponseForReleaseTest(1, 2)
	encodeData, err := ngap.Encoder(pdu)
	if err != nil {
		t.Error(err.Error())
		t.Error("Encode is failed")
	} else {
		decodeData, err := ngap.Decoder(encodeData)
		if err != nil {
			t.Error(err.Error())
			t.Error("Decode is Failed")
		} else if reflect.DeepEqual(pdu, *decodeData) {
			fmt.Println("[PASS]")
		}
	}
}
