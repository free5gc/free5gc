package pcf_producer

import (
	"context"
	"free5gc/lib/Nudr_DataRepository"
	"free5gc/lib/openapi/models"
	"free5gc/src/pcf/logger"
	"free5gc/src/pcf/pcf_context"
	"free5gc/src/pcf/pcf_handler/pcf_message"
	"free5gc/src/pcf/pcf_util"
	"net/http"
	"time"

	"github.com/antihax/optional"
)

var BdtPolicyStore = make(map[string]*models.BdtPolicy) // key is aspid

func GetBDTPolicyContext(httpChannel chan pcf_message.HttpResponseMessage, BdtPolicyId string) {

	var problem models.ProblemDetails
	// check BdtRefId from pcfUecontext
	for key := range BdtPolicyStore {
		if BdtPolicyStore[key] == nil {
			continue
		}
		if BdtPolicyId == BdtPolicyStore[key].BdtPolData.BdtRefId {
			pcf_message.SendHttpResponseMessage(httpChannel, nil, 200, BdtPolicyStore[key])
			return
		}
	}
	// can not found
	problem.Status = 404
	problem.Cause = "CONTEXT_NOT_FOUND"
	pcf_message.SendHttpResponseMessage(httpChannel, nil, 404, problem)

}

// UpdateBDTPolicy - Update an Individual BDT policy
func UpdateBDTPolicyContext(httpChannel chan pcf_message.HttpResponseMessage, ReqURI string, body models.BdtPolicyDataPatch) {
	var problem models.ProblemDetails
	URI := ReqURI
	bdtPolicyDataPatch := body
	// check BdtRefId from pcfUeContext
	for key := range BdtPolicyStore {
		if BdtPolicyStore[key] == nil {
			continue
		}
		if URI == BdtPolicyStore[key].BdtPolData.BdtRefId {
			if bdtPolicyDataPatch.SelTransPolicyId == BdtPolicyStore[key].BdtPolData.SelTransPolicyId {
				BdtPolicyStore[key].BdtPolData.SelTransPolicyId = bdtPolicyDataPatch.SelTransPolicyId
				pcf_message.SendHttpResponseMessage(httpChannel, nil, 204, BdtPolicyStore[key])
				return
			}
			for i := 0; i < len(BdtPolicyStore[key].BdtPolData.TransfPolicies); i++ {
				if bdtPolicyDataPatch.SelTransPolicyId == BdtPolicyStore[key].BdtPolData.TransfPolicies[i].TransPolicyId {
					transPolicy := BdtPolicyStore[key].BdtPolData.TransfPolicies[i]
					BdtPolicyStore[key].BdtPolData.SelTransPolicyId = bdtPolicyDataPatch.SelTransPolicyId
					// update bdtdata to udr
					client := pcf_util.GetNudrClient("https://localhost:29504")
					var bdtData models.BdtData
					bdtData.AspId = key
					bdtData.BdtRefId = BdtPolicyStore[key].BdtPolData.BdtRefId
					bdtData.TransPolicy = transPolicy
					var PolicyDataBdtData optional.Interface
					PolicyDataBdtData.Default(bdtData)
					data := Nudr_DataRepository.PolicyDataBdtDataBdtReferenceIdPutParamOpts{
						BdtData: PolicyDataBdtData,
					}
					_, err := client.DefaultApi.PolicyDataBdtDataBdtReferenceIdPut(context.Background(), BdtPolicyStore[key].BdtPolData.BdtRefId, &data)
					if err != nil {
						logger.Bdtpolicylog.Warnln("UDR Create bdtdate error")
					}

					pcf_message.SendHttpResponseMessage(httpChannel, nil, 200, BdtPolicyStore[key])
					return
				}
			}
		}
	}
	// not found
	problem.Status = 404
	problem.Cause = "CONTEXT_NOT_FOUND"
	pcf_message.SendHttpResponseMessage(httpChannel, nil, 404, problem)
}

