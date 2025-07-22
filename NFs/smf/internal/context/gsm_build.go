package context

import (
	"encoding/hex"

	"github.com/free5gc/nas"
	"github.com/free5gc/nas/nasConvert"
	"github.com/free5gc/nas/nasMessage"
	"github.com/free5gc/nas/nasType"
	"github.com/free5gc/smf/internal/logger"
)

func BuildGSMPDUSessionEstablishmentAccept(smContext *SMContext) ([]byte, error) {
	m := nas.NewMessage()
	m.GsmMessage = nas.NewGsmMessage()
	m.GsmHeader.SetMessageType(nas.MsgTypePDUSessionEstablishmentAccept)
	m.GsmHeader.SetExtendedProtocolDiscriminator(nasMessage.Epd5GSSessionManagementMessage)
	m.PDUSessionEstablishmentAccept = nasMessage.NewPDUSessionEstablishmentAccept(0x0)
	pDUSessionEstablishmentAccept := m.PDUSessionEstablishmentAccept

	sessRule := smContext.SelectedSessionRule()
	authDefQos := sessRule.AuthDefQos

	pDUSessionEstablishmentAccept.SetPDUSessionID(uint8(smContext.PDUSessionID))
	pDUSessionEstablishmentAccept.SetMessageType(nas.MsgTypePDUSessionEstablishmentAccept)
	pDUSessionEstablishmentAccept.SetExtendedProtocolDiscriminator(nasMessage.Epd5GSSessionManagementMessage)
	pDUSessionEstablishmentAccept.SetPTI(smContext.Pti)

	if v := smContext.EstAcceptCause5gSMValue; v != 0 {
		pDUSessionEstablishmentAccept.Cause5GSM = nasType.NewCause5GSM(nasMessage.PDUSessionEstablishmentAcceptCause5GSMType)
		pDUSessionEstablishmentAccept.Cause5GSM.SetCauseValue(v)
	}
	pDUSessionEstablishmentAccept.SetPDUSessionType(smContext.SelectedPDUSessionType)

	pDUSessionEstablishmentAccept.SetSSCMode(1)
	pDUSessionEstablishmentAccept.SessionAMBR = nasConvert.ModelsToSessionAMBR(sessRule.AuthSessAmbr)
	pDUSessionEstablishmentAccept.SessionAMBR.SetLen(uint8(len(pDUSessionEstablishmentAccept.SessionAMBR.Octet)))

	qoSRules := nasType.QoSRules{
		{
			Identifier: smContext.defRuleID,
			DQR:        true,
			Operation:  nasType.OperationCodeCreateNewQoSRule,
			Precedence: 255,
			QFI:        sessRule.DefQosQFI,
			PacketFilterList: nasType.PacketFilterList{
				{
					Identifier: 1,
					Direction:  nasType.PacketFilterDirectionBidirectional,
					Components: nasType.PacketFilterComponentList{
						&nasType.PacketFilterMatchAll{},
					},
				},
			},
		},
	}

	for _, pccRule := range smContext.PCCRules {
		if qosRule, err1 := pccRule.BuildNasQoSRule(smContext,
			nasType.OperationCodeCreateNewQoSRule); err1 != nil {
			logger.GsmLog.Warnln("Create QoS rule from pcc error ", err1)
		} else {
			if ruleID, err2 := smContext.QoSRuleIDGenerator.Allocate(); err2 != nil {
				return nil, err2
			} else {
				qosRule.Identifier = uint8(ruleID)
				smContext.PCCRuleIDToQoSRuleID[pccRule.PccRuleId] = uint8(ruleID)
			}
			qoSRules = append(qoSRules, *qosRule)
		}
	}

	qosRulesBytes, errMarshalBinary := qoSRules.MarshalBinary()
	if errMarshalBinary != nil {
		return nil, errMarshalBinary
	}

	pDUSessionEstablishmentAccept.AuthorizedQosRules.SetLen(uint16(len(qosRulesBytes)))
	pDUSessionEstablishmentAccept.AuthorizedQosRules.SetQosRule(qosRulesBytes)

	if smContext.PDUAddress != nil {
		addr, addrLen := smContext.PDUAddressToNAS()
		pDUSessionEstablishmentAccept.PDUAddress = nasType.
			NewPDUAddress(nasMessage.PDUSessionEstablishmentAcceptPDUAddressType)
		pDUSessionEstablishmentAccept.PDUAddress.SetLen(addrLen)
		pDUSessionEstablishmentAccept.PDUAddress.SetPDUSessionTypeValue(smContext.SelectedPDUSessionType)
		pDUSessionEstablishmentAccept.PDUAddress.SetPDUAddressInformation(addr)
	}

	authDescs := nasType.QoSFlowDescs{}
	dafaultAuthDesc := nasType.QoSFlowDesc{}
	dafaultAuthDesc.QFI = sessRule.DefQosQFI
	dafaultAuthDesc.OperationCode = nasType.OperationCodeCreateNewQoSFlowDescription
	parameter := new(nasType.QoSFlow5QI)
	parameter.FiveQI = uint8(authDefQos.Var5qi)
	dafaultAuthDesc.Parameters = append(dafaultAuthDesc.Parameters, parameter)
	authDescs = append(authDescs, dafaultAuthDesc)
	for _, qosFlow := range smContext.AdditonalQosFlows {
		if qosDesc, e := qosFlow.BuildNasQoSDesc(nasType.OperationCodeCreateNewQoSFlowDescription); e != nil {
			logger.GsmLog.Warnf("Create QoS Desc from qos flow error: %s\n", e)
		} else {
			authDescs = append(authDescs, qosDesc)
		}
	}
	qosDescBytes, errMarshalBinary := authDescs.MarshalBinary()
	if errMarshalBinary != nil {
		return nil, errMarshalBinary
	}
	pDUSessionEstablishmentAccept.AuthorizedQosFlowDescriptions = nasType.
		NewAuthorizedQosFlowDescriptions(nasMessage.PDUSessionEstablishmentAcceptAuthorizedQosFlowDescriptionsType)
	pDUSessionEstablishmentAccept.AuthorizedQosFlowDescriptions.SetLen(uint16(len(qosDescBytes)))
	pDUSessionEstablishmentAccept.SetQoSFlowDescriptions(qosDescBytes)

	var sd [3]uint8

	if byteArray, errDecodeString := hex.DecodeString(smContext.SNssai.Sd); errDecodeString != nil {
		return nil, errDecodeString
	} else {
		copy(sd[:], byteArray)
	}

	pDUSessionEstablishmentAccept.SNSSAI = nasType.NewSNSSAI(nasMessage.ULNASTransportSNSSAIType)
	pDUSessionEstablishmentAccept.SNSSAI.SetLen(4)
	pDUSessionEstablishmentAccept.SNSSAI.SetSST(uint8(smContext.SNssai.Sst))
	pDUSessionEstablishmentAccept.SNSSAI.SetSD(sd)

	pDUSessionEstablishmentAccept.DNN = nasType.NewDNN(nasMessage.ULNASTransportDNNType)
	pDUSessionEstablishmentAccept.DNN.SetDNN(smContext.Dnn)

	if smContext.ProtocolConfigurationOptions.DNSIPv4Request ||
		smContext.ProtocolConfigurationOptions.DNSIPv6Request ||
		smContext.ProtocolConfigurationOptions.PCSCFIPv4Request ||
		smContext.ProtocolConfigurationOptions.IPv4LinkMTURequest {
		pDUSessionEstablishmentAccept.ExtendedProtocolConfigurationOptions = nasType.NewExtendedProtocolConfigurationOptions(
			nasMessage.PDUSessionEstablishmentAcceptExtendedProtocolConfigurationOptionsType,
		)
		protocolConfigurationOptions := nasConvert.NewProtocolConfigurationOptions()

		// IPv4 DNS
		if smContext.ProtocolConfigurationOptions.DNSIPv4Request {
			errAddDNSServerIPv4Address := protocolConfigurationOptions.AddDNSServerIPv4Address(smContext.DNNInfo.DNS.IPv4Addr)
			if errAddDNSServerIPv4Address != nil {
				logger.GsmLog.Warnln("Error while adding DNS IPv4 Addr: ", errAddDNSServerIPv4Address)
			}
		}

		// IPv6 DNS
		if smContext.ProtocolConfigurationOptions.DNSIPv6Request {
			errAddDNSServerIPv6Address := protocolConfigurationOptions.AddDNSServerIPv6Address(smContext.DNNInfo.DNS.IPv6Addr)
			if errAddDNSServerIPv6Address != nil {
				logger.GsmLog.Warnln("Error while adding DNS IPv6 Addr: ", errAddDNSServerIPv6Address)
			}
		}

		// IPv4 PCSCF (need for ims DNN)
		if smContext.ProtocolConfigurationOptions.PCSCFIPv4Request {
			errAddPCSCFIPv4Address := protocolConfigurationOptions.AddPCSCFIPv4Address(smContext.DNNInfo.PCSCF.IPv4Addr)
			if errAddPCSCFIPv4Address != nil {
				logger.GsmLog.Warnln("Error while adding PCSCF IPv4 Addr: ", errAddPCSCFIPv4Address)
			}
		}

		// MTU
		if smContext.ProtocolConfigurationOptions.IPv4LinkMTURequest {
			errAddIPv4LinkMTU := protocolConfigurationOptions.AddIPv4LinkMTU(1400)
			if errAddIPv4LinkMTU != nil {
				logger.GsmLog.Warnln("Error while adding MTU: ", errAddIPv4LinkMTU)
			}
		}

		pcoContents := protocolConfigurationOptions.Marshal()
		pcoContentsLength := len(pcoContents)
		pDUSessionEstablishmentAccept.
			ExtendedProtocolConfigurationOptions.
			SetLen(uint16(pcoContentsLength))
		pDUSessionEstablishmentAccept.
			ExtendedProtocolConfigurationOptions.
			SetExtendedProtocolConfigurationOptionsContents(pcoContents)
	}
	return m.PlainNasEncode()
}

