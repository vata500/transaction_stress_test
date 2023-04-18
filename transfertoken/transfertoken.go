package transfertoken

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"math/big"
	"os"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"golang.org/x/crypto/sha3"
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
		result(checkStartTime)
	}
}



func TransferErc20token(h Host, t Transfererctoken, receiveAddr ...string) {
	client, err := ethclient.Dial(h.Url)
	if err != nil {
		log.Fatal(err)
	}

	privateKey, err := crypto.HexToECDSA(h.PrivateKey)
	if err != nil {
		log.Fatal(err)
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatal(err)
	}

	value := big.NewInt(0) // in wei (0 eth)
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	var toAddress common.Address

	// account 인자로 전달받았을 경우, toAddress에 account address를 삽입.
	if len(receiveAddr) > 0 {
		toAddress = common.HexToAddress(receiveAddr[0])
	} else {
		toAddress = common.HexToAddress(h.Address)
	}
	tokenAddress := common.HexToAddress(Conf.Transfererctoken.Tokenaddress)

	transferFnSignature := []byte("transfer(address,uint256)")
	hash := sha3.NewLegacyKeccak256()
	hash.Write(transferFnSignature)
	methodID := hash.Sum(nil)[:4]
	// fmt.Println(hexutil.Encode(methodID)) // 0xa9059cbb

	paddedAddress := common.LeftPadBytes(toAddress.Bytes(), 32)
	// fmt.Println(hexutil.Encode(paddedAddress)) 

	toStringValue := fmt.Sprintf("%f",t.Value)
	amount := new(big.Int)
	amount.SetString(toStringValue, 10) // sets the value to 1000 tokens, in the token denomination

	paddedAmount := common.LeftPadBytes(amount.Bytes(), 32)
	// fmt.Println(hexutil.Encode(paddedAmount))

	var data []byte
	data = append(data, methodID...)
	data = append(data, paddedAddress...)
	data = append(data, paddedAmount...)

	gasLimit, err := client.EstimateGas(context.Background(), ethereum.CallMsg{
		To:   &tokenAddress,
		Data: data,
	})
	if err != nil {
		log.Fatal(err)
	}
	// fmt.Println(gasLimit) // 23256

	tx := types.NewTransaction(nonce, tokenAddress, value, gasLimit, gasPrice, data)

	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		log.Fatal(err)
	}

	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		log.Fatal(err)
	}
	Total_tx += 1
	fmt.Printf("tx sent: %s\n", signedTx.Hash().Hex())
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
