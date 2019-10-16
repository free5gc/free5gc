//go:binary-only-package

package nasMessage_test

import (
	"bytes"
	"free5gc/lib/nas"
	"free5gc/lib/nas/logger"
	"free5gc/lib/nas/nasMessage"
	"free5gc/lib/nas/nasType"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

type nasMessageServiceRequestData struct {
	inExtendedProtocolDiscriminator uint8
	inSecurityHeader                uint8
	inSpareHalfOctet                uint8
	inServiceRequestMessageIdentity uint8
	inTMSI5GS                       nasType.TMSI5GS
	inUplinkDataStatus              nasType.UplinkDataStatus
	inPDUSessionStatus              nasType.PDUSessionStatus
	inAllowedPDUSessionStatus       nasType.AllowedPDUSessionStatus
	inNASMessageContainer           nasType.NASMessageContainer
}

var nasMessageServiceRequestTable = []nasMessageServiceRequestData{
	{
		inExtendedProtocolDiscriminator: nasMessage.Epd5GSMobilityManagementMessage,
		inSecurityHeader:                0x01,
		inSpareHalfOctet:                0x01,
		inServiceRequestMessageIdentity: nas.MsgTypeServiceRequest,
		inTMSI5GS: nasType.TMSI5GS{
			Len:   7,
			Octet: [7]uint8{0x01, 0x01},
		},
		inUplinkDataStatus: nasType.UplinkDataStatus{
			Iei:    nasMessage.ServiceRequestUplinkDataStatusType,
			Len:    2,
			Buffer: []uint8{0x01, 0x01},
		},
		inPDUSessionStatus: nasType.PDUSessionStatus{
			Iei:    nasMessage.ServiceRequestPDUSessionStatusType,
			Len:    2,
			Buffer: []uint8{0x01, 0x01},
		},
		inAllowedPDUSessionStatus: nasType.AllowedPDUSessionStatus{
			Iei:    nasMessage.ServiceRequestAllowedPDUSessionStatusType,
			Len:    2,
			Buffer: []uint8{0x01, 0x01},
		},
		inNASMessageContainer: nasType.NASMessageContainer{
			Iei:    nasMessage.ServiceRequestNASMessageContainerType,
			Len:    2,
			Buffer: []uint8{0x01, 0x01},
		},
	},
}

func TestNasTypeNewServiceRequest(t *testing.T) {}

func TestNasTypeNewServiceRequestMessage(t *testing.T) {}
