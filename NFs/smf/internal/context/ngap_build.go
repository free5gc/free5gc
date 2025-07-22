package context

import (
	"encoding/binary"
	"fmt"

	"github.com/free5gc/aper"
	"github.com/free5gc/ngap/ngapConvert"
	"github.com/free5gc/ngap/ngapType"
	"github.com/free5gc/openapi/models"
)

const DefaultNonGBR5QI = 9

func BuildPDUSessionResourceSetupRequestTransfer(ctx *SMContext) ([]byte, error) {
	ANUPF := ctx.Tunnel.DataPathPool.GetDefaultPath().FirstDPNode
	UpNode := ANUPF.UPF
	teidOct := make([]byte, 4)
	teidOctForSplitPDUSession := make([]byte, 4)
	binary.BigEndian.PutUint32(teidOct, ctx.LocalULTeid)
	binary.BigEndian.PutUint32(teidOctForSplitPDUSession, ctx.LocalULTeidForSplitPDUSession)

	resourceSetupRequestTransfer := ngapType.PDUSessionResourceSetupRequestTransfer{}

	// PDU Session Aggregate Maximum Bit Rate
	// This IE is Conditional and shall be present when at least one NonGBR QoS flow is being setup.
	// TODO: should check if there is at least one NonGBR QoS flow
	ie := ngapType.PDUSessionResourceSetupRequestTransferIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDPDUSessionAggregateMaximumBitRate
	ie.Criticality.Value = ngapType.CriticalityPresentReject
	sessRule := ctx.SelectedSessionRule()
	if sessRule == nil || sessRule.AuthSessAmbr == nil {
		return nil, fmt.Errorf("no PDU Session AMBR")
	}
	ie.Value = ngapType.PDUSessionResourceSetupRequestTransferIEsValue{
		Present: ngapType.PDUSessionResourceSetupRequestTransferIEsPresentPDUSessionAggregateMaximumBitRate,
		PDUSessionAggregateMaximumBitRate: &ngapType.PDUSessionAggregateMaximumBitRate{
			PDUSessionAggregateMaximumBitRateDL: ngapType.BitRate{
				Value: ngapConvert.UEAmbrToInt64(sessRule.AuthSessAmbr.Downlink),
			},
			PDUSessionAggregateMaximumBitRateUL: ngapType.BitRate{
				Value: ngapConvert.UEAmbrToInt64(sessRule.AuthSessAmbr.Uplink),
			},
		},
	}
	resourceSetupRequestTransfer.ProtocolIEs.List = append(resourceSetupRequestTransfer.ProtocolIEs.List, ie)

	// UL NG-U UP TNL Information
	ie = ngapType.PDUSessionResourceSetupRequestTransferIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDULNGUUPTNLInformation
	ie.Criticality.Value = ngapType.CriticalityPresentReject
	if n3IP, err := UpNode.N3Interfaces[0].IP(ctx.SelectedPDUSessionType); err != nil {
		return nil, err
	} else {
		ie.Value = ngapType.PDUSessionResourceSetupRequestTransferIEsValue{
			Present: ngapType.PDUSessionResourceSetupRequestTransferIEsPresentULNGUUPTNLInformation,
			ULNGUUPTNLInformation: &ngapType.UPTransportLayerInformation{
				Present: ngapType.UPTransportLayerInformationPresentGTPTunnel,
				GTPTunnel: &ngapType.GTPTunnel{
					TransportLayerAddress: ngapType.TransportLayerAddress{
						Value: aper.BitString{
							Bytes:     n3IP,
							BitLength: uint64(len(n3IP) * 8),
						},
					},
					GTPTEID: ngapType.GTPTEID{Value: teidOct},
				},
			},
		}
	}

	resourceSetupRequestTransfer.ProtocolIEs.List = append(resourceSetupRequestTransfer.ProtocolIEs.List, ie)

	// Additional UL NG-U UP TNL Information
	ie = ngapType.PDUSessionResourceSetupRequestTransferIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDAdditionalULNGUUPTNLInformation
	ie.Criticality.Value = ngapType.CriticalityPresentIgnore
	if n3IP, err := UpNode.N3Interfaces[0].IP(ctx.SelectedPDUSessionType); err != nil {
		return nil, err
	} else {
		ie.Value = ngapType.PDUSessionResourceSetupRequestTransferIEsValue{
			Present: ngapType.PDUSessionResourceSetupRequestTransferIEsPresentAdditionalULNGUUPTNLInformation,
			AdditionalULNGUUPTNLInformation: &ngapType.UPTransportLayerInformationList{
				List: []ngapType.UPTransportLayerInformationItem{
					{
						NGUUPTNLInformation: ngapType.UPTransportLayerInformation{
							Present: ngapType.UPTransportLayerInformationPresentGTPTunnel,
							GTPTunnel: &ngapType.GTPTunnel{
								TransportLayerAddress: ngapType.TransportLayerAddress{
									Value: aper.BitString{
										Bytes:     n3IP,
										BitLength: uint64(len(n3IP) * 8),
									},
								},
								GTPTEID: ngapType.GTPTEID{Value: teidOctForSplitPDUSession},
							},
						},
					},
				},
			},
		}
	}
	resourceSetupRequestTransfer.ProtocolIEs.List = append(resourceSetupRequestTransfer.ProtocolIEs.List, ie)

	// PDU Session Type
	ie = ngapType.PDUSessionResourceSetupRequestTransferIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDPDUSessionType
	ie.Criticality.Value = ngapType.CriticalityPresentReject
	ie.Value = ngapType.PDUSessionResourceSetupRequestTransferIEsValue{
		Present: ngapType.PDUSessionResourceSetupRequestTransferIEsPresentPDUSessionType,
		PDUSessionType: &ngapType.PDUSessionType{
			Value: ngapType.PDUSessionTypePresentIpv4,
		},
	}
	resourceSetupRequestTransfer.ProtocolIEs.List = append(resourceSetupRequestTransfer.ProtocolIEs.List, ie)

	// QoS Flow Setup Request List
	// use Default 5qi, arp

	authDefQos := sessRule.AuthDefQos
	ie = ngapType.PDUSessionResourceSetupRequestTransferIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDQosFlowSetupRequestList
	ie.Criticality.Value = ngapType.CriticalityPresentReject
	ie.Value = ngapType.PDUSessionResourceSetupRequestTransferIEsValue{
		Present: ngapType.PDUSessionResourceSetupRequestTransferIEsPresentQosFlowSetupRequestList,
		QosFlowSetupRequestList: &ngapType.QosFlowSetupRequestList{
			List: []ngapType.QosFlowSetupRequestItem{
				{
					QosFlowIdentifier: ngapType.QosFlowIdentifier{
						Value: int64(sessRule.DefQosQFI),
					},
					QosFlowLevelQosParameters: ngapType.QosFlowLevelQosParameters{
						QosCharacteristics: ngapType.QosCharacteristics{
							Present: ngapType.QosCharacteristicsPresentNonDynamic5QI,
							NonDynamic5QI: &ngapType.NonDynamic5QIDescriptor{
								FiveQI: ngapType.FiveQI{
									Value: int64(authDefQos.Var5qi),
								},
							},
						},
						AllocationAndRetentionPriority: ngapType.AllocationAndRetentionPriority{
							PriorityLevelARP: ngapType.PriorityLevelARP{
								Value: int64(authDefQos.Arp.PriorityLevel),
							},
							PreEmptionCapability: ngapType.PreEmptionCapability{
								Value: ngapType.PreEmptionCapabilityPresentShallNotTriggerPreEmption,
							},
							PreEmptionVulnerability: ngapType.PreEmptionVulnerability{
								Value: ngapType.PreEmptionVulnerabilityPresentNotPreEmptable,
							},
						},
					},
				},
			},
		},
	}

	for _, qosFlow := range ctx.AdditonalQosFlows {
		if qosDesc, err := qosFlow.BuildNgapQosFlowSetupRequestItem(); err != nil {
			return nil, fmt.Errorf("encode BuildNgapQosFlowSetupRequestItem failed: %s", err)
		} else {
			ie.Value.QosFlowSetupRequestList.List = append(ie.Value.QosFlowSetupRequestList.List, qosDesc)
		}
	}

	resourceSetupRequestTransfer.ProtocolIEs.List = append(resourceSetupRequestTransfer.ProtocolIEs.List, ie)

	// Security Indication to NG-RAN (optional) TS 38.413 9.3.1.27
	// Only over 3GPP access TS 23.501 5.10.3
	if ctx.AnType == models.AccessType__3_GPP_ACCESS && ctx.UpSecurity != nil {
		upSecurity := ctx.UpSecurity
		maximumDataRatePerUEForUserPlaneIntegrityProtectionForUpLink := ctx.
			MaximumDataRatePerUEForUserPlaneIntegrityProtectionForUpLink

		ie = ngapType.PDUSessionResourceSetupRequestTransferIEs{}
		ie.Id.Value = ngapType.ProtocolIEIDSecurityIndication
		ie.Criticality.Value = ngapType.CriticalityPresentReject

		securityIndication := new(ngapType.SecurityIndication)

		switch upSecurity.UpIntegr {
		case models.UpIntegrity_REQUIRED:
			securityIndication.IntegrityProtectionIndication.Value = ngapType.
				IntegrityProtectionIndicationPresentRequired
		case models.UpIntegrity_PREFERRED:
			securityIndication.IntegrityProtectionIndication.Value = ngapType.
				IntegrityProtectionIndicationPresentPreferred
		case models.UpIntegrity_NOT_NEEDED:
			securityIndication.IntegrityProtectionIndication.Value = ngapType.
				IntegrityProtectionIndicationPresentNotNeeded
		}
		switch upSecurity.UpConfid {
		case models.UpConfidentiality_REQUIRED:
			securityIndication.ConfidentialityProtectionIndication.Value = ngapType.
				ConfidentialityProtectionIndicationPresentRequired
		case models.UpConfidentiality_PREFERRED:
			securityIndication.ConfidentialityProtectionIndication.Value = ngapType.
				ConfidentialityProtectionIndicationPresentPreferred
		case models.UpConfidentiality_NOT_NEEDED:
			securityIndication.ConfidentialityProtectionIndication.Value = ngapType.
				ConfidentialityProtectionIndicationPresentNotNeeded
		}
		// Present only when Integrity Indication within the Security Indication is set to "required" or "preferred"
		integrityProtectionInd := securityIndication.IntegrityProtectionIndication.Value
		if integrityProtectionInd == ngapType.IntegrityProtectionIndicationPresentRequired ||
			integrityProtectionInd ==
				ngapType.IntegrityProtectionIndicationPresentPreferred {
			securityIndication.MaximumIntegrityProtectedDataRateUL = new(ngapType.MaximumIntegrityProtectedDataRate)
			switch maximumDataRatePerUEForUserPlaneIntegrityProtectionForUpLink {
			case models.MaxIntegrityProtectedDataRate_MAX_UE_RATE:
				securityIndication.MaximumIntegrityProtectedDataRateUL.Value = ngapType.
					MaximumIntegrityProtectedDataRatePresentMaximumUERate
			case models.MaxIntegrityProtectedDataRate__64_KBPS:
				securityIndication.MaximumIntegrityProtectedDataRateUL.Value = ngapType.
					MaximumIntegrityProtectedDataRatePresentBitrate64kbs
			}
		}
		ie.Value = ngapType.PDUSessionResourceSetupRequestTransferIEsValue{
			Present:            ngapType.PDUSessionResourceSetupRequestTransferIEsPresentSecurityIndication,
			SecurityIndication: securityIndication,
		}
		resourceSetupRequestTransfer.ProtocolIEs.List = append(resourceSetupRequestTransfer.ProtocolIEs.List, ie)
	}

	if buf, err := aper.MarshalWithParams(resourceSetupRequestTransfer, "valueExt"); err != nil {
		return nil, fmt.Errorf("encode resourceSetupRequestTransfer failed: %s", err)
	} else {
		return buf, nil
	}
}

