package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"net/rpc"
	"strconv"
	"strings"
)

type Call struct {
	List fulllist `json:"list"`
}

type tsc struct {
	Type     string `json:"-"`
	Begining int32  `json:"-"`
	Number   int32  `json:"-"`
}

type fulllist struct {
	Meta      tsc     `json:"-"`
	Resources []frscs `json:"resources"`
}

type frscs struct {
	Resource frsc `json:"resource"`
}

type frsc struct {
	Name   string `json:"classname"`
	Fields fd     `json:"fields"`
}

type fd struct {
	Price  string `json:"price"`
	Symbol string `json:"symbol"`
}

type StockDetails struct {
	Symbol string
	Budget float64
}

type EstimatedAmount struct {
	Stocks           string  `json:"stocksymbol"`
	AvailaibleAmount float64 `json:"stockvaluefetch"`
	TradeId          int     `json:"id"`
}

type Id struct {
	TradeId int `json:"id"`
}

type Update struct {
	Stocks           string  `json:"stocksymbol"`
	AvailaibleAmount float64 `json:"stockvaluefetch"`
}

type SttockCall int

var M map[int]EstimatedAmount

func (t *SttockCall) StockValueFetch(userInput *StockDetails, quote *EstimatedAmount) error {
	quote.TradeId++
	a := string(userInput.Symbol[:])

	a = strings.Replace(a, ":", ",", -1)
	a = strings.Replace(a, "%", ",", -1)
	a = strings.Replace(a, ",,", ",", -1)
	a = strings.Trim(a, " ")
	a = strings.Replace(a, "\"", "", -1)
	a = strings.TrimSpace(a)
	a = strings.TrimSuffix(a, ",")
	StackArray := strings.Split(a, ",")

	Total := 0.0
	var Url string

	for i := 0; i < len(StackArray); i++ {
		i = i + 1

		storevalue, _ := strconv.ParseFloat(StackArray[i], 64)
		Total = (storevalue * userInput.Budget * 0.01)
		fmt.Println(StackArray[i-1], Total)
		Url = Url + (StackArray[i-1] + ",")

	}

	Url = strings.TrimSuffix(Url, ",")

	CompleteUrl := "http://finance.yahoo.com/webservice/v1/symbols/" + Url + "/quote?format=json"

	client := &http.Client{}

	respondvalue, _ := client.Get(CompleteUrl)
	rplaced, _ := http.NewRequest("GET", CompleteUrl, nil)

	rplaced.Header.Add("If-None-Match", "application/json")
	rplaced.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	// make request
	respondvalue, _ = client.Do(rplaced)
	if respondvalue.StatusCode >= 200 && respondvalue.StatusCode < 300 {
		var C Call
		body, _ := ioutil.ReadAll(respondvalue.Body)

		err := json.Unmarshal(body, &C)

		n := len(StackArray)

		Est := make([]float64, n, n)

		for i := 0; i < n; i++ {
			i = i + 1
			Temp, _ := strconv.ParseFloat(StackArray[i], 64)
			Est[i] = (Temp * userInput.Budget * 0.01)
			fmt.Println(Est)
			fmt.Println(StackArray[i-1], Est[i])
		}

		var buffer bytes.Buffer
		q := 0
		for _, Sample := range C.List.Resources {

			getsample := Sample.Resource.Fields.Symbol
			floatparse, _ := strconv.ParseFloat(Sample.Resource.Fields.Price, 64)
			typecast := (int)(Est[q+1] / floatparse)
			mathsymbol := math.Mod(Est[q+1], floatparse)
			q = q + 2

			quote.Stocks = fmt.Sprintf("%s:%g:%d", getsample, floatparse, typecast)
			quote.AvailaibleAmount = quote.AvailaibleAmount + mathsymbol
			buffer.WriteString(quote.Stocks)
			buffer.WriteString(",")
		}

		quote.Stocks = (buffer.String())
		quote.Stocks = strings.TrimSuffix(quote.Stocks, ",")
		fmt.Println(quote.Stocks)
		fmt.Println(quote.AvailaibleAmount)

		M = map[int]EstimatedAmount{
			quote.TradeId: {quote.Stocks, quote.AvailaibleAmount, quote.TradeId},
		}

		fmt.Print(M)
		if err == nil {
			fmt.Println("Completed")
		}
	} else {
		fmt.Println(respondvalue.Status)

	}

	return nil
}

func (t *SttockCall) Portfo(id *Id, qt *EstimatedAmount) error {
	if id.TradeId == 1 {
		*qt = M[id.TradeId]
	}

	return nil

}

func main() {

	SttockCall := new(SttockCall)
	rpc.Register(SttockCall)

	rpc.HandleHTTP()

	err := http.ListenAndServe(":1528", nil)
	if err != nil {
		fmt.Println(err.Error())
	}
}
