package ike

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1" // #nosec G505
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"math"
	"net"
	"runtime/debug"
	"sync/atomic"
	"time"

	"github.com/pkg/errors"
	"github.com/vishvananda/netlink"
	"golang.org/x/sys/unix"

	ike_message "github.com/free5gc/ike/message"
	ike_security "github.com/free5gc/ike/security"
	"github.com/free5gc/ike/security/dh"
	"github.com/free5gc/ike/security/encr"
	"github.com/free5gc/ike/security/integ"
	"github.com/free5gc/ike/security/prf"
	n3iwf_context "github.com/free5gc/n3iwf/internal/context"
	"github.com/free5gc/n3iwf/internal/ike/xfrm"
	"github.com/free5gc/n3iwf/internal/logger"
)

func (s *Server) HandleIKESAINIT(
	udpConn *net.UDPConn,
	n3iwfAddr, ueAddr *net.UDPAddr,
	message *ike_message.IKEMessage,
	realMessage1 []byte,
) {
	ikeLog := logger.IKELog
	ikeLog.Infoln("Handle IKE_SA_INIT")

	// Used to receive value from peer
	var securityAssociation *ike_message.SecurityAssociation
	var keyExcahge *ike_message.KeyExchange
	var nonce *ike_message.Nonce
	var notifications []*ike_message.Notification

	n3iwfCtx := s.Context()
	cfg := s.Config()

	// For response or needed data
	var responseIKEPayload ike_message.IKEPayloadContainer
	var localNonce, concatenatedNonce []byte
	// Chosen transform from peer's proposal
	var chooseProposal ike_message.ProposalContainer
	var localPublicValue []byte
	var chosenDiffieHellmanGroup uint16

	for _, ikePayload := range message.Payloads {
		switch ikePayload.Type() {
		case ike_message.TypeSA:
			securityAssociation = ikePayload.(*ike_message.SecurityAssociation)
		case ike_message.TypeKE:
			keyExcahge = ikePayload.(*ike_message.KeyExchange)
		case ike_message.TypeNiNr:
			nonce = ikePayload.(*ike_message.Nonce)
		case ike_message.TypeN:
			notifications = append(notifications, ikePayload.(*ike_message.Notification))
		default:
			ikeLog.Warnf(
				"Get IKE payload (type %d) in IKE_SA_INIT message, this payload will not be handled by IKE handler",
				ikePayload.Type())
		}
	}

	if securityAssociation != nil {
		responseSecurityAssociation := responseIKEPayload.BuildSecurityAssociation()
		chooseProposal = SelectProposal(securityAssociation.Proposals)
		responseSecurityAssociation.Proposals = append(responseSecurityAssociation.Proposals, chooseProposal...)

		if len(responseSecurityAssociation.Proposals) == 0 {
			ikeLog.Warn("No proposal chosen")
			// Respond NO_PROPOSAL_CHOSEN to UE
			responseIKEPayload.Reset()
			responseIKEPayload.BuildNotification(ike_message.TypeNone, ike_message.NO_PROPOSAL_CHOSEN, nil, nil)
			responseIKEMessage := ike_message.NewMessage(message.InitiatorSPI, message.ResponderSPI,
				ike_message.IKE_SA_INIT, true, false, message.MessageID, responseIKEPayload)

			err := SendIKEMessageToUE(udpConn, n3iwfAddr, ueAddr, responseIKEMessage, nil)
			if err != nil {
				ikeLog.Errorf("HandleIKESAINIT(): %v", err)
			}
			return
		}
	} else {
		ikeLog.Error("The security association field is nil")
		// TODO: send error message to UE
		return
	}

	if keyExcahge != nil {
		chosenDiffieHellmanGroup = chooseProposal[0].DiffieHellmanGroup[0].TransformID
		if chosenDiffieHellmanGroup != keyExcahge.DiffieHellmanGroup {
			ikeLog.Warn("The Diffie-Hellman group defined in key exchange payload not matches the one in chosen proposal")
			// send INVALID_KE_PAYLOAD to UE
			responseIKEPayload.Reset()

			notificationData := make([]byte, 2)
			binary.BigEndian.PutUint16(notificationData, chosenDiffieHellmanGroup)
			responseIKEPayload.BuildNotification(
				ike_message.TypeNone, ike_message.INVALID_KE_PAYLOAD, nil, notificationData)

			responseIKEMessage := ike_message.NewMessage(message.InitiatorSPI, message.ResponderSPI,
				ike_message.IKE_SA_INIT, true, false, message.MessageID, responseIKEPayload)

			err := SendIKEMessageToUE(udpConn, n3iwfAddr, ueAddr, responseIKEMessage, nil)
			if err != nil {
				ikeLog.Errorf("HandleIKESAINIT(): %v", err)
			}
			return
		}
	} else {
		ikeLog.Error("The key exchange field is nil")
		// TODO: send error message to UE
		return
	}

	if nonce != nil {
		localNonceBigInt, err := ike_security.GenerateRandomNumber()
		if err != nil {
			ikeLog.Errorf("HandleIKESAINIT: %v", err)
			return
		}
		localNonce = localNonceBigInt.Bytes()
		concatenatedNonce = append(nonce.NonceData, localNonce...)

		responseIKEPayload.BuildNonce(localNonce)
	} else {
		ikeLog.Error("The nonce field is nil")
		// TODO: send error message to UE
		return
	}

	ueBehindNAT, n3iwfBehindNAT, err := s.handleNATDetect(
		message.InitiatorSPI, message.ResponderSPI,
		notifications, ueAddr, n3iwfAddr)
	if err != nil {
		ikeLog.Errorf("Handle IKE_SA_INIT: %v", err)
		return
	}

	// Create new IKE security association
	ikeSecurityAssociation := n3iwfCtx.NewIKESecurityAssociation()
	ikeSecurityAssociation.RemoteSPI = message.InitiatorSPI
	ikeSecurityAssociation.InitiatorMessageID = message.MessageID

	ikeSecurityAssociation.IKESAKey, localPublicValue, err = ike_security.NewIKESAKey(chooseProposal[0],
		keyExcahge.KeyExchangeData, concatenatedNonce,
		ikeSecurityAssociation.RemoteSPI, ikeSecurityAssociation.LocalSPI)
	if err != nil {
		ikeLog.Errorf("Handle IKE_SA_INIT: %v", err)
		return
	}

	ikeLog.Debugln(ikeSecurityAssociation.String())

	// Record concatenated nonce
	ikeSecurityAssociation.ConcatenatedNonce = append(
		ikeSecurityAssociation.ConcatenatedNonce, concatenatedNonce...)
	ikeSecurityAssociation.UeBehindNAT = ueBehindNAT
	ikeSecurityAssociation.N3iwfBehindNAT = n3iwfBehindNAT

	responseIKEPayload.BUildKeyExchange(chosenDiffieHellmanGroup, localPublicValue)
	err = s.buildNATDetectNotifPayload(
		ikeSecurityAssociation, &responseIKEPayload, ueAddr, n3iwfAddr)
	if err != nil {
		ikeLog.Warnf("Handle IKE_SA_INIT: %v", err)
		return
	}

	// IKE response to UE
	responseIKEMessage := ike_message.NewMessage(message.InitiatorSPI, ikeSecurityAssociation.LocalSPI,
		ike_message.IKE_SA_INIT, true, false, message.MessageID, responseIKEPayload)

	// Prepare authentication data - InitatorSignedOctet
	// InitatorSignedOctet = RealMessage1 | NonceRData | MACedIDForI
	// MACedIDForI is acquired in IKE_AUTH exchange
	ikeSecurityAssociation.InitiatorSignedOctets = append(realMessage1, localNonce...)

	// Prepare authentication data - ResponderSignedOctet
	// ResponderSignedOctet = RealMessage2 | NonceIData | MACedIDForR
	responseIKEMessageData, err := responseIKEMessage.Encode()
	if err != nil {
		ikeLog.Errorf("Encoding IKE message failed: %v", err)
		return
	}
	ikeSecurityAssociation.ResponderSignedOctets = append(responseIKEMessageData, nonce.NonceData...)
	// MACedIDForR
	var idPayload ike_message.IKEPayloadContainer
	idPayload.BuildIdentificationResponder(ike_message.ID_FQDN, []byte(cfg.GetFQDN()))
	idPayloadData, err := idPayload.Encode()
	if err != nil {
		ikeLog.Errorf("Encode IKE payload failed: %v", err)
		return
	}

	ikeSecurityAssociation.Prf_r.Reset()
	_, err = ikeSecurityAssociation.Prf_r.Write(idPayloadData[4:])
	if err != nil {
		ikeLog.Errorf("Pseudorandom function write error: %v", err)
		return
	}

	ikeSecurityAssociation.ResponderSignedOctets = append(ikeSecurityAssociation.ResponderSignedOctets,
		ikeSecurityAssociation.Prf_r.Sum(nil)...)

	ikeLog.Tracef("Local unsigned authentication data:\n%s", hex.Dump(ikeSecurityAssociation.ResponderSignedOctets))

	// Send response to UE
	err = SendIKEMessageToUE(udpConn, n3iwfAddr, ueAddr, responseIKEMessage, nil)
	if err != nil {
		ikeLog.Errorf("HandleIKESAINIT(): %v", err)
	}
}

const (
	// IKE_AUTH state
	PreSignalling = iota
	EAPSignalling
	PostSignalling
	EndSignalling

	// CREATE_CHILDSA
	HandleCreateChildSA
)

