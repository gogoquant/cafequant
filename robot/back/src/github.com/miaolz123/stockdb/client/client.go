package client

import (
	"encoding/base64"
	"net/http"

	"github.com/miaolz123/stockdb/types"
	"github.com/hprose/hprose-golang/io"
	"github.com/hprose/hprose-golang/rpc"
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

// New can create a StockDB Client
func New(uri, auth string) (client *Client) {
	client = &Client{
		uri:    uri,
		Hprose: rpc.NewHTTPClient(uri),
	}
	if auth != "" {
		client.Hprose.Header = make(http.Header)
		client.Hprose.Header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(auth)))
	}
	client.Hprose.UseService(&client)
	return
}
