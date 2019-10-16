//go:binary-only-package

package nasType

// RequestType 9.11.3.47
// Iei Row, sBit, len = [0, 0], 8 , 4
// RequestTypeValue Row, sBit, len = [0, 0], 3 , 3
type RequestType struct {
	Octet uint8
}

func NewRequestType(iei uint8) (requestType *RequestType) {}

// RequestType 9.11.3.47
// Iei Row, sBit, len = [0, 0], 8 , 4
func (a *RequestType) GetIei() (iei uint8) {}

// RequestType 9.11.3.47
// Iei Row, sBit, len = [0, 0], 8 , 4
func (a *RequestType) SetIei(iei uint8) {}

// RequestType 9.11.3.47
// RequestTypeValue Row, sBit, len = [0, 0], 3 , 3
func (a *RequestType) GetRequestTypeValue() (requestTypeValue uint8) {}

// RequestType 9.11.3.47
// RequestTypeValue Row, sBit, len = [0, 0], 3 , 3
func (a *RequestType) SetRequestTypeValue(requestTypeValue uint8) {}
