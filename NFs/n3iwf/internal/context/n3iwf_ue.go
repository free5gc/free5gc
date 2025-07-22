package context

import (
	"net"

	"github.com/pkg/errors"

	"github.com/free5gc/ngap/ngapConvert"
	"github.com/free5gc/ngap/ngapType"
)

type N3IWFRanUe struct {
	RanUeSharedCtx

	// Temporary cached NAS message
	// Used when NAS registration accept arrived before
	// UE setup NAS TCP connection with N3IWF, and
	// Forward pduSessionEstablishmentAccept to UE after
	// UE send CREATE_CHILD_SA response
	TemporaryCachedNASMessage []byte

	// NAS TCP Connection Established
	IsNASTCPConnEstablished         bool
	IsNASTCPConnEstablishedComplete bool

	// NAS TCP Connection
	TCPConnection net.Conn
}

func (n3iwfUe *N3IWFRanUe) init(ranUeNgapId int64) {
	n3iwfUe.RanUeNgapId = ranUeNgapId
	n3iwfUe.AmfUeNgapId = AmfUeNgapIdUnspecified
	n3iwfUe.PduSessionList = make(map[int64]*PDUSession)
	n3iwfUe.TemporaryPDUSessionSetupData = new(PDUSessionSetupTemporaryData)
	n3iwfUe.IsNASTCPConnEstablished = false
	n3iwfUe.IsNASTCPConnEstablishedComplete = false
}

func (ranUe *N3IWFRanUe) Remove() error {
	// remove from AMF context
	ranUe.DetachAMF()

	// remove from RAN UE context
	n3iwfCtx := ranUe.N3iwfCtx
	n3iwfCtx.DeleteRanUe(ranUe.RanUeNgapId)

	for _, pduSession := range ranUe.PduSessionList {
		n3iwfCtx.DeleteTEID(pduSession.GTPConnInfo.IncomingTEID)
	}

	if ranUe.TCPConnection != nil {
		if err := ranUe.TCPConnection.Close(); err != nil {
			return errors.Errorf("Close TCP conn error : %v", err)
		}
	}

	return nil
}

func (n3iwfUe *N3IWFRanUe) AttachAMF(sctpAddr string) bool {
	if amf, ok := n3iwfUe.N3iwfCtx.AMFPoolLoad(sctpAddr); ok {
		amf.N3iwfRanUeList[n3iwfUe.RanUeNgapId] = n3iwfUe
		n3iwfUe.AMF = amf
		return true
	} else {
		return false
	}
}

func (n3iwfUe *N3IWFRanUe) DetachAMF() {
	if n3iwfUe.AMF == nil {
		return
	}
	delete(n3iwfUe.AMF.N3iwfRanUeList, n3iwfUe.RanUeNgapId)
}

// Implement RanUe interface
func (n3iwfUe *N3IWFRanUe) GetUserLocationInformation() *ngapType.UserLocationInformation {
	userLocationInformation := new(ngapType.UserLocationInformation)

	userLocationInformation.Present = ngapType.UserLocationInformationPresentUserLocationInformationN3IWF
	userLocationInformation.UserLocationInformationN3IWF = new(ngapType.UserLocationInformationN3IWF)

	userLocationInfoN3IWF := userLocationInformation.UserLocationInformationN3IWF
	userLocationInfoN3IWF.IPAddress = ngapConvert.IPAddressToNgap(n3iwfUe.IPAddrv4, n3iwfUe.IPAddrv6)
	userLocationInfoN3IWF.PortNumber = ngapConvert.PortNumberToNgap(n3iwfUe.PortNumber)

	return userLocationInformation
}
