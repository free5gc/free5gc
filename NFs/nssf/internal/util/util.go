/*
 * NSSF Utility
 */

package util

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/free5gc/nssf/internal/logger"
	"github.com/free5gc/nssf/pkg/factory"
	"github.com/free5gc/openapi"
	"github.com/free5gc/openapi/models"
)

// Title in Problem Details for NSSF HTTP APIs
const (
	INTERNAL_ERROR        = "Internal server error"
	INVALID_REQUEST       = "Invalid request message framing"
	MANDATORY_IE_MISSING  = "Mandatory IEs are missing"
	MALFORMED_REQUEST     = "Malformed request syntax"
	UNAUTHORIZED_CONSUMER = "Unauthorized NF service consumer"
	UNSUPPORTED_RESOURCE  = "Unsupported request resources"
)

// Check if a slice contains an element
func Contain(target interface{}, slice interface{}) bool {
	arr := reflect.ValueOf(slice)
	if arr.Kind() == reflect.Slice {
		for i := 0; i < arr.Len(); i++ {
			if reflect.DeepEqual(arr.Index(i).Interface(), target) {
				return true
			}
		}
	}
	return false
}

func SnssaiEqualFold(s models.ExtSnssai, t models.Snssai) bool {
	// TODO: Compare SdRanges and WildcardSd
	if s.Sst == t.Sst && strings.EqualFold(s.Sd, t.Sd) {
		return true
	}

	return false
}

// Check whether UE's Home PLMN is configured/supported
func CheckSupportedHplmn(homePlmnId models.PlmnId) bool {
	factory.NssfConfig.RLock()
	defer factory.NssfConfig.RUnlock()
	for _, mappingFromPlmn := range factory.NssfConfig.Configuration.MappingListFromPlmn {
		if *mappingFromPlmn.HomePlmnId == homePlmnId {
			return true
		}
	}
	logger.UtilLog.Warnf("No Home PLMN %+v in NSSF configuration", homePlmnId)
	return false
}

// Check whether UE's current TA is configured/supported
func CheckSupportedTa(tai models.Tai) bool {
	factory.NssfConfig.RLock()
	defer factory.NssfConfig.RUnlock()
	for _, taConfig := range factory.NssfConfig.Configuration.TaList {
		if reflect.DeepEqual(*taConfig.Tai, tai) {
			return true
		}
	}
	e, err := json.Marshal(tai)
	if err != nil {
		logger.UtilLog.Errorf("Marshal error in CheckSupportedTa: %+v", err)
	}
	logger.UtilLog.Warnf("No TA %s in NSSF configuration", e)
	return false
}

// Check whether the given S-NSSAI is supported or not in PLMN
func CheckSupportedSnssaiInPlmn(snssai models.Snssai, plmnId models.PlmnId) bool {
	factory.NssfConfig.RLock()
	defer factory.NssfConfig.RUnlock()
	if CheckStandardSnssai(snssai) {
		return true
	}

	for _, supportedNssaiInPlmn := range factory.NssfConfig.Configuration.SupportedNssaiInPlmnList {
		if *supportedNssaiInPlmn.PlmnId == plmnId {
			for _, supportedSnssai := range supportedNssaiInPlmn.SupportedSnssaiList {
				if openapi.SnssaiEqualFold(snssai, supportedSnssai) {
					return true
				}
			}
			return false
		}
	}
	logger.UtilLog.Warnf("No supported S-NSSAI list of PLMNID %+v in NSSF configuration", plmnId)
	return false
}

// Check whether S-NSSAIs in NSSAI are supported or not in PLMN
func CheckSupportedNssaiInPlmn(nssai any, plmnId models.PlmnId) bool {
	factory.NssfConfig.RLock()
	defer factory.NssfConfig.RUnlock()
	for _, supportedNssaiInPlmn := range factory.NssfConfig.Configuration.SupportedNssaiInPlmnList {
		if *supportedNssaiInPlmn.PlmnId == plmnId {
			if n, ok := nssai.([]models.ExtSnssai); ok {
				for _, snssai := range n {
					// Standard S-NSSAIs are supposed to be supported
					// If not, disable following check and be sure to add supported standard S-NSSAI(s) in configuration
					if CheckStandardSnssai(models.Snssai{Sst: snssai.Sst, Sd: snssai.Sd}) {
						continue
					}

					hitSupportedNssai := false
					for _, supportedSnssai := range supportedNssaiInPlmn.SupportedSnssaiList {
						if SnssaiEqualFold(snssai, supportedSnssai) {
							hitSupportedNssai = true
							break
						}
					}

					if !hitSupportedNssai {
						return false
					}
				}
				return true
			} else if n, ok := nssai.([]models.Snssai); ok {
				for _, snssai := range n {
					if CheckStandardSnssai(snssai) {
						continue
					}

					hitSupportedNssai := false
					for _, supportedSnssai := range supportedNssaiInPlmn.SupportedSnssaiList {
						if openapi.SnssaiEqualFold(snssai, supportedSnssai) {
							hitSupportedNssai = true
							break
						}
					}

					if !hitSupportedNssai {
						return false
					}
				}
				return true
			} else {
				logger.UtilLog.Warnf("Unsupported type of NSSAI: %+v", nssai)
				return false
			}
		}
	}
	logger.UtilLog.Warnf("No supported S-NSSAI list of PLMNID %+v in NSSF configuration", plmnId)
	return false
}

