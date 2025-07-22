package message

import (
	"encoding/hex"

	"github.com/pkg/errors"
	gtpMsg "github.com/wmnsk/go-gtp/gtpv1/message"

	"github.com/free5gc/n3iwf/internal/logger"
)

// [TS 38.415] 5.5.2 Frame format for the PDU Session user plane protocol
const (
	DL_PDU_SESSION_INFORMATION_TYPE = 0x00
	UL_PDU_SESSION_INFORMATION_TYPE = 0x10
)

type QoSTPDUPacket struct {
	tPDU *gtpMsg.TPDU
	qos  bool
	rqi  bool
	qfi  uint8
}

func (p *QoSTPDUPacket) GetPayload() []byte {
	return p.tPDU.Payload
}

func (p *QoSTPDUPacket) GetTEID() uint32 {
	return p.tPDU.TEID()
}

func (p *QoSTPDUPacket) GetExtensionHeader() []*gtpMsg.ExtensionHeader {
	return p.tPDU.ExtensionHeaders
}

func (p *QoSTPDUPacket) HasQoS() bool {
	return p.qos
}

func (p *QoSTPDUPacket) GetQoSParameters() (uint8, bool) {
	return p.qfi, p.rqi
}

func (p *QoSTPDUPacket) Unmarshal(pdu *gtpMsg.TPDU) error {
	p.tPDU = pdu
	if p.tPDU.HasExtensionHeader() {
		if err := p.unmarshalExtensionHeader(); err != nil {
			return err
		}
	}

	return nil
}

// [TS 29.281] [TS 38.415]
// Define GTP extension header
// [TS 38.415]
// Define PDU Session User Plane protocol
func (p *QoSTPDUPacket) unmarshalExtensionHeader() error {
	gtpLog := logger.GTPLog

	for _, eh := range p.tPDU.ExtensionHeaders {
		switch eh.Type {
		case gtpMsg.ExtHeaderTypePDUSessionContainer:
			p.qos = true
			p.rqi = ((int(eh.Content[1]) >> 6) & 0x1) == 1
			p.qfi = eh.Content[1] & 0x3F
			gtpLog.Tracef("Parsed Extension Header: Len=%d, Next Type=%d, Content Dump:\n%s",
				eh.Length, eh.NextType, hex.Dump(eh.Content))
		default:
			gtpLog.Warningf("Unsupported Extension Header Field Value: %x", eh.Type)
		}
	}

	if !p.qos {
		return errors.Errorf("unmarshalExtensionHeader err: no PDUSessionContainer in ExtensionHeaders.")
	}

	return nil
}

func BuildQoSGTPPacket(teid uint32, qfi uint8, payload []byte) ([]byte, error) {
	header := gtpMsg.NewHeader(0x34, gtpMsg.MsgTypeTPDU, teid, 0x00, payload).WithExtensionHeaders(
		gtpMsg.NewExtensionHeader(
			gtpMsg.ExtHeaderTypePDUSessionContainer,
			[]byte{UL_PDU_SESSION_INFORMATION_TYPE, qfi},
			gtpMsg.ExtHeaderTypeNoMoreExtensionHeaders,
		),
	)

	b := make([]byte, header.MarshalLen())
	if err := header.MarshalTo(b); err != nil {
		return nil, errors.Wrapf(err, "go-gtp Marshal failed")
	}

	return b, nil
}
