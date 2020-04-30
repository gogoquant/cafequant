package api

//base exchange
type BaseExchange struct {
	IOMode           int                // io mode for exchange
	contractType     string             // contractType
	direction        string             // trade type
	stockType        string             // stockType
	lever            int                // lever
	recordsPeriodMap map[string]int64   // recordsPeriod support
	minAmountMap     map[string]float64 // minAmount of trade
}

// SetIO set IO mode
func (e *BaseExchange) SetIO(mode int) {
	e.IOMode = mode
}

// GetIO get IO mode
func (e *BaseExchange) GetIO() int {
	return e.IOMode
}

// SetContractType set the limit calls amount per second of this exchange
func (e *BaseExchange) SetContractType(contractType string) {
	e.contractType = contractType
}

// GetContractType set the limit calls amount per second of this exchange
func (e *BaseExchange) GetContractType() string {
	return e.contractType
}

// SetContractType set the limit calls amount per second of this exchange
func (e *BaseExchange) SetDirection(direction string) {
	e.direction = direction
}

// SetContractType set the limit calls amount per second of this exchange
func (e *BaseExchange) GetDirection() string {
	return e.direction
}

// SetContractType set the limit calls amount per second of this exchange
func (e *BaseExchange) SetMarginLevel(lever int) {
	e.lever = lever
}

// SetContractType set the limit calls amount per second of this exchange
func (e *BaseExchange) GetMarginLevel() int {
	return e.lever
}

// SetContractType set the limit calls amount per second of this exchange
func (e *BaseExchange) GetStockType() string {
	return e.stockType
}

// SetContractType set the limit calls amount per second of this exchange
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
