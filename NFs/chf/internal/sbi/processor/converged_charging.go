package processor

import (
	"context"
	"encoding/json"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/fiorix/go-diameter/diam/datatype"
	"github.com/gin-gonic/gin"
	"golang.org/x/exp/constraints"

	charging_datatype "github.com/free5gc/chf/ccs_diameter/datatype"
	"github.com/free5gc/chf/cdr/asn"
	"github.com/free5gc/chf/cdr/cdrConvert"
	"github.com/free5gc/chf/cdr/cdrType"
	"github.com/free5gc/chf/internal/abmf"
	"github.com/free5gc/chf/internal/cgf"
	chf_context "github.com/free5gc/chf/internal/context"
	"github.com/free5gc/chf/internal/logger"
	"github.com/free5gc/chf/internal/rating"
	"github.com/free5gc/chf/internal/util"
	Nchf_ConvergedCharging "github.com/free5gc/openapi/chf/ConvergedCharging"
	"github.com/free5gc/openapi/models"
)

func min[T constraints.Ordered](a, b T) T {
	if a < b {
		return a
	}
	return b
}

func (p *Processor) NotifyRecharge(ueId string, rg int32) {
	var reauthorizationDetails []models.ReauthorizationDetails

	self := chf_context.GetSelf()
	ue, ok := self.ChfUeFindBySupi(ueId)
	if !ok {
		logger.NotifyEventLog.Errorf("Do not find charging data for UE: %s", ueId)
		return
	}

	// If it is previosly set to debit mode due to quota exhausted, need to reverse to the reserve mode
	ue.RatingType[rg] = charging_datatype.REQ_SUBTYPE_RESERVE
	reauthorizationDetails = append(reauthorizationDetails, models.ReauthorizationDetails{
		RatingGroup: rg,
	})

	notifyRequest := models.ChargingNotifyRequest{
		ReauthorizationDetails: reauthorizationDetails,
	}

	p.SendChargingNotification(ue.NotifyUri, notifyRequest)
}

func (p *Processor) SendChargingNotification(notifyUri string, notifyRequest models.ChargingNotifyRequest) {
	client := util.GetNchfChargingNotificationCallbackClient()
	logger.NotifyEventLog.Warn("Send Charging Notification  to SMF: uri: ", notifyUri)
	chargingNotifyRequest := Nchf_ConvergedCharging.NewPostChargingNotificationRequest()
	chargingNotifyRequest.SetChargingNotifyRequest(notifyRequest)
	_, err := client.DefaultApi.PostChargingNotification(context.Background(), notifyUri, chargingNotifyRequest)
	if err != nil {
		logger.NotifyEventLog.Warnf("Charging Notification Failed[%s]", err.Error())
		return
	}

	logger.NotifyEventLog.Tracef("Charging Notification Success")
}

func (p *Processor) HandleChargingdataInitial(
	c *gin.Context,
	chargingdata models.ChfConvergedChargingChargingDataRequest,
) {
	logger.ChargingdataPostLog.Infof("HandleChargingdataInitial")
	response, locationURI, problemDetails := p.ChargingDataCreate(chargingdata)

	if response != nil {
		c.Header("Location", locationURI)
		c.JSON(http.StatusCreated, response)
		return
	} else if problemDetails != nil {
		c.JSON(int(problemDetails.Status), problemDetails)
		return
	}
	problemDetails = &models.ProblemDetails{
		Status: http.StatusForbidden,
		Cause:  "UNSPECIFIED",
	}
	c.JSON(int(problemDetails.Status), problemDetails)
}

