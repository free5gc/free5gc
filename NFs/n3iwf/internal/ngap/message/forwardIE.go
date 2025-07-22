package message

import (
	"github.com/free5gc/aper"
	"github.com/free5gc/ngap/ngapType"
)

func AppendPDUSessionResourceSetupListCxtRes(
	list *ngapType.PDUSessionResourceSetupListCxtRes, pduSessionID int64, transfer []byte,
) {
	item := ngapType.PDUSessionResourceSetupItemCxtRes{}
	item.PDUSessionID.Value = pduSessionID
	item.PDUSessionResourceSetupResponseTransfer = transfer
	list.List = append(list.List, item)
}

func AppendPDUSessionResourceFailedToSetupListCxtRes(
	list *ngapType.PDUSessionResourceFailedToSetupListCxtRes, pduSessionID int64, transfer []byte,
) {
	item := ngapType.PDUSessionResourceFailedToSetupItemCxtRes{}
	item.PDUSessionID.Value = pduSessionID
	item.PDUSessionResourceSetupUnsuccessfulTransfer = transfer
	list.List = append(list.List, item)
}

func AppendPDUSessionResourceFailedToSetupListCxtfail(
	list *ngapType.PDUSessionResourceFailedToSetupListCxtFail, pduSessionID int64, transfer []byte,
) {
	item := ngapType.PDUSessionResourceFailedToSetupItemCxtFail{}
	item.PDUSessionID.Value = pduSessionID
	item.PDUSessionResourceSetupUnsuccessfulTransfer = transfer
	list.List = append(list.List, item)
}

func AppendPDUSessionResourceSetupListSURes(
	list *ngapType.PDUSessionResourceSetupListSURes, pduSessionID int64, transfer []byte,
) {
	item := ngapType.PDUSessionResourceSetupItemSURes{}
	item.PDUSessionID.Value = pduSessionID
	item.PDUSessionResourceSetupResponseTransfer = transfer
	list.List = append(list.List, item)
}

func AppendPDUSessionResourceFailedToSetupListSURes(
	list *ngapType.PDUSessionResourceFailedToSetupListSURes, pduSessionID int64, transfer []byte,
) {
	item := ngapType.PDUSessionResourceFailedToSetupItemSURes{}
	item.PDUSessionID.Value = pduSessionID
	item.PDUSessionResourceSetupUnsuccessfulTransfer = transfer
	list.List = append(list.List, item)
}

func AppendPDUSessionResourceModifyListModRes(
	list *ngapType.PDUSessionResourceModifyListModRes, pduSessionID int64, transfer []byte,
) {
	var pduSessionResourceModifyResponseTransfer aper.OctetString = transfer
	item := ngapType.PDUSessionResourceModifyItemModRes{}
	item.PDUSessionID.Value = pduSessionID
	item.PDUSessionResourceModifyResponseTransfer = pduSessionResourceModifyResponseTransfer
	list.List = append(list.List, item)
}

func AppendPDUSessionResourceFailedToModifyListModRes(
	list *ngapType.PDUSessionResourceFailedToModifyListModRes, pduSessionID int64, transfer []byte,
) {
	item := ngapType.PDUSessionResourceFailedToModifyItemModRes{}
	item.PDUSessionID.Value = pduSessionID
	item.PDUSessionResourceModifyUnsuccessfulTransfer = transfer
	list.List = append(list.List, item)
}
