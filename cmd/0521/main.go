package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/jinzhu/now"
)

func getFriday(date time.Time) time.Time {
	return now.New(date).Monday().AddDate(0, 0, 4)
}

// NOTE: 传入的日期在周 5，6，7 的时候，此函数工作正常
func getNextFriday(date time.Time) time.Time {
	return now.New(date).Monday().AddDate(0, 0, 11)
}

var symbols = []string{
	"BTC-USDT",
	"ETH-USDT",
	"XRP-USDT",
	"EOS-USDT",
	"LTC-USDT",
	"ETC-USDT",
	"BCH-USDT",
	"TRX-USDT",
	"BSV-USDT",
}
var format = "2006-01-02"

func main() {
	for _, symbol := range symbols {
		fmt.Println("开始处理交易对", symbol)

		// 开始时间是 2019-02-21
		t := time.Date(2019, 2, 21, 00, 00, 01, 123456789, time.Now().Location())

		for {
			// 计算下一个周五
			t = getNextFriday(t)

			// 计算出来的周五超过今天，说明是未来的日子，直接中断
			if t.After(time.Now()) {
				break
			}
			getURL(symbol, t.Format(format), t.AddDate(0, 0, 1).Format(format))
			time.Sleep(time.Millisecond * 500)
		}
	}
}

func getURL(symbol, start, end string) {
	url := "https://www.okex.com/api/spot/v3/instruments/" + symbol + "/candles?granularity=3600&start=" + start + "T08%3A28%3A48.899Z&end=" + end + "T09%3A28%3A48.899Z"
	// fmt.Println(url)
	contents := getResponseData(url)
	var datas [][]string
	json.Unmarshal(contents, &datas) //Parse JSON data and stores the result in var users
	for _, data := range datas {
		if data[0] == start+"T16:00:00.000Z" {
			fmt.Println(data[0], data[1])
		}
	}
}

func getResponseData(url string) []byte {
	response, err := http.Get(url)
	if err != nil {
		fmt.Print(err)
	}
	defer response.Body.Close()
	contents, _ := ioutil.ReadAll(response.Body)
	return contents
}
