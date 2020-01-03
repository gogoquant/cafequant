package constant

// error constants
const (
	Banner                     = "QuantBot"
	Version                    = "0.0.3"
	ErrAuthorizationError      = "Authorization Error"
	ErrInsufficientPermissions = "Insufficient Permissions"
)

// exchange types
const (
	HuoBiDm = "HuoBiDm"
	Fmex    = "Fmex"
	HuoBi   = "HuoBi"
)

// log types
const (
	ERROR  = "ERROR"
	INFO   = "INFO"
	PROFIT = "PROFIT"
)

// trade types
const (
	TradeTypeBuy        = "buy"
	TradeTypeSell       = "sell"
	TradeTypeLong       = "buy"
	TradeTypeShort      = "sell"
	TradeTypeLongClose  = "closeBuy"
	TradeTypeShortClose = "closeSell"
	TradeTypeCancel     = "cancel"
	TradeTypeHold       = "hold"
)

// some variables
var (
	ExchangeTypes = []string{HuoBiDm, Fmex, HuoBi}
)

// future userInfo string
const (
	Currency      = "Currency"
	AccountRights = "AccountRights" //账户权益
	KeepDeposit   = "KeepDeposit"   //保证金
	ProfitReal    = "ProfitReal"    //已实现盈亏
	ProfitUnreal  = "ProfitUnreal"
	RiskRate      = "RiskRate" //保证金率
)
