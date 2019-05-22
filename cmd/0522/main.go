package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/jinzhu/now"
)

const (
	endPoint = "https://hist-quote.1tokentrade.cn"
	format   = "2006-01-02"
	otKey    = "XXX"
)

var exchanges = []string{
	"okex",
	// "huobip",
}

var symbols = []string{
	"btc.usdt",
	// "eth.usdt",
	// "xrp.usdt",
	// "eos.usdt",
	// "ltc.usdt",
	// "etc.usdt",
	// "bch.usdt",
	// "trx.usdt",
	// "bsv.usdt",
}

type Candle struct {
	Timestamp int64   `json:"timestamp"`
	Open      float64 `json:"open"`
	High      float64 `json:"high"`
	Low       float64 `json:"low"`
	Close     float64 `json:"close"`
	Volume    float64 `json:"volume"`
}

func main() {
	for _, exchange := range exchanges {
		fmt.Println("开始处理交易所", exchange)
		for _, symbol := range symbols {
			fmt.Println("开始处理交易对", symbol)

			// 开始时间是 2018-04-01
			// t := time.Date(2018, 3, 27, 00, 00, 01, 123456789, time.Now().Location())
			t := time.Date(2019, 4, 27, 00, 00, 01, 123456789, time.Now().Location())

			// 缓存需要写入文件的字符串
			// var tmpSlice []string

			for {
				// 计算下一个周五
				t = getNextFriday(t)

				// 计算出来的周五超过今天，说明是未来的日子，直接中断
				if t.After(time.Now()) {
					// fmt.Println(tmpSlice)
					// writeLines(tmpSlice, exchange+"-"+symbol+".csv")
					break
				}

				// fmt.Println(t.Format(format))
				result := getCandleClose(exchange+"/"+symbol, t.Format(format), t.AddDate(0, 0, 1).Format(format))
				// tmpSlice = append(tmpSlice, result)
				fmt.Println(result)
				time.Sleep(time.Millisecond * 500)
			}
		}
	}
}

func getCandleClose(contract, since, until string) (result string) {
	var candles []Candle
	contents, _ := getCandles(contract, since, until)
	json.Unmarshal(contents, &candles)

	fmt.Println(contents)
	fmt.Println(candles)

	for _, candle := range candles {
		hour := time.Unix(candle.Timestamp, 0).Hour()
		min := time.Unix(candle.Timestamp, 0).Minute()
		if hour == 16 && min == 2 {
			fmt.Println(candle)
			// fmt.Println(printCandle(candle))
			result = printCandle(candle)
			break
		}
	}
	return result
}

func printCandle(c Candle) string {
	date := time.Unix(c.Timestamp, 0).Format("2006-01-02T15:04:05Z07:00")
	s := fmt.Sprintf("%s,%g", date, c.Close)
	// return date + "," + strconv.FormatFloat(c.Close, 'f', 6, 64)
	fmt.Println(s)
	return s
}

func getCandles(contract, since, until string) ([]byte, error) {
	url := endPoint + packagePATH(contract, since, until)

	client := &http.Client{}
	fmt.Println(url)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return []byte{}, err
	}

	req.Header.Set("ot-key", otKey)

	resp, err := client.Do(req)

	if err != nil {
		log.Println(err)
		return []byte{}, err
	}

	defer resp.Body.Close() //关闭
	body, err := ioutil.ReadAll(resp.Body)
	// fmt.Println(string(body))
	return body, nil
}

func packagePATH(contract, since, until string) string {
	return "/candles?contract=" + contract + "&since=" + since + "&until=" + until + "&duration=1m&format=json"
}

// NOTE: 传入的日期在周 5，6，7 的时候，此函数工作正常
func getNextFriday(date time.Time) time.Time {
	return now.New(date).Monday().AddDate(0, 0, 11)
}

func writeLines(lines []string, outputFile string) error {
	// overwrite file if it exists
	file, err := os.OpenFile(outputFile, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	check(err)
	defer file.Close()

	// new writer w/ default 4096 buffer size
	w := bufio.NewWriter(file)

	fmt.Println(lines)
	for _, line := range lines {
		fmt.Println(line)
		_, err := w.WriteString(line + "\n")
		check(err)
	}

	// flush outstanding data
	return w.Flush()
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
