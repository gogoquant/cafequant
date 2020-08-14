package api

import (
	"fmt"
	"time"

	dbsdk "snack.com/xiyanxiyan10/stockdb/sdk"
	dbtypes "snack.com/xiyanxiyan10/stockdb/types"
	"snack.com/xiyanxiyan10/stocktrader/constant"
	"snack.com/xiyanxiyan10/stocktrader/model"
	"snack.com/xiyanxiyan10/stocktrader/util"
)

// BaseExchange ...
type BaseExchange struct {
	BaseExchangeCachePool // cache for exchange
	id                    int
	ioMode                int                // io mode for exchange
	contractType          string             // contractType
	direction             string             // trade type
	stockType             string             // stockType
	lever                 float64            // lever
	recordsPeriodMap      map[string]int64   // recordsPeriod support
	minAmountMap          map[string]float64 // minAmount of trade
	limit                 float64
	lastSleep             int64
	lastTimes             int64
	subscribeMap          map[string][]string
	host                  string
	logger                model.Logger
	option                constant.Option
}

// Subscribe ...
func (e *BaseExchange) Subscribe(source, action string) interface{} {
	e.subscribeMap[source] = append(e.subscribeMap[source], action)
	return "success"
}

// SetLimit set the limit calls amount per second of this exchange
func (e *BaseExchange) SetLimit(times interface{}) float64 {
	e.limit = util.Float64Must(times)
	return e.limit
}

// AutoSleep auto sleep to achieve the limit calls amount per second of this exchange
func (e *BaseExchange) AutoSleep() {
	now := time.Now().UnixNano()
	interval := 1e+9/e.limit*util.Float64Must(e.lastTimes) - util.Float64Must(now-e.lastSleep)
	if interval > 0.0 {
		time.Sleep(time.Duration(util.Int64Must(interval)))
	}
	e.lastTimes = 0
	e.lastSleep = now
}

// BackGetStats ...
func (e *BaseExchange) BackGetStats() interface{} {
	if !e.option.BackTest {
		e.logger.Log(constant.ERROR, e.GetStockType(), 0.0, 0.0, fmt.Sprint("GetStats error, the error number not in backtest"))
		return nil
	}
	client := dbsdk.NewClient(constant.STOCKDBURL, constant.STOCKDBAUTH)
	ohlc := client.GetStats()
	if !ohlc.Success {
		e.logger.Log(constant.ERROR, e.GetStockType(), 0.0, 0.0, fmt.Sprint("GetStats error, the error number is %s"+ohlc.Message))
	}
	return ohlc.Data
}

// BackGetMarkets ...
func (e *BaseExchange) BackGetMarkets() interface{} {
	if !e.option.BackTest {
		e.logger.Log(constant.ERROR, e.GetStockType(), 0.0, 0.0, fmt.Sprint("GetMarkets error, the error number not in backtest"))
		return nil
	}
	client := dbsdk.NewClient(constant.STOCKDBURL, constant.STOCKDBAUTH)
	ohlc := client.GetStats()
	if !ohlc.Success {
		e.logger.Log(constant.ERROR, e.GetStockType(), 0.0, 0.0, fmt.Sprint("GetMarkets error, the error number is %s"+ohlc.Message))
	}
	return ohlc.Data
}

// BackGetSymbols  ...
func (e *BaseExchange) BackGetSymbols(market string) interface{} {
	if !e.option.BackTest {
		e.logger.Log(constant.ERROR, e.GetStockType(), 0.0, 0.0, fmt.Sprint("GetSymbols error, the error number not in backtest"))
		return nil
	}
	client := dbsdk.NewClient(constant.STOCKDBURL, constant.STOCKDBAUTH)
	ohlc := client.GetStats()
	if !ohlc.Success {
		e.logger.Log(constant.ERROR, e.GetStockType(), 0.0, 0.0, fmt.Sprint("GetSymbols error, the error number is %s"+ohlc.Message))
	}
	return ohlc.Data
}

// BackGetOHLCs ...
func (e *BaseExchange) BackGetOHLCs(begin, end, period int64) interface{} {
	var opt dbtypes.Option
	if !e.option.BackTest {
		e.logger.Log(constant.ERROR, e.GetStockType(), 0.0, 0.0, fmt.Sprint("GetOHLCs error, the error number not in backtest"))
		return nil
	}
	opt.Market = e.option.Type
	opt.Symbol = e.GetStockType()
	opt.Period = period
	opt.BeginTime = begin
	opt.EndTime = end
	client := dbsdk.NewClient(constant.STOCKDBURL, constant.STOCKDBAUTH)
	ohlc := client.GetOHLCs(opt)
	if !ohlc.Success {
		e.logger.Log(constant.ERROR, e.GetStockType(), 0.0, 0.0, fmt.Sprint("GetOHLCs error, the error number is %s"+ohlc.Message))
	}
	return ohlc.Data
}

