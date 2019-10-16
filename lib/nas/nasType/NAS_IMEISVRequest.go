//go:binary-only-package

package nasType

// IMEISVRequest 9.11.3.28
// Iei Row, sBit, len = [0, 0], 8 , 4
// IMEISVRequestValue Row, sBit, len = [0, 0], 3 , 3
type IMEISVRequest struct {
	Octet uint8
}

func NewIMEISVRequest(iei uint8) (iMEISVRequest *IMEISVRequest) {}

// IMEISVRequest 9.11.3.28
// Iei Row, sBit, len = [0, 0], 8 , 4
func (a *IMEISVRequest) GetIei() (iei uint8) {}

// IMEISVRequest 9.11.3.28
// Iei Row, sBit, len = [0, 0], 8 , 4
func (a *IMEISVRequest) SetIei(iei uint8) {}

// IMEISVRequest 9.11.3.28
// IMEISVRequestValue Row, sBit, len = [0, 0], 3 , 3
func (a *IMEISVRequest) GetIMEISVRequestValue() (iMEISVRequestValue uint8) {}

// IMEISVRequest 9.11.3.28
// IMEISVRequestValue Row, sBit, len = [0, 0], 3 , 3
func (a *IMEISVRequest) SetIMEISVRequestValue(iMEISVRequestValue uint8) {}
