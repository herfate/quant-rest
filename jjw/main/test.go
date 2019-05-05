package main

import (
	"fmt"
	"time"

	huobiService "github.com/findthefirst/quant-rest/exchange/huobi/services"
	"github.com/findthefirst/quant-rest/jjw/utils/log"
	"github.com/findthefirst/quant-rest/jjw/core"
	"github.com/findthefirst/quant-rest/jjw/core/strategy"
	"github.com/findthefirst/quant-rest/jjw/model"

	"github.com/syndtr/goleveldb/leveldb"
	"github.com/findthefirst/quant-rest/jjw/core/client"
)

func main() {
	fmt.Println("主线程print")


	t := &strategy.MartingaleGrid{CoinSymbol: "ont_usdt", StrategyName: "grid", Exchange: client.NewHuobiClient(), BeginAsset: 0.1, AllAsset: 0.5, Orders: make([]model.Order, 0)}
	core.GoTask(t)
	for {
		if core.IsExit() {
			break
		}
		time.Sleep(time.Duration(1)*time.Second)
	}
	log.Info("应用结束...")
}

func main2()  {
	kline := huobiService.GetKLine("ontusdt", "15min", 100)
	fmt.Println(kline)

	db, err := leveldb.OpenFile("db", nil)
	defer db.Close()
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println( db)

	_ = db.Put([]byte("key1"), []byte("好好检查"), nil)
	_ = db.Put([]byte("key2"), []byte("天天向上"), nil)
	_ = db.Put([]byte("key:3"), []byte("就会一个本事"), nil)

	log.ErrorAndWrap("测试log")
	log.WarningAndWrap("测试log")
}

