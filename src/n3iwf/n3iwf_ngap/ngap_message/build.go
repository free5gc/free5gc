package ngap_message

import (
	"encoding/hex"

	"free5gc/lib/aper"
	"free5gc/lib/ngap"
	"free5gc/lib/ngap/ngapConvert"
	"free5gc/lib/ngap/ngapType"
	"free5gc/src/n3iwf/n3iwf_context"
	"free5gc/src/n3iwf/n3iwf_util"
)

func BuildNGSetupRequest() ([]byte, error) {

	n3iwfSelf := n3iwf_context.N3IWFSelf()
	var pdu ngapType.NGAPPDU
	pdu.Present = ngapType.NGAPPDUPresentInitiatingMessage
	pdu.InitiatingMessage = new(ngapType.InitiatingMessage)

	initiatingMessage := pdu.InitiatingMessage
	initiatingMessage.ProcedureCode.Value = ngapType.ProcedureCodeNGSetup
	initiatingMessage.Criticality.Value = ngapType.CriticalityPresentReject

	initiatingMessage.Value.Present = ngapType.InitiatingMessagePresentNGSetupRequest
	initiatingMessage.Value.NGSetupRequest = new(ngapType.NGSetupRequest)

	nGSetupRequest := initiatingMessage.Value.NGSetupRequest
	nGSetupRequestIEs := &nGSetupRequest.ProtocolIEs

	// GlobalRANNodeID
	ie := ngapType.NGSetupRequestIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDGlobalRANNodeID
	ie.Criticality.Value = ngapType.CriticalityPresentReject
	ie.Value.Present = ngapType.NGSetupRequestIEsPresentGlobalRANNodeID
	ie.Value.GlobalRANNodeID = new(ngapType.GlobalRANNodeID)

	globalRANNodeID := ie.Value.GlobalRANNodeID
	globalRANNodeID.Present = ngapType.GlobalRANNodeIDPresentGlobalN3IWFID
	globalRANNodeID.GlobalN3IWFID = new(ngapType.GlobalN3IWFID)

	globalN3IWFID := globalRANNodeID.GlobalN3IWFID
	globalN3IWFID.PLMNIdentity = n3iwf_util.PlmnIdToNgap(n3iwfSelf.NFInfo.GlobalN3IWFID.PLMNID)
	globalN3IWFID.N3IWFID.Present = ngapType.N3IWFIDPresentN3IWFID
	globalN3IWFID.N3IWFID.N3IWFID = n3iwf_util.N3iwfIdToNgap(n3iwfSelf.NFInfo.GlobalN3IWFID.N3IWFID)
	nGSetupRequestIEs.List = append(nGSetupRequestIEs.List, ie)

	// RANNodeName
	ie = ngapType.NGSetupRequestIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDRANNodeName
	ie.Criticality.Value = ngapType.CriticalityPresentIgnore
	ie.Value.Present = ngapType.NGSetupRequestIEsPresentRANNodeName
	ie.Value.RANNodeName = new(ngapType.RANNodeName)

	rANNodeName := ie.Value.RANNodeName
	rANNodeName.Value = n3iwfSelf.NFInfo.RanNodeName
	nGSetupRequestIEs.List = append(nGSetupRequestIEs.List, ie)
	// SupportedTAList
	ie = ngapType.NGSetupRequestIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDSupportedTAList
	ie.Criticality.Value = ngapType.CriticalityPresentReject
	ie.Value.Present = ngapType.NGSetupRequestIEsPresentSupportedTAList
	ie.Value.SupportedTAList = new(ngapType.SupportedTAList)

	supportedTAList := ie.Value.SupportedTAList

	for _, supportedTAItemLocal := range n3iwfSelf.NFInfo.SupportedTAList {
		// SupportedTAItem in SupportedTAList
		supportedTAItem := ngapType.SupportedTAItem{}
		supportedTAItem.TAC.Value, _ = hex.DecodeString(supportedTAItemLocal.TAC)

		broadcastPLMNList := &supportedTAItem.BroadcastPLMNList

		for _, broadcastPLMNListLocal := range supportedTAItemLocal.BroadcastPLMNList {
			// BroadcastPLMNItem in BroadcastPLMNList
			broadcastPLMNItem := ngapType.BroadcastPLMNItem{}
			broadcastPLMNItem.PLMNIdentity = n3iwf_util.PlmnIdToNgap(broadcastPLMNListLocal.PLMNID)

			sliceSupportList := &broadcastPLMNItem.TAISliceSupportList

			for _, sliceSupportItemLocal := range broadcastPLMNListLocal.TAISliceSupportList {
				// SliceSupportItem in SliceSupportList
				sliceSupportItem := ngapType.SliceSupportItem{}
				sliceSupportItem.SNSSAI.SST.Value, _ = hex.DecodeString(sliceSupportItemLocal.SNSSAI.SST)

				if sliceSupportItemLocal.SNSSAI.SD != "" {
					sliceSupportItem.SNSSAI.SD = new(ngapType.SD)
					sliceSupportItem.SNSSAI.SD.Value, _ = hex.DecodeString(sliceSupportItemLocal.SNSSAI.SD)
				}

				sliceSupportList.List = append(sliceSupportList.List, sliceSupportItem)
			}

			broadcastPLMNList.List = append(broadcastPLMNList.List, broadcastPLMNItem)
		}

		supportedTAList.List = append(supportedTAList.List, supportedTAItem)
	}

	nGSetupRequestIEs.List = append(nGSetupRequestIEs.List, ie)

	/*
		* The reason PagingDRX ie was commented is that in TS23.501
		* PagingDRX was mentioned to be used only for 3GPP access.
		* However, the question that if the paging function for N3IWF
		* is needed requires verification.

			// PagingDRX
			ie = ngapType.NGSetupRequestIEs{}
			ie.Id.Value = ngapType.ProtocolIEIDDefaultPagingDRX
			ie.Criticality.Value = ngapType.CriticalityPresentIgnore
			ie.Value.Present = ngapType.NGSetupRequestIEsPresentDefaultPagingDRX
			ie.Value.DefaultPagingDRX = new(ngapType.PagingDRX)

			pagingDRX := ie.Value.DefaultPagingDRX
			pagingDRX.Value = ngapType.PagingDRXPresentV128
			nGSetupRequestIEs.List = append(nGSetupRequestIEs.List, ie)
	*/

	return ngap.Encoder(pdu)
}

