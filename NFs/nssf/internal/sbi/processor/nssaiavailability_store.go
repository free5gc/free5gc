/*
 * NSSF NSSAI Availability
 *
 * NSSF NSSAI Availability Service
 */

package processor

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"

	jsonpatch "github.com/evanphx/json-patch"
	"github.com/gin-gonic/gin"

	"github.com/free5gc/nssf/internal/logger"
	"github.com/free5gc/nssf/internal/plugin"
	"github.com/free5gc/nssf/internal/util"
	"github.com/free5gc/nssf/pkg/factory"
	"github.com/free5gc/openapi/models"
)

func (p *Processor) NssaiAvailabilityNfInstanceDelete(c *gin.Context, nfId string) {
	var problemDetails *models.ProblemDetails
	for i, amfConfig := range factory.NssfConfig.Configuration.AmfList {
		if amfConfig.NfId == nfId {
			factory.NssfConfig.Configuration.AmfList = append(
				factory.NssfConfig.Configuration.AmfList[:i],
				factory.NssfConfig.Configuration.AmfList[i+1:]...)

			c.Status(http.StatusNoContent)
			return
		}
	}

	problemDetails = &models.ProblemDetails{
		Title:  util.UNSUPPORTED_RESOURCE,
		Status: http.StatusNotFound,
		Detail: fmt.Sprintf("AMF ID '%s' does not exist", nfId),
	}
	util.GinProblemJson(c, problemDetails)
}

func (p *Processor) NssaiAvailabilityNfInstancePatch(
	c *gin.Context,
	nssaiAvailabilityUpdateInfo plugin.PatchDocument, nfId string,
) {
	var (
		response       = &models.AuthorizedNssaiAvailabilityInfo{}
		problemDetails *models.ProblemDetails
	)

	var amfIdx int
	var original []byte
	hitAmf := false
	factory.NssfConfig.RLock()
	for amfIdx, amfConfig := range factory.NssfConfig.Configuration.AmfList {
		if amfConfig.NfId == nfId {
			// Since json-patch package does not have idea of optional field of datatype,
			// provide with null or empty value instead of omitting the field
			temp := factory.NssfConfig.Configuration.AmfList[amfIdx].SupportedNssaiAvailabilityData
			const dummyString string = "DUMMY"
			for i := range temp {
				for j := range temp[i].SupportedSnssaiList {
					if temp[i].SupportedSnssaiList[j].Sd == "" {
						temp[i].SupportedSnssaiList[j].Sd = dummyString
					}
				}
			}
			var err error
			original, err = json.Marshal(temp)
			if err != nil {
				logger.NssaiavailLog.Errorf("Marshal error in NSSAIAvailabilityPatchProcedure: %+v", err)
			}
			original = bytes.ReplaceAll(original, []byte(dummyString), []byte(""))

			// original, _ = json.Marshal(factory.NssfConfig.Configuration.AmfList[amfIdx].SupportedNssaiAvailabilityData)

			hitAmf = true
			break
		}
	}
	factory.NssfConfig.RUnlock()
	if !hitAmf {
		problemDetails = &models.ProblemDetails{
			Title:  util.UNSUPPORTED_RESOURCE,
			Status: http.StatusNotFound,
			Detail: fmt.Sprintf("AMF ID '%s' does not exist", nfId),
		}
		util.GinProblemJson(c, problemDetails)
		return
	}

	// TODO: Check if returned HTTP status codes or problem details are proper when errors occur

	// Provide JSON string with null or empty value in `Value` of `PatchItem`
	for i, patchItem := range nssaiAvailabilityUpdateInfo {
		if reflect.ValueOf(patchItem.Value).Kind() == reflect.Map {
			_, exist := patchItem.Value.(map[string]interface{})["sst"]
			_, notExist := patchItem.Value.(map[string]interface{})["sd"]
			if exist && !notExist {
				nssaiAvailabilityUpdateInfo[i].Value.(map[string]interface{})["sd"] = ""
			}
		}
	}
	patchJSON, err := json.Marshal(nssaiAvailabilityUpdateInfo)
	if err != nil {
		logger.NssaiavailLog.Errorf("Marshal error in NSSAIAvailabilityPatchProcedure: %+v", err)
	}

	patch, err := jsonpatch.DecodePatch(patchJSON)
	if err != nil {
		problemDetails = &models.ProblemDetails{
			Title:  util.MALFORMED_REQUEST,
			Status: http.StatusBadRequest,
			Detail: err.Error(),
		}
		util.GinProblemJson(c, problemDetails)
		return
	}

	modified, err := patch.Apply(original)
	if err != nil {
		problemDetails = &models.ProblemDetails{
			Title:  util.INVALID_REQUEST,
			Status: http.StatusConflict,
			Detail: err.Error(),
		}
		util.GinProblemJson(c, problemDetails)
		return
	}

	factory.NssfConfig.Lock()
	err = json.Unmarshal(modified, &factory.NssfConfig.Configuration.AmfList[amfIdx].SupportedNssaiAvailabilityData)
	factory.NssfConfig.Unlock()
	if err != nil {
		problemDetails = &models.ProblemDetails{
			Title:  util.INVALID_REQUEST,
			Status: http.StatusBadRequest,
			Detail: err.Error(),
		}
		util.GinProblemJson(c, problemDetails)
		return
	}

	// Return all authorized NSSAI availability information
	response.AuthorizedNssaiAvailabilityData, err = util.AuthorizeOfAmfFromConfig(nfId)
	if err != nil {
		logger.NssaiavailLog.Errorf("util AuthorizeOfAmfFromConfig error in NSSAIAvailabilityPatchProcedure: %+v", err)
	}

	// TODO: Return authorized NSSAI availability information of updated TAI only

	c.JSON(http.StatusOK, response)
}