func BuildGSMPDUSessionEstablishmentReject(smContext *SMContext, cause uint8) ([]byte, error) {
	m := nas.NewMessage()
	m.GsmMessage = nas.NewGsmMessage()
	m.GsmHeader.SetMessageType(nas.MsgTypePDUSessionEstablishmentReject)
	m.GsmHeader.SetExtendedProtocolDiscriminator(nasMessage.Epd5GSSessionManagementMessage)
	m.PDUSessionEstablishmentReject = nasMessage.NewPDUSessionEstablishmentReject(0x0)
	pDUSessionEstablishmentReject := m.PDUSessionEstablishmentReject

	pDUSessionEstablishmentReject.SetMessageType(nas.MsgTypePDUSessionEstablishmentReject)
	pDUSessionEstablishmentReject.SetPTI(smContext.Pti)
	pDUSessionEstablishmentReject.SetExtendedProtocolDiscriminator(nasMessage.Epd5GSSessionManagementMessage)
	pDUSessionEstablishmentReject.SetPDUSessionID(uint8(smContext.PDUSessionID))
	pDUSessionEstablishmentReject.SetCauseValue(cause)

	return m.PlainNasEncode()
}

// BuildGSMPDUSessionReleaseCommand makes a plain NAS message.
//
// If isTriggeredByUE is true, the PTI field of the constructed NAS message is
// the value of smContext.Pti which is received from UE, otherwise it is 0.
// ref. 6.3.3.2 Network-requested PDU session release procedure initiation in TS24.501.
func BuildGSMPDUSessionReleaseCommand(smContext *SMContext, cause uint8, isTriggeredByUE bool) ([]byte, error) {
	m := nas.NewMessage()
	m.GsmMessage = nas.NewGsmMessage()
	m.GsmHeader.SetMessageType(nas.MsgTypePDUSessionReleaseCommand)
	m.GsmHeader.SetExtendedProtocolDiscriminator(nasMessage.Epd5GSSessionManagementMessage)
	m.PDUSessionReleaseCommand = nasMessage.NewPDUSessionReleaseCommand(0x0)
	pDUSessionReleaseCommand := m.PDUSessionReleaseCommand

	pDUSessionReleaseCommand.SetMessageType(nas.MsgTypePDUSessionReleaseCommand)
	pDUSessionReleaseCommand.SetExtendedProtocolDiscriminator(nasMessage.Epd5GSSessionManagementMessage)
	pDUSessionReleaseCommand.SetPDUSessionID(uint8(smContext.PDUSessionID))

	if isTriggeredByUE {
		pDUSessionReleaseCommand.SetPTI(smContext.Pti)
	} else {
		pDUSessionReleaseCommand.SetPTI(0x00)
	}
	pDUSessionReleaseCommand.SetCauseValue(cause)

	return m.PlainNasEncode()
}

