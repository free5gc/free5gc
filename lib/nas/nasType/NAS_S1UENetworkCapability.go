//go:binary-only-package

package nasType

// S1UENetworkCapability 9.11.3.48
// EEA0 Row, sBit, len = [0, 0], 8 , 1
// EEA1_128 Row, sBit, len = [0, 0], 7 , 1
// EEA2_128 Row, sBit, len = [0, 0], 6 , 1
// EEA3_128 Row, sBit, len = [0, 0], 5 , 1
// EEA4 Row, sBit, len = [0, 0], 4 , 1
// EEA5 Row, sBit, len = [0, 0], 3 , 1
// EEA6 Row, sBit, len = [0, 0], 2 , 1
// EEA7 Row, sBit, len = [0, 0], 1 , 1
// EIA0 Row, sBit, len = [1, 1], 8 , 1
// EIA1_128 Row, sBit, len = [1, 1], 7 , 1
// EIA2_128 Row, sBit, len = [1, 1], 6 , 1
// EIA3_128 Row, sBit, len = [1, 1], 5 , 1
// EIA4 Row, sBit, len = [1, 1], 4 , 1
// EIA5 Row, sBit, len = [1, 1], 3 , 1
// EIA6 Row, sBit, len = [1, 1], 2 , 1
// EIA7 Row, sBit, len = [1, 1], 1 , 1
// UEA0 Row, sBit, len = [2, 2], 8 , 1
// UEA1 Row, sBit, len = [2, 2], 7 , 1
// UEA2 Row, sBit, len = [2, 2], 6 , 1
// UEA3 Row, sBit, len = [2, 2], 5 , 1
// UEA4 Row, sBit, len = [2, 2], 4 , 1
// UEA5 Row, sBit, len = [2, 2], 3 , 1
// UEA6 Row, sBit, len = [2, 2], 2 , 1
// UEA7 Row, sBit, len = [2, 2], 1 , 1
// UCS2 Row, sBit, len = [3, 3], 8 , 1
// UIA1 Row, sBit, len = [3, 3], 7 , 1
// UIA2 Row, sBit, len = [3, 3], 6 , 1
// UIA3 Row, sBit, len = [3, 3], 5 , 1
// UIA4 Row, sBit, len = [3, 3], 4 , 1
// UIA5 Row, sBit, len = [3, 3], 3 , 1
// UIA6 Row, sBit, len = [3, 3], 2 , 1
// UIA7 Row, sBit, len = [3, 3], 1 , 1
// ProSedd Row, sBit, len = [4, 4], 8 , 1
// ProSe Row, sBit, len = [4, 4], 7 , 1
// H245ASH Row, sBit, len = [4, 4], 6 , 1
// ACCCSFB Row, sBit, len = [4, 4], 5 , 1
// LPP Row, sBit, len = [4, 4], 4 , 1
// LCS Row, sBit, len = [4, 4], 3 , 1
// xSRVCC Row, sBit, len = [4, 4], 2 , 1
// NF Row, sBit, len = [4, 4], 1 , 1
// EPCO Row, sBit, len = [5, 5], 8 , 1
// HCCPCIOT Row, sBit, len = [5, 5], 7 , 1
// ERwoPDN Row, sBit, len = [5, 5], 6 , 1
// S1UData Row, sBit, len = [5, 5], 5 , 1
// UPCIot Row, sBit, len = [5, 5], 4 , 1
// CPCIot Row, sBit, len = [5, 5], 3 , 1
// Proserelay Row, sBit, len = [5, 5], 2 , 1
// ProSedc Row, sBit, len = [5, 5], 1 , 1
// Bearer15 Row, sBit, len = [6, 6], 8 , 1
// SGC Row, sBit, len = [6, 6], 7 , 1
// N1mode Row, sBit, len = [6, 6], 6 , 1
// DCNR Row, sBit, len = [6, 6], 5 , 1
// CPbackoff Row, sBit, len = [6, 6], 4 , 1
// RestrictEC Row, sBit, len = [6, 6], 3 , 1
// V2XPC5 Row, sBit, len = [6, 6], 2 , 1
// MulitpeDRB Row, sBit, len = [6, 6], 1 , 1
// Spare Row, sBit, len = [7, 12], 8 , INF
type S1UENetworkCapability struct {
	Iei    uint8
	Len    uint8
	Buffer []uint8
}