func BuildNGReset() ([]byte, error) {
	var pdu ngapType.NGAPPDU
	return ngap.Encoder(pdu)
}

func BuildNGResetAcknowledge() ([]byte, error) {
	var pdu ngapType.NGAPPDU
	return ngap.Encoder(pdu)
}

func BuildInitialContextSetupResponse(
	ue *n3iwf_context.N3IWFUe,
	responseList *ngapType.PDUSessionResourceSetupListCxtRes,
	failedList *ngapType.PDUSessionResourceFailedToSetupListCxtRes,
	criticalityDiagnostics *ngapType.CriticalityDiagnostics) ([]byte, error) {

	var pdu ngapType.NGAPPDU
	pdu.Present = ngapType.NGAPPDUPresentSuccessfulOutcome
	pdu.SuccessfulOutcome = new(ngapType.SuccessfulOutcome)

	successfulOutcome := pdu.SuccessfulOutcome
	successfulOutcome.ProcedureCode.Value = ngapType.ProcedureCodeInitialContextSetup
	successfulOutcome.Criticality.Value = ngapType.CriticalityPresentReject

	successfulOutcome.Value.Present = ngapType.SuccessfulOutcomePresentInitialContextSetupResponse
	successfulOutcome.Value.InitialContextSetupResponse = new(ngapType.InitialContextSetupResponse)

	initialContextSetupResponse := successfulOutcome.Value.InitialContextSetupResponse
	initialContextSetupResponseIEs := &initialContextSetupResponse.ProtocolIEs

	// AMF UE NGAP ID
	ie := ngapType.InitialContextSetupResponseIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDAMFUENGAPID
	ie.Criticality.Value = ngapType.CriticalityPresentIgnore
	ie.Value.Present = ngapType.InitialContextSetupResponseIEsPresentAMFUENGAPID
	ie.Value.AMFUENGAPID = new(ngapType.AMFUENGAPID)

	aMFUENGAPID := ie.Value.AMFUENGAPID
	aMFUENGAPID.Value = ue.AmfUeNgapId

	initialContextSetupResponseIEs.List = append(initialContextSetupResponseIEs.List, ie)

	// RAN UE NGAP ID
	ie = ngapType.InitialContextSetupResponseIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDRANUENGAPID
	ie.Criticality.Value = ngapType.CriticalityPresentIgnore
	ie.Value.Present = ngapType.InitialContextSetupResponseIEsPresentRANUENGAPID
	ie.Value.RANUENGAPID = new(ngapType.RANUENGAPID)

	rANUENGAPID := ie.Value.RANUENGAPID
	rANUENGAPID.Value = ue.RanUeNgapId

	initialContextSetupResponseIEs.List = append(initialContextSetupResponseIEs.List, ie)

	// PDU Session Resource Setup Response List (optional)
	if responseList != nil && len(responseList.List) > 0 {
		ie = ngapType.InitialContextSetupResponseIEs{}
		ie.Id.Value = ngapType.ProtocolIEIDPDUSessionResourceSetupListCxtRes
		ie.Criticality.Value = ngapType.CriticalityPresentIgnore
		ie.Value.Present = ngapType.InitialContextSetupResponseIEsPresentPDUSessionResourceSetupListCxtRes
		ie.Value.PDUSessionResourceSetupListCxtRes = responseList
		initialContextSetupResponseIEs.List = append(initialContextSetupResponseIEs.List, ie)
	}

	// PDU Session Resource Failed to Setup List (optional)
	if failedList != nil && len(failedList.List) > 0 {
		ie = ngapType.InitialContextSetupResponseIEs{}
		ie.Id.Value = ngapType.ProtocolIEIDPDUSessionResourceFailedToSetupListCxtRes
		ie.Criticality.Value = ngapType.CriticalityPresentIgnore
		ie.Value.Present = ngapType.InitialContextSetupResponseIEsPresentPDUSessionResourceFailedToSetupListCxtRes
		ie.Value.PDUSessionResourceFailedToSetupListCxtRes = failedList
		initialContextSetupResponseIEs.List = append(initialContextSetupResponseIEs.List, ie)
	}

	// Criticality Diagnostics (optional)
	if criticalityDiagnostics != nil {
		ie = ngapType.InitialContextSetupResponseIEs{}
		ie.Id.Value = ngapType.ProtocolIEIDCriticalityDiagnostics
		ie.Criticality.Value = ngapType.CriticalityPresentIgnore
		ie.Value.CriticalityDiagnostics = criticalityDiagnostics
		initialContextSetupResponseIEs.List = append(initialContextSetupResponseIEs.List, ie)
	}

	return ngap.Encoder(pdu)
}

