package sbi

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) getImsSDMRoutes() []Route {
	return []Route{
		{
			"Index",
			"GET",
			"/",
			Index,
		},

		{
			Name:        "DeleteRepositoryDataServInd",
			Method:      "Delete",
			Pattern:     "/:imsUeId/repository-data/:serviceIndication",
			HandlerFunc: s.HTTPDeleteRepositoryDataServInd,
		},

		{
			Name:        "DeleteSmsRegistrationInfo",
			Method:      "Delete",
			Pattern:     "/:imsUeId/service-data/sms-registration-info",
			HandlerFunc: s.HTTPDeleteSmsRegistrationInfo,
		},

		{
			Name:        "GetChargingInfo",
			Method:      "Get",
			Pattern:     "/:imsUeId/ims-data/profile-data/charging-info",
			HandlerFunc: s.HTTPGetChargingInfo,
		},

		{
			Name:        "GetCsUserStateInfo",
			Method:      "Get",
			Pattern:     "/:imsUeId/access-data/cs-domain/user-state",
			HandlerFunc: s.HTTPGetCsUserStateInfo,
		},

		{
			Name:        "GetCsrn",
			Method:      "Get",
			Pattern:     "/:imsUeId/access-data/cs-domain/csrn",
			HandlerFunc: s.HTTPGetCsrn,
		},

		{
			Name:        "GetDsaiInfo",
			Method:      "Get",
			Pattern:     "/:imsUeId/service-data/dsai",
			HandlerFunc: s.HTTPGetDsaiInfo,
		},

		{
			Name:        "GetIMEISVInfo",
			Method:      "Get",
			Pattern:     "/:imsUeId/identities/imeisv",
			HandlerFunc: s.HTTPGetIMEISVInfo,
		},

		{
			Name:        "GetIfcs",
			Method:      "Get",
			Pattern:     "/:imsUeId/ims-data/profile-data/ifcs",
			HandlerFunc: s.HTTPGetIfcs,
		},

		{
			Name:        "GetImsAssocIds",
			Method:      "Get",
			Pattern:     "/:imsUeId/identities/ims-associated-identities",
			HandlerFunc: s.HTTPGetImsAssocIds,
		},

		{
			Name:        "GetImsPrivateIds",
			Method:      "Get",
			Pattern:     "/:imsUeId/identities/private-identities",
			HandlerFunc: s.HTTPGetImsPrivateIds,
		},

		{
			Name:        "GetIpAddressInfo",
			Method:      "Get",
			Pattern:     "/:imsUeId/access-data/ps-domain/ip-address",
			HandlerFunc: s.HTTPGetIpAddressInfo,
		},

		{
			Name:        "GetLocCsDomain",
			Method:      "Get",
			Pattern:     "/:imsUeId/access-data/cs-domain/location-data",
			HandlerFunc: s.HTTPGetLocCsDomain,
		},

		{
			Name:        "GetLocPsDomain",
			Method:      "Get",
			Pattern:     "/:imsUeId/access-data/ps-domain/location-data",
			HandlerFunc: s.HTTPGetLocPsDomain,
		},

		{
			Name:        "GetMsisdns",
			Method:      "Get",
			Pattern:     "/:imsUeId/identities/msisdns",
			HandlerFunc: s.HTTPGetMsisdns,
		},

		{
			Name:        "GetPriorityInfo",
			Method:      "Get",
			Pattern:     "/:imsUeId/ims-data/profile-data/priority-levels",
			HandlerFunc: s.HTTPGetPriorityInfo,
		},

		{
			Name:        "GetProfileData",
			Method:      "Get",
			Pattern:     "/:imsUeId/ims-data/profile-data",
			HandlerFunc: s.HTTPGetProfileData,
		},

		{
			Name:        "GetPsUserStateInfo",
			Method:      "Get",
			Pattern:     "/:imsUeId/access-data/ps-domain/user-state",
			HandlerFunc: s.HTTPGetPsUserStateInfo,
		},

		{
			Name:        "GetPsiState",
			Method:      "Get",
			Pattern:     "/:imsUeId/service-data/psi-status",
			HandlerFunc: s.HTTPGetPsiState,
		},

		{
			Name:        "GetReferenceLocationInfo",
			Method:      "Get",
			Pattern:     "/:imsUeId/access-data/wireline-domain/reference-location",
			HandlerFunc: s.HTTPGetReferenceLocationInfo,
		},

		{
			Name:        "GetRegistrationStatus",
			Method:      "Get",
			Pattern:     "/:imsUeId/ims-data/registration-status",
			HandlerFunc: s.HTTPGetRegistrationStatus,
		},

		{
			Name:        "GetRepositoryDataServInd",
			Method:      "Get",
			Pattern:     "/:imsUeId/repository-data/:serviceIndication",
			HandlerFunc: s.HTTPGetRepositoryDataServInd,
		},

		{
			Name:        "GetRepositoryDataServIndList",
			Method:      "Get",
			Pattern:     "/:imsUeId/repository-data",
			HandlerFunc: s.HTTPGetRepositoryDataServIndList,
		},

		{
			Name:        "GetScscfCapabilities",
			Method:      "Get",
			Pattern:     "/:imsUeId/ims-data/location-data/scscf-capabilities",
			HandlerFunc: s.HTTPGetScscfCapabilities,
		},

		{
			Name:        "GetScscfSelectionAssistanceInfo",
			Method:      "Get",
			Pattern:     "/:imsUeId/ims-data/location-data/scscf-selection-assistance-info",
			HandlerFunc: s.HTTPGetScscfSelectionAssistanceInfo,
		},

		{
			Name:        "GetServerName",
			Method:      "Get",
			Pattern:     "/:imsUeId/ims-data/location-data/server-name",
			HandlerFunc: s.HTTPGetServerName,
		},

		{
			Name:        "GetServiceTraceInfo",
			Method:      "Get",
			Pattern:     "/:imsUeId/ims-data/profile-data/service-level-trace-information",
			HandlerFunc: s.HTTPGetServiceTraceInfo,
		},

		{
			Name:        "GetSharedData",
			Method:      "Get",
			Pattern:     "/shared-data",
			HandlerFunc: s.HTTPGetSharedData,
		},

		{
			Name:        "GetSmsRegistrationInfo",
			Method:      "Get",
			Pattern:     "/:imsUeId/service-data/sms-registration-info",
			HandlerFunc: s.HTTPGetSmsRegistrationInfo,
		},

		{
			Name:        "GetSrvccData",
			Method:      "Get",
			Pattern:     "/:imsUeId/srvcc-data",
			HandlerFunc: s.HTTPGetSrvccData,
		},

		{
			Name:        "GetTadsInfo",
			Method:      "Get",
			Pattern:     "/:imsUeId/access-data/ps-domain/tads-info",
			HandlerFunc: s.HTTPGetTadsInfo,
		},

		{
			Name:        "ImsSdmSubsModify",
			Method:      "Patch",
			Pattern:     "/:imsUeId/subscriptions/:subscriptionId",
			HandlerFunc: s.HTTPImsSdmSubsModify,
		},

		{
			Name:        "ImsSdmSubscribe",
			Method:      "Post",
			Pattern:     "/:imsUeId/subscriptions",
			HandlerFunc: s.HTTPImsSdmSubscribe,
		},

		{
			Name:        "ImsSdmUnsubscribe",
			Method:      "Delete",
			Pattern:     "/:imsUeId/subscriptions/:subscriptionId",
			HandlerFunc: s.HTTPImsSdmUnsubscribe,
		},

		{
			Name:        "ModifySharedDataSubs",
			Method:      "Patch",
			Pattern:     "/shared-data-subscriptions/:subscriptionId",
			HandlerFunc: s.HTTPModifySharedDataSubs,
		},

		{
			Name:        "SubscribeToSharedData",
			Method:      "Post",
			Pattern:     "/shared-data-subscriptions",
			HandlerFunc: s.HTTPSubscribeToSharedData,
		},

		{
			Name:        "UeReachIpSubscribe",
			Method:      "Post",
			Pattern:     "/:imsUeId/access-data/ps-domain/ue-reach-subscriptions",
			HandlerFunc: s.HTTPUeReachIpSubscribe,
		},

		{
			Name:        "UeReachSubsModify",
			Method:      "Patch",
			Pattern:     "/:imsUeId/access-data/ps-domain/ue-reach-subscriptions/:subscriptionId",
			HandlerFunc: s.HTTPUeReachSubsModify,
		},

		{
			Name:        "UeReachUnsubscribe",
			Method:      "Delete",
			Pattern:     "/:imsUeId/access-data/ps-domain/ue-reach-subscriptions/:subscriptionId",
			HandlerFunc: s.HTTPUeReachUnsubscribe,
		},

		{
			Name:        "UnsubscribeForSharedData",
			Method:      "Delete",
			Pattern:     "/shared-data-subscriptions/:subscriptionId",
			HandlerFunc: s.HTTPUnsubscribeForSharedData,
		},

		{
			Name:        "UpdateDsaiState",
			Method:      "Patch",
			Pattern:     "/:imsUeId/service-data/dsai",
			HandlerFunc: s.HTTPUpdateDsaiState,
		},

		{
			Name:        "UpdatePsiState",
			Method:      "Patch",
			Pattern:     "/:imsUeId/service-data/psi-status",
			HandlerFunc: s.HTTPUpdatePsiState,
		},

		{
			Name:        "UpdateRepositoryDataServInd",
			Method:      "Put",
			Pattern:     "/:imsUeId/repository-data/:serviceIndication",
			HandlerFunc: s.HTTPUpdateRepositoryDataServInd,
		},

		{
			Name:        "UpdateSmsRegistrationInfo",
			Method:      "Put",
			Pattern:     "/:imsUeId/service-data/sms-registration-info",
			HandlerFunc: s.HTTPUpdateSmsRegistrationInfo,
		},

		{
			Name:        "UpdateSrvccData",
			Method:      "Patch",
			Pattern:     "/:imsUeId/srvcc-data",
			HandlerFunc: s.HTTPUpdateSrvccData,
		},
	}
}

