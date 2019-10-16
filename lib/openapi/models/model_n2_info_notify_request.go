/*
 * Namf_Communication
 *
 * AMF Communication Service
 *
 * API version: 1.0.0
 * Manually Created
 */

package models

type N2InfoNotifyRequest struct {
	JsonData                *N2InformationNotification `json:"jsonData,omitempty" multipart:"contentType:application/json"`
	BinaryDataN1Message     []byte                     `json:"binaryDataN1Message,omitempty" multipart:"contentType:application/vnd.3gpp.5gnas,ref:{N1Message}"`
	BinaryDataN2Information []byte                     `json:"binaryDataN2Information,omitempty" multipart:"contentType:application/vnd.3gpp.ngap,class:JsonData.N2InfoContainer.N2InformationClass,ref:(N2InfoContent).NgapData.ContentId"`
}
