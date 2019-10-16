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

type nasMessageStatus5GSMData struct {
	inExtendedProtocolDiscriminator uint8
	inPDUSessionID                  nasType.PDUSessionID
	inPTI                           nasType.PTI
	inStatus5GSMMessageIdentity     uint8
	inCause5GSM                     nasType.Cause5GSM
}

var nasMessageStatus5GSMTable = []nasMessageStatus5GSMData{
	{
		inExtendedProtocolDiscriminator: nasMessage.Epd5GSSessionManagementMessage,
		inPDUSessionID: nasType.PDUSessionID{
			Octet: 0x01,
		},
		inPTI: nasType.PTI{
			Octet: 0x01,
		},
		inStatus5GSMMessageIdentity: nas.MsgTypeStatus5GSM,
		inCause5GSM: nasType.Cause5GSM{
			Octet: 0x01,
		},
	},
}

func TestNasTypeNewStatus5GSM(t *testing.T) {}

func TestNasTypeNewStatus5GSMMessage(t *testing.T) {}