// NSSAIAvailability PUT method
func (p *Processor) NssaiAvailabilityNfInstanceUpdate(
	c *gin.Context,
	nssaiAvailabilityInfo models.NssaiAvailabilityInfo, nfId string,
) {
	var (
		response       = &models.AuthorizedNssaiAvailabilityInfo{}
		problemDetails *models.ProblemDetails
	)

	for _, s := range nssaiAvailabilityInfo.SupportedNssaiAvailabilityData {
		if !util.CheckSupportedNssaiInPlmn(s.SupportedSnssaiList, *s.Tai.PlmnId) {
			problemDetails = &models.ProblemDetails{
				Title:  util.UNSUPPORTED_RESOURCE,
				Status: http.StatusForbidden,
				Detail: "S-NSSAI in Requested NSSAI is not supported in PLMN",
				Cause:  "SNSSAI_NOT_SUPPORTED",
			}

			util.GinProblemJson(c, problemDetails)
			return
		}
	}

	// TODO: Currently authorize all the provided S-NSSAIs
	//       Take some issue into consideration e.g. operator policies

	hitAmf := false
	// Find AMF configuration of given NfId
	// If found, then update the SupportedNssaiAvailabilityData
	factory.NssfConfig.Lock()
	for i, amfConfig := range factory.NssfConfig.Configuration.AmfList {
		if amfConfig.NfId == nfId {
			factory.NssfConfig.Configuration.AmfList[i].SupportedNssaiAvailabilityData = nssaiAvailabilityInfo.
				SupportedNssaiAvailabilityData

			hitAmf = true
			break
		}
	}
	factory.NssfConfig.Unlock()

	// If no AMF record is found, create a new one
	if !hitAmf {
		var amfConfig factory.AmfConfig
		amfConfig.NfId = nfId
		amfConfig.SupportedNssaiAvailabilityData = nssaiAvailabilityInfo.SupportedNssaiAvailabilityData
		factory.NssfConfig.Lock()
		factory.NssfConfig.Configuration.AmfList = append(
			factory.NssfConfig.Configuration.AmfList,
			amfConfig)
		factory.NssfConfig.Unlock()
	}

	// Return all authorized NSSAI availability information
	// a.AuthorizedNssaiAvailabilityData, _ = authorizeOfAmfFromConfig(nfId)

	// Return authorized NSSAI availability information of updated TAI only
	for _, s := range nssaiAvailabilityInfo.SupportedNssaiAvailabilityData {
		authorizedNssaiAvailabilityData, err := util.AuthorizeOfAmfTaFromConfig(nfId, *s.Tai)
		if err == nil {
			response.AuthorizedNssaiAvailabilityData = append(
				response.AuthorizedNssaiAvailabilityData,
				authorizedNssaiAvailabilityData)
		} else {
			logger.NssaiavailLog.Warn(err)
		}
	}

	c.JSON(http.StatusOK, response)
}
