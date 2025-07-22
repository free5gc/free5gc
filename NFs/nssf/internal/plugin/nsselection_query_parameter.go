/*
 * NSSF Plugin
 */

package plugin

import (
	"github.com/free5gc/openapi/models"
)

type NsselectionQueryParameter struct {
	NfType *models.NrfNfManagementNfType `json:"nf-type"`

	NfId string `json:"nf-id"`

	SliceInfoRequestForRegistration *models.SliceInfoForRegistration `json:"slice-info-request-for-registration,omitempty"`

	SliceInfoRequestForPduSession *models.SliceInfoForPduSession `json:"slice-info-request-for-pdu-session,omitempty"`

	HomePlmnId *models.PlmnId `json:"home-plmn-id,omitempty"`

	Tai *models.Tai `json:"tai,omitempty"`

	SupportedFeatures string `json:"supported-features,omitempty"`
}
