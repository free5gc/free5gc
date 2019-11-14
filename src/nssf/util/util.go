/*
 * NSSF Utility
 */

package util

import (
	"encoding/json"
	"fmt"
	"reflect"

	. "free5gc/lib/openapi/models"
	"free5gc/lib/path_util"
	"free5gc/src/nssf/factory"
	"free5gc/src/nssf/logger"
)

// Title in Problem Details for NSSF HTTP APIs
const (
	INVALID_REQUEST       = "Invalid request message framing"
	MALFORMED_REQUEST     = "Malformed request syntax"
	UNAUTHORIZED_CONSUMER = "Unauthorized NF service consumer"
	UNSUPPORTED_RESOURCE  = "Unsupported request resources"
)

// Path of HTTP2 key and log file
var (
	NSSF_LOG_PATH = path_util.Gofree5gcPath("free5gc/nssfsslkey.log")
	NSSF_PEM_PATH = path_util.Gofree5gcPath("free5gc/support/TLS/nssf.pem")
	NSSF_KEY_PATH = path_util.Gofree5gcPath("free5gc/support/TLS/nssf.key")
)

// Default configuration file
var DEFAULT_CONFIG string = "free5gc/config/nssfcfg.conf"

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

// Check whether UE's Home PLMN is configured/supported
func CheckSupportedHplmn(homePlmnId PlmnId) bool {
	for _, mappingFromPlmn := range factory.NssfConfig.Configuration.MappingListFromPlmn {
		if *mappingFromPlmn.HomePlmnId == homePlmnId {
			return true
		}
	}
	logger.Util.Warnf("No Home PLMN %+v in NSSF configuration", homePlmnId)
	return false
}

// Check whether UE's current TA is configured/supported
func CheckSupportedTa(tai Tai) bool {
	for _, taConfig := range factory.NssfConfig.Configuration.TaList {
		if reflect.DeepEqual(*taConfig.Tai, tai) {
			return true
		}
	}
	e, _ := json.Marshal(tai)
	logger.Util.Warnf("No TA %s in NSSF configuration", e)
	return false
}

// Check whether the given S-NSSAI is supported or not in PLMN
func CheckSupportedSnssaiInPlmn(snssai Snssai, plmnId PlmnId) bool {
	if CheckStandardSnssai(snssai) {
		return true
	}

	for _, supportedNssaiInPlmn := range factory.NssfConfig.Configuration.SupportedNssaiInPlmnList {
		if *supportedNssaiInPlmn.PlmnId == plmnId {
			for _, supportedSnssai := range supportedNssaiInPlmn.SupportedSnssaiList {
				if snssai == supportedSnssai {
					return true
				}
			}
			return false
		}
	}
	logger.Util.Warnf("No supported S-NSSAI list of PLMNID %+v in NSSF configuration", plmnId)
	return false
}

// Check whether S-NSSAIs in NSSAI are supported or not in PLMN
func CheckSupportedNssaiInPlmn(nssai []Snssai, plmnId PlmnId) bool {
	for _, supportedNssaiInPlmn := range factory.NssfConfig.Configuration.SupportedNssaiInPlmnList {
		if *supportedNssaiInPlmn.PlmnId == plmnId {
			for _, snssai := range nssai {
				// Standard S-NSSAIs are supposed to be supported
				// If not, disable following check and be sure to add supported standard S-NSSAI(s) in configuration
				if CheckStandardSnssai(snssai) {
					continue
				}

				hitSupportedNssai := false
				for _, supportedSnssai := range supportedNssaiInPlmn.SupportedSnssaiList {
					if snssai == supportedSnssai {
						hitSupportedNssai = true
						break
					}
				}

				if !hitSupportedNssai {
					return false
				}
			}
			return true
		}
	}
	logger.Util.Warnf("No supported S-NSSAI list of PLMNID %+v in NSSF configuration", plmnId)
	return false
}

