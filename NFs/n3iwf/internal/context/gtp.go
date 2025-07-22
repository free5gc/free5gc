package context

type GtpEventType int64

// GTP Event Type
const (
	ForwardUL GtpEventType = iota
)

type GtpEvt interface {
	Type() GtpEventType
}

type ForwardULEvt struct {
	GtpConnInfo *GTPConnectionInfo
	QFI         *uint8
	Payload     []byte
}

func (forwardDLEvt *ForwardULEvt) Type() GtpEventType {
	return ForwardUL
}

func NewForwardULEvt(gtpConnInfo *GTPConnectionInfo, qfi *uint8, payload []byte) *ForwardULEvt {
	return &ForwardULEvt{
		GtpConnInfo: gtpConnInfo,
		QFI:         qfi,
		Payload:     payload,
	}
}
