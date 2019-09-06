package main

//Binding routes and procedures/functions in handler.go
var routes = Routes{
	Route{
		"CreateInvoice",
		"POST",
		"/invoice/create/",
		CreateInvoice,
	},
	Route{
		"DeleteInvoice",
		"DELETE",
		"/invoice/delete/",
		DeleteInvoice,
	},
	Route{
		"PayInvoice",
		"POST",
		"/invoice/pay/",
		PayInvoice,
	},
	Route{
		"ListInvoices",
		"GET",
		"/invoice/list/",
		ListInvoices,
	},
}
