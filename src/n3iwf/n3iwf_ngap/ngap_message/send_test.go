package ngap_message_test

import (
	"fmt"
	"git.cs.nctu.edu.tw/calee/sctp"
	"github.com/stretchr/testify/assert"
	"free5gc/lib/aper"
	"free5gc/lib/nas/nasMessage"
	"free5gc/lib/nas/nasTestpacket"
	"free5gc/lib/nas/nasType"
	"free5gc/lib/ngap/ngapConvert"
	"free5gc/lib/ngap/ngapSctp"
	"free5gc/lib/ngap/ngapType"
	"free5gc/lib/path_util"
	"free5gc/src/n3iwf/factory"
	"free5gc/src/n3iwf/n3iwf_context"
	"free5gc/src/n3iwf/n3iwf_ngap/n3iwf_sctp"
	"free5gc/src/n3iwf/n3iwf_ngap/ngap_message"
	"free5gc/src/n3iwf/n3iwf_util"
	"free5gc/src/test/ngapTestpacket"
	"log"
	"sync"
	"testing"
	"time"
)

var listenConn *sctp.SCTPListener
var testConfig *factory.Configuration
var testAddrAndPort string
var testAmf *n3iwf_context.N3IWFAMF
var n3iwfSelf *n3iwf_context.N3IWFContext

func init() {

	// parse config file
	configFile := path_util.Gofree5gcPath("free5gc/src/n3iwf/n3iwf_util/test/testN3iwfcfg.conf")
	factory.InitConfigFactory(configFile)
	testConfig = factory.N3iwfConfig.Configuration

	testAddrAndPort = fmt.Sprintf("%s:38412", testConfig.AMFAddress[0].NetworkAddress)
	listenConn = ngapSctp.Server(testConfig.AMFAddress[0].NetworkAddress)
	go func() {
		for {
			readChan := make(chan ngapSctp.ConnData, 1024)
			conn, err := ngapSctp.Accept(listenConn)
			if err != nil {
				log.Panicf("Accept error [%s]", err.Error())
			}
			fmt.Printf("SCTP Accept from: %s", conn.RemoteAddr().String())
			go ngapSctp.Start(conn, readChan)
		}
	}()
	time.Sleep(100 * time.Millisecond)

	// init n3iwf context
	n3iwf_util.InitN3IWFContext()

	// set up sctp conntection by n3iwf
	wg := sync.WaitGroup{}
	n3iwf_sctp.SetupSCTPConnection(testConfig.AMFAddress, &wg)

	n3iwfSelf = n3iwf_context.N3IWFSelf()
	testAmf = n3iwfSelf.NewN3iwfAmf(testAddrAndPort)

}

func TestSendErrorIndication(t *testing.T) {

	cause := buildCause(ngapType.CausePresentProtocol, ngapType.CauseProtocolPresentAbstractSyntaxErrorReject)

	procedureCode := ngapType.ProcedureCodeNGSetup
	triggeringMessage := ngapType.TriggeringMessagePresentSuccessfulOutcome
	procedureCriticality := ngapType.CriticalityPresentReject
	var iesCriticalityDiagnostics ngapType.CriticalityDiagnosticsIEList
	item := buildCriticalityDiagnosticsIEItem(ngapType.CriticalityPresentReject, 13, ngapType.TypeOfErrorPresentMissing)
	iesCriticalityDiagnostics.List = append(iesCriticalityDiagnostics.List, item)
	criticalityDiagnostics := buildCriticalityDiagnostics(&procedureCode, &triggeringMessage, &procedureCriticality, &iesCriticalityDiagnostics)

	ngap_message.SendErrorIndication(testAmf, nil, nil, cause, &criticalityDiagnostics)

	time.Sleep(100 * time.Millisecond)

}

func TestSendPDUSessionResourceModifyIndication(t *testing.T) {

	ue := initUe(t, "10.0.0.1", "", 8080, 1)
	modifyList := []ngapType.PDUSessionResourceModifyItemModInd{
		{
			PDUSessionID: ngapType.PDUSessionID{
				Value: 10,
			},
			PDUSessionResourceModifyIndicationTransfer: ngapTestpacket.GetPDUSessionResourceModifyIndicationTransfer(),
		},
	}

	ngap_message.SendPDUSessionResourceModifyIndication(testAmf, ue, modifyList)

	time.Sleep(100 * time.Millisecond)
	ue.Remove()

}

