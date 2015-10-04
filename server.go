package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net"
	"net/http"
	"net/rpc"
	"strconv"
	"strings"
)

type query struct {
	Query stockQuery
}

type stockQuery struct {
	Count   int
	Created string
	Lang    string
	Results result
}

type result struct {
	Quote []quote
}

type quote struct {
	LastTradePriceOnly string
	Symbol             string
	realPrice          float64
}

type stocks struct {
	eachStock []stock
}

type stock struct {
	name         string
	buyPrice     float64
	currentPrice float64
	shares       int
}

type Person struct {
	ID               int
	stocks           []stock
	capital          float64
	totolInvested    float64
	uninvestedAmount float64
}

type Buy struct {
	TradeId        string
	Stocks         string
	UnvestedAmount string
}

type Check struct {
	Stocks             string
	CurrentMarketValue float64
	UnvestedAmount     float64
}

type Echo int

var Client1 Person
var quotes []quote
var share int
var price float64

func (t *Echo) Buy(args string, reply *Buy) error {
	strArray := strings.Split(args, "#")
	str := strArray[0]
	strCap := strArray[1]

	strUrl := formQuery(str)       // Form a string that will be used in url(%22YHOO%22)
	quotes := GetStockInfo(strUrl) // Get the stock information in []quote format

	names := make([]stock, len(quotes)) // Declare a slice, each unit contains a stock type element
	for i := 0; i < len(quotes); i++ {  // Assign the set of stocks with quotes' information
		names[i].name = quotes[i].Symbol
		names[i].buyPrice = quotes[i].realPrice
		names[i].shares = 0
		names[i].currentPrice = quotes[i].realPrice
	}

	percentages := splitPercentage(str)              // Pull the percentage of every stock
	capital, error := strconv.ParseFloat(strCap, 64) // Convert string into float
	if error != nil {
		fmt.Println(error)
	}

	for i := 0; i < len(quotes); i++ { // Calculate the shares of each stock
		names[i].shares = int((capital * percentages[i]) / names[i].buyPrice)
	}

	var totalInvested float64
	for i := 0; i < len(quotes); i++ {
		totalInvested += float64(names[i].shares) * names[i].buyPrice
	}

	tradeID := rand.Intn(100)
	Client1.ID = tradeID
	Client1.stocks = names
	Client1.capital = capital
	Client1.totolInvested = totalInvested
	Client1.uninvestedAmount = capital - totalInvested

	var resultBuy string
	for i := 0; i < len(Client1.stocks); i++ {
		share = Client1.stocks[i].shares
		price = Client1.stocks[i].buyPrice * float64(Client1.stocks[i].shares)
		resultBuy += Client1.stocks[i].name + ":" + strconv.Itoa(share) + ":$" + strconv.FormatFloat(price, 'f', 2, 64) + ", "
	}

	strTradeID := strconv.Itoa(tradeID)
	strUnInvested := strconv.FormatFloat(capital-totalInvested, 'f', 2, 64)

	reply.Stocks = resultBuy
	reply.TradeId = strTradeID
	reply.UnvestedAmount = strUnInvested

	return nil
}

func (t *Echo) CheckPortfolio(args string, reply *Check) error {
	id := strconv.Itoa(Client1.ID)
	if args != id {
		fmt.Println("no user")
		return nil
	} else {
		quotesCheck := GetStockInfo(formQuery_Check(Client1))
		fmt.Println("length of quoteCHeck", len(quotesCheck))
		count := len(quotesCheck)
		for i := 0; i < count; i++ {
			Client1.stocks[i].currentPrice = quotes[i].realPrice
		}
		fmt.Println(Client1)
		var currentTotal = make([]float64, count)
		var buyTotal = make([]float64, count)
		var signs = make([]string, count)
		for i := 0; i < count; i++ {
			currentTotal[i] = float64(Client1.stocks[i].shares) * Client1.stocks[i].currentPrice
			buyTotal[i] = float64(Client1.stocks[i].shares) * Client1.stocks[i].buyPrice
		}
		// Set up the sign
		for i := 0; i < count; i++ {
			if currentTotal[i] >= buyTotal[i] {
				signs[i] = "+"
			} else if currentTotal[i] < buyTotal[i] {
				signs[i] = "-"
			}
		}
		var resultCheck string
		for i := 0; i < count; i++ {
			share = Client1.stocks[i].shares
			price = Client1.stocks[i].buyPrice * float64(Client1.stocks[i].shares)
			resultCheck += Client1.stocks[i].name + ":" + strconv.Itoa(share) + ":$" + strconv.FormatFloat(price, 'f', 2, 64) + ", "
		}
		var currentMV float64
		for i := 0; i < count; i++ {
			currentMV += currentTotal[i]
		}
		reply.CurrentMarketValue = currentMV
		reply.Stocks = resultCheck
		reply.UnvestedAmount = Client1.uninvestedAmount
	}
	return nil
}

