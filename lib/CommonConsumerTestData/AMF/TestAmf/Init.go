package TestAmf

import (
	"free5gc/lib/path_util"
	amf_ngap_sctp "free5gc/src/amf/amf_ngap/ngap_sctp"
	"free5gc/src/amf/amf_util"
	"free5gc/src/amf/factory"
	"free5gc/src/amf/gmm"
	"free5gc/src/amf/logger"
	"log"
	"net"
	"strings"
	"time"

	"github.com/ishidawataru/sctp"

	"github.com/davecgh/go-spew/spew"

	"free5gc/lib/openapi/models"
	"free5gc/src/amf/amf_context"
)

var TestAmf = amf_context.AMF_Self()
var Config = spew.NewDefaultConfig()
var AmfLogPath = path_util.Gofree5gcPath("free5gc/amfsslkey.log")
var AmfPemPath = path_util.Gofree5gcPath("free5gc/support/TLS/amf.pem")
var AmfKeyPath = path_util.Gofree5gcPath("free5gc/support/TLS/amf.key")
var Laddr *sctp.SCTPAddr
var Laddr2 *sctp.SCTPAddr
var Conn *sctp.SCTPConn
var Conn2 *sctp.SCTPConn
var ServerAddr *sctp.SCTPAddr
var configFile = path_util.Gofree5gcPath("free5gc/lib/CommonConsumerTestData/AMF/TestAmf/testAmfCfg.conf")

func AmfInit() {

	Config.DisablePointerAddresses = true
	Config.Indent = "\t"
	Conn = nil
	Conn2 = nil
	TestAmf.Reset()
	// init AMFcontext with test config file
	factory.InitConfigFactory(configFile)
	amf_util.InitAmfContext(TestAmf)

	// add Ue info
	supi := "imsi-2089300007487"
	ue := TestAmf.NewAmfUe(supi)
	if err := gmm.InitAmfUeSm(ue); err != nil {
		log.Printf("InitAmfUeSm error: %v", err)
	}
	ue.GroupID = "12121212-208-93-01010101"
	ue.TimeZone = "+08:00+1h"
	ue.Location = models.UserLocation{
		NrLocation: &models.NrLocation{
			Tai: &models.Tai{
				PlmnId: &models.PlmnId{
					Mcc: "208",
					Mnc: "93",
				},
				Tac: "000001",
			},
			Ncgi: &models.Ncgi{
				PlmnId: &models.PlmnId{
					Mcc: "208",
					Mnc: "93",
				},
				NrCellId: "000000001",
			},
		},
	}
	ue.Tai = *ue.Location.NrLocation.Tai
	ue.AllowedNssai[models.AccessType__3_GPP_ACCESS] = []models.Snssai{
		{
			Sst: 1,
			Sd:  "010203",
		},
	}
	ue.SmfSelectionData = &models.SmfSelectionSubscriptionData{
		SubscribedSnssaiInfos: make(map[string]models.SnssaiInfo),
	}
	ue.SmfSelectionData.SubscribedSnssaiInfos["01010203"] = models.SnssaiInfo{
		DnnInfos: []models.DnnInfo{
			{
				Dnn: "internet",
			},
		},
	}
	ue.AccessAndMobilitySubscriptionData = &models.AccessAndMobilitySubscriptionData{
		SubscribedUeAmbr: &models.AmbrRm{
			Uplink:   "800 Kbps",
			Downlink: "1000 Kbps",
		},
		RatRestrictions: []models.RatType{models.RatType_EUTRA},
	}
	ue.RatType = models.RatType_NR
	ue.Kamf = strings.Repeat("1", 64)
	ue.SecurityCapabilities.NRIntegrityProtectionAlgorithms = [2]byte{0x40, 0x00}
	ue.SecurityCapabilities.NREncryptionAlgorithms = [2]byte{0x40, 0x00}
	ue.CipheringAlg = amf_context.ALG_CIPHERING_128_NEA2
	ue.IntegrityAlg = amf_context.ALG_INTEGRITY_128_NIA2
	ue.DerivateAnKey(models.AccessType__3_GPP_ACCESS)
	ue.DerivateNH(ue.Kgnb)
	ue.DerivateAlgKey()
}

