package sdkctp

import (
	"fmt"
	"github.com/mayiweb/goctp"
	"log"
	"os"
	"sync"
)

// CtpMaster
type CtpMaster struct {
	// 行情模块 api
	MdApi goctp.CThostFtdcMdApi

	// 交易模块 api
	TraderApi goctp.CThostFtdcTraderApi

	// ctp 句柄及配置项
	Client *CtpClient

	// 行情模块函数 句柄
	MdSpi *FtdcMdSpi

	// 交易模块函数 句柄
	TraderSpi *FtdcTraderSpi

	// 交易所合约详情列表 InstrumentInfoStruct
	MapInstrumentInfos sync.Map

	// 报单列表（已成交、未成交、撤单等状态）的列表数据 OrderListStruct
	MapOrderList sync.Map

	// ctp 服务器，及交易账号
	MdFront     []string
	TraderFront []string

	BrokerID   string
	InvestorID string
	Password   string

	// 客户端认证
	AppID    string
	AuthCode string

	// ctp 流文件，绝对路径
	// StreamFile string = GetCurrentDirectory() + "/StreamFile/"
	StreamFile string

	// 买卖方向：买
	//OrderBuy byte  = '0'
	OrderBuy byte

	// 买卖方向：卖
	//OrderSell byte = '1'
	OrderSell byte

	// 运行模式（prod 生产，test 标准环境测试，dev 24小时测试）
	RunMode string
}

/**
 * 计算盈亏
 *
 * @param   InstrumentID  string  合约
 * @param   OpenPrice     float64 开仓价格
 * @param   LastPrice     float64 最新价|平仓价格
 * @param   Number        int     数量
 * @param   PosiDirection string  持仓方向[2：买，3：卖]
 */
func (ctp *CtpMaster) GetPositionProfit(InstrumentID string, OpenPrice float64, LastPrice float64, Number int, PosiDirection string) float64 {

	InstrumentInfo, _ := ctp.GetInstrumentInfo(InstrumentID)
	if PosiDirection == "2" {
		return ((LastPrice - OpenPrice) * float64(InstrumentInfo.VolumeMultiple)) * float64(Number)
	} else {
		return ((OpenPrice - LastPrice) * float64(InstrumentInfo.VolumeMultiple)) * float64(Number)
	}
}

// 获得持仓结构体数据
func (ctp *CtpMaster) GetInvestorPositionStruct(pInvestorPosition goctp.CThostFtdcInvestorPositionField) InvestorPositionStruct {

	var mInvestorPosition InvestorPositionStruct

	// 检查合约详情是否存在
	mInstrumentInfo, mapRes := ctp.GetInstrumentInfo(pInvestorPosition.GetInstrumentID())
	if !mapRes {
		fmt.Printf("合约 %v 不存在！\n", pInvestorPosition.GetInstrumentID())
		return mInvestorPosition
	}

	// 合约乘数
	var VolumeMultiple int = mInstrumentInfo.VolumeMultiple

	// 开仓成本
	var OpenCost float64 = pInvestorPosition.GetOpenCost() / float64(pInvestorPosition.GetPosition()*VolumeMultiple)

	mInvestorPosition.BrokerID = pInvestorPosition.GetBrokerID()
	mInvestorPosition.InvestorID = pInvestorPosition.GetInvestorID()
	mInvestorPosition.InstrumentID = pInvestorPosition.GetInstrumentID()
	mInvestorPosition.InstrumentName = mInstrumentInfo.InstrumentName
	mInvestorPosition.PosiDirection = string(pInvestorPosition.GetPosiDirection())
	mInvestorPosition.PosiDirectionTitle = GetPosiDirectionTitle(mInvestorPosition.PosiDirection)
	mInvestorPosition.HedgeFlag = string(pInvestorPosition.GetHedgeFlag())
	mInvestorPosition.HedgeFlagTitle = GetHedgeFlagTitle(mInvestorPosition.HedgeFlag)
	mInvestorPosition.PositionDate = string(pInvestorPosition.GetPositionDate())
	mInvestorPosition.PositionDateTitle = GetPositionDateTitle(mInvestorPosition.PositionDate)
	mInvestorPosition.Position = pInvestorPosition.GetPosition()
	mInvestorPosition.YdPosition = pInvestorPosition.GetYdPosition()
	mInvestorPosition.TodayPosition = pInvestorPosition.GetTodayPosition()
	mInvestorPosition.LongFrozen = pInvestorPosition.GetLongFrozen()
	mInvestorPosition.ShortFrozen = pInvestorPosition.GetShortFrozen()

	// 冻结的持仓量（多空并成一个字段）
	if mInvestorPosition.PosiDirection == string(goctp.THOST_FTDC_PD_Long) {
		// 多头冻结的持仓量
		mInvestorPosition.ShortVolume = pInvestorPosition.GetShortFrozen()
	} else {
		// 空头冻结的持仓量
		mInvestorPosition.ShortVolume = pInvestorPosition.GetLongFrozen()
	}

	mInvestorPosition.OpenVolume = pInvestorPosition.GetOpenVolume()
	mInvestorPosition.CloseVolume = pInvestorPosition.GetCloseVolume()
	mInvestorPosition.PositionCost = Decimal(pInvestorPosition.GetPositionCost(), 2)
	mInvestorPosition.Commission = Decimal(pInvestorPosition.GetCommission(), 2)
	mInvestorPosition.CloseProfit = pInvestorPosition.GetCloseProfit()
	mInvestorPosition.PositionProfit = Decimal(pInvestorPosition.GetPositionProfit(), 2)
	mInvestorPosition.PreSettlementPrice = pInvestorPosition.GetPreSettlementPrice()
	mInvestorPosition.SettlementPrice = Decimal(pInvestorPosition.GetSettlementPrice(), 2)
	mInvestorPosition.SettlementID = pInvestorPosition.GetSettlementID()
	mInvestorPosition.OpenCost = Decimal(OpenCost, 2)
	mInvestorPosition.ExchangeID = pInvestorPosition.GetExchangeID()

	return mInvestorPosition
}

