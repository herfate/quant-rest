package constants

// trade types
const (
	OrderTypeBuy        = "buy"
	OrderTypeSell       = "sell"
)

//
const (
	SpotOrder         = iota
	MarginOrder
)


// ex
const (
	ExchangeNameHuobi     = "huobi"
	ExchangeNameGate      = "gate"
)

// order tx status
const (
	OrderStatusSubmitted  		= "submitted"
	OrderStatusPartialFilled    = "partial-filled"
	OrderStatusPartialCanceled  = "partial-canceled"
	OrderStatusFilled	  	    = "filled"
	OrderStatusCanceled		    = "canceled"
)

// k line type
const (
	KlineType1min  		= "1min"
	KlineType5min       = "5min"
	KlineType15min      = "15min"
	KlineType30min	    = "30min"
	KlineType60min      = "60min"
)

// strategy name 网格
const (
	StrategyNameTest  		= "test_strategy"
)



// my definition
const (
	StrategyConfigMartingaleGridSleepTime   	= "martingale_grid_sleep_time"
	StrategyConfigMartingaleGridInterval  		= "martingale_grid_interval"
	StrategyConfigMartingaleGridOpenOrClose 	= "martingale_grid_open_or_close"
	StrategyConfigMartingaleGridBeginAsset   	= "martingale_grid_begin_asset"
)