package ngap_handler

import (
	"encoding/hex"
	"free5gc/lib/aper"
	"free5gc/lib/nas/nasMessage"
	"free5gc/lib/ngap/ngapConvert"
	"free5gc/lib/ngap/ngapType"
	"free5gc/lib/openapi/models"
	"free5gc/src/amf/amf_consumer"
	"free5gc/src/amf/amf_context"
	"free5gc/src/amf/amf_nas"
	"free5gc/src/amf/amf_ngap/ngap_message"
	"free5gc/src/amf/gmm"
	"free5gc/src/amf/gmm/gmm_message"
	"free5gc/src/amf/gmm/gmm_state"
	"free5gc/src/amf/logger"

	"github.com/sirupsen/logrus"
)

var Ngaplog *logrus.Entry

func init() {
	Ngaplog = logger.NgapLog
}

func HandleNGSetupRequest(ran *amf_context.AmfRan, message *ngapType.NGAPPDU) {
	var globalRANNodeID *ngapType.GlobalRANNodeID
	var rANNodeName *ngapType.RANNodeName
	var supportedTAList *ngapType.SupportedTAList
	var pagingDRX *ngapType.PagingDRX

	var cause ngapType.Cause

	if ran == nil {
		logger.NgapLog.Error("ran is nil")
		return
	}
	if message == nil {
		logger.NgapLog.Error("NGAP Message is nil")
		return
	}
	initiatingMessage := message.InitiatingMessage
	if initiatingMessage == nil {
		logger.NgapLog.Error("Initiating Message is nil")
		return
	}
	nGSetupRequest := initiatingMessage.Value.NGSetupRequest
	if nGSetupRequest == nil {
		logger.NgapLog.Error("NGSetupRequest is nil")
		return
	}
	logger.NgapLog.Info("[AMF] NG Setup request")
	for i := 0; i < len(nGSetupRequest.ProtocolIEs.List); i++ {
		ie := nGSetupRequest.ProtocolIEs.List[i]
		switch ie.Id.Value {
		case ngapType.ProtocolIEIDGlobalRANNodeID:
			globalRANNodeID = ie.Value.GlobalRANNodeID
			logger.NgapLog.Trace("[NGAP] Decode IE GlobalRANNodeID")
			if globalRANNodeID == nil {
				logger.NgapLog.Error("GlobalRANNodeID is nil")
				return
			}
		case ngapType.ProtocolIEIDSupportedTAList:
			supportedTAList = ie.Value.SupportedTAList
			logger.NgapLog.Trace("[NGAP] Decode IE SupportedTAList")
			if supportedTAList == nil {
				logger.NgapLog.Error("SupportedTAList is nil")
				return
			}
		case ngapType.ProtocolIEIDRANNodeName:
			rANNodeName = ie.Value.RANNodeName
			logger.NgapLog.Trace("[NGAP] Decode IE RANNodeName")
			if rANNodeName == nil {
				logger.NgapLog.Error("RANNodeName is nil")
				return
			}
		case ngapType.ProtocolIEIDDefaultPagingDRX:
			pagingDRX = ie.Value.DefaultPagingDRX
			logger.NgapLog.Trace("[NGAP] Decode IE DefaultPagingDRX")
			if pagingDRX == nil {
				logger.NgapLog.Error("DefaultPagingDRX is nil")
				return
			}
		}
	}

	ran.SetRanId(globalRANNodeID)
	if rANNodeName != nil {
		ran.Name = rANNodeName.Value
	}
	if pagingDRX != nil {
		logger.NgapLog.Tracef("PagingDRX[%d]", pagingDRX.Value)
	}

	for i := 0; i < len(supportedTAList.List); i++ {
		supportedTAItem := supportedTAList.List[i]
		tac := hex.EncodeToString(supportedTAItem.TAC.Value)
		capOfSupportTai := cap(ran.SupportedTAList)
		for j := 0; j < len(supportedTAItem.BroadcastPLMNList.List); j++ {
			supportedTAI := amf_context.NewSupportedTAI()
			supportedTAI.Tai.Tac = tac
			broadcastPLMNItem := supportedTAItem.BroadcastPLMNList.List[j]
			plmnId := ngapConvert.PlmnIdToModels(broadcastPLMNItem.PLMNIdentity)
			supportedTAI.Tai.PlmnId = &plmnId
			capOfSNssaiList := cap(supportedTAI.SNssaiList)
			for k := 0; k < len(broadcastPLMNItem.TAISliceSupportList.List); k++ {
				tAISliceSupportItem := broadcastPLMNItem.TAISliceSupportList.List[k]
				if len(supportedTAI.SNssaiList) < capOfSNssaiList {
					supportedTAI.SNssaiList = append(supportedTAI.SNssaiList, ngapConvert.SNssaiToModels(tAISliceSupportItem.SNSSAI))
				} else {
					break
				}
			}
			logger.NgapLog.Tracef("PLMN_ID[MCC:%s MNC:%s] TAC[%s]", plmnId.Mcc, plmnId.Mnc, tac)
			if len(ran.SupportedTAList) < capOfSupportTai {
				ran.SupportedTAList = append(ran.SupportedTAList, supportedTAI)

			} else {
				break
			}
		}

	}

	if len(ran.SupportedTAList) == 0 {
		logger.NgapLog.Warn("NG-Setup failure: No supported TA exist in NG-Setup request")
		cause.Present = ngapType.CausePresentMisc
		cause.Misc = &ngapType.CauseMisc{
			Value: ngapType.CauseMiscPresentUnspecified,
		}
	} else {
		var found bool
		for i, tai := range ran.SupportedTAList {
			if amf_context.InTaiList(tai.Tai, amf_context.AMF_Self().SupportTaiLists) {
				logger.NgapLog.Tracef("SERVED_TAI_INDEX[%d]", i)
				found = true
				break
			}
		}
		if !found {
			logger.NgapLog.Warn("NG-Setup failure: Cannot find Served TAI in AMF")
			cause.Present = ngapType.CausePresentMisc
			cause.Misc = &ngapType.CauseMisc{
				Value: ngapType.CauseMiscPresentUnknownPLMN,
			}
		}
	}

	if cause.Present == ngapType.CausePresentNothing {
		ngap_message.SendNGSetupResponse(ran)
	} else {
		ngap_message.SendNGSetupFailure(ran, cause)
	}
}

func HandleUplinkNasTransport(ran *amf_context.AmfRan, message *ngapType.NGAPPDU) {

	var aMFUENGAPID *ngapType.AMFUENGAPID
	var rANUENGAPID *ngapType.RANUENGAPID
	var nASPDU *ngapType.NASPDU
	var userLocationInformation *ngapType.UserLocationInformation

	if ran == nil {
		logger.NgapLog.Error("ran is nil")
		return
	}
	if message == nil {
		logger.NgapLog.Error("NGAP Message is nil")
		return
	}

	initiatingMessage := message.InitiatingMessage
	if initiatingMessage == nil {
		logger.NgapLog.Error("Initiating Message is nil")
		return
	}

	uplinkNasTransport := initiatingMessage.Value.UplinkNASTransport
	if uplinkNasTransport == nil {
		logger.NgapLog.Error("UplinkNasTransport is nil")
		return
	}
	logger.NgapLog.Info("[AMF] Uplink Nas Transport")

	for i := 0; i < len(uplinkNasTransport.ProtocolIEs.List); i++ {
		ie := uplinkNasTransport.ProtocolIEs.List[i]
		switch ie.Id.Value {
		case ngapType.ProtocolIEIDAMFUENGAPID:
			aMFUENGAPID = ie.Value.AMFUENGAPID
			logger.NgapLog.Trace("[NGAP] Decode IE AmfUeNgapID")
			if aMFUENGAPID == nil {
				logger.NgapLog.Error("AmfUeNgapID is nil")
				return
			}
		case ngapType.ProtocolIEIDRANUENGAPID:
			rANUENGAPID = ie.Value.RANUENGAPID
			logger.NgapLog.Trace("[NGAP] Decode IE RanUeNgapID")
			if rANUENGAPID == nil {
				logger.NgapLog.Error("RanUeNgapID is nil")
				return
			}
		case ngapType.ProtocolIEIDNASPDU:
			nASPDU = ie.Value.NASPDU
			logger.NgapLog.Trace("[NGAP] Decode IE NasPdu")
			if nASPDU == nil {
				logger.NgapLog.Error("nASPDU is nil")
				return
			}
		case ngapType.ProtocolIEIDUserLocationInformation:
			userLocationInformation = ie.Value.UserLocationInformation
			logger.NgapLog.Trace("[NGAP] Decode IE UserLocationInformation")
			if userLocationInformation == nil {
				logger.NgapLog.Error("UserLocationInformation is nil")
				return
			}
		}
	}

	printRanInfo(ran)

	ranUe := ran.RanUeFindByRanUeNgapID(rANUENGAPID.Value)
	if ranUe == nil {
		logger.NgapLog.Errorf("No UE Context[RanUeNgapID: %d]", rANUENGAPID.Value)
		return
	}
	amfUe := ranUe.AmfUe
	if amfUe == nil {
		err := ranUe.Remove()
		if err != nil {
			logger.NgapLog.Errorf(err.Error())
		}
		logger.NgapLog.Errorf("No UE Context of RanUe with RANUENGAPID[%d] AMFUENGAPID[%d] ", rANUENGAPID.Value, aMFUENGAPID.Value)
		return
	}

	logger.NgapLog.Tracef("RANUENGAPID[%d] AMFUENGAPID[%d]", ranUe.RanUeNgapId, ranUe.AmfUeNgapId)

	if userLocationInformation != nil {
		ranUe.UpdateLocation(userLocationInformation)
	}

	amf_nas.HandleNAS(ranUe, ngapType.ProcedureCodeUplinkNASTransport, nASPDU.Value)
}

func HandleNGReset(ran *amf_context.AmfRan, message *ngapType.NGAPPDU) {

	var cause *ngapType.Cause
	var resetType *ngapType.ResetType

	if ran == nil {
		logger.NgapLog.Error("ran is nil")
		return
	}
	if message == nil {
		logger.NgapLog.Error("NGAP Message is nil")
		return
	}
	initiatingMessage := message.InitiatingMessage
	if initiatingMessage == nil {
		logger.NgapLog.Error("Initiating Message is nil")
		return
	}
	nGReset := initiatingMessage.Value.NGReset
	if nGReset == nil {
		logger.NgapLog.Error("NGReset is nil")
		return
	}

	logger.NgapLog.Info("[AMF] NG Reset")

	for _, ie := range nGReset.ProtocolIEs.List {
		switch ie.Id.Value {
		case ngapType.ProtocolIEIDCause:
			cause = ie.Value.Cause
			logger.NgapLog.Trace("[NGAP] Decode IE Cause")
			if cause == nil {
				logger.NgapLog.Error("Cause is nil")
				return
			}
		case ngapType.ProtocolIEIDResetType:
			resetType = ie.Value.ResetType
			logger.NgapLog.Trace("[NGAP] Decode IE ResetType")
			if resetType == nil {
				logger.NgapLog.Error("ResetType is nil")
				return
			}
		}
	}

	printRanInfo(ran)

	printAndGetCause(cause)

	switch resetType.Present {
	case ngapType.ResetTypePresentNGInterface:
		logger.NgapLog.Trace("ResetType Present: NG Interface")
		ran.RemoveAllUeInRan()
		ngap_message.SendNGResetAcknowledge(ran, nil, nil)
	case ngapType.ResetTypePresentPartOfNGInterface:
		logger.NgapLog.Trace("ResetType Present: Part of NG Interface")

		partOfNGInterface := resetType.PartOfNGInterface
		if partOfNGInterface == nil {
			logger.NgapLog.Error("PartOfNGInterface is nil")
			return
		}

		var ranUe *amf_context.RanUe

		for _, ueAssociatedLogicalNGConnectionItem := range partOfNGInterface.List {
			if ueAssociatedLogicalNGConnectionItem.AMFUENGAPID != nil {
				logger.NgapLog.Tracef("AmfUeNgapID[%d]", ueAssociatedLogicalNGConnectionItem.AMFUENGAPID.Value)
				for _, ue := range ran.RanUeList {
					if ue.AmfUeNgapId == ueAssociatedLogicalNGConnectionItem.AMFUENGAPID.Value {
						ranUe = ue
						break
					}
				}
			} else if ueAssociatedLogicalNGConnectionItem.RANUENGAPID != nil {
				logger.NgapLog.Tracef("RanUeNgapID[%d]", ueAssociatedLogicalNGConnectionItem.RANUENGAPID.Value)
				for _, ue := range ran.RanUeList {
					if ue.RanUeNgapId == ueAssociatedLogicalNGConnectionItem.RANUENGAPID.Value {
						ranUe = ue
						break
					}
				}
			}

			if ranUe == nil {
				logger.NgapLog.Warn("Cannot not find UE Context")
				if ueAssociatedLogicalNGConnectionItem.AMFUENGAPID != nil {
					logger.NgapLog.Warnf("AmfUeNgapID[%d]", ueAssociatedLogicalNGConnectionItem.AMFUENGAPID.Value)
				}
				if ueAssociatedLogicalNGConnectionItem.RANUENGAPID != nil {
					logger.NgapLog.Warnf("RanUeNgapID[%d]", ueAssociatedLogicalNGConnectionItem.RANUENGAPID.Value)
				}
			}

			err := ranUe.Remove()
			if err != nil {
				logger.NgapLog.Error(err.Error())
			}
		}
		ngap_message.SendNGResetAcknowledge(ran, partOfNGInterface, nil)
	default:
		logger.NgapLog.Warnf("Invalid ResetType[%d]", resetType.Present)
	}
}

func HandleNGResetAcknowledge(ran *amf_context.AmfRan, message *ngapType.NGAPPDU) {

	var uEAssociatedLogicalNGConnectionList *ngapType.UEAssociatedLogicalNGConnectionList
	var criticalityDiagnostics *ngapType.CriticalityDiagnostics

	if ran == nil {
		logger.NgapLog.Error("ran is nil")
		return
	}
	if message == nil {
		logger.NgapLog.Error("NGAP Message is nil")
		return
	}
	successfulOutcome := message.SuccessfulOutcome
	if successfulOutcome == nil {
		logger.NgapLog.Error("SuccessfulOutcome is nil")
		return
	}
	nGResetAcknowledge := successfulOutcome.Value.NGResetAcknowledge
	if nGResetAcknowledge == nil {
		logger.NgapLog.Error("NGResetAcknowledge is nil")
		return
	}

	logger.NgapLog.Info("[AMF] NG Reset Acknowledge")

	for _, ie := range nGResetAcknowledge.ProtocolIEs.List {
		switch ie.Id.Value {
		case ngapType.ProtocolIEIDUEAssociatedLogicalNGConnectionList:
			uEAssociatedLogicalNGConnectionList = ie.Value.UEAssociatedLogicalNGConnectionList
		case ngapType.ProtocolIEIDCriticalityDiagnostics:
			criticalityDiagnostics = ie.Value.CriticalityDiagnostics
		}
	}

	printRanInfo(ran)

	if uEAssociatedLogicalNGConnectionList != nil {
		logger.NgapLog.Tracef("%d UE association(s) has been reset", len(uEAssociatedLogicalNGConnectionList.List))
		for i, item := range uEAssociatedLogicalNGConnectionList.List {
			if item.AMFUENGAPID != nil && item.RANUENGAPID != nil {
				logger.NgapLog.Tracef("%d: AmfUeNgapID[%d] RanUeNgapID[%d]", i+1, item.AMFUENGAPID.Value, item.RANUENGAPID.Value)
			} else if item.AMFUENGAPID != nil {
				logger.NgapLog.Tracef("%d: AmfUeNgapID[%d] RanUeNgapID[-1]", i+1, item.AMFUENGAPID.Value)
			} else if item.RANUENGAPID != nil {
				logger.NgapLog.Tracef("%d: AmfUeNgapID[-1] RanUeNgapID[%d]", i+1, item.RANUENGAPID.Value)
			}
		}
	}

	if criticalityDiagnostics != nil {
		printCriticalityDiagnostics(criticalityDiagnostics)
	}
}

func HandleUEContextReleaseComplete(ran *amf_context.AmfRan, message *ngapType.NGAPPDU) {

	var aMFUENGAPID *ngapType.AMFUENGAPID
	var rANUENGAPID *ngapType.RANUENGAPID
	var userLocationInformation *ngapType.UserLocationInformation
	var infoOnRecommendedCellsAndRANNodesForPaging *ngapType.InfoOnRecommendedCellsAndRANNodesForPaging
	var pDUSessionResourceList *ngapType.PDUSessionResourceListCxtRelCpl
	var criticalityDiagnostics *ngapType.CriticalityDiagnostics

	if ran == nil {
		logger.NgapLog.Error("ran is nil")
		return
	}
	if message == nil {
		logger.NgapLog.Error("NGAP Message is nil")
		return
	}
	successfulOutcome := message.SuccessfulOutcome
	if successfulOutcome == nil {
		logger.NgapLog.Error("SuccessfulOutcome is nil")
		return
	}
	uEContextReleaseComplete := successfulOutcome.Value.UEContextReleaseComplete
	if uEContextReleaseComplete == nil {
		logger.NgapLog.Error("NGResetAcknowledge is nil")
		return
	}

	logger.NgapLog.Info("[AMF] UE Context Release Complete")

	for _, ie := range uEContextReleaseComplete.ProtocolIEs.List {
		switch ie.Id.Value {
		case ngapType.ProtocolIEIDAMFUENGAPID:
			aMFUENGAPID = ie.Value.AMFUENGAPID
			logger.NgapLog.Trace("[NGAP] Decode IE AmfUeNgapID")
			if aMFUENGAPID == nil {
				logger.NgapLog.Error("AmfUeNgapID is nil")
				return
			}
		case ngapType.ProtocolIEIDRANUENGAPID:
			rANUENGAPID = ie.Value.RANUENGAPID
			logger.NgapLog.Trace("[NGAP] Decode IE RanUeNgapID")
			if rANUENGAPID == nil {
				logger.NgapLog.Error("RanUeNgapID is nil")
				return
			}
		case ngapType.ProtocolIEIDUserLocationInformation:
			userLocationInformation = ie.Value.UserLocationInformation
			logger.NgapLog.Trace("[NGAP] Decode IE UserLocationInformation")
		case ngapType.ProtocolIEIDInfoOnRecommendedCellsAndRANNodesForPaging:
			infoOnRecommendedCellsAndRANNodesForPaging = ie.Value.InfoOnRecommendedCellsAndRANNodesForPaging
			logger.NgapLog.Trace("[NGAP] Decode IE InfoOnRecommendedCellsAndRANNodesForPaging")
			if infoOnRecommendedCellsAndRANNodesForPaging != nil {
				logger.NgapLog.Warn("IE infoOnRecommendedCellsAndRANNodesForPaging is not support")
			}
		case ngapType.ProtocolIEIDPDUSessionResourceListCxtRelCpl:
			pDUSessionResourceList = ie.Value.PDUSessionResourceListCxtRelCpl
			logger.NgapLog.Trace("[NGAP] Decode IE PDUSessionResourceList")
		case ngapType.ProtocolIEIDCriticalityDiagnostics:
			criticalityDiagnostics = ie.Value.CriticalityDiagnostics
			logger.NgapLog.Trace("[NGAP] Decode IE CriticalityDiagnostics")
		}
	}

	printRanInfo(ran)

	ranUe := amf_context.AMF_Self().RanUeFindByAmfUeNgapID(aMFUENGAPID.Value)
	if ranUe == nil {
		logger.NgapLog.Errorf("No RanUe Context[AmfUeNgapID: %d]", aMFUENGAPID.Value)
		cause := ngapType.Cause{
			Present: ngapType.CausePresentRadioNetwork,
			RadioNetwork: &ngapType.CauseRadioNetwork{
				Value: ngapType.CauseRadioNetworkPresentUnknownLocalUENGAPID,
			},
		}
		ngap_message.SendErrorIndication(ran, nil, nil, &cause, nil)
		return
	}

	if userLocationInformation != nil {
		ranUe.UpdateLocation(userLocationInformation)
	}
	if criticalityDiagnostics != nil {
		printCriticalityDiagnostics(criticalityDiagnostics)
	}

	amfUe := ranUe.AmfUe
	if amfUe == nil {
		logger.NgapLog.Infof("Release UE Context : RanUe[AmfUeNgapId: %d]", ranUe.AmfUeNgapId)
		err := ranUe.Remove()
		if err != nil {
			logger.NgapLog.Errorln(err.Error())
		}
		return
	}
	// TODO: AMF shall, if supported, store it and may use it for subsequent paging
	if infoOnRecommendedCellsAndRANNodesForPaging != nil {
		amfUe.InfoOnRecommendedCellsAndRanNodesForPaging = new(amf_context.InfoOnRecommendedCellsAndRanNodesForPaging)

		recommendedCells := amfUe.InfoOnRecommendedCellsAndRanNodesForPaging.RecommendedCells
		for _, item := range infoOnRecommendedCellsAndRANNodesForPaging.RecommendedCellsForPaging.RecommendedCellList.List {
			recommendedCell := amf_context.RecommendedCell{}

			switch item.NGRANCGI.Present {
			case ngapType.NGRANCGIPresentNRCGI:
				recommendedCell.NgRanCGI.Present = amf_context.NgRanCgiPresentNRCGI
				recommendedCell.NgRanCGI.NRCGI = new(models.Ncgi)
				plmnID := ngapConvert.PlmnIdToModels(item.NGRANCGI.NRCGI.PLMNIdentity)
				recommendedCell.NgRanCGI.NRCGI.PlmnId = &plmnID
				recommendedCell.NgRanCGI.NRCGI.NrCellId = ngapConvert.BitStringToHex(&item.NGRANCGI.NRCGI.NRCellIdentity.Value)
			case ngapType.NGRANCGIPresentEUTRACGI:
				recommendedCell.NgRanCGI.Present = amf_context.NgRanCgiPresentEUTRACGI
				recommendedCell.NgRanCGI.EUTRACGI = new(models.Ecgi)
				plmnID := ngapConvert.PlmnIdToModels(item.NGRANCGI.EUTRACGI.PLMNIdentity)
				recommendedCell.NgRanCGI.EUTRACGI.PlmnId = &plmnID
				recommendedCell.NgRanCGI.EUTRACGI.EutraCellId = ngapConvert.BitStringToHex(&item.NGRANCGI.EUTRACGI.EUTRACellIdentity.Value)
			}

			if item.TimeStayedInCell != nil {
				recommendedCell.TimeStayedInCell = new(int64)
				*recommendedCell.TimeStayedInCell = *item.TimeStayedInCell
			}

			recommendedCells = append(recommendedCells, recommendedCell)
		}

		recommendedRanNodes := amfUe.InfoOnRecommendedCellsAndRanNodesForPaging.RecommendedRanNodes
		for _, item := range infoOnRecommendedCellsAndRANNodesForPaging.RecommendRANNodesForPaging.RecommendedRANNodeList.List {
			recommendedRanNode := amf_context.RecommendRanNode{}

			switch item.AMFPagingTarget.Present {
			case ngapType.AMFPagingTargetPresentGlobalRANNodeID:
				recommendedRanNode.Present = amf_context.RecommendRanNodePresentRanNode
				recommendedRanNode.GlobalRanNodeId = new(models.GlobalRanNodeId)
				// TODO: recommendedRanNode.GlobalRanNodeId = ngapConvert.RanIdToModels(item.AMFPagingTarget.GlobalRANNodeID)
			case ngapType.AMFPagingTargetPresentTAI:
				recommendedRanNode.Present = amf_context.RecommendRanNodePresentTAI
				tai := ngapConvert.TaiToModels(*item.AMFPagingTarget.TAI)
				recommendedRanNode.Tai = &tai
			}
			recommendedRanNodes = append(recommendedRanNodes, recommendedRanNode)
		}
	}

	// for each pduSessionID invoke Nsmf_PDUSession_UpdateSMContext Request
	var cause amf_context.CauseAll
	if tmp, exist := amfUe.ReleaseCause[ran.AnType]; exist {
		cause = *tmp
	}
	if amfUe.Sm[ran.AnType].Check(gmm_state.REGISTERED) {
		Ngaplog.Info("[NGAP] Rel Ue Context in GMM-Registered")
		if pDUSessionResourceList != nil {
			for _, pduSessionReourceItem := range pDUSessionResourceList.List {
				pduSessionID := int32(pduSessionReourceItem.PDUSessionID.Value)
				response, _, _, err := amf_consumer.SendUpdateSmContextDeactivateUpCnxState(amfUe, pduSessionID, cause)
				if err != nil {
					logger.NgapLog.Errorf("Send Update SmContextDeactivate UpCnxState Error[%s]", err.Error())
				} else if response == nil {
					logger.NgapLog.Errorln("Send Update SmContextDeactivate UpCnxState Error")
				}
			}
		}
	}

	// Remove UE N2 Connection
	amfUe.ReleaseCause[ran.AnType] = nil
	switch ranUe.ReleaseAction {
	case amf_context.UeContextN2NormalRelease:
		logger.NgapLog.Infof("Release UE[%s] Context : N2 Connection Release", amfUe.Supi)
		// amfUe.DetachRanUe(ran.AnType)
		err := ranUe.Remove()
		if err != nil {
			logger.NgapLog.Errorln(err.Error())
		}
	case amf_context.UeContextReleaseUeContext:
		logger.NgapLog.Infof("Release UE[%s] Context : Release Ue Context", amfUe.Supi)
		err := ranUe.Remove()
		if err != nil {
			logger.NgapLog.Errorln(err.Error())
		}
		amfUe.Remove()
	case amf_context.UeContextReleaseHandover:
		logger.NgapLog.Infof("Release UE[%s] Context : Release for Handover", amfUe.Supi)
		amf_context.DetachSourceUeTargetUe(ranUe)
		err := ranUe.Remove()
		if err != nil {
			logger.NgapLog.Errorln(err.Error())
		}
		// Todo: remove indirect tunnel
	default:
		logger.NgapLog.Errorf("Invalid Release Action[%d]", ranUe.ReleaseAction)

	}

}

