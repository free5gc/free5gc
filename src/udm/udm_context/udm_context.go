package udm_context

import (
	// "fmt"
	"fmt"
	"free5gc/lib/Nnrf_NFDiscovery"
	"free5gc/lib/openapi/models"
	"strconv"
	"strings"
)

var udmContext UDMContext

const (
	LocationUriAmf3GppAccessRegistration int = iota
	LocationUriAmfNon3GppAccessRegistration
	LocationUriSmfRegistration
)

func init() {
	UDM_Self().UdmUePool = make(map[string]*UdmUeContext)
	UDM_Self().NfService = make(map[models.ServiceName]models.NfService)
}

type UDMContext struct {
	Name              string
	NfId              string
	GroupId           string
	HttpIpv4Port      int
	HttpIPv4Address   string
	UriScheme         models.UriScheme
	NfService         map[models.ServiceName]models.NfService
	NFDiscoveryClient *Nnrf_NFDiscovery.APIClient
	UdmUePool         map[string]*UdmUeContext // supi as key
	NrfUri            string
}

type UdmUeContext struct {
	Supi                         string
	Amf3GppAccessRegistration    *models.Amf3GppAccessRegistration
	AmfNon3GppAccessRegistration *models.AmfNon3GppAccessRegistration
	PduSessionID                 string
	UdrUri                       string
}

func CreateUdmUe(Supi string) (udmUe *UdmUeContext) {
	udmUe = new(UdmUeContext)
	udmUe.Supi = Supi
	UDM_Self().UdmUePool[Supi] = udmUe
	return
}

func UdmAmf3gppRegContextNotExists(Supi string) bool {
	udmUe := UDM_Self().UdmUePool[Supi]
	if udmUe != nil {
		return udmUe.Amf3GppAccessRegistration == nil
	}
	return true
}

func UdmAmfNon3gppRegContextNotExists(Supi string) bool {
	udmUe := UDM_Self().UdmUePool[Supi]
	if udmUe != nil {
		return udmUe.AmfNon3GppAccessRegistration == nil
	}
	return true
}

func UdmSmfRegContextNotExists(Supi string) bool {
	udmUe := UDM_Self().UdmUePool[Supi]
	if udmUe != nil {
		return udmUe.PduSessionID == ""
	}
	return true
}

func CreateAmf3gppRegContext(Supi string, body models.Amf3GppAccessRegistration) {
	udmUe := UDM_Self().UdmUePool[Supi]
	if udmUe == nil {
		udmUe = CreateUdmUe(Supi)
	}
	udmUe.Amf3GppAccessRegistration = &body
}

func CreateAmfNon3gppRegContext(Supi string, body models.AmfNon3GppAccessRegistration) {
	udmUe := UDM_Self().UdmUePool[Supi]
	if udmUe == nil {
		udmUe = CreateUdmUe(Supi)
	}
	udmUe.AmfNon3GppAccessRegistration = &body
}

func CreateSmfRegContext(Supi string, pduSessionID string) {
	udmUe := UDM_Self().UdmUePool[Supi]
	if udmUe == nil {
		udmUe = CreateUdmUe(Supi)
	}
	if udmUe.PduSessionID == "" {
		udmUe.PduSessionID = pduSessionID
	}
}

func GetAmf3gppRegContext(Supi string) *models.Amf3GppAccessRegistration {
	udmUe := UDM_Self().UdmUePool[Supi]
	if udmUe != nil {
		return udmUe.Amf3GppAccessRegistration
	}
	return nil
}

func GetAmfNon3gppRegContext(Supi string) *models.AmfNon3GppAccessRegistration {
	udmUe := UDM_Self().UdmUePool[Supi]
	if udmUe != nil {
		return udmUe.AmfNon3GppAccessRegistration
	}
	return nil
}
func GetSmfRegContext(Supi string) string {
	udmUe := UDM_Self().UdmUePool[Supi]
	if udmUe != nil {
		return udmUe.PduSessionID
	}
	return ""
}

func (ue *UdmUeContext) GetLocationURI(types int) string {
	switch types {
	case LocationUriAmf3GppAccessRegistration:
		return UDM_Self().GetIPv4Uri() + "/nudm-uecm/v1/" + ue.Supi + "/registrations/amf-3gpp-access"
	case LocationUriAmfNon3GppAccessRegistration:
		return UDM_Self().GetIPv4Uri() + "/nudm-uecm/v1/" + ue.Supi + "/registrations/amf-non-3gpp-access"
	case LocationUriSmfRegistration:
		return UDM_Self().GetIPv4Uri() + "/nudm-uecm/v1/" + ue.Supi + "/registrations/smf-registrations/" + ue.PduSessionID
	}
	return ""
}

func (ue *UdmUeContext) SameAsStoredGUAMI3gpp(inGuami models.Guami) bool {
	if ue.Amf3GppAccessRegistration == nil {
		return false
	}
	ug := ue.Amf3GppAccessRegistration.Guami
	if ug != nil {
		if (ug.PlmnId == nil) == (inGuami.PlmnId == nil) {
			if ug.PlmnId != nil && ug.PlmnId.Mcc == inGuami.PlmnId.Mcc && ug.PlmnId.Mnc == inGuami.PlmnId.Mnc {
				if ug.AmfId == inGuami.AmfId {
					return true
				}
			}
		}
	}
	return false
}

func (ue *UdmUeContext) SameAsStoredGUAMINon3gpp(inGuami models.Guami) bool {
	if ue.AmfNon3GppAccessRegistration == nil {
		return false
	}
	ug := ue.AmfNon3GppAccessRegistration.Guami
	if ug != nil {
		if (ug.PlmnId == nil) == (inGuami.PlmnId == nil) {
			if ug.PlmnId != nil && ug.PlmnId.Mcc == inGuami.PlmnId.Mcc && ug.PlmnId.Mnc == inGuami.PlmnId.Mnc {
				if ug.AmfId == inGuami.AmfId {
					return true
				}
			}
		}
	}
	return false
}

func (context *UDMContext) GetIPv4Uri() string {
	return fmt.Sprintf("%s://%s:%d", context.UriScheme, context.HttpIPv4Address, context.HttpIpv4Port)
}

func (context *UDMContext) InitNFService(serviceName []string, version string) {
	tmpVersion := strings.Split(version, ".")
	versionUri := "v" + tmpVersion[0]
	for index, nameString := range serviceName {
		name := models.ServiceName(nameString)
		context.NfService[name] = models.NfService{
			ServiceInstanceId: strconv.Itoa(index),
			ServiceName:       name,
			Versions: &[]models.NfServiceVersion{
				{
					ApiFullVersion:  version,
					ApiVersionInUri: versionUri,
				},
			},
			Scheme:          context.UriScheme,
			NfServiceStatus: models.NfServiceStatus_REGISTERED,
			ApiPrefix:       context.GetIPv4Uri(),
			IpEndPoints: &[]models.IpEndPoint{
				{
					Ipv4Address: context.HttpIPv4Address,
					Transport:   models.TransportProtocol_TCP,
					Port:        int32(context.HttpIpv4Port),
				},
			},
		}
	}
}

func UDM_Self() *UDMContext {
	return &udmContext
}
