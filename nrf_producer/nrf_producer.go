package nrf_producer

import (
	"context"
	"free5gc/lib/Nnrf_NFManagement"
	"free5gc/lib/http_wrapper"
	"free5gc/lib/openapi/models"
	"free5gc/src/nrf/logger"
	"free5gc/src/nrf/nrf_handler/nrf_message"
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
		logger.ManagementLog.Info("Notify fail")
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
