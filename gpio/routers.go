package swagger

import (
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/peramic/utils"
)

//NewRouter Returns all routes
func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		var handler http.Handler
		handler = route.HandlerFunc
		handler = Logger(handler, route.Name)

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}

	return router
}

//AddRoutes add new  routes
func AddRoutes(myroute []utils.Route) {
	routes = append(routes, myroute...)
}

var routes = []utils.Route{
	utils.Route{
		Name:        "GetAPI",
		Method:      strings.ToUpper("Get"),
		Pattern:     "/rest/gpio/apidocs",
		HandlerFunc: getAPI,
	},
	utils.Route{
		Name:        "GetDevices",
		Method:      strings.ToUpper("Get"),
		Pattern:     "/rest/gpio/devices",
		HandlerFunc: getDevices,
	},
	utils.Route{
		Name:        "GetDevice",
		Method:      strings.ToUpper("Get"),
		Pattern:     "/rest/gpio/devices/{device}",
		HandlerFunc: getDevice,
	},

	utils.Route{
		Name:        "GetLabel",
		Method:      strings.ToUpper("Get"),
		Pattern:     "/rest/gpio/devices/{device}/label",
		HandlerFunc: getLabel,
	},
	utils.Route{
		Name:        "DeleteLabel",
		Method:      strings.ToUpper("Delete"),
		Pattern:     "/rest/gpio/devices/{device}/label",
		HandlerFunc: deleteLabel,
	},
	utils.Route{
		Name:        "GetProperties",
		Method:      strings.ToUpper("Get"),
		Pattern:     "/rest/gpio/devices/{device}/properties",
		HandlerFunc: getProperties,
	},
	utils.Route{
		Name:        "GetProperty",
		Method:      strings.ToUpper("Get"),
		Pattern:     "/rest/gpio/devices/{device}/properties/{name}",
		HandlerFunc: getProperty,
	},
	utils.Route{
		Name:        "SetProperty",
		Method:      strings.ToUpper("Put"),
		Pattern:     "/rest/gpio/devices/{device}/properties/{name}",
		HandlerFunc: setProperty,
	},
	utils.Route{
		Name:        "SetLabel",
		Method:      strings.ToUpper("Put"),
		Pattern:     "/rest/gpio/devices/{device}/label",
		HandlerFunc: setLabel,
	},
	utils.Route{
		Name:        "GetFields",
		Method:      strings.ToUpper("Get"),
		Pattern:     "/rest/gpio/devices/{device}/fields",
		HandlerFunc: getFields,
	},
	utils.Route{
		Name:        "GetField",
		Method:      strings.ToUpper("Get"),
		Pattern:     "/rest/gpio/devices/{device}/fields/{field}",
		HandlerFunc: getField,
	},
	utils.Route{
		Name:        "SetField",
		Method:      strings.ToUpper("Put"),
		Pattern:     "/rest/gpio/devices/{device}/fields/{field}",
		HandlerFunc: setField,
	},
	utils.Route{
		Name:        "GetFieldProperties",
		Method:      strings.ToUpper("Get"),
		Pattern:     "/rest/gpio/devices/{device}/fields/{field}/properties",
		HandlerFunc: getFieldProperties,
	},
	utils.Route{
		Name:        "SetFieldProperties",
		Method:      strings.ToUpper("PUT"),
		Pattern:     "/rest/gpio/devices/{device}/fields/{field}/properties",
		HandlerFunc: setFieldProperties,
	},
	utils.Route{
		Name:        "SetFieldProperty",
		Method:      strings.ToUpper("Put"),
		Pattern:     "/rest/gpio/devices/{device}/fields/{field}/properties/{name}",
		HandlerFunc: setFieldProperty,
	},
	utils.Route{
		Name:        "GetFieldProperty",
		Method:      strings.ToUpper("Get"),
		Pattern:     "/rest/gpio/devices/{device}/fields/{field}/properties/{name}",
		HandlerFunc: getFieldProperty,
	},
	utils.Route{
		Name:        "SetFieldValue",
		Method:      strings.ToUpper("Put"),
		Pattern:     "/rest/gpio/devices/{device}/fields/{field}/value",
		HandlerFunc: setFieldValue,
	},
	utils.Route{
		Name:        "GetFieldValue",
		Method:      strings.ToUpper("Get"),
		Pattern:     "/rest/gpio/devices/{device}/fields/{field}/value",
		HandlerFunc: getFieldValue,
	},
	utils.Route{
		Name:        "GetFieldLabel",
		Method:      strings.ToUpper("Get"),
		Pattern:     "/rest/gpio/devices/{device}/fields/{field}/label",
		HandlerFunc: getFieldLabel,
	},
	utils.Route{
		Name:        "SetFieldLabel",
		Method:      strings.ToUpper("Put"),
		Pattern:     "/rest/gpio/devices/{device}/fields/{field}/label",
		HandlerFunc: setFieldLabel,
	},
	utils.Route{
		Name:        "DeleteFieldLabel",
		Method:      strings.ToUpper("Delete"),
		Pattern:     "/rest/gpio/devices/{device}/fields/{field}/label",
		HandlerFunc: deleteFieldLabel,
	},
	utils.Route{
		Name:        "cors",
		Method:      "GET",
		Pattern:     "/ws",
		HandlerFunc: handleConnections,
	},
	//Reports
	utils.Route{
		Name:        "GetReports",
		Method:      strings.ToUpper("Get"),
		Pattern:     "/rest/gpio/reports",
		HandlerFunc: getReports,
	},
	utils.Route{
		Name:        "AddReport",
		Method:      strings.ToUpper("Post"),
		Pattern:     "/rest/gpio/reports",
		HandlerFunc: addReport,
	},
	utils.Route{
		Name:        "GetReport",
		Method:      strings.ToUpper("Get"),
		Pattern:     "/rest/gpio/reports/{id}",
		HandlerFunc: getReport,
	},
	utils.Route{
		Name:        "SetReport",
		Method:      strings.ToUpper("Put"),
		Pattern:     "/rest/gpio/reports/{id}",
		HandlerFunc: setReport,
	},
	utils.Route{
		Name:        "DeleteReport",
		Method:      strings.ToUpper("Delete"),
		Pattern:     "/rest/gpio/reports/{id}",
		HandlerFunc: deleteReport,
	},
	utils.Route{
		Name:        "GetSubscriptors",
		Method:      strings.ToUpper("Get"),
		Pattern:     "/rest/gpio/reports/{id}/subscriptions",
		HandlerFunc: getSubscriptors,
	},
	utils.Route{
		Name:        "AddSubscriptor",
		Method:      strings.ToUpper("Post"),
		Pattern:     "/rest/gpio/reports/{id}/subscriptions",
		HandlerFunc: addSubscriptor,
	},
	utils.Route{
		Name:        "GetSubscriptor",
		Method:      strings.ToUpper("Get"),
		Pattern:     "/rest/gpio/reports/{id}/subscriptions/{subscriptorId}",
		HandlerFunc: getSubscriptor,
	},
	utils.Route{
		Name:        "SetSubscriptor",
		Method:      strings.ToUpper("Put"),
		Pattern:     "/rest/gpio/reports/{id}/subscriptions/{subscriptorId}",
		HandlerFunc: setSubscriptor,
	},
	utils.Route{
		Name:        "DeleteSubscriptor",
		Method:      strings.ToUpper("Delete"),
		Pattern:     "/rest/gpio/reports/{id}/subscriptions/{subscriptorId}",
		HandlerFunc: deleteSubscriptor,
	},
	utils.Route{
		Name:        "SetPinValue",
		Method:      strings.ToUpper("Get"),
		Pattern:     "/rest/gpio/setpinvalue/{id}/value/{value}",
		HandlerFunc: mockPinValue,
	},
}
