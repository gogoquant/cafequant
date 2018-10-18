package handler

import (
	"fmt"
	"github.com/xiyanxiyan10/samaritan/model"
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
	self, err := model.GetUser(username)
	if err != nil {
		resp.Message = fmt.Sprint(err)
		return
	}
	if _, err = self.GetTrader(id); err != nil {
		resp.Message = fmt.Sprint(err)
		return
	}
	master := trader.GlobalDataGram()
	if master == nil{
		resp.Message = "datagram not found"
		return
	}
	cmd := fmt.Sprintf("select * from data_%d",ID)
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

func (datagram) Delete(id string, ctx rpc.Context) (resp response) {
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
	self, err := model.GetUser(username)
	if err != nil {
		resp.Message = fmt.Sprint(err)
		return
	}
	if _, err = self.GetTrader(id); err != nil {
		resp.Message = fmt.Sprint(err)
		return
	}
	master := trader.GlobalDataGram()
	if master == nil{
		resp.Message = "datagram not found"
		return
	}
	cmd := fmt.Sprintf("drop measurement  data_%d",ID)
	_, _, err = master.QueryDB(cmd)
	if err != nil{
		resp.Message = err.Error()
		return
	}

	resp.Data = "reset success"
	resp.Success = true
	return
}