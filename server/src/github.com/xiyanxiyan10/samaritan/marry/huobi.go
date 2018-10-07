package marry

import (
	"errors"
	goback "github.com/xiyanxiyan10/gobacktest"
)

func init() {
	//@todo more coin type register
	marryStore["datxbtc"] = new(MarryHuobi)
}

type MarryHuobi struct {

}

// Marry
func (bt *MarryHuobi) Marry(back *goback.Backtest, data goback.DataEvent) (bool, error) {
	stockType := data.Symbol()
	orders, ok := back.OrdersBySymbol(stockType)
	if ok != true {
		return true, nil
	}
	if ok != true {
		return false, errors.New("get latest fail")
	}
	for _, order := range orders {
		status := order.Status()
		if status == goback.OrderCanceled || status == goback.OrderCancelPending {
			continue
		}
		dir := order.Direction()
		var err error
		switch dir {
		case goback.BOT:
			if order.FQty() >= data.High() {
				_, err = back.CommitOrder(order.ID())
			}
		case goback.SLD:
			if order.FQty() <= data.Low() {
				_, err = back.CommitOrder(order.ID())
			}
		default:
			return false, errors.New("unknown dir")
		}
		if err != nil {
			return false, err
		}
	}
	return true, nil
}


