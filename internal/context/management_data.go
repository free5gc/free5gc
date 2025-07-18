package context

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/mitchellh/mapstructure"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/free5gc/nrf/internal/logger"
	"github.com/free5gc/nrf/pkg/factory"
	"github.com/free5gc/openapi"
	"github.com/free5gc/openapi/models"
	"github.com/free5gc/util/mongoapi"
)

const NRF_NFINST_RES_URI_PREFIX = factory.NrfNfmResUriPrefix + "/nf-instances/"

func NnrfNFManagementDataModel(nf *models.NrfNfManagementNfProfile, nfprofile *models.NrfNfManagementNfProfile) error {
	if nfprofile.NfInstanceId != "" {
		nf.NfInstanceId = nfprofile.NfInstanceId
	} else {
		return fmt.Errorf("NfInstanceId field is required")
	}

	if nfprofile.NfType != "" {
		nf.NfType = nfprofile.NfType
	} else {
		return fmt.Errorf("NfType field is required")
	}

	if nfprofile.NfStatus != "" {
		nf.NfStatus = nfprofile.NfStatus
	} else {
		return fmt.Errorf("NfStatus field is required")
	}

	nnrfNFManagementCondition(nf, nfprofile)
	nnrfNFManagementOption(nf, nfprofile)

	return nil
}

func SetsubscriptionId() (string, error) {
	subscriptionIdSize := 16

	buffer := make([]byte, subscriptionIdSize)
	_, err := rand.Read(buffer)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(buffer), nil
}

func nnrfNFManagementCondition(nf *models.NrfNfManagementNfProfile, nfprofile *models.NrfNfManagementNfProfile) {
	// HeartBeatTimer
	if nfprofile.HeartBeatTimer >= 0 {
		nf.HeartBeatTimer = nfprofile.HeartBeatTimer
	}
	// PlmnList
	if nfprofile.PlmnList != nil {
		plmnList := make([]models.PlmnId, len(nfprofile.PlmnList))
		copy(plmnList, nfprofile.PlmnList)
		nf.PlmnList = plmnList
	} else {
		nf.PlmnList = []models.PlmnId{
			factory.NrfConfig.Configuration.DefaultPlmnId,
		}
	}
	// fqdn
	if nfprofile.Fqdn != "" {
		nf.Fqdn = nfprofile.Fqdn
	}
	// interPlmnFqdn
	if nfprofile.InterPlmnFqdn != "" {
		nf.InterPlmnFqdn = nfprofile.InterPlmnFqdn
	}
	// ipv4Addresses
	if nfprofile.Ipv4Addresses != nil {
		// fmt.Println("NsiList")
		ipv4 := make([]string, len(nfprofile.Ipv4Addresses))
		copy(ipv4, nfprofile.Ipv4Addresses)
		nf.Ipv4Addresses = ipv4
	}
	// ipv6Addresses
	if nfprofile.Ipv6Addresses != nil {
		// fmt.Println("NsiList")
		ipv6 := make([]string, len(nfprofile.Ipv6Addresses))
		copy(ipv6, nfprofile.Ipv6Addresses)
		nf.Ipv6Addresses = ipv6
	}
	// DefaultNotificationSubscription
	if nfprofile.DefaultNotificationSubscriptions != nil {
		nf.DefaultNotificationSubscriptions = append(nf.DefaultNotificationSubscriptions,
			nfprofile.DefaultNotificationSubscriptions...)
	}
}

