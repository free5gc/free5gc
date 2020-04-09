package Namf_Callback_test

import (
	"flag"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli"
	"free5gc/lib/CommonConsumerTestData/AMF/TestAmf"
	"free5gc/lib/http2_util"
	"free5gc/lib/openapi/models"
	"free5gc/src/amf/Communication"
	"free5gc/src/amf/HttpCallback"
	"free5gc/src/amf/amf_consumer"
	"free5gc/src/amf/amf_handler"
	"free5gc/src/nrf/nrf_service"
	"free5gc/src/pcf/pcf_context"
	"free5gc/src/pcf/pcf_producer"
	"free5gc/src/pcf/pcf_service"
	"free5gc/src/udr/udr_service"

	"testing"
	"time"
)

func pcfInit() {
	flags := flag.FlagSet{}
	c := cli.NewContext(nil, &flags, nil)
	pcf := &pcf_service.PCF{}
	pcf.Initialize(c)
	go pcf.Start()
	time.Sleep(100 * time.Millisecond)
}

func nrfInit() {
	flags := flag.FlagSet{}
	c := cli.NewContext(nil, &flags, nil)
	nrf := &nrf_service.NRF{}
	nrf.Initialize(c)
	go nrf.Start()
	time.Sleep(100 * time.Millisecond)
}

func udrInit() {
	flags := flag.FlagSet{}
	c := cli.NewContext(nil, &flags, nil)
	udr := &udr_service.UDR{}
	udr.Initialize(c)
	go udr.Start()
	time.Sleep(100 * time.Millisecond)
}

func TestAmPolicyControlUpdateNotifyUpdate(t *testing.T) {
	nrfInit()
	pcfInit()
	udrInit()

	go func() {
		router := gin.Default()
		Namf_Callback.AddService(router)
		Communication.AddService(router)

		server, err := http2_util.NewServer(":29518", TestAmf.AmfLogPath, router)
		if err == nil && server != nil {
			err = server.ListenAndServeTLS(TestAmf.AmfPemPath, TestAmf.AmfKeyPath)
		}
		assert.True(t, err == nil, err.Error())
	}()
	go amf_handler.Handle()

	TestAmf.AmfInit()
	TestAmf.SctpSever()
	TestAmf.SctpConnectToServer(models.AccessType__3_GPP_ACCESS)
	time.Sleep(100 * time.Millisecond)

	ue := TestAmf.TestAmf.UePool["imsi-2089300007487"]
	ue.PlmnId = models.PlmnId{
		Mcc: "208",
		Mnc: "93",
	}
	ue.PcfUri = "https://localhost:29507"
	ue.AccessAndMobilitySubscriptionData = &models.AccessAndMobilitySubscriptionData{
		RfspIndex: 1,
	}
	problemDetails, err := amf_consumer.AMPolicyControlCreate(ue, models.AccessType__3_GPP_ACCESS)
	if err != nil {
		t.Error(err)
	} else if problemDetails != nil {
		t.Logf("problemDetail: %+v", problemDetails)
	}

	time.Sleep(100 * time.Millisecond)

	pcfUe := pcf_context.UeContext{}
	pcfUe.AMPolicyData = make(map[string]*pcf_context.UeAMPolicyData)
	pcfUe.AMPolicyData[ue.PolicyAssociationId] = new(pcf_context.UeAMPolicyData)
	amPolicyData := pcfUe.AMPolicyData[ue.PolicyAssociationId]
	amPolicyData.NotificationUri = ue.AmPolicyAssociation.Request.NotificationUri + ue.PolicyAssociationId

	req := models.PolicyUpdate{}
	req.Rfsp = 2
	pcf_producer.SendAMPolicyUpdateNotification(&pcfUe, ue.PolicyAssociationId, req)

	time.Sleep(200 * time.Millisecond)
}

func TestAmPolicyControlUpdateNotifyTerminate(t *testing.T) {
	nrfInit()
	pcfInit()
	udrInit()

	go func() {
		router := gin.Default()
		Namf_Callback.AddService(router)
		Communication.AddService(router)

		server, err := http2_util.NewServer(":29518", TestAmf.AmfLogPath, router)
		if err == nil && server != nil {
			err = server.ListenAndServeTLS(TestAmf.AmfPemPath, TestAmf.AmfKeyPath)
		}
		assert.True(t, err == nil, err.Error())
	}()
	go amf_handler.Handle()

	TestAmf.AmfInit()
	TestAmf.SctpSever()
	TestAmf.SctpConnectToServer(models.AccessType__3_GPP_ACCESS)
	time.Sleep(100 * time.Millisecond)

	ue := TestAmf.TestAmf.UePool["imsi-2089300007487"]
	ue.PlmnId = models.PlmnId{
		Mcc: "208",
		Mnc: "93",
	}
	ue.PcfUri = "https://localhost:29507"
	ue.AccessAndMobilitySubscriptionData = &models.AccessAndMobilitySubscriptionData{
		RfspIndex: 1,
	}
	problemDetails, err := amf_consumer.AMPolicyControlCreate(ue, models.AccessType__3_GPP_ACCESS)
	if err != nil {
		t.Error(err)
	} else if problemDetails != nil {
		t.Logf("problemDetail: %+v", problemDetails)
	}

	time.Sleep(100 * time.Millisecond)

	pcfUe := pcf_context.UeContext{}
	pcfUe.AMPolicyData = make(map[string]*pcf_context.UeAMPolicyData)
	pcfUe.AMPolicyData[ue.PolicyAssociationId] = new(pcf_context.UeAMPolicyData)
	amPolicyData := pcfUe.AMPolicyData[ue.PolicyAssociationId]
	amPolicyData.NotificationUri = ue.AmPolicyAssociation.Request.NotificationUri + ue.PolicyAssociationId

	req := models.TerminationNotification{}
	req.Cause = models.PolicyAssociationReleaseCause_UNSPECIFIED
	pcf_producer.SendAMPolicyTerminationRequestNotification(&pcfUe, ue.PolicyAssociationId, req)

	time.Sleep(200 * time.Millisecond)
}
