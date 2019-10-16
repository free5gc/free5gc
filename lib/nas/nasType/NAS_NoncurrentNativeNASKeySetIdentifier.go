//go:binary-only-package

package nasType

// NoncurrentNativeNASKeySetIdentifier 9.11.3.32
// Iei  Row, sBit, len = [0, 0], 8 , 4
// Tsc  Row, sBit, len = [0, 0], 4 , 1
// NasKeySetIdentifiler  Row, sBit, len = [0, 0], 3 , 3
type NoncurrentNativeNASKeySetIdentifier struct {
	Octet uint8
}

func NewNoncurrentNativeNASKeySetIdentifier(iei uint8) (noncurrentNativeNASKeySetIdentifier *NoncurrentNativeNASKeySetIdentifier) {}

// NoncurrentNativeNASKeySetIdentifier 9.11.3.32
// Iei Row, sBit, len = [0, 0], 8 , 4
func (a *NoncurrentNativeNASKeySetIdentifier) GetIei() (iei uint8) {}

// NoncurrentNativeNASKeySetIdentifier 9.11.3.32
// Iei Row, sBit, len = [0, 0], 8 , 4
func (a *NoncurrentNativeNASKeySetIdentifier) SetIei(iei uint8) {}

// NoncurrentNativeNASKeySetIdentifier 9.11.3.32
// Tsc Row, sBit, len = [0, 0], 4 , 1
func (a *NoncurrentNativeNASKeySetIdentifier) GetTsc() (tsc uint8) {}

// NoncurrentNativeNASKeySetIdentifier 9.11.3.32
// Tsc Row, sBit, len = [0, 0], 4 , 1
func (a *NoncurrentNativeNASKeySetIdentifier) SetTsc(tsc uint8) {}

// NoncurrentNativeNASKeySetIdentifier 9.11.3.32
// NasKeySetIdentifiler Row, sBit, len = [0, 0], 3 , 3
func (a *NoncurrentNativeNASKeySetIdentifier) GetNasKeySetIdentifiler() (nasKeySetIdentifiler uint8) {}

// NoncurrentNativeNASKeySetIdentifier 9.11.3.32
// NasKeySetIdentifiler Row, sBit, len = [0, 0], 3 , 3
func (a *NoncurrentNativeNASKeySetIdentifier) SetNasKeySetIdentifiler(nasKeySetIdentifiler uint8) {}
