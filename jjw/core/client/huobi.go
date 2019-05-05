package client

import (
	"fmt"
	"errors"
	"strings"
	"strconv"
	"encoding/json"

	"github.com/findthefirst/quant-rest/jjw/model"
	hb "github.com/findthefirst/quant-rest/exchange/huobi/services"
	"github.com/findthefirst/quant-rest/jjw/utils"
	"github.com/findthefirst/quant-rest/jjw/utils/log"
	"github.com/findthefirst/quant-rest/exchange/huobi/models"
	"github.com/findthefirst/quant-rest/jjw/core/constants"
)

var coinMap = make(map[string]string)
var spotAccountId string

type HuobiClient struct {
	name string
}

func NewHuobiClient() *HuobiClient {
	return &HuobiClient{name: constants.ExchangeNameHuobi}
}

func (c *HuobiClient) QueryKline(coinSymbol string, klineType string, size int) (tickers []*model.Ticker, err error) {
	defer func() {
		if err := recover(); err != nil {
			err = log.ErrorAndWrap(fmt.Sprintf("---huobi--- query kline error %v", err))
		}
	}()
	kline := hb.GetKLine(convCoinSymbol(coinSymbol), klineType, size)
	if kline.Data == nil {
		return tickers, log.ErrorAndWrap("---huobi--- query kline no data")
	}
	//err = utils.DeepCopy(kline.Data, &tickers)
	for _, it := range kline.Data {
		ticker := &model.Ticker{Open : it.Open, Close : it.Close, High : it.High, Low : it.Low, Vol : it.Vol, Time : it.ID, FmtTime : utils.SecFormat(it.ID)}
		tickers = append(tickers, ticker)
	}
	for i, j := 0, len(tickers)-1; i < j; i, j = i+1, j-1 {
		tickers[i], tickers[j] = tickers[j], tickers[i]
	}
 	return tickers, err
}

func (c *HuobiClient) QueryOrderBook(coinSymbol string, step int) (book *model.OrderBook, err error) {
	depth := hb.GetMarketDepth(convCoinSymbol(coinSymbol), "step0")
	asks := depth.Tick.Asks
	bids := depth.Tick.Bids
	book = &model.OrderBook{}
	for _, it := range asks { //book.Asks = make([]model.SingleBook, len(depth.Tick.Asks))
		book.Asks = append(book.Asks, model.SingleBook{Price: it[0], Amount:it[1]})
	}
	for _, it := range bids {
		book.Bids = append(book.Bids, model.SingleBook{Price: it[0], Amount:it[1]})
	}
	if len(book.Asks) > 0 && len(book.Bids) > 0 {
		book.Buy = bids[0][0]
		book.Sell = asks[0][0]
		book.Avg = (book.Buy + book.Sell) / 2
	} else {
		err = errors.New("---huobi--- query OrderBook error, empty data")
	}
	return book, err
}

func (c *HuobiClient) CreateOrder(createOrder *model.Order) (err error) {
	txCoinSymbol := convCoinSymbol(createOrder.CoinSymbol)
	if createOrder == nil {
		return log.ErrorAndWrap("---huobi--- order is nil")
	}
	if !c.IsSupport(txCoinSymbol) {
		return log.ErrorAndWrap("---huobi--- unsupport token")
	}
	req := models.PlaceRequestParams{}
	switch createOrder.SpotOrMargin {
	case 1:
		req.AccountID = ""
	default:
		req.AccountID = spotAccountId
	}
	req.Amount = fmtAmount(createOrder.StockAmount, txCoinSymbol)
	req.Price = fmtPrice(createOrder.StockPrice, txCoinSymbol)
	req.Symbol = txCoinSymbol
	req.Type = createOrder.OrderType + "-limit"//constants.OrderTypeBuy
	if createOrder.ExtraOrderType != "" { req.Type = createOrder.ExtraOrderType }
	req.Source = "api"
	log.Debug(fmt.Sprintf("---huobi--- prepare createOrder : %s %s price %s amount %s", req.Type, req.Symbol, req.Price, req.Amount))
	placeRtn := hb.Place(req)
	if placeRtn.Data != "" {
		now := utils.Now()
		randomId := utils.RandomId()
		createOrder.OrderId = placeRtn.Data
		createOrder.RandomId = randomId
		createOrder.CreateTime = now
		createOrder.FormatCreateTime = utils.MilliSecFormat(now)
		createOrder.ExchangeName = constants.ExchangeNameHuobi
		createOrder.Status = constants.OrderStatusSubmitted
		log.Debug(fmt.Sprintf("---huobi--- success createOrder : %s %s price %s amount %s", req.Type, req.Symbol, req.Price, req.Amount))
	} else {
		err = log.WarningAndWrap(fmt.Sprintf("---huobi--- failure createOrder : %s %s  errormsg %s", req.Type, req.Symbol, placeRtn.ErrMsg))
	}
	return err
}

