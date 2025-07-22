package ngap

import (
	"encoding/binary"
	"errors"

	"github.com/free5gc/aper"
	"github.com/free5gc/n3iwf/internal/logger"
	"github.com/free5gc/n3iwf/internal/ngap/message"
	"github.com/free5gc/ngap/ngapType"
)

// 3GPP specified EAP-5G

// Access Network Parameters
type ANParameters struct {
	GUAMI              *ngapType.GUAMI
	SelectedPLMNID     *ngapType.PLMNIdentity
	RequestedNSSAI     *ngapType.AllowedNSSAI
	EstablishmentCause *ngapType.RRCEstablishmentCause
}

func UnmarshalEAP5GData(
	codedData []byte,
) (
	anParameters *ANParameters,
	nasPDU []byte,
	err error,
) {
	ngapLog := logger.NgapLog
	if len(codedData) < 2 {
		return nil, nil, errors.New("no data to decode")
	}

	ngapLog.Debug("===== Unmarshal EAP5G Data (Ref: TS24.502 Fig. 9.3.2.2.2-1) =====")

	codedData = codedData[2:]

	// [TS 24.502 f30] 9.3.2.2.2.3
	// AN-parameter value field in GUAMI, PLMN ID and NSSAI is coded as value part
	// Therefore, IEI of AN-parameter is not needed to be included.
	// anParameter = AN-parameter Type | AN-parameter Length | Value part of IE

	if len(codedData) < 2 {
		ngapLog.Error("No AN-Parameter type or length specified")
		return nil, nil, errors.New("error formatting")
	}

	// Length of the AN-Parameter field
	anParameterLength := binary.BigEndian.Uint16(codedData[:2])
	ngapLog.Debugf("AN-parameters length: %d", anParameterLength)

	if anParameterLength != 0 {
		anParameterField := codedData[2:]

		// Bound checking
		if len(anParameterField) < int(anParameterLength) {
			ngapLog.Error("Packet contained error length of value")
			return nil, nil, errors.New("error formatting")
		}
		anParameterField = anParameterField[:anParameterLength]

		ngapLog.Debugf("Parsing AN-parameters...: % v", anParameterField)

		anParameters = new(ANParameters)

		// Parse AN-Parameters
		for len(anParameterField) >= 2 {
			parameterType := anParameterField[0]
			// The AN-parameter length field indicates the length of the AN-parameter value field.
			parameterLength := anParameterField[1]

			switch parameterType {
			case message.ANParametersTypeGUAMI:
				ngapLog.Debugf("-> Parameter type: GUAMI")
				if parameterLength != 0 {
					parameterValue := anParameterField[2:]

					if len(parameterValue) < int(parameterLength) {
						return nil, nil, errors.New("error formatting")
					}
					parameterValue = parameterValue[:parameterLength]

					if len(parameterValue) != message.ANParametersLenGUAMI {
						return nil, nil, errors.New("unmatched GUAMI length")
					}

					guamiField := make([]byte, 1)
					guamiField = append(guamiField, parameterValue...)
					// Decode GUAMI using aper
					ngapGUAMI := new(ngapType.GUAMI)
					err = aper.UnmarshalWithParams(guamiField, ngapGUAMI, "valueExt")
					if err != nil {
						ngapLog.Errorf("APER unmarshal with parameter failed: %+v", err)
						return nil, nil, errors.New("unmarshal failed when decoding GUAMI")
					}
					anParameters.GUAMI = ngapGUAMI
					ngapLog.Debugf("Unmarshal GUAMI: % x", guamiField)
					ngapLog.Debugf("\tGUAMI: PLMNIdentity[% x], "+
						"AMFRegionID[% x], AMFSetID[% x], AMFPointer[% x]",
						anParameters.GUAMI.PLMNIdentity, anParameters.GUAMI.AMFRegionID,
						anParameters.GUAMI.AMFSetID, anParameters.GUAMI.AMFPointer)
				} else {
					ngapLog.Warn("AN-Parameter GUAMI field empty")
				}
			case message.ANParametersTypeSelectedPLMNID:
				ngapLog.Debugf("-> Parameter type: ANParametersTypeSelectedPLMNID")
				if parameterLength != 0 {
					parameterValue := anParameterField[2:]

					if len(parameterValue) < int(parameterLength) {
						return nil, nil, errors.New("error formatting")
					}
					parameterValue = parameterValue[:parameterLength]

					if len(parameterValue) != message.ANParametersLenPLMNID {
						return nil, nil, errors.New("unmatched PLMN ID length")
					}

					plmnField := make([]byte, 1)
					plmnField = append(plmnField, parameterValue...)
					// Decode PLMN using aper
					ngapPLMN := new(ngapType.PLMNIdentity)
					err = aper.UnmarshalWithParams(plmnField, ngapPLMN, "valueExt")
					if err != nil {
						ngapLog.Errorf("APER unmarshal with parameter failed: %v", err)
						return nil, nil, errors.New("unmarshal failed when decoding PLMN")
					}
					anParameters.SelectedPLMNID = ngapPLMN
					ngapLog.Debugf("Unmarshal SelectedPLMNID: % x", plmnField)
					ngapLog.Debugf("\tSelectedPLMNID: % x", anParameters.SelectedPLMNID.Value)
				} else {
					ngapLog.Warn("AN-Parameter PLMN field empty")
				}
			case message.ANParametersTypeRequestedNSSAI:
				ngapLog.Debugf("-> Parameter type: ANParametersTypeRequestedNSSAI")
				if parameterLength != 0 {
					parameterValue := anParameterField[2:]

					if len(parameterValue) < int(parameterLength) {
						return nil, nil, errors.New("error formatting")
					}
					parameterValue = parameterValue[:parameterLength]

					ngapNSSAI := new(ngapType.AllowedNSSAI)

					// [TS 24501 f30] 9.11.2.8 S-NSSAI
					// s-nssai(LV) consists of
					// len(1 byte) | SST(1) | SD(3,opt) | Mapped HPLMN SST (1,opt) | Mapped HPLMN SD (3,opt)
					// The length of minimum s-nssai comprised of a length and a SST is 2 bytes.

					for len(parameterValue) >= 2 {
						snssaiLength := parameterValue[0]
						snssaiValue := parameterValue[1:]

						if len(snssaiValue) < int(snssaiLength) {
							ngapLog.Error("SNSSAI length error")
							return nil, nil, errors.New("error formatting")
						}
						snssaiValue = snssaiValue[:snssaiLength]

						ngapSNSSAIItem := ngapType.AllowedNSSAIItem{}

						if len(snssaiValue) == 1 {
							ngapSNSSAIItem.SNSSAI = ngapType.SNSSAI{
								SST: ngapType.SST{
									Value: []byte{snssaiValue[0]},
								},
							}
						} else if len(snssaiValue) == 4 {
							ngapSNSSAIItem.SNSSAI = ngapType.SNSSAI{
								SST: ngapType.SST{
									Value: []byte{snssaiValue[0]},
								},
								SD: &ngapType.SD{
									Value: []byte{snssaiValue[1], snssaiValue[2], snssaiValue[3]},
								},
							}
						} else {
							ngapLog.Error("Empty SNSSAI value")
							return nil, nil, errors.New("error formatting")
						}

						ngapNSSAI.List = append(ngapNSSAI.List, ngapSNSSAIItem)

						ngapLog.Debugf("Unmarshal SNSSAI: % x", parameterValue[:1+snssaiLength])
						ngapLog.Debugf("\t\t\tSST: % x", ngapSNSSAIItem.SNSSAI.SST.Value)
						sd := ngapSNSSAIItem.SNSSAI.SD
						if sd == nil {
							ngapLog.Debugf("\t\t\tSD: nil")
						} else {
							ngapLog.Debugf("\t\t\tSD: % x", sd.Value)
						}

						// shift parameterValue for parsing next s-nssai
						parameterValue = parameterValue[1+snssaiLength:]
					}
					anParameters.RequestedNSSAI = ngapNSSAI
				} else {
					ngapLog.Warn("AN-Parameter NSSAI is empty")
				}
			case message.ANParametersTypeEstablishmentCause:
				ngapLog.Debugf("-> Parameter type: ANParametersTypeEstablishmentCause")
				if parameterLength != 0 {
					parameterValue := anParameterField[2:]

					if len(parameterValue) < int(parameterLength) {
						return nil, nil, errors.New("error formatting")
					}
					parameterValue = parameterValue[:parameterLength]

					if len(parameterValue) != message.ANParametersLenEstCause {
						return nil, nil, errors.New("unmatched Establishment Cause length")
					}

					ngapLog.Debugf("Unmarshal ANParametersTypeEstablishmentCause: % x", parameterValue)

					establishmentCause := parameterValue[0] & 0x0f
					switch establishmentCause {
					case message.EstablishmentCauseEmergency:
						ngapLog.Trace("AN-Parameter establishment cause: Emergency")
					case message.EstablishmentCauseHighPriorityAccess:
						ngapLog.Trace("AN-Parameter establishment cause: High Priority Access")
					case message.EstablishmentCauseMO_Signalling:
						ngapLog.Trace("AN-Parameter establishment cause: MO Signalling")
					case message.EstablishmentCauseMO_Data:
						ngapLog.Trace("AN-Parameter establishment cause: MO Data")
					case message.EstablishmentCauseMPS_PriorityAccess:
						ngapLog.Trace("AN-Parameter establishment cause: MPS Priority Access")
					case message.EstablishmentCauseMCS_PriorityAccess:
						ngapLog.Trace("AN-Parameter establishment cause: MCS Priority Access")
					default:
						ngapLog.Trace("AN-Parameter establishment cause: Unknown. Treat as mo-Data")
						establishmentCause = message.EstablishmentCauseMO_Data
					}

					ngapEstablishmentCause := new(ngapType.RRCEstablishmentCause)
					ngapEstablishmentCause.Value = aper.Enumerated(establishmentCause)

					anParameters.EstablishmentCause = ngapEstablishmentCause
				} else {
					ngapLog.Warn("AN-Parameter establishment cause field empty")
				}
			default:
				ngapLog.Warn("Unsopprted AN-Parameter. Ignore.")
			}

			// shift anParameterField
			anParameterField = anParameterField[2+parameterLength:]
		}
	}

	// shift codedData
	codedData = codedData[2+anParameterLength:]

	if len(codedData) < 2 {
		ngapLog.Error("No NASPDU length specified")
		return nil, nil, errors.New("error formatting")
	}
	// Length of the NASPDU field
	nasPDULength := binary.BigEndian.Uint16(codedData[:2])
	ngapLog.Debugf("nasPDULength: %d", nasPDULength)

	if nasPDULength == 0 {
		ngapLog.Error("No NAS PDU included in EAP-5G packet")
		return nil, nil, errors.New("no NAS PDU")
	}

	nasPDUField := codedData[2:]
	// Bound checking
	if len(nasPDUField) < int(nasPDULength) {
		return nil, nil, errors.New("error formatting")
	} else {
		nasPDUField = nasPDUField[:nasPDULength]
	}
	ngapLog.Debugf("nasPDUField: % v", nasPDUField)

	nasPDU = append(nasPDU, nasPDUField...)
	return anParameters, nasPDU, err
}