// DeleteRepositoryDataServInd - delete the Repository Data for a Service Indication
func (s *Server) HTTPDeleteRepositoryDataServInd(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{})
}

// DeleteSmsRegistrationInfo - delete the SMS registration information
func (s *Server) HTTPDeleteSmsRegistrationInfo(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{})
}

// GetChargingInfo - Retrieve the charging information for to the user
func (s *Server) HTTPGetChargingInfo(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{})
}

// GetCsUserStateInfo - Retrieve the user state information in CS domain
func (s *Server) HTTPGetCsUserStateInfo(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{})
}

// GetCsrn - Retrieve the routeing number in CS domain
func (s *Server) HTTPGetCsrn(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{})
}

// GetDsaiInfo - Retrieve the DSAI information associated to an Application Server
func (s *Server) HTTPGetDsaiInfo(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{})
}

// GetIMEISVInfo - Retrieve the IMEISV information
func (s *Server) HTTPGetIMEISVInfo(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{})
}

// GetIfcs - Retrieve the Initial Filter Criteria for the associated IMS subscription
func (s *Server) HTTPGetIfcs(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{})
}

// GetImsAssocIds - Retrieve the associated identities to the IMS public identity included in the service request
func (s *Server) HTTPGetImsAssocIds(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{})
}

// GetImsPrivateIds - Retrieve the associated private identities
// to the IMS public identity included in the service request
func (s *Server) HTTPGetImsPrivateIds(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{})
}

