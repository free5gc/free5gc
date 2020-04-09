package smf_consumer

import (
	"context"
	"free5gc/lib/Nsmf_PDUSession"
	"free5gc/lib/openapi/models"
	"free5gc/src/smf/logger"
	"net/http"
)

func SendSMContextStatusNotification(uri string) {
	if uri != "" {
		request := models.SmContextStatusNotification{}
		request.StatusInfo = &models.StatusInfo{
			ResourceStatus: models.ResourceStatus_RELEASED,
		}
		configuration := Nsmf_PDUSession.NewConfiguration()
		client := Nsmf_PDUSession.NewAPIClient(configuration)
		httpResponse, err := client.IndividualSMContextNotificationApi.SMContextNotification(context.Background(), uri, request)
		if err != nil {
			if httpResponse != nil {
				logger.PduSessLog.Warnf("Send SMContextStatus Notification Error[%s]", httpResponse.Status)
			} else {
				logger.PduSessLog.Warnf("Send SMContextStatus Notification Failed[%s]", err.Error())
			}
			return
		} else if httpResponse == nil {
			logger.PduSessLog.Warnln("Send SMContextStatus Notification Failed[HTTP Response is nil]")
			return
		}
		if httpResponse.StatusCode != http.StatusNoContent {
			logger.PduSessLog.Warnf("Send SMContextStatus Notification Failed")
		} else {
			logger.PduSessLog.Tracef("Send SMContextStatus Notification Success")
		}
	}
}
