//go:binary-only-package

package ngapTestpacket

import (
	"encoding/hex"
	"free5gc/lib/aper"
	"free5gc/lib/ngap/ngapConvert"
	"free5gc/lib/ngap/ngapType"
)

// TODO: check test data
var TestPlmn ngapType.PLMNIdentity

func init() {}

func BuildNGSetupRequest() (pdu ngapType.NGAPPDU) {}

func BuildNGReset() (pdu ngapType.NGAPPDU) {}

func BuildNGResetAcknowledge() (pdu ngapType.NGAPPDU) {}

func BuildInitialUEMessage(ranUeNgapID int64, nasPdu []byte, fiveGSTmsi string) (pdu ngapType.NGAPPDU) {}

func BuildErrorIndication() (pdu ngapType.NGAPPDU) {}

func BuildUEContextReleaseRequest(amfUeNgapID, ranUeNgapID int64) (pdu ngapType.NGAPPDU) {}

func BuildUEContextReleaseComplete(amfUeNgapID, ranUeNgapID int64) (pdu ngapType.NGAPPDU) {}

func BuildUEContextModificationResponse(amfUeNgapID, ranUeNgapID int64) (pdu ngapType.NGAPPDU) {}

func BuildUplinkNasTransport(amfUeNgapID, ranUeNgapID int64, nasPdu []byte) (pdu ngapType.NGAPPDU) {}

func BuildInitialContextSetupResponse(amfUeNgapID, ranUeNgapID int64, pduSessionFailedList *ngapType.PDUSessionResourceFailedToSetupListCxtRes) (pdu ngapType.NGAPPDU) {}

func BuildInitialContextSetupFailure(amfUeNgapID, ranUeNgapID int64) (pdu ngapType.NGAPPDU) {}

func BuildPathSwitchRequest(sourceAmfUeNgapID, ranUeNgapID int64) (pdu ngapType.NGAPPDU) {}

func BuildHandoverRequestAcknowledge(amfUeNgapID, ranUeNgapID int64) (pdu ngapType.NGAPPDU) {}

func BuildHandoverFailure(amfUeNgapID int64) (pdu ngapType.NGAPPDU) {}

func BuildPDUSessionResourceReleaseResponse() (pdu ngapType.NGAPPDU) {}

func BuildAMFConfigurationUpdateFailure() (pdu ngapType.NGAPPDU) {}

func BuildUERadioCapabilityCheckResponse() (pdu ngapType.NGAPPDU) {}

func BuildHandoverCancel() (pdu ngapType.NGAPPDU) {}
func BuildLocationReportingFailureIndication() (pdu ngapType.NGAPPDU) {}

func BuildPDUSessionResourceSetupResponse(amfUeNgapID, ranUeNgapID int64) (pdu ngapType.NGAPPDU) {}

func BuildPDUSessionResourceModifyResponse(amfUeNgapID, ranUeNgapID int64) (pdu ngapType.NGAPPDU) {}

func BuildPDUSessionResourceNotify() (pdu ngapType.NGAPPDU) {}

func BuildPDUSessionResourceModifyIndication(amfUeNgapID, ranUeNgapID int64) (pdu ngapType.NGAPPDU) {}

func BuildUEContextModificationFailure(amfUeNgapID, ranUeNgapID int64) (pdu ngapType.NGAPPDU) {}

func BuildRRCInactiveTransitionReport() (pdu ngapType.NGAPPDU) {}

func BuildHandoverNotify() (pdu ngapType.NGAPPDU) {}

func BuildUplinkRanStatusTransfer(amfUeNgapID, ranUeNgapID int64) (pdu ngapType.NGAPPDU) {}

func BuildNasNonDeliveryIndication(amfUeNgapID, ranUeNgapID int64, naspdu aper.OctetString) (pdu ngapType.NGAPPDU) {}

func BuildRanConfigurationUpdate() (pdu ngapType.NGAPPDU) {}

func BuildAMFStatusIndication() (pdu ngapType.NGAPPDU) {}

func BuildUplinkRanConfigurationTransfer() (pdu ngapType.NGAPPDU) {}

func BuildUplinkUEAssociatedNRPPATransport() (pdu ngapType.NGAPPDU) {}

func BuildUplinkNonUEAssociatedNRPPATransport() (pdu ngapType.NGAPPDU) {}

func BuildLocationReport() (pdu ngapType.NGAPPDU) {}

func BuildUETNLABindingReleaseRequest() (pdu ngapType.NGAPPDU) {}

func BuildUERadioCapabilityInfoIndication() (pdu ngapType.NGAPPDU) {}

func BuildAMFConfigurationUpdateAcknowledge() (pdu ngapType.NGAPPDU) {}

func BuildHandoverRequired() (pdu ngapType.NGAPPDU) {}

func BuildCellTrafficTrace(amfUeNgapID, ranUeNgapID int64) (pdu ngapType.NGAPPDU) {}

func buildPDUSessionResourceSetupResponseTransfer() (data ngapType.PDUSessionResourceSetupResponseTransfer) {}

func buildPDUSessionResourceModifyResponseTransfer() (data ngapType.PDUSessionResourceModifyResponseTransfer) {}

func buildPDUSessionResourceSetupUnsucessfulTransfer() (data ngapType.PDUSessionResourceSetupUnsuccessfulTransfer) {}

func buildPDUSessionResourceModifyUnsuccessfulTransfer() (data ngapType.PDUSessionResourceModifyUnsuccessfulTransfer) {}

func buildPDUSessionResourceReleaseResponseTransfer() (data ngapType.PDUSessionResourceReleaseResponse) {}

func buildPathSwitchRequestTransfer() (data ngapType.PathSwitchRequestTransfer) {}

func buildPDUSessionResourceModifyIndicationTransfer() (data ngapType.PDUSessionResourceModifyIndicationTransfer) {}

func buildPathSwitchRequestSetupFailedTransfer() (data ngapType.PathSwitchRequestSetupFailedTransfer) {}

func buildHandoverRequestAcknowledgeTransfer() (data ngapType.HandoverRequestAcknowledgeTransfer) {}

func buildHandoverResourceAllocationUnsuccessfulTransfer() (data ngapType.HandoverResourceAllocationUnsuccessfulTransfer) {}

func GetPDUSessionResourceSetupResponseTransfer() []byte {}

func GetPDUSessionResourceModifyResponseTransfer() []byte {}

func GetPDUSessionResourceSetupUnsucessfulTransfer() []byte {}

func GetPDUSessionResourceModifyUnsuccessfulTransfer() []byte {}

func GetPathSwitchRequestTransfer() []byte {}

func GetPathSwitchRequestSetupFailedTransfer() []byte {}

func GetPDUSessionResourceModifyIndicationTransfer() []byte {}

func GetPDUSessionResourceReleaseResponseTransfer() []byte {}

func GetHandoverRequestAcknowledgeTransfer() []byte {}

func GetHandoverResourceAllocationUnsuccessfulTransfer() []byte {}

func BuildInitialContextSetupResponseForRegistraionTest(amfUeNgapID, ranUeNgapID int64) (pdu ngapType.NGAPPDU) {}

func BuildPDUSessionResourceSetupResponseForRegistrationTest(amfUeNgapID, ranUeNgapID int64) (pdu ngapType.NGAPPDU) {}

func BuildPDUSessionResourceReleaseResponseForReleaseTest(amfUeNgapID, ranUeNgapID int64) (pdu ngapType.NGAPPDU) {}
