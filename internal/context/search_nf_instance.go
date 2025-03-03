package context

import (
	"github.com/free5gc/openapi/models"
)

type SearchNFInstances struct {
	TargetNFType            models.NrfNfManagementNfType `form:"target-nf-type" binding:"required"`
	RequesterNFType         models.NrfNfManagementNfType `form:"requester-nf-type" binding:"required"`
	ServiceNames            []models.ServiceName         `form:"service-names" `
	RequesterNfInstanceFqdn string                       `form:"requester-nf-instance-fqdn" `
	TargetPlmnList          []models.PlmnId              `form:"target-plmn-list" `
	RequesterPlmnList       []models.PlmnId              `form:"requester-plmn-list" `
	TargetNfInstanceID      string                       `form:"target-nf-instance-id" `
	TargetNfFqdn            string                       `form:"target-nf-type" `
	HnrfURI                 models.UriScheme             `form:"hnrf-uri" `
	Snssais                 []models.Snssai              `form:"snssais" `
	PlmnSpecificSnssaiList  []models.PlmnSnssai          `form:"plmn-specific-snssai-list"`
	Dnn                     string                       `form:"dnn" `
	Ipv4Index               []string                     `form:"ipv4-index"`
	Ipv6Index               []string                     `form:"ipv6-index"`
	NsiList                 []string                     `form:"nsi-list" `
	SmfServingArea          string                       `form:"smf-serving-area" `
	MbSmfServingArea        string                       `form:"mbsmf-serving-area" `
	Tai                     models.Tai                   `form:"tai" `
	AmfRegionID             string                       `form:"amf-region-id" `
	AmfSetID                string                       `form:"amf-set-id" `
	Guami                   models.Guami                 `form:"guami" `
	Supi                    string                       `form:"supi" `
	UeIpv4Address           string                       `form:"ue-ipv4-address" `
	IPDomain                string                       `form:"ip-domain" `
	UeIpv6Prefix            string                       `form:"ue-ipv6-prefix" `
	PgwInd                  bool                         `form:"pgw-ind" `
	Pgw                     string                       `form:"pgw" `
	PgwIP                   string                       `form:"pgw-ip" `
	Gpsi                    string                       `form:"gpsi" `
	ExternalGroupIdentity   string                       `form:"external-group-identity" `
	InternalGroupIdentity   string                       `form:"internal-group-identity" `
	PfdData                 string                       `form:"pfd-data" `
	DataSet                 models.DataSetId             `form:"data-set" `
	RoutingIndicator        string                       `form:"routing-indicator" `
	GroupIDList             []string                     `form:"group-id-list" `
	DnaiList                []string                     `form:"dnai-list" `
	PduSessionType          []models.PduSessionType      `form:"pdu-session-typs" `
	EventIdList             []string                     `form:"event-id-list" `
	NwdafEventList          []models.NwdafEvent          `form:"nwdaf-event-list" `
	SupportedFeatures       []string                     `form:"supported-features" `
	UpfIwkEpsInd            bool                         `form:"upf-iwk-eps-ind" `
	ChfSupportedPlmn        models.PlmnId                `form:"chf-supported-plmn" `
	PreferredLocality       string                       `form:"preferred-locality" `
	AccessType              models.AccessType            `form:"access-type" `
	// IfNoneMatch             string            `form:"target-nf-type" `
}