const (
	KState 			= "state"
	KCompleteAmount = "field-amount"
	KCompleteCash   = "field-cash-amount"
)
func (c *HuobiClient) UpdateOrder(oldOrder *model.Order) (err error) {
	if oldOrder != nil && oldOrder.OrderId != "" {
		rtn := hb.GetOrderDetail(oldOrder.OrderId)
		//{"status":"ok","data":{"id":31317074460,"symbol":"thetausdt","account-id":1111,"amount":"2.800000000000000000","price":"0.096000000000000000","created-at":1556610593363,"type":"sell-limit","field-amount":"0.0","field-cash-amount":"0.0","field-fees":"0.0","finished-at":0,"source":"api","state":"submitted","canceled-at":0}}
		if rtn.Data != nil {
			state, ok := rtn.Data[KState].(string)
			if ok {
				oldOrder.Status = state
			}
			completeAmount, ok := rtn.Data[KCompleteAmount].(string)
			if ok {
				ca, err := strconv.ParseFloat(completeAmount, 64)
				if err != nil {
					oldOrder.CompleteAmount = ca
					completeCash, ok := rtn.Data[KCompleteCash].(string)
					if ok {
						cc, err := strconv.ParseFloat(completeCash, 64)
						if err != nil {
							if ca == 0 {
								oldOrder.CompletePrice = 0
							} else {
								oldOrder.CompletePrice = cc / ca
							}
						}
					}
				}
			}
			//fmt.Println(rtn.Data)
			//oldOrder.Status =
		} else {
			err = log.WarningAndWrap(fmt.Sprintf("---huobi--- failure updateOrder : %v  ", oldOrder))
		}
	}
	return err
}

func (c *HuobiClient) CancelOrder(order *model.Order) (err error) {
	if order != nil && order.OrderId != "" {
		hb.SubmitCancel(order.OrderId)
		utils.Sleep(200)
		oldComplete := order.CompleteAmount
		c.UpdateOrder(order)
		if order.CompleteAmount / oldComplete > 1.01 {
			err = log.ErrorAndWrap("---huobi---  cancel order failure : before %f, after %f", oldComplete,  order.CompleteAmount)
		}
	}
	return err
}

func (c *HuobiClient) GetPricePrecision(coinSymbol string) (precision int) {
	ap := coinMap[convCoinSymbol(coinSymbol)]
	pfmt, err := strconv.Atoi(strings.Split(ap, ",")[1])
	if err != nil {
		return -1
	}
	return pfmt
}

func (c *HuobiClient) GetFee() (fee float64) {
	return 0.00049 // 0.002 * 0.245
}
func (c *HuobiClient) GetExchangeName() (name string) {
	return c.name
}

func (c *HuobiClient) IsSupport(coinSymbol string) (isSp bool) {
	return coinMap[coinSymbol] != ""
}

/*(c *HuobiClient)*/
func querySymbols() (symbols []models.SymbolsData, err error) {
	r := hb.GetSymbols()
	if len(r.Data) < 1 {
		return nil, log.ErrorAndWrap("---huobi--- get symbol empty data error")
	}
	return hb.GetSymbols().Data, err
}

func init()  {
	// 获取精确度信息
	symbols, err := querySymbols()
	if err != nil {
		log.ErrorAndWrap("---huobi--- get symbols error")
	}
	for _, it := range symbols {
		coinMap[it.BaseCurrency + it.QuoteCurrency] = fmt.Sprintf("%d,%d", it.AmountPrecision, it.PricePrecision)
	}
	json, _ := json.Marshal(coinMap)
	log.Debug(string(json))

	//spotAccountId
	acs := hb.GetAccounts().Data
	for _, it := range acs {
		if it.Type == "spot" {
			spotAccountId = strconv.FormatInt(it.ID, 10)
		}
	}
}

func convCoinSymbol(coinSymbol string) (txSymbol string) {
	return strings.ToLower(strings.Replace(coinSymbol, "_", "", -1))
}

func fmtPrice(price float64, coinSymbol string) (priceStr string) {
	ap := coinMap[convCoinSymbol(coinSymbol)]
	//afmt, err1 := strconv.Atoi(strings.Split(ap, ",")[0])
	pfmt, err := strconv.Atoi(strings.Split(ap, ",")[1])
	if err != nil {
		log.ErrorAndWrap("fmt price  error %v" , err)
	}
	return strconv.FormatFloat(price, 'f', pfmt, 64)
}

func fmtAmount(amount float64, coinSymbol string) (amountStr string) {
	ap := coinMap[convCoinSymbol(coinSymbol)]
	afmt, err := strconv.Atoi(strings.Split(ap, ",")[0])
	//pfmt, err := strconv.Atoi(strings.Split(ap, ",")[0])
	if err != nil {
		log.ErrorAndWrap("fmt amount  error %v" , err)
	}
	return strconv.FormatFloat(amount, 'f', afmt, 64)
}