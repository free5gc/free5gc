/*
 * NRF UriList
 */

package context

import (
	"github.com/free5gc/openapi/models"
)

type UriList struct {
	NfType models.NrfNfManagementNfType `json:"nfType,omitempty" bson:"nfType,omitempty"`
	Link   Links                        `json:"_link" bson:"_link" mapstructure:"_link"`
}
