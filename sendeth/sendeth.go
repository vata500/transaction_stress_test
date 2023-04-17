package sendeth

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)


type Host struct {
	Url string `toml:"url"`
	Address    string `toml:"address"`
	PrivateKey string `toml:"privatekey"`
}

type Sendeth struct {
	Value float64 `toml:"value"`
	Interval      float64 `toml:"interval"`
	Minute int     `toml:"minute"`
}

type Config struct {
	Host    Host    `toml:"host"`
	Sendeth      Sendeth      `toml:"sendeth"`
}

var conf Config

func SendEther(){
	if _, err := toml.DecodeFile("config.toml", &conf); err != nil {
		log.Println("hey! let's create config.toml")
		log.Fatal(err)
	}
	
	ticker := time.NewTicker(time.Duration(conf.Sendeth.Interval*1000) * time.Millisecond)
	done := make(chan bool)

	go func() {
		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				SendTransaction(conf.Host, conf.Sendeth)
				fmt.Println("Run your method here...")
			}
		}
	}()

	time.Sleep(time.Duration(conf.Sendeth.Minute) * time.Minute)
	ticker.Stop()
	done <- true
	fmt.Println("Ticker stopped.")
}


func SendTransaction(h Host, s Sendeth)(string, error){
	// RPC URL 지정
	client, err := ethclient.Dial(h.Url)
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum client: %v", err)
	}

	// private key를 지정합니다.
	privateKey, err := crypto.HexToECDSA(h.PrivateKey)
	if err != nil {
		log.Fatalf("Failed to retrieve private key: %v", err)
	}

	// 계정 주소 가져오기
	fromAddress := crypto.PubkeyToAddress(privateKey.PublicKey)

	// nonce 값 가져오기
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatalf("Failed to retrieve nonce: %v", err)
	}

	// 수신자 주소 가져오기
	toAddress := common.HexToAddress(h.Address)

	// 전송할 이더 양 설정
	value := ConvertToWei(s.Value)	

	// data 값 설정 (일반적으로 이 값은 빈 바이트 배열입니다)
	var data []byte


	// chainID 받아오기
	chainID, err := client.NetworkID(context.Background())
	if err != nil {
  		log.Fatal(err)
	}	

	// 출력big.NewInt(0)
	fmt.Printf("\n##### Nonce: %d, toAddress: %s, value: %d, chainID: %d #####\n", nonce, toAddress, value, chainID)

	// 새로운 전송 트랜잭션 생성
	tx := types.NewTransaction(nonce, toAddress, value, uint64(21000), big.NewInt(0), data)

	// 서명합니다.
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		log.Fatalf("Failed to sign transaction: %v", err)
	}

	// 트랜잭션을 블록체인에 전송합니다.
	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		log.Fatalf("Failed to send transaction: %v", err)
	}

	// 트랜잭션 해시 출력
	txHash := signedTx.Hash().Hex()
	fmt.Printf("Transaction sent: %s", txHash)

	return txHash, nil
}

func ConvertToWei(ether float64) *big.Int {
	wei := big.NewInt(0)
	wei.Exp(big.NewInt(10), big.NewInt(18), nil) // 10^18
	wei.Mul(wei, big.NewInt(int64(ether)))

	return wei
}