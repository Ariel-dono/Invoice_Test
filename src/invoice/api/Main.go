package main

import (
"log"
"net/http"
routing "../routing"
)

func main() {
	router := routing.NewRouter();
	log.Fatal(http.ListenAndServe(":9000", router));
}
