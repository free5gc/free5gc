/*
 * NRF UriList
 */

package context

import (
	"github.com/free5gc/openapi/models"
)

type UriList struct {
	NfType models.NfType `json:"nfType" bson:"nfType"`
	Link   Links         `json:"_link" bson:"_link" mapstructure:"_link"`
}
