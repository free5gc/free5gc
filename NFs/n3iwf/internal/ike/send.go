package ike

import (
	"math"
	"net"

	"github.com/pkg/errors"

	"github.com/free5gc/ike"
	ike_message "github.com/free5gc/ike/message"
	"github.com/free5gc/ike/security"
	n3iwf_context "github.com/free5gc/n3iwf/internal/context"
	"github.com/free5gc/n3iwf/internal/logger"
)

func SendIKEMessageToUE(
	udpConn *net.UDPConn,
	srcAddr, dstAddr *net.UDPAddr,
	message *ike_message.IKEMessage,
	ikeSAKey *security.IKESAKey,
) error {
	ikeLog := logger.IKELog
	ikeLog.Trace("Send IKE message to UE")
	ikeLog.Trace("Encoding...")
	pkt, err := ike.EncodeEncrypt(message, ikeSAKey, ike_message.Role_Responder)
	if err != nil {
		return errors.Wrapf(err, "SendIKEMessageToUE")
	}
	// As specified in RFC 7296 section 3.1, the IKE message send from/to UDP port 4500
	// should prepend a 4 bytes zero
	if srcAddr.Port == 4500 {
		prependZero := make([]byte, 4)
		pkt = append(prependZero, pkt...)
	}

	ikeLog.Trace("Sending...")
	n, err := udpConn.WriteToUDP(pkt, dstAddr)
	if err != nil {
		return errors.Wrapf(err, "SendIKEMessageToUE")
	}
	if n != len(pkt) {
		return errors.Errorf("SendIKEMessageToUE Not all of the data is sent. Total length: %d. Sent: %d.",
			len(pkt), n)
	}
	return nil
}

func SendUEInformationExchange(
	ikeSA *n3iwf_context.IKESecurityAssociation,
	ikeSAKey *security.IKESAKey,
	payload *ike_message.IKEPayloadContainer, initiator bool,
	response bool, messageID uint32, conn *net.UDPConn,
	ueAddr *net.UDPAddr, n3iwfAddr *net.UDPAddr,
) {
	ikeLog := logger.IKELog

	// Build IKE message
	responseIKEMessage := ike_message.NewMessage(ikeSA.RemoteSPI, ikeSA.LocalSPI,
		ike_message.INFORMATIONAL, response, initiator, messageID, nil)

	if payload != nil && len(*payload) > 0 {
		responseIKEMessage.Payloads = append(responseIKEMessage.Payloads, *payload...)
	}

	err := SendIKEMessageToUE(conn, n3iwfAddr, ueAddr, responseIKEMessage, ikeSAKey)
	if err != nil {
		ikeLog.Errorf("SendUEInformationExchange err: %+v", err)
		return
	}
}

func SendIKEDeleteRequest(n3iwfCtx *n3iwf_context.N3IWFContext, localSPI uint64) {
	ikeLog := logger.IKELog
	ikeUe, ok := n3iwfCtx.IkeUePoolLoad(localSPI)
	if !ok {
		ikeLog.Errorf("Cannot get IkeUE from SPI : %+v", localSPI)
		return
	}

	var deletePayload ike_message.IKEPayloadContainer
	deletePayload.BuildDeletePayload(ike_message.TypeIKE, 0, 0, nil)
	SendUEInformationExchange(ikeUe.N3IWFIKESecurityAssociation, ikeUe.N3IWFIKESecurityAssociation.IKESAKey,
		&deletePayload, false, false, ikeUe.N3IWFIKESecurityAssociation.ResponderMessageID,
		ikeUe.IKEConnection.Conn, ikeUe.IKEConnection.UEAddr, ikeUe.IKEConnection.N3IWFAddr)
}

func SendChildSADeleteRequest(
	ikeUe *n3iwf_context.N3IWFIkeUe,
	relaseList []int64,
) {
	ikeLog := logger.IKELog
	var deleteSPIs []uint32
	spiLen := uint16(0)
	for _, releaseItem := range relaseList {
		for _, childSA := range ikeUe.N3IWFChildSecurityAssociation {
			if childSA.PDUSessionIds[0] == releaseItem {
				spi := childSA.XfrmStateList[0].Spi
				if spi < 0 || spi > math.MaxUint32 {
					ikeLog.Errorf("SendChildSADeleteRequest spi out of uint32 range : %d", spi)
					return
				}
				deleteSPIs = append(deleteSPIs, uint32(spi))
				spiLen += 1
				err := ikeUe.DeleteChildSA(childSA)
				if err != nil {
					ikeLog.Errorf("Delete Child SA error : %v", err)
					return
				}
			}
		}
	}

	var deletePayload ike_message.IKEPayloadContainer
	deletePayload.BuildDeletePayload(ike_message.TypeESP, 4, spiLen, deleteSPIs)
	SendUEInformationExchange(ikeUe.N3IWFIKESecurityAssociation, ikeUe.N3IWFIKESecurityAssociation.IKESAKey,
		&deletePayload, false, false, ikeUe.N3IWFIKESecurityAssociation.ResponderMessageID,
		ikeUe.IKEConnection.Conn, ikeUe.IKEConnection.UEAddr, ikeUe.IKEConnection.N3IWFAddr)
}
