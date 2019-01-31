package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/sirishamca/stockrestapi/config"
	"github.com/sirishamca/stockrestapi/getStockData"
	"github.com/sirishamca/stockrestapi/srslogger"
	"go/build"
	"log"
	"net/http"
)

var l srslogger.Logger

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
	data, err := getStockData.GetStockData(key, queryParam, l)
	if err != nil {
		fmt.Fprintf(w, err.Error())
	} else {
		w.Write(data)
	}
}

func main() {
	var logpath = build.Default.GOPATH + config.LOGFILEPATH
	l.Init(logpath)
	handleRequests()
}
