package context

import (
	"encoding/json"
	"fmt"
	"free5gc/lib/MongoDBLibrary"

	// "free5gc/lib/TimeDecode"
	// "free5gc/lib/http_wrapper"
	"free5gc/lib/openapi"
	"free5gc/lib/openapi/models"
	"free5gc/src/nrf/factory"
	"free5gc/src/nrf/logger"
	"math/rand"

	// "net"
	// "reflect"
	"strconv"
	"time"

	//"github.com/gin-gonic/gin"
	"github.com/mitchellh/mapstructure"
	"go.mongodb.org/mongo-driver/bson"
)

func NnrfNFManagementDataModel(nf *models.NfProfile, nfprofile models.NfProfile) error {

	if nfprofile.NfInstanceId != "" {
		nf.NfInstanceId = nfprofile.NfInstanceId
	} else {
		return fmt.Errorf("NfInstanceId field is required\n")
	}

	if nfprofile.NfType != "" {
		nf.NfType = nfprofile.NfType
	} else {
		return fmt.Errorf("NfType field is required\n")
	}

	if nfprofile.NfStatus != "" {
		nf.NfStatus = nfprofile.NfStatus
	} else {
		return fmt.Errorf("NfStatus field is required\n")
	}

	nnrfNFManagementCondition(nf, nfprofile)
	nnrfNFManagementOption(nf, nfprofile)

	return nil

}

func SetsubscriptionId() string {
	rand.Seed(time.Now().UnixNano())
	x := rand.Intn(100)
	return strconv.Itoa(x)
}

