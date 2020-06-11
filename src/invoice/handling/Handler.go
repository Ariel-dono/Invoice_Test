package handling

import (
	"encoding/json"
	"net/http"
	"io/ioutil"
	"encoding/xml"
	"strings"
	"time"
	"fmt"
	storing "../storing"
	model "../model"
)

//It's a bridge to Storage.go like an "event handler"
func CreateInvoice(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var instance model.Invoice
	var response *model.Response = new(model.Response)
	err := decoder.Decode(&instance)
	if err != nil {
		response.Code = -5
		response.State.Message = "decoding-body"
	}else{
		defer r.Body.Close()
		storing.CreateWithUuid(setupInvoice(instance, response))
	}
	responseMessage(w,response)
}

//Fullfilling the invoice and making initial calculations
func setupInvoice(invoice model.Invoice, response *model.Response)(model.Invoice,*model.Response){
	exchangeUSDBCCR := getUSDExchangeRateBCCR(response)
	for _, line := range invoice.Lines {
		if strings.Compare(line.Currency, "USD") == 0{
			line.Price = line.Price*exchangeUSDBCCR
		}else if strings.Compare(line.Currency, "CRC") != 0{
			response.Code = -4
			response.State.Message = "unsupported-currency"
		}
		line_total := line.Price * line.Quantity
		invoice.Subtotal += line_total
		invoice.TaxTotal += line_total*line.TaxRate
		invoice.DiscountTotal += line_total*line.DiscountRate
	}
	invoice.Total = invoice.Subtotal - invoice.DiscountTotal + invoice.TaxTotal
	invoice.Balance = invoice.Total
	invoice.Payments = []model.Payment{}
	return invoice, response
}

func DeleteInvoice(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	var invoice model.Request
	var response *model.Response = new(model.Response)
	err := decoder.Decode(&invoice)
	if err != nil {
		response.Code = -5
		response.State.Message = "decoding-body"
	}else{
		defer r.Body.Close()
		storing.Remove(invoice.Id, response)
	}
	responseMessage(w,response)
}

func PayInvoice(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	var payment model.Payment
	var response *model.Response = new(model.Response)
	err := decoder.Decode(&payment)
	if err != nil {
		response.Code = -5
		response.State.Message = "decoding-body"
	}else{
		defer r.Body.Close()
		storing.MutateWithUuid(payment, response)
	}
	responseMessage(w,response)
}

func ListInvoices(w http.ResponseWriter, r *http.Request) {
	var response *model.Response = new(model.Response)
	var invoices *model.InvoicesResponse = new(model.InvoicesResponse)
	defer r.Body.Close()
	storing.RetrieveInvoices(response, invoices)
	responseList(w, response, invoices)
}

//Setting up the response configuration for simple messages
func responseMessage (w http.ResponseWriter,  response *model.Response) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if response.Code < 0{
		w.WriteHeader(http.StatusInternalServerError)
	}else{
		w.WriteHeader(http.StatusOK)
	}
	err := json.NewEncoder(w).Encode(&response)
	if err != nil {
		panic(err)
	}
}

//Setting up the response configuration for invoice list format
func responseList (w http.ResponseWriter,  response *model.Response, invoices *model.InvoicesResponse) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	var err error
	if response.Code < 0{
		w.WriteHeader(http.StatusInternalServerError)
		err = json.NewEncoder(w).Encode(&response)
	}else{
		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode(&invoices)
	}
	if err != nil {
		panic(err)
	}
}

//Provides the exchange rates from BCCR
func getUSDExchangeRateBCCR(response *model.Response)float32{
	var exchangeRateBCCR model.BCCRGaugeStructure
	currentDateTime := time.Now()
	//Format DD/MM/YYYY
	today := currentDateTime.Format("02/01/2006")
	resp, err := http.Get(fmt.Sprintf("http://indicadoreseconomicos.bccr.fi.cr/indicadoreseconomicos/WebServices/wsIndicadoresEconomicos.asmx/ObtenerIndicadoresEconomicosXML?tcIndicador=318&tcFechaInicio=%s&tcFechaFinal=%s&tcNombre=arielherrera&tnSubNiveles=N",today,today))
	if(err != nil){
		response.Code = -5
		response.State.Message = "connection-error"
	}else{
		body, err := ioutil.ReadAll(resp.Body)
		if(err != nil){
			response.Code = -5
			response.State.Message = "error-loading-exchange-rate"
		}else{
			exchange := string(body)
			exchange = strings.Replace(exchange, "&lt;", "<", -1)
			exchange = strings.Replace(exchange, "&gt;", ">", -1)
			err := xml.Unmarshal([]byte(exchange), &exchangeRateBCCR)
			if(err != nil){
				response.Code = -5
				response.State.Message = "error-loading-exchange-rate"
			}
			defer resp.Body.Close()
		}
	}
	return exchangeRateBCCR.GaugeData.Exchange.Valor
}
