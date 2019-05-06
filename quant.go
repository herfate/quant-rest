package main

import (
	"github.com/findthefirst/quant-rest/jjw/core/strategy"
	"github.com/findthefirst/quant-rest/jjw/core"
	"github.com/findthefirst/quant-rest/jjw/utils/log"
	"github.com/findthefirst/quant-rest/jjw/model"
	"github.com/findthefirst/quant-rest/jjw/core/client"
	"github.com/findthefirst/quant-rest/jjw/core/constants"
	"github.com/findthefirst/quant-rest/jjw/utils"
)

func main() {
	log.Debug("应用开始...")

	t := &strategy.MartingaleGrid{CoinSymbol: "ont_usdt", StrategyName: constants.StrategyConfigMartingaleGrid, Exchange: client.NewHuobiClient(), BeginAsset: 0.1, AllAsset: 0.5, Orders: make([]model.Order, 0)}
	core.GoTask(t)
	for {
		if core.IsExit() {
			break
		}
		utils.Sleep(3000)
	}
	log.WarningAndWrap("应用结束...")
}
