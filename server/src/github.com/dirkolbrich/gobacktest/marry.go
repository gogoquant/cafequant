package gobacktest

// MarryHandler .
type MarryHandler interface {
	Marry(bt *Backtest, stockType string) (bool, error)
}
