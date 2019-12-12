//go:binary-only-package

package nasTestpacket

import (
	"bytes"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"free5gc/lib/nas"
	"free5gc/lib/nas/nasMessage"
	"free5gc/lib/nas/nasType"
	"free5gc/lib/openapi/models"
)

const (
	PDUSesModiReq    string = "PDU Session Modification Request"
	PDUSesModiCmp    string = "PDU Session Modification Complete"
	PDUSesModiCmdRej string = "PDU Session Modification Command Reject"
	PDUSesRelReq     string = "PDU Session Release Request"
	PDUSesRelCmp     string = "PDU Session Release Complete"
	PDUSesRelRej     string = "PDU Session Release Reject"
	PDUSesAuthCmp    string = "PDU Session Authentication Complete"
)

func GetRegistrationRequest(registrationType uint8, mobileIdentity nasType.MobileIdentity5GS, requestedNSSAI *nasType.RequestedNSSAI, uplinkDataStatus *nasType.UplinkDataStatus) (nasPdu []byte) {}

func GetPduSessionEstablishmentRequest(pduSessionId uint8) (nasPdu []byte) {}

func GetUlNasTransport_PduSessionEstablishmentRequest(pduSessionId uint8, requestType uint8, dnnString string, sNssai *models.Snssai) (nasPdu []byte) {}

func GetUlNasTransport_PduSessionModificationRequest(pduSessionId uint8, requestType uint8, dnnString string, sNssai *models.Snssai) (nasPdu []byte) {}

func GetPduSessionModificationRequest(pduSessionId uint8) (nasPdu []byte) {}
func GetPduSessionModificationComplete(pduSessionId uint8) (nasPdu []byte) {}
func GetPduSessionModificationCommandReject(pduSessionId uint8) (nasPdu []byte) {}

func GetPduSessionReleaseRequest(pduSessionId uint8) (nasPdu []byte) {}

func GetPduSessionReleaseComplete(pduSessionId uint8) (nasPdu []byte) {}

func GetPduSessionReleaseReject(pduSessionId uint8) (nasPdu []byte) {}

func GetPduSessionReleaseCommand(pduSessionId uint8) (nasPdu []byte) {}

func GetPduSessionAuthenticationComplete(pduSessionId uint8) (nasPdu []byte) {}

func GetUlNasTransport_PduSessionCommonData(pduSessionId uint8, types string) (nasPdu []byte) {}

func GetIdentityResponse(mobileIdentity nasType.MobileIdentity) (nasPdu []byte) {}

func GetNotificationResponse(pDUSessionStatus []uint8) (nasPdu []byte) {}

func GetConfigurationUpdateComplete() (nasPdu []byte) {}

func GetServiceRequest(serviceType uint8) (nasPdu []byte) {}

func GetAuthenticationResponse(authenticationResponseParam []uint8, eapMsg string) (nasPdu []byte) {}

func GetAuthenticationFailure(cause5GMM uint8, authenticationFailureParam []uint8) (nasPdu []byte) {}

func GetRegistrationComplete(sorTransparentContainer []uint8) (nasPdu []byte) {}

// TODO: finish it; TS 24.501 8.2.26
func GetSecurityModeComplete() (nasPdu []byte) {}

func GetSecurityModeReject(cause5GMM uint8) (nasPdu []byte) {}

func GetDeregistrationRequest(accessType uint8, switchOff uint8, ngKsi uint8, mobileIdentity5GS nasType.MobileIdentity5GS) (nasPdu []byte) {}

func GetDeregistrationAccept() (nasPdu []byte) {}

func GetStatus5GMM(cause uint8) (nasPdu []byte) {}

func GetStatus5GSM(pduSessionId uint8, cause uint8) (nasPdu []byte) {}

func GetUlNasTransport_Status5GSM(pduSessionId uint8, cause uint8) (nasPdu []byte) {}

func GetUlNasTransport_PduSessionReleaseRequest(pduSessionId uint8) (nasPdu []byte) {}

func GetUlNasTransport_PduSessionReleaseCommand(pduSessionId uint8, requestType uint8, dnnString string, sNssai *models.Snssai) (nasPdu []byte) {}
