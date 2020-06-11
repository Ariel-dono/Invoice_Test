package model

import (
	"net/http"
	"encoding/xml"
)

//Routing model
type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

//Response messaging
type Response struct {
	State      	State	`json:"state"`
	Code 				int   `json:"code"`
}

//Response message content
type State struct {
	Message 	string	`json:"message"`
}

//Request by key
type Request  struct {
	Id 		string	`json:"id"`
}

type InvoicesResponse struct {
	Response []Invoice `json:"response"`
}

type Invoice struct {
	Id 							string				`json:"id",omitempty`
	Lines						[]LineDetail	`json:"lines"`
	Client 					Client 			 	`json:"client"`
	TaxTotal				float32				`json:"tax_total",omitempty`
	DiscountTotal		float32				`json:"discount_total",omitempty`
	Subtotal				float32				`json:"subtotal",omitempty`
	Total 					float32				`json:"total",omitempty`
	Payments				[]Payment			`json:"payments",omitempty`
	Balance					float32				`json:"balance",omitempty`
}

type Payment struct{
	Id		string 	`json:"id"`
	Total	float32	`json:"total"`
}

type LineDetail struct {
	Product 			string	`json:"product"`
	Quantity 			float32	`json:"quantity"`
	Price 				float32	`json:"price"`
	TaxRate 			float32	`json:"tax_rate"`
	DiscountRate 	float32	`json:"discount_rate"`
	Currency 			string	`json:"currency"`
}

type Client struct {
	Name		string	`json:"name"`
	Id			string	`json:"id"`
}

type BCCRGaugeStructure struct{
	XMLName     xml.Name
	GaugeData 	BCCRGaugeData `xml:"Datos_de_INGC011_CAT_INDICADORECONOMIC"`
}

type BCCRGaugeData struct{
	XMLName     xml.Name
	Exchange 	BCCRExchange	`xml:"INGC011_CAT_INDICADORECONOMIC"`
}

type BCCRExchange struct{
	Valor 		float32 `xml:"NUM_VALOR"`
	GaugeCode int 		`xml:"COD_INDICADORINTERNO"`
	Date 			string 	`xml:"DES_FECHA"`
}
