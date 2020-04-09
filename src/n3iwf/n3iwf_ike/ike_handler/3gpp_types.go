package ike_handler

import (
	"encoding/binary"
	"errors"

	"free5gc/lib/aper"
	"free5gc/lib/ngap/ngapType"
	"free5gc/src/n3iwf/n3iwf_ike/ike_message"
)

// 3GPP specified EAP-5G

// Access Network Parameters
type ANParameters struct {
	GUAMI              *ngapType.GUAMI
	SelectedPLMNID     *ngapType.PLMNIdentity
	RequestedNSSAI     *ngapType.AllowedNSSAI
	EstablishmentCause *ngapType.RRCEstablishmentCause
}

func UnmarshalEAP5GData(codedData []byte) (eap5GMessageID uint8, anParameters *ANParameters, nasPDU []byte, err error) {
	if len(codedData) >= 2 {
		eap5GMessageID = codedData[0]

		if eap5GMessageID == ike_message.EAP5GType5GStop {
			return
		}

		codedData = codedData[2:]

		if len(codedData) >= 2 {
			// Length of the AN-Parameter field
			anParameterLength := binary.BigEndian.Uint16(codedData[:2])

			if anParameterLength != 0 {
				anParameterField := codedData[2:]

				// Bound checking
				if len(anParameterField) < int(anParameterLength) {
					ikeLog.Error("[IKE] Packet contained error length of value")
					return 0, nil, nil, errors.New("Error formatting")
				} else {
					anParameterField = anParameterField[:anParameterLength]
				}

				anParameters = new(ANParameters)

				// Parse AN-Parameters
				for len(anParameterField) >= 2 {
					parameterType := anParameterField[0]
					// The AN-parameter length field indicates the length of the AN-parameter value field.
					parameterLength := anParameterField[1]

					switch parameterType {
					case ike_message.ANParametersTypeGUAMI:
						if parameterLength != 0 {
							parameterValue := anParameterField[2:]

							if len(parameterValue) < int(parameterLength) {
								return 0, nil, nil, errors.New("Error formatting")
							} else {
								parameterValue = parameterValue[:parameterLength]
							}

							if len(parameterValue) != 7 {
								return 0, nil, nil, errors.New("Unmatched GUAMI length")
							}

							guamiField := make([]byte, 1)
							guamiField = append(guamiField, parameterValue[1:]...)
							// Decode GUAMI using aper
							ngapGUAMI := new(ngapType.GUAMI)
							err := aper.UnmarshalWithParams(guamiField, ngapGUAMI, "valueExt")
							if err != nil {
								ikeLog.Errorf("[IKE] APER unmarshal with parameter failed: %+v", err)
								return 0, nil, nil, errors.New("Unmarshal failed when decoding GUAMI")
							}
							anParameters.GUAMI = ngapGUAMI
						} else {
							ikeLog.Warn("[IKE] AN-Parameter GUAMI field empty")
						}
					case ike_message.ANParametersTypeSelectedPLMNID:
						if parameterLength != 0 {
							parameterValue := anParameterField[2:]

							if len(parameterValue) < int(parameterLength) {
								return 0, nil, nil, errors.New("Error formatting")
							} else {
								parameterValue = parameterValue[:parameterLength]
							}

							if len(parameterValue) != 5 {
								return 0, nil, nil, errors.New("Unmatched PLMN ID length")
							}

							plmnField := make([]byte, 1)
							plmnField = append(plmnField, parameterValue[2:]...)
							// Decode PLMN using aper
							ngapPLMN := new(ngapType.PLMNIdentity)
							err := aper.UnmarshalWithParams(plmnField, ngapPLMN, "valueExt")
							if err != nil {
								ikeLog.Errorf("[IKE] APER unmarshal with parameter failed: %v", err)
								return 0, nil, nil, errors.New("Unmarshal failed when decoding PLMN")
							}
							anParameters.SelectedPLMNID = ngapPLMN
						} else {
							ikeLog.Warn("[IKE] AN-Parameter PLMN field empty")
						}
					case ike_message.ANParametersTypeRequestedNSSAI:
						if parameterLength != 0 {
							parameterValue := anParameterField[2:]

							if len(parameterValue) < int(parameterLength) {
								return 0, nil, nil, errors.New("Error formatting")
							} else {
								parameterValue = parameterValue[:parameterLength]
							}

							if len(parameterValue) >= 2 {
								nssaiLength := parameterValue[1]

								if nssaiLength != 0 {
									nssaiValue := parameterValue[2:]

									if len(nssaiValue) < int(nssaiLength) {
										return 0, nil, nil, errors.New("Error formatting")
									} else {
										nssaiValue = nssaiValue[:nssaiLength]
									}

									ngapNSSAI := new(ngapType.AllowedNSSAI)

									for len(nssaiValue) >= 2 {
										snssaiLength := nssaiValue[1]

										if snssaiLength != 0 {
											snssaiValue := nssaiValue[2:]

											if len(snssaiValue) < int(snssaiLength) {
												return 0, nil, nil, errors.New("Error formatting")
											} else {
												snssaiValue = snssaiValue[:snssaiLength]
											}

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
												ikeLog.Error("[IKE] SNSSAI length error")
												return 0, nil, nil, errors.New("Error formatting")
											}

											ngapNSSAI.List = append(ngapNSSAI.List, ngapSNSSAIItem)
										} else {
											ikeLog.Error("[IKE] Empty SNSSAI value")
											return 0, nil, nil, errors.New("Error formatting")
										}

										// shift nssaiValue
										nssaiValue = nssaiValue[2+snssaiLength:]
									}

									anParameters.RequestedNSSAI = ngapNSSAI
								} else {
									ikeLog.Error("[IKE] Empty NSSAI value")
									return 0, nil, nil, errors.New("Error formatting")
								}
							} else {
								ikeLog.Error("[IKE] No NSSAI type or length specified")
								return 0, nil, nil, errors.New("Error formatting")
							}

						} else {
							ikeLog.Warn("[IKE] AN-Parameter value for NSSAI empty")
						}
					case ike_message.ANParametersTypeEstablishmentCause:
						if parameterLength != 0 {
							parameterValue := anParameterField[2:]

							if len(parameterValue) < int(parameterLength) {
								return 0, nil, nil, errors.New("Error formatting")
							} else {
								parameterValue = parameterValue[:parameterLength]
							}

							if len(parameterValue) != 2 {
								return 0, nil, nil, errors.New("Unmatched Establishment Cause length")
							}

							establishmentCause := parameterValue[1] & 0x0f
							switch establishmentCause {
							case ike_message.EstablishmentCauseEmergency:
								ikeLog.Trace("[IKE] AN-Parameter establishment cause: Emergency")
							case ike_message.EstablishmentCauseHighPriorityAccess:
								ikeLog.Trace("[IKE] AN-Parameter establishment cause: High Priority Access")
							case ike_message.EstablishmentCauseMO_Signalling:
								ikeLog.Trace("[IKE] AN-Parameter establishment cause: MO Signalling")
							case ike_message.EstablishmentCauseMO_Data:
								ikeLog.Trace("[IKE] AN-Parameter establishment cause: MO Data")
							case ike_message.EstablishmentCauseMPS_PriorityAccess:
								ikeLog.Trace("[IKE] AN-Parameter establishment cause: MPS Priority Access")
							case ike_message.EstablishmentCauseMCS_PriorityAccess:
								ikeLog.Trace("[IKE] AN-Parameter establishment cause: MCS Priority Access")
							default:
								ikeLog.Trace("[IKE] AN-Parameter establishment cause: Unknown. Treat as mo-Data")
								establishmentCause = ike_message.EstablishmentCauseMO_Data
							}

							ngapEstablishmentCause := new(ngapType.RRCEstablishmentCause)
							ngapEstablishmentCause.Value = aper.Enumerated(establishmentCause)

							anParameters.EstablishmentCause = ngapEstablishmentCause
						} else {
							ikeLog.Warn("[IKE] AN-Parameter establishment cause field empty")
						}
					default:
						ikeLog.Warn("[IKE] Unsopprted AN-Parameter. Ignore.")
					}

					// shift anParameterField
					anParameterField = anParameterField[2+parameterLength:]
				}
			}

			// shift codedData
			codedData = codedData[2+anParameterLength:]
		} else {
			ikeLog.Error("[IKE] No AN-Parameter type or length specified")
			return 0, nil, nil, errors.New("Error formatting")
		}

		if len(codedData) >= 2 {
			// Length of the NASPDU field
			nasPDULength := binary.BigEndian.Uint16(codedData[:2])

			if nasPDULength != 0 {
				nasPDUField := codedData[2:]

				// Bound checking
				if len(nasPDUField) < int(nasPDULength) {
					return 0, nil, nil, errors.New("Error formatting")
				} else {
					nasPDUField = nasPDUField[:nasPDULength]
				}

				nasPDU = append(nasPDU, nasPDUField...)
			} else {
				ikeLog.Error("[IKE] No NAS PDU included in EAP-5G packet")
				return 0, nil, nil, errors.New("No NAS PDU")
			}
		} else {
			ikeLog.Error("[IKE] No NASPDU length specified")
			return 0, nil, nil, errors.New("Error formatting")
		}

		return
	} else {
		return 0, nil, nil, errors.New("No data to decode")
	}
}
