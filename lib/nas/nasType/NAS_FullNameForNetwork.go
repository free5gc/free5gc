//go:binary-only-package

package nasType

// FullNameForNetwork 9.11.3.35
// Ext Row, sBit, len = [0, 0], 8 ,1
// CodingScheme Row, sBit, len = [0, 0], 7 , 3
// AddCI Row, sBit, len = [0, 0], 4 , 1
// NumberOfSpareBitsInLastOctet Row, sBit, len = [0, 0], 3 , 3
// TextString Row, sBit, len = [1, 1], 4 , INF
type FullNameForNetwork struct {
	Iei    uint8
	Len    uint8
	Buffer []uint8
}

func NewFullNameForNetwork(iei uint8) (fullNameForNetwork *FullNameForNetwork) {}

// FullNameForNetwork 9.11.3.35
// Iei Row, sBit, len = [], 8, 8
func (a *FullNameForNetwork) GetIei() (iei uint8) {}

// FullNameForNetwork 9.11.3.35
// Iei Row, sBit, len = [], 8, 8
func (a *FullNameForNetwork) SetIei(iei uint8) {}

// FullNameForNetwork 9.11.3.35
// Len Row, sBit, len = [], 8, 8
func (a *FullNameForNetwork) GetLen() (len uint8) {}

// FullNameForNetwork 9.11.3.35
// Len Row, sBit, len = [], 8, 8
func (a *FullNameForNetwork) SetLen(len uint8) {}

// FullNameForNetwork 9.11.3.35
// Ext Row, sBit, len = [0, 0], 8 ,1
func (a *FullNameForNetwork) GetExt() (ext uint8) {}

// FullNameForNetwork 9.11.3.35
// Ext Row, sBit, len = [0, 0], 8 ,1
func (a *FullNameForNetwork) SetExt(ext uint8) {}

// FullNameForNetwork 9.11.3.35
// CodingScheme Row, sBit, len = [0, 0], 7 , 3
func (a *FullNameForNetwork) GetCodingScheme() (codingScheme uint8) {}

// FullNameForNetwork 9.11.3.35
// CodingScheme Row, sBit, len = [0, 0], 7 , 3
func (a *FullNameForNetwork) SetCodingScheme(codingScheme uint8) {}

// FullNameForNetwork 9.11.3.35
// AddCI Row, sBit, len = [0, 0], 4 , 1
func (a *FullNameForNetwork) GetAddCI() (addCI uint8) {}

// FullNameForNetwork 9.11.3.35
// AddCI Row, sBit, len = [0, 0], 4 , 1
func (a *FullNameForNetwork) SetAddCI(addCI uint8) {}

// FullNameForNetwork 9.11.3.35
// NumberOfSpareBitsInLastOctet Row, sBit, len = [0, 0], 3 , 3
func (a *FullNameForNetwork) GetNumberOfSpareBitsInLastOctet() (numberOfSpareBitsInLastOctet uint8) {}

// FullNameForNetwork 9.11.3.35
// NumberOfSpareBitsInLastOctet Row, sBit, len = [0, 0], 3 , 3
func (a *FullNameForNetwork) SetNumberOfSpareBitsInLastOctet(numberOfSpareBitsInLastOctet uint8) {}

// FullNameForNetwork 9.11.3.35
// TextString Row, sBit, len = [1, 1], 4 , INF
func (a *FullNameForNetwork) GetTextString() (textString []uint8) {}

// FullNameForNetwork 9.11.3.35
// TextString Row, sBit, len = [1, 1], 4 , INF
func (a *FullNameForNetwork) SetTextString(textString []uint8) {}
