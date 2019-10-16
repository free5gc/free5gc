//go:binary-only-package

package nasMessage_test

import (
	"bytes"
	"free5gc/lib/nas"
	"free5gc/lib/nas/logger"
	"free5gc/lib/nas/nasMessage"
	"free5gc/lib/nas/nasType"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

type nasMessageULNASTransportData struct {
	inExtendedProtocolDiscriminator         uint8
	inSecurityHeader                        uint8
	inSpareHalfOctet                        uint8
	inULNASTRANSPORTMessageIdentity         uint8
	inSpareHalfOctetAndPayloadContainerType nasType.SpareHalfOctetAndPayloadContainerType
	inPayloadContainer                      nasType.PayloadContainer
	inPduSessionID2Value                    nasType.PduSessionID2Value
	inOldPDUSessionID                       nasType.OldPDUSessionID
	inRequestType                           nasType.RequestType
	inSNSSAI                                nasType.SNSSAI
	inDNN                                   nasType.DNN
	inAdditionalInformation                 nasType.AdditionalInformation
}

var nasMessageULNASTransportTable = []nasMessageULNASTransportData{
	{
		inExtendedProtocolDiscriminator: nasMessage.Epd5GSMobilityManagementMessage,
		inSecurityHeader:                0x01,
		inSpareHalfOctet:                0x01,
		inULNASTRANSPORTMessageIdentity: nas.MsgTypeULNASTransport,
		inSpareHalfOctetAndPayloadContainerType: nasType.SpareHalfOctetAndPayloadContainerType{
			Octet: 0x01,
		},
		inPayloadContainer: nasType.PayloadContainer{
			Len:    2,
			Buffer: []uint8{0x01, 0x01},
		},
		inPduSessionID2Value: nasType.PduSessionID2Value{
			Iei:   nasMessage.ULNASTransportPduSessionID2ValueType,
			Octet: 0x01,
		},
		inOldPDUSessionID: nasType.OldPDUSessionID{
			Iei:   nasMessage.ULNASTransportOldPDUSessionIDType,
			Octet: 0x01,
		},
		inRequestType: nasType.RequestType{
			Octet: 0x80,
		},
		inDNN: nasType.DNN{
			Iei:    nasMessage.ULNASTransportDNNType,
			Len:    2,
			Buffer: []uint8{0x01, 0x01},
		},
		inSNSSAI: nasType.SNSSAI{
			Iei:   nasMessage.ULNASTransportSNSSAIType,
			Len:   8,
			Octet: [8]uint8{0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01},
		},
		inAdditionalInformation: nasType.AdditionalInformation{
			Iei:    nasMessage.ULNASTransportAdditionalInformationType,
			Len:    2,
			Buffer: []uint8{0x01, 0x01},
		},
	},
}

func TestNasTypeNewULNASTransport(t *testing.T) {}

func TestNasTypeNewULNASTransportMessage(t *testing.T) {}
