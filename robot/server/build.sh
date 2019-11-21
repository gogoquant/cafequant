gopath=`cd $(dirname $0); pwd -P`
printf "prepare gopath is %s\n" $gopath

export GOPATH=$gopath
go build --tags netgo src/github.com/xiyanxiyan10/quantcore/QuantBot.go
