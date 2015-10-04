package main

import (
	"bufio"

	"fmt"

	"log"

	"net/rpc"

	"os"
)

type StockDetails struct {
	Symbol string

	Budget float64
}

type EstimatedAmount struct {
	MarketStocks string `json:"stocksymbol"`

	AvailaibleAmount float64 `json:"stockvaluefetch"`

	TradeId int `json:"tss"`
}

type UniqueId struct {
	TradeId int `json:"tss"`
}

type UpdEstimatedAmount struct {
	MarketStocks string `json:"stocksymbol"`

	AvailaibleAmount float64 `json:"stockvaluefetch"`
}

func main() {

	client, err := rpc.DialHTTP("tcp", "127.0.0.1:1528")

	var choice int

	fmt.Println("Virtual Stock Market Analyser")

	fmt.Println("-----------------------------")

	fmt.Println("Select an option in the following format: YHOO:100")

	fmt.Println("-----------------------------")

	fmt.Println("1. I want to buy stocks")

	fmt.Println("2. Show my Portfolio")

	fmt.Println("-----------------------------")

	fmt.Print("Enter 1 or 2::  ")

	fmt.Scanf("%d", &choice)

	switch choice {

	case 1:

		reader := bufio.NewReader(os.Stdin)

		fmt.Print("Enter Stock Symbol: ")

		Symbol, _ := reader.ReadString('\n')

		fmt.Println(Symbol)

		fmt.Print("Enter the amount you plan to invest: ")

		var Budget float64

		fmt.Scan(&Budget)

		userInput := StockDetails{Symbol, Budget}

		var estimate EstimatedAmount

		err = client.Call("SttockCall.StockValueFetch", userInput, &estimate)

		if err != nil {

			log.Fatal("arith error:", err)

		}

		fmt.Print("Stocks : ")

		fmt.Println(estimate.MarketStocks)

		fmt.Print("Amount: ")

		fmt.Println(estimate.AvailaibleAmount)

		fmt.Print("ID: ")

		fmt.Println(estimate.TradeId)

		break

	case 2:

		fmt.Println("Enter your TradeID: ")

		var tss int

		fmt.Scanf("%d", &tss)

		tradeID := UniqueId{tss}

		var estimate EstimatedAmount

		err = client.Call("SttockCall.DisplayPortfolio", tradeID, &estimate)

		if err != nil {

			log.Fatal("arith error:", err)

		}

		fmt.Println(estimate.MarketStocks)

		fmt.Println(estimate.AvailaibleAmount)

	}

}
