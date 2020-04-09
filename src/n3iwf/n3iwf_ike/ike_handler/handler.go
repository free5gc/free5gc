package ike_handler

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"net"
	"strings"

	"free5gc/lib/ngap/ngapType"
	"free5gc/src/n3iwf/logger"
	"free5gc/src/n3iwf/n3iwf_context"
	"free5gc/src/n3iwf/n3iwf_data_relay"
	"free5gc/src/n3iwf/n3iwf_handler/n3iwf_message"
	"free5gc/src/n3iwf/n3iwf_ike/ike_message"
	"free5gc/src/n3iwf/n3iwf_ngap/ngap_message"

	"github.com/sirupsen/logrus"
	"golang.org/x/sys/unix"
)

// Log
var ikeLog *logrus.Entry

func init() {
	ikeLog = logger.IKELog
}

func HandleIKESAINIT(ueSendInfo *n3iwf_message.UDPSendInfoGroup, message *ike_message.IKEMessage) {
	ikeLog.Infoln("[IKE] Handle IKE_SA_INIT")

	var securityAssociation *ike_message.SecurityAssociation
	var keyExcahge *ike_message.KeyExchange
	var nonce *ike_message.Nonce

	n3iwfSelf := n3iwf_context.N3IWFSelf()

	var responseIKEMessage *ike_message.IKEMessage
	var responseSecurityAssociation *ike_message.SecurityAssociation
	var responseKeyExchange *ike_message.KeyExchange
	var responseNonce *ike_message.Nonce

	var sharedKeyData, concatenatedNonce []byte

	if message == nil {
		ikeLog.Error("[IKE] IKE Message is nil")
		return
	}

	// parse IKE header and setup IKE context
	// check major version
	majorVersion := ((message.Version & 0xf0) >> 4)
	if majorVersion > 2 {
		ikeLog.Warn("[IKE] Received an IKE message with higher major version")
		// send INFORMATIONAL type message with INVALID_MAJOR_VERSION Notify payload
		responseIKEMessage = ike_message.BuildIKEHeader(message.InitiatorSPI, message.ResponderSPI, ike_message.INFORMATIONAL, ike_message.ResponseBitCheck, message.MessageID)
		notificationPayload := ike_message.BuildNotification(ike_message.TypeNone, ike_message.INVALID_MAJOR_VERSION, nil, nil)

		responseIKEMessage.IKEPayload = append(responseIKEMessage.IKEPayload, notificationPayload)

		SendIKEMessageToUE(ueSendInfo, responseIKEMessage)

		return
	}

	for _, ikePayload := range message.IKEPayload {
		switch ikePayload.Type() {
		case ike_message.TypeSA:
			securityAssociation = ikePayload.(*ike_message.SecurityAssociation)
		case ike_message.TypeKE:
			keyExcahge = ikePayload.(*ike_message.KeyExchange)
		case ike_message.TypeNiNr:
			nonce = ikePayload.(*ike_message.Nonce)
		default:
			ikeLog.Warnf("[IKE] Get IKE payload (type %d) in IKE_SA_INIT message, this payload will not be handled by IKE handler", ikePayload.Type())
		}
	}

	if securityAssociation != nil {
		for _, proposal := range securityAssociation.Proposals {
			chosenProposal := new(ike_message.Proposal)

			if len(proposal.EncryptionAlgorithm) > 0 {
				for _, transform := range proposal.EncryptionAlgorithm {
					if is_supported(ike_message.TypeEncryptionAlgorithm, transform.TransformID, transform.AttributePresent, transform.AttributeValue) {
						chosenProposal.EncryptionAlgorithm = append(chosenProposal.EncryptionAlgorithm, transform)
						break
					}
				}
				if len(chosenProposal.EncryptionAlgorithm) == 0 {
					continue
				}
			} else {
				continue // mandatory
			}
			if len(proposal.PseudorandomFunction) > 0 {
				for _, transform := range proposal.PseudorandomFunction {
					if is_supported(ike_message.TypePseudorandomFunction, transform.TransformID, transform.AttributePresent, transform.AttributeValue) {
						chosenProposal.PseudorandomFunction = append(chosenProposal.PseudorandomFunction, transform)
						break
					}
				}
				if len(chosenProposal.PseudorandomFunction) == 0 {
					continue
				}
			} else {
				continue // mandatory
			}
			if len(proposal.IntegrityAlgorithm) > 0 {
				for _, transform := range proposal.IntegrityAlgorithm {
					if is_supported(ike_message.TypeIntegrityAlgorithm, transform.TransformID, transform.AttributePresent, transform.AttributeValue) {
						chosenProposal.IntegrityAlgorithm = append(chosenProposal.IntegrityAlgorithm, transform)
						break
					}
				}
				if len(chosenProposal.IntegrityAlgorithm) == 0 {
					continue
				}
			} else {
				continue // mandatory
			}
			if len(proposal.DiffieHellmanGroup) > 0 {
				for _, transform := range proposal.DiffieHellmanGroup {
					if is_supported(ike_message.TypeDiffieHellmanGroup, transform.TransformID, transform.AttributePresent, transform.AttributeValue) {
						chosenProposal.DiffieHellmanGroup = append(chosenProposal.DiffieHellmanGroup, transform)
						break
					}
				}
				if len(chosenProposal.DiffieHellmanGroup) == 0 {
					continue
				}
			} else {
				continue // mandatory
			}
			if len(proposal.ExtendedSequenceNumbers) > 0 {
				continue // No ESN
			}

			chosenProposal.ProposalNumber = proposal.ProposalNumber
			chosenProposal.ProtocolID = proposal.ProtocolID

			responseSecurityAssociation = &ike_message.SecurityAssociation{
				Proposals: []*ike_message.Proposal{
					chosenProposal,
				},
			}

			break
		}

		if responseSecurityAssociation == nil {
			ikeLog.Warn("[IKE] No proposal chosen")
			// Respond NO_PROPOSAL_CHOSEN to UE
			responseIKEMessage = ike_message.BuildIKEHeader(message.InitiatorSPI, message.ResponderSPI, ike_message.IKE_SA_INIT, ike_message.ResponseBitCheck, message.MessageID)
			notificationPayload := ike_message.BuildNotification(ike_message.TypeNone, ike_message.NO_PROPOSAL_CHOSEN, nil, nil)

			responseIKEMessage.IKEPayload = append(responseIKEMessage.IKEPayload, notificationPayload)

			SendIKEMessageToUE(ueSendInfo, responseIKEMessage)

			return
		}
	} else {
		ikeLog.Error("[IKE] The security association field is nil")
		// TODO: send error message to UE
		return
	}

	if keyExcahge != nil {
		chosenDiffieHellmanGroup := responseSecurityAssociation.Proposals[0].DiffieHellmanGroup[0].TransformID
		if chosenDiffieHellmanGroup != keyExcahge.DiffieHellmanGroup {
			ikeLog.Warn("[IKE] The Diffie-Hellman group defined in key exchange payload not matches the one in chosen proposal")
			// send INVALID_KE_PAYLOAD to UE
			responseIKEMessage = ike_message.BuildIKEHeader(message.InitiatorSPI, message.ResponderSPI, ike_message.IKE_SA_INIT, ike_message.ResponseBitCheck, message.MessageID)

			notificationData := make([]byte, 2)
			binary.BigEndian.PutUint16(notificationData, chosenDiffieHellmanGroup)

			notificationPayload := ike_message.BuildNotification(ike_message.TypeNone, ike_message.INVALID_KE_PAYLOAD, nil, notificationData)

			responseIKEMessage.IKEPayload = append(responseIKEMessage.IKEPayload, notificationPayload)

			SendIKEMessageToUE(ueSendInfo, responseIKEMessage)

			return
		}

		var localPublicValue []byte

		localPublicValue, sharedKeyData = CalculateDiffieHellmanMaterials(GenerateRandomNumber(), keyExcahge.KeyExchangeData, chosenDiffieHellmanGroup)
		responseKeyExchange = ike_message.BUildKeyExchange(chosenDiffieHellmanGroup, localPublicValue)
	} else {
		ikeLog.Error("[IKE] The key exchange field is nil")
		// TODO: send error message to UE
		return
	}

	if nonce != nil {
		localNonce := GenerateRandomNumber().Bytes()
		concatenatedNonce = append(nonce.NonceData, localNonce...)

		responseNonce = ike_message.BuildNonce(localNonce)
	} else {
		ikeLog.Error("[IKE] The nonce field is nil")
		// TODO: send error message to UE
		return
	}

	// Create new IKE security association
	ikeSecurityAssociation := n3iwfSelf.NewIKESecurityAssociation()
	ikeSecurityAssociation.RemoteSPI = message.InitiatorSPI

	// Record algorithm in context
	ikeSecurityAssociation.EncryptionAlgorithm = responseSecurityAssociation.Proposals[0].EncryptionAlgorithm[0]
	ikeSecurityAssociation.IntegrityAlgorithm = responseSecurityAssociation.Proposals[0].IntegrityAlgorithm[0]
	ikeSecurityAssociation.PseudorandomFunction = responseSecurityAssociation.Proposals[0].PseudorandomFunction[0]
	ikeSecurityAssociation.DiffieHellmanGroup = responseSecurityAssociation.Proposals[0].DiffieHellmanGroup[0]

	// Record concatenated nonce
	ikeSecurityAssociation.ConcatenatedNonce = append(ikeSecurityAssociation.ConcatenatedNonce, concatenatedNonce...)
	// Record Diffie-Hellman shared key
	ikeSecurityAssociation.DiffieHellmanSharedKey = append(ikeSecurityAssociation.DiffieHellmanSharedKey, sharedKeyData...)

	if err := GenerateKeyForIKESA(ikeSecurityAssociation); err != nil {
		ikeLog.Errorf("Generate key for IKE SA failed: %+v", err)
		return
	}

	// Send response to UE
	responseIKEMessage = ike_message.BuildIKEHeader(ikeSecurityAssociation.RemoteSPI, ikeSecurityAssociation.LocalSPI, ike_message.IKE_SA_INIT, ike_message.ResponseBitCheck, message.MessageID)
	responseIKEMessage.IKEPayload = append(responseIKEMessage.IKEPayload, responseSecurityAssociation)
	responseIKEMessage.IKEPayload = append(responseIKEMessage.IKEPayload, responseKeyExchange)
	responseIKEMessage.IKEPayload = append(responseIKEMessage.IKEPayload, responseNonce)

	// Prepare authentication data - InitatorSignedOctet
	// InitatorSignedOctet = RealMessage1 | NonceRData | MACedIDForI
	// MACedIDForI is acquired in IKE_AUTH exchange
	receivedIKEMessageData, err := ike_message.Encode(message)
	if err != nil {
		ikeLog.Errorln(err)
		ikeLog.Error("[IKE] Encode message failed.")
		return
	}
	ikeSecurityAssociation.RemoteUnsignedAuthentication = append(receivedIKEMessageData, responseNonce.NonceData...)

	// Prepare authentication data - ResponderSignedOctet
	// ResponderSignedOctet = RealMessage2 | NonceIData | MACedIDForR
	responseIKEMessageData, err := ike_message.Encode(responseIKEMessage)
	if err != nil {
		ikeLog.Errorln(err)
		ikeLog.Error("[IKE] Encoding IKE message failed")
		return
	}
	ikeSecurityAssociation.LocalUnsignedAuthentication = append(responseIKEMessageData, nonce.NonceData...)
	// MACedIDForR
	idPayload := []ike_message.IKEPayloadType{
		ike_message.BuildIdentificationResponder(ike_message.ID_FQDN, []byte(n3iwfSelf.FQDN)),
	}
	idPayloadData, err := ike_message.EncodePayload(idPayload)
	if err != nil {
		ikeLog.Errorln(err)
		ikeLog.Error("[IKE] Encode IKE payload failed.")
		return
	}
	pseudorandomFunction, ok := NewPseudorandomFunction(ikeSecurityAssociation.SK_pr, ikeSecurityAssociation.PseudorandomFunction.TransformID)
	if !ok {
		ikeLog.Error("[IKE] Get an unsupported pseudorandom funcion. This may imply an unsupported transform is chosen.")
		return
	}
	if _, err := pseudorandomFunction.Write(idPayloadData[4:]); err != nil {
		ikeLog.Errorf("[IKE] Pseudorandom function write error: %+v", err)
		return
	}
	ikeSecurityAssociation.LocalUnsignedAuthentication = append(ikeSecurityAssociation.LocalUnsignedAuthentication, pseudorandomFunction.Sum(nil)...)

	ikeLog.Tracef("Local unsigned authentication data:\n%s", hex.Dump(ikeSecurityAssociation.LocalUnsignedAuthentication))

	SendIKEMessageToUE(ueSendInfo, responseIKEMessage)
}

