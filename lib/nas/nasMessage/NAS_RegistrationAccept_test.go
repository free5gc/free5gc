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

type nasMessageRegistrationAcceptData struct {
	inExtendedProtocolDiscriminator            uint8
	inSecurityHeader                           uint8
	inSpareHalfOctet                           uint8
	inRegistrationAcceptMessageIdentity        uint8
	inRegistrationResult5GS                    nasType.RegistrationResult5GS
	inGUTI5G                                   nasType.GUTI5G
	inEquivalentPlmns                          nasType.EquivalentPlmns
	inTAIList                                  nasType.TAIList
	inAllowedNSSAI                             nasType.AllowedNSSAI
	inRejectedNSSAI                            nasType.RejectedNSSAI
	inConfiguredNSSAI                          nasType.ConfiguredNSSAI
	inNetworkFeatureSupport5GS                 nasType.NetworkFeatureSupport5GS
	inPDUSessionStatus                         nasType.PDUSessionStatus
	inPDUSessionReactivationResult             nasType.PDUSessionReactivationResult
	inPDUSessionReactivationResultErrorCause   nasType.PDUSessionReactivationResultErrorCause
	inLADNInformation                          nasType.LADNInformation
	inMICOIndication                           nasType.MICOIndication
	inNetworkSlicingIndication                 nasType.NetworkSlicingIndication
	inServiceAreaList                          nasType.ServiceAreaList
	inT3512Value                               nasType.T3512Value
	inNon3GppDeregistrationTimerValue          nasType.Non3GppDeregistrationTimerValue
	inT3502Value                               nasType.T3502Value
	inEmergencyNumberList                      nasType.EmergencyNumberList
	inExtendedEmergencyNumberList              nasType.ExtendedEmergencyNumberList
	inSORTransparentContainer                  nasType.SORTransparentContainer
	inEAPMessage                               nasType.EAPMessage
	inNSSAIInclusionMode                       nasType.NSSAIInclusionMode
	inOperatordefinedAccessCategoryDefinitions nasType.OperatordefinedAccessCategoryDefinitions
	inNegotiatedDRXParameters                  nasType.NegotiatedDRXParameters
}