func (s *Server) HandleIKEAUTH(
	udpConn *net.UDPConn,
	n3iwfAddr, ueAddr *net.UDPAddr,
	message *ike_message.IKEMessage,
	ikeSecurityAssociation *n3iwf_context.IKESecurityAssociation,
) {
	ikeLog := logger.IKELog
	ikeLog.Infoln("Handle IKE_AUTH")

	n3iwfCtx := s.Context()
	cfg := s.Config()
	ipsecGwAddr := cfg.GetIPSecGatewayAddr()

	// Used for response
	var responseIKEPayload ike_message.IKEPayloadContainer

	// Parse payloads
	var initiatorID *ike_message.IdentificationInitiator
	var certificateRequest *ike_message.CertificateRequest
	var certificate *ike_message.Certificate
	var securityAssociation *ike_message.SecurityAssociation
	var trafficSelectorInitiator *ike_message.TrafficSelectorInitiator
	var trafficSelectorResponder *ike_message.TrafficSelectorResponder
	var eap *ike_message.EAP
	var authentication *ike_message.Authentication
	var configuration *ike_message.Configuration
	var ok bool

	for _, ikePayload := range message.Payloads {
		switch ikePayload.Type() {
		case ike_message.TypeIDi:
			initiatorID = ikePayload.(*ike_message.IdentificationInitiator)
		case ike_message.TypeCERTreq:
			certificateRequest = ikePayload.(*ike_message.CertificateRequest)
		case ike_message.TypeCERT:
			certificate = ikePayload.(*ike_message.Certificate)
		case ike_message.TypeSA:
			securityAssociation = ikePayload.(*ike_message.SecurityAssociation)
		case ike_message.TypeTSi:
			trafficSelectorInitiator = ikePayload.(*ike_message.TrafficSelectorInitiator)
		case ike_message.TypeTSr:
			trafficSelectorResponder = ikePayload.(*ike_message.TrafficSelectorResponder)
		case ike_message.TypeEAP:
			eap = ikePayload.(*ike_message.EAP)
		case ike_message.TypeAUTH:
			authentication = ikePayload.(*ike_message.Authentication)
		case ike_message.TypeCP:
			configuration = ikePayload.(*ike_message.Configuration)
		default:
			ikeLog.Warnf(
				"Get IKE payload (type %d) in IKE_AUTH message, this payload will not be handled by IKE handler",
				ikePayload.Type())
		}
	}

	ikeSecurityAssociation.InitiatorMessageID = message.MessageID

	switch ikeSecurityAssociation.State {
	case PreSignalling:
		if initiatorID != nil {
			ikeLog.Info("Ecoding initiator for later IKE authentication")
			ikeSecurityAssociation.InitiatorID = initiatorID

			// Record maced identification for authentication
			idPayload := ike_message.IKEPayloadContainer{
				initiatorID,
			}
			idPayloadData, err := idPayload.Encode()
			if err != nil {
				ikeLog.Errorln(err)
				ikeLog.Error("Encoding ID payload message failed.")
				return
			}
			ikeSecurityAssociation.Prf_i.Reset()
			if _, err := ikeSecurityAssociation.Prf_i.Write(idPayloadData[4:]); err != nil {
				ikeLog.Errorf("Pseudorandom function write error: %v", err)
				return
			}
			ikeSecurityAssociation.InitiatorSignedOctets = append(
				ikeSecurityAssociation.InitiatorSignedOctets,
				ikeSecurityAssociation.Prf_i.Sum(nil)...)
		} else {
			ikeLog.Error("The initiator identification field is nil")
			// TODO: send error message to UE
			return
		}

		// Certificate request and prepare coresponding certificate
		// RFC 7296 section 3.7:
		// The Certificate Request payload is processed by inspecting the
		// Cert Encoding field to determine whether the processor has any
		// certificates of this type.  If so, the Certification Authority field
		// is inspected to determine if the processor has any certificates that
		// can be validated up to one of the specified certification
		// authorities.  This can be a chain of certificates.
		if certificateRequest != nil {
			ikeLog.Info("UE request N3IWF certificate")
			if ike_security.CompareRootCertificate(
				n3iwfCtx.CertificateAuthority,
				certificateRequest.CertificateEncoding,
				certificateRequest.CertificationAuthority) {
				// TODO: Complete N3IWF Certificate/Certificate Authority related procedure
				ikeLog.Info("Certificate Request sent from UE matches N3IWF CA")
			}
		}

		if certificate != nil {
			ikeLog.Info("UE send its certficate")
			ikeSecurityAssociation.InitiatorCertificate = certificate
		}

		if securityAssociation != nil {
			ikeLog.Info("Parsing security association")
			responseSecurityAssociation := new(ike_message.SecurityAssociation)

			for _, proposal := range securityAssociation.Proposals {
				var encryptionAlgorithmTransform *ike_message.Transform = nil
				var integrityAlgorithmTransform *ike_message.Transform = nil
				var diffieHellmanGroupTransform *ike_message.Transform = nil
				var extendedSequenceNumbersTransform *ike_message.Transform = nil

				if len(proposal.SPI) != 4 {
					continue // The SPI of ESP must be 32-bit
				}

				if len(proposal.EncryptionAlgorithm) > 0 {
					for _, transform := range proposal.EncryptionAlgorithm {
						if isTransformKernelSupported(ike_message.TypeEncryptionAlgorithm, transform.TransformID,
							transform.AttributePresent, transform.AttributeValue) {
							encryptionAlgorithmTransform = transform
							break
						}
					}
					if encryptionAlgorithmTransform == nil {
						continue
					}
				} else {
					continue // mandatory
				}
				if len(proposal.PseudorandomFunction) > 0 {
					continue // Pseudorandom function is not used by ESP
				}
				if len(proposal.IntegrityAlgorithm) > 0 {
					for _, transform := range proposal.IntegrityAlgorithm {
						if isTransformKernelSupported(ike_message.TypeIntegrityAlgorithm, transform.TransformID,
							transform.AttributePresent, transform.AttributeValue) {
							integrityAlgorithmTransform = transform
							break
						}
					}
					if integrityAlgorithmTransform == nil {
						continue
					}
				} // Optional
				if len(proposal.DiffieHellmanGroup) > 0 {
					for _, transform := range proposal.DiffieHellmanGroup {
						if isTransformKernelSupported(ike_message.TypeDiffieHellmanGroup, transform.TransformID,
							transform.AttributePresent, transform.AttributeValue) {
							diffieHellmanGroupTransform = transform
							break
						}
					}
					if diffieHellmanGroupTransform == nil {
						continue
					}
				} // Optional
				if len(proposal.ExtendedSequenceNumbers) > 0 {
					for _, transform := range proposal.ExtendedSequenceNumbers {
						if isTransformKernelSupported(ike_message.TypeExtendedSequenceNumbers, transform.TransformID,
							transform.AttributePresent, transform.AttributeValue) {
							extendedSequenceNumbersTransform = transform
							break
						}
					}
					if extendedSequenceNumbersTransform == nil {
						continue
					}
				} else {
					continue // Mandatory
				}

				chosenProposal := responseSecurityAssociation.Proposals.BuildProposal(
					proposal.ProposalNumber, proposal.ProtocolID, proposal.SPI)
				chosenProposal.EncryptionAlgorithm = append(chosenProposal.EncryptionAlgorithm, encryptionAlgorithmTransform)
				chosenProposal.ExtendedSequenceNumbers = append(
					chosenProposal.ExtendedSequenceNumbers, extendedSequenceNumbersTransform)
				if integrityAlgorithmTransform != nil {
					chosenProposal.IntegrityAlgorithm = append(chosenProposal.IntegrityAlgorithm, integrityAlgorithmTransform)
				}
				if diffieHellmanGroupTransform != nil {
					chosenProposal.DiffieHellmanGroup = append(chosenProposal.DiffieHellmanGroup, diffieHellmanGroupTransform)
				}
				break
			}

			if len(responseSecurityAssociation.Proposals) == 0 {
				ikeLog.Warn("No proposal chosen")
				// Respond NO_PROPOSAL_CHOSEN to UE
				// Notification
				responseIKEPayload.BuildNotification(
					ike_message.TypeNone, ike_message.NO_PROPOSAL_CHOSEN, nil, nil)

				// Build IKE message
				responseIKEMessage := ike_message.NewMessage(message.InitiatorSPI, message.ResponderSPI,
					ike_message.IKE_AUTH, true, false, message.MessageID, responseIKEPayload)

				// Send IKE message to UE
				err := SendIKEMessageToUE(udpConn, n3iwfAddr, ueAddr, responseIKEMessage,
					ikeSecurityAssociation.IKESAKey)
				if err != nil {
					ikeLog.Errorf("HandleIKEAUTH(): %v", err)
				}
				return
			}

			ikeSecurityAssociation.IKEAuthResponseSA = responseSecurityAssociation
		} else {
			ikeLog.Error("The security association field is nil")
			// TODO: send error message to UE
			return
		}

		if trafficSelectorInitiator == nil {
			ikeLog.Error("The initiator traffic selector field is nil")
			// TODO: send error message to UE
			return
		}
		ikeLog.Info("Received traffic selector initiator from UE")
		ikeSecurityAssociation.TrafficSelectorInitiator = trafficSelectorInitiator

		if trafficSelectorResponder == nil {
			ikeLog.Error("The initiator traffic selector field is nil")
			// TODO: send error message to UE
			return
		}
		ikeLog.Info("Received traffic selector initiator from UE")
		ikeSecurityAssociation.TrafficSelectorResponder = trafficSelectorResponder

		responseIKEPayload.Reset()
		// Identification
		responseIKEPayload.BuildIdentificationResponder(ike_message.ID_FQDN, []byte(cfg.GetFQDN()))

		// Certificate
		responseIKEPayload.BuildCertificate(
			ike_message.X509CertificateSignature,
			n3iwfCtx.N3IWFCertificate)

		// Authentication Data
		ikeLog.Tracef("Local authentication data:\n%s", hex.Dump(ikeSecurityAssociation.ResponderSignedOctets))
		sha1HashFunction := sha1.New() // #nosec G401
		_, err := sha1HashFunction.Write(ikeSecurityAssociation.ResponderSignedOctets)
		if err != nil {
			ikeLog.Errorf("Hash function write error: %v", err)
			return
		}

		var signedAuth []byte

		signedAuth, err = rsa.SignPKCS1v15(
			rand.Reader, n3iwfCtx.N3IWFPrivateKey,
			crypto.SHA1, sha1HashFunction.Sum(nil))
		if err != nil {
			ikeLog.Errorf("Sign authentication data failed: %v", err)
		}

		responseIKEPayload.BuildAuthentication(ike_message.RSADigitalSignature, signedAuth)

		// EAP expanded 5G-Start
		var identifier uint8
		for {
			identifier, err = ike_security.GenerateRandomUint8()
			if err != nil {
				ikeLog.Errorf("Random number failed: %v", err)
				return
			}
			if identifier != ikeSecurityAssociation.LastEAPIdentifier {
				ikeSecurityAssociation.LastEAPIdentifier = identifier
				break
			}
		}
		responseIKEPayload.BuildEAP5GStart(identifier)

		// Build IKE message
		responseIKEMessage := ike_message.NewMessage(message.InitiatorSPI, message.ResponderSPI,
			ike_message.IKE_AUTH, true, false, message.MessageID, responseIKEPayload)

		// Shift state
		ikeSecurityAssociation.State++

		// Send IKE message to UE
		err = SendIKEMessageToUE(udpConn, n3iwfAddr, ueAddr, responseIKEMessage,
			ikeSecurityAssociation.IKESAKey)
		if err != nil {
			ikeLog.Errorf("HandleIKEAUTH(): %v", err)
			return
		}

	case EAPSignalling:
		// If success, N3IWF will send an UPLinkNASTransport to AMF
		if eap != nil {
			if eap.Code != ike_message.EAPCodeResponse {
				ikeLog.Error("[EAP] Received an EAP payload with code other than response. Drop the payload.")
				return
			}
			if eap.Identifier != ikeSecurityAssociation.LastEAPIdentifier {
				ikeLog.Error("[EAP] Received an EAP payload with unmatched identifier. Drop the payload.")
				return
			}

			eapTypeData := eap.EAPTypeData[0]
			var eapExpanded *ike_message.EAPExpanded

			switch eapTypeData.Type() {
			// TODO: handle
			// case ike_message.EAPTypeIdentity:
			// case ike_message.EAPTypeNotification:
			// case ike_message.EAPTypeNak:
			case ike_message.EAPTypeExpanded:
				eapExpanded = eapTypeData.(*ike_message.EAPExpanded)
			default:
				ikeLog.Errorf("[EAP] Received EAP packet with type other than EAP expanded type: %d", eapTypeData.Type())
				return
			}

			if eapExpanded.VendorID != ike_message.VendorID3GPP {
				ikeLog.Error("The peer sent EAP expended packet with wrong vendor ID. Drop the packet.")
				return
			}
			if eapExpanded.VendorType != ike_message.VendorTypeEAP5G {
				ikeLog.Error("The peer sent EAP expanded packet with wrong vendor type. Drop the packet.")
				return
			}

			eap5GMessageID := eapExpanded.VendorData[0]
			ikeLog.Infof("EAP5G MessageID : %+v", eap5GMessageID)

			if eap5GMessageID == ike_message.EAP5GType5GStop {
				// Send EAP failure
				responseIKEPayload.Reset()

				// EAP
				identifier, err := ike_security.GenerateRandomUint8()
				if err != nil {
					ikeLog.Errorf("Generate random uint8 failed: %v", err)
					return
				}
				responseIKEPayload.BuildEAPfailure(identifier)

				// Build IKE message
				responseIKEMessage := ike_message.NewMessage(message.InitiatorSPI, message.ResponderSPI,
					ike_message.IKE_AUTH, true, false, message.MessageID, responseIKEPayload)

				// Send IKE message to UE
				err = SendIKEMessageToUE(udpConn, n3iwfAddr, ueAddr, responseIKEMessage,
					ikeSecurityAssociation.IKESAKey)
				if err != nil {
					ikeLog.Errorf("HandleIKEAUTH(): %v", err)
				}
				return
			}

			var ranNgapId int64
			ranNgapId, ok = n3iwfCtx.NgapIdLoad(ikeSecurityAssociation.LocalSPI)
			if !ok {
				ranNgapId = 0
			}

			s.SendNgapEvt(n3iwf_context.NewUnmarshalEAP5GDataEvt(
				ikeSecurityAssociation.LocalSPI,
				eapExpanded.VendorData,
				ikeSecurityAssociation.IkeUE != nil,
				ranNgapId,
			))

			ikeSecurityAssociation.IKEConnection = &n3iwf_context.UDPSocketInfo{
				Conn:      udpConn,
				N3IWFAddr: n3iwfAddr,
				UEAddr:    ueAddr,
			}

			ikeSecurityAssociation.InitiatorMessageID = message.MessageID
		} else {
			ikeLog.Error("EAP is nil")
		}

	case PostSignalling:
		// Load needed information
		ikeUE := ikeSecurityAssociation.IkeUE

		// Prepare pseudorandom function for calculating/verifying authentication data
		pseudorandomFunction := ikeSecurityAssociation.PrfInfo.Init(ikeUE.Kn3iwf)
		_, err := pseudorandomFunction.Write([]byte("Key Pad for IKEv2"))
		if err != nil {
			ikeLog.Errorf("Pseudorandom function write error: %v", err)
			return
		}
		secret := pseudorandomFunction.Sum(nil)
		pseudorandomFunction = ikeSecurityAssociation.PrfInfo.Init(secret)

		if authentication != nil {
			// Verifying remote AUTH
			pseudorandomFunction.Reset()
			_, err = pseudorandomFunction.Write(ikeSecurityAssociation.InitiatorSignedOctets)
			if err != nil {
				ikeLog.Errorf("Pseudorandom function write error: %v", err)
				return
			}
			expectedAuthenticationData := pseudorandomFunction.Sum(nil)

			ikeLog.Tracef("Kn3iwf:\n%s", hex.Dump(ikeUE.Kn3iwf))
			ikeLog.Tracef("secret:\n%s", hex.Dump(secret))
			ikeLog.Tracef("InitiatorSignedOctets:\n%s", hex.Dump(ikeSecurityAssociation.InitiatorSignedOctets))
			ikeLog.Tracef("Expected Authentication Data:\n%s", hex.Dump(expectedAuthenticationData))
			if !bytes.Equal(authentication.AuthenticationData, expectedAuthenticationData) {
				ikeLog.Warn("Peer authentication failed.")
				// Inform UE the authentication has failed
				responseIKEPayload.Reset()

				// Notification
				responseIKEPayload.BuildNotification(
					ike_message.TypeNone, ike_message.AUTHENTICATION_FAILED, nil, nil)

				// Build IKE message
				responseIKEMessage := ike_message.NewMessage(message.InitiatorSPI, message.ResponderSPI,
					ike_message.IKE_AUTH, true, false, message.MessageID, responseIKEPayload)

				// Send IKE message to UE
				err = SendIKEMessageToUE(udpConn, n3iwfAddr, ueAddr, responseIKEMessage,
					ikeSecurityAssociation.IKESAKey)
				if err != nil {
					ikeLog.Errorf("HandleIKEAUTH(): %v", err)
				}
				return
			} else {
				ikeLog.Tracef("Peer authentication success")
			}
		} else {
			ikeLog.Warn("Peer authentication failed.")
			// Inform UE the authentication has failed
			responseIKEPayload.Reset()

			// Notification
			responseIKEPayload.BuildNotification(ike_message.TypeNone, ike_message.AUTHENTICATION_FAILED, nil, nil)

			// Build IKE message
			responseIKEMessage := ike_message.NewMessage(message.InitiatorSPI, message.ResponderSPI,
				ike_message.IKE_AUTH, true, false, message.MessageID, responseIKEPayload)

			// Send IKE message to UE
			err = SendIKEMessageToUE(udpConn, n3iwfAddr, ueAddr, responseIKEMessage,
				ikeSecurityAssociation.IKESAKey)
			if err != nil {
				ikeLog.Errorf("HandleIKEAUTH(): %v", err)
			}
			return
		}

		// Parse configuration request to get if the UE has requested internal address,
		// and prepare configuration payload to UE
		addrRequest := false

		if configuration != nil {
			ikeLog.Tracef("Received configuration payload with type: %d", configuration.ConfigurationType)

			var attribute *ike_message.IndividualConfigurationAttribute
			for _, attribute = range configuration.ConfigurationAttribute {
				switch attribute.Type {
				case ike_message.INTERNAL_IP4_ADDRESS:
					addrRequest = true
					if len(attribute.Value) != 0 {
						ikeLog.Tracef("Got client requested address: %d.%d.%d.%d",
							attribute.Value[0], attribute.Value[1], attribute.Value[2], attribute.Value[3])
					}
				default:
					ikeLog.Warnf("Receive other type of configuration request: %d", attribute.Type)
				}
			}
		} else {
			ikeLog.Warn("Configuration is nil. UE did not sent any configuration request.")
		}

		responseIKEPayload.Reset()

		// Calculate local AUTH
		pseudorandomFunction.Reset()
		_, err = pseudorandomFunction.Write(ikeSecurityAssociation.ResponderSignedOctets)
		if err != nil {
			ikeLog.Errorf("Pseudorandom function write error: %v", err)
			return
		}

		// Authentication
		responseIKEPayload.BuildAuthentication(
			ike_message.SharedKeyMesageIntegrityCode, pseudorandomFunction.Sum(nil))

		// Prepare configuration payload and traffic selector payload for initiator and responder
		var ueIPAddr, n3iwfIPAddr net.IP
		if addrRequest {
			// IP addresses (IPSec)
			var ueIp net.IP
			ueIp, err = n3iwfCtx.NewIPsecInnerUEIP(ikeUE)
			if err != nil {
				ikeLog.Errorf("HandleIKEAUTH(): %v", err)
				return
			}
			ueIPAddr = ueIp.To4()
			n3iwfIPAddr = net.ParseIP(ipsecGwAddr).To4()

			responseConfiguration := responseIKEPayload.BuildConfiguration(
				ike_message.CFG_REPLY)
			responseConfiguration.ConfigurationAttribute.BuildConfigurationAttribute(
				ike_message.INTERNAL_IP4_ADDRESS, ueIPAddr)
			responseConfiguration.ConfigurationAttribute.BuildConfigurationAttribute(
				ike_message.INTERNAL_IP4_NETMASK, n3iwfCtx.IPSecInnerIPPool.IPSubnet.Mask)

			var ipsecInnerIPAddr *net.IPAddr
			ikeUE.IPSecInnerIP = ueIPAddr
			ipsecInnerIPAddr, err = net.ResolveIPAddr("ip", ueIPAddr.String())
			if err != nil {
				ikeLog.Errorf("Resolve UE inner IP address failed: %v", err)
				return
			}
			ikeUE.IPSecInnerIPAddr = ipsecInnerIPAddr
			ikeLog.Tracef("ueIPAddr: %+v", ueIPAddr)
		} else {
			ikeLog.Error("UE did not send any configuration request for its IP address.")
			return
		}

		// Security Association
		responseIKEPayload = append(responseIKEPayload, ikeSecurityAssociation.IKEAuthResponseSA)

		// Traffic Selectors initiator/responder
		responseTrafficSelectorInitiator := responseIKEPayload.BuildTrafficSelectorInitiator()
		responseTrafficSelectorInitiator.TrafficSelectors.BuildIndividualTrafficSelector(
			ike_message.TS_IPV4_ADDR_RANGE, ike_message.IPProtocolAll, 0, 65535, ueIPAddr.To4(), ueIPAddr.To4())
		responseTrafficSelectorResponder := responseIKEPayload.BuildTrafficSelectorResponder()
		responseTrafficSelectorResponder.TrafficSelectors.BuildIndividualTrafficSelector(
			ike_message.TS_IPV4_ADDR_RANGE, ike_message.IPProtocolAll, 0, 65535, n3iwfIPAddr.To4(), n3iwfIPAddr.To4())

		// Record traffic selector to IKE security association
		ikeSecurityAssociation.TrafficSelectorInitiator = responseTrafficSelectorInitiator
		ikeSecurityAssociation.TrafficSelectorResponder = responseTrafficSelectorResponder

		// Get data needed by xfrm

		// Allocate N3IWF inbound SPI
		var inboundSPI uint32
		inboundSPIByte := make([]byte, 4)
		for {
			buf := make([]byte, 4)
			_, err = rand.Read(buf)
			if err != nil {
				ikeLog.Errorf("Handle IKE_AUTH Generate ChildSA inboundSPI: %v", err)
				return
			}
			randomUint32 := binary.BigEndian.Uint32(buf)
			// check if the inbound SPI havn't been allocated by N3IWF
			if _, ok1 := n3iwfCtx.ChildSA.Load(randomUint32); !ok1 {
				inboundSPI = randomUint32
				break
			}
		}
		binary.BigEndian.PutUint32(inboundSPIByte, inboundSPI)

		outboundSPI := binary.BigEndian.Uint32(ikeSecurityAssociation.IKEAuthResponseSA.Proposals[0].SPI)
		ikeLog.Infof("Inbound SPI: 0x%08x, Outbound SPI: 0x%08x", inboundSPI, outboundSPI)

		// SPI field of IKEAuthResponseSA is used to save outbound SPI temporarily.
		// After N3IWF produced its inbound SPI, the field will be overwritten with the SPI.
		ikeSecurityAssociation.IKEAuthResponseSA.Proposals[0].SPI = inboundSPIByte

		// Consider 0x01 as the speicified index for IKE_AUTH exchange
		ikeUE.CreateHalfChildSA(0x01, inboundSPI, -1)
		childSecurityAssociationContext, err := ikeUE.CompleteChildSA(
			0x01, outboundSPI, ikeSecurityAssociation.IKEAuthResponseSA)
		if err != nil {
			ikeLog.Errorf("HandleIKEAUTH(): Create child security association context failed: %v", err)
			return
		}
		err = s.parseIPAddressInformationToChildSecurityAssociation(
			childSecurityAssociationContext, ueAddr.IP,
			ikeSecurityAssociation.TrafficSelectorResponder.TrafficSelectors[0],
			ikeSecurityAssociation.TrafficSelectorInitiator.TrafficSelectors[0])
		if err != nil {
			ikeLog.Errorf("Parse IP address to child security association failed: %v", err)
			return
		}
		// Select TCP traffic
		childSecurityAssociationContext.SelectedIPProtocol = unix.IPPROTO_TCP

		errGen := childSecurityAssociationContext.ChildSAKey.GenerateKeyForChildSA(ikeSecurityAssociation.IKESAKey,
			ikeSecurityAssociation.ConcatenatedNonce)
		if errGen != nil {
			ikeLog.Errorf("Generate key for child SA failed: %v", errGen)
			return
		}
		// NAT-T concern
		if ikeSecurityAssociation.UeBehindNAT || ikeSecurityAssociation.N3iwfBehindNAT {
			childSecurityAssociationContext.EnableEncapsulate = true
			childSecurityAssociationContext.N3IWFPort = n3iwfAddr.Port
			childSecurityAssociationContext.NATPort = ueAddr.Port
		}

		// Notification(NAS_IP_ADDRESS)
		responseIKEPayload.BuildNotifyNAS_IP4_ADDRESS(ipsecGwAddr)

		// Notification(NSA_TCP_PORT)
		responseIKEPayload.BuildNotifyNAS_TCP_PORT(cfg.GetNasTcpPort())

		// Build IKE message
		responseIKEMessage := ike_message.NewMessage(message.InitiatorSPI, message.ResponderSPI,
			ike_message.IKE_AUTH, true, false, message.MessageID, responseIKEPayload)

		childSecurityAssociationContext.LocalIsInitiator = false
		// Aplly XFRM rules
		// IPsec for CP always use default XFRM interface
		err = xfrm.ApplyXFRMRule(false, cfg.GetXfrmIfaceId(), childSecurityAssociationContext)
		if err != nil {
			ikeLog.Errorf("Applying XFRM rules failed: %v", err)
			return
		}
		ikeLog.Debugln(childSecurityAssociationContext.String(cfg.GetXfrmIfaceId()))

		// Send IKE message to UE
		err = SendIKEMessageToUE(udpConn, n3iwfAddr, ueAddr, responseIKEMessage,
			ikeSecurityAssociation.IKESAKey)
		if err != nil {
			ikeLog.Errorf("HandleIKEAUTH(): %v", err)
			return
		}

		ranNgapId, ok := n3iwfCtx.NgapIdLoad(ikeUE.N3IWFIKESecurityAssociation.LocalSPI)
		if !ok {
			ikeLog.Errorf("Cannot get RanNgapId from SPI : %+v",
				ikeUE.N3IWFIKESecurityAssociation.LocalSPI)
			return
		}

		ikeSecurityAssociation.State++

		// After this, N3IWF will forward NAS with Child SA (IPSec SA)
		s.SendNgapEvt(n3iwf_context.NewStartTCPSignalNASMsgEvt(ranNgapId))

		// Get TempPDUSessionSetupData from NGAP to setup PDU session if needed
		s.SendNgapEvt(n3iwf_context.NewGetNGAPContextEvt(
			ranNgapId, []int64{n3iwf_context.CxtTempPDUSessionSetupData},
		))
	}
}

