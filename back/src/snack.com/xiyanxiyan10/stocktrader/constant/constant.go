package constant

const (
	// CacheTicker ...
	CacheTicker = "ticker"

	// CacheDepth ...
	CacheDepth = "depth"

	// CacheTrader ...
	CacheTrader = "trader"

	// CacheKLine ...
	CacheKLine = "kline"

	// CacheOrder ...
	CacheOrder = "order"

	// CacheRecord ...
	CacheRecord = "record"

	// CachePosition ...
	CachePosition = "position"
)

const (
	// IONONE ...
	IONONE = "online"
	// IOBLOCK get from sync
	IOBLOCK = "block"
	// IOCACHE get from cache
	IOCACHE = "cache"
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
	Version                    = "1.0.0"
	ErrAuthorizationError      = "Authorization Error"
	ErrInsufficientPermissions = "Insufficient Permissions"
)

// exchange types
const (
	HuoBiDm    = "HuoBiDm"
	HuoBi      = "HuoBi"
	SZ         = "SZ"
	SpotBack   = "SpotBack"
	FutureBack = "FutureBack"
)

// log types
const (
	ERROR  = "ERROR"
	INFO   = "INFO"
	PROFIT = "PROFIT"
)

const (
	// STOCKDBURL ...
	STOCKDBURL = "stockdburl"
	// STOCKDBAUTH ...
	STOCKDBAUTH = "stockdbpwd"
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
	ExchangeTypes = []string{HuoBiDm, FutureBack, HuoBi, SpotBack, SZ}
	ScriptTypes   = []string{ScriptJs, ScriptGo}
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
	DepthSize = 10
	// GoHandler ...
	GoHandler = "NewHandler"
	// DefaultTimeOut ...
	DefaultTimeOut = 2
)

const (
	// ScriptJs ...
	ScriptJs = "js"
	// ScriptGo ... @todo change as go
	ScriptGo = "go"
)