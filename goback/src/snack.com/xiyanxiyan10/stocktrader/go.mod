module snack.com/xiyanxiyan10/stocktrader

go 1.13

require (
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/go-echarts/go-echarts v0.0.0-20190915064101-cbb3b43ade5d
	github.com/go-gomail/gomail v0.0.0-20160411212932-81ebce5c23df
	github.com/go-ini/ini v1.55.0
	github.com/go-sql-driver/mysql v1.5.0 // indirect
	github.com/gobuffalo/packr v1.30.1 // indirect
	github.com/gorilla/websocket v1.4.2 // indirect
	github.com/hprose/hprose-golang v2.0.5+incompatible
	github.com/jinzhu/gorm v1.9.12
	github.com/json-iterator/go v1.1.10 // indirect
	github.com/klauspost/compress v1.10.3 // indirect
	github.com/lib/pq v1.2.0 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.1 // indirect
	github.com/nntaoli-project/goex v1.2.2
	github.com/pkg/errors v0.9.1 // indirect
	github.com/robertkrimen/otto v0.0.0-20191219234010-c382bd3c16ff
	github.com/sbinet/go-python v0.1.0
	github.com/sirupsen/logrus v1.6.0 // indirect
	github.com/stretchr/testify v1.6.1 // indirect
	golang.org/x/crypto v0.0.0-20200323165209-0ec3e9974c59 // indirect
	golang.org/x/net v0.0.0-20200501053045-e0ff5e5a1de5 // indirect
	golang.org/x/sys v0.0.0-20200501052902-10377860bb8e // indirect
	gopkg.in/alexcesaro/quotedprintable.v3 v3.0.0-20150716171945-2caba252f4dc // indirect
	gopkg.in/check.v1 v1.0.0-20190902080502-41f04d3bba15 // indirect
	gopkg.in/gomail.v2 v2.0.0-20160411212932-81ebce5c23df // indirect
	gopkg.in/logger.v1 v1.0.1
	gopkg.in/sourcemap.v1 v1.0.5 // indirect
	snack.com/xiyanxiyan10/stockdb v0.0.0-00010101000000-000000000000
)

replace (
	github.com/nntaoli-project/goex v1.1.0 => github.com/xiyanxiyan10/goex v1.0.9-0.20200806110227-7c44f68a1375
	snack.com/xiyanxiyan10/conver => ../conver
	snack.com/xiyanxiyan10/stockdb => ../stockdb
)
