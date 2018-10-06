package gobacktest

// MarryHandler .
type MarryHandler interface {
	Marry(bt *Backtest, data DataEvent) (end bool, err error)
}
