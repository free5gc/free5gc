package ngap

import (
	"encoding/hex"
	"fmt"
	"net"
	"testing"

	"github.com/google/uuid"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/require"

	amf_context "github.com/free5gc/amf/internal/context"
	"github.com/free5gc/amf/internal/logger"
	nastesting "github.com/free5gc/amf/internal/nas/testing"
	ngaptesting "github.com/free5gc/amf/internal/ngap/testing"
	"github.com/free5gc/amf/pkg/factory"
	"github.com/free5gc/aper"
	"github.com/free5gc/nas/nasMessage"
	"github.com/free5gc/nas/nasType"
	"github.com/free5gc/ngap"
	"github.com/free5gc/ngap/ngapType"
	"github.com/free5gc/openapi/models"
)

func NewAmfRan(conn net.Conn) *amf_context.AmfRan {
	ran := amf_context.AmfRan{
		RanPresent: 1,
		RanId: &models.GlobalRanNodeId{
			PlmnId: &models.PlmnId{
				Mcc: "208",
				Mnc: "93",
			},
			GNbId: &models.GNbId{
				BitLength: 24,
				GNBValue:  "000102",
			},
		},
		Name:   "free5gc",
		AnType: "3GPP_ACCESS",

		/* socket Connect*/
		Conn: conn,

		/* Supported TA List */
		SupportedTAList: []amf_context.SupportedTAI{
			{
				Tai: models.Tai{
					PlmnId: &models.PlmnId{
						Mcc: "208",
						Mnc: "93",
					},
					Tac: "000001",
				},
				SNssaiList: []models.Snssai{
					{
						Sst: 1,
						Sd:  "010203",
					},
				},
			},
		},

		/* logger */
		Log: logger.NgapLog.WithField(logger.FieldRanAddr, "127.0.0.1"),
	}
	return &ran
}

func NewAmfContext(amfCtx *amf_context.AMFContext) {
	*amfCtx = amf_context.AMFContext{
		NfId: uuid.New().String(),

		NgapIpList:   []string{"127.0.0.1"},
		NgapPort:     38412,
		UriScheme:    "http",
		RegisterIPv4: "127.0.0.18",
		BindingIPv4:  "127.0.0.18",
		SBIPort:      8000,
		ServedGuamiList: []models.Guami{
			{
				PlmnId: &models.PlmnIdNid{
					Mcc: "208",
					Mnc: "93",
				},
				AmfId: "cafe00",
			},
		},
		SupportTaiLists: []models.Tai{
			{
				PlmnId: &models.PlmnId{
					Mcc: "208",
					Mnc: "93",
				},
				Tac: "000001",
			},
		},
		PlmnSupportList: []factory.PlmnSupportItem{
			{
				PlmnId: &models.PlmnId{
					Mcc: "208",
					Mnc: "93",
				},
				SNssaiList: []models.Snssai{
					{
						Sst: 1,
						Sd:  "010203",
					},
					{
						Sst: 1,
						Sd:  "112233",
					},
				},
			},
		},
		SupportDnnLists: []string{
			"internet",
		},
		NrfUri: "http://127.0.0.10:8000",
		SecurityAlgorithm: amf_context.SecurityAlgorithm{
			IntegrityOrder: []uint8{0x02},
			CipheringOrder: []uint8{0x00},
		},
		NetworkName: factory.NetworkName{
			Full:  "free5GC",
			Short: "free",
		},
		T3502Value:             720,
		T3512Value:             3600,
		Non3gppDeregTimerValue: 3240,
		T3513Cfg: factory.TimerValue{
			Enable:        true,
			ExpireTime:    6000000000,
			MaxRetryTimes: 4,
		},
		T3522Cfg: factory.TimerValue{
			Enable:        true,
			ExpireTime:    6000000000,
			MaxRetryTimes: 4,
		},
		T3550Cfg: factory.TimerValue{
			Enable:        true,
			ExpireTime:    6000000000,
			MaxRetryTimes: 4,
		},
		T3560Cfg: factory.TimerValue{
			Enable:        true,
			ExpireTime:    6000000000,
			MaxRetryTimes: 4,
		},
		T3565Cfg: factory.TimerValue{
			Enable:        true,
			ExpireTime:    6000000000,
			MaxRetryTimes: 4,
		},
	}
}

