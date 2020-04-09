package smf_handler

import (
	"free5gc/lib/http_wrapper"
	"free5gc/lib/openapi/models"
	"free5gc/src/smf/smf_handler/smf_message"
	"free5gc/src/smf/smf_pfcp"
	"free5gc/src/smf/smf_producer"
	"net/http"
	"time"
)

func Handle() {

	for {
		select {
		case msg, ok := <-smf_message.SmfChannel:
			if ok {
				switch msg.Event {
				case smf_message.PFCPMessage:
					smf_pfcp.Dispatch(msg.PFCPRequest)
				case smf_message.PDUSessionSMContextCreate:
					smf_producer.HandlePDUSessionSMContextCreate(msg.ResponseChan, msg.HTTPRequest.Body.(models.PostSmContextsRequest))
				case smf_message.PDUSessionSMContextUpdate:
					smContextRef := msg.HTTPRequest.Params["smContextRef"]
					seqNum, ResBody := smf_producer.HandlePDUSessionSMContextUpdate(
						msg.ResponseChan, smContextRef, msg.HTTPRequest.Body.(models.UpdateSmContextRequest))
					response := http_wrapper.Response{
						Status: http.StatusOK,
						Body:   ResBody,
					}
					smf_message.RspQueue.PutItem(seqNum, msg.ResponseChan, response)
				case smf_message.PDUSessionSMContextRelease:
					smContextRef := msg.HTTPRequest.Params["smContextRef"]
					seqNum := smf_producer.HandlePDUSessionSMContextRelease(
						msg.ResponseChan, smContextRef, msg.HTTPRequest.Body.(models.ReleaseSmContextRequest))
					response := http_wrapper.Response{
						Status: http.StatusNoContent,
						Body:   nil,
					}
					smf_message.RspQueue.PutItem(seqNum, msg.ResponseChan, response)
				case smf_message.OAMGetUEPDUSessionInfo:
					smContextRef := msg.HTTPRequest.Params["smContextRef"]
					smf_producer.HandleOAMGetUEPDUSessionInfo(msg.ResponseChan, smContextRef)
				}
			}
		case <-time.After(time.Second * 1):
		}
	}
}