func nnrfNFManagementCondition(nf *models.NfProfile, nfprofile models.NfProfile) {

	//HeartBeatTimer
	if nfprofile.HeartBeatTimer >= 0 {
		nf.HeartBeatTimer = nfprofile.HeartBeatTimer
	}
	//PlmnList
	if nfprofile.PlmnList != nil {
		a := make([]models.PlmnId, len(*nfprofile.PlmnList))
		for i := 0; i < len(*nfprofile.PlmnList); i++ {
			a[i] = (*nfprofile.PlmnList)[i]
		}
		nf.PlmnList = &a
	} else {
		nf.PlmnList = &[]models.PlmnId{
			factory.NrfConfig.Configuration.DefaultPlmnId,
		}
	}
	//fqdn
	if nfprofile.Fqdn != "" {
		nf.Fqdn = nfprofile.Fqdn
	}
	//interPlmnFqdn
	if nfprofile.InterPlmnFqdn != "" {
		nf.InterPlmnFqdn = nfprofile.InterPlmnFqdn
	}
	//ipv4Addresses
	if nfprofile.Ipv4Addresses != nil {
		//fmt.Println("NsiList")
		a := make([]string, len(nfprofile.Ipv4Addresses))
		for i := 0; i < len(nfprofile.Ipv4Addresses); i++ {
			a[i] = (nfprofile.Ipv4Addresses)[i]
		}
		nf.Ipv4Addresses = a
	}
	//ipv6Addresses
	if nfprofile.Ipv6Addresses != nil {
		//fmt.Println("NsiList")
		a := make([]string, len(nfprofile.Ipv6Addresses))
		for i := 0; i < len(nfprofile.Ipv6Addresses); i++ {
			a[i] = (nfprofile.Ipv6Addresses)[i]
		}
		nf.Ipv6Addresses = a
	}

}
func nnrfNFManagementOption(nf *models.NfProfile, nfprofile models.NfProfile) {

	//sNssais
	if nfprofile.SNssais != nil {
		//fmt.Println("SNssais")
		a := make([]models.Snssai, len(*nfprofile.SNssais))
		for i := 0; i < len(*nfprofile.SNssais); i++ {
			a[i] = (*nfprofile.SNssais)[i]
		}
		nf.SNssais = &a
	}

	//nsiList
	if nfprofile.NsiList != nil {
		//fmt.Println("NsiList")
		a := make([]string, len(nfprofile.NsiList))
		for i := 0; i < len(nfprofile.NsiList); i++ {
			a[i] = (nfprofile.NsiList)[i]
		}
		nf.NsiList = a
	}
	//allowedPlmns
	if nfprofile.AllowedPlmns != nil {
		a := make([]models.PlmnId, len(*nfprofile.AllowedPlmns))
		for i := 0; i < len(*nfprofile.AllowedPlmns); i++ {
			a[i] = (*nfprofile.AllowedPlmns)[i]
		}
		nf.AllowedPlmns = &a
	}
	//allowedNfTypes
	if nfprofile.AllowedNfTypes != nil {
		a := make([]models.NfType, len(nfprofile.AllowedNfTypes))
		for i := 0; i < len(nfprofile.AllowedNfTypes); i++ {
			a[i] = (nfprofile.AllowedNfTypes)[i]
		}
		nf.AllowedNfTypes = a
	}
	//allowedNfDomains
	if nfprofile.AllowedNfDomains != nil {
		a := make([]string, len(nfprofile.AllowedNfDomains))
		for i := 0; i < len(nfprofile.AllowedNfDomains); i++ {
			a[i] = (nfprofile.AllowedNfDomains)[i]
		}
		nf.AllowedNfDomains = a
	}

	//allowedNssais
	if nfprofile.AllowedNssais != nil {
		//fmt.Println("SNssais")
		a := make([]models.Snssai, len(*nfprofile.AllowedNssais))
		for i := 0; i < len(*nfprofile.AllowedNssais); i++ {
			a[i] = (*nfprofile.AllowedNssais)[i]
		}
		nf.AllowedNssais = &a
	}
	//Priority
	if nfprofile.Priority > 0 && nfprofile.Priority <= 65535 {
		nf.Priority = nfprofile.Priority
	}
	//Capacity
	if nfprofile.Capacity > 0 && nfprofile.Capacity <= 65535 {
		nf.Capacity = nfprofile.Capacity
	}
	//Load
	if nfprofile.Load > 0 && nfprofile.Load <= 100 {
		nf.Load = nfprofile.Load
	}
	//Locality
	if nfprofile.Locality != "" {
		nf.Locality = nfprofile.Locality
	}
	//udrInfo
	if nfprofile.UdrInfo != nil {

		var a models.UdrInfo

		if nfprofile.UdrInfo.GroupId != "" {
			a.GroupId = nfprofile.UdrInfo.GroupId
		}

		if nfprofile.UdrInfo.SupiRanges != nil {
			a.SupiRanges = nfprofile.UdrInfo.SupiRanges
		}

		if nfprofile.UdrInfo.GpsiRanges != nil {
			a.GpsiRanges = nfprofile.UdrInfo.GpsiRanges
		}

		if nfprofile.UdrInfo.ExternalGroupIdentifiersRanges != nil {
			a.ExternalGroupIdentifiersRanges = nfprofile.UdrInfo.ExternalGroupIdentifiersRanges
		}

		if nfprofile.UdrInfo.SupportedDataSets != nil {
			a.SupportedDataSets = nfprofile.UdrInfo.SupportedDataSets
		}

		nf.UdrInfo = &a
	}
	//udmInfo
	if nfprofile.UdmInfo != nil {

		var a models.UdmInfo

		if nfprofile.UdmInfo.GroupId != "" {
			a.GroupId = nfprofile.UdmInfo.GroupId
		}

		if nfprofile.UdmInfo.SupiRanges != nil {
			a.SupiRanges = nfprofile.UdmInfo.SupiRanges
		}

		if nfprofile.UdmInfo.GpsiRanges != nil {
			a.GpsiRanges = nfprofile.UdmInfo.GpsiRanges
		}

		if nfprofile.UdmInfo.ExternalGroupIdentifiersRanges != nil {
			a.ExternalGroupIdentifiersRanges = nfprofile.UdmInfo.ExternalGroupIdentifiersRanges
		}

		if nfprofile.UdmInfo.RoutingIndicators != nil {
			a.RoutingIndicators = nfprofile.UdmInfo.RoutingIndicators
		}

		nf.UdmInfo = &a

	}
	//ausfInfo
	if nfprofile.AusfInfo != nil {

		var a models.AusfInfo

		if nfprofile.AusfInfo.GroupId != "" {
			a.GroupId = nfprofile.AusfInfo.GroupId
		}

		if nfprofile.AusfInfo.SupiRanges != nil {
			a.SupiRanges = nfprofile.AusfInfo.SupiRanges
		}

		if nfprofile.AusfInfo.RoutingIndicators != nil {
			a.RoutingIndicators = nfprofile.AusfInfo.RoutingIndicators
		}

		nf.AusfInfo = &a
	}
	//amfInfo
	if nfprofile.AmfInfo != nil {

		var a models.AmfInfo

		if nfprofile.AmfInfo.AmfSetId != "" {
			a.AmfSetId = nfprofile.AmfInfo.AmfSetId
		}

		if nfprofile.AmfInfo.AmfRegionId != "" {
			a.AmfRegionId = nfprofile.AmfInfo.AmfRegionId
		}

		if nfprofile.AmfInfo.GuamiList != nil {
			a.GuamiList = nfprofile.AmfInfo.GuamiList
		}

		if nfprofile.AmfInfo.TaiList != nil {
			a.TaiList = nfprofile.AmfInfo.TaiList
		}

		if nfprofile.AmfInfo.TaiRangeList != nil {
			a.TaiRangeList = nfprofile.AmfInfo.TaiRangeList
		}

		if nfprofile.AmfInfo.BackupInfoAmfFailure != nil {
			a.BackupInfoAmfFailure = nfprofile.AmfInfo.BackupInfoAmfFailure
		}

		if nfprofile.AmfInfo.BackupInfoAmfRemoval != nil {
			a.BackupInfoAmfRemoval = nfprofile.AmfInfo.BackupInfoAmfRemoval
		}

		if nfprofile.AmfInfo.N2InterfaceAmfInfo != nil {
			a.N2InterfaceAmfInfo = nfprofile.AmfInfo.N2InterfaceAmfInfo
		}
		nf.AmfInfo = &a
	}
	//smfInfo
	if nfprofile.SmfInfo != nil {

		var a models.SmfInfo

		if nfprofile.SmfInfo.SNssaiSmfInfoList != nil {
			a.SNssaiSmfInfoList = nfprofile.SmfInfo.SNssaiSmfInfoList
		}
		if nfprofile.SmfInfo.TaiList != nil {
			a.TaiList = nfprofile.SmfInfo.TaiList
		}
		if nfprofile.SmfInfo.TaiRangeList != nil {
			a.TaiRangeList = nfprofile.SmfInfo.TaiRangeList
		}
		if nfprofile.SmfInfo.PgwFqdn != "" {
			a.PgwFqdn = nfprofile.SmfInfo.PgwFqdn
		}
		if nfprofile.SmfInfo.AccessType != nil {
			a.AccessType = nfprofile.SmfInfo.AccessType
		}
		nf.SmfInfo = &a
	}
	//upfInfo
	if nfprofile.UpfInfo != nil {

		var a models.UpfInfo

		if nfprofile.UpfInfo.SNssaiUpfInfoList != nil {
			a.SNssaiUpfInfoList = nfprofile.UpfInfo.SNssaiUpfInfoList
		}
		if nfprofile.UpfInfo.SmfServingArea != nil {
			a.SmfServingArea = nfprofile.UpfInfo.SmfServingArea
		}
		if nfprofile.UpfInfo.InterfaceUpfInfoList != nil {
			a.InterfaceUpfInfoList = nfprofile.UpfInfo.InterfaceUpfInfoList
		}

		a.IwkEpsInd = nfprofile.UpfInfo.IwkEpsInd

		nf.UpfInfo = &a
	}
	//pcfInfo
	if nfprofile.PcfInfo != nil {

		var a models.PcfInfo

		if nfprofile.PcfInfo.DnnList != nil {
			a.DnnList = nfprofile.PcfInfo.DnnList
		}
		if nfprofile.PcfInfo.SupiRanges != nil {
			a.SupiRanges = nfprofile.PcfInfo.SupiRanges
		}
		if nfprofile.PcfInfo.RxDiamHost != "" {
			a.RxDiamHost = nfprofile.PcfInfo.RxDiamHost
		}
		if nfprofile.PcfInfo.RxDiamRealm != "" {
			a.RxDiamRealm = nfprofile.PcfInfo.RxDiamRealm
		}
		nf.PcfInfo = &a
	}
	//bsfInfo
	if nfprofile.BsfInfo != nil {

		var a models.BsfInfo

		if nfprofile.BsfInfo.DnnList != nil {
			a.DnnList = nfprofile.BsfInfo.DnnList
		}
		if nfprofile.BsfInfo.IpDomainList != nil {
			a.IpDomainList = nfprofile.BsfInfo.IpDomainList
		}
		if nfprofile.BsfInfo.Ipv4AddressRanges != nil {

			b := make([]models.Ipv4AddressRange, len(*nfprofile.BsfInfo.Ipv4AddressRanges))
			for i := 0; i < len(*nfprofile.BsfInfo.Ipv4AddressRanges); i++ {
				b[i].Start = strconv.Itoa(int(Ipv4ToInt((*nfprofile.BsfInfo.Ipv4AddressRanges)[i].Start)))
				b[i].End = strconv.Itoa(int(Ipv4ToInt((*nfprofile.BsfInfo.Ipv4AddressRanges)[i].End)))
			}
			a.Ipv4AddressRanges = &b
		}
		if nfprofile.BsfInfo.Ipv6PrefixRanges != nil {

			b := make([]models.Ipv6PrefixRange, len(*nfprofile.BsfInfo.Ipv6PrefixRanges))
			for i := 0; i < len(*nfprofile.BsfInfo.Ipv6PrefixRanges); i++ {
				b[i].Start = Ipv6ToInt(((*nfprofile.BsfInfo.Ipv6PrefixRanges)[i].Start)).String()
				b[i].End = Ipv6ToInt(((*nfprofile.BsfInfo.Ipv6PrefixRanges)[i].End)).String()
			}
			a.Ipv6PrefixRanges = &b
		}
		nf.BsfInfo = &a
	}
	//chfInfo
	if nfprofile.ChfInfo != nil {

		var a models.ChfInfo

		if nfprofile.ChfInfo.SupiRangeList != nil {
			a.SupiRangeList = nfprofile.ChfInfo.SupiRangeList
		}
		if nfprofile.ChfInfo.GpsiRangeList != nil {
			a.GpsiRangeList = nfprofile.ChfInfo.GpsiRangeList
		}
		if nfprofile.ChfInfo.PlmnRangeList != nil {
			a.PlmnRangeList = nfprofile.ChfInfo.PlmnRangeList
		}
		nf.ChfInfo = &a
	}
	//nrfInfo
	if nfprofile.NrfInfo != nil {
		nf.NrfInfo = nfprofile.NrfInfo
	}
	//recoveryTime
	if nfprofile.RecoveryTime != nil {
		// Update when restart (Setting by NF itself)
		nf.RecoveryTime = nfprofile.RecoveryTime
	}

	//nfServicePersistence
	if nfprofile.NfServicePersistence {
		nf.NfServicePersistence = true
	} else {
		nf.NfServicePersistence = false
	}

	//nfServices
	if nfprofile.NfServices != nil {

		a := make([]models.NfService, len(*nfprofile.NfServices))
		for i := 0; i < len(*nfprofile.NfServices); i++ {
			a[i] = (*nfprofile.NfServices)[i]
		}
		nf.NfServices = &a
	}
	//

}

