
This is a RESTful API written in Go lang to get realtime Stock data by using Stock Symbol and Stock Exchange list.
Used APIs from www.worldtradingdata.com to get stock data.

Downloading the Go Dependencies

To download the Go dependencies for this project, execute the following from the command line:
go get github.com/gorilla/mux
The mux package is a more feature rich option to the already existing net/http package for Go. It will be used for managing endpoints.

Install and Run
go get github.com/sirishamca/stockrestapi
After this command if you see in command line warning : github.com\sirishamca\stockrestapi\main.go:6:2: build constraints exclude all Go files
Please ignore.

Go to folder: cd $GOPATH/src/github.com/sirishamca/stockrestapi

Build code:
go build -tags dev

Run code:
go run -tags dev main.go


In case of production deployment
Please update config file under config folder with proper APIToken and log file path. 

Build code:
go build -tags prod

Run code:
go run -tags prod main.go

without updating APIToken if we run will get error in form API:

Invalid API Key.

API Endpoint:
http://localhost:10000/stock/{symbol”}

•	symbol: endpoint variable its mandatory.

•	stock_exchange– This is query parameter and its optional. Filter by a comma separated list of stock exchanges.

Example: 
http://localhost:10000/stock/AAPL?stock_exchange=NASDAQ
http://localhost:10000/stock/AAP
