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

type nasMessageSecurityModeCommandData struct {
	inExtendedProtocolDiscriminator      uint8
	inSecurityHeader                     uint8
	inSpareHalfOctet                     uint8
	inSecurityModeCommandMessageIdentity uint8
	inSelectedNASSecurityAlgorithms      nasType.SelectedNASSecurityAlgorithms
	inNgksi                              uint8
	inReplayedUESecurityCapabilities     nasType.ReplayedUESecurityCapabilities
	inIMEISVRequest                      nasType.IMEISVRequest
	inSelectedEPSNASSecurityAlgorithms   nasType.SelectedEPSNASSecurityAlgorithms
	inAdditional5GSecurityInformation    nasType.Additional5GSecurityInformation
	inEAPMessage                         nasType.EAPMessage
	inABBA                               nasType.ABBA
	inReplayedS1UESecurityCapabilities   nasType.ReplayedS1UESecurityCapabilities
}

var nasMessageSecurityModeCommandTable = []nasMessageSecurityModeCommandData{
	{
		inExtendedProtocolDiscriminator:      nasMessage.Epd5GSMobilityManagementMessage,
		inSecurityHeader:                     0x01,
		inSpareHalfOctet:                     0x01,
		inSecurityModeCommandMessageIdentity: nas.MsgTypeSecurityModeCommand,
		inSelectedNASSecurityAlgorithms: nasType.SelectedNASSecurityAlgorithms{
			Octet: 0x01,
		},
		inNgksi: 0x01,
		inReplayedUESecurityCapabilities: nasType.ReplayedUESecurityCapabilities{
			Len:    8,
			Buffer: []uint8{0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01},
		},
		inIMEISVRequest: nasType.IMEISVRequest{
			Octet: 0xE0,
		},
		inSelectedEPSNASSecurityAlgorithms: nasType.SelectedEPSNASSecurityAlgorithms{
			Iei:   nasMessage.SecurityModeCommandSelectedEPSNASSecurityAlgorithmsType,
			Octet: 0x01,
		},
		inAdditional5GSecurityInformation: nasType.Additional5GSecurityInformation{
			Iei:   nasMessage.SecurityModeCommandAdditional5GSecurityInformationType,
			Len:   2,
			Octet: 0x01,
		},
		inEAPMessage: nasType.EAPMessage{
			Iei:    nasMessage.SecurityModeCommandEAPMessageType,
			Len:    2,
			Buffer: []uint8{0x01, 0x01},
		},
		inABBA: nasType.ABBA{
			Iei:    nasMessage.SecurityModeCommandABBAType,
			Len:    2,
			Buffer: []uint8{0x01, 0x01},
		},
		inReplayedS1UESecurityCapabilities: nasType.ReplayedS1UESecurityCapabilities{
			Iei:    nasMessage.SecurityModeCommandReplayedS1UESecurityCapabilitiesType,
			Len:    5,
			Buffer: []uint8{0x01, 0x01, 0x01, 0x01, 0x01},
		},
	},
}

func TestNasTypeNewSecurityModeCommand(t *testing.T) {}

func TestNasTypeNewSecurityModeCommandMessage(t *testing.T) {}