// Check whether S-NSSAI is supported or not at UE's current TA
func CheckSupportedSnssaiInTa(snssai Snssai, tai Tai) bool {
	for _, taConfig := range factory.NssfConfig.Configuration.TaList {
		if reflect.DeepEqual(*taConfig.Tai, tai) {
			for _, supportedSnssai := range taConfig.SupportedSnssaiList {
				if supportedSnssai == snssai {
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
func CheckSupportedNssaiAvailabilityData(snssai Snssai, tai Tai, s []SupportedNssaiAvailabilityData) bool {
	for _, supportedNssaiAvailabilityData := range s {
		if reflect.DeepEqual(*supportedNssaiAvailabilityData.Tai, tai) &&
			CheckSnssaiInNssai(snssai, supportedNssaiAvailabilityData.SupportedSnssaiList) {
			return true
		}
	}
	return false
}

// Check whether S-NSSAI is supported or not by the AMF at UE's current TA
func CheckSupportedSnssaiInAmfTa(snssai Snssai, nfId string, tai Tai) bool {
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

	logger.Util.Warnf("No AMF %s in NSSF configuration", nfId)
	return false
}

// Check whether all S-NSSAIs in Allowed NSSAI is supported by the AMF at UE's current TA
func CheckAllowedNssaiInAmfTa(allowedNssaiList []AllowedNssai, nfId string, tai Tai) bool {
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
func CheckStandardSnssai(snssai Snssai) bool {
	if snssai.Sst >= 1 && snssai.Sst <= 3 && snssai.Sd == "" {
		return true
	}
	return false
}

// Check whether the NSSAI contains the specific S-NSSAI
func CheckSnssaiInNssai(targetSnssai Snssai, nssai []Snssai) bool {
	for _, snssai := range nssai {
		if snssai == targetSnssai {
			return true
		}
	}
	return false
}

// Get S-NSSAI mappings of the given Home PLMN ID from configuration
func GetMappingOfPlmnFromConfig(homePlmnId PlmnId) []MappingOfSnssai {
	for _, mappingFromPlmn := range factory.NssfConfig.Configuration.MappingListFromPlmn {
		if *mappingFromPlmn.HomePlmnId == homePlmnId {
			return mappingFromPlmn.MappingOfSnssai
		}
	}
	return nil
}

// Get NSI information list of the given S-NSSAI from configuration
func GetNsiInformationListFromConfig(snssai Snssai) []NsiInformation {
	for _, nsiConfig := range factory.NssfConfig.Configuration.NsiList {
		if *nsiConfig.Snssai == snssai {
			return nsiConfig.NsiInformationList
		}
	}
	return nil
}

// Get Access Type of the given TAI from configuraion
func GetAccessTypeFromConfig(tai Tai) AccessType {
	for _, taConfig := range factory.NssfConfig.Configuration.TaList {
		if reflect.DeepEqual(*taConfig.Tai, tai) {
			return *taConfig.AccessType
		}
	}
	e, _ := json.Marshal(tai)
	logger.Util.Warnf("No TA %s in NSSF configuration", e)
	return AccessType__3_GPP_ACCESS
}

// Get restricted S-NSSAI list of the given TAI from configuration
func GetRestrictedSnssaiListFromConfig(tai Tai) []RestrictedSnssai {
	for _, taConfig := range factory.NssfConfig.Configuration.TaList {
		if reflect.DeepEqual(*taConfig.Tai, tai) {
			if taConfig.RestrictedSnssaiList != nil && len(taConfig.RestrictedSnssaiList) != 0 {
				return taConfig.RestrictedSnssaiList
			} else {
				return nil
			}
		}
	}
	e, _ := json.Marshal(tai)
	logger.Util.Warnf("No TA %s in NSSF configuration", e)
	return nil
}

// Get authorized NSSAI availability data of the given NF ID and TAI from configuration
func AuthorizeOfAmfTaFromConfig(nfId string, tai Tai) (AuthorizedNssaiAvailabilityData, error) {
	var a AuthorizedNssaiAvailabilityData
	a.Tai = new(Tai)
	*a.Tai = tai

	for _, amfConfig := range factory.NssfConfig.Configuration.AmfList {
		if amfConfig.NfId == nfId {
			for _, supportedNssaiAvailabilityData := range amfConfig.SupportedNssaiAvailabilityData {
				if reflect.DeepEqual(*supportedNssaiAvailabilityData.Tai, tai) {
					a.SupportedSnssaiList = supportedNssaiAvailabilityData.SupportedSnssaiList
					a.RestrictedSnssaiList = GetRestrictedSnssaiListFromConfig(tai)

					// TODO: Sort the returned slice
					return a, nil
				}
			}
			e, _ := json.Marshal(tai)
			err := fmt.Errorf("No supported S-NSSAI list by AMF %s under TAI %s in NSSF configuration", nfId, e)
			return a, err
		}
	}
	err := fmt.Errorf("No AMF configuration of %s", nfId)
	return a, err
}

// Get all authorized NSSAI availability data of the given NF ID from configuration
func AuthorizeOfAmfFromConfig(nfId string) ([]AuthorizedNssaiAvailabilityData, error) {
	var s []AuthorizedNssaiAvailabilityData

	for _, amfConfig := range factory.NssfConfig.Configuration.AmfList {
		if amfConfig.NfId == nfId {
			for _, supportedNssaiAvailabilityData := range amfConfig.SupportedNssaiAvailabilityData {
				var a AuthorizedNssaiAvailabilityData
				a.Tai = new(Tai)
				*a.Tai = *supportedNssaiAvailabilityData.Tai
				a.SupportedSnssaiList = supportedNssaiAvailabilityData.SupportedSnssaiList
				a.RestrictedSnssaiList = GetRestrictedSnssaiListFromConfig(*a.Tai)

				s = append(s, a)
			}
			return s, nil
		}
	}
	err := fmt.Errorf("No AMF configuration of %s", nfId)
	return s, err
}

// Get authorized NSSAI availability data of the given TAI list from configuration
func AuthorizeOfTaListFromConfig(taiList []Tai) []AuthorizedNssaiAvailabilityData {
	var s []AuthorizedNssaiAvailabilityData

	for _, taConfig := range factory.NssfConfig.Configuration.TaList {
		for _, tai := range taiList {
			if reflect.DeepEqual(*taConfig.Tai, tai) {
				var a AuthorizedNssaiAvailabilityData
				a.Tai = new(Tai)
				*a.Tai = tai
				a.SupportedSnssaiList = taConfig.SupportedSnssaiList
				a.RestrictedSnssaiList = GetRestrictedSnssaiListFromConfig(tai)

				s = append(s, a)
			}
		}
	}
	return s
}

// Get supported S-NSSAI list of the given NF ID and TAI from configuration
func GetSupportedSnssaiListFromConfig(nfId string, tai Tai) []Snssai {
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
func FindMappingWithServingSnssai(snssai Snssai, mappings []MappingOfSnssai) (MappingOfSnssai, bool) {
	for _, mapping := range mappings {
		if *mapping.ServingSnssai == snssai {
			return mapping, true
		}
	}
	return MappingOfSnssai{}, false
}

// Find target S-NSSAI mapping with home S-NSSAIs from mapping of S-NSSAI(s)
func FindMappingWithHomeSnssai(snssai Snssai, mappings []MappingOfSnssai) (MappingOfSnssai, bool) {
	for _, mapping := range mappings {
		if *mapping.HomeSnssai == snssai {
			return mapping, true
		}
	}
	return MappingOfSnssai{}, false
}

// Add Allowed S-NSSAI to Authorized Network Slice Info
func AddAllowedSnssai(allowedSnssai AllowedSnssai, accessType AccessType, a *AuthorizedNetworkSliceInfo) {
	hitAllowedNssai := false
	for i := range a.AllowedNssaiList {
		if a.AllowedNssaiList[i].AccessType == accessType {
			hitAllowedNssai = true
			if len(a.AllowedNssaiList[i].AllowedSnssaiList) == 8 {
				logger.Util.Infof("Unable to add a new Allowed S-NSSAI since already eight S-NSSAIs in Allowed NSSAI")
			} else {
				a.AllowedNssaiList[i].AllowedSnssaiList = append(a.AllowedNssaiList[i].AllowedSnssaiList, allowedSnssai)
			}
			break
		}
	}

	if !hitAllowedNssai {
		var allowedNssaiElement AllowedNssai
		allowedNssaiElement.AllowedSnssaiList = append(allowedNssaiElement.AllowedSnssaiList, allowedSnssai)
		allowedNssaiElement.AccessType = accessType

		a.AllowedNssaiList = append(a.AllowedNssaiList, allowedNssaiElement)
	}
}

// Add AMF information to Authorized Network Slice Info
func AddAmfInformation(tai Tai, a *AuthorizedNetworkSliceInfo) {
	if a.AllowedNssaiList == nil || len(a.AllowedNssaiList) == 0 {
		return
	}

	// Check if any AMF can serve the UE
	// That is, whether NSSAI of all Allowed S-NSSAIs is a subset of NSSAI supported by AMF

	// Find AMF Set that could serve UE from AMF Set list in configuration
	// Simply use the first applicable AMF set
	// TODO: Policies of AMF selection (e.g. load balance between AMF instances)
	for _, amfSetConfig := range factory.NssfConfig.Configuration.AmfSetList {
		hitAllowedNssai := true
		for _, allowedNssai := range a.AllowedNssaiList {
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
			if amfSetConfig.AmfList != nil && len(amfSetConfig.AmfList) != 0 {
				// List of candidate AMF(s) provided in configuration
				a.CandidateAmfList = append(a.CandidateAmfList, amfSetConfig.AmfList...)
			} else {
				// TODO: Possibly querying the NRF
				a.TargetAmfSet = amfSetConfig.AmfSetId
				// The API URI of the NRF may be included if target AMF Set is included
				a.NrfAmfSet = amfSetConfig.NrfAmfSet
			}
			return
		}
	}

	// No AMF Set in configuration can serve the UE
	// Find all candidate AMFs that could serve UE from AMF list in configuration
	hitAmf := false
	for _, amfConfig := range factory.NssfConfig.Configuration.AmfList {
		hitAllowedNssai := true
		for _, allowedNssai := range a.AllowedNssaiList {
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
			a.CandidateAmfList = append(a.CandidateAmfList, amfConfig.NfId)
			hitAmf = true
		}
	}

	if !hitAmf {
		logger.Util.Warnf("No candidate AMF or AMF Set can serve the UE")
	}
}