func BuildPDUSessionResourceModifyRequestTransfer(ctx *SMContext) ([]byte, error) {
	resourceModifyRequestTransfer := ngapType.PDUSessionResourceModifyRequestTransfer{}
	ie := ngapType.PDUSessionResourceModifyRequestTransferIEs{}

	ie.Id.Value = ngapType.ProtocolIEIDQosFlowAddOrModifyRequestList
	ie.Criticality.Value = ngapType.CriticalityPresentReject
	ie.Value.Present = ngapType.PDUSessionResourceModifyRequestTransferIEsPresentQosFlowAddOrModifyRequestList
	ie.Value.QosFlowAddOrModifyRequestList = new(ngapType.QosFlowAddOrModifyRequestList)
	qosFlowAddOrModifyRequestList := ie.Value.QosFlowAddOrModifyRequestList

	for _, qos := range ctx.AdditonalQosFlows {
		if qos.State == QoSFlowUnset || qos.State == QoSFlowToBeModify {
			if qosDesc, err := qos.BuildNgapQosFlowAddOrModifyRequestItem(); err != nil {
				return nil, fmt.Errorf("BuildNgapQosFlowSetupRequestItem failed: %s", err)
			} else {
				qosFlowAddOrModifyRequestList.List = append(qosFlowAddOrModifyRequestList.List, qosDesc)
			}
		}
	}

	resourceModifyRequestTransfer.ProtocolIEs.List = append(resourceModifyRequestTransfer.ProtocolIEs.List, ie)

	if buf, err := aper.MarshalWithParams(resourceModifyRequestTransfer, "valueExt"); err != nil {
		return nil, fmt.Errorf("encode resourceModifyRequestTransfer failed: %s", err)
	} else {
		return buf, nil
	}
}