func HandlePDUSessionResourceReleaseResponse(ran *amf_context.AmfRan, message *ngapType.NGAPPDU) {

	var aMFUENGAPID *ngapType.AMFUENGAPID
	var rANUENGAPID *ngapType.RANUENGAPID
	var pDUSessionResourceReleasedList *ngapType.PDUSessionResourceReleasedListRelRes
	var userLocationInformation *ngapType.UserLocationInformation
	var criticalityDiagnostics *ngapType.CriticalityDiagnostics

	if ran == nil {
		logger.NgapLog.Error("ran is nil")
		return
	}
	if message == nil {
		logger.NgapLog.Error("NGAP Message is nil")
		return
	}
	successfulOutcome := message.SuccessfulOutcome
	if successfulOutcome == nil {
		logger.NgapLog.Error("SuccessfulOutcome is nil")
		return
	}
	pDUSessionResourceReleaseResponse := successfulOutcome.Value.PDUSessionResourceReleaseResponse
	if pDUSessionResourceReleaseResponse == nil {
		logger.NgapLog.Error("PDUSessionResourceReleaseResponse is nil")
		return
	}

	logger.NgapLog.Info("[AMF] PDU Session Resource Release Response")

	for _, ie := range pDUSessionResourceReleaseResponse.ProtocolIEs.List {
		switch ie.Id.Value {
		case ngapType.ProtocolIEIDAMFUENGAPID:
			aMFUENGAPID = ie.Value.AMFUENGAPID
			logger.NgapLog.Trace("[NGAP] Decode IE AmfUeNgapID")
			if aMFUENGAPID == nil {
				logger.NgapLog.Error("AmfUeNgapID is nil")
				return
			}
		case ngapType.ProtocolIEIDRANUENGAPID:
			rANUENGAPID = ie.Value.RANUENGAPID
			logger.NgapLog.Trace("[NGAP] Decode IE RanUENgapID")
			if rANUENGAPID == nil {
				logger.NgapLog.Error("RanUeNgapID is nil")
				return
			}
		case ngapType.ProtocolIEIDPDUSessionResourceReleasedListRelRes:
			pDUSessionResourceReleasedList = ie.Value.PDUSessionResourceReleasedListRelRes
			logger.NgapLog.Trace("[NGAP] Decode IE PDUSessionResourceReleasedList")
			if pDUSessionResourceReleasedList == nil {
				logger.NgapLog.Error("PDUSessionResourceReleasedList is nil")
				return
			}
		case ngapType.ProtocolIEIDUserLocationInformation:
			userLocationInformation = ie.Value.UserLocationInformation
			logger.NgapLog.Trace("[NGAP] Decode IE UserLocationInformation")
		case ngapType.ProtocolIEIDCriticalityDiagnostics:
			criticalityDiagnostics = ie.Value.CriticalityDiagnostics
			logger.NgapLog.Trace("[NGAP] Decode IE CriticalityDiagnostics")
		}
	}

	printRanInfo(ran)

	ranUe := ran.RanUeFindByRanUeNgapID(rANUENGAPID.Value)
	if ranUe == nil {
		logger.NgapLog.Errorf("No UE Context[RanUeNgapID: %d]", rANUENGAPID.Value)
		return
	}

	if userLocationInformation != nil {
		ranUe.UpdateLocation(userLocationInformation)
	}

	if criticalityDiagnostics != nil {
		printCriticalityDiagnostics(criticalityDiagnostics)
	}

	amfUe := ranUe.AmfUe
	if amfUe == nil {
		Ngaplog.Error("amfUe is nil")
		return
	}
	if pDUSessionResourceReleasedList != nil {
		Ngaplog.Trace("[NGAP] Send PDUSessionResourceReleaseResponseTransfer to SMF")

		for _, item := range pDUSessionResourceReleasedList.List {
			pduSessionID := int32(item.PDUSessionID.Value)
			transfer := item.PDUSessionResourceReleaseResponseTransfer

			_, responseErr, problemDetail, err := amf_consumer.SendUpdateSmContextN2Info(amfUe, pduSessionID, models.N2SmInfoType_PDU_RES_REL_RSP, transfer)
			// TODO: error handling
			if err != nil {
				Ngaplog.Errorf("SendUpdateSmContextN2Info[PDUSessionResourceReleaseResponse] Error:\n%s", err.Error())
			} else if responseErr != nil && responseErr.JsonData.Error != nil {
				Ngaplog.Errorf("SendUpdateSmContextN2Info[PDUSessionResourceReleaseResponse] Error:\n%s", responseErr.JsonData.Error.Cause)
			} else if problemDetail != nil {
				Ngaplog.Errorf("SendUpdateSmContextN2Info[PDUSessionResourceReleaseResponse] Error:\n%s", problemDetail.Cause)
			}
		}
	}
}

func HandleUERadioCapabilityCheckResponse(ran *amf_context.AmfRan, message *ngapType.NGAPPDU) {

	var aMFUENGAPID *ngapType.AMFUENGAPID
	var rANUENGAPID *ngapType.RANUENGAPID
	var iMSVoiceSupportIndicator *ngapType.IMSVoiceSupportIndicator
	var criticalityDiagnostics *ngapType.CriticalityDiagnostics
	var ranUe *amf_context.RanUe

	logger.SetLogLevel(logrus.TraceLevel)
	logger.SetReportCaller(false)

	if ran == nil {
		logger.NgapLog.Error("ran is nil")
		return
	}
	if message == nil {
		logger.NgapLog.Error("NGAP Message is nil")
		return
	}
	successfulOutcome := message.SuccessfulOutcome
	if successfulOutcome == nil {
		logger.NgapLog.Error("SuccessfulOutcome is nil")
		return
	}

	uERadioCapabilityCheckResponse := successfulOutcome.Value.UERadioCapabilityCheckResponse
	if uERadioCapabilityCheckResponse == nil {
		logger.NgapLog.Error("UERadioCapabilityCheckResponse is nil")
		return
	}
	logger.NgapLog.Info("[AMF] UE Radio Capability Check Response")

	for i := 0; i < len(uERadioCapabilityCheckResponse.ProtocolIEs.List); i++ {
		ie := uERadioCapabilityCheckResponse.ProtocolIEs.List[i]
		switch ie.Id.Value {
		case ngapType.ProtocolIEIDAMFUENGAPID:
			aMFUENGAPID = ie.Value.AMFUENGAPID
			logger.NgapLog.Trace("[NGAP] Decode IE AmfUeNgapID")
			if aMFUENGAPID == nil {
				logger.NgapLog.Error("AmfUeNgapID is nil")
				return
			}
		case ngapType.ProtocolIEIDRANUENGAPID:
			rANUENGAPID = ie.Value.RANUENGAPID
			logger.NgapLog.Trace("[NGAP] Decode IE RanUeNgapID")
			if rANUENGAPID == nil {
				logger.NgapLog.Error("RanUeNgapID is nil")
				return
			}
		case ngapType.ProtocolIEIDIMSVoiceSupportIndicator:
			iMSVoiceSupportIndicator = ie.Value.IMSVoiceSupportIndicator
			logger.NgapLog.Trace("[NGAP] Decode IE IMSVoiceSupportIndicator")
			if iMSVoiceSupportIndicator == nil {
				logger.NgapLog.Error("iMSVoiceSupportIndicator is nil")
				return
			}
		case ngapType.ProtocolIEIDCriticalityDiagnostics:
			criticalityDiagnostics = ie.Value.CriticalityDiagnostics
			logger.NgapLog.Trace("[NGAP] Decode IE CriticalityDiagnostics")
		}
	}

	printRanInfo(ran)

	for i := range ran.RanUeList {
		if ran.RanUeList[i].RanUeNgapId == rANUENGAPID.Value {
			ranUe = ran.RanUeList[i]
		}
	}
	if ranUe == nil {
		logger.NgapLog.Errorf("No UE Context[RanUeNgapID: %d]", rANUENGAPID.Value)
		return
	}

	// TODO: handle iMSVoiceSupportIndicator

	if criticalityDiagnostics != nil {
		printCriticalityDiagnostics(criticalityDiagnostics)
	}
}

func HandleLocationReportingFailureIndication(ran *amf_context.AmfRan, message *ngapType.NGAPPDU) {

	var aMFUENGAPID *ngapType.AMFUENGAPID
	var rANUENGAPID *ngapType.RANUENGAPID
	var ranUe *amf_context.RanUe

	var cause *ngapType.Cause

	if ran == nil {
		logger.NgapLog.Error("ran is nil")
		return
	}
	if message == nil {
		logger.NgapLog.Error("NGAP Message is nil")
		return
	}
	initiatingMessage := message.InitiatingMessage
	if initiatingMessage == nil {
		logger.NgapLog.Error("Initiating Message is nil")
		return
	}
	locationReportingFailureIndication := initiatingMessage.Value.LocationReportingFailureIndication
	if locationReportingFailureIndication == nil {
		logger.NgapLog.Error("LocationReportingFailureIndication is nil")
		return
	}

	logger.NgapLog.Info("[AMF] Location Reporting Failure Indication")

	for i := 0; i < len(locationReportingFailureIndication.ProtocolIEs.List); i++ {
		ie := locationReportingFailureIndication.ProtocolIEs.List[i]
		switch ie.Id.Value {
		case ngapType.ProtocolIEIDAMFUENGAPID:
			aMFUENGAPID = ie.Value.AMFUENGAPID
			logger.NgapLog.Trace("[NGAP] Decode IE AmfUeNgapID")
			if aMFUENGAPID == nil {
				logger.NgapLog.Error("AmfUeNgapID is nil")
				return
			}
		case ngapType.ProtocolIEIDRANUENGAPID:
			rANUENGAPID = ie.Value.RANUENGAPID
			logger.NgapLog.Trace("[NGAP] Decode IE RanUeNgapID")
			if rANUENGAPID == nil {
				logger.NgapLog.Error("RanUeNgapID is nil")
				return
			}
		case ngapType.ProtocolIEIDCause:
			cause = ie.Value.Cause
			logger.NgapLog.Trace("[NGAP] Decode IE Cause")
			if cause == nil {
				logger.NgapLog.Error("Cause is nil")
				return
			}
		}
	}

	printRanInfo(ran)

	printAndGetCause(cause)

	for i := range ran.RanUeList {
		if ran.RanUeList[i].RanUeNgapId == rANUENGAPID.Value {
			ranUe = ran.RanUeList[i]
		}
	}
	if ranUe == nil {
		logger.NgapLog.Errorf("No UE Context[RanUeNgapID: %d]", rANUENGAPID.Value)
		return
	}
}

func HandleInitialUEMessage(ran *amf_context.AmfRan, message *ngapType.NGAPPDU) {

	amfSelf := amf_context.AMF_Self()

	var rANUENGAPID *ngapType.RANUENGAPID
	var nASPDU *ngapType.NASPDU
	var userLocationInformation *ngapType.UserLocationInformation
	var rRCEstablishmentCause *ngapType.RRCEstablishmentCause
	var fiveGSTMSI *ngapType.FiveGSTMSI
	var aMFSetID *ngapType.AMFSetID
	var uEContextRequest *ngapType.UEContextRequest
	var allowedNSSAI *ngapType.AllowedNSSAI

	var iesCriticalityDiagnostics ngapType.CriticalityDiagnosticsIEList

	if message == nil {
		logger.NgapLog.Error("NGAP Message is nil")
		return
	}

	initiatingMessage := message.InitiatingMessage
	if initiatingMessage == nil {
		logger.NgapLog.Error("Initiating Message is nil")
		return
	}
	initialUEMessage := initiatingMessage.Value.InitialUEMessage
	if initialUEMessage == nil {
		logger.NgapLog.Error("InitialUEMessage is nil")
		return
	}
	logger.NgapLog.Info("[AMF] Initial UE Message")

	for _, ie := range initialUEMessage.ProtocolIEs.List {
		switch ie.Id.Value {
		case ngapType.ProtocolIEIDRANUENGAPID: // reject
			rANUENGAPID = ie.Value.RANUENGAPID
			Ngaplog.Trace("[NGAP] Decode IE RanUeNgapID")
			if rANUENGAPID == nil {
				Ngaplog.Error("RanUeNgapID is nil")
				item := buildCriticalityDiagnosticsIEItem(ngapType.CriticalityPresentReject, ngapType.ProtocolIEIDRANUENGAPID, ngapType.TypeOfErrorPresentMissing)
				iesCriticalityDiagnostics.List = append(iesCriticalityDiagnostics.List, item)
			}
		case ngapType.ProtocolIEIDNASPDU: // reject
			nASPDU = ie.Value.NASPDU
			logger.NgapLog.Trace("[NGAP] Decode IE NasPdu")
			if nASPDU == nil {
				Ngaplog.Error("NasPdu is nil")
				item := buildCriticalityDiagnosticsIEItem(ngapType.CriticalityPresentReject, ngapType.ProtocolIEIDNASPDU, ngapType.TypeOfErrorPresentMissing)
				iesCriticalityDiagnostics.List = append(iesCriticalityDiagnostics.List, item)
			}
		case ngapType.ProtocolIEIDUserLocationInformation: // reject
			userLocationInformation = ie.Value.UserLocationInformation
			Ngaplog.Trace("[NGAP] Decode IE UserLocationInformation")
			if userLocationInformation == nil {
				Ngaplog.Error("UserLocationInformation is nil")
				item := buildCriticalityDiagnosticsIEItem(ngapType.CriticalityPresentReject, ngapType.ProtocolIEIDUserLocationInformation, ngapType.TypeOfErrorPresentMissing)
				iesCriticalityDiagnostics.List = append(iesCriticalityDiagnostics.List, item)
			}
		case ngapType.ProtocolIEIDRRCEstablishmentCause: // ignore
			rRCEstablishmentCause = ie.Value.RRCEstablishmentCause
			Ngaplog.Trace("[NGAP] Decode IE RRCEstablishmentCause")
		case ngapType.ProtocolIEIDFiveGSTMSI: // optional, reject
			fiveGSTMSI = ie.Value.FiveGSTMSI
			Ngaplog.Trace("[NGAP] Decode IE 5G-S-TMSI")
		case ngapType.ProtocolIEIDAMFSetID: // optional, ignore
			aMFSetID = ie.Value.AMFSetID
			Ngaplog.Trace("[NGAP] Decode IE AmfSetID")
		case ngapType.ProtocolIEIDUEContextRequest: // optional, ignore
			uEContextRequest = ie.Value.UEContextRequest
			Ngaplog.Trace("[NGAP] Decode IE UEContextRequest")
		case ngapType.ProtocolIEIDAllowedNSSAI: // optional, reject
			allowedNSSAI = ie.Value.AllowedNSSAI
			Ngaplog.Trace("[NGAP] Decode IE Allowed NSSAI")
		}
	}

	if len(iesCriticalityDiagnostics.List) > 0 {
		Ngaplog.Trace("Has missing reject IE(s)")

		procedureCode := ngapType.ProcedureCodeInitialUEMessage
		triggeringMessage := ngapType.TriggeringMessagePresentInitiatingMessage
		procedureCriticality := ngapType.CriticalityPresentIgnore
		criticalityDiagnostics := buildCriticalityDiagnostics(&procedureCode, &triggeringMessage, &procedureCriticality, &iesCriticalityDiagnostics)
		ngap_message.SendErrorIndication(ran, nil, nil, nil, &criticalityDiagnostics)
	}

	printRanInfo(ran)

	ranUe := ran.RanUeFindByRanUeNgapID(rANUENGAPID.Value)
	if ranUe != nil && ranUe.AmfUe == nil {
		err := ranUe.Remove()
		if err != nil {
			Ngaplog.Errorln(err.Error())
		}
		ranUe = nil
	}
	if ranUe == nil {
		ranUe = ran.NewRanUe()
		ranUe.RanUeNgapId = rANUENGAPID.Value
		Ngaplog.Debugf("New RanUe [RanUeNgapID: %d]", ranUe.RanUeNgapId)

		if fiveGSTMSI != nil {
			Ngaplog.Debug("Receive 5G-S-TMSI")

			servedGuami := amfSelf.ServedGuamiList[0]

			// <5G-S-TMSI> := <AMF Set ID><AMF Pointer><5G-TMSI>
			// GUAMI := <MCC><MNC><AMF Region ID><AMF Set ID><AMF Pointer>
			// 5G-GUTI := <GUAMI><5G-TMSI>
			tmpReginID, _, _ := ngapConvert.AmfIdToNgap(servedGuami.AmfId)
			amfID := ngapConvert.AmfIdToModels(tmpReginID, fiveGSTMSI.AMFSetID.Value, fiveGSTMSI.AMFPointer.Value)

			tmsi := hex.EncodeToString(fiveGSTMSI.FiveGTMSI.Value)

			guti := servedGuami.PlmnId.Mcc + servedGuami.PlmnId.Mnc + amfID + tmsi

			// TODO: invoke Namf_Communication_UEContextTransfer if serving AMF has changed since last Registration Request procedure
			// Described in TS 23.502 4.2.2.2.2 step 4 (without UDSF deployment)

			amfUe := amfSelf.AmfUeFindByGuti(guti)
			if amfUe == nil {
				Ngaplog.Warnf("Unknown UE [GUTI: %s]", guti)
			} else {
				Ngaplog.Tracef("find AmfUe [GUTI: %s]", guti)

				if amfUe.CmConnect(ran.AnType) {
					Ngaplog.Debug("Implicit Deregistration")
					Ngaplog.Tracef("AmfUeNgapID[%d] RanUeNgapID[%d]", amfUe.RanUe[ran.AnType].AmfUeNgapId, amfUe.RanUe[ran.AnType].RanUeNgapId)
					amfUe.DetachRanUe(ran.AnType)
				}
				// TODO: stop Implicit Deregistration timer
				Ngaplog.Debugf("AmfUe Attach RanUe [RanUeNgapID: %d]", ranUe.RanUeNgapId)
				amfUe.AttachRanUe(ranUe)
			}
		} else {
			ranUe.AmfUe = amfSelf.NewAmfUe("")
			if err := gmm.InitAmfUeSm(ranUe.AmfUe); err != nil {
				logger.NgapLog.Errorf("InitAmfUeSm error: %v", err.Error())
			}
			ranUe.AmfUe.AttachRanUe(ranUe)
		}
	} else {
		ranUe.AmfUe.AttachRanUe(ranUe)
	}

	if userLocationInformation != nil {
		ranUe.UpdateLocation(userLocationInformation)
	}

	if rRCEstablishmentCause != nil {
		Ngaplog.Tracef("[Initial UE Message] RRC Establishment Cause[%d]", rRCEstablishmentCause.Value)
	}

	if uEContextRequest != nil {
		Ngaplog.Debug("Trigger initial Context Setup procedure")
		// TODO: Trigger Initial Context Setup procedure
	}

	// TS 23.502 4.2.2.2.3 step 6a Nnrf_NFDiscovery_Request (NF type, AMF Set)
	if aMFSetID != nil {
		// TODO: This is a rerouted message
		// TS 38.413: AMF shall, if supported, use the IE as described in TS 23.502
	}

	// ng-ran propagate allowedNssai in the rerouted initial ue message (TS 38.413 8.6.5)
	// TS 23.502 4.2.2.2.3 step 4a Nnssf_NSSelection_Get
	if allowedNSSAI != nil {
		// TODO: AMF should use it as defined in TS 23.502
	}

	amf_nas.HandleNAS(ranUe, ngapType.ProcedureCodeInitialUEMessage, nASPDU.Value)
}