var nasMessageRegistrationAcceptTable = []nasMessageRegistrationAcceptData{
	{
		inExtendedProtocolDiscriminator:     nasMessage.Epd5GSMobilityManagementMessage,
		inSecurityHeader:                    0x01,
		inSpareHalfOctet:                    0x01,
		inRegistrationAcceptMessageIdentity: nas.MsgTypeRegistrationAccept,
		inRegistrationResult5GS: nasType.RegistrationResult5GS{
			Len:   1,
			Octet: 0x01,
		},
		inGUTI5G: nasType.GUTI5G{
			Iei:   nasMessage.RegistrationAcceptGUTI5GType,
			Len:   11,
			Octet: [11]uint8{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A, 0x0B},
		},
		inEquivalentPlmns: nasType.EquivalentPlmns{
			Iei:   nasMessage.RegistrationAcceptEquivalentPlmnsType,
			Len:   45,
			Octet: [45]uint8{0x01, 0x01},
		},
		inTAIList: nasType.TAIList{
			Iei:    nasMessage.RegistrationAcceptTAIListType,
			Len:    2,
			Buffer: []uint8{0x01, 0x01},
		},
		inAllowedNSSAI: nasType.AllowedNSSAI{
			Iei:    nasMessage.RegistrationAcceptAllowedNSSAIType,
			Len:    2,
			Buffer: []uint8{0x01, 0x01},
		},
		inRejectedNSSAI: nasType.RejectedNSSAI{
			Iei:    nasMessage.RegistrationAcceptRejectedNSSAIType,
			Len:    2,
			Buffer: []uint8{0x01, 0x01},
		},
		inConfiguredNSSAI: nasType.ConfiguredNSSAI{
			Iei:    nasMessage.RegistrationAcceptConfiguredNSSAIType,
			Len:    2,
			Buffer: []uint8{0x01, 0x01},
		},
		inNetworkFeatureSupport5GS: nasType.NetworkFeatureSupport5GS{
			Iei:   nasMessage.RegistrationAcceptNetworkFeatureSupport5GSType,
			Len:   3,
			Octet: [3]uint8{0x01, 0x01, 0x01},
		},
		inPDUSessionStatus: nasType.PDUSessionStatus{
			Iei:    nasMessage.RegistrationAcceptPDUSessionStatusType,
			Len:    2,
			Buffer: []uint8{0x01, 0x01},
		},
		inPDUSessionReactivationResult: nasType.PDUSessionReactivationResult{
			Iei:    nasMessage.RegistrationAcceptPDUSessionReactivationResultType,
			Len:    2,
			Buffer: []uint8{0x01, 0x01},
		},
		inPDUSessionReactivationResultErrorCause: nasType.PDUSessionReactivationResultErrorCause{
			Iei:    nasMessage.RegistrationAcceptPDUSessionReactivationResultErrorCauseType,
			Len:    2,
			Buffer: []uint8{0x01, 0x01},
		},
		inLADNInformation: nasType.LADNInformation{
			Iei:    nasMessage.RegistrationAcceptLADNInformationType,
			Len:    2,
			Buffer: []uint8{0x01, 0x01},
		},
		inMICOIndication: nasType.MICOIndication{
			Octet: 0xB0,
		},
		inNetworkSlicingIndication: nasType.NetworkSlicingIndication{
			Octet: 0x90,
		},
		inServiceAreaList: nasType.ServiceAreaList{
			Iei:    nasMessage.RegistrationAcceptServiceAreaListType,
			Len:    2,
			Buffer: []uint8{0x01, 0x01},
		},
		inT3512Value: nasType.T3512Value{
			Iei:   nasMessage.RegistrationAcceptT3512ValueType,
			Len:   1,
			Octet: 0x01,
		},
		inNon3GppDeregistrationTimerValue: nasType.Non3GppDeregistrationTimerValue{
			Iei:   nasMessage.RegistrationAcceptNon3GppDeregistrationTimerValueType,
			Len:   1,
			Octet: 0x01,
		},
		inT3502Value: nasType.T3502Value{
			Iei:   nasMessage.RegistrationAcceptT3502ValueType,
			Len:   1,
			Octet: 0x01,
		},
		inEmergencyNumberList: nasType.EmergencyNumberList{
			Iei:    nasMessage.RegistrationAcceptEmergencyNumberListType,
			Len:    2,
			Buffer: []uint8{0x01, 0x01},
		},
		inExtendedEmergencyNumberList: nasType.ExtendedEmergencyNumberList{
			Iei:    nasMessage.RegistrationAcceptExtendedEmergencyNumberListType,
			Len:    2,
			Buffer: []uint8{0x01, 0x01},
		},
		inSORTransparentContainer: nasType.SORTransparentContainer{
			Iei:    nasMessage.RegistrationAcceptSORTransparentContainerType,
			Len:    2,
			Buffer: []uint8{0x01, 0x01},
		},
		inEAPMessage: nasType.EAPMessage{
			Iei:    nasMessage.RegistrationAcceptEAPMessageType,
			Len:    2,
			Buffer: []uint8{0x01, 0x01},
		},
		inNSSAIInclusionMode: nasType.NSSAIInclusionMode{
			Octet: 0xA0,
		},
		inOperatordefinedAccessCategoryDefinitions: nasType.OperatordefinedAccessCategoryDefinitions{
			Iei:    nasMessage.RegistrationAcceptOperatordefinedAccessCategoryDefinitionsType,
			Len:    2,
			Buffer: []uint8{0x01, 0x01},
		},
		inNegotiatedDRXParameters: nasType.NegotiatedDRXParameters{
			Iei:   nasMessage.RegistrationAcceptNegotiatedDRXParametersType,
			Len:   1,
			Octet: 0x01,
		},
	},
}

func TestNasTypeNewRegistrationAccept(t *testing.T) {}

func TestNasTypeNewRegistrationAcceptMessage(t *testing.T) {}
