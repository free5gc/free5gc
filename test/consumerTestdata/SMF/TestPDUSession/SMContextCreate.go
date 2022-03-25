package TestPDUSession

import (
	"bytes"

	"github.com/google/uuid"

	"github.com/free5gc/nas"
	"github.com/free5gc/nas/nasMessage"
	"github.com/free5gc/nas/nasType"
	"github.com/free5gc/openapi/models"
)

const (
	SERVICE_REQUEST = "Service Request"
	ACTIVATING      = "ACTIVATING"
)

type nasMessagePDUSessionEstablishmentRequestData struct {
	inExtendedProtocolDiscriminator                 uint8
	inPDUSessionID                                  uint8
	inPTI                                           uint8
	inPDUSESSIONESTABLISHMENTREQUESTMessageIdentity uint8
	inIntegrityProtectionMaximumDataRate            nasType.IntegrityProtectionMaximumDataRate
	inPDUSessionType                                nasType.PDUSessionType
	inSSCMode                                       nasType.SSCMode
	inCapability5GSM                                nasType.Capability5GSM
	inMaximumNumberOfSupportedPacketFilters         nasType.MaximumNumberOfSupportedPacketFilters
	inAlwaysonPDUSessionRequested                   nasType.AlwaysonPDUSessionRequested
	inSMPDUDNRequestContainer                       nasType.SMPDUDNRequestContainer
	inExtendedProtocolConfigurationOptions          nasType.ExtendedProtocolConfigurationOptions
}

var NasMessagePDUSessionEstablishmentRequestTable = make(map[string]nasMessagePDUSessionEstablishmentRequestData)

func init() {
	NasMessagePDUSessionEstablishmentRequestTable[SERVICE_REQUEST] = nasMessagePDUSessionEstablishmentRequestData{
		inExtendedProtocolDiscriminator: nasMessage.Epd5GSSessionManagementMessage,
		inPDUSessionID:                  0x01,
		inPTI:                           0x01,
		inPDUSESSIONESTABLISHMENTREQUESTMessageIdentity: nas.MsgTypePDUSessionEstablishmentRequest,
		inIntegrityProtectionMaximumDataRate: nasType.IntegrityProtectionMaximumDataRate{
			Iei:   0,
			Octet: [2]uint8{0x01, 0x01},
		},
		inPDUSessionType: nasType.PDUSessionType{
			Octet: 0x90,
		},
		inSSCMode: nasType.SSCMode{
			Octet: 0xA0,
		},
		inCapability5GSM: nasType.Capability5GSM{
			Iei:   nasMessage.PDUSessionEstablishmentRequestCapability5GSMType,
			Len:   2,
			Octet: [13]uint8{0x01, 0x01},
		},
		inMaximumNumberOfSupportedPacketFilters: nasType.MaximumNumberOfSupportedPacketFilters{
			Iei:   nasMessage.PDUSessionEstablishmentRequestMaximumNumberOfSupportedPacketFiltersType,
			Octet: [2]uint8{0x01, 0x01},
		},
		inAlwaysonPDUSessionRequested: nasType.AlwaysonPDUSessionRequested{
			Octet: 0xB0,
		},
		inSMPDUDNRequestContainer: nasType.SMPDUDNRequestContainer{
			Iei:    nasMessage.PDUSessionEstablishmentRequestSMPDUDNRequestContainerType,
			Len:    2,
			Buffer: []uint8{0x01, 0x01},
		},
		inExtendedProtocolConfigurationOptions: nasType.ExtendedProtocolConfigurationOptions{
			Iei:    nasMessage.PDUSessionEstablishmentRequestExtendedProtocolConfigurationOptionsType,
			Len:    2,
			Buffer: []uint8{0x01, 0x01},
		},
	}
}