func (p *Processor) HandleChargingdataUpdate(
	c *gin.Context,
	chargingdata models.ChfConvergedChargingChargingDataRequest,
	chargingSessionId string,
) {
	logger.ChargingdataPostLog.Infof("HandleChargingdataUpdate")
	response, problemDetails := p.ChargingDataUpdate(chargingdata, chargingSessionId)

	if response != nil {
		c.JSON(http.StatusOK, response)
		return
	} else if problemDetails != nil {
		c.JSON(int(problemDetails.Status), problemDetails)
		return
	}
	problemDetails = &models.ProblemDetails{
		Status: http.StatusForbidden,
		Cause:  "UNSPECIFIED",
	}
	c.JSON(int(problemDetails.Status), problemDetails)
}

func (p *Processor) HandleChargingdataRelease(
	c *gin.Context,
	chargingdata models.ChfConvergedChargingChargingDataRequest,
	chargingSessionId string,
) {
	logger.ChargingdataPostLog.Infof("HandleChargingdateRelease")

	problemDetails := p.ChargingDataRelease(chargingdata, chargingSessionId)
	if problemDetails == nil {
		c.Status(http.StatusBadRequest)
		return
	}
	c.JSON(int(problemDetails.Status), problemDetails)
}

func (p *Processor) ChargingDataCreate(
	chargingData models.ChfConvergedChargingChargingDataRequest,
) (
	*models.ChfConvergedChargingChargingDataResponse,
	string, *models.ProblemDetails,
) {
	var responseBody models.ChfConvergedChargingChargingDataResponse
	var chargingSessionId string

	self := chf_context.GetSelf()
	ueId := chargingData.SubscriberIdentifier

	// Open CDR
	// ChargingDataRef(charging session id):
	// A unique identifier for a charging data resource in a PLMN
	// TODO determine charging session id(string type) supi+consumerid+localseq?
	ue, err := self.NewCHFUe(ueId)
	if err != nil {
		logger.ChargingdataPostLog.Errorf("New CHFUe error %s", err)
		problemDetails := &models.ProblemDetails{
			Status: http.StatusBadRequest,
		}
		return nil, "", problemDetails
	}

	ue.CULock.Lock()
	ue.NotifyUri = chargingData.NotifyUri

	consumerId := chargingData.NfConsumerIdentification.NFName
	if !chargingData.OneTimeEvent {
		chargingSessionId = ueId + consumerId + strconv.Itoa(int(self.LocalRecordSequenceNumber))
	}
	cdr, err := p.OpenCDR(chargingData, ue, chargingSessionId, false)
	if err != nil {
		// Lock in line 158
		ue.CULock.Unlock()
		problemDetails := &models.ProblemDetails{
			Status: http.StatusBadRequest,
		}
		return nil, "", problemDetails
	}

	err = p.UpdateCDR(cdr, chargingData)
	if err != nil {
		// Lock in line 158
		ue.CULock.Unlock()
		problemDetails := &models.ProblemDetails{
			Status: http.StatusBadRequest,
		}
		return nil, "", problemDetails
	}

	ue.Cdr[chargingSessionId] = cdr
	ue.Records = append(ue.Records, ue.Cdr[chargingSessionId])
	ue.CULock.Unlock()

	if chargingData.OneTimeEvent {
		err = p.CloseCDR(cdr, false)
		if err != nil {
			problemDetails := &models.ProblemDetails{
				Status: http.StatusBadRequest,
			}
			return nil, "", problemDetails
		}
	}

	// CDR Transfer
	err = cgf.SendCDR(chargingData.SubscriberIdentifier)
	if err != nil {
		logger.ChargingdataPostLog.Errorf("Charging gateway fail to send CDR to billing domain %v", err)
	}

	logger.ChargingdataPostLog.Infof("Open CDR for UE %s", ueId)

	// build response
	logger.ChargingdataPostLog.Infof("NewChfUe %s", ueId)
	locationURI := self.Url + "/nchf-convergedcharging/v3/chargingdata/" + chargingSessionId
	timeStamp := time.Now()

	responseBody.InvocationTimeStamp = &timeStamp
	responseBody.InvocationSequenceNumber = chargingData.InvocationSequenceNumber

	return &responseBody, locationURI, nil
}

