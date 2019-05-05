package client

import "github.com/findthefirst/quant-rest/jjw/model"

type TxClient interface {
	QueryKline(coinSymbol string, klineType string, size int) (tickers []*model.Ticker, err error)
	QueryOrderBook(coinSymbol string, step int) (book *model.OrderBook, err error)
	CreateOrder(createOrder *model.Order) (err error)
	UpdateOrder(oldOrder *model.Order) (err error)
	CancelOrder(order *model.Order) (err error)
	GetPricePrecision(coinSymbol string) (precision int)
	GetFee() (fee float64)
	IsSupport(coinSymbol string) (isSp bool)
	GetExchangeName() (name string)
}