// BackGetTimeRange ...
func (e *BaseExchange) BackGetTimeRange() interface{} {
	var opt dbtypes.Option
	if !e.option.BackTest {
		e.logger.Log(constant.ERROR, e.GetStockType(), 0.0, 0.0, fmt.Sprint("GetTimeRange error, the error number not in backtest"))
		return nil
	}
	opt.Market = e.option.Type
	opt.Symbol = e.GetStockType()
	client := dbsdk.NewClient(constant.STOCKDBURL, constant.STOCKDBAUTH)
	timeRange := client.GetTimeRange(opt)
	if !timeRange.Success {
		e.logger.Log(constant.ERROR, e.GetStockType(), 0.0, 0.0, fmt.Sprint("GetTimeRange, the error number is %s"+timeRange.Message))
	}
	return timeRange.Data
}

// BackGetPeriodRange ...
func (e *BaseExchange) BackGetPeriodRange() interface{} {
	var opt dbtypes.Option
	if !e.option.BackTest {
		e.logger.Log(constant.ERROR, e.GetStockType(), 0.0, 0.0, fmt.Sprint("GetPeriodRange error, the error number not in backtest"))
		return nil
	}
	opt.Market = e.option.Type
	opt.Symbol = e.GetStockType()
	client := dbsdk.NewClient(constant.STOCKDBURL, constant.STOCKDBAUTH)
	timeRange := client.GetPeriodRange(opt)
	if !timeRange.Success {
		e.logger.Log(constant.ERROR, e.GetStockType(), 0.0, 0.0, fmt.Sprint("GetPeriodRange, the error number is %s"+timeRange.Message))
	}
	return timeRange.Data
}

// BackGetDepth ...
func (e *BaseExchange) BackGetDepth(begin, end, period int64) interface{} {
	var opt dbtypes.Option
	if !e.option.BackTest {
		e.logger.Log(constant.ERROR, e.GetStockType(), 0.0, 0.0, fmt.Sprint(" GetDepth error, the error number not in backtest"))
		return nil
	}
	opt.Market = e.option.Type
	opt.Symbol = e.GetStockType()
	opt.Period = period
	opt.BeginTime = begin
	opt.EndTime = end
	client := dbsdk.NewClient(constant.STOCKDBURL, constant.STOCKDBAUTH)
	depth := client.GetDepth(opt)
	if !depth.Success {
		e.logger.Log(constant.ERROR, e.GetStockType(), 0.0, 0.0, fmt.Sprint("GetDepth error, the error number is %s"+depth.Message))
	}
	return depth.Data
}

// Init ...
func (e *BaseExchange) Init(opt constant.Option) error {
	e.logger = model.Logger{TraderID: opt.TraderID, ExchangeType: opt.Type}
	e.option = opt
	e.limit = opt.Limit
	e.lastSleep = time.Now().UnixNano()
	return nil
}

// SetID set ID
func (e *BaseExchange) SetID(id int) {
	e.id = id
}

// GetID get ID
func (e *BaseExchange) GetID() int {
	return e.id
}

// SetIO set IO mode
func (e *BaseExchange) SetIO(mode int) {
	e.ioMode = mode
}

// GetIO get IO mode
func (e *BaseExchange) GetIO() int {
	return e.ioMode
}

// SetContractType set the limit calls amount per second of this exchange
func (e *BaseExchange) SetContractType(contractType string) {
	e.contractType = contractType
}

// GetContractType set the limit calls amount per second of this exchange
func (e *BaseExchange) GetContractType() string {
	return e.contractType
}

// SetDirection set the limit calls amount per second of this exchange
func (e *BaseExchange) SetDirection(direction string) {
	e.direction = direction
}

// GetDirection set the limit calls amount per second of this exchange
func (e *BaseExchange) GetDirection() string {
	return e.direction
}

// SetMarginLevel set the limit calls amount per second of this exchange
func (e *BaseExchange) SetMarginLevel(lever float64) {
	e.lever = lever
}

// GetMarginLevel set the limit calls amount per second of this exchange
func (e *BaseExchange) GetMarginLevel() float64 {
	return e.lever
}

// GetStockType set the limit calls amount per second of this exchange
func (e *BaseExchange) GetStockType() string {
	return e.stockType
}

// SetStockType set the limit calls amount per second of this exchange
func (e *BaseExchange) SetStockType(stockType string) {
	e.stockType = stockType
}

// SetMinAmountMap ...
func (e *BaseExchange) SetMinAmountMap(m map[string]float64) {
	e.minAmountMap = m
}

// GetMinAmountMap ...
func (e *BaseExchange) GetMinAmountMap() map[string]float64 {
	return e.minAmountMap
}

// SetRecordsPeriodMap ...
func (e *BaseExchange) SetRecordsPeriodMap(m map[string]int64) {
	e.recordsPeriodMap = m
}

// GetRecordsPeriodMap ...
func (e *BaseExchange) GetRecordsPeriodMap() map[string]int64 {
	return e.recordsPeriodMap
}
