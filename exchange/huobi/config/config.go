package config

import (
	"github.com/findthefirst/quant-rest/jjw/utils/config"
	"github.com/findthefirst/quant-rest/jjw/utils/log"
)

// API KEY
var (
	// todo: replace with your own AccessKey and Secret Key
	ACCESS_KEY string = "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
	SECRET_KEY string = "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"

	// default to be disabled, please DON'T enable it unless it's officially announced.
	ENABLE_PRIVATE_SIGNATURE bool = false

	// generated the key by: openssl ecparam -name prime256v1 -genkey -noout -out privatekey.pem
	// only required when Private Signature is enabled
	// todo: replace with your own PrivateKey from privatekey.pem
	PRIVATE_KEY_PRIME_256 string = `xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx`

)

// API请求地址, 不要带最后的/
const (
	//todo: replace with real URLs and HostName
	MARKET_URL string = "https://api.huobi.br.com"
	TRADE_URL  string = "https://api.huobi.br.com"
	HOST_NAME  string = "api.huobi.br.com"
)

func init() {
	conf, err := config.GetConfig()
	if err != nil {
		log.ErrorAndWrap("load config error")
	}
	ACCESS_KEY = conf.HuobiKey
	SECRET_KEY = conf.HuobiSecret
}
