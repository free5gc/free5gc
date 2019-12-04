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

type nasMessageSecurityModeCompleteData struct {
	inExtendedProtocolDiscriminator       uint8
	inSecurityHeader                      uint8
	inSpareHalfOctet                      uint8
	inSecurityModeCompleteMessageIdentity uint8
	inIMEISV                              nasType.IMEISV
	inNASMessageContainer                 nasType.NASMessageContainer
}

var nasMessageSecurityModeCompleteTable = []nasMessageSecurityModeCompleteData{
	{
		inExtendedProtocolDiscriminator:       nasMessage.Epd5GSMobilityManagementMessage,
		inSecurityHeader:                      0x01,
		inSpareHalfOctet:                      0x01,
		inSecurityModeCompleteMessageIdentity: nas.MsgTypeSecurityModeComplete,
		inIMEISV: nasType.IMEISV{
			Iei:   nasMessage.SecurityModeCompleteIMEISVType,
			Len:   2,
			Octet: [9]uint8{0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01},
		},
		inNASMessageContainer: nasType.NASMessageContainer{
			Iei:    nasMessage.SecurityModeCompleteNASMessageContainerType,
			Len:    2,
			Buffer: []uint8{0x01, 0x01},
		},
	},
}

func TestNasTypeNewSecurityModeComplete(t *testing.T) {}

func TestNasTypeNewSecurityModeCompleteMessage(t *testing.T) {}