func (s *Server) HandleCREATECHILDSA(
	udpConn *net.UDPConn,
	n3iwfAddr, ueAddr *net.UDPAddr,
	message *ike_message.IKEMessage,
	ikeSecurityAssociation *n3iwf_context.IKESecurityAssociation,
) {
	ikeLog := logger.IKELog
	ikeLog.Infoln("Handle CREATE_CHILD_SA")

	n3iwfCtx := s.Context()

	if !ikeSecurityAssociation.IKEConnection.UEAddr.IP.Equal(ueAddr.IP) ||
		!ikeSecurityAssociation.IKEConnection.N3IWFAddr.IP.Equal(n3iwfAddr.IP) {
		ikeLog.Warnf("Get unexpteced IP in SPI: %016x", ikeSecurityAssociation.LocalSPI)
		return
	}

	// Parse payloads
	var securityAssociation *ike_message.SecurityAssociation
	var nonce *ike_message.Nonce
	var trafficSelectorInitiator *ike_message.TrafficSelectorInitiator
	var trafficSelectorResponder *ike_message.TrafficSelectorResponder

	for _, ikePayload := range message.Payloads {
		switch ikePayload.Type() {
		case ike_message.TypeSA:
			securityAssociation = ikePayload.(*ike_message.SecurityAssociation)
		case ike_message.TypeNiNr:
			nonce = ikePayload.(*ike_message.Nonce)
		case ike_message.TypeTSi:
			trafficSelectorInitiator = ikePayload.(*ike_message.TrafficSelectorInitiator)
		case ike_message.TypeTSr:
			trafficSelectorResponder = ikePayload.(*ike_message.TrafficSelectorResponder)
		default:
			ikeLog.Warnf(
				"Get IKE payload (type %d) in CREATE_CHILD_SA message, this payload will not be handled by IKE handler",
				ikePayload.Type())
		}
	}

	// Check received message
	if securityAssociation == nil {
		ikeLog.Error("The security association field is nil")
		return
	}

	if trafficSelectorInitiator == nil {
		ikeLog.Error("The traffic selector initiator field is nil")
		return
	}

	if trafficSelectorResponder == nil {
		ikeLog.Error("The traffic selector responder field is nil")
		return
	}

	// Nonce
	if nonce == nil {
		ikeLog.Error("The nonce field is nil")
		// TODO: send error message to UE
		return
	}
	ikeSecurityAssociation.ConcatenatedNonce = append(
		ikeSecurityAssociation.ConcatenatedNonce, nonce.NonceData...)

	ikeSecurityAssociation.TemporaryIkeMsg = &n3iwf_context.IkeMsgTemporaryData{
		SecurityAssociation:      securityAssociation,
		TrafficSelectorInitiator: trafficSelectorInitiator,
		TrafficSelectorResponder: trafficSelectorResponder,
	}

	ranNgapId, ok := n3iwfCtx.NgapIdLoad(ikeSecurityAssociation.LocalSPI)
	if !ok {
		ikeLog.Errorf("Cannot get RanNgapID from SPI : %+v",
			ikeSecurityAssociation.LocalSPI)
		return
	}

	ngapCxtReqNumlist := []int64{n3iwf_context.CxtTempPDUSessionSetupData}

	s.SendNgapEvt(n3iwf_context.NewGetNGAPContextEvt(ranNgapId, ngapCxtReqNumlist))
}