func GetEstablishmentRequestData(testType string) (n1SmBytes []byte) {
	table := NasMessagePDUSessionEstablishmentRequestTable[testType]
	m := nas.NewMessage()
	m.GsmMessage = nas.NewGsmMessage()
	m.PDUSessionEstablishmentRequest = nasMessage.NewPDUSessionEstablishmentRequest(0x0)
	n1SmBuf := bytes.Buffer{}
	m.PDUSessionEstablishmentRequest.ExtendedProtocolDiscriminator.SetExtendedProtocolDiscriminator(table.inExtendedProtocolDiscriminator)
	m.PDUSessionEstablishmentRequest.PDUSessionID.SetPDUSessionID(table.inPDUSessionID)
	m.PDUSessionEstablishmentRequest.PTI.SetPTI(table.inPTI)
	m.PDUSessionEstablishmentRequest.PDUSESSIONESTABLISHMENTREQUESTMessageIdentity.SetMessageType(table.inPDUSESSIONESTABLISHMENTREQUESTMessageIdentity)
	m.PDUSessionEstablishmentRequest.IntegrityProtectionMaximumDataRate = table.inIntegrityProtectionMaximumDataRate

	m.PDUSessionEstablishmentRequest.PDUSessionType = nasType.NewPDUSessionType(nasMessage.PDUSessionEstablishmentRequestPDUSessionTypeType)
	m.PDUSessionEstablishmentRequest.PDUSessionType = &table.inPDUSessionType

	m.PDUSessionEstablishmentRequest.SSCMode = nasType.NewSSCMode(nasMessage.PDUSessionEstablishmentRequestSSCModeType)
	m.PDUSessionEstablishmentRequest.SSCMode = &table.inSSCMode

	m.PDUSessionEstablishmentRequest.Capability5GSM = nasType.NewCapability5GSM(nasMessage.PDUSessionEstablishmentRequestCapability5GSMType)
	m.PDUSessionEstablishmentRequest.Capability5GSM = &table.inCapability5GSM

	m.PDUSessionEstablishmentRequest.MaximumNumberOfSupportedPacketFilters = nasType.NewMaximumNumberOfSupportedPacketFilters(nasMessage.PDUSessionEstablishmentRequestMaximumNumberOfSupportedPacketFiltersType)
	m.PDUSessionEstablishmentRequest.MaximumNumberOfSupportedPacketFilters = &table.inMaximumNumberOfSupportedPacketFilters

	m.PDUSessionEstablishmentRequest.AlwaysonPDUSessionRequested = nasType.NewAlwaysonPDUSessionRequested(nasMessage.PDUSessionEstablishmentRequestAlwaysonPDUSessionRequestedType)
	m.PDUSessionEstablishmentRequest.AlwaysonPDUSessionRequested = &table.inAlwaysonPDUSessionRequested

	m.PDUSessionEstablishmentRequest.SMPDUDNRequestContainer = nasType.NewSMPDUDNRequestContainer(nasMessage.PDUSessionEstablishmentRequestSMPDUDNRequestContainerType)
	m.PDUSessionEstablishmentRequest.SMPDUDNRequestContainer = &table.inSMPDUDNRequestContainer

	m.PDUSessionEstablishmentRequest.ExtendedProtocolConfigurationOptions = nasType.NewExtendedProtocolConfigurationOptions(nasMessage.PDUSessionEstablishmentRequestExtendedProtocolConfigurationOptionsType)
	m.PDUSessionEstablishmentRequest.ExtendedProtocolConfigurationOptions = &table.inExtendedProtocolConfigurationOptions
	m.PDUSessionEstablishmentRequest.EncodePDUSessionEstablishmentRequest(&n1SmBuf)

	n1SmBytes = n1SmBuf.Bytes()
	return n1SmBytes
}

var ConsumerSMFPDUSessionSMContextCreateTable = make(map[string]models.SmContextCreateData)

func init() {
	ConsumerSMFPDUSessionSMContextCreateTable[SERVICE_REQUEST] = models.SmContextCreateData{
		Supi:                "imsi-2089300007487",
		UnauthenticatedSupi: false,
		PduSessionId:        2,
		Dnn:                 "internet",
		ServingNfId:         uuid.New().String(),
		Guami: &models.Guami{
			PlmnId: &models.PlmnId{
				Mcc: "208",
				Mnc: "93",
			},
			AmfId: "cafe00",
		},
		ServingNetwork: &models.PlmnId{
			Mcc: "208",
			Mnc: "93",
		},
		RequestType: models.RequestType_INITIAL_REQUEST,
		N1SmMsg: &models.RefToBinaryData{
			ContentId: "NGAP",
		},
		AnType:  models.AccessType__3_GPP_ACCESS,
		RatType: models.RatType_NR,
		SelMode: models.DnnSelectionMode_VERIFIED,
	}
}