// 获得合约详情信息
func (ctp *CtpMaster) GetInstrumentInfo(InstrumentID string) (InstrumentInfoStruct, bool) {
	if v, ok := ctp.MapInstrumentInfos.Load(InstrumentID); ok {
		return v.(InstrumentInfoStruct), true
	} else {
		var mInstrumentInfo InstrumentInfoStruct
		return mInstrumentInfo, false
	}
}

// Ctp 行情 spi 回调函数
type FtdcMdSpi struct {
	Master *CtpMaster
}

// Ctp 交易 spi 回调函数
type FtdcTraderSpi struct {
	Master *CtpMaster
}

// Ctp 客户端 行情、交易模块 全局变量
type CtpClient struct {
	BrokerID   string
	InvestorID string
	Password   string

	// 客户端认证
	AppID    string
	AuthCode string

	// 当前交易日期
	TradingDay string

	// 当前交易月份
	TradeMonth string

	// 行情请求编号
	MdRequestId int

	// 交易请求编号
	TraderRequestId int

	// 交易系统是否已经初始化了
	IsTraderInit bool

	// 交易程序是否初始化完成（自动完成如下动作：交易账号登陆、结算单确认、查询合约、查询资金账户、查询用户报单、查询用户持仓 后算完成）
	IsTraderInitFinish bool

	// 交易程序是否已登录过
	IsTraderLogin bool

	// 行情程序是否已登录过
	IsMdLogin bool
}

// 获得交易请求编号
func (client *CtpClient) GetTraderRequestId() int {
	client.TraderRequestId += 1
	return client.TraderRequestId
}

// SetTradeAccount ...
func (ctp *CtpMaster) SetTradeAccount(MdFront, TraderFront []string, BrokerID, InvestorID, Password, AppID, AuthCode string) {
	ctp.MdFront = append(ctp.MdFront, MdFront...)
	ctp.TraderFront = append(ctp.TraderFront, TraderFront...)
	ctp.BrokerID = BrokerID
	ctp.InvestorID = InvestorID
	ctp.Password = Password
	ctp.AppID = AppID
	ctp.AuthCode = AuthCode
	ctp.OrderBuy = '0'
	ctp.OrderSell = '1'
	//ctp.MdSpi = FtdcMdSpi{}
	//ctp.TraderSpi = FtdcTraderSpi{}
}

func (ctp *CtpMaster) Start() error {
	log.Println("启动交易程序")

	// 检查流文件目录是否存在
	fileExists, _ := PathExists(ctp.StreamFile)
	if !fileExists {
		err := os.Mkdir(ctp.StreamFile, os.ModePerm)
		if err != nil {
			fmt.Println("创建目录失败，请检查是否有操作权限")
		}
	}
	ctp.MdApi = goctp.CThostFtdcMdApiCreateFtdcMdApi(ctp.StreamFile)
	ctp.TraderApi = goctp.CThostFtdcTraderApiCreateFtdcTraderApi(ctp.StreamFile)

	client := &CtpClient{
		//MdApi:              goctp.CThostFtdcMdApiCreateFtdcMdApi(ctp.StreamFile),
		//TraderApi:          goctp.CThostFtdcTraderApiCreateFtdcTraderApi(ctp.StreamFile),
		BrokerID:           ctp.BrokerID,
		InvestorID:         ctp.InvestorID,
		Password:           ctp.Password,
		AppID:              ctp.AppID,
		AuthCode:           ctp.AuthCode,
		MdRequestId:        1,
		TraderRequestId:    1,
		IsTraderInit:       false,
		IsTraderInitFinish: false,
		IsMdLogin:          false,
		IsTraderLogin:      false,
	}
	ctp.Client = client
	ctp.MdSpi = &FtdcMdSpi{Master: ctp}
	ctp.MdApi.RegisterSpi(goctp.NewDirectorCThostFtdcMdSpi(ctp.MdSpi))
	for _, val := range ctp.MdFront {
		ctp.MdApi.RegisterFront(val)
	}
	ctp.MdApi.Init()

	ctp.TraderSpi = &FtdcTraderSpi{Master: ctp}
	ctp.TraderApi.RegisterSpi(goctp.NewDirectorCThostFtdcTraderSpi(ctp.TraderSpi))

	for _, val := range ctp.TraderFront {
		ctp.TraderApi.RegisterFront(val)
	}

	ctp.TraderApi.SubscribePublicTopic(goctp.THOST_TERT_QUICK)
	ctp.TraderApi.SubscribePrivateTopic(goctp.THOST_TERT_QUICK)

	ctp.TraderApi.Init()
	ctp.TraderApi.Join()
	return nil
	// .Join() 如果后面有其它需要处理的功能可以不写，但必须保证程序不能退出，Join 就是保证程序不退出的
}

