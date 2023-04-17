package main

import (
	"log"

	"github.com/BurntSushi/toml"
)

type Host struct {
	Url string `toml:"url"`
	Address    string `toml:"address"`
	PrivateKey string `toml:"privatekey"`
}

type Tx struct {
	Value int `toml:"value"`
	Interval      int `toml:"interval"`
	Minute int     `toml:"minute"`
}

type Transfertoken struct {
	Value 	int `toml:"value"`
	Interval	 int `toml:"inteval"`
	Minute 	int `toml:"minute"`
	Tokenaddress 	string `toml:"tokenaddress"`
}

type Config struct {
	Host    Host    `toml:"host"`
	Tx      Tx      `toml:"tx"`
	Transfertoken Transfertoken `toml:"transfertoken"`
}

var conf Config

func main() {

	if _, err := toml.DecodeFile("config.toml", &conf); err != nil {
		log.Println("hey! let's create config.toml")
		log.Fatal(err)
	}

	
 	//erc20deploy.Erc20deploy(conf.Host.PrivateKey, conf.Host.Url)
	//sendtx.SendEth(conf.Tx.TransactionPerSecond, conf.Tx.Time)

}

