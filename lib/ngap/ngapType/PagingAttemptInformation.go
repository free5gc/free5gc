//go:binary-only-package

package ngapType

// Need to import "free5gc/lib/aper" if it uses "aper"

type PagingAttemptInformation struct {
	PagingAttemptCount             PagingAttemptCount
	IntendedNumberOfPagingAttempts IntendedNumberOfPagingAttempts
	NextPagingAreaScope            *NextPagingAreaScope                                      `aper:"optional"`
	IEExtensions                   *ProtocolExtensionContainerPagingAttemptInformationExtIEs `aper:"optional"`
}
