package api

//base exchange
type BaseExchange struct {
	contractType string
	direction    string
	stockType    string
	lever        int
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