func HandlePDUSessionResourceSetupResponse(ran *amf_context.AmfRan, message *ngapType.NGAPPDU) {

	var aMFUENGAPID *ngapType.AMFUENGAPID
	var rANUENGAPID *ngapType.RANUENGAPID
	var pDUSessionResourceSetupResponseList *ngapType.PDUSessionResourceSetupListSURes
	var pDUSessionResourceFailedToSetupList *ngapType.PDUSessionResourceFailedToSetupListSURes
	var criticalityDiagnostics *ngapType.CriticalityDiagnostics

	var ranUe *amf_context.RanUe

	if ran == nil {
		Ngaplog.Error("ran is nil")
		return
	}
	if message == nil {
		Ngaplog.Error("NGAP Message is nil")
		return
	}
	successfulOutcome := message.SuccessfulOutcome
	if successfulOutcome == nil {
		Ngaplog.Error("SuccessfulOutcome is nil")
		return
	}
	pDUSessionResourceSetupResponse := successfulOutcome.Value.PDUSessionResourceSetupResponse
	if pDUSessionResourceSetupResponse == nil {
		Ngaplog.Error("PDUSessionResourceSetupResponse is nil")
		return
	}

	Ngaplog.Info("[AMF] PDU Session Resource Setup Response")

	for _, ie := range pDUSessionResourceSetupResponse.ProtocolIEs.List {
		switch ie.Id.Value {
		case ngapType.ProtocolIEIDAMFUENGAPID: // ignore
			aMFUENGAPID = ie.Value.AMFUENGAPID
			Ngaplog.Trace("[NGAP] Decode IE AmfUeNgapID")
		case ngapType.ProtocolIEIDRANUENGAPID: // ignore
			rANUENGAPID = ie.Value.RANUENGAPID
			Ngaplog.Trace("[NGAP] Decode IE RanUeNgapID")
		case ngapType.ProtocolIEIDPDUSessionResourceSetupListSURes: // ignore
			pDUSessionResourceSetupResponseList = ie.Value.PDUSessionResourceSetupListSURes
			Ngaplog.Trace("[NGAP] Decode IE PDUSessionResourceSetupListSURes")
		case ngapType.ProtocolIEIDPDUSessionResourceFailedToSetupListSURes: // ignore
			pDUSessionResourceFailedToSetupList = ie.Value.PDUSessionResourceFailedToSetupListSURes
			Ngaplog.Trace("[NGAP] Decode IE PDUSessionResourceFailedToSetupListSURes")
		case ngapType.ProtocolIEIDCriticalityDiagnostics: // optional, ignore
			criticalityDiagnostics = ie.Value.CriticalityDiagnostics
			Ngaplog.Trace("[NGAP] Decode IE CriticalityDiagnostics")
		}
	}

	printRanInfo(ran)

	if rANUENGAPID != nil {
		ranUe = ran.RanUeFindByRanUeNgapID(rANUENGAPID.Value)
		if ranUe == nil {
			Ngaplog.Warnf("No UE Context[RanUeNgapID: %d]", rANUENGAPID.Value)
		}
	}

	if aMFUENGAPID != nil {
		ranUe = amf_context.AMF_Self().RanUeFindByAmfUeNgapID(aMFUENGAPID.Value)
		if ranUe == nil {
			Ngaplog.Warnf("No UE Context[AmfUeNgapID: %d]", aMFUENGAPID.Value)
			return
		}
	}

	if ranUe != nil {
		Ngaplog.Tracef("AmfUeNgapID[%d] RanUeNgapID[%d]", ranUe.AmfUeNgapId, ranUe.RanUeNgapId)
		amfUe := ranUe.AmfUe
		if amfUe == nil {
			Ngaplog.Error("amfUe is nil")
			return
		}

		if pDUSessionResourceSetupResponseList != nil {
			Ngaplog.Trace("[NGAP] Send PDUSessionResourceSetupResponseTransfer to SMF")

			for _, item := range pDUSessionResourceSetupResponseList.List {
				pduSessionID := int32(item.PDUSessionID.Value)
				transfer := item.PDUSessionResourceSetupResponseTransfer

				response, _, _, err := amf_consumer.SendUpdateSmContextN2Info(amfUe, pduSessionID, models.N2SmInfoType_PDU_RES_SETUP_RSP, transfer)
				if err != nil {
					Ngaplog.Errorf("SendUpdateSmContextN2Info[PDUSessionResourceSetupResponseTransfer] Error:\n%s", err.Error())
				}
				// RAN initiated QoS Flow Mobility in subclause 5.2.2.3.7
				if response != nil && response.BinaryDataN2SmInformation != nil {
					// TODO: n2SmInfo send to RAN
				} else if response == nil {
					// TODO: error handling
				}
			}
		}

		if pDUSessionResourceFailedToSetupList != nil {
			Ngaplog.Trace("[NGAP] Send PDUSessionResourceSetupUnsuccessfulTransfer to SMF")

			for _, item := range pDUSessionResourceFailedToSetupList.List {
				pduSessionID := int32(item.PDUSessionID.Value)
				transfer := item.PDUSessionResourceSetupUnsuccessfulTransfer

				response, _, _, err := amf_consumer.SendUpdateSmContextN2Info(amfUe, pduSessionID, models.N2SmInfoType_PDU_RES_SETUP_FAIL, transfer)
				if err != nil {
					Ngaplog.Errorf("SendUpdateSmContextN2Info[PDUSessionResourceSetupUnsuccessfulTransfer] Error:\n%s", err.Error())
				}

				if response != nil && response.BinaryDataN2SmInformation != nil {
					// TODO: n2SmInfo send to RAN
				} else if response == nil {
					// TODO: error handling
				}
			}
		}
	}

	if criticalityDiagnostics != nil {
		printCriticalityDiagnostics(criticalityDiagnostics)
	}
}

func HandlePDUSessionResourceModifyResponse(ran *amf_context.AmfRan, message *ngapType.NGAPPDU) {

	var aMFUENGAPID *ngapType.AMFUENGAPID
	var rANUENGAPID *ngapType.RANUENGAPID
	var pduSessionResourceModifyResponseList *ngapType.PDUSessionResourceModifyListModRes
	var pduSessionResourceFailedToModifyList *ngapType.PDUSessionResourceFailedToModifyListModRes
	var userLocationInformation *ngapType.UserLocationInformation
	var criticalityDiagnostics *ngapType.CriticalityDiagnostics

	var ranUe *amf_context.RanUe

	if ran == nil {
		Ngaplog.Error("ran is nil")
		return
	}
	if message == nil {
		Ngaplog.Error("NGAP Message is nil")
		return
	}
	successfulOutcome := message.SuccessfulOutcome
	if successfulOutcome == nil {
		Ngaplog.Error("SuccessfulOutcome is nil")
		return
	}
	pDUSessionResourceModifyResponse := successfulOutcome.Value.PDUSessionResourceModifyResponse
	if pDUSessionResourceModifyResponse == nil {
		Ngaplog.Error("PDUSessionResourceModifyResponse is nil")
		return
	}

	Ngaplog.Info("[AMF] PDU Session Resource Modify Response")

	for _, ie := range pDUSessionResourceModifyResponse.ProtocolIEs.List {
		switch ie.Id.Value {
		case ngapType.ProtocolIEIDAMFUENGAPID: // ignore
			aMFUENGAPID = ie.Value.AMFUENGAPID
			Ngaplog.Trace("[NGAP] Decode IE AmfUeNgapID")
		case ngapType.ProtocolIEIDRANUENGAPID: //ignore
			rANUENGAPID = ie.Value.RANUENGAPID
			Ngaplog.Trace("[NGAP] Decode IE RanUeNgapID")
		case ngapType.ProtocolIEIDPDUSessionResourceModifyListModRes: // ignore
			pduSessionResourceModifyResponseList = ie.Value.PDUSessionResourceModifyListModRes
			Ngaplog.Trace("[NGAP] Decode IE PDUSessionResourceModifyListModRes")
		case ngapType.ProtocolIEIDPDUSessionResourceFailedToModifyListModRes: // ignore
			pduSessionResourceFailedToModifyList = ie.Value.PDUSessionResourceFailedToModifyListModRes
			Ngaplog.Trace("[NGAP] Decode IE PDUSessionResourceFailedToModifyListModRes")
		case ngapType.ProtocolIEIDUserLocationInformation: // optional, ignore
			userLocationInformation = ie.Value.UserLocationInformation
			Ngaplog.Trace("[NGAP] Decode IE UserLocationInformation")
		case ngapType.ProtocolIEIDCriticalityDiagnostics: // optional, ignore
			criticalityDiagnostics = ie.Value.CriticalityDiagnostics
			Ngaplog.Trace("[NGAP] Decode IE CriticalityDiagnostics")
		}
	}

	printRanInfo(ran)

	if rANUENGAPID != nil {
		ranUe = ran.RanUeFindByRanUeNgapID(rANUENGAPID.Value)
		if ranUe == nil {
			Ngaplog.Warnf("No UE Context[RanUeNgapID: %d]", rANUENGAPID.Value)
		}
	}

	if aMFUENGAPID != nil {
		ranUe = amf_context.AMF_Self().RanUeFindByAmfUeNgapID(aMFUENGAPID.Value)
		if ranUe == nil {
			Ngaplog.Warnf("No UE Context[AmfUeNgapID: %d]", aMFUENGAPID.Value)
			return
		}
	}

	if ranUe != nil {
		Ngaplog.Tracef("AmfUeNgapID[%d] RanUeNgapID[%d]", ranUe.AmfUeNgapId, ranUe.RanUeNgapId)
		amfUe := ranUe.AmfUe
		if amfUe == nil {
			Ngaplog.Error("amfUe is nil")
			return
		}

		if pduSessionResourceModifyResponseList != nil {
			Ngaplog.Trace("[NGAP] Send PDUSessionResourceModifyResponseTransfer to SMF")

			for _, item := range pduSessionResourceModifyResponseList.List {
				pduSessionID := int32(item.PDUSessionID.Value)
				transfer := item.PDUSessionResourceModifyResponseTransfer

				response, _, _, err := amf_consumer.SendUpdateSmContextN2Info(amfUe, pduSessionID, models.N2SmInfoType_PDU_RES_MOD_RSP, *transfer)
				if err != nil {
					Ngaplog.Errorf("SendUpdateSmContextN2Info[PDUSessionResourceModifyResponseTransfer] Error:\n%s", err.Error())
				}
				if response != nil && response.BinaryDataN2SmInformation != nil {
					// TODO: n2SmInfo send to RAN
				} else if response == nil {
					// TODO: error handling
				}
			}
		}

		if pduSessionResourceFailedToModifyList != nil {
			Ngaplog.Trace("[NGAP] Send PDUSessionResourceModifyUnsuccessfulTransfer to SMF")

			for _, item := range pduSessionResourceFailedToModifyList.List {
				pduSessionID := int32(item.PDUSessionID.Value)
				transfer := item.PDUSessionResourceModifyUnsuccessfulTransfer

				response, _, _, err := amf_consumer.SendUpdateSmContextN2Info(amfUe, pduSessionID, models.N2SmInfoType_PDU_RES_MOD_FAIL, transfer)
				if err != nil {
					Ngaplog.Errorf("SendUpdateSmContextN2Info[PDUSessionResourceModifyUnsuccessfulTransfer] Error:\n%s", err.Error())
				}
				if response != nil && response.BinaryDataN2SmInformation != nil {
					// TODO: n2SmInfo send to RAN
				} else if response == nil {
					// TODO: error handling
				}
			}
		}

		if userLocationInformation != nil {
			ranUe.UpdateLocation(userLocationInformation)
		}
	}

	if criticalityDiagnostics != nil {
		printCriticalityDiagnostics(criticalityDiagnostics)
	}
}

func HandlePDUSessionResourceNotify(ran *amf_context.AmfRan, message *ngapType.NGAPPDU) {

	var aMFUENGAPID *ngapType.AMFUENGAPID
	var rANUENGAPID *ngapType.RANUENGAPID
	var pDUSessionResourceNotifyList *ngapType.PDUSessionResourceNotifyList
	var pDUSessionResourceReleasedListNot *ngapType.PDUSessionResourceReleasedListNot
	var userLocationInformation *ngapType.UserLocationInformation

	var ranUe *amf_context.RanUe

	if ran == nil {
		Ngaplog.Error("ran is nil")
		return
	}
	if message == nil {
		Ngaplog.Error("NGAP Message is nil")
		return
	}
	initiatingMessage := message.InitiatingMessage
	if initiatingMessage == nil {
		Ngaplog.Error("InitiatingMessage is nil")
		return
	}
	PDUSessionResourceNotify := initiatingMessage.Value.PDUSessionResourceNotify
	if PDUSessionResourceNotify == nil {
		Ngaplog.Error("PDUSessionResourceNotify is nil")
		return
	}

	for _, ie := range PDUSessionResourceNotify.ProtocolIEs.List {
		switch ie.Id.Value {
		case ngapType.ProtocolIEIDAMFUENGAPID:
			aMFUENGAPID = ie.Value.AMFUENGAPID // reject
			Ngaplog.Trace("[NGAP] Decode IE AmfUeNgapID")
		case ngapType.ProtocolIEIDRANUENGAPID:
			rANUENGAPID = ie.Value.RANUENGAPID // reject
			Ngaplog.Trace("[NGAP] Decode IE RanUeNgapID")
		case ngapType.ProtocolIEIDPDUSessionResourceNotifyList: // reject
			pDUSessionResourceNotifyList = ie.Value.PDUSessionResourceNotifyList
			Ngaplog.Trace("[NGAP] Decode IE pDUSessionResourceNotifyList")
			if pDUSessionResourceNotifyList == nil {
				Ngaplog.Error("pDUSessionResourceNotifyList is nil")
			}
		case ngapType.ProtocolIEIDPDUSessionResourceReleasedListNot: // ignore
			pDUSessionResourceReleasedListNot = ie.Value.PDUSessionResourceReleasedListNot
			Ngaplog.Trace("[NGAP] Decode IE PDUSessionResourceReleasedListNot")
			if pDUSessionResourceReleasedListNot == nil {
				Ngaplog.Error("PDUSessionResourceReleasedListNot is nil")
			}
		case ngapType.ProtocolIEIDUserLocationInformation: // optional, ignore
			userLocationInformation = ie.Value.UserLocationInformation
			Ngaplog.Trace("[NGAP] Decode IE userLocationInformation")
			if userLocationInformation == nil {
				Ngaplog.Warn("userLocationInformation is nil [optional]")
			}
		}
	}

	printRanInfo(ran)

	ranUe = ran.RanUeFindByRanUeNgapID(rANUENGAPID.Value)
	if ranUe == nil {
		Ngaplog.Warnf("No UE Context[RanUeNgapID: %d]", rANUENGAPID.Value)
	}

	ranUe = amf_context.AMF_Self().RanUeFindByAmfUeNgapID(aMFUENGAPID.Value)
	if ranUe == nil {
		Ngaplog.Warnf("No UE Context[AmfUeNgapID: %d]", aMFUENGAPID.Value)
		return
	}

	Ngaplog.Tracef("AmfUeNgapID[%d] RanUeNgapID[%d]", ranUe.AmfUeNgapId, ranUe.RanUeNgapId)
	amfUe := ranUe.AmfUe
	if amfUe == nil {
		Ngaplog.Error("amfUe is nil")
		return
	}

	if userLocationInformation != nil {
		ranUe.UpdateLocation(userLocationInformation)
	}

	Ngaplog.Trace("[NGAP] Send PDUSessionResourceNotifyTransfer to SMF")

	for _, item := range pDUSessionResourceNotifyList.List {
		pduSessionID := int32(item.PDUSessionID.Value)
		transfer := item.PDUSessionResourceNotifyTransfer

		response, errResponse, problemDetail, err := amf_consumer.SendUpdateSmContextN2Info(amfUe, pduSessionID, models.N2SmInfoType_PDU_RES_NTY, transfer)
		if err != nil {
			Ngaplog.Errorf("SendUpdateSmContextN2Info[PDUSessionResourceNotifyTransfer] Error:\n%s", err.Error())
		}

		if response != nil {
			responseData := response.JsonData
			n2Info := response.BinaryDataN1SmMessage
			n1Msg := response.BinaryDataN2SmInformation
			if n2Info != nil {
				switch responseData.N2SmInfoType {
				case models.N2SmInfoType_PDU_RES_MOD_REQ:
					logger.HttpLog.Debugln("AMF Transfer NGAP PDU Resource Modify Req from SMF")
					var nasPdu []byte
					if n1Msg != nil {
						pduSessionId := uint8(pduSessionID)
						nasPdu, err = gmm_message.BuildDLNASTransport(amfUe, nasMessage.PayloadContainerTypeN1SMInfo, n1Msg, &pduSessionId, nil, nil, 0)
					}
					list := ngapType.PDUSessionResourceModifyListModReq{}
					ngap_message.AppendPDUSessionResourceModifyListModReq(&list, pduSessionID, nasPdu, n2Info)
					ngap_message.SendPDUSessionResourceModifyRequest(ranUe, list)
				}
			}
		} else if errResponse != nil {
			errJSON := errResponse.JsonData
			n1Msg := errResponse.BinaryDataN2SmInformation
			logger.HttpLog.Warnf("PDU Session Modification is rejected by SMF[pduSessionId:%d], Error[%s]\n", pduSessionID, errJSON.Error.Cause)
			if n1Msg != nil {
				gmm_message.SendDLNASTransport(ranUe, nasMessage.PayloadContainerTypeN1SMInfo, errResponse.BinaryDataN1SmMessage, &pduSessionID, 0, nil, 0)
			}
			// TODO: handle n2 info transfer
		} else if err != nil {
			return
		} else {
			// TODO: error handling
			logger.HttpLog.Errorf("Failed to Update smContext[pduSessionID: %d], Error[%v]", pduSessionID, problemDetail)
			return
		}
	}

	if pDUSessionResourceReleasedListNot != nil {
		Ngaplog.Trace("[NGAP] Send PDUSessionResourceNotifyReleasedTransfer to SMF")
		for _, item := range pDUSessionResourceReleasedListNot.List {
			pduSessionID := int32(item.PDUSessionID.Value)
			transfer := item.PDUSessionResourceNotifyReleasedTransfer

			response, errResponse, problemDetail, err := amf_consumer.SendUpdateSmContextN2Info(amfUe, pduSessionID, models.N2SmInfoType_PDU_RES_NTY_REL, transfer)
			if err != nil {
				Ngaplog.Errorf("SendUpdateSmContextN2Info[PDUSessionResourceNotifyReleasedTransfer] Error:\n%s", err.Error())
			}
			if response != nil {
				responseData := response.JsonData
				n2Info := response.BinaryDataN1SmMessage
				n1Msg := response.BinaryDataN2SmInformation
				if n2Info != nil {
					switch responseData.N2SmInfoType {
					case models.N2SmInfoType_PDU_RES_REL_CMD:
						logger.GmmLog.Debugln("AMF Transfer NGAP PDU Session Resource Rel Co from SMF")
						var nasPdu []byte
						if n1Msg != nil {
							pduSessionId := uint8(pduSessionID)
							nasPdu, err = gmm_message.BuildDLNASTransport(amfUe, nasMessage.PayloadContainerTypeN1SMInfo, n1Msg, &pduSessionId, nil, nil, 0)
						}
						list := ngapType.PDUSessionResourceToReleaseListRelCmd{}
						ngap_message.AppendPDUSessionResourceToReleaseListRelCmd(&list, pduSessionID, n2Info)
						ngap_message.SendPDUSessionResourceReleaseCommand(ranUe, nasPdu, list)
					}
				}
			} else if errResponse != nil {
				errJSON := errResponse.JsonData
				n1Msg := errResponse.BinaryDataN2SmInformation
				logger.HttpLog.Warnf("PDU Session Release is rejected by SMF[pduSessionId:%d], Error[%s]\n", pduSessionID, errJSON.Error.Cause)
				if n1Msg != nil {
					gmm_message.SendDLNASTransport(ranUe, nasMessage.PayloadContainerTypeN1SMInfo, errResponse.BinaryDataN1SmMessage, &pduSessionID, 0, nil, 0)
				}
			} else if err != nil {
				return
			} else {
				// TODO: error handling
				logger.HttpLog.Errorf("Failed to Update smContext[pduSessionID: %d], Error[%v]", pduSessionID, problemDetail)
				return
			}
		}
	}

}

func HandlePDUSessionResourceModifyIndication(ran *amf_context.AmfRan, message *ngapType.NGAPPDU) {

	var aMFUENGAPID *ngapType.AMFUENGAPID
	var rANUENGAPID *ngapType.RANUENGAPID
	var pduSessionResourceModifyIndicationList *ngapType.PDUSessionResourceModifyListModInd

	var iesCriticalityDiagnostics ngapType.CriticalityDiagnosticsIEList

	var ranUe *amf_context.RanUe

	if ran == nil {
		Ngaplog.Error("ran is nil")
		return
	}
	if message == nil {
		Ngaplog.Error("NGAP Message is nil")
		return
	}
	initiatingMessage := message.InitiatingMessage // reject
	if initiatingMessage == nil {
		Ngaplog.Error("InitiatingMessage is nil")
		cause := ngapType.Cause{
			Present: ngapType.CausePresentProtocol,
			Protocol: &ngapType.CauseProtocol{
				Value: ngapType.CauseProtocolPresentAbstractSyntaxErrorReject,
			},
		}
		ngap_message.SendErrorIndication(ran, nil, nil, &cause, nil)
		return
	}
	pDUSessionResourceModifyIndication := initiatingMessage.Value.PDUSessionResourceModifyIndication
	if pDUSessionResourceModifyIndication == nil {
		Ngaplog.Error("PDUSessionResourceModifyIndication is nil")
		cause := ngapType.Cause{
			Present: ngapType.CausePresentProtocol,
			Protocol: &ngapType.CauseProtocol{
				Value: ngapType.CauseProtocolPresentAbstractSyntaxErrorReject,
			},
		}
		ngap_message.SendErrorIndication(ran, nil, nil, &cause, nil)
		return
	}

	Ngaplog.Info("[AMF] PDU Session Resource Modify Indication")

	for _, ie := range pDUSessionResourceModifyIndication.ProtocolIEs.List {
		switch ie.Id.Value {
		case ngapType.ProtocolIEIDAMFUENGAPID: // reject
			aMFUENGAPID = ie.Value.AMFUENGAPID
			Ngaplog.Trace("[NGAP] Decode IE AmfUeNgapID")
			if aMFUENGAPID == nil {
				Ngaplog.Error("AmfUeNgapID is nil")
				item := buildCriticalityDiagnosticsIEItem(ngapType.CriticalityPresentReject, ngapType.ProtocolIEIDAMFUENGAPID, ngapType.TypeOfErrorPresentMissing)
				iesCriticalityDiagnostics.List = append(iesCriticalityDiagnostics.List, item)
			}
		case ngapType.ProtocolIEIDRANUENGAPID: // reject
			rANUENGAPID = ie.Value.RANUENGAPID
			Ngaplog.Trace("[NGAP] Decode IE RanUeNgapID")
			if rANUENGAPID == nil {
				Ngaplog.Error("RanUeNgapID is nil")
				item := buildCriticalityDiagnosticsIEItem(ngapType.CriticalityPresentReject, ngapType.ProtocolIEIDRANUENGAPID, ngapType.TypeOfErrorPresentMissing)
				iesCriticalityDiagnostics.List = append(iesCriticalityDiagnostics.List, item)
			}
		case ngapType.ProtocolIEIDPDUSessionResourceModifyListModInd: // reject
			pduSessionResourceModifyIndicationList = ie.Value.PDUSessionResourceModifyListModInd
			Ngaplog.Trace("[NGAP] Decode IE PDUSessionResourceModifyListModInd")
			if pduSessionResourceModifyIndicationList == nil {
				Ngaplog.Error("PDUSessionResourceModifyListModInd is nil")
				item := buildCriticalityDiagnosticsIEItem(ngapType.CriticalityPresentReject, ngapType.ProtocolIEIDPDUSessionResourceModifyListModInd, ngapType.TypeOfErrorPresentMissing)
				iesCriticalityDiagnostics.List = append(iesCriticalityDiagnostics.List, item)
			}
		}
	}

	if len(iesCriticalityDiagnostics.List) > 0 {
		Ngaplog.Error("Has missing reject IE(s)")

		procedureCode := ngapType.ProcedureCodePDUSessionResourceModifyIndication
		triggeringMessage := ngapType.TriggeringMessagePresentInitiatingMessage
		procedureCriticality := ngapType.CriticalityPresentReject
		criticalityDiagnostics := buildCriticalityDiagnostics(&procedureCode, &triggeringMessage, &procedureCriticality, &iesCriticalityDiagnostics)
		ngap_message.SendErrorIndication(ran, nil, nil, nil, &criticalityDiagnostics)
		return
	}

	printRanInfo(ran)

	ranUe = ran.RanUeFindByRanUeNgapID(rANUENGAPID.Value)
	if ranUe == nil {
		Ngaplog.Errorf("No UE Context[RanUeNgapID: %d]", rANUENGAPID.Value)
		cause := ngapType.Cause{
			Present: ngapType.CausePresentRadioNetwork,
			RadioNetwork: &ngapType.CauseRadioNetwork{
				Value: ngapType.CauseRadioNetworkPresentUnknownLocalUENGAPID,
			},
		}
		ngap_message.SendErrorIndication(ran, nil, nil, &cause, nil)
		return
	}

	Ngaplog.Tracef("UE Context AmfUeNgapID[%d] RanUeNgapID[%d]", ranUe.AmfUeNgapId, ranUe.RanUeNgapId)

	amfUe := ranUe.AmfUe
	if amfUe == nil {
		Ngaplog.Error("AmfUe is nil")
		return
	}

	pduSessionResourceModifyListModCfm := ngapType.PDUSessionResourceModifyListModCfm{}
	pduSessionResourceFailedToModifyListModCfm := ngapType.PDUSessionResourceFailedToModifyListModCfm{}

	Ngaplog.Trace("[NGAP] Send PDUSessionResourceModifyIndicationTransfer to SMF")
	for _, item := range pduSessionResourceModifyIndicationList.List {
		pduSessionID := item.PDUSessionID.Value
		transfer := item.PDUSessionResourceModifyIndicationTransfer

		response, errResponse, _, err := amf_consumer.SendUpdateSmContextN2Info(amfUe, int32(pduSessionID), models.N2SmInfoType_PDU_RES_MOD_IND, transfer)

		if err != nil {
			Ngaplog.Errorf("SendUpdateSmContextN2Info Error:\n%s", err.Error())
		}

		if response != nil && response.BinaryDataN2SmInformation != nil {
			ngap_message.AppendPDUSessionResourceModifyListModCfm(&pduSessionResourceModifyListModCfm, pduSessionID, response.BinaryDataN2SmInformation)
		}
		if errResponse != nil && errResponse.BinaryDataN2SmInformation != nil {
			ngap_message.AppendPDUSessionResourceFailedToModifyListModCfm(&pduSessionResourceFailedToModifyListModCfm, pduSessionID, errResponse.BinaryDataN2SmInformation)
		}
	}

	ngap_message.SendPDUSessionResourceModifyConfirm(ranUe, pduSessionResourceModifyListModCfm, pduSessionResourceFailedToModifyListModCfm, nil)
}

