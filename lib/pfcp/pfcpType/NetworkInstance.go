//go:binary-only-package

package pfcpType

type NetworkInstance struct {
	NetworkInstance []byte
}

func (n *NetworkInstance) MarshalBinary() (data []byte, err error) {}

func (n *NetworkInstance) UnmarshalBinary(data []byte) error {}
