//go:binary-only-package

package ngapType

import "free5gc/lib/aper"

// Need to import "free5gc/lib/aper" if it uses "aper"

const (
	CauseRadioNetworkPresentUnspecified                                              aper.Enumerated = 0
	CauseRadioNetworkPresentTxnrelocoverallExpiry                                    aper.Enumerated = 1
	CauseRadioNetworkPresentSuccessfulHandover                                       aper.Enumerated = 2
	CauseRadioNetworkPresentReleaseDueToNgranGeneratedReason                         aper.Enumerated = 3
	CauseRadioNetworkPresentReleaseDueTo5gcGeneratedReason                           aper.Enumerated = 4
	CauseRadioNetworkPresentHandoverCancelled                                        aper.Enumerated = 5
	CauseRadioNetworkPresentPartialHandover                                          aper.Enumerated = 6
	CauseRadioNetworkPresentHoFailureInTarget5GCNgranNodeOrTargetSystem              aper.Enumerated = 7
	CauseRadioNetworkPresentHoTargetNotAllowed                                       aper.Enumerated = 8
	CauseRadioNetworkPresentTngrelocoverallExpiry                                    aper.Enumerated = 9
	CauseRadioNetworkPresentTngrelocprepExpiry                                       aper.Enumerated = 10
	CauseRadioNetworkPresentCellNotAvailable                                         aper.Enumerated = 11
	CauseRadioNetworkPresentUnknownTargetID                                          aper.Enumerated = 12
	CauseRadioNetworkPresentNoRadioResourcesAvailableInTargetCell                    aper.Enumerated = 13
	CauseRadioNetworkPresentUnknownLocalUENGAPID                                     aper.Enumerated = 14
	CauseRadioNetworkPresentInconsistentRemoteUENGAPID                               aper.Enumerated = 15
	CauseRadioNetworkPresentHandoverDesirableForRadioReason                          aper.Enumerated = 16
	CauseRadioNetworkPresentTimeCriticalHandover                                     aper.Enumerated = 17
	CauseRadioNetworkPresentResourceOptimisationHandover                             aper.Enumerated = 18
	CauseRadioNetworkPresentReduceLoadInServingCell                                  aper.Enumerated = 19
	CauseRadioNetworkPresentUserInactivity                                           aper.Enumerated = 20
	CauseRadioNetworkPresentRadioConnectionWithUeLost                                aper.Enumerated = 21
	CauseRadioNetworkPresentRadioResourcesNotAvailable                               aper.Enumerated = 22
	CauseRadioNetworkPresentInvalidQosCombination                                    aper.Enumerated = 23
	CauseRadioNetworkPresentFailureInRadioInterfaceProcedure                         aper.Enumerated = 24
	CauseRadioNetworkPresentInteractionWithOtherProcedure                            aper.Enumerated = 25
	CauseRadioNetworkPresentUnknownPDUSessionID                                      aper.Enumerated = 26
	CauseRadioNetworkPresentUnkownQosFlowID                                          aper.Enumerated = 27
	CauseRadioNetworkPresentMultiplePDUSessionIDInstances                            aper.Enumerated = 28
	CauseRadioNetworkPresentMultipleQosFlowIDInstances                               aper.Enumerated = 29
	CauseRadioNetworkPresentEncryptionAndOrIntegrityProtectionAlgorithmsNotSupported aper.Enumerated = 30
	CauseRadioNetworkPresentNgIntraSystemHandoverTriggered                           aper.Enumerated = 31
	CauseRadioNetworkPresentNgInterSystemHandoverTriggered                           aper.Enumerated = 32
	CauseRadioNetworkPresentXnHandoverTriggered                                      aper.Enumerated = 33
	CauseRadioNetworkPresentNotSupported5QIValue                                     aper.Enumerated = 34
	CauseRadioNetworkPresentUeContextTransfer                                        aper.Enumerated = 35
	CauseRadioNetworkPresentImsVoiceEpsFallbackOrRatFallbackTriggered                aper.Enumerated = 36
	CauseRadioNetworkPresentUpIntegrityProtectionNotPossible                         aper.Enumerated = 37
	CauseRadioNetworkPresentUpConfidentialityProtectionNotPossible                   aper.Enumerated = 38
	CauseRadioNetworkPresentSliceNotSupported                                        aper.Enumerated = 39
	CauseRadioNetworkPresentUeInRrcInactiveStateNotReachable                         aper.Enumerated = 40
	CauseRadioNetworkPresentRedirection                                              aper.Enumerated = 41
	CauseRadioNetworkPresentResourcesNotAvailableForTheSlice                         aper.Enumerated = 42
	CauseRadioNetworkPresentUeMaxIntegrityProtectedDataRateReason                    aper.Enumerated = 43
	CauseRadioNetworkPresentReleaseDueToCnDetectedMobility                           aper.Enumerated = 44
)

type CauseRadioNetwork struct {
	Value aper.Enumerated `aper:"valueExt,valueLB:0,valueUB:44"`
}