func nnrfNFManagementOption(nf *models.NrfNfManagementNfProfile, nfprofile *models.NrfNfManagementNfProfile) {
	// sNssais
	if nfprofile.SNssais != nil {
		// fmt.Println("SNssais")
		snssais := make([]models.ExtSnssai, len(nfprofile.SNssais))
		copy(snssais, nfprofile.SNssais)
		nf.SNssais = snssais
	}

	// nsiList
	if nfprofile.NsiList != nil {
		// fmt.Println("NsiList")
		nsiList := make([]string, len(nfprofile.NsiList))
		copy(nsiList, nfprofile.NsiList)
		nf.NsiList = nsiList
	}
	// allowedPlmns
	if nfprofile.AllowedPlmns != nil {
		allowedPlmns := make([]models.PlmnId, len(nfprofile.AllowedPlmns))
		copy(allowedPlmns, nfprofile.AllowedPlmns)
		nf.AllowedPlmns = allowedPlmns
	}
	// allowedNfTypes
	if nfprofile.AllowedNfTypes != nil {
		allowedNfType := make([]models.NrfNfManagementNfType, len(nfprofile.AllowedNfTypes))
		copy(allowedNfType, nfprofile.AllowedNfTypes)
		nf.AllowedNfTypes = allowedNfType
	}
	// allowedNfDomains
	if nfprofile.AllowedNfDomains != nil {
		allowedNfDomains := make([]string, len(nfprofile.AllowedNfDomains))
		copy(allowedNfDomains, nfprofile.AllowedNfDomains)
		nf.AllowedNfDomains = allowedNfDomains
	}

	// allowedNssais
	if nfprofile.AllowedNssais != nil {
		// fmt.Println("SNssais")
		allowedNssais := make([]models.ExtSnssai, len(nfprofile.AllowedNssais))
		copy(allowedNssais, nfprofile.AllowedNssais)
		nf.AllowedNssais = allowedNssais
	}
	// Priority
	if nfprofile.Priority > 0 && nfprofile.Priority <= 65535 {
		nf.Priority = nfprofile.Priority
	}
	// Capacity
	if nfprofile.Capacity > 0 && nfprofile.Capacity <= 65535 {
		nf.Capacity = nfprofile.Capacity
	}
	// Load
	if nfprofile.Load > 0 && nfprofile.Load <= 100 {
		nf.Load = nfprofile.Load
	}
	// Locality
	if nfprofile.Locality != "" {
		nf.Locality = nfprofile.Locality
	}
	// udrInfo
	if nfprofile.UdrInfo != nil {
		var udrInfo models.UdrInfo

		if nfprofile.UdrInfo.GroupId != "" {
			udrInfo.GroupId = nfprofile.UdrInfo.GroupId
		}

		if nfprofile.UdrInfo.SupiRanges != nil {
			udrInfo.SupiRanges = nfprofile.UdrInfo.SupiRanges
		}

		if nfprofile.UdrInfo.GpsiRanges != nil {
			udrInfo.GpsiRanges = nfprofile.UdrInfo.GpsiRanges
		}

		if nfprofile.UdrInfo.ExternalGroupIdentifiersRanges != nil {
			udrInfo.ExternalGroupIdentifiersRanges = nfprofile.UdrInfo.ExternalGroupIdentifiersRanges
		}

		if nfprofile.UdrInfo.SupportedDataSets != nil {
			udrInfo.SupportedDataSets = nfprofile.UdrInfo.SupportedDataSets
		}

		nf.UdrInfo = &udrInfo
	}
	// udmInfo
	if nfprofile.UdmInfo != nil {
		var udmInfo models.UdmInfo

		if nfprofile.UdmInfo.GroupId != "" {
			udmInfo.GroupId = nfprofile.UdmInfo.GroupId
		}

		if nfprofile.UdmInfo.SupiRanges != nil {
			udmInfo.SupiRanges = nfprofile.UdmInfo.SupiRanges
		}

		if nfprofile.UdmInfo.GpsiRanges != nil {
			udmInfo.GpsiRanges = nfprofile.UdmInfo.GpsiRanges
		}

		if nfprofile.UdmInfo.ExternalGroupIdentifiersRanges != nil {
			udmInfo.ExternalGroupIdentifiersRanges = nfprofile.UdmInfo.ExternalGroupIdentifiersRanges
		}

		if nfprofile.UdmInfo.RoutingIndicators != nil {
			udmInfo.RoutingIndicators = nfprofile.UdmInfo.RoutingIndicators
		}

		nf.UdmInfo = &udmInfo
	}
	// ausfInfo
	if nfprofile.AusfInfo != nil {
		var ausfInfo models.AusfInfo

		if nfprofile.AusfInfo.GroupId != "" {
			ausfInfo.GroupId = nfprofile.AusfInfo.GroupId
		}

		if nfprofile.AusfInfo.SupiRanges != nil {
			ausfInfo.SupiRanges = nfprofile.AusfInfo.SupiRanges
		}

		if nfprofile.AusfInfo.RoutingIndicators != nil {
			ausfInfo.RoutingIndicators = nfprofile.AusfInfo.RoutingIndicators
		}

		nf.AusfInfo = &ausfInfo
	}
	// amfInfo
	if nfprofile.AmfInfo != nil {
		var amfInfo models.NrfNfManagementAmfInfo

		if nfprofile.AmfInfo.AmfSetId != "" {
			amfInfo.AmfSetId = nfprofile.AmfInfo.AmfSetId
		}

		if nfprofile.AmfInfo.AmfRegionId != "" {
			amfInfo.AmfRegionId = nfprofile.AmfInfo.AmfRegionId
		}

		if nfprofile.AmfInfo.GuamiList != nil {
			amfInfo.GuamiList = nfprofile.AmfInfo.GuamiList
		}

		if nfprofile.AmfInfo.TaiList != nil {
			amfInfo.TaiList = nfprofile.AmfInfo.TaiList
		}

		if nfprofile.AmfInfo.TaiRangeList != nil {
			amfInfo.TaiRangeList = nfprofile.AmfInfo.TaiRangeList
		}

		if nfprofile.AmfInfo.BackupInfoAmfFailure != nil {
			amfInfo.BackupInfoAmfFailure = nfprofile.AmfInfo.BackupInfoAmfFailure
		}

		if nfprofile.AmfInfo.BackupInfoAmfRemoval != nil {
			amfInfo.BackupInfoAmfRemoval = nfprofile.AmfInfo.BackupInfoAmfRemoval
		}

		if nfprofile.AmfInfo.N2InterfaceAmfInfo != nil {
			amfInfo.N2InterfaceAmfInfo = nfprofile.AmfInfo.N2InterfaceAmfInfo
		}
		nf.AmfInfo = &amfInfo
	}
	// smfInfo
	if nfprofile.SmfInfo != nil {
		var smfInfo models.SmfInfo

		if nfprofile.SmfInfo.SNssaiSmfInfoList != nil {
			smfInfo.SNssaiSmfInfoList = nfprofile.SmfInfo.SNssaiSmfInfoList
		}
		if nfprofile.SmfInfo.TaiList != nil {
			smfInfo.TaiList = nfprofile.SmfInfo.TaiList
		}
		if nfprofile.SmfInfo.TaiRangeList != nil {
			smfInfo.TaiRangeList = nfprofile.SmfInfo.TaiRangeList
		}
		if nfprofile.SmfInfo.PgwFqdn != "" {
			smfInfo.PgwFqdn = nfprofile.SmfInfo.PgwFqdn
		}
		if nfprofile.SmfInfo.AccessType != nil {
			smfInfo.AccessType = nfprofile.SmfInfo.AccessType
		}
		nf.SmfInfo = &smfInfo
	}
	// upfInfo
	if nfprofile.UpfInfo != nil {
		var upfInfo models.UpfInfo

		if nfprofile.UpfInfo.SNssaiUpfInfoList != nil {
			upfInfo.SNssaiUpfInfoList = nfprofile.UpfInfo.SNssaiUpfInfoList
		}
		if nfprofile.UpfInfo.SmfServingArea != nil {
			upfInfo.SmfServingArea = nfprofile.UpfInfo.SmfServingArea
		}
		if nfprofile.UpfInfo.InterfaceUpfInfoList != nil {
			upfInfo.InterfaceUpfInfoList = nfprofile.UpfInfo.InterfaceUpfInfoList
		}

		upfInfo.IwkEpsInd = nfprofile.UpfInfo.IwkEpsInd

		nf.UpfInfo = &upfInfo
	}
	// pcfInfo
	if nfprofile.PcfInfo != nil {
		var pcfInfo models.PcfInfo

		if nfprofile.PcfInfo.DnnList != nil {
			pcfInfo.DnnList = nfprofile.PcfInfo.DnnList
		}
		if nfprofile.PcfInfo.SupiRanges != nil {
			pcfInfo.SupiRanges = nfprofile.PcfInfo.SupiRanges
		}
		if nfprofile.PcfInfo.RxDiamHost != "" {
			pcfInfo.RxDiamHost = nfprofile.PcfInfo.RxDiamHost
		}
		if nfprofile.PcfInfo.RxDiamRealm != "" {
			pcfInfo.RxDiamRealm = nfprofile.PcfInfo.RxDiamRealm
		}
		nf.PcfInfo = &pcfInfo
	}
	// bsfInfo
	if nfprofile.BsfInfo != nil {
		var bsfInfo models.NrfNfManagementBsfInfo

		if nfprofile.BsfInfo.DnnList != nil {
			bsfInfo.DnnList = nfprofile.BsfInfo.DnnList
		}
		if nfprofile.BsfInfo.IpDomainList != nil {
			bsfInfo.IpDomainList = nfprofile.BsfInfo.IpDomainList
		}
		if nfprofile.BsfInfo.Ipv4AddressRanges != nil {
			Ipv4Range := make([]models.NrfNfManagementIpv4AddressRange, len(nfprofile.BsfInfo.Ipv4AddressRanges))
			for i := 0; i < len(nfprofile.BsfInfo.Ipv4AddressRanges); i++ {
				Ipv4Range[i].Start = strconv.Itoa(int(Ipv4ToInt((nfprofile.BsfInfo.Ipv4AddressRanges)[i].Start)))
				Ipv4Range[i].End = strconv.Itoa(int(Ipv4ToInt((nfprofile.BsfInfo.Ipv4AddressRanges)[i].End)))
			}
			bsfInfo.Ipv4AddressRanges = Ipv4Range
		}
		if nfprofile.BsfInfo.Ipv6PrefixRanges != nil {
			Ipv6Range := make([]models.NrfNfManagementIpv6PrefixRange, len(nfprofile.BsfInfo.Ipv6PrefixRanges))
			for i := 0; i < len(nfprofile.BsfInfo.Ipv6PrefixRanges); i++ {
				Ipv6Range[i].Start = Ipv6ToInt(nfprofile.BsfInfo.Ipv6PrefixRanges[i].Start).String()
				Ipv6Range[i].End = Ipv6ToInt(nfprofile.BsfInfo.Ipv6PrefixRanges[i].End).String()
			}
			bsfInfo.Ipv6PrefixRanges = Ipv6Range
		}
		nf.BsfInfo = &bsfInfo
	}
	// chfInfo
	if nfprofile.ChfInfo != nil {
		var chfInfo models.ChfInfo

		if nfprofile.ChfInfo.SupiRangeList != nil {
			chfInfo.SupiRangeList = nfprofile.ChfInfo.SupiRangeList
		}
		if nfprofile.ChfInfo.GpsiRangeList != nil {
			chfInfo.GpsiRangeList = nfprofile.ChfInfo.GpsiRangeList
		}
		if nfprofile.ChfInfo.PlmnRangeList != nil {
			chfInfo.PlmnRangeList = nfprofile.ChfInfo.PlmnRangeList
		}
		nf.ChfInfo = &chfInfo
	}
	// nrfInfo
	if nfprofile.NrfInfo != nil {
		nf.NrfInfo = nfprofile.NrfInfo
	}
	// recoveryTime
	if nfprofile.RecoveryTime != nil {
		// Update when restart (Setting by NF itself)
		nf.RecoveryTime = nfprofile.RecoveryTime
	}

	// nfServicePersistence
	if nfprofile.NfServicePersistence {
		nf.NfServicePersistence = true
	} else {
		nf.NfServicePersistence = false
	}

	// nfServices
	if nfprofile.NfServices != nil {
		nfServices := make([]models.NrfNfManagementNfService, len(nfprofile.NfServices))
		copy(nfServices, nfprofile.NfServices)
		nf.NfServices = nfServices
	}
	// CustomerInfo
	if nfprofile.CustomInfo != nil {
		nf.CustomInfo = nfprofile.CustomInfo
	} else {
		nf.CustomInfo = make(map[string]interface{})
	}
	nf.CustomInfo["oauth2"] = factory.NrfConfig.GetOAuth()
}

