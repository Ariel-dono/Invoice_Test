package routing

import (
	"../handling"
	"../model"
)

//Binding routes and procedures/functions in handler.go
var routes = model.Routes{
	model.Route{
		"CreateInvoice",
		"POST",
		"/invoice/create/",
		handling.CreateInvoice,
	},
	model.Route{
		"DeleteInvoice",
		"DELETE",
		"/invoice/delete/",
		handling.DeleteInvoice,
	},
	model.Route{
		"PayInvoice",
		"POST",
		"/invoice/pay/",
		handling.PayInvoice,
	},
	model.Route{
		"ListInvoices",
		"GET",
		"/invoice/list/",
		handling.ListInvoices,
	},
}