func BuildInitialUEMessage(ranUeNgapID int64, nasPdu []byte, fiveGSTmsi string) ngapType.NGAPPDU {
	var TestPlmn ngapType.PLMNIdentity
	var tmsi aper.OctetString
	var amfSetID, amfPointer []byte
	var err error

	pdu := ngapType.NGAPPDU{
		Present: ngapType.NGAPPDUPresentInitiatingMessage,
		InitiatingMessage: &ngapType.InitiatingMessage{
			ProcedureCode: ngapType.ProcedureCode{
				Value: ngapType.ProcedureCodeInitialUEMessage,
			},
			Criticality: ngapType.Criticality{
				Value: ngapType.CriticalityPresentIgnore,
			},
			Value: ngapType.InitiatingMessageValue{
				Present: ngapType.InitiatingMessagePresentInitialUEMessage,
				InitialUEMessage: &ngapType.InitialUEMessage{
					ProtocolIEs: ngapType.ProtocolIEContainerInitialUEMessageIEs{
						List: []ngapType.InitialUEMessageIEs{
							{
								Id: ngapType.ProtocolIEID{
									Value: ngapType.ProtocolIEIDRANUENGAPID,
								},
								Criticality: ngapType.Criticality{
									Value: ngapType.CriticalityPresentReject,
								},
								Value: ngapType.InitialUEMessageIEsValue{
									Present: ngapType.InitialUEMessageIEsPresentRANUENGAPID,
									RANUENGAPID: &ngapType.RANUENGAPID{
										Value: ranUeNgapID,
									},
								},
							},
							{
								Id: ngapType.ProtocolIEID{
									Value: ngapType.ProtocolIEIDNASPDU,
								},
								Criticality: ngapType.Criticality{
									Value: ngapType.CriticalityPresentReject,
								},
								Value: ngapType.InitialUEMessageIEsValue{
									Present: ngapType.InitialUEMessageIEsPresentNASPDU,
									NASPDU: &ngapType.NASPDU{
										Value: nasPdu,
									},
								},
							},
							{
								Id: ngapType.ProtocolIEID{
									Value: ngapType.ProtocolIEIDUserLocationInformation,
								},
								Criticality: ngapType.Criticality{
									Value: ngapType.CriticalityPresentReject,
								},
								Value: ngapType.InitialUEMessageIEsValue{
									Present: ngapType.InitialUEMessageIEsPresentUserLocationInformation,
									UserLocationInformation: &ngapType.UserLocationInformation{
										UserLocationInformationNR: &ngapType.UserLocationInformationNR{
											NRCGI: ngapType.NRCGI{
												PLMNIdentity: ngapType.PLMNIdentity{
													Value: TestPlmn.Value,
												},
												NRCellIdentity: ngapType.NRCellIdentity{
													Value: aper.BitString{
														Bytes:     []byte{0x00, 0x00, 0x00, 0x00, 0x10},
														BitLength: 36,
													},
												},
											},
											TAI: ngapType.TAI{
												PLMNIdentity: ngapType.PLMNIdentity{
													Value: TestPlmn.Value,
												},
												TAC: ngapType.TAC{
													Value: aper.OctetString("\x00\x00\x01"),
												},
											},
										},
									},
								},
							},
							{
								Id: ngapType.ProtocolIEID{
									Value: ngapType.ProtocolIEIDRRCEstablishmentCause,
								},
								Criticality: ngapType.Criticality{
									Value: ngapType.CriticalityPresentIgnore,
								},
								Value: ngapType.InitialUEMessageIEsValue{
									Present: ngapType.InitialUEMessageIEsPresentRRCEstablishmentCause,
									RRCEstablishmentCause: &ngapType.RRCEstablishmentCause{
										Value: ngapType.RRCEstablishmentCausePresentMtAccess,
									},
								},
							},
							{
								Id: ngapType.ProtocolIEID{
									Value: ngapType.ProtocolIEIDFiveGSTMSI,
								},
								Criticality: ngapType.Criticality{
									Value: ngapType.CriticalityPresentReject,
								},
							},
						},
					},
				},
			},
		},
	}

	if fiveGSTmsi != "" {
		ie := ngapType.InitialUEMessageIEs{}
		ie.Id.Value = ngapType.ProtocolIEIDFiveGSTMSI
		ie.Criticality.Value = ngapType.CriticalityPresentReject
		ie.Value.Present = ngapType.InitialUEMessageIEsPresentFiveGSTMSI
		ie.Value.FiveGSTMSI = new(ngapType.FiveGSTMSI)

		fiveGSTMSI := ie.Value.FiveGSTMSI

		amfSetID, err = hex.DecodeString(fiveGSTmsi[:4])
		if err != nil {
			fmt.Println("DecodeString AMFSetID error in BuildInitialUEMessage")
		}

		fiveGSTMSI.AMFSetID.Value = aper.BitString{
			Bytes:     amfSetID,
			BitLength: 10,
		}

		amfPointer, err = hex.DecodeString(fiveGSTmsi[2:4])
		if err != nil {
			fmt.Println("DecodeString AMFPointer error in BuildInitialUEMessage")
		}

		fiveGSTMSI.AMFPointer.Value = aper.BitString{
			Bytes:     amfPointer,
			BitLength: 6,
		}

		tmsi, err = hex.DecodeString(fiveGSTmsi[4:])
		if err != nil {
			fmt.Println("DecodeString 5G-S-TMSI error in BuildInitialUEMessage")
		}

		fiveGSTMSI.FiveGTMSI.Value = tmsi

		pdu.InitiatingMessage.Value.InitialUEMessage.ProtocolIEs.List = append(
			pdu.InitiatingMessage.Value.InitialUEMessage.ProtocolIEs.List, ie)
	}

	return pdu
}