func GetNfInstanceURI(nfInstID string) string {
	return factory.NrfConfig.GetSbiUri() + NRF_NFINST_RES_URI_PREFIX + nfInstID
}

func SetLocationHeader(nfprofile *models.NrfNfManagementNfProfile) string {
	var modifyUL UriList
	var locationHeader []string

	// set nfprofile location
	locationHeader = append(locationHeader, GetNfInstanceURI(nfprofile.NfInstanceId))

	collName := "urilist"
	nfType := nfprofile.NfType
	filter := bson.M{"nfType": nfType}

	ul, err := mongoapi.RestfulAPIGetOne(collName, filter)
	if err != nil {
		logger.NfmLog.Errorf("SetLocationHeader err: %+v", err)
	}

	var originalUL UriList
	if err1 := mapstructure.Decode(ul, &originalUL); err1 != nil {
		logger.NfmLog.Errorf("SetLocationHeader err: %+v", err1)
	}

	// obtain location header = NF URI
	nnrfUriList(&originalUL, &modifyUL, locationHeader)
	modifyUL.NfType = nfprofile.NfType

	tmp, err := json.Marshal(modifyUL)
	if err != nil {
		logger.NfmLog.Errorf("SetLocationHeader err: %+v", err)
	}
	putData := bson.M{}
	err = json.Unmarshal(tmp, &putData)
	if err != nil {
		logger.NfmLog.Errorf("SetLocationHeader err: %+v", err)
	}

	existed, err := mongoapi.RestfulAPIPutOne(collName, filter, putData)
	if err != nil {
		logger.NfmLog.Errorf("SetLocationHeader err: %+v", err)
	} else {
		if existed {
			logger.NfmLog.Info("urilist update")
		} else {
			logger.NfmLog.Info("urilist create")
		}
	}

	return locationHeader[0]
}

