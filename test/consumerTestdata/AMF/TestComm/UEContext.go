package TestComm

import (
	"github.com/free5gc/openapi/models"
)

const (
	CreateUEContext403           = "CreateUEContext403"
	CreateUEContext201           = "CreateUEContext201"
	UeContextRelease404          = "UeContextRelease404"
	UeContextRelease201          = "UeContextRelease201"
	UeContextTransfer404         = "UeContextTransfer404"
	UeContextTransferINIT_REG200 = "UeContextTransferINIT_REG200"
	UeContextTransferMOBI_REG200 = "UeContextTransferMOBI_REG200"
	AssignEbiData403             = "AssignEbiData403"
	AssignEbiData200             = "AssignEbiData200"
	RegistrationStatusUpdate404  = "RegistrationStatusUpdate404"
	RegistrationStatusUpdate200  = "RegistrationStatusUpdate200"
)

var ConsumerAMFCreateUEContextRequsetTable = make(map[string]models.CreateUeContextRequest)

func init() {
	ConsumerAMFCreateUEContextRequsetTable[CreateUEContext403] = models.CreateUeContextRequest{
		JsonData: &models.UeContextCreateData{
			UeContext: &models.UeContext{
				Supi: "imsi-2089300007487",
			},
			TargetId:           &models.NgRanTargetId{},
			SourceToTargetData: &models.N2InfoContent{},
			PduSessionList:     []models.N2SmInformation{},
			N2NotifyUri:        "127.0.0.1",
			UeRadioCapability:  nil,
			NgapCause:          nil,
			SupportedFeatures:  "",
		},
	}
	ConsumerAMFCreateUEContextRequsetTable[CreateUEContext201] = models.CreateUeContextRequest{
		JsonData: &models.UeContextCreateData{
			UeContext: &models.UeContext{
				Supi: "imsi-2089300007487",
				RestrictedRatList: []models.RatType{
					models.RatType_NR,
				},
			},
			TargetId: &models.NgRanTargetId{
				RanNodeId: &models.GlobalRanNodeId{
					PlmnId: &models.PlmnId{
						Mcc: "208",
						Mnc: "93",
					},
					N3IwfId: "123",
					GNbId: &models.GNbId{
						BitLength: 123,
						GNBValue:  "string",
					},
					NgeNbId: "string",
				},
				Tai: &models.Tai{
					PlmnId: &models.PlmnId{
						Mcc: "208",
						Mnc: "93",
					},
					Tac: "000001",
				},
			},
			SourceToTargetData: &models.N2InfoContent{
				NgapMessageType: 0,
				NgapIeType:      "NgapIeType_PDU_RES_SETUP_REQ",
				NgapData: &models.RefToBinaryData{
					ContentId: "N2SmInfo",
				},
			},
			PduSessionList: []models.N2SmInformation{
				{
					PduSessionId: 10,
					N2InfoContent: &models.N2InfoContent{
						NgapIeType: models.NgapIeType_PDU_RES_REL_CMD,
						NgapData: &models.RefToBinaryData{
							ContentId: "N2SmInfo",
						},
					},
				},
			},
			N2NotifyUri:       "127.0.0.1",
			UeRadioCapability: nil,
			NgapCause:         nil,
			SupportedFeatures: "",
		},
	}
}

var ConsumerAMFReleaseUEContextRequestTable = make(map[string]models.UeContextRelease)

func init() {
	ConsumerAMFReleaseUEContextRequestTable[UeContextRelease404] = models.UeContextRelease{
		Supi:                "",
		UnauthenticatedSupi: false,
		NgapCause: &models.NgApCause{
			Group: 0,
			Value: 0,
		},
	}
	ConsumerAMFReleaseUEContextRequestTable[UeContextRelease201] = models.UeContextRelease{
		Supi:                "imsi-2089300007487",
		UnauthenticatedSupi: true,
		NgapCause: &models.NgApCause{
			Group: 0,
			Value: 0,
		},
	}

}

var ConsumerAMFUEContextTransferRequestTable = make(map[string]models.UeContextTransferRequest)

func init() {
	ConsumerAMFUEContextTransferRequestTable[UeContextTransfer404] = models.UeContextTransferRequest{
		JsonData: &models.UeContextTransferReqData{
			Reason:            "",
			AccessType:        "",
			PlmnId:            nil,
			RegRequest:        nil,
			SupportedFeatures: "",
		},
	}
	ConsumerAMFUEContextTransferRequestTable[UeContextTransferINIT_REG200] = models.UeContextTransferRequest{
		JsonData: &models.UeContextTransferReqData{
			Reason:            models.TransferReason_INIT_REG,
			AccessType:        models.AccessType__3_GPP_ACCESS,
			PlmnId:            nil,
			RegRequest:        nil,
			SupportedFeatures: "",
		},
	}
	ConsumerAMFUEContextTransferRequestTable[UeContextTransferMOBI_REG200] = models.UeContextTransferRequest{
		JsonData: &models.UeContextTransferReqData{
			Reason:            models.TransferReason_MOBI_REG,
			AccessType:        models.AccessType__3_GPP_ACCESS,
			PlmnId:            nil,
			RegRequest:        nil,
			SupportedFeatures: "",
		},
	}
}

var ConsumerAMFUEContextEBIAssignmentTable = make(map[string]models.AssignEbiData)

func init() {
	ConsumerAMFUEContextEBIAssignmentTable[AssignEbiData403] = models.AssignEbiData{
		PduSessionId:    0,
		ArpList:         nil,
		ReleasedEbiList: nil,
	}
	ConsumerAMFUEContextEBIAssignmentTable[AssignEbiData200] = models.AssignEbiData{
		PduSessionId:    10,
		ArpList:         nil,
		ReleasedEbiList: nil,
	}
}

var ConsumerRegistrationStatusUpdateTable = make(map[string]models.UeRegStatusUpdateReqData)

func init() {
	ConsumerRegistrationStatusUpdateTable[RegistrationStatusUpdate200] = models.UeRegStatusUpdateReqData{
		TransferStatus:       models.UeContextTransferStatus_TRANSFERRED,
		ToReleaseSessionList: nil,
		PcfReselectedInd:     false,
	}
	ConsumerRegistrationStatusUpdateTable[RegistrationStatusUpdate404] = models.UeRegStatusUpdateReqData{
		TransferStatus:       "",
		ToReleaseSessionList: nil,
		PcfReselectedInd:     false,
	}

}
