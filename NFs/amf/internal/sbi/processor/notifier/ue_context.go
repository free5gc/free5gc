package callback

import (
	"context"
	"fmt"

	amf_context "github.com/free5gc/amf/internal/context"
	Namf_Communication "github.com/free5gc/openapi/amf/Communication"
	"github.com/free5gc/openapi/models"
)

func SendN2InfoNotifyN2Handover(ue *amf_context.AmfUe, releaseList []int32) error {
	if ue.HandoverNotifyUri == "" {
		return fmt.Errorf("N2 Info Notify N2Handover failed(uri dose not exist)")
	}
	configuration := Namf_Communication.NewConfiguration()
	client := Namf_Communication.NewAPIClient(configuration)

	n2InformationNotification := models.N2InformationNotification{
		N2NotifySubscriptionId: ue.Supi,
		ToReleaseSessionList:   releaseList,
		NotifyReason:           models.N2InfoNotifyReason_HANDOVER_COMPLETED,
	}

	n2InformationNotificationReq := Namf_Communication.N2InfoNotifyHandoverCompleteRequest{
		N2InformationNotification: &n2InformationNotification,
	}

	_, err := client.IndividualUeContextDocumentApi.
		N2InfoNotifyHandoverComplete(context.Background(), ue.HandoverNotifyUri, &n2InformationNotificationReq)

	if err == nil {
		// TODO: handle Msg
	} else {
		HttpLog.Errorln(err.Error())
		return err
	}
	return nil
}
