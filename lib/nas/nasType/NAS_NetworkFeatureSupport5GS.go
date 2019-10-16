//go:binary-only-package

package nasType

// NetworkFeatureSupport5GS 9.11.3.5
// MPSI Row, sBit, len = [0, 0], 8 , 1
// IWKN26 Row, sBit, len = [0, 0], 7 , 1
// EMF Row, sBit, len = [0, 0], 6 , 2
// EMC Row, sBit, len = [0, 0], 4 , 2
// IMSVoPSN3GPP Row, sBit, len = [0, 0], 2 , 1
// IMSVoPS3GPP Row, sBit, len = [0, 0], 1 , 1
// MCSI Row, sBit, len = [1, 1], 2 , 1
// EMCN Row, sBit, len = [1, 1], 1 , 1
// Spare Row, sBit, len = [2, 2], 8 , 8
type NetworkFeatureSupport5GS struct {
	Iei   uint8
	Len   uint8
	Octet [3]uint8
}

func NewNetworkFeatureSupport5GS(iei uint8) (networkFeatureSupport5GS *NetworkFeatureSupport5GS) {}

// NetworkFeatureSupport5GS 9.11.3.5
// Iei Row, sBit, len = [], 8, 8
func (a *NetworkFeatureSupport5GS) GetIei() (iei uint8) {}

// NetworkFeatureSupport5GS 9.11.3.5
// Iei Row, sBit, len = [], 8, 8
func (a *NetworkFeatureSupport5GS) SetIei(iei uint8) {}

// NetworkFeatureSupport5GS 9.11.3.5
// Len Row, sBit, len = [], 8, 8
func (a *NetworkFeatureSupport5GS) GetLen() (len uint8) {}

// NetworkFeatureSupport5GS 9.11.3.5
// Len Row, sBit, len = [], 8, 8
func (a *NetworkFeatureSupport5GS) SetLen(len uint8) {}

// NetworkFeatureSupport5GS 9.11.3.5
// MPSI Row, sBit, len = [0, 0], 8 , 1
func (a *NetworkFeatureSupport5GS) GetMPSI() (mPSI uint8) {}

// NetworkFeatureSupport5GS 9.11.3.5
// MPSI Row, sBit, len = [0, 0], 8 , 1
func (a *NetworkFeatureSupport5GS) SetMPSI(mPSI uint8) {}

// NetworkFeatureSupport5GS 9.11.3.5
// IWKN26 Row, sBit, len = [0, 0], 7 , 1
func (a *NetworkFeatureSupport5GS) GetIWKN26() (iWKN26 uint8) {}

// NetworkFeatureSupport5GS 9.11.3.5
// IWKN26 Row, sBit, len = [0, 0], 7 , 1
func (a *NetworkFeatureSupport5GS) SetIWKN26(iWKN26 uint8) {}

// NetworkFeatureSupport5GS 9.11.3.5
// EMF Row, sBit, len = [0, 0], 6 , 2
func (a *NetworkFeatureSupport5GS) GetEMF() (eMF uint8) {}

// NetworkFeatureSupport5GS 9.11.3.5
// EMF Row, sBit, len = [0, 0], 6 , 2
func (a *NetworkFeatureSupport5GS) SetEMF(eMF uint8) {}

// NetworkFeatureSupport5GS 9.11.3.5
// EMC Row, sBit, len = [0, 0], 4 , 2
func (a *NetworkFeatureSupport5GS) GetEMC() (eMC uint8) {}

// NetworkFeatureSupport5GS 9.11.3.5
// EMC Row, sBit, len = [0, 0], 4 , 2
func (a *NetworkFeatureSupport5GS) SetEMC(eMC uint8) {}

// NetworkFeatureSupport5GS 9.11.3.5
// IMSVoPSN3GPP Row, sBit, len = [0, 0], 2 , 1
func (a *NetworkFeatureSupport5GS) GetIMSVoPSN3GPP() (iMSVoPSN3GPP uint8) {}

// NetworkFeatureSupport5GS 9.11.3.5
// IMSVoPSN3GPP Row, sBit, len = [0, 0], 2 , 1
func (a *NetworkFeatureSupport5GS) SetIMSVoPSN3GPP(iMSVoPSN3GPP uint8) {}

// NetworkFeatureSupport5GS 9.11.3.5
// IMSVoPS3GPP Row, sBit, len = [0, 0], 1 , 1
func (a *NetworkFeatureSupport5GS) GetIMSVoPS3GPP() (iMSVoPS3GPP uint8) {}

// NetworkFeatureSupport5GS 9.11.3.5
// IMSVoPS3GPP Row, sBit, len = [0, 0], 1 , 1
func (a *NetworkFeatureSupport5GS) SetIMSVoPS3GPP(iMSVoPS3GPP uint8) {}

// NetworkFeatureSupport5GS 9.11.3.5
// MCSI Row, sBit, len = [1, 1], 2 , 1
func (a *NetworkFeatureSupport5GS) GetMCSI() (mCSI uint8) {}

// NetworkFeatureSupport5GS 9.11.3.5
// MCSI Row, sBit, len = [1, 1], 2 , 1
func (a *NetworkFeatureSupport5GS) SetMCSI(mCSI uint8) {}

// NetworkFeatureSupport5GS 9.11.3.5
// EMCN Row, sBit, len = [1, 1], 1 , 1
func (a *NetworkFeatureSupport5GS) GetEMCN() (eMCN uint8) {}

// NetworkFeatureSupport5GS 9.11.3.5
// EMCN Row, sBit, len = [1, 1], 1 , 1
func (a *NetworkFeatureSupport5GS) SetEMCN(eMCN uint8) {}

// NetworkFeatureSupport5GS 9.11.3.5
// Spare Row, sBit, len = [2, 2], 8 , 8
func (a *NetworkFeatureSupport5GS) GetSpare() (spare uint8) {}

// NetworkFeatureSupport5GS 9.11.3.5
// Spare Row, sBit, len = [2, 2], 8 , 8
func (a *NetworkFeatureSupport5GS) SetSpare(spare uint8) {}
