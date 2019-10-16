/*
 * NSSF NS Selection
 *
 * NSSF Network Slice Selection Service
 */

package nssf_producer

import (
	"fmt"
	"math/rand"
	"net/http"

	. "free5gc/lib/openapi/models"
	. "free5gc/src/nssf/plugin"
	"free5gc/src/nssf/util"
)

func selectNsiInformation(nsiInformationList []NsiInformation) NsiInformation {
	// TODO: Algorithm to select Network Slice Instance
	//       Take roaming indication into consideration

	// Randomly select a Network Slice Instance
	idx := rand.Intn(len(nsiInformationList))
	return nsiInformationList[idx]
}

// Network slice selection for PDU session
// The function is executed when the IE, `slice-info-for-pdu-session`, is provided in query parameters
func nsselectionForPduSession(p NsselectionQueryParameter,
	a *AuthorizedNetworkSliceInfo,
	d *ProblemDetails) (status int) {
	if p.HomePlmnId != nil {
		// Check whether UE's Home PLMN is supported when UE is a roamer
		if !util.CheckSupportedHplmn(*p.HomePlmnId) {
			a.RejectedNssaiInPlmn = append(a.RejectedNssaiInPlmn, *p.SliceInfoRequestForPduSession.SNssai)

			status = http.StatusOK
			return
		}
	}

	if p.Tai != nil {
		// Check whether UE's current TA is supported when UE provides TAI
		if !util.CheckSupportedTa(*p.Tai) {
			a.RejectedNssaiInTa = append(a.RejectedNssaiInTa, *p.SliceInfoRequestForPduSession.SNssai)

			status = http.StatusOK
			return
		}
	}

	if p.Tai != nil && !util.CheckSupportedSnssaiInPlmn(*p.SliceInfoRequestForPduSession.SNssai, *p.Tai.PlmnId) {
		// Return ProblemDetails indicating S-NSSAI is not supported
		// TODO: Based on TS 23.501 V15.2.0, if the Requested NSSAI includes an S-NSSAI that is not valid in the
		//       Serving PLMN, the NSSF may derive the Configured NSSAI for Serving PLMN
		*d = ProblemDetails{
			Title:  util.UNSUPPORTED_RESOURCE,
			Status: http.StatusForbidden,
			Detail: "S-NSSAI in Requested NSSAI is not supported in PLMN",
			Cause:  "SNSSAI_NOT_SUPPORTED",
		}

		status = http.StatusForbidden
		return
	}

	if p.HomePlmnId != nil {
		if p.SliceInfoRequestForPduSession.RoamingIndication == RoamingIndication_NON_ROAMING {
			problemDetail := "`home-plmn-id` is provided, which contradicts `roamingIndication`:'NON_ROAMING'"
			*d = ProblemDetails{
				Title:  util.INVALID_REQUEST,
				Status: http.StatusBadRequest,
				Detail: problemDetail,
				InvalidParams: []InvalidParam{
					{
						Param:  "home-plmn-id",
						Reason: problemDetail,
					},
				},
			}

			status = http.StatusBadRequest
			return
		}
	} else {
		if p.SliceInfoRequestForPduSession.RoamingIndication != RoamingIndication_NON_ROAMING {
			problemDetail := fmt.Sprintf("`home-plmn-id` is not provided, which contradicts `roamingIndication`:'%s'",
				string(p.SliceInfoRequestForPduSession.RoamingIndication))
			*d = ProblemDetails{
				Title:  util.INVALID_REQUEST,
				Status: http.StatusBadRequest,
				Detail: problemDetail,
				InvalidParams: []InvalidParam{
					{
						Param:  "home-plmn-id",
						Reason: problemDetail,
					},
				},
			}

			status = http.StatusBadRequest
			return
		}
	}

	if p.Tai != nil && !util.CheckSupportedSnssaiInTa(*p.SliceInfoRequestForPduSession.SNssai, *p.Tai) {
		// Requested S-NSSAI does not supported in UE's current TA
		// Add it to Rejected NSSAI in TA
		a.RejectedNssaiInTa = append(a.RejectedNssaiInTa, *p.SliceInfoRequestForPduSession.SNssai)
		status = http.StatusOK
		return
	}

	nsiInformationList := util.GetNsiInformationListFromConfig(*p.SliceInfoRequestForPduSession.SNssai)

	nsiInformation := selectNsiInformation(nsiInformationList)

	a.NsiInformation = new(NsiInformation)
	*a.NsiInformation = nsiInformation

	return http.StatusOK
}
