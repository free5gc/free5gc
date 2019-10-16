/*
 * NSSF Plugin
 */

package plugin

import (
	. "free5gc/lib/openapi/models"
)

type NsselectionQueryParameter struct {
	NfType *NfType `json:"nf-type"`

	NfId string `json:"nf-id"`

	SliceInfoRequestForRegistration *SliceInfoForRegistration `json:"slice-info-request-for-registration,omitempty"`

	SliceInfoRequestForPduSession *SliceInfoForPduSession `json:"slice-info-request-for-pdu-session,omitempty"`

	HomePlmnId *PlmnId `json:"home-plmn-id,omitempty"`

	Tai *Tai `json:"tai,omitempty"`

	SupportedFeatures string `json:"supported-features,omitempty"`
}
