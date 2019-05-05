package calculate

import "github.com/findthefirst/quant-rest/jjw/model"

//  calculate macd
func CalculateMACD(tickers []*model.Ticker) {
	ema12 := 0.0
	ema26 := 0.0
	dif := 0.0
	dea := 0.0
	macd := 0.0

	for i := 0; i < len(tickers); i++ {
		point := tickers[i]
		closePrice := point.Close
		if i == 0 {
			ema12 = closePrice
			ema26 = closePrice
		} else {
			// EMA（12） = 前一日EMA（12） X 11/13 + 今日收盘价 X 2/13
			ema12 = ema12 * 11 / 13 + closePrice * 2 / 13
			// EMA（26） = 前一日EMA（26） X 25/27 + 今日收盘价 X 2/27
			ema26 = ema26 * 25 / 27 + closePrice * 2 / 27
		}
		// DIF = EMA（12） - EMA（26） 。
		// 今日DEA = （前一日DEA X 8/10 + 今日DIF X 2/10）
		// 用（DIF-DEA）*2即为MACD柱状图。
		dif = ema12 - ema26;
		dea = dea * 8 / 10 + dif * 2 / 10
		macd = (dif - dea) * 2
		point.Dif = dif
		point.Dea = dea
		point.Macd = macd
	}
}

