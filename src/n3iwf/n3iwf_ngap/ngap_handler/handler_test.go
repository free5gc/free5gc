package ngap_handler_test

import (
	"encoding/hex"
	"fmt"
	"git.cs.nctu.edu.tw/calee/sctp"
	"github.com/stretchr/testify/assert"
	"free5gc/lib/aper"
	"free5gc/lib/ngap"
	"free5gc/lib/ngap/ngapConvert"
	"free5gc/lib/ngap/ngapSctp"
	"free5gc/lib/ngap/ngapType"
	"free5gc/lib/path_util"
	"free5gc/src/n3iwf/factory"
	"free5gc/src/n3iwf/n3iwf_context"
	"free5gc/src/n3iwf/n3iwf_handler"
	"free5gc/src/n3iwf/n3iwf_ngap/n3iwf_sctp"
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
var testConns *sctp.SCTPConn
var n3iwfSelf *n3iwf_context.N3IWFContext
var testAmf *n3iwf_context.N3IWFAMF

func init() {

	// parse config file
	configFile := path_util.Gofree5gcPath("free5gc/src/n3iwf/n3iwf_util/test/testN3iwfcfg.conf")
	factory.InitConfigFactory(configFile)
	testConfig = factory.N3iwfConfig.Configuration

	// setup sctp server(AMF Side)
	testAddrAndPort = fmt.Sprintf("%s:38412", testConfig.AMFAddress[0].NetworkAddress)
	listenConn = ngapSctp.Server(testConfig.AMFAddress[0].NetworkAddress)
	go func() {
		for {
			readChan := make(chan ngapSctp.ConnData, 1024)
			conn, err := ngapSctp.Accept(listenConn)
			if err != nil {
				log.Panicf("Accept error [%s]", err.Error())
			}
			testConns = conn
			fmt.Printf("SCTP Accept from: %s", conn.RemoteAddr().String())
			go ngapSctp.Start(conn, readChan)
		}
	}()

	// init n3iwf context
	n3iwf_util.InitN3IWFContext()

	go n3iwf_handler.Handle()

	n3iwfSelf = n3iwf_context.N3IWFSelf()

	wg := sync.WaitGroup{}
	n3iwf_sctp.SetupSCTPConnection(testConfig.AMFAddress, &wg)
	testAmf = n3iwfSelf.NewN3iwfAmf(testAddrAndPort)

	time.Sleep(100 * time.Millisecond)

}

func TestHandleNGSetupResponse(t *testing.T) {

	plmnId := []byte("\x02\xf8\x39")
	guami := newGuami(plmnId, "cafe00")

	guamiList := []ngapType.ServedGUAMIItem{
		{
			GUAMI: guami,
		},
	}
	plmnList := []ngapType.PLMNSupportItem{
		{
			PLMNIdentity: ngapType.PLMNIdentity{
				Value: plmnId,
			},
			SliceSupportList: ngapType.SliceSupportList{
				List: []ngapType.SliceSupportItem{
					{
						SNSSAI: newNssai(1, ""),
					},
				},
			},
		},
	}

	pdu := ngapTestpacket.BuildNGSetupResponse("amf", guamiList, plmnList, 0xff)
	message, err := ngap.Encoder(pdu)

	assert.True(t, err == nil)

	_, err = testConns.Write(message)
	assert.True(t, err == nil)

	time.Sleep(200 * time.Millisecond)

}

func TestHandlePDUSessionResourceModifyConfirm(t *testing.T) {
	// init Ue Context and pdu session
	ue := initUe(t, "10.0.0.1", "", 8080, 1)

	snssai := newNssai(1, "010203")
	{
		sess, err := ue.CreatePDUSession(10, snssai)
		assert.True(t, err == nil)
		sess.QosFlows[1] = &n3iwf_context.QosFlow{
			Identifier: 1,
		}
	}
	{
		_, err := ue.CreatePDUSession(11, snssai)
		assert.True(t, err == nil)
	}

	pduSessionResourceModifyConfirmList := ngapType.PDUSessionResourceModifyListModCfm{}
	{
		item := ngapType.PDUSessionResourceModifyItemModCfm{
			PDUSessionID: ngapType.PDUSessionID{
				Value: 10,
			},
			PDUSessionResourceModifyConfirmTransfer: ngapTestpacket.GetPDUSessionResourceModifyConfirmTransfer([]int64{1}),
		}
		pduSessionResourceModifyConfirmList.List = append(pduSessionResourceModifyConfirmList.List, item)
	}
	pduSessionResourceFailedToModifyListModCfm := ngapType.PDUSessionResourceFailedToModifyListModCfm{}
	{
		item := ngapType.PDUSessionResourceFailedToModifyItemModCfm{
			PDUSessionID: ngapType.PDUSessionID{
				Value: 11,
			},
			PDUSessionResourceModifyIndicationUnsuccessfulTransfer: ngapTestpacket.GetPDUSessionResourceModifyIndicationUnsuccessfulTransfer(),
		}
		pduSessionResourceFailedToModifyListModCfm.List = append(pduSessionResourceFailedToModifyListModCfm.List, item)
	}
	pdu := ngapTestpacket.BuildPDUSessionResourceModifyConfirm(ue.AmfUeNgapId, ue.RanUeNgapId, pduSessionResourceModifyConfirmList, pduSessionResourceFailedToModifyListModCfm, nil)
	message, err := ngap.Encoder(pdu)

	assert.True(t, err == nil)

	_, err = testConns.Write(message)
	assert.True(t, err == nil)

	time.Sleep(200 * time.Millisecond)

	_, exist := ue.PduSessionList[11]
	assert.True(t, !exist)

	ue.Remove()

}

func TestHandlePDUSessionResourceRelCmd(t *testing.T) {
	// init Ue Context and pdu session
	ue := initUe(t, "10.0.0.1", "", 8080, 1)

	snssai := newNssai(1, "010203")
	{
		_, err := ue.CreatePDUSession(10, snssai)
		assert.True(t, err == nil)
	}
	pDUSessionResourceToReleaseListRelCmd := ngapType.PDUSessionResourceToReleaseListRelCmd{}
	{
		item := ngapType.PDUSessionResourceToReleaseItemRelCmd{
			PDUSessionID: ngapType.PDUSessionID{
				Value: 10,
			},
			PDUSessionResourceReleaseCommandTransfer: ngapTestpacket.GetPDUSessionResourceReleaseCommandTransfer(),
		}
		pDUSessionResourceToReleaseListRelCmd.List = append(pDUSessionResourceToReleaseListRelCmd.List, item)
	}
	pdu := ngapTestpacket.BuildPDUSessionResourceReleaseCommand(ue.AmfUeNgapId, ue.RanUeNgapId, nil, nil, pDUSessionResourceToReleaseListRelCmd)
	message, err := ngap.Encoder(pdu)

	assert.True(t, err == nil)

	_, err = testConns.Write(message)
	assert.True(t, err == nil)

	time.Sleep(200 * time.Millisecond)

	_, exist := ue.PduSessionList[10]
	assert.True(t, !exist)

	ue.Remove()

}

func TestHandleNGReset(t *testing.T) {
	// init Ue Context and pdu session
	ue := initUe(t, "10.0.0.1", "", 8080, 1)
	ue2 := initUe(t, "10.0.0.2", "", 8081, 2)

	{
		partOfNGInterface := &ngapType.UEAssociatedLogicalNGConnectionList{
			List: []ngapType.UEAssociatedLogicalNGConnectionItem{
				{
					AMFUENGAPID: &ngapType.AMFUENGAPID{
						Value: ue2.AmfUeNgapId,
					},
					RANUENGAPID: &ngapType.RANUENGAPID{
						Value: ue2.RanUeNgapId,
					},
				},
			},
		}
		pdu := ngapTestpacket.BuildNGReset(partOfNGInterface)
		message, err := ngap.Encoder(pdu)

		assert.True(t, err == nil)
		_, err = testConns.Write(message)
		assert.True(t, err == nil)

		time.Sleep(100 * time.Millisecond)

		_, exist := n3iwfSelf.UePool[ue.RanUeNgapId]
		assert.True(t, exist)
		_, exist = n3iwfSelf.UePool[ue2.RanUeNgapId]
		assert.True(t, !exist)
	}
	{
		pdu := ngapTestpacket.BuildNGReset(nil)
		message, err := ngap.Encoder(pdu)

		assert.True(t, err == nil)
		_, err = testConns.Write(message)
		assert.True(t, err == nil)

		time.Sleep(100 * time.Millisecond)

		_, exist := n3iwfSelf.UePool[ue.RanUeNgapId]
		assert.True(t, !exist)
	}

	ue.Remove()

}

func TestHandleNGResetAcknowledge(t *testing.T) {

	pdu := ngapTestpacket.BuildNGResetAcknowledge()
	message, err := ngap.Encoder(pdu)

	assert.True(t, err == nil)
	_, err = testConns.Write(message)
	assert.True(t, err == nil)

	time.Sleep(100 * time.Millisecond)

}

func TestHandleAMFConfigurationUpdate(t *testing.T) {

	plmnId := []byte("\x02\xf8\x39")
	guami := newGuami(plmnId, "cafe00")

	guamiList := []ngapType.ServedGUAMIItem{
		{
			GUAMI: guami,
		},
	}
	plmnList := []ngapType.PLMNSupportItem{
		{
			PLMNIdentity: ngapType.PLMNIdentity{
				Value: plmnId,
			},
			SliceSupportList: ngapType.SliceSupportList{
				List: []ngapType.SliceSupportItem{
					{
						SNSSAI: newNssai(1, ""),
					},
				},
			},
		},
	}
	addr := ngapConvert.IPAddressToNgap("10.0.0.2", "")
	addList := &ngapType.AMFTNLAssociationToAddList{}
	addItem := ngapType.AMFTNLAssociationToAddItem{}
	addItem.AMFTNLAssociationAddress.Present = ngapType.CPTransportLayerInformationPresentEndpointIPAddress
	addItem.AMFTNLAssociationAddress.EndpointIPAddress = &addr
	addList.List = append(addList.List, addItem)

	pdu := ngapTestpacket.BuildAMFConfigurationUpdate("amf", guamiList, plmnList, 0xff, addList, nil, nil)
	message, err := ngap.Encoder(pdu)

	assert.True(t, err == nil)
	_, err = testConns.Write(message)
	assert.True(t, err == nil)

	time.Sleep(100 * time.Millisecond)

}

func TestHandleRANConfigurationUpdateAcknowledge(t *testing.T) {

	pdu := ngapTestpacket.BuildRanConfigurationUpdateAck(nil)
	message, err := ngap.Encoder(pdu)

	assert.True(t, err == nil)
	_, err = testConns.Write(message)
	assert.True(t, err == nil)

	time.Sleep(100 * time.Millisecond)

}

func TestHandleRanConfigurationUpdateFailure(t *testing.T) {

	wait := ngapType.TimeToWait{
		Value: ngapType.TimeToWaitPresentV1s,
	}
	pdu := ngapTestpacket.BuildRanConfigurationUpdateFailure(&wait, nil)
	message, err := ngap.Encoder(pdu)

	assert.True(t, err == nil)
	_, err = testConns.Write(message)
	assert.True(t, err == nil)

	time.Sleep(1 * time.Second)

	time.Sleep(200 * time.Millisecond)

}

func TestHandleOverload(t *testing.T) {

	action := ngapType.OverloadAction{
		Value: ngapType.OverloadActionPresentRejectRrcCrSignalling,
	}
	var ind int64 = 50
	snssai1 := newNssai(1, "010203")
	snssai2 := newNssai(1, "112233")

	item := ngapType.OverloadStartNSSAIItem{}
	sliceItem := ngapType.SliceOverloadItem{
		SNSSAI: snssai1,
	}
	item.SliceOverloadList.List = append(item.SliceOverloadList.List, sliceItem)
	sliceItem = ngapType.SliceOverloadItem{
		SNSSAI: snssai2,
	}
	item.SliceOverloadList.List = append(item.SliceOverloadList.List, sliceItem)
	item.SliceOverloadResponse = &ngapType.OverloadResponse{
		Present: ngapType.OverloadResponsePresentOverloadAction,
		OverloadAction: &ngapType.OverloadAction{
			Value: ngapType.OverloadActionPresentRejectRrcCrSignalling,
		},
	}
	item.SliceTrafficLoadReductionIndication = &ngapType.TrafficLoadReductionIndication{
		Value: 30,
	}
	list := []ngapType.OverloadStartNSSAIItem{item}
	{
		pdu := ngapTestpacket.BuildOverloadStart(&action, &ind, list)
		message, err := ngap.Encoder(pdu)
		assert.True(t, err == nil)

		_, err = testConns.Write(message)
		assert.True(t, err == nil)

		time.Sleep(100 * time.Millisecond)
		assert.True(t, testAmf.AMFOverloadContent != nil)
	}
	{
		pdu := ngapTestpacket.BuildOverloadStop()
		message, err := ngap.Encoder(pdu)
		assert.True(t, err == nil)

		_, err = testConns.Write(message)
		assert.True(t, err == nil)

		time.Sleep(200 * time.Millisecond)
		assert.True(t, testAmf.AMFOverloadContent == nil)
	}

}

func TestHandleUERadioCapabilityCheckRequest(t *testing.T) {

	ue := initUe(t, "10.0.0.1", "", 8080, 1)

	pdu := ngapTestpacket.BuildUERadioCapabilityCheckRequest(ue.AmfUeNgapId, ue.RanUeNgapId)
	message, err := ngap.Encoder(pdu)

	assert.True(t, err == nil)
	_, err = testConns.Write(message)
	assert.True(t, err == nil)

	time.Sleep(100 * time.Millisecond)

}

func newNssai(sst int32, sd string) (snssai ngapType.SNSSAI) {
	snssai.SST.Value = []byte{byte(sst)}
	if sd != "" {
		snssai.SD = new(ngapType.SD)
		snssai.SD.Value, _ = hex.DecodeString(sd)
	}
	return
}
func newGuami(plmnId aper.OctetString, amfId string) ngapType.GUAMI {
	regId, setId, ptrId := ngapConvert.AmfIdToNgap(amfId)
	return ngapType.GUAMI{
		PLMNIdentity: ngapType.PLMNIdentity{
			Value: plmnId,
		},
		AMFRegionID: ngapType.AMFRegionID{
			Value: regId,
		},
		AMFSetID: ngapType.AMFSetID{
			Value: setId,
		},
		AMFPointer: ngapType.AMFPointer{
			Value: ptrId,
		},
	}
}

func initUe(t *testing.T, v4, v6 string, port int32, amfUeNgapId int64) *n3iwf_context.N3IWFUe {
	ue := n3iwfSelf.NewN3iwfUe()
	err := ue.AttachAMF(testAddrAndPort)
	assert.True(t, err == nil)
	ue.IPAddrv4, ue.IPAddrv6, ue.PortNumber = v4, v6, port
	ue.AmfUeNgapId = amfUeNgapId
	return ue

}
