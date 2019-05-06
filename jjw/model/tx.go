package model



type Ticker struct {
	Open      float64	`json:"open"`
	Close     float64	`json:"close"`
	High      float64	`json:"high"`
	Low       float64	`json:"low"`
	Vol       float64	`json:"vol"`
	Time      int64     `json:"time"`
	FmtTime   string 	`json:"fmtTime"`
	Rsi		  float64   `json:"rsi"`
	Macd	  float64   `json:"macd"`
	Dif	      float64   `json:"dif"`
	Dea	      float64   `json:"dea"`
	MA5Price  float64   `json:"ma5Price"`
	MA10Price float64   `json:"ma10Price"`
	MA20Price float64   `json:"ma20Price"`
	MA30Price float64   `json:"ma30Price"`
	MA60Price float64   `json:"ma60Price"`
	Up        float64   `json:"up"`
	Mb        float64   `json:"mb"`
	Dn        float64   `json:"dn"`
	K         float64   `json:"k"`
	D         float64   `json:"d"`
	J         float64   `json:"j"`
}


type Order struct {
	RandomId  string
	OrderId   string
	OrderType string
	ExtraOrderType string
	CoinSymbol string
	ExchangeName string
	StockPrice float64
	StockAmount float64
	CompletePrice float64
	CompleteAmount float64
	Status string
	SpotOrMargin int
	StrategyName string
	ExtraMsg string
	CreateTime int64
	FormatCreateTime string
	TargetOrderId string
}

func (o *Order) SimplePlace(orderType string, coinSymbol string, exchangeName string, stockPrice float64, stockAmount float64, status string, spotOrMargin int, strategyName string, extraMsg string) {
	o.OrderType = orderType
	o.CoinSymbol = coinSymbol
	o.ExchangeName = exchangeName
	o.StockPrice = stockPrice
	o.StockAmount = stockAmount
	o.Status = status
	o.SpotOrMargin = spotOrMargin
	o.StrategyName = strategyName
	o.ExtraMsg = extraMsg
}

type SingleBook struct {
	Price  float64 //价格
	Amount float64 //市场深度量
}

type OrderBook struct {
	Bids []SingleBook //买单市场深度列表
	Buy  float64     //买一价, Bids[0].Price
	Avg  float64     //(Buy + Sell) / 2
	Sell float64     //卖一价, Asks[0].Price
	Asks []SingleBook //卖单市场深度列表
}