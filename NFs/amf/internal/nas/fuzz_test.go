//go:build go1.18
// +build go1.18

package nas_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	amf_context "github.com/free5gc/amf/internal/context"
	"github.com/free5gc/amf/internal/logger"
	amf_nas "github.com/free5gc/amf/internal/nas"
	"github.com/free5gc/amf/internal/sbi/consumer"
	"github.com/free5gc/amf/pkg/service"
	"github.com/free5gc/nas"
	"github.com/free5gc/nas/nasMessage"
	"github.com/free5gc/nas/nasType"
	"github.com/free5gc/ngap/ngapType"
	"github.com/free5gc/openapi/models"
)

func FuzzHandleNAS(f *testing.F) {
	amfSelf := amf_context.GetSelf()
	amfSelf.ServedGuamiList = []models.Guami{{
		PlmnId: &models.PlmnIdNid{
			Mcc: "208",
			Mnc: "93",
		},
		AmfId: "cafe00",
	}}
	tai := models.Tai{
		PlmnId: &models.PlmnId{
			Mcc: "208",
			Mnc: "93",
		},
		Tac: "1",
	}
	amfSelf.SupportTaiLists = []models.Tai{tai}

	msg := nas.NewMessage()
	msg.GmmMessage = nas.NewGmmMessage()
	msg.GmmMessage.GmmHeader.SetMessageType(nas.MsgTypeRegistrationRequest)
	msg.GmmMessage.RegistrationRequest = nasMessage.NewRegistrationRequest(nas.MsgTypeRegistrationRequest)
	reg := msg.GmmMessage.RegistrationRequest
	reg.ExtendedProtocolDiscriminator.SetExtendedProtocolDiscriminator(nasMessage.Epd5GSMobilityManagementMessage)
	reg.SpareHalfOctetAndSecurityHeaderType.SetSecurityHeaderType(nas.SecurityHeaderTypePlainNas)
	reg.RegistrationRequestMessageIdentity.SetMessageType(nas.MsgTypeRegistrationRequest)
	reg.NgksiAndRegistrationType5GS.SetTSC(nasMessage.TypeOfSecurityContextFlagNative)
	reg.NgksiAndRegistrationType5GS.SetNasKeySetIdentifiler(7)
	reg.NgksiAndRegistrationType5GS.SetFOR(1)
	reg.NgksiAndRegistrationType5GS.SetRegistrationType5GS(nasMessage.RegistrationType5GSInitialRegistration)
	id := []uint8{0x01, 0x02, 0xf8, 0x39, 0xf0, 0xff, 0x00, 0x00, 0x00, 0x00, 0x00, 0x10}
	reg.MobileIdentity5GS.SetLen(uint16(len(id)))
	reg.MobileIdentity5GS.SetMobileIdentity5GSContents(id)
	reg.UESecurityCapability = nasType.NewUESecurityCapability(nasMessage.RegistrationRequestUESecurityCapabilityType)
	reg.UESecurityCapability.SetLen(2)
	reg.UESecurityCapability.SetEA0_5G(1)
	reg.UESecurityCapability.SetIA2_128_5G(1)
	buf, err := msg.PlainNasEncode()
	require.NoError(f, err)
	f.Add(buf)

	msg = nas.NewMessage()
	msg.GmmMessage = nas.NewGmmMessage()
	msg.GmmMessage.GmmHeader.SetMessageType(nas.MsgTypeDeregistrationRequestUEOriginatingDeregistration)
	deReg := nasMessage.NewDeregistrationRequestUEOriginatingDeregistration(
		nas.MsgTypeDeregistrationRequestUEOriginatingDeregistration)
	msg.GmmMessage.DeregistrationRequestUEOriginatingDeregistration = deReg
	deReg.ExtendedProtocolDiscriminator.SetExtendedProtocolDiscriminator(nasMessage.Epd5GSMobilityManagementMessage)
	deReg.SpareHalfOctetAndSecurityHeaderType.SetSecurityHeaderType(nas.SecurityHeaderTypePlainNas)
	deReg.DeregistrationRequestMessageIdentity.SetMessageType(nas.MsgTypeDeregistrationRequestUEOriginatingDeregistration)
	deReg.NgksiAndDeregistrationType.SetTSC(nasMessage.TypeOfSecurityContextFlagNative)
	deReg.NgksiAndDeregistrationType.SetNasKeySetIdentifiler(7)
	deReg.NgksiAndDeregistrationType.SetSwitchOff(0)
	deReg.NgksiAndDeregistrationType.SetAccessType(nasMessage.AccessType3GPP)
	deReg.MobileIdentity5GS.SetLen(uint16(len(id)))
	deReg.MobileIdentity5GS.SetMobileIdentity5GSContents(id)
	buf, err = msg.PlainNasEncode()
	require.NoError(f, err)
	f.Add(buf)

	msg = nas.NewMessage()
	msg.GmmMessage = nas.NewGmmMessage()
	msg.GmmMessage.GmmHeader.SetMessageType(nas.MsgTypeServiceRequest)
	msg.GmmMessage.ServiceRequest = nasMessage.NewServiceRequest(nas.MsgTypeServiceRequest)
	sr := msg.GmmMessage.ServiceRequest
	sr.ExtendedProtocolDiscriminator.SetExtendedProtocolDiscriminator(nasMessage.Epd5GSMobilityManagementMessage)
	sr.SpareHalfOctetAndSecurityHeaderType.SetSecurityHeaderType(nas.SecurityHeaderTypePlainNas)
	sr.ServiceRequestMessageIdentity.SetMessageType(nas.MsgTypeServiceRequest)
	sr.ServiceTypeAndNgksi.SetTSC(nasMessage.TypeOfSecurityContextFlagNative)
	sr.ServiceTypeAndNgksi.SetNasKeySetIdentifiler(0)
	sr.ServiceTypeAndNgksi.SetServiceTypeValue(nasMessage.ServiceTypeSignalling)
	sr.TMSI5GS.SetLen(7)
	buf, err = msg.PlainNasEncode()
	require.NoError(f, err)
	buf = append([]uint8{
		nasMessage.Epd5GSMobilityManagementMessage,
		nas.SecurityHeaderTypeIntegrityProtected,
		0, 0, 0, 0, 0,
	},
		buf...)
	f.Add(buf)

	f.Fuzz(func(t *testing.T, d []byte) {
		ue := new(amf_context.RanUe)
		ue.Ran = new(amf_context.AmfRan)
		ue.Ran.AnType = models.AccessType__3_GPP_ACCESS
		ue.Ran.Log = logger.NgapLog
		ue.Log = logger.NgapLog
		ue.Tai = tai
		ue.AmfUe = amfSelf.NewAmfUe("")
		amf_nas.HandleNAS(ue, ngapType.ProcedureCodeInitialUEMessage, d, true)
	})
}