func SetLocationHeader(nfprofile models.NfProfile, IP string) string {

	var modifyUL UriList
	var locationHeader []string

	//set nfprofile location
	uri := fmt.Sprintf("%s://%s%s%s%s", factory.NrfConfig.Configuration.Sbi.Scheme, IP, ":29510",
		"/nnrf-nfm/v1/nf-instances/", nfprofile.NfInstanceId)
	locationHeader = append(locationHeader, uri)

	collName := "urilist"
	nfType := nfprofile.NfType
	filter := bson.M{"nfType": nfType}

	ul := MongoDBLibrary.RestfulAPIGetOne(collName, filter)

	var originalUL UriList
	err := mapstructure.Decode(ul, &originalUL)
	if err != nil {
		panic(err)
	}

	//obtain location header = NF URI
	nnrfUriList(&originalUL, &modifyUL, locationHeader)
	modifyUL.NfType = nfprofile.NfType

	tmp, err := json.Marshal(modifyUL)
	if err != nil {
		logger.ManagementLog.Error(err)
	}
	var putData = bson.M{}
	err = json.Unmarshal(tmp, &putData)
	if err != nil {
		logger.ManagementLog.Error(err)
	}

	if MongoDBLibrary.RestfulAPIPutOne(collName, filter, putData) {
		logger.ManagementLog.Info("urilist update")
	} else {
		logger.ManagementLog.Info("urilist create")
	}

	return locationHeader[0]

}

