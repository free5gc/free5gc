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

type nasMessageConfigurationUpdateCommandData struct {
	inExtendedProtocolDiscriminator             uint8
	inSecurityHeaderType                        uint8
	inSpareHalfOctet                            uint8
	inConfigurationUpdateCommandMessageIdentity uint8
	inConfigurationUpdateIndication             nasType.ConfigurationUpdateIndication
	inGUTI5G                                    nasType.GUTI5G
	inTAIList                                   nasType.TAIList
	inAllowedNSSAI                              nasType.AllowedNSSAI
	inServiceAreaList                           nasType.ServiceAreaList
	inFullNameForNetwork                        nasType.FullNameForNetwork
	inShortNameForNetwork                       nasType.ShortNameForNetwork
	inLocalTimeZone                             nasType.LocalTimeZone
	inUniversalTimeAndLocalTimeZone             nasType.UniversalTimeAndLocalTimeZone
	inNetworkDaylightSavingTime                 nasType.NetworkDaylightSavingTime
	inLADNInformation                           nasType.LADNInformation
	inMICOIndication                            nasType.MICOIndication
	inNetworkSlicingIndication                  nasType.NetworkSlicingIndication
	inConfiguredNSSAI                           nasType.ConfiguredNSSAI
	inRejectedNSSAI                             nasType.RejectedNSSAI
	inOperatordefinedAccessCategoryDefinitions  nasType.OperatordefinedAccessCategoryDefinitions
	inSMSIndication                             nasType.SMSIndication
}

var nasMessageConfigurationUpdateCommandTable = []nasMessageConfigurationUpdateCommandData{
	{
		inExtendedProtocolDiscriminator:             nasMessage.Epd5GSSessionManagementMessage,
		inSecurityHeaderType:                        0x01,
		inSpareHalfOctet:                            0x01,
		inConfigurationUpdateCommandMessageIdentity: nas.MsgTypeConfigurationUpdateCommand,
		inConfigurationUpdateIndication: nasType.ConfigurationUpdateIndication{
			Octet: 0xD0,
		},
		inGUTI5G: nasType.GUTI5G{
			Iei:   nasMessage.ConfigurationUpdateCommandGUTI5GType,
			Len:   11,
			Octet: [11]uint8{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A, 0x0B},
		},
		inTAIList: nasType.TAIList{
			Iei:    nasMessage.ConfigurationUpdateCommandTAIListType,
			Len:    2,
			Buffer: []uint8{0x01, 0x01},
		},
		inAllowedNSSAI: nasType.AllowedNSSAI{
			Iei:    nasMessage.ConfigurationUpdateCommandAllowedNSSAIType,
			Len:    2,
			Buffer: []uint8{0x01, 0x01},
		},
		inServiceAreaList: nasType.ServiceAreaList{
			Iei:    nasMessage.ConfigurationUpdateCommandServiceAreaListType,
			Len:    2,
			Buffer: []uint8{0x01, 0x01},
		},
		inFullNameForNetwork: nasType.FullNameForNetwork{
			Iei:    nasMessage.ConfigurationUpdateCommandFullNameForNetworkType,
			Len:    2,
			Buffer: []uint8{0x01, 0x01},
		},
		inShortNameForNetwork: nasType.ShortNameForNetwork{
			Iei:    nasMessage.ConfigurationUpdateCommandShortNameForNetworkType,
			Len:    2,
			Buffer: []uint8{0x01, 0x01},
		},
		inLocalTimeZone: nasType.LocalTimeZone{
			Iei:   nasMessage.ConfigurationUpdateCommandLocalTimeZoneType,
			Octet: 0x01,
		},
		inUniversalTimeAndLocalTimeZone: nasType.UniversalTimeAndLocalTimeZone{
			Iei:   nasMessage.ConfigurationUpdateCommandUniversalTimeAndLocalTimeZoneType,
			Octet: [7]uint8{0x01},
		},
		inNetworkDaylightSavingTime: nasType.NetworkDaylightSavingTime{
			Iei:   nasMessage.ConfigurationUpdateCommandNetworkDaylightSavingTimeType,
			Len:   2,
			Octet: 0x01,
		},
		inLADNInformation: nasType.LADNInformation{
			Iei:    nasMessage.ConfigurationUpdateCommandLADNInformationType,
			Len:    2,
			Buffer: []uint8{0x01, 0x01},
		},
		inMICOIndication: nasType.MICOIndication{
			Octet: 0xB0,
		},
		inNetworkSlicingIndication: nasType.NetworkSlicingIndication{
			Octet: 0x90,
		},
		inConfiguredNSSAI: nasType.ConfiguredNSSAI{
			Iei:    nasMessage.ConfigurationUpdateCommandConfiguredNSSAIType,
			Len:    2,
			Buffer: []uint8{0x01, 0x01},
		},
		inRejectedNSSAI: nasType.RejectedNSSAI{
			Iei:    nasMessage.ConfigurationUpdateCommandRejectedNSSAIType,
			Len:    2,
			Buffer: []uint8{0x01, 0x01},
		},
		inOperatordefinedAccessCategoryDefinitions: nasType.OperatordefinedAccessCategoryDefinitions{
			Iei:    nasMessage.ConfigurationUpdateCommandOperatordefinedAccessCategoryDefinitionsType,
			Len:    2,
			Buffer: []uint8{0x01, 0x01},
		},
		inSMSIndication: nasType.SMSIndication{
			Octet: 0xF0,
		},
	},
}

func TestNasTypeNewConfigurationUpdateCommand(t *testing.T) {}

func TestNasTypeNewConfigurationUpdateCommandMessage(t *testing.T) {}
