package producer

import (
	"encoding/json"
	"math/big"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/free5gc/nrf/internal/context"
	"github.com/free5gc/nrf/internal/logger"
	"github.com/free5gc/openapi/models"
	"github.com/free5gc/util/httpwrapper"
	timedecode "github.com/free5gc/util/mapstruct"
	"github.com/free5gc/util/mongoapi"
)

func HandleNFDiscoveryRequest(request *httpwrapper.Request) *httpwrapper.Response {
	// Get all query parameters
	logger.DiscoveryLog.Infoln("Handle NFDiscoveryRequest")

	response, problemDetails := NFDiscoveryProcedure(request.Query)
	// Send Response
	// step 4: process the return value from step 3
	if response != nil {
		// status code is based on SPEC, and option headers
		return httpwrapper.NewResponse(http.StatusOK, nil, response)
	} else if problemDetails != nil {
		return httpwrapper.NewResponse(int(problemDetails.Status), nil, problemDetails)
	}
	problemDetails = &models.ProblemDetails{
		Status: http.StatusForbidden,
		Cause:  "UNSPECIFIED",
	}
	return httpwrapper.NewResponse(http.StatusForbidden, nil, problemDetails)
}

func NFDiscoveryProcedure(
	queryParameters url.Values,
) (response *models.SearchResult, problemDetails *models.ProblemDetails) {
	if queryParameters["target-nf-type"] == nil || queryParameters["requester-nf-type"] == nil {
		problemDetails := &models.ProblemDetails{
			Title:  "Invalid Parameter",
			Status: http.StatusBadRequest,
			Cause:  "Loss mandatory parameter",
		}
		return nil, problemDetails
	}

	if queryParameters["complexQuery"] != nil {
		// IF SUPPORT COMPLEX QUERY
		// translate raw data to complexQuery structure
		complexQuery := queryParameters["complexQuery"][0]
		complexQueryStruct := &models.ComplexQuery{}
		err := json.Unmarshal([]byte(complexQuery), complexQueryStruct)
		if err != nil {
			logger.DiscoveryLog.Warnln("UnMasrhal complexQuery Error: ", err)
		}
		// Check either CNF or DNF
		if complexQueryStruct.CNf != nil && complexQueryStruct.DNf != nil {
			problemDetails := &models.ProblemDetails{
				Title:  "Invalid Parameter",
				Status: http.StatusBadRequest,
				Cause:  "EITHER CNF OR DNF",
				InvalidParams: []models.InvalidParam{
					{Param: "complexQuery"},
				},
			}
			return nil, problemDetails
		}
	}

	// Check ComplexQuery (FOR REPORT PROBLEM!)

	// Build Query Filter
	var filter bson.M = buildFilter(queryParameters)
	logger.DiscoveryLog.Traceln("Query filter: ", filter)

	// Use the filter to find documents
	nfProfilesRaw, err := mongoapi.RestfulAPIGetMany("NfProfile", filter)
	if err != nil {
		logger.DiscoveryLog.Errorf("NFDiscoveryProcedure err: %+v", err)
		problemDetails := &models.ProblemDetails{
			Title:  "System failure",
			Status: http.StatusInternalServerError,
			Detail: err.Error(),
			Cause:  "SYSTEM_FAILURE",
		}
		return nil, problemDetails
	}

	// nfProfile data for response
	var nfProfilesStruct []models.NfProfile
	if err := timedecode.Decode(nfProfilesRaw, &nfProfilesStruct); err != nil {
		logger.DiscoveryLog.Errorf("NF Profile Raw decode error: %+v", err)
		problemDetails := &models.ProblemDetails{
			Title:  "System failure",
			Status: http.StatusInternalServerError,
			Detail: err.Error(),
			Cause:  "SYSTEM_FAILURE",
		}
		return nil, problemDetails
	}

	// handle ipv4 & ipv6
	if queryParameters["target-nf-type"][0] == "BSF" {
		for i, nfProfile := range nfProfilesStruct {
			if nfProfile.BsfInfo.Ipv4AddressRanges != nil {
				for j := range *nfProfile.BsfInfo.Ipv4AddressRanges {
					ipv4IntStart, err := strconv.Atoi((((*(*nfProfilesStruct[i].BsfInfo).Ipv4AddressRanges)[j]).Start))
					if err != nil {
						logger.DiscoveryLog.Warnln("ipv4IntStart Atoi Error: ", err)
					}
					((*(*nfProfilesStruct[i].BsfInfo).Ipv4AddressRanges)[j]).Start = context.Ipv4IntToIpv4String(int64(ipv4IntStart))
					ipv4IntEnd, err := strconv.Atoi((((*(*nfProfilesStruct[i].BsfInfo).Ipv4AddressRanges)[j]).End))
					if err != nil {
						logger.DiscoveryLog.Warnln("ipv4IntEnd Atoi Error: ", err)
					}
					((*(*nfProfilesStruct[i].BsfInfo).Ipv4AddressRanges)[j]).End = context.Ipv4IntToIpv4String(int64(ipv4IntEnd))
				}
			}
			if nfProfile.BsfInfo.Ipv6PrefixRanges != nil {
				for j := range *nfProfile.BsfInfo.Ipv6PrefixRanges {
					ipv6IntStart := new(big.Int)
					ipv6IntStart.SetString(((*(*nfProfilesStruct[i].BsfInfo).Ipv6PrefixRanges)[j]).Start, 10)
					((*(*nfProfilesStruct[i].BsfInfo).Ipv6PrefixRanges)[j]).Start = context.Ipv6IntToIpv6String(ipv6IntStart)

					ipv6IntEnd := new(big.Int)
					ipv6IntEnd.SetString(((*(*nfProfilesStruct[i].BsfInfo).Ipv6PrefixRanges)[j]).End, 10)
					((*(*nfProfilesStruct[i].BsfInfo).Ipv6PrefixRanges)[j]).End = context.Ipv6IntToIpv6String(ipv6IntEnd)
				}
			}
		}
	}
	// Build SearchResult model
	searchResult := &models.SearchResult{
		ValidityPeriod: 100,
		NfInstances:    nfProfilesStruct,
	}

	return searchResult, nil
}