func setUriListByFilter(filter bson.M, uriList *[]string) {
	filterNfTypeResultsRaw := MongoDBLibrary.RestfulAPIGetMany("Subscriptions", filter)
	var filterNfTypeResults []models.NrfSubscriptionData
	err := openapi.Convert(filterNfTypeResultsRaw, &filterNfTypeResults)

	if err != nil {
		logger.ManagementLog.Error(err)
	}

	for _, subscr := range filterNfTypeResults {
		*uriList = append(*uriList, subscr.NfStatusNotificationUri)
	}
}

func nnrfUriList(originalUL *UriList, UL *UriList, location []string) {

	var i int
	var b *Links
	var flag bool
	var c []Item
	flag = true
	b = new(Links)
	size := len(location) + len(originalUL.Link.Item)

	//check duplicate
	for i = 0; i < len(originalUL.Link.Item); i++ {
		if originalUL.Link.Item[i].Href == location[0] {
			flag = false
		}
	}

	if flag {
		c = make([]Item, size)
		for i = 0; i < len(originalUL.Link.Item); i++ {
			c[i].Href = originalUL.Link.Item[i].Href
		}
		for i = len(originalUL.Link.Item); i < len(location)+len(originalUL.Link.Item); i++ {
			c[i].Href = location[i-len(originalUL.Link.Item)]
		}
	} else {
		c = make([]Item, size-1)
		for i = 0; i < len(originalUL.Link.Item); i++ {
			c[i].Href = originalUL.Link.Item[i].Href
		}
	}

	b.Item = c
	UL.Link = *b
}