// GetIpAddressInfo - Retrieve the IP address information
func (s *Server) HTTPGetIpAddressInfo(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{})
}

// GetLocCsDomain - Retrieve the location data in CS domain
func (s *Server) HTTPGetLocCsDomain(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{})
}

// GetLocPsDomain - Retrieve the location data in PS domain
func (s *Server) HTTPGetLocPsDomain(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{})
}

// GetMsisdns - retrieve the Msisdns associated to requested identity
func (s *Server) HTTPGetMsisdns(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{})
}

// GetPriorityInfo - Retrieve the service priority levels associated to the user
func (s *Server) HTTPGetPriorityInfo(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{})
}

// GetProfileData - Retrieve the complete IMS profile
// for a given IMS public identity (and public identities in the same IRS)
func (s *Server) HTTPGetProfileData(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{})
}

// GetPsUserStateInfo - Retrieve the user state information in PS domain
func (s *Server) HTTPGetPsUserStateInfo(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{})
}

// GetPsiState - Retrieve the PSI activation state data
func (s *Server) HTTPGetPsiState(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{})
}

// GetReferenceLocationInfo - Retrieve the reference location information
func (s *Server) HTTPGetReferenceLocationInfo(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{})
}

// GetRegistrationStatus - Retrieve the registration status of a user
func (s *Server) HTTPGetRegistrationStatus(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{})
}