func buildFilter(queryParameters url.Values) bson.M {
	// build the filter
	filter := bson.M{
		"$and": []bson.M{},
	}

	// [Query-1] target-nf-type
	targetNfType := queryParameters["target-nf-type"][0]
	if targetNfType != "" {
		targetNfTypeFilter := bson.M{
			"nfType": targetNfType,
		}
		filter["$and"] = append(filter["$and"].([]bson.M), targetNfTypeFilter)
	}

	// [Query-2] request-nf-type
	requesterNfType := queryParameters["requester-nf-type"][0]
	if requesterNfType != "" {
		requesterNfTypeFilter := bson.M{
			"$or": []bson.M{
				{"allowedNfTypes": requesterNfType},
				{"allowedNfTypes": bson.M{
					"$exists": false,
				}},
			},
		}
		filter["$and"] = append(filter["$and"].([]bson.M), requesterNfTypeFilter)
	}

	// [Query-3] service-names
	// TODO: return exist service name
	if queryParameters["service-names"] != nil {
		serviceNames := queryParameters["service-names"][0]
		serviceNamesSplit := strings.Split(serviceNames, ",")
		var serviceNamesBsonArray bson.A

		for _, v := range serviceNamesSplit {
			serviceNamesBsonArray = append(serviceNamesBsonArray, v)
		}
		serviceNamesFilter := bson.M{
			"nfServices": bson.M{
				"$elemMatch": bson.M{
					"serviceName": bson.M{
						// get all service in array
						"$in": serviceNamesBsonArray,
					},
					// the service need to be registered
					"nfServiceStatus": "REGISTERED",
				},
			},
		}
		filter["$and"] = append(filter["$and"].([]bson.M), serviceNamesFilter)
	}

	// [Query-4] requester-nfinstance-fqdn
	if queryParameters["requester-nf-instance-fqdn"] != nil {
		requesterNfinstanceFqdn := queryParameters["requester-nf-instance-fqdn"][0]

		requesterNfinstanceFqdnFilter := bson.M{
			"$or": []bson.M{
				{
					"nfServices": bson.M{
						"$elemMatch": bson.M{
							"allowedNfDomains": requesterNfinstanceFqdn,
						},
					},
				},
				{ // if not provided, allow any.
					"nfServices": bson.M{
						"$elemMatch": bson.M{
							"allowedNfDomains": bson.M{
								"$exists": false,
							},
						},
					},
				},
			},
		}
		filter["$and"] = append(filter["$and"].([]bson.M), requesterNfinstanceFqdnFilter)
	}

	// [Query-5] target-plmn-list [C] = Mcc + Mnc
	// Mcc: Pattern: '^[0-9]{3}$'
	// Mnc: Pattern: '^[0-9]{2,3}$'
	if queryParameters["target-plmn-list"] != nil {
		targetPlmnList := queryParameters["target-plmn-list"][0]
		targetPlmnListSplit := strings.Split(targetPlmnList, ",")
		var targetPlmnListBsonArray bson.A

		var temptargetPlmn string
		for i, v := range targetPlmnListSplit {
			if i%2 == 0 {
				temptargetPlmn = v
			} else {
				temptargetPlmn += ","
				temptargetPlmn += v

				targetPlmnListtruct := &models.PlmnId{}
				err := json.Unmarshal([]byte(temptargetPlmn), targetPlmnListtruct)
				if err != nil {
					logger.DiscoveryLog.Warnln("Unmarshal Error in targetPlmnListtruct: ", err)
				}

				targetPlmnByteArray, err := bson.Marshal(targetPlmnListtruct)
				if err != nil {
					logger.DiscoveryLog.Warnln("Unmarshal Error in targetPlmnListtruct: ", err)
				}

				targetPlmnBsonM := bson.M{}
				err = bson.Unmarshal(targetPlmnByteArray, &targetPlmnBsonM)
				if err != nil {
					logger.DiscoveryLog.Warnln("Unmarshal Error in targetPlmnBsonM: ", err)
				}
				logger.DiscoveryLog.Traceln("temp target Plmn:", temptargetPlmn)

				targetPlmnListBsonArray = append(targetPlmnListBsonArray, bson.M{"plmnList": bson.M{"$elemMatch": targetPlmnBsonM}})
			}
		}

		targetPlmnListFilter := bson.M{
			"$or": targetPlmnListBsonArray,
		}

		filter["$and"] = append(filter["$and"].([]bson.M), targetPlmnListFilter)
	}

	// [Query-6] requester-plmn-list
	// if queryParameters["requester-plmn-list"] != nil {
	// requesterPlmnPist := queryParameters["requester-plmn-list"][0]
	// TODO
	// }

	// [Query-7] target-nf-instance-id
	if queryParameters["target-nf-instance-id"] != nil {
		targetNfInstanceid := queryParameters["target-nf-instance-id"][0]
		nfInstanceIdFilter := bson.M{
			"nfInstanceId": targetNfInstanceid,
		}
		filter["$and"] = append(filter["$and"].([]bson.M), nfInstanceIdFilter)
	}

	// [Query-8] target-nf-fqdn
	if queryParameters["target-nf-fqdn"] != nil {
		targetNfFqdn := queryParameters["target-nf-fqdn"][0]
		fqdnFilter := bson.M{
			"fqdn": targetNfFqdn,
		}
		filter["$and"] = append(filter["$and"].([]bson.M), fqdnFilter)
	}

	// [Query-9] hnrf-uri
	// for Roaming

	// [Query-10] snssais
	// Pattern: '^[A-Fa-f0-9]{6}$'
	if queryParameters["snssais"] != nil {
		snssais := queryParameters["snssais"][0]
		snssaisSplit := strings.Split(snssais, ",")
		var snssaisBsonArray bson.A

		var tempSnssai string
		for i, v := range snssaisSplit {
			if i%2 == 0 {
				tempSnssai = v
			} else {
				tempSnssai += ","
				tempSnssai += v

				snssaiStruct := &models.Snssai{}
				err := json.Unmarshal([]byte(tempSnssai), snssaiStruct)
				if err != nil {
					logger.DiscoveryLog.Warnln("Unmarshal Error in snssaiStruct", err)
				}

				snssaiByteArray, err := bson.Marshal(snssaiStruct)
				if err != nil {
					logger.DiscoveryLog.Warnln("Unmarshal Error in snssaiStruct", err)
				}

				snssaiBsonM := bson.M{}
				err = bson.Unmarshal(snssaiByteArray, &snssaiBsonM)
				if err != nil {
					logger.DiscoveryLog.Warnln("Unmarshal Error in snssaiBsonM", err)
				}

				snssaisBsonArray = append(snssaisBsonArray, bson.M{"sNssais": bson.M{"$elemMatch": snssaiBsonM}})
			}
		}

		// if not assign, serve all NF
		snssaisBsonArray = append(snssaisBsonArray, bson.M{"sNssais": bson.M{"$exists": false}})

		snssaisFilter := bson.M{
			"$or": snssaisBsonArray,
		}

		filter["$and"] = append(filter["$and"].([]bson.M), snssaisFilter)
	}

	// [Query-11] nsi-list
	if queryParameters["nsi-list"] != nil {
		nsiList := queryParameters["nsi-list"][0]
		nsiListSplit := strings.Split(nsiList, ",")
		var nsiListBsonArray bson.A
		for _, v := range nsiListSplit {
			nsiListBsonArray = append(nsiListBsonArray, v)
		}
		nsiListFilter := bson.M{
			"nsiList": bson.M{
				"$all": nsiListBsonArray,
			},
		}
		filter["$and"] = append(filter["$and"].([]bson.M), nsiListFilter)
	}

	// [Query-12] dnn
	if queryParameters["dnn"] != nil {
		dnn := queryParameters["dnn"][0]
		var dnnFilter bson.M
		if targetNfType == "SMF" {
			dnnFilter = bson.M{
				"smfInfo.sNssaiSmfInfoList": bson.M{
					"$elemMatch": bson.M{
						"dnnSmfInfoList": bson.M{
							"$elemMatch": bson.M{
								"dnn": dnn,
							},
						},
					},
				},
			}
		} else if targetNfType == "UPF" {
			dnnFilter = bson.M{
				"upfInfo.sNssaiUpfInfoList": bson.M{
					"$elemMatch": bson.M{
						"dnnUpfInfoList": bson.M{
							"$elemMatch": bson.M{
								"dnn": dnn,
							},
						},
					},
				},
			}
		} else if targetNfType == "BSF" {
			dnnFilter = bson.M{
				"$or": []bson.M{
					{
						"bsfInfo.dnnList": dnn,
					},
					{
						"bsfInfo.dnnList": bson.M{
							"$exists": false,
						},
					},
				},
			}
		} else if targetNfType == "PCF" {
			dnnFilter = bson.M{
				"$or": []bson.M{
					{
						"pcfInfo.dnnList": dnn,
					},
					{
						"pcfInfo.dnnList": bson.M{
							"$exists": false,
						},
					},
				},
			}
		}
		filter["$and"] = append(filter["$and"].([]bson.M), dnnFilter)
	}

	// [Query-13] smf-serving-area
	if queryParameters["smf-serving-area"] != nil {
		var smfServingAreaFilter bson.M
		smfServingArea := queryParameters["smf-serving-area"][0]
		if targetNfType == "UPF" {
			smfServingAreaFilter = bson.M{
				"$or": []bson.M{
					{
						"upfInfo.smfServingArea": smfServingArea,
					},
					{
						"upfInfo.smfServingArea": bson.M{
							"$exists": false,
						},
					},
				},
			}
		}
		filter["$and"] = append(filter["$and"].([]bson.M), smfServingAreaFilter)
	}

	// [Query-14] tai
	if queryParameters["tai"] != nil {
		var taiFilter bson.M
		tai := queryParameters["tai"][0]

		taiStruct := &models.Tai{}
		err := json.Unmarshal([]byte(tai), taiStruct)
		if err != nil {
			logger.DiscoveryLog.Warnln("Unmarshal Error in taiStruct: ", err)
		}

		taiByteArray, err := bson.Marshal(taiStruct)
		if err != nil {
			logger.DiscoveryLog.Warnln("Unmarshal Error in taiByteArray: ", err)
		}

		taiBsonM := bson.M{}
		err = bson.Unmarshal(taiByteArray, &taiBsonM)
		if err != nil {
			logger.DiscoveryLog.Warnln("Unmarshal Error in taiByteArray: ", err)
		}
		if targetNfType == "SMF" {
			taiFilter = bson.M{
				"smfInfo.taiList": bson.M{
					"$elemMatch": taiBsonM,
				},
			}
		} else if targetNfType == "AMF" {
			taiFilter = bson.M{
				"amfInfo.taiList": bson.M{
					"$elemMatch": taiBsonM,
				},
			}
		}
		filter["$and"] = append(filter["$and"].([]bson.M), taiFilter)
	}

	// [Query-15] amf-region-id
	if queryParameters["amf-region-id"] != nil {
		if targetNfType == "AMF" {
			amfRegionId := queryParameters["amf-region-id"][0]
			amfRegionIdFilter := bson.M{
				"amfInfo.amfRegionId": amfRegionId,
			}
			filter["$and"] = append(filter["$and"].([]bson.M), amfRegionIdFilter)
		}
	}

	// [Query-16] amf-set-id
	if queryParameters["amf-set-id"] != nil {
		if targetNfType == "AMF" {
			amfSetId := queryParameters["amf-set-id"][0]
			amfSetIdFilter := bson.M{
				"amfInfo.amfSetId": amfSetId,
			}
			filter["$and"] = append(filter["$and"].([]bson.M), amfSetIdFilter)
		}
	}

	// Query-17: guami
	// TODO: NOTE[1]
	if queryParameters["guami"] != nil {
		if targetNfType == "AMF" {
			guami := queryParameters["guami"][0]

			guamiStruct := &models.Guami{}
			err := json.Unmarshal([]byte(guami), guamiStruct)
			if err != nil {
				logger.DiscoveryLog.Warnln("Unmarshal Error in guamiStruct: ", err)
			}

			guamiByteArray, err := bson.Marshal(guamiStruct)
			if err != nil {
				logger.DiscoveryLog.Warnln("Unmarshal Error in guamiByteArray: ", err)
			}

			guamiBsonM := bson.M{}
			err = bson.Unmarshal(guamiByteArray, &guamiBsonM)
			if err != nil {
				logger.DiscoveryLog.Warnln("Unmarshal Error in guamiByteArray: ", err)
			}

			guamiFilter := bson.M{
				"amfInfo.guamiList": bson.M{
					"$elemMatch": guamiBsonM,
				},
			}

			filter["$and"] = append(filter["$and"].([]bson.M), guamiFilter)
		}
	}

	// [Query-18] supi
	var supi string
	if queryParameters["supi"] != nil {
		var supiFilter bson.M
		supi = queryParameters["supi"][0]
		supi = supi[5:]
		if targetNfType == "PCF" {
			supiFilter = bson.M{
				"$or": []bson.M{
					{
						"pcfInfo.supiRanges": bson.M{
							"$elemMatch": bson.M{
								"start": bson.M{
									"$lte": supi,
								},
								"end": bson.M{
									"$gte": supi,
								},
							},
						},
					},
					{
						"pcfInfo.supiRanges": bson.M{
							"$exists": false,
						},
					},
				},
			}
		} else if targetNfType == "CHF" {
			supiFilter = bson.M{
				"$or": []bson.M{
					{
						"chfInfo.supiRangeList": bson.M{
							"$elemMatch": bson.M{
								"start": bson.M{
									"$lte": supi,
								},
								"end": bson.M{
									"$gte": supi,
								},
							},
						},
					},
					{
						"chfInfo.supiRangeList": bson.M{
							"$exists": false,
						},
					},
				},
			}
		} else if targetNfType == "AUSF" {
			supiFilter = bson.M{
				"$or": []bson.M{
					{
						"ausfInfo.supiRanges": bson.M{
							"$elemMatch": bson.M{
								"start": bson.M{
									"$lte": supi,
								},
								"end": bson.M{
									"$gte": supi,
								},
							},
						},
					},
					{
						"ausfInfo.supiRanges": bson.M{
							"$exists": false,
						},
					},
				},
			}
		} else if targetNfType == "UDM" {
			supiFilter = bson.M{
				"$or": []bson.M{
					{
						"udmInfo.supiRanges": bson.M{
							"$elemMatch": bson.M{
								"start": bson.M{
									"$lte": supi,
								},
								"end": bson.M{
									"$gte": supi,
								},
							},
						},
					},
					{
						"udmInfo.supiRanges": bson.M{
							"$exists": false,
						},

						"udmInfo.gpsiRanges": bson.M{
							"$exists": false,
						},

						"udmInfo.externalGroupIdentifiersRanges": bson.M{
							"$exists": false,
						},
					},
				},
			}
		} else if targetNfType == "UDR" {
			supiFilter = bson.M{
				"$or": []bson.M{
					{
						"udrInfo.supiRanges": bson.M{
							"$elemMatch": bson.M{
								"start": bson.M{
									"$lte": supi,
								},
								"end": bson.M{
									"$gte": supi,
								},
							},
						},
					},
					{
						"udrInfo.supiRanges": bson.M{
							"$exists": false,
						},

						"udrInfo.gpsiRanges": bson.M{
							"$exists": false,
						},

						"udrInfo.externalGroupIdentifiersRanges": bson.M{
							"$exists": false,
						},
					},
				},
			}
		}
		filter["$and"] = append(filter["$and"].([]bson.M), supiFilter)
	}

	// [Query-19] ue-ipv4-address
	if queryParameters["ue-ipv4-address"] != nil {
		var ueIpv4AddressFilter bson.M
		if targetNfType == "BSF" {
			ueIpv4Address := queryParameters["ue-ipv4-address"][0]
			ueIpv4AddressNumber := context.Ipv4ToInt(ueIpv4Address)
			ueIpv4AddressFilter = bson.M{
				"$or": []bson.M{
					{
						"bsfInfo.ipv4AddressRanges": bson.M{
							"$elemMatch": bson.M{
								"start": bson.M{
									"$lte": strconv.Itoa(int(ueIpv4AddressNumber)),
								},
								"end": bson.M{
									"$gte": strconv.Itoa(int(ueIpv4AddressNumber)),
								},
							},
						},
					},
					{
						"bsfInfo.ipv4AddressRanges": bson.M{
							"$exists": false,
						},
					},
				},
			}
		}
		filter["$and"] = append(filter["$and"].([]bson.M), ueIpv4AddressFilter)
	}

	// [Query-20] ip-domain
	if queryParameters["ip-domain"] != nil {
		var ipDomainFilter bson.M
		if targetNfType == "BSF" {
			ipDomain := queryParameters["ip-domain"][0]
			ipDomainFilter = bson.M{
				"$or": []bson.M{
					{
						"bsfInfo.ipDomainList": ipDomain,
					},
					{
						"bsfInfo.ipDomainList": bson.M{
							"$exists": false,
						},
					},
				},
			}
		}
		filter["$and"] = append(filter["$and"].([]bson.M), ipDomainFilter)
	}

	// [Query-21] ue-ipv6-prefix
	if queryParameters["ue-ipv6-prefix"] != nil {
		var ueIpv6PrefixFilter bson.M
		if targetNfType == "BSF" {
			ueIpv6Prefix := queryParameters["ue-ipv6-prefix"][0]
			ueIpv6PrefixNumber := context.Ipv6ToInt(ueIpv6Prefix)
			ueIpv6PrefixFilter = bson.M{
				"$or": []bson.M{
					{
						"bsfInfo.ipv6PrefixRanges": bson.M{
							"$elemMatch": bson.M{
								"start": bson.M{
									"$lte": ueIpv6PrefixNumber.String(),
								},
								"end": bson.M{
									"$gte": ueIpv6PrefixNumber.String(),
								},
							},
						},
					},
					{
						"bsfInfo.ipv6PrefixRanges": bson.M{
							"$exists": false,
						},
					},
				},
			}
		}
		filter["$and"] = append(filter["$and"].([]bson.M), ueIpv6PrefixFilter)
	}

	// [Query-22] pgw-ind
	if queryParameters["pgw-ind"] != nil {
		pgwInd := queryParameters["pgw-ind"][0]
		if pgwInd == "true" {
			pgwIndFilter := bson.M{
				"smfInfo.pgwFqdn": bson.M{
					"$exists": true,
				},
			}
			filter["$and"] = append(filter["$and"].([]bson.M), pgwIndFilter)
		}
	}

	// [Query-23] pgw
	if queryParameters["pgw"] != nil {
		pgw := queryParameters["pgw"][0]
		pgwFilter := bson.M{
			"smfInfo.pgwFqdn": pgw,
		}
		filter["$and"] = append(filter["$and"].([]bson.M), pgwFilter)
	}

	// [Query-24] gpsi
	if queryParameters["gpsi"] != nil {
		var gpsiFilter bson.M
		gpsi := queryParameters["gpsi"][0]
		gpsi = gpsi[7:]
		if targetNfType == "CHF" {
			gpsiFilter = bson.M{
				"$or": []bson.M{
					{
						"chfInfo.gpsiRangeList": bson.M{
							"$elemMatch": bson.M{
								"start": bson.M{
									"$lte": gpsi,
								},
								"end": bson.M{
									"$gte": gpsi,
								},
							},
						},
					},
					{
						"chfInfo.gpsiRangeList": bson.M{
							"$exists": false,
						},
					},
				},
			}
		} else if targetNfType == "UDM" {
			gpsiFilter = bson.M{
				"$or": []bson.M{
					{
						"udmInfo.gpsiRanges": bson.M{
							"$elemMatch": bson.M{
								"start": bson.M{
									"$lte": gpsi,
								},
								"end": bson.M{
									"$gte": gpsi,
								},
							},
						},
					},
					{
						"udmInfo.supiRanges": bson.M{
							"$exists": false,
						},

						"udmInfo.gpsiRanges": bson.M{
							"$exists": false,
						},

						"udmInfo.externalGroupIdentifiersRanges": bson.M{
							"$exists": false,
						},
					},
				},
			}
		} else if targetNfType == "UDR" {
			gpsiFilter = bson.M{
				"$or": []bson.M{
					{
						"udrInfo.gpsiRanges": bson.M{
							"$elemMatch": bson.M{
								"start": bson.M{
									"$lte": gpsi,
								},
								"end": bson.M{
									"$gte": gpsi,
								},
							},
						},
					},
					{
						"udrInfo.supiRanges": bson.M{
							"$exists": false,
						},

						"udrInfo.gpsiRanges": bson.M{
							"$exists": false,
						},

						"udrInfo.externalGroupIdentifiersRanges": bson.M{
							"$exists": false,
						},
					},
				},
			}
		}
		filter["$and"] = append(filter["$and"].([]bson.M), gpsiFilter)
	}

	// [Query-25] external-group-identity
	if queryParameters["external-group-identity"] != nil {
		var externalGroupIdentityFilter bson.M
		externalGroupIdentity := queryParameters["external-group-identity"][0]

		encodedGroupId := context.EncodeGroupId(externalGroupIdentity)

		if targetNfType == "UDM" {
			externalGroupIdentityFilter = bson.M{
				"$or": []bson.M{
					{
						"udmInfo.externalGroupIdentifiersRanges": bson.M{
							"$elemMatch": bson.M{
								"start": bson.M{
									"$lte": encodedGroupId,
								},
								"end": bson.M{
									"$gte": encodedGroupId,
								},
							},
						},
					},
					{
						"udmInfo.supiRanges": bson.M{
							"$exists": false,
						},

						"udmInfo.gpsiRanges": bson.M{
							"$exists": false,
						},

						"udmInfo.externalGroupIdentifiersRanges": bson.M{
							"$exists": false,
						},
					},
				},
			}
		} else if targetNfType == "UDR" {
			externalGroupIdentityFilter = bson.M{
				"$or": []bson.M{
					{
						"udrInfo.externalGroupIdentifiersRanges": bson.M{
							"$elemMatch": bson.M{
								"start": bson.M{
									"$lte": encodedGroupId,
								},
								"end": bson.M{
									"$gte": encodedGroupId,
								},
							},
						},
					},
					{
						"udrInfo.supiRanges": bson.M{
							"$exists": false,
						},

						"udrInfo.gpsiRanges": bson.M{
							"$exists": false,
						},

						"udrInfo.externalGroupIdentifiersRanges": bson.M{
							"$exists": false,
						},
					},
				},
			}
		}
		filter["$and"] = append(filter["$and"].([]bson.M), externalGroupIdentityFilter)
	}

	// [Query-26] data-set
	if queryParameters["data-set"] != nil {
		var dataSetFilter bson.M
		dataSet := queryParameters["data-set"]
		if targetNfType == "UDR" {
			dataSetFilter = bson.M{
				"$or": []bson.M{
					{
						"udrInfo.supportedDataSets": dataSet,
					},
					{
						"udrInfo.supportedDataSets": bson.M{
							"$exists": false,
						},
					},
				},
			}
		}
		filter["$and"] = append(filter["$and"].([]bson.M), dataSetFilter)
	}

	// [Query-27] routing-indicator
	if queryParameters["routing-indicator"] != nil {
		var routingIndicatorFilter bson.M
		routingIndicator := queryParameters["routing-indicator"][0]
		if targetNfType == "AUSF" {
			routingIndicatorFilter = bson.M{
				"$or": []bson.M{
					{
						"ausfInfo.routingIndicators": routingIndicator,
					},
					{
						"ausfInfo.routingIndicators": bson.M{
							"$exists": false,
						},
					},
				},
			}
		} else if targetNfType == "UDM" {
			routingIndicatorFilter = bson.M{
				"$or": []bson.M{
					{
						"udmInfo.routingIndicators": routingIndicator,
					},
					{
						"udmInfo.routingIndicators": bson.M{
							"$exists": false,
						},
					},
				},
			}
		}
		filter["$and"] = append(filter["$and"].([]bson.M), routingIndicatorFilter)
	}

	// [Query-28] group-id-list
	if queryParameters["group-id-list"] != nil {
		var groupIdListFilter bson.M

		groupIdList := queryParameters["group-id-list"][0]
		groupIdListSplit := strings.Split(groupIdList, ",")
		var groupIdListBsonArray bson.A

		for _, v := range groupIdListSplit {
			groupIdListBsonArray = append(groupIdListBsonArray, v)
		}

		if targetNfType == "UDR" {
			groupIdListFilter = bson.M{
				"udrInfo.groupId": bson.M{
					"$in": groupIdListBsonArray,
				},
			}
		} else if targetNfType == "UDM" {
			groupIdListFilter = bson.M{
				"udmInfo.groupId": bson.M{
					"$in": groupIdListBsonArray,
				},
			}
		} else if targetNfType == "AUSF" {
			groupIdListFilter = bson.M{
				"ausfInfo.groupId": bson.M{
					"$in": groupIdListBsonArray,
				},
			}
		}
		filter["$and"] = append(filter["$and"].([]bson.M), groupIdListFilter)
	}

	// [Query-29] dnai-list
	if queryParameters["dnai-list"] != nil {
		var dnaiFilter bson.M
		dnaiList := queryParameters["dnai-list"][0]
		dnaiListSplit := strings.Split(dnaiList, ",")
		var dnaiListBsonArray bson.A

		for _, v := range dnaiListSplit {
			dnaiListBsonArray = append(dnaiListBsonArray, v)
		}
		if targetNfType == "UPF" {
			dnaiFilter = bson.M{
				"upfInfo.sNssaiUpfInfoList": bson.M{
					"$elemMatch": bson.M{
						"dnnUpfInfoList": bson.M{
							"$elemMatch": bson.M{
								"dnaiList": bson.M{
									"$in": dnaiListBsonArray,
								},
							},
						},
					},
				},
			}
		}
		filter["$and"] = append(filter["$and"].([]bson.M), dnaiFilter)
	}

	// [Query-30] upf-iwk-eps-ind
	if queryParameters["upf-iwk-eps-ind"] != nil {
		var upfIwkEpsIndFilter bson.M
		// upfIwkEpsInd := queryParameters["upf-iwk-eps-ind"][0]
		if targetNfType == "UPF" {
			upfIwkEpsIndFilter = bson.M{
				"upfInfo.iwkEpsInd": true,
			}
		}
		filter["$and"] = append(filter["$and"].([]bson.M), upfIwkEpsIndFilter)
	}

	// [Query-31] chf-supported-plmn
	if queryParameters["chf-supported-plmn"] != nil {
		var chfSupportedPlmnFilter bson.M
		chfSupportedPlmn := queryParameters["chf-supported-plmn"][0]
		chfSupportedPlmnStruct := &models.PlmnId{}
		err := json.Unmarshal([]byte(chfSupportedPlmn), chfSupportedPlmnStruct)
		if err != nil {
			logger.DiscoveryLog.Warnln("Unmarshal Error in chfSupportedPlmnStruct: ", err)
		}

		encodedchfSupportedPlmn := chfSupportedPlmnStruct.Mcc + chfSupportedPlmnStruct.Mnc

		if targetNfType == "CHF" {
			chfSupportedPlmnFilter = bson.M{
				"$or": []bson.M{
					{
						"chfInfo.plmnRangeList": bson.M{
							"$elemMatch": bson.M{
								"start": bson.M{
									"$lte": encodedchfSupportedPlmn,
								},
								"end": bson.M{
									"$gte": encodedchfSupportedPlmn,
								},
							},
						},
					},
					{
						"chfInfo.plmnRangeList": bson.M{
							"$exists": false,
						},
					},
				},
			}
		}
		filter["$and"] = append(filter["$and"].([]bson.M), chfSupportedPlmnFilter)
	}

	// [Query-32]  preferred-locality
	// TODO: if no match
	if queryParameters["preferred-locality"] != nil {
		preferredLocality := queryParameters["preferred-locality"][0]
		preferredLocalityFilter := bson.M{
			"locality": preferredLocality,
		}
		filter["$and"] = append(filter["$and"].([]bson.M), preferredLocalityFilter)
	}

	// [Query-33] access-type
	if queryParameters["access-type"] != nil {
		accessType := queryParameters["access-type"][0]
		accessTypeFilter := bson.M{
			"$or": []bson.M{
				{
					"smfInfo.accessType": accessType,
				},
				{
					"smfInfo.accessType": bson.M{
						"$exists": false,
					},
				},
			},
		}
		filter["$and"] = append(filter["$and"].([]bson.M), accessTypeFilter)
	}

	// [Query-34] supported-features
	if queryParameters["supported-features"] != nil {
		supportedFeatures := queryParameters["supported-features"][0]
		supportedFeaturesFilter := bson.M{
			"nfServices": bson.M{
				"$elemMatch": bson.M{
					"supportedFeatures": supportedFeatures,
				},
			},
		}
		filter["$and"] = append(filter["$and"].([]bson.M), supportedFeaturesFilter)
	}

	// [Query-35] complexQuery
	if queryParameters["complexQuery"] != nil {
		// translate raw data to complexQuery structure
		complexQuery := queryParameters["complexQuery"][0]
		complexQueryStruct := &models.ComplexQuery{}
		err := json.Unmarshal([]byte(complexQuery), complexQueryStruct)
		if err != nil {
			logger.DiscoveryLog.Warnln("Unmarshal Error in complexQuery: ", err)
		}
		complexQueryFilter := complexQueryFilter(complexQueryStruct)
		filter["$and"] = append(filter["$and"].([]bson.M), complexQueryFilter)
	}
	return filter
}