func NewS1UENetworkCapability(iei uint8) (s1UENetworkCapability *S1UENetworkCapability) {}

// S1UENetworkCapability 9.11.3.48
// Iei Row, sBit, len = [], 8, 8
func (a *S1UENetworkCapability) GetIei() (iei uint8) {}

// S1UENetworkCapability 9.11.3.48
// Iei Row, sBit, len = [], 8, 8
func (a *S1UENetworkCapability) SetIei(iei uint8) {}

// S1UENetworkCapability 9.11.3.48
// Len Row, sBit, len = [], 8, 8
func (a *S1UENetworkCapability) GetLen() (len uint8) {}

// S1UENetworkCapability 9.11.3.48
// Len Row, sBit, len = [], 8, 8
func (a *S1UENetworkCapability) SetLen(len uint8) {}

// S1UENetworkCapability 9.11.3.48
// EEA0 Row, sBit, len = [0, 0], 8 , 1
func (a *S1UENetworkCapability) GetEEA0() (eEA0 uint8) {}

// S1UENetworkCapability 9.11.3.48
// EEA0 Row, sBit, len = [0, 0], 8 , 1
func (a *S1UENetworkCapability) SetEEA0(eEA0 uint8) {}

// S1UENetworkCapability 9.11.3.48
// EEA1_128 Row, sBit, len = [0, 0], 7 , 1
func (a *S1UENetworkCapability) GetEEA1_128() (eEA1_128 uint8) {}

// S1UENetworkCapability 9.11.3.48
// EEA1_128 Row, sBit, len = [0, 0], 7 , 1
func (a *S1UENetworkCapability) SetEEA1_128(eEA1_128 uint8) {}

// S1UENetworkCapability 9.11.3.48
// EEA2_128 Row, sBit, len = [0, 0], 6 , 1
func (a *S1UENetworkCapability) GetEEA2_128() (eEA2_128 uint8) {}

// S1UENetworkCapability 9.11.3.48
// EEA2_128 Row, sBit, len = [0, 0], 6 , 1
func (a *S1UENetworkCapability) SetEEA2_128(eEA2_128 uint8) {}

// S1UENetworkCapability 9.11.3.48
// EEA3_128 Row, sBit, len = [0, 0], 5 , 1
func (a *S1UENetworkCapability) GetEEA3_128() (eEA3_128 uint8) {}

// S1UENetworkCapability 9.11.3.48
// EEA3_128 Row, sBit, len = [0, 0], 5 , 1
func (a *S1UENetworkCapability) SetEEA3_128(eEA3_128 uint8) {}

// S1UENetworkCapability 9.11.3.48
// EEA4 Row, sBit, len = [0, 0], 4 , 1
func (a *S1UENetworkCapability) GetEEA4() (eEA4 uint8) {}

// S1UENetworkCapability 9.11.3.48
// EEA4 Row, sBit, len = [0, 0], 4 , 1
func (a *S1UENetworkCapability) SetEEA4(eEA4 uint8) {}

// S1UENetworkCapability 9.11.3.48
// EEA5 Row, sBit, len = [0, 0], 3 , 1
func (a *S1UENetworkCapability) GetEEA5() (eEA5 uint8) {}

// S1UENetworkCapability 9.11.3.48
// EEA5 Row, sBit, len = [0, 0], 3 , 1
func (a *S1UENetworkCapability) SetEEA5(eEA5 uint8) {}

// S1UENetworkCapability 9.11.3.48
// EEA6 Row, sBit, len = [0, 0], 2 , 1
func (a *S1UENetworkCapability) GetEEA6() (eEA6 uint8) {}

// S1UENetworkCapability 9.11.3.48
// EEA6 Row, sBit, len = [0, 0], 2 , 1
func (a *S1UENetworkCapability) SetEEA6(eEA6 uint8) {}

// S1UENetworkCapability 9.11.3.48
// EEA7 Row, sBit, len = [0, 0], 1 , 1
func (a *S1UENetworkCapability) GetEEA7() (eEA7 uint8) {}

