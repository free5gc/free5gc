/*
 * Namf_Communication
 *
 * AMF Communication Service
 *
 * API version: 1.0.0
 * Manually Created
 */

package models

type N2InfoNotifyResponse struct {
	JsonData                *N2InfoNotifyRspData `json:"jsonData,omitempty" multipart:"contentType:application/json"`
	BinaryDataN2Information []byte               `json:"binaryDataN2Information,omitempty" multipart:"contentType:application/vnd.3gpp.ngap,ref:JsonData.N2InfoContent.NgapData.ContentId"`
}
