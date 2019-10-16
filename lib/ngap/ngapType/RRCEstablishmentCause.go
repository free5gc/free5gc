//go:binary-only-package

package ngapType

import "free5gc/lib/aper"

// Need to import "free5gc/lib/aper" if it uses "aper"

const (
	RRCEstablishmentCausePresentEmergency          aper.Enumerated = 0
	RRCEstablishmentCausePresentHighPriorityAccess aper.Enumerated = 1
	RRCEstablishmentCausePresentMtAccess           aper.Enumerated = 2
	RRCEstablishmentCausePresentMoSignalling       aper.Enumerated = 3
	RRCEstablishmentCausePresentMoData             aper.Enumerated = 4
	RRCEstablishmentCausePresentMoVoiceCall        aper.Enumerated = 5
	RRCEstablishmentCausePresentMoVideoCall        aper.Enumerated = 6
	RRCEstablishmentCausePresentMoSMS              aper.Enumerated = 7
	RRCEstablishmentCausePresentMpsPriorityAccess  aper.Enumerated = 8
	RRCEstablishmentCausePresentMcsPriorityAccess  aper.Enumerated = 9
)

type RRCEstablishmentCause struct {
	Value aper.Enumerated `aper:"valueExt,valueLB:0,valueUB:9"`
}
