package smf_context

import (
	"free5gc/lib/nas"
	"free5gc/lib/nas/nasMessage"
	// "free5gc/lib/nas/nasType"
)

func BuildGSMPDUSessionEstablishmentAccept(smContext *SMContext) ([]byte, error) {
	m := nas.NewMessage()
	m.GsmMessage = nas.NewGsmMessage()
	m.GsmHeader.SetMessageType(nas.MsgTypePDUSessionEstablishmentAccept)
	m.GsmHeader.SetExtendedProtocolDiscriminator(nasMessage.Epd5GSSessionManagementMessage)
	m.PDUSessionEstablishmentAccept = nasMessage.NewPDUSessionEstablishmentAccept(0x0)
	pDUSessionEstablishmentAccept := m.PDUSessionEstablishmentAccept

	pDUSessionEstablishmentAccept.SetPDUSessionID(uint8(10))
	pDUSessionEstablishmentAccept.SetMessageType(nas.MsgTypePDUSessionEstablishmentAccept)
	pDUSessionEstablishmentAccept.SetExtendedProtocolDiscriminator(nasMessage.Epd5GSSessionManagementMessage)
	pDUSessionEstablishmentAccept.SetPTI(0x00)
	pDUSessionEstablishmentAccept.SetPDUSessionType(1)
	pDUSessionEstablishmentAccept.SetSSCMode(1)
	// pDUSessionEstablishmentAccept.SetQosRule()
	// pDUSessionEstablishmentAccept.AuthorizedQosRules.SetLen()
	pDUSessionEstablishmentAccept.SessionAMBR.SetSessionAMBRForDownlink([2]uint8{0x11, 0x11})
	pDUSessionEstablishmentAccept.SessionAMBR.SetSessionAMBRForUplink([2]uint8{0x11, 0x11})
	pDUSessionEstablishmentAccept.SessionAMBR.SetUnitForSessionAMBRForDownlink(10)
	pDUSessionEstablishmentAccept.SessionAMBR.SetUnitForSessionAMBRForUplink(10)
	pDUSessionEstablishmentAccept.SessionAMBR.SetLen(uint8(len(pDUSessionEstablishmentAccept.SessionAMBR.Octet)))

	// pDUSessionEstablishmentAccept.AuthorizedQosFlowDescriptions = nasType.NewAuthorizedQosFlowDescriptions(nasMessage.PDUSessionEstablishmentAcceptAuthorizedQosFlowDescriptionsType)
	// pDUSessionEstablishmentAccept.AuthorizedQosFlowDescriptions.SetLen(6)
	// pDUSessionEstablishmentAccept.SetQoSFlowDescriptions([]uint8{0x09, 0x20, 0x41, 0x01, 0x01, 0x09})
	return m.PlainNasEncode()
}

func BuildGSMPDUSessionReleaseCommand(smContext *SMContext) ([]byte, error) {

	m := nas.NewMessage()
	m.GsmMessage = nas.NewGsmMessage()
	m.GsmHeader.SetMessageType(nas.MsgTypePDUSessionReleaseCommand)
	m.GsmHeader.SetExtendedProtocolDiscriminator(nasMessage.Epd5GSSessionManagementMessage)
	m.PDUSessionReleaseCommand = nasMessage.NewPDUSessionReleaseCommand(0x0)
	pDUSessionReleaseCommand := m.PDUSessionReleaseCommand

	pDUSessionReleaseCommand.SetMessageType(nas.MsgTypePDUSessionReleaseCommand)
	pDUSessionReleaseCommand.SetExtendedProtocolDiscriminator(nasMessage.Epd5GSSessionManagementMessage)
	pDUSessionReleaseCommand.SetPDUSessionID(uint8(smContext.PDUSessionID))
	pDUSessionReleaseCommand.SetPTI(0x00)
	pDUSessionReleaseCommand.SetCauseValue(0x0)

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
	pDUSessionModificationCommand.SetPTI(0x00)
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
