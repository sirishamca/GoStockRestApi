package getStockData

import (
	"encoding/json"
	"errors"
	"github.com/sirishamca/stockrestapi/config"
	"github.com/sirishamca/stockrestapi/model"
	"github.com/sirishamca/stockrestapi/srslogger"
	"net/http"
	"strings"
)

var Symbol, StockExchange, BaseUrl string

func GetStockData(symbol string, stockExchange string, logger srslogger.Logger) ([]byte, error) {

	defaultFlag := false
	if len(stockExchange) == 0 {
		//setting to default Stock Exchange
		StockExchange = "AMEX"
		defaultFlag = true
		logger.Info("Using default Stock Exchange AMEX")
	} else {
		StockExchange = strings.ToUpper(stockExchange)
	}
	Symbol = strings.ToUpper(symbol)

	BaseUrl = "https://www.worldtradingdata.com/api/v1/"
	//Calling this URL to validate Stock
	url := BaseUrl + "stock_search?search_term=" + Symbol + "&search_by=symbol&stock_exchange=" + StockExchange + "&api_token=" + config.APITOKEN
	logger.Info("Seraching StockData URL=" + url)
	// Build the request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		logger.Fatal("NewRequest: ", err)
		return nil, err
	}
	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		logger.Fatal("Do: ", err)
		return nil, err
	}

	defer resp.Body.Close()

	var record model.StockSearch

	// Use json.Decode for reading response data
	if err := json.NewDecoder(resp.Body).Decode(&record); err != nil {
		logger.Fatal(err)
		return nil, err
	}

	//Validate Symbol, if Data present then onlly call other API to fetch data.
	if len(record.Data) > 0 && record.Data[0].Symbol == Symbol {
		return fetchStockData(client, logger)
	} else if page != 0{
		if defaultFlag {
			return nil, errors.New("Message:No data found for the stock symbol:\"" + Symbol + "\" with default stock exchange \"" + StockExchange + "\".")
		}
		return nil, errors.New("Message:No data found for the stock symbol:\"" + Symbol + "\" with stock exchange(s) \"" + StockExchange + "\".")
	} {
		return nil, errors.New(record.MESSAGE)
	}
}

func fetchStockData(httpcl *http.Client, logger srslogger.Logger) ([]byte, error) {
	url := BaseUrl + "stock?symbol=" + Symbol + "&api_token=" + config.APITOKEN
	// Build the request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		logger.Fatal("NewRequest: ", err)
		return nil, err
	}
	logger.Info("Fetch Realtime StockData url=" + url)
	resp, err := httpcl.Do(req)
	if err != nil {
		logger.Fatal("Do: ", err)
		return nil, err
	}

	defer resp.Body.Close()

	var record model.StockRequestedData

	if err := json.NewDecoder(resp.Body).Decode(&record); err != nil {
		logger.Error(err)
		return nil, err
	}

	datas := make(map[string]*model.StockFinalResponseData)
	//Splitting StockExchange since its accepts multiple values with comma seperated
	stockExchangeStr := strings.Split(StockExchange, ",")
	for _, value := range record.Data {
		if contains(stockExchangeStr, strings.ToUpper(value.StockExchangeShort)) {
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
