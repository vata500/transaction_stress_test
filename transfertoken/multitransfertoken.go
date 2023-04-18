package transfertoken

import (
	"crypto/ecdsa"
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

type Account struct {
	Url string
	Address    string
	PrivateKey string
}

func MultiTransferToken(n int) {
	Accounts := NewAccounts(n)
	fmt.Println(Accounts)

	// for i := 0; i < len(Accounts); i++ {
	// 	account := Accounts[i]
	// 	// 구조체의 변수를 메서드 인자로 사용
	// 	receiveToken(account)
	// 	fmt.Printf("Account %d: %s\n", i, account.Address)
	// 	// ...
	// }


}

func receiveToken(a Account) {
	// TransferErc20token(a, Conf.Transfererctoken, )

}

func NewAccounts(n int) []Account {
	accounts := make([]Account, 0, n)
	pvk, pbk := GenerateKeyPair()
	for i := 1; i <= n; i++ {
        account := Account{Url: Conf.Host.Url, Address: pbk, PrivateKey: pvk}
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