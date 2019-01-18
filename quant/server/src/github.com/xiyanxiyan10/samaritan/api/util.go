package api

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	goback "github.com/xiyanxiyan10/gobacktest"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

var client = http.DefaultClient

// Position struct
type Position struct {
	Price         float64
	Leverage      int
	Amount        float64
	ConfirmAmount float64
	FrozenAmount  float64
	Profit        float64
	ContractType  string
	TradeType     string
	StockType     string
}

// Order struct
type Order struct {
	ID         string
	Price      float64
	Amount     float64
	DealAmount float64
	Fee        float64
	TradeType  string
	StockType  string
}

// Record struct
type Record struct {
	Time   int64
	Open   float64
	High   float64
	Low    float64
	Close  float64
	Volume float64
}

// OrderBook struct
type OrderBook struct {
	Price  float64
	Amount float64
}

// Ticker struct
type Ticker struct {
	goback.Event

	high   float64 // 最高价
	mid    float64 // 中间价
	low    float64 // 最低价
	amount float64 // 成交量
	count  int64   // 成交笔数
	open   float64 // 开盘价
	close  float64 // 收盘价
	vol    float64 // 成交额
	asks   []OrderBook
	bids   []OrderBook
}

func (t Ticker) Mid() float64 {
	return t.mid
}

func (t Ticker) Open() float64 {
	return t.open
}

func (t Ticker) Close() float64 {
	return t.close
}

func (t Ticker) Low() float64 {
	return t.low
}

func (t Ticker) High() float64 {
	return t.high
}

func (t *Ticker) SetOpen(open float64) {
	t.open = open
}

func (t *Ticker) SetClose(close float64) {
	t.close = close
}

func (t *Ticker) SetMid(mid float64) {
	t.mid = mid
}

func (t *Ticker) SetLow(low float64) {
	t.low = low
}

func (t *Ticker) SetHigh(high float64) {
	t.high = high
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

// IncomingHandler ...
type IncomingHandler struct {
	// incomeing data for user api
	incoming chan goback.EventHandler
}

// NewIncomingHandler ...
func NewIncomingHandler(len int) *IncomingHandler {
	return &IncomingHandler{incoming: make(chan goback.EventHandler, len)}
}

// receiver event
func (h *IncomingHandler) SyncReceive() (data goback.EventHandler, err error) {
	data = <-h.incoming
	return data, nil
}

// Send event
func (h *IncomingHandler) SyncSend(data goback.EventHandler) error {
	h.incoming <- data
	return nil
}

// Receive event
func (h *IncomingHandler) Receive() (event goback.EventHandler, err error) {
	select {
	case event = <-h.incoming:
		err = nil
		return
	default:
		err = errors.New("nonblock receive fail")
		return
	}

}

// Send event
func (h *IncomingHandler) Send(data goback.EventHandler) (err error) {
	select {
	case h.incoming <- data:
		err = nil
		return
	default:
		err = errors.New("nonblock send fail")
		return
	}
}

// TimeoutReceive event
func (h *IncomingHandler) TimeoutReceive(sec int64) (event goback.EventHandler, err error) {
	select {
	case event = <-h.incoming:
		err = nil
		return
	case <-time.After(time.Second * time.Duration(sec)):
		err = errors.New("timeout receive fail")
		return
	}

}

// TimeoutSend event
func (h *IncomingHandler) TimeoutSend(data goback.EventHandler, sec int64) (err error) {
	select {
	case h.incoming <- data:
		err = nil
		return
	case <-time.After(time.Second * time.Duration(sec)):
		err = errors.New("timeout send fail")
		return
	}
}