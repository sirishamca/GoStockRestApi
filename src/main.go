package main

import (
	"fmt"
	"getStockData"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func handleRequests() {

	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/stock/{symbol}", returnStockData)
	log.Fatal(http.ListenAndServe(":10000", myRouter))
}

func returnStockData(w http.ResponseWriter, r *http.Request) {
	queryValues := r.URL.Query()
	queryParam := queryValues.Get("stock_exchange")

	vars := mux.Vars(r)
	key := vars["symbol"]
	if len(queryParam) < 1 {
		//setting to default Stock Exchange
		queryParam = "AMEX"
	}
	data, err := getStockData.GetStockData(key, queryParam)
	if err != nil {
		fmt.Fprintf(w, err.Error())
	} else {
		w.Write(data)
	}
}

func main() {
	handleRequests()
}