// S1UENetworkCapability 9.11.3.48
// EEA7 Row, sBit, len = [0, 0], 1 , 1
func (a *S1UENetworkCapability) SetEEA7(eEA7 uint8) {}

// S1UENetworkCapability 9.11.3.48
// EIA0 Row, sBit, len = [1, 1], 8 , 1
func (a *S1UENetworkCapability) GetEIA0() (eIA0 uint8) {}

// S1UENetworkCapability 9.11.3.48
// EIA0 Row, sBit, len = [1, 1], 8 , 1
func (a *S1UENetworkCapability) SetEIA0(eIA0 uint8) {}

// S1UENetworkCapability 9.11.3.48
// EIA1_128 Row, sBit, len = [1, 1], 7 , 1
func (a *S1UENetworkCapability) GetEIA1_128() (eIA1_128 uint8) {}

// S1UENetworkCapability 9.11.3.48
// EIA1_128 Row, sBit, len = [1, 1], 7 , 1
func (a *S1UENetworkCapability) SetEIA1_128(eIA1_128 uint8) {}

// S1UENetworkCapability 9.11.3.48
// EIA2_128 Row, sBit, len = [1, 1], 6 , 1
func (a *S1UENetworkCapability) GetEIA2_128() (eIA2_128 uint8) {}

// S1UENetworkCapability 9.11.3.48
// EIA2_128 Row, sBit, len = [1, 1], 6 , 1
func (a *S1UENetworkCapability) SetEIA2_128(eIA2_128 uint8) {}

// S1UENetworkCapability 9.11.3.48
// EIA3_128 Row, sBit, len = [1, 1], 5 , 1
func (a *S1UENetworkCapability) GetEIA3_128() (eIA3_128 uint8) {}

// S1UENetworkCapability 9.11.3.48
// EIA3_128 Row, sBit, len = [1, 1], 5 , 1
func (a *S1UENetworkCapability) SetEIA3_128(eIA3_128 uint8) {}

// S1UENetworkCapability 9.11.3.48
// EIA4 Row, sBit, len = [1, 1], 4 , 1
func (a *S1UENetworkCapability) GetEIA4() (eIA4 uint8) {}

// S1UENetworkCapability 9.11.3.48
// EIA4 Row, sBit, len = [1, 1], 4 , 1
func (a *S1UENetworkCapability) SetEIA4(eIA4 uint8) {}

// S1UENetworkCapability 9.11.3.48
// EIA5 Row, sBit, len = [1, 1], 3 , 1
func (a *S1UENetworkCapability) GetEIA5() (eIA5 uint8) {}

// S1UENetworkCapability 9.11.3.48
// EIA5 Row, sBit, len = [1, 1], 3 , 1
func (a *S1UENetworkCapability) SetEIA5(eIA5 uint8) {}

// S1UENetworkCapability 9.11.3.48
// EIA6 Row, sBit, len = [1, 1], 2 , 1
func (a *S1UENetworkCapability) GetEIA6() (eIA6 uint8) {}

// S1UENetworkCapability 9.11.3.48
// EIA6 Row, sBit, len = [1, 1], 2 , 1
func (a *S1UENetworkCapability) SetEIA6(eIA6 uint8) {}

// S1UENetworkCapability 9.11.3.48
// EIA7 Row, sBit, len = [1, 1], 1 , 1
func (a *S1UENetworkCapability) GetEIA7() (eIA7 uint8) {}

// S1UENetworkCapability 9.11.3.48
// EIA7 Row, sBit, len = [1, 1], 1 , 1
func (a *S1UENetworkCapability) SetEIA7(eIA7 uint8) {}

// S1UENetworkCapability 9.11.3.48
// UEA0 Row, sBit, len = [2, 2], 8 , 1
func (a *S1UENetworkCapability) GetUEA0() (uEA0 uint8) {}

// S1UENetworkCapability 9.11.3.48
// UEA0 Row, sBit, len = [2, 2], 8 , 1
func (a *S1UENetworkCapability) SetUEA0(uEA0 uint8) {}

