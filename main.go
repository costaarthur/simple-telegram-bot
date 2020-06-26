package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/yanzay/tbot"
)

type TickerData struct {
	Symbol     string `json:"symbol"`
	Name       string `json:"name"`
	Region     string `json:"region"`
	Currency   string `json:"currency"`
	MarketTime struct {
		Open     string `json:"open"`
		Close    string `json:"close"`
		Timezone int    `json:"timezone"`
	} `json:"market_time"`
	MarketCap     float64 `json:"market_cap"`
	Price         float64 `json:"price"`
	ChangePercent float64 `json:"change_percent"`
	UpdatedAt     string  `json:"updated_at"`
}

type StockData struct {
	By            string                `json:"by"`
	ValidKey      bool                  `json:"valid_key"`
	Results       map[string]TickerData `json:"results"`
	ExecutionTime float64               `json:"execution_time"`
	FromCache     bool                  `json:"from_cache"`
}

func getStock(s *tbot.Message) {
	getTicker := s.Vars["ticker"]
	stringTicker := string(getTicker)
	stringUpperTicker := strings.ToUpper(stringTicker)

	resp, err := http.Get("http://api.hgbrasil.com/finance/stock_price?key=" + os.Getenv("HGBRASIL_TOKEN") + "&symbol=" + getTicker)
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	bodyString := string(bodyBytes)
	fmt.Println("API Response as String:\n" + bodyString)

	var stockStruct StockData
	json.Unmarshal(bodyBytes, &stockStruct)
	fmt.Printf("API Response as struct %+v\n", stockStruct)

	s.Reply(stockStruct.Results[stringUpperTicker].Name)
	s.Reply("R$ " + fmt.Sprintf("%.2f", stockStruct.Results[stringUpperTicker].Price))
}

func timerHandler(m *tbot.Message) {
	// m.Vars contains all variables, parsed during routing
	secondsStr := m.Vars["seconds"]
	// Convert string variable to integer seconds value
	seconds, err := strconv.Atoi(secondsStr)
	if err != nil {
		m.Reply("Invalid number of seconds")
		return
	}
	m.Replyf("Timer for %d seconds started", seconds)
	time.Sleep(time.Duration(seconds) * time.Second)
	m.Reply("Time out!")
}

func main() {
	fmt.Println("Welcome home sir")
	bot, err := tbot.NewServer(os.Getenv("TELEGRAM_TOKEN"))
	if err != nil {
		log.Fatal(err)
	}
	bot.Handle("/answer", "42")
	bot.HandleFunc("/timer {seconds}", timerHandler)
	bot.HandleFunc("/stock {ticker}", getStock)

	bot.ListenAndServe()
}
