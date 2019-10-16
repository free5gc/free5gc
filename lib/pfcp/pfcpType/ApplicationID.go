//go:binary-only-package

package pfcpType

type ApplicationID struct {
	ApplicationIdentifier []byte
}

func (a *ApplicationID) MarshalBinary() (data []byte, err error) {}

func (a *ApplicationID) UnmarshalBinary(data []byte) error {}
