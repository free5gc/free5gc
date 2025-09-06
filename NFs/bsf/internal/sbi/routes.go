package sbi

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Route defines a single API route
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

// applyRoutes applies the given routes to the router group
func applyRoutes(group *gin.RouterGroup, routes Routes) {
	for _, route := range routes {
		switch route.Method {
		case http.MethodGet:
			group.GET(route.Pattern, route.HandlerFunc)
		case http.MethodPost:
			group.POST(route.Pattern, route.HandlerFunc)
		case http.MethodPut:
			group.PUT(route.Pattern, route.HandlerFunc)
		case http.MethodPatch:
			group.PATCH(route.Pattern, route.HandlerFunc)
		case http.MethodDelete:
			group.DELETE(route.Pattern, route.HandlerFunc)
		}
	}
}