func BuildInitialContextSetupFailure(
	ue *n3iwf_context.N3IWFUe,
	cause ngapType.Cause,
	failedList *ngapType.PDUSessionResourceFailedToSetupListCxtFail,
	criticalityDiagnostics *ngapType.CriticalityDiagnostics) ([]byte, error) {

	var pdu ngapType.NGAPPDU
	pdu.Present = ngapType.NGAPPDUPresentUnsuccessfulOutcome
	pdu.UnsuccessfulOutcome = new(ngapType.UnsuccessfulOutcome)

	unsuccessfulOutcome := pdu.UnsuccessfulOutcome
	unsuccessfulOutcome.ProcedureCode.Value = ngapType.ProcedureCodeInitialContextSetup
	unsuccessfulOutcome.Criticality.Value = ngapType.CriticalityPresentReject

	unsuccessfulOutcome.Value.Present = ngapType.UnsuccessfulOutcomePresentInitialContextSetupFailure
	unsuccessfulOutcome.Value.InitialContextSetupFailure = new(ngapType.InitialContextSetupFailure)

	initialContextSetupFailure := unsuccessfulOutcome.Value.InitialContextSetupFailure
	initialContextSetupFailureIEs := &initialContextSetupFailure.ProtocolIEs

	// AMF UE NGAP ID
	ie := ngapType.InitialContextSetupFailureIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDAMFUENGAPID
	ie.Criticality.Value = ngapType.CriticalityPresentIgnore
	ie.Value.Present = ngapType.InitialContextSetupFailureIEsPresentAMFUENGAPID
	ie.Value.AMFUENGAPID = new(ngapType.AMFUENGAPID)

	aMFUENGAPID := ie.Value.AMFUENGAPID
	aMFUENGAPID.Value = ue.AmfUeNgapId

	initialContextSetupFailureIEs.List = append(initialContextSetupFailureIEs.List, ie)

	// RAN UE NGAP ID
	ie = ngapType.InitialContextSetupFailureIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDRANUENGAPID
	ie.Criticality.Value = ngapType.CriticalityPresentIgnore
	ie.Value.Present = ngapType.InitialContextSetupFailureIEsPresentRANUENGAPID
	ie.Value.RANUENGAPID = new(ngapType.RANUENGAPID)

	rANUENGAPID := ie.Value.RANUENGAPID
	rANUENGAPID.Value = ue.RanUeNgapId

	initialContextSetupFailureIEs.List = append(initialContextSetupFailureIEs.List, ie)

	// PDU Session Resource Failed to Setup List
	if failedList != nil && len(failedList.List) > 0 {
		ie = ngapType.InitialContextSetupFailureIEs{}
		ie.Id.Value = ngapType.ProtocolIEIDPDUSessionResourceFailedToSetupListCxtFail
		ie.Criticality.Value = ngapType.CriticalityPresentIgnore
		ie.Value.Present = ngapType.InitialContextSetupFailureIEsPresentPDUSessionResourceFailedToSetupListCxtFail
		ie.Value.PDUSessionResourceFailedToSetupListCxtFail = failedList
		initialContextSetupFailureIEs.List = append(initialContextSetupFailureIEs.List, ie)
	}

	// Cause
	ie = ngapType.InitialContextSetupFailureIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDCause
	ie.Criticality.Value = ngapType.CriticalityPresentIgnore
	ie.Value.Present = ngapType.InitialContextSetupFailureIEsPresentCause
	ie.Value.Cause = &cause
	initialContextSetupFailureIEs.List = append(initialContextSetupFailureIEs.List, ie)

	// Criticality Diagnostics (optional)
	if criticalityDiagnostics != nil {
		ie = ngapType.InitialContextSetupFailureIEs{}
		ie.Id.Value = ngapType.ProtocolIEIDCriticalityDiagnostics
		ie.Criticality.Value = ngapType.CriticalityPresentIgnore
		ie.Value.CriticalityDiagnostics = criticalityDiagnostics
		initialContextSetupFailureIEs.List = append(initialContextSetupFailureIEs.List, ie)
	}

	return ngap.Encoder(pdu)
}

