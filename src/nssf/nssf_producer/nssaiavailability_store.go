/*
 * NSSF NSSAI Availability
 *
 * NSSF NSSAI Availability Service
 */

package nssf_producer

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"

	jsonpatch "github.com/evanphx/json-patch"

	. "free5gc/lib/openapi/models"
	"free5gc/src/nssf/factory"
	"free5gc/src/nssf/logger"
	. "free5gc/src/nssf/plugin"
	"free5gc/src/nssf/util"
)

// NSSAIAvailability DELETE method
func nssaiavailabilityDelete(nfId string, d *ProblemDetails) (status int) {
	for i, amfConfig := range factory.NssfConfig.Configuration.AmfList {
		if amfConfig.NfId == nfId {
			factory.NssfConfig.Configuration.AmfList = append(
				factory.NssfConfig.Configuration.AmfList[:i],
				factory.NssfConfig.Configuration.AmfList[i+1:]...)

			status = http.StatusNoContent
			return
		}
	}

	problemDetail := fmt.Sprintf("AMF ID '%s' does not exist", nfId)
	*d = ProblemDetails{
		Title:  util.UNSUPPORTED_RESOURCE,
		Status: http.StatusNotFound,
		Detail: problemDetail,
	}

	status = http.StatusNotFound
	return
}

// NSSAIAvailability PATCH method
func nssaiavailabilityPatch(nfId string,
	p PatchDocument,
	a *AuthorizedNssaiAvailabilityInfo,
	d *ProblemDetails) (status int) {
	var amfIdx int
	var original []byte
	hitAmf := false
	for amfIdx, amfConfig := range factory.NssfConfig.Configuration.AmfList {
		if amfConfig.NfId == nfId {
			// Since json-patch package does not have idea of optional field of datatype,
			// provide with null or empty value instead of omitting the field
			temp := factory.NssfConfig.Configuration.AmfList[amfIdx].SupportedNssaiAvailabilityData
			const DUMMY_STRING string = "DUMMY"
			for i := range temp {
				for j := range temp[i].SupportedSnssaiList {
					if temp[i].SupportedSnssaiList[j].Sd == "" {
						temp[i].SupportedSnssaiList[j].Sd = DUMMY_STRING
					}
				}
			}
			original, _ = json.Marshal(temp)
			original = bytes.ReplaceAll(original, []byte(DUMMY_STRING), []byte(""))

			// original, _ = json.Marshal(factory.NssfConfig.Configuration.AmfList[amfIdx].SupportedNssaiAvailabilityData)

			hitAmf = true
			break
		}
	}
	if !hitAmf {
		problemDetail := fmt.Sprintf("AMF ID '%s' does not exist", nfId)
		*d = ProblemDetails{
			Title:  util.UNSUPPORTED_RESOURCE,
			Status: http.StatusNotFound,
			Detail: problemDetail,
		}

		status = http.StatusNotFound
		return
	}

	// TODO: Check if returned HTTP status codes or problem details are proper when errors occur

	// Provide JSON string with null or empty value in `Value` of `PatchItem`
	for i, patchItem := range p {
		if reflect.ValueOf(patchItem.Value).Kind() == reflect.Map {
			_, exist := patchItem.Value.(map[string]interface{})["sst"]
			_, notExist := patchItem.Value.(map[string]interface{})["sd"]
			if exist && !notExist {
				p[i].Value.(map[string]interface{})["sd"] = ""
			}
		}
	}
	patchJson, _ := json.Marshal(p)

	patch, err := jsonpatch.DecodePatch(patchJson)
	if err != nil {
		*d = ProblemDetails{
			Title:  util.MALFORMED_REQUEST,
			Status: http.StatusBadRequest,
			Detail: err.Error(),
		}

		status = http.StatusBadRequest
		return
	}

	modified, err := patch.Apply(original)
	if err != nil {
		*d = ProblemDetails{
			Title:  util.INVALID_REQUEST,
			Status: http.StatusConflict,
			Detail: err.Error(),
		}

		status = http.StatusConflict
		return
	}

	err = json.Unmarshal(modified, &factory.NssfConfig.Configuration.AmfList[amfIdx].SupportedNssaiAvailabilityData)
	if err != nil {
		*d = ProblemDetails{
			Title:  util.INVALID_REQUEST,
			Status: http.StatusBadRequest,
			Detail: err.Error(),
		}

		status = http.StatusBadRequest
		return
	}

	// Return all authorized NSSAI availability information
	a.AuthorizedNssaiAvailabilityData, _ = util.AuthorizeOfAmfFromConfig(nfId)

	// TODO: Return authorized NSSAI availability information of updated TAI only

	return http.StatusOK
}

// NSSAIAvailability PUT method
func nssaiavailabilityPut(nfId string,
	n NssaiAvailabilityInfo,
	a *AuthorizedNssaiAvailabilityInfo,
	d *ProblemDetails) (status int) {
	for _, s := range n.SupportedNssaiAvailabilityData {
		if !util.CheckSupportedNssaiInPlmn(s.SupportedSnssaiList, *s.Tai.PlmnId) {
			*d = ProblemDetails{
				Title:  util.UNSUPPORTED_RESOURCE,
				Status: http.StatusForbidden,
				Detail: "S-NSSAI in Requested NSSAI is not supported in PLMN",
				Cause:  "SNSSAI_NOT_SUPPORTED",
			}

			status = http.StatusForbidden
			return
		}
	}

	// TODO: Currently authorize all the provided S-NSSAIs
	//       Take some issue into consideration e.g. operator policies

	hitAmf := false
	// Find AMF configuration of given NfId
	// If found, then update the SupportedNssaiAvailabilityData
	for i, amfConfig := range factory.NssfConfig.Configuration.AmfList {
		if amfConfig.NfId == nfId {
			factory.NssfConfig.Configuration.AmfList[i].SupportedNssaiAvailabilityData = n.SupportedNssaiAvailabilityData

			hitAmf = true
			break
		}
	}

	// If no AMF record is found, create a new one
	if !hitAmf {
		var amfConfig factory.AmfConfig
		amfConfig.NfId = nfId
		amfConfig.SupportedNssaiAvailabilityData = n.SupportedNssaiAvailabilityData
		factory.NssfConfig.Configuration.AmfList = append(factory.NssfConfig.Configuration.AmfList,
			amfConfig)
	}

	// Return all authorized NSSAI availability information
	// a.AuthorizedNssaiAvailabilityData, _ = authorizeOfAmfFromConfig(nfId)

	// Return authorized NSSAI availability information of updated TAI only
	for _, s := range n.SupportedNssaiAvailabilityData {
		authorizedNssaiAvailabilityData, err := util.AuthorizeOfAmfTaFromConfig(nfId, *s.Tai)
		if err == nil {
			a.AuthorizedNssaiAvailabilityData = append(a.AuthorizedNssaiAvailabilityData, authorizedNssaiAvailabilityData)
		} else {
			logger.Nssaiavailability.Warnf(err.Error())
		}
	}

	return http.StatusOK
}
