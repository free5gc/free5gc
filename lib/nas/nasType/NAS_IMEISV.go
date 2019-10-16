//go:binary-only-package

package nasType

// IMEISV 9.11.3.4
// IdentityDigit1 Row, sBit, len = [0, 0], 8 , 4
// OddEvenIdic Row, sBit, len = [0, 0], 4 , 1
// TypeOfIdentity Row, sBit, len = [0, 0], 3 , 3
// IdentityDigitP_1 Row, sBit, len = [1, 1], 8 , 4
// IdentityDigitP Row, sBit, len = [1, 1], 4 , 4
// IdentityDigitP_3 Row, sBit, len = [2, 2], 8 , 4
// IdentityDigitP_2 Row, sBit, len = [2, 2], 4 , 4
// IdentityDigitP_5 Row, sBit, len = [3, 3], 8 , 4
// IdentityDigitP_4 Row, sBit, len = [3, 3], 4 , 4
// IdentityDigitP_7 Row, sBit, len = [4, 4], 8 , 4
// IdentityDigitP_6 Row, sBit, len = [4, 4], 4 , 4
// IdentityDigitP_9 Row, sBit, len = [5, 5], 8 , 4
// IdentityDigitP_8 Row, sBit, len = [5, 5], 4 , 4
// IdentityDigitP_11 Row, sBit, len = [6, 6], 8 , 4
// IdentityDigitP_10 Row, sBit, len = [6, 6], 4 , 4
// IdentityDigitP_13 Row, sBit, len = [7, 7], 8 , 4
// IdentityDigitP_12 Row, sBit, len = [7, 7], 4 , 4
// IdentityDigitP_15 Row, sBit, len = [8, 8], 8 , 4
// IdentityDigitP_14 Row, sBit, len = [8, 8], 4 , 4
type IMEISV struct {
	Iei   uint8
	Len   uint16
	Octet [9]uint8
}

func NewIMEISV(iei uint8) (iMEISV *IMEISV) {}

// IMEISV 9.11.3.4
// Iei Row, sBit, len = [], 8, 8
func (a *IMEISV) GetIei() (iei uint8) {}

// IMEISV 9.11.3.4
// Iei Row, sBit, len = [], 8, 8
func (a *IMEISV) SetIei(iei uint8) {}

// IMEISV 9.11.3.4
// Len Row, sBit, len = [], 8, 16
func (a *IMEISV) GetLen() (len uint16) {}

// IMEISV 9.11.3.4
// Len Row, sBit, len = [], 8, 16
func (a *IMEISV) SetLen(len uint16) {}

// IMEISV 9.11.3.4
// IdentityDigit1 Row, sBit, len = [0, 0], 8 , 4
func (a *IMEISV) GetIdentityDigit1() (identityDigit1 uint8) {}

// IMEISV 9.11.3.4
// IdentityDigit1 Row, sBit, len = [0, 0], 8 , 4
func (a *IMEISV) SetIdentityDigit1(identityDigit1 uint8) {}

// IMEISV 9.11.3.4
// OddEvenIdic Row, sBit, len = [0, 0], 4 , 1
func (a *IMEISV) GetOddEvenIdic() (oddEvenIdic uint8) {}

// IMEISV 9.11.3.4
// OddEvenIdic Row, sBit, len = [0, 0], 4 , 1
func (a *IMEISV) SetOddEvenIdic(oddEvenIdic uint8) {}

// IMEISV 9.11.3.4
// TypeOfIdentity Row, sBit, len = [0, 0], 3 , 3
func (a *IMEISV) GetTypeOfIdentity() (typeOfIdentity uint8) {}

// IMEISV 9.11.3.4
// TypeOfIdentity Row, sBit, len = [0, 0], 3 , 3
func (a *IMEISV) SetTypeOfIdentity(typeOfIdentity uint8) {}

// IMEISV 9.11.3.4
// IdentityDigitP_1 Row, sBit, len = [1, 1], 8 , 4
func (a *IMEISV) GetIdentityDigitP_1() (identityDigitP_1 uint8) {}

// IMEISV 9.11.3.4
// IdentityDigitP_1 Row, sBit, len = [1, 1], 8 , 4
func (a *IMEISV) SetIdentityDigitP_1(identityDigitP_1 uint8) {}

// IMEISV 9.11.3.4
// IdentityDigitP Row, sBit, len = [1, 1], 4 , 4
func (a *IMEISV) GetIdentityDigitP() (identityDigitP uint8) {}

// IMEISV 9.11.3.4
// IdentityDigitP Row, sBit, len = [1, 1], 4 , 4
func (a *IMEISV) SetIdentityDigitP(identityDigitP uint8) {}

// IMEISV 9.11.3.4
// IdentityDigitP_3 Row, sBit, len = [2, 2], 8 , 4
func (a *IMEISV) GetIdentityDigitP_3() (identityDigitP_3 uint8) {}

// IMEISV 9.11.3.4
// IdentityDigitP_3 Row, sBit, len = [2, 2], 8 , 4
func (a *IMEISV) SetIdentityDigitP_3(identityDigitP_3 uint8) {}

// IMEISV 9.11.3.4
// IdentityDigitP_2 Row, sBit, len = [2, 2], 4 , 4
func (a *IMEISV) GetIdentityDigitP_2() (identityDigitP_2 uint8) {}