func main() {
	rpc.Register(new(Echo))
	rpc.HandleHTTP()
	l, e := net.Listen("tcp", ":1234")
	if e != nil {
		fmt.Println("listen error:", e)
	}
	http.Serve(l, nil)
}

func formQuery_Check(client Person) string {
	var rslt string
	if len(client.stocks) == 1 {
		rslt = "%22" + client.stocks[0].name + "%22"
	} else if len(client.stocks) > 1 {
		for i := 0; i < len(client.stocks)-1; i++ {
			rslt += "%22" + client.stocks[i].name + "%22%2C"
		}
		rslt += "%22" + client.stocks[len(client.stocks)-1].name + "%22"
	}
	return rslt
}

func formQuery(str string) string {
	strArray := strings.Split(str, ",")
	temp := concateStockQuery(splitSymbol01(strArray))
	return temp
}

func splitSymbol01(str []string) []string {
	var temp string
	for i := 0; i < len(str)-1; i++ {
		temp += splitSymbol02(str[i]) + ","
	}
	temp += splitSymbol02(str[len(str)-1])
	result := strings.Split(temp, ",")
	//fmt.Println(result)
	return result
}

func splitSymbol02(str string) string {
	result := strings.Split(str, ":")
	return result[0]
}

func concateStockQuery(names []string) string {
	length := len(names)
	str := ""

	if length == 1 {
		str += "%22" + names[0] + "%22"
	}
	if length > 1 {
		for i := 0; i < length-1; i++ {
			str += "%22" + names[i] + "%22%2C"
		}
		str += "%22" + names[length-1] + "%22"
	}
	return str
}

func GetStockInfo(str string) []quote {

	firstUrl := "https://query.yahooapis.com/v1/public/yql?q=select%20Symbol%2C%20LastTradePriceOnly%20from%20yahoo.finance.quote%20where%20symbol%20in%20("
	lastUrl := ")&format=json&env=store%3A%2F%2Fdatatables.org%2Falltableswithkeys&callback="
	finalUrl := firstUrl + str + lastUrl

	resp, err := http.Get(finalUrl)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	content, err1 := ioutil.ReadAll(resp.Body)
	if err1 != nil {
		fmt.Println(err1)
	}
	//fmt.Println("content is: ", string(content))
	var json_rslt query
	//var json_rslt interface{}

	err2 := json.Unmarshal(content, &json_rslt)
	if err2 != nil {
		panic(err2)
	}
	//m := f.(map[string]interface{})

	var count = len(json_rslt.Query.Results.Quote)
	quotes := json_rslt.Query.Results.Quote
	if count != 0 {
		for i := 0; i < count; i++ {
			fmt.Println("Symbol is: ", json_rslt.Query.Results.Quote[i].Symbol)
			fmt.Println("Price is: ", json_rslt.Query.Results.Quote[i].LastTradePriceOnly)
		}
	}

	for i := 0; i < len(quotes); i++ {
		temp, error := strconv.ParseFloat(quotes[i].LastTradePriceOnly, 64)
		if error != nil {
			fmt.Println(error)
			break
		}
		quotes[i].realPrice = temp
	}

	return quotes
}

func splitPercentage(str string) []float64 {
	temp1 := strings.Split(str, ",")
	percentages := make([]float64, len(temp1))
	for i := 0; i < len(temp1); i++ {
		percentages[i] = splitPercentage2(splitPercentage1(temp1[i]))
	}
	return percentages
}

func splitPercentage1(str string) string {
	s := strings.Split(str, ":")
	return s[1]
}

func splitPercentage2(strNum string) float64 {
	b := strings.Split(strNum, "%")
	result, error := strconv.ParseFloat(b[0], 64)
	if error != nil {
		fmt.Println(error)
		result = -1.0
	}
	return result / 100
}