func TestSendPDUSessionResourceReleaseResponse(t *testing.T) {

	ue := initUe(t, "10.0.0.1", "", 8080, 1)
	relList := ngapType.PDUSessionResourceReleasedListRelRes{
		List: []ngapType.PDUSessionResourceReleasedItemRelRes{
			{
				PDUSessionID: ngapType.PDUSessionID{
					Value: 10,
				},
				PDUSessionResourceReleaseResponseTransfer: ngapTestpacket.GetPDUSessionResourceReleaseResponseTransfer(),
			},
		},
	}

	ngap_message.SendPDUSessionResourceReleaseResponse(testAmf, ue, relList, nil)

	time.Sleep(100 * time.Millisecond)
	ue.Remove()

}

func TestSendPDUSessionResourceNotify(t *testing.T) {

	ue := initUe(t, "10.0.0.1", "", 8080, 1)
	notiList := ngapType.PDUSessionResourceNotifyList{
		List: []ngapType.PDUSessionResourceNotifyItem{
			{
				PDUSessionID: ngapType.PDUSessionID{
					Value: 10,
				},
				PDUSessionResourceNotifyTransfer: ngapTestpacket.GetPDUSessionResourceNotifyTransfer([]int64{1}, []uint64{1}, []int64{2}),
			},
		},
	}
	relList := ngapType.PDUSessionResourceReleasedListNot{
		List: []ngapType.PDUSessionResourceReleasedItemNot{
			{
				PDUSessionID: ngapType.PDUSessionID{
					Value: 11,
				},
				PDUSessionResourceNotifyReleasedTransfer: ngapTestpacket.GetPDUSessionResourceNotifyReleasedTransfer(),
			},
		},
	}

	ngap_message.SendPDUSessionResourceNotify(testAmf, ue, &notiList, &relList)

	time.Sleep(100 * time.Millisecond)
	ue.Remove()

}

func TestSendNGReset(t *testing.T) {

	ue := initUe(t, "10.0.0.1", "", 8080, 1)

	cause := buildCause(ngapType.CausePresentProtocol, ngapType.CauseProtocolPresentAbstractSyntaxErrorReject)

	// ResetAll
	ngap_message.SendNGReset(testAmf, *cause, nil)
	time.Sleep(100 * time.Millisecond)

	// Reset Partof
	partOfNGInterface := &ngapType.UEAssociatedLogicalNGConnectionList{
		List: []ngapType.UEAssociatedLogicalNGConnectionItem{
			{
				AMFUENGAPID: &ngapType.AMFUENGAPID{
					Value: 1,
				},
				RANUENGAPID: &ngapType.RANUENGAPID{
					Value: 1,
				},
			},
		},
	}
	ngap_message.SendNGReset(testAmf, *cause, partOfNGInterface)
	time.Sleep(100 * time.Millisecond)
	ue.Remove()
}
func TestSendNGResetAcknowledge(t *testing.T) {

	ue := initUe(t, "10.0.0.1", "", 8080, 1)

	partOfNGInterface := &ngapType.UEAssociatedLogicalNGConnectionList{
		List: []ngapType.UEAssociatedLogicalNGConnectionItem{
			{
				AMFUENGAPID: &ngapType.AMFUENGAPID{
					Value: 1,
				},
				RANUENGAPID: &ngapType.RANUENGAPID{
					Value: 1,
				},
			},
		},
	}
	ngap_message.SendNGResetAcknowledge(testAmf, partOfNGInterface, nil)
	time.Sleep(100 * time.Millisecond)
	ue.Remove()
}

func TestSendAMFConfigurationUpdateAcknowledge(t *testing.T) {

	aMFTNLAssociationSetupList := new(ngapType.AMFTNLAssociationSetupList)

	//	AMF TNL Association Setup Item
	addr := ngapConvert.IPAddressToNgap("10.0.0.2", "")
	aMFTNLAssociationSetupItem := ngapType.AMFTNLAssociationSetupItem{}
	aMFTNLAssociationSetupItem.AMFTNLAssociationAddress.Present = ngapType.CPTransportLayerInformationPresentEndpointIPAddress
	aMFTNLAssociationSetupItem.AMFTNLAssociationAddress.EndpointIPAddress = &addr
	aMFTNLAssociationSetupList.List = append(aMFTNLAssociationSetupList.List, aMFTNLAssociationSetupItem)

	ngap_message.SendAMFConfigurationUpdateAcknowledge(testAmf, aMFTNLAssociationSetupList, nil, nil)
	time.Sleep(100 * time.Millisecond)
}

func TestSendAMFConfigurationUpdateFailure(t *testing.T) {

	cause := buildCause(ngapType.CausePresentProtocol, ngapType.CauseProtocolPresentAbstractSyntaxErrorReject)
	ngap_message.SendAMFConfigurationUpdateFailure(testAmf, *cause, nil, nil)
	time.Sleep(100 * time.Millisecond)
}