func TestHandleInitialUEMessage(t *testing.T) {
	var message *ngapType.NGAPPDU
	var ranUeNgapID int64 = 1
	fiveGSTmsi := "fe0000000001"
	var msg ngapType.NGAPPDU

	testCases := []struct {
		amfUENGAPID      int
		nasPdu           []byte
		paramStr         string
		resultStr        string
		expectedResponse aper.OctetString
	}{
		{
			amfUENGAPID:      1,
			nasPdu:           nastesting.GetServiceRequest(nasMessage.ServiceTypeData),
			paramStr:         "Service Request after NGSetup. AMF can not recognize UE.",
			resultStr:        "DownlinkNASTransport with Service Reject",
			expectedResponse: aper.OctetString{0x7e, 0x00, 0x4d, 0x0a},
		},
		{
			amfUENGAPID: 2,
			nasPdu: nastesting.GetRegistrationRequest(
				nasMessage.RegistrationType5GSPeriodicRegistrationUpdating,
				nasType.MobileIdentity5GS{
					Len:    12, // suci
					Buffer: []uint8{0x01, 0x02, 0xf8, 0x39, 0xf0, 0xff, 0x00, 0x00, 0x00, 0x00, 0x47, 0x78},
				},
				nil, nil, nil, nil, nil),
			paramStr:         "Periodic Registration. AMF can not recognize UE.",
			resultStr:        "DownlinkNASTransport with Registration Reject",
			expectedResponse: aper.OctetString{0x7e, 0x00, 0x44, 0x0a, 0x16, 0x01, 0x2C},
		},
		{
			amfUENGAPID: 3,
			nasPdu: nastesting.GetRegistrationRequest(
				nasMessage.RegistrationType5GSMobilityRegistrationUpdating,
				nasType.MobileIdentity5GS{
					Len:    12, // suci
					Buffer: []uint8{0x01, 0x02, 0xf8, 0x39, 0xf0, 0xff, 0x00, 0x00, 0x00, 0x00, 0x47, 0x78},
				},
				nil, nil, nil, nil, nil),
			paramStr:         "Mobility Registration. AMF can not recognize UE.",
			resultStr:        "DownlinkNASTransport with Registration Reject",
			expectedResponse: aper.OctetString{0x7e, 0x00, 0x44, 0x0a, 0x16, 0x01, 0x2C},
		},
	}

	for i, testcase := range testCases {
		infoStr := fmt.Sprintf("testcase[%d]: ", i)

		// Set up fake connection
		connStub := new(ngaptesting.SctpConnStub)

		// Init AMF context
		amf_self := amf_context.GetSelf()
		NewAmfContext(amf_self)

		// Set up AmfRan
		ran := NewAmfRan(connStub)

		// Set up message
		msg = BuildInitialUEMessage(ranUeNgapID, testcase.nasPdu, fiveGSTmsi)
		message = &msg
		initiatingMessage := message.InitiatingMessage
		require.NotNil(t, initiatingMessage)

		handlerInitialUEMessage(ran, &msg, initiatingMessage)
		Convey(infoStr, t, func() {
			Convey(testcase.paramStr, func() {
				Convey(testcase.resultStr, func() {
					rcv, err := ngap.Decoder(connStub.MsgList[0])
					if err != nil {
						fmt.Println("decode ngap message error")
					}
					ieListDownlinkNASTransport := rcv.InitiatingMessage.Value.DownlinkNASTransport.ProtocolIEs.List
					for _, ie := range ieListDownlinkNASTransport {
						switch ie.Value.Present {
						case ngapType.DownlinkNASTransportIEsPresentAMFUENGAPID:
							Convey("AMFUENGAPID", func() {
								So(ie.Value.AMFUENGAPID.Value, ShouldEqual, testcase.amfUENGAPID)
							})
						case ngapType.DownlinkNASTransportIEsPresentRANUENGAPID:
							Convey("RANUENGAPID", func() {
								So(ie.Value.RANUENGAPID.Value, ShouldEqual, 1)
							})
						case ngapType.DownlinkNASTransportIEsPresentNASPDU:
							Convey("Cause5GMMImplicitlyDeregistered", func() {
								So(ie.Value.NASPDU.Value, ShouldResemble, testcase.expectedResponse)
							})
						}
					}
				})

				Convey("UEContextReleaseCommand ", func() {
					rcv, err := ngap.Decoder(connStub.MsgList[1])
					if err != nil {
						fmt.Println("decode ngap message error")
					}

					ieListUEContextReleaseCommand := rcv.InitiatingMessage.Value.UEContextReleaseCommand.ProtocolIEs.List
					for _, ie := range ieListUEContextReleaseCommand {
						switch ie.Value.Present {
						case ngapType.UEContextReleaseCommandIEsPresentUENGAPIDs:
							Convey("AMFUENGAPID", func() {
								So(ie.Value.UENGAPIDs.UENGAPIDPair.AMFUENGAPID.Value, ShouldEqual, testcase.amfUENGAPID)
							})
							Convey("RANUENGAPID", func() {
								So(ie.Value.UENGAPIDs.UENGAPIDPair.RANUENGAPID.Value, ShouldEqual, 1)
							})
						case ngapType.UEContextReleaseCommandIEsPresentCause:
							Convey("UEContextReleaseCommandIEsPresentCause", func() {
								So(ie.Value.Cause.Nas.Value, ShouldEqual, ngapType.CauseNasPresentNormalRelease)
							})
						}
					}
				})
			})
		})
	}
}
