package api

import (
	"errors"
	"github.com/axgle/mahonia"
	simplejson "github.com/bitly/go-simplejson"
	goex "github.com/nntaoli-project/goex"
	"io/ioutil"
	"net/http"
	"snack.com/xiyanxiyan10/conver"
	"snack.com/xiyanxiyan10/stocktrader/constant"
	"strconv"
	"strings"
	"time"
)

const (
	timeTemplate1 = "2006-01-02 15:04:05"
	tickerURL     = "http://hq.sinajs.cn/list="
	//depthURL      = "http://hq.sinajs.cn/list="
	recordURL = "http://money.finance.sina.com.cn/quotes_service/api/json_v2.php/CN_MarketData.getKLineData?"
)

// getRecords ...
func getRecords(symbol string, period, ma, size int) (string, error) {
	client := &http.Client{}
	url := recordURL
	url = url + "symbol=" + symbol + "&scale=" + strconv.Itoa(period) + "&ma=" + strconv.Itoa(ma) + "&datalen=" + strconv.Itoa(size)
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
			record.Open = conver.Float64Must(maps["open"])
			record.Close = conver.Float64Must(maps["close"])
			record.High = conver.Float64Must(maps["high"])
			record.Low = conver.Float64Must(maps["low"])
			record.Volume = conver.Float64Must(maps["volume"])

			record.MaPrice = conver.Float64Must(maps["ma_price"+strconv.Itoa(ma)])
			record.MaVolume = conver.Float64Must(maps["ma_volume"+strconv.Itoa(ma)])
			records = append(records, record)
		}
	}
	return records
}

