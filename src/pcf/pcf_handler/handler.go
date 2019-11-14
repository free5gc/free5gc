package pcf_handler

import (
	"free5gc/lib/openapi/models"
	"free5gc/src/pcf/logger"
	"free5gc/src/pcf/pcf_handler/pcf_message"
	"free5gc/src/pcf/pcf_producer"
	"time"

	"github.com/sirupsen/logrus"
)

var HandlerLog *logrus.Entry

func init() {
	// init Pool
	HandlerLog = logger.HandlerLog
}

func Handle() {
	for {
		select {
		case msg, ok := <-pcf_message.PCFChannel:
			if ok {
				switch msg.Event {
				case pcf_message.EventBDTPolicyCreate:
					pcf_producer.CreateBDTPolicyContext(msg.HttpChannel, msg.HTTPRequest.Body.(models.BdtReqData))
				case pcf_message.EventBDTPolicyGet:
					BdtPolicyId := msg.HTTPRequest.Params["bdtPolicyId"]
					pcf_producer.GetBDTPolicyContext(msg.HttpChannel, BdtPolicyId)
				case pcf_message.EventBDTPolicyUpdate:
					BdtPolicyId := msg.HTTPRequest.Params["bdtPolicyId"]
					pcf_producer.UpdateBDTPolicyContext(msg.HttpChannel, BdtPolicyId, msg.HTTPRequest.Body.(models.BdtPolicyDataPatch))
				case pcf_message.EventPostAppSessions:
					ReqURI := msg.HTTPRequest.Params["ReqURI"]
					pcf_producer.PostAppSessionsContext(msg.HttpChannel, ReqURI, msg.HTTPRequest.Body.(models.AppSessionContext))
				case pcf_message.EventGetAppSession:
					ReqURI := msg.HTTPRequest.Params["ReqURI"]
					pcf_producer.GetAppSessionContext(msg.HttpChannel, ReqURI)
				case pcf_message.EventDeleteAppSession:
					ReqURI := msg.HTTPRequest.Params["ReqURI"]
					pcf_producer.DeleteAppSessionContext(msg.HttpChannel, ReqURI)
				case pcf_message.EventModAppSession:
					ReqURI := msg.HTTPRequest.Params["ReqURI"]
					pcf_producer.ModAppSessionContext(msg.HttpChannel, ReqURI, msg.HTTPRequest.Body.(models.AppSessionContextUpdateData))
				case pcf_message.EventDeleteEventsSubsc:
					ReqURI := msg.HTTPRequest.Params["ReqURI"]
					pcf_producer.DeleteEventsSubscContext(msg.HttpChannel, ReqURI)
				case pcf_message.EventUpdateEventsSubsc:
					ReqURI := msg.HTTPRequest.Params["ReqURI"]
					pcf_producer.UpdateEventsSubscContext(msg.HttpChannel, ReqURI, msg.HTTPRequest.Body.(models.EventsSubscReqData))
				case pcf_message.EventAMPolicyGet:
					PolAssoId := msg.HTTPRequest.Params["polAssoId"]
					pcf_producer.GetPoliciesPolAssoId(msg.HttpChannel, PolAssoId)
				case pcf_message.EventAMPolicyDelete:
					PolAssoId := msg.HTTPRequest.Params["polAssoId"]
					pcf_producer.DeletePoliciesPolAssoId(msg.HttpChannel, PolAssoId)
				case pcf_message.EventAMPolicyCreate:
					ReqURI := msg.HTTPRequest.Params["ReqURI"]
					pcf_producer.PostPolicies(msg.HttpChannel, ReqURI, msg.HTTPRequest.Body.(models.PolicyAssociationRequest))
				case pcf_message.EventAMPolicyUpdate:
					ReqURI := msg.HTTPRequest.Params["ReqURI"]
					PolAssoId := msg.HTTPRequest.Params["polAssoId"]
					pcf_producer.UpdatePostPoliciesPolAssoId(msg.HttpChannel, ReqURI, PolAssoId, msg.HTTPRequest.Body.(models.PolicyAssociationUpdateRequest))
				case pcf_message.EventSMPolicyCreate:
					ReqURI := msg.HTTPRequest.Params["ReqURI"]
					pcf_producer.CreateSmPolicyContext(msg.HttpChannel, ReqURI, msg.HTTPRequest.Body.(models.SmPolicyContextData))
				case pcf_message.EventSMPolicyGet:
					ReqURI := msg.HTTPRequest.Params["ReqURI"]
					pcf_producer.GetSmPolicyContext(msg.HttpChannel, ReqURI)
				case pcf_message.EventSMPolicyUpdate:
					ReqURI := msg.HTTPRequest.Params["ReqURI"]
					pcf_producer.UpdateSmPolicyContext(msg.HttpChannel, ReqURI, msg.HTTPRequest.Body.(models.SmPolicyUpdateContextData))
				case pcf_message.EventSMPolicyDelete:
					ReqURI := msg.HTTPRequest.Params["ReqURI"]
					pcf_producer.DeleteSmPolicyContext(msg.HttpChannel, ReqURI)
				case pcf_message.EventSMPolicyNotify:
					ReqURI := msg.HTTPRequest.Params["ReqURI"]
					pcf_producer.HandleSmPolicyNotify(msg.HttpChannel, ReqURI, msg.HTTPRequest.Body.(models.PolicyDataChangeNotification))
				// TODO: http event dispatcher
				default:
					HandlerLog.Warnf("Event[%s] has not implemented", msg.Event)
				}
			} else {
				HandlerLog.Errorln("Channel closed!")
			}

		case <-time.After(time.Second * 1):

		}
	}
}