// IKE_AUTH state
const (
	PreSignalling = iota
	EAPSignalling
	PostSignalling
)

func HandleIKEAUTH(ueSendInfo *n3iwf_message.UDPSendInfoGroup, message *ike_message.IKEMessage) {
	ikeLog.Infoln("[IKE] Handle IKE_AUTH")

	var encryptedPayload *ike_message.Encrypted

	n3iwfSelf := n3iwf_context.N3IWFSelf()

	// {response}
	var responseIKEMessage *ike_message.IKEMessage

	if message == nil {
		ikeLog.Error("[IKE] IKE Message is nil")
		return
	}

	// parse IKE header and setup IKE context
	// check major version
	majorVersion := ((message.Version & 0xf0) >> 4)
	if majorVersion > 2 {
		ikeLog.Warn("[IKE] Received an IKE message with higher major version")
		// send INFORMATIONAL type message with INVALID_MAJOR_VERSION Notify payload ( OUTSIDE IKE SA )

		// IKEHDR-{response}
		responseNotification := ike_message.BuildNotification(ike_message.TypeNone, ike_message.INVALID_MAJOR_VERSION, nil, nil)

		responseIKEMessage = ike_message.BuildIKEHeader(message.InitiatorSPI, message.ResponderSPI, ike_message.INFORMATIONAL, ike_message.ResponseBitCheck, message.MessageID)
		responseIKEMessage.IKEPayload = append(responseIKEMessage.IKEPayload, responseNotification)

		SendIKEMessageToUE(ueSendInfo, responseIKEMessage)

		return
	}

	// Find corresponding IKE security association
	localSPI := message.ResponderSPI
	ikeSecurityAssociation := n3iwfSelf.FindIKESecurityAssociationBySPI(localSPI)
	if ikeSecurityAssociation == nil {
		ikeLog.Warn("[IKE] Unrecognized SPI")
		// send INFORMATIONAL type message with INVALID_IKE_SPI Notify payload ( OUTSIDE IKE SA )

		// IKEHDR-{response}
		responseNotification := ike_message.BuildNotification(ike_message.TypeNone, ike_message.INVALID_IKE_SPI, nil, nil)

		responseIKEMessage = ike_message.BuildIKEHeader(message.InitiatorSPI, 0, ike_message.INFORMATIONAL, ike_message.ResponseBitCheck, message.MessageID)
		responseIKEMessage.IKEPayload = append(responseIKEMessage.IKEPayload, responseNotification)

		SendIKEMessageToUE(ueSendInfo, responseIKEMessage)

		return
	}

	for _, ikePayload := range message.IKEPayload {
		switch ikePayload.Type() {
		case ike_message.TypeSK:
			encryptedPayload = ikePayload.(*ike_message.Encrypted)
		default:
			ikeLog.Warnf("[IKE] Get IKE payload (type %d) in IKE_SA_INIT message, this payload will not be handled by IKE handler", ikePayload.Type())
		}
	}

	decryptedIKEPayload, err := DecryptProcedure(ikeSecurityAssociation, message, encryptedPayload)
	if err != nil {
		ikeLog.Errorf("[IKE] Decrypt IKE message failed: %+v", err)
		return
	}

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

	for _, ikePayload := range decryptedIKEPayload {
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
			ikeLog.Warnf("[IKE] Get IKE payload (type %d) in IKE_AUTH message, this payload will not be handled by IKE handler", ikePayload.Type())
		}
	}

	// NOTE: tune it
	transformPseudorandomFunction := ikeSecurityAssociation.PseudorandomFunction

	switch ikeSecurityAssociation.State {
	case PreSignalling:
		// IKEHDR-SK-{response}
		var responseIdentification *ike_message.IdentificationResponder
		var responseCertificate *ike_message.Certificate
		var responseAuthentication *ike_message.Authentication
		var requestEAPPayload *ike_message.EAP

		if initiatorID != nil {
			ikeLog.Info("Ecoding initiator for later IKE authentication")
			ikeSecurityAssociation.InitiatorID = initiatorID

			// Record maced identification for authentication
			idPayload := []ike_message.IKEPayloadType{
				initiatorID,
			}
			idPayloadData, err := ike_message.EncodePayload(idPayload)
			if err != nil {
				ikeLog.Errorln(err)
				ikeLog.Error("[IKE] Encoding ID payload message failed.")
				return
			}
			pseudorandomFunction, ok := NewPseudorandomFunction(ikeSecurityAssociation.SK_pr, transformPseudorandomFunction.TransformID)
			if !ok {
				ikeLog.Error("[IKE] Get an unsupported pseudorandom funcion. This may imply an unsupported transform is chosen.")
				return
			}
			if _, err := pseudorandomFunction.Write(idPayloadData[4:]); err != nil {
				ikeLog.Errorf("[IKE] Pseudorandom function write error: %+v", err)
				return
			}
			ikeSecurityAssociation.RemoteUnsignedAuthentication = append(ikeSecurityAssociation.RemoteUnsignedAuthentication, pseudorandomFunction.Sum(nil)...)
		} else {
			ikeLog.Error("[IKE] The initiator identification field is nil")
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
			if CompareRootCertificate(certificateRequest.CertificateEncoding, certificateRequest.CertificationAuthority) {
				responseCertificate = ike_message.BuildCertificate(ike_message.X509CertificateSignature, n3iwfSelf.N3IWFCertificate)
			}
		}

		if certificate != nil {
			ikeLog.Info("UE send its certficate")
			ikeSecurityAssociation.InitiatorCertificate = certificate
		}

		if securityAssociation != nil {
			ikeLog.Info("Parsing security association")
			var chosenSecurityAssociation *ike_message.SecurityAssociation

			for _, proposal := range securityAssociation.Proposals {
				chosenProposal := new(ike_message.Proposal)

				if len(proposal.SPI) != 4 {
					continue // The SPI of ESP must be 32-bit
				}

				// check SPI
				spi := binary.BigEndian.Uint32(proposal.SPI)
				if _, ok := n3iwfSelf.ChildSA[spi]; ok {
					continue
				}

				if len(proposal.EncryptionAlgorithm) > 0 {
					for _, transform := range proposal.EncryptionAlgorithm {
						if is_Kernel_Supported(ike_message.TypeEncryptionAlgorithm, transform.TransformID, transform.AttributePresent, transform.AttributeValue) {
							chosenProposal.EncryptionAlgorithm = append(chosenProposal.EncryptionAlgorithm, transform)
							break
						}
					}
					if len(chosenProposal.EncryptionAlgorithm) == 0 {
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
						if is_Kernel_Supported(ike_message.TypeIntegrityAlgorithm, transform.TransformID, transform.AttributePresent, transform.AttributeValue) {
							chosenProposal.IntegrityAlgorithm = append(chosenProposal.IntegrityAlgorithm, transform)
							break
						}
					}
					if len(chosenProposal.IntegrityAlgorithm) == 0 {
						continue
					}
				} // Optional
				if len(proposal.DiffieHellmanGroup) > 0 {
					for _, transform := range proposal.DiffieHellmanGroup {
						if is_Kernel_Supported(ike_message.TypeDiffieHellmanGroup, transform.TransformID, transform.AttributePresent, transform.AttributeValue) {
							chosenProposal.DiffieHellmanGroup = append(chosenProposal.DiffieHellmanGroup, transform)
							break
						}
					}
					if len(chosenProposal.DiffieHellmanGroup) == 0 {
						continue
					}
				} // Optional
				if len(proposal.ExtendedSequenceNumbers) > 0 {
					for _, transform := range proposal.ExtendedSequenceNumbers {
						if is_Kernel_Supported(ike_message.TypeExtendedSequenceNumbers, transform.TransformID, transform.AttributePresent, transform.AttributeValue) {
							chosenProposal.ExtendedSequenceNumbers = append(chosenProposal.ExtendedSequenceNumbers, transform)
							break
						}
					}
					if len(chosenProposal.ExtendedSequenceNumbers) == 0 {
						continue
					}
				} else {
					continue // Mandatory
				}

				chosenProposal.ProposalNumber = proposal.ProposalNumber
				chosenProposal.ProtocolID = proposal.ProtocolID
				chosenProposal.SPI = append(chosenProposal.SPI, proposal.SPI...)

				chosenSecurityAssociation = &ike_message.SecurityAssociation{
					Proposals: []*ike_message.Proposal{
						chosenProposal,
					},
				}

				break
			}

			if chosenSecurityAssociation == nil {
				ikeLog.Warn("[IKE] No proposal chosen")
				// Respond NO_PROPOSAL_CHOSEN to UE
				// Build IKE message
				responseIKEMessage = ike_message.BuildIKEHeader(message.InitiatorSPI, message.ResponderSPI, ike_message.IKE_AUTH, ike_message.ResponseBitCheck, message.MessageID)

				// Build response
				var ikePayload []ike_message.IKEPayloadType

				// Notification
				notificationPayload := ike_message.BuildNotification(ike_message.TypeNone, ike_message.NO_PROPOSAL_CHOSEN, nil, nil)
				ikePayload = append(ikePayload, notificationPayload)

				if err := EncryptProcedure(ikeSecurityAssociation, ikePayload, responseIKEMessage); err != nil {
					ikeLog.Errorf("Encrypting IKE message failed: %+v", err)
					return
				}

				// Send IKE message to UE
				SendIKEMessageToUE(ueSendInfo, responseIKEMessage)

				return
			}

			ikeSecurityAssociation.IKEAuthResponseSA = chosenSecurityAssociation
		} else {
			ikeLog.Error("[IKE] The security association field is nil")
			// TODO: send error message to UE
			return
		}

		if trafficSelectorInitiator != nil {
			ikeLog.Info("Received traffic selector initiator from UE")
			ikeSecurityAssociation.TrafficSelectorInitiator = trafficSelectorInitiator
		} else {
			ikeLog.Error("[IKE] The initiator traffic selector field is nil")
			// TODO: send error message to UE
			return
		}

		if trafficSelectorResponder != nil {
			ikeLog.Info("Received traffic selector initiator from UE")
			ikeSecurityAssociation.TrafficSelectorResponder = trafficSelectorResponder
		} else {
			ikeLog.Error("[IKE] The initiator traffic selector field is nil")
			// TODO: send error message to UE
			return
		}

		// Build IKE message
		responseIKEMessage = ike_message.BuildIKEHeader(message.InitiatorSPI, message.ResponderSPI, ike_message.IKE_AUTH, ike_message.ResponseBitCheck, message.MessageID)

		// Build response
		var ikePayload []ike_message.IKEPayloadType

		// Identification
		responseIdentification = ike_message.BuildIdentificationResponder(ike_message.ID_FQDN, []byte(n3iwfSelf.FQDN))
		ikePayload = append(ikePayload, responseIdentification)

		// Certificate
		if responseCertificate == nil {
			responseCertificate = ike_message.BuildCertificate(ike_message.X509CertificateSignature, n3iwfSelf.N3IWFCertificate)
		}
		ikePayload = append(ikePayload, responseCertificate)

		// Authentication Data
		ikeLog.Tracef("Local authentication data:\n%s", hex.Dump(ikeSecurityAssociation.LocalUnsignedAuthentication))
		sha1HashFunction := sha1.New()
		if _, err := sha1HashFunction.Write(ikeSecurityAssociation.LocalUnsignedAuthentication); err != nil {
			ikeLog.Errorf("[IKE] Hash function write error: %+v", err)
			return
		}

		signedAuth, err := rsa.SignPKCS1v15(rand.Reader, n3iwfSelf.N3IWFPrivateKey, crypto.SHA1, sha1HashFunction.Sum(nil))
		if err != nil {
			ikeLog.Errorf("[IKE] Sign authentication data failed: %+v", err)
		}

		responseAuthentication = ike_message.BuildAuthentication(ike_message.RSADigitalSignature, signedAuth)
		ikePayload = append(ikePayload, responseAuthentication)

		// EAP expanded 5G-Start
		var identifier uint8
		for {
			identifier, err = GenerateRandomUint8()
			if err != nil {
				ikeLog.Errorf("[IKE] Random number failed: %+v", err)
				return
			}
			if identifier != ikeSecurityAssociation.LastEAPIdentifier {
				ikeSecurityAssociation.LastEAPIdentifier = identifier
				break
			}
		}
		requestEAPPayload = ike_message.BuildEAP5GStart(identifier)
		ikePayload = append(ikePayload, requestEAPPayload)

		if err := EncryptProcedure(ikeSecurityAssociation, ikePayload, responseIKEMessage); err != nil {
			ikeLog.Errorf("Encrypting IKE message failed: %+v", err)
			return
		}

		// Send IKE message to UE
		SendIKEMessageToUE(ueSendInfo, responseIKEMessage)

		// Shift state
		ikeSecurityAssociation.State++

	case EAPSignalling:
		// If success, N3IWF will send an UPLinkNASTransport to AMF
		if eap != nil {
			if eap.Code != ike_message.EAPCodeResponse {
				ikeLog.Error("[IKE][EAP] Received an EAP payload with code other than response. Drop the payload.")
				return
			}
			if eap.Identifier != ikeSecurityAssociation.LastEAPIdentifier {
				ikeLog.Error("[IKE][EAP] Received an EAP payload with unmatched identifier. Drop the payload.")
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
				ikeLog.Errorf("[IKE][EAP] Received EAP packet with type other than EAP expanded type: %d", eapTypeData.Type())
				return
			}

			if eapExpanded.VendorID != ike_message.VendorID3GPP {
				ikeLog.Error("[IKE] The peer sent EAP expended packet with wrong vendor ID. Drop the packet.")
				return
			}
			if eapExpanded.VendorType != ike_message.VendorTypeEAP5G {
				ikeLog.Error("[IKE] The peer sent EAP expanded packet with wrong vendor type. Drop the packet.")
				return
			}

			eap5GMessageID, anParameters, nasPDU, err := UnmarshalEAP5GData(eapExpanded.VendorData)
			if err != nil {
				ikeLog.Errorf("[IKE] Unmarshalling EAP-5G packet failed: %+v", err)
				return
			}

			if eap5GMessageID == ike_message.EAP5GType5GStop {
				// IKEHDR-SK-{response}
				var responseEAP *ike_message.EAP

				// Send EAP failure
				// Build IKE message
				responseIKEMessage = ike_message.BuildIKEHeader(message.InitiatorSPI, message.ResponderSPI, ike_message.IKE_AUTH, ike_message.ResponseBitCheck, message.MessageID)

				// Build response
				var ikePayload []ike_message.IKEPayloadType

				// EAP
				identifier, err := GenerateRandomUint8()
				if err != nil {
					ikeLog.Errorf("[IKE] Generate random uint8 failed: %+v", err)
					return
				}
				responseEAP = ike_message.BuildEAPfailure(identifier)
				ikePayload = append(ikePayload, responseEAP)

				if err := EncryptProcedure(ikeSecurityAssociation, ikePayload, responseIKEMessage); err != nil {
					ikeLog.Errorf("Encrypting IKE message failed: %+v", err)
					return
				}

				// Send IKE message to UE
				SendIKEMessageToUE(ueSendInfo, responseIKEMessage)
				return
			}

			// Send Initial UE Message or Uplink NAS Transport
			if anParameters != nil {
				// AMF selection
				selectedAMF := n3iwfSelf.AMFSelection(anParameters.GUAMI)
				if selectedAMF == nil {
					ikeLog.Warn("[IKE] No avalible AMF for this UE")
					return
				}

				// Create UE context
				ue := n3iwfSelf.NewN3iwfUe()

				// Relative context
				ikeSecurityAssociation.ThisUE = ue
				ue.N3IWFIKESecurityAssociation = ikeSecurityAssociation
				ue.AMF = selectedAMF

				// Store some information in conext
				ikeSecurityAssociation.MessageID = message.MessageID

				ue.UDPSendInfoGroup = ueSendInfo
				networkAddrStringSlice := strings.Split(ueSendInfo.Addr.String(), ":")
				ue.IPAddrv4 = networkAddrStringSlice[0]
				ue.PortNumber = int32(ueSendInfo.Addr.Port)
				ue.RRCEstablishmentCause = int16(anParameters.EstablishmentCause.Value)

				// Send Initial UE Message
				ngap_message.SendInitialUEMessage(selectedAMF, ue, nasPDU)
			} else {
				ue := ikeSecurityAssociation.ThisUE
				amf := ue.AMF

				// Store some information in context
				ikeSecurityAssociation.MessageID = message.MessageID

				ue.UDPSendInfoGroup = ueSendInfo

				// Send Uplink NAS Transport
				ngap_message.SendUplinkNASTransport(amf, ue, nasPDU)
			}
		} else {
			ikeLog.Error("EAP is nil")
		}

	case PostSignalling:
		// Load needed information
		thisUE := ikeSecurityAssociation.ThisUE

		// IKEHDR-SK-{response}
		var responseConfiguration *ike_message.Configuration
		var responseAuthentication *ike_message.Authentication
		var responseSecurityAssociation *ike_message.SecurityAssociation
		var responseTrafficSelectorInitiator *ike_message.TrafficSelectorInitiator
		var responseTrafficSelectorResponder *ike_message.TrafficSelectorResponder
		var responseNotification *ike_message.Notification

		if authentication != nil {
			// Verifying remote AUTH
			pseudorandomFunction, ok := NewPseudorandomFunction(thisUE.Kn3iwf, transformPseudorandomFunction.TransformID)
			if !ok {
				ikeLog.Error("[IKE] Get an unsupported pseudorandom funcion. This may imply an unsupported transform is chosen.")
				return
			}
			if _, err := pseudorandomFunction.Write([]byte("Key Pad for IKEv2")); err != nil {
				ikeLog.Errorf("[IKE] Pseudorandom function write error: %+v", err)
				return
			}
			secret := pseudorandomFunction.Sum(nil)
			pseudorandomFunction, ok = NewPseudorandomFunction(secret, transformPseudorandomFunction.TransformID)
			if !ok {
				ikeLog.Error("[IKE] Get an unsupported pseudorandom funcion. This may imply an unsupported transform is chosen.")
				return
			}
			if _, err := pseudorandomFunction.Write(ikeSecurityAssociation.RemoteUnsignedAuthentication); err != nil {
				ikeLog.Errorf("[IKE] Pseudorandom function write error: %+v", err)
				return
			}
			expectedAuthenticationData := pseudorandomFunction.Sum(nil)

			ikeLog.Tracef("Expected Authentication Data:\n%s", hex.Dump(expectedAuthenticationData))
			// TODO: Finish authentication test for UE and N3IWF
			/*
				if !bytes.Equal(authentication.AuthenticationData, expectedAuthenticationData) {
					ikeLog.Warn("[IKE] Peer authentication failed.")
					// Inform UE the authentication has failed
					// IKEHDR-SK-{response}
					var notification *ike_message.Notification

					// Build IKE message
					responseIKEMessage = ike_message.BuildIKEHeader(message.InitiatorSPI, message.ResponderSPI, ike_message.IKE_AUTH, ike_message.ResponseBitCheck, message.MessageID)

					// Build response
					var ikePayload []ike_message.IKEPayloadType

					// Notification
					notification = ike_message.BuildNotification(ike_message.TypeNone, ike_message.AUTHENTICATION_FAILED, nil, nil)
					ikePayload = append(ikePayload, notification)

					if err := EncryptProcedure(ikeSecurityAssociation, ikePayload, responseIKEMessage); err != nil {
						ikeLog.Errorf("Encrypting IKE message failed: %+v", err)
						return
					}

					// Send IKE message to UE
					SendIKEMessageToUE(ueSendInfo, responseIKEMessage)
					return
				}
			*/
		} else {
			ikeLog.Warn("[IKE] Peer authentication failed.")
			// Inform UE the authentication has failed
			// IKEHDR-SK-{response}
			var responseNotification *ike_message.Notification

			// Build IKE message
			responseIKEMessage = ike_message.BuildIKEHeader(message.InitiatorSPI, message.ResponderSPI, ike_message.IKE_AUTH, ike_message.ResponseBitCheck, message.MessageID)

			// Build response
			var ikePayload []ike_message.IKEPayloadType

			// Notification
			responseNotification = ike_message.BuildNotification(ike_message.TypeNone, ike_message.AUTHENTICATION_FAILED, nil, nil)
			ikePayload = append(ikePayload, responseNotification)

			if err := EncryptProcedure(ikeSecurityAssociation, ikePayload, responseIKEMessage); err != nil {
				ikeLog.Errorf("Encrypting IKE message failed: %+v", err)
				return
			}

			// Send IKE message to UE
			SendIKEMessageToUE(ueSendInfo, responseIKEMessage)
			return
		}

		// Parse configuration request to get if the UE has requested internal address,
		// and prepare configuration payload to UE
		var addrRequest bool = false

		if configuration != nil {
			ikeLog.Tracef("[IKE] Received configuration payload with type: %d", configuration.ConfigurationType)

			var attribute *ike_message.IndividualConfigurationAttribute
			for _, attribute = range configuration.ConfigurationAttribute {
				switch attribute.Type {
				case ike_message.INTERNAL_IP4_ADDRESS:
					addrRequest = true
					if len(attribute.Value) != 0 {
						ikeLog.Tracef("[IKE] Got client requested address: %d.%d.%d.%d", attribute.Value[0], attribute.Value[1], attribute.Value[2], attribute.Value[3])
					}
				default:
					ikeLog.Warnf("[IKE] Receive other type of configuration request: %d", attribute.Type)
				}
			}
		} else {
			ikeLog.Warn("[IKE] Configuration is nil. UE did not sent any configuration request.")
		}

		// Prepare configuration payload and traffic selector payload for initiator and responder
		if addrRequest {
			var attributes []*ike_message.IndividualConfigurationAttribute
			var ueIPAddr net.IP

			n3iwfIPAddr := net.ParseIP(n3iwfSelf.IPSecGatewayAddress)

			// UE internal IP address
			for {
				ueIPAddr = GenerateRandomIPinRange(n3iwfSelf.Subnet)
				if ueIPAddr != nil {
					if ueIPAddr.String() == n3iwfSelf.IPSecGatewayAddress {
						continue
					}
					if _, ok := n3iwfSelf.AllocatedUEIPAddress[ueIPAddr.String()]; !ok {
						// Should be release if there is any error occur
						n3iwfSelf.AllocatedUEIPAddress[ueIPAddr.String()] = thisUE
						break
					}
				}
			}
			attributes = append(attributes, ike_message.BuildConfigurationAttribute(ike_message.INTERNAL_IP4_ADDRESS, ueIPAddr))
			attributes = append(attributes, ike_message.BuildConfigurationAttribute(ike_message.INTERNAL_IP4_NETMASK, n3iwfSelf.Subnet.Mask))

			thisUE.IPSecInnerIP = ueIPAddr.String()
			ikeLog.Tracef("ueIPAddr: %+v", ueIPAddr)

			// Prepare individual traffic selectors
			individualTrafficSelectorInitiator := ike_message.BuildIndividualTrafficSelector(ike_message.TS_IPV4_ADDR_RANGE, ike_message.IPProtocolAll,
				0, 65535, ueIPAddr.To4(), ueIPAddr.To4())
			individualTrafficSelectorResponder := ike_message.BuildIndividualTrafficSelector(ike_message.TS_IPV4_ADDR_RANGE, ike_message.IPProtocolAll,
				0, 65535, n3iwfIPAddr.To4(), n3iwfIPAddr.To4())

			responseTrafficSelectorInitiator = ike_message.BuildTrafficSelectorInitiator([]*ike_message.IndividualTrafficSelector{individualTrafficSelectorInitiator})
			responseTrafficSelectorResponder = ike_message.BuildTrafficSelectorResponder([]*ike_message.IndividualTrafficSelector{individualTrafficSelectorResponder})

			// Record traffic selector to IKE security association
			ikeSecurityAssociation.TrafficSelectorInitiator = responseTrafficSelectorInitiator
			ikeSecurityAssociation.TrafficSelectorResponder = responseTrafficSelectorResponder

			responseConfiguration = ike_message.BuildConfiguration(ike_message.CFG_REPLY, attributes)
		} else {
			ikeLog.Error("[IKE] UE did not send any configuration request for its IP address.")
			return
		}

		// Calculate local AUTH
		pseudorandomFunction, ok := NewPseudorandomFunction(thisUE.Kn3iwf, transformPseudorandomFunction.TransformID)
		if !ok {
			ikeLog.Error("[IKE] Get an unsupported pseudorandom funcion. This may imply an unsupported transform is chosen.")
			return
		}
		if _, err := pseudorandomFunction.Write([]byte("Key Pad for IKEv2")); err != nil {
			ikeLog.Errorf("[IKE] Pseudorandom function write error: %+v", err)
			return
		}
		secret := pseudorandomFunction.Sum(nil)
		pseudorandomFunction, ok = NewPseudorandomFunction(secret, transformPseudorandomFunction.TransformID)
		if !ok {
			ikeLog.Error("[IKE] Get an unsupported pseudorandom funcion. This may imply an unsupported transform is chosen.")
			return
		}
		if _, err := pseudorandomFunction.Write(ikeSecurityAssociation.LocalUnsignedAuthentication); err != nil {
			ikeLog.Errorf("[IKE] Pseudorandom function write error: %+v", err)
			return
		}

		// Get xfrm needed data
		// As specified in RFC 7296, ESP negotiate two child security association (pair) in one IKE_AUTH
		childSecurityAssociationContext, err := thisUE.CreateIKEChildSecurityAssociation(ikeSecurityAssociation.IKEAuthResponseSA)
		if err != nil {
			ikeLog.Errorf("[IKE] Create child security association context failed: %+v", err)
			return
		}
		err = parseIPAddressInformationToChildSecurityAssociation(childSecurityAssociationContext, ueSendInfo.Addr.IP, ikeSecurityAssociation.TrafficSelectorResponder.TrafficSelectors[0], ikeSecurityAssociation.TrafficSelectorInitiator.TrafficSelectors[0])
		if err != nil {
			ikeLog.Errorf("[IKE] Parse IP address to child security association failed: %+v", err)
			return
		}
		// Select TCP traffic
		childSecurityAssociationContext.SelectedIPProtocol = unix.IPPROTO_TCP

		if err := GenerateKeyForChildSA(ikeSecurityAssociation, childSecurityAssociationContext); err != nil {
			ikeLog.Errorf("[IKE] Generate key for child SA failed: %+v", err)
			return
		}

		// Build IKE message
		responseIKEMessage = ike_message.BuildIKEHeader(message.InitiatorSPI, message.ResponderSPI, ike_message.IKE_AUTH, ike_message.ResponseBitCheck, message.MessageID)

		// Build response
		var ikePayload []ike_message.IKEPayloadType

		// Configuration
		ikePayload = append(ikePayload, responseConfiguration)

		// Authentication
		responseAuthentication = ike_message.BuildAuthentication(ike_message.SharedKeyMesageIntegrityCode, pseudorandomFunction.Sum(nil))
		ikePayload = append(ikePayload, responseAuthentication)

		// Security Association
		responseSecurityAssociation = ikeSecurityAssociation.IKEAuthResponseSA
		ikePayload = append(ikePayload, responseSecurityAssociation)

		// Traffic Selector Initiator and Responder
		ikePayload = append(ikePayload, responseTrafficSelectorInitiator)
		ikePayload = append(ikePayload, responseTrafficSelectorResponder)

		// Notification(NAS_IP_ADDRESS)
		responseNotification = ike_message.BuildNotifyNAS_IP4_ADDRESS(n3iwfSelf.IPSecGatewayAddress)
		ikePayload = append(ikePayload, responseNotification)

		// Notification(NSA_TCP_PORT)
		responseNotification = ike_message.BuildNotifyNAS_TCP_PORT(n3iwfSelf.TCPPort)
		ikePayload = append(ikePayload, responseNotification)

		if err := EncryptProcedure(ikeSecurityAssociation, ikePayload, responseIKEMessage); err != nil {
			ikeLog.Errorf("Encrypting IKE message failed: %+v", err)
			return
		}

		// Send IKE message to UE
		SendIKEMessageToUE(ueSendInfo, responseIKEMessage)

		// Aplly XFRM rules
		if err = ApplyXFRMRule(false, childSecurityAssociationContext); err != nil {
			ikeLog.Errorf("[IKE] Applying XFRM rules failed: %+v", err)
			return
		}

		// If needed, setup PDU session
		if thisUE.TemporaryPDUSessionSetupData != nil {
			for {
				if len(thisUE.TemporaryPDUSessionSetupData.UnactivatedPDUSession) != 0 {
					pduSessionID := thisUE.TemporaryPDUSessionSetupData.UnactivatedPDUSession[0]
					pduSession := thisUE.PduSessionList[pduSessionID]

					ikeSecurityAssociation := thisUE.N3IWFIKESecurityAssociation

					// Send CREATE_CHILD_SA to UE
					// Add MessageID for IKE security association
					ikeSecurityAssociation.MessageID++
					ikeMessage := ike_message.BuildIKEHeader(ikeSecurityAssociation.LocalSPI, ikeSecurityAssociation.RemoteSPI, ike_message.CREATE_CHILD_SA, ike_message.InitiatorBitCheck, ikeSecurityAssociation.MessageID)

					// IKE payload
					var ikePayload []ike_message.IKEPayloadType

					// Build SA
					// Proposals
					var proposals []*ike_message.Proposal

					// Allocate SPI
					var spi uint32
					spiByte := make([]byte, 4)
					for {
						randomUint64 := GenerateRandomNumber().Uint64()
						if _, ok := n3iwfSelf.ChildSA[uint32(randomUint64)]; !ok {
							spi = uint32(randomUint64)
							break
						}
					}
					binary.BigEndian.PutUint32(spiByte, spi)

					// First Proposal - Proposal No.1
					proposal := ike_message.BuildProposal(1, ike_message.TypeESP, spiByte)

					// Encryption transform
					var attributeType uint16 = ike_message.AttributeTypeKeyLength
					var attributeValue uint16 = 256
					encryptionTransform := ike_message.BuildTransform(ike_message.TypeEncryptionAlgorithm, ike_message.ENCR_AES_CBC, &attributeType, &attributeValue, nil)
					if ok := ike_message.AppendTransformToProposal(proposal, encryptionTransform); !ok {
						ikeLog.Error("Generate IKE message failed: Cannot append to proposal")
						thisUE.TemporaryPDUSessionSetupData.UnactivatedPDUSession = thisUE.TemporaryPDUSessionSetupData.UnactivatedPDUSession[1:]
						cause := ngapType.Cause{
							Present: ngapType.CausePresentTransport,
							Transport: &ngapType.CauseTransport{
								Value: ngapType.CauseTransportPresentTransportResourceUnavailable,
							},
						}
						transfer, err := ngap_message.BuildPDUSessionResourceSetupUnsuccessfulTransfer(cause, nil)
						if err != nil {
							ikeLog.Errorf("Build PDU Session Resource Setup Unsuccessful Transfer Failed: %+v", err)
							continue
						}
						ngap_message.AppendPDUSessionResourceFailedToSetupListCxtRes(thisUE.TemporaryPDUSessionSetupData.FailedListCxtRes, pduSessionID, transfer)
						continue
					}
					// Integrity transform
					if pduSession.SecurityIntegrity {
						integrityTransform := ike_message.BuildTransform(ike_message.TypeIntegrityAlgorithm, ike_message.AUTH_HMAC_SHA1_96, nil, nil, nil)
						if ok := ike_message.AppendTransformToProposal(proposal, integrityTransform); !ok {
							ikeLog.Error("Generate IKE message failed: Cannot append to proposal")
							thisUE.TemporaryPDUSessionSetupData.UnactivatedPDUSession = thisUE.TemporaryPDUSessionSetupData.UnactivatedPDUSession[1:]
							cause := ngapType.Cause{
								Present: ngapType.CausePresentTransport,
								Transport: &ngapType.CauseTransport{
									Value: ngapType.CauseTransportPresentTransportResourceUnavailable,
								},
							}
							transfer, err := ngap_message.BuildPDUSessionResourceSetupUnsuccessfulTransfer(cause, nil)
							if err != nil {
								ikeLog.Errorf("Build PDU Session Resource Setup Unsuccessful Transfer Failed: %+v", err)
								continue
							}
							ngap_message.AppendPDUSessionResourceFailedToSetupListCxtRes(thisUE.TemporaryPDUSessionSetupData.FailedListCxtRes, pduSessionID, transfer)
							continue
						}
					}
					// ESN transform
					esnTransform := ike_message.BuildTransform(ike_message.TypeExtendedSequenceNumbers, ike_message.ESN_NO, nil, nil, nil)
					if ok := ike_message.AppendTransformToProposal(proposal, esnTransform); !ok {
						ikeLog.Error("Generate IKE message failed: Cannot append to proposal")
						thisUE.TemporaryPDUSessionSetupData.UnactivatedPDUSession = thisUE.TemporaryPDUSessionSetupData.UnactivatedPDUSession[1:]
						cause := ngapType.Cause{
							Present: ngapType.CausePresentTransport,
							Transport: &ngapType.CauseTransport{
								Value: ngapType.CauseTransportPresentTransportResourceUnavailable,
							},
						}
						transfer, err := ngap_message.BuildPDUSessionResourceSetupUnsuccessfulTransfer(cause, nil)
						if err != nil {
							ikeLog.Errorf("Build PDU Session Resource Setup Unsuccessful Transfer Failed: %+v", err)
							continue
						}
						ngap_message.AppendPDUSessionResourceFailedToSetupListCxtRes(thisUE.TemporaryPDUSessionSetupData.FailedListCxtRes, pduSessionID, transfer)
						continue
					}

					proposals = append(proposals, proposal)

					securityAssociation := ike_message.BuildSecurityAssociation(proposals)

					ikePayload = append(ikePayload, securityAssociation)

					// Build Nonce
					nonceData := GenerateRandomNumber().Bytes()
					nonce := ike_message.BuildNonce(nonceData)

					// Store nonce into context
					ikeSecurityAssociation.ConcatenatedNonce = nonceData

					ikePayload = append(ikePayload, nonce)

					// TSi
					ueIPAddr := net.ParseIP(thisUE.IPSecInnerIP)
					individualTrafficSelector := ike_message.BuildIndividualTrafficSelector(ike_message.TS_IPV4_ADDR_RANGE, ike_message.IPProtocolAll,
						0, 65535, ueIPAddr, ueIPAddr)
					trafficSelectorInitiator := ike_message.BuildTrafficSelectorInitiator([]*ike_message.IndividualTrafficSelector{individualTrafficSelector})

					ikePayload = append(ikePayload, trafficSelectorInitiator)

					// TSr
					n3iwfIPAddr := net.ParseIP(n3iwfSelf.IPSecGatewayAddress)
					individualTrafficSelector = ike_message.BuildIndividualTrafficSelector(ike_message.TS_IPV4_ADDR_RANGE, ike_message.IPProtocolAll,
						0, 65535, n3iwfIPAddr, n3iwfIPAddr)
					trafficSelectorResponder := ike_message.BuildTrafficSelectorResponder([]*ike_message.IndividualTrafficSelector{individualTrafficSelector})

					ikePayload = append(ikePayload, trafficSelectorResponder)

					// Notify-Qos
					notifyQos := ike_message.BuildNotify5G_QOS_INFO(uint8(pduSessionID), pduSession.QFIList, true)

					ikePayload = append(ikePayload, notifyQos)

					// Notify-UP_IP_ADDRESS
					notifyUPIPAddr := ike_message.BuildNotifyUP_IP4_ADDRESS(n3iwfSelf.IPSecGatewayAddress)

					ikePayload = append(ikePayload, notifyUPIPAddr)

					if err := EncryptProcedure(thisUE.N3IWFIKESecurityAssociation, ikePayload, ikeMessage); err != nil {
						ikeLog.Errorf("Encrypting IKE message failed: %+v", err)
						thisUE.TemporaryPDUSessionSetupData.UnactivatedPDUSession = thisUE.TemporaryPDUSessionSetupData.UnactivatedPDUSession[1:]
						cause := ngapType.Cause{
							Present: ngapType.CausePresentTransport,
							Transport: &ngapType.CauseTransport{
								Value: ngapType.CauseTransportPresentTransportResourceUnavailable,
							},
						}
						transfer, err := ngap_message.BuildPDUSessionResourceSetupUnsuccessfulTransfer(cause, nil)
						if err != nil {
							ikeLog.Errorf("Build PDU Session Resource Setup Unsuccessful Transfer Failed: %+v", err)
							continue
						}
						ngap_message.AppendPDUSessionResourceFailedToSetupListCxtRes(thisUE.TemporaryPDUSessionSetupData.FailedListCxtRes, pduSessionID, transfer)
						continue
					}

					SendIKEMessageToUE(ueSendInfo, ikeMessage)
					break
				} else {
					// Send Initial Context Setup Response to AMF
					ngap_message.SendInitialContextSetupResponse(thisUE.AMF, thisUE, thisUE.TemporaryPDUSessionSetupData.SetupListCxtRes, thisUE.TemporaryPDUSessionSetupData.FailedListCxtRes, nil)
					break
				}
			}
		} else {
			// Send Initial Context Setup Response to AMF
			ngap_message.SendInitialContextSetupResponse(thisUE.AMF, thisUE, nil, nil, nil)
		}
	}
}

