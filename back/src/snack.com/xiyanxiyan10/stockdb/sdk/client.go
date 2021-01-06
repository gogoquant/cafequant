package sdk

import (
	"encoding/base64"
	"github.com/hprose/hprose-golang/rpc"
	"net/http"
	"snack.com/xiyanxiyan10/stockdb/types"
)

// Client Client of StockDB
type Client struct {
	uri    string
	auth   string
	hprose *rpc.HTTPClient

	PutOHLC        func(datum types.OHLC, opt types.Option) types.BaseResponse
	PutOHLCs       func(data []types.OHLC, opt types.Option) types.BaseResponse
	PutOrder       func(datum types.Order, opt types.Option) types.BaseResponse
	PutOrders      func(data []types.Order, opt types.Option) types.BaseResponse
	GetStats       func() types.StatsResponse
	GetMarkets     func() types.StringsResponse
	GetSymbols     func(market string) types.StringsResponse
	GetTimeRange   func(opt types.Option) types.TimeRangeResponse
	GetPeriodRange func(opt types.Option) types.TimeRangeResponse
	GetOHLCs       func(opt types.Option) types.OHLCResponse
	GetDepth       func(opt types.Option) types.DepthResponse
}

func (c *Client) init(uri, auth string) {
	c.uri = uri
	c.auth = auth
	c.hprose = rpc.NewHTTPClient(c.uri)
	if auth != "" {
		c.hprose.Header = make(http.Header)
		c.hprose.Header.Set("content-type", "application/json")
		c.hprose.Header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(auth)))
	}
	c.hprose.UseService(c)
}

// local client on instance
var stockClient *Client = nil

// New can create a StockDB Client
func NewClient(uri, auth string) (client *Client) {
	if stockClient == nil {
		stockClient = new(Client)
		stockClient.init(uri, auth)
	}
	return stockClient
}
