package handler

import (
	"fmt"

	"github.com/hprose/hprose-golang/rpc"
	"github.com/xiyanxiyan10/quantcore/constant"
	"github.com/xiyanxiyan10/quantcore/model"
)

type traderExt struct{}

// List ...
func (traderExt) List(traderID int64, ctx rpc.Context) (resp response) {
	username := ctx.GetString("username")
	if username == "" {
		resp.Message = constant.ErrAuthorizationError
		return
	}

	var um model.User
	exts, err := um.ListParameters(traderID)
	if err != nil {
		resp.Message = fmt.Sprint(err)
		return
	}
	resp.Data = struct {
		Total int64
		List  []model.TraderExt
	}{
		Total: int64(len(exts)),
		List:  exts,
	}
	resp.Success = true
	return
}

// Put ...
func (traderExt) Put(req model.TraderExt, traderID int64, ctx rpc.Context) (resp response) {
	username := ctx.GetString("username")
	if username == "" {
		resp.Message = constant.ErrAuthorizationError
		return
	}
	self, err := model.GetUser(username)
	if err != nil {
		resp.Message = fmt.Sprint(err)
		return
	}
	db, err := model.NewOrm()
	if err != nil {
		resp.Message = fmt.Sprint(err)
		return
	}
	if _, err := self.GetTrader(traderID); err != nil {
		resp.Message = fmt.Sprint(err)
		return
	}
	defer db.Close()
	db = db.Begin()

	if req.ID > 0 {
		var ext model.TraderExt
		if err := model.DB.First(&ext, req.ID).Error; err != nil {
			resp.Message = fmt.Sprint(err)
			return
		}
		ext.Content = req.Content
		ext.Desc = req.Desc
		ext.Val = req.Val
		if err := model.DB.Save(&ext).Error; err != nil {
			resp.Message = fmt.Sprint(err)
			return
		}
		resp.Success = true
		return
	}

	req.TraderID = traderID
	if err := db.Create(&req).Error; err != nil {
		resp.Message = fmt.Sprint(err)
		db.Rollback()
		return
	}

	if err := db.Commit().Error; err != nil {
		db.Rollback()
		resp.Message = fmt.Sprint(err)
		return
	}
	resp.Success = true
	return
}

// Delete ...
func (traderExt) Delete(ext model.TraderExt, ctx rpc.Context) (resp response) {
	username := ctx.GetString("username")
	if username == "" {
		resp.Message = constant.ErrAuthorizationError
		return
	}
	self, err := model.GetUser(username)
	if err != nil {
		resp.Message = fmt.Sprint(err)
		return
	}
	if _, err := self.GetTrader(ext.TraderID); err != nil {
		resp.Message = fmt.Sprint(err)
		return
	}
	if err := self.DeleteParameter(ext.ID); err != nil {
		resp.Message = fmt.Sprint(err)
		return
	}
	resp.Success = true
	return
}