func HandleCREATECHILDSA(ueSendInfo *n3iwf_message.UDPSendInfoGroup, message *ike_message.IKEMessage) {
	ikeLog.Infoln("[IKE] Handle CREATE_CHILD_SA")

	var encryptedPayload *ike_message.Encrypted

	n3iwfSelf := n3iwf_context.N3IWFSelf()

	// {response}
	var responseIKEMessage *ike_message.IKEMessage

	if message == nil {
		ikeLog.Error("[IKE] IKE Message is nil")
		return
	}

	// parse IKE header and setup IKE context
	// check major version
	majorVersion := ((message.Version & 0xf0) >> 4)
	if majorVersion > 2 {
		ikeLog.Warn("[IKE] Received an IKE message with higher major version")
		// send INFORMATIONAL type message with INVALID_MAJOR_VERSION Notify payload ( OUTSIDE IKE SA )

		// IKEHDR-{response}
		responseNotification := ike_message.BuildNotification(ike_message.TypeNone, ike_message.INVALID_MAJOR_VERSION, nil, nil)

		responseIKEMessage = ike_message.BuildIKEHeader(message.InitiatorSPI, message.ResponderSPI, ike_message.INFORMATIONAL, ike_message.ResponseBitCheck, message.MessageID)
		responseIKEMessage.IKEPayload = append(responseIKEMessage.IKEPayload, responseNotification)

		SendIKEMessageToUE(ueSendInfo, responseIKEMessage)

		return
	}

	// Find corresponding IKE security association
	localSPI := message.InitiatorSPI
	ikeSecurityAssociation := n3iwfSelf.FindIKESecurityAssociationBySPI(localSPI)
	if ikeSecurityAssociation == nil {
		ikeLog.Warn("[IKE] Unrecognized SPI")
		// send INFORMATIONAL type message with INVALID_IKE_SPI Notify payload ( OUTSIDE IKE SA )

		// IKEHDR-{response}
		responseNotification := ike_message.BuildNotification(ike_message.TypeNone, ike_message.INVALID_IKE_SPI, nil, nil)

		responseIKEMessage = ike_message.BuildIKEHeader(message.InitiatorSPI, 0, ike_message.INFORMATIONAL, ike_message.ResponseBitCheck, message.MessageID)
		responseIKEMessage.IKEPayload = append(responseIKEMessage.IKEPayload, responseNotification)

		SendIKEMessageToUE(ueSendInfo, responseIKEMessage)

		return
	}

	for _, ikePayload := range message.IKEPayload {
		switch ikePayload.Type() {
		case ike_message.TypeSK:
			encryptedPayload = ikePayload.(*ike_message.Encrypted)
		default:
			ikeLog.Warnf("[IKE] Get IKE payload (type %d) in IKE_SA_INIT message, this payload will not be handled by IKE handler", ikePayload.Type())
		}
	}

	decryptedIKEPayload, err := DecryptProcedure(ikeSecurityAssociation, message, encryptedPayload)
	if err != nil {
		ikeLog.Errorf("[IKE] Decrypt IKE message failed: %+v", err)
		return
	}

	// Parse payloads
	var securityAssociation *ike_message.SecurityAssociation
	var nonce *ike_message.Nonce
	var trafficSelectorInitiator *ike_message.TrafficSelectorInitiator
	var trafficSelectorResponder *ike_message.TrafficSelectorResponder

	for _, ikePayload := range decryptedIKEPayload {
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
			ikeLog.Warnf("[IKE] Get IKE payload (type %d) in IKE_AUTH message, this payload will not be handled by IKE handler", ikePayload.Type())
		}
	}

	// UE context
	thisUE := ikeSecurityAssociation.ThisUE
	if thisUE == nil {
		ikeLog.Error("UE context is nil")
		return
	}
	// PDU session information
	if thisUE.TemporaryPDUSessionSetupData == nil {
		ikeLog.Error("No PDU session information")
		return
	}
	temporaryPDUSessionSetupData := thisUE.TemporaryPDUSessionSetupData
	if len(temporaryPDUSessionSetupData.UnactivatedPDUSession) == 0 {
		ikeLog.Error("No unactivated PDU session information")
		return
	}
	pduSessionID := temporaryPDUSessionSetupData.UnactivatedPDUSession[0]
	pduSession, ok := thisUE.PduSessionList[pduSessionID]
	if !ok {
		ikeLog.Errorf("No such PDU session [PDU session ID: %d]", pduSessionID)
		return
	}

	// Check received message
	if securityAssociation == nil {
		ikeLog.Error("[IKE] The security association field is nil")
		return
	}

	if trafficSelectorInitiator == nil {
		ikeLog.Error("[IKE] The traffic selector initiator field is nil")
		return
	}

	if trafficSelectorResponder == nil {
		ikeLog.Error("[IKE] The traffic selector responder field is nil")
		return
	}

	// Nonce
	if nonce != nil {
		ikeSecurityAssociation.ConcatenatedNonce = append(ikeSecurityAssociation.ConcatenatedNonce, nonce.NonceData...)
	} else {
		ikeLog.Error("[IKE] The nonce field is nil")
		// TODO: send error message to UE
		return
	}

	// Get xfrm needed data
	// As specified in RFC 7296, ESP negotiate two child security association (pair) in one IKE_AUTH
	childSecurityAssociationContext, err := thisUE.CreateIKEChildSecurityAssociation(securityAssociation)
	if err != nil {
		ikeLog.Errorf("[IKE] Create child security association context failed: %+v", err)
		return
	}
	err = parseIPAddressInformationToChildSecurityAssociation(childSecurityAssociationContext, ueSendInfo.Addr.IP, trafficSelectorInitiator.TrafficSelectors[0], trafficSelectorResponder.TrafficSelectors[0])
	if err != nil {
		ikeLog.Errorf("[IKE] Parse IP address to child security association failed: %+v", err)
		return
	}
	// Select GRE traffic
	childSecurityAssociationContext.SelectedIPProtocol = unix.IPPROTO_GRE

	if err := GenerateKeyForChildSA(ikeSecurityAssociation, childSecurityAssociationContext); err != nil {
		ikeLog.Errorf("[IKE] Generate key for child SA failed: %+v", err)
		return
	}

	// Aplly XFRM rules
	if err = ApplyXFRMRule(true, childSecurityAssociationContext); err != nil {
		ikeLog.Errorf("[IKE] Applying XFRM rules failed: %+v", err)
		return
	}

	// Setup GTP tunnel for UE
	ueAssociatedGTPConnection := pduSession.GTPConnection
	if userPlaneConnection, ok := n3iwfSelf.GTPConnectionWithUPF[ueAssociatedGTPConnection.UPFIPAddr]; ok {
		// UPF UDP address
		upfUDPAddr, err := net.ResolveUDPAddr("udp", ueAssociatedGTPConnection.UPFIPAddr+":2152")
		if err != nil {
			ikeLog.Errorf("Resolve UDP address failed: %+v", err)
			return
		}

		// UE TEID
		ueTEID := n3iwfSelf.NewTEID(thisUE)

		// Set UE associated GTP connection
		ueAssociatedGTPConnection.UPFUDPAddr = upfUDPAddr
		ueAssociatedGTPConnection.IncomingTEID = ueTEID
		ueAssociatedGTPConnection.UserPlaneConnection = userPlaneConnection

		// Append NGAP PDU session resource setup response transfer
		transfer, err := ngap_message.BuildPDUSessionResourceSetupResponseTransfer(pduSession)
		if err != nil {
			ikeLog.Errorf("Build PDU session resource setup response transfer failed: %+v", err)
			return
		}
		ngap_message.AppendPDUSessionResourceSetupListSURes(temporaryPDUSessionSetupData.SetupListSURes, pduSessionID, transfer)
	} else {
		// Setup GTP connection with UPF
		userPlaneConnection, upfUDPAddr, err := n3iwf_data_relay.SetupGTPTunnelWithUPF(ueAssociatedGTPConnection.UPFIPAddr)
		if err != nil {
			ikeLog.Errorf("Setup GTP connection with UPF failed: %+v", err)
			return
		}
		// Listen GTP tunnel
		if err := n3iwf_data_relay.ListenGTP(userPlaneConnection); err != nil {
			ikeLog.Errorf("Listening GTP tunnel failed: %+v", err)
			return
		}

		// UE TEID
		ueTEID := n3iwfSelf.NewTEID(thisUE)

		// Setup GTP connection with UPF
		ueAssociatedGTPConnection.UPFUDPAddr = upfUDPAddr
		ueAssociatedGTPConnection.IncomingTEID = ueTEID
		ueAssociatedGTPConnection.UserPlaneConnection = userPlaneConnection

		// Store GTP connection with UPF into N3IWF context
		n3iwfSelf.GTPConnectionWithUPF[ueAssociatedGTPConnection.UPFIPAddr] = userPlaneConnection

		// Append NGAP PDU session resource setup response transfer
		transfer, err := ngap_message.BuildPDUSessionResourceSetupResponseTransfer(pduSession)
		if err != nil {
			ikeLog.Errorf("Build PDU session resource setup response transfer failed: %+v", err)
			return
		}
		ngap_message.AppendPDUSessionResourceSetupListSURes(temporaryPDUSessionSetupData.SetupListSURes, pduSessionID, transfer)
	}

	temporaryPDUSessionSetupData.UnactivatedPDUSession = temporaryPDUSessionSetupData.UnactivatedPDUSession[1:]

	for {
		if len(temporaryPDUSessionSetupData.UnactivatedPDUSession) != 0 {
			ngapProcedure := temporaryPDUSessionSetupData.NGAPProcedureCode.Value
			pduSessionID := temporaryPDUSessionSetupData.UnactivatedPDUSession[0]
			pduSession := thisUE.PduSessionList[pduSessionID]

			ikeSecurityAssociation := thisUE.N3IWFIKESecurityAssociation

			// Send CREATE_CHILD_SA to UE
			// Add MessageID for IKE security association
			ikeSecurityAssociation.MessageID++
			ikeMessage := ike_message.BuildIKEHeader(ikeSecurityAssociation.LocalSPI, ikeSecurityAssociation.RemoteSPI, ike_message.CREATE_CHILD_SA, ike_message.InitiatorBitCheck, ikeSecurityAssociation.MessageID)

			// IKE payload
			var ikePayload []ike_message.IKEPayloadType

			// Build SA
			// Proposals
			var proposals []*ike_message.Proposal

			// Allocate SPI
			var spi uint32
			spiByte := make([]byte, 4)
			for {
				randomUint64 := GenerateRandomNumber().Uint64()
				if _, ok := n3iwfSelf.ChildSA[uint32(randomUint64)]; !ok {
					spi = uint32(randomUint64)
					break
				}
			}
			binary.BigEndian.PutUint32(spiByte, spi)

			// First Proposal - Proposal No.1
			proposal := ike_message.BuildProposal(1, ike_message.TypeESP, spiByte)

			// Encryption transform
			var attributeType uint16 = ike_message.AttributeTypeKeyLength
			var attributeValue uint16 = 256
			encryptionTransform := ike_message.BuildTransform(ike_message.TypeEncryptionAlgorithm, ike_message.ENCR_AES_CBC, &attributeType, &attributeValue, nil)
			if ok := ike_message.AppendTransformToProposal(proposal, encryptionTransform); !ok {
				ikeLog.Error("Generate IKE message failed: Cannot append to proposal")
				temporaryPDUSessionSetupData.UnactivatedPDUSession = temporaryPDUSessionSetupData.UnactivatedPDUSession[1:]
				cause := ngapType.Cause{
					Present: ngapType.CausePresentTransport,
					Transport: &ngapType.CauseTransport{
						Value: ngapType.CauseTransportPresentTransportResourceUnavailable,
					},
				}
				transfer, err := ngap_message.BuildPDUSessionResourceSetupUnsuccessfulTransfer(cause, nil)
				if err != nil {
					ikeLog.Errorf("Build PDU Session Resource Setup Unsuccessful Transfer Failed: %+v", err)
					continue
				}
				if ngapProcedure == ngapType.ProcedureCodeInitialContextSetup {
					ngap_message.AppendPDUSessionResourceFailedToSetupListCxtRes(temporaryPDUSessionSetupData.FailedListCxtRes, pduSessionID, transfer)
				} else {
					ngap_message.AppendPDUSessionResourceFailedToSetupListSURes(temporaryPDUSessionSetupData.FailedListSURes, pduSessionID, transfer)
				}
				continue
			}
			// Integrity transform
			if pduSession.SecurityIntegrity {
				integrityTransform := ike_message.BuildTransform(ike_message.TypeIntegrityAlgorithm, ike_message.AUTH_HMAC_MD5_96, nil, nil, nil)
				if ok := ike_message.AppendTransformToProposal(proposal, integrityTransform); !ok {
					ikeLog.Error("Generate IKE message failed: Cannot append to proposal")
					temporaryPDUSessionSetupData.UnactivatedPDUSession = temporaryPDUSessionSetupData.UnactivatedPDUSession[1:]
					cause := ngapType.Cause{
						Present: ngapType.CausePresentTransport,
						Transport: &ngapType.CauseTransport{
							Value: ngapType.CauseTransportPresentTransportResourceUnavailable,
						},
					}
					transfer, err := ngap_message.BuildPDUSessionResourceSetupUnsuccessfulTransfer(cause, nil)
					if err != nil {
						ikeLog.Errorf("Build PDU Session Resource Setup Unsuccessful Transfer Failed: %+v", err)
						continue
					}
					if ngapProcedure == ngapType.ProcedureCodeInitialContextSetup {
						ngap_message.AppendPDUSessionResourceFailedToSetupListCxtRes(temporaryPDUSessionSetupData.FailedListCxtRes, pduSessionID, transfer)
					} else {
						ngap_message.AppendPDUSessionResourceFailedToSetupListSURes(temporaryPDUSessionSetupData.FailedListSURes, pduSessionID, transfer)
					}
					continue
				}
			}
			// ESN transform
			esnTransform := ike_message.BuildTransform(ike_message.TypeExtendedSequenceNumbers, ike_message.ESN_NO, nil, nil, nil)
			if ok := ike_message.AppendTransformToProposal(proposal, esnTransform); !ok {
				ikeLog.Error("Generate IKE message failed: Cannot append to proposal")
				temporaryPDUSessionSetupData.UnactivatedPDUSession = temporaryPDUSessionSetupData.UnactivatedPDUSession[1:]
				cause := ngapType.Cause{
					Present: ngapType.CausePresentTransport,
					Transport: &ngapType.CauseTransport{
						Value: ngapType.CauseTransportPresentTransportResourceUnavailable,
					},
				}
				transfer, err := ngap_message.BuildPDUSessionResourceSetupUnsuccessfulTransfer(cause, nil)
				if err != nil {
					ikeLog.Errorf("Build PDU Session Resource Setup Unsuccessful Transfer Failed: %+v", err)
					continue
				}
				if ngapProcedure == ngapType.ProcedureCodeInitialContextSetup {
					ngap_message.AppendPDUSessionResourceFailedToSetupListCxtRes(temporaryPDUSessionSetupData.FailedListCxtRes, pduSessionID, transfer)
				} else {
					ngap_message.AppendPDUSessionResourceFailedToSetupListSURes(temporaryPDUSessionSetupData.FailedListSURes, pduSessionID, transfer)
				}
				continue
			}

			proposals = append(proposals, proposal)

			securityAssociation := ike_message.BuildSecurityAssociation(proposals)

			ikePayload = append(ikePayload, securityAssociation)

			// Build Nonce
			nonceData := GenerateRandomNumber().Bytes()
			nonce := ike_message.BuildNonce(nonceData)

			// Store nonce into context
			ikeSecurityAssociation.ConcatenatedNonce = nonceData

			ikePayload = append(ikePayload, nonce)

			// TSi
			ueIPAddr := net.ParseIP(thisUE.IPSecInnerIP)
			individualTrafficSelector := ike_message.BuildIndividualTrafficSelector(ike_message.TS_IPV4_ADDR_RANGE, ike_message.IPProtocolAll,
				0, 65535, ueIPAddr, ueIPAddr)
			trafficSelectorInitiator := ike_message.BuildTrafficSelectorInitiator([]*ike_message.IndividualTrafficSelector{individualTrafficSelector})

			ikePayload = append(ikePayload, trafficSelectorInitiator)

			// TSr
			n3iwfIPAddr := net.ParseIP(n3iwfSelf.IPSecGatewayAddress)
			individualTrafficSelector = ike_message.BuildIndividualTrafficSelector(ike_message.TS_IPV4_ADDR_RANGE, ike_message.IPProtocolAll,
				0, 65535, n3iwfIPAddr, n3iwfIPAddr)
			trafficSelectorResponder := ike_message.BuildTrafficSelectorResponder([]*ike_message.IndividualTrafficSelector{individualTrafficSelector})

			ikePayload = append(ikePayload, trafficSelectorResponder)

			// Notify-Qos
			notifyQos := ike_message.BuildNotify5G_QOS_INFO(uint8(pduSessionID), pduSession.QFIList, true)

			ikePayload = append(ikePayload, notifyQos)

			// Notify-UP_IP_ADDRESS
			notifyUPIPAddr := ike_message.BuildNotifyUP_IP4_ADDRESS(n3iwfSelf.IPSecGatewayAddress)

			ikePayload = append(ikePayload, notifyUPIPAddr)

			if err := EncryptProcedure(thisUE.N3IWFIKESecurityAssociation, ikePayload, ikeMessage); err != nil {
				ikeLog.Errorf("Encrypting IKE message failed: %+v", err)
				temporaryPDUSessionSetupData.UnactivatedPDUSession = temporaryPDUSessionSetupData.UnactivatedPDUSession[1:]
				cause := ngapType.Cause{
					Present: ngapType.CausePresentTransport,
					Transport: &ngapType.CauseTransport{
						Value: ngapType.CauseTransportPresentTransportResourceUnavailable,
					},
				}
				transfer, err := ngap_message.BuildPDUSessionResourceSetupUnsuccessfulTransfer(cause, nil)
				if err != nil {
					ikeLog.Errorf("Build PDU Session Resource Setup Unsuccessful Transfer Failed: %+v", err)
					continue
				}
				if ngapProcedure == ngapType.ProcedureCodeInitialContextSetup {
					ngap_message.AppendPDUSessionResourceFailedToSetupListCxtRes(temporaryPDUSessionSetupData.FailedListCxtRes, pduSessionID, transfer)
				} else {
					ngap_message.AppendPDUSessionResourceFailedToSetupListSURes(temporaryPDUSessionSetupData.FailedListSURes, pduSessionID, transfer)
				}
				continue
			}

			SendIKEMessageToUE(ueSendInfo, ikeMessage)
			break
		} else {
			// Send Response to AMF
			ngapProcedure := temporaryPDUSessionSetupData.NGAPProcedureCode.Value
			if ngapProcedure == ngapType.ProcedureCodeInitialContextSetup {
				ngap_message.SendInitialContextSetupResponse(thisUE.AMF, thisUE, temporaryPDUSessionSetupData.SetupListCxtRes, temporaryPDUSessionSetupData.FailedListCxtRes, nil)
			} else {
				ngap_message.SendPDUSessionResourceSetupResponse(thisUE.AMF, thisUE, temporaryPDUSessionSetupData.SetupListSURes, temporaryPDUSessionSetupData.FailedListSURes, nil)
			}
			break
		}
	}
}

