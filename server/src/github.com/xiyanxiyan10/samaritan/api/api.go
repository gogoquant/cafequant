package api

// Option is an exchange option
type Option struct {
	TraderID  int64
	Type      string
	Name      string
	AccessKey string
	SecretKey string
	// Ctx       *otto.Otto
}

// Exchange interface
type Exchange interface {
	Log(...interface{})
	GetType() string
	GetName() string
	SetLimit(times interface{}) float64
	AutoSleep()
	GetMinAmount(stock string) float64
	GetAccount() interface{}
	Trade(tradeType string, stockType string, price, amount interface{}, msgs ...interface{}) interface{}
	GetOrder(stockType, id string) interface{}
	GetOrders(stockType string) interface{}
	GetTrades(stockType string) interface{}
	CancelOrder(order Order) bool
	GetTicker(stockType string, sizes ...interface{}) interface{}
	GetRecords(stockType, period string, sizes ...interface{}) interface{}
}

//Entity transaction
type Etrade interface {
	Log(...interface{})
	GetList(name string, company string, minprice float64, maxprice float64) interface{}
}

var (
	constructor = map[string]func(Option) Exchange{}
)