// Check whether S-NSSAI is supported or not at UE's current TA
func CheckSupportedSnssaiInTa(snssai models.Snssai, tai models.Tai) bool {
	factory.NssfConfig.RLock()
	defer factory.NssfConfig.RUnlock()
	for _, taConfig := range factory.NssfConfig.Configuration.TaList {
		if reflect.DeepEqual(*taConfig.Tai, tai) {
			for _, supportedSnssai := range taConfig.SupportedSnssaiList {
				if SnssaiEqualFold(supportedSnssai, snssai) {
					return true
				}
			}
			return false
		}
	}
	return false

	// // Check supported S-NSSAI in AmfList instead of TaList
	// for _, amfConfig := range factory.NssfConfig.Configuration.AmfList {
	//     if checkSupportedNssaiAvailabilityData(snssai, tai, amfConfig.SupportedNssaiAvailabilityData) == true {
	//         return true
	//     }
	// }
	// return false
}

// Check whether S-NSSAI is in SupportedNssaiAvailabilityData under the given TAI
func CheckSupportedNssaiAvailabilityData(
	snssai models.Snssai, tai models.Tai, s []models.SupportedNssaiAvailabilityData,
) bool {
	for _, supportedNssaiAvailabilityData := range s {
		if reflect.DeepEqual(*supportedNssaiAvailabilityData.Tai, tai) &&
			CheckSnssaiInNssai(snssai, supportedNssaiAvailabilityData.SupportedSnssaiList) {
			return true
		}
	}
	return false
}

// Check whether S-NSSAI is supported or not by the AMF at UE's current TA
func CheckSupportedSnssaiInAmfTa(snssai models.Snssai, nfId string, tai models.Tai) bool {
	// Uncomment following lines if supported S-NSSAI lists of AMF Sets are independent of those of AMFs
	// for _, amfSetConfig := range factory.NssfConfig.Configuration.AmfSetList {
	//     if amfSetConfig.AmfList != nil && len(amfSetConfig.AmfList) != 0 && Contain(nfId, amfSetConfig.AmfList) {
	//         return checkSupportedNssaiAvailabilityData(snssai, tai, amfSetConfig.SupportedNssaiAvailabilityData)
	//     }
	// }

	for _, amfConfig := range factory.NssfConfig.Configuration.AmfList {
		if amfConfig.NfId == nfId {
			return CheckSupportedNssaiAvailabilityData(snssai, tai, amfConfig.SupportedNssaiAvailabilityData)
		}
	}

	logger.UtilLog.Warnf("No AMF %s in NSSF configuration", nfId)
	return false
}

// Check whether all S-NSSAIs in Allowed NSSAI is supported by the AMF at UE's current TA
func CheckAllowedNssaiInAmfTa(allowedNssaiList []models.AllowedNssai, nfId string, tai models.Tai) bool {
	for _, allowedNssai := range allowedNssaiList {
		for _, allowedSnssai := range allowedNssai.AllowedSnssaiList {
			if CheckSupportedSnssaiInAmfTa(*allowedSnssai.AllowedSnssai, nfId, tai) {
				continue
			} else {
				return false
			}
		}
	}
	return true
}

// Check whether S-NSSAI is standard or non-standard value
// A standard S-NSSAI is only comprised of a standardized SST value and no SD
func CheckStandardSnssai(snssai models.Snssai) bool {
	if snssai.Sst >= 1 && snssai.Sst <= 3 && snssai.Sd == "" {
		return true
	}
	return false
}