func BuildUEContextModificationResponse(ue *n3iwf_context.N3IWFUe, criticalityDiagnostics *ngapType.CriticalityDiagnostics) ([]byte, error) {
	var pdu ngapType.NGAPPDU
	pdu.Present = ngapType.NGAPPDUPresentSuccessfulOutcome
	pdu.SuccessfulOutcome = new(ngapType.SuccessfulOutcome)

	successfulOutcome := pdu.SuccessfulOutcome
	successfulOutcome.ProcedureCode.Value = ngapType.ProcedureCodeUEContextModification
	successfulOutcome.Criticality.Value = ngapType.CriticalityPresentReject

	successfulOutcome.Value.Present = ngapType.SuccessfulOutcomePresentUEContextModificationResponse
	successfulOutcome.Value.UEContextModificationResponse = new(ngapType.UEContextModificationResponse)

	uEContextModificationResponse := successfulOutcome.Value.UEContextModificationResponse
	uEContextModificationResponseIEs := &uEContextModificationResponse.ProtocolIEs

	// AMF UE NGAP ID
	ie := ngapType.UEContextModificationResponseIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDAMFUENGAPID
	ie.Criticality.Value = ngapType.CriticalityPresentIgnore
	ie.Value.Present = ngapType.UEContextModificationResponseIEsPresentAMFUENGAPID
	ie.Value.AMFUENGAPID = new(ngapType.AMFUENGAPID)

	aMFUENGAPID := ie.Value.AMFUENGAPID
	aMFUENGAPID.Value = ue.AmfUeNgapId

	uEContextModificationResponseIEs.List = append(uEContextModificationResponseIEs.List, ie)

	// RAN UE NGAP ID
	ie = ngapType.UEContextModificationResponseIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDRANUENGAPID
	ie.Criticality.Value = ngapType.CriticalityPresentIgnore
	ie.Value.Present = ngapType.UEContextModificationResponseIEsPresentRANUENGAPID
	ie.Value.RANUENGAPID = new(ngapType.RANUENGAPID)

	rANUENGAPID := ie.Value.RANUENGAPID
	rANUENGAPID.Value = ue.RanUeNgapId

	uEContextModificationResponseIEs.List = append(uEContextModificationResponseIEs.List, ie)

	// Criticality Diagnostics (optional)
	ie = ngapType.UEContextModificationResponseIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDCriticalityDiagnostics
	ie.Criticality.Value = ngapType.CriticalityPresentIgnore
	ie.Value.CriticalityDiagnostics = criticalityDiagnostics
	uEContextModificationResponseIEs.List = append(uEContextModificationResponseIEs.List, ie)

	return ngap.Encoder(pdu)
}

