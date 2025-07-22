package message

import (
	"encoding/hex"

	"github.com/free5gc/amf/internal/context"
	"github.com/free5gc/amf/internal/logger"
	"github.com/free5gc/amf/internal/util"
	"github.com/free5gc/ngap/ngapConvert"
	"github.com/free5gc/ngap/ngapType"
	"github.com/free5gc/openapi/models"
)

func AppendPDUSessionResourceSetupListSUReq(list *ngapType.PDUSessionResourceSetupListSUReq,
	pduSessionId int32, snssai models.Snssai, nasPDU []byte, transfer []byte,
) {
	var item ngapType.PDUSessionResourceSetupItemSUReq
	item.PDUSessionID.Value = int64(pduSessionId)
	item.SNSSAI = ngapConvert.SNssaiToNgap(snssai)
	item.PDUSessionResourceSetupRequestTransfer = transfer
	if nasPDU != nil {
		item.PDUSessionNASPDU = new(ngapType.NASPDU)
		item.PDUSessionNASPDU.Value = nasPDU
	}
	list.List = append(list.List, item)
}

func AppendPDUSessionResourceSetupListHOReq(list *ngapType.PDUSessionResourceSetupListHOReq,
	pduSessionId int32, snssai models.Snssai, transfer []byte,
) {
	var item ngapType.PDUSessionResourceSetupItemHOReq
	item.PDUSessionID.Value = int64(pduSessionId)
	item.SNSSAI = ngapConvert.SNssaiToNgap(snssai)
	item.HandoverRequestTransfer = transfer
	list.List = append(list.List, item)
}

func AppendPDUSessionResourceSetupListCxtReq(list *ngapType.PDUSessionResourceSetupListCxtReq,
	pduSessionId int32, snssai models.Snssai, nasPDU []byte, transfer []byte,
) {
	var item ngapType.PDUSessionResourceSetupItemCxtReq
	item.PDUSessionID.Value = int64(pduSessionId)
	item.SNSSAI = ngapConvert.SNssaiToNgap(snssai)
	if nasPDU != nil {
		item.NASPDU = new(ngapType.NASPDU)
		item.NASPDU.Value = nasPDU
	}
	item.PDUSessionResourceSetupRequestTransfer = transfer
	list.List = append(list.List, item)
}

func ConvertPDUSessionResourceSetupListCxtReqToSUReq(
	listCxtReq *ngapType.PDUSessionResourceSetupListCxtReq,
) *ngapType.PDUSessionResourceSetupListSUReq {
	if listCxtReq == nil {
		return nil
	}
	listSUReq := ngapType.PDUSessionResourceSetupListSUReq{}
	for _, itemCxt := range listCxtReq.List {
		var itemSU ngapType.PDUSessionResourceSetupItemSUReq
		itemSU.PDUSessionID = itemCxt.PDUSessionID
		itemSU.PDUSessionNASPDU = itemCxt.NASPDU
		itemSU.SNSSAI = itemCxt.SNSSAI
		itemSU.PDUSessionResourceSetupRequestTransfer = itemCxt.PDUSessionResourceSetupRequestTransfer
		listSUReq.List = append(listSUReq.List, itemSU)
	}
	return &listSUReq
}

func AppendPDUSessionResourceModifyListModReq(list *ngapType.PDUSessionResourceModifyListModReq,
	pduSessionId int32, nasPDU []byte, transfer []byte,
) {
	var item ngapType.PDUSessionResourceModifyItemModReq
	item.PDUSessionID.Value = int64(pduSessionId)
	item.PDUSessionResourceModifyRequestTransfer = transfer
	if nasPDU != nil {
		item.NASPDU = new(ngapType.NASPDU)
		item.NASPDU.Value = nasPDU
	}
	list.List = append(list.List, item)
}

func AppendPDUSessionResourceModifyListModCfm(list *ngapType.PDUSessionResourceModifyListModCfm,
	pduSessionId int64, transfer []byte,
) {
	var item ngapType.PDUSessionResourceModifyItemModCfm
	item.PDUSessionID.Value = pduSessionId
	item.PDUSessionResourceModifyConfirmTransfer = transfer
	list.List = append(list.List, item)
}

func AppendPDUSessionResourceFailedToModifyListModCfm(list *ngapType.PDUSessionResourceFailedToModifyListModCfm,
	pduSessionId int64, transfer []byte,
) {
	var item ngapType.PDUSessionResourceFailedToModifyItemModCfm
	item.PDUSessionID.Value = pduSessionId
	item.PDUSessionResourceModifyIndicationUnsuccessfulTransfer = transfer
	list.List = append(list.List, item)
}

