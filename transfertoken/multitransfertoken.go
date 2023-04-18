package transfertoken

import (
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha3"
	"encoding/hex"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

type Account struct {
	Url string
	Address    string
	PrivateKey string
}

func MultiTransferToken(n int) {
	NewAccounts(n)
}

func NewAccounts(n int) []Account {
	accounts := make([]Account, 0, n)
	for i := 1; i <= n; i++ {
		pvk, err := GenerateKeyPair()
		if err != nil {
			panic(err)
		}

		fmt.Printf("pvk: %s\n", pvk)
		// fmt.Printf("pvk: %s, pbk: %s: \n", pvk, pbk)
        // account := Account{Url: Conf.Host.Url, Address: pbk, PrivateKey: pvk}
        // accounts = append(accounts, account)
    }
	return accounts
}


func GenerateKeyPair() (common.Address, error) {
    privateKey, _ := ecdsa.GenerateKey(crypto.S256(), rand.Reader)
	Pvkey := PrivateKeyToAddress(privateKey)
	return Pvkey, nil
}

func PrivateKeyToAddress(privateKey *ecdsa.PrivateKey) common.Address {
    publicKey := privateKey.Public()
    publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
    if !ok {
        return common.Address{}
    }
    publicKeyBytes := append(publicKeyECDSA.X.Bytes(), publicKeyECDSA.Y.Bytes()...)
    hash := sha3.NewLegacyKeccak256()
    hash.Write(publicKeyBytes)
    address := hash.Sum(nil)[12:]
    return common.BytesToAddress(address)
}