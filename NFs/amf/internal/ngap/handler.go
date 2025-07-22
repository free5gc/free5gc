package ngap

import (
	"encoding/hex"
	"fmt"
	"strconv"

	"github.com/free5gc/amf/internal/context"
	gmm_common "github.com/free5gc/amf/internal/gmm/common"
	gmm_message "github.com/free5gc/amf/internal/gmm/message"
	amf_nas "github.com/free5gc/amf/internal/nas"
	"github.com/free5gc/amf/internal/nas/nas_security"
	ngap_message "github.com/free5gc/amf/internal/ngap/message"
	"github.com/free5gc/amf/internal/sbi/consumer"
	"github.com/free5gc/amf/pkg/factory"
	"github.com/free5gc/aper"
	"github.com/free5gc/nas"
	"github.com/free5gc/nas/nasMessage"
	libngap "github.com/free5gc/ngap"
	"github.com/free5gc/ngap/ngapConvert"
	"github.com/free5gc/ngap/ngapType"
	"github.com/free5gc/openapi/models"
)

func handleNGSetupRequestMain(ran *context.AmfRan,
	globalRANNodeID *ngapType.GlobalRANNodeID,
	rANNodeName *ngapType.RANNodeName,
	supportedTAList *ngapType.SupportedTAList,
	pagingDRX *ngapType.PagingDRX,
) {
	var cause ngapType.Cause

	ran.SetRanId(globalRANNodeID)
	if rANNodeName != nil {
		ran.Name = rANNodeName.Value
	}
	if pagingDRX != nil {
		ran.Log.Tracef("PagingDRX[%d]", pagingDRX.Value)
	}

	for i := 0; i < len(supportedTAList.List); i++ {
		supportedTAItem := supportedTAList.List[i]
		tac := hex.EncodeToString(supportedTAItem.TAC.Value)
		capOfSupportTai := cap(ran.SupportedTAList)
		for j := 0; j < len(supportedTAItem.BroadcastPLMNList.List); j++ {
			supportedTAI := context.NewSupportedTAI()
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
			ran.Log.Tracef("PLMN_ID[MCC:%s MNC:%s] TAC[%s]", plmnId.Mcc, plmnId.Mnc, tac)
			if len(ran.SupportedTAList) < capOfSupportTai {
				ran.SupportedTAList = append(ran.SupportedTAList, supportedTAI)
			} else {
				break
			}
		}
	}

	if len(ran.SupportedTAList) == 0 {
		ran.Log.Warn("NG-Setup failure: No supported TA exist in NG-Setup request")
		cause.Present = ngapType.CausePresentMisc
		cause.Misc = &ngapType.CauseMisc{
			Value: ngapType.CauseMiscPresentUnspecified,
		}
	} else {
		var found bool
		for i, tai := range ran.SupportedTAList {
			if context.InTaiList(tai.Tai, context.GetSelf().SupportTaiLists) {
				ran.Log.Tracef("SERVED_TAI_INDEX[%d]", i)
				found = true
				break
			}
		}
		if !found {
			ran.Log.Warn("NG-Setup failure: Cannot find Served TAI in AMF")
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

func handleUplinkNASTransportMain(ran *context.AmfRan,
	ranUe *context.RanUe,
	nASPDU *ngapType.NASPDU,
	userLocationInformation *ngapType.UserLocationInformation,
) {
	amfUe := ranUe.AmfUe
	if amfUe == nil {
		err := ranUe.Remove()
		if err != nil {
			ran.Log.Error(err)
		}
		ran.Log.Errorf("No UE Context of RanUe with RANUENGAPID[%d] AMFUENGAPID[%d] ",
			ranUe.RanUeNgapId, ranUe.AmfUeNgapId)
		return
	}

	if userLocationInformation != nil {
		ranUe.UpdateLocation(userLocationInformation)
	}

	amf_nas.HandleNAS(ranUe, ngapType.ProcedureCodeUplinkNASTransport, nASPDU.Value, false)
}

func handleNGResetMain(ran *context.AmfRan,
	cause *ngapType.Cause,
	resetType *ngapType.ResetType,
) {
	if cause != nil {
		printAndGetCause(ran, cause)
	}

	switch resetType.Present {
	case ngapType.ResetTypePresentNGInterface:
		ran.Log.Trace("ResetType Present: NG Interface")
		ran.RemoveAllRanUe(false)
		ngap_message.SendNGResetAcknowledge(ran, nil, nil)
	case ngapType.ResetTypePresentPartOfNGInterface:
		ran.Log.Trace("ResetType Present: Part of NG Interface")

		partOfNGInterface := resetType.PartOfNGInterface
		if partOfNGInterface == nil {
			ran.Log.Error("PartOfNGInterface is nil")
			return
		}

		var ranUe *context.RanUe

		for _, ueAssociatedLogicalNGConnectionItem := range partOfNGInterface.List {
			if ueAssociatedLogicalNGConnectionItem.AMFUENGAPID != nil {
				ran.Log.Tracef("AmfUeNgapID[%d]", ueAssociatedLogicalNGConnectionItem.AMFUENGAPID.Value)
				ranUe = ran.FindRanUeByAmfUeNgapID(ueAssociatedLogicalNGConnectionItem.AMFUENGAPID.Value)
			} else if ueAssociatedLogicalNGConnectionItem.RANUENGAPID != nil {
				ran.Log.Tracef("RanUeNgapID[%d]", ueAssociatedLogicalNGConnectionItem.RANUENGAPID.Value)
				ranUe = ran.RanUeFindByRanUeNgapID(ueAssociatedLogicalNGConnectionItem.RANUENGAPID.Value)
			}

			if ranUe == nil {
				ran.Log.Warn("Cannot not find UE Context")
				if ueAssociatedLogicalNGConnectionItem.AMFUENGAPID != nil {
					ran.Log.Warnf("AmfUeNgapID[%d]", ueAssociatedLogicalNGConnectionItem.AMFUENGAPID.Value)
				}
				if ueAssociatedLogicalNGConnectionItem.RANUENGAPID != nil {
					ran.Log.Warnf("RanUeNgapID[%d]", ueAssociatedLogicalNGConnectionItem.RANUENGAPID.Value)
				}
			}

			err := ranUe.Remove()
			if err != nil {
				ran.Log.Error(err.Error())
			}
		}
		ngap_message.SendNGResetAcknowledge(ran, partOfNGInterface, nil)
	default:
		ran.Log.Warnf("Invalid ResetType[%d]", resetType.Present)
	}
}

func handleNGResetAcknowledgeMain(ran *context.AmfRan,
	uEAssociatedLogicalNGConnectionList *ngapType.UEAssociatedLogicalNGConnectionList,
	criticalityDiagnostics *ngapType.CriticalityDiagnostics,
) {
	if uEAssociatedLogicalNGConnectionList != nil {
		ran.Log.Tracef("%d UE association(s) has been reset", len(uEAssociatedLogicalNGConnectionList.List))
		for i, item := range uEAssociatedLogicalNGConnectionList.List {
			if item.AMFUENGAPID != nil && item.RANUENGAPID != nil {
				ran.Log.Tracef("%d: AmfUeNgapID[%d] RanUeNgapID[%d]", i+1, item.AMFUENGAPID.Value, item.RANUENGAPID.Value)
			} else if item.AMFUENGAPID != nil {
				ran.Log.Tracef("%d: AmfUeNgapID[%d] RanUeNgapID[-1]", i+1, item.AMFUENGAPID.Value)
			} else if item.RANUENGAPID != nil {
				ran.Log.Tracef("%d: AmfUeNgapID[-1] RanUeNgapID[%d]", i+1, item.RANUENGAPID.Value)
			}
		}
	}

	if criticalityDiagnostics != nil {
		printCriticalityDiagnostics(ran, criticalityDiagnostics)
	}
}

func handleUEContextReleaseCompleteMain(ran *context.AmfRan,
	ranUe *context.RanUe,
	userLocationInformation *ngapType.UserLocationInformation,
	infoOnRecommendedCellsAndRANNodesForPaging *ngapType.InfoOnRecommendedCellsAndRANNodesForPaging,
	pDUSessionResourceList *ngapType.PDUSessionResourceListCxtRelCpl,
	criticalityDiagnostics *ngapType.CriticalityDiagnostics,
) {
	if ranUe == nil {
		ran.Log.Error("ranUe is nil")
		return
	}

	if userLocationInformation != nil {
		ranUe.UpdateLocation(userLocationInformation)
	}
	if criticalityDiagnostics != nil {
		printCriticalityDiagnostics(ran, criticalityDiagnostics)
	}

	amfUe := ranUe.AmfUe
	if amfUe == nil {
		ran.Log.Infof("Release UE Context : RanUe[AmfUeNgapId: %d]", ranUe.AmfUeNgapId)
		err := ranUe.Remove()
		if err != nil {
			ran.Log.Errorln(err.Error())
		}
		return
	}
	// TODO: AMF shall, if supported, store it and may use it for subsequent paging
	if infoOnRecommendedCellsAndRANNodesForPaging != nil {
		amfUe.InfoOnRecommendedCellsAndRanNodesForPaging = new(context.InfoOnRecommendedCellsAndRanNodesForPaging)

		recommendedCells := &amfUe.InfoOnRecommendedCellsAndRanNodesForPaging.RecommendedCells
		for _, item := range infoOnRecommendedCellsAndRANNodesForPaging.RecommendedCellsForPaging.RecommendedCellList.List {
			recommendedCell := context.RecommendedCell{}

			switch item.NGRANCGI.Present {
			case ngapType.NGRANCGIPresentNRCGI:
				recommendedCell.NgRanCGI.Present = context.NgRanCgiPresentNRCGI
				recommendedCell.NgRanCGI.NRCGI = new(models.Ncgi)
				plmnID := ngapConvert.PlmnIdToModels(item.NGRANCGI.NRCGI.PLMNIdentity)
				recommendedCell.NgRanCGI.NRCGI.PlmnId = &plmnID
				recommendedCell.NgRanCGI.NRCGI.NrCellId = ngapConvert.BitStringToHex(&item.NGRANCGI.NRCGI.NRCellIdentity.Value)
			case ngapType.NGRANCGIPresentEUTRACGI:
				recommendedCell.NgRanCGI.Present = context.NgRanCgiPresentEUTRACGI
				recommendedCell.NgRanCGI.EUTRACGI = new(models.Ecgi)
				plmnID := ngapConvert.PlmnIdToModels(item.NGRANCGI.EUTRACGI.PLMNIdentity)
				recommendedCell.NgRanCGI.EUTRACGI.PlmnId = &plmnID
				recommendedCell.NgRanCGI.EUTRACGI.EutraCellId = ngapConvert.BitStringToHex(
					&item.NGRANCGI.EUTRACGI.EUTRACellIdentity.Value)
			}

			if item.TimeStayedInCell != nil {
				recommendedCell.TimeStayedInCell = new(int64)
				*recommendedCell.TimeStayedInCell = *item.TimeStayedInCell
			}

			*recommendedCells = append(*recommendedCells, recommendedCell)
		}

		recommendedRanNodes := &amfUe.InfoOnRecommendedCellsAndRanNodesForPaging.RecommendedRanNodes
		ranNodeList := infoOnRecommendedCellsAndRANNodesForPaging.RecommendRANNodesForPaging.RecommendedRANNodeList.List
		for _, item := range ranNodeList {
			recommendedRanNode := context.RecommendRanNode{}

			switch item.AMFPagingTarget.Present {
			case ngapType.AMFPagingTargetPresentGlobalRANNodeID:
				recommendedRanNode.Present = context.RecommendRanNodePresentRanNode
				recommendedRanNode.GlobalRanNodeId = new(models.GlobalRanNodeId)
				// TODO: recommendedRanNode.GlobalRanNodeId = ngapConvert.RanIdToModels(item.AMFPagingTarget.GlobalRANNodeID)
			case ngapType.AMFPagingTargetPresentTAI:
				recommendedRanNode.Present = context.RecommendRanNodePresentTAI
				tai := ngapConvert.TaiToModels(*item.AMFPagingTarget.TAI)
				recommendedRanNode.Tai = &tai
			}
			*recommendedRanNodes = append(*recommendedRanNodes, recommendedRanNode)
		}
	}

	// for each pduSessionID invoke Nsmf_PDUSession_UpdateSMContext Request
	var cause context.CauseAll
	if tmp, exist := amfUe.ReleaseCause[ran.AnType]; exist {
		cause = *tmp
	}
	if amfUe.State[ran.AnType].Is(context.Registered) {
		ranUe.Log.Info("Release Ue Context in GMM-Registered")
		// If this release cause by handover, no needs deactivate CN tunnel
		if cause.NgapCause != nil && pDUSessionResourceList != nil {
			for _, pduSessionReourceItem := range pDUSessionResourceList.List {
				pduSessionID := int32(pduSessionReourceItem.PDUSessionID.Value)
				smContext, ok := amfUe.SmContextFindByPDUSessionID(pduSessionID)
				if !ok {
					ranUe.Log.Warnf("SmContext[PDU Session ID:%d] not found", pduSessionID)
					// TODO: Check if doing error handling here
					continue
				}
				response, _, _, err := consumer.GetConsumer().SendUpdateSmContextDeactivateUpCnxState(amfUe, smContext, cause)
				if err != nil {
					ran.Log.Errorf("Send Update SmContextDeactivate UpCnxState Error[%s]", err.Error())
				} else if response == nil {
					ran.Log.Errorln("Send Update SmContextDeactivate UpCnxState Error")
				}
			}
		}
	}

	// TODO: stop timer and release RanUe context
	// Remove UE N2 Connection
	delete(amfUe.ReleaseCause, ran.AnType)
	switch ranUe.ReleaseAction {
	case context.UeContextN2NormalRelease:
		ran.Log.Infof("Release UE[%s] Context : N2 Connection Release", amfUe.Supi)
		// amfUe.DetachRanUe(ran.AnType)
		err := ranUe.Remove()
		if err != nil {
			ran.Log.Errorln(err.Error())
		}
	case context.UeContextReleaseUeContext:
		ran.Log.Infof("Release UE[%s] Context : Release Ue Context", amfUe.Supi)
		amfUe.Lock.Lock()
		gmm_common.RemoveAmfUe(amfUe, false)
		amfUe.Lock.Unlock()
	case context.UeContextReleaseHandover:
		ran.Log.Infof("Release UE[%s] Context : Release for Handover", amfUe.Supi)
		// TODO: it's a workaround, need to fix it.
		targetRanUe := context.GetSelf().RanUeFindByAmfUeNgapID(ranUe.TargetUe.AmfUeNgapId)

		context.DetachSourceUeTargetUe(ranUe)
		err := ranUe.Remove()
		if err != nil {
			ran.Log.Errorln(err.Error())
		}
		gmm_common.AttachRanUeToAmfUeAndReleaseOldIfAny(amfUe, targetRanUe)
		// Todo: remove indirect tunnel
	default:
		ran.Log.Errorf("Invalid Release Action[%d]", ranUe.ReleaseAction)
	}
}

func handlePDUSessionResourceReleaseResponseMain(ran *context.AmfRan,
	ranUe *context.RanUe,
	pDUSessionResourceReleasedList *ngapType.PDUSessionResourceReleasedListRelRes,
	userLocationInformation *ngapType.UserLocationInformation,
	criticalityDiagnostics *ngapType.CriticalityDiagnostics,
) {
	if ranUe == nil {
		ran.Log.Error("ranUe is nil")
		return
	}

	if userLocationInformation != nil {
		ranUe.UpdateLocation(userLocationInformation)
	}

	if criticalityDiagnostics != nil {
		printCriticalityDiagnostics(ran, criticalityDiagnostics)
	}

	amfUe := ranUe.AmfUe
	if amfUe == nil {
		ranUe.Log.Error("amfUe is nil")
		return
	}
	if pDUSessionResourceReleasedList != nil {
		ranUe.Log.Infof("Send PDUSessionResourceReleaseResponseTransfer to SMF")

		for _, item := range pDUSessionResourceReleasedList.List {
			pduSessionID := int32(item.PDUSessionID.Value)
			transfer := item.PDUSessionResourceReleaseResponseTransfer
			smContext, ok := amfUe.SmContextFindByPDUSessionID(pduSessionID)
			if !ok {
				// TODO: Check if NAS (PDU Session Release Complete) comes before PDUSesstionResourceRelease
				ranUe.Log.Warnf("SmContext[PDU Session ID:%d] not found", pduSessionID)
				// TODO: Check if doing error handling here
				continue
			}

			_, responseErr, problemDetail, err := consumer.GetConsumer().SendUpdateSmContextN2Info(amfUe, smContext,
				models.N2SmInfoType_PDU_RES_REL_RSP, transfer)
			// TODO: error handling
			if err != nil {
				ranUe.Log.Errorf("SendUpdateSmContextN2Info[PDUSessionResourceReleaseResponse] Error: %+v", err)
			} else if responseErr != nil && responseErr.JsonData.Error != nil {
				ranUe.Log.Errorf("SendUpdateSmContextN2Info[PDUSessionResourceReleaseResponse] Error: %+v",
					responseErr.JsonData.Error.Cause)
			} else if problemDetail != nil {
				ranUe.Log.Errorf("SendUpdateSmContextN2Info[PDUSessionResourceReleaseResponse] Failed: %+v", problemDetail)
			}
		}
	}
}

func handleUERadioCapabilityCheckResponseMain(ran *context.AmfRan,
	ranUe *context.RanUe,
	criticalityDiagnostics *ngapType.CriticalityDiagnostics,
) {
	// TODO: handle iMSVoiceSupportIndicator

	if criticalityDiagnostics != nil {
		printCriticalityDiagnostics(ran, criticalityDiagnostics)
	}
}

func handleLocationReportingFailureIndicationMain(ran *context.AmfRan,
	ranUe *context.RanUe,
	cause *ngapType.Cause,
) {
	if cause != nil {
		printAndGetCause(ran, cause)
	}
}

func handleInitialUEMessageMain(ran *context.AmfRan,
	message *ngapType.NGAPPDU,
	rANUENGAPID *ngapType.RANUENGAPID,
	nASPDU *ngapType.NASPDU,
	userLocationInformation *ngapType.UserLocationInformation,
	rRCEstablishmentCause *ngapType.RRCEstablishmentCause,
	fiveGSTMSI *ngapType.FiveGSTMSI,
	uEContextRequest *ngapType.UEContextRequest,
) {
	ranUe := ran.RanUeFindByRanUeNgapID(rANUENGAPID.Value)
	if ranUe != nil {
		amfUe := ranUe.AmfUe
		if amfUe != nil {
			// The fact that an amfUe having N2 connection (ranUE) is receiving
			// an Initial UE Message indicates there is something wrong,
			// so the ranUe with wrong RAN-UE-NGAP-IP should be cleared and detached from the amfUe.
			gmm_common.StopAll5GSMMTimers(amfUe)
			amfUe.DetachRanUe(ran.AnType)
			ranUe.DetachAmfUe()
		}
		err := ranUe.Remove()
		if err != nil {
			ran.Log.Errorln(err.Error())
		}
	}

	var err error
	ranUe, err = ran.NewRanUe(rANUENGAPID.Value)
	if err != nil {
		ran.Log.Errorf("NewRanUe Error: %+v", err)
	}
	ran.Log.Debugf("New RanUe [RanUeNgapID: %d]", ranUe.RanUeNgapId)

	// Try to get identity from 5G-S-TMSI IE first; if not available, try to get identity from the plain NAS.
	var id, idType string
	var gmmMessage *nas.GmmMessage
	var nasMsgType, regReqType uint8
	// Get nasMsgType to send corresponding NAS reject to UE when amfUe is not found.
	nasMsg, err := nas_security.DecodePlainNasNoIntegrityCheck(nASPDU.Value)
	if err == nil && nasMsg.GmmMessage != nil {
		gmmMessage = nasMsg.GmmMessage
		nasMsgType = gmmMessage.GmmHeader.GetMessageType()
		if gmmMessage.RegistrationRequest != nil {
			regReqType = gmmMessage.RegistrationRequest.NgksiAndRegistrationType5GS.GetRegistrationType5GS()
		}
	}

	if fiveGSTMSI != nil {
		// <5G-S-TMSI> := <AMF Set ID><AMF Pointer><5G-TMSI>
		// GUAMI := <MCC><MNC><AMF Region ID><AMF Set ID><AMF Pointer>
		// 5G-GUTI := <GUAMI><5G-TMSI>
		amfSetPtrID := hex.EncodeToString([]byte{
			fiveGSTMSI.AMFSetID.Value.Bytes[0],
			(fiveGSTMSI.AMFSetID.Value.Bytes[1] & 0xc0) | (fiveGSTMSI.AMFPointer.Value.Bytes[0] >> 2),
		})
		tmsi := hex.EncodeToString(fiveGSTMSI.FiveGTMSI.Value)

		id = amfSetPtrID + tmsi
		idType = "5G-S-TMSI"
		ranUe.Log.Infof("Find 5G-S-TMSI [%q] in InitialUEMessage", id)
	} else if regReqType == nasMessage.RegistrationType5GSInitialRegistration {
		// NGAP 5G-S-TMSI IE might not be present in InitialUEMessage carrying Initial Registration.
		// Need to get 5GSMobileIdentity from Initial Registration.

		id, idType, err = amf_nas.GetNas5GSMobileIdentity(gmmMessage)
		ran.Log.Infof("5GSMobileIdentity [%q:%q, err: %v]", idType, id, err)
	} else {
		// Missing NGAP 5G-S-TMSI IE
		var iesCriticalityDiagnostics ngapType.CriticalityDiagnosticsIEList
		ranUe.Log.Warnf("Missing 5G-S-TMSI IE in InitialUEMessage; send ErrorIndication")
		item := buildCriticalityDiagnosticsIEItem(ngapType.CriticalityPresentReject,
			ngapType.ProtocolIEIDFiveGSTMSI, ngapType.TypeOfErrorPresentMissing)
		iesCriticalityDiagnostics.List = append(iesCriticalityDiagnostics.List, item)
		sendErrorMessage(ran, nil, rANUENGAPID, iesCriticalityDiagnostics)

		ngap_message.SendUEContextReleaseCommand(ranUe, context.UeContextN2NormalRelease,
			ngapType.CausePresentProtocol, ngapType.CauseProtocolPresentUnspecified)
		return
	}

	// If id type is GUTI, since MAC can't be checked here (no amfUe context), the GUTI may not direct to the right amfUe.
	// In this case, create a new amfUe to handle the following registration procedure.
	isInvalidGUTI := (idType == "5G-GUTI")
	amfUe, ok := findAmfUe(ran, id, idType)
	if ok && !isInvalidGUTI {
		// TODO: invoke Namf_Communication_UEContextTransfer if serving AMF has changed since
		// last Registration Request procedure
		// Described in TS 23.502 4.2.2.2.2 step 4 (without UDSF deployment)
		ranUe.Log.Infof("find AmfUe [%q:%q]", idType, id)
		ranUe.Log.Debugf("AmfUe Attach RanUe [RanUeNgapID: %d]", ranUe.RanUeNgapId)
		ranUe.HoldingAmfUe = amfUe
	} else if regReqType != nasMessage.RegistrationType5GSInitialRegistration {
		if regReqType == nasMessage.RegistrationType5GSPeriodicRegistrationUpdating ||
			regReqType == nasMessage.RegistrationType5GSMobilityRegistrationUpdating {
			gmm_message.SendRegistrationReject(
				ranUe, nasMessage.Cause5GMMImplicitlyDeregistered, "")
			ranUe.Log.Warn("Send RegistrationReject [Cause5GMMImplicitlyDeregistered]")
		} else if nasMsgType == nas.MsgTypeServiceRequest {
			gmm_message.SendServiceReject(
				ranUe, nil, nasMessage.Cause5GMMImplicitlyDeregistered)
			ranUe.Log.Warn("Send ServiceReject [Cause5GMMImplicitlyDeregistered]")
		}

		ngap_message.SendUEContextReleaseCommand(ranUe, context.UeContextN2NormalRelease,
			ngapType.CausePresentNas, ngapType.CauseNasPresentNormalRelease)
		return
	}

	if userLocationInformation != nil {
		ranUe.UpdateLocation(userLocationInformation)
	}

	if rRCEstablishmentCause != nil {
		ranUe.Log.Tracef("[Initial UE Message] RRC Establishment Cause[%d]", rRCEstablishmentCause.Value)
		ranUe.RRCEstablishmentCause = strconv.Itoa(int(rRCEstablishmentCause.Value))
	}

	if uEContextRequest != nil {
		ran.Log.Debug("Trigger initial Context Setup procedure")
		ranUe.UeContextRequest = true
		// TODO: Trigger Initial Context Setup procedure
	} else {
		ranUe.UeContextRequest = factory.AmfConfig.Configuration.DefaultUECtxReq
	}

	// TS 23.502 4.2.2.2.3 step 6a Nnrf_NFDiscovery_Request (NF type, AMF Set)
	// if aMFSetID != nil {
	// TODO: This is a rerouted message
	// TS 38.413: AMF shall, if supported, use the IE as described in TS 23.502
	// }

	// ng-ran propagate allowedNssai in the rerouted initial ue message (TS 38.413 8.6.5)
	// TS 23.502 4.2.2.2.3 step 4a Nnssf_NSSelection_Get
	// if allowedNSSAI != nil {
	// TODO: AMF should use it as defined in TS 23.502
	// }

	pdu, err := libngap.Encoder(*message)
	if err != nil {
		ran.Log.Errorf("libngap Encoder Error: %+v", err)
	}
	ranUe.InitialUEMessage = pdu
	amf_nas.HandleNAS(ranUe, ngapType.ProcedureCodeInitialUEMessage, nASPDU.Value, true)
}

func findAmfUe(ran *context.AmfRan, id, idType string) (*context.AmfUe, bool) {
	var amfUe *context.AmfUe
	var ok bool

	amfSelf := context.GetSelf()
	servedGuami := amfSelf.ServedGuamiList[0]
	tmpRegionID, _, _ := ngapConvert.AmfIdToNgap(servedGuami.AmfId)

	switch idType {
	case "SUPI":
		ran.Log.Debugf("SUPI %s", id)
		amfUe, ok = amfSelf.AmfUeFindBySupi(id)
	case "SUCI":
		ran.Log.Debugf("SUCI %s", id)
		amfUe, ok = amfSelf.AmfUeFindBySuci(id)
	case "5G-GUTI":
		ran.Log.Debugf("5G-GUTI %s", id)
		amfUe, ok = amfSelf.AmfUeFindByGuti(id)
	case "5G-S-TMSI":
		id = servedGuami.PlmnId.Mcc + servedGuami.PlmnId.Mnc + ngapConvert.BitStringToHex(&tmpRegionID) + id
		ran.Log.Debugf("5G-S-TMSI %s", id)
		amfUe, ok = amfSelf.AmfUeFindByGuti(id)
	}
	return amfUe, ok
}

func sendErrorMessage(ran *context.AmfRan, amfUeNgapId *ngapType.AMFUENGAPID, ranUeNgapId *ngapType.RANUENGAPID,
	iesCriticalityDiagnostics ngapType.CriticalityDiagnosticsIEList,
) {
	ran.Log.Trace("Has missing reject IE(s)")

	procedureCode := ngapType.ProcedureCodeInitialUEMessage
	triggeringMessage := ngapType.TriggeringMessagePresentInitiatingMessage
	procedureCriticality := ngapType.CriticalityPresentIgnore
	criticalityDiagnostics := buildCriticalityDiagnostics(&procedureCode, &triggeringMessage, &procedureCriticality,
		&iesCriticalityDiagnostics)
	ngap_message.SendErrorIndication(ran, amfUeNgapId, ranUeNgapId, nil, &criticalityDiagnostics)
}

func handlePDUSessionResourceSetupResponseMain(ran *context.AmfRan,
	ranUe *context.RanUe,
	pDUSessionResourceSetupResponseList *ngapType.PDUSessionResourceSetupListSURes,
	pDUSessionResourceFailedToSetupList *ngapType.PDUSessionResourceFailedToSetupListSURes,
	criticalityDiagnostics *ngapType.CriticalityDiagnostics,
) {
	if ranUe == nil {
		ran.Log.Error("ranUe is nil")
		return
	}

	amfUe := ranUe.AmfUe
	if amfUe == nil {
		ranUe.Log.Error("amfUe is nil")
		return
	}

	if pDUSessionResourceSetupResponseList != nil {
		ranUe.Log.Trace("Send PDUSessionResourceSetupResponseTransfer to SMF")

		for _, item := range pDUSessionResourceSetupResponseList.List {
			pduSessionID := int32(item.PDUSessionID.Value)
			transfer := item.PDUSessionResourceSetupResponseTransfer
			smContext, ok := amfUe.SmContextFindByPDUSessionID(pduSessionID)
			if !ok {
				ranUe.Log.Errorf("SmContext[PDU Session ID:%d] not found", pduSessionID)
				continue
			}
			_, _, _, err := consumer.GetConsumer().SendUpdateSmContextN2Info(amfUe, smContext,
				models.N2SmInfoType_PDU_RES_SETUP_RSP, transfer)
			if err != nil {
				ranUe.Log.Errorf("SendUpdateSmContextN2Info[PDUSessionResourceSetupResponseTransfer] Error: %+v", err)
			}
			// RAN initiated QoS Flow Mobility in subclause 5.2.2.3.7
			// if response != nil && response.BinaryDataN2SmInformation != nil {
			// TODO: n2SmInfo send to RAN
			// } else if response == nil {
			// TODO: error handling
			// }
		}
	}

	if pDUSessionResourceFailedToSetupList != nil {
		ranUe.Log.Trace("Send PDUSessionResourceSetupUnsuccessfulTransfer to SMF")

		for _, item := range pDUSessionResourceFailedToSetupList.List {
			pduSessionID := int32(item.PDUSessionID.Value)
			transfer := item.PDUSessionResourceSetupUnsuccessfulTransfer
			smContext, ok := amfUe.SmContextFindByPDUSessionID(pduSessionID)
			if !ok {
				ranUe.Log.Errorf("SmContext[PDU Session ID:%d] not found", pduSessionID)
				continue
			}
			_, _, _, err := consumer.GetConsumer().SendUpdateSmContextN2Info(amfUe, smContext,
				models.N2SmInfoType_PDU_RES_SETUP_FAIL, transfer)
			if err != nil {
				ranUe.Log.Errorf("SendUpdateSmContextN2Info[PDUSessionResourceSetupUnsuccessfulTransfer] Error: %+v", err)
			}

			// if response != nil && response.BinaryDataN2SmInformation != nil {
			// TODO: n2SmInfo send to RAN
			// } else if response == nil {
			// TODO: error handling
			// }
		}
	}

	if criticalityDiagnostics != nil {
		printCriticalityDiagnostics(ran, criticalityDiagnostics)
	}
}

func handlePDUSessionResourceModifyResponseMain(ran *context.AmfRan,
	ranUe *context.RanUe,
	pduSessionResourceModifyResponseList *ngapType.PDUSessionResourceModifyListModRes,
	pduSessionResourceFailedToModifyList *ngapType.PDUSessionResourceFailedToModifyListModRes,
	userLocationInformation *ngapType.UserLocationInformation,
	criticalityDiagnostics *ngapType.CriticalityDiagnostics,
) {
	if ranUe == nil {
		ran.Log.Error("ranUe is nil")
		return
	}

	amfUe := ranUe.AmfUe
	if amfUe == nil {
		ranUe.Log.Error("amfUe is nil")
		return
	}

	if pduSessionResourceModifyResponseList != nil {
		ranUe.Log.Trace("Send PDUSessionResourceModifyResponseTransfer to SMF")

		for _, item := range pduSessionResourceModifyResponseList.List {
			pduSessionID := int32(item.PDUSessionID.Value)
			transfer := item.PDUSessionResourceModifyResponseTransfer
			smContext, ok := amfUe.SmContextFindByPDUSessionID(pduSessionID)
			if !ok {
				ranUe.Log.Errorf("SmContext[PDU Session ID:%d] not found", pduSessionID)
				continue
			}
			_, _, _, err := consumer.GetConsumer().SendUpdateSmContextN2Info(amfUe, smContext,
				models.N2SmInfoType_PDU_RES_MOD_RSP, transfer)
			if err != nil {
				ranUe.Log.Errorf("SendUpdateSmContextN2Info[PDUSessionResourceModifyResponseTransfer] Error: %+v", err)
			}
			// if response != nil && response.BinaryDataN2SmInformation != nil {
			// TODO: n2SmInfo send to RAN
			// } else if response == nil {
			// TODO: error handling
			// }
		}
	}

	if pduSessionResourceFailedToModifyList != nil {
		ranUe.Log.Trace("Send PDUSessionResourceModifyUnsuccessfulTransfer to SMF")

		for _, item := range pduSessionResourceFailedToModifyList.List {
			pduSessionID := int32(item.PDUSessionID.Value)
			transfer := item.PDUSessionResourceModifyUnsuccessfulTransfer
			smContext, ok := amfUe.SmContextFindByPDUSessionID(pduSessionID)
			if !ok {
				ranUe.Log.Errorf("SmContext[PDU Session ID:%d] not found", pduSessionID)
				continue
			}
			_, _, _, err := consumer.GetConsumer().SendUpdateSmContextN2Info(amfUe, smContext,
				models.N2SmInfoType_PDU_RES_MOD_FAIL, transfer)
			if err != nil {
				ranUe.Log.Errorf("SendUpdateSmContextN2Info[PDUSessionResourceModifyUnsuccessfulTransfer] Error: %+v", err)
			}
			// if response != nil && response.BinaryDataN2SmInformation != nil {
			// TODO: n2SmInfo send to RAN
			// } else if response == nil {
			// TODO: error handling
			// }
		}
	}

	if userLocationInformation != nil {
		ranUe.UpdateLocation(userLocationInformation)
	}

	if criticalityDiagnostics != nil {
		printCriticalityDiagnostics(ran, criticalityDiagnostics)
	}
}

func handlePDUSessionResourceNotifyMain(ran *context.AmfRan,
	ranUe *context.RanUe,
	pDUSessionResourceNotifyList *ngapType.PDUSessionResourceNotifyList,
	pDUSessionResourceReleasedListNot *ngapType.PDUSessionResourceReleasedListNot,
	userLocationInformation *ngapType.UserLocationInformation,
) {
	amfUe := ranUe.AmfUe
	if amfUe == nil {
		ranUe.Log.Error("amfUe is nil")
		return
	}

	if userLocationInformation != nil {
		ranUe.UpdateLocation(userLocationInformation)
	}

	if pDUSessionResourceNotifyList != nil {
		ranUe.Log.Infof("Send PDUSessionResourceNotifyTransfer to SMF")
		for _, item := range pDUSessionResourceNotifyList.List {
			pduSessionID := int32(item.PDUSessionID.Value)
			transfer := item.PDUSessionResourceNotifyTransfer
			smContext, ok := amfUe.SmContextFindByPDUSessionID(pduSessionID)
			if !ok {
				ranUe.Log.Errorf("SmContext[PDU Session ID:%d] not found", pduSessionID)
				continue
			}
			response, errResponse, problemDetail, err := consumer.GetConsumer().SendUpdateSmContextN2Info(amfUe, smContext,
				models.N2SmInfoType_PDU_RES_NTY, transfer)
			if err != nil {
				ranUe.Log.Errorf("SendUpdateSmContextN2Info[PDUSessionResourceNotifyTransfer] Error: %+v", err)
			}

			if response != nil {
				responseData := response.JsonData
				n2Info := response.BinaryDataN1SmMessage
				n1Msg := response.BinaryDataN2SmInformation
				if n2Info != nil {
					switch responseData.N2SmInfoType {
					case models.N2SmInfoType_PDU_RES_MOD_REQ:
						ranUe.Log.Debugln("AMF Transfer NGAP PDU Resource Modify Req from SMF")
						var nasPdu []byte
						if n1Msg != nil {
							pduSessionId := uint8(pduSessionID)
							nasPdu, err = gmm_message.BuildDLNASTransport(amfUe, ran.AnType, nasMessage.PayloadContainerTypeN1SMInfo,
								n1Msg, pduSessionId, nil, nil, 0)
							if err != nil {
								ranUe.Log.Warnf("GMM Message build DL NAS Transport filaed: %v", err)
							}
						}
						list := ngapType.PDUSessionResourceModifyListModReq{}
						ngap_message.AppendPDUSessionResourceModifyListModReq(&list, pduSessionID, nasPdu, n2Info)
						ngap_message.SendPDUSessionResourceModifyRequest(ranUe, list)
					default:
					}
				}
			} else if errResponse != nil {
				errJSON := errResponse.JsonData
				n1Msg := errResponse.BinaryDataN2SmInformation
				ranUe.Log.Warnf("PDU Session Modification is rejected by SMF[pduSessionId:%d], Error[%s]\n",
					pduSessionID, errJSON.Error.Cause)
				if n1Msg != nil {
					gmm_message.SendDLNASTransport(
						ranUe, nasMessage.PayloadContainerTypeN1SMInfo, errResponse.BinaryDataN1SmMessage, pduSessionID, 0, nil, 0)
				}
				// TODO: handle n2 info transfer
			} else if err != nil {
				return
			} else {
				// TODO: error handling
				ranUe.Log.Errorf("Failed to Update smContext[pduSessionID: %d], Error[%v]", pduSessionID, problemDetail)
				return
			}
		}
	}

	if pDUSessionResourceReleasedListNot != nil {
		ranUe.Log.Infof("Send PDUSessionResourceNotifyReleasedTransfer to SMF")
		for _, item := range pDUSessionResourceReleasedListNot.List {
			pduSessionID := int32(item.PDUSessionID.Value)
			transfer := item.PDUSessionResourceNotifyReleasedTransfer
			smContext, ok := amfUe.SmContextFindByPDUSessionID(pduSessionID)
			if !ok {
				ranUe.Log.Warnf("SmContext[PDU Session ID:%d] not found", pduSessionID)
				// TODO: Check if doing error handling here
				continue
			}
			response, errResponse, problemDetail, err := consumer.GetConsumer().SendUpdateSmContextN2Info(amfUe, smContext,
				models.N2SmInfoType_PDU_RES_NTY_REL, transfer)
			if err != nil {
				ranUe.Log.Errorf("SendUpdateSmContextN2Info[PDUSessionResourceNotifyReleasedTransfer] Error: %+v", err)
			}
			if response != nil {
				responseData := response.JsonData
				n2Info := response.BinaryDataN1SmMessage
				n1Msg := response.BinaryDataN2SmInformation
				if n2Info != nil {
					if responseData.N2SmInfoType == models.N2SmInfoType_PDU_RES_REL_CMD {
						ranUe.Log.Debugln("AMF Transfer NGAP PDU Session Resource Rel Co from SMF")
						var nasPdu []byte
						if n1Msg != nil {
							nasPdu, err = gmm_message.BuildDLNASTransport(
								amfUe, ran.AnType, nasMessage.PayloadContainerTypeN1SMInfo, n1Msg,
								uint8(pduSessionID), nil, nil, 0)
							if err != nil {
								ranUe.Log.Warnf("GMM Message build DL NAS Transport filaed: %v", err)
							}
						}
						list := ngapType.PDUSessionResourceToReleaseListRelCmd{}
						ngap_message.AppendPDUSessionResourceToReleaseListRelCmd(&list, pduSessionID, n2Info)
						ngap_message.SendPDUSessionResourceReleaseCommand(ranUe, nasPdu, list)
					}
				}
			} else if errResponse != nil {
				errJSON := errResponse.JsonData
				n1Msg := errResponse.BinaryDataN2SmInformation
				ranUe.Log.Warnf("PDU Session Release is rejected by SMF[pduSessionID:%d], Error[%s]\n",
					pduSessionID, errJSON.Error.Cause)
				if n1Msg != nil {
					gmm_message.SendDLNASTransport(
						ranUe, nasMessage.PayloadContainerTypeN1SMInfo, errResponse.BinaryDataN1SmMessage, pduSessionID, 0, nil, 0)
				}
			} else if err != nil {
				return
			} else {
				// TODO: error handling
				ranUe.Log.Errorf("Failed to Update smContext[pduSessionID: %d], Error[%v]", pduSessionID, problemDetail)
				return
			}
		}
	}
}

func handlePDUSessionResourceModifyIndicationMain(ran *context.AmfRan,
	ranUe *context.RanUe,
	pduSessionResourceModifyIndicationList *ngapType.PDUSessionResourceModifyListModInd,
) {
	amfUe := ranUe.AmfUe
	if amfUe == nil {
		ran.Log.Error("AmfUe is nil")
		return
	}

	pduSessionResourceModifyListModCfm := ngapType.PDUSessionResourceModifyListModCfm{}
	pduSessionResourceFailedToModifyListModCfm := ngapType.PDUSessionResourceFailedToModifyListModCfm{}

	if pduSessionResourceModifyIndicationList != nil {
		ran.Log.Infof("Send PDUSessionResourceModifyIndicationTransfer to SMF")
		for _, item := range pduSessionResourceModifyIndicationList.List {
			pduSessionID := int32(item.PDUSessionID.Value)
			transfer := item.PDUSessionResourceModifyIndicationTransfer
			smContext, ok := amfUe.SmContextFindByPDUSessionID(pduSessionID)
			if !ok {
				ranUe.Log.Warnf("SmContext[PDU Session ID:%d] not found", pduSessionID)
				// TODO: Check if doing error handling here
				continue
			}
			response, errResponse, _, err := consumer.GetConsumer().SendUpdateSmContextN2Info(amfUe, smContext,
				models.N2SmInfoType_PDU_RES_MOD_IND, transfer)
			if err != nil {
				ran.Log.Errorf("SendUpdateSmContextN2Info Error:\n%s", err.Error())
			}

			if response != nil && response.BinaryDataN2SmInformation != nil {
				ngap_message.AppendPDUSessionResourceModifyListModCfm(
					&pduSessionResourceModifyListModCfm,
					int64(pduSessionID), response.BinaryDataN2SmInformation)
			}
			if errResponse != nil && errResponse.BinaryDataN2SmInformation != nil {
				ngap_message.AppendPDUSessionResourceFailedToModifyListModCfm(
					&pduSessionResourceFailedToModifyListModCfm,
					int64(pduSessionID), errResponse.BinaryDataN2SmInformation)
			}
		}
	}

	ngap_message.SendPDUSessionResourceModifyConfirm(ranUe, pduSessionResourceModifyListModCfm,
		pduSessionResourceFailedToModifyListModCfm, nil)
}

func handleInitialContextSetupResponseMain(ran *context.AmfRan,
	ranUe *context.RanUe,
	pDUSessionResourceSetupResponseList *ngapType.PDUSessionResourceSetupListCxtRes,
	pDUSessionResourceFailedToSetupList *ngapType.PDUSessionResourceFailedToSetupListCxtRes,
	criticalityDiagnostics *ngapType.CriticalityDiagnostics,
) {
	if ranUe == nil {
		ran.Log.Error("ranUe is nil")
		return
	}

	amfUe := ranUe.AmfUe
	if amfUe == nil {
		ran.Log.Error("amfUe is nil")
		return
	}

	ran.Log.Tracef("RanUeNgapID[%d] AmfUeNgapID[%d]", ranUe.RanUeNgapId, ranUe.AmfUeNgapId)
	ranUe.InitialContextSetup = true

	if pDUSessionResourceSetupResponseList != nil {
		ranUe.Log.Infof("Send PDUSessionResourceSetupResponseTransfer to SMF")

		for _, item := range pDUSessionResourceSetupResponseList.List {
			pduSessionID := int32(item.PDUSessionID.Value)
			transfer := item.PDUSessionResourceSetupResponseTransfer
			smContext, ok := amfUe.SmContextFindByPDUSessionID(pduSessionID)
			if !ok {
				ranUe.Log.Warnf("SmContext[PDU Session ID:%d] not found", pduSessionID)
				// TODO: Check if doing error handling here
				continue
			}
			_, _, _, err := consumer.GetConsumer().SendUpdateSmContextN2Info(amfUe, smContext,
				models.N2SmInfoType_PDU_RES_SETUP_RSP, transfer)
			if err != nil {
				ranUe.Log.Errorf("SendUpdateSmContextN2Info[PDUSessionResourceSetupResponseTransfer] Error: %+v", err)
			}
			// RAN initiated QoS Flow Mobility in subclause 5.2.2.3.7
			// if response != nil && response.BinaryDataN2SmInformation != nil {
			// TODO: n2SmInfo send to RAN
			// } else if response == nil {
			// TODO: error handling
			// }
		}
	}

	if pDUSessionResourceFailedToSetupList != nil {
		ranUe.Log.Infof("Send PDUSessionResourceSetupUnsuccessfulTransfer to SMF")

		for _, item := range pDUSessionResourceFailedToSetupList.List {
			pduSessionID := int32(item.PDUSessionID.Value)
			transfer := item.PDUSessionResourceSetupUnsuccessfulTransfer
			smContext, ok := amfUe.SmContextFindByPDUSessionID(pduSessionID)
			if !ok {
				ranUe.Log.Warnf("SmContext[PDU Session ID:%d] not found", pduSessionID)
				// TODO: Check if doing error handling here
				continue
			}
			_, _, _, err := consumer.GetConsumer().SendUpdateSmContextN2Info(amfUe, smContext,
				models.N2SmInfoType_PDU_RES_SETUP_FAIL, transfer)
			if err != nil {
				ranUe.Log.Errorf("SendUpdateSmContextN2Info[PDUSessionResourceSetupUnsuccessfulTransfer] Error: %+v", err)
			}

			// if response != nil && response.BinaryDataN2SmInformation != nil {
			// TODO: n2SmInfo send to RAN
			// } else if response == nil {
			// TODO: error handling
			// }
		}
	}

	if ranUe.Ran.AnType == models.AccessType_NON_3_GPP_ACCESS {
		ngap_message.SendDownlinkNasTransport(ranUe, amfUe.RegistrationAcceptForNon3GPPAccess, nil)
	}

	if criticalityDiagnostics != nil {
		printCriticalityDiagnostics(ran, criticalityDiagnostics)
	}
}

func handleInitialContextSetupFailureMain(ran *context.AmfRan,
	ranUe *context.RanUe,
	pDUSessionResourceFailedToSetupList *ngapType.PDUSessionResourceFailedToSetupListCxtFail,
	cause *ngapType.Cause,
	criticalityDiagnostics *ngapType.CriticalityDiagnostics,
) {
	if cause != nil {
		printAndGetCause(ran, cause)
	}

	if criticalityDiagnostics != nil {
		printCriticalityDiagnostics(ran, criticalityDiagnostics)
	}

	if ranUe == nil {
		ran.Log.Error("ranUe is nil")
		return
	}

	amfUe := ranUe.AmfUe
	if amfUe == nil {
		ran.Log.Error("amfUe is nil")
		return
	}

	if pDUSessionResourceFailedToSetupList != nil {
		ranUe.Log.Infof("Send PDUSessionResourceSetupUnsuccessfulTransfer to SMF")

		for _, item := range pDUSessionResourceFailedToSetupList.List {
			pduSessionID := int32(item.PDUSessionID.Value)
			transfer := item.PDUSessionResourceSetupUnsuccessfulTransfer
			smContext, ok := amfUe.SmContextFindByPDUSessionID(pduSessionID)
			if !ok {
				ranUe.Log.Warnf("SmContext[PDU Session ID:%d] not found", pduSessionID)
				// TODO: Check if doing error handling here
				continue
			}
			_, _, _, err := consumer.GetConsumer().SendUpdateSmContextN2Info(amfUe, smContext,
				models.N2SmInfoType_PDU_RES_SETUP_FAIL, transfer)
			if err != nil {
				ranUe.Log.Errorf("SendUpdateSmContextN2Info[PDUSessionResourceSetupUnsuccessfulTransfer] Error: %+v", err)
			}

			// if response != nil && response.BinaryDataN2SmInformation != nil {
			// TODO: n2SmInfo send to RAN
			// } else if response == nil {
			// TODO: error handling
			// }
		}
	}
}

func handleUEContextReleaseRequestMain(ran *context.AmfRan,
	ranUe *context.RanUe,
	pDUSessionResourceList *ngapType.PDUSessionResourceListCxtRelReq,
	cause *ngapType.Cause,
) {
	causeGroup := ngapType.CausePresentRadioNetwork
	causeValue := ngapType.CauseRadioNetworkPresentUnspecified
	if cause != nil {
		causeGroup, causeValue = printAndGetCause(ran, cause)
	}

	amfUe := ranUe.AmfUe
	if amfUe != nil {
		if !isLatestAmfUe(amfUe) {
			amfUe.DetachRanUe(ran.AnType)
			ranUe.DetachAmfUe()
			gmm_common.StopAll5GSMMTimers(amfUe)
			causeValue = ngapType.CauseRadioNetworkPresentReleaseDueToNgranGeneratedReason
			ngap_message.SendUEContextReleaseCommand(ranUe, context.UeContextReleaseUeContext, causeGroup, causeValue)
			return
		}
		gmm_common.StopAll5GSMMTimers(amfUe)
		causeAll := context.CauseAll{
			NgapCause: &models.NgApCause{
				Group: int32(causeGroup),
				Value: int32(causeValue),
			},
		}
		if amfUe.State[ran.AnType].Is(context.Registered) {
			ranUe.Log.Info("Ue Context in GMM-Registered")
			if pDUSessionResourceList != nil {
				for _, pduSessionReourceItem := range pDUSessionResourceList.List {
					pduSessionID := int32(pduSessionReourceItem.PDUSessionID.Value)
					smContext, ok := amfUe.SmContextFindByPDUSessionID(pduSessionID)
					if !ok {
						ranUe.Log.Warnf("SmContext[PDU Session ID:%d] not found", pduSessionID)
						// TODO: Check if doing error handling here
						continue
					}
					rsp, _, _, err := consumer.GetConsumer().SendUpdateSmContextDeactivateUpCnxState(amfUe, smContext, causeAll)
					if err != nil {
						ranUe.Log.Errorf("Send Update SmContextDeactivate UpCnxState Error[%s]", err.Error())
					} else if rsp == nil {
						ranUe.Log.Errorln("Send Update SmContextDeactivate UpCnxState Error")
					}
				}
			}
		} else {
			ranUe.Log.Info("Ue Context in Non GMM-Registered")
			amfUe.SmContextList.Range(func(key, value interface{}) bool {
				smContext := value.(*context.SmContext)
				detail, err := consumer.GetConsumer().SendReleaseSmContextRequest(amfUe, smContext, &causeAll, "", nil)
				if err != nil {
					ranUe.Log.Errorf("Send ReleaseSmContextRequest Error[%s]", err.Error())
				} else if detail != nil {
					ranUe.Log.Errorf("Send ReleaseSmContextRequeste Error[%s]", detail.Cause)
				}
				return true
			})
			ngap_message.SendUEContextReleaseCommand(ranUe, context.UeContextReleaseUeContext, causeGroup, causeValue)
			// TODO: start timer to release RanUe context
			return
		}
	}
	ngap_message.SendUEContextReleaseCommand(ranUe, context.UeContextN2NormalRelease, causeGroup, causeValue)
	// TODO: start timer to release RanUe context
}

func handleUEContextModificationResponseMain(ran *context.AmfRan,
	ranUe *context.RanUe,
	rRCState *ngapType.RRCState,
	userLocationInformation *ngapType.UserLocationInformation,
	criticalityDiagnostics *ngapType.CriticalityDiagnostics,
) {
	if ranUe == nil {
		ran.Log.Error("ranUe is nil")
		return
	}

	if rRCState != nil {
		switch rRCState.Value {
		case ngapType.RRCStatePresentInactive:
			ranUe.Log.Trace("UE RRC State: Inactive")
		case ngapType.RRCStatePresentConnected:
			ranUe.Log.Trace("UE RRC State: Connected")
		}
	}

	if userLocationInformation != nil {
		ranUe.UpdateLocation(userLocationInformation)
	}

	if criticalityDiagnostics != nil {
		printCriticalityDiagnostics(ran, criticalityDiagnostics)
	}
}

func handleUEContextModificationFailureMain(ran *context.AmfRan,
	ranUe *context.RanUe,
	cause *ngapType.Cause,
	criticalityDiagnostics *ngapType.CriticalityDiagnostics,
) {
	if cause != nil {
		printAndGetCause(ran, cause)
	}

	if criticalityDiagnostics != nil {
		printCriticalityDiagnostics(ran, criticalityDiagnostics)
	}
}

func handleRRCInactiveTransitionReportMain(ran *context.AmfRan,
	ranUe *context.RanUe,
	rRCState *ngapType.RRCState,
	userLocationInformation *ngapType.UserLocationInformation,
) {
	if rRCState != nil {
		switch rRCState.Value {
		case ngapType.RRCStatePresentInactive:
			ran.Log.Trace("UE RRC State: Inactive")
		case ngapType.RRCStatePresentConnected:
			ran.Log.Trace("UE RRC State: Connected")
		}
	}
	ranUe.UpdateLocation(userLocationInformation)
}

func handleHandoverNotifyMain(ran *context.AmfRan,
	targetUe *context.RanUe,
	userLocationInformation *ngapType.UserLocationInformation,
) {
	targetUe.Log.Info("Handle Handover notification")

	if userLocationInformation != nil {
		targetUe.UpdateLocation(userLocationInformation)
	}
	amfUe := targetUe.AmfUe
	if amfUe == nil {
		ran.Log.Error("AmfUe is nil")
		return
	}
	sourceUe := targetUe.SourceUe
	if sourceUe == nil {
		// TODO: Send to S-AMF
		// Desciibed in (23.502 4.9.1.3.3) [conditional] 6a.Namf_Communication_N2InfoNotify.
		ran.Log.Error("N2 Handover between AMF has not been implemented yet")
	} else {
		ran.Log.Info("Handle Handover notification Finshed")
		for _, pduSessionID := range targetUe.SuccessPduSessionId {
			smContext, ok := amfUe.SmContextFindByPDUSessionID(pduSessionID)
			if !ok {
				sourceUe.Log.Warnf("SmContext[PDU Session ID:%d] not found", pduSessionID)
				// TODO: Check if doing error handling here
				continue
			}
			_, _, _, err := consumer.GetConsumer().SendUpdateSmContextN2HandoverComplete(amfUe, smContext, "", nil)
			if err != nil {
				ran.Log.Errorf("Send UpdateSmContextN2HandoverComplete Error[%s]", err.Error())
			}
		}

		gmm_common.AttachRanUeToAmfUeAndReleaseOldHandover(amfUe, sourceUe, targetUe)
	}

	// TODO: The UE initiates Mobility Registration Update procedure as described in clause 4.2.2.2.2.
}

// TS 23.502 4.9.1
func handlePathSwitchRequestMain(ran *context.AmfRan,
	rANUENGAPID *ngapType.RANUENGAPID,
	sourceAMFUENGAPID *ngapType.AMFUENGAPID,
	userLocationInformation *ngapType.UserLocationInformation,
	uESecurityCapabilities *ngapType.UESecurityCapabilities,
	pduSessionResourceToBeSwitchedInDLList *ngapType.PDUSessionResourceToBeSwitchedDLList,
	pduSessionResourceFailedToSetupList *ngapType.PDUSessionResourceFailedToSetupListPSReq,
) {
	ranUe := context.GetSelf().RanUeFindByAmfUeNgapID(sourceAMFUENGAPID.Value)
	if ranUe == nil {
		ran.Log.Errorf("Cannot find UE from sourceAMfUeNgapID[%d]", sourceAMFUENGAPID.Value)
		ngap_message.SendPathSwitchRequestFailure(ran, sourceAMFUENGAPID.Value, rANUENGAPID.Value, nil, nil)
		return
	}

	ran.Log.Tracef("AmfUeNgapID[%d] RanUeNgapID[%d]", ranUe.AmfUeNgapId, ranUe.RanUeNgapId)

	amfUe := ranUe.AmfUe
	if amfUe == nil {
		ranUe.Log.Error("AmfUe is nil")
		ngap_message.SendPathSwitchRequestFailure(ran, sourceAMFUENGAPID.Value, rANUENGAPID.Value, nil, nil)
		return
	}

	if amfUe.SecurityContextIsValid() {
		// Update NH
		amfUe.UpdateNH()
	} else {
		ranUe.Log.Errorf("No Security Context : SUPI[%s]", amfUe.Supi)
		ngap_message.SendPathSwitchRequestFailure(ran, sourceAMFUENGAPID.Value, rANUENGAPID.Value, nil, nil)
		return
	}

	if uESecurityCapabilities != nil {
		amfUe.UESecurityCapability.SetEA1_128_5G(uESecurityCapabilities.NRencryptionAlgorithms.Value.Bytes[0] & 0x80)
		amfUe.UESecurityCapability.SetEA2_128_5G(uESecurityCapabilities.NRencryptionAlgorithms.Value.Bytes[0] & 0x40)
		amfUe.UESecurityCapability.SetEA3_128_5G(uESecurityCapabilities.NRencryptionAlgorithms.Value.Bytes[0] & 0x20)
		amfUe.UESecurityCapability.SetIA1_128_5G(uESecurityCapabilities.NRintegrityProtectionAlgorithms.Value.Bytes[0] & 0x80)
		amfUe.UESecurityCapability.SetIA2_128_5G(uESecurityCapabilities.NRintegrityProtectionAlgorithms.Value.Bytes[0] & 0x40)
		amfUe.UESecurityCapability.SetIA3_128_5G(uESecurityCapabilities.NRintegrityProtectionAlgorithms.Value.Bytes[0] & 0x20)
		// not support any E-UTRA algorithms
	}

	if rANUENGAPID != nil {
		ranUe.RanUeNgapId = rANUENGAPID.Value
	}

	ranUe.UpdateLocation(userLocationInformation)

	var pduSessionResourceSwitchedList ngapType.PDUSessionResourceSwitchedList
	var pduSessionResourceReleasedListPSAck ngapType.PDUSessionResourceReleasedListPSAck
	var pduSessionResourceReleasedListPSFail ngapType.PDUSessionResourceReleasedListPSFail

	if pduSessionResourceToBeSwitchedInDLList != nil {
		ranUe.Log.Infof("Send PathSwitchRequestTransfer to SMF")
		for _, item := range pduSessionResourceToBeSwitchedInDLList.List {
			pduSessionID := int32(item.PDUSessionID.Value)
			transfer := item.PathSwitchRequestTransfer
			smContext, ok := amfUe.SmContextFindByPDUSessionID(pduSessionID)
			if !ok {
				ranUe.Log.Warnf("SmContext[PDU Session ID:%d] not found", pduSessionID)
				// TODO: Check if doing error handling here
				continue
			}
			response, errResponse, _, err := consumer.GetConsumer().SendUpdateSmContextXnHandover(amfUe, smContext,
				models.N2SmInfoType_PATH_SWITCH_REQ, transfer)
			if err != nil {
				ranUe.Log.Errorf("SendUpdateSmContextXnHandover[PathSwitchRequestTransfer] Error:\n%s", err.Error())
			}
			if response != nil && response.BinaryDataN2SmInformation != nil {
				pduSessionResourceSwitchedItem := ngapType.PDUSessionResourceSwitchedItem{}
				pduSessionResourceSwitchedItem.PDUSessionID.Value = int64(pduSessionID)
				pduSessionResourceSwitchedItem.PathSwitchRequestAcknowledgeTransfer = response.BinaryDataN2SmInformation
				pduSessionResourceSwitchedList.List = append(pduSessionResourceSwitchedList.List, pduSessionResourceSwitchedItem)
			}
			if errResponse != nil && errResponse.BinaryDataN2SmInformation != nil {
				pduSessionResourceReleasedItem := ngapType.PDUSessionResourceReleasedItemPSFail{}
				pduSessionResourceReleasedItem.PDUSessionID.Value = int64(pduSessionID)
				pduSessionResourceReleasedItem.PathSwitchRequestUnsuccessfulTransfer = errResponse.BinaryDataN2SmInformation
				pduSessionResourceReleasedListPSFail.List = append(pduSessionResourceReleasedListPSFail.List,
					pduSessionResourceReleasedItem)
			}
		}
	}

	if pduSessionResourceFailedToSetupList != nil {
		ranUe.Log.Infof("Send PathSwitchRequestSetupFailedTransfer to SMF")
		for _, item := range pduSessionResourceFailedToSetupList.List {
			pduSessionID := int32(item.PDUSessionID.Value)
			transfer := item.PathSwitchRequestSetupFailedTransfer
			smContext, ok := amfUe.SmContextFindByPDUSessionID(pduSessionID)
			if !ok {
				ranUe.Log.Warnf("SmContext[PDU Session ID:%d] not found", pduSessionID)
				// TODO: Check if doing error handling here
				continue
			}
			response, errResponse, _, err := consumer.GetConsumer().SendUpdateSmContextXnHandoverFailed(amfUe, smContext,
				models.N2SmInfoType_PATH_SWITCH_SETUP_FAIL, transfer)
			if err != nil {
				ranUe.Log.Errorf("SendUpdateSmContextXnHandoverFailed[PathSwitchRequestSetupFailedTransfer] Error: %+v", err)
			}
			if response != nil && response.BinaryDataN2SmInformation != nil {
				pduSessionResourceReleasedItem := ngapType.PDUSessionResourceReleasedItemPSAck{}
				pduSessionResourceReleasedItem.PDUSessionID.Value = int64(pduSessionID)
				pduSessionResourceReleasedItem.PathSwitchRequestUnsuccessfulTransfer = response.BinaryDataN2SmInformation
				pduSessionResourceReleasedListPSAck.List = append(pduSessionResourceReleasedListPSAck.List,
					pduSessionResourceReleasedItem)
			}
			if errResponse != nil && errResponse.BinaryDataN2SmInformation != nil {
				pduSessionResourceReleasedItem := ngapType.PDUSessionResourceReleasedItemPSFail{}
				pduSessionResourceReleasedItem.PDUSessionID.Value = int64(pduSessionID)
				pduSessionResourceReleasedItem.PathSwitchRequestUnsuccessfulTransfer = errResponse.BinaryDataN2SmInformation
				pduSessionResourceReleasedListPSFail.List = append(pduSessionResourceReleasedListPSFail.List,
					pduSessionResourceReleasedItem)
			}
		}
	}

	// TS 23.502 4.9.1.2.2 step 7: send ack to Target NG-RAN. If none of the requested PDU Sessions have been switched
	// successfully, the AMF shall send an N2 Path Switch Request Failure message to the Target NG-RAN
	if len(pduSessionResourceSwitchedList.List) > 0 {
		// TODO: set newSecurityContextIndicator to true if there is a new security context
		err := ranUe.SwitchToRan(ran, rANUENGAPID.Value)
		if err != nil {
			ranUe.Log.Error(err.Error())
			return
		}
		ngap_message.SendPathSwitchRequestAcknowledge(ranUe, pduSessionResourceSwitchedList,
			pduSessionResourceReleasedListPSAck, false, nil, nil, nil)
	} else if len(pduSessionResourceReleasedListPSFail.List) > 0 {
		ngap_message.SendPathSwitchRequestFailure(ran, sourceAMFUENGAPID.Value, rANUENGAPID.Value,
			&pduSessionResourceReleasedListPSFail, nil)
	} else {
		ngap_message.SendPathSwitchRequestFailure(ran, sourceAMFUENGAPID.Value, rANUENGAPID.Value, nil, nil)
	}
}

func handleHandoverRequestAcknowledgeMain(ran *context.AmfRan,
	targetUe *context.RanUe,
	rANUENGAPID *ngapType.RANUENGAPID,
	pDUSessionResourceAdmittedList *ngapType.PDUSessionResourceAdmittedList,
	pDUSessionResourceFailedToSetupListHOAck *ngapType.PDUSessionResourceFailedToSetupListHOAck,
	targetToSourceTransparentContainer *ngapType.TargetToSourceTransparentContainer,
	criticalityDiagnostics *ngapType.CriticalityDiagnostics,
) {
	if criticalityDiagnostics != nil {
		printCriticalityDiagnostics(ran, criticalityDiagnostics)
	}

	if targetUe == nil {
		ran.Log.Errorf("Target Ue is missing")
		return
	}

	if rANUENGAPID != nil {
		targetUe.RanUeNgapId = rANUENGAPID.Value
		ran.RanUeList.Store(targetUe.RanUeNgapId, targetUe)
	}
	ran.Log.Debugf("Target Ue RanUeNgapID[%d] AmfUeNgapID[%d]", targetUe.RanUeNgapId, targetUe.AmfUeNgapId)

	amfUe := targetUe.AmfUe
	if amfUe == nil {
		targetUe.Log.Error("amfUe is nil")
		return
	}

	var pduSessionResourceHandoverList ngapType.PDUSessionResourceHandoverList
	var pduSessionResourceToReleaseList ngapType.PDUSessionResourceToReleaseListHOCmd

	// describe in 23.502 4.9.1.3.2 step11
	if pDUSessionResourceAdmittedList != nil {
		targetUe.Log.Infof("Send HandoverRequestAcknowledgeTransfer to SMF")
		for _, item := range pDUSessionResourceAdmittedList.List {
			pduSessionID := int32(item.PDUSessionID.Value)
			transfer := item.HandoverRequestAcknowledgeTransfer
			smContext, ok := amfUe.SmContextFindByPDUSessionID(pduSessionID)
			if !ok {
				targetUe.Log.Warnf("SmContext[PDU Session ID:%d] not found", pduSessionID)
				// TODO: Check if doing error handling here
				continue
			}
			resp, errResponse, problemDetails, err := consumer.GetConsumer().SendUpdateSmContextN2HandoverPrepared(amfUe,
				smContext, models.N2SmInfoType_HANDOVER_REQ_ACK, transfer)
			if err != nil {
				targetUe.Log.Errorf("Send HandoverRequestAcknowledgeTransfer error: %v", err)
			}
			if problemDetails != nil {
				targetUe.Log.Warnf("ProblemDetails[status: %d, Cause: %s]", problemDetails.Status, problemDetails.Cause)
			}
			if resp != nil && resp.BinaryDataN2SmInformation != nil {
				handoverItem := ngapType.PDUSessionResourceHandoverItem{}
				handoverItem.PDUSessionID = item.PDUSessionID
				handoverItem.HandoverCommandTransfer = resp.BinaryDataN2SmInformation
				pduSessionResourceHandoverList.List = append(pduSessionResourceHandoverList.List, handoverItem)
				targetUe.SuccessPduSessionId = append(targetUe.SuccessPduSessionId, pduSessionID)
			}
			if errResponse != nil && errResponse.BinaryDataN2SmInformation != nil {
				releaseItem := ngapType.PDUSessionResourceToReleaseItemHOCmd{}
				releaseItem.PDUSessionID = item.PDUSessionID
				releaseItem.HandoverPreparationUnsuccessfulTransfer = errResponse.BinaryDataN2SmInformation
				pduSessionResourceToReleaseList.List = append(pduSessionResourceToReleaseList.List, releaseItem)
			}
		}
	}

	if pDUSessionResourceFailedToSetupListHOAck != nil {
		targetUe.Log.Infof("Send HandoverResourceAllocationUnsuccessfulTransfer to SMF")
		for _, item := range pDUSessionResourceFailedToSetupListHOAck.List {
			pduSessionID := int32(item.PDUSessionID.Value)
			transfer := item.HandoverResourceAllocationUnsuccessfulTransfer
			smContext, ok := amfUe.SmContextFindByPDUSessionID(pduSessionID)
			if !ok {
				targetUe.Log.Warnf("SmContext[PDU Session ID:%d] not found", pduSessionID)
				// TODO: Check if doing error handling here
				continue
			}
			_, _, problemDetails, err := consumer.GetConsumer().SendUpdateSmContextN2HandoverPrepared(amfUe, smContext,
				models.N2SmInfoType_HANDOVER_RES_ALLOC_FAIL, transfer)
			if err != nil {
				targetUe.Log.Errorf("Send HandoverResourceAllocationUnsuccessfulTransfer error: %v", err)
			}
			if problemDetails != nil {
				targetUe.Log.Warnf("ProblemDetails[status: %d, Cause: %s]", problemDetails.Status, problemDetails.Cause)
			}
		}
	}

	sourceUe := targetUe.SourceUe
	if sourceUe == nil {
		// TODO: Send Namf_Communication_CreateUEContext Response to S-AMF
		ran.Log.Error("handover between different Ue has not been implement yet")
	} else {
		ran.Log.Tracef("Source: RanUeNgapID[%d] AmfUeNgapID[%d]", sourceUe.RanUeNgapId, sourceUe.AmfUeNgapId)
		ran.Log.Tracef("Target: RanUeNgapID[%d] AmfUeNgapID[%d]", targetUe.RanUeNgapId, targetUe.AmfUeNgapId)
		if len(pduSessionResourceHandoverList.List) == 0 {
			targetUe.Log.Info("Handle Handover Preparation Failure [HoFailure In Target5GC NgranNode Or TargetSystem]")
			cause := &ngapType.Cause{
				Present: ngapType.CausePresentRadioNetwork,
				RadioNetwork: &ngapType.CauseRadioNetwork{
					Value: ngapType.CauseRadioNetworkPresentHoFailureInTarget5GCNgranNodeOrTargetSystem,
				},
			}
			ngap_message.SendHandoverPreparationFailure(sourceUe, *cause, nil)
			return
		}
		ngap_message.SendHandoverCommand(sourceUe, pduSessionResourceHandoverList, pduSessionResourceToReleaseList,
			*targetToSourceTransparentContainer, nil)
	}
}

func handleHandoverFailureMain(ran *context.AmfRan,
	targetUe *context.RanUe,
	cause *ngapType.Cause,
	criticalityDiagnostics *ngapType.CriticalityDiagnostics,
) {
	causePresent := ngapType.CausePresentRadioNetwork
	causeValue := ngapType.CauseRadioNetworkPresentHoFailureInTarget5GCNgranNodeOrTargetSystem
	if cause != nil {
		causePresent, causeValue = printAndGetCause(ran, cause)
	}

	if criticalityDiagnostics != nil {
		printCriticalityDiagnostics(ran, criticalityDiagnostics)
	}

	if targetUe == nil {
		ran.Log.Errorf("Target Ue is missing")
		return
	}

	targetUe.Log.Info("Handle Handover Failure")

	sourceUe := targetUe.SourceUe
	if sourceUe == nil {
		// TODO: handle N2 Handover between AMF
		ran.Log.Error("N2 Handover between AMF has not been implemented yet")
	} else {
		amfUe := targetUe.AmfUe
		if amfUe != nil {
			amfUe.SmContextList.Range(func(key, value interface{}) bool {
				pduSessionID := key.(int32)
				smContext := value.(*context.SmContext)
				causeAll := context.CauseAll{
					NgapCause: &models.NgApCause{
						Group: int32(causePresent),
						Value: int32(causeValue),
					},
				}
				_, _, _, err := consumer.GetConsumer().SendUpdateSmContextN2HandoverCanceled(amfUe, smContext, causeAll)
				if err != nil {
					ran.Log.Errorf("Send UpdateSmContextN2HandoverCanceled Error for pduSessionID[%d]", pduSessionID)
				}
				return true
			})
		}
		sendCause := cause
		if sendCause == nil {
			sendCause = &ngapType.Cause{
				Present: ngapType.CausePresentRadioNetwork,
				RadioNetwork: &ngapType.CauseRadioNetwork{
					Value: ngapType.CauseRadioNetworkPresentHoFailureInTarget5GCNgranNodeOrTargetSystem,
				},
			}
		}
		ngap_message.SendHandoverPreparationFailure(sourceUe, *sendCause, criticalityDiagnostics)
	}

	ngap_message.SendUEContextReleaseCommand(targetUe, context.UeContextReleaseHandover, causePresent, causeValue)
}

func handleHandoverRequiredMain(ran *context.AmfRan,
	sourceUe *context.RanUe,
	handoverType *ngapType.HandoverType,
	cause *ngapType.Cause,
	targetID *ngapType.TargetID,
	pDUSessionResourceListHORqd *ngapType.PDUSessionResourceListHORqd,
	sourceToTargetTransparentContainer *ngapType.SourceToTargetTransparentContainer,
) {
	amfUe := sourceUe.AmfUe
	if amfUe == nil {
		ran.Log.Error("Cannot find amfUE from sourceUE")
		return
	}

	if targetID.Present != ngapType.TargetIDPresentTargetRANNodeID {
		ran.Log.Errorf("targetID type[%d] is not supported", targetID.Present)
		return
	}
	amfUe.SetOnGoing(sourceUe.Ran.AnType, &context.OnGoing{
		Procedure: context.OnGoingProcedureN2Handover,
	})
	if !amfUe.SecurityContextIsValid() {
		sourceUe.Log.Info("Handle Handover Preparation Failure [Authentication Failure]")
		cause = &ngapType.Cause{
			Present: ngapType.CausePresentNas,
			Nas: &ngapType.CauseNas{
				Value: ngapType.CauseNasPresentAuthenticationFailure,
			},
		}
		ngap_message.SendHandoverPreparationFailure(sourceUe, *cause, nil)
		return
	}
	aMFSelf := context.GetSelf()
	targetRanNodeId := ngapConvert.RanIdToModels(targetID.TargetRANNodeID.GlobalRANNodeID)
	targetRan, ok := aMFSelf.AmfRanFindByRanID(targetRanNodeId)
	if !ok {
		// handover between different AMF
		sourceUe.Log.Warnf("Handover required : cannot find target Ran Node Id[%+v] in this AMF", targetRanNodeId)
		sourceUe.Log.Error("Handover between different AMF has not been implemented yet")
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

		if pDUSessionResourceListHORqd != nil {
			sourceUe.Log.Infof("Send HandoverRequiredTransfer to SMF")
			for _, pDUSessionResourceHoItem := range pDUSessionResourceListHORqd.List {
				pduSessionID := int32(pDUSessionResourceHoItem.PDUSessionID.Value)
				smContext, okSmContextFindByPDUSessionID := amfUe.SmContextFindByPDUSessionID(pduSessionID)
				if !okSmContextFindByPDUSessionID {
					sourceUe.Log.Warnf("SmContext[PDU Session ID:%d] not found", pduSessionID)
					// TODO: Check if doing error handling here
					continue
				}

				response, _, _, err := consumer.GetConsumer().SendUpdateSmContextN2HandoverPreparing(amfUe, smContext,
					models.N2SmInfoType_HANDOVER_REQUIRED, pDUSessionResourceHoItem.HandoverRequiredTransfer, "", &targetId)
				if err != nil {
					sourceUe.Log.Errorf("consumer.GetConsumer().SendUpdateSmContextN2HandoverPreparing Error: %+v", err)
				}
				if response == nil {
					sourceUe.Log.Errorf("SendUpdateSmContextN2HandoverPreparing Error for pduSessionID[%d]", pduSessionID)
					continue
				} else if response.BinaryDataN2SmInformation != nil {
					ngap_message.AppendPDUSessionResourceSetupListHOReq(&pduSessionReqList, pduSessionID,
						smContext.Snssai(), response.BinaryDataN2SmInformation)
				}
			}
		}
		if len(pduSessionReqList.List) == 0 {
			sourceUe.Log.Info("Handle Handover Preparation Failure [HoFailure In Target5GC NgranNode Or TargetSystem]")
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
		if cause == nil {
			sourceUe.Log.Warnf("Cause is nil")
			cause = &ngapType.Cause{
				Present: ngapType.CausePresentMisc,
				Misc: &ngapType.CauseMisc{
					Value: ngapType.CauseMiscPresentUnspecified,
				},
			}
		}
		ngap_message.SendHandoverRequest(sourceUe, targetRan, *cause, pduSessionReqList,
			*sourceToTargetTransparentContainer, false)
	}
}

func handleHandoverCancelMain(ran *context.AmfRan,
	sourceUe *context.RanUe,
	cause *ngapType.Cause,
) {
	causePresent := ngapType.CausePresentRadioNetwork
	causeValue := ngapType.CauseRadioNetworkPresentHoFailureInTarget5GCNgranNodeOrTargetSystem
	if cause != nil {
		causePresent, causeValue = printAndGetCause(ran, cause)
	}
	targetUe := sourceUe.TargetUe
	if targetUe == nil {
		// Described in (23.502 4.11.1.2.3) step 2
		// Todo : send to T-AMF invoke Namf_UeContextReleaseRequest(targetUe)
		ran.Log.Error("N2 Handover between AMF has not been implemented yet")
	} else {
		ran.Log.Tracef("Target : RAN_UE_NGAP_ID[%d] AMF_UE_NGAP_ID[%d]", targetUe.RanUeNgapId, targetUe.AmfUeNgapId)
		amfUe := sourceUe.AmfUe
		if amfUe != nil {
			amfUe.SmContextList.Range(func(key, value interface{}) bool {
				pduSessionID := key.(int32)
				smContext := value.(*context.SmContext)
				causeAll := context.CauseAll{
					NgapCause: &models.NgApCause{
						Group: int32(causePresent),
						Value: int32(causeValue),
					},
				}
				_, _, _, err := consumer.GetConsumer().SendUpdateSmContextN2HandoverCanceled(amfUe, smContext, causeAll)
				if err != nil {
					sourceUe.Log.Errorf("Send UpdateSmContextN2HandoverCanceled Error for pduSessionID[%d]", pduSessionID)
				}
				return true
			})
		}
		ngap_message.SendUEContextReleaseCommand(targetUe, context.UeContextReleaseHandover, causePresent, causeValue)
		ngap_message.SendHandoverCancelAcknowledge(sourceUe, nil)
	}
}

func handleUplinkRANStatusTransferMain(ran *context.AmfRan,
	ranUe *context.RanUe,
) {
	amfUe := ranUe.AmfUe
	if amfUe == nil {
		ranUe.Log.Error("AmfUe is nil")
		return
	}
	// send to T-AMF using N1N2MessageTransfer (R16)
}

func handleNASNonDeliveryIndicationMain(ran *context.AmfRan,
	ranUe *context.RanUe,
	nASPDU *ngapType.NASPDU,
	cause *ngapType.Cause,
) {
	if cause != nil {
		printAndGetCause(ran, cause)
	}

	if nASPDU != nil {
		amf_nas.HandleNAS(ranUe, ngapType.ProcedureCodeNASNonDeliveryIndication, nASPDU.Value, false)
	}
}

func handleRANConfigurationUpdateMain(ran *context.AmfRan,
	supportedTAList *ngapType.SupportedTAList,
) {
	var cause ngapType.Cause

	if supportedTAList != nil {
		for i := 0; i < len(supportedTAList.List); i++ {
			supportedTAItem := supportedTAList.List[i]
			tac := hex.EncodeToString(supportedTAItem.TAC.Value)
			capOfSupportTai := cap(ran.SupportedTAList)
			for j := 0; j < len(supportedTAItem.BroadcastPLMNList.List); j++ {
				supportedTAI := context.NewSupportedTAI()
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
				ran.Log.Tracef("PLMN_ID[MCC:%s MNC:%s] TAC[%s]", plmnId.Mcc, plmnId.Mnc, tac)
				if len(ran.SupportedTAList) < capOfSupportTai {
					ran.SupportedTAList = append(ran.SupportedTAList, supportedTAI)
				} else {
					break
				}
			}
		}
	}

	if len(ran.SupportedTAList) == 0 {
		ran.Log.Warn("RanConfigurationUpdate failure: No supported TA exist in RanConfigurationUpdate")
		cause.Present = ngapType.CausePresentMisc
		cause.Misc = &ngapType.CauseMisc{
			Value: ngapType.CauseMiscPresentUnspecified,
		}
	} else {
		var found bool
		for i, tai := range ran.SupportedTAList {
			if context.InTaiList(tai.Tai, context.GetSelf().SupportTaiLists) {
				ran.Log.Tracef("SERVED_TAI_INDEX[%d]", i)
				found = true
				break
			}
		}
		if !found {
			ran.Log.Warn("RanConfigurationUpdate failure: Cannot find Served TAI in AMF")
			cause.Present = ngapType.CausePresentMisc
			cause.Misc = &ngapType.CauseMisc{
				Value: ngapType.CauseMiscPresentUnknownPLMN,
			}
		}
	}

	if cause.Present == ngapType.CausePresentNothing {
		ran.Log.Info("Handle RanConfigurationUpdateAcknowledge")
		ngap_message.SendRanConfigurationUpdateAcknowledge(ran, nil)
	} else {
		ran.Log.Info("Handle RanConfigurationUpdateAcknowledgeFailure")
		ngap_message.SendRanConfigurationUpdateFailure(ran, cause, nil)
	}
}

func handleUplinkRANConfigurationTransferMain(ran *context.AmfRan,
	sONConfigurationTransferUL *ngapType.SONConfigurationTransfer,
) {
	if sONConfigurationTransferUL != nil {
		targetRanNodeID := ngapConvert.RanIdToModels(sONConfigurationTransferUL.TargetRANNodeID.GlobalRANNodeID)

		if targetRanNodeID.GNbId != nil && targetRanNodeID.GNbId.GNBValue != "" {
			ran.Log.Tracef("targerRanID [%s]", targetRanNodeID.GNbId.GNBValue)
		}

		aMFSelf := context.GetSelf()

		targetRan, ok := aMFSelf.AmfRanFindByRanID(targetRanNodeID)
		if !ok {
			ran.Log.Warn("targetRan is nil")
		}

		ngap_message.SendDownlinkRanConfigurationTransfer(targetRan, sONConfigurationTransferUL)
	}
}

func handleUplinkUEAssociatedNRPPaTransportMain(ran *context.AmfRan,
	ranUe *context.RanUe,
	routingID *ngapType.RoutingID,
) {
	ranUe.RoutingID = hex.EncodeToString(routingID.Value)

	// TODO: Forward NRPPaPDU to LMF
}

func handleUplinkNonUEAssociatedNRPPaTransportMain(ran *context.AmfRan,
	routingID *ngapType.RoutingID,
	nRPPaPDU *ngapType.NRPPaPDU,
) {
	// Forward routingID to LMF
	// Described in (23.502 4.13.5.6)

	// TODO: Forward NRPPaPDU to LMF
}

func handleLocationReportMain(ran *context.AmfRan,
	ranUe *context.RanUe,
	userLocationInformation *ngapType.UserLocationInformation,
	uEPresenceInAreaOfInterestList *ngapType.UEPresenceInAreaOfInterestList,
	locationReportingRequestType *ngapType.LocationReportingRequestType,
) {
	ranUe.UpdateLocation(userLocationInformation)

	if locationReportingRequestType != nil {
		ranUe.Log.Tracef("Report Area[%d]", locationReportingRequestType.ReportArea.Value)

		switch locationReportingRequestType.EventType.Value {
		case ngapType.EventTypePresentDirect:
			ranUe.Log.Trace("To report directly")

		case ngapType.EventTypePresentChangeOfServeCell:
			ranUe.Log.Trace("To report upon change of serving cell")

		case ngapType.EventTypePresentUePresenceInAreaOfInterest:
			ranUe.Log.Trace("To report UE presence in the area of interest")
			if uEPresenceInAreaOfInterestList != nil {
				for _, uEPresenceInAreaOfInterestItem := range uEPresenceInAreaOfInterestList.List {
					uEPresence := uEPresenceInAreaOfInterestItem.UEPresence.Value
					referenceID := uEPresenceInAreaOfInterestItem.LocationReportingReferenceID.Value

					for _, AOIitem := range locationReportingRequestType.AreaOfInterestList.List {
						if referenceID == AOIitem.LocationReportingReferenceID.Value {
							ran.Log.Tracef("uEPresence[%d], presence AOI ReferenceID[%d]", uEPresence, referenceID)
						}
					}
				}
			}

		case ngapType.EventTypePresentStopChangeOfServeCell:
			ranUe.Log.Trace("To stop reporting at change of serving cell")
			ngap_message.SendLocationReportingControl(ranUe, nil, 0, locationReportingRequestType.EventType)
			// TODO: Clear location report

		case ngapType.EventTypePresentStopUePresenceInAreaOfInterest:
			ranUe.Log.Trace("To stop reporting UE presence in the area of interest")
			ranUe.Log.Tracef("ReferenceID To Be Canceled[%d]",
				locationReportingRequestType.LocationReportingReferenceIDToBeCancelled.Value)
			// TODO: Clear location report

		case ngapType.EventTypePresentCancelLocationReportingForTheUe:
			ranUe.Log.Trace("To cancel location reporting for the UE")
			// TODO: Clear location report
		}
	}
}

func handleUERadioCapabilityInfoIndicationMain(ran *context.AmfRan,
	ranUe *context.RanUe,
	uERadioCapability *ngapType.UERadioCapability,
	uERadioCapabilityForPaging *ngapType.UERadioCapabilityForPaging,
) {
	amfUe := ranUe.AmfUe

	if amfUe == nil {
		ranUe.Log.Errorln("amfUe is nil")
		return
	}
	if uERadioCapability != nil {
		amfUe.UeRadioCapability = hex.EncodeToString(uERadioCapability.Value)
	}
	if uERadioCapabilityForPaging != nil {
		amfUe.UeRadioCapabilityForPaging = &context.UERadioCapabilityForPaging{}
		if uERadioCapabilityForPaging.UERadioCapabilityForPagingOfNR != nil {
			amfUe.UeRadioCapabilityForPaging.NR = hex.EncodeToString(
				uERadioCapabilityForPaging.UERadioCapabilityForPagingOfNR.Value)
		}
		if uERadioCapabilityForPaging.UERadioCapabilityForPagingOfEUTRA != nil {
			amfUe.UeRadioCapabilityForPaging.EUTRA = hex.EncodeToString(
				uERadioCapabilityForPaging.UERadioCapabilityForPagingOfEUTRA.Value)
		}
	}

	// TS 38.413 8.14.1.2/TS 23.502 4.2.8a step5/TS 23.501, clause 5.4.4.1.
	// send its most up to date UE Radio Capability information to the RAN in the N2 REQUEST message.
}

func handleAMFConfigurationUpdateFailureMain(ran *context.AmfRan,
	cause *ngapType.Cause,
	criticalityDiagnostics *ngapType.CriticalityDiagnostics,
) {
	if cause != nil {
		printAndGetCause(ran, cause)
	}

	//	TODO: Time To Wait

	if criticalityDiagnostics != nil {
		printCriticalityDiagnostics(ran, criticalityDiagnostics)
	}
}

func handleAMFConfigurationUpdateAcknowledgeMain(ran *context.AmfRan,
	criticalityDiagnostics *ngapType.CriticalityDiagnostics,
) {
	if criticalityDiagnostics != nil {
		printCriticalityDiagnostics(ran, criticalityDiagnostics)
	}
}

func handleErrorIndicationMain(ran *context.AmfRan,
	aMFUENGAPID *ngapType.AMFUENGAPID,
	rANUENGAPID *ngapType.RANUENGAPID,
	cause *ngapType.Cause,
	criticalityDiagnostics *ngapType.CriticalityDiagnostics,
) {
	ran.Log.Infof("Handle Error Indication: RAN_UE_NGAP_ID:%v AMF_UE_NGAP_ID:%v", rANUENGAPID, aMFUENGAPID)

	if cause == nil && criticalityDiagnostics == nil {
		ran.Log.Error("[ErrorIndication] both Cause IE and CriticalityDiagnostics IE are nil, should have at least one")
		return
	}

	if cause != nil {
		printAndGetCause(ran, cause)
	}

	if criticalityDiagnostics != nil {
		printCriticalityDiagnostics(ran, criticalityDiagnostics)
	}

	// TODO: handle error based on cause/criticalityDiagnostics

	if cause != nil &&
		cause.Present == ngapType.CausePresentRadioNetwork &&
		(cause.RadioNetwork.Value == ngapType.CauseRadioNetworkPresentUnknownLocalUENGAPID ||
			cause.RadioNetwork.Value == ngapType.CauseRadioNetworkPresentInconsistentRemoteUENGAPID) {
		// Implement invalid AP ID behavior in TS 38.413
		// These is in "10.6 Handling of AP ID" in TS 38.413
		//  > if this message is not the last message for this UE-associated logical connection, the node
		//  > shall initiate an Error Indication procedure with inclusion of the received AP ID(s) from the
		//  > peer node and an appropriate cause value. Both nodes shall initiate a local release of any
		//  > established UE-associated logical connection (for the same NG interface) having the erroneous
		//  > AP ID as either the local or remote identifier.
		// So we think that these Cause codes that represent incorrect AP ID(s) need to trigger local release.
		if aMFUENGAPID != nil {
			ranUe := context.GetSelf().RanUeFindByAmfUeNgapID(aMFUENGAPID.Value)
			if ranUe != nil && ranUe.Ran == ran {
				removeRanUeByInvalidId(ran, ranUe, fmt.Sprintf("ErrorIndication (AmfUeNgapID: %d)", aMFUENGAPID.Value))
			}
		}
		if rANUENGAPID != nil {
			ranUe := ran.RanUeFindByRanUeNgapID(rANUENGAPID.Value)
			removeRanUeByInvalidId(ran, ranUe, fmt.Sprintf("ErrorIndication (RanUeNgapID: %d)", rANUENGAPID.Value))
		}
	}
}

func handleCellTrafficTraceMain(ran *context.AmfRan,
	ranUe *context.RanUe,
	nGRANTraceID *ngapType.NGRANTraceID,
	nGRANCGI *ngapType.NGRANCGI,
	traceCollectionEntityIPAddress *ngapType.TransportLayerAddress,
) {
	if nGRANTraceID != nil {
		ranUe.Trsr = hex.EncodeToString(nGRANTraceID.Value[6:])

		ranUe.Log.Tracef("TRSR[%s]", ranUe.Trsr)
	}

	if nGRANCGI != nil {
		switch nGRANCGI.Present {
		case ngapType.NGRANCGIPresentNRCGI:
			plmnID := ngapConvert.PlmnIdToModels(nGRANCGI.NRCGI.PLMNIdentity)
			cellID := ngapConvert.BitStringToHex(&nGRANCGI.NRCGI.NRCellIdentity.Value)
			ranUe.Log.Debugf("NRCGI[plmn: %s, cellID: %s]", plmnID, cellID)
		case ngapType.NGRANCGIPresentEUTRACGI:
			plmnID := ngapConvert.PlmnIdToModels(nGRANCGI.EUTRACGI.PLMNIdentity)
			cellID := ngapConvert.BitStringToHex(&nGRANCGI.EUTRACGI.EUTRACellIdentity.Value)
			ranUe.Log.Debugf("EUTRACGI[plmn: %s, cellID: %s]", plmnID, cellID)
		}
	}

	if traceCollectionEntityIPAddress != nil {
		tceIpv4, tceIpv6 := ngapConvert.IPAddressToString(*traceCollectionEntityIPAddress)
		if tceIpv4 != "" {
			ranUe.Log.Debugf("TCE IP Address[v4: %s]", tceIpv4)
		}
		if tceIpv6 != "" {
			ranUe.Log.Debugf("TCE IP Address[v6: %s]", tceIpv6)
		}
	}

	// TODO: TS 32.422 4.2.2.10
	// When AMF receives this new NG signaling message containing the Trace Recording Session Reference (TRSR)
	// and Trace Reference (TR), the AMF shall look up the SUPI/IMEI(SV) of the given call from its database and
	// shall send the SUPI/IMEI(SV) numbers together with the Trace Recording Session Reference and Trace Reference
	// to the Trace Collection Entity.
}

func printAndGetCause(ran *context.AmfRan, cause *ngapType.Cause) (present int, value aper.Enumerated) {
	present = cause.Present
	switch cause.Present {
	case ngapType.CausePresentRadioNetwork:
		ran.Log.Warnf("Cause RadioNetwork[%d]", cause.RadioNetwork.Value)
		value = cause.RadioNetwork.Value
	case ngapType.CausePresentTransport:
		ran.Log.Warnf("Cause Transport[%d]", cause.Transport.Value)
		value = cause.Transport.Value
	case ngapType.CausePresentProtocol:
		ran.Log.Warnf("Cause Protocol[%d]", cause.Protocol.Value)
		value = cause.Protocol.Value
	case ngapType.CausePresentNas:
		ran.Log.Warnf("Cause Nas[%d]", cause.Nas.Value)
		value = cause.Nas.Value
	case ngapType.CausePresentMisc:
		ran.Log.Warnf("Cause Misc[%d]", cause.Misc.Value)
		value = cause.Misc.Value
	default:
		ran.Log.Errorf("Invalid Cause group[%d]", cause.Present)
	}
	return
}

func printCriticalityDiagnostics(ran *context.AmfRan, criticalityDiagnostics *ngapType.CriticalityDiagnostics) {
	ran.Log.Trace("Criticality Diagnostics")

	if criticalityDiagnostics.ProcedureCriticality != nil {
		switch criticalityDiagnostics.ProcedureCriticality.Value {
		case ngapType.CriticalityPresentReject:
			ran.Log.Trace("Procedure Criticality: Reject")
		case ngapType.CriticalityPresentIgnore:
			ran.Log.Trace("Procedure Criticality: Ignore")
		case ngapType.CriticalityPresentNotify:
			ran.Log.Trace("Procedure Criticality: Notify")
		}
	}

	if criticalityDiagnostics.IEsCriticalityDiagnostics != nil {
		for _, ieCriticalityDiagnostics := range criticalityDiagnostics.IEsCriticalityDiagnostics.List {
			ran.Log.Tracef("IE ID: %d", ieCriticalityDiagnostics.IEID.Value)

			switch ieCriticalityDiagnostics.IECriticality.Value {
			case ngapType.CriticalityPresentReject:
				ran.Log.Trace("Criticality Reject")
			case ngapType.CriticalityPresentNotify:
				ran.Log.Trace("Criticality Notify")
			}

			switch ieCriticalityDiagnostics.TypeOfError.Value {
			case ngapType.TypeOfErrorPresentNotUnderstood:
				ran.Log.Trace("Type of error: Not understood")
			case ngapType.TypeOfErrorPresentMissing:
				ran.Log.Trace("Type of error: Missing")
			}
		}
	}
}

func buildCriticalityDiagnostics(
	procedureCode *int64,
	triggeringMessage *aper.Enumerated,
	procedureCriticality *aper.Enumerated,
	iesCriticalityDiagnostics *ngapType.CriticalityDiagnosticsIEList) (
	criticalityDiagnostics ngapType.CriticalityDiagnostics,
) {
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

func buildCriticalityDiagnosticsIEItem(ieCriticality aper.Enumerated, ieID int64, typeOfErr aper.Enumerated) (
	item ngapType.CriticalityDiagnosticsIEItem,
) {
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

func isLatestAmfUe(amfUe *context.AmfUe) bool {
	if latestAmfUe, ok := context.GetSelf().AmfUeFindByUeContextID(amfUe.Supi); ok {
		if amfUe == latestAmfUe {
			return true
		}
	}
	return false
}

func removeRanUeByInvalidId(ran *context.AmfRan, ranUe *context.RanUe, reason string) {
	if ranUe == nil {
		return
	}

	ranUe.Log.Errorf("Remove RanUe by %s", reason)
	amfUe := ranUe.AmfUe
	if amfUe != nil && amfUe.RanUe[ran.AnType] == ranUe {
		if amfUe.T3550 != nil {
			amfUe.State[ranUe.Ran.AnType].Set(context.Registered)
		}
		gmm_common.StopAll5GSMMTimers(amfUe)
		amfUe.DetachRanUe(ran.AnType)
	}
	ranUe.DetachAmfUe()
	if err := ranUe.Remove(); err != nil {
		ran.Log.Errorf("Remove ranUe error: %s", err)
	}
}

// Implementation "10.6 Handling of AP ID" in TS 38.413
// This function is implementation of these sentences.
//
//	> a local release of any established UE-associated logical connection (for the same NG interface)
//	> having the erroneous AP ID as either the local or remote identifier.
func removeRanUeByInvalidUE(ran *context.AmfRan, aMFUENGAPID *ngapType.AMFUENGAPID, rANUENGAPID *ngapType.RANUENGAPID) {
	if aMFUENGAPID != nil {
		ranUe := context.GetSelf().RanUeFindByAmfUeNgapID(aMFUENGAPID.Value)
		if ranUe != nil && ranUe.Ran == ran {
			removeRanUeByInvalidId(ran, ranUe, fmt.Sprintf("Invalid UE ID (AmfUeNgapID: %d)", aMFUENGAPID.Value))
		}
	}
	if rANUENGAPID != nil {
		ranUe := ran.RanUeFindByRanUeNgapID(rANUENGAPID.Value)
		removeRanUeByInvalidId(ran, ranUe, fmt.Sprintf("Invalid UE ID (RanUeNgapID: %d)", rANUENGAPID.Value))
	}
}

// Implementation "10.6 Handling of AP ID" in TS 38.413
//
// firstReturnedMessage: This argument is true in case of the search for “first returned message”.
// sendErrorIndication: If this is true, this function sends the ErrorIndication message for
// invalid AP ID. In case of “last message”, this argument muse be set false.
func ranUeFind(ran *context.AmfRan,
	aMFUENGAPID *ngapType.AMFUENGAPID, rANUENGAPID *ngapType.RANUENGAPID,
	firstReturnedMessage bool, sendErrorIndication bool,
) (ranUe *context.RanUe, err error) {
	if ran == nil {
		return nil, fmt.Errorf("ran is nil")
	}
	if aMFUENGAPID == nil {
		return nil, fmt.Errorf("AmfUeNgapID is nil")
	}
	var rANUENGAPID_string string
	if rANUENGAPID == nil {
		rANUENGAPID_string = "none"
	} else {
		rANUENGAPID_string = fmt.Sprintf("%d", rANUENGAPID.Value)
	}

	ranUe = context.GetSelf().RanUeFindByAmfUeNgapID(aMFUENGAPID.Value)
	if ranUe == nil {
		cause := &ngapType.Cause{
			Present: ngapType.CausePresentRadioNetwork,
			RadioNetwork: &ngapType.CauseRadioNetwork{
				Value: ngapType.CauseRadioNetworkPresentUnknownLocalUENGAPID,
			},
		}
		if sendErrorIndication {
			ngap_message.SendErrorIndication(ran, aMFUENGAPID, rANUENGAPID, cause, nil)
		}
		removeRanUeByInvalidUE(ran, aMFUENGAPID, rANUENGAPID)
		return nil, fmt.Errorf("no RanUe Context[AmfUeNgapID: %d, RanUeNgapID: %s]",
			aMFUENGAPID.Value, rANUENGAPID_string)
	}
	if ranUe.Ran != ran {
		cause := &ngapType.Cause{
			Present: ngapType.CausePresentRadioNetwork,
			RadioNetwork: &ngapType.CauseRadioNetwork{
				Value: ngapType.CauseRadioNetworkPresentUnknownLocalUENGAPID,
			},
		}
		if sendErrorIndication {
			ngap_message.SendErrorIndication(ran, aMFUENGAPID, rANUENGAPID, cause, nil)
		}
		removeRanUeByInvalidUE(ran, aMFUENGAPID, rANUENGAPID)
		return nil, fmt.Errorf("RanUe Context is not in Ran[AmfUeNgapID: %d, RanUeNgapID: %s]",
			aMFUENGAPID.Value, rANUENGAPID_string)
	}

	if rANUENGAPID == nil || firstReturnedMessage {
		if ranUe.RanUeNgapId != context.RanUeNgapIdUnspecified {
			cause := &ngapType.Cause{
				Present: ngapType.CausePresentRadioNetwork,
				RadioNetwork: &ngapType.CauseRadioNetwork{
					Value: ngapType.CauseRadioNetworkPresentInconsistentRemoteUENGAPID,
				},
			}
			if sendErrorIndication {
				ngap_message.SendErrorIndication(ran, aMFUENGAPID, rANUENGAPID, cause, nil)
			}
			removeRanUeByInvalidUE(ran, aMFUENGAPID, rANUENGAPID)
			return nil, fmt.Errorf("first returned message, but local RanUeNgapID is exist"+
				"[AmfUeNgapID: %d, remote RanUeNgapID: %s, local RanUeNgapID: %d]",
				aMFUENGAPID.Value, rANUENGAPID_string, ranUe.RanUeNgapId)
		}
	} else {
		if ranUe.RanUeNgapId != rANUENGAPID.Value {
			cause := &ngapType.Cause{
				Present: ngapType.CausePresentRadioNetwork,
				RadioNetwork: &ngapType.CauseRadioNetwork{
					Value: ngapType.CauseRadioNetworkPresentInconsistentRemoteUENGAPID,
				},
			}
			if sendErrorIndication {
				ngap_message.SendErrorIndication(ran, aMFUENGAPID, rANUENGAPID, cause, nil)
			}
			removeRanUeByInvalidUE(ran, aMFUENGAPID, rANUENGAPID)
			return nil, fmt.Errorf("inconsistent RanUe ID[AmfUeNgapID: %d, remote RanUeNgapID: %s, local RanUeNgapID: %d]",
				aMFUENGAPID.Value, rANUENGAPID_string, ranUe.RanUeNgapId)
		}
	}
	return ranUe, err
}