func (p *Processor) ChargingDataUpdate(
	chargingData models.ChfConvergedChargingChargingDataRequest, chargingSessionId string,
) (*models.ChfConvergedChargingChargingDataResponse, *models.ProblemDetails) {
	self := chf_context.GetSelf()
	ueId := chargingData.SubscriberIdentifier
	ue, ok := self.ChfUeFindBySupi(ueId)
	if !ok {
		logger.ChargingdataPostLog.Errorf("CHFUe[%s] not found", ueId)
		problemDetails := &models.ProblemDetails{
			Status: http.StatusBadRequest,
		}
		return nil, problemDetails
	}

	ue.CULock.Lock()
	defer ue.CULock.Unlock()

	// Online charging: Rate, Account, Reservation
	responseBody, partialRecord := p.BuildConvergedChargingDataUpdateResopone(chargingData)

	cdr := ue.Cdr[chargingSessionId]

	if len(ue.Records) > 1 {
		cdr = ue.Records[len(ue.Records)-1]
	}

	cdrBytes, errCdrBer := asn.BerMarshalWithParams(&cdr, "explicit,choice")
	if errCdrBer != nil {
		logger.ChargingdataPostLog.Error(errCdrBer)
		problemDetails := &models.ProblemDetails{
			Status: http.StatusBadRequest,
			Detail: errCdrBer.Error(),
		}
		return nil, problemDetails
	}

	var chgDataBytes []byte
	var errChgDataBer error
	if len(chargingData.MultipleUnitUsage) != 0 {
		cdrMultiUnitUsage := cdrConvert.MultiUnitUsageToCdr(chargingData.MultipleUnitUsage)
		chgDataBytes, errChgDataBer = asn.BerMarshalWithParams(&cdrMultiUnitUsage, "explicit,choice")
		if errChgDataBer != nil {
			logger.ChargingdataPostLog.Error(errChgDataBer)
			problemDetails := &models.ProblemDetails{
				Status: http.StatusBadRequest,
				Detail: errChgDataBer.Error(),
			}
			return nil, problemDetails
		}
	}

	if len(cdrBytes)+len(chgDataBytes) > math.MaxUint16 {
		var newRecord *cdrType.CHFRecord
		cdrJson, err := json.Marshal(cdr)
		if err != nil {
			logger.ChargingdataPostLog.Error(err)
		}
		err = json.Unmarshal(cdrJson, &newRecord)
		if err != nil {
			logger.ChargingdataPostLog.Error(err)
		}

		newRecord.ChargingFunctionRecord.ListOfMultipleUnitUsage = []cdrType.MultipleUnitUsage{}
		cdr = newRecord
		ue.Records = append(ue.Records, cdr)
	}

	err := p.UpdateCDR(cdr, chargingData)
	if err != nil {
		problemDetails := &models.ProblemDetails{
			Status: http.StatusBadRequest,
		}
		return nil, problemDetails
	}

	if partialRecord {
		ueId = chargingData.SubscriberIdentifier

		close_err := p.CloseCDR(cdr, partialRecord)
		if close_err != nil {
			logger.ChargingdataPostLog.Error("CloseCDR error:", close_err)
		}
		err = dumpCdrFile(ueId, []*cdrType.CHFRecord{cdr})
		if err != nil {
			problemDetails := &models.ProblemDetails{
				Status: http.StatusBadRequest,
			}
			return nil, problemDetails
		}

		_, oper_err := p.OpenCDR(chargingData, ue, chargingSessionId, partialRecord)
		if oper_err != nil {
			logger.ChargingdataPostLog.Error("OpenCDR error:", oper_err)
		}
		logger.ChargingdataPostLog.Tracef(
			"CDR Record Sequence Number after Reopen %+v", *cdr.ChargingFunctionRecord.RecordSequenceNumber)
	}

	err = dumpCdrFile(ueId, ue.Records)
	if err != nil {
		problemDetails := &models.ProblemDetails{
			Status: http.StatusBadRequest,
		}
		return nil, problemDetails
	}

	err = cgf.SendCDR(chargingData.SubscriberIdentifier)
	if err != nil {
		logger.ChargingdataPostLog.Errorf("Charging gateway fail to send CDR to billing domain %v", err)
	}

	timeStamp := time.Now()
	responseBody.InvocationTimeStamp = &timeStamp
	responseBody.InvocationSequenceNumber = chargingData.InvocationSequenceNumber

	return &responseBody, nil
}