func HandleInitialContextSetupResponse(ran *amf_context.AmfRan, message *ngapType.NGAPPDU) {

	var aMFUENGAPID *ngapType.AMFUENGAPID
	var rANUENGAPID *ngapType.RANUENGAPID
	var pDUSessionResourceSetupResponseList *ngapType.PDUSessionResourceSetupListCxtRes
	var pDUSessionResourceFailedToSetupList *ngapType.PDUSessionResourceFailedToSetupListCxtRes
	var criticalityDiagnostics *ngapType.CriticalityDiagnostics

	if ran == nil {
		Ngaplog.Error("ran is nil")
		return
	}
	if message == nil {
		Ngaplog.Error("NGAP Message is nil")
		return
	}
	successfulOutcome := message.SuccessfulOutcome
	if successfulOutcome == nil {
		Ngaplog.Error("SuccessfulOutcome is nil")
		return
	}
	initialContextSetupResponse := successfulOutcome.Value.InitialContextSetupResponse
	if initialContextSetupResponse == nil {
		Ngaplog.Error("InitialContextSetupResponse is nil")
		return
	}

	Ngaplog.Info("[AMF] Initial Context Setup Response")

	for _, ie := range initialContextSetupResponse.ProtocolIEs.List {
		switch ie.Id.Value {
		case ngapType.ProtocolIEIDAMFUENGAPID:
			aMFUENGAPID = ie.Value.AMFUENGAPID
			Ngaplog.Trace("[NGAP] Decode IE AmfUeNgapID")
			if aMFUENGAPID == nil {
				Ngaplog.Warn("AmfUeNgapID is nil")
			}
		case ngapType.ProtocolIEIDRANUENGAPID:
			rANUENGAPID = ie.Value.RANUENGAPID
			Ngaplog.Trace("[NGAP] Decode IE RanUeNgapID")
			if rANUENGAPID == nil {
				Ngaplog.Warn("RanUeNgapID is nil")
			}
		case ngapType.ProtocolIEIDPDUSessionResourceSetupListCxtRes:
			pDUSessionResourceSetupResponseList = ie.Value.PDUSessionResourceSetupListCxtRes
			Ngaplog.Trace("[NGAP] Decode IE PDUSessionResourceSetupResponseList")
			if pDUSessionResourceSetupResponseList == nil {
				Ngaplog.Warn("PDUSessionResourceSetupResponseList is nil")
			}
		case ngapType.ProtocolIEIDPDUSessionResourceFailedToSetupListCxtRes:
			pDUSessionResourceFailedToSetupList = ie.Value.PDUSessionResourceFailedToSetupListCxtRes
			Ngaplog.Trace("[NGAP] Decode IE PDUSessionResourceFailedToSetupList")
			if pDUSessionResourceFailedToSetupList == nil {
				Ngaplog.Warn("PDUSessionResourceFailedToSetupList is nil")
			}
		case ngapType.ProtocolIEIDCriticalityDiagnostics:
			criticalityDiagnostics = ie.Value.CriticalityDiagnostics
			Ngaplog.Trace("[NGAP] Decode IE Criticality Diagnostics")
			if criticalityDiagnostics == nil {
				Ngaplog.Warn("Criticality Diagnostics is nil")
			}
		}
	}

	printRanInfo(ran)

	ranUe := ran.RanUeFindByRanUeNgapID(rANUENGAPID.Value)
	if ranUe == nil {
		Ngaplog.Errorf("No UE Context[RanUeNgapID: %d]", rANUENGAPID.Value)
		return
	}
	amfUe := ranUe.AmfUe
	if amfUe == nil {
		Ngaplog.Error("amfUe is nil")
		return
	}

	Ngaplog.Tracef("RanUeNgapID[%d] AmfUeNgapID[%d]", ranUe.RanUeNgapId, ranUe.AmfUeNgapId)

	if pDUSessionResourceSetupResponseList != nil {
		Ngaplog.Trace("[NGAP] Send PDUSessionResourceSetupResponseTransfer to SMF")

		for _, item := range pDUSessionResourceSetupResponseList.List {
			pduSessionID := int32(item.PDUSessionID.Value)
			transfer := item.PDUSessionResourceSetupResponseTransfer

			response, _, _, err := amf_consumer.SendUpdateSmContextN2Info(amfUe, pduSessionID, models.N2SmInfoType_PDU_RES_SETUP_RSP, transfer)
			if err != nil {
				Ngaplog.Errorf("SendUpdateSmContextN2Info[PDUSessionResourceSetupResponseTransfer] Error:\n%s", err.Error())
			}
			// RAN initiated QoS Flow Mobility in subclause 5.2.2.3.7
			if response != nil && response.BinaryDataN2SmInformation != nil {
				// TODO: n2SmInfo send to RAN
			} else if response == nil {
				// TODO: error handling
			}
		}
	}

	if pDUSessionResourceFailedToSetupList != nil {
		Ngaplog.Trace("[NGAP] Send PDUSessionResourceSetupUnsuccessfulTransfer to SMF")

		for _, item := range pDUSessionResourceFailedToSetupList.List {
			pduSessionID := int32(item.PDUSessionID.Value)
			transfer := item.PDUSessionResourceSetupUnsuccessfulTransfer

			response, _, _, err := amf_consumer.SendUpdateSmContextN2Info(amfUe, pduSessionID, models.N2SmInfoType_PDU_RES_SETUP_FAIL, transfer)
			if err != nil {
				Ngaplog.Errorf("SendUpdateSmContextN2Info[PDUSessionResourceSetupUnsuccessfulTransfer] Error:\n%s", err.Error())
			}

			if response != nil && response.BinaryDataN2SmInformation != nil {
				// TODO: n2SmInfo send to RAN
			} else if response == nil {
				// TODO: error handling
			}
		}
	}

	if criticalityDiagnostics != nil {
		Ngaplog.Trace("Criticality Diagnostics")
		if criticalityDiagnostics.ProcedureCriticality != nil {
			switch criticalityDiagnostics.ProcedureCriticality.Value {
			case ngapType.CriticalityPresentReject:
				Ngaplog.Trace("Procedure Criticality: Reject")
			case ngapType.CriticalityPresentIgnore:
				Ngaplog.Trace("Procedure Criticality: Ignore")
			case ngapType.CriticalityPresentNotify:
				Ngaplog.Trace("Procedure Criticality: Notify")
			}
		}
		if criticalityDiagnostics.IEsCriticalityDiagnostics != nil {
			for _, ieCriticalityDiagnostics := range criticalityDiagnostics.IEsCriticalityDiagnostics.List {
				Ngaplog.Tracef("IE ID: %d", ieCriticalityDiagnostics.IEID.Value)

				switch ieCriticalityDiagnostics.IECriticality.Value {
				case ngapType.CriticalityPresentReject:
					Ngaplog.Trace("Criticality Reject")
				case ngapType.CriticalityPresentNotify:
					Ngaplog.Trace("Criticality Notify")
				}

				switch ieCriticalityDiagnostics.TypeOfError.Value {
				case ngapType.TypeOfErrorPresentNotUnderstood:
					Ngaplog.Trace("Type of error: Not understood")
				case ngapType.TypeOfErrorPresentMissing:
					Ngaplog.Trace("Type of error: Missing")
				}
			}
		}
	}
}

func HandleInitialContextSetupFailure(ran *amf_context.AmfRan, message *ngapType.NGAPPDU) {

	var aMFUENGAPID *ngapType.AMFUENGAPID
	var rANUENGAPID *ngapType.RANUENGAPID
	var pDUSessionResourceFailedToSetupList *ngapType.PDUSessionResourceFailedToSetupListCxtFail
	var cause *ngapType.Cause
	var criticalityDiagnostics *ngapType.CriticalityDiagnostics

	if ran == nil {
		Ngaplog.Error("ran is nil")
		return
	}
	if message == nil {
		Ngaplog.Error("NGAP Message is nil")
		return
	}
	unsuccessfulOutcome := message.UnsuccessfulOutcome
	if unsuccessfulOutcome == nil {
		Ngaplog.Error("UnsuccessfulOutcome is nil")
		return
	}
	initialContextSetupFailure := unsuccessfulOutcome.Value.InitialContextSetupFailure
	if initialContextSetupFailure == nil {
		Ngaplog.Error("InitialContextSetupFailure is nil")
		return
	}

	Ngaplog.Info("[AMF] Initial Context Setup Failure")

	for _, ie := range initialContextSetupFailure.ProtocolIEs.List {
		switch ie.Id.Value {
		case ngapType.ProtocolIEIDAMFUENGAPID:
			aMFUENGAPID = ie.Value.AMFUENGAPID
			Ngaplog.Trace("[NGAP] Decode IE AmfUeNgapID")
			if aMFUENGAPID == nil {
				Ngaplog.Warn("AmfUeNgapID is nil")
			}
		case ngapType.ProtocolIEIDRANUENGAPID:
			rANUENGAPID = ie.Value.RANUENGAPID
			Ngaplog.Trace("[NGAP] Decode IE RanUeNgapID")
			if rANUENGAPID == nil {
				Ngaplog.Warn("RanUeNgapID is nil")
			}
		case ngapType.ProtocolIEIDPDUSessionResourceFailedToSetupListCxtFail:
			pDUSessionResourceFailedToSetupList = ie.Value.PDUSessionResourceFailedToSetupListCxtFail
			Ngaplog.Trace("[NGAP] Decode IE PDUSessionResourceFailedToSetupList")
			if pDUSessionResourceFailedToSetupList == nil {
				Ngaplog.Warn("PDUSessionResourceFailedToSetupList is nil")
			}
		case ngapType.ProtocolIEIDCause:
			cause = ie.Value.Cause
			Ngaplog.Trace("[NGAP] Decode IE Cause")
			if cause == nil {
				Ngaplog.Warn("Cause is nil")
			}
		case ngapType.ProtocolIEIDCriticalityDiagnostics:
			criticalityDiagnostics = ie.Value.CriticalityDiagnostics
			Ngaplog.Trace("[NGAP] Decode IE Criticality Diagnostics")
			if criticalityDiagnostics == nil {
				Ngaplog.Warn("CriticalityDiagnostics is nil")
			}
		}
	}

	printRanInfo(ran)

	printAndGetCause(cause)

	if criticalityDiagnostics != nil {
		Ngaplog.Trace("Criticality Diagnostics")
		if criticalityDiagnostics.ProcedureCriticality != nil {
			switch criticalityDiagnostics.ProcedureCriticality.Value {
			case ngapType.CriticalityPresentReject:
				Ngaplog.Trace("Procedure Criticality: Reject")
			case ngapType.CriticalityPresentIgnore:
				Ngaplog.Trace("Procedure Criticality: Ignore")
			case ngapType.CriticalityPresentNotify:
				Ngaplog.Trace("Procedure Criticality: Notify")
			}
		}
		if criticalityDiagnostics.IEsCriticalityDiagnostics != nil {
			for _, ieCriticalityDiagnostics := range criticalityDiagnostics.IEsCriticalityDiagnostics.List {
				Ngaplog.Tracef("IE ID: %d", ieCriticalityDiagnostics.IEID.Value)

				switch ieCriticalityDiagnostics.IECriticality.Value {
				case ngapType.CriticalityPresentReject:
					Ngaplog.Trace("Criticality Reject")
				case ngapType.CriticalityPresentNotify:
					Ngaplog.Trace("Criticality Notify")
				}

				switch ieCriticalityDiagnostics.TypeOfError.Value {
				case ngapType.TypeOfErrorPresentNotUnderstood:
					Ngaplog.Trace("Type of error: Not understood")
				case ngapType.TypeOfErrorPresentMissing:
					Ngaplog.Trace("Type of error: Missing")
				}
			}
		}
	}
	ranUe := ran.RanUeFindByRanUeNgapID(rANUENGAPID.Value)
	if ranUe == nil {
		Ngaplog.Errorf("No UE Context[RanUeNgapID: %d]", rANUENGAPID.Value)
		return
	}
	amfUe := ranUe.AmfUe
	if amfUe == nil {
		Ngaplog.Error("amfUe is nil")
		return
	}

	if pDUSessionResourceFailedToSetupList != nil {
		Ngaplog.Trace("[NGAP] Send PDUSessionResourceSetupUnsuccessfulTransfer to SMF")

		for _, item := range pDUSessionResourceFailedToSetupList.List {
			pduSessionID := int32(item.PDUSessionID.Value)
			transfer := item.PDUSessionResourceSetupUnsuccessfulTransfer

			response, _, _, err := amf_consumer.SendUpdateSmContextN2Info(amfUe, pduSessionID, models.N2SmInfoType_PDU_RES_SETUP_FAIL, transfer)
			if err != nil {
				Ngaplog.Errorf("SendUpdateSmContextN2Info[PDUSessionResourceSetupUnsuccessfulTransfer] Error:\n%s", err.Error())
			}

			if response != nil && response.BinaryDataN2SmInformation != nil {
				// TODO: n2SmInfo send to RAN
			} else if response == nil {
				// TODO: error handling
			}
		}
	}
}

func HandleUEContextReleaseRequest(ran *amf_context.AmfRan, message *ngapType.NGAPPDU) {

	var aMFUENGAPID *ngapType.AMFUENGAPID
	var rANUENGAPID *ngapType.RANUENGAPID
	var pDUSessionResourceList *ngapType.PDUSessionResourceListCxtRelReq
	var cause *ngapType.Cause

	if ran == nil {
		Ngaplog.Error("ran is nil")
		return
	}
	if message == nil {
		Ngaplog.Error("NGAP Message is nil")
		return
	}
	initiatingMessage := message.InitiatingMessage
	if initiatingMessage == nil {
		Ngaplog.Error("InitiatingMessage is nil")
		return
	}
	uEContextReleaseRequest := initiatingMessage.Value.UEContextReleaseRequest
	if uEContextReleaseRequest == nil {
		Ngaplog.Error("UEContextReleaseRequest is nil")
		return
	}

	Ngaplog.Info("[AMF] UE Context Release Request")

	for _, ie := range uEContextReleaseRequest.ProtocolIEs.List {
		switch ie.Id.Value {
		case ngapType.ProtocolIEIDAMFUENGAPID:
			aMFUENGAPID = ie.Value.AMFUENGAPID
			Ngaplog.Trace("[NGAP] Decode IE AmfUeNgapID")
			if aMFUENGAPID == nil {
				Ngaplog.Error("AmfUeNgapID is nil")
				return
			}
		case ngapType.ProtocolIEIDRANUENGAPID:
			rANUENGAPID = ie.Value.RANUENGAPID
			Ngaplog.Trace("[NGAP] Decode IE RanUeNgapID")
			if rANUENGAPID == nil {
				Ngaplog.Error("RanUeNgapID is nil")
				return
			}
		case ngapType.ProtocolIEIDPDUSessionResourceListCxtRelReq:
			pDUSessionResourceList = ie.Value.PDUSessionResourceListCxtRelReq
			Ngaplog.Trace("[NGAP] Decode IE Pdu Session Resource List")
		case ngapType.ProtocolIEIDCause:
			cause = ie.Value.Cause
			Ngaplog.Trace("[NGAP] Decode IE Cause")
			if cause == nil {
				Ngaplog.Warn("Cause is nil")
			}
		}
	}

	printRanInfo(ran)

	ranUe := amf_context.AMF_Self().RanUeFindByAmfUeNgapID(aMFUENGAPID.Value)
	if ranUe == nil {
		Ngaplog.Errorf("No RanUe Context[AmfUeNgapID: %d]", aMFUENGAPID.Value)
		cause = &ngapType.Cause{
			Present: ngapType.CausePresentRadioNetwork,
			RadioNetwork: &ngapType.CauseRadioNetwork{
				Value: ngapType.CauseRadioNetworkPresentUnknownLocalUENGAPID,
			},
		}
		ngap_message.SendErrorIndication(ran, nil, nil, cause, nil)
		return
	}

	Ngaplog.Tracef("RanUeNgapID[%d] AmfUeNgapID[%d]", ranUe.RanUeNgapId, ranUe.AmfUeNgapId)

	causeGroup := ngapType.CausePresentRadioNetwork
	causeValue := ngapType.CauseRadioNetworkPresentUnspecified
	if cause != nil {
		causeGroup, causeValue = printAndGetCause(cause)
	}

	amfUe := ranUe.AmfUe
	if amfUe != nil {
		causeAll := amf_context.CauseAll{
			NgapCause: &models.NgApCause{
				Group: int32(causeGroup),
				Value: int32(causeValue),
			},
		}
		if amfUe.Sm[ran.AnType].Check(gmm_state.REGISTERED) {
			Ngaplog.Info("[NGAP] Ue Context in GMM-Registered")
			if pDUSessionResourceList != nil {
				for _, pduSessionReourceItem := range pDUSessionResourceList.List {
					pduSessionID := int32(pduSessionReourceItem.PDUSessionID.Value)
					response, _, _, err := amf_consumer.SendUpdateSmContextDeactivateUpCnxState(amfUe, pduSessionID, causeAll)
					if err != nil {
						logger.NgapLog.Errorf("Send Update SmContextDeactivate UpCnxState Error[%s]", err.Error())
					} else if response == nil {
						logger.NgapLog.Errorln("Send Update SmContextDeactivate UpCnxState Error")
					}
				}
			}
		} else {
			Ngaplog.Info("[NGAP] Ue Context in Non GMM-Registered")
			for pduSessionId := range amfUe.SmContextList {
				releaseData := amf_consumer.BuildReleaseSmContextRequest(amfUe, &causeAll, "", nil)
				detail, err := amf_consumer.SendReleaseSmContextRequest(amfUe, pduSessionId, releaseData)
				if err != nil {
					logger.NgapLog.Errorf("Send ReleaseSmContextRequest Error[%s]", err.Error())
				} else if detail != nil {
					logger.NgapLog.Errorf("Send ReleaseSmContextRequeste Error[%s]", detail.Cause)
				}
			}
			ngap_message.SendUEContextReleaseCommand(ranUe, amf_context.UeContextReleaseUeContext, causeGroup, causeValue)
			return
		}
	}
	ngap_message.SendUEContextReleaseCommand(ranUe, amf_context.UeContextN2NormalRelease, causeGroup, causeValue)

}

func HandleUEContextModificationResponse(ran *amf_context.AmfRan, message *ngapType.NGAPPDU) {

	var aMFUENGAPID *ngapType.AMFUENGAPID
	var rANUENGAPID *ngapType.RANUENGAPID
	var rRCState *ngapType.RRCState
	var userLocationInformation *ngapType.UserLocationInformation
	var criticalityDiagnostics *ngapType.CriticalityDiagnostics

	var ranUe *amf_context.RanUe

	if ran == nil {
		Ngaplog.Error("ran is nil")
		return
	}
	if message == nil {
		Ngaplog.Error("NGAP Message is nil")
		return
	}
	successfulOutcome := message.SuccessfulOutcome
	if successfulOutcome == nil {
		Ngaplog.Error("SuccessfulOutcome is nil")
		return
	}
	uEContextModificationResponse := successfulOutcome.Value.UEContextModificationResponse
	if uEContextModificationResponse == nil {
		Ngaplog.Error("UEContextModificationResponse is nil")
		return
	}

	Ngaplog.Info("[AMF] UE Context Modification Response")

	for _, ie := range uEContextModificationResponse.ProtocolIEs.List {
		switch ie.Id.Value {
		case ngapType.ProtocolIEIDAMFUENGAPID: // ignore
			aMFUENGAPID = ie.Value.AMFUENGAPID
			Ngaplog.Trace("[NGAP] Decode IE AmfUeNgapID")
			if aMFUENGAPID == nil {
				Ngaplog.Warn("AmfUeNgapID is nil")
			}
		case ngapType.ProtocolIEIDRANUENGAPID: // ignore
			rANUENGAPID = ie.Value.RANUENGAPID
			Ngaplog.Trace("[NGAP] Decode IE RanUeNgapID")
			if rANUENGAPID == nil {
				Ngaplog.Warn("RanUeNgapID is nil")
			}
		case ngapType.ProtocolIEIDRRCState: // optional, ignore
			rRCState = ie.Value.RRCState
			Ngaplog.Trace("[NGAP] Decode IE RRCState")
		case ngapType.ProtocolIEIDUserLocationInformation: // optional, ignore
			userLocationInformation = ie.Value.UserLocationInformation
			Ngaplog.Trace("[NGAP] Decode IE UserLocationInformation")
		case ngapType.ProtocolIEIDCriticalityDiagnostics: // optional, ignore
			criticalityDiagnostics = ie.Value.CriticalityDiagnostics
			Ngaplog.Trace("[NGAP] Decode IE CriticalityDiagnostics")
		}
	}

	printRanInfo(ran)

	if rANUENGAPID != nil {
		ranUe = ran.RanUeFindByRanUeNgapID(rANUENGAPID.Value)
		if ranUe == nil {
			Ngaplog.Warnf("No UE Context[RanUeNgapID: %d]", rANUENGAPID.Value)
		}
	}

	if aMFUENGAPID != nil {
		ranUe = amf_context.AMF_Self().RanUeFindByAmfUeNgapID(aMFUENGAPID.Value)
		if ranUe == nil {
			Ngaplog.Warnf("No UE Context[AmfUeNgapID: %d]", aMFUENGAPID.Value)
			return
		}
	}

	if ranUe != nil {
		Ngaplog.Tracef("AmfUeNgapID[%d] RanUeNgapID[%d]", ranUe.AmfUeNgapId, ranUe.RanUeNgapId)

		if rRCState != nil {
			switch rRCState.Value {
			case ngapType.RRCStatePresentInactive:
				Ngaplog.Trace("UE RRC State: Inactive")
			case ngapType.RRCStatePresentConnected:
				Ngaplog.Trace("UE RRC State: Connected")
			}
		}

		if userLocationInformation != nil {
			ranUe.UpdateLocation(userLocationInformation)
		}
	}

	if criticalityDiagnostics != nil {
		printCriticalityDiagnostics(criticalityDiagnostics)
	}
}