const (
	COMPLEX_QUERY_TYPE_CNF string = "CNF"
	COMPLEX_QUERY_TYPE_DNF string = "DNF"
)

type AtomElem struct {
	value    string
	negative bool
}

func complexQueryFilter(complexQueryParameter *models.ComplexQuery) bson.M {
	complexQueryType := ""
	if complexQueryParameter.CNf != nil {
		complexQueryType = COMPLEX_QUERY_TYPE_CNF
	} else {
		complexQueryType = COMPLEX_QUERY_TYPE_DNF
	}

	// build the filter
	var filter bson.M

	if complexQueryType == COMPLEX_QUERY_TYPE_CNF {
		filter = bson.M{
			"$and": []bson.M{},
		}
		for _, cnfUnit := range complexQueryParameter.CNf.CnfUnits {
			var queryParameters map[string]*AtomElem = make(map[string]*AtomElem)
			var cnfUnitFilter bson.M
			for _, atom := range cnfUnit.CnfUnit {
				queryParameters[atom.Attr] = &AtomElem{value: atom.Value, negative: atom.Negative}
			}
			cnfUnitFilter = complexQueryFilterSubprocess(queryParameters, complexQueryType)

			filter["$and"] = append(filter["$and"].([]bson.M), cnfUnitFilter)
		}
	} else {
		filter = bson.M{
			"$or": []bson.M{},
		}
	}
	return filter
}