// IMEISV 9.11.3.4
// IdentityDigitP_2 Row, sBit, len = [2, 2], 4 , 4
func (a *IMEISV) SetIdentityDigitP_2(identityDigitP_2 uint8) {}

// IMEISV 9.11.3.4
// IdentityDigitP_5 Row, sBit, len = [3, 3], 8 , 4
func (a *IMEISV) GetIdentityDigitP_5() (identityDigitP_5 uint8) {}

// IMEISV 9.11.3.4
// IdentityDigitP_5 Row, sBit, len = [3, 3], 8 , 4
func (a *IMEISV) SetIdentityDigitP_5(identityDigitP_5 uint8) {}

// IMEISV 9.11.3.4
// IdentityDigitP_4 Row, sBit, len = [3, 3], 4 , 4
func (a *IMEISV) GetIdentityDigitP_4() (identityDigitP_4 uint8) {}

// IMEISV 9.11.3.4
// IdentityDigitP_4 Row, sBit, len = [3, 3], 4 , 4
func (a *IMEISV) SetIdentityDigitP_4(identityDigitP_4 uint8) {}

// IMEISV 9.11.3.4
// IdentityDigitP_7 Row, sBit, len = [4, 4], 8 , 4
func (a *IMEISV) GetIdentityDigitP_7() (identityDigitP_7 uint8) {}

// IMEISV 9.11.3.4
// IdentityDigitP_7 Row, sBit, len = [4, 4], 8 , 4
func (a *IMEISV) SetIdentityDigitP_7(identityDigitP_7 uint8) {}

// IMEISV 9.11.3.4
// IdentityDigitP_6 Row, sBit, len = [4, 4], 4 , 4
func (a *IMEISV) GetIdentityDigitP_6() (identityDigitP_6 uint8) {}

// IMEISV 9.11.3.4
// IdentityDigitP_6 Row, sBit, len = [4, 4], 4 , 4
func (a *IMEISV) SetIdentityDigitP_6(identityDigitP_6 uint8) {}

// IMEISV 9.11.3.4
// IdentityDigitP_9 Row, sBit, len = [5, 5], 8 , 4
func (a *IMEISV) GetIdentityDigitP_9() (identityDigitP_9 uint8) {}

// IMEISV 9.11.3.4
// IdentityDigitP_9 Row, sBit, len = [5, 5], 8 , 4
func (a *IMEISV) SetIdentityDigitP_9(identityDigitP_9 uint8) {}

// IMEISV 9.11.3.4
// IdentityDigitP_8 Row, sBit, len = [5, 5], 4 , 4
func (a *IMEISV) GetIdentityDigitP_8() (identityDigitP_8 uint8) {}

// IMEISV 9.11.3.4
// IdentityDigitP_8 Row, sBit, len = [5, 5], 4 , 4
func (a *IMEISV) SetIdentityDigitP_8(identityDigitP_8 uint8) {}

// IMEISV 9.11.3.4
// IdentityDigitP_11 Row, sBit, len = [6, 6], 8 , 4
func (a *IMEISV) GetIdentityDigitP_11() (identityDigitP_11 uint8) {}

// IMEISV 9.11.3.4
// IdentityDigitP_11 Row, sBit, len = [6, 6], 8 , 4
func (a *IMEISV) SetIdentityDigitP_11(identityDigitP_11 uint8) {}

// IMEISV 9.11.3.4
// IdentityDigitP_10 Row, sBit, len = [6, 6], 4 , 4
func (a *IMEISV) GetIdentityDigitP_10() (identityDigitP_10 uint8) {}

// IMEISV 9.11.3.4
// IdentityDigitP_10 Row, sBit, len = [6, 6], 4 , 4
func (a *IMEISV) SetIdentityDigitP_10(identityDigitP_10 uint8) {}

// IMEISV 9.11.3.4
// IdentityDigitP_13 Row, sBit, len = [7, 7], 8 , 4
func (a *IMEISV) GetIdentityDigitP_13() (identityDigitP_13 uint8) {}

// IMEISV 9.11.3.4
// IdentityDigitP_13 Row, sBit, len = [7, 7], 8 , 4
func (a *IMEISV) SetIdentityDigitP_13(identityDigitP_13 uint8) {}

// IMEISV 9.11.3.4
// IdentityDigitP_12 Row, sBit, len = [7, 7], 4 , 4
func (a *IMEISV) GetIdentityDigitP_12() (identityDigitP_12 uint8) {}

// IMEISV 9.11.3.4
// IdentityDigitP_12 Row, sBit, len = [7, 7], 4 , 4
func (a *IMEISV) SetIdentityDigitP_12(identityDigitP_12 uint8) {}

// IMEISV 9.11.3.4
// IdentityDigitP_15 Row, sBit, len = [8, 8], 8 , 4
func (a *IMEISV) GetIdentityDigitP_15() (identityDigitP_15 uint8) {}

// IMEISV 9.11.3.4
// IdentityDigitP_15 Row, sBit, len = [8, 8], 8 , 4
func (a *IMEISV) SetIdentityDigitP_15(identityDigitP_15 uint8) {}

// IMEISV 9.11.3.4
// IdentityDigitP_14 Row, sBit, len = [8, 8], 4 , 4
func (a *IMEISV) GetIdentityDigitP_14() (identityDigitP_14 uint8) {}

// IMEISV 9.11.3.4
// IdentityDigitP_14 Row, sBit, len = [8, 8], 4 , 4
func (a *IMEISV) SetIdentityDigitP_14(identityDigitP_14 uint8) {}
