package calculate

import (
	"github.com/findthefirst/quant-rest/jjw/model"
	"math"
)

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
			ema12 = ema12 * 11.0 / 13.0 + closePrice * 2.0 / 13.0
			// EMA（26） = 前一日EMA（26） X 25/27 + 今日收盘价 X 2/27
			ema26 = ema26 * 25.0 / 27.0 + closePrice * 2.0 / 27.0
		}
		// DIF = EMA（12） - EMA（26） 。
		// 今日DEA = （前一日DEA X 8/10 + 今日DIF X 2/10）
		// 用（DIF-DEA）*2即为MACD柱状图。
		dif = ema12 - ema26
		dea = dea * 8.0 / 10.0 + dif * 2.0 / 10.0
		macd = (dif - dea) * 2.0
		point.Dif = dif
		point.Dea = dea
		point.Macd = macd
	}
}


func CalculateRSI(tickers []*model.Ticker) {
	rsi := 0.0
	rsiABSEma := 0.0
	rsiMaxEma := 0.0
	for i := 0; i < len(tickers); i++ {
		point := tickers[i]
		closePrice := point.Close
		if i == 0 {
			rsi = 0.0
			rsiABSEma = 0.0
			rsiMaxEma = 0.0
		} else {
			Rmax := math.Max(0, closePrice - tickers[i - 1].Close)
			RAbs := math.Abs(closePrice - tickers[i - 1].Close)

			rsiMaxEma = (Rmax + (14.0 - 1.0) * rsiMaxEma) / 14.0
			rsiABSEma = (RAbs + (14.0 - 1.0) * rsiABSEma) / 14.0
			rsi = (rsiMaxEma / rsiABSEma) * 100
		}
		if i < 13 {
			rsi = 0.0
		}
		if math.IsNaN(rsi) {
			rsi = 0.0
		}
		tickers[i].Rsi = rsi
	}
}

func CalculateMA(tickers []*model.Ticker) {
	ma5 := 0.0
	ma10 := 0.0
	ma20 := 0.0
	ma30 := 0.0
	ma60 := 0.0
	for i := 0; i < len(tickers); i++ {
		point := tickers[i]
		closePrice := point.Close

		ma5 += closePrice
		ma10 += closePrice
		ma20 += closePrice
		ma30 += closePrice
		ma60 += closePrice

		if i == 4 {
			point.MA5Price = ma5 / 5.0
		} else if i >= 5 {
			ma5 -= tickers[i - 5].Close
			point.MA5Price = ma5 / 5.0
		} else {
			point.MA5Price = 0.0
		}

		if i == 9 {
			point.MA10Price = ma10 / 10.0
		} else if i >= 10 {
			ma10 -= tickers[i - 10].Close
			point.MA10Price = ma10 / 10.0
		} else {
			point.MA10Price = 0.0
		}

		if i == 19 {
			point.MA20Price = ma20 / 20.0
		} else if i >= 20 {
			ma20 -= tickers[i - 20].Close
			point.MA20Price = ma20 / 20.0
		} else {
			point.MA20Price = 0.0
		}

		if i == 29 {
			point.MA30Price = ma30 / 30.0
		} else if i >= 30 {
			ma30 -= tickers[i - 30].Close
			point.MA30Price = ma30 / 30.0
		} else {
			point.MA30Price = 0.0
		}
		if i == 59 {
			point.MA60Price = ma60 / 60.0
		} else if i >= 60 {
			ma60 -= tickers[i - 60].Close
			point.MA60Price = ma60 / 60.0
		} else {
			point.MA60Price = 0.0
		}
	}
}

func CalculateBOLL(tickers []*model.Ticker) {
	// must calculate MA before ...
	if len(tickers) > 0 && tickers[0].MA5Price == 0 {
		CalculateMA(tickers)
	}
	for i := 0; i < len(tickers); i++ {
		point := tickers[i]
		if i < 19 {
			point.Mb = 0.0
			point.Up = 0.0
			point.Dn = 0.0
		} else {
			n := 20
			md := 0.0
			for j := i - n + 1; j <= i; j++ {
				c := tickers[j].Close
				m := point.MA20Price
				value := c - m
				md += value * value
			}
			md = md / (float64(n) - 1.0)
			md = math.Sqrt(md)
			point.Mb = point.MA20Price
			point.Up = point.Mb + 2.0 * md
			point.Dn = point.Mb - 2.0 * md
		}
	}
}


func CalculateKDJ(tickers []*model.Ticker) {
	k := 0.0
	d := 0.0
	for i := 0; i < len(tickers); i++ {
		point := tickers[i]
		closePrice := point.Close
		startIndex := i - 13
		if startIndex < 0 {
			startIndex = 0
		}
		max14 := -math.MaxFloat32
		min14 := math.MaxFloat32
		for index := startIndex; index <= i; index++ {
			max14 = math.Max(max14, tickers[index].High)
			min14 = math.Min(min14, tickers[index].Low)
		}
		rsv := 100.0 * (closePrice - min14) / (max14 - min14)
		if math.IsNaN(rsv) {
			rsv = 0.0
		}
		if i == 0 {
			k = 50.0
			d = 50.0
		} else {
			k = (rsv + 2.0 * k) / 3.0
			d = (k + 2.0 * d) / 3.0
		}

		if i < 13 {
			point.K = 0.0
			point.D = 0.0
			point.J = 0.0
		} else if i == 13 || i == 14 {
			point.K = k
			point.D = 0.0
			point.J = 0.0
		} else {
			point.K = k
			point.D = d
			point.J = 3.0 * k - 2.0 * d
		}
	}
}