func complexQueryFilterSubprocess(queryParameters map[string]*AtomElem, complexQueryType string) bson.M {
	var filter bson.M
	var logicalOperator string

	switch complexQueryType {
	case COMPLEX_QUERY_TYPE_CNF:
		logicalOperator = "$or"
	case COMPLEX_QUERY_TYPE_DNF:
		logicalOperator = "$and"
	}

	filter = bson.M{
		logicalOperator: []bson.M{},
	}

	// [Query-1] target-nf-type
	var targetNfType string
	if targetNfType != "" {
		var targetNfTypeFilter bson.M
		targetNfType = queryParameters["target-nf-type"].value
		negative := queryParameters["target-nf-type"].negative
		if negative {
			targetNfTypeFilter = bson.M{
				"nfType": bson.M{
					"$ne": targetNfType,
				},
			}
		} else if !negative {
			targetNfTypeFilter = bson.M{
				"nfType": targetNfType,
			}
		}
		filter[logicalOperator] = append(filter[logicalOperator].([]bson.M), targetNfTypeFilter)
	}

	// [Query-2] requester-nf-type
	// requesterNfType := queryParameters["requester-nf-type"].value
	// TODO

	// [Query-3] service-names
	// TODO: return exist service name
	if queryParameters["service-names"] != nil {
		var serviceNamesFilter bson.M
		serviceNames := queryParameters["service-names"].value
		serviceNamesSplit := strings.Split(serviceNames, ",")
		var serviceNamesBsonArray bson.A

		for _, v := range serviceNamesSplit {
			serviceNamesBsonArray = append(serviceNamesBsonArray, v)
		}

		negative := queryParameters["service-names"].negative
		if negative {
			serviceNamesFilter = bson.M{
				"nfServices": bson.M{
					"$elemMatch": bson.M{
						"serviceName": bson.M{
							// get all service in array
							"$nin": serviceNamesBsonArray,
						},
						// the service need to be registered
						"nfServiceStatus": "REGISTERED",
					},
				},
			}
		} else if !negative {
			serviceNamesFilter = bson.M{
				"nfServices": bson.M{
					"$elemMatch": bson.M{
						"serviceName": bson.M{
							// get all service in array
							"$in": serviceNamesBsonArray,
						},
						// the service need to be registered
						"nfServiceStatus": "REGISTERED",
					},
				},
			}
		}
		filter[logicalOperator] = append(filter[logicalOperator].([]bson.M), serviceNamesFilter)
	}

	// [Query-4] requester-nfinstance-fqdn
	if queryParameters["requester-nfinstance-fqdn"] != nil {
		var requesterNfinstanceFqdnFilter bson.M
		requesterNfinstanceFqdn := queryParameters["requester-nfinstance-fqdn"].value

		negative := queryParameters["requester-nfinstance-fqdn"].negative
		if negative {
			requesterNfinstanceFqdnFilter = bson.M{
				"nfServices": bson.M{
					"$elemMatch": bson.M{
						"allowedNfDomains": requesterNfinstanceFqdn,
					},
				},
			}
		} else if !negative {
			requesterNfinstanceFqdnFilter = bson.M{
				"nfServices": bson.M{
					"$elemMatch": bson.M{
						"allowedNfDomains": bson.M{
							"$ne": requesterNfinstanceFqdn,
						},
					},
				},
			}
		}
		filter[logicalOperator] = append(filter[logicalOperator].([]bson.M), requesterNfinstanceFqdnFilter)
	}

	// [Query-5] target-plmn-list [C] = Mcc + Mnc
	// Mcc: Pattern: '^[0-9]{3}$'
	// Mnc: Pattern: '^[0-9]{2,3}$'
	if queryParameters["target-plmn-list"] != nil {
		targetPlmnList := queryParameters["target-plmn-list"].value
		targetPlmnListSplit := strings.Split(targetPlmnList, ",")
		var targetPlmnListBsonArray bson.A

		var temptargetPlmn string
		for i, v := range targetPlmnListSplit {
			if i%2 == 0 {
				temptargetPlmn = v
			} else {
				temptargetPlmn += ","
				temptargetPlmn += v

				targetPlmnListtruct := &models.PlmnId{}
				err := json.Unmarshal([]byte(temptargetPlmn), targetPlmnListtruct)
				if err != nil {
					logger.DiscoveryLog.Warnln("Unmarshal Error in targetPlmnListstruct: ", err)
				}

				targetPlmnByteArray, err := bson.Marshal(targetPlmnListtruct)
				if err != nil {
					logger.DiscoveryLog.Warnln("Unmarshal Error in targetPlmnByteArray: ", err)
				}

				targetPlmnBsonM := bson.M{}
				err = bson.Unmarshal(targetPlmnByteArray, &targetPlmnBsonM)
				if err != nil {
					logger.DiscoveryLog.Warnln("Unmarshal Error in targetPlmnBsonM: ", err)
				}

				targetPlmnListBsonArray = append(targetPlmnListBsonArray, targetPlmnBsonM)
			}
		}

		var targetPlmnListFilter bson.M
		negative := queryParameters["target-plmn-list"].negative
		if negative {
			targetPlmnListFilter = bson.M{
				"PlmnList": bson.M{
					"$nin": targetPlmnListBsonArray,
				},
			}
		} else if !negative {
			targetPlmnListFilter = bson.M{
				"PlmnList": bson.M{
					"$in": targetPlmnListBsonArray,
				},
			}
		}
		filter[logicalOperator] = append(filter[logicalOperator].([]bson.M), targetPlmnListFilter)
	}

	// [Query-6] requester-plmn-list
	// if queryParameters["requester-plmn-list"] != nil {
	// requesterPlmnPist := queryParameters["requester-plmn-list"].value
	// TODO
	// }

	// [Query-7] target-nf-instanceid
	if queryParameters["target-nf-instanceid"] != nil {
		targetNfInstanceid := queryParameters["target-nf-instanceid"].value
		var nfInstanceIdFilter bson.M

		negative := queryParameters["target-nf-instanceid"].negative
		if negative {
			nfInstanceIdFilter = bson.M{
				"nfInstanceId": bson.M{
					"$ne": targetNfInstanceid,
				},
			}
		} else if !negative {
			nfInstanceIdFilter = bson.M{
				"nfInstanceId": targetNfInstanceid,
			}
		}
		filter[logicalOperator] = append(filter[logicalOperator].([]bson.M), nfInstanceIdFilter)
	}

	// [Query-8] target-nf-fqdn
	if queryParameters["target-nf-fqdn"] != nil {
		targetNfFqdn := queryParameters["target-nf-fqdn"].value
		fqdnFilter := bson.M{
			"fqdn": targetNfFqdn,
		}
		if queryParameters["target-nf-fqdn"].negative {
			fqdnFilter = bson.M{
				"$not": fqdnFilter,
			}
		}
		filter[logicalOperator] = append(filter[logicalOperator].([]bson.M), fqdnFilter)
	}

	// [Query-9] hnrf-uri
	// for Roaming

	// [Query-10] snssais
	// Pattern: '^[A-Fa-f0-9]{6}$'
	if queryParameters["snssais"] != nil {
		snssais := queryParameters["snssais"].value
		snssaisSplit := strings.Split(snssais, ",")
		var snssaisBsonArray bson.A

		var tempSnssai string
		for i, v := range snssaisSplit {
			if i%2 == 0 {
				tempSnssai = v
			} else {
				tempSnssai += ","
				tempSnssai += v

				snssaiStruct := &models.Snssai{}
				err := json.Unmarshal([]byte(tempSnssai), snssaiStruct)
				if err != nil {
					logger.DiscoveryLog.Warnln("Unmarshal Error in snssaiStruct: ", err)
				}

				snssaiByteArray, err := bson.Marshal(snssaiStruct)
				if err != nil {
					logger.DiscoveryLog.Warnln("Unmarshal Error in snssaiByteArray: ", err)
				}

				snssaiBsonM := bson.M{}
				err = bson.Unmarshal(snssaiByteArray, &snssaiBsonM)
				if err != nil {
					logger.DiscoveryLog.Warnln("Unmarshal Error in snssaiBsonM: ", err)
				}

				snssaisBsonArray = append(snssaisBsonArray, snssaiBsonM)
			}
		}

		snssaisFilter := bson.M{
			"snssais": bson.M{
				"$elemMatch": snssaisBsonArray,
			},
		}
		if queryParameters["snssais"].negative {
			snssaisFilter = bson.M{
				"$not": snssaisFilter,
			}
		}
		filter[logicalOperator] = append(filter[logicalOperator].([]bson.M), snssaisFilter)
	}

	// [Query-11] nsi-list
	if queryParameters["nsi-list"] != nil {
		nsiList := queryParameters["nsi-list"].value
		nsiListSplit := strings.Split(nsiList, ",")
		var nsiListBsonArray bson.A
		for _, v := range nsiListSplit {
			nsiListBsonArray = append(nsiListBsonArray, v)
		}
		nsiListFilter := bson.M{
			"nsiList": bson.M{
				"$all": nsiListBsonArray,
			},
		}
		if queryParameters["nsi-list"].negative {
			nsiListFilter = bson.M{
				"$not": nsiListFilter,
			}
		}
		filter[logicalOperator] = append(filter[logicalOperator].([]bson.M), nsiListFilter)
	}

	// [Query-12] dnn
	if queryParameters["dnn"] != nil {
		dnn := queryParameters["dnn"].value
		var dnnFilter bson.M
		if targetNfType == "SMF" {
			dnnFilter = bson.M{
				"smfInfo": bson.M{
					"$elemMatch": bson.M{
						"sNssaiSmfInfoList": bson.M{
							"$elemMatch": bson.M{
								"dnnSmfInfoList": bson.M{
									"$elemMatch": bson.M{
										"dnn": dnn[0],
									},
								},
							},
						},
					},
				},
			}
		} else if targetNfType == "UPF" {
			dnnFilter = bson.M{
				"upfInfo": bson.M{
					"$elemMatch": bson.M{
						"sNssaiUpfInfoList": bson.M{
							"$elemMatch": bson.M{
								"dnnUpfInfoList": bson.M{
									"$elemMatch": bson.M{
										"dnn": dnn,
									},
								},
							},
						},
					},
				},
			}
		} else if targetNfType == "BSF" {
			dnnFilter = bson.M{
				"bsfInfo": bson.M{
					"$elemMatch": bson.M{
						"dnnList": dnn[0],
					},
				},
			}
		}
		if queryParameters["dnn"].negative {
			dnnFilter = bson.M{
				"$not": dnnFilter,
			}
		}
		filter[logicalOperator] = append(filter[logicalOperator].([]bson.M), dnnFilter)
	}

	// [Query-13] smf-serving-area
	if queryParameters["smf-serving-area"] != nil {
		var smfServingAreaFilter bson.M
		smfServingArea := queryParameters["smf-serving-area"].value
		if targetNfType == "UPF" {
			smfServingAreaFilter = bson.M{
				"upfInfo": bson.M{
					"$elemMatch": bson.M{
						"smfServingArea": smfServingArea,
					},
				},
			}
		}
		if queryParameters["smf-serving-area"].negative {
			smfServingAreaFilter = bson.M{
				"$not": smfServingAreaFilter,
			}
		}
		filter[logicalOperator] = append(filter[logicalOperator].([]bson.M), smfServingAreaFilter)
	}

	// [Query-14] tai
	if queryParameters["tai"] != nil {
		var taiFilter bson.M
		tai := queryParameters["tai"].value
		taiSplit := strings.Split(tai, ",")
		tempTai := taiSplit[0] + "," + taiSplit[1]

		taiStruct := &models.Tai{}
		err := json.Unmarshal([]byte(tempTai), taiStruct)
		if err != nil {
			logger.DiscoveryLog.Warnln("Unmarshal Error in taiStruct: ", err)
		}

		taiByteArray, err := bson.Marshal(taiStruct)
		if err != nil {
			logger.DiscoveryLog.Warnln("Unmarshal Error in taiByteArray: ", err)
		}

		taiBsonM := bson.M{}
		err = bson.Unmarshal(taiByteArray, &taiBsonM)
		if err != nil {
			logger.DiscoveryLog.Warnln("Unmarshal Error in taiByteArray: ", err)
		}
		if targetNfType == "SMF" {
			taiFilter = bson.M{
				"smfInfo": bson.M{
					"$elemMatch": bson.M{
						"taiList": taiBsonM,
					},
				},
			}
		} else if targetNfType == "AMF" {
			taiFilter = bson.M{
				"amfInfo": bson.M{
					"$elemMatch": bson.M{
						"taiList": taiBsonM,
					},
				},
			}
		}
		if queryParameters["tai"].negative {
			taiFilter = bson.M{
				"$not": taiFilter,
			}
		}
		filter[logicalOperator] = append(filter[logicalOperator].([]bson.M), taiFilter)
	}

	// [Query-15] amf-region-id
	if queryParameters["amf-region-id"] != nil {
		var amfRegionIdFilter bson.M
		if targetNfType == "AMF" {
			amfRegionId := queryParameters["amf-region-id"].value
			amfRegionIdFilter = bson.M{
				"amfInfo": bson.M{
					"$elemMatch": bson.M{
						"amfRegionId": amfRegionId[0],
					},
				},
			}
		}
		if queryParameters["amf-region-id"].negative {
			amfRegionIdFilter = bson.M{
				"$not": amfRegionIdFilter,
			}
		}
		filter[logicalOperator] = append(filter[logicalOperator].([]bson.M), amfRegionIdFilter)
	}

	// [Query-16] amf-set-id
	if queryParameters["amf-set-id"] != nil {
		var amfSetIdFilter bson.M
		if targetNfType == "AMF" {
			amfSetId := queryParameters["amf-set-id"].value
			amfSetIdFilter = bson.M{
				"amfInfo": bson.M{
					"$elemMatch": bson.M{ // TOCHECK : elemMatch
						"amfSetId": amfSetId[0],
					},
				},
			}
		}
		if queryParameters["amf-set-id"].negative {
			amfSetIdFilter = bson.M{
				"$not": amfSetIdFilter,
			}
		}
		filter[logicalOperator] = append(filter[logicalOperator].([]bson.M), amfSetIdFilter)
	}

	// Query-17: guami
	// TODO: NOTE[1]
	if queryParameters["guami"] != nil {
		var guamiFilter bson.M
		if targetNfType == "AMF" {
			guami := queryParameters["guami"].value
			guamiSplit := strings.Split(guami, ",")
			tempguami := guamiSplit[0] + "," + guamiSplit[1]

			guamiStruct := &models.Guami{}
			err := json.Unmarshal([]byte(tempguami), guamiStruct)
			if err != nil {
				logger.DiscoveryLog.Warnln("Unmarshal Error in guamiStruct: ", err)
			}

			guamiByteArray, err := bson.Marshal(guamiStruct)
			if err != nil {
				logger.DiscoveryLog.Warnln("Unmarshal Error in guamiByteArray: ", err)
			}

			guamiBsonM := bson.M{}
			err = bson.Unmarshal(guamiByteArray, &guamiBsonM)
			if err != nil {
				logger.DiscoveryLog.Warnln("Unmarshal Error in guamiByteArray: ", err)
			}

			guamiFilter = bson.M{
				"amfInfo": bson.M{
					"$elemMatch": bson.M{ // TOCHECK : elemMatch
						"guamiList": bson.M{
							"$elemMatch": guamiBsonM,
						},
					},
				},
			}
		}
		if queryParameters["guami"].negative {
			guamiFilter = bson.M{
				"$not": guamiFilter,
			}
		}
		filter[logicalOperator] = append(filter[logicalOperator].([]bson.M), guamiFilter)
	}

	// [Query-18] supi
	var supi string
	if queryParameters["supi"] != nil {
		var supiFilter bson.M
		supi = queryParameters["supi"].value
		if targetNfType == "PCF" {
			supiFilter = bson.M{
				"pcfInfo": bson.M{
					"$elemMatch": bson.M{
						"supiRanges": bson.M{
							"$elemMatch": bson.M{
								"start": bson.M{
									"$lte": supi[0],
								},
								"end": bson.M{
									"$gte": supi[0],
								},
							},
						},
					},
				},
			}
		} else if targetNfType == "CHF" {
			supiFilter = bson.M{
				"chfInfo": bson.M{
					"$elemMatch": bson.M{
						"supiRanges": bson.M{
							"$elemMatch": bson.M{
								"start": bson.M{
									"$lte": supi[0],
								},
								"end": bson.M{
									"$gte": supi[0],
								},
							},
						},
					},
				},
			}
		} else if targetNfType == "AUSF" {
			supiFilter = bson.M{
				"ausfInfo": bson.M{
					"$elemMatch": bson.M{
						"supiRanges": bson.M{
							"$elemMatch": bson.M{
								"start": bson.M{
									"$lte": supi[0],
								},
								"end": bson.M{
									"$gte": supi[0],
								},
							},
						},
					},
				},
			}
		} else if targetNfType == "UDM" {
			supiFilter = bson.M{
				"udmInfo": bson.M{
					"$elemMatch": bson.M{
						"supiRanges": bson.M{
							"$elemMatch": bson.M{
								"start": bson.M{
									"$lte": supi[0],
								},
								"end": bson.M{
									"$gte": supi[0],
								},
							},
						},
					},
				},
			}
		} else if targetNfType == "UDR" {
			supiFilter = bson.M{
				"udrInfo": bson.M{
					"$elemMatch": bson.M{
						"supiRanges": bson.M{
							"$elemMatch": bson.M{
								"start": bson.M{
									"$lte": supi[0],
								},
								"end": bson.M{
									"$gte": supi[0],
								},
							},
						},
					},
				},
			}
		}
		if queryParameters["supi"].negative {
			supiFilter = bson.M{
				"$not": supiFilter,
			}
		}
		filter[logicalOperator] = append(filter[logicalOperator].([]bson.M), supiFilter)
	}

	// [Query-19] ue-ipv4-address
	if queryParameters["ue-ipv4-address"] != nil {
		var ueIpv4AddressFilter bson.M
		if targetNfType == "BSF" {
			ueIpv4Address := queryParameters["ue-ipv4-address"].value
			ueIpv4AddressNumber := context.Ipv4ToInt(ueIpv4Address)
			ueIpv4AddressFilter = bson.M{
				"bsfInfo": bson.M{
					"$elemMatch": bson.M{
						"ipv4AddressNumberRanges": bson.M{
							"$elemMatch": bson.M{
								"start": bson.M{
									"$lte": ueIpv4AddressNumber,
								},
								"end": bson.M{
									"$gte": ueIpv4AddressNumber,
								},
							},
						},
					},
				},
			}
		}
		if queryParameters["ue-ipv4-address"].negative {
			ueIpv4AddressFilter = bson.M{
				"$not": ueIpv4AddressFilter,
			}
		}
		filter[logicalOperator] = append(filter[logicalOperator].([]bson.M), ueIpv4AddressFilter)
	}

	// [Query-20] ip-domain
	if queryParameters["ip-domain"] != nil {
		var ipDomainFilter bson.M
		if targetNfType == "BSF" {
			ipDomain := queryParameters["ip-domain"].value
			ipDomainFilter = bson.M{
				"bsfInfo": bson.M{
					"$elemMatch": bson.M{
						"ipDomain": ipDomain[0],
					},
				},
			}
		}
		if queryParameters["ip-domain"].negative {
			ipDomainFilter = bson.M{
				"$not": ipDomainFilter,
			}
		}
		filter[logicalOperator] = append(filter[logicalOperator].([]bson.M), ipDomainFilter)
	}

	// [Query-21] ue-ipv6-prefix
	if queryParameters["ue-ipv6-prefix"] != nil {
		var ueIpv6PrefixFilter bson.M
		if targetNfType == "BSF" {
			ueIpv6Prefix := queryParameters["ue-ipv6-prefix"].value
			ueIpv6PrefixNumber := context.Ipv6ToInt(ueIpv6Prefix)
			ueIpv6PrefixFilter = bson.M{
				"bsfInfo": bson.M{
					"$elemMatch": bson.M{
						"ipv6PrefixRanges": bson.M{
							"$elemMatch": bson.M{
								"start": bson.M{
									"$lte": ueIpv6PrefixNumber,
								},
								"end": bson.M{
									"$gte": ueIpv6PrefixNumber,
								},
							},
						},
					},
				},
			}
		}
		if queryParameters["ue-ipv6-prefix"].negative {
			ueIpv6PrefixFilter = bson.M{
				"$not": ueIpv6PrefixFilter,
			}
		}
		filter[logicalOperator] = append(filter[logicalOperator].([]bson.M), ueIpv6PrefixFilter)
	}

	// [Query-22] pgw-ind
	if queryParameters["pgw-ind"] != nil {
		var pgwIndFilter bson.M
		pgwInd := queryParameters["pgw-ind"].value
		if pgwInd == "true" {
			pgwIndFilter = bson.M{
				"smfInfo": bson.M{
					"$elemMatch": bson.M{
						"pgwFqdn": bson.M{
							"$ne": "",
						},
					},
				},
			}
		}
		if queryParameters["pgw-ind"].negative {
			pgwIndFilter = bson.M{
				"$not": pgwIndFilter,
			}
		}
		filter[logicalOperator] = append(filter[logicalOperator].([]bson.M), pgwIndFilter)
	}

	// [Query-23] pgw
	if queryParameters["pgw"] != nil {
		pgw := queryParameters["pgw"].value
		pgwFilter := bson.M{
			"smfInfo": bson.M{
				"$elemMatch": bson.M{
					"pgwFqdn": pgw[0],
				},
			},
		}
		if queryParameters["pgw"].negative {
			pgwFilter = bson.M{
				"$not": pgwFilter,
			}
		}
		filter[logicalOperator] = append(filter[logicalOperator].([]bson.M), pgwFilter)
	}

	// [Query-24] gpsi
	if queryParameters["gpsi"] != nil {
		var gpsiFilter bson.M
		gpsi := queryParameters["gpsi"].value
		if targetNfType == "CHF" {
			gpsiFilter = bson.M{
				"chfInfo": bson.M{
					"$elemMatch": bson.M{
						"gpsiRangeList": bson.M{
							"$elemMatch": bson.M{
								"start": bson.M{
									"$lte": gpsi,
								},
								"end": bson.M{
									"$gte": supi,
								},
							},
						},
					},
				},
			}
		} else if targetNfType == "UDM" {
			gpsiFilter = bson.M{
				"udmInfo": bson.M{
					"$elemMatch": bson.M{
						"gpsiRangeList": bson.M{
							"$elemMatch": bson.M{
								"start": bson.M{
									"$lte": gpsi[0],
								},
								"end": bson.M{
									"$gte": supi[0],
								},
							},
						},
					},
				},
			}
		} else if targetNfType == "UDR" {
			gpsiFilter = bson.M{
				"udrInfo": bson.M{
					"$elemMatch": bson.M{
						"gpsiRangeList": bson.M{
							"$elemMatch": bson.M{
								"start": bson.M{
									"$lte": gpsi[0],
								},
								"end": bson.M{
									"$gte": supi[0],
								},
							},
						},
					},
				},
			}
		}
		if queryParameters["gpsi"].negative {
			gpsiFilter = bson.M{
				"$not": gpsiFilter,
			}
		}
		filter[logicalOperator] = append(filter[logicalOperator].([]bson.M), gpsiFilter)
	}

	// [Query-25] external-group-identity
	if queryParameters["external-group-identity"] != nil {
		var externalGroupIdentityFilter bson.M
		externalGroupIdentity := queryParameters["external-group-identity"].value
		if targetNfType == "UDM" {
			externalGroupIdentityFilter = bson.M{
				"udmInfo": bson.M{
					"$elemMatch": bson.M{
						"groupId": externalGroupIdentity,
					},
				},
			}
		} else if targetNfType == "UDR" {
			externalGroupIdentityFilter = bson.M{
				"udrInfo": bson.M{
					"$elemMatch": bson.M{
						"groupId": externalGroupIdentity,
					},
				},
			}
		}
		if queryParameters["external-group-identity"].negative {
			externalGroupIdentityFilter = bson.M{
				"$not": externalGroupIdentityFilter,
			}
		}
		filter[logicalOperator] = append(filter[logicalOperator].([]bson.M), externalGroupIdentityFilter)
	}

	// [Query-26] data-set
	if queryParameters["data-set"] != nil {
		var dataSetFilter bson.M
		dataSet := queryParameters["data-set"]
		if targetNfType == "UDR" {
			dataSetFilter = bson.M{
				"udrInfo": bson.M{
					"$elemMatch": bson.M{
						"SupportedDataSets": dataSet,
					},
				},
			}
		}
		if queryParameters["data-set"].negative {
			dataSetFilter = bson.M{
				"$not": dataSetFilter,
			}
		}
		filter[logicalOperator] = append(filter[logicalOperator].([]bson.M), dataSetFilter)
	}

	// [Query-27] routing-indicator
	if queryParameters["routing-indicator"] != nil {
		var routingIndicatorFilter bson.M
		routingIndicator := queryParameters["routing-indicator"].value
		if targetNfType == "AUSF" {
			routingIndicatorFilter = bson.M{
				"ausfInfo": bson.M{
					"$elemMatch": bson.M{
						"routingIndicators": routingIndicator,
					},
				},
			}
		} else if targetNfType == "UDM" {
			routingIndicatorFilter = bson.M{
				"udmInfo": bson.M{
					"$elemMatch": bson.M{
						"routingIndicators": routingIndicator,
					},
				},
			}
		}
		if queryParameters["routing-indicator"].negative {
			routingIndicatorFilter = bson.M{
				"$not": routingIndicatorFilter,
			}
		}
		filter[logicalOperator] = append(filter[logicalOperator].([]bson.M), routingIndicatorFilter)
	}

	// [Query-28] group-id-list
	if queryParameters["group-id-list"] != nil {
		var groupIdListFilter bson.M

		groupIdList := queryParameters["group-id-list"].value
		groupIdListSplit := strings.Split(groupIdList, ",")
		var groupIdListBsonArray bson.A

		for _, v := range groupIdListSplit {
			groupIdListBsonArray = append(groupIdListBsonArray, v)
		}

		if targetNfType == "UDR" {
			groupIdListFilter = bson.M{
				"udrInfo": bson.M{
					"$elemMatch": bson.M{
						"groupId": bson.M{
							"$in": groupIdListBsonArray,
						},
					},
				},
			}
		} else if targetNfType == "UDM" {
			groupIdListFilter = bson.M{
				"udmInfo": bson.M{
					"$elemMatch": bson.M{
						"groupId": bson.M{
							"$in": groupIdListBsonArray,
						},
					},
				},
			}
		} else if targetNfType == "AUSF" {
			groupIdListFilter = bson.M{
				"ausfInfo": bson.M{
					"$elemMatch": bson.M{
						"groupId": bson.M{
							"$in": groupIdListBsonArray,
						},
					},
				},
			}
		}
		if queryParameters["group-id-list"].negative {
			groupIdListFilter = bson.M{
				"$not": groupIdListFilter,
			}
		}
		filter[logicalOperator] = append(filter[logicalOperator].([]bson.M), groupIdListFilter)
	}

	// [Query-29] dnai-list
	if queryParameters["dnai-list"] != nil {
		var dnaiFilter bson.M
		dnaiList := queryParameters["dnai-list"].value
		dnaiListSplit := strings.Split(dnaiList, ",")
		var dnaiListBsonArray bson.A

		for _, v := range dnaiListSplit {
			dnaiListBsonArray = append(dnaiListBsonArray, v)
		}
		if targetNfType == "UPF" {
			dnaiFilter = bson.M{
				"upfInfo": bson.M{
					"$elemMatch": bson.M{
						"sNssaiUpfInfoList": bson.M{
							"$elemMatch": bson.M{
								"dnnUpfInfoList": bson.M{
									"$elemMatch": bson.M{
										"dnaiList": dnaiListBsonArray,
									},
								},
							},
						},
					},
				},
			}
		}
		if queryParameters["dnai-list"].negative {
			dnaiFilter = bson.M{
				"$not": dnaiFilter,
			}
		}
		filter[logicalOperator] = append(filter[logicalOperator].([]bson.M), dnaiFilter)
	}

	// [Query-30] upf-iwk-eps-ind
	if queryParameters["upf-iwk-eps-ind"] != nil {
		var upfIwkEpsIndFilter bson.M
		// upfIwkEpsInd := queryParameters["upf-iwk-eps-ind"].value
		if targetNfType == "UPF" {
			upfIwkEpsIndFilter = bson.M{
				"upfInfo": bson.M{
					"$elemMatch": bson.M{
						"iwkEpsInd": true,
					},
				},
			}
		}
		if queryParameters["upf-iwk-eps-ind"].negative {
			upfIwkEpsIndFilter = bson.M{
				"$not": upfIwkEpsIndFilter,
			}
		}
		filter[logicalOperator] = append(filter[logicalOperator].([]bson.M), upfIwkEpsIndFilter)
	}

	// [Query-31] chf-supported-plmn
	if queryParameters["chf-supported-plmn"] != nil {
		var chfSupportedPlmnFilter bson.M
		chfSupportedPlmn := queryParameters["chf-supported-plmn"].value
		if targetNfType == "CHF" {
			chfSupportedPlmnFilter = bson.M{
				"$or": []bson.M{
					{
						"chfInfo": bson.M{
							"$elemMatch": bson.M{
								"plmnRangeList": bson.M{
									"$elemMatch": bson.M{
										"start": bson.M{
											"$lte": chfSupportedPlmn,
										},
										"end": bson.M{
											"$gte": chfSupportedPlmn,
										},
									},
								},
							},
						},
					},
					{
						"chfInfo.plmnRangeList": bson.M{
							"$exists": false,
						},
					},
				},
			}
		}
		if queryParameters["chf-supported-plmn"].negative {
			chfSupportedPlmnFilter = bson.M{
				"$not": chfSupportedPlmnFilter,
			}
		}
		filter[logicalOperator] = append(filter[logicalOperator].([]bson.M), chfSupportedPlmnFilter)
	}

	// [Query-32]  preferred-locality
	// TODO: if no match
	if queryParameters["preferred-locality"] != nil {
		preferredLocality := queryParameters["preferred-locality"].value
		preferredLocalityFilter := bson.M{
			"locality": preferredLocality,
		}
		if queryParameters["preferred-locality"].negative {
			preferredLocalityFilter = bson.M{
				"$not": preferredLocalityFilter,
			}
		}
		filter[logicalOperator] = append(filter[logicalOperator].([]bson.M), preferredLocalityFilter)
	}

	// [Query-33] access-type
	if queryParameters["access-type"] != nil {
		accessType := queryParameters["access-type"].value
		accessTypeFilter := bson.M{
			"smfInfo": bson.M{
				"$elemMatch": bson.M{
					"accessType": accessType[0],
				},
			},
		}
		if queryParameters["access-type"].negative {
			accessTypeFilter = bson.M{
				"$not": accessTypeFilter,
			}
		}
		filter[logicalOperator] = append(filter[logicalOperator].([]bson.M), accessTypeFilter)
	}

	// [Query-34] supported-features
	if queryParameters["supported-features"] != nil {
		supportedFeatures := queryParameters["supported-features"].value
		supportedFeaturesFilter := bson.M{
			"nfServices": bson.M{
				"$elemMatch": bson.M{
					"supportedFeatures": supportedFeatures,
				},
			},
		}
		if queryParameters["supported-features"].negative {
			supportedFeaturesFilter = bson.M{
				"$not": supportedFeaturesFilter,
			}
		}
		filter[logicalOperator] = append(filter[logicalOperator].([]bson.M), supportedFeaturesFilter)
	}

	return filter
}
