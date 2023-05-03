package transfertoken

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/BurntSushi/toml"
)

type Host struct {
	Url string `toml:"url"`
	Address    string `toml:"address"`
	PrivateKey string `toml:"privatekey"`
}

type Transfererctoken struct {
	Value 	float64 `toml:"value"`
	Interval 	float64 `toml:"interval"`
	Minute 	int `toml:"minute"`
	Tokenaddress 	string `toml:"tokenaddress"`
	Log_path string `toml:"log_path"`
	Accounts int `toml:"accounts"`
}

type Config struct {
	Host    Host    `toml:"host"`
	Transfererctoken     Transfererctoken     `toml:"transfererctoken"`
}

var Conf Config
var Total_tx = 0

func Start(checkStartTime time.Time){
	createTimeNowFile(checkStartTime)

	// config.toml create instance
	if _, err := toml.DecodeFile("config.toml", &Conf); err != nil {
		log.Fatal(err)
	}

	// account가 있을 경우, go routine으로 병렬 트랜잭션 처리
	if Conf.Transfererctoken.Accounts > 0 {
		MultiTransferToken(Conf.Transfererctoken.Accounts)
	} else {
		ticker := time.NewTicker(time.Duration(Conf.Transfererctoken.Interval*1000) * time.Millisecond)
		done := make(chan bool)

		go func() {
		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				TransferErc20token(Conf.Host, Conf.Transfererctoken)
			}
		}
		}()

		time.Sleep(time.Duration(Conf.Transfererctoken.Minute) * time.Minute)
		ticker.Stop()
		done <- true

		fmt.Println("erc20 token transfer stopped")
	}
	result(checkStartTime)
}


func createTimeNowFile(checkStartTime time.Time){
	f, err := os.Create("timeNow")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	_, err = fmt.Fprintf(f, "%s", checkStartTime.Format("2006-01-02 15:04:05"))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println()
}

func result(startTime time.Time){
	Testing_Time := Conf.Transfererctoken.Minute *60
	endTime := time.Now()
	duration := endTime.Sub(startTime)
	Average_TPS := float64(Total_tx) / duration.Seconds()
	fmt.Printf("Total_Tx : %d, Average_TPS : %f, Testing_Time : %ds\n", Total_tx, Average_TPS, Testing_Time)
}