// 设置交易账号
/*
func SetTradeAccount() {

	switch RunMode {

	// 迈科期货
	case "prod":
		MdFront = []string{}
		TraderFront = []string{}
		BrokerID = ""
		InvestorID = ""
		Password = ""
		AppID = ""
		AuthCode = ""

	// 与实际生产环境保持一致
	case "test":
		MdFront = []string{"tcp://180.168.146.187:10110", "tcp://180.168.146.187:10111", "tcp://218.202.237.33:10112"}
		TraderFront = []string{"tcp://180.168.146.187:10100", "tcp://180.168.146.187:10101", "tcp://218.202.237.33:10102"}
		BrokerID = "9999"
		InvestorID = ""
		Password = ""
		AppID = ""
		AuthCode = ""

	// 7*24 服务器，交易日，16：00～次日09：00；非交易日，16：00～次日15：00
	case "dev":
		MdFront = []string{"tcp://180.168.146.187:10131"}
		TraderFront = []string{"tcp://180.168.146.187:10130"}
		BrokerID = "9999"
		InvestorID = ""
		Password = ""
		AppID = ""
		AuthCode = ""

	default:
		Println("该模式未设置交易账号信息")
		os.Exit(1)
	}
}
*/

/*
func init() {
	// 全局 行情、交易 函数句柄
	MdSpi = FtdcMdSpi{}
	TraderSpi = FtdcTraderSpi{}
}
*/
/*
func main() {

	// 运行模式【运行程序时带上参数可设置】
	if len(os.Args) < 2 {
		RunMode = "test"
	} else {
		RunMode = os.Args[1]
	}

	// 设置交易账号
	SetTradeAccount()

	log.Println("启动交易程序")

	// 检查流文件目录是否存在
	fileExists, _ := PathExists(StreamFile)
	if !fileExists {
		err := os.Mkdir(StreamFile, os.ModePerm)
		if err != nil {
			fmt.Println("创建目录失败，请检查是否有操作权限")
		}
	}

	Ctp = CtpClient{
		MdApi:              goctp.CThostFtdcMdApiCreateFtdcMdApi(StreamFile),
		TraderApi:          goctp.CThostFtdcTraderApiCreateFtdcTraderApi(StreamFile),
		BrokerID:           BrokerID,
		InvestorID:         InvestorID,
		Password:           Password,
		AppID:              AppID,
		AuthCode:           AuthCode,
		MdRequestId:        1,
		TraderRequestId:    1,
		IsTraderInit:       false,
		IsTraderInitFinish: false,
		IsMdLogin:          false,
		IsTraderLogin:      false,
	}

	Ctp.MdApi.RegisterSpi(goctp.NewDirectorCThostFtdcMdSpi(&FtdcMdSpi{Client: Ctp}))

	for _, val := range MdFront {
		Ctp.MdApi.RegisterFront(val)
	}
	Ctp.MdApi.Init()

	Ctp.TraderApi.RegisterSpi(goctp.NewDirectorCThostFtdcTraderSpi(&FtdcTraderSpi{Client: Ctp}))

	for _, val := range TraderFront {
		Ctp.TraderApi.RegisterFront(val)
	}

	Ctp.TraderApi.SubscribePublicTopic(goctp.THOST_TERT_QUICK)
	Ctp.TraderApi.SubscribePrivateTopic(goctp.THOST_TERT_QUICK)
	Ctp.TraderApi.Init()
	Ctp.TraderApi.Join()

	// .Join() 如果后面有其它需要处理的功能可以不写，但必须保证程序不能退出，Join 就是保证程序不退出的
}
*/

type CtpHandler interface {
	SetTradeAccount(MdFront, TraderFront []string, BrokerID, InvestorID, Password, AppID, AuthCode string)
	Start() error
}

// NewCtp ...
func NewCtp() CtpHandler {
	return &CtpMaster{}
}
