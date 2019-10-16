//go:binary-only-package

package nasMessage

import (
	"bytes"
	"encoding/binary"
	"free5gc/lib/nas/nasType"
)

type ConfigurationUpdateCommand struct {
	nasType.ExtendedProtocolDiscriminator
	nasType.SpareHalfOctetAndSecurityHeaderType
	nasType.ConfigurationUpdateCommandMessageIdentity
	*nasType.ConfigurationUpdateIndication
	*nasType.GUTI5G
	*nasType.TAIList
	*nasType.AllowedNSSAI
	*nasType.ServiceAreaList
	*nasType.FullNameForNetwork
	*nasType.ShortNameForNetwork
	*nasType.LocalTimeZone
	*nasType.UniversalTimeAndLocalTimeZone
	*nasType.NetworkDaylightSavingTime
	*nasType.LADNInformation
	*nasType.MICOIndication
	*nasType.NetworkSlicingIndication
	*nasType.ConfiguredNSSAI
	*nasType.RejectedNSSAI
	*nasType.OperatordefinedAccessCategoryDefinitions
	*nasType.SMSIndication
}

func NewConfigurationUpdateCommand(iei uint8) (configurationUpdateCommand *ConfigurationUpdateCommand) {}

const (
	ConfigurationUpdateCommandConfigurationUpdateIndicationType            uint8 = 0x0D
	ConfigurationUpdateCommandGUTI5GType                                   uint8 = 0x77
	ConfigurationUpdateCommandTAIListType                                  uint8 = 0x54
	ConfigurationUpdateCommandAllowedNSSAIType                             uint8 = 0x15
	ConfigurationUpdateCommandServiceAreaListType                          uint8 = 0x27
	ConfigurationUpdateCommandFullNameForNetworkType                       uint8 = 0x43
	ConfigurationUpdateCommandShortNameForNetworkType                      uint8 = 0x45
	ConfigurationUpdateCommandLocalTimeZoneType                            uint8 = 0x46
	ConfigurationUpdateCommandUniversalTimeAndLocalTimeZoneType            uint8 = 0x47
	ConfigurationUpdateCommandNetworkDaylightSavingTimeType                uint8 = 0x49
	ConfigurationUpdateCommandLADNInformationType                          uint8 = 0x79
	ConfigurationUpdateCommandMICOIndicationType                           uint8 = 0x0B
	ConfigurationUpdateCommandNetworkSlicingIndicationType                 uint8 = 0x09
	ConfigurationUpdateCommandConfiguredNSSAIType                          uint8 = 0x31
	ConfigurationUpdateCommandRejectedNSSAIType                            uint8 = 0x11
	ConfigurationUpdateCommandOperatordefinedAccessCategoryDefinitionsType uint8 = 0x76
	ConfigurationUpdateCommandSMSIndicationType                            uint8 = 0x0F
)

func (a *ConfigurationUpdateCommand) EncodeConfigurationUpdateCommand(buffer *bytes.Buffer) {}

func (a *ConfigurationUpdateCommand) DecodeConfigurationUpdateCommand(byteArray *[]byte) {}
