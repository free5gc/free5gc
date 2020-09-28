/*
 * NRF UriList
 */

package context

import (
	"free5gc/lib/openapi/models"
)

type UriList struct {
	NfType models.NfType `json:"nfType" bson:"nfType"`
	Link   Links         `json:"_link" bson:"_link" mapstructure:"_link"`
}