func setUriListByFilter(filter bson.M, uriList *[]string) {
	filterNfTypeResultsRaw, err := mongoapi.RestfulAPIGetMany("Subscriptions", filter)
	if err != nil {
		logger.NfmLog.Errorf("setUriListByFilter err: %+v", err)
	}

	var filterNfTypeResults []models.NrfNfManagementSubscriptionData
	if err = openapi.Convert(filterNfTypeResultsRaw, &filterNfTypeResults); err != nil {
		logger.NfmLog.Errorf("setUriListByFilter err: %+v", err)
	}

	for _, subscr := range filterNfTypeResults {
		*uriList = append(*uriList, subscr.NfStatusNotificationUri)
	}
}

func nnrfUriList(originalUL *UriList, ul *UriList, location []string) {
	var i int
	var links *Links
	var flag bool
	var item []Item
	flag = true
	links = new(Links)
	size := len(location) + len(originalUL.Link.Item)

	// check duplicate
	for i = 0; i < len(originalUL.Link.Item); i++ {
		if originalUL.Link.Item[i].Href == location[0] {
			flag = false
		}
	}

	if flag {
		item = make([]Item, size)
		for i = 0; i < len(originalUL.Link.Item); i++ {
			item[i].Href = originalUL.Link.Item[i].Href
		}
		for i = len(originalUL.Link.Item); i < len(location)+len(originalUL.Link.Item); i++ {
			item[i].Href = location[i-len(originalUL.Link.Item)]
		}
	} else {
		item = make([]Item, size-1)
		for i = 0; i < len(originalUL.Link.Item); i++ {
			item[i].Href = originalUL.Link.Item[i].Href
		}
	}

	links.Item = item
	ul.Link = *links
}