func BuildGSMPDUSessionModificationCommand(smContext *SMContext) ([]byte, error) {
	m := nas.NewMessage()
	m.GsmMessage = nas.NewGsmMessage()
	m.GsmHeader.SetMessageType(nas.MsgTypePDUSessionModificationCommand)
	m.GsmHeader.SetExtendedProtocolDiscriminator(nasMessage.Epd5GSSessionManagementMessage)
	m.PDUSessionModificationCommand = nasMessage.NewPDUSessionModificationCommand(0x0)
	pDUSessionModificationCommand := m.PDUSessionModificationCommand

	pDUSessionModificationCommand.SetExtendedProtocolDiscriminator(nasMessage.Epd5GSSessionManagementMessage)
	pDUSessionModificationCommand.SetPDUSessionID(uint8(smContext.PDUSessionID))
	pDUSessionModificationCommand.SetPTI(smContext.Pti)
	pDUSessionModificationCommand.SetMessageType(nas.MsgTypePDUSessionModificationCommand)
	// pDUSessionModificationCommand.SetQosRule()
	// pDUSessionModificationCommand.AuthorizedQosRules.SetLen()
	// pDUSessionModificationCommand.SessionAMBR.SetSessionAMBRForDownlink([2]uint8{0x11, 0x11})
	// pDUSessionModificationCommand.SessionAMBR.SetSessionAMBRForUplink([2]uint8{0x11, 0x11})
	// pDUSessionModificationCommand.SessionAMBR.SetUnitForSessionAMBRForDownlink(10)
	// pDUSessionModificationCommand.SessionAMBR.SetUnitForSessionAMBRForUplink(10)
	// pDUSessionModificationCommand.SessionAMBR.SetLen(uint8(len(pDUSessionModificationCommand.SessionAMBR.Octet)))

	return m.PlainNasEncode()
}