func is_supported(transformType uint8, transformID uint16, attributePresent bool, attributeValue uint16) bool {
	switch transformType {
	case ike_message.TypeEncryptionAlgorithm:
		switch transformID {
		case ike_message.ENCR_DES_IV64:
			return false
		case ike_message.ENCR_DES:
			return false
		case ike_message.ENCR_3DES:
			return false
		case ike_message.ENCR_RC5:
			return false
		case ike_message.ENCR_IDEA:
			return false
		case ike_message.ENCR_CAST:
			return false
		case ike_message.ENCR_BLOWFISH:
			return false
		case ike_message.ENCR_3IDEA:
			return false
		case ike_message.ENCR_DES_IV32:
			return false
		case ike_message.ENCR_NULL:
			return false
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
			return false
		default:
			return false
		}
	case ike_message.TypePseudorandomFunction:
		switch transformID {
		case ike_message.PRF_HMAC_MD5:
			return true
		case ike_message.PRF_HMAC_SHA1:
			return true
		case ike_message.PRF_HMAC_TIGER:
			return false
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
			return false
		default:
			return false
		}
	case ike_message.TypeDiffieHellmanGroup:
		switch transformID {
		case ike_message.DH_NONE:
			return false
		case ike_message.DH_768_BIT_MODP:
			return false
		case ike_message.DH_1024_BIT_MODP:
			return true
		case ike_message.DH_1536_BIT_MODP:
			return false
		case ike_message.DH_2048_BIT_MODP:
			return true
		case ike_message.DH_3072_BIT_MODP:
			return false
		case ike_message.DH_4096_BIT_MODP:
			return false
		case ike_message.DH_6144_BIT_MODP:
			return false
		case ike_message.DH_8192_BIT_MODP:
			return false
		default:
			return false
		}
	default:
		return false
	}
}