func BuildUEContextModificationFailure(ue *n3iwf_context.N3IWFUe, cause ngapType.Cause, criticalityDiagnostics *ngapType.CriticalityDiagnostics) ([]byte, error) {
	var pdu ngapType.NGAPPDU
	pdu.Present = ngapType.NGAPPDUPresentUnsuccessfulOutcome
	pdu.UnsuccessfulOutcome = new(ngapType.UnsuccessfulOutcome)

	unsuccessfulOutcome := pdu.UnsuccessfulOutcome
	unsuccessfulOutcome.ProcedureCode.Value = ngapType.ProcedureCodeUEContextModification
	unsuccessfulOutcome.Criticality.Value = ngapType.CriticalityPresentReject

	unsuccessfulOutcome.Value.Present = ngapType.UnsuccessfulOutcomePresentUEContextModificationFailure
	unsuccessfulOutcome.Value.UEContextModificationFailure = new(ngapType.UEContextModificationFailure)

	uEContextModificationFailure := unsuccessfulOutcome.Value.UEContextModificationFailure
	uEContextModificationFailureIEs := &uEContextModificationFailure.ProtocolIEs

	// AMF UE NGAP ID
	ie := ngapType.UEContextModificationFailureIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDAMFUENGAPID
	ie.Criticality.Value = ngapType.CriticalityPresentIgnore
	ie.Value.Present = ngapType.UEContextModificationFailureIEsPresentAMFUENGAPID
	ie.Value.AMFUENGAPID = new(ngapType.AMFUENGAPID)

	aMFUENGAPID := ie.Value.AMFUENGAPID
	aMFUENGAPID.Value = ue.AmfUeNgapId

	uEContextModificationFailureIEs.List = append(uEContextModificationFailureIEs.List, ie)

	// RAN UE NGAP ID
	ie = ngapType.UEContextModificationFailureIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDRANUENGAPID
	ie.Criticality.Value = ngapType.CriticalityPresentIgnore
	ie.Value.Present = ngapType.UEContextModificationFailureIEsPresentRANUENGAPID
	ie.Value.RANUENGAPID = new(ngapType.RANUENGAPID)

	rANUENGAPID := ie.Value.RANUENGAPID
	rANUENGAPID.Value = ue.RanUeNgapId

	uEContextModificationFailureIEs.List = append(uEContextModificationFailureIEs.List, ie)

	// Cause
	ie = ngapType.UEContextModificationFailureIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDCause
	ie.Criticality.Value = ngapType.CriticalityPresentIgnore
	ie.Value.Present = ngapType.UEContextModificationFailureIEsPresentCause
	ie.Value.Cause = &cause
	uEContextModificationFailureIEs.List = append(uEContextModificationFailureIEs.List, ie)

	// Criticality Diagnostics (optional)
	ie = ngapType.UEContextModificationFailureIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDCriticalityDiagnostics
	ie.Criticality.Value = ngapType.CriticalityPresentIgnore
	ie.Value.CriticalityDiagnostics = criticalityDiagnostics
	uEContextModificationFailureIEs.List = append(uEContextModificationFailureIEs.List, ie)

	return ngap.Encoder(pdu)
}

func BuildUEContextReleaseComplete(ue *n3iwf_context.N3IWFUe, criticalityDiagnostics *ngapType.CriticalityDiagnostics) ([]byte, error) {
	var pdu ngapType.NGAPPDU
	pdu.Present = ngapType.NGAPPDUPresentSuccessfulOutcome
	pdu.SuccessfulOutcome = new(ngapType.SuccessfulOutcome)

	successfulOutcome := pdu.SuccessfulOutcome
	successfulOutcome.ProcedureCode.Value = ngapType.ProcedureCodeUEContextRelease
	successfulOutcome.Criticality.Value = ngapType.CriticalityPresentReject

	successfulOutcome.Value.Present = ngapType.SuccessfulOutcomePresentUEContextReleaseComplete
	successfulOutcome.Value.UEContextReleaseComplete = new(ngapType.UEContextReleaseComplete)

	uEContextReleaseComplete := successfulOutcome.Value.UEContextReleaseComplete
	uEContextReleaseCompleteIEs := &uEContextReleaseComplete.ProtocolIEs

	// AMF UE NGAP ID
	ie := ngapType.UEContextReleaseCompleteIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDAMFUENGAPID
	ie.Criticality.Value = ngapType.CriticalityPresentIgnore
	ie.Value.Present = ngapType.UEContextReleaseCompleteIEsPresentAMFUENGAPID
	ie.Value.AMFUENGAPID = new(ngapType.AMFUENGAPID)

	aMFUENGAPID := ie.Value.AMFUENGAPID
	aMFUENGAPID.Value = ue.AmfUeNgapId

	uEContextReleaseCompleteIEs.List = append(uEContextReleaseCompleteIEs.List, ie)

	// RAN UE NGAP ID
	ie = ngapType.UEContextReleaseCompleteIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDRANUENGAPID
	ie.Criticality.Value = ngapType.CriticalityPresentIgnore
	ie.Value.Present = ngapType.UEContextReleaseCompleteIEsPresentRANUENGAPID
	ie.Value.RANUENGAPID = new(ngapType.RANUENGAPID)

	rANUENGAPID := ie.Value.RANUENGAPID
	rANUENGAPID.Value = ue.RanUeNgapId

	uEContextReleaseCompleteIEs.List = append(uEContextReleaseCompleteIEs.List, ie)

	// User Location Information (optional)
	ie = ngapType.UEContextReleaseCompleteIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDUserLocationInformation
	ie.Criticality.Value = ngapType.CriticalityPresentIgnore
	ie.Value.Present = ngapType.UEContextReleaseCompleteIEsPresentUserLocationInformation
	ie.Value.UserLocationInformation = new(ngapType.UserLocationInformation)

	userLocationInformation := ie.Value.UserLocationInformation
	userLocationInformation.Present = ngapType.UserLocationInformationPresentUserLocationInformationN3IWF
	userLocationInformation.UserLocationInformationN3IWF = new(ngapType.UserLocationInformationN3IWF)

	userLocationInfoN3IWF := userLocationInformation.UserLocationInformationN3IWF
	userLocationInfoN3IWF.IPAddress = ngapConvert.IPAddressToNgap(ue.IPAddrv4, ue.IPAddrv6)
	userLocationInfoN3IWF.PortNumber = ngapConvert.PortNumberToNgap(ue.PortNumber)

	uEContextReleaseCompleteIEs.List = append(uEContextReleaseCompleteIEs.List, ie)

	// PDU Session Resource List
	ie = ngapType.UEContextReleaseCompleteIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDPDUSessionResourceListCxtRelCpl
	ie.Criticality.Value = ngapType.CriticalityPresentReject
	ie.Value.Present = ngapType.UEContextReleaseCompleteIEsPresentPDUSessionResourceListCxtRelCpl
	ie.Value.PDUSessionResourceListCxtRelCpl = new(ngapType.PDUSessionResourceListCxtRelCpl)

	pDUSessionResourceListCxtRelCpl := ie.Value.PDUSessionResourceListCxtRelCpl

	// PDU Session Resource Item (in PDU Session Resource List)
	for _, pduSession := range ue.PduSessionList {
		pDUSessionResourceItemCxtRelCpl := ngapType.PDUSessionResourceItemCxtRelCpl{}
		pDUSessionResourceItemCxtRelCpl.PDUSessionID.Value = pduSession.Id
		pDUSessionResourceListCxtRelCpl.List = append(pDUSessionResourceListCxtRelCpl.List, pDUSessionResourceItemCxtRelCpl)
	}

	uEContextReleaseCompleteIEs.List = append(uEContextReleaseCompleteIEs.List, ie)

	// Criticality Diagnostics (optional)
	if criticalityDiagnostics != nil {
		ie = ngapType.UEContextReleaseCompleteIEs{}
		ie.Id.Value = ngapType.ProtocolIEIDCriticalityDiagnostics
		ie.Criticality.Value = ngapType.CriticalityPresentIgnore
		ie.Value.CriticalityDiagnostics = criticalityDiagnostics
		uEContextReleaseCompleteIEs.List = append(uEContextReleaseCompleteIEs.List, ie)
	}

	return ngap.Encoder(pdu)
}

