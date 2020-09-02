package constant

const (
	// CacheTicker ...
	CacheTicker = "ticker"

	// CacheDepth ...
	CacheDepth = "depth"

	// CacheTrader ...
	CacheTrader = "trader"

	// CacheKline ...
	CacheKline = "kline"

	// CacheOrder ...
	CacheOrder = "order"

	// CachePosition ...
	CachePosition = "position"
)

const (
	// IONONE ...
	IONONE = iota
	// IOONLINE  get from online
	IOONLINE
	// IOBLOCK get from sync
	IOBLOCK
	// IOCACHE get from cache
	IOCACHE
)

const (
	// RUNNORMAIL ..
	RUNNORMAIL = iota
	// RUNBACK ...
	RUNBACK
)

// error constants
const (
	Banner                     = "QuantBot"
	Version                    = "0.0.3"
	ErrAuthorizationError      = "Authorization Error"
	ErrInsufficientPermissions = "Insufficient Permissions"
)

// exchange types
const (
	HuoBiDm     = "HuoBiDm"
	HuoBiDmBack = "HuoBiDmBack"
	HuoBi       = "HuoBi"
	HuoBiBack   = "HuoBiBack"
)

// log types
const (
	ERROR  = "ERROR"
	INFO   = "INFO"
	PROFIT = "PROFIT"
)

const (
	// STOCKDBURL ...
	STOCKDBURL = ""
	// STOCKDBAUTH ...
	STOCKDBAUTH = ""
)

// trade types
const (
	TradeTypeBuy        = "buy"
	TradeTypeSell       = "sell"
	TradeTypeLong       = "buy"
	TradeTypeShort      = "sell"
	TradeTypeLongClose  = "closebuy"
	TradeTypeShortClose = "closesell"
	TradeTypeCancel     = "cancel"
	TradeTypeHold       = "hold"
)

// some variables
var (
	ExchangeTypes = []string{HuoBiDm, HuoBiDmBack, HuoBi, HuoBiBack}
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

const (
	// FilePath ...
	FilePath = "files"
	// GoPluginPath ...
	GoPluginPath = "goplugin"
	// KLineSize ...
	KLineSize = 100
)

const (
	// RecordSize ...
	RecordSize = 100
	// DepthSize ...
	DepthSize = 100
	// GoHandler ...
	GoHandler = "NewHandler"
)