func HandleUEContextModificationFailure(ran *amf_context.AmfRan, message *ngapType.NGAPPDU) {

	var aMFUENGAPID *ngapType.AMFUENGAPID
	var rANUENGAPID *ngapType.RANUENGAPID
	var cause *ngapType.Cause
	var criticalityDiagnostics *ngapType.CriticalityDiagnostics

	var ranUe *amf_context.RanUe

	if ran == nil {
		Ngaplog.Error("ran is nil")
		return
	}
	if message == nil {
		Ngaplog.Error("NGAP Message is nil")
		return
	}
	unsuccessfulOutcome := message.UnsuccessfulOutcome
	if unsuccessfulOutcome == nil {
		Ngaplog.Error("UnsuccessfulOutcome is nil")
		return
	}
	uEContextModificationFailure := unsuccessfulOutcome.Value.UEContextModificationFailure
	if uEContextModificationFailure == nil {
		Ngaplog.Error("UEContextModificationFailure is nil")
		return
	}

	Ngaplog.Info("[AMF] UE Context Modification Failure")

	for _, ie := range uEContextModificationFailure.ProtocolIEs.List {
		switch ie.Id.Value {
		case ngapType.ProtocolIEIDAMFUENGAPID: // ignore
			aMFUENGAPID = ie.Value.AMFUENGAPID
			Ngaplog.Trace("[NGAP] Decode IE AmfUeNgapID")
			if aMFUENGAPID == nil {
				Ngaplog.Warn("AmfUeNgapID is nil")
			}
		case ngapType.ProtocolIEIDRANUENGAPID: // ignore
			rANUENGAPID = ie.Value.RANUENGAPID
			Ngaplog.Trace("[NGAP] Decode IE RanUeNgapID")
			if rANUENGAPID == nil {
				Ngaplog.Warn("RanUeNgapID is nil")
			}
		case ngapType.ProtocolIEIDCause: // ignore
			cause = ie.Value.Cause
			Ngaplog.Trace("[NGAP] Decode IE Cause")
			if cause == nil {
				Ngaplog.Warn("Cause is nil")
			}
		case ngapType.ProtocolIEIDCriticalityDiagnostics: // optional, ignore
			criticalityDiagnostics = ie.Value.CriticalityDiagnostics
			Ngaplog.Trace("[NGAP] Decode IE CriticalityDiagnostics")
		}
	}

	printRanInfo(ran)

	if rANUENGAPID != nil {
		ranUe = ran.RanUeFindByRanUeNgapID(rANUENGAPID.Value)
		if ranUe == nil {
			Ngaplog.Warnf("No UE Context[RanUeNgapID: %d]", rANUENGAPID.Value)
		}
	}

	if aMFUENGAPID != nil {
		ranUe = amf_context.AMF_Self().RanUeFindByAmfUeNgapID(aMFUENGAPID.Value)
		if ranUe == nil {
			Ngaplog.Warnf("No UE Context[AmfUeNgapID: %d]", aMFUENGAPID.Value)
		}
	}

	if ranUe != nil {
		Ngaplog.Tracef("AmfUeNgapID[%d] RanUeNgapID[%d]", ranUe.AmfUeNgapId, ranUe.RanUeNgapId)
	}

	if cause != nil {
		printAndGetCause(cause)
	}

	if criticalityDiagnostics != nil {
		printCriticalityDiagnostics(criticalityDiagnostics)
	}
}

func HandleRRCInactiveTransitionReport(ran *amf_context.AmfRan, message *ngapType.NGAPPDU) {

	var aMFUENGAPID *ngapType.AMFUENGAPID
	var rANUENGAPID *ngapType.RANUENGAPID
	var rRCState *ngapType.RRCState
	var userLocationInformation *ngapType.UserLocationInformation

	logger.SetLogLevel(logrus.TraceLevel)
	logger.SetReportCaller(false)

	if ran == nil {
		logger.NgapLog.Error("ran is nil")
		return
	}
	if message == nil {
		logger.NgapLog.Error("NGAP Message is nil")
		return
	}

	initiatingMessage := message.InitiatingMessage
	if initiatingMessage == nil {
		logger.NgapLog.Error("Initiating Message is nil")
		return
	}

	rRCInactiveTransitionReport := initiatingMessage.Value.RRCInactiveTransitionReport
	if rRCInactiveTransitionReport == nil {
		logger.NgapLog.Error("RRCInactiveTransitionReport is nil")
		return
	}
	logger.NgapLog.Info("[AMF] RRC Inactive Transition Report")

	for i := 0; i < len(rRCInactiveTransitionReport.ProtocolIEs.List); i++ {
		ie := rRCInactiveTransitionReport.ProtocolIEs.List[i]
		switch ie.Id.Value {
		case ngapType.ProtocolIEIDAMFUENGAPID: //reject
			aMFUENGAPID = ie.Value.AMFUENGAPID
			logger.NgapLog.Trace("[NGAP] Decode IE AmfUeNgapID")
			if aMFUENGAPID == nil {
				logger.NgapLog.Error("AmfUeNgapID is nil")
				return
			}
		case ngapType.ProtocolIEIDRANUENGAPID: //reject
			rANUENGAPID = ie.Value.RANUENGAPID
			logger.NgapLog.Trace("[NGAP] Decode IE RanUeNgapID")
			if rANUENGAPID == nil {
				logger.NgapLog.Error("RanUeNgapID is nil")
				return
			}
		case ngapType.ProtocolIEIDRRCState: //ignore
			rRCState = ie.Value.RRCState
			logger.NgapLog.Trace("[NGAP] Decode IE RRCState")
			if rRCState == nil {
				logger.NgapLog.Error("RRCState is nil")
				return
			}
		case ngapType.ProtocolIEIDUserLocationInformation: //ignore
			userLocationInformation = ie.Value.UserLocationInformation
			logger.NgapLog.Trace("[NGAP] Decode IE UserLocationInformation")
			if userLocationInformation == nil {
				logger.NgapLog.Error("UserLocationInformation is nil")
				return
			}
		}
	}

	printRanInfo(ran)

	ranUe := ran.RanUeFindByRanUeNgapID(rANUENGAPID.Value)
	if ranUe == nil {
		Ngaplog.Warnf("No UE Context[RanUeNgapID: %d]", rANUENGAPID.Value)
	} else {
		logger.NgapLog.Tracef("RANUENGAPID[%d] AMFUENGAPID[%d]", ranUe.RanUeNgapId, ranUe.AmfUeNgapId)

		if rRCState != nil {
			switch rRCState.Value {
			case ngapType.RRCStatePresentInactive:
				Ngaplog.Trace("UE RRC State: Inactive")
			case ngapType.RRCStatePresentConnected:
				Ngaplog.Trace("UE RRC State: Connected")
			}
		}
		ranUe.UpdateLocation(userLocationInformation)

	}
}

// TODO
func HandleHandoverNotify(ran *amf_context.AmfRan, message *ngapType.NGAPPDU) {

	var aMFUENGAPID *ngapType.AMFUENGAPID
	var rANUENGAPID *ngapType.RANUENGAPID
	var userLocationInformation *ngapType.UserLocationInformation

	if ran == nil {
		logger.NgapLog.Error("ran is nil")
		return
	}
	if message == nil {
		logger.NgapLog.Error("NGAP Message is nil")
		return
	}

	initiatingMessage := message.InitiatingMessage
	if initiatingMessage == nil {
		logger.NgapLog.Error("Initiating Message is nil")
		return
	}
	HandoverNotify := initiatingMessage.Value.HandoverNotify
	if HandoverNotify == nil {
		logger.NgapLog.Error("HandoverNotify is nil")
		return
	}

	logger.NgapLog.Info("[AMF] Handover notification")

	for i := 0; i < len(HandoverNotify.ProtocolIEs.List); i++ {
		ie := HandoverNotify.ProtocolIEs.List[i]
		switch ie.Id.Value {
		case ngapType.ProtocolIEIDAMFUENGAPID:
			aMFUENGAPID = ie.Value.AMFUENGAPID
			logger.NgapLog.Trace("[NGAP] Decode IE AmfUeNgapID")
			if aMFUENGAPID == nil {
				logger.NgapLog.Error("AMFUENGAPID is nil")
				return
			}
		case ngapType.ProtocolIEIDRANUENGAPID:
			rANUENGAPID = ie.Value.RANUENGAPID
			logger.NgapLog.Trace("[NGAP] Decode IE RanUeNgapID")
			if rANUENGAPID == nil {
				logger.NgapLog.Error("RANUENGAPID is nil")
				return
			}
		case ngapType.ProtocolIEIDUserLocationInformation:
			userLocationInformation = ie.Value.UserLocationInformation
			logger.NgapLog.Trace("[NGAP] Decode IE userLocationInformation")
			if userLocationInformation == nil {
				logger.NgapLog.Error("userLocationInformation is nil")
				return
			}
		}
	}

	printRanInfo(ran)

	targetUe := ran.RanUeFindByRanUeNgapID(rANUENGAPID.Value)

	if targetUe == nil {
		logger.NgapLog.Errorf("No RanUe Context[AmfUeNgapID: %d]", aMFUENGAPID.Value)
		cause := ngapType.Cause{
			Present: ngapType.CausePresentRadioNetwork,
			RadioNetwork: &ngapType.CauseRadioNetwork{
				Value: ngapType.CauseRadioNetworkPresentUnknownLocalUENGAPID,
			},
		}
		ngap_message.SendErrorIndication(ran, nil, nil, &cause, nil)
		return
	}

	if userLocationInformation != nil {
		targetUe.UpdateLocation(userLocationInformation)
	}
	amfUe := targetUe.AmfUe
	if amfUe == nil {
		Ngaplog.Error("AmfUe is nil")
		return
	}
	sourceUe := targetUe.SourceUe
	if sourceUe == nil {
		// TODO: Send to S-AMF
		// Desciibed in (23.502 4.9.1.3.3) [conditional] 6a.Namf_Communication_N2InfoNotify.
		Ngaplog.Error("N2 Handover between AMF has not been implemented yet")
	} else {
		logger.NgapLog.Info("[AMF] Handover notification Finshed ")
		for _, pduSessionid := range targetUe.SuccessPduSessionId {
			_, _, _, err := amf_consumer.SendUpdateSmContextN2HandoverComplete(amfUe, pduSessionid, "", nil)
			if err != nil {
				Ngaplog.Errorf("Send UpdateSmContextN2HandoverComplete Error[%s]", err.Error())
			}
		}
		amfUe.AttachRanUe(targetUe)
		ngap_message.SendUEContextReleaseCommand(sourceUe, amf_context.UeContextReleaseHandover, ngapType.CausePresentNas, ngapType.CauseNasPresentNormalRelease)

	}

	// TODO: The UE initiates Mobility Registration Update procedure as described in clause 4.2.2.2.2.

}

// TS 23.502 4.9.1
func HandlePathSwitchRequest(ran *amf_context.AmfRan, message *ngapType.NGAPPDU) {

	var rANUENGAPID *ngapType.RANUENGAPID
	var sourceAMFUENGAPID *ngapType.AMFUENGAPID
	var userLocationInformation *ngapType.UserLocationInformation
	var uESecurityCapabilities *ngapType.UESecurityCapabilities
	var pduSessionResourceToBeSwitchedInDLList *ngapType.PDUSessionResourceToBeSwitchedDLList
	var pduSessionResourceFailedToSetupList *ngapType.PDUSessionResourceFailedToSetupListPSReq

	var ranUe *amf_context.RanUe

	if ran == nil {
		Ngaplog.Error("ran is nil")
		return
	}
	if message == nil {
		Ngaplog.Error("NGAP Message is nil")
		return
	}
	initiatingMessage := message.InitiatingMessage
	if initiatingMessage == nil {
		Ngaplog.Error("InitiatingMessage is nil")
		return
	}
	pathSwitchRequest := initiatingMessage.Value.PathSwitchRequest
	if pathSwitchRequest == nil {
		Ngaplog.Error("PathSwitchRequest is nil")
		return
	}

	Ngaplog.Info("[AMF] Path Switch Request")

	for _, ie := range pathSwitchRequest.ProtocolIEs.List {
		switch ie.Id.Value {
		case ngapType.ProtocolIEIDRANUENGAPID: // reject
			rANUENGAPID = ie.Value.RANUENGAPID
			Ngaplog.Trace("[NGAP] Decode IE RanUeNgapID")
			if rANUENGAPID == nil {
				Ngaplog.Error("RanUeNgapID is nil")
				return
			}
		case ngapType.ProtocolIEIDSourceAMFUENGAPID: // reject
			sourceAMFUENGAPID = ie.Value.SourceAMFUENGAPID
			Ngaplog.Trace("[NGAP] Decode IE SourceAmfUeNgapID")
			if sourceAMFUENGAPID == nil {
				Ngaplog.Error("SourceAmfUeNgapID is nil")
				return
			}
		case ngapType.ProtocolIEIDUserLocationInformation: // ignore
			userLocationInformation = ie.Value.UserLocationInformation
			Ngaplog.Trace("[NGAP] Decode IE UserLocationInformation")
		case ngapType.ProtocolIEIDUESecurityCapabilities: // ignore
			uESecurityCapabilities = ie.Value.UESecurityCapabilities
			Ngaplog.Trace("[NGAP] Decode IE UESecurityCapabilities")
		case ngapType.ProtocolIEIDPDUSessionResourceToBeSwitchedDLList: // reject
			pduSessionResourceToBeSwitchedInDLList = ie.Value.PDUSessionResourceToBeSwitchedDLList
			Ngaplog.Trace("[NGAP] Decode IE PDUSessionResourceToBeSwitchedDLList")
			if pduSessionResourceToBeSwitchedInDLList == nil {
				Ngaplog.Error("PDUSessionResourceToBeSwitchedDLList is nil")
				return
			}
		case ngapType.ProtocolIEIDPDUSessionResourceFailedToSetupListPSReq: // ignore
			pduSessionResourceFailedToSetupList = ie.Value.PDUSessionResourceFailedToSetupListPSReq
			Ngaplog.Trace("[NGAP] Decode IE PDUSessionResourceFailedToSetupListPSReq")
		}
	}

	printRanInfo(ran)

	if sourceAMFUENGAPID == nil {
		Ngaplog.Error("SourceAmfUeNgapID is nil")
		return
	}
	ranUe = amf_context.AMF_Self().RanUeFindByAmfUeNgapID(sourceAMFUENGAPID.Value)
	if ranUe == nil {
		Ngaplog.Errorf("Cannot find UE from sourceAMfUeNgapID[%d]", sourceAMFUENGAPID.Value)
		ngap_message.SendPathSwitchRequestFailure(ran, sourceAMFUENGAPID.Value, rANUENGAPID.Value, nil, nil)
		return
	}

	Ngaplog.Tracef("AmfUeNgapID[%d] RanUeNgapID[%d]", ranUe.AmfUeNgapId, ranUe.RanUeNgapId)

	amfUe := ranUe.AmfUe
	if amfUe == nil {
		Ngaplog.Error("AmfUe is nil")
		ngap_message.SendPathSwitchRequestFailure(ran, sourceAMFUENGAPID.Value, rANUENGAPID.Value, nil, nil)
		return
	}

	if amfUe.SecurityContextIsValid() {
		// Update NH
		amfUe.UpdateNH()
	} else {
		Ngaplog.Errorf("No Security Context : SUPI[%s]", amfUe.Supi)
		ngap_message.SendPathSwitchRequestFailure(ran, sourceAMFUENGAPID.Value, rANUENGAPID.Value, nil, nil)
		return
	}

	if uESecurityCapabilities != nil {
		copy(amfUe.SecurityCapabilities.NREncryptionAlgorithms[:2], uESecurityCapabilities.NRencryptionAlgorithms.Value.Bytes)
		copy(amfUe.SecurityCapabilities.NRIntegrityProtectionAlgorithms[:2], uESecurityCapabilities.NRintegrityProtectionAlgorithms.Value.Bytes)
		copy(amfUe.SecurityCapabilities.EUTRAEncryptionAlgorithms[:2], uESecurityCapabilities.EUTRAencryptionAlgorithms.Value.Bytes)
		copy(amfUe.SecurityCapabilities.EUTRAIntegrityProtectionAlgorithms[:2], uESecurityCapabilities.EUTRAintegrityProtectionAlgorithms.Value.Bytes)
	}

	if rANUENGAPID != nil {
		ranUe.RanUeNgapId = rANUENGAPID.Value
	}

	ranUe.UpdateLocation(userLocationInformation)

	var pduSessionResourceSwitchedList ngapType.PDUSessionResourceSwitchedList
	var pduSessionResourceReleasedListPSAck ngapType.PDUSessionResourceReleasedListPSAck
	var pduSessionResourceReleasedListPSFail ngapType.PDUSessionResourceReleasedListPSFail

	if pduSessionResourceToBeSwitchedInDLList != nil {
		for _, item := range pduSessionResourceToBeSwitchedInDLList.List {
			pduSessionID := item.PDUSessionID.Value
			transfer := item.PathSwitchRequestTransfer

			response, errResponse, _, err := amf_consumer.SendUpdateSmContextXnHandover(amfUe, int32(pduSessionID), models.N2SmInfoType_PATH_SWITCH_REQ, transfer)
			if err != nil {
				Ngaplog.Errorf("SendUpdateSmContextXnHandover[PathSwitchRequestTransfer] Error:\n%s", err.Error())
			}
			if response != nil && response.BinaryDataN2SmInformation != nil {
				pduSessionResourceSwitchedItem := ngapType.PDUSessionResourceSwitchedItem{}
				pduSessionResourceSwitchedItem.PDUSessionID.Value = pduSessionID
				pduSessionResourceSwitchedItem.PathSwitchRequestAcknowledgeTransfer = response.BinaryDataN2SmInformation
				pduSessionResourceSwitchedList.List = append(pduSessionResourceSwitchedList.List, pduSessionResourceSwitchedItem)
			}
			if errResponse != nil && errResponse.BinaryDataN2SmInformation != nil {
				pduSessionResourceReleasedItem := ngapType.PDUSessionResourceReleasedItemPSFail{}
				pduSessionResourceReleasedItem.PDUSessionID.Value = pduSessionID
				pduSessionResourceReleasedItem.PathSwitchRequestUnsuccessfulTransfer = errResponse.BinaryDataN2SmInformation
				pduSessionResourceReleasedListPSFail.List = append(pduSessionResourceReleasedListPSFail.List, pduSessionResourceReleasedItem)
			}
		}
	}

	if pduSessionResourceFailedToSetupList != nil {
		for _, item := range pduSessionResourceFailedToSetupList.List {
			pduSessionID := item.PDUSessionID.Value
			transfer := item.PathSwitchRequestSetupFailedTransfer

			response, errResponse, _, err := amf_consumer.SendUpdateSmContextXnHandoverFailed(amfUe, int32(pduSessionID), models.N2SmInfoType_PATH_SWITCH_SETUP_FAIL, transfer)
			if err != nil {
				Ngaplog.Errorf("SendUpdateSmContextXnHandoverFailed[PathSwitchRequestSetupFailedTransfer] Error:\n%s", err.Error())
			}
			if response != nil && response.BinaryDataN2SmInformation != nil {
				pduSessionResourceReleasedItem := ngapType.PDUSessionResourceReleasedItemPSAck{}
				pduSessionResourceReleasedItem.PDUSessionID.Value = pduSessionID
				pduSessionResourceReleasedItem.PathSwitchRequestUnsuccessfulTransfer = response.BinaryDataN2SmInformation
				pduSessionResourceReleasedListPSAck.List = append(pduSessionResourceReleasedListPSAck.List, pduSessionResourceReleasedItem)
			}
			if errResponse != nil && errResponse.BinaryDataN2SmInformation != nil {
				pduSessionResourceReleasedItem := ngapType.PDUSessionResourceReleasedItemPSFail{}
				pduSessionResourceReleasedItem.PDUSessionID.Value = pduSessionID
				pduSessionResourceReleasedItem.PathSwitchRequestUnsuccessfulTransfer = errResponse.BinaryDataN2SmInformation
				pduSessionResourceReleasedListPSFail.List = append(pduSessionResourceReleasedListPSFail.List, pduSessionResourceReleasedItem)
			}
		}
	}

	// TS 23.502 4.9.1.2.2 step 7: send ack to Target NG-RAN. If none of the requested PDU Sessions have been switched successfully,
	// the AMF shall send an N2 Path Switch Request Failure message to the Target NG-RAN
	if len(pduSessionResourceSwitchedList.List) > 0 {
		// TODO: set newSecurityContextIndicator to true if there is a new security context
		err := ranUe.SwitchToRan(ran, rANUENGAPID.Value)
		if err != nil {
			Ngaplog.Error(err.Error())
			return
		}
		ngap_message.SendPathSwitchRequestAcknowledge(ranUe, pduSessionResourceSwitchedList, pduSessionResourceReleasedListPSAck, false, nil, nil, nil)
	} else if len(pduSessionResourceReleasedListPSFail.List) > 0 {
		ngap_message.SendPathSwitchRequestFailure(ran, sourceAMFUENGAPID.Value, rANUENGAPID.Value, &pduSessionResourceReleasedListPSFail, nil)
	} else {
		ngap_message.SendPathSwitchRequestFailure(ran, sourceAMFUENGAPID.Value, rANUENGAPID.Value, nil, nil)
	}
}

