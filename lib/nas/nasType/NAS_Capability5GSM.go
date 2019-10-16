//go:binary-only-package

package nasType

// Capability5GSM 9.11.4.1
// MH6PDU Row, sBit, len = [0, 0], 2 , 1
// RqoS Row, sBit, len = [0, 0], 1 , 1
// Spare Row, sBit, len = [1, 12], 8 , 96
type Capability5GSM struct {
	Iei   uint8
	Len   uint8
	Octet [13]uint8
}

func NewCapability5GSM(iei uint8) (capability5GSM *Capability5GSM) {}

// Capability5GSM 9.11.4.1
// Iei Row, sBit, len = [], 8, 8
func (a *Capability5GSM) GetIei() (iei uint8) {}

// Capability5GSM 9.11.4.1
// Iei Row, sBit, len = [], 8, 8
func (a *Capability5GSM) SetIei(iei uint8) {}

// Capability5GSM 9.11.4.1
// Len Row, sBit, len = [], 8, 8
func (a *Capability5GSM) GetLen() (len uint8) {}

// Capability5GSM 9.11.4.1
// Len Row, sBit, len = [], 8, 8
func (a *Capability5GSM) SetLen(len uint8) {}

// Capability5GSM 9.11.4.1
// MH6PDU Row, sBit, len = [0, 0], 2 , 1
func (a *Capability5GSM) GetMH6PDU() (mH6PDU uint8) {}

// Capability5GSM 9.11.4.1
// MH6PDU Row, sBit, len = [0, 0], 2 , 1
func (a *Capability5GSM) SetMH6PDU(mH6PDU uint8) {}

// Capability5GSM 9.11.4.1
// RqoS Row, sBit, len = [0, 0], 1 , 1
func (a *Capability5GSM) GetRqoS() (rqoS uint8) {}

// Capability5GSM 9.11.4.1
// RqoS Row, sBit, len = [0, 0], 1 , 1
func (a *Capability5GSM) SetRqoS(rqoS uint8) {}

// Capability5GSM 9.11.4.1
// Spare Row, sBit, len = [1, 12], 8 , 96
func (a *Capability5GSM) GetSpare() (spare [12]uint8) {}

// Capability5GSM 9.11.4.1
// Spare Row, sBit, len = [1, 12], 8 , 96
func (a *Capability5GSM) SetSpare(spare [12]uint8) {}
