//go:binary-only-package

package nasType

// SMSIndication 9.10.3.50A
// Iei Row, sBit, len = [0, 0], 8 , 4
// SAI Row, sBit, len = [0, 0], 1 , 1
type SMSIndication struct {
	Octet uint8
}

func NewSMSIndication(iei uint8) (sMSIndication *SMSIndication) {}

// SMSIndication 9.10.3.50A
// Iei Row, sBit, len = [0, 0], 8 , 4
func (a *SMSIndication) GetIei() (iei uint8) {}

// SMSIndication 9.10.3.50A
// Iei Row, sBit, len = [0, 0], 8 , 4
func (a *SMSIndication) SetIei(iei uint8) {}

// SMSIndication 9.10.3.50A
// SAI Row, sBit, len = [0, 0], 1 , 1
func (a *SMSIndication) GetSAI() (sAI uint8) {}

// SMSIndication 9.10.3.50A
// SAI Row, sBit, len = [0, 0], 1 , 1
func (a *SMSIndication) SetSAI(sAI uint8) {}