func BuildUEContextReleaseRequest(ue *n3iwf_context.N3IWFUe, cause ngapType.Cause) ([]byte, error) {
	var pdu ngapType.NGAPPDU
	pdu.Present = ngapType.NGAPPDUPresentInitiatingMessage
	pdu.InitiatingMessage = new(ngapType.InitiatingMessage)

	initiatingMessage := pdu.InitiatingMessage
	initiatingMessage.ProcedureCode.Value = ngapType.ProcedureCodeUEContextReleaseRequest
	initiatingMessage.Criticality.Value = ngapType.CriticalityPresentIgnore

	initiatingMessage.Value.Present = ngapType.InitiatingMessagePresentUEContextReleaseRequest
	initiatingMessage.Value.UEContextReleaseRequest = new(ngapType.UEContextReleaseRequest)

	uEContextReleaseRequest := initiatingMessage.Value.UEContextReleaseRequest
	uEContextReleaseRequestIEs := &uEContextReleaseRequest.ProtocolIEs

	// AMF UE NGAP ID
	ie := ngapType.UEContextReleaseRequestIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDAMFUENGAPID
	ie.Criticality.Value = ngapType.CriticalityPresentReject
	ie.Value.Present = ngapType.UEContextReleaseRequestIEsPresentAMFUENGAPID
	ie.Value.AMFUENGAPID = new(ngapType.AMFUENGAPID)

	aMFUENGAPID := ie.Value.AMFUENGAPID
	aMFUENGAPID.Value = ue.AmfUeNgapId

	uEContextReleaseRequestIEs.List = append(uEContextReleaseRequestIEs.List, ie)

	// RAN UE NGAP ID
	ie = ngapType.UEContextReleaseRequestIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDRANUENGAPID
	ie.Criticality.Value = ngapType.CriticalityPresentReject
	ie.Value.Present = ngapType.UEContextReleaseRequestIEsPresentRANUENGAPID
	ie.Value.RANUENGAPID = new(ngapType.RANUENGAPID)

	rANUENGAPID := ie.Value.RANUENGAPID
	rANUENGAPID.Value = ue.RanUeNgapId

	uEContextReleaseRequestIEs.List = append(uEContextReleaseRequestIEs.List, ie)

	// PDU Session Resource List
	ie = ngapType.UEContextReleaseRequestIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDPDUSessionResourceListCxtRelReq
	ie.Criticality.Value = ngapType.CriticalityPresentReject
	ie.Value.Present = ngapType.UEContextReleaseRequestIEsPresentPDUSessionResourceListCxtRelReq
	ie.Value.PDUSessionResourceListCxtRelReq = new(ngapType.PDUSessionResourceListCxtRelReq)

	pDUSessionResourceListCxtRelReq := ie.Value.PDUSessionResourceListCxtRelReq

	// PDU Session Resource Item in PDU session Resource List
	for _, pduSession := range ue.PduSessionList {
		pDUSessionResourceItem := ngapType.PDUSessionResourceItemCxtRelReq{}
		pDUSessionResourceItem.PDUSessionID.Value = pduSession.Id
		pDUSessionResourceListCxtRelReq.List = append(pDUSessionResourceListCxtRelReq.List, pDUSessionResourceItem)
	}
	uEContextReleaseRequestIEs.List = append(uEContextReleaseRequestIEs.List, ie)

	// Cause
	ie = ngapType.UEContextReleaseRequestIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDCause
	ie.Criticality.Value = ngapType.CriticalityPresentIgnore
	ie.Value.Present = ngapType.UEContextReleaseRequestIEsPresentCause
	ie.Value.Cause = &cause
	uEContextReleaseRequestIEs.List = append(uEContextReleaseRequestIEs.List, ie)

	return ngap.Encoder(pdu)
}

