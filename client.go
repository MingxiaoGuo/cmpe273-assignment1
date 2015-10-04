package main

import (
	"fmt"
	"net/rpc"
)

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

func main() {
	client, err := rpc.DialHTTP("tcp", "127.0.0.1:1234")
	if err != nil {
		fmt.Println("dialing:", err)
	}
	// buy
	var mystock string
	var myBudget string
	fmt.Println("please input your stockName")
	fmt.Scanln(&mystock)
	fmt.Println("please input your budget")
	fmt.Scanln(&myBudget)
	buyClient := mystock + "#" + myBudget
	var buy Buy
	err = client.Call("Echo.Buy", buyClient, &buy)
	fmt.Println("tradeID:", buy.TradeId)
	fmt.Println("stocks:", buy.Stocks)
	fmt.Println("unvestedAmount:", buy.UnvestedAmount)
	// check
	var checkP Check
	var id string
	fmt.Println("input your id")
	fmt.Scanln(&id)
	err = client.Call("Echo.CheckProfolio", id, &checkP)
	fmt.Println(checkP)
}