func (p *Processor) ChargingDataRelease(
	chargingData models.ChfConvergedChargingChargingDataRequest, chargingSessionId string,
) *models.ProblemDetails {
	self := chf_context.GetSelf()
	ueId := chargingData.SubscriberIdentifier
	ue, ok := self.ChfUeFindBySupi(ueId)
	if !ok {
		logger.ChargingdataPostLog.Errorf("Do not find CHFUe[%s] error", ueId)
		problemDetails := &models.ProblemDetails{
			Status: http.StatusBadRequest,
		}
		return problemDetails
	}

	ue.CULock.Lock()
	defer ue.CULock.Unlock()

	sessionChargingReservation(chargingData)

	cdr := ue.Cdr[chargingSessionId]

	err := p.UpdateCDR(cdr, chargingData)
	if err != nil {
		problemDetails := &models.ProblemDetails{
			Status: http.StatusBadRequest,
			Cause:  err.Error(),
		}
		return problemDetails
	}

	err = p.CloseCDR(cdr, false)
	if err != nil {
		problemDetails := &models.ProblemDetails{
			Status: http.StatusBadRequest,
		}
		return problemDetails
	}

	err = dumpCdrFile(ueId, []*cdrType.CHFRecord{cdr})
	if err != nil {
		problemDetails := &models.ProblemDetails{
			Status: http.StatusBadRequest,
		}
		return problemDetails
	}

	return nil
}

func (p *Processor) BuildOnlineChargingDataCreateResopone(
	ue *chf_context.ChfUe, chargingData models.ChfConvergedChargingChargingDataRequest,
) models.ChfConvergedChargingChargingDataResponse {
	logger.ChargingdataPostLog.Info("In Build Online Charging Data Create Resopone")
	ue.NotifyUri = chargingData.NotifyUri

	multipleUnitInformation, _ := sessionChargingReservation(chargingData)

	responseBody := models.ChfConvergedChargingChargingDataResponse{
		MultipleUnitInformation: multipleUnitInformation,
	}

	return responseBody
}

func (p *Processor) BuildConvergedChargingDataUpdateResopone(
	chargingData models.ChfConvergedChargingChargingDataRequest,
) (models.ChfConvergedChargingChargingDataResponse, bool) {
	var partialRecord bool

	logger.ChargingdataPostLog.Info("In BuildConvergedChargingDataUpdateResopone")

	multipleUnitInformation, partialRecord := sessionChargingReservation(chargingData)

	responseBody := models.ChfConvergedChargingChargingDataResponse{
		MultipleUnitInformation: multipleUnitInformation,
	}

	return responseBody, partialRecord
}

func getUnitCost(ue *chf_context.ChfUe, rg int32, sur *charging_datatype.ServiceUsageRequest) uint32 {
	if sur == nil {
		logger.ChargingdataPostLog.Errorln("ServiceUsageRequest is nil, set unitCost to 1")
		return 1
	}

	sur.ServiceRating = &charging_datatype.ServiceRating{
		ServiceIdentifier: datatype.Unsigned32(rg),
		MonetaryQuota:     datatype.Unsigned32(0), // dummy
		RequestSubType:    charging_datatype.REQ_SUBTYPE_RESERVE,
	}

	serviceUsageRsp, err := rating.SendServiceUsageRequest(ue, sur)
	if err != nil {
		logger.ChargingdataPostLog.Errorf("err: %+v", err)
		logger.ChargingdataPostLog.Errorln("cannot get unitCost by SendServiceUsageRequest, set unitCost to 1")
		return 1
	}

	return uint32(serviceUsageRsp.ServiceRating.MonetaryTariff.RateElement.UnitCost.ValueDigits) *
		uint32(math.Pow10(int(serviceUsageRsp.ServiceRating.MonetaryTariff.RateElement.UnitCost.Exponent)))
}

