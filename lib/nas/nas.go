//go:binary-only-package

package nas

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"free5gc/lib/nas/nasMessage"
)

// Message TODO：description
type Message struct {
	SecurityHeader
	*GmmMessage
	*GsmMessage
}

// SecurityHeader TODO：description
type SecurityHeader struct {
	ProtocolDiscriminator     uint8
	SecurityHeaderType        uint8
	MessageAuthenticationCode uint32
	SequenceNumber            uint8
}

const (
	SecurityHeaderTypePlainNas                                                 uint8 = 0x00
	SecurityHeaderTypeIntegrityProtected                                       uint8 = 0x01
	SecurityHeaderTypeIntegrityProtectedAndCiphered                            uint8 = 0x02
	SecurityHeaderTypeIntegrityProtectedWithNew5gNasSecurityContext            uint8 = 0x03
	SecurityHeaderTypeIntegrityProtectedAndCipheredWithNew5gNasSecurityContext uint8 = 0x04
)

// NewMessage TODO:desc
func NewMessage() *Message {}

// NewGmmMessage TODO:desc
func NewGmmMessage() *GmmMessage {}

// NewGmmMessage TODO:desc
func NewGsmMessage() *GsmMessage {}

// GmmHeader Octet1 protocolDiscriminator securityHeaderType
//           Octet2 MessageType
type GmmHeader struct {
	Octet [3]uint8
}

type GsmHeader struct {
	Octet [4]uint8
}

// GetMessageType 9.8
func (a *GmmHeader) GetMessageType() (messageType uint8) {}

// GetMessageType 9.8
func (a *GmmHeader) SetMessageType(messageType uint8) {}

func (a *GmmHeader) GetExtendedProtocolDiscriminator() uint8 {}

func (a *GmmHeader) SetExtendedProtocolDiscriminator(epd uint8) {}

func (a *GsmHeader) GetExtendedProtocolDiscriminator() uint8 {}

func (a *GsmHeader) SetExtendedProtocolDiscriminator(epd uint8) {}

// GetMessageType 9.8
func (a *GsmHeader) GetMessageType() (messageType uint8) {}

// GetMessageType 9.8
func (a *GsmHeader) SetMessageType(messageType uint8) {}

func GetEPD(byteArray []byte) uint8 {}

func GetSecurityHeaderType(byteArray []byte) uint8 {}

type GmmMessage struct {
	GmmHeader
	*nasMessage.AuthenticationRequest                            //8.2.1
	*nasMessage.AuthenticationResponse                           //8.2.2
	*nasMessage.AuthenticationResult                             //8.2.3
	*nasMessage.AuthenticationFailure                            //8.2.4
	*nasMessage.AuthenticationReject                             //8.2.5
	*nasMessage.RegistrationRequest                              //8.2.6
	*nasMessage.RegistrationAccept                               //8.2.7
	*nasMessage.RegistrationComplete                             //8.2.8
	*nasMessage.RegistrationReject                               //8.2.9
	*nasMessage.ULNASTransport                                   //8.2.10
	*nasMessage.DLNASTransport                                   //8.2.11
	*nasMessage.DeregistrationRequestUEOriginatingDeregistration //8.2.12
	*nasMessage.DeregistrationAcceptUEOriginatingDeregistration  //8.2.13
	*nasMessage.DeregistrationRequestUETerminatedDeregistration  //8.2.14
	*nasMessage.DeregistrationAcceptUETerminatedDeregistration   //8.2.15
	*nasMessage.ServiceRequest                                   //8.2.16
	*nasMessage.ServiceAccept                                    //8.2.17
	*nasMessage.ServiceReject                                    //8.2.18
	*nasMessage.ConfigurationUpdateCommand                       //8.2.19
	*nasMessage.ConfigurationUpdateComplete                      //8.2.20
	*nasMessage.IdentityRequest                                  //8.2.21
	*nasMessage.IdentityResponse                                 //8.2.22
	*nasMessage.Notification                                     //8.2.23
	*nasMessage.NotificationResponse                             //8.2.24
	*nasMessage.SecurityModeCommand                              //8.2.25
	*nasMessage.SecurityModeComplete                             //8.2.26
	*nasMessage.SecurityModeReject                               //8.2.27
	*nasMessage.SecurityProtected5GSNASMessage                   //8.2.28
	*nasMessage.Status5GMM                                       //8.2.29
}