// Check whether the NSSAI contains the specific S-NSSAI
func CheckSnssaiInNssai(targetSnssai models.Snssai, nssai []models.ExtSnssai) bool {
	for _, snssai := range nssai {
		if SnssaiEqualFold(snssai, targetSnssai) {
			return true
		}
	}
	return false
}

// Get S-NSSAI mappings of the given Home PLMN ID from configuration
func GetMappingOfPlmnFromConfig(homePlmnId models.PlmnId) []models.MappingOfSnssai {
	factory.NssfConfig.RLock()
	defer factory.NssfConfig.RUnlock()
	for _, mappingFromPlmn := range factory.NssfConfig.Configuration.MappingListFromPlmn {
		if *mappingFromPlmn.HomePlmnId == homePlmnId {
			return mappingFromPlmn.MappingOfSnssai
		}
	}
	return nil
}

// Get NSI information list of the given S-NSSAI from configuration
func GetNsiInformationListFromConfig(snssai models.Snssai) []models.NsiInformation {
	factory.NssfConfig.RLock()
	defer factory.NssfConfig.RUnlock()
	for _, nsiConfig := range factory.NssfConfig.Configuration.NsiList {
		if openapi.SnssaiEqualFold(*nsiConfig.Snssai, snssai) {
			return nsiConfig.NsiInformationList
		}
	}
	return nil
}

// Get Access Type of the given TAI from configuraion
func GetAccessTypeFromConfig(tai models.Tai) models.AccessType {
	factory.NssfConfig.RLock()
	defer factory.NssfConfig.RUnlock()
	for _, taConfig := range factory.NssfConfig.Configuration.TaList {
		if reflect.DeepEqual(*taConfig.Tai, tai) {
			return *taConfig.AccessType
		}
	}
	e, err := json.Marshal(tai)
	if err != nil {
		logger.UtilLog.Errorf("Marshal error in GetAccessTypeFromConfig: %+v", err)
	}
	logger.UtilLog.Warnf("No TA %s in NSSF configuration", e)
	return models.AccessType__3_GPP_ACCESS
}

// Get restricted S-NSSAI list of the given TAI from configuration
func GetRestrictedSnssaiListFromConfig(tai models.Tai) []models.RestrictedSnssai {
	factory.NssfConfig.RLock()
	defer factory.NssfConfig.RUnlock()
	for _, taConfig := range factory.NssfConfig.Configuration.TaList {
		if reflect.DeepEqual(*taConfig.Tai, tai) {
			if len(taConfig.RestrictedSnssaiList) != 0 {
				return taConfig.RestrictedSnssaiList
			} else {
				return nil
			}
		}
	}
	e, err := json.Marshal(tai)
	if err != nil {
		logger.UtilLog.Errorf("Marshal error in GetRestrictedSnssaiListFromConfig: %+v", err)
	}
	logger.UtilLog.Warnf("No TA %s in NSSF configuration", e)
	return nil
}

// Get authorized NSSAI availability data of the given NF ID and TAI from configuration
func AuthorizeOfAmfTaFromConfig(nfId string, tai models.Tai) (models.AuthorizedNssaiAvailabilityData, error) {
	var authorizedNssaiAvailabilityData models.AuthorizedNssaiAvailabilityData
	authorizedNssaiAvailabilityData.Tai = new(models.Tai)
	*authorizedNssaiAvailabilityData.Tai = tai

	for _, amfConfig := range factory.NssfConfig.Configuration.AmfList {
		if amfConfig.NfId == nfId {
			for _, supportedNssaiAvailabilityData := range amfConfig.SupportedNssaiAvailabilityData {
				if reflect.DeepEqual(*supportedNssaiAvailabilityData.Tai, tai) {
					authorizedNssaiAvailabilityData.SupportedSnssaiList = supportedNssaiAvailabilityData.SupportedSnssaiList
					authorizedNssaiAvailabilityData.RestrictedSnssaiList = GetRestrictedSnssaiListFromConfig(tai)

					// TODO: Sort the returned slice
					return authorizedNssaiAvailabilityData, nil
				}
			}
			e, err1 := json.Marshal(tai)
			if err1 != nil {
				logger.UtilLog.Errorf("Marshal error in AuthorizeOfAmfTaFromConfig: %+v", err1)
			}
			err := fmt.Errorf("no supported S-NSSAI list by AMF %s under TAI %s in NSSF configuration", nfId, e)
			return authorizedNssaiAvailabilityData, err
		}
	}
	err := fmt.Errorf("no AMF configuration of %s", nfId)
	return authorizedNssaiAvailabilityData, err
}

