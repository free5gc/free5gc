//go:binary-only-package

package ngapType

// Need to import "free5gc/lib/aper" if it uses "aper"

/* Sequence of = 35, FULL Name = struct AMF_TNLAssociationSetupList */
/* AMFTNLAssociationSetupItem */
type AMFTNLAssociationSetupList struct {
	List []AMFTNLAssociationSetupItem `aper:"valueExt,sizeLB:1,sizeUB:32"`
}
