package context

type IkeEventType int64

// IKE Event type
const (
	UnmarshalEAP5GDataResponse IkeEventType = iota
	SendEAP5GFailureMsg
	SendEAPNASMsg
	SendEAPSuccessMsg
	CreatePDUSession
	IKEDeleteRequest
	SendChildSADeleteRequest
	IKEContextUpdate
	GetNGAPContextResponse
)

type IkeEvt interface {
	Type() IkeEventType
}

type UnmarshalEAP5GDataResponseEvt struct {
	LocalSPI    uint64
	RanUeNgapId int64
	NasPDU      []byte
}

func (unmarshalEAP5GDataResponseEvt *UnmarshalEAP5GDataResponseEvt) Type() IkeEventType {
	return UnmarshalEAP5GDataResponse
}

func NewUnmarshalEAP5GDataResponseEvt(localSPI uint64, ranUeNgapId int64, nasPDU []byte,
) *UnmarshalEAP5GDataResponseEvt {
	return &UnmarshalEAP5GDataResponseEvt{
		LocalSPI:    localSPI,
		RanUeNgapId: ranUeNgapId,
		NasPDU:      nasPDU,
	}
}

type SendEAP5GFailureMsgEvt struct {
	LocalSPI uint64
	ErrMsg   EvtError
}

func (sendEAP5GFailureMsgEvt *SendEAP5GFailureMsgEvt) Type() IkeEventType {
	return SendEAP5GFailureMsg
}

func NewSendEAP5GFailureMsgEvt(localSPI uint64, errMsg EvtError,
) *SendEAP5GFailureMsgEvt {
	return &SendEAP5GFailureMsgEvt{
		LocalSPI: localSPI,
		ErrMsg:   errMsg,
	}
}

type SendEAPNASMsgEvt struct {
	LocalSPI uint64
	NasPDU   []byte
}

func (sendEAPNASMsgEvt *SendEAPNASMsgEvt) Type() IkeEventType {
	return SendEAPNASMsg
}

func NewSendEAPNASMsgEvt(localSPI uint64, nasPDU []byte,
) *SendEAPNASMsgEvt {
	return &SendEAPNASMsgEvt{
		LocalSPI: localSPI,
		NasPDU:   nasPDU,
	}
}

type SendEAPSuccessMsgEvt struct {
	LocalSPI          uint64
	Kn3iwf            []byte
	PduSessionListLen int
}

func (SendEAPSuccessMsgEvt *SendEAPSuccessMsgEvt) Type() IkeEventType {
	return SendEAPSuccessMsg
}

func NewSendEAPSuccessMsgEvt(localSPI uint64, kn3iwf []byte, pduSessionListLen int,
) *SendEAPSuccessMsgEvt {
	return &SendEAPSuccessMsgEvt{
		LocalSPI:          localSPI,
		Kn3iwf:            kn3iwf,
		PduSessionListLen: pduSessionListLen,
	}
}

type CreatePDUSessionEvt struct {
	LocalSPI                uint64
	PduSessionListLen       int
	TempPDUSessionSetupData *PDUSessionSetupTemporaryData
}

func (createPDUSessionEvt *CreatePDUSessionEvt) Type() IkeEventType {
	return CreatePDUSession
}

func NewCreatePDUSessionEvt(localSPI uint64, pduSessionListLen int,
	tempPDUSessionSetupData *PDUSessionSetupTemporaryData,
) *CreatePDUSessionEvt {
	return &CreatePDUSessionEvt{
		LocalSPI:                localSPI,
		PduSessionListLen:       pduSessionListLen,
		TempPDUSessionSetupData: tempPDUSessionSetupData,
	}
}

type IKEDeleteRequestEvt struct {
	LocalSPI uint64
}

func (ikeDeleteRequestEvt *IKEDeleteRequestEvt) Type() IkeEventType {
	return IKEDeleteRequest
}

func NewIKEDeleteRequestEvt(localSPI uint64,
) *IKEDeleteRequestEvt {
	return &IKEDeleteRequestEvt{
		LocalSPI: localSPI,
	}
}

type SendChildSADeleteRequestEvt struct {
	LocalSPI      uint64
	ReleaseIdList []int64
}

func (sendChildSADeleteRequestEvt *SendChildSADeleteRequestEvt) Type() IkeEventType {
	return SendChildSADeleteRequest
}

func NewSendChildSADeleteRequestEvt(localSPI uint64, releaseIdList []int64,
) *SendChildSADeleteRequestEvt {
	return &SendChildSADeleteRequestEvt{
		LocalSPI:      localSPI,
		ReleaseIdList: releaseIdList,
	}
}

type IKEContextUpdateEvt struct {
	LocalSPI uint64
	Kn3iwf   []byte
}

func (ikeContextUpdateEvt *IKEContextUpdateEvt) Type() IkeEventType {
	return IKEContextUpdate
}

func NewIKEContextUpdateEvt(localSPI uint64, kn3iwf []byte,
) *IKEContextUpdateEvt {
	return &IKEContextUpdateEvt{
		LocalSPI: localSPI,
		Kn3iwf:   kn3iwf,
	}
}

type GetNGAPContextRepEvt struct {
	LocalSPI          uint64
	NgapCxtReqNumlist []int64
	NgapCxt           []interface{}
}

func (getNGAPContextRepEvt *GetNGAPContextRepEvt) Type() IkeEventType {
	return GetNGAPContextResponse
}

func NewGetNGAPContextRepEvt(localSPI uint64, ngapCxtReqNumlist []int64,
	ngapCxt []interface{},
) *GetNGAPContextRepEvt {
	return &GetNGAPContextRepEvt{
		LocalSPI:          localSPI,
		NgapCxtReqNumlist: ngapCxtReqNumlist,
		NgapCxt:           ngapCxt,
	}
}