func FuzzHandleNAS2(f *testing.F) {
	amfSelf := amf_context.GetSelf()
	amfSelf.ServedGuamiList = []models.Guami{{
		PlmnId: &models.PlmnIdNid{
			Mcc: "208",
			Mnc: "93",
		},
		AmfId: "cafe00",
	}}
	tai := models.Tai{
		PlmnId: &models.PlmnId{
			Mcc: "208",
			Mnc: "93",
		},
		Tac: "1",
	}
	amfSelf.SupportTaiLists = []models.Tai{tai}
	amfSelf.NrfUri = "test"

	msg := nas.NewMessage()
	msg.GmmMessage = nas.NewGmmMessage()
	msg.GmmMessage.GmmHeader.SetMessageType(nas.MsgTypeRegistrationRequest)
	msg.GmmMessage.RegistrationRequest = nasMessage.NewRegistrationRequest(nas.MsgTypeRegistrationRequest)
	reg := msg.GmmMessage.RegistrationRequest
	reg.ExtendedProtocolDiscriminator.SetExtendedProtocolDiscriminator(nasMessage.Epd5GSMobilityManagementMessage)
	reg.SpareHalfOctetAndSecurityHeaderType.SetSecurityHeaderType(nas.SecurityHeaderTypePlainNas)
	reg.RegistrationRequestMessageIdentity.SetMessageType(nas.MsgTypeRegistrationRequest)
	reg.NgksiAndRegistrationType5GS.SetTSC(nasMessage.TypeOfSecurityContextFlagNative)
	reg.NgksiAndRegistrationType5GS.SetNasKeySetIdentifiler(7)
	reg.NgksiAndRegistrationType5GS.SetFOR(1)
	reg.NgksiAndRegistrationType5GS.SetRegistrationType5GS(nasMessage.RegistrationType5GSInitialRegistration)
	id := []uint8{0x01, 0x02, 0xf8, 0x39, 0xf0, 0xff, 0x00, 0x00, 0x00, 0x00, 0x00, 0x10}
	reg.MobileIdentity5GS.SetLen(uint16(len(id)))
	reg.MobileIdentity5GS.SetMobileIdentity5GSContents(id)
	reg.UESecurityCapability = nasType.NewUESecurityCapability(nasMessage.RegistrationRequestUESecurityCapabilityType)
	reg.UESecurityCapability.SetLen(2)
	reg.UESecurityCapability.SetEA0_5G(1)
	reg.UESecurityCapability.SetIA2_128_5G(1)
	regPkt, err := msg.PlainNasEncode()
	require.NoError(f, err)

	msg = nas.NewMessage()
	msg.GmmMessage = nas.NewGmmMessage()
	msg.GmmMessage.GmmHeader.SetMessageType(nas.MsgTypeIdentityResponse)
	msg.GmmMessage.IdentityResponse = nasMessage.NewIdentityResponse(nas.MsgTypeIdentityResponse)
	ir := msg.GmmMessage.IdentityResponse
	ir.ExtendedProtocolDiscriminator.SetExtendedProtocolDiscriminator(nasMessage.Epd5GSMobilityManagementMessage)
	ir.SpareHalfOctetAndSecurityHeaderType.SetSecurityHeaderType(nas.SecurityHeaderTypePlainNas)
	ir.IdentityResponseMessageIdentity.SetMessageType(nas.MsgTypeIdentityResponse)
	ir.MobileIdentity.SetLen(uint16(len(id)))
	ir.MobileIdentity.SetMobileIdentityContents(id)
	buf, err := msg.PlainNasEncode()
	require.NoError(f, err)
	f.Add(buf)

	msg = nas.NewMessage()
	msg.GmmMessage = nas.NewGmmMessage()
	msg.GmmMessage.GmmHeader.SetMessageType(nas.MsgTypeAuthenticationResponse)
	msg.GmmMessage.AuthenticationResponse = nasMessage.NewAuthenticationResponse(nas.MsgTypeAuthenticationResponse)
	ar := msg.GmmMessage.AuthenticationResponse
	ar.ExtendedProtocolDiscriminator.SetExtendedProtocolDiscriminator(nasMessage.Epd5GSMobilityManagementMessage)
	ar.SpareHalfOctetAndSecurityHeaderType.SetSecurityHeaderType(nas.SecurityHeaderTypePlainNas)
	ar.AuthenticationResponseMessageIdentity.SetMessageType(nas.MsgTypeAuthenticationResponse)
	ar.AuthenticationResponseParameter = nasType.NewAuthenticationResponseParameter(
		nasMessage.AuthenticationResponseAuthenticationResponseParameterType)
	ar.AuthenticationResponseParameter.SetLen(16)
	buf, err = msg.PlainNasEncode()
	require.NoError(f, err)
	f.Add(buf)

	msg = nas.NewMessage()
	msg.GmmMessage = nas.NewGmmMessage()
	msg.GmmMessage.GmmHeader.SetMessageType(nas.MsgTypeAuthenticationFailure)
	msg.GmmMessage.AuthenticationFailure = nasMessage.NewAuthenticationFailure(nas.MsgTypeAuthenticationFailure)
	af := msg.GmmMessage.AuthenticationFailure
	af.ExtendedProtocolDiscriminator.SetExtendedProtocolDiscriminator(nasMessage.Epd5GSMobilityManagementMessage)
	af.SpareHalfOctetAndSecurityHeaderType.SetSecurityHeaderType(nas.SecurityHeaderTypePlainNas)
	af.AuthenticationFailureMessageIdentity.SetMessageType(nas.MsgTypeAuthenticationFailure)
	af.Cause5GMM.SetCauseValue(nasMessage.Cause5GMMSynchFailure)
	af.AuthenticationFailureParameter = nasType.NewAuthenticationFailureParameter(
		nasMessage.AuthenticationFailureAuthenticationFailureParameterType)
	af.AuthenticationFailureParameter.SetLen(14)
	buf, err = msg.PlainNasEncode()
	require.NoError(f, err)
	f.Add(buf)

	msg = nas.NewMessage()
	msg.GmmMessage = nas.NewGmmMessage()
	msg.GmmMessage.GmmHeader.SetMessageType(nas.MsgTypeStatus5GMM)
	msg.GmmMessage.Status5GMM = nasMessage.NewStatus5GMM(nas.MsgTypeStatus5GMM)
	st := msg.GmmMessage.Status5GMM
	st.ExtendedProtocolDiscriminator.SetExtendedProtocolDiscriminator(nasMessage.Epd5GSMobilityManagementMessage)
	st.SpareHalfOctetAndSecurityHeaderType.SetSecurityHeaderType(nas.SecurityHeaderTypePlainNas)
	st.STATUSMessageIdentity5GMM.SetMessageType(nas.MsgTypeStatus5GMM)
	st.Cause5GMM.SetCauseValue(nasMessage.Cause5GMMProtocolErrorUnspecified)
	buf, err = msg.PlainNasEncode()
	require.NoError(f, err)
	f.Add(buf)

	f.Fuzz(func(t *testing.T, d []byte) {
		ctrl := gomock.NewController(t)
		// m := app.NewMockApp(ctrl)
		m := service.NewMockAmfAppInterface(ctrl)
		c, errc := consumer.NewConsumer(m)
		service.AMF = m
		require.NoError(t, errc)
		m.EXPECT().
			Consumer().
			AnyTimes().
			Return(c)

		ue := new(amf_context.RanUe)
		ue.Ran = new(amf_context.AmfRan)
		ue.Ran.AnType = models.AccessType__3_GPP_ACCESS
		ue.Ran.Log = logger.NgapLog
		ue.Log = logger.NgapLog
		ue.Tai = tai
		ue.AmfUe = amfSelf.NewAmfUe("")
		amf_nas.HandleNAS(ue, ngapType.ProcedureCodeInitialUEMessage, regPkt, true)
		amfUe := ue.AmfUe
		amfUe.State[models.AccessType__3_GPP_ACCESS].Set(amf_context.Authentication)
		amfUe.RequestIdentityType = nasMessage.MobileIdentity5GSTypeSuci
		amfUe.AuthenticationCtx = &models.UeAuthenticationCtx{
			AuthType: models.AusfUeAuthenticationAuthType__5_G_AKA,
		}
		amf_nas.HandleNAS(ue, ngapType.ProcedureCodeUplinkNASTransport, d, false)
	})
}