// S1UENetworkCapability 9.11.3.48
// UEA1 Row, sBit, len = [2, 2], 7 , 1
func (a *S1UENetworkCapability) GetUEA1() (uEA1 uint8) {}

// S1UENetworkCapability 9.11.3.48
// UEA1 Row, sBit, len = [2, 2], 7 , 1
func (a *S1UENetworkCapability) SetUEA1(uEA1 uint8) {}

// S1UENetworkCapability 9.11.3.48
// UEA2 Row, sBit, len = [2, 2], 6 , 1
func (a *S1UENetworkCapability) GetUEA2() (uEA2 uint8) {}

// S1UENetworkCapability 9.11.3.48
// UEA2 Row, sBit, len = [2, 2], 6 , 1
func (a *S1UENetworkCapability) SetUEA2(uEA2 uint8) {}

// S1UENetworkCapability 9.11.3.48
// UEA3 Row, sBit, len = [2, 2], 5 , 1
func (a *S1UENetworkCapability) GetUEA3() (uEA3 uint8) {}

// S1UENetworkCapability 9.11.3.48
// UEA3 Row, sBit, len = [2, 2], 5 , 1
func (a *S1UENetworkCapability) SetUEA3(uEA3 uint8) {}

// S1UENetworkCapability 9.11.3.48
// UEA4 Row, sBit, len = [2, 2], 4 , 1
func (a *S1UENetworkCapability) GetUEA4() (uEA4 uint8) {}

// S1UENetworkCapability 9.11.3.48
// UEA4 Row, sBit, len = [2, 2], 4 , 1
func (a *S1UENetworkCapability) SetUEA4(uEA4 uint8) {}

// S1UENetworkCapability 9.11.3.48
// UEA5 Row, sBit, len = [2, 2], 3 , 1
func (a *S1UENetworkCapability) GetUEA5() (uEA5 uint8) {}

// S1UENetworkCapability 9.11.3.48
// UEA5 Row, sBit, len = [2, 2], 3 , 1
func (a *S1UENetworkCapability) SetUEA5(uEA5 uint8) {}

// S1UENetworkCapability 9.11.3.48
// UEA6 Row, sBit, len = [2, 2], 2 , 1
func (a *S1UENetworkCapability) GetUEA6() (uEA6 uint8) {}

// S1UENetworkCapability 9.11.3.48
// UEA6 Row, sBit, len = [2, 2], 2 , 1
func (a *S1UENetworkCapability) SetUEA6(uEA6 uint8) {}

// S1UENetworkCapability 9.11.3.48
// UEA7 Row, sBit, len = [2, 2], 1 , 1
func (a *S1UENetworkCapability) GetUEA7() (uEA7 uint8) {}

// S1UENetworkCapability 9.11.3.48
// UEA7 Row, sBit, len = [2, 2], 1 , 1
func (a *S1UENetworkCapability) SetUEA7(uEA7 uint8) {}

// S1UENetworkCapability 9.11.3.48
// UCS2 Row, sBit, len = [3, 3], 8 , 1
func (a *S1UENetworkCapability) GetUCS2() (uCS2 uint8) {}

// S1UENetworkCapability 9.11.3.48
// UCS2 Row, sBit, len = [3, 3], 8 , 1
func (a *S1UENetworkCapability) SetUCS2(uCS2 uint8) {}

// S1UENetworkCapability 9.11.3.48
// UIA1 Row, sBit, len = [3, 3], 7 , 1
func (a *S1UENetworkCapability) GetUIA1() (uIA1 uint8) {}

// S1UENetworkCapability 9.11.3.48
// UIA1 Row, sBit, len = [3, 3], 7 , 1
func (a *S1UENetworkCapability) SetUIA1(uIA1 uint8) {}

// S1UENetworkCapability 9.11.3.48
// UIA2 Row, sBit, len = [3, 3], 6 , 1
func (a *S1UENetworkCapability) GetUIA2() (uIA2 uint8) {}

// S1UENetworkCapability 9.11.3.48
// UIA2 Row, sBit, len = [3, 3], 6 , 1
func (a *S1UENetworkCapability) SetUIA2(uIA2 uint8) {}

