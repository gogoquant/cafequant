package gobacktest

// MarryHandler .
type MarryHandler interface {
	Marry(handler PortfolioHandler) (bool, error)
}