type nasMessageULNASTransportData struct {
	inExtendedProtocolDiscriminator         uint8
	inSpareHalfOctetAndSecurityHeaderType   uint8
	inULNASTRANSPORTMessageIdentity         uint8
	inSpareHalfOctetAndPayloadContainerType uint8
	inPayloadContainer                      nasType.PayloadContainer
	inPduSessionID2Value                    nasType.PduSessionID2Value
	//inOldPDUSessionID                       nasType.OldPDUSessionID
	inRequestType nasType.RequestType
	inSNSSAI      nasType.SNSSAI
	//inAdditionalInformation                 nasType.AdditionalInformation
}

var NasMessageNasMessageULNASTransportDataTable = make(map[string]nasMessageULNASTransportData)

func init() {
	NasMessageNasMessageULNASTransportDataTable[SERVICE_REQUEST] = nasMessageULNASTransportData{
		inExtendedProtocolDiscriminator:         nasMessage.Epd5GSMobilityManagementMessage,
		inSpareHalfOctetAndSecurityHeaderType:   0x01,
		inULNASTRANSPORTMessageIdentity:         0x00,
		inSpareHalfOctetAndPayloadContainerType: nasMessage.PayloadContainerTypeN1SMInfo,
		inPayloadContainer: nasType.PayloadContainer{
			Iei:    (1), // n1Sminfo
			Len:    uint16(len(GetEstablishmentRequestData(SERVICE_REQUEST))),
			Buffer: GetEstablishmentRequestData(SERVICE_REQUEST),
		},
		inPduSessionID2Value: nasType.PduSessionID2Value{
			Iei:   nasMessage.ULNASTransportPduSessionID2ValueType,
			Octet: 10,
		},

		//inOldPDUSessionID

		inRequestType: nasType.RequestType{
			Octet: nasMessage.ULNASTransportRequestTypeType<<4 | nasMessage.ULNASTransportRequestTypeType,
		},
		inSNSSAI: nasType.SNSSAI{
			Iei:   nasMessage.ULNASTransportSNSSAIType,
			Len:   4,
			Octet: [8]uint8{0x01, 0x02, 0x03, 0x01},
		},
		//inAdditionalInformation
	}
}

func GetUlNasTransportData(testType string) (ulNasTransport nasMessage.ULNASTransport) {
	table := NasMessageNasMessageULNASTransportDataTable[testType]

	ulNasTransport.SetMessageType(nas.MsgTypeULNASTransport)
	ulNasTransport.SetExtendedProtocolDiscriminator(table.inExtendedProtocolDiscriminator) // 5GS MM
	ulNasTransport.SetPayloadContainerType(table.inSpareHalfOctetAndPayloadContainerType)  // n1SmInfo
	ulNasTransport.PayloadContainer.SetLen(table.inPayloadContainer.Len)
	ulNasTransport.SetPayloadContainerContents(table.inPayloadContainer.Buffer)
	ulNasTransport.PduSessionID2Value = &table.inPduSessionID2Value
	ulNasTransport.RequestType = &table.inRequestType
	ulNasTransport.SNSSAI = &table.inSNSSAI

	return ulNasTransport
}

var ConsumerSMFPDUSessionUpdateContextTable = make(map[string]models.UpdateSmContextRequest)

func init() {
	ConsumerSMFPDUSessionUpdateContextTable[ACTIVATING] = models.UpdateSmContextRequest{
		JsonData: &models.SmContextUpdateData{
			UpCnxState:  ACTIVATING,
			ServingNfId: uuid.New().String(),
			Guami: &models.Guami{
				PlmnId: &models.PlmnId{
					Mcc: "208",
					Mnc: "93",
				},
				AmfId: "cafe00",
			},
			ServingNetwork: &models.PlmnId{
				Mcc: "208",
				Mnc: "93",
			},
			N1SmMsg: &models.RefToBinaryData{
				ContentId: "NGAP",
			},
			AnType:  models.AccessType__3_GPP_ACCESS,
			RatType: models.RatType_NR,
		},
		BinaryDataN1SmMessage:     nil,
		BinaryDataN2SmInformation: nil,
	}
}
