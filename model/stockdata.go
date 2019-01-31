package model

type StockSearch struct {
	MESSAGE string      `json:"message"`
	Data    []StockData `json:"data"`
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