func HandleHandoverRequestAcknowledge(ran *amf_context.AmfRan, message *ngapType.NGAPPDU) {

	var aMFUENGAPID *ngapType.AMFUENGAPID
	var rANUENGAPID *ngapType.RANUENGAPID
	var pDUSessionResourceAdmittedList *ngapType.PDUSessionResourceAdmittedList
	var pDUSessionResourceFailedToSetupListHOAck *ngapType.PDUSessionResourceFailedToSetupListHOAck
	var targetToSourceTransparentContainer *ngapType.TargetToSourceTransparentContainer
	var criticalityDiagnostics *ngapType.CriticalityDiagnostics

	var iesCriticalityDiagnostics ngapType.CriticalityDiagnosticsIEList

	if ran == nil {
		logger.NgapLog.Error("ran is nil")
		return
	}
	if message == nil {
		logger.NgapLog.Error("NGAP Message is nil")
		return
	}
	successfulOutcome := message.SuccessfulOutcome
	if successfulOutcome == nil {
		logger.NgapLog.Error("SuccessfulOutcome is nil")
		return
	}
	handoverRequestAcknowledge := successfulOutcome.Value.HandoverRequestAcknowledge // reject
	if handoverRequestAcknowledge == nil {
		logger.NgapLog.Error("HandoverRequestAcknowledge is nil")
		return
	}

	logger.NgapLog.Info("[AMF] Handover Request Acknowledge")

	for _, ie := range handoverRequestAcknowledge.ProtocolIEs.List {
		switch ie.Id.Value {
		case ngapType.ProtocolIEIDAMFUENGAPID: // ignore
			aMFUENGAPID = ie.Value.AMFUENGAPID
			Ngaplog.Trace("[NGAP] Decode IE AmfUeNgapID")
		case ngapType.ProtocolIEIDRANUENGAPID: // ignore
			rANUENGAPID = ie.Value.RANUENGAPID
			Ngaplog.Trace("[NGAP] Decode IE RanUeNgapID")
		case ngapType.ProtocolIEIDPDUSessionResourceAdmittedList: // ignore
			pDUSessionResourceAdmittedList = ie.Value.PDUSessionResourceAdmittedList
			Ngaplog.Trace("[NGAP] Decode IE PduSessionResourceAdmittedList")
		case ngapType.ProtocolIEIDPDUSessionResourceFailedToSetupListHOAck: // ignore
			pDUSessionResourceFailedToSetupListHOAck = ie.Value.PDUSessionResourceFailedToSetupListHOAck
			Ngaplog.Trace("[NGAP] Decode IE PduSessionResourceFailedToSetupListHOAck")
		case ngapType.ProtocolIEIDTargetToSourceTransparentContainer: // reject
			targetToSourceTransparentContainer = ie.Value.TargetToSourceTransparentContainer
			Ngaplog.Trace("[NGAP] Decode IE TargetToSourceTransparentContainer")

		case ngapType.ProtocolIEIDCriticalityDiagnostics: // ignore
			criticalityDiagnostics = ie.Value.CriticalityDiagnostics
			Ngaplog.Trace("[NGAP] Decode IE CriticalityDiagnostics")
		}
	}
	if targetToSourceTransparentContainer == nil {
		Ngaplog.Error("TargetToSourceTransparentContainer is nil")
		item := buildCriticalityDiagnosticsIEItem(ngapType.CriticalityPresentReject, ngapType.ProtocolIEIDTargetToSourceTransparentContainer, ngapType.TypeOfErrorPresentMissing)
		iesCriticalityDiagnostics.List = append(iesCriticalityDiagnostics.List, item)
	}
	if len(iesCriticalityDiagnostics.List) > 0 {
		Ngaplog.Error("Has missing reject IE(s)")

		procedureCode := ngapType.ProcedureCodeHandoverResourceAllocation
		triggeringMessage := ngapType.TriggeringMessagePresentSuccessfulOutcome
		procedureCriticality := ngapType.CriticalityPresentReject
		criticalityDiagnostics := buildCriticalityDiagnostics(&procedureCode, &triggeringMessage, &procedureCriticality, &iesCriticalityDiagnostics)
		ngap_message.SendErrorIndication(ran, nil, nil, nil, &criticalityDiagnostics)
	}

	printRanInfo(ran)

	if criticalityDiagnostics != nil {
		printCriticalityDiagnostics(criticalityDiagnostics)
	}

	targetUe := amf_context.AMF_Self().RanUeFindByAmfUeNgapID(aMFUENGAPID.Value)
	if targetUe == nil {
		Ngaplog.Errorf("No UE Context[AMFUENGAPID: %d]", aMFUENGAPID.Value)
		return
	}

	if rANUENGAPID != nil {
		targetUe.RanUeNgapId = rANUENGAPID.Value
	}
	Ngaplog.Debugf("Target Ue RanUeNgapID[%d] AmfUeNgapID[%d]", targetUe.RanUeNgapId, targetUe.AmfUeNgapId)

	amfUe := targetUe.AmfUe
	if amfUe == nil {
		Ngaplog.Error("amfUe is nil")
		return
	}

	var pduSessionResourceHandoverList ngapType.PDUSessionResourceHandoverList
	var pduSessionResourceToReleaseList ngapType.PDUSessionResourceToReleaseListHOCmd

	// describe in 23.502 4.9.1.3.2 step11
	Ngaplog.Debugf("[AMF] Update Sm Context Request")

	for _, item := range pDUSessionResourceAdmittedList.List {
		pduSessionID := item.PDUSessionID.Value
		transfer := item.HandoverRequestAcknowledgeTransfer
		pduSessionId := int32(pduSessionID)
		if _, exist := amfUe.SmContextList[pduSessionId]; exist {
			response, errResponse, problemDetails, err := amf_consumer.SendUpdateSmContextN2HandoverPrepared(amfUe, pduSessionId, models.N2SmInfoType_HANDOVER_REQ_ACK, transfer)
			if err != nil {
				Ngaplog.Errorf("Send HandoverRequestAcknowledgeTransfer error: %v", err)
			}
			if problemDetails != nil {
				Ngaplog.Warnf("ProblemDetails[status: %d, Cause: %s]", problemDetails.Status, problemDetails.Cause)
			}
			if response != nil && response.BinaryDataN2SmInformation != nil {
				handoverItem := ngapType.PDUSessionResourceHandoverItem{}
				handoverItem.PDUSessionID = item.PDUSessionID
				handoverItem.HandoverCommandTransfer = response.BinaryDataN2SmInformation
				pduSessionResourceHandoverList.List = append(pduSessionResourceHandoverList.List, handoverItem)
				targetUe.SuccessPduSessionId = append(targetUe.SuccessPduSessionId, pduSessionId)
			}
			if errResponse != nil && errResponse.BinaryDataN2SmInformation != nil {
				releaseItem := ngapType.PDUSessionResourceToReleaseItemHOCmd{}
				releaseItem.PDUSessionID = item.PDUSessionID
				releaseItem.HandoverPreparationUnsuccessfulTransfer = errResponse.BinaryDataN2SmInformation
				pduSessionResourceToReleaseList.List = append(pduSessionResourceToReleaseList.List, releaseItem)
			}
		}

	}

	for _, item := range pDUSessionResourceFailedToSetupListHOAck.List {
		pduSessionID := item.PDUSessionID.Value
		transfer := item.HandoverResourceAllocationUnsuccessfulTransfer
		pduSessionId := int32(pduSessionID)
		if _, exist := amfUe.SmContextList[pduSessionId]; exist {
			_, _, problemDetails, err := amf_consumer.SendUpdateSmContextN2HandoverPrepared(amfUe, pduSessionId, models.N2SmInfoType_HANDOVER_RES_ALLOC_FAIL, transfer)
			if err != nil {
				Ngaplog.Errorf("Send HandoverResourceAllocationUnsuccessfulTransfer error: %v", err)
			}
			if problemDetails != nil {
				Ngaplog.Warnf("ProblemDetails[status: %d, Cause: %s]", problemDetails.Status, problemDetails.Cause)
			}
		}
	}

	sourceUe := targetUe.SourceUe
	if sourceUe == nil {
		// TODO: Send Namf_Communication_CreateUEContext Response to S-AMF
		Ngaplog.Error("handover between different Ue has not been implement yet")
	} else {

		Ngaplog.Tracef("Source: RanUeNgapID[%d] AmfUeNgapID[%d]", sourceUe.RanUeNgapId, sourceUe.AmfUeNgapId)
		Ngaplog.Tracef("Target: RanUeNgapID[%d] AmfUeNgapID[%d]", targetUe.RanUeNgapId, targetUe.AmfUeNgapId)
		if len(pduSessionResourceHandoverList.List) == 0 {
			logger.NgapLog.Info("[AMF] Handover Preparation Failure [HoFailure In Target5GC NgranNode Or TargetSystem]")
			cause := &ngapType.Cause{
				Present: ngapType.CausePresentRadioNetwork,
				RadioNetwork: &ngapType.CauseRadioNetwork{
					Value: ngapType.CauseRadioNetworkPresentHoFailureInTarget5GCNgranNodeOrTargetSystem,
				},
			}
			ngap_message.SendHandoverPreparationFailure(sourceUe, *cause, nil)
			return
		}
		ngap_message.SendHandoverCommand(sourceUe, pduSessionResourceHandoverList, pduSessionResourceToReleaseList, *targetToSourceTransparentContainer, nil)
	}
}

func HandleHandoverFailure(ran *amf_context.AmfRan, message *ngapType.NGAPPDU) {

	var aMFUENGAPID *ngapType.AMFUENGAPID
	var cause *ngapType.Cause
	var targetUe *amf_context.RanUe
	var criticalityDiagnostics *ngapType.CriticalityDiagnostics

	if ran == nil {
		logger.NgapLog.Error("ran is nil")
		return
	}
	if message == nil {
		logger.NgapLog.Error("NGAP Message is nil")
		return
	}

	unsuccessfulOutcome := message.UnsuccessfulOutcome // reject
	if unsuccessfulOutcome == nil {
		logger.NgapLog.Error("Unsuccessful Message is nil")
		return
	}

	handoverFailure := unsuccessfulOutcome.Value.HandoverFailure
	if handoverFailure == nil {
		logger.NgapLog.Error("HandoverFailure is nil")
		return
	}

	for _, ie := range handoverFailure.ProtocolIEs.List {
		switch ie.Id.Value {
		case ngapType.ProtocolIEIDAMFUENGAPID: //ignore
			aMFUENGAPID = ie.Value.AMFUENGAPID
			Ngaplog.Trace("[NGAP] Decode IE AmfUeNgapID")
		case ngapType.ProtocolIEIDCause: //ignore
			cause = ie.Value.Cause
			Ngaplog.Trace("[NGAP] Decode IE Cause")
		case ngapType.ProtocolIEIDCriticalityDiagnostics: // ignore
			criticalityDiagnostics = ie.Value.CriticalityDiagnostics
			Ngaplog.Trace("[NGAP] Decode IE CriticalityDiagnostics")
		}
	}

	printRanInfo(ran)

	causePresent := ngapType.CausePresentRadioNetwork
	causeValue := ngapType.CauseRadioNetworkPresentHoFailureInTarget5GCNgranNodeOrTargetSystem
	if cause != nil {
		causePresent, causeValue = printAndGetCause(cause)
	}

	if criticalityDiagnostics != nil {
		printCriticalityDiagnostics(criticalityDiagnostics)
	}

	targetUe = amf_context.AMF_Self().RanUeFindByAmfUeNgapID(aMFUENGAPID.Value)

	if targetUe == nil {
		logger.NgapLog.Errorf("No UE Context[AmfUENGAPID: %d]", aMFUENGAPID.Value)
		cause := ngapType.Cause{
			Present: ngapType.CausePresentRadioNetwork,
			RadioNetwork: &ngapType.CauseRadioNetwork{
				Value: ngapType.CauseRadioNetworkPresentUnknownLocalUENGAPID,
			},
		}
		ngap_message.SendErrorIndication(ran, nil, nil, &cause, nil)
		return
	}

	sourceUe := targetUe.SourceUe
	if sourceUe == nil {
		// TODO: handle N2 Handover between AMF
		Ngaplog.Error("N2 Handover between AMF has not been implemented yet")
	} else {
		amfUe := targetUe.AmfUe
		if amfUe != nil {
			for pduSessionId := range amfUe.SmContextList {
				causeAll := amf_context.CauseAll{
					NgapCause: &models.NgApCause{
						Group: int32(causePresent),
						Value: int32(causeValue),
					},
				}
				_, _, _, err := amf_consumer.SendUpdateSmContextN2HandoverCanceled(amfUe, pduSessionId, causeAll)
				if err != nil {
					logger.NgapLog.Errorf("Send UpdateSmContextN2HandoverCanceled Error for PduSessionId[%d]", pduSessionId)
				}
			}
		}
		ngap_message.SendHandoverPreparationFailure(sourceUe, *cause, criticalityDiagnostics)
	}

	ngap_message.SendUEContextReleaseCommand(targetUe, amf_context.UeContextReleaseHandover, causePresent, causeValue)
}

func HandleHandoverRequired(ran *amf_context.AmfRan, message *ngapType.NGAPPDU) {

	var aMFUENGAPID *ngapType.AMFUENGAPID
	var rANUENGAPID *ngapType.RANUENGAPID
	var handoverType *ngapType.HandoverType
	var cause *ngapType.Cause
	var targetID *ngapType.TargetID
	var pDUSessionResourceListHORqd *ngapType.PDUSessionResourceListHORqd
	var sourceToTargetTransparentContainer *ngapType.SourceToTargetTransparentContainer
	var iesCriticalityDiagnostics ngapType.CriticalityDiagnosticsIEList

	if ran == nil {
		logger.NgapLog.Error("ran is nil")
		return
	}
	if message == nil {
		logger.NgapLog.Error("NGAP Message is nil")
		return
	}

	initiatingMessage := message.InitiatingMessage
	if initiatingMessage == nil {
		logger.NgapLog.Error("Initiating Message is nil")
		return
	}
	HandoverRequired := initiatingMessage.Value.HandoverRequired
	if HandoverRequired == nil {
		logger.NgapLog.Error("HandoverRequired is nil")
		return
	}

	logger.NgapLog.Info("[AMF] HandoverRequired\n")
	for i := 0; i < len(HandoverRequired.ProtocolIEs.List); i++ {
		ie := HandoverRequired.ProtocolIEs.List[i]
		switch ie.Id.Value {
		case ngapType.ProtocolIEIDAMFUENGAPID:
			aMFUENGAPID = ie.Value.AMFUENGAPID //reject
			logger.NgapLog.Trace("[NGAP] Decode IE AmfUeNgapID")
		case ngapType.ProtocolIEIDRANUENGAPID: //reject
			rANUENGAPID = ie.Value.RANUENGAPID
			logger.NgapLog.Trace("[NGAP] Decode IE RanUeNgapID")
		case ngapType.ProtocolIEIDHandoverType: //reject
			handoverType = ie.Value.HandoverType
			logger.NgapLog.Trace("[NGAP] Decode IE HandoverType")
		case ngapType.ProtocolIEIDCause: //ignore
			cause = ie.Value.Cause
			logger.NgapLog.Trace("[NGAP] Decode IE Cause")
		case ngapType.ProtocolIEIDTargetID: //reject
			targetID = ie.Value.TargetID
			logger.NgapLog.Trace("[NGAP] Decode IE TargetID")
		case ngapType.ProtocolIEIDPDUSessionResourceListHORqd: //reject
			pDUSessionResourceListHORqd = ie.Value.PDUSessionResourceListHORqd
			logger.NgapLog.Trace("[NGAP] Decode IE PDUSessionResourceListHORqd")
		case ngapType.ProtocolIEIDSourceToTargetTransparentContainer: //reject
			sourceToTargetTransparentContainer = ie.Value.SourceToTargetTransparentContainer
			logger.NgapLog.Trace("[NGAP] Decode IE SourceToTargetTransparentContainer")
		}
	}

	printRanInfo(ran)

	if aMFUENGAPID == nil {
		Ngaplog.Error("AmfUeNgapID is nil")
		item := buildCriticalityDiagnosticsIEItem(ngapType.CriticalityPresentReject, ngapType.ProtocolIEIDAMFUENGAPID, ngapType.TypeOfErrorPresentMissing)
		iesCriticalityDiagnostics.List = append(iesCriticalityDiagnostics.List, item)
	}
	if rANUENGAPID == nil {
		Ngaplog.Error("RanUeNgapID is nil")
		item := buildCriticalityDiagnosticsIEItem(ngapType.CriticalityPresentReject, ngapType.ProtocolIEIDRANUENGAPID, ngapType.TypeOfErrorPresentMissing)
		iesCriticalityDiagnostics.List = append(iesCriticalityDiagnostics.List, item)
	}

	if handoverType == nil {
		Ngaplog.Error("handoverType is nil")
		item := buildCriticalityDiagnosticsIEItem(ngapType.CriticalityPresentReject, ngapType.ProtocolIEIDHandoverType, ngapType.TypeOfErrorPresentMissing)
		iesCriticalityDiagnostics.List = append(iesCriticalityDiagnostics.List, item)
	}
	if targetID == nil {
		Ngaplog.Error("targetID is nil")
		item := buildCriticalityDiagnosticsIEItem(ngapType.CriticalityPresentReject, ngapType.ProtocolIEIDTargetID, ngapType.TypeOfErrorPresentMissing)
		iesCriticalityDiagnostics.List = append(iesCriticalityDiagnostics.List, item)
	}
	if pDUSessionResourceListHORqd == nil {
		Ngaplog.Error("pDUSessionResourceListHORqd is nil")
		item := buildCriticalityDiagnosticsIEItem(ngapType.CriticalityPresentReject, ngapType.ProtocolIEIDPDUSessionResourceListHORqd, ngapType.TypeOfErrorPresentMissing)
		iesCriticalityDiagnostics.List = append(iesCriticalityDiagnostics.List, item)
	}
	if sourceToTargetTransparentContainer == nil {
		Ngaplog.Error("sourceToTargetTransparentContainer is nil")
		item := buildCriticalityDiagnosticsIEItem(ngapType.CriticalityPresentReject, ngapType.ProtocolIEIDSourceToTargetTransparentContainer, ngapType.TypeOfErrorPresentMissing)
		iesCriticalityDiagnostics.List = append(iesCriticalityDiagnostics.List, item)
	}

	if len(iesCriticalityDiagnostics.List) > 0 {
		procedureCode := ngapType.ProcedureCodeHandoverPreparation
		triggeringMessage := ngapType.TriggeringMessagePresentInitiatingMessage
		procedureCriticality := ngapType.CriticalityPresentReject
		criticalityDiagnostics := buildCriticalityDiagnostics(&procedureCode, &triggeringMessage, &procedureCriticality, &iesCriticalityDiagnostics)
		ngap_message.SendErrorIndication(ran, nil, nil, nil, &criticalityDiagnostics)
		return
	}

	sourceUe := ran.RanUeFindByRanUeNgapID(rANUENGAPID.Value)
	if sourceUe == nil {
		logger.NgapLog.Errorf("Cannot find UE for RAN_UE_NGAP_ID[%d] ", rANUENGAPID.Value)
		cause := ngapType.Cause{
			Present: ngapType.CausePresentRadioNetwork,
			RadioNetwork: &ngapType.CauseRadioNetwork{
				Value: ngapType.CauseRadioNetworkPresentUnknownLocalUENGAPID,
			},
		}
		ngap_message.SendErrorIndication(ran, nil, nil, &cause, nil)
		return
	}
	amfUe := sourceUe.AmfUe
	if amfUe == nil {
		logger.NgapLog.Error("Cannot find amfUE from sourceUE")
		return
	}

	if targetID.Present != ngapType.TargetIDPresentTargetRANNodeID {
		logger.NgapLog.Errorf("targetID type[%d] is not supported", targetID.Present)
		return
	}
	amfUe.OnGoing[sourceUe.Ran.AnType].Procedure = amf_context.OnGoingProcedureN2Handover
	if !amfUe.SecurityContextIsValid() {
		logger.NgapLog.Info("[AMF] Handover Preparation Failure [Authentication Failure]")
		cause = &ngapType.Cause{
			Present: ngapType.CausePresentNas,
			Nas: &ngapType.CauseNas{
				Value: ngapType.CauseNasPresentAuthenticationFailure,
			},
		}
		ngap_message.SendHandoverPreparationFailure(sourceUe, *cause, nil)
		return
	}
	var aMFSelf = amf_context.AMF_Self()
	targetRanNodeId := ngapConvert.RanIdToModels(targetID.TargetRANNodeID.GlobalRANNodeID)
	targetRan := aMFSelf.AmfRanFindByRanId(targetRanNodeId)
	if targetRan == nil {
		// handover between different AMF
		logger.NgapLog.Warnf("Handover required : cannot find target Ran Node Id[%+v] in this AMF", targetRanNodeId)
		logger.NgapLog.Error("Handover between different AMF has not been implemented yet")
		return
		// TODO: Send to T-AMF
		// Described in (23.502 4.9.1.3.2) step 3.Namf_Communication_CreateUEContext Request

	} else {
		// Handover in same AMF
		sourceUe.HandOverType.Value = handoverType.Value
		tai := ngapConvert.TaiToModels(targetID.TargetRANNodeID.SelectedTAI)
		targetId := models.NgRanTargetId{
			RanNodeId: &targetRanNodeId,
			Tai:       &tai,
		}
		var pduSessionReqList ngapType.PDUSessionResourceSetupListHOReq
		for _, pDUSessionResourceHoItem := range pDUSessionResourceListHORqd.List {
			pduSessionId := int32(pDUSessionResourceHoItem.PDUSessionID.Value)
			if smContext, exist := amfUe.SmContextList[pduSessionId]; exist {
				response, _, _, _ := amf_consumer.SendUpdateSmContextN2HandoverPreparing(amfUe, pduSessionId, models.N2SmInfoType_HANDOVER_REQUIRED, pDUSessionResourceHoItem.HandoverRequiredTransfer, "", &targetId)
				if response == nil {
					logger.NgapLog.Errorf("SendUpdateSmContextN2HandoverPreparing Error for PduSessionId[%d]", pduSessionId)
					continue
				} else if response.BinaryDataN2SmInformation != nil {
					ngap_message.AppendPDUSessionResourceSetupListHOReq(&pduSessionReqList, pduSessionId, *smContext.PduSessionContext.SNssai, response.BinaryDataN2SmInformation)
				}
			}

		}
		if len(pduSessionReqList.List) == 0 {
			logger.NgapLog.Info("[AMF] Handover Preparation Failure [HoFailure In Target5GC NgranNode Or TargetSystem]")
			cause = &ngapType.Cause{
				Present: ngapType.CausePresentRadioNetwork,
				RadioNetwork: &ngapType.CauseRadioNetwork{
					Value: ngapType.CauseRadioNetworkPresentHoFailureInTarget5GCNgranNodeOrTargetSystem,
				},
			}
			ngap_message.SendHandoverPreparationFailure(sourceUe, *cause, nil)
			return
		}
		// Update NH
		amfUe.UpdateNH()
		ngap_message.SendHandoverRequest(sourceUe, targetRan, *cause, pduSessionReqList, *sourceToTargetTransparentContainer, false)
	}

}

func HandleHandoverCancel(ran *amf_context.AmfRan, message *ngapType.NGAPPDU) {

	var aMFUENGAPID *ngapType.AMFUENGAPID
	var rANUENGAPID *ngapType.RANUENGAPID
	var cause *ngapType.Cause

	if ran == nil {
		logger.NgapLog.Error("ran is nil")
		return
	}
	if message == nil {
		logger.NgapLog.Error("NGAP Message is nil")
		return
	}

	initiatingMessage := message.InitiatingMessage
	if initiatingMessage == nil {
		logger.NgapLog.Error("Initiating Message is nil")
		return
	}
	HandoverCancel := initiatingMessage.Value.HandoverCancel
	if HandoverCancel == nil {
		logger.NgapLog.Error("Handover Cancel is nil")
		return
	}

	logger.NgapLog.Info("[AMF] Handover Cancel")
	for i := 0; i < len(HandoverCancel.ProtocolIEs.List); i++ {
		ie := HandoverCancel.ProtocolIEs.List[i]
		switch ie.Id.Value {
		case ngapType.ProtocolIEIDAMFUENGAPID:
			aMFUENGAPID = ie.Value.AMFUENGAPID
			logger.NgapLog.Trace("[NGAP] Decode IE AmfUeNgapID")
			if aMFUENGAPID == nil {
				logger.NgapLog.Error("AMFUENGAPID is nil")
				return
			}
		case ngapType.ProtocolIEIDRANUENGAPID:
			rANUENGAPID = ie.Value.RANUENGAPID
			logger.NgapLog.Trace("[NGAP] Decode IE RanUeNgapID")
			if rANUENGAPID == nil {
				logger.NgapLog.Error("RANUENGAPID is nil")
				return
			}
		case ngapType.ProtocolIEIDCause:
			cause = ie.Value.Cause
			logger.NgapLog.Trace("[NGAP] Decode IE Cause")
			if cause == nil {
				logger.NgapLog.Error(cause, "cause is nil")
				return
			}
		}
	}

	printRanInfo(ran)

	sourceUe := ran.RanUeFindByRanUeNgapID(rANUENGAPID.Value)
	if sourceUe == nil {
		logger.NgapLog.Errorf("No UE Context[RanUeNgapID: %d]", rANUENGAPID.Value)
		cause := ngapType.Cause{
			Present: ngapType.CausePresentRadioNetwork,
			RadioNetwork: &ngapType.CauseRadioNetwork{
				Value: ngapType.CauseRadioNetworkPresentUnknownLocalUENGAPID,
			},
		}
		ngap_message.SendErrorIndication(ran, nil, nil, &cause, nil)
		return
	}

	if sourceUe.AmfUeNgapId != aMFUENGAPID.Value {
		logger.NgapLog.Warnf("Conflict AMF_UE_NGAP_ID : %d != %d", sourceUe.AmfUeNgapId, aMFUENGAPID.Value)
	}
	logger.NgapLog.Tracef("Source: RAN_UE_NGAP_ID[%d] AMF_UE_NGAP_ID[%d]", sourceUe.RanUeNgapId, sourceUe.AmfUeNgapId)

	causePresent := ngapType.CausePresentRadioNetwork
	causeValue := ngapType.CauseRadioNetworkPresentHoFailureInTarget5GCNgranNodeOrTargetSystem
	if cause != nil {
		causePresent, causeValue = printAndGetCause(cause)
	}
	targetUe := sourceUe.TargetUe
	if targetUe == nil {
		// Described in (23.502 4.11.1.2.3) step 2
		// Todo : send to T-AMF invoke Namf_UeContextReleaseRequest(targetUe)
		Ngaplog.Error("N2 Handover between AMF has not been implemented yet")

	} else {
		logger.NgapLog.Tracef("Target : RAN_UE_NGAP_ID[%d] AMF_UE_NGAP_ID[%d]", targetUe.RanUeNgapId, targetUe.AmfUeNgapId)
		amfUe := sourceUe.AmfUe
		if amfUe != nil {
			for pduSessionId := range amfUe.SmContextList {
				causeAll := amf_context.CauseAll{
					NgapCause: &models.NgApCause{
						Group: int32(causePresent),
						Value: int32(causeValue),
					},
				}
				_, _, _, err := amf_consumer.SendUpdateSmContextN2HandoverCanceled(amfUe, pduSessionId, causeAll)
				if err != nil {
					logger.NgapLog.Errorf("Send UpdateSmContextN2HandoverCanceled Error for PduSessionId[%d]", pduSessionId)
				}
			}
		}
		ngap_message.SendUEContextReleaseCommand(targetUe, amf_context.UeContextReleaseHandover, causePresent, causeValue)
		ngap_message.SendHandoverCancelAcknowledge(sourceUe, nil)
	}

}