func is_Kernel_Supported(transformType uint8, transformID uint16, attributePresent bool, attributeValue uint16) bool {
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
		default:
			return false
		}
	case ike_message.TypeDiffieHellmanGroup:
		switch transformID {
		case ike_message.DH_NONE:
			return false
		case ike_message.DH_768_BIT_MODP:
			return false
		case ike_message.DH_1024_BIT_MODP:
			return false
		case ike_message.DH_1536_BIT_MODP:
			return false
		case ike_message.DH_2048_BIT_MODP:
			return false
		case ike_message.DH_3072_BIT_MODP:
			return false
		case ike_message.DH_4096_BIT_MODP:
			return false
		case ike_message.DH_6144_BIT_MODP:
			return false
		case ike_message.DH_8192_BIT_MODP:
			return false
		default:
			return false
		}
	case ike_message.TypeExtendedSequenceNumbers:
		switch transformID {
		case ike_message.ESN_NO:
			return true
		case ike_message.ESN_NEED:
			return true
		default:
			return false
		}
	default:
		return false
	}
}

func GenerateRandomIPinRange(subnet *net.IPNet) net.IP {
	ipAddr := make([]byte, 4)

	// TODO: elimenate network name, gateway, and broadcast
	for i := 0; i < 4; i++ {
		randomNumber, err := GenerateRandomUint8()
		if err != nil {
			ikeLog.Errorf("[IKE] Generate random number for IP address failed: %+v", err)
			return nil
		}
		alter := byte(randomNumber) & (subnet.Mask[i] ^ 255)
		ipAddr[i] = subnet.IP[i] + alter
	}

	return net.IPv4(ipAddr[0], ipAddr[1], ipAddr[2], ipAddr[3])
}

func parseIPAddressInformationToChildSecurityAssociation(
	childSecurityAssociation *n3iwf_context.ChildSecurityAssociation,
	uePublicIPAddr net.IP,
	trafficSelectorLocal *ike_message.IndividualTrafficSelector,
	trafficSelectorRemote *ike_message.IndividualTrafficSelector) error {

	if childSecurityAssociation == nil {
		return errors.New("childSecurityAssociation is nil")
	}

	childSecurityAssociation.PeerPublicIPAddr = uePublicIPAddr
	childSecurityAssociation.LocalPublicIPAddr = net.ParseIP(n3iwf_context.N3IWFSelf().IKEBindAddress)

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