func TestSendInitialUEMessage(t *testing.T) {

	ue := initUe(t, "10.0.0.1", "", 8080, 1)

	mobileIdentity5GS := nasType.MobileIdentity5GS{
		Len:    12, // suci
		Buffer: []uint8{0x01, 0x02, 0xf8, 0x39, 0xf0, 0x00, 0x00, 0x00, 0x00, 0x00, 0x47, 0x78},
	}
	nasPdu := nasTestpacket.GetRegistrationRequest(nasMessage.RegistrationType5GSInitialRegistration, mobileIdentity5GS, nil, nil)
	ngap_message.SendInitialUEMessage(testAmf, ue, nasPdu)

	time.Sleep(100 * time.Millisecond)
	ue.Remove()

}

func TestSendRanConfigurationUpdate(t *testing.T) {

	ngap_message.SendRANConfigurationUpdate(testAmf)
	time.Sleep(100 * time.Millisecond)

}

func TestSendUERadioCapabilityCheckResponse(t *testing.T) {

	ue := initUe(t, "10.0.0.1", "", 8080, 1)

	ngap_message.SendUERadioCapabilityCheckResponse(testAmf, ue, nil)
}

func TestSendNASNonDeliveryIndication(t *testing.T) {

	ue := initUe(t, "10.0.0.1", "", 8080, 1)
	cause := buildCause(ngapType.CausePresentRadioNetwork, ngapType.CauseRadioNetworkPresentNgIntraSystemHandoverTriggered)
	ngap_message.SendNASNonDeliveryIndication(testAmf, ue, aper.OctetString("\x01\x02\x03"), *cause)
	time.Sleep(100 * time.Millisecond)
}

func buildCriticalityDiagnostics(
	procedureCode *int64,
	triggeringMessage *aper.Enumerated,
	procedureCriticality *aper.Enumerated,
	iesCriticalityDiagnostics *ngapType.CriticalityDiagnosticsIEList) (criticalityDiagnostics ngapType.CriticalityDiagnostics) {

	if procedureCode != nil {
		criticalityDiagnostics.ProcedureCode = new(ngapType.ProcedureCode)
		criticalityDiagnostics.ProcedureCode.Value = *procedureCode
	}

	if triggeringMessage != nil {
		criticalityDiagnostics.TriggeringMessage = new(ngapType.TriggeringMessage)
		criticalityDiagnostics.TriggeringMessage.Value = *triggeringMessage
	}

	if procedureCriticality != nil {
		criticalityDiagnostics.ProcedureCriticality = new(ngapType.Criticality)
		criticalityDiagnostics.ProcedureCriticality.Value = *procedureCriticality
	}

	if iesCriticalityDiagnostics != nil {
		criticalityDiagnostics.IEsCriticalityDiagnostics = iesCriticalityDiagnostics
	}

	return criticalityDiagnostics
}

func buildCriticalityDiagnosticsIEItem(ieCriticality aper.Enumerated, ieID int64, typeOfErr aper.Enumerated) (item ngapType.CriticalityDiagnosticsIEItem) {

	item = ngapType.CriticalityDiagnosticsIEItem{
		IECriticality: ngapType.Criticality{
			Value: ieCriticality,
		},
		IEID: ngapType.ProtocolIEID{
			Value: ieID,
		},
		TypeOfError: ngapType.TypeOfError{
			Value: typeOfErr,
		},
	}

	return item
}

func buildCause(present int, value aper.Enumerated) (cause *ngapType.Cause) {
	cause = new(ngapType.Cause)
	cause.Present = present

	switch present {
	case ngapType.CausePresentRadioNetwork:
		cause.RadioNetwork = new(ngapType.CauseRadioNetwork)
		cause.RadioNetwork.Value = value
	case ngapType.CausePresentTransport:
		cause.Transport = new(ngapType.CauseTransport)
		cause.Transport.Value = value
	case ngapType.CausePresentNas:
		cause.Nas = new(ngapType.CauseNas)
		cause.Nas.Value = value
	case ngapType.CausePresentProtocol:
		cause.Protocol = new(ngapType.CauseProtocol)
		cause.Protocol.Value = value
	case ngapType.CausePresentMisc:
		cause.Misc = new(ngapType.CauseMisc)
		cause.Misc.Value = value
	case ngapType.CausePresentNothing:
	}

	return
}

func initUe(t *testing.T, v4, v6 string, port int32, amfUeNgapId int64) *n3iwf_context.N3IWFUe {
	ue := n3iwfSelf.NewN3iwfUe()
	err := ue.AttachAMF(testAddrAndPort)
	assert.True(t, err == nil)
	ue.IPAddrv4, ue.IPAddrv6, ue.PortNumber = v4, v6, port
	ue.AmfUeNgapId = amfUeNgapId
	return ue

}
