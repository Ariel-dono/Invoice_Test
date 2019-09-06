package main

import (
	"gopkg.in/couchbase/gocb.v1"
	"github.com/google/uuid"
)

// Storage layer

const bucketName string = "Invoices"
//getting couchbase bucket connetion
func initialize() gocb.Bucket{
	cluster, _ := gocb.Connect("couchbase://localhost")
	cluster.Authenticate(gocb.PasswordAuthenticator{
		Username: "root",
		Password: "$2quantumDot$2",
	})

	bucket, _ := cluster.OpenBucket(bucketName, "")
	bucket.Manager("", "").CreatePrimaryIndex("", true, false) //loading bucket's primary index
	return *bucket
}

//generating uuid by invoice
func createWithUuid (instance Invoice, response *Response){
  if response.Code > -4{
		uuid_invoice, errUuid := uuid.NewRandom()
		if errUuid == nil{
			uuid_invoice_str := uuid_invoice.String()
			instance.Id = uuid_invoice_str
			create(uuid_invoice_str, instance, response)
		}else{
			response.Code = -3
			response.State.Message = "error-generating-uuid"
		}
	}
}

//Creating full invoice
func create (uuid string, instance Invoice, response *Response){
	bucket := initialize()
	_, err := bucket.Insert(uuid, instance, 0)
	if err != nil {
		response.Code = -2
		response.State.Message = "error-mutating-invoice"
	}else {
		response.Code = 1
		response.State.Message = "invoice-was-created"
	}
	bucket.Close()
}

//Removing full invoice
func remove(key string, response *Response){
	bucket := initialize()
	var cas gocb.Cas
	_,err := bucket.Remove(key,cas)
	if err != nil {
		response.Code = -1
		response.State.Message = "error-deleting-invoice"
	}else{
		response.Code = 1
		response.State.Message = "invoice-was-deleted"
	}
	bucket.Close()
}

//generating uuid by payment
func mutateWithUuid (payment Payment, response *Response){
	uuid_invoice, errUuid := uuid.NewRandom()
	if errUuid == nil {
		current_key := payment.Id
		payment.Id = uuid_invoice.String()
		mutate(current_key, payment, response)
	}else{
		response.Code = -3
		response.State.Message = "error-generating-uuid"
	}
}

//Updating only keys balance and payments from full invoice
func mutate(key string, payment Payment, response *Response){
	bucket := initialize()
	var invoice Invoice
	_,errRetrieving := bucket.Get(key, &invoice)
	if errRetrieving != nil {
		response.Code = -1
		response.State.Message = "error-invoice-unavailable"
	}else{
		if invoice.Balance - payment.Total >= 0{
			_, errPayments := bucket.MutateIn(key, 0, 0).ArrayAppend("payments", payment, false).Execute()
			_, errBalance := bucket.MutateIn(key, 0, 0).Upsert("balance", invoice.Balance - payment.Total, false).Execute()
			if errPayments != nil {
				response.Code = -2
				response.State.Message = "error-mutating-invoice"
			}else if errBalance != nil {
				response.Code = -3
				response.State.Message = "error-mutating-balance"
			}else {
				response.Code = 1
				response.State.Message = "payment-was-created"
			}
		}else{
			response.Code = -4
			response.State.Message = "error-negative-balance"
		}
	}
	bucket.Close()
}

//It uses Couchbase N1QL Engine
func retrieveInvoices(response *Response, invoices *InvoicesResponse){
	bucket := initialize()
	invoicesQuery := gocb.NewN1qlQuery("SELECT * FROM Invoices")
	rows, err := bucket.ExecuteN1qlQuery(invoicesQuery, nil)
	var row map[string]Invoice
	for rows.Next(&row) {
		invoices.Response = append(invoices.Response, row["Invoices"])
	}
	if err != nil {
		response.Code = -1;
		response.State.Message = "error-retrieving-invoice";
	}else{
		response.Code = 1;
		response.State.Message = "retrieving-invoices";
	}
	bucket.Close();
}