// TS 38.413 9.3.4.9
func BuildPathSwitchRequestAcknowledgeTransfer(ctx *SMContext) ([]byte, error) {
	ANUPF := ctx.Tunnel.DataPathPool.GetDefaultPath().FirstDPNode
	UpNode := ANUPF.UPF
	teidOct := make([]byte, 4)
	binary.BigEndian.PutUint32(teidOct, ANUPF.UpLinkTunnel.TEID)

	pathSwitchRequestAcknowledgeTransfer := ngapType.PathSwitchRequestAcknowledgeTransfer{}

	// UL NG-U UP TNL Information(optional) TS 38.413 9.3.2.2
	pathSwitchRequestAcknowledgeTransfer.
		ULNGUUPTNLInformation = new(ngapType.UPTransportLayerInformation)

	ULNGUUPTNLInformation := pathSwitchRequestAcknowledgeTransfer.ULNGUUPTNLInformation
	ULNGUUPTNLInformation.Present = ngapType.UPTransportLayerInformationPresentGTPTunnel
	ULNGUUPTNLInformation.GTPTunnel = new(ngapType.GTPTunnel)

	if n3IP, err := UpNode.N3Interfaces[0].IP(ctx.SelectedPDUSessionType); err != nil {
		return nil, err
	} else {
		gtpTunnel := ULNGUUPTNLInformation.GTPTunnel
		gtpTunnel.GTPTEID.Value = teidOct
		gtpTunnel.TransportLayerAddress.Value = aper.BitString{
			Bytes:     n3IP,
			BitLength: uint64(len(n3IP) * 8),
		}
	}

	// Received UP security policy mismatch from SMF locally stored TS 33.501 6.6.1
	// Security Indication(optional) TS 38.413 9.3.1.27
	if !ctx.UpSecurityFromPathSwitchRequestSameAsLocalStored {
		pathSwitchRequestAcknowledgeTransfer.SecurityIndication = new(ngapType.SecurityIndication)
		securityIndication := pathSwitchRequestAcknowledgeTransfer.SecurityIndication

		upSecurity := ctx.UpSecurity
		maximumDataRatePerUEForUserPlaneIntegrityProtectionForUpLink := ctx.
			MaximumDataRatePerUEForUserPlaneIntegrityProtectionForUpLink

		switch upSecurity.UpIntegr {
		case models.UpIntegrity_REQUIRED:
			securityIndication.IntegrityProtectionIndication.Value = ngapType.
				IntegrityProtectionIndicationPresentRequired
		case models.UpIntegrity_PREFERRED:
			securityIndication.IntegrityProtectionIndication.Value = ngapType.
				IntegrityProtectionIndicationPresentPreferred
		case models.UpIntegrity_NOT_NEEDED:
			securityIndication.IntegrityProtectionIndication.Value = ngapType.
				IntegrityProtectionIndicationPresentNotNeeded
		}
		switch upSecurity.UpConfid {
		case models.UpConfidentiality_REQUIRED:
			securityIndication.ConfidentialityProtectionIndication.Value = ngapType.
				ConfidentialityProtectionIndicationPresentRequired
		case models.UpConfidentiality_PREFERRED:
			securityIndication.ConfidentialityProtectionIndication.Value = ngapType.
				ConfidentialityProtectionIndicationPresentPreferred
		case models.UpConfidentiality_NOT_NEEDED:
			securityIndication.ConfidentialityProtectionIndication.Value = ngapType.
				ConfidentialityProtectionIndicationPresentNotNeeded
		}
		// Present only when Integrity Indication within the
		// Security Indication is set to "required" or "preferred"
		integrityProtectionInd := securityIndication.
			IntegrityProtectionIndication.Value
		if integrityProtectionInd == ngapType.IntegrityProtectionIndicationPresentRequired ||
			integrityProtectionInd == ngapType.
				IntegrityProtectionIndicationPresentPreferred {
			securityIndication.MaximumIntegrityProtectedDataRateUL = new(ngapType.MaximumIntegrityProtectedDataRate)
			switch maximumDataRatePerUEForUserPlaneIntegrityProtectionForUpLink {
			case models.MaxIntegrityProtectedDataRate_MAX_UE_RATE:
				securityIndication.MaximumIntegrityProtectedDataRateUL.Value = ngapType.
					MaximumIntegrityProtectedDataRatePresentMaximumUERate
			case models.MaxIntegrityProtectedDataRate__64_KBPS:
				securityIndication.MaximumIntegrityProtectedDataRateUL.Value = ngapType.
					MaximumIntegrityProtectedDataRatePresentBitrate64kbs
			}
		}
	}
	if buf, err := aper.MarshalWithParams(pathSwitchRequestAcknowledgeTransfer, "valueExt"); err != nil {
		return nil, err
	} else {
		return buf, nil
	}
}

