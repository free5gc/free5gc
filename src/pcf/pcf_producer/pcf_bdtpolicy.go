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

func GetBDTPolicyContext(httpChannel chan pcf_message.HttpResponseMessage, ReqURI string) {

	var problem models.ProblemDetails
	URI := ReqURI
	pcfUeContext := pcf_context.GetPCFUeContext()
	// check BdtRefId from pcfUecontext
	for key := range pcfUeContext {
		if pcfUeContext[key].BdtPolicyStore == nil {
			continue
		}
		BdtrefId_temp := pcf_context.BdtPolicyUri + pcfUeContext[key].BdtPolicyStore.BdtPolData.BdtRefId
		if URI == BdtrefId_temp {
			// check transpolicy stoptime
			for i := 0; i < len(pcfUeContext[key].BdtPolicyStore.BdtPolData.TransfPolicies); i++ {
				if pcfUeContext[key].BdtPolicyStore.BdtPolData.TransfPolicies[i].TransPolicyId == pcfUeContext[key].BdtPolicyStore.BdtPolData.SelTransPolicyId {
					if pcfUeContext[key].BdtPolicyStore.BdtPolData.TransfPolicies[i].RecTimeInt.StopTime.Before(time.Now()) {
						pcfUeContext[key].BdtPolicyTimeout = true
						break
					} else {
						pcfUeContext[key].BdtPolicyTimeout = false
						break
					}
				}
			}
			pcf_message.SendHttpResponseMessage(httpChannel, nil, 200, pcfUeContext[key].BdtPolicyStore)
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
	pcfUeContext := pcf_context.GetPCFUeContext()
	// check BdtRefId from pcfUeContext
	for key := range pcfUeContext {
		if pcfUeContext[key].BdtPolicyStore == nil {
			continue
		}
		BdtrefId_temp := pcf_context.BdtPolicyUri + pcfUeContext[key].BdtPolicyStore.BdtPolData.BdtRefId
		if URI == BdtrefId_temp {
			if bdtPolicyDataPatch.SelTransPolicyId == pcfUeContext[key].BdtPolicyStore.BdtPolData.SelTransPolicyId {
				pcfUeContext[key].BdtPolicyStore.BdtPolData.SelTransPolicyId = bdtPolicyDataPatch.SelTransPolicyId
				pcf_message.SendHttpResponseMessage(httpChannel, nil, 204, pcfUeContext[key].BdtPolicyStore)
				return
			}
			for i := 0; i < len(pcfUeContext[key].BdtPolicyStore.BdtPolData.TransfPolicies); i++ {
				if bdtPolicyDataPatch.SelTransPolicyId == pcfUeContext[key].BdtPolicyStore.BdtPolData.TransfPolicies[i].TransPolicyId {
					transPolicy := pcfUeContext[key].BdtPolicyStore.BdtPolData.TransfPolicies[i]
					if pcfUeContext[key].BdtPolicyStore.BdtPolData.TransfPolicies[i].RecTimeInt.StopTime.Before(time.Now()) {
						pcfUeContext[key].BdtPolicyTimeout = true
					} else {
						pcfUeContext[key].BdtPolicyTimeout = false
					}
					pcfUeContext[key].BdtPolicyStore.BdtPolData.SelTransPolicyId = bdtPolicyDataPatch.SelTransPolicyId
					// update bdtdata to udr
					client := pcf_util.GetNudrClient()
					var bdtData models.BdtData
					bdtData.AspId = pcfUeContext[key].AspId
					bdtData.BdtRefId = pcfUeContext[key].BdtPolicyStore.BdtPolData.BdtRefId
					bdtData.TransPolicy = transPolicy
					var PolicyDataBdtData optional.Interface
					PolicyDataBdtData.Default(bdtData)
					data := Nudr_DataRepository.PolicyDataBdtDataBdtReferenceIdPutParamOpts{
						BdtData: PolicyDataBdtData,
					}
					_, err := client.DefaultApi.PolicyDataBdtDataBdtReferenceIdPut(context.Background(), pcfUeContext[key].BdtPolicyStore.BdtPolData.BdtRefId, &data)
					if err != nil {
						logger.Bdtpolicylog.Warnln("UDR Create bdtdate error")
					}

					pcf_message.SendHttpResponseMessage(httpChannel, nil, 200, pcfUeContext[key].BdtPolicyStore)
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
func CreateBDTPolicyContext(httpChannel chan pcf_message.HttpResponseMessage, ReqURI string, body models.BdtReqData) {
	var bdtReqData models.BdtReqData = body
	logger.Bdtpolicylog.Traceln("bdtReqData to store: ", bdtReqData)
	var bdtPolicyData models.BdtPolicyData
	var problem models.ProblemDetails
	var bdtPolicy models.BdtPolicy
	pcfUeContext := pcf_context.GetPCFUeContext()
	client := pcf_util.GetNudrClient()
	NeedPolicy := true
	// check request essential IE
	if (bdtReqData.AspId != "") && (bdtReqData.NumOfUes != 0) {
		bdtPolicyData.BdtRefId = pcf_context.DefaultBdtRefId + bdtReqData.AspId

		// check aspid on pcfuecontext
		key, err := pcf_context.CheckAspidOnPcfUeContext(bdtReqData.AspId)
		if err != nil {
			problem.Status = 404
			problem.Cause = "CONTEXT_NOT_FOUND"
			pcf_message.SendHttpResponseMessage(httpChannel, nil, 404, problem)
			return
		}

		//query pcfUeContext bdtpolicy and check pcfUeContext avoid rePolicy
		if pcfUeContext[key].BdtPolicyStore != nil {
			for i := 0; i < len(pcfUeContext[key].BdtPolicyStore.BdtPolData.TransfPolicies); i++ {
				if pcfUeContext[key].BdtPolicyStore.BdtPolData.SelTransPolicyId == pcfUeContext[key].BdtPolicyStore.BdtPolData.TransfPolicies[i].TransPolicyId {
					//check query pcfUeContext bdtpolicy stoptime
					if pcfUeContext[key].BdtPolicyStore.BdtPolData.TransfPolicies[i].RecTimeInt.StopTime.Before(time.Now()) {
						pcfUeContext[key].BdtPolicyTimeout = true
					} else {
						pcfUeContext[key].BdtPolicyTimeout = false
					}
				}
			}
			if !pcfUeContext[key].BdtPolicyTimeout {
				var uri = pcf_util.PCF_BASIC_PATH + pcf_context.BdtPolicyUri + pcfUeContext[key].BdtPolicyStore.BdtPolData.BdtRefId
				respHeader := make(http.Header)
				respHeader.Set("Location", uri)
				pcf_message.SendHttpResponseMessage(httpChannel, respHeader, 303, nil)
				return
			}
		}
		// pcf buffer has no data about this request and query udr bdtpolicy data avoid re policy
		bdtdata, _, err := client.DefaultApi.PolicyDataBdtDataGet(context.Background())
		if err == nil || bdtdata != nil {
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
			if bdtPolicyData.TransfPolicies[0].RecTimeInt.StopTime.Before(time.Now()) {
				NeedPolicy = false
			}

			if !NeedPolicy {
				// query data no problem and save bdtpolicy data for pcfUeContext
				bdtPolicy.BdtPolData = &bdtPolicyData
				bdtPolicy.BdtReqData = &bdtReqData
				pcfUeContext[key].BdtPolicyStore = &bdtPolicy
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
		pcfUeContext[key].BdtPolicyStore = &bdtPolicy
		pcf_message.SendHttpResponseMessage(httpChannel, nil, 201, bdtPolicyData)

		// Udr Create bdtpolicy
		client := pcf_util.GetNudrClient()
		var bdtData models.BdtData
		bdtData.AspId = pcfUeContext[key].AspId
		bdtData.BdtRefId = pcfUeContext[key].BdtPolicyStore.BdtPolData.BdtRefId
		bdtData.TransPolicy = pcfUeContext[key].BdtPolicyStore.BdtPolData.TransfPolicies[0]
		var PolicyDataBdtData optional.Interface
		PolicyDataBdtData.Default(bdtData)
		data := Nudr_DataRepository.PolicyDataBdtDataBdtReferenceIdPutParamOpts{
			BdtData: PolicyDataBdtData,
		}
		_, err = client.DefaultApi.PolicyDataBdtDataBdtReferenceIdPut(context.Background(), pcfUeContext[key].BdtPolicyStore.BdtPolData.BdtRefId, &data)
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