// func addSupportDnnList() {
// 	dnn := "internet"
// 	TestAmf.SupportDnnLists = append(TestAmf.SupportDnnLists, dnn)
// }
// func addSupportTaiList() {
// 	tai := models.Tai{
// 		PlmnId: &models.PlmnId{
// 			Mcc: "208",
// 			Mnc: "93",
// 		},
// 		Tac: "001122",
// 	}
// 	TestAmf.SupportTaiLists = append(TestAmf.SupportTaiLists, tai)
// }
// func addServedGuamiList() {
// 	guami := models.Guami{
// 		PlmnId: &models.PlmnId{
// 			Mcc: "208",
// 			Mnc: "93",
// 		},
// 		AmfId: "CAFE00",
// 	}
// 	if len(TestAmf.ServedGuamiList) < cap(TestAmf.ServedGuamiList) {
// 		TestAmf.ServedGuamiList = append(TestAmf.ServedGuamiList, guami)
// 	}
// }
// func addPlmnSupportList() {
// 	plmnSupportItem := amf_context.NewPlmnSupportItem()
// 	plmnSupportItem.PlmnId = models.PlmnId{
// 		Mcc: "208",
// 		Mnc: "93",
// 	}
// 	addSNssaiList(&plmnSupportItem)
// 	if len(TestAmf.PlmnSupportList) < cap(TestAmf.PlmnSupportList) {
// 		TestAmf.PlmnSupportList = append(TestAmf.PlmnSupportList, plmnSupportItem)
// 	}
// }
// func addSNssaiList(item *amf_context.PlmnSupportItem) {
// 	snssai := models.Snssai{
// 		Sst: 1,
// 		Sd:  "010203",
// 	}
// 	if len(item.SNssaiList) < cap(item.SNssaiList) {
// 		item.SNssaiList = append(item.SNssaiList, snssai)
// 	}
// }

func SctpSever() {
	ipStr := "127.0.0.1"
	ips := []net.IPAddr{}
	if ip, err := net.ResolveIPAddr("ip", ipStr); err != nil {
		logger.NgapLog.Errorf("Error resolving address '%s': %v", ipStr, err)
	} else {
		ips = append(ips, *ip)
	}
	Laddr = &sctp.SCTPAddr{
		IPAddrs: ips,
		Port:    9487,
	}
	Laddr2 = &sctp.SCTPAddr{
		IPAddrs: ips,
		Port:    9488,
	}
	ServerAddr = &sctp.SCTPAddr{
		IPAddrs: ips,
		Port:    38412,
	}
	amf_ngap_sctp.Server("127.0.0.1")
}

func SctpConnectToServer(anType models.AccessType) {

	var err error
	if anType == models.AccessType__3_GPP_ACCESS {
		Conn, err = sctp.DialSCTP("sctp", Laddr, ServerAddr)
		if err != nil {
			logger.NgapLog.Errorf("failed to dial: %v\n", err)
		}
		logger.NgapLog.Printf("Dail LocalAddr: %s; RemoteAddr: %s", Conn.LocalAddr(), Conn.RemoteAddr())
	} else {
		Conn2, err = sctp.DialSCTP("sctp", Laddr2, ServerAddr)
		if err != nil {
			logger.NgapLog.Errorf("failed to dial: %v\n", err)
		}
		logger.NgapLog.Printf("Dail LocalAddr: %s; RemoteAddr: %s", Conn2.LocalAddr(), Conn2.RemoteAddr())

	}

	time.Sleep(10 * time.Millisecond)
	UeAttach(anType)
}

func UeAttach(anType models.AccessType) {
	testConn := &TestConn{
		LAddr: &TestAddr{
			Net:   "sctp",
			Value: "192.188.2.2:12345",
		},
		RAddr: &TestAddr{
			Net:   "sctp",
			Value: "127.0.0.1:9487",
		},
	}
	var ran *amf_context.AmfRan
	if anType == models.AccessType__3_GPP_ACCESS {
		if Conn == nil {
			ran = TestAmf.NewAmfRan(testConn)
		} else {
			ran = TestAmf.AmfRanPool[Laddr.String()]
		}
	} else {
		if Conn2 == nil {
			testConn.RAddr.Value = "127.0.0.1:9488"
			ran = TestAmf.NewAmfRan(testConn)
		} else {
			ran = TestAmf.AmfRanPool[Laddr2.String()]
		}
	}
	ran.AnType = anType
	ranUe := ran.NewRanUe()
	ue := TestAmf.UePool["imsi-2089300007487"]
	ran.SupportedTAList = []amf_context.SupportedTAI{
		{
			Tai:        ue.Tai,
			SNssaiList: ue.AllowedNssai[anType],
		},
	}
	ue.AttachRanUe(ranUe)
}