func BuildInitialUEMessage() ([]byte, error) {
	var pdu ngapType.NGAPPDU
	return ngap.Encoder(pdu)
}

func BuildUplinkNASTransport(ue *n3iwf_context.N3IWFUe, nasPdu []byte) ([]byte, error) {
	var pdu ngapType.NGAPPDU
	pdu.Present = ngapType.NGAPPDUPresentInitiatingMessage
	pdu.InitiatingMessage = new(ngapType.InitiatingMessage)

	initiatingMessage := pdu.InitiatingMessage
	initiatingMessage.ProcedureCode.Value = ngapType.ProcedureCodeUplinkNASTransport
	initiatingMessage.Criticality.Value = ngapType.CriticalityPresentIgnore

	initiatingMessage.Value.Present = ngapType.InitiatingMessagePresentUplinkNASTransport
	initiatingMessage.Value.UplinkNASTransport = new(ngapType.UplinkNASTransport)

	uplinkNasTransport := initiatingMessage.Value.UplinkNASTransport
	uplinkNasTransportIEs := &uplinkNasTransport.ProtocolIEs

	// AMF UE NGAP ID
	ie := ngapType.UplinkNASTransportIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDAMFUENGAPID
	ie.Criticality.Value = ngapType.CriticalityPresentReject
	ie.Value.Present = ngapType.UplinkNASTransportIEsPresentAMFUENGAPID
	ie.Value.AMFUENGAPID = new(ngapType.AMFUENGAPID)

	aMFUENGAPID := ie.Value.AMFUENGAPID
	aMFUENGAPID.Value = ue.AmfUeNgapId

	uplinkNasTransportIEs.List = append(uplinkNasTransportIEs.List, ie)

	// RAN UE NGAP ID
	ie = ngapType.UplinkNASTransportIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDRANUENGAPID
	ie.Criticality.Value = ngapType.CriticalityPresentReject
	ie.Value.Present = ngapType.UplinkNASTransportIEsPresentRANUENGAPID
	ie.Value.RANUENGAPID = new(ngapType.RANUENGAPID)

	rANUENGAPID := ie.Value.RANUENGAPID
	rANUENGAPID.Value = ue.RanUeNgapId

	uplinkNasTransportIEs.List = append(uplinkNasTransportIEs.List, ie)

	// NAS-PDU
	ie = ngapType.UplinkNASTransportIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDNASPDU
	ie.Criticality.Value = ngapType.CriticalityPresentReject
	ie.Value.Present = ngapType.UplinkNASTransportIEsPresentNASPDU
	ie.Value.NASPDU = new(ngapType.NASPDU)
	nASPDU := ie.Value.NASPDU
	nASPDU.Value = nasPdu
	uplinkNasTransportIEs.List = append(uplinkNasTransportIEs.List, ie)

	// User Location Information
	ie = ngapType.UplinkNASTransportIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDUserLocationInformation
	ie.Criticality.Value = ngapType.CriticalityPresentIgnore
	ie.Value.Present = ngapType.UplinkNASTransportIEsPresentUserLocationInformation
	ie.Value.UserLocationInformation = new(ngapType.UserLocationInformation)

	userLocationInformation := ie.Value.UserLocationInformation
	userLocationInformation.Present = ngapType.UserLocationInformationPresentUserLocationInformationN3IWF
	userLocationInformation.UserLocationInformationN3IWF = new(ngapType.UserLocationInformationN3IWF)
	userLocationInformationN3IWF := userLocationInformation.UserLocationInformationN3IWF
	userLocationInformationN3IWF.IPAddress = ngapConvert.IPAddressToNgap(ue.IPAddrv4, ue.IPAddrv6)
	userLocationInformationN3IWF.PortNumber = ngapConvert.PortNumberToNgap(ue.PortNumber)

	uplinkNasTransportIEs.List = append(uplinkNasTransportIEs.List, ie)

	return ngap.Encoder(pdu)
}

