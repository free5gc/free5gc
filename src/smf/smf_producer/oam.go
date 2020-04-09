package smf_producer

import (
	"free5gc/lib/http_wrapper"
	"free5gc/lib/openapi/models"
	"free5gc/src/smf/smf_context"
	"free5gc/src/smf/smf_handler/smf_message"
	"net/http"
	"strconv"
)

type PDUSessionInfo struct {
	Supi         string
	PDUSessionID string
	LocalSEID    string
	RemoteSEID   string
	Dnn          string
	Sst          string
	Sd           string
	AnType       models.AccessType
	PDUAddress   string
	SessionRule  models.SessionRule
	UpCnxState   models.UpCnxState
	Tunnel       smf_context.UPTunnel
}

func HandleOAMGetUEPDUSessionInfo(rspChan chan smf_message.HandlerResponseMessage, smContextRef string) {
	smContext := smf_context.GetSMContext(smContextRef)
	if smContext == nil {
		rspChan <- smf_message.HandlerResponseMessage{
			HTTPResponse: &http_wrapper.Response{
				Header: nil,
				Status: http.StatusNotFound,
				Body:   nil,
			},
		}
		return
	}

	rspChan <- smf_message.HandlerResponseMessage{
		HTTPResponse: &http_wrapper.Response{
			Header: nil,
			Status: http.StatusOK,
			Body: PDUSessionInfo{
				Supi:         smContext.Supi,
				PDUSessionID: strconv.Itoa(int(smContext.PDUSessionID)),
				LocalSEID:    strconv.FormatUint(smContext.LocalSEID, 10),
				RemoteSEID:   strconv.FormatUint(smContext.RemoteSEID, 10),
				Dnn:          smContext.Dnn,
				Sst:          strconv.Itoa(int(smContext.Snssai.Sst)),
				Sd:           smContext.Snssai.Sd,
				AnType:       smContext.AnType,
				PDUAddress:   smContext.PDUAddress.String(),
				SessionRule:  smContext.SessionRule,
				UpCnxState:   smContext.UpCnxState,
				Tunnel: smf_context.UPTunnel{
					//UpfRoot:  smContext.Tunnel.UpfRoot,
					ULCLRoot: smContext.Tunnel.ULCLRoot,
				},
			},
		},
	}
}