//CreateBDTPolicy - Create a new Individual BDT policy
func CreateBDTPolicyContext(httpChannel chan pcf_message.HttpResponseMessage, body models.BdtReqData) {
	var bdtReqData models.BdtReqData = body
	logger.Bdtpolicylog.Traceln("bdtReqData to store: ", bdtReqData)
	var bdtPolicyData models.BdtPolicyData
	var problem models.ProblemDetails
	var bdtPolicy models.BdtPolicy
	client := pcf_util.GetNudrClient("https://localhost:29504")
	NeedPolicy := true
	// check request essential IE
	if (bdtReqData.AspId != "") && (bdtReqData.NumOfUes != 0) {
		bdtPolicyData.BdtRefId = pcf_context.DefaultBdtRefId + bdtReqData.AspId

		//query pcfUeContext bdtpolicy and check pcfUeContext avoid rePolicy
		key := bdtReqData.AspId
		if BdtPolicyStore[key] != nil {
			for i := 0; i < len(BdtPolicyStore[key].BdtPolData.TransfPolicies); i++ {
				if BdtPolicyStore[key].BdtPolData.SelTransPolicyId == BdtPolicyStore[key].BdtPolData.TransfPolicies[i].TransPolicyId {
					//check query pcfUeContext bdtpolicy stoptime
					if BdtPolicyStore[key].BdtPolData.TransfPolicies[i].RecTimeInt.StopTime.After(time.Now()) {
						NeedPolicy = true
					} else {
						NeedPolicy = false
					}
				}
			}
			if !NeedPolicy {
				var uri = pcf_util.PCF_BASIC_PATH + pcf_context.BdtPolicyUri + BdtPolicyStore[key].BdtPolData.BdtRefId
				respHeader := make(http.Header)
				respHeader.Set("Location", uri)
				pcf_message.SendHttpResponseMessage(httpChannel, respHeader, 303, nil)
				return
			}
		}
		// pcf buffer has no data about this request and query udr bdtpolicy data avoid re policy
		bdtdata, _, err := client.DefaultApi.PolicyDataBdtDataGet(context.Background())
		if err == nil && bdtdata != nil {
			// query found bdtpolicy data
			bdtPolicyData.SelTransPolicyId = 1
			bdtPolicyData.SuppFeat = bdtReqData.SuppFeat
			for i := 0; i < len(bdtdata); i++ {
				if bdtReqData.AspId == bdtdata[i].AspId {
					// found
					bdtPolicyData.SelTransPolicyId = 1
					bdtPolicyData.SuppFeat = bdtReqData.SuppFeat
					bdtPolicyData.BdtRefId = bdtdata[i].BdtRefId
					bdtPolicyData.TransfPolicies[0] = bdtdata[i].TransPolicy
					break
				}
			}
			// check query data stoptime
			if bdtPolicyData.TransfPolicies[0].RecTimeInt.StopTime.After(time.Now()) {
				NeedPolicy = false
			}

			if !NeedPolicy {
				// query data no problem and save bdtpolicy data for pcfUeContext
				bdtPolicy.BdtPolData = &bdtPolicyData
				bdtPolicy.BdtReqData = &bdtReqData
				BdtPolicyStore[key] = &bdtPolicy
				pcf_message.SendHttpResponseMessage(httpChannel, nil, 201, bdtPolicyData)

				return
			}
		}
		// not found or timeout,make policy
		if NeedPolicy {
			if bdtReqData.DesTimeInt == nil {
				*bdtReqData.DesTimeInt = pcf_util.GetDefaultTime()
			}
			if bdtReqData.VolPerUe == nil {
				*bdtReqData.VolPerUe = pcf_util.GetDefaultDataRate()
			}
			StartTime, _ := pcf_util.TimeParse(*bdtReqData.DesTimeInt.StartTime)
			StopTime, _ := pcf_util.TimeParse(*bdtReqData.DesTimeInt.StopTime)
			bdtPolicyData.SelTransPolicyId = 1
			bdtPolicyData.SuppFeat = bdtReqData.SuppFeat

			bdtPolicyData.TransfPolicies = []models.TransferPolicy{
				{
					//MaxBitRateUl: pcf_util.Convert(bdtReqData.VolPerUe.UplinkVolume), //option
					//MaxBitRateDl: pcf_util.Convert(bdtReqData.VolPerUe.DownlinkVolume), //option
					RatingGroup: 1,
					RecTimeInt: &models.TimeWindow{
						StartTime: &StartTime,
						StopTime:  &StopTime,
					},
					TransPolicyId: 1,
				},
			}
		}

		//save bdtpolicy data for pcfUeContext
		bdtPolicy.BdtPolData = &bdtPolicyData
		bdtPolicy.BdtReqData = &bdtReqData
		BdtPolicyStore[key] = &bdtPolicy
		pcf_message.SendHttpResponseMessage(httpChannel, nil, 201, bdtPolicyData)

		// Udr Create bdtpolicy
		client := pcf_util.GetNudrClient("https://localhost:29504")
		var bdtData models.BdtData
		bdtData.AspId = bdtReqData.AspId
		bdtData.BdtRefId = BdtPolicyStore[key].BdtPolData.BdtRefId
		bdtData.TransPolicy = BdtPolicyStore[key].BdtPolData.TransfPolicies[0]
		var PolicyDataBdtData optional.Interface
		PolicyDataBdtData.Default(bdtData)
		data := Nudr_DataRepository.PolicyDataBdtDataBdtReferenceIdPutParamOpts{
			BdtData: PolicyDataBdtData,
		}
		_, err = client.DefaultApi.PolicyDataBdtDataBdtReferenceIdPut(context.Background(), BdtPolicyStore[key].BdtPolData.BdtRefId, &data)
		if err != nil {
			logger.Bdtpolicylog.Warnln("UDR Create bdtdate error")
		}

		return

	} else {
		problem.Status = 404
		problem.Cause = "CONTEXT_NOT_FOUND"
		pcf_message.SendHttpResponseMessage(httpChannel, nil, 404, problem)
		return
	}
}
