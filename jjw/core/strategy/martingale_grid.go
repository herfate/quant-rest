package strategy

import (
	"fmt"
	"strconv"

	"github.com/findthefirst/quant-rest/jjw/utils"
	"github.com/findthefirst/quant-rest/jjw/model"
	"github.com/findthefirst/quant-rest/jjw/core/constants"
	"github.com/findthefirst/quant-rest/jjw/core"
	"github.com/findthefirst/quant-rest/jjw/core/client"
	"github.com/findthefirst/quant-rest/jjw/core/calculate"
	"github.com/findthefirst/quant-rest/jjw/utils/log"
	"math"
)

type MartingaleGrid struct {
	CoinSymbol string
	StrategyName string
	Exchange client.TxClient
	Profit float64
	BeginAsset float64  // TODO-jjw 暂时用usdt计算
	AllAsset float64
	Orders []model.Order
	Exit bool
}

func (s *MartingaleGrid) GetDesc() string {
	return " MartingaleGrid run " + s.CoinSymbol + " in " + s.Exchange.GetExchangeName() + " exchange" //s.Exchange.
}
func (s *MartingaleGrid) IsExit() bool {
	return s.Exit
}
func (s *MartingaleGrid) Run() {
	// c
	client := s.Exchange
	for {
		sleepTimeStr := getStrConfig(constants.StrategyConfigMartingaleGridSleepTime)
		// 网格
		interval := getFloatConfig(constants.StrategyConfigMartingaleGridInterval)
		open := getStrConfig(constants.StrategyConfigMartingaleGridOpenOrClose) == "open"
		beginAsset := getFloatConfig(constants.StrategyConfigMartingaleGridBeginAsset)

		tickers, err := client.QueryKline(s.CoinSymbol, constants.KlineType15min, 100)
		if err != nil {
			continue
		}
		//fmt.Println(" ticker: " + utils.JsonMarshal(tickers))
		nowTicker := tickers[len(tickers) - 1]
		nowAvgPrice := nowTicker.Close
		//update orders
		orders := core.ImportantData.Orders[s.CoinSymbol + s.StrategyName + s.Exchange.GetExchangeName()]
		for i:=0; i < len(orders); i++ {
			it := orders[i]
			if constants.OrderStatusFilled != it.Status && constants.OrderStatusPartialCanceled != it.Status && constants.OrderStatusCanceled != it.Status {
				if err := client.UpdateOrder(it); err != nil {
					continue
				}
			}
			if it.SpotOrMargin == constants.SpotOrder && it.Status == constants.OrderStatusSubmitted && it.OrderType == constants.OrderTypeBuy {
				if nowAvgPrice/it.StockPrice > (1 + 0.004) {
					log.Debug("-------Martingale Grid------ too high,cancel order ")
					client.CancelOrder(it)
				}
			}
			if constants.OrderStatusCanceled == it.Status {
				orders = append(orders[:i], orders[i+1:]...)
				i--
			}
		}
		core.ImportantData.Orders[s.CoinSymbol + s.StrategyName + s.Exchange.GetExchangeName()] = orders
		fmt.Println("orders:   " + utils.JsonMarshal(orders))

		calculate.CalculateMACD(tickers)

		minPrice, maxPrice := getMinAndMaxPrice(orders)
		fmt.Printf( "minPrice: %f, maxPrice: %f, nowAvgPrice: %f", minPrice, maxPrice, nowAvgPrice)

		//---------------------- analysis in goods  --------------------------------
		if isBuy(tickers, orders) && open {
			orderPrice := nowAvgPrice / (1 + interval)
			if len(orders) == 0 {
				book, err := client.QueryOrderBook(s.CoinSymbol, 0)
				if err == nil {
					orderPrice = book.Sell
				}
			}
			log.Debug("------- Martingale Grid ---预加仓现货----- %s buy %f", s.CoinSymbol, orderPrice)
		//	BigDecimal amount = BigDecimal.valueOf(calAmountHalf(beginAsset / orderPrice.doubleValue(), spotBuyOrders.size() + 1))
			readyBuys := len(getOrders(orders, constants.OrderTypeBuy, constants.SpotOrder))
			orderAmount := beginAsset/nowAvgPrice * math.Pow(2, float64(readyBuys))
			if orderAmount*nowAvgPrice < 25 {
				extraMsg := "profitPrice.toString()"
				newOrder := &model.Order{OrderType: constants.OrderTypeBuy, CoinSymbol:s.CoinSymbol, ExchangeName: s.Exchange.GetExchangeName(),
								StockPrice:orderPrice, StockAmount:orderAmount, StrategyName:s.StrategyName, SpotOrMargin:constants.SpotOrder, ExtraMsg: extraMsg}
				err := client.CreateOrder(newOrder)
				if err == nil && newOrder.OrderId != "" {
					orders = append(orders, newOrder)
					core.ImportantData.Orders[s.CoinSymbol + s.StrategyName + s.Exchange.GetExchangeName()] = orders
				} else {
					log.ErrorAndWrap("------- Martingale Grid --- create order error")
				}
			}
		}

		sleepTime, err := strconv.Atoi(sleepTimeStr)
		if err != nil {
			sleepTime = 3
		}
		utils.Sleep(sleepTime * 1000)
	}
	defer func() {
		if err := recover(); err != nil {
			log.ErrorAndWrap("------- Martingale Grid --- panic err %v", err)
			s.Exit = true
		}
	}()
	s.Exit = true
}

func isBuy(tickers []*model.Ticker, orders []*model.Order) (ok bool) {
	if tickers[len(tickers) -1].Macd > 0 {
		for _, it := range orders {
			if it.OrderType == constants.OrderTypeBuy && (it.Status == constants.OrderStatusSubmitted || it.Status == constants.OrderStatusPartialFilled) {
				return false
			}
		}
		return true
	}
	return false
}

func getMinAndMaxPrice(os []*model.Order) (minPrice float64, maxPrice float64) {
	minPrice = 9999999.9
	maxPrice = -0.1
	for _, o := range os {
		if o.CompletePrice > maxPrice {
			maxPrice = o.CompletePrice
		}
		if o.CompletePrice < minPrice {
			minPrice = o.CompletePrice
		}
	}
	return minPrice, maxPrice
}

func getOrders(os []*model.Order, orderType string, spotOrMargin int) (needOs []*model.Order) {
	for _, o := range os {
		if o.OrderType == orderType && o.SpotOrMargin == spotOrMargin {
			needOs = append(needOs, o)
		}
	}
	return needOs
}

func getStrConfig(key string) (rtn string) {
	switch key {
	case "123":
		rtn = "undef"
	}
	v := core.ImportantData.GetConfig()[key]
	if v != "" {
		rtn = v
	}
	return rtn
}
func getFloatConfig(key string) (rtn float64) {
	switch key {
	case constants.StrategyConfigMartingaleGridInterval:
		rtn = 0.006
	}
	configMap := core.ImportantData.GetConfig()
	v, err := strconv.ParseFloat(configMap[key], 64)
	if err == nil {
		rtn = v
	}
	return rtn
}