func (s *Server) continueCreateChildSA(
	ikeSecurityAssociation *n3iwf_context.IKESecurityAssociation,
	temporaryPDUSessionSetupData *n3iwf_context.PDUSessionSetupTemporaryData,
) {
	ikeLog := logger.IKELog
	n3iwfCtx := s.Context()
	cfg := s.Config()
	ipsecGwAddr := cfg.GetIPSecGatewayAddr()

	// UE context
	ikeUe := ikeSecurityAssociation.IkeUE
	if ikeUe == nil {
		ikeLog.Error("UE context is nil")
		return
	}

	// PDU session information
	if temporaryPDUSessionSetupData == nil {
		ikeLog.Error("No PDU session information")
		return
	}

	if len(temporaryPDUSessionSetupData.UnactivatedPDUSession) == 0 {
		ikeLog.Error("No unactivated PDU session information")
		return
	}

	temporaryIkeMsg := ikeSecurityAssociation.TemporaryIkeMsg
	ikeConnection := ikeSecurityAssociation.IKEConnection

	// Get xfrm needed data
	// As specified in RFC 7296, ESP negotiate two child security association (pair) in one exchange
	// Message ID is used to be a index to pair two SPI in serveral IKE messages.
	outboundSPI := binary.BigEndian.Uint32(temporaryIkeMsg.SecurityAssociation.Proposals[0].SPI)
	childSecurityAssociationContext, err := ikeUe.CompleteChildSA(
		ikeSecurityAssociation.ResponderMessageID, outboundSPI, temporaryIkeMsg.SecurityAssociation)
	if err != nil {
		ikeLog.Errorf("Create child security association context failed: %v", err)
		return
	}

	// Build TSi if there is no one in the response
	if len(temporaryIkeMsg.TrafficSelectorInitiator.TrafficSelectors) == 0 {
		ikeLog.Warnf("There is no TSi in CREATE_CHILD_SA response.")
		n3iwfIPAddr := net.ParseIP(ipsecGwAddr)
		temporaryIkeMsg.TrafficSelectorInitiator.TrafficSelectors.BuildIndividualTrafficSelector(
			ike_message.TS_IPV4_ADDR_RANGE, ike_message.IPProtocolAll,
			0, 65535, n3iwfIPAddr, n3iwfIPAddr)
	}

	// Build TSr if there is no one in the response
	if len(temporaryIkeMsg.TrafficSelectorResponder.TrafficSelectors) == 0 {
		ikeLog.Warnf("There is no TSr in CREATE_CHILD_SA response.")
		ueIPAddr := ikeUe.IPSecInnerIP
		temporaryIkeMsg.TrafficSelectorResponder.TrafficSelectors.BuildIndividualTrafficSelector(
			ike_message.TS_IPV4_ADDR_RANGE, ike_message.IPProtocolAll,
			0, 65535, ueIPAddr, ueIPAddr)
	}

	err = s.parseIPAddressInformationToChildSecurityAssociation(
		childSecurityAssociationContext,
		ikeConnection.UEAddr.IP,
		temporaryIkeMsg.TrafficSelectorInitiator.TrafficSelectors[0],
		temporaryIkeMsg.TrafficSelectorResponder.TrafficSelectors[0])
	if err != nil {
		ikeLog.Errorf("Parse IP address to child security association failed: %v", err)
		return
	}
	// Select GRE traffic
	childSecurityAssociationContext.SelectedIPProtocol = unix.IPPROTO_GRE

	err = childSecurityAssociationContext.ChildSAKey.GenerateKeyForChildSA(ikeSecurityAssociation.IKESAKey,
		ikeSecurityAssociation.ConcatenatedNonce)
	if err != nil {
		ikeLog.Errorf("Generate key for child SA failed: %v", err)
		return
	}
	// NAT-T concern
	if ikeSecurityAssociation.UeBehindNAT || ikeSecurityAssociation.N3iwfBehindNAT {
		childSecurityAssociationContext.EnableEncapsulate = true
		childSecurityAssociationContext.N3IWFPort = ikeConnection.N3IWFAddr.Port
		childSecurityAssociationContext.NATPort = ikeConnection.UEAddr.Port
	}

	newXfrmiId := cfg.GetXfrmIfaceId()

	pduSessionListLen := ikeUe.PduSessionListLen

	// The additional PDU session will be separated from default xfrm interface
	// to avoid SPD entry collision
	if pduSessionListLen > 1 {
		// Setup XFRM interface for ipsec
		var linkIPSec netlink.Link
		n3iwfIPAddr := net.ParseIP(ipsecGwAddr).To4()
		n3iwfIPAddrAndSubnet := net.IPNet{IP: n3iwfIPAddr, Mask: n3iwfCtx.IPSecInnerIPPool.IPSubnet.Mask}
		newXfrmiId += cfg.GetXfrmIfaceId() + n3iwfCtx.XfrmIfaceIdOffsetForUP
		newXfrmiName := fmt.Sprintf("%s-%d", cfg.GetXfrmIfaceName(), newXfrmiId)

		linkIPSec, err = xfrm.SetupIPsecXfrmi(
			newXfrmiName, n3iwfCtx.XfrmParentIfaceName,
			newXfrmiId, n3iwfIPAddrAndSubnet)
		if err != nil {
			ikeLog.Errorf("Setup XFRM interface %s fail: %v", newXfrmiName, err)
			return
		}

		ikeLog.Infof("Setup XFRM interface: %s", newXfrmiName)
		n3iwfCtx.XfrmIfaces.LoadOrStore(newXfrmiId, linkIPSec)
		childSecurityAssociationContext.XfrmIface = linkIPSec
		n3iwfCtx.XfrmIfaceIdOffsetForUP++
	} else {
		linkIPSec, ok := n3iwfCtx.XfrmIfaces.Load(newXfrmiId)
		if !ok {
			ikeLog.Warnf("Cannot find the XFRM interface with if_id: %d", newXfrmiId)
			return
		}
		childSecurityAssociationContext.XfrmIface = linkIPSec.(netlink.Link)
	}

	// Aplly XFRM rules
	childSecurityAssociationContext.LocalIsInitiator = true
	err = xfrm.ApplyXFRMRule(true, newXfrmiId, childSecurityAssociationContext)
	if err != nil {
		ikeLog.Errorf("Applying XFRM rules failed: %v", err)
		return
	}
	ikeLog.Debugln(childSecurityAssociationContext.String(newXfrmiId))

	ranNgapId, ok := n3iwfCtx.NgapIdLoad(ikeSecurityAssociation.LocalSPI)
	if !ok {
		ikeLog.Errorf("Cannot get RanNgapId from SPI : %+v",
			ikeSecurityAssociation.LocalSPI)
		return
	}
	// Forward PDU Seesion Establishment Accept to UE
	s.SendNgapEvt(n3iwf_context.NewSendNASMsgEvt(ranNgapId))

	temporaryPDUSessionSetupData.FailedErrStr = append(temporaryPDUSessionSetupData.FailedErrStr, n3iwf_context.ErrNil)

	ikeSecurityAssociation.ResponderMessageID++

	// If needed, setup another PDU session
	s.CreatePDUSessionChildSA(ikeUe, temporaryPDUSessionSetupData)
}

func (s *Server) HandleInformational(
	udpConn *net.UDPConn,
	n3iwfAddr, ueAddr *net.UDPAddr,
	message *ike_message.IKEMessage,
	ikeSecurityAssociation *n3iwf_context.IKESecurityAssociation,
) {
	ikeLog := logger.IKELog
	ikeLog.Infoln("Handle Informational")

	var deletePayload *ike_message.Delete
	var err error
	responseIKEPayload := new(ike_message.IKEPayloadContainer)

	n3iwfIke := ikeSecurityAssociation.IkeUE

	if n3iwfIke.N3IWFIKESecurityAssociation.DPDReqRetransTimer != nil {
		n3iwfIke.N3IWFIKESecurityAssociation.DPDReqRetransTimer.Stop()
		n3iwfIke.N3IWFIKESecurityAssociation.DPDReqRetransTimer = nil
		atomic.StoreInt32(&n3iwfIke.N3IWFIKESecurityAssociation.CurrentRetryTimes, 0)
	}

	for _, ikePayload := range message.Payloads {
		switch ikePayload.Type() {
		case ike_message.TypeD:
			deletePayload = ikePayload.(*ike_message.Delete)
		default:
			ikeLog.Warnf(
				"Get IKE payload (type %d) in Inoformational message, this payload will not be handled by IKE handler",
				ikePayload.Type())
		}
	}

	if deletePayload != nil {
		responseIKEPayload, err = s.handleDeletePayload(deletePayload, message.IsResponse(), ikeSecurityAssociation)
		if err != nil {
			ikeLog.Errorf("HandleInformational(): %v", err)
			return
		}
	}

	if message.IsResponse() {
		ikeSecurityAssociation.ResponderMessageID++
	} else { // Get Request message
		SendUEInformationExchange(ikeSecurityAssociation, ikeSecurityAssociation.IKESAKey,
			responseIKEPayload, false, true, message.MessageID,
			udpConn, ueAddr, n3iwfAddr)
	}
}

func (s *Server) HandleEvent(ikeEvt n3iwf_context.IkeEvt) {
	ikeLog := logger.IKELog
	ikeLog.Infof("Handle IKE event")

	switch ikeEvt.Type() {
	case n3iwf_context.UnmarshalEAP5GDataResponse:
		s.HandleUnmarshalEAP5GDataResponse(ikeEvt)
	case n3iwf_context.SendEAP5GFailureMsg:
		s.HandleSendEAP5GFailureMsg(ikeEvt)
	case n3iwf_context.SendEAPSuccessMsg:
		s.HandleSendEAPSuccessMsg(ikeEvt)
	case n3iwf_context.SendEAPNASMsg:
		s.HandleSendEAPNASMsg(ikeEvt)
	case n3iwf_context.CreatePDUSession:
		s.HandleCreatePDUSession(ikeEvt)
	case n3iwf_context.IKEDeleteRequest:
		s.HandleIKEDeleteEvt(ikeEvt)
	case n3iwf_context.SendChildSADeleteRequest:
		s.HandleSendChildSADeleteRequest(ikeEvt)
	case n3iwf_context.IKEContextUpdate:
		s.HandleIKEContextUpdate(ikeEvt)
	case n3iwf_context.GetNGAPContextResponse:
		s.HandleGetNGAPContextResponse(ikeEvt)
	default:
		ikeLog.Errorf("Undefine IKE event type : %d", ikeEvt.Type())
		return
	}
}

func (s *Server) HandleUnmarshalEAP5GDataResponse(ikeEvt n3iwf_context.IkeEvt) {
	ikeLog := logger.IKELog
	ikeLog.Infof("Handle UnmarshalEAP5GDataResponse event")

	unmarshalEAP5GDataResponseEvt := ikeEvt.(*n3iwf_context.UnmarshalEAP5GDataResponseEvt)
	localSPI := unmarshalEAP5GDataResponseEvt.LocalSPI
	ranUeNgapId := unmarshalEAP5GDataResponseEvt.RanUeNgapId
	nasPDU := unmarshalEAP5GDataResponseEvt.NasPDU

	n3iwfCtx := s.Context()
	ikeSecurityAssociation, _ := n3iwfCtx.IKESALoad(localSPI)

	// Create UE context
	ikeUe := n3iwfCtx.NewN3iwfIkeUe(localSPI)

	// Relative context
	ikeSecurityAssociation.IkeUE = ikeUe
	ikeUe.N3IWFIKESecurityAssociation = ikeSecurityAssociation
	ikeUe.IKEConnection = ikeSecurityAssociation.IKEConnection

	n3iwfCtx.IkeSpiNgapIdMapping(ikeUe.N3IWFIKESecurityAssociation.LocalSPI, ranUeNgapId)

	s.SendNgapEvt(n3iwf_context.NewSendInitialUEMessageEvt(
		ranUeNgapId,
		ikeSecurityAssociation.IKEConnection.UEAddr.IP.To4().String(),
		ikeSecurityAssociation.IKEConnection.UEAddr.Port,
		nasPDU,
	))
}

func (s *Server) HandleSendEAP5GFailureMsg(ikeEvt n3iwf_context.IkeEvt) {
	ikeLog := logger.IKELog
	ikeLog.Infof("Handle SendEAP5GFailureMsg event")

	sendEAP5GFailureMsgEvt := ikeEvt.(*n3iwf_context.SendEAP5GFailureMsgEvt)
	errMsg := sendEAP5GFailureMsgEvt.ErrMsg
	localSPI := sendEAP5GFailureMsgEvt.LocalSPI

	n3iwfCtx := s.Context()
	ikeSecurityAssociation, _ := n3iwfCtx.IKESALoad(localSPI)
	ikeLog.Warnf("EAP Failure : %s", errMsg.Error())

	var responseIKEPayload ike_message.IKEPayloadContainer
	// Send EAP failure

	// EAP
	identifier, err := ike_security.GenerateRandomUint8()
	if err != nil {
		ikeLog.Errorf("Generate random uint8 failed: %v", err)
		return
	}
	responseIKEPayload.BuildEAPfailure(identifier)

	// Build IKE message
	responseIKEMessage := ike_message.NewMessage(ikeSecurityAssociation.RemoteSPI, ikeSecurityAssociation.LocalSPI,
		ike_message.IKE_AUTH, true, false, ikeSecurityAssociation.InitiatorMessageID, responseIKEPayload)

	// Send IKE message to UE
	err = SendIKEMessageToUE(ikeSecurityAssociation.IKEConnection.Conn,
		ikeSecurityAssociation.IKEConnection.N3IWFAddr, ikeSecurityAssociation.IKEConnection.UEAddr,
		responseIKEMessage, ikeSecurityAssociation.IKESAKey)
	if err != nil {
		ikeLog.Errorf("HandleSendEAP5GFailureMsg(): %v", err)
	}
}

func (s *Server) HandleSendEAPSuccessMsg(ikeEvt n3iwf_context.IkeEvt) {
	ikeLog := logger.IKELog
	ikeLog.Infof("Handle SendEAPSuccessMsg event")

	sendEAPSuccessMsgEvt := ikeEvt.(*n3iwf_context.SendEAPSuccessMsgEvt)
	localSPI := sendEAPSuccessMsgEvt.LocalSPI
	kn3iwf := sendEAPSuccessMsgEvt.Kn3iwf
	pduSessionListLen := sendEAPSuccessMsgEvt.PduSessionListLen

	n3iwfCtx := s.Context()
	ikeSecurityAssociation, _ := n3iwfCtx.IKESALoad(localSPI)

	if kn3iwf != nil {
		ikeSecurityAssociation.IkeUE.Kn3iwf = kn3iwf
	}

	ikeSecurityAssociation.IkeUE.PduSessionListLen = pduSessionListLen

	var responseIKEPayload ike_message.IKEPayloadContainer

	responseIKEPayload.Reset()

	var identifier uint8
	var err error
	for {
		identifier, err = ike_security.GenerateRandomUint8()
		if err != nil {
			ikeLog.Errorf("HandleSendEAPSuccessMsg() rand : %v", err)
			return
		}
		if identifier != ikeSecurityAssociation.LastEAPIdentifier {
			ikeSecurityAssociation.LastEAPIdentifier = identifier
			break
		}
	}

	responseIKEPayload.BuildEAPSuccess(identifier)

	// Build IKE message
	responseIKEMessage := ike_message.NewMessage(ikeSecurityAssociation.RemoteSPI,
		ikeSecurityAssociation.LocalSPI, ike_message.IKE_AUTH, true, false,
		ikeSecurityAssociation.InitiatorMessageID, responseIKEPayload)

	// Send IKE message to UE
	err = SendIKEMessageToUE(ikeSecurityAssociation.IKEConnection.Conn,
		ikeSecurityAssociation.IKEConnection.N3IWFAddr,
		ikeSecurityAssociation.IKEConnection.UEAddr, responseIKEMessage,
		ikeSecurityAssociation.IKESAKey)
	if err != nil {
		ikeLog.Errorf("HandleSendEAPSuccessMsg(): %v", err)
		return
	}

	ikeSecurityAssociation.State++
}

func (s *Server) HandleSendEAPNASMsg(ikeEvt n3iwf_context.IkeEvt) {
	ikeLog := logger.IKELog
	ikeLog.Infof("Handle SendEAPNASMsg event")

	sendEAPNASMsgEvt := ikeEvt.(*n3iwf_context.SendEAPNASMsgEvt)
	localSPI := sendEAPNASMsgEvt.LocalSPI
	nasPDU := sendEAPNASMsgEvt.NasPDU

	n3iwfCtx := s.Context()
	ikeSecurityAssociation, _ := n3iwfCtx.IKESALoad(localSPI)

	var responseIKEPayload ike_message.IKEPayloadContainer
	responseIKEPayload.Reset()

	var identifier uint8
	var err error
	for {
		identifier, err = ike_security.GenerateRandomUint8()
		if err != nil {
			ikeLog.Errorf("HandleSendEAPNASMsg() rand : %v", err)
			return
		}
		if identifier != ikeSecurityAssociation.LastEAPIdentifier {
			ikeSecurityAssociation.LastEAPIdentifier = identifier
			break
		}
	}

	err = responseIKEPayload.BuildEAP5GNAS(identifier, nasPDU)
	if err != nil {
		ikeLog.Errorf("HandleSendEAPNASMsg() BuildEAP5GNAS: %v", err)
		return
	}

	// Build IKE message
	responseIKEMessage := ike_message.NewMessage(ikeSecurityAssociation.RemoteSPI,
		ikeSecurityAssociation.LocalSPI, ike_message.IKE_AUTH, true, false,
		ikeSecurityAssociation.InitiatorMessageID, responseIKEPayload)

	// Send IKE message to UE
	err = SendIKEMessageToUE(ikeSecurityAssociation.IKEConnection.Conn,
		ikeSecurityAssociation.IKEConnection.N3IWFAddr,
		ikeSecurityAssociation.IKEConnection.UEAddr, responseIKEMessage,
		ikeSecurityAssociation.IKESAKey)
	if err != nil {
		ikeLog.Errorf("HandleSendEAPNASMsg(): %v", err)
	}
}

func (s *Server) HandleCreatePDUSession(ikeEvt n3iwf_context.IkeEvt) {
	ikeLog := logger.IKELog
	ikeLog.Infof("Handle CreatePDUSession event")

	createPDUSessionEvt := ikeEvt.(*n3iwf_context.CreatePDUSessionEvt)
	localSPI := createPDUSessionEvt.LocalSPI
	pduSessionListLen := createPDUSessionEvt.PduSessionListLen
	temporaryPDUSessionSetupData := createPDUSessionEvt.TempPDUSessionSetupData

	n3iwfCtx := s.Context()
	ikeSecurityAssociation, _ := n3iwfCtx.IKESALoad(localSPI)

	ikeSecurityAssociation.IkeUE.PduSessionListLen = pduSessionListLen

	s.CreatePDUSessionChildSA(ikeSecurityAssociation.IkeUE, temporaryPDUSessionSetupData)
}

func (s *Server) HandleIKEDeleteEvt(ikeEvt n3iwf_context.IkeEvt) {
	ikeLog := logger.IKELog
	ikeLog.Infof("Handle IKEDeleteRequest event")

	n3iwfCtx := s.Context()
	ikeDeleteRequest := ikeEvt.(*n3iwf_context.IKEDeleteRequestEvt)
	localSPI := ikeDeleteRequest.LocalSPI

	SendIKEDeleteRequest(n3iwfCtx, localSPI)

	// In normal case, should wait response and then remove ikeUe.
	// Remove ikeUe here to prevent no response received.
	// Even response replied, it will be discarded.
	err := s.removeIkeUe(localSPI)
	if err != nil {
		ikeLog.Errorf("HandleIKEDeleteEvt(): %v", err)
	}
}

func (s *Server) removeIkeUe(localSPI uint64) error {
	n3iwfCtx := s.Context()
	ikeUe, ok := n3iwfCtx.IkeUePoolLoad(localSPI)
	if !ok {
		return errors.Errorf("Cannot get IkeUE from SPI : %016x", localSPI)
	}
	err := ikeUe.Remove()
	if err != nil {
		return errors.Wrapf(err, "Delete IkeUe error")
	}
	return nil
}

func (s *Server) HandleSendChildSADeleteRequest(ikeEvt n3iwf_context.IkeEvt) {
	ikeLog := logger.IKELog
	ikeLog.Infof("Handle SendChildSADeleteRequest event")

	sendChildSADeleteRequestEvt := ikeEvt.(*n3iwf_context.SendChildSADeleteRequestEvt)
	localSPI := sendChildSADeleteRequestEvt.LocalSPI
	releaseIdList := sendChildSADeleteRequestEvt.ReleaseIdList

	ikeUe, ok := s.Context().IkeUePoolLoad(localSPI)
	if !ok {
		ikeLog.Errorf("Cannot get IkeUE from SPI : %+v", localSPI)
		return
	}
	SendChildSADeleteRequest(ikeUe, releaseIdList)
}

func (s *Server) HandleIKEContextUpdate(ikeEvt n3iwf_context.IkeEvt) {
	ikeLog := logger.IKELog
	ikeLog.Infof("Handle IKEContextUpdate event")

	ikeContextUpdateEvt := ikeEvt.(*n3iwf_context.IKEContextUpdateEvt)
	localSPI := ikeContextUpdateEvt.LocalSPI
	kn3iwf := ikeContextUpdateEvt.Kn3iwf

	ikeUe, ok := s.Context().IkeUePoolLoad(localSPI)
	if !ok {
		ikeLog.Errorf("Cannot get IkeUE from SPI : %+v", localSPI)
		return
	}

	if kn3iwf != nil {
		ikeUe.Kn3iwf = kn3iwf
	}
}

func (s *Server) HandleGetNGAPContextResponse(ikeEvt n3iwf_context.IkeEvt) {
	ikeLog := logger.IKELog
	ikeLog.Infof("Handle GetNGAPContextResponse event")

	getNGAPContextRepEvt := ikeEvt.(*n3iwf_context.GetNGAPContextRepEvt)
	localSPI := getNGAPContextRepEvt.LocalSPI
	ngapCxtReqNumlist := getNGAPContextRepEvt.NgapCxtReqNumlist
	ngapCxt := getNGAPContextRepEvt.NgapCxt

	n3iwfCtx := s.Context()
	ikeSecurityAssociation, _ := n3iwfCtx.IKESALoad(localSPI)

	var tempPDUSessionSetupData *n3iwf_context.PDUSessionSetupTemporaryData

	for i, num := range ngapCxtReqNumlist {
		switch num {
		case n3iwf_context.CxtTempPDUSessionSetupData:
			tempPDUSessionSetupData = ngapCxt[i].(*n3iwf_context.PDUSessionSetupTemporaryData)
		default:
			ikeLog.Errorf("Receive undefine NGAP Context Request number : %d", num)
		}
	}

	switch ikeSecurityAssociation.State {
	case EndSignalling:
		s.CreatePDUSessionChildSA(ikeSecurityAssociation.IkeUE, tempPDUSessionSetupData)
		ikeSecurityAssociation.State++
		go s.StartDPD(ikeSecurityAssociation.IkeUE)
	case HandleCreateChildSA:
		s.continueCreateChildSA(ikeSecurityAssociation, tempPDUSessionSetupData)
	}
}

