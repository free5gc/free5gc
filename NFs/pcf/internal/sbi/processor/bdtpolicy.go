package processor

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/mohae/deepcopy"

	"github.com/free5gc/openapi/models"
	"github.com/free5gc/openapi/nrf/NFDiscovery"
	"github.com/free5gc/openapi/udr/DataRepository"
	pcf_context "github.com/free5gc/pcf/internal/context"
	"github.com/free5gc/pcf/internal/logger"
	"github.com/free5gc/pcf/internal/util"
)

func (p *Processor) HandleGetBDTPolicyContextRequest(
	c *gin.Context,
	bdtPolicyID string,
) {
	// step 1: log
	logger.BdtPolicyLog.Infof("Handle GetBDTPolicyContext")

	// step 2: handle the message
	logger.BdtPolicyLog.Traceln("Handle BDT Policy GET")
	// check bdtPolicyID from pcfUeContext
	if value, ok := p.Context().BdtPolicyPool.Load(bdtPolicyID); ok {
		bdtPolicy := value.(*models.BdtPolicy)
		c.JSON(http.StatusOK, bdtPolicy)
		return
	} else {
		// not found
		problemDetails := util.GetProblemDetail("Can't find bdtPolicyID related resource", util.CONTEXT_NOT_FOUND)
		logger.BdtPolicyLog.Warn(problemDetails.Detail)
		c.JSON(int(problemDetails.Status), problemDetails)
		return
	}
}

// UpdateBDTPolicy - Update an Individual BDT policy (choose policy data)
func (p *Processor) HandleUpdateBDTPolicyContextProcedure(
	c *gin.Context,
	bdtPolicyID string,
	bdtPolicyDataPatch models.PcfBdtPolicyControlBdtPolicyDataPatch,
) {
	// step 1: log
	logger.BdtPolicyLog.Infof("Handle UpdateBDTPolicyContext")

	// step 2: handle the message
	logger.BdtPolicyLog.Infoln("Handle BDTPolicyUpdate")
	// check bdtPolicyID from pcfUeContext
	pcfSelf := p.Context()

	var bdtPolicy *models.BdtPolicy
	if value, ok := pcfSelf.BdtPolicyPool.Load(bdtPolicyID); ok {
		bdtPolicy = value.(*models.BdtPolicy)
	} else {
		// not found
		problemDetail := util.GetProblemDetail("Can't find bdtPolicyID related resource", util.CONTEXT_NOT_FOUND)
		logger.BdtPolicyLog.Warn(problemDetail.Detail)
		c.JSON(int(problemDetail.Status), problemDetail)
		return
	}

	for _, policy := range bdtPolicy.BdtPolData.TransfPolicies {
		if policy.TransPolicyId == bdtPolicyDataPatch.SelTransPolicyId {
			polData := bdtPolicy.BdtPolData
			polReq := bdtPolicy.BdtReqData
			polData.SelTransPolicyId = bdtPolicyDataPatch.SelTransPolicyId
			bdtData := models.BdtData{
				AspId:       polReq.AspId,
				TransPolicy: &policy,
				BdtRefId:    polData.BdtRefId,
			}
			if polReq.NwAreaInfo != nil {
				bdtData.NwAreaInfo = polReq.NwAreaInfo
			}

			udrUri := p.getDefaultUdrUri(pcfSelf)
			if udrUri == "" {
				// Can't find any UDR support this Ue
				pd := &models.ProblemDetails{
					Status: http.StatusServiceUnavailable,
					Detail: "Can't find any UDR which supported to this PCF",
				}
				logger.BdtPolicyLog.Warn(pd.Detail)
				c.JSON(int(pd.Status), pd)
				return
			}
			pd, err := p.Consumer().CreateBdtData(udrUri, &bdtData)
			if err != nil {
				logger.BdtPolicyLog.Warnf("UDR Put BdtDate error[%s]", err.Error())
				c.JSON(http.StatusInternalServerError, err.Error())
				return
			} else if pd != nil {
				logger.BdtPolicyLog.Warnf("UDR Put BdtDate fault[%s]", pd.Detail)
				c.JSON(int(pd.Status), pd)
				return
			}

			logger.BdtPolicyLog.Tracef("bdtPolicyID[%s] has Updated with SelTransPolicyId[%d]",
				bdtPolicyID, bdtPolicyDataPatch.SelTransPolicyId)
			c.JSON(http.StatusOK, bdtPolicy)
			return
		}
	}
	problemDetail := util.GetProblemDetail(
		fmt.Sprintf("Can't find TransPolicyId[%d] in TransfPolicies with bdtPolicyID[%s]",
			bdtPolicyDataPatch.SelTransPolicyId, bdtPolicyID),
		util.CONTEXT_NOT_FOUND)
	logger.BdtPolicyLog.Warn(problemDetail.Detail)
	c.JSON(int(problemDetail.Status), problemDetail)
}

