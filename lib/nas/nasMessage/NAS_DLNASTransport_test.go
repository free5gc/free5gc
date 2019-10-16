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

type nasMessageDLNASTransportData struct {
	inExtendedProtocolDiscriminator uint8
	inSecurityHeaderType            uint8
	inSpareHalfOctet1               uint8
	inDLNASTRANSPORTMessageIdentity uint8
	inPayloadContainerType          uint8
	inSpareHalfOctet2               uint8
	inPayloadContainer              nasType.PayloadContainer
	inPduSessionID2Value            nasType.PduSessionID2Value
	inAdditionalInformation         nasType.AdditionalInformation
	inCause5GMM                     nasType.Cause5GMM
	inBackoffTimerValue             nasType.BackoffTimerValue
}

var nasMessageDLNASTransportTable = []nasMessageDLNASTransportData{
	{
		inExtendedProtocolDiscriminator: nas.MsgTypeDLNASTransport,
		inSecurityHeaderType:            0x01,
		inSpareHalfOctet1:               0x01,
		inDLNASTRANSPORTMessageIdentity: 0x01,
		inPayloadContainerType:          0x01,
		inSpareHalfOctet2:               0x01,
		inPayloadContainer: nasType.PayloadContainer{
			Iei:    0,
			Len:    2,
			Buffer: []uint8{0x01, 0x01},
		},
		inPduSessionID2Value: nasType.PduSessionID2Value{
			Iei:   nasMessage.DLNASTransportPduSessionID2ValueType,
			Octet: 0x01,
		},
		inAdditionalInformation: nasType.AdditionalInformation{
			Iei:    nasMessage.DLNASTransportAdditionalInformationType,
			Len:    2,
			Buffer: []uint8{0x01, 0x01},
		},
		inCause5GMM: nasType.Cause5GMM{
			Iei:   nasMessage.DLNASTransportCause5GMMType,
			Octet: 0xF0,
		},
		inBackoffTimerValue: nasType.BackoffTimerValue{
			Iei:   nasMessage.DLNASTransportBackoffTimerValueType,
			Len:   2,
			Octet: 0x01,
		},
	},
}

func TestNasTypeNewDLNASTransport(t *testing.T) {}

func TestNasTypeNewDLNASTransportMessage(t *testing.T) {}