func GetNofificationUri(nfProfile *models.NrfNfManagementNfProfile) []string {
	var uriList []string

	// nfTypeCond
	nfTypeCond := bson.M{
		"subscrCond": bson.M{
			"nfType": nfProfile.NfType,
		},
	}
	setUriListByFilter(nfTypeCond, &uriList)

	// NfInstanceIdCond
	nfInstanceIDCond := bson.M{
		"subscrCond": bson.M{
			"nfInstanceId": nfProfile.NfInstanceId,
		},
	}
	setUriListByFilter(nfInstanceIDCond, &uriList)

	// ServiceNameCond
	if nfProfile.NfServices != nil {
		var ServiceNameCond bson.M
		var serviceNames bson.A
		for _, nfService := range nfProfile.NfServices {
			serviceNames = append(serviceNames, string(nfService.ServiceName))
		}
		ServiceNameCond = bson.M{
			"subscrCond.serviceName": bson.M{
				"$in": serviceNames,
			},
		}
		setUriListByFilter(ServiceNameCond, &uriList)
	}

	// AmfCond
	if nfProfile.AmfInfo != nil {
		amfCond := bson.M{
			"subscrCond": bson.M{
				"amfSetId":    nfProfile.AmfInfo.AmfSetId,
				"amfRegionId": nfProfile.AmfInfo.AmfRegionId,
			},
		}
		setUriListByFilter(amfCond, &uriList)
	}

	// GuamiListCond
	if nfProfile.AmfInfo != nil {
		var guamiListFilter bson.M
		if nfProfile.AmfInfo.GuamiList != nil {
			var guamiListBsonArray bson.A
			for _, guami := range nfProfile.AmfInfo.GuamiList {
				tmp, err := json.Marshal(guami)
				if err != nil {
					logger.NfmLog.Error(err)
				}
				guamiMarshal := bson.M{}
				err = json.Unmarshal(tmp, &guamiMarshal)
				if err != nil {
					logger.NfmLog.Error(err)
				}

				guamiListBsonArray = append(guamiListBsonArray, bson.M{"subscrCond": bson.M{"$elemMatch": guamiMarshal}})
			}
			guamiListFilter = bson.M{
				"$or": guamiListBsonArray,
			}
		}
		setUriListByFilter(guamiListFilter, &uriList)
	}

	// NetworkSliceCond
	if nfProfile.SNssais != nil {
		var networkSliceFilter bson.M
		var snssaisBsonArray bson.A
		for _, snssai := range nfProfile.SNssais {
			tmp, err := json.Marshal(snssai)
			if err != nil {
				logger.NfmLog.Error(err)
			}
			snssaiMarshal := bson.M{}
			err = json.Unmarshal(tmp, &snssaiMarshal)
			if err != nil {
				logger.NfmLog.Error(err)
			}

			snssaisBsonArray = append(snssaisBsonArray, bson.M{"subscrCond": bson.M{"$elemMatch": snssaiMarshal}})
		}

		var nsiListBsonArray bson.A
		if nfProfile.NsiList != nil {
			for _, nsi := range nfProfile.NsiList {
				nsiListBsonArray = append(nsiListBsonArray, nsi)
			}
		}

		if nsiListBsonArray != nil {
			networkSliceFilter = bson.M{
				"$and": bson.A{
					bson.M{
						"subscrCond.nsiList": bson.M{
							"$in": nsiListBsonArray,
						},
					},
					bson.M{
						"$or": snssaisBsonArray,
					},
				},
			}
		} else {
			networkSliceFilter = bson.M{
				"$and": bson.A{
					bson.M{
						"$or": snssaisBsonArray,
					},
				},
			}
		}
		setUriListByFilter(networkSliceFilter, &uriList)
	}

	// NfGroupCond
	if nfProfile.UdrInfo != nil {
		nfGroupCond := bson.M{
			"subscrCond": bson.M{
				"nfType":    nfProfile.NfType,
				"nfGroupId": nfProfile.UdrInfo.GroupId,
			},
		}
		setUriListByFilter(nfGroupCond, &uriList)
	} else if nfProfile.UdmInfo != nil {
		nfGroupCond := bson.M{
			"subscrCond": bson.M{
				"nfType":    nfProfile.NfType,
				"nfGroupId": nfProfile.UdmInfo.GroupId,
			},
		}
		setUriListByFilter(nfGroupCond, &uriList)
	} else if nfProfile.AusfInfo != nil {
		nfGroupCond := bson.M{
			"subscrCond": bson.M{
				"nfType":    nfProfile.NfType,
				"nfGroupId": nfProfile.AusfInfo.GroupId,
			},
		}
		setUriListByFilter(nfGroupCond, &uriList)
	}

	return uriList
}

func NnrfUriListLimit(originalUL *UriList, limit int) {
	// response limit

	if limit <= 0 {
		return
	}

	if limit < len(originalUL.Link.Item) {
		var i int
		links := new(Links)
		item := make([]Item, limit)
		for i = 0; i < limit; i++ {
			item[i].Href = originalUL.Link.Item[i].Href
		}
		links.Item = item
		originalUL.Link = *links
	}
}
