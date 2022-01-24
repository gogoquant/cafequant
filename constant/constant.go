package constant

const (
	StepLine   = "stepline"
	SmoothLine = "smoothline"
	BrokeLine  = "brokeline"
	AreaLine   = "arealine"
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
	SpotBack   = "SpotaBack"
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
	ScriptTypes   = []string{ScriptJs}
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
	ORDER_UNFINISH    = 0
	ORDER_PART_FINISH = 1
	ORDER_FINISH      = 2
	ORDER_CANCEL      = 3
	ORDER_REJECT      = 4
	ORDER_CANCEL_ING  = 5
	ORDER_FAIL        = 6
)

// Period constant
const (
	Second    int64 = 1
	Minute    int64 = 60 * Second
	Hour      int64 = 60 * Minute
	Day       int64 = 24 * Hour
	Week      int64 = 7 * Day
	Month     int64 = 30 * Day
	Quarter   int64 = 3 * Month
	Year      int64 = 365 * Day
	MinPeriod int64 = 3
)

const (
	// FilePath ...
	FilePath = "files"
	// GoPluginPath ...
	GoPluginPath = "goplugin"
	// KLineSize ...
	KLineSize = 100
	//CacheSize ...
	CacheSize = 10
	// RecordSize ...
	RecordSize = 1000
	// DepthSize ...
	DepthSize = 10
	// GoHandler ...
	GoHandler = "GoHandler"
	// DefaultTimeOut ...
	DefaultTimeOut = 2
	// ScriptJs ...
	ScriptJs = "js"
	// ScriptGo ... @todo change as go
	//ScriptGo = "go"
	// Pending ...
	Pending = -1
	// Running ...
	Running = 1
	// Stop ...
	Stop = 0
	// Enable ...
	Enable = 1
	// Disable ...
	Disable = 0
	// BackEnd ...
	BackEnd = "backend"
)