const (
	MsgTypeRegistrationRequest                              uint8 = 65
	MsgTypeRegistrationAccept                               uint8 = 66
	MsgTypeRegistrationComplete                             uint8 = 67
	MsgTypeRegistrationReject                               uint8 = 68
	MsgTypeDeregistrationRequestUEOriginatingDeregistration uint8 = 69
	MsgTypeDeregistrationAcceptUEOriginatingDeregistration  uint8 = 70
	MsgTypeDeregistrationRequestUETerminatedDeregistration  uint8 = 71
	MsgTypeDeregistrationAcceptUETerminatedDeregistration   uint8 = 72
	MsgTypeServiceRequest                                   uint8 = 76
	MsgTypeServiceReject                                    uint8 = 77
	MsgTypeServiceAccept                                    uint8 = 78
	MsgTypeConfigurationUpdateCommand                       uint8 = 84
	MsgTypeConfigurationUpdateComplete                      uint8 = 85
	MsgTypeAuthenticationRequest                            uint8 = 86
	MsgTypeAuthenticationResponse                           uint8 = 87
	MsgTypeAuthenticationReject                             uint8 = 88
	MsgTypeAuthenticationFailure                            uint8 = 89
	MsgTypeAuthenticationResult                             uint8 = 90
	MsgTypeIdentityRequest                                  uint8 = 91
	MsgTypeIdentityResponse                                 uint8 = 92
	MsgTypeSecurityModeCommand                              uint8 = 93
	MsgTypeSecurityModeComplete                             uint8 = 94
	MsgTypeSecurityModeReject                               uint8 = 95
	MsgTypeStatus5GMM                                       uint8 = 100
	MsgTypeNotification                                     uint8 = 101
	MsgTypeNotificationResponse                             uint8 = 102
	MsgTypeULNASTransport                                   uint8 = 103
	MsgTypeDLNASTransport                                   uint8 = 104
)

func (a *Message) PlainNasDecode(byteArray *[]byte) error {}
func (a *Message) PlainNasEncode() ([]byte, error) {}

func (a *Message) GmmMessageDecode(byteArray *[]byte) error {}

func (a *Message) GmmMessageEncode(buffer *bytes.Buffer) error {}

type GsmMessage struct {
	GsmHeader
	*nasMessage.PDUSessionEstablishmentRequest      //8.3.1
	*nasMessage.PDUSessionEstablishmentAccept       //8.3.2
	*nasMessage.PDUSessionEstablishmentReject       //8.3.3
	*nasMessage.PDUSessionAuthenticationCommand     //8.3.4
	*nasMessage.PDUSessionAuthenticationComplete    //8.3.5
	*nasMessage.PDUSessionAuthenticationResult      //8.3.6
	*nasMessage.PDUSessionModificationRequest       //8.3.7
	*nasMessage.PDUSessionModificationReject        //8.3.8
	*nasMessage.PDUSessionModificationCommand       //8.3.9
	*nasMessage.PDUSessionModificationComplete      //8.3.10
	*nasMessage.PDUSessionModificationCommandReject //8.3.11
	*nasMessage.PDUSessionReleaseRequest            //8.3.12
	*nasMessage.PDUSessionReleaseReject             //8.3.13
	*nasMessage.PDUSessionReleaseCommand            //8.3.14
	*nasMessage.PDUSessionReleaseComplete           //8.3.15
	*nasMessage.Status5GSM                          //8.3.16
}

const (
	MsgTypePDUSessionEstablishmentRequest      uint8 = 193
	MsgTypePDUSessionEstablishmentAccept       uint8 = 194
	MsgTypePDUSessionEstablishmentReject       uint8 = 195
	MsgTypePDUSessionAuthenticationCommand     uint8 = 197
	MsgTypePDUSessionAuthenticationComplete    uint8 = 198
	MsgTypePDUSessionAuthenticationResult      uint8 = 199
	MsgTypePDUSessionModificationRequest       uint8 = 201
	MsgTypePDUSessionModificationReject        uint8 = 202
	MsgTypePDUSessionModificationCommand       uint8 = 203
	MsgTypePDUSessionModificationComplete      uint8 = 204
	MsgTypePDUSessionModificationCommandReject uint8 = 205
	MsgTypePDUSessionReleaseRequest            uint8 = 209
	MsgTypePDUSessionReleaseReject             uint8 = 210
	MsgTypePDUSessionReleaseCommand            uint8 = 211
	MsgTypePDUSessionReleaseComplete           uint8 = 212
	MsgTypeStatus5GSM                          uint8 = 214
)

func (a *Message) GsmMessageDecode(byteArray *[]byte) error {}

func (a *Message) GsmMessageEncode(buffer *bytes.Buffer) error {}
