package config

import (
	"os"
	"log"
	"io/ioutil"

	"github.com/BurntSushi/toml"
)

type Config struct {
	LogPath string `toml:"log_path"`
	HuobiKey string `toml:"huobi_key"`
	HuobiSecret string `toml:"huobi_secret"`
}

const CfgPath = "/usr/jjw/config.toml"

func GetConfig() (conf *Config, err error) {
	file, err := os.Open(CfgPath)
	defer file.Close()
	if err != nil {
		log.Println("load file -> config.toml error", err)
		return nil, err
	}
	fileByte, err := ioutil.ReadAll(file)
	if err != nil {
		log.Println("read file -> config.toml error", err)
		return nil, err
	}
	if _, err := toml.Decode(string(fileByte), &conf); err != nil {
		log.Println("decode config file -> config.toml error", err)
		return nil, err
	}
	return conf, err
}