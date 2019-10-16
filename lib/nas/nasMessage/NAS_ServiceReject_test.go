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

type nasMessageServiceRejectData struct {
	inExtendedProtocolDiscriminator uint8
	inSecurityHeader                uint8
	inSpareHalfOctet                uint8
	inServiceRejectMessageIdentity  uint8
	inCause5GMM                     nasType.Cause5GMM
	inPDUSessionStatus              nasType.PDUSessionStatus
	inT3346Value                    nasType.T3346Value
	inEAPMessage                    nasType.EAPMessage
}

var nasMessageServiceRejectTable = []nasMessageServiceRejectData{
	{
		inExtendedProtocolDiscriminator: nasMessage.Epd5GSMobilityManagementMessage,
		inSecurityHeader:                0x01,
		inSpareHalfOctet:                0x01,
		inServiceRejectMessageIdentity:  nas.MsgTypeServiceReject,
		inCause5GMM: nasType.Cause5GMM{
			Octet: 0x01,
		},
		inPDUSessionStatus: nasType.PDUSessionStatus{
			Iei:    nasMessage.ServiceRejectPDUSessionStatusType,
			Len:    2,
			Buffer: []uint8{0x01, 0x01},
		},
		inT3346Value: nasType.T3346Value{
			Iei:   nasMessage.ServiceRejectT3346ValueType,
			Len:   2,
			Octet: 0x01,
		},
		inEAPMessage: nasType.EAPMessage{
			Iei:    nasMessage.ServiceRejectEAPMessageType,
			Len:    2,
			Buffer: []uint8{0x01, 0x01},
		},
	},
}

func TestNasTypeNewServiceReject(t *testing.T) {}

func TestNasTypeNewServiceRejectMessage(t *testing.T) {}
