package util

import (
	"fmt"

	"github.com/free5gc/openapi/models"
)

// SearchNFServiceUri returns NF Uri derived from NfProfile with corresponding service
func SearchNFServiceUri(nfProfile models.NrfNfDiscoveryNfProfile, serviceName models.ServiceName,
	nfServiceStatus models.NfServiceStatus,
) (nfUri string) {
	if nfProfile.NfServices != nil {
		for _, service := range nfProfile.NfServices {
			if service.ServiceName == serviceName && service.NfServiceStatus == nfServiceStatus {
				if nfProfile.Fqdn != "" {
					nfUri = nfProfile.Fqdn
				} else if service.Fqdn != "" {
					nfUri = service.Fqdn
				} else if service.ApiPrefix != "" {
					nfUri = service.ApiPrefix
				} else if service.IpEndPoints != nil {
					point := (service.IpEndPoints)[0]
					if point.Ipv4Address != "" {
						nfUri = getSbiUri(service.Scheme, point.Ipv4Address, point.Port)
					} else if len(nfProfile.Ipv4Addresses) != 0 {
						nfUri = getSbiUri(service.Scheme, nfProfile.Ipv4Addresses[0], point.Port)
					}
				}
			}
			if nfUri != "" {
				break
			}
		}
	}

	return
}

func getSbiUri(scheme models.UriScheme, ipv4Address string, port int32) (uri string) {
	if port != 0 {
		uri = fmt.Sprintf("%s://%s:%d", scheme, ipv4Address, port)
	} else {
		switch scheme {
		case models.UriScheme_HTTP:
			uri = fmt.Sprintf("%s://%s:80", scheme, ipv4Address)
		case models.UriScheme_HTTPS:
			uri = fmt.Sprintf("%s://%s:443", scheme, ipv4Address)
		}
	}
	return
}