func (s *Server) CreatePDUSessionChildSA(
	ikeUe *n3iwf_context.N3IWFIkeUe,
	temporaryPDUSessionSetupData *n3iwf_context.PDUSessionSetupTemporaryData,
) {
	ikeLog := logger.IKELog
	n3iwfCtx := s.Context()
	cfg := s.Config()
	ipsecGwAddr := cfg.GetIPSecGatewayAddr()

	ikeSecurityAssociation := ikeUe.N3IWFIKESecurityAssociation

	ranNgapId, ok := n3iwfCtx.NgapIdLoad(ikeUe.N3IWFIKESecurityAssociation.LocalSPI)
	if !ok {
		ikeLog.Errorf("Cannot get RanNgapId from SPI : %+v",
			ikeUe.N3IWFIKESecurityAssociation.LocalSPI)
		return
	}

	for {
		if len(temporaryPDUSessionSetupData.UnactivatedPDUSession) > temporaryPDUSessionSetupData.Index {
			pduSession := temporaryPDUSessionSetupData.UnactivatedPDUSession[temporaryPDUSessionSetupData.Index]
			pduSessionID := pduSession.Id

			// Send CREATE_CHILD_SA to UE
			var responseIKEPayload ike_message.IKEPayloadContainer
			errStr := n3iwf_context.ErrNil

			responseIKEPayload.Reset()

			// Build SA
			requestSA := responseIKEPayload.BuildSecurityAssociation()

			// Allocate SPI
			var spi uint32
			spiByte := make([]byte, 4)
			for {
				var err error
				buf := make([]byte, 4)
				_, err = rand.Read(buf)
				if err != nil {
					ikeLog.Errorf("CreatePDUSessionChildSA Generate SPI: %v", err)
					return
				}
				randomUint32 := binary.BigEndian.Uint32(buf)
				if _, ok := n3iwfCtx.ChildSA.Load(randomUint32); !ok {
					spi = randomUint32
					break
				}
			}
			binary.BigEndian.PutUint32(spiByte, spi)

			// First Proposal - Proposal No.1
			proposal := requestSA.Proposals.BuildProposal(1, ike_message.TypeESP, spiByte)

			// Encryption transform
			encrTranform, err := encr.ToTransform(ikeSecurityAssociation.EncrInfo)
			if err != nil {
				ikeLog.Errorf("encr ToTransform error: %v", err)
				break
			}

			proposal.EncryptionAlgorithm = append(proposal.EncryptionAlgorithm,
				encrTranform)
			// Integrity transform
			if pduSession.SecurityIntegrity {
				proposal.IntegrityAlgorithm = append(proposal.IntegrityAlgorithm,
					integ.ToTransform(ikeSecurityAssociation.IntegInfo))
			}

			// RFC 7296
			// Diffie-Hellman transform is optional in CREATE_CHILD_SA
			// proposal.DiffieHellmanGroup.BuildTransform(
			// 	ike_message.TypeDiffieHellmanGroup, ike_message.DH_1024_BIT_MODP, nil, nil, nil)

			// ESN transform
			proposal.ExtendedSequenceNumbers.BuildTransform(
				ike_message.TypeExtendedSequenceNumbers, ike_message.ESN_DISABLE, nil, nil, nil)

			ikeUe.CreateHalfChildSA(ikeSecurityAssociation.ResponderMessageID, spi, pduSessionID)

			// Build Nonce
			nonceDataBigInt, errGen := ike_security.GenerateRandomNumber()
			if errGen != nil {
				ikeLog.Errorf("CreatePDUSessionChildSA Build Nonce: %v", errGen)
				return
			}
			nonceData := nonceDataBigInt.Bytes()
			responseIKEPayload.BuildNonce(nonceData)

			// Store nonce into context
			ikeSecurityAssociation.ConcatenatedNonce = nonceData

			// TSi
			n3iwfIPAddr := net.ParseIP(ipsecGwAddr)
			tsi := responseIKEPayload.BuildTrafficSelectorInitiator()
			tsi.TrafficSelectors.BuildIndividualTrafficSelector(
				ike_message.TS_IPV4_ADDR_RANGE, ike_message.IPProtocolAll,
				0, 65535, n3iwfIPAddr.To4(), n3iwfIPAddr.To4())

			// TSr
			ueIPAddr := ikeUe.IPSecInnerIP
			tsr := responseIKEPayload.BuildTrafficSelectorResponder()
			tsr.TrafficSelectors.BuildIndividualTrafficSelector(
				ike_message.TS_IPV4_ADDR_RANGE, ike_message.IPProtocolAll,
				0, 65535, ueIPAddr.To4(), ueIPAddr.To4())

			if pduSessionID < 0 || pduSessionID > math.MaxUint8 {
				ikeLog.Errorf("CreatePDUSessionChildSA pduSessionID exceeds uint8 range: %d", pduSessionID)
				break
			}
			// Notify-Qos
			err = responseIKEPayload.BuildNotify5G_QOS_INFO(uint8(pduSessionID), pduSession.QFIList, true, false, 0)
			if err != nil {
				ikeLog.Errorf("CreatePDUSessionChildSA error : %v", err)
				break
			}

			// Notify-UP_IP_ADDRESS
			responseIKEPayload.BuildNotifyUP_IP4_ADDRESS(ipsecGwAddr)

			temporaryPDUSessionSetupData.Index++

			// Build IKE message
			ikeMessage := ike_message.NewMessage(ikeSecurityAssociation.RemoteSPI, ikeSecurityAssociation.LocalSPI,
				ike_message.CREATE_CHILD_SA, false, false, ikeSecurityAssociation.ResponderMessageID,
				responseIKEPayload)

			err = SendIKEMessageToUE(ikeSecurityAssociation.IKEConnection.Conn,
				ikeSecurityAssociation.IKEConnection.N3IWFAddr,
				ikeSecurityAssociation.IKEConnection.UEAddr, ikeMessage,
				ikeSecurityAssociation.IKESAKey)
			if err != nil {
				ikeLog.Errorf("CreatePDUSessionChildSA error : %v", err)
				errStr = n3iwf_context.ErrTransportResourceUnavailable
				temporaryPDUSessionSetupData.FailedErrStr = append(temporaryPDUSessionSetupData.FailedErrStr,
					errStr)
			} else {
				temporaryPDUSessionSetupData.FailedErrStr = append(temporaryPDUSessionSetupData.FailedErrStr,
					errStr)
				break
			}
		} else {
			s.SendNgapEvt(n3iwf_context.NewSendPDUSessionResourceSetupResEvt(ranNgapId))
			break
		}
	}
}

func (s *Server) StartDPD(ikeUe *n3iwf_context.N3IWFIkeUe) {
	ikeLog := logger.IKELog
	defer func() {
		if p := recover(); p != nil {
			// Print stack for panic to log. Fatalf() will let program exit.
			ikeLog.Errorf("panic: %v\n%s", p, string(debug.Stack()))
		}
	}()

	ikeUe.N3IWFIKESecurityAssociation.IKESAClosedCh = make(chan struct{})

	n3iwfCtx := s.Context()
	cfg := s.Config()
	ikeSA := ikeUe.N3IWFIKESecurityAssociation

	liveness := cfg.GetLivenessCheck()
	if liveness.Enable {
		ikeSA.IsUseDPD = true
		timer := time.NewTicker(liveness.TransFreq)
		for {
			select {
			case <-ikeSA.IKESAClosedCh:
				close(ikeSA.IKESAClosedCh)
				timer.Stop()
				return
			case <-timer.C:
				var payload *ike_message.IKEPayloadContainer
				SendUEInformationExchange(ikeSA, ikeSA.IKESAKey, payload, false, false,
					ikeSA.ResponderMessageID, ikeUe.IKEConnection.Conn, ikeUe.IKEConnection.UEAddr,
					ikeUe.IKEConnection.N3IWFAddr)

				DPDReqRetransTime := 2 * time.Second // TODO: make it configurable
				ikeSA.DPDReqRetransTimer = n3iwf_context.NewDPDPeriodicTimer(
					DPDReqRetransTime, liveness.MaxRetryTimes, ikeSA,
					func() {
						ikeLog.Errorf("UE is down")
						ranNgapId, ok := n3iwfCtx.NgapIdLoad(ikeSA.LocalSPI)
						if !ok {
							ikeLog.Infof("Cannot find ranNgapId form SPI : %+v",
								ikeSA.LocalSPI)
							return
						}

						s.SendNgapEvt(n3iwf_context.NewSendUEContextReleaseRequestEvt(
							ranNgapId, n3iwf_context.ErrRadioConnWithUeLost,
						))

						ikeSA.DPDReqRetransTimer = nil
						timer.Stop()
					})
			}
		}
	}
}

func (s *Server) handleNATDetect(
	initiatorSPI, responderSPI uint64,
	notifications []*ike_message.Notification,
	ueAddr, n3iwfAddr *net.UDPAddr,
) (bool, bool, error) {
	ikeLog := logger.IKELog
	ueBehindNAT := false
	n3iwfBehindNAT := false

	srcNatDData, err := s.generateNATDetectHash(initiatorSPI, responderSPI, ueAddr)
	if err != nil {
		return false, false, errors.Wrapf(err, "handle NATD")
	}

	dstNatDData, err := s.generateNATDetectHash(initiatorSPI, responderSPI, n3iwfAddr)
	if err != nil {
		return false, false, errors.Wrapf(err, "handle NATD")
	}

	for _, notification := range notifications {
		switch notification.NotifyMessageType {
		case ike_message.NAT_DETECTION_SOURCE_IP:
			ikeLog.Tracef("Received IKE Notify: NAT_DETECTION_SOURCE_IP")
			if !bytes.Equal(notification.NotificationData, srcNatDData) {
				ikeLog.Tracef("UE(SPI: %016x) is behind NAT", responderSPI)
				ueBehindNAT = true
			}
		case ike_message.NAT_DETECTION_DESTINATION_IP:
			ikeLog.Tracef("Received IKE Notify: NAT_DETECTION_DESTINATION_IP")
			if !bytes.Equal(notification.NotificationData, dstNatDData) {
				ikeLog.Tracef("N3IWF is behind NAT")
				n3iwfBehindNAT = true
			}
		default:
		}
	}
	return ueBehindNAT, n3iwfBehindNAT, nil
}

func (s *Server) generateNATDetectHash(
	initiatorSPI, responderSPI uint64,
	addr *net.UDPAddr,
) ([]byte, error) {
	// Calculate NAT_DETECTION hash for NAT-T
	// : sha1(ispi | rspi | ip | port)
	natdData := make([]byte, 22)
	binary.BigEndian.PutUint64(natdData[0:8], initiatorSPI)
	binary.BigEndian.PutUint64(natdData[8:16], responderSPI)
	copy(natdData[16:20], addr.IP.To4())
	binary.BigEndian.PutUint16(natdData[20:22], uint16(addr.Port)) // #nosec G115

	sha1HashFunction := sha1.New() // #nosec G401
	_, err := sha1HashFunction.Write(natdData)
	if err != nil {
		return nil, errors.Wrapf(err, "generate NATD Hash")
	}
	return sha1HashFunction.Sum(nil), nil
}

func (s *Server) buildNATDetectNotifPayload(
	ikeSA *n3iwf_context.IKESecurityAssociation,
	payload *ike_message.IKEPayloadContainer,
	ueAddr, n3iwfAddr *net.UDPAddr,
) error {
	srcNatDHash, err := s.generateNATDetectHash(ikeSA.RemoteSPI, ikeSA.LocalSPI, n3iwfAddr)
	if err != nil {
		return errors.Wrapf(err, "build NATD")
	}
	// Build and append notify payload for NAT_DETECTION_SOURCE_IP
	payload.BuildNotification(
		ike_message.TypeNone, ike_message.NAT_DETECTION_SOURCE_IP, nil, srcNatDHash)

	dstNatDHash, err := s.generateNATDetectHash(ikeSA.RemoteSPI, ikeSA.LocalSPI, ueAddr)
	if err != nil {
		return errors.Wrapf(err, "build NATD")
	}
	// Build and append notify payload for NAT_DETECTION_DESTINATION_IP
	payload.BuildNotification(
		ike_message.TypeNone, ike_message.NAT_DETECTION_DESTINATION_IP, nil, dstNatDHash)

	return nil
}