// Get all authorized NSSAI availability data of the given NF ID from configuration
func AuthorizeOfAmfFromConfig(nfId string) ([]models.AuthorizedNssaiAvailabilityData, error) {
	var authorizedNssaiAvailabilityDataList []models.AuthorizedNssaiAvailabilityData

	factory.NssfConfig.RLock()
	defer factory.NssfConfig.RUnlock()
	for _, amfConfig := range factory.NssfConfig.Configuration.AmfList {
		if amfConfig.NfId == nfId {
			for _, supportedNssaiAvailabilityData := range amfConfig.SupportedNssaiAvailabilityData {
				var authorizedNssaiAvailabilityData models.AuthorizedNssaiAvailabilityData
				authorizedNssaiAvailabilityData.Tai = new(models.Tai)
				*authorizedNssaiAvailabilityData.Tai = *supportedNssaiAvailabilityData.Tai
				authorizedNssaiAvailabilityData.SupportedSnssaiList = supportedNssaiAvailabilityData.SupportedSnssaiList
				authorizedNssaiAvailabilityData.RestrictedSnssaiList = GetRestrictedSnssaiListFromConfig(
					*authorizedNssaiAvailabilityData.Tai)

				authorizedNssaiAvailabilityDataList = append(
					authorizedNssaiAvailabilityDataList,
					authorizedNssaiAvailabilityData)
			}
			return authorizedNssaiAvailabilityDataList, nil
		}
	}
	err := fmt.Errorf("no AMF configuration of %s", nfId)
	return authorizedNssaiAvailabilityDataList, err
}

// Get authorized NSSAI availability data of the given TAI list from configuration
func AuthorizeOfTaListFromConfig(taiList []models.Tai) []models.AuthorizedNssaiAvailabilityData {
	var authorizedNssaiAvailabilityDataList []models.AuthorizedNssaiAvailabilityData

	for _, taConfig := range factory.NssfConfig.Configuration.TaList {
		for _, tai := range taiList {
			if reflect.DeepEqual(*taConfig.Tai, tai) {
				var authorizedNssaiAvailabilityData models.AuthorizedNssaiAvailabilityData
				authorizedNssaiAvailabilityData.Tai = new(models.Tai)
				*authorizedNssaiAvailabilityData.Tai = tai
				authorizedNssaiAvailabilityData.SupportedSnssaiList = taConfig.SupportedSnssaiList
				authorizedNssaiAvailabilityData.RestrictedSnssaiList = GetRestrictedSnssaiListFromConfig(tai)

				authorizedNssaiAvailabilityDataList = append(authorizedNssaiAvailabilityDataList, authorizedNssaiAvailabilityData)
			}
		}
	}
	return authorizedNssaiAvailabilityDataList
}

// Get supported S-NSSAI list of the given NF ID and TAI from configuration
func GetSupportedSnssaiListFromConfig(nfId string, tai models.Tai) []models.ExtSnssai {
	for _, amfConfig := range factory.NssfConfig.Configuration.AmfList {
		if amfConfig.NfId == nfId {
			for _, supportedNssaiAvailabilityData := range amfConfig.SupportedNssaiAvailabilityData {
				if reflect.DeepEqual(*supportedNssaiAvailabilityData.Tai, tai) {
					return supportedNssaiAvailabilityData.SupportedSnssaiList
				}
			}
			return nil
		}
	}
	return nil
}

// Find target S-NSSAI mapping with serving S-NSSAIs from mapping of S-NSSAI(s)
func FindMappingWithServingSnssai(
	snssai models.Snssai, mappings []models.MappingOfSnssai,
) (models.MappingOfSnssai, bool) {
	for _, mapping := range mappings {
		if openapi.SnssaiEqualFold(*mapping.ServingSnssai, snssai) {
			return mapping, true
		}
	}
	return models.MappingOfSnssai{}, false
}

// Find target S-NSSAI mapping with home S-NSSAIs from mapping of S-NSSAI(s)
func FindMappingWithHomeSnssai(snssai models.Snssai, mappings []models.MappingOfSnssai) (models.MappingOfSnssai, bool) {
	for _, mapping := range mappings {
		if openapi.SnssaiEqualFold(*mapping.HomeSnssai, snssai) {
			return mapping, true
		}
	}
	return models.MappingOfSnssai{}, false
}