func BuildNASNonDeliveryIndication() ([]byte, error) {
	var pdu ngapType.NGAPPDU
	return ngap.Encoder(pdu)
}

func BuildRerouteNASRequest() ([]byte, error) {
	var pdu ngapType.NGAPPDU
	return ngap.Encoder(pdu)
}

func BuildPDUSessionResourceSetupResponse() ([]byte, error) {
	var pdu ngapType.NGAPPDU
	return ngap.Encoder(pdu)
}

func BuildPDUSessionResourceModifyResponse() ([]byte, error) {
	var pdu ngapType.NGAPPDU
	return ngap.Encoder(pdu)
}

func BuildPDUSessionResourceModifyIndication() ([]byte, error) {
	var pdu ngapType.NGAPPDU
	return ngap.Encoder(pdu)
}

func BuildPDUSessionResourceNotify() ([]byte, error) {
	var pdu ngapType.NGAPPDU
	return ngap.Encoder(pdu)
}

func BuildPDUSessionResourceReleaseResponse() ([]byte, error) {
	var pdu ngapType.NGAPPDU
	return ngap.Encoder(pdu)
}

func BuildErrorIndication() ([]byte, error) {
	var pdu ngapType.NGAPPDU
	return ngap.Encoder(pdu)
}

func BuildUERadioCapabilityInfoIndication() ([]byte, error) {
	var pdu ngapType.NGAPPDU
	return ngap.Encoder(pdu)
}

func BuildUERadioCapabilityCheckResponse() ([]byte, error) {
	var pdu ngapType.NGAPPDU
	return ngap.Encoder(pdu)
}

func BuildAMFConfigurationUpdateAcknowledge() ([]byte, error) {
	var pdu ngapType.NGAPPDU
	return ngap.Encoder(pdu)
}

func BuildAMFConfigurationUpdateFailure() ([]byte, error) {
	var pdu ngapType.NGAPPDU
	return ngap.Encoder(pdu)
}

func BuildRANConfigurationUpdate() ([]byte, error) {
	var pdu ngapType.NGAPPDU
	return ngap.Encoder(pdu)
}

func BuildUplinkRANConfigurationTransfer() ([]byte, error) {
	var pdu ngapType.NGAPPDU
	return ngap.Encoder(pdu)
}

func BuildUplinkRANStatusTransfer() ([]byte, error) {
	var pdu ngapType.NGAPPDU
	return ngap.Encoder(pdu)
}

func BuildLocationReportingFailureIndication() ([]byte, error) {
	var pdu ngapType.NGAPPDU
	return ngap.Encoder(pdu)
}

func BuildLocationReport() ([]byte, error) {
	var pdu ngapType.NGAPPDU
	return ngap.Encoder(pdu)
}

func BuildRRCInactiveTransitionReport() ([]byte, error) {
	var pdu ngapType.NGAPPDU
	return ngap.Encoder(pdu)
}

func BuildPDUSessionResourceSetupResponseTransfer(pduSession *n3iwf_context.PDUSession) ([]byte, error) {

	transfer := ngapType.PDUSessionResourceSetupResponseTransfer{}

	// TODO: use tunnel info allocated by n3iwf
	// QOS Flow Per TNL Information
	qosFlowPerTNLInformation := &transfer.QosFlowPerTNLInformation
	qosFlowPerTNLInformation.UPTransportLayerInformation.Present = ngapType.UPTransportLayerInformationPresentGTPTunnel
	qosFlowPerTNLInformation.UPTransportLayerInformation.GTPTunnel = new(ngapType.GTPTunnel)
	qosFlowPerTNLInformation.UPTransportLayerInformation.GTPTunnel.GTPTEID.Value = aper.OctetString(pduSession.TEID)
	qosFlowPerTNLInformation.UPTransportLayerInformation.GTPTunnel.TransportLayerAddress = ngapConvert.IPAddressToNgap(pduSession.GTPEndpointIPv4, pduSession.GTPEndpointIPv6)

	return aper.MarshalWithParams(transfer, "valueExt")
}

func BuildPDUSessionResourceSetupUnsuccessfulTransfer(cause ngapType.Cause, criticalityDiagnostics *ngapType.CriticalityDiagnostics) ([]byte, error) {

	transfer := ngapType.PDUSessionResourceSetupUnsuccessfulTransfer{}

	// Cause
	transfer.Cause = cause

	// Criticality Diagnostics (optional)
	if criticalityDiagnostics != nil {
		transfer.CriticalityDiagnostics = criticalityDiagnostics
	}

	return aper.MarshalWithParams(transfer, "valueExt")
}