func HandleUplinkRanStatusTransfer(ran *amf_context.AmfRan, message *ngapType.NGAPPDU) {
	var aMFUENGAPID *ngapType.AMFUENGAPID
	var rANUENGAPID *ngapType.RANUENGAPID
	var rANStatusTransferTransparentContainer *ngapType.RANStatusTransferTransparentContainer
	var ranUe *amf_context.RanUe

	if ran == nil {
		Ngaplog.Error("ran is nil")
		return
	}
	if message == nil {
		Ngaplog.Error("NGAP Message is nil")
		return
	}
	initiatingMessage := message.InitiatingMessage // ignore
	if initiatingMessage == nil {
		Ngaplog.Error("InitiatingMessage is nil")
		return
	}
	uplinkRanStatusTransfer := initiatingMessage.Value.UplinkRANStatusTransfer
	if uplinkRanStatusTransfer == nil {
		Ngaplog.Error("UplinkRanStatusTransfer is nil")
		return
	}

	Ngaplog.Info("[AMF] Uplink Ran Status Transfer")

	for _, ie := range uplinkRanStatusTransfer.ProtocolIEs.List {
		switch ie.Id.Value {
		case ngapType.ProtocolIEIDAMFUENGAPID: // reject
			aMFUENGAPID = ie.Value.AMFUENGAPID
			Ngaplog.Trace("[NGAP] Decode IE AmfUeNgapID")
			if aMFUENGAPID == nil {
				Ngaplog.Error("AmfUeNgapID is nil")

			}
		case ngapType.ProtocolIEIDRANUENGAPID: // reject
			rANUENGAPID = ie.Value.RANUENGAPID
			Ngaplog.Trace("[NGAP] Decode IE RanUeNgapID")
			if rANUENGAPID == nil {
				Ngaplog.Error("RanUeNgapID is nil")

			}
		case ngapType.ProtocolIEIDRANStatusTransferTransparentContainer: // reject
			rANStatusTransferTransparentContainer = ie.Value.RANStatusTransferTransparentContainer
			Ngaplog.Trace("[NGAP] Decode IE RANStatusTransferTransparentContainer")
			if rANStatusTransferTransparentContainer == nil {
				Ngaplog.Error("RANStatusTransferTransparentContainer is nil")

			}
		}
	}

	ranUe = ran.RanUeFindByRanUeNgapID(rANUENGAPID.Value)
	if ranUe == nil {
		Ngaplog.Errorf("No UE Context[RanUeNgapID: %d]", rANUENGAPID.Value)
		return
	}

	Ngaplog.Tracef("UE Context AmfUeNgapID[%d] RanUeNgapID[%d]", ranUe.AmfUeNgapId, ranUe.RanUeNgapId)

	amfUe := ranUe.AmfUe
	if amfUe == nil {
		Ngaplog.Error("AmfUe is nil")
		return
	}
	// send to T-AMF using N1N2MessageTransfer (R16)
}

func HandleNasNonDeliveryIndication(ran *amf_context.AmfRan, message *ngapType.NGAPPDU) {

	var aMFUENGAPID *ngapType.AMFUENGAPID
	var rANUENGAPID *ngapType.RANUENGAPID
	var nASPDU *ngapType.NASPDU
	var cause *ngapType.Cause

	if ran == nil {
		Ngaplog.Error("ran is nil")
		return
	}
	if message == nil {
		Ngaplog.Error("NGAP Message is nil")
		return
	}
	initiatingMessage := message.InitiatingMessage
	if initiatingMessage == nil {
		Ngaplog.Error("InitiatingMessage is nil")
		return
	}
	nASNonDeliveryIndication := initiatingMessage.Value.NASNonDeliveryIndication
	if nASNonDeliveryIndication == nil {
		Ngaplog.Error("NASNonDeliveryIndication is nil")
		return
	}

	Ngaplog.Info("[AMF] Nas Non Delivery Indication")

	for _, ie := range nASNonDeliveryIndication.ProtocolIEs.List {
		switch ie.Id.Value {
		case ngapType.ProtocolIEIDAMFUENGAPID:
			aMFUENGAPID = ie.Value.AMFUENGAPID
			if aMFUENGAPID == nil {
				Ngaplog.Error("AmfUeNgapID is nil")
				return
			}
		case ngapType.ProtocolIEIDRANUENGAPID:
			rANUENGAPID = ie.Value.RANUENGAPID
			if rANUENGAPID == nil {
				Ngaplog.Error("RanUeNgapID is nil")
				return
			}
		case ngapType.ProtocolIEIDNASPDU:
			nASPDU = ie.Value.NASPDU
			if nASPDU == nil {
				Ngaplog.Error("NasPdu is nil")
				return
			}
		case ngapType.ProtocolIEIDCause:
			cause = ie.Value.Cause
			if cause == nil {
				Ngaplog.Error("Cause is nil")
				return
			}
		}
	}

	printRanInfo(ran)

	ranUe := ran.RanUeFindByRanUeNgapID(rANUENGAPID.Value)
	if ranUe == nil {
		Ngaplog.Errorf("No UE Context[RanUeNgapID: %d]", rANUENGAPID.Value)
		return
	}

	Ngaplog.Tracef("RanUeNgapID[%d] AmfUeNgapID[%d]", ranUe.RanUeNgapId, ranUe.AmfUeNgapId)

	printAndGetCause(cause)

	amf_nas.HandleNAS(ranUe, ngapType.ProcedureCodeNASNonDeliveryIndication, nASPDU.Value)
}

func HandleRanConfigurationUpdate(ran *amf_context.AmfRan, message *ngapType.NGAPPDU) {

	var rANNodeName *ngapType.RANNodeName
	var supportedTAList *ngapType.SupportedTAList
	var pagingDRX *ngapType.PagingDRX

	var cause ngapType.Cause

	logger.SetReportCaller(false)

	if ran == nil {
		logger.NgapLog.Error("ran is nil")
		return
	}

	if message == nil {
		logger.NgapLog.Error("NGAP Message is nil")
		return
	}

	initiatingMessage := message.InitiatingMessage
	if initiatingMessage == nil {
		logger.NgapLog.Error("Initiating Message is nil")
		return
	}
	rANConfigurationUpdate := initiatingMessage.Value.RANConfigurationUpdate
	if rANConfigurationUpdate == nil {
		logger.NgapLog.Error("RAN Configuration is nil")
		return
	}
	logger.NgapLog.Info("[AMF] Ran Configuration Update")
	for i := 0; i < len(rANConfigurationUpdate.ProtocolIEs.List); i++ {
		ie := rANConfigurationUpdate.ProtocolIEs.List[i]
		switch ie.Id.Value {
		case ngapType.ProtocolIEIDRANNodeName:
			rANNodeName = ie.Value.RANNodeName
			logger.NgapLog.Tracef("[NGAP] Decode IE RANNodeName = [%s]", rANNodeName.Value)
			if rANNodeName == nil {
				logger.NgapLog.Error("RAN Node Name is nil")
				return
			}
		case ngapType.ProtocolIEIDSupportedTAList:
			supportedTAList = ie.Value.SupportedTAList
			logger.NgapLog.Trace("[NGAP] Decode IE SupportedTAList")
			if supportedTAList == nil {
				logger.NgapLog.Error("Supported TA List is nil")
				return
			}
		case ngapType.ProtocolIEIDDefaultPagingDRX:
			pagingDRX = ie.Value.DefaultPagingDRX
			logger.NgapLog.Tracef("[NGAP] Decode IE PagingDRX = [%d]", pagingDRX.Value)
			if pagingDRX == nil {
				logger.NgapLog.Error("PagingDRX is nil")
				return
			}
		}
	}

	for i := 0; i < len(supportedTAList.List); i++ {
		supportedTAItem := supportedTAList.List[i]
		tac := hex.EncodeToString(supportedTAItem.TAC.Value)
		capOfSupportTai := cap(ran.SupportedTAList)
		for j := 0; j < len(supportedTAItem.BroadcastPLMNList.List); j++ {
			supportedTAI := amf_context.NewSupportedTAI()
			supportedTAI.Tai.Tac = tac
			broadcastPLMNItem := supportedTAItem.BroadcastPLMNList.List[j]
			plmnId := ngapConvert.PlmnIdToModels(broadcastPLMNItem.PLMNIdentity)
			supportedTAI.Tai.PlmnId = &plmnId
			capOfSNssaiList := cap(supportedTAI.SNssaiList)
			for k := 0; k < len(broadcastPLMNItem.TAISliceSupportList.List); k++ {
				tAISliceSupportItem := broadcastPLMNItem.TAISliceSupportList.List[k]
				if len(supportedTAI.SNssaiList) < capOfSNssaiList {
					supportedTAI.SNssaiList = append(supportedTAI.SNssaiList, ngapConvert.SNssaiToModels(tAISliceSupportItem.SNSSAI))
				} else {
					break
				}
			}
			logger.NgapLog.Tracef("PLMN_ID[MCC:%s MNC:%s] TAC[%s]", plmnId.Mcc, plmnId.Mnc, tac)
			if len(ran.SupportedTAList) < capOfSupportTai {
				ran.SupportedTAList = append(ran.SupportedTAList, supportedTAI)

			} else {
				break
			}
		}

	}

	if len(ran.SupportedTAList) == 0 {
		logger.NgapLog.Warn("RanConfigurationUpdate failure: No supported TA exist in RanConfigurationUpdate")
		cause.Present = ngapType.CausePresentMisc
		cause.Misc = &ngapType.CauseMisc{
			Value: ngapType.CauseMiscPresentUnspecified,
		}
	} else {
		var found bool
		for i, tai := range ran.SupportedTAList {
			if amf_context.InTaiList(tai.Tai, amf_context.AMF_Self().SupportTaiLists) {
				logger.NgapLog.Tracef("SERVED_TAI_INDEX[%d]", i)
				found = true
				break
			}
		}
		if !found {
			logger.NgapLog.Warn("RanConfigurationUpdate failure: Cannot find Served TAI in AMF")
			cause.Present = ngapType.CausePresentMisc
			cause.Misc = &ngapType.CauseMisc{
				Value: ngapType.CauseMiscPresentUnknownPLMN,
			}
		}
	}

	if cause.Present == ngapType.CausePresentNothing {
		logger.NgapLog.Info("[AMF] RanConfigurationUpdateAcknowledge")
		ngap_message.SendRanConfigurationUpdateAcknowledge(ran, nil)
	} else {
		logger.NgapLog.Info("[AMF] RanConfigurationUpdateAcknowledgeFailure")
		ngap_message.SendRanConfigurationUpdateFailure(ran, cause, nil)
	}
}

func HandleAMFStatusIndication(ran *amf_context.AmfRan, message *ngapType.NGAPPDU) {
}

func HandleUplinkRanConfigurationTransfer(ran *amf_context.AmfRan, message *ngapType.NGAPPDU) {

	var sONConfigurationTransferUL *ngapType.SONConfigurationTransfer

	if ran == nil {
		Ngaplog.Error("ran is nil")
		return
	}
	if message == nil {
		Ngaplog.Error("NGAP Message is nil")
		return
	}
	initiatingMessage := message.InitiatingMessage
	if initiatingMessage == nil {
		Ngaplog.Error("InitiatingMessage is nil")
		return
	}
	uplinkRANConfigurationTransfer := initiatingMessage.Value.UplinkRANConfigurationTransfer
	if uplinkRANConfigurationTransfer == nil {
		Ngaplog.Error("ErrorIndication is nil")
		return
	}

	for _, ie := range uplinkRANConfigurationTransfer.ProtocolIEs.List {
		switch ie.Id.Value {
		case ngapType.ProtocolIEIDSONConfigurationTransferUL: // optional, ignore
			sONConfigurationTransferUL = ie.Value.SONConfigurationTransferUL
			Ngaplog.Trace("[NGAP] Decode IE SONConfigurationTransferUL")
			if sONConfigurationTransferUL == nil {
				Ngaplog.Warn("sONConfigurationTransferUL is nil")
			}
		}
	}

	if sONConfigurationTransferUL != nil {

		targetRanNodeID := ngapConvert.RanIdToModels(sONConfigurationTransferUL.TargetRANNodeID.GlobalRANNodeID)

		if targetRanNodeID.GNbId.GNBValue != "" {
			logger.NgapLog.Tracef("targerRanID [%s]", targetRanNodeID.GNbId.GNBValue)
		}

		aMFSelf := amf_context.AMF_Self()

		targetRan := aMFSelf.AmfRanFindByRanId(targetRanNodeID)
		if targetRan == nil {
			logger.NgapLog.Warn("targetRan is nil")
		}

		ngap_message.SendDownlinkRanConfigurationTransfer(targetRan, sONConfigurationTransferUL)
	}
}

func HandleUplinkUEAssociatedNRPPATransport(ran *amf_context.AmfRan, message *ngapType.NGAPPDU) {

	var aMFUENGAPID *ngapType.AMFUENGAPID
	var rANUENGAPID *ngapType.RANUENGAPID
	var routingID *ngapType.RoutingID
	var nRPPaPDU *ngapType.NRPPaPDU

	if ran == nil {
		Ngaplog.Error("ran is nil")
		return
	}
	if message == nil {
		Ngaplog.Error("NGAP Message is nil")
		return
	}
	initiatingMessage := message.InitiatingMessage
	if initiatingMessage == nil {
		Ngaplog.Error("InitiatingMessage is nil")
		return
	}
	uplinkUEAssociatedNRPPaTransport := initiatingMessage.Value.UplinkUEAssociatedNRPPaTransport
	if uplinkUEAssociatedNRPPaTransport == nil {
		Ngaplog.Error("uplinkUEAssociatedNRPPaTransport is nil")
		return
	}

	Ngaplog.Info("[AMF] Uplink UE Associated NRPPA Transpor")

	for _, ie := range uplinkUEAssociatedNRPPaTransport.ProtocolIEs.List {
		switch ie.Id.Value {
		case ngapType.ProtocolIEIDAMFUENGAPID: // reject
			aMFUENGAPID = ie.Value.AMFUENGAPID
			Ngaplog.Trace("[NGAP] Decode IE aMFUENGAPID")
			if aMFUENGAPID == nil {
				Ngaplog.Error("AmfUeNgapID is nil")
				return
			}
		case ngapType.ProtocolIEIDRANUENGAPID: // reject
			rANUENGAPID = ie.Value.RANUENGAPID
			Ngaplog.Trace("[NGAP] Decode IE rANUENGAPID")
			if rANUENGAPID == nil {
				Ngaplog.Error("RanUeNgapID is nil")
				return
			}
		case ngapType.ProtocolIEIDRoutingID: // reject
			routingID = ie.Value.RoutingID
			Ngaplog.Trace("[NGAP] Decode IE routingID")
			if routingID == nil {
				Ngaplog.Error("routingID is nil")
				return
			}
		case ngapType.ProtocolIEIDNRPPaPDU: // reject
			nRPPaPDU = ie.Value.NRPPaPDU
			Ngaplog.Trace("[NGAP] Decode IE nRPPaPDU")
			if nRPPaPDU == nil {
				Ngaplog.Error("nRPPaPDU is nil")
				return
			}
		}
	}

	printRanInfo(ran)

	ranUe := ran.RanUeFindByRanUeNgapID(rANUENGAPID.Value)
	if ranUe == nil {
		Ngaplog.Errorf("No UE Context[RanUeNgapID: %d]", rANUENGAPID.Value)
		return
	}

	Ngaplog.Tracef("RanUeNgapId[%d] AmfUeNgapId[%d]", ranUe.RanUeNgapId, ranUe.AmfUeNgapId)

	ranUe.RoutingID = hex.EncodeToString(routingID.Value)

	// TODO: Forward NRPPaPDU to LMF
}

func HandleUplinkNonUEAssociatedNRPPATransport(ran *amf_context.AmfRan, message *ngapType.NGAPPDU) {
	var routingID *ngapType.RoutingID
	var nRPPaPDU *ngapType.NRPPaPDU

	if ran == nil {
		logger.NgapLog.Error("ran is nil")
		return
	}
	if message == nil {
		logger.NgapLog.Error("NGAP Message is nil")
		return
	}
	initiatingMessage := message.InitiatingMessage
	if initiatingMessage == nil {
		logger.NgapLog.Error("Initiating Message is nil")
		return
	}
	uplinkNonUEAssociatedNRPPATransport := initiatingMessage.Value.UplinkNonUEAssociatedNRPPaTransport
	if uplinkNonUEAssociatedNRPPATransport == nil {
		logger.NgapLog.Error("Uplink Non UE Associated NRPPA Transport is nil")
		return
	}

	logger.NgapLog.Info("[AMF] Uplink Non UE Associated NRPPA Transport")

	for i := 0; i < len(uplinkNonUEAssociatedNRPPATransport.ProtocolIEs.List); i++ {
		ie := uplinkNonUEAssociatedNRPPATransport.ProtocolIEs.List[i]
		switch ie.Id.Value {

		case ngapType.ProtocolIEIDRoutingID:
			routingID = ie.Value.RoutingID
			logger.NgapLog.Trace("[NGAP] Decode IE RoutingID")

		case ngapType.ProtocolIEIDNRPPaPDU:
			nRPPaPDU = ie.Value.NRPPaPDU
			logger.NgapLog.Trace("[NGAP] Decode IE NRPPaPDU")
		}
	}

	if routingID == nil {
		logger.NgapLog.Error("RoutingID is nil")
		return
	}
	// Forward routingID to LMF
	// Described in (23.502 4.13.5.6)

	if nRPPaPDU == nil {
		logger.NgapLog.Error("NRPPaPDU is nil")
		return
	}
	// TODO: Forward NRPPaPDU to LMF
}

func HandleLocationReport(ran *amf_context.AmfRan, message *ngapType.NGAPPDU) {

	var aMFUENGAPID *ngapType.AMFUENGAPID
	var rANUENGAPID *ngapType.RANUENGAPID
	var userLocationInformation *ngapType.UserLocationInformation
	var uEPresenceInAreaOfInterestList *ngapType.UEPresenceInAreaOfInterestList
	var locationReportingRequestType *ngapType.LocationReportingRequestType

	if ran == nil {
		Ngaplog.Error("ran is nil")
		return
	}
	if message == nil {
		Ngaplog.Error("NGAP Message is nil")
		return
	}
	initiatingMessage := message.InitiatingMessage
	if initiatingMessage == nil {
		Ngaplog.Error("InitiatingMessage is nil")
		return
	}
	locationReport := initiatingMessage.Value.LocationReport
	if locationReport == nil {
		Ngaplog.Error("LocationReport is nil")
		return
	}

	logger.NgapLog.Info("[AMF] Location Report")
	for _, ie := range locationReport.ProtocolIEs.List {
		switch ie.Id.Value {
		case ngapType.ProtocolIEIDAMFUENGAPID: // reject
			aMFUENGAPID = ie.Value.AMFUENGAPID
			Ngaplog.Trace("[NGAP] Decode IE AmfUeNgapID")
			if aMFUENGAPID == nil {
				Ngaplog.Error("AmfUeNgapID is nil")
			}
		case ngapType.ProtocolIEIDRANUENGAPID: // reject
			rANUENGAPID = ie.Value.RANUENGAPID
			Ngaplog.Trace("[NGAP] Decode IE RanUeNgapID")
			if rANUENGAPID == nil {
				Ngaplog.Error("RanUeNgapID is nil")
			}
		case ngapType.ProtocolIEIDUserLocationInformation: // ignore
			userLocationInformation = ie.Value.UserLocationInformation
			Ngaplog.Trace("[NGAP] Decode IE userLocationInformation")
			if userLocationInformation == nil {
				Ngaplog.Warn("userLocationInformation is nil")
			}
		case ngapType.ProtocolIEIDUEPresenceInAreaOfInterestList: // optional, ignore
			uEPresenceInAreaOfInterestList = ie.Value.UEPresenceInAreaOfInterestList
			Ngaplog.Trace("[NGAP] Decode IE uEPresenceInAreaOfInterestList")
			if uEPresenceInAreaOfInterestList == nil {
				Ngaplog.Warn("uEPresenceInAreaOfInterestList is nil [optional]")
			}
		case ngapType.ProtocolIEIDLocationReportingRequestType: // ignore
			locationReportingRequestType = ie.Value.LocationReportingRequestType
			Ngaplog.Trace("[NGAP] Decode IE LocationReportingRequestType")
			if locationReportingRequestType == nil {
				Ngaplog.Warn("LocationReportingRequestType is nil")
			}
		}
	}

	printRanInfo(ran)

	ranUe := ran.RanUeFindByRanUeNgapID(rANUENGAPID.Value)
	if ranUe == nil {
		Ngaplog.Errorf("No UE Context[RanUeNgapID: %d]", rANUENGAPID.Value)
		return
	}

	ranUe.UpdateLocation(userLocationInformation)

	Ngaplog.Tracef("Report Area[%d]", locationReportingRequestType.ReportArea.Value)

	switch locationReportingRequestType.EventType.Value {
	case ngapType.EventTypePresentDirect:
		Ngaplog.Trace("To report directly")

	case ngapType.EventTypePresentChangeOfServeCell:
		Ngaplog.Trace("To report upon change of serving cell")

	case ngapType.EventTypePresentUePresenceInAreaOfInterest:
		Ngaplog.Trace("To report UE presence in the area of interest")
		for _, uEPresenceInAreaOfInterestItem := range uEPresenceInAreaOfInterestList.List {
			uEPresence := uEPresenceInAreaOfInterestItem.UEPresence.Value
			referenceID := uEPresenceInAreaOfInterestItem.LocationReportingReferenceID.Value

			for _, AOIitem := range locationReportingRequestType.AreaOfInterestList.List {
				if referenceID == AOIitem.LocationReportingReferenceID.Value {
					Ngaplog.Tracef("uEPresence[%d], presence AOI ReferenceID[%d]", uEPresence, referenceID)
				}
			}
		}

	case ngapType.EventTypePresentStopChangeOfServeCell:
		Ngaplog.Trace("To stop reporting at change of serving cell")
		ngap_message.SendLocationReportingControl(ranUe, nil, 0, locationReportingRequestType.EventType)
		// TODO: Clear location report

	case ngapType.EventTypePresentStopUePresenceInAreaOfInterest:
		Ngaplog.Trace("To stop reporting UE presence in the area of interest")
		Ngaplog.Tracef("ReferenceID To Be Cancelled[%d]", locationReportingRequestType.LocationReportingReferenceIDToBeCancelled.Value)
		// TODO: Clear location report

	case ngapType.EventTypePresentCancelLocationReportingForTheUe:
		Ngaplog.Trace("To cancel location reporting for the UE")
		// TODO: Clear location report
	}
}