// 32.296 6.2.2.3.1: Service usage request method with reservation
func sessionChargingReservation(
	chargingData models.ChfConvergedChargingChargingDataRequest,
) ([]models.MultipleUnitInformation, bool) {
	var multipleUnitInformation []models.MultipleUnitInformation
	var partialRecord bool
	var subscriberIdentifier *charging_datatype.SubscriptionId

	self := chf_context.GetSelf()
	supi := chargingData.SubscriberIdentifier

	ue, ok := self.ChfUeFindBySupi(supi)
	if !ok {
		logger.ChargingdataPostLog.Warnf("Do not find UE[%s]", supi)
		return nil, false
	}

	supiType := strings.Split(supi, "-")[0]
	switch supiType {
	case "imsi":
		subscriberIdentifier = &charging_datatype.SubscriptionId{
			SubscriptionIdType: charging_datatype.END_USER_IMSI,
			SubscriptionIdData: datatype.UTF8String(supi[5:]),
		}
	case "nai":
		subscriberIdentifier = &charging_datatype.SubscriptionId{
			SubscriptionIdType: charging_datatype.END_USER_NAI,
			SubscriptionIdData: datatype.UTF8String(supi[4:]),
		}
	case "gci":
		subscriberIdentifier = &charging_datatype.SubscriptionId{
			SubscriptionIdType: charging_datatype.END_USER_NAI,
			SubscriptionIdData: datatype.UTF8String(supi[4:]),
		}
	case "gli":
		subscriberIdentifier = &charging_datatype.SubscriptionId{
			SubscriptionIdType: charging_datatype.END_USER_NAI,
			SubscriptionIdData: datatype.UTF8String(supi[4:]),
		}
	}

	for unitUsageNum, unitUsage := range chargingData.MultipleUnitUsage {
		var totalUsedUnit uint32
		var finalUnitIndication models.FinalUnitIndication
		creditControl := false

		rg := unitUsage.RatingGroup
		if !ue.FindRatingGroup(rg) {
			ue.RatingGroups = append(ue.RatingGroups, rg)
			ue.RatingType[rg] = charging_datatype.REQ_SUBTYPE_RESERVE
		}

		unitInformation := models.MultipleUnitInformation{
			UPFID:               unitUsage.UPFID,
			FinalUnitIndication: &finalUnitIndication,
			RatingGroup:         rg,
		}

		for _, usedUnit := range unitUsage.UsedUnitContainer {
			switch usedUnit.QuotaManagementIndicator {
			case models.QuotaManagementIndicator_OFFLINE_CHARGING:
				unitInformation.Triggers = append(unitInformation.Triggers,
					models.ChfConvergedChargingTrigger{
						TriggerType:     models.ChfConvergedChargingTriggerType_QUOTA_THRESHOLD,
						TriggerCategory: models.TriggerCategory_IMMEDIATE_REPORT,
					},
				)

				unitInformation.VolumeQuotaThreshold = int32(30000000)
				continue
			case models.QuotaManagementIndicator_ONLINE_CHARGING:
				creditControl = true

				for _, trigger := range chargingData.Triggers {
					// Check if partial record is needed
					partialRecord = true
					switch t := trigger; {
					case t == models.ChfConvergedChargingTrigger{
						TriggerType:     models.ChfConvergedChargingTriggerType_VOLUME_LIMIT,
						TriggerCategory: models.TriggerCategory_IMMEDIATE_REPORT,
					}:
					case t.TriggerType == models.ChfConvergedChargingTriggerType_MAX_NUMBER_OF_CHANGES_IN_CHARGING_CONDITIONS:
					case t.TriggerType == models.ChfConvergedChargingTriggerType_MANAGEMENT_INTERVENTION:
					case t.TriggerType == models.ChfConvergedChargingTriggerType_FINAL:
						ue.RatingType[rg] = charging_datatype.REQ_SUBTYPE_DEBIT
						partialRecord = false
					}
				}
				// calculate total used unit
				totalUsedUnit += uint32(usedUnit.TotalVolume)
			case models.QuotaManagementIndicator_QUOTA_MANAGEMENT_SUSPENDED:
				logger.ChargingdataPostLog.Errorf("Current do not support QUOTA MANAGEMENT SUSPENDED")
			}
		}
		if !creditControl {
			logger.ChargingdataPostLog.Infof("Credit Control are not required for rating group: %d", rg)
			continue
		}
		// Only online charging with request unit or used unit need to perform credit control

		ccr := &charging_datatype.AccountDebitRequest{
			SessionId:       datatype.UTF8String(strconv.Itoa(int(ue.AcctSessionId))),
			OriginHost:      datatype.DiameterIdentity(self.AbmfCfg.OriginHost),
			OriginRealm:     datatype.DiameterIdentity(self.AbmfCfg.OriginRealm),
			EventTimestamp:  datatype.Time(time.Now()),
			SubscriptionId:  subscriberIdentifier,
			UserName:        datatype.OctetString(self.Name),
			CcRequestNumber: datatype.Unsigned32(ue.AcctRequestNum[rg]),
		}

		sur := &charging_datatype.ServiceUsageRequest{
			SessionId:      datatype.UTF8String(strconv.Itoa(int(ue.RateSessionId))),
			OriginHost:     datatype.DiameterIdentity(self.RatingCfg.OriginHost),
			OriginRealm:    datatype.DiameterIdentity(self.RatingCfg.OriginRealm),
			ActualTime:     datatype.Time(time.Now()),
			SubscriptionId: subscriberIdentifier,
			UserName:       datatype.OctetString(self.Name),
		}

		switch ue.RatingType[rg] {
		case charging_datatype.REQ_SUBTYPE_RESERVE:
			var requestedQuota uint64

			ue.UnitCost[rg] = getUnitCost(ue, rg, sur)

			usedQuota := uint64(totalUsedUnit * ue.UnitCost[rg])
			requestedQuota = uint64(uint32(unitUsage.RequestedUnit.TotalVolume) * ue.UnitCost[rg])
			ue.ReservedQuota[rg] -= int64(usedQuota)
			NeedReserveQuota := ue.ReservedQuota[rg] <= 0

			if NeedReserveQuota {
				reserveQuota := -uint64(ue.ReservedQuota[rg]) + requestedQuota
				ccr.CcRequestType = charging_datatype.UPDATE_REQUEST
				ccr.RequestedAction = charging_datatype.DIRECT_DEBITING
				ccr.MultipleServicesCreditControl = &charging_datatype.MultipleServicesCreditControl{
					RatingGroup: datatype.Unsigned32(rg),
					RequestedServiceUnit: &charging_datatype.RequestedServiceUnit{
						CCTotalOctets: datatype.Unsigned64(reserveQuota),
					},
				}

				acctDebitRsp, err := abmf.SendAccountDebitRequest(ue, ccr)
				if err != nil {
					logger.ChargingdataPostLog.Errorf("SendAccountDebitRequest err: %+v", err)
					continue
				}

				ue.ReservedQuota[rg] += int64(acctDebitRsp.MultipleServicesCreditControl.GrantedServiceUnit.CCTotalOctets)

				// Deduct the reserved quota from the account
				if acctDebitRsp.MultipleServicesCreditControl.FinalUnitIndication != nil {
					switch acctDebitRsp.MultipleServicesCreditControl.FinalUnitIndication.FinalUnitAction {
					case charging_datatype.TERMINATE:
						logger.ChargingdataPostLog.Tracef("Last granted quota")
						finalUnitIndication = models.FinalUnitIndication{
							FinalUnitAction: models.FinalUnitAction_TERMINATE,
						}
						ue.RatingType[rg] = charging_datatype.REQ_SUBTYPE_DEBIT
					}
				}
			}

			sur.ServiceRating = &charging_datatype.ServiceRating{
				ServiceIdentifier: datatype.Unsigned32(rg),
				MonetaryQuota:     datatype.Unsigned32(requestedQuota),
				RequestSubType:    charging_datatype.REQ_SUBTYPE_RESERVE,
			}

			// Retrieve and save the tarrif for pricing the next usage
			serviceUsageRsp, err := rating.SendServiceUsageRequest(ue, sur)
			if err != nil {
				logger.ChargingdataPostLog.Errorf("SendServiceUsageRequest err: %+v", err)
				continue
			}

			ue.UnitCost[rg] = getUnitCost(ue, rg, sur)

			grantedUnit := min(uint32(serviceUsageRsp.ServiceRating.AllowedUnits), uint32(unitUsage.RequestedUnit.TotalVolume))

			if ue.RatingType[rg] == charging_datatype.REQ_SUBTYPE_RESERVE {
				unitInformation.Triggers = append(unitInformation.Triggers,
					models.ChfConvergedChargingTrigger{
						TriggerType:     models.ChfConvergedChargingTriggerType_QUOTA_THRESHOLD,
						TriggerCategory: models.TriggerCategory_IMMEDIATE_REPORT,
					},
				)

				unitInformation.VolumeQuotaThreshold = int32(float32(grantedUnit) * ue.VolumeThresholdRate)
			}

			unitInformation.Triggers = append(unitInformation.Triggers,
				models.ChfConvergedChargingTrigger{
					TriggerType:     models.ChfConvergedChargingTriggerType_QUOTA_EXHAUSTED,
					TriggerCategory: models.TriggerCategory_IMMEDIATE_REPORT,
				},
			)

			unitInformation.GrantedUnit = &models.GrantedUnit{
				TotalVolume:    int32(grantedUnit),
				DownlinkVolume: int32(grantedUnit),
				UplinkVolume:   int32(grantedUnit),
			}
			logger.ChargingdataPostLog.Tracef("granted Unit: %d", unitInformation.GrantedUnit.TotalVolume)

			// The timer of VolumeLimit is remain in SMF
			if ue.VolumeLimit != 0 {
				unitInformation.Triggers = append(unitInformation.Triggers,
					models.ChfConvergedChargingTrigger{
						TriggerType:     models.ChfConvergedChargingTriggerType_VOLUME_LIMIT,
						TriggerCategory: models.TriggerCategory_DEFERRED_REPORT,
						VolumeLimit:     ue.VolumeLimit,
					},
				)
			}

			// VolumeLimit for PDU session only need to add once
			if ue.VolumeLimitPDU != 0 && unitUsageNum == 0 {
				unitInformation.Triggers = append(unitInformation.Triggers,
					models.ChfConvergedChargingTrigger{
						TriggerType:     models.ChfConvergedChargingTriggerType_VOLUME_LIMIT,
						TriggerCategory: models.TriggerCategory_IMMEDIATE_REPORT,
						VolumeLimit:     ue.VolumeLimitPDU,
					},
				)
			}

			// The timer of QuotaValidityTime is remain in UPF
			if ue.QuotaValidityTime != 0 {
				unitInformation.Triggers = append(unitInformation.Triggers,
					models.ChfConvergedChargingTrigger{
						TriggerType:     models.ChfConvergedChargingTriggerType_VALIDITY_TIME,
						TriggerCategory: models.TriggerCategory_IMMEDIATE_REPORT,
					},
				)
				unitInformation.ValidityTime = ue.QuotaValidityTime
			}

		case charging_datatype.REQ_SUBTYPE_DEBIT:
			logger.ChargingdataPostLog.Info("Debit mode, will not grant unit")
			// retrieved tarrif for final pricing
			sur.ServiceRating = &charging_datatype.ServiceRating{
				ServiceIdentifier: datatype.Unsigned32(rg),
				ConsumedUnits:     datatype.Unsigned32(totalUsedUnit),
				RequestSubType:    charging_datatype.REQ_SUBTYPE_DEBIT,
			}

			serviceUsageRsp, err := rating.SendServiceUsageRequest(ue, sur)
			if err != nil {
				logger.ChargingdataPostLog.Errorf("SendServiceUsageRequest err: %+v", err)
				continue
			}
			logger.ChargingdataPostLog.Tracef(
				"price %+v, ue.ReservedQuota[rg]: %+v", serviceUsageRsp.ServiceRating.Price, ue.ReservedQuota[rg])

			if int64(serviceUsageRsp.ServiceRating.Price) < ue.ReservedQuota[rg] {
				// The final consumed quota is smaller than the reserved quota
				// Therefore, return the extra reserved quota back to the user account
				reservedRemained := ue.ReservedQuota[rg] - int64(serviceUsageRsp.ServiceRating.Price)
				ccr.RequestedAction = charging_datatype.REFUND_ACCOUNT
				ccr.MultipleServicesCreditControl = &charging_datatype.MultipleServicesCreditControl{
					RatingGroup: datatype.Unsigned32(rg),
					RequestedServiceUnit: &charging_datatype.RequestedServiceUnit{
						CCTotalOctets: datatype.Unsigned64(reservedRemained),
					},
				}
				// Typically, the reserved quota will be exhausted for the flow (or PDU session)
				// However, for the case the flow quota  and PDU session's quota is both last granted quota
				// and the PDU session's quota is larger than the flow's quota
				// PDU session's quota should be refund and set to reserved mode in order to reserve the quota for other flow
				ue.RatingType[rg] = charging_datatype.REQ_SUBTYPE_RESERVE
			} else {
				// The final consumed quota exceed the reserved quota
				// Deduct the extra consumed quota from the user account
				extraConsumed := int64(serviceUsageRsp.ServiceRating.Price) - ue.ReservedQuota[rg]
				ccr.RequestedAction = charging_datatype.DIRECT_DEBITING
				ccr.CcRequestType = charging_datatype.TERMINATION_REQUEST
				ccr.MultipleServicesCreditControl = &charging_datatype.MultipleServicesCreditControl{
					RatingGroup: datatype.Unsigned32(rg),
					UsedServiceUnit: &charging_datatype.UsedServiceUnit{
						CCTotalOctets: datatype.Unsigned64(extraConsumed),
					},
				}
			}

			_, err = abmf.SendAccountDebitRequest(ue, ccr)
			if err != nil {
				logger.ChargingdataPostLog.Errorf("SendAccountDebitRequest err: %+v", err)
				continue
			}
			ue.ReservedQuota[rg] = 0

			unitInformation.Triggers = append(unitInformation.Triggers,
				models.ChfConvergedChargingTrigger{
					TriggerType:     models.ChfConvergedChargingTriggerType_QUOTA_EXHAUSTED,
					TriggerCategory: models.TriggerCategory_IMMEDIATE_REPORT,
				},
			)
			unitInformation.GrantedUnit = &models.GrantedUnit{
				TotalVolume:    int32(0),
				DownlinkVolume: int32(0),
				UplinkVolume:   int32(0),
			}
		}
		multipleUnitInformation = append(multipleUnitInformation, unitInformation)

		ue.AcctRequestNum[rg]++
	}

	return multipleUnitInformation, partialRecord
}