func BuildGSMPDUSessionReleaseReject(smContext *SMContext) ([]byte, error) {
	m := nas.NewMessage()
	m.GsmMessage = nas.NewGsmMessage()
	m.GsmHeader.SetMessageType(nas.MsgTypePDUSessionReleaseReject)
	m.GsmHeader.SetExtendedProtocolDiscriminator(nasMessage.Epd5GSSessionManagementMessage)
	m.PDUSessionReleaseReject = nasMessage.NewPDUSessionReleaseReject(0x0)
	pDUSessionReleaseReject := m.PDUSessionReleaseReject

	pDUSessionReleaseReject.SetMessageType(nas.MsgTypePDUSessionReleaseReject)
	pDUSessionReleaseReject.SetExtendedProtocolDiscriminator(nasMessage.Epd5GSSessionManagementMessage)

	pDUSessionReleaseReject.SetPDUSessionID(uint8(smContext.PDUSessionID))

	pDUSessionReleaseReject.SetPTI(smContext.Pti)
	// TODO: fix to real value
	pDUSessionReleaseReject.SetCauseValue(nasMessage.Cause5GSMRequestRejectedUnspecified)

	return m.PlainNasEncode()
}

func BuildGSMPDUSessionModificationReject(smContext *SMContext) ([]byte, error) {
	m := nas.NewMessage()
	m.GsmMessage = nas.NewGsmMessage()
	m.GsmHeader.SetMessageType(nas.MsgTypePDUSessionModificationReject)
	m.GsmHeader.SetExtendedProtocolDiscriminator(nasMessage.Epd5GSSessionManagementMessage)
	m.PDUSessionModificationReject = nasMessage.NewPDUSessionModificationReject(0x0)
	pDUSessionModificationReject := m.PDUSessionModificationReject

	pDUSessionModificationReject.SetMessageType(nas.MsgTypePDUSessionModificationReject)
	pDUSessionModificationReject.SetExtendedProtocolDiscriminator(nasMessage.Epd5GSSessionManagementMessage)
	pDUSessionModificationReject.SetPDUSessionID(uint8(smContext.PDUSessionID))
	pDUSessionModificationReject.SetPTI(smContext.Pti)
	pDUSessionModificationReject.SetCauseValue(nasMessage.Cause5GSMMessageTypeNonExistentOrNotImplemented)

	return m.PlainNasEncode()
}