func HandleUETNLABindingReleaseRequest(ran *amf_context.AmfRan, message *ngapType.NGAPPDU) {
}

func HandleUERadioCapabilityInfoIndication(ran *amf_context.AmfRan, message *ngapType.NGAPPDU) {

	var aMFUENGAPID *ngapType.AMFUENGAPID
	var rANUENGAPID *ngapType.RANUENGAPID

	var uERadioCapability *ngapType.UERadioCapability
	var uERadioCapabilityForPaging *ngapType.UERadioCapabilityForPaging

	if ran == nil {
		logger.NgapLog.Error("ran is nil")
		return
	}
	if message == nil {
		logger.NgapLog.Error("NGAP Message is nil")
		return
	}
	initiatingMessage := message.InitiatingMessage
	if initiatingMessage == nil {
		logger.NgapLog.Error("Initiating Message is nil")
		return
	}
	uERadioCapabilityInfoIndication := initiatingMessage.Value.UERadioCapabilityInfoIndication
	if uERadioCapabilityInfoIndication == nil {
		logger.NgapLog.Error("UERadioCapabilityInfoIndication is nil")
		return
	}

	logger.NgapLog.Info("[AMF] UE Radio Capability Info Indication")

	for i := 0; i < len(uERadioCapabilityInfoIndication.ProtocolIEs.List); i++ {
		ie := uERadioCapabilityInfoIndication.ProtocolIEs.List[i]
		switch ie.Id.Value {
		case ngapType.ProtocolIEIDAMFUENGAPID:
			aMFUENGAPID = ie.Value.AMFUENGAPID
			logger.NgapLog.Trace("[NGAP] Decode IE AmfUeNgapID")
			if aMFUENGAPID == nil {
				logger.NgapLog.Error("AmfUeNgapID is nil")
				return
			}
		case ngapType.ProtocolIEIDRANUENGAPID:
			rANUENGAPID = ie.Value.RANUENGAPID
			logger.NgapLog.Trace("[NGAP] Decode IE RanUeNgapID")
			if rANUENGAPID == nil {
				logger.NgapLog.Error("RanUeNgapID is nil")
				return
			}
		case ngapType.ProtocolIEIDUERadioCapability:
			uERadioCapability = ie.Value.UERadioCapability
			logger.NgapLog.Trace("[NGAP] Decode IE UERadioCapability")
			if uERadioCapability == nil {
				logger.NgapLog.Error("UERadioCapability is nil")
				return
			}
		case ngapType.ProtocolIEIDUERadioCapabilityForPaging:
			uERadioCapabilityForPaging = ie.Value.UERadioCapabilityForPaging
			logger.NgapLog.Trace("[NGAP] Decode IE UERadioCapabilityForPaging")
			if uERadioCapabilityForPaging == nil {
				logger.NgapLog.Error("UERadioCapabilityForPaging is nil")
				return
			}
		}
	}

	printRanInfo(ran)

	ranUe := ran.RanUeFindByRanUeNgapID(rANUENGAPID.Value)
	if ranUe == nil {
		Ngaplog.Errorf("No UE Context[RanUeNgapID: %d]", rANUENGAPID.Value)
		return
	}
	Ngaplog.Tracef("RanUeNgapID[%d] AmfUeNgapID[%d]", ranUe.RanUeNgapId, ranUe.AmfUeNgapId)
	amfUe := ranUe.AmfUe

	if amfUe == nil {
		Ngaplog.Errorln("amfUe is nil")
		return
	}
	if uERadioCapability != nil {
		amfUe.UeRadioCapability = hex.EncodeToString(uERadioCapability.Value)
	}
	if uERadioCapabilityForPaging != nil {
		amfUe.UeRadioCapabilityForPaging = &amf_context.UERadioCapabilityForPaging{}
		if uERadioCapabilityForPaging.UERadioCapabilityForPagingOfNR != nil {
			amfUe.UeRadioCapabilityForPaging.NR = hex.EncodeToString(uERadioCapabilityForPaging.UERadioCapabilityForPagingOfNR.Value)
		}
		if uERadioCapabilityForPaging.UERadioCapabilityForPagingOfEUTRA != nil {
			amfUe.UeRadioCapabilityForPaging.EUTRA = hex.EncodeToString(uERadioCapabilityForPaging.UERadioCapabilityForPagingOfEUTRA.Value)
		}
	}

	// TS 38.413 8.14.1.2/TS 23.502 4.2.8a step5/TS 23.501, clause 5.4.4.1.
	//send its most up to date UE Radio Capability information to the RAN in the N2 REQUEST message.
}

func HandleAMFconfigurationUpdateFailure(ran *amf_context.AmfRan, message *ngapType.NGAPPDU) {

	var cause *ngapType.Cause
	var criticalityDiagnostics *ngapType.CriticalityDiagnostics
	if ran == nil {
		logger.NgapLog.Error("ran is nil")
		return
	}
	if message == nil {
		logger.NgapLog.Error("NGAP Message is nil")
		return
	}
	unsuccessfulOutcome := message.UnsuccessfulOutcome
	if unsuccessfulOutcome == nil {
		logger.NgapLog.Error("Unsuccessful Message is nil")
		return
	}

	AMFconfigurationUpdateFailure := unsuccessfulOutcome.Value.AMFConfigurationUpdateFailure
	if AMFconfigurationUpdateFailure == nil {
		logger.NgapLog.Error("AMFConfigurationUpdateFailure is nil")
		return
	}

	logger.NgapLog.Info(3, " [AMF] AMF Confioguration Update Failure")

	for _, ie := range AMFconfigurationUpdateFailure.ProtocolIEs.List {
		switch ie.Id.Value {
		case ngapType.ProtocolIEIDCause:
			cause = ie.Value.Cause
			logger.NgapLog.Trace("[NGAP] Decode IE Cause")
			if cause == nil {
				logger.NgapLog.Error("Cause is nil")
				return
			}
		case ngapType.ProtocolIEIDCriticalityDiagnostics:
			criticalityDiagnostics = ie.Value.CriticalityDiagnostics
			Ngaplog.Trace("[NGAP] Decode IE CriticalityDiagnostics")
		}
	}

	//	TODO: Time To Wait

	if criticalityDiagnostics != nil {
		printCriticalityDiagnostics(criticalityDiagnostics)
	}
}

func HandleAMFconfigurationUpdateAcknowledge(ran *amf_context.AmfRan, message *ngapType.NGAPPDU) {

	var aMFTNLAssociationSetupList *ngapType.AMFTNLAssociationSetupList
	var criticalityDiagnostics *ngapType.CriticalityDiagnostics
	var aMFTNLAssociationFailedToSetupList *ngapType.TNLAssociationList
	if ran == nil {
		logger.NgapLog.Error("ran is nil")
		return
	}
	if message == nil {
		logger.NgapLog.Error("NGAP Message is nil")
		return
	}
	successfulOutcome := message.SuccessfulOutcome
	if successfulOutcome == nil {
		logger.NgapLog.Error("SuccessfulOutcome is nil")
		return
	}
	aMFConfigurationUpdateAcknowledge := successfulOutcome.Value.AMFConfigurationUpdateAcknowledge
	if aMFConfigurationUpdateAcknowledge == nil {
		logger.NgapLog.Error("AMFConfigurationUpdateAcknowledge is nil")
		return
	}

	logger.NgapLog.Info("[AMF] AMF Configuration Update Acknowledge")

	for i := 0; i < len(aMFConfigurationUpdateAcknowledge.ProtocolIEs.List); i++ {
		ie := aMFConfigurationUpdateAcknowledge.ProtocolIEs.List[i]
		switch ie.Id.Value {
		case ngapType.ProtocolIEIDAMFTNLAssociationSetupList:
			aMFTNLAssociationSetupList = ie.Value.AMFTNLAssociationSetupList
			logger.NgapLog.Trace("[NGAP] Decode IE AMFTNLAssociationSetupList")
			if aMFTNLAssociationSetupList == nil {
				logger.NgapLog.Error("AMFTNLAssociationSetupList is nil")
				return
			}
		case ngapType.ProtocolIEIDCriticalityDiagnostics:
			criticalityDiagnostics = ie.Value.CriticalityDiagnostics
			logger.NgapLog.Trace("[NGAP] Decode IE Criticality Diagnostics")

		case ngapType.ProtocolIEIDAMFTNLAssociationFailedToSetupList:
			aMFTNLAssociationFailedToSetupList = ie.Value.AMFTNLAssociationFailedToSetupList
			logger.NgapLog.Trace("[NGAP] Decode IE AMFTNLAssociationFailedToSetupList")
			if aMFTNLAssociationFailedToSetupList == nil {
				logger.NgapLog.Error("AMFTNLAssociationFailedToSetupList is nil")
				return
			}
		}
	}

	if criticalityDiagnostics != nil {
		printCriticalityDiagnostics(criticalityDiagnostics)
	}
}

func HandleErrorIndication(ran *amf_context.AmfRan, message *ngapType.NGAPPDU) {

	var aMFUENGAPID *ngapType.AMFUENGAPID
	var rANUENGAPID *ngapType.RANUENGAPID
	var cause *ngapType.Cause
	var criticalityDiagnostics *ngapType.CriticalityDiagnostics

	if ran == nil {
		Ngaplog.Error("ran is nil")
		return
	}
	if message == nil {
		Ngaplog.Error("NGAP Message is nil")
		return
	}
	initiatingMessage := message.InitiatingMessage
	if initiatingMessage == nil {
		Ngaplog.Error("InitiatingMessage is nil")
		return
	}
	errorIndication := initiatingMessage.Value.ErrorIndication
	if errorIndication == nil {
		Ngaplog.Error("ErrorIndication is nil")
		return
	}

	for _, ie := range errorIndication.ProtocolIEs.List {
		switch ie.Id.Value {
		case ngapType.ProtocolIEIDAMFUENGAPID:
			aMFUENGAPID = ie.Value.AMFUENGAPID
			Ngaplog.Trace("[NGAP] Decode IE AmfUeNgapID")
			if aMFUENGAPID == nil {
				Ngaplog.Error("AmfUeNgapID is nil")
			}
		case ngapType.ProtocolIEIDRANUENGAPID:
			rANUENGAPID = ie.Value.RANUENGAPID
			Ngaplog.Trace("[NGAP] Decode IE RanUeNgapID")
			if rANUENGAPID == nil {
				Ngaplog.Error("RanUeNgapID is nil")
			}
		case ngapType.ProtocolIEIDCause:
			cause = ie.Value.Cause
			Ngaplog.Trace("[NGAP] Decode IE Cause")
		case ngapType.ProtocolIEIDCriticalityDiagnostics:
			criticalityDiagnostics = ie.Value.CriticalityDiagnostics
			Ngaplog.Trace("[NGAP] Decode IE CriticalityDiagnostics")
		}
	}

	printRanInfo(ran)

	if cause == nil && criticalityDiagnostics == nil {
		Ngaplog.Error("[ErrorIndication] both Cause IE and CriticalityDiagnostics IE are nil, should have at least one")
		return
	}

	if cause != nil {
		printAndGetCause(cause)
	}

	if criticalityDiagnostics != nil {
		if criticalityDiagnostics.IEsCriticalityDiagnostics != nil {
			for _, ieCriticalityDiagnostics := range criticalityDiagnostics.IEsCriticalityDiagnostics.List {
				Ngaplog.Tracef("IE ID: %d", ieCriticalityDiagnostics.IEID.Value)

				switch ieCriticalityDiagnostics.IECriticality.Value {
				case ngapType.CriticalityPresentReject:
					Ngaplog.Trace("Criticality Reject")
				case ngapType.CriticalityPresentNotify:
					Ngaplog.Trace("Criticality Notify")
				}

				switch ieCriticalityDiagnostics.TypeOfError.Value {
				case ngapType.TypeOfErrorPresentNotUnderstood:
					Ngaplog.Trace("Type of error: Not understood")
				case ngapType.TypeOfErrorPresentMissing:
					Ngaplog.Trace("Type of error: Missing")
				}
			}
		}
	}

	// TODO: handle error based on cause/criticalityDiagnostics
}

func HandleCellTrafficTrace(ran *amf_context.AmfRan, message *ngapType.NGAPPDU) {

	var aMFUENGAPID *ngapType.AMFUENGAPID
	var rANUENGAPID *ngapType.RANUENGAPID
	var nGRANTraceID *ngapType.NGRANTraceID
	var nGRANCGI *ngapType.NGRANCGI
	var traceCollectionEntityIPAddress *ngapType.TransportLayerAddress

	var ranUe *amf_context.RanUe

	var iesCriticalityDiagnostics ngapType.CriticalityDiagnosticsIEList

	if ran == nil {
		Ngaplog.Error("ran is nil")
		return
	}
	if message == nil {
		Ngaplog.Error("NGAP Message is nil")
		return
	}
	initiatingMessage := message.InitiatingMessage // ignore
	if initiatingMessage == nil {
		Ngaplog.Error("InitiatingMessage is nil")
		return
	}
	cellTrafficTrace := initiatingMessage.Value.CellTrafficTrace
	if cellTrafficTrace == nil {
		Ngaplog.Error("CellTrafficTrace is nil")
		return
	}

	Ngaplog.Info("[AMF] Cell Traffic Trace")

	for _, ie := range cellTrafficTrace.ProtocolIEs.List {
		switch ie.Id.Value {
		case ngapType.ProtocolIEIDAMFUENGAPID: // reject
			aMFUENGAPID = ie.Value.AMFUENGAPID
			Ngaplog.Trace("[NGAP] Decode IE AmfUeNgapID")
		case ngapType.ProtocolIEIDRANUENGAPID: // reject
			rANUENGAPID = ie.Value.RANUENGAPID
			Ngaplog.Trace("[NGAP] Decode IE RanUeNgapID")

		case ngapType.ProtocolIEIDNGRANTraceID: // ignore
			nGRANTraceID = ie.Value.NGRANTraceID
			Ngaplog.Trace("[NGAP] Decode IE NGRANTraceID")
		case ngapType.ProtocolIEIDNGRANCGI: // ignore
			nGRANCGI = ie.Value.NGRANCGI
			Ngaplog.Trace("[NGAP] Decode IE NGRANCGI")
		case ngapType.ProtocolIEIDTraceCollectionEntityIPAddress: // ignore
			traceCollectionEntityIPAddress = ie.Value.TraceCollectionEntityIPAddress
			Ngaplog.Trace("[NGAP] Decode IE TraceCollectionEntityIPAddress")
		}
	}
	if aMFUENGAPID == nil {
		Ngaplog.Error("AmfUeNgapID is nil")
		item := buildCriticalityDiagnosticsIEItem(ngapType.CriticalityPresentReject, ngapType.ProtocolIEIDAMFUENGAPID, ngapType.TypeOfErrorPresentMissing)
		iesCriticalityDiagnostics.List = append(iesCriticalityDiagnostics.List, item)
	}
	if rANUENGAPID == nil {
		Ngaplog.Error("RanUeNgapID is nil")
		item := buildCriticalityDiagnosticsIEItem(ngapType.CriticalityPresentReject, ngapType.ProtocolIEIDRANUENGAPID, ngapType.TypeOfErrorPresentMissing)
		iesCriticalityDiagnostics.List = append(iesCriticalityDiagnostics.List, item)
	}

	if len(iesCriticalityDiagnostics.List) > 0 {
		procedureCode := ngapType.ProcedureCodeCellTrafficTrace
		triggeringMessage := ngapType.TriggeringMessagePresentInitiatingMessage
		procedureCriticality := ngapType.CriticalityPresentIgnore
		criticalityDiagnostics := buildCriticalityDiagnostics(&procedureCode, &triggeringMessage, &procedureCriticality, &iesCriticalityDiagnostics)
		ngap_message.SendErrorIndication(ran, nil, nil, nil, &criticalityDiagnostics)
		return
	}

	if aMFUENGAPID != nil {
		ranUe = amf_context.AMF_Self().RanUeFindByAmfUeNgapID(aMFUENGAPID.Value)
		if ranUe == nil {
			Ngaplog.Errorf("No UE Context[AmfUeNgapID: %d]", aMFUENGAPID.Value)
			cause := ngapType.Cause{
				Present: ngapType.CausePresentRadioNetwork,
				RadioNetwork: &ngapType.CauseRadioNetwork{
					Value: ngapType.CauseRadioNetworkPresentUnknownLocalUENGAPID,
				},
			}
			ngap_message.SendErrorIndication(ran, nil, nil, &cause, nil)
			return
		}
	}

	Ngaplog.Debugf("UE: AmfUeNgapID[%d], RanUeNgapID[%d]", ranUe.AmfUeNgapId, ranUe.RanUeNgapId)

	ranUe.Trsr = hex.EncodeToString(nGRANTraceID.Value[6:])

	Ngaplog.Tracef("TRSR[%s]", ranUe.Trsr)

	switch nGRANCGI.Present {
	case ngapType.NGRANCGIPresentNRCGI:
		plmnID := ngapConvert.PlmnIdToModels(nGRANCGI.NRCGI.PLMNIdentity)
		cellID := ngapConvert.BitStringToHex(&nGRANCGI.NRCGI.NRCellIdentity.Value)
		Ngaplog.Debugf("NRCGI[plmn: %s, cellID: %s]", plmnID, cellID)
	case ngapType.NGRANCGIPresentEUTRACGI:
		plmnID := ngapConvert.PlmnIdToModels(nGRANCGI.EUTRACGI.PLMNIdentity)
		cellID := ngapConvert.BitStringToHex(&nGRANCGI.EUTRACGI.EUTRACellIdentity.Value)
		Ngaplog.Debugf("EUTRACGI[plmn: %s, cellID: %s]", plmnID, cellID)

	}

	tceIpv4, tceIpv6 := ngapConvert.IPAddressToString(*traceCollectionEntityIPAddress)
	if tceIpv4 != "" {
		Ngaplog.Debugf("TCE IP Address[v4: %s]", tceIpv4)
	}
	if tceIpv6 != "" {
		Ngaplog.Debugf("TCE IP Address[v6: %s]", tceIpv6)
	}

	// TODO: TS 32.422 4.2.2.10
	// When AMF receives this new NG signalling message containing the Trace Recording Session Reference (TRSR)
	// and Trace Reference (TR), the AMF shall look up the SUPI/IMEI(SV) of the given call from its database and
	// shall send the SUPI/IMEI(SV) numbers together with the Trace Recording Session Reference and Trace Reference
	// to the Trace Collection Entity.
}

func printAndGetCause(cause *ngapType.Cause) (present int, value aper.Enumerated) {

	present = cause.Present
	switch cause.Present {
	case ngapType.CausePresentRadioNetwork:
		Ngaplog.Warnf("Cause RadioNetwork[%d]", cause.RadioNetwork.Value)
		value = cause.RadioNetwork.Value
	case ngapType.CausePresentTransport:
		Ngaplog.Warnf("Cause Transport[%d]", cause.Transport.Value)
		value = cause.Transport.Value
	case ngapType.CausePresentProtocol:
		Ngaplog.Warnf("Cause Protocol[%d]", cause.Protocol.Value)
		value = cause.Protocol.Value
	case ngapType.CausePresentNas:
		Ngaplog.Warnf("Cause Nas[%d]", cause.Nas.Value)
		value = cause.Nas.Value
	case ngapType.CausePresentMisc:
		Ngaplog.Warnf("Cause Misc[%d]", cause.Misc.Value)
		value = cause.Misc.Value
	default:
		Ngaplog.Errorf("Invalid Cause group[%d]", cause.Present)
	}
	return
}

func printCriticalityDiagnostics(criticalityDiagnostics *ngapType.CriticalityDiagnostics) {

	if criticalityDiagnostics.IEsCriticalityDiagnostics != nil {
		for _, ieCriticalityDiagnostics := range criticalityDiagnostics.IEsCriticalityDiagnostics.List {
			Ngaplog.Tracef("IE ID: %d", ieCriticalityDiagnostics.IEID.Value)

			switch ieCriticalityDiagnostics.IECriticality.Value {
			case ngapType.CriticalityPresentReject:
				Ngaplog.Trace("Criticality Reject")
			case ngapType.CriticalityPresentNotify:
				Ngaplog.Trace("Criticality Notify")
			}

			switch ieCriticalityDiagnostics.TypeOfError.Value {
			case ngapType.TypeOfErrorPresentNotUnderstood:
				Ngaplog.Trace("Type of error: Not understood")
			case ngapType.TypeOfErrorPresentMissing:
				Ngaplog.Trace("Type of error: Missing")
			}
		}
	}
}

func buildCriticalityDiagnostics(
	procedureCode *int64,
	triggeringMessage *aper.Enumerated,
	procedureCriticality *aper.Enumerated,
	iesCriticalityDiagnostics *ngapType.CriticalityDiagnosticsIEList) (criticalityDiagnostics ngapType.CriticalityDiagnostics) {

	if procedureCode != nil {
		criticalityDiagnostics.ProcedureCode = new(ngapType.ProcedureCode)
		criticalityDiagnostics.ProcedureCode.Value = *procedureCode
	}

	if triggeringMessage != nil {
		criticalityDiagnostics.TriggeringMessage = new(ngapType.TriggeringMessage)
		criticalityDiagnostics.TriggeringMessage.Value = *triggeringMessage
	}

	if procedureCriticality != nil {
		criticalityDiagnostics.ProcedureCriticality = new(ngapType.Criticality)
		criticalityDiagnostics.ProcedureCriticality.Value = *procedureCriticality
	}

	if iesCriticalityDiagnostics != nil {
		criticalityDiagnostics.IEsCriticalityDiagnostics = iesCriticalityDiagnostics
	}

	return criticalityDiagnostics
}

func buildCriticalityDiagnosticsIEItem(ieCriticality aper.Enumerated, ieID int64, typeOfErr aper.Enumerated) (item ngapType.CriticalityDiagnosticsIEItem) {

	item = ngapType.CriticalityDiagnosticsIEItem{
		IECriticality: ngapType.Criticality{
			Value: ieCriticality,
		},
		IEID: ngapType.ProtocolIEID{
			Value: ieID,
		},
		TypeOfError: ngapType.TypeOfError{
			Value: typeOfErr,
		},
	}

	return item
}

func printRanInfo(ran *amf_context.AmfRan) {
	switch ran.RanPresent {
	case amf_context.RanPresentGNbId:
		Ngaplog.Tracef("IP[%s] GNbId[%s]", ran.Conn.RemoteAddr().String(), ran.RanId.GNbId.GNBValue)
	case amf_context.RanPresentNgeNbId:
		Ngaplog.Tracef("IP[%s] NgeNbId[%s]", ran.Conn.RemoteAddr().String(), ran.RanId.NgeNbId)
	case amf_context.RanPresentN3IwfId:
		Ngaplog.Tracef("IP[%s] N3IwfId[%s]", ran.Conn.RemoteAddr().String(), ran.RanId.N3IwfId)
	}
}
