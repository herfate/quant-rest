package client

import (
	"testing"
	"fmt"
	"encoding/json"
	"strings"
	"strconv"
	"github.com/findthefirst/quant-rest/jjw/model"
	"github.com/findthefirst/quant-rest/jjw/core/constants"
	"time"
)

func TestQueryKline(t *testing.T)  {
	hb := &HuobiClient{}
	kline, err := hb.QueryKline("btc_usdt", "15min", 10)
	if err != nil {
		t.Errorf(" get kline error %v", err)
	} else {
		v, _ := json.Marshal(kline)
		fmt.Printf(" %s   ", v)
	}
}

func TestQueryOrderBook(t *testing.T) {
	hb := &HuobiClient{}
	book, err := hb.QueryOrderBook("btC_usdt", 0)
	if err != nil {
		t.Errorf(" get orderbook error %v", err)
	} else {
		v, _ := json.Marshal(book)
		fmt.Printf(" %s   ", v)
	}
}

func TestFmtPrice(t *testing.T) {
	amountStr := fmtAmount(0.12345678, "ont_usdt")
	l, err := strconv.Atoi(strings.Split(coinMap["ontusdt"], ",")[1])
	if len(amountStr) != (l + 2) || err != nil {
		t.Errorf(" format price error ")
	}
}

func TestQuerySymbols(t *testing.T) {
	infos, _ := querySymbols()
	if len(infos) < 1 {
		t.Errorf(" TestQuerySymbols no data error ")
	}

}


func TestCreateOrder(t *testing.T) {
	hb := &HuobiClient{}
	createO := &model.Order{}
	createO.SimplePlace(constants.OrderTypeSell,  "theta_usdt",  constants.ExchangeNameHuobi,  0.096,  2.8,  constants.OrderStatusSubmitted,  constants.SpotOrder,   "",   "")
	err := hb.CreateOrder(createO)
	if err != nil {
		t.Errorf(" TestCreateOrder place order error %v", createO)
	} else {
		j, _ := json.Marshal(createO)
		fmt.Println(string(j))
	}
}


func TestUpdateOrder(t *testing.T) {
	hb := &HuobiClient{}
	o := &model.Order{}
	o.OrderId = "31317074460"
	err := hb.UpdateOrder(o)
	if err != nil {
		t.Errorf(" TestCreateOrder update order error %v", o)
	} else {
		j, _ := json.Marshal(o)
		fmt.Println(string(j))
	}
}


func TestOther(t *testing.T) {
	fmt.Println(time.Now().Unix())
	fmt.Println(time.Now().UnixNano()/1e6)
}