// Add Allowed S-NSSAI to Authorized Network Slice Info
func AddAllowedSnssai(allowedSnssai models.AllowedSnssai, accessType models.AccessType,
	authorizedNetworkSliceInfo *models.AuthorizedNetworkSliceInfo,
) {
	hitAllowedNssai := false
	for i := range authorizedNetworkSliceInfo.AllowedNssaiList {
		if authorizedNetworkSliceInfo.AllowedNssaiList[i].AccessType == accessType {
			hitAllowedNssai = true
			const MAX_ALLOWED_SNSSAI = 8
			if len(authorizedNetworkSliceInfo.AllowedNssaiList[i].AllowedSnssaiList) == MAX_ALLOWED_SNSSAI {
				logger.UtilLog.Infof("Unable to add a new Allowed S-NSSAI since already eight S-NSSAIs in Allowed NSSAI")
			} else {
				authorizedNetworkSliceInfo.AllowedNssaiList[i].AllowedSnssaiList = append(
					authorizedNetworkSliceInfo.AllowedNssaiList[i].AllowedSnssaiList,
					allowedSnssai)
			}
			break
		}
	}

	if !hitAllowedNssai {
		var allowedNssaiElement models.AllowedNssai
		allowedNssaiElement.AllowedSnssaiList = append(allowedNssaiElement.AllowedSnssaiList, allowedSnssai)
		allowedNssaiElement.AccessType = accessType

		authorizedNetworkSliceInfo.AllowedNssaiList = append(authorizedNetworkSliceInfo.AllowedNssaiList, allowedNssaiElement)
	}
}

// Add AMF information to Authorized Network Slice Info
func AddAmfInformation(tai models.Tai, authorizedNetworkSliceInfo *models.AuthorizedNetworkSliceInfo) {
	factory.NssfConfig.RLock()
	defer factory.NssfConfig.RUnlock()
	if len(authorizedNetworkSliceInfo.AllowedNssaiList) == 0 {
		return
	}

	// Check if any AMF can serve the UE
	// That is, whether NSSAI of all Allowed S-NSSAIs is a subset of NSSAI supported by AMF

	// Find AMF Set that could serve UE from AMF Set list in configuration
	// Simply use the first applicable AMF set
	// TODO: Policies of AMF selection (e.g. load balance between AMF instances)
	for _, amfSetConfig := range factory.NssfConfig.Configuration.AmfSetList {
		hitAllowedNssai := true
		for _, allowedNssai := range authorizedNetworkSliceInfo.AllowedNssaiList {
			for _, allowedSnssai := range allowedNssai.AllowedSnssaiList {
				if CheckSupportedNssaiAvailabilityData(*allowedSnssai.AllowedSnssai,
					tai, amfSetConfig.SupportedNssaiAvailabilityData) {
					continue
				} else {
					hitAllowedNssai = false
					break
				}
			}
			if !hitAllowedNssai {
				break
			}
		}

		if !hitAllowedNssai {
			continue
		} else {
			// Add AMF Set to Authorized Network Slice Info
			if len(amfSetConfig.AmfList) != 0 {
				// List of candidate AMF(s) provided in configuration
				authorizedNetworkSliceInfo.CandidateAmfList = append(
					authorizedNetworkSliceInfo.CandidateAmfList,
					amfSetConfig.AmfList...)
			} else {
				// TODO: Possibly querying the NRF
				authorizedNetworkSliceInfo.TargetAmfSet = amfSetConfig.AmfSetId
				// The API URI of the NRF may be included if target AMF Set is included
				authorizedNetworkSliceInfo.NrfAmfSet = amfSetConfig.NrfAmfSet
			}
			return
		}
	}

	// No AMF Set in configuration can serve the UE
	// Find all candidate AMFs that could serve UE from AMF list in configuration
	hitAmf := false
	for _, amfConfig := range factory.NssfConfig.Configuration.AmfList {
		hitAllowedNssai := true
		for _, allowedNssai := range authorizedNetworkSliceInfo.AllowedNssaiList {
			for _, allowedSnssai := range allowedNssai.AllowedSnssaiList {
				if CheckSupportedNssaiAvailabilityData(*allowedSnssai.AllowedSnssai,
					tai, amfConfig.SupportedNssaiAvailabilityData) {
					continue
				} else {
					hitAllowedNssai = false
					break
				}
			}
			if !hitAllowedNssai {
				break
			}
		}

		if !hitAllowedNssai {
			continue
		} else {
			// Add AMF Set to Authorized Network Slice Info
			authorizedNetworkSliceInfo.CandidateAmfList = append(authorizedNetworkSliceInfo.CandidateAmfList, amfConfig.NfId)
			hitAmf = true
		}
	}

	if !hitAmf {
		logger.UtilLog.Warnf("No candidate AMF or AMF Set can serve the UE")
	}
}