func GetNofificationUri(nfProfile models.NfProfile) []string {
	var uriList []string

	// nfTypeCond
	var nfTypeCond = bson.M{
		"subscrCond": bson.M{
			"nfType": nfProfile.NfType,
		},
	}
	setUriListByFilter(nfTypeCond, &uriList)

	// NfInstanceIdCond
	var nfInstanceIDCond = bson.M{
		"subscrCond": bson.M{
			"nfInstanceId": nfProfile.NfInstanceId,
		},
	}
	setUriListByFilter(nfInstanceIDCond, &uriList)

	// ServiceNameCond
	if nfProfile.NfServices != nil {
		var ServiceNameCond bson.M
		var serviceNames bson.A
		for _, nfService := range *nfProfile.NfServices {
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
		var amfCond = bson.M{
			"subscrCond": bson.M{
				"amfSetId":    (*nfProfile.AmfInfo).AmfSetId,
				"amfRegionId": (*nfProfile.AmfInfo).AmfRegionId,
			},
		}
		setUriListByFilter(amfCond, &uriList)
	}

	// GuamiListCond
	if nfProfile.AmfInfo != nil {
		var guamiListFilter bson.M
		if (*nfProfile.AmfInfo).GuamiList != nil {
			var guamiListBsonArray bson.A
			for _, guami := range *(*nfProfile.AmfInfo).GuamiList {
				tmp, err := json.Marshal(guami)
				if err != nil {
					logger.ManagementLog.Error(err)
				}
				var guamiMarshal = bson.M{}
				err = json.Unmarshal(tmp, &guamiMarshal)
				if err != nil {
					logger.ManagementLog.Error(err)
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
		for _, snssai := range *nfProfile.SNssais {
			tmp, err := json.Marshal(snssai)
			if err != nil {
				logger.ManagementLog.Error(err)
			}
			var snssaiMarshal = bson.M{}
			err = json.Unmarshal(tmp, &snssaiMarshal)
			if err != nil {
				logger.ManagementLog.Error(err)
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
		var nfGroupCond = bson.M{
			"subscrCond": bson.M{
				"nfType":    nfProfile.NfType,
				"nfGroupId": (*nfProfile.UdrInfo).GroupId,
			},
		}
		setUriListByFilter(nfGroupCond, &uriList)
	} else if nfProfile.UdmInfo != nil {
		var nfGroupCond = bson.M{
			"subscrCond": bson.M{
				"nfType":    nfProfile.NfType,
				"nfGroupId": (*nfProfile.UdmInfo).GroupId,
			},
		}
		setUriListByFilter(nfGroupCond, &uriList)
	} else if nfProfile.AusfInfo != nil {
		var nfGroupCond = bson.M{
			"subscrCond": bson.M{
				"nfType":    nfProfile.NfType,
				"nfGroupId": (*nfProfile.AusfInfo).GroupId,
			},
		}
		setUriListByFilter(nfGroupCond, &uriList)
	}

	return uriList
}

func NnrfUriListLimit(originalUL *UriList, limit int) {
	//response limit

	if limit < len(originalUL.Link.Item) {
		var i int
		var b *Links = new(Links)
		var c []Item = make([]Item, limit)
		for i = 0; i < limit; i++ {
			c[i].Href = originalUL.Link.Item[i].Href
		}
		b.Item = c
		originalUL.Link = *b
	}

}
