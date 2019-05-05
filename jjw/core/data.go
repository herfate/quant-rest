package core

import (
	"github.com/findthefirst/quant-rest/jjw/model"
	"github.com/findthefirst/quant-rest/jjw/utils"
	"github.com/findthefirst/quant-rest/jjw/utils/config"
	"github.com/findthefirst/quant-rest/jjw/utils/log"
	"github.com/BurntSushi/toml"
	"os"
	"io/ioutil"
)

var ImportantData *StaticData

//
type StaticData struct {
	config 			map[string]string
	lastUpdateTime  int64
	Orders 			map[string] []*model.Order
}

func init() {
	ImportantData = &StaticData{}
	ImportantData.Orders = make(map[string] []*model.Order)
}

func (s *StaticData) GetConfig() map[string]string {
	now := utils.Now()
	if now - s.lastUpdateTime > 10 * 1000 {
		conf := getConfigMap()
		if conf != nil {
			s.config = conf
			s.lastUpdateTime = now
			log.Debug("log config update")
			return s.config
		}
	}
	return  s.config
}

func getConfigMap() (conf map[string]string) {
	file, err := os.Open(config.CfgPath)
	defer file.Close()
	if err == nil {
		//log.Println("load file -> config.toml error", err)
		//return nil, err
		fileByte, err := ioutil.ReadAll(file)
		if err == nil {
			//log.Println("read file -> config.toml error", err)
			//return nil, err
			if _, err := toml.Decode(string(fileByte), &conf); err == nil {
				//log.Println("decode config file -> config.toml error", err)
				//return nil, err
				//core.ImportantData.LastUpdateTime = now
				//core.ImportantData.Config = conf
				return conf
			}
		}
	}
	return nil
}