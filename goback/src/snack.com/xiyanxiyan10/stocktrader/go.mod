module snack.com/xiyanxiyan10/stocktrader

go 1.13

require (
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/go-echarts/go-echarts v0.0.0-20190915064101-cbb3b43ade5d
	github.com/go-gomail/gomail v0.0.0-20160411212932-81ebce5c23df
	github.com/go-ini/ini v1.55.0
	github.com/gobuffalo/packr v1.30.1 // indirect
	github.com/hankchan/goctp v0.0.0-20200420075836-7b6ba0f25997
	github.com/hprose/hprose-golang v2.0.5+incompatible
	github.com/jinzhu/gorm v1.9.12
	github.com/nntaoli-project/goex v1.1.0
	github.com/robertkrimen/otto v0.0.0-20191219234010-c382bd3c16ff
	github.com/sbinet/go-python v0.1.0
	github.com/smartystreets/goconvey v1.6.4 // indirect
	gonum.org/v1/gonum v0.7.0
	gopkg.in/alexcesaro/quotedprintable.v3 v3.0.0-20150716171945-2caba252f4dc // indirect
	gopkg.in/gomail.v2 v2.0.0-20160411212932-81ebce5c23df // indirect
	gopkg.in/ini.v1 v1.55.0 // indirect
	gopkg.in/logger.v1 v1.0.1
	gopkg.in/sourcemap.v1 v1.0.5 // indirect
)

replace (
	github.com/hankchan/goctp v0.0.0-20200420075836-7b6ba0f25997 => github.com/xiyanxiyan10/goctp-1 v0.0.0-20200420075836-7b6ba0f25997
	github.com/nntaoli-project/goex v1.1.0 => github.com/xiyanxiyan10/goex v1.0.9-0.20200806110227-7c44f68a1375
)
