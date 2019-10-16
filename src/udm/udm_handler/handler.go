package udm_handler

import (
	// "fmt"
	"free5gc/lib/openapi/models"
	"free5gc/src/udm/logger"
	"free5gc/src/udm/udm_handler/udm_message"
	"free5gc/src/udm/udm_producer"
	"time"

	"github.com/sirupsen/logrus"
)

const (
	MaxChannel int = 100000
)

var udmChannel chan udm_message.HandlerMessage
var HandlerLog *logrus.Entry

func init() {

	HandlerLog = logger.Handlelog
	udmChannel = make(chan udm_message.HandlerMessage, MaxChannel)
}

func SendMessage(msg udm_message.HandlerMessage) {
	udmChannel <- msg
}

func Handle() {
	for {
		select {
		case msg, ok := <-udmChannel:
			if ok {
				switch msg.Event {
				case udm_message.EventGenerateAuthData:
					supiOrSuci := msg.HTTPRequest.Params["supiOrSuci"]
					udm_producer.HandleGenerateAuthData(msg.ResponseChan, supiOrSuci,
						msg.HTTPRequest.Body.(models.AuthenticationInfoRequest))
				case udm_message.EventConfirmAuth:
					supi := msg.HTTPRequest.Params["supi"]
					udm_producer.HandleConfirmAuthData(msg.ResponseChan, supi,
						msg.HTTPRequest.Body.(models.AuthEvent))
				case udm_message.EventGetAmData:
					supi := msg.HTTPRequest.Params["supi"]
					plmnID := msg.HTTPRequest.Query.Get("plmn-id")
					udm_producer.HandleGetAmData(msg.ResponseChan, supi, plmnID)
				case udm_message.EventGetIdTranslationResult:
					gpsi := msg.HTTPRequest.Params["gpsi"]
					supportedFeatures := msg.HTTPRequest.Query.Get("supported-features")
					_ = supportedFeatures
					udm_producer.HandleGetIdTranslationResult(msg.ResponseChan, gpsi)
				case udm_message.EventGetNssai:
					supi := msg.HTTPRequest.Params["supi"]
					plmnID := msg.HTTPRequest.Query.Get("plmn-id")
					udm_producer.HandleGetNssai(msg.ResponseChan, supi, plmnID)

				case udm_message.EventGetSharedData:
					sharedDataIds := msg.HTTPRequest.Query["sharedDataIds"]
					udm_producer.HandleGetSharedData(msg.ResponseChan, sharedDataIds)

				case udm_message.EventGetSmData:
					supi := msg.HTTPRequest.Params["supi"]
					plmnID := msg.HTTPRequest.Query.Get("plmn-id")
					udm_producer.HandleGetSmData(msg.ResponseChan, supi, plmnID)
				case udm_message.EventGetSmfSelectData:
					supi := msg.HTTPRequest.Params["supi"]
					plmnID := msg.HTTPRequest.Query.Get("plmn-id")
					udm_producer.HandleGetSmfSelectData(msg.ResponseChan, supi, plmnID)

				case udm_message.EventGetSupi:
					supi := msg.HTTPRequest.Params["supi"]
					plmnID := msg.HTTPRequest.Query.Get("plmn-id")
					// dataSetNames := msg.HTTPRequest.Query["dataset-names"]
					udm_producer.HandleGetSupi(msg.ResponseChan, supi, plmnID)
				case udm_message.EventGetTraceData:
					supi := msg.HTTPRequest.Params["supi"]
					plmnID := msg.HTTPRequest.Query.Get("plmn-id")
					udm_producer.HandleGetTraceData(msg.ResponseChan, supi, plmnID)
				case udm_message.EventGetUeContextInSmfData:
					supi := msg.HTTPRequest.Params["supi"]
					udm_producer.HandleGetUeContextInSmfData(msg.ResponseChan, supi)
				case udm_message.EventSubscribe:
					supi := msg.HTTPRequest.Params["supi"]
					subscriptionID := msg.HTTPRequest.Params["subscriptionId"]
					udm_producer.HandleSubscribe(msg.ResponseChan, supi, subscriptionID, msg.HTTPRequest.Body.(models.SdmSubscription))
				case udm_message.EventSubscribeToSharedData:
					udm_producer.HandleSubscribeToSharedData(msg.ResponseChan, msg.HTTPRequest.Body.(models.SdmSubscription))
				case udm_message.EventUnsubscribe:
					supi := msg.HTTPRequest.Params["supi"]
					subscriptionID := msg.HTTPRequest.Params["subscriptionId"]
					udm_producer.HandleUnsubscribe(msg.ResponseChan, supi, subscriptionID)
				case udm_message.EventUnsubscribeForSharedData:
					subscriptionID := msg.HTTPRequest.Params["subscriptionId"]
					udm_producer.HandleUnsubscribeForSharedData(msg.ResponseChan, subscriptionID)
				case udm_message.EventModify:
					supi := msg.HTTPRequest.Params["supi"]
					subscriptionID := msg.HTTPRequest.Params["subscriptionId"]
					udm_producer.HandleModify(msg.ResponseChan, supi, subscriptionID, msg.HTTPRequest.Body.(models.SdmSubscription))
				case udm_message.EventModifyForSharedData:
					supi := msg.HTTPRequest.Params["supi"]
					subscriptionID := msg.HTTPRequest.Params["subscriptionId"]
					udm_producer.HandleModifyForSharedData(msg.ResponseChan, supi, subscriptionID)
				case udm_message.EventCreateEeSubscription:
					ueIdentity := msg.HTTPRequest.Params["ueIdentity"]
					subscriptionID := msg.HTTPRequest.Params["subscriptionId"]
					udm_producer.HandleCreateEeSubscription(msg.ResponseChan, ueIdentity, subscriptionID, msg.HTTPRequest.Body.(models.EeSubscription))
				case udm_message.EventDeleteEeSubscription:
					ueIdentity := msg.HTTPRequest.Params["ueIdentity"]
					subscriptionID := msg.HTTPRequest.Params["subscriptionId"]
					udm_producer.HandleDeleteEeSubscription(msg.ResponseChan, ueIdentity, subscriptionID)
				case udm_message.EventUpdateEeSubscription:
					ueIdentity := msg.HTTPRequest.Params["ueIdentity"]
					subscriptionID := msg.HTTPRequest.Params["subscriptionId"]
					udm_producer.HandleUpdateEeSubscription(msg.ResponseChan, ueIdentity, subscriptionID)
				case udm_message.EventGetAmf3gppAccess:
					ueID := msg.HTTPRequest.Params["ueId"]
					supportedFeatures := msg.HTTPRequest.Query.Get("supported-features")
					udm_producer.HandleGetAmf3gppAccess(msg.ResponseChan, ueID, supportedFeatures)
				case udm_message.EventGetAmfNon3gppAccess:
					ueID := msg.HTTPRequest.Params["ueId"]
					supportedFeatures := msg.HTTPRequest.Query.Get("supported-features")
					udm_producer.HandleGetAmfNon3gppAccess(msg.ResponseChan, ueID, supportedFeatures)
				case udm_message.EventRegistrationAmf3gppAccess:
					ueID := msg.HTTPRequest.Params["ueId"]
					udm_producer.HandleRegistrationAmf3gppAccess(msg.ResponseChan, ueID, msg.HTTPRequest.Body.(models.Amf3GppAccessRegistration))
				case udm_message.EventRegisterAmfNon3gppAccess:
					ueID := msg.HTTPRequest.Params["ueId"]
					udm_producer.HandleRegisterAmfNon3gppAccess(msg.ResponseChan, ueID, msg.HTTPRequest.Body.(models.AmfNon3GppAccessRegistration))
				case udm_message.EventUpdateAmf3gppAccess:
					ueID := msg.HTTPRequest.Params["ueId"]
					udm_producer.HandleUpdateAmf3gppAccess(msg.ResponseChan, ueID, msg.HTTPRequest.Body.(models.Amf3GppAccessRegistrationModification))
				case udm_message.EventUpdateAmfNon3gppAccess:
					ueID := msg.HTTPRequest.Params["ueId"]
					udm_producer.HandleUpdateAmfNon3gppAccess(msg.ResponseChan, ueID, msg.HTTPRequest.Body.(models.AmfNon3GppAccessRegistrationModification))
				case udm_message.EventDeregistrationSmfRegistrations:
					ueID := msg.HTTPRequest.Params["ueId"]
					pduSessionID := msg.HTTPRequest.Params["pduSessionId"]
					udm_producer.HandleDeregistrationSmfRegistrations(msg.ResponseChan, ueID, pduSessionID)
				case udm_message.EventRegistrationSmfRegistrations:
					ueID := msg.HTTPRequest.Params["ueId"]
					pduSessionID := msg.HTTPRequest.Params["pduSessionId"]
					udm_producer.HandleRegistrationSmfRegistrations(msg.ResponseChan, ueID, pduSessionID, msg.HTTPRequest.Body.(models.SmfRegistration))
				case udm_message.EventUpdate:
					gpsi := msg.HTTPRequest.Params["gpsi"]
					udm_producer.HandleUpdate(msg.ResponseChan, gpsi, msg.HTTPRequest.Body.(models.PpData))
				default:
					HandlerLog.Warnf("Event[%d] has not implemented", msg.Event)
				}

			} else {
				HandlerLog.Errorln("UDM Channel closed!")
			}

		case <-time.After(time.Second * 1):
		}
	}
}
