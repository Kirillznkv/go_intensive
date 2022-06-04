package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

const myUrl string = "http://localhost:3333/buy_candy"

type BuyRequest struct {
	Money      int    `json:"money"`
	CandyType  string `json:"candyType"`
	CandyCount int    `json:"candyCount"`
}

type BuyResponce struct {
	Change int    `json:"change"`
	Thanks string `json:"thanks"`
}

func getJsonRequest(r io.ReadCloser) (b BuyRequest) {
	if err := json.NewDecoder(r).Decode(&b); err != nil {
		log.Fatal(err)
	}
	return b
}

func checkBuy(b BuyRequest) (j BuyResponce, err error) {
	prices := map[string]int{"CE": 10, "AA": 15, "NT": 17, "DE": 21, "YR": 23}
	p := 0
	if v, f := prices[b.CandyType]; f {
		p = v * b.CandyCount
	} else {
		return j, fmt.Errorf("no such position")
	}
	if p > b.Money {
		return j, fmt.Errorf("not enough money")
	}
	j.Change = b.Money - p
	j.Thanks = "Thank you!"
	return j, nil
}

func buyCandy(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	default:
		fmt.Fprintf(w, "Sorry, only POST methods are supported.")
	case "POST":
		reqData := getJsonRequest(r.Body)
		res, err := checkBuy(reqData)
		if err == nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(res)
		} else if err.Error() == "not enough money" {
			w.WriteHeader(http.StatusPaymentRequired)
			w.Write([]byte(err.Error()))
		} else if err.Error() == "no such position" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
		}
	}
}

func main() {
	http.HandleFunc("/buy_candy", buyCandy)
	http.ListenAndServe(":3333", nil)
}