// S1UENetworkCapability 9.11.3.48
// UIA3 Row, sBit, len = [3, 3], 5 , 1
func (a *S1UENetworkCapability) GetUIA3() (uIA3 uint8) {}

// S1UENetworkCapability 9.11.3.48
// UIA3 Row, sBit, len = [3, 3], 5 , 1
func (a *S1UENetworkCapability) SetUIA3(uIA3 uint8) {}

// S1UENetworkCapability 9.11.3.48
// UIA4 Row, sBit, len = [3, 3], 4 , 1
func (a *S1UENetworkCapability) GetUIA4() (uIA4 uint8) {}

// S1UENetworkCapability 9.11.3.48
// UIA4 Row, sBit, len = [3, 3], 4 , 1
func (a *S1UENetworkCapability) SetUIA4(uIA4 uint8) {}

// S1UENetworkCapability 9.11.3.48
// UIA5 Row, sBit, len = [3, 3], 3 , 1
func (a *S1UENetworkCapability) GetUIA5() (uIA5 uint8) {}

// S1UENetworkCapability 9.11.3.48
// UIA5 Row, sBit, len = [3, 3], 3 , 1
func (a *S1UENetworkCapability) SetUIA5(uIA5 uint8) {}

// S1UENetworkCapability 9.11.3.48
// UIA6 Row, sBit, len = [3, 3], 2 , 1
func (a *S1UENetworkCapability) GetUIA6() (uIA6 uint8) {}

// S1UENetworkCapability 9.11.3.48
// UIA6 Row, sBit, len = [3, 3], 2 , 1
func (a *S1UENetworkCapability) SetUIA6(uIA6 uint8) {}

// S1UENetworkCapability 9.11.3.48
// UIA7 Row, sBit, len = [3, 3], 1 , 1
func (a *S1UENetworkCapability) GetUIA7() (uIA7 uint8) {}

// S1UENetworkCapability 9.11.3.48
// UIA7 Row, sBit, len = [3, 3], 1 , 1
func (a *S1UENetworkCapability) SetUIA7(uIA7 uint8) {}

// S1UENetworkCapability 9.11.3.48
// ProSedd Row, sBit, len = [4, 4], 8 , 1
func (a *S1UENetworkCapability) GetProSedd() (proSedd uint8) {}

// S1UENetworkCapability 9.11.3.48
// ProSedd Row, sBit, len = [4, 4], 8 , 1
func (a *S1UENetworkCapability) SetProSedd(proSedd uint8) {}

// S1UENetworkCapability 9.11.3.48
// ProSe Row, sBit, len = [4, 4], 7 , 1
func (a *S1UENetworkCapability) GetProSe() (proSe uint8) {}

// S1UENetworkCapability 9.11.3.48
// ProSe Row, sBit, len = [4, 4], 7 , 1
func (a *S1UENetworkCapability) SetProSe(proSe uint8) {}

// S1UENetworkCapability 9.11.3.48
// H245ASH Row, sBit, len = [4, 4], 6 , 1
func (a *S1UENetworkCapability) GetH245ASH() (h245ASH uint8) {}

// S1UENetworkCapability 9.11.3.48
// H245ASH Row, sBit, len = [4, 4], 6 , 1
func (a *S1UENetworkCapability) SetH245ASH(h245ASH uint8) {}

// S1UENetworkCapability 9.11.3.48
// ACCCSFB Row, sBit, len = [4, 4], 5 , 1
func (a *S1UENetworkCapability) GetACCCSFB() (aCCCSFB uint8) {}

// S1UENetworkCapability 9.11.3.48
// ACCCSFB Row, sBit, len = [4, 4], 5 , 1
func (a *S1UENetworkCapability) SetACCCSFB(aCCCSFB uint8) {}

// S1UENetworkCapability 9.11.3.48
// LPP Row, sBit, len = [4, 4], 4 , 1
func (a *S1UENetworkCapability) GetLPP() (lPP uint8) {}

// S1UENetworkCapability 9.11.3.48
// LPP Row, sBit, len = [4, 4], 4 , 1
func (a *S1UENetworkCapability) SetLPP(lPP uint8) {}

