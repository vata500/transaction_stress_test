package transfertoken

import (
	"crypto/ecdsa"
	"fmt"
	"log"
	"time"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

func MultiTransferToken(n int) {
	Accounts := NewAccounts(n)
	fmt.Println(Accounts)

	for i := 0; i < len(Accounts); i++ {
		account := Accounts[i]
		receive(account) // Host account로 부터 erc20 토큰 수령 
		fmt.Printf("Account %d: %s\n", i, account.Address)
	}

	MultiStart(Accounts)
}

func MultiStart(a []Host){

	n := len(a)

	// ticker := time.NewTicker(time.Duration(Conf.Transfererctoken.Interval*1000) * time.Millisecond)
	done := make(chan bool)

	for i := 0; i < n; i++ {
		host := a[i]
		go func() {
			for {
				select {
				case <-done:
					return
				default:
					TransferErc20token(host, Conf.Transfererctoken)
					time.Sleep(time.Duration(Conf.Transfererctoken.Interval*1000) * time.Millisecond)
				}
			}
		}()
	}

	time.Sleep(time.Duration(Conf.Transfererctoken.Minute) * time.Minute)
	// ticker.Stop()
	done <- true

	fmt.Println("erc20 token transfer stopped")
}

func receive(a Host) {
	receiveAddress := a.Address
	ReceiveETH(Conf.Host, receiveAddress, 100)
	ReceiveToken(Conf.Host, Conf.Transfererctoken, receiveAddress)
}

func NewAccounts(n int) []Host {
	accounts := make([]Host, 0, n)
	for i := 1; i <= n; i++ {
		pvk, pbk := GenerateKeyPair()
        account := Host{Url: Conf.Host.Url, Address: pbk, PrivateKey: pvk}
        accounts = append(accounts, account)
    }
	return accounts
}

func GenerateKeyPair() (string, string) {
   // 키 생성
   privateKey, err := crypto.GenerateKey()
   if err != nil {
	   log.Fatal(err)
   }

   // 개인 키 출력
   privateKeyBytes := crypto.FromECDSA(privateKey)
   fmt.Println("Private Key:", hexutil.Encode(privateKeyBytes)[2:])
   validPrivateKey := hexutil.Encode(privateKeyBytes)[2:]


   // 공개 키 추출
   pvK, err := crypto.HexToECDSA(validPrivateKey)
   if err != nil {
	   log.Fatal(err)
   }

   publicKey := pvK.Public()
   publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
   if !ok {
	   log.Fatal("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
   }
   Address := crypto.PubkeyToAddress(*publicKeyECDSA)
   AddressStr := Address.Hex()

   return validPrivateKey, AddressStr
}