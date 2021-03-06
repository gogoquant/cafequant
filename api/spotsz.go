package api

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/axgle/mahonia"
	simplejson "github.com/bitly/go-simplejson"
	"snack.com/xiyanxiyan10/stocktrader/constant"
	"snack.com/xiyanxiyan10/stocktrader/util"
)

const (
	timeTemplate1 = "2006-01-02 15:04:05"
	tickerURL     = "http://hq.sinajs.cn/list="
	//depthURL      = "http://hq.sinajs.cn/list="
	recordURL = "http://money.finance.sina.com.cn/quotes_service/api/json_v2.php/CN_MarketData.getKLineData?"
)

// NewSZExchange create an exchange struct of futureExchange.com
func NewSZExchange(opt constant.Option) (Exchange, error) {
	exchange := NewSZSpotExchange(opt)
	exchange.SetRecordsPeriodMap(map[string]int64{
		"M5":  5,
		"M15": 15,
		"M30": 30,
	})
	if err := exchange.Init(opt); err != nil {
		return nil, err
	}
	return exchange, nil
}

// getRecords ...
func getRecords(symbol string, period, ma, size int) (string, error) {
	client := &http.Client{}
	url := recordURL
	url = url + "symbol=" + symbol + "&scale=" + strconv.Itoa(period) + "&ma=" + strconv.Itoa(ma) + "&datalen=" + strconv.Itoa(size)
	fmt.Printf("call address:%s\n", url)
	reqest, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}
	reqest.Header.Add("Accept-Language", "zh-CN,zh;q=0.8")
	reqest.Header.Add("Content-Type", "text/html; charset=utf-8")
	resp, err := client.Do(reqest)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

// pareseRecords ...
func pareseRecords(str string, ma int) []constant.Record {
	var records []constant.Record
	res, err := simplejson.NewJson([]byte(str))
	if err != nil {
		return records
	}
	rows, err := res.Array()
	if err != nil {
		return records
	}
	for _, row := range rows {
		if maps, ok := row.(map[string]interface{}); ok {
			var record constant.Record
			record.Open = util.Float64Must(maps["open"])
			record.Close = util.Float64Must(maps["close"])
			record.High = util.Float64Must(maps["high"])
			record.Low = util.Float64Must(maps["low"])
			record.Volume = util.Float64Must(maps["volume"])

			//record.MaPrice = util.Float64Must(maps["ma_price"+strconv.Itoa(ma)])
			//record.MaVolume = util.Float64Must(maps["ma_volume"+strconv.Itoa(ma)])
			records = append(records, record)
		}
	}
	return records
}

// getTickerAndDepth ...
func getTickerAndDepth(no string) (string, error) {
	client := &http.Client{}
	url := tickerURL + no
	fmt.Printf("ticker call address:%s\n", url)
	reqest, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return "", err
	}
	reqest.Header.Add("Accept-Language", "zh-CN,zh;q=0.8")
	reqest.Header.Add("Content-Type", "text/html; charset=utf-8")
	resp, err := client.Do(reqest)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	dec := mahonia.NewDecoder("GBK")
	s := dec.ConvertString(string(data))
	s = s[strings.Index(s, "\"")+1 : strings.LastIndex(s, "\"")]
	//s = tickerFormat(s)
	return s, nil
}

// parseTicker ...
func parseTicker(data string) (*constant.Ticker, error) {
	var ticker constant.Ticker
	arr := strings.Split(data, ",")
	if len(arr) < 32 {
		return nil, fmt.Errorf("arr len too smail:%d", len(arr))
	}
	ticker.Open = util.Float64Must(arr[1])
	ticker.Close = util.Float64Must(arr[2])
	ticker.Last = util.Float64Must(arr[3])
	ticker.High = util.Float64Must(arr[4])
	ticker.Low = util.Float64Must(arr[5])
	ticker.Buy = util.Float64Must(arr[6])
	ticker.Sell = util.Float64Must(arr[7])
	ticker.Vol = util.Float64Must(arr[8])
	t := arr[30] + " " + arr[31]
	stamp, err := time.ParseInLocation(timeTemplate1, t, time.Local)
	if err != nil {
		return nil, err
	}
	ticker.Time = stamp.Unix()
	return &ticker, nil
}