// S1UENetworkCapability 9.11.3.48
// LCS Row, sBit, len = [4, 4], 3 , 1
func (a *S1UENetworkCapability) GetLCS() (lCS uint8) {}

// S1UENetworkCapability 9.11.3.48
// LCS Row, sBit, len = [4, 4], 3 , 1
func (a *S1UENetworkCapability) SetLCS(lCS uint8) {}

// S1UENetworkCapability 9.11.3.48
// xSRVCC Row, sBit, len = [4, 4], 2 , 1
func (a *S1UENetworkCapability) GetxSRVCC() (xSRVCC uint8) {}

// S1UENetworkCapability 9.11.3.48
// xSRVCC Row, sBit, len = [4, 4], 2 , 1
func (a *S1UENetworkCapability) SetxSRVCC(xSRVCC uint8) {}

// S1UENetworkCapability 9.11.3.48
// NF Row, sBit, len = [4, 4], 1 , 1
func (a *S1UENetworkCapability) GetNF() (nF uint8) {}

// S1UENetworkCapability 9.11.3.48
// NF Row, sBit, len = [4, 4], 1 , 1
func (a *S1UENetworkCapability) SetNF(nF uint8) {}

// S1UENetworkCapability 9.11.3.48
// EPCO Row, sBit, len = [5, 5], 8 , 1
func (a *S1UENetworkCapability) GetEPCO() (ePCO uint8) {}

// S1UENetworkCapability 9.11.3.48
// EPCO Row, sBit, len = [5, 5], 8 , 1
func (a *S1UENetworkCapability) SetEPCO(ePCO uint8) {}

// S1UENetworkCapability 9.11.3.48
// HCCPCIOT Row, sBit, len = [5, 5], 7 , 1
func (a *S1UENetworkCapability) GetHCCPCIOT() (hCCPCIOT uint8) {}

// S1UENetworkCapability 9.11.3.48
// HCCPCIOT Row, sBit, len = [5, 5], 7 , 1
func (a *S1UENetworkCapability) SetHCCPCIOT(hCCPCIOT uint8) {}

// S1UENetworkCapability 9.11.3.48
// ERwoPDN Row, sBit, len = [5, 5], 6 , 1
func (a *S1UENetworkCapability) GetERwoPDN() (eRwoPDN uint8) {}

// S1UENetworkCapability 9.11.3.48
// ERwoPDN Row, sBit, len = [5, 5], 6 , 1
func (a *S1UENetworkCapability) SetERwoPDN(eRwoPDN uint8) {}

// S1UENetworkCapability 9.11.3.48
// S1UData Row, sBit, len = [5, 5], 5 , 1
func (a *S1UENetworkCapability) GetS1UData() (s1UData uint8) {}

// S1UENetworkCapability 9.11.3.48
// S1UData Row, sBit, len = [5, 5], 5 , 1
func (a *S1UENetworkCapability) SetS1UData(s1UData uint8) {}

// S1UENetworkCapability 9.11.3.48
// UPCIot Row, sBit, len = [5, 5], 4 , 1
func (a *S1UENetworkCapability) GetUPCIot() (uPCIot uint8) {}

// S1UENetworkCapability 9.11.3.48
// UPCIot Row, sBit, len = [5, 5], 4 , 1
func (a *S1UENetworkCapability) SetUPCIot(uPCIot uint8) {}

// S1UENetworkCapability 9.11.3.48
// CPCIot Row, sBit, len = [5, 5], 3 , 1
func (a *S1UENetworkCapability) GetCPCIot() (cPCIot uint8) {}

// S1UENetworkCapability 9.11.3.48
// CPCIot Row, sBit, len = [5, 5], 3 , 1
func (a *S1UENetworkCapability) SetCPCIot(cPCIot uint8) {}

// S1UENetworkCapability 9.11.3.48
// Proserelay Row, sBit, len = [5, 5], 2 , 1
func (a *S1UENetworkCapability) GetProserelay() (proserelay uint8) {}

// S1UENetworkCapability 9.11.3.48
// Proserelay Row, sBit, len = [5, 5], 2 , 1
func (a *S1UENetworkCapability) SetProserelay(proserelay uint8) {}

