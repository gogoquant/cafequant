module snack.com/xiyanxiyan10/stocktrader

go 1.15

require (
	github.com/axgle/mahonia v0.0.0-20180208002826-3358181d7394
	github.com/bitly/go-simplejson v0.5.0
	github.com/blinkbean/dingtalk v0.0.0-20200822153748-8cf931f926ab
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/go-echarts/go-echarts v1.0.0
	github.com/go-echarts/go-echarts/v2 v2.2.3
	github.com/go-gomail/gomail v0.0.0-20160411212932-81ebce5c23df
	github.com/go-ini/ini v1.55.0
	github.com/go-python/gopy v0.3.4
	github.com/go-sql-driver/mysql v1.5.0 // indirect
	github.com/gofinance/ib v0.0.0-20190131202149-a7abd0c5d772
	github.com/gorilla/websocket v1.4.2 // indirect
	github.com/hprose/hprose-golang v2.0.5+incompatible
	github.com/jinzhu/gorm v1.9.12
	github.com/json-iterator/go v1.1.10 // indirect
	github.com/kirinlabs/HttpRequest v1.0.5
	github.com/klauspost/compress v1.10.3 // indirect
	github.com/lib/pq v1.2.0 // indirect
	github.com/mayiweb/goctp v0.0.0-20190917081845-ed4a7d3f7e3e
	github.com/nntaoli-project/goex v1.2.5
	github.com/pkg/errors v0.9.1 // indirect
	github.com/robertkrimen/otto v0.0.0-20191219234010-c382bd3c16ff
	github.com/senseyeio/roger v0.0.0-20191009211040-43e330bee47f
	github.com/sirupsen/logrus v1.6.0 // indirect
	github.com/stretchr/testify v1.6.1 // indirect
	golang.org/x/net v0.0.0-20201021035429-f5854403a974 // indirect
	gopkg.in/alexcesaro/quotedprintable.v3 v3.0.0-20150716171945-2caba252f4dc // indirect
	gopkg.in/check.v1 v1.0.0-20190902080502-41f04d3bba15 // indirect
	gopkg.in/gomail.v2 v2.0.0-20160411212932-81ebce5c23df // indirect
	gopkg.in/logger.v1 v1.0.1
	gopkg.in/sourcemap.v1 v1.0.5 // indirect
	snack.com/xiyanxiyan10/conver v0.0.0
	snack.com/xiyanxiyan10/stockdb v0.0.0
)

replace (
	github.com/qiniu/py v1.2.2 => github.com/xiyanxiyan10/py v0.0.0-20200907052829-5727a2a1895d
	snack.com/xiyanxiyan10/conver v0.0.0 => ../conver
	snack.com/xiyanxiyan10/stockdb => ../stockdb
)
