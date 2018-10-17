package handler

import (
	"fmt"
	"strconv"

	"github.com/hprose/hprose-golang/rpc"
	"github.com/xiyanxiyan10/samaritan/constant"
	"github.com/xiyanxiyan10/samaritan/trader"
)

type datagram struct{}

func (datagram) List(id string, mode string, ctx rpc.Context) (resp response) {
	username := ctx.GetString("username")
	if username == "" {
		resp.Message = constant.ErrAuthorizationError
		return
	}
	ID, err := strconv.ParseInt(id, 10, 64)
	if err != nil{
		resp.Message = err.Error()
		return
	}
	global, err := trader.GetTrader(ID)
	if err != nil{
		resp.Message = "trader not found"
		return
	}
	master := global.Datagram()
	if master == nil{
		resp.Message = "datagram not found"
		return
	}
	cmd := fmt.Sprintf("select * from name_%d",ID)
	items, tables, err := master.QueryDB(cmd)
	if err != nil{
		resp.Message = err.Error()
		return
	}

	resp.Data = struct {
		List  interface{}
		Col []string
		Mode string
	}{
		List:  items,
		Col: tables,
		Mode: mode,
	}
	resp.Success = true
	return
}