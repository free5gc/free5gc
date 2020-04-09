package amf_consumer_test

import (
	"flag"
	"github.com/urfave/cli"
	"free5gc/lib/CommonConsumerTestData/AMF/TestAmf"
	"free5gc/lib/openapi/models"
	"free5gc/src/amf/amf_consumer"
	"free5gc/src/nssf/nssf_service"
	"testing"
	"time"
)

func nssfInit() {
	flags := flag.FlagSet{}
	c := cli.NewContext(nil, &flags, nil)
	nssf := &nssf_service.NSSF{}
	nssf.Initialize(c)
	go nssf.Start()
	time.Sleep(100 * time.Millisecond)
}

func TestNSSelectionGetForRegistration(t *testing.T) {
	nrfInit()
	nssfInit()

	TestAmf.AmfInit()
	TestAmf.UeAttach(models.AccessType__3_GPP_ACCESS)
	ue := TestAmf.TestAmf.UePool["imsi-2089300007487"]

	ue.NssfUri = "https://localhost:29531"
	requestNssai := []models.Snssai{
		{
			Sst: 1,
			Sd:  "010203",
		},
	}
	ue.SubscribedNssai = []models.SubscribedSnssai{
		{
			SubscribedSnssai: &models.Snssai{
				Sst: 1,
				Sd:  "010203",
			},
			DefaultIndication: true,
		},
		{
			SubscribedSnssai: &models.Snssai{
				Sst: 1,
				Sd:  "112233",
			},
			DefaultIndication: true,
		},
		{
			SubscribedSnssai: &models.Snssai{
				Sst: 1,
				Sd:  "010203",
			},
			DefaultIndication: false,
		},
		{
			SubscribedSnssai: &models.Snssai{
				Sst: 1,
				Sd:  "112233",
			},
			DefaultIndication: false,
		},
	}
	problemDetails, err := amf_consumer.NSSelectionGetForRegistration(ue, requestNssai)
	if err != nil {
		t.Error(err.Error())
	} else if problemDetails != nil {
		t.Logf("ProblemDetails: %+v", problemDetails)
	} else {
		for _, allowedSnssai := range ue.AllowedNssai[models.AccessType__3_GPP_ACCESS] {
			t.Logf("AllowedSnssai[3GPP_ACCESS]: %+v", *allowedSnssai.AllowedSnssai)
		}
	}
}

func TestNSSelectionGetForPduSession(t *testing.T) {
	nrfInit()
	nssfInit()

	TestAmf.AmfInit()
	TestAmf.UeAttach(models.AccessType__3_GPP_ACCESS)
	ue := TestAmf.TestAmf.UePool["imsi-2089300007487"]

	ue.NssfUri = "https://localhost:29531"
	requestNssai := []models.Snssai{
		{
			Sst: 1,
			Sd:  "010203",
		},
	}
	ue.SubscribedNssai = []models.SubscribedSnssai{
		{
			SubscribedSnssai: &models.Snssai{
				Sst: 1,
				Sd:  "010203",
			},
			DefaultIndication: true,
		},
		{
			SubscribedSnssai: &models.Snssai{
				Sst: 1,
				Sd:  "112233",
			},
			DefaultIndication: true,
		},
		{
			SubscribedSnssai: &models.Snssai{
				Sst: 1,
				Sd:  "010203",
			},
			DefaultIndication: false,
		},
		{
			SubscribedSnssai: &models.Snssai{
				Sst: 1,
				Sd:  "112233",
			},
			DefaultIndication: false,
		},
	}
	problemDetails, err := amf_consumer.NSSelectionGetForRegistration(ue, requestNssai)
	if err != nil {
		t.Error(err.Error())
	} else if problemDetails != nil {
		t.Logf("ProblemDetails: %+v", problemDetails)
	} else {
		for _, allowedSnssai := range ue.AllowedNssai[models.AccessType__3_GPP_ACCESS] {
			t.Logf("AllowedSnssai[3GPP_ACCESS]: %+v", *allowedSnssai.AllowedSnssai)
		}
	}

	snssai := models.Snssai{
		Sst: ue.AllowedNssai[models.AccessType__3_GPP_ACCESS][0].AllowedSnssai.Sst,
		Sd:  ue.AllowedNssai[models.AccessType__3_GPP_ACCESS][0].AllowedSnssai.Sd,
	}
	res, problemDetails, err := amf_consumer.NSSelectionGetForPduSession(ue, snssai)
	if err != nil {
		t.Error(err.Error())
	} else if problemDetails != nil {
		t.Logf("ProblemDetails: %+v", problemDetails)
	} else {
		t.Logf("response: %+v", res)
	}
}