// CreateBDTPolicy - Create a new Individual BDT policy
func (p *Processor) HandleCreateBDTPolicyContextRequest(
	c *gin.Context,
	requestMsg models.BdtReqData,
) {
	// step 1: log
	logger.BdtPolicyLog.Infof("Handle CreateBdtPolicyContext")

	var problemDetails *models.ProblemDetails

	// step 2: retrieve request and check mandatory contents
	if requestMsg.AspId == "" || requestMsg.DesTimeInt == nil || requestMsg.NumOfUes == 0 || requestMsg.VolPerUe == nil {
		logger.BdtPolicyLog.Errorf("Required BdtReqData not found: AspId[%+v], DesTimeInt[%+v], NumOfUes[%+v], VolPerUe[%+v]",
			requestMsg.AspId, requestMsg.DesTimeInt, requestMsg.NumOfUes, requestMsg.VolPerUe)
		c.JSON(http.StatusNotFound, nil)
		return
	}

	// // step 3: handle the message

	response := &models.BdtPolicy{}
	logger.BdtPolicyLog.Traceln("Handle BDT Policy Create")

	pcfSelf := p.Context()
	udrUri := p.getDefaultUdrUri(pcfSelf)
	if udrUri == "" {
		// Can't find any UDR support this Ue
		problemDetails = &models.ProblemDetails{
			Status: http.StatusServiceUnavailable,
			Detail: "Can't find any UDR which supported to this PCF",
		}
		logger.BdtPolicyLog.Warn(problemDetails.Detail)
		c.JSON(int(problemDetails.Status), problemDetails)
		return
	}
	pcfSelf.DefaultUdrURI = udrUri

	// Query BDT DATA array from UDR
	req := DataRepository.ReadBdtDataRequest{}
	resp, problemDetails, err := p.Consumer().CreateBdtPolicyContext(udrUri, &req)
	if err != nil {
		problemDetails = &models.ProblemDetails{
			Status: http.StatusServiceUnavailable,
			Detail: "Query to UDR failed",
		}
		logger.BdtPolicyLog.Warn("Query to UDR failed")
		c.JSON(int(problemDetails.Status), problemDetails)
		return
	} else if problemDetails != nil {
		c.JSON(int(problemDetails.Status), problemDetails)
		return
	}

	bdtDatas := resp.BdtData
	// TODO: decide BDT Policy from other bdt policy data
	response.BdtReqData = deepcopy.Copy(requestMsg).(*models.BdtReqData)
	var bdtData *models.BdtData
	var bdtPolicyData models.PcfBdtPolicyControlBdtPolicyData
	for _, data := range bdtDatas {
		// If ASP has exist, use its background data policy
		if requestMsg.AspId == data.AspId {
			bdtData = &data
			break
		}
	}
	// Only support one bdt policy, TODO: more policy for decision
	if bdtData != nil {
		// found
		// modify policy according to new request
		bdtData.TransPolicy.RecTimeInt = requestMsg.DesTimeInt
	} else {
		// use default bdt policy, TODO: decide bdt transfer data policy
		bdtData = &models.BdtData{
			AspId:       requestMsg.AspId,
			BdtRefId:    uuid.New().String(),
			TransPolicy: getDefaultTransferPolicy(1, *requestMsg.DesTimeInt),
		}
	}
	if requestMsg.NwAreaInfo != nil {
		bdtData.NwAreaInfo = requestMsg.NwAreaInfo
	}
	bdtPolicyData.SelTransPolicyId = bdtData.TransPolicy.TransPolicyId
	// no support feature in subclause 5.8 of TS29554
	bdtPolicyData.BdtRefId = bdtData.BdtRefId
	bdtPolicyData.TransfPolicies = append(bdtPolicyData.TransfPolicies, *bdtData.TransPolicy)
	response.BdtPolData = &bdtPolicyData
	bdtPolicyID, err := pcfSelf.AllocBdtPolicyID()
	if err != nil {
		problemDetails = &models.ProblemDetails{
			Status: http.StatusServiceUnavailable,
			Detail: "Allocate bdtPolicyID failed",
		}
		logger.BdtPolicyLog.Warn("Allocate bdtPolicyID failed")
		c.JSON(int(problemDetails.Status), problemDetails)
		return
	}

	pcfSelf.BdtPolicyPool.Store(bdtPolicyID, response)

	// Update UDR BDT Data(PUT)
	problemDetails, err = p.Consumer().CreateBdtData(udrUri, bdtData)
	if err != nil {
		logger.BdtPolicyLog.Warnf("UDR Put BdtDate error[%s]", err.Error())
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	} else if problemDetails != nil {
		logger.BdtPolicyLog.Warnf("UDR Put BdtDate fault[%s]", problemDetails.Detail)
		c.JSON(int(problemDetails.Status), problemDetails)
		return
	}

	locationHeader := util.GetResourceUri(models.ServiceName_NPCF_BDTPOLICYCONTROL, bdtPolicyID)
	logger.BdtPolicyLog.Tracef("BDT Policy Id[%s] Create", bdtPolicyID)

	c.Header("Location", locationHeader)
	c.JSON(http.StatusCreated, response)
}

func (p *Processor) getDefaultUdrUri(context *pcf_context.PCFContext) string {
	context.DefaultUdrURILock.RLock()
	defer context.DefaultUdrURILock.RUnlock()
	if context.DefaultUdrURI != "" {
		return context.DefaultUdrURI
	}
	param := NFDiscovery.SearchNFInstancesRequest{
		ServiceNames: []models.ServiceName{models.ServiceName_NUDR_DR},
	}
	resp, err := p.Consumer().SendSearchNFInstances(
		context.NrfUri,
		models.NrfNfManagementNfType_UDR,
		models.NrfNfManagementNfType_PCF,
		param,
	)
	if err != nil {
		return ""
	}
	for _, nfProfile := range resp.NfInstances {
		udruri := util.SearchNFServiceUri(
			nfProfile,
			models.ServiceName_NUDR_DR,
			models.NfServiceStatus_REGISTERED,
		)
		if udruri != "" {
			return udruri
		}
	}
	return ""
}

// get default background data transfer policy
func getDefaultTransferPolicy(
	transferPolicyId int32,
	timeWindow models.TimeWindow,
) *models.PcfBdtPolicyControlTransferPolicy {
	return &models.PcfBdtPolicyControlTransferPolicy{
		TransPolicyId: transferPolicyId,
		RecTimeInt:    &timeWindow,
		RatingGroup:   1,
	}
}
