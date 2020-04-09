package pcf_producer

import (
	"free5gc/lib/openapi/models"
	"free5gc/src/pcf/logger"
	"free5gc/src/pcf/pcf_handler/pcf_message"
	"net/http"
)

func HandleAmfStatusChangeNotify(httpChannel chan pcf_message.HttpResponseMessage, notification models.AmfStatusChangeNotification) {
	logger.CallbackLog.Warnf("[PCF] Handle Amf Status Change Notify is not implemented.")
	logger.CallbackLog.Debugf("receive AMF status change notification[%+v]", notification)
	// TODO: handle AMF Status Change Notify
	pcf_message.SendHttpResponseMessage(httpChannel, nil, http.StatusNoContent, nil)
}