// getTickerAndDepth ...
func getTickerAndDepth(no string) (string, error) {
	client := &http.Client{}
	url := ""
	if strings.HasPrefix(no, "60") || strings.HasPrefix(no, "51") {
		url = tickerURL + "sh" + no
	} else if strings.HasPrefix(no, "00") || strings.HasPrefix(no, "30") {
		url = tickerURL + "sz" + no
	} else {
		return "", errors.New("Unknown stock")
	}

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
func parseTicker(data string) *constant.Ticker {
	var ticker constant.Ticker
	arr := strings.Split(data, ",")
	if len(arr) < 32 {
		return nil
	}
	ticker.Open = conver.Float64Must(arr[1])
	ticker.Close = conver.Float64Must(arr[2])
	ticker.Last = conver.Float64Must(arr[3])
	ticker.High = conver.Float64Must(arr[4])
	ticker.Low = conver.Float64Must(arr[5])
	ticker.Buy = conver.Float64Must(arr[6])
	ticker.Sell = conver.Float64Must(arr[7])
	ticker.Vol = conver.Float64Must(arr[8])
	t := arr[30] + " " + arr[31]
	stamp, _ := time.ParseInLocation(timeTemplate1, t, time.Local)
	ticker.Time = stamp.Unix()
	return &ticker
}

// parseDepth ...
func parseDepth(data string) *constant.Depth {
	var depth constant.Depth
	arr := strings.Split(data, ",")
	if len(arr) < 31 {
		return nil
	}
	for i := 0 + 10; i < 5; i++ {
		var record constant.DepthRecord
		record.Amount = conver.Float64Must(arr[10+2*i])
		record.Price = conver.Float64Must(arr[10+2*i+1])
		depth.Bids = append(depth.Bids, record)

		record.Amount = conver.Float64Must(arr[20+2*i])
		record.Price = conver.Float64Must(arr[20+2*i+1])
		depth.Bids = append(depth.Bids, record)
	}
	t := arr[30] + " " + arr[31]
	stamp, _ := time.ParseInLocation(timeTemplate1, t, time.Local)
	depth.Time = stamp.Unix()
	return &depth
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
	spotExchange.SetRecordsPeriodMap(map[string]int64{
		"M1":  60,
		"M5":  300,
		"M15": 900,
		"M30": 1800,
		"H1":  3600,
		"H2":  7200,
	})
	spotExchange.SetRecordsPeriodMap(map[string]int64{
		"M1":  goex.KLINE_PERIOD_1MIN,
		"M5":  goex.KLINE_PERIOD_5MIN,
		"M15": goex.KLINE_PERIOD_15MIN,
		"M30": goex.KLINE_PERIOD_30MIN,
		"H1":  goex.KLINE_PERIOD_1H,
		"H2":  goex.KLINE_PERIOD_4H,
		"H4":  goex.KLINE_PERIOD_4H,
		"D1":  goex.KLINE_PERIOD_1DAY,
		"W1":  goex.KLINE_PERIOD_1WEEK,
	})
	spotExchange.SetMinAmountMap(map[string]float64{
		"BTC/USD": 0.001,
	})
	spotExchange.back = false
	return &spotExchange
}

// Ready ...
func (e *SZExchange) Ready() interface{} {
	return "success"
}

// Init get the type of this exchange
func (e *SZExchange) Init(opt constant.Option) error {
	e.BaseExchange.Init(opt)
	return nil
}

// Log print something to console
func (e *SZExchange) Log(msgs ...interface{}) {
	e.logger.Log(constant.INFO, e.GetStockType(), 0.0, 0.0, msgs...)
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
func (e *SZExchange) GetPosition() interface{} {
	return nil
}

// GetAccount get the account detail of this exchange
func (e *SZExchange) GetAccount() interface{} {
	return nil
}

// Buy ...
func (e *SZExchange) Buy(price, amount string, msg ...interface{}) interface{} {
	return nil
}

// Sell ...
func (e *SZExchange) Sell(price, amount string, msg ...interface{}) interface{} {
	return nil
}

// GetOrder get details of an order
func (e *SZExchange) GetOrder(id string) interface{} {
	return nil
}

// GetOrders get all unfilled orders
func (e *SZExchange) GetOrders() interface{} {
	return nil
}

// CancelOrder cancel an order
func (e *SZExchange) CancelOrder(orderID string) interface{} {
	return nil
}

// GetTicker get market ticker
func (e *SZExchange) GetTicker() interface{} {
	stockType := e.GetStockType()
	res, err := getTickerAndDepth(stockType)
	if err != nil {
		e.logger.Log(constant.ERROR, e.GetStockType(), 0.0, 0.0, "GetTicker() error, the error number is ", err.Error())
		return nil
	}
	ticker := parseTicker(res)
	if ticker == nil {
		e.logger.Log(constant.ERROR, e.GetStockType(), 0.0, 0.0, "GetTicker() error, the error number is ticker parse fail")
		return nil
	}
	return ticker
}

// GetRecords get candlestick data
func (e *SZExchange) GetRecords(periodStr, maStr string) interface{} {
	exchangeStockType := e.GetStockType()
	var period int64 = -1
	var size = constant.RecordSize
	period, ok := e.recordsPeriodMap[periodStr]
	if !ok {
		e.logger.Log(constant.ERROR, e.GetStockType(), 0, 0, "GetRecords() error, the error number is stockType")
		return nil
	}
	ma, err := strconv.Atoi(maStr)
	if err != nil {
		e.logger.Log(constant.ERROR, e.GetStockType(), 0, 0, "GetRecords() error, the error number is stockType")
		return nil
	}
	res, err := getRecords(exchangeStockType, int(period), ma, size)
	if err != nil {
		e.logger.Log(constant.ERROR, e.GetStockType(), 0.0, 0.0, "GetRecords() error, the error number is ", err.Error())
		return nil
	}
	return pareseRecords(res, ma)
}

// GetDepth ...
func (e *SZExchange) GetDepth() interface{} {
	stockType := e.GetStockType()
	res, err := getTickerAndDepth(stockType)
	if err != nil {
		e.logger.Log(constant.ERROR, e.GetStockType(), 0.0, 0.0, "GetDepth() error, the error number is ", err.Error())
		return nil
	}
	depth := parseDepth(res)
	if depth == nil {
		e.logger.Log(constant.ERROR, e.GetStockType(), 0.0, 0.0, "GetDepth() error, the error number is depth parse fail")
		return nil
	}
	return depth
}