func (s *Server) handleDeletePayload(payload *ike_message.Delete, isResponse bool,
	ikeSecurityAssociation *n3iwf_context.IKESecurityAssociation) (
	*ike_message.IKEPayloadContainer, error,
) {
	var evt n3iwf_context.NgapEvt
	var err error
	n3iwfCtx := s.Context()
	n3iwfIke := ikeSecurityAssociation.IkeUE
	responseIKEPayload := new(ike_message.IKEPayloadContainer)

	ranNgapId, ok := n3iwfCtx.NgapIdLoad(n3iwfIke.N3IWFIKESecurityAssociation.LocalSPI)
	if !ok {
		return nil, errors.Errorf("handleDeletePayload: Cannot get RanNgapId from SPI : %+v",
			n3iwfIke.N3IWFIKESecurityAssociation.LocalSPI)
	}

	switch payload.ProtocolID {
	case ike_message.TypeIKE:
		if !isResponse {
			err = n3iwfIke.Remove()
			if err != nil {
				return nil, errors.Wrapf(err, "handleDeletePayload: Delete IkeUe Context error")
			}
		}

		evt = n3iwf_context.NewSendUEContextReleaseEvt(ranNgapId)
	case ike_message.TypeESP:
		var deletSPIs []uint32
		var deletPduIds []int64
		if !isResponse {
			deletSPIs, deletPduIds, err = s.deleteChildSAFromSPIList(n3iwfIke, payload.SPIs)
			if err != nil {
				return nil, errors.Wrapf(err, "handleDeletePayload")
			}
			responseIKEPayload.BuildDeletePayload(ike_message.TypeESP, 4, uint16(len(deletSPIs)), deletSPIs)
		}

		evt = n3iwf_context.NewendPDUSessionResourceReleaseEvt(ranNgapId, deletPduIds)
	default:
		return nil, errors.Errorf("Get Protocol ID %d in Informational delete payload, "+
			"this payload will not be handled by IKE handler", payload.ProtocolID)
	}
	s.SendNgapEvt(evt)
	return responseIKEPayload, nil
}

func isTransformKernelSupported(
	transformType uint8,
	transformID uint16,
	attributePresent bool,
	attributeValue uint16,
) bool {
	switch transformType {
	case ike_message.TypeEncryptionAlgorithm:
		switch transformID {
		case ike_message.ENCR_DES_IV64:
			return false
		case ike_message.ENCR_DES:
			return true
		case ike_message.ENCR_3DES:
			return true
		case ike_message.ENCR_RC5:
			return false
		case ike_message.ENCR_IDEA:
			return false
		case ike_message.ENCR_CAST:
			if attributePresent {
				switch attributeValue {
				case 128:
					return true
				case 256:
					return false
				default:
					return false
				}
			} else {
				return false
			}
		case ike_message.ENCR_BLOWFISH:
			return true
		case ike_message.ENCR_3IDEA:
			return false
		case ike_message.ENCR_DES_IV32:
			return false
		case ike_message.ENCR_NULL:
			return true
		case ike_message.ENCR_AES_CBC:
			if attributePresent {
				switch attributeValue {
				case 128:
					return true
				case 192:
					return true
				case 256:
					return true
				default:
					return false
				}
			} else {
				return false
			}
		case ike_message.ENCR_AES_CTR:
			if attributePresent {
				switch attributeValue {
				case 128:
					return true
				case 192:
					return true
				case 256:
					return true
				default:
					return false
				}
			} else {
				return false
			}
		default:
			return false
		}
	case ike_message.TypeIntegrityAlgorithm:
		switch transformID {
		case ike_message.AUTH_NONE:
			return false
		case ike_message.AUTH_HMAC_MD5_96:
			return true
		case ike_message.AUTH_HMAC_SHA1_96:
			return true
		case ike_message.AUTH_DES_MAC:
			return false
		case ike_message.AUTH_KPDK_MD5:
			return false
		case ike_message.AUTH_AES_XCBC_96:
			return true
		case ike_message.AUTH_HMAC_SHA2_256_128:
			return true
		default:
			return false
		}
	case ike_message.TypeDiffieHellmanGroup:
		switch transformID {
		// case ike_message.DH_NONE:
		// 	return false
		// case ike_message.DH_768_BIT_MODP:
		// 	return false
		// case ike_message.DH_1024_BIT_MODP:
		// 	return false
		// case ike_message.DH_1536_BIT_MODP:
		// 	return false
		// case ike_message.DH_2048_BIT_MODP:
		// 	return false
		// case ike_message.DH_3072_BIT_MODP:
		// 	return false
		// case ike_message.DH_4096_BIT_MODP:
		// 	return false
		// case ike_message.DH_6144_BIT_MODP:
		// 	return false
		// case ike_message.DH_8192_BIT_MODP:
		// 	return false
		default:
			return false
		}
	case ike_message.TypeExtendedSequenceNumbers:
		switch transformID {
		case ike_message.ESN_ENABLE:
			return true
		case ike_message.ESN_DISABLE:
			return true
		default:
			return false
		}
	default:
		return false
	}
}

func (s *Server) parseIPAddressInformationToChildSecurityAssociation(
	childSecurityAssociation *n3iwf_context.ChildSecurityAssociation,
	uePublicIPAddr net.IP,
	trafficSelectorLocal *ike_message.IndividualTrafficSelector,
	trafficSelectorRemote *ike_message.IndividualTrafficSelector,
) error {
	ikeLog := logger.IKELog
	if childSecurityAssociation == nil {
		return errors.New("childSecurityAssociation is nil")
	}

	n3iwfCtx := s.Context()
	cfg := n3iwfCtx.Config()
	childSecurityAssociation.PeerPublicIPAddr = uePublicIPAddr
	childSecurityAssociation.LocalPublicIPAddr = net.ParseIP(cfg.GetIKEBindAddr())

	ikeLog.Tracef("Local TS: %+v", trafficSelectorLocal.StartAddress)
	ikeLog.Tracef("Remote TS: %+v", trafficSelectorRemote.StartAddress)

	childSecurityAssociation.TrafficSelectorLocal = net.IPNet{
		IP:   trafficSelectorLocal.StartAddress,
		Mask: []byte{255, 255, 255, 255},
	}

	childSecurityAssociation.TrafficSelectorRemote = net.IPNet{
		IP:   trafficSelectorRemote.StartAddress,
		Mask: []byte{255, 255, 255, 255},
	}

	return nil
}

func SelectProposal(proposals ike_message.ProposalContainer) ike_message.ProposalContainer {
	var chooseProposal ike_message.ProposalContainer

	for _, proposal := range proposals {
		// We need ENCR, PRF, INTEG, DH, but not ESN

		var encryptionAlgorithmTransform, pseudorandomFunctionTransform *ike_message.Transform
		var integrityAlgorithmTransform, diffieHellmanGroupTransform *ike_message.Transform
		var chooseDH dh.DHType
		var chooseEncr encr.ENCRType
		var chooseInte integ.INTEGType
		var choosePrf prf.PRFType

		for _, transform := range proposal.DiffieHellmanGroup {
			dhType := dh.DecodeTransform(transform)
			if dhType != nil {
				if diffieHellmanGroupTransform == nil {
					diffieHellmanGroupTransform = transform
					chooseDH = dhType
				}
			}
		}
		if chooseDH == nil {
			continue // mandatory
		}

		for _, transform := range proposal.EncryptionAlgorithm {
			encrType := encr.DecodeTransform(transform)
			if encrType != nil {
				if encryptionAlgorithmTransform == nil {
					encryptionAlgorithmTransform = transform
					chooseEncr = encrType
				}
			}
		}
		if chooseEncr == nil {
			continue // mandatory
		}

		for _, transform := range proposal.IntegrityAlgorithm {
			integType := integ.DecodeTransform(transform)
			if integType != nil {
				if integrityAlgorithmTransform == nil {
					integrityAlgorithmTransform = transform
					chooseInte = integType
				}
			}
		}
		if chooseInte == nil {
			continue // mandatory
		}

		for _, transform := range proposal.PseudorandomFunction {
			prfType := prf.DecodeTransform(transform)
			if prfType != nil {
				if pseudorandomFunctionTransform == nil {
					pseudorandomFunctionTransform = transform
					choosePrf = prfType
				}
			}
		}
		if choosePrf == nil {
			continue // mandatory
		}
		if len(proposal.ExtendedSequenceNumbers) > 0 {
			continue // No ESN
		}

		// Construct chosen proposal, with ENCR, PRF, INTEG, DH, and each
		// contains one transform expectively
		chosenProposal := chooseProposal.BuildProposal(proposal.ProposalNumber, proposal.ProtocolID, nil)
		chosenProposal.EncryptionAlgorithm = append(chosenProposal.EncryptionAlgorithm, encryptionAlgorithmTransform)
		chosenProposal.IntegrityAlgorithm = append(chosenProposal.IntegrityAlgorithm, integrityAlgorithmTransform)
		chosenProposal.PseudorandomFunction = append(chosenProposal.PseudorandomFunction, pseudorandomFunctionTransform)
		chosenProposal.DiffieHellmanGroup = append(chosenProposal.DiffieHellmanGroup, diffieHellmanGroupTransform)
		break
	}
	return chooseProposal
}

func (s *Server) deleteChildSAFromSPIList(ikeUe *n3iwf_context.N3IWFIkeUe, spiList []uint32) (
	[]uint32, []int64, error,
) {
	ikeLog := logger.IKELog
	var deleteSPIs []uint32
	var deletePduIds []int64

	for _, spi := range spiList {
		found := false
		for _, childSA := range ikeUe.N3IWFChildSecurityAssociation {
			if childSA.OutboundSPI == spi {
				found = true
				deleteSPIs = append(deleteSPIs, childSA.InboundSPI)

				if len(childSA.PDUSessionIds) == 0 {
					return nil, nil, errors.Errorf("Child_SA SPI: 0x%08x doesn't have PDU Session ID",
						spi)
				}
				deletePduIds = append(deletePduIds, childSA.PDUSessionIds[0])

				err := ikeUe.DeleteChildSA(childSA)
				if err != nil {
					return nil, nil, errors.Wrapf(err, "DeleteChildSAFromSPIList")
				}
				break
			}
		}
		if !found {
			ikeLog.Warnf("deleteChildSAFromSPIList(): Get unknown Child_SA with SPI: 0x%08x", spi)
		}
	}

	return deleteSPIs, deletePduIds, nil
}