// S1UENetworkCapability 9.11.3.48
// ProSedc Row, sBit, len = [5, 5], 1 , 1
func (a *S1UENetworkCapability) GetProSedc() (proSedc uint8) {}

// S1UENetworkCapability 9.11.3.48
// ProSedc Row, sBit, len = [5, 5], 1 , 1
func (a *S1UENetworkCapability) SetProSedc(proSedc uint8) {}

// S1UENetworkCapability 9.11.3.48
// Bearer15 Row, sBit, len = [6, 6], 8 , 1
func (a *S1UENetworkCapability) GetBearer15() (bearer15 uint8) {}

// S1UENetworkCapability 9.11.3.48
// Bearer15 Row, sBit, len = [6, 6], 8 , 1
func (a *S1UENetworkCapability) SetBearer15(bearer15 uint8) {}

// S1UENetworkCapability 9.11.3.48
// SGC Row, sBit, len = [6, 6], 7 , 1
func (a *S1UENetworkCapability) GetSGC() (sGC uint8) {}

// S1UENetworkCapability 9.11.3.48
// SGC Row, sBit, len = [6, 6], 7 , 1
func (a *S1UENetworkCapability) SetSGC(sGC uint8) {}

// S1UENetworkCapability 9.11.3.48
// N1mode Row, sBit, len = [6, 6], 6 , 1
func (a *S1UENetworkCapability) GetN1mode() (n1mode uint8) {}

// S1UENetworkCapability 9.11.3.48
// N1mode Row, sBit, len = [6, 6], 6 , 1
func (a *S1UENetworkCapability) SetN1mode(n1mode uint8) {}

// S1UENetworkCapability 9.11.3.48
// DCNR Row, sBit, len = [6, 6], 5 , 1
func (a *S1UENetworkCapability) GetDCNR() (dCNR uint8) {}

// S1UENetworkCapability 9.11.3.48
// DCNR Row, sBit, len = [6, 6], 5 , 1
func (a *S1UENetworkCapability) SetDCNR(dCNR uint8) {}

// S1UENetworkCapability 9.11.3.48
// CPbackoff Row, sBit, len = [6, 6], 4 , 1
func (a *S1UENetworkCapability) GetCPbackoff() (cPbackoff uint8) {}

// S1UENetworkCapability 9.11.3.48
// CPbackoff Row, sBit, len = [6, 6], 4 , 1
func (a *S1UENetworkCapability) SetCPbackoff(cPbackoff uint8) {}

// S1UENetworkCapability 9.11.3.48
// RestrictEC Row, sBit, len = [6, 6], 3 , 1
func (a *S1UENetworkCapability) GetRestrictEC() (restrictEC uint8) {}

// S1UENetworkCapability 9.11.3.48
// RestrictEC Row, sBit, len = [6, 6], 3 , 1
func (a *S1UENetworkCapability) SetRestrictEC(restrictEC uint8) {}

// S1UENetworkCapability 9.11.3.48
// V2XPC5 Row, sBit, len = [6, 6], 2 , 1
func (a *S1UENetworkCapability) GetV2XPC5() (v2XPC5 uint8) {}

// S1UENetworkCapability 9.11.3.48
// V2XPC5 Row, sBit, len = [6, 6], 2 , 1
func (a *S1UENetworkCapability) SetV2XPC5(v2XPC5 uint8) {}

// S1UENetworkCapability 9.11.3.48
// MulitpeDRB Row, sBit, len = [6, 6], 1 , 1
func (a *S1UENetworkCapability) GetMulitpeDRB() (mulitpeDRB uint8) {}

// S1UENetworkCapability 9.11.3.48
// MulitpeDRB Row, sBit, len = [6, 6], 1 , 1
func (a *S1UENetworkCapability) SetMulitpeDRB(mulitpeDRB uint8) {}

// S1UENetworkCapability 9.11.3.48
// Spare Row, sBit, len = [7, 12], 8 , INF
func (a *S1UENetworkCapability) GetSpare() (spare []uint8) {}

// S1UENetworkCapability 9.11.3.48
// Spare Row, sBit, len = [7, 12], 8 , INF
func (a *S1UENetworkCapability) SetSpare(spare []uint8) {}
