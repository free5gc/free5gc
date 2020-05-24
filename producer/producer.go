package producer

import (
	"context"
	"free5gc/lib/http_wrapper"
	"free5gc/lib/openapi/Nnrf_NFManagement"
	"free5gc/lib/openapi/models"
	nrf_message "free5gc/src/nrf/handler/message"
	"free5gc/src/nrf/logger"
	"log"
	"net/http"
)

func HandleNotification(rspChan chan nrf_message.HandlerResponseMessage, url string, body models.NotificationData) {

	// Set client and set url
	configuration := Nnrf_NFManagement.NewConfiguration()
	//url = fmt.Sprintf("%s%s", url, "/notification")

	configuration.SetBasePathNoGroup(url)
	client := Nnrf_NFManagement.NewAPIClient(configuration)

	res, err := client.NotificationApi.NotificationPost(context.TODO(), body)
	if err != nil {
		logger.ManagementLog.Infof("Notify fail: %v", err)
		rspChan <- nrf_message.HandlerResponseMessage{
			HTTPResponse: &http_wrapper.Response{
				Header: nil,
				Status: http.StatusNoContent,
				Body:   "",
			},
		}
	}
	if res != nil {
		if status := res.StatusCode; status != http.StatusNoContent {
			log.Println("error: ", status)
		} else {
			rspChan <- nrf_message.HandlerResponseMessage{
				HTTPResponse: &http_wrapper.Response{
					Header: nil,
					Status: status,
					Body:   "",
				},
			}
		}
	}

}
