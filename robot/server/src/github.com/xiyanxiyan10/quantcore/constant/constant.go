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
	HuobiDm = "huobidm"
	Fmex    = "fmex"
)

// log types
const (
	ERROR = "ERROR"
	INFO  = "INFO"
)

// trade types
const (
	TradeTypeBuy         = "buy"
	TradeTypeSell        = "sell"
	TradeTypeLong        = "buy"
	TradeTypeShort       = "sell"
	TradeTypeLongClose   = "closebuy"
	TradeTypeShortClose  = "closesell"
	TradeTypeCancelOrder = "cancel"
)

// some variables
var (
	Consts        = []string{"M", "M5", "M15", "M30", "H", "D", "W"}
	ExchangeTypes = []string{HuobiDm, Fmex}
)

// future userinfo string
const (
	Currency      = "Currency"
	AccountRights = "AccountRights" //账户权益
	KeepDeposit   = "KeepDeposit"   //保证金
	ProfitReal    = "ProfitReal"    //已实现盈亏
	ProfitUnreal  = "ProfitUnreal"
	RiskRate      = "RiskRate" //保证金率
)