func BuildPathSwitchRequestUnsuccessfulTransfer(causePresent int, causeValue aper.Enumerated) (buf []byte, err error) {
	pathSwitchRequestUnsuccessfulTransfer := ngapType.PathSwitchRequestUnsuccessfulTransfer{}

	pathSwitchRequestUnsuccessfulTransfer.Cause.Present = causePresent
	cause := &pathSwitchRequestUnsuccessfulTransfer.Cause

	switch causePresent {
	case ngapType.CausePresentRadioNetwork:
		cause.RadioNetwork = new(ngapType.CauseRadioNetwork)
		cause.RadioNetwork.Value = causeValue
	case ngapType.CausePresentTransport:
		cause.Transport = new(ngapType.CauseTransport)
		cause.Transport.Value = causeValue
	case ngapType.CausePresentNas:
		cause.Nas = new(ngapType.CauseNas)
		cause.Nas.Value = causeValue
	case ngapType.CausePresentProtocol:
		cause.Protocol = new(ngapType.CauseProtocol)
		cause.Protocol.Value = causeValue
	case ngapType.CausePresentMisc:
		cause.Misc = new(ngapType.CauseMisc)
		cause.Misc.Value = causeValue
	}

	buf, err = aper.MarshalWithParams(pathSwitchRequestUnsuccessfulTransfer, "valueExt")
	if err != nil {
		return nil, err
	}
	return
}