// parseDepth ...
func parseDepth(data string) (*constant.Depth, error) {
	var depth constant.Depth
	arr := strings.Split(data, ",")
	if len(arr) < 31 {
		return nil, fmt.Errorf("len too samall")
	}
	for i := 0; i < 5; i++ {
		var record constant.DepthRecord
		record.Amount = util.Float64Must(arr[10+2*i])
		record.Price = util.Float64Must(arr[10+2*i+1])
		depth.Bids = append(depth.Bids, record)

		record.Amount = util.Float64Must(arr[20+2*i])
		record.Price = util.Float64Must(arr[20+2*i+1])
		depth.Asks = append(depth.Asks, record)
	}
	t := arr[30] + " " + arr[31]
	stamp, err := time.ParseInLocation(timeTemplate1, t, time.Local)
	if err != nil {
		return nil, err
	}
	depth.Time = stamp.Unix()
	return &depth, nil
}

// SZExchange the exchange struct of futureExchange.com
type SZExchange struct {
	BaseExchange

	tradeTypeMap        map[int]string
	tradeTypeMapReverse map[string]int
	exchangeTypeMap     map[string]string

	records map[string][]constant.Record
}

// NewSZSpotExchange create an exchange struct of futureExchange.com
func NewSZSpotExchange(opt constant.Option) *SZExchange {
	spotExchange := SZExchange{
		records: make(map[string][]constant.Record),
	}
	spotExchange.SetMinAmountMap(map[string]float64{
		"BTC/USD": 0.001,
	})
	return &spotExchange
}

// Stop ...
func (e *SZExchange) Stop() error {
	return nil
}

// Start ...
func (e *SZExchange) Start() error {
	return nil
}

// Init get the type of this exchange
func (e *SZExchange) Init(opt constant.Option) error {
	e.BaseExchange.Init(opt)
	return nil
}

// GetType get the type of this exchange
func (e *SZExchange) GetType() string {
	return e.option.Type
}

// GetName get the name of this exchange
func (e *SZExchange) GetName() string {
	return e.option.Name
}

// GetPosition ...
func (e *SZExchange) GetPosition() ([]constant.Position, error) {
	return nil, fmt.Errorf("not support")
}

// GetAccount get the account detail of this exchange
func (e *SZExchange) GetAccount() (*constant.Account, error) {
	return nil, fmt.Errorf("not support")
}

// Buy ...
func (e *SZExchange) Buy(price, amount string, msg string) (string, error) {
	return "", fmt.Errorf("not support")
}

// Sell ...
func (e *SZExchange) Sell(price, amount string, msg string) (string, error) {
	return "", fmt.Errorf("not support")
}

// GetOrder get details of an order
func (e *SZExchange) GetOrder(id string) (*constant.Order, error) {
	return nil, fmt.Errorf("not support")
}

// GetOrders get all unfilled orders
func (e *SZExchange) GetOrders() ([]constant.Order, error) {
	return nil, fmt.Errorf("not support")
}

// CancelOrder cancel an order
func (e *SZExchange) CancelOrder(orderID string) (bool, error) {
	return false, fmt.Errorf("not support")
}

// GetTicker get market ticker
func (e *SZExchange) GetTicker() (*constant.Ticker, error) {
	stockType := e.GetStockType()
	res, err := getTickerAndDepth(stockType)
	if err != nil {
		e.logger.Log(constant.ERROR, e.GetStockType(), 0.0, 0.0, "GetTicker() error, the error number is ", err.Error())
		return nil, fmt.Errorf("GetTicker() error, the error number is %s ", err.Error())
	}
	ticker, err := parseTicker(res)
	if err != nil {
		e.logger.Log(constant.ERROR, e.GetStockType(), 0.0, 0.0, "GetTicker() error, the error number is ", err.Error())
		return nil, fmt.Errorf("GetTicker() error, the error number is %s ", err.Error())
	}
	if ticker == nil {
		e.logger.Log(constant.ERROR, e.GetStockType(), 0.0, 0.0, "GetTicker() error, the error number is ticker parse fail")
		return nil, fmt.Errorf("GetTicker() error, the error number is ticker parse fail")
	}
	return ticker, nil
}

// GetDepth ...
func (e *SZExchange) GetDepth() (*constant.Depth, error) {
	stockType := e.GetStockType()
	res, err := getTickerAndDepth(stockType)
	if err != nil {
		e.logger.Log(constant.ERROR, e.GetStockType(), 0.0, 0.0, "GetDepth() error, the error number is:%s\n", err.Error())
		return nil, fmt.Errorf("GetDepth() error, the error number is %s", err.Error())
	}
	depth, err := parseDepth(res)
	if err != nil {
		e.logger.Log(constant.ERROR, e.GetStockType(), 0.0, 0.0, "GetDepth() error, the error number is:%s\n", err.Error())
		return nil, fmt.Errorf("GetDepth() error, the error number is depth parse fail")
	}
	return depth, nil
}
