package getStockData

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
)

type StockSearch struct {
	TotalReturned int         `json:"total_returned"`
	TotalResults  int         `json:"total_results"`
	TotalPages    int         `json:"total_pages"`
	Limit         int         `json:"limit"`
	Page          int         `json:"page"`
	Data          []StockData `json:"data"`
}

type StockRequestedData struct {
	SymbolsRequested int         `json:"symbols_requested"`
	SymbolsReturned  int         `json:"symbols_returned"`
	Data             []StockData `json:"data"`
}

type StockData struct {
	*StockFinalResponseData
	Currency           string `json:"currency"`
	PriceOpen          string `json:"price_open"`
	DayHigh            string `json:"day_high"`
	DayLow             string `json:"day_low"`
	Five2WeekHigh      string `json:"52_week_high"`
	Five2WeekLow       string `json:"52_week_low"`
	DayChange          string `json:"day_change"`
	ChangePct          string `json:"change_pct"`
	VolumeAvg          string `json:"volume_avg"`
	Shares             string `json:"shares"`
	StockExchangeLong  string `json:"stock_exchange_long"`
	StockExchangeShort string `json:"stock_exchange_short"`
}

type StockFinalResponseData struct {
	Symbol         string `json:"symbol"`
	Name           string `json:"name"`
	Price          string `json:"price"`
	CloseYesterday string `json:"close_yesterday"`
	MarketCap      string `json:"market_cap"`
	Volume         string `json:"volume"`
	Timezone       string `json:"timezone"`
	TimezoneName   string `json:"timezone_name"`
	GmtOffset      string `json:"gmt_offset"`
	LastTradeTime  string `json:"last_trade_time"`
}

var Symbol, StockExchange, BaseUrl, Api_Token string

func GetStockData(symbol string, stockExchange string) ([]byte, error) {

	Api_Token = "GeExVZllXEKNUCMc2MobfeftCj2mX8NVXa3AiH78IMzq69QbMtqU767l7fXX"
	Symbol = symbol
	StockExchange = stockExchange
	BaseUrl = "https://www.worldtradingdata.com/api/v1/"
	//Calling this URL to validate Stock
	url := BaseUrl + "stock_search?search_term=" + symbol + "&search_by=symbol&stock_exchange=" + stockExchange + "&api_token=" + Api_Token
	// Build the request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal("NewRequest: ", err)
		return nil, err
	}
	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Do: ", err)
		return nil, err
	}

	defer resp.Body.Close()

	var record StockSearch

	// Use json.Decode for reading response data
	if err := json.NewDecoder(resp.Body).Decode(&record); err != nil {
		log.Println(err)
		return nil, err
	}
	//Validate Symbol, if Data present then onlly call other API to fetch data.
	if len(record.Data) > 0 && record.Data[0].Symbol == symbol {
		return fetchStockData(client)
	} else {
		return nil, errors.New("Message:Error! The requested data could not be found.")
	}
}

func fetchStockData(httpcl *http.Client) ([]byte, error) {
	url := BaseUrl + "stock?symbol=" + Symbol + "&api_token=" + Api_Token
	// Build the request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal("NewRequest: ", err)
		return nil, err
	}
	
	resp, err := httpcl.Do(req)
	if err != nil {
		log.Fatal("Do: ", err)
		return nil, err
	}

	defer resp.Body.Close()

	var record StockRequestedData

	if err := json.NewDecoder(resp.Body).Decode(&record); err != nil {
		log.Println(err)
		return nil, err
	}
	datas := make(map[string]*StockFinalResponseData)
	StockExchangeStr := strings.Split(strings.ToUpper(StockExchange), ",")
	for _, value := range record.Data {
		if contains(StockExchangeStr, strings.ToUpper(value.StockExchangeShort)) {
			datas[value.StockExchangeShort] = value.StockFinalResponseData
		}
	}
	data, _ := json.Marshal(datas)
	return data, nil

}
func contains(str []string, str1 string) bool {
	for _, s := range str {
		if s == str1 {
			return true
		}
	}
	return false
}
