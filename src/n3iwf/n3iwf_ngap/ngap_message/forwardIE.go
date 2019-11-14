package ngap_message

import (
	"free5gc/lib/ngap/ngapType"
)

func AppendPDUSessionResourceSetupListCxtRes(list *ngapType.PDUSessionResourceSetupListCxtRes, pduSessionID int64, transfer []byte) {
	item := ngapType.PDUSessionResourceSetupItemCxtRes{}
	item.PDUSessionID.Value = pduSessionID
	item.PDUSessionResourceSetupResponseTransfer = transfer
	list.List = append(list.List, item)
}

func AppendPDUSessionResourceFailedToSetupListCxtRes(list *ngapType.PDUSessionResourceFailedToSetupListCxtRes, pduSessionID int64, transfer []byte) {
	item := ngapType.PDUSessionResourceFailedToSetupItemCxtRes{}
	item.PDUSessionID.Value = pduSessionID
	item.PDUSessionResourceSetupUnsuccessfulTransfer = transfer
	list.List = append(list.List, item)
}

func AppendPDUSessionResourceFailedToSetupListCxtfail(list *ngapType.PDUSessionResourceFailedToSetupListCxtFail, pduSessionID int64, transfer []byte) {
	item := ngapType.PDUSessionResourceFailedToSetupItemCxtFail{}
	item.PDUSessionID.Value = pduSessionID
	item.PDUSessionResourceSetupUnsuccessfulTransfer = transfer
	list.List = append(list.List, item)
}
