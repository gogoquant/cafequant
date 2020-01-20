package sdk

import (
	"encoding/base64"
	"github.com/hprose/hprose-golang/io"
	"github.com/hprose/hprose-golang/rpc"
	"github.com/xiyanxiyan10/stockdb/types"
	"net/http"
)

func init() {
	io.Register(types.Option{}, "Option", "json")
	io.Register(types.OHLC{}, "OHLC", "json")
	io.Register(types.Order{}, "Order", "json")
	io.Register(types.OrderBook{}, "OrderBook", "json")
	io.Register(types.Depth{}, "Depth", "json")
	io.Register(types.BaseResponse{}, "BaseResponse", "json")
	io.Register(types.Stats{}, "Stats", "json")
	io.Register(types.StatsResponse{}, "StatsResponse", "json")
	io.Register(types.StringsResponse{}, "StringsResponse", "json")
	io.Register(types.TimeRangeResponse{}, "TimeRangeResponse", "json")
	io.Register(types.OHLCResponse{}, "OHLCResponse", "json")
	io.Register(types.DepthResponse{}, "DepthResponse", "json")
}

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
		c.hprose.Header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(auth)))
	}
	c.hprose.UseService(c)
}

// New can create a StockDB Client
func NewClient(uri, auth string) (client *Client) {
	var stockClient Client
	stockClient.init(uri, auth)
	return &stockClient
}
