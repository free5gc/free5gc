//go:binary-only-package

package tlv

type fragments map[int][][]byte

func (f fragments) Add(tag int, buf []byte) {}

func (f fragments) Get(tag int) ([][]byte, bool) {}