// GetRepositoryDataServInd - Retrieve the repository data associated to an IMPU and service indication
func (s *Server) HTTPGetRepositoryDataServInd(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{})
}

// GetRepositoryDataServIndList - Retrieve the repository data associated to an IMPU and service indication list
func (s *Server) HTTPGetRepositoryDataServIndList(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{})
}

// GetScscfCapabilities - Retrieve the S-CSCF capabilities for the associated IMS subscription
func (s *Server) HTTPGetScscfCapabilities(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{})
}

// GetScscfSelectionAssistanceInfo - Retrieve the S-CSCF selection assistance info
func (s *Server) HTTPGetScscfSelectionAssistanceInfo(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{})
}

// GetServerName - Retrieve the server name for the associated user
func (s *Server) HTTPGetServerName(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{})
}

// GetServiceTraceInfo - Retrieve the IMS service level trace information for the associated user
func (s *Server) HTTPGetServiceTraceInfo(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{})
}

// GetSharedData - retrieve shared data
func (s *Server) HTTPGetSharedData(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{})
}

// GetSmsRegistrationInfo - Retrieve the SMS registration information associated to a user
func (s *Server) HTTPGetSmsRegistrationInfo(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{})
}

// GetSrvccData - Retrieve the srvcc data
func (s *Server) HTTPGetSrvccData(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{})
}

// GetTadsInfo - Retrieve the T-ADS information
func (s *Server) HTTPGetTadsInfo(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{})
}

// ImsSdmSubsModify - modify the subscription
func (s *Server) HTTPImsSdmSubsModify(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{})
}

// ImsSdmSubscribe - subscribe to notifications
func (s *Server) HTTPImsSdmSubscribe(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{})
}

// ImsSdmUnsubscribe - unsubscribe from notifications
func (s *Server) HTTPImsSdmUnsubscribe(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{})
}

// ModifySharedDataSubs - modify the subscription
func (s *Server) HTTPModifySharedDataSubs(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{})
}

// SubscribeToSharedData - subscribe to notifications for shared data
func (s *Server) HTTPSubscribeToSharedData(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{})
}

// UeReachIpSubscribe - subscribe to notifications of UE reachability
func (s *Server) HTTPUeReachIpSubscribe(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{})
}

// UeReachSubsModify - modify the subscription
func (s *Server) HTTPUeReachSubsModify(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{})
}

// UeReachUnsubscribe - unsubscribe from notifications to UE reachability
func (s *Server) HTTPUeReachUnsubscribe(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{})
}

// UnsubscribeForSharedData - unsubscribe from notifications for shared data
func (s *Server) HTTPUnsubscribeForSharedData(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{})
}

// UpdateDsaiState - Patch
func (s *Server) HTTPUpdateDsaiState(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{})
}

// UpdatePsiState - Patch
func (s *Server) HTTPUpdatePsiState(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{})
}

// UpdateRepositoryDataServInd - Update the repository data associated to an IMPU and service indication
func (s *Server) HTTPUpdateRepositoryDataServInd(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{})
}

// UpdateSmsRegistrationInfo - Update the SMS registration information associated to a user
func (s *Server) HTTPUpdateSmsRegistrationInfo(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{})
}

// UpdateSrvccData - Patch
func (s *Server) HTTPUpdateSrvccData(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{})
}
