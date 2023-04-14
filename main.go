package main

import (
	"l2_testing_tool/erc20deploy"
	"log"

	"github.com/BurntSushi/toml"
)

type Host struct {
	Url string `toml:"url"`
	Address    string `toml:"address"`
	PrivateKey string `toml:"privatekey"`
}

type Tx struct {
	Eth_value int `toml:"eth_value"`
	TransactionPerSecond      int `toml:"transactionPerSecond"`
	Time int     `toml:"time"`
}

type Config struct {
	Host    Host    `toml:"host"`
	Tx      Tx      `toml:"tx"`
}

var conf Config

func main() {

	if _, err := toml.DecodeFile("config.toml", &conf); err != nil {
		log.Println("hey! let's create config.toml")
		log.Fatal(err)
	}

	
 	erc20deploy.Erc20deploy(conf.Host.PrivateKey, conf.Host.Url)
	//sendtx.SendEth(conf.Tx.TransactionPerSecond, conf.Tx.Time)

}

