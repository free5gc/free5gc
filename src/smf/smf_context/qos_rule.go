package smf_context

import (
	"bytes"
)

const (
	OperationCodeCreateNewQoSRule                                   uint8 = 1
	OperationCodeDeleteExistingQoSRule                              uint8 = 2
	OperationCodeModifyExistingQoSRuleAndAddPacketFilters           uint8 = 3
	OperationCodeModifyExistingQoSRuleAndReplaceAllPacketFilters    uint8 = 4
	OperationCodeModifyExistingQoSRuleAndDeletePacketFilters        uint8 = 5
	OperationCodeModifyExistingQoSRuleWithoutModifyingPacketFilters uint8 = 6
)

const (
	PacketFilterDirectionDownlink      uint8 = 1
	PacketFilterDirectionUplink        uint8 = 2
	PacketFilterDirectionBidirectional uint8 = 3
)

// TS 24.501 Table 9.11.4.13.1
const (
	PacketFilterComponentTypeMatchAll                       uint8 = 0x01
	PacketFilterComponentTypeIPv4RemoteAddress              uint8 = 0x10
	PacketFilterComponentTypeIPv4LocalAddress               uint8 = 0x11
	PacketFilterComponentTypeIPv6RemoteAddress              uint8 = 0x21
	PacketFilterComponentTypeIPv6LocalAddress               uint8 = 0x23
	PacketFilterComponentTypeProtocolIdentifierOrNextHeader uint8 = 0x30
	PacketFilterComponentTypeSingleLocalPort                uint8 = 0x40
	PacketFilterComponentTypeLocalPortRange                 uint8 = 0x41
	PacketFilterComponentTypeSingleRemotePort               uint8 = 0x50
	PacketFilterComponentTypeRemotePortRange                uint8 = 0x51
	PacketFilterComponentTypeSecurityParameterIndex         uint8 = 0x60
	PacketFilterComponentTypeTypeOfServiceOrTrafficClass    uint8 = 0x70
	PacketFilterComponentTypeFlowLabel                      uint8 = 0x80
	PacketFilterComponentTypeDestinationMACAddress          uint8 = 0x81
	PacketFilterComponentTypeSourceMACAddress               uint8 = 0x82
	PacketFilterComponentType8021Q_CTAG_VID                 uint8 = 0x83
	PacketFilterComponentType8021Q_STAG_VID                 uint8 = 0x84
	PacketFilterComponentType8021Q_CTAG_PCPOrDEI            uint8 = 0x85
	PacketFilterComponentType8021Q_STAG_PCPOrDEI            uint8 = 0x86
	PacketFilterComponentTypeEthertype                      uint8 = 0x87
)

type PacketFilter struct {
	Direction     uint8
	Identifier    uint8
	ComponentType uint8
	Component     []byte
}

func (pf *PacketFilter) MarshalBinary() (data []byte, err error) {
	packetFilterBuffer := bytes.NewBuffer(nil)
	header := 0 | pf.Direction<<4 | pf.Identifier
	// write header
	err = packetFilterBuffer.WriteByte(header)
	if err != nil {
		return nil, err
	}
	// write length of packet filter
	err = packetFilterBuffer.WriteByte(uint8(1 + len(pf.Component)))
	if err != nil {
		return nil, err
	}

	err = packetFilterBuffer.WriteByte(pf.ComponentType)
	if err != nil {
		return nil, err
	}

	if pf.ComponentType == PacketFilterComponentTypeMatchAll || pf.Component == nil {
		_, err = packetFilterBuffer.Write(pf.Component)
		if err != nil {
			return nil, err
		}
	}

	return packetFilterBuffer.Bytes(), nil
}

type QoSRule struct {
	Identifier       uint8
	OperationCode    uint8
	DQR              uint8
	Segregation      uint8
	PacketFilterList []PacketFilter
	Precedence       uint8
	QFI              uint8
}

func (r *QoSRule) MarshalBinary() (data []byte, err error) {
	ruleContentBuffer := bytes.NewBuffer(nil)

	// write rule content Header
	ruleContentHeader := r.OperationCode<<5 | r.DQR<<4 | uint8(len(r.PacketFilterList))
	ruleContentBuffer.WriteByte(ruleContentHeader)
	if err != nil {
		return nil, err
	}

	packetFilterListBuffer := &bytes.Buffer{}
	for _, pf := range r.PacketFilterList {
		packetFilterBuffer, err := pf.MarshalBinary()
		if err != nil {
			return nil, err
		}
		_, err = packetFilterListBuffer.Write(packetFilterBuffer)
		if err != nil {
			return nil, err
		}
	}

	// write QoS
	_, err = ruleContentBuffer.ReadFrom(packetFilterListBuffer)
	if err != nil {
		return nil, err
	}

	// write precedence
	err = ruleContentBuffer.WriteByte(r.Precedence)
	if err != nil {
		return nil, err
	}

	// write Segregation and QFI
	segregationAndQFIByte := r.Segregation<<6 | r.QFI
	err = ruleContentBuffer.WriteByte(segregationAndQFIByte)
	if err != nil {
		return nil, err
	}

	ruleBuffer := bytes.NewBuffer(nil)
	// write QoS rule identifier
	err = ruleBuffer.WriteByte(r.Identifier)
	if err != nil {
		return nil, err
	}

	// write QoS rule length
	err = ruleBuffer.WriteByte(uint8(ruleContentBuffer.Len()))
	if err != nil {
		return nil, err
	}

	// write QoS rule Content
	_, err = ruleBuffer.ReadFrom(ruleContentBuffer)
	if err != nil {
		return nil, err
	}

	return ruleBuffer.Bytes(), nil
}

type QoSRules []QoSRule

func (rs QoSRules) MarshalBinary() (data []byte, err error) {
	qosRulesBuffer := bytes.NewBuffer(nil)

	for _, rule := range rs {
		ruleBytes, err := rule.MarshalBinary()
		if err != nil {
			return nil, err
		}
		_, err = qosRulesBuffer.Write(ruleBytes)
		if err != nil {
			return nil, err
		}
	}
	return qosRulesBuffer.Bytes(), nil
}