func AppendPDUSessionResourceToReleaseListRelCmd(list *ngapType.PDUSessionResourceToReleaseListRelCmd,
	pduSessionId int32, transfer []byte,
) {
	var item ngapType.PDUSessionResourceToReleaseItemRelCmd
	item.PDUSessionID.Value = int64(pduSessionId)
	item.PDUSessionResourceReleaseCommandTransfer = transfer
	list.List = append(list.List, item)
}

func BuildIEMobilityRestrictionList(ue *context.AmfUe) ngapType.MobilityRestrictionList {
	mobilityRestrictionList := ngapType.MobilityRestrictionList{}
	mobilityRestrictionList.ServingPLMN = ngapConvert.PlmnIdToNgap(ue.PlmnId)

	if ue.AccessAndMobilitySubscriptionData != nil && len(ue.AccessAndMobilitySubscriptionData.RatRestrictions) > 0 {
		mobilityRestrictionList.RATRestrictions = new(ngapType.RATRestrictions)
		ratRestrictions := mobilityRestrictionList.RATRestrictions
		for _, ratType := range ue.AccessAndMobilitySubscriptionData.RatRestrictions {
			item := ngapType.RATRestrictionsItem{}
			item.PLMNIdentity = ngapConvert.PlmnIdToNgap(ue.PlmnId)
			item.RATRestrictionInformation = ngapConvert.RATRestrictionInformationToNgap(ratType)
			ratRestrictions.List = append(ratRestrictions.List, item)
		}
	}

	if ue.AccessAndMobilitySubscriptionData != nil && len(ue.AccessAndMobilitySubscriptionData.ForbiddenAreas) > 0 {
		mobilityRestrictionList.ForbiddenAreaInformation = new(ngapType.ForbiddenAreaInformation)
		forbiddenAreaInformation := mobilityRestrictionList.ForbiddenAreaInformation
		for _, info := range ue.AccessAndMobilitySubscriptionData.ForbiddenAreas {
			item := ngapType.ForbiddenAreaInformationItem{}
			item.PLMNIdentity = ngapConvert.PlmnIdToNgap(ue.PlmnId)
			for _, tac := range info.Tacs {
				tacBytes, err := hex.DecodeString(tac)
				if err != nil {
					logger.NgapLog.Errorf(
						"[Error] DecodeString tac error: %+v", err)
				}
				tacNgap := ngapType.TAC{}
				tacNgap.Value = tacBytes
				item.ForbiddenTACs.List = append(item.ForbiddenTACs.List, tacNgap)
			}
			forbiddenAreaInformation.List = append(forbiddenAreaInformation.List, item)
		}
	}

	if ue.AmPolicyAssociation.ServAreaRes != nil {
		mobilityRestrictionList.ServiceAreaInformation = new(ngapType.ServiceAreaInformation)
		serviceAreaInformation := mobilityRestrictionList.ServiceAreaInformation

		item := ngapType.ServiceAreaInformationItem{}
		item.PLMNIdentity = ngapConvert.PlmnIdToNgap(ue.PlmnId)
		var tacList []ngapType.TAC
		for _, area := range ue.AmPolicyAssociation.ServAreaRes.Areas {
			for _, tac := range area.Tacs {
				tacBytes, err := hex.DecodeString(tac)
				if err != nil {
					logger.NgapLog.Errorf(
						"[Error] DecodeString tac error: %+v", err)
				}
				tacNgap := ngapType.TAC{}
				tacNgap.Value = tacBytes
				tacList = append(tacList, tacNgap)
			}
		}
		if ue.AmPolicyAssociation.ServAreaRes.RestrictionType == models.RestrictionType_ALLOWED_AREAS {
			item.AllowedTACs = new(ngapType.AllowedTACs)
			item.AllowedTACs.List = append(item.AllowedTACs.List, tacList...)
		} else {
			item.NotAllowedTACs = new(ngapType.NotAllowedTACs)
			item.NotAllowedTACs.List = append(item.NotAllowedTACs.List, tacList...)
		}
		serviceAreaInformation.List = append(serviceAreaInformation.List, item)
	}
	return mobilityRestrictionList
}

func BuildUnavailableGUAMIList(guamiList []models.Guami) (unavailableGUAMIList ngapType.UnavailableGUAMIList) {
	for _, guami := range guamiList {
		item := ngapType.UnavailableGUAMIItem{}
		item.GUAMI.PLMNIdentity = ngapConvert.PlmnIdToNgap(util.PlmnIdNidToModelsPlmnId(*guami.PlmnId))
		regionId, setId, ptrId := ngapConvert.AmfIdToNgap(guami.AmfId)
		item.GUAMI.AMFRegionID.Value = regionId
		item.GUAMI.AMFSetID.Value = setId
		item.GUAMI.AMFPointer.Value = ptrId
		// TODO: item.TimerApproachForGUAMIRemoval and item.BackupAMFName not support yet
		unavailableGUAMIList.List = append(unavailableGUAMIList.List, item)
	}
	return
}
