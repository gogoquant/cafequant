package constant

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

var client = http.DefaultClient

// BackTime backtest time set
type BackTime struct {
	Start  int64
	End    int64
	Period string
}

// Position struct
type Position struct {
	Price        float64 //价格
	MarginLevel  float64 //杠杆比例
	Amount       float64 //总合约数量
	Available    float64 //可平仓量
	FrozenAmount float64 //冻结的合约数量
	Profit       float64 //收益
	ProfitRate   float64 //收益率
	ContractType string  //合约类型
	TradeType    string  //交易类型
	Margin       float64 //仓位占用的保证金
	StockType    string  //货币类型
	ForcePrice   float64 //强制平仓价格
}

// DepthRecord ...
type DepthRecord struct {
	Price  float64
	Amount float64
}

// SubAccount ...
type SubAccount struct {
	StockType     string
	AccountRights float64 //账户权益
	KeepDeposit   float64 //保证金
	ProfitReal    float64 //已实现盈亏
	ProfitUnreal  float64
	RiskRate      float64 //保证金率
	Amount        float64
	FrozenAmount  float64
	LoanAmount    float64
}

// TradeStatus ..
type TradeStatus int

func (ts TradeStatus) String() string {
	return tradeStatusSymbol[ts]
}

var tradeStatusSymbol = [...]string{"UNFINISH", "PART_FINISH", "FINISH", "CANCEL", "REJECT", "CANCEL_ING", "FAIL"}

type Account struct {
	SubAccounts map[string]SubAccount
}

type DepthRecords []DepthRecord

type Depth struct {
	ContractType string //for future
	StockType    string
	Time         int64
	Asks         DepthRecords // Descending order
	Bids         DepthRecords // Descending order
}

// Order struct
type Order struct {
	Id        string  //订单ID
	Price     float64 //价格
	OpenPrice float64 //open price
	AvgPrice  float64

	Amount     float64 //总量
	DealAmount float64 //成交量
	Fee        float64 //这个订单的交易费
	TradeType  string  //交易类型
	StockType  string  //货币类型
	//ContractUnit int64   //对应张数

	Time         int64
	FinishedTime int64

	Status TradeStatus // trader status
}

// OHLC is a candlestick struct
type OHLC struct {
	Time   int64   `json:"Time"`
	Open   float64 `json:"Open"`
	High   float64 `json:"High"`
	Low    float64 `json:"Low"`
	Close  float64 `json:"Close"`
	Volume float64 `json:"Volume"`
}

type Record OHLC

// Option is an exchange option
type Option struct {
	Index     int
	TraderID  int64
	Type      string
	Name      string
	AccessKey string
	SecretKey string

	Limit     int64
	LastSleep int64
	LastTimes int64
	WatchList []string

	host string

	BackTest bool // 是否开启回测
	BackLog  bool // 是否将日志输出到终端，而不是数据库
}

// OrderBook struct
/*
type OrderBook struct {
	Price  float64 //价格
	Amount float64 //市场深度量
}
*/

// Ticker  ...
type Ticker struct {
	Last  float64
	Buy   float64
	Sell  float64
	Open  float64
	Close float64
	High  float64
	Low   float64
	Vol   float64
	Time  int64
}

// Trader ...
type Trader struct {
	Id        int64
	TradeType string
	Amount    float64
	Price     float64
	StockType string
	Time      int64
}

func base64Encode(data string) string {
	return base64.StdEncoding.EncodeToString([]byte(data))
}

func signMd5(params []string) string {
	m := md5.New()
	m.Write([]byte(strings.Join(params, "&")))
	return hex.EncodeToString(m.Sum(nil))
}

func signSha512(params []string, key string) string {
	h := hmac.New(sha512.New, []byte(key))
	h.Write([]byte(strings.Join(params, "&")))
	return hex.EncodeToString(h.Sum(nil))
}

func signSha1(params []string, key string) string {
	h := hmac.New(sha1.New, []byte(key))
	h.Write([]byte(strings.Join(params, "&")))
	return hex.EncodeToString(h.Sum(nil))
}

func signChbtc(params []string, key string) string {
	sha := sha1.New()
	sha.Write([]byte(key))
	secret := hex.EncodeToString(sha.Sum(nil))
	h := hmac.New(md5.New, []byte(secret))
	h.Write([]byte(strings.Join(params, "&")))
	return hex.EncodeToString(h.Sum(nil))
}

func post_gateio(url string, data []string, key string, sign string) (ret []byte, err error) {
	req, err := http.NewRequest("POST", url, strings.NewReader(strings.Join(data, "&")))
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("key", key)
	req.Header.Set("sign", sign)
	resp, err := client.Do(req)
	if resp == nil {
		err = fmt.Errorf("[POST %s] HTTP Error Info: %v", url, err)
	} else if resp.StatusCode == 200 {
		ret, _ = ioutil.ReadAll(resp.Body)
		resp.Body.Close()
	} else {
		err = fmt.Errorf("[POST %s] HTTP Status: %d, Info: %v", url, resp.StatusCode, err)
	}
	return ret, err
}

func post(url string, data []string) (ret []byte, err error) {
	req, err := http.NewRequest("POST", url, strings.NewReader(strings.Join(data, "&")))
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := client.Do(req)
	if resp == nil {
		err = fmt.Errorf("[POST %s] HTTP Error Info: %v", url, err)
	} else if resp.StatusCode == 200 {
		ret, _ = ioutil.ReadAll(resp.Body)
		resp.Body.Close()
	} else {
		err = fmt.Errorf("[POST %s] HTTP Status: %d, Info: %v", url, resp.StatusCode, err)
	}
	return ret, err
}

func get(url string) (ret []byte, err error) {
	req, err := http.NewRequest("GET", url, strings.NewReader(""))
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := client.Do(req)
	if resp == nil {
		err = fmt.Errorf("[GET %s] HTTP Error Info: %v", url, err)
	} else if resp.StatusCode == 200 {
		ret, _ = ioutil.ReadAll(resp.Body)
		resp.Body.Close()
	} else {
		err = fmt.Errorf("[GET %s] HTTP Status: %d, Info: %v", url, resp.StatusCode, err)
	}
	return ret, err
}
