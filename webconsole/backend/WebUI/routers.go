package WebUI

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Route is the information for every URI.
type Route struct {
	// Name is the name of this Route.
	Name string
	// Method is the string for the HTTP method. ex) GET, POST etc..
	Method string
	// Pattern is the pattern of the URI.
	Pattern string
	// HandlerFunc is the handler function of this route.
	HandlerFunc gin.HandlerFunc
}

// Routes is the list of the generated Route.
type Routes []Route

// NewRouter returns a new router.
func NewRouter() *gin.Engine {
	router := gin.Default()
	AddService(router)
	return router
}

func AddService(engine *gin.Engine) *gin.RouterGroup {
	group := engine.Group("/api")

	for _, route := range routes {
		switch route.Method {
		case http.MethodGet:
			group.GET(route.Pattern, route.HandlerFunc)
		case http.MethodPost:
			group.POST(route.Pattern, route.HandlerFunc)
		case http.MethodPut:
			group.PUT(route.Pattern, route.HandlerFunc)
		case http.MethodDelete:
			group.DELETE(route.Pattern, route.HandlerFunc)
		case http.MethodPatch:
			group.PATCH(route.Pattern, route.HandlerFunc)
		}
	}

	return group
}

var routes = Routes{
	{
		"GetExample",
		http.MethodGet,
		"/sample",
		GetSampleJSON,
	},

	{
		"GetSubscribers",
		http.MethodGet,
		"/subscriber",
		GetSubscribers,
	},

	{
		"GetSubscriberByID",
		http.MethodGet,
		"/subscriber/:ueId/:servingPlmnId",
		GetSubscriberByID,
	},

	{
		"PostSubscriberByID",
		http.MethodPost,
		"/subscriber/:ueId/:servingPlmnId",
		PostSubscriberByID,
	},

	{
		"PutSubscriberByID",
		http.MethodPut,
		"/subscriber/:ueId/:servingPlmnId",
		PutSubscriberByID,
	},

	{
		"DeleteSubscriberByID",
		http.MethodDelete,
		"/subscriber/:ueId/:servingPlmnId",
		DeleteSubscriberByID,
	},

	{
		"PatchSubscriberByID",
		http.MethodPatch,
		"/subscriber/:ueId/:servingPlmnId",
		PatchSubscriberByID,
	},

	{
		"Registered UE Context",
		http.MethodGet,
		"/registered-ue-context",
		GetRegisteredUEContext,
	},

	{
		"Individual Registered UE Context",
		http.MethodGet,
		"/registered-ue-context/:supi",
		GetRegisteredUEContext,
	},

	{
		"UE PDU Session Info",
		http.MethodGet,
		"/ue-pdu-session-info/:smContextRef",
		GetUEPDUSessionInfo,
	},
}