func BuildPDUSessionResourceReleaseCommandTransfer(ctx *SMContext) (buf []byte, err error) {
	resourceReleaseCommandTransfer := ngapType.PDUSessionResourceReleaseCommandTransfer{
		Cause: ngapType.Cause{
			Present: ngapType.CausePresentNas,
			Nas: &ngapType.CauseNas{
				Value: ngapType.CauseNasPresentNormalRelease,
			},
		},
	}
	buf, err = aper.MarshalWithParams(resourceReleaseCommandTransfer, "valueExt")
	if err != nil {
		return nil, err
	}
	return
}

func BuildHandoverCommandTransfer(ctx *SMContext) ([]byte, error) {
	handoverCommandTransfer := ngapType.HandoverCommandTransfer{}

	switch ctx.DLForwardingType {
	case IndirectForwarding:
		ANUPF := ctx.Tunnel.DataPathPool.GetDefaultPath().FirstDPNode
		UpNode := ANUPF.UPF
		teidOct := make([]byte, 4)
		binary.BigEndian.PutUint32(teidOct, ctx.IndirectForwardingTunnel.FirstDPNode.UpLinkTunnel.TEID)

		handoverCommandTransfer.DLForwardingUPTNLInformation = new(ngapType.UPTransportLayerInformation)
		handoverCommandTransfer.DLForwardingUPTNLInformation.Present = ngapType.UPTransportLayerInformationPresentGTPTunnel
		handoverCommandTransfer.DLForwardingUPTNLInformation.GTPTunnel = new(ngapType.GTPTunnel)

		if n3IP, err := UpNode.N3Interfaces[0].IP(ctx.SelectedPDUSessionType); err != nil {
			return nil, err
		} else {
			gtpTunnel := handoverCommandTransfer.DLForwardingUPTNLInformation.GTPTunnel
			gtpTunnel.GTPTEID.Value = teidOct
			gtpTunnel.TransportLayerAddress.Value = aper.BitString{
				Bytes:     n3IP,
				BitLength: uint64(len(n3IP) * 8),
			}
		}
	case DirectForwarding:
		handoverCommandTransfer.DLForwardingUPTNLInformation = ctx.DLDirectForwardingTunnel
	}

	handoverCommandTransfer.QosFlowToBeForwardedList = &ngapType.QosFlowToBeForwardedList{
		List: []ngapType.QosFlowToBeForwardedItem{
			{
				QosFlowIdentifier: ngapType.QosFlowIdentifier{
					Value: DefaultNonGBR5QI,
				},
			},
		},
	}

	if buf, err := aper.MarshalWithParams(handoverCommandTransfer, "valueExt"); err != nil {
		return nil, err
	} else {
		return buf, nil
	}
}
