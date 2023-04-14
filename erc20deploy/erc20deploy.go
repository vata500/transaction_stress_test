package erc20deploy

import (
	"crypto/ecdsa"
	"fmt"
	"log"
	"math/big"

	erc20 "l2_testing_tool/src"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

func Erc20deploy(pvkey string, url string) {
	client, err := ethclient.Dial(url) // 이더리움 노드에 연결합니다.
    if err != nil {
        log.Fatal(err)
    }

    privateKey, err := crypto.HexToECDSA(pvkey) // 배포 계정의 개인 키를 가져옵니다.
    if err != nil {
        log.Fatal(err)
    }

    publicKey := privateKey.Public() // 개인 키에서 공개 키를 추출합니다.
    publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
    if !ok {
        log.Fatal("Error casting public key to ECDSA")
    }

    fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA) // 공개 키에서 주소를 추출합니다.

    // 배포할 ERC20 컨트랙트의 인자를 설정합니다.
    name := "MyToken"
    symbol := "MT"
    decimals := uint8(18)
    totalSupply := big.NewInt(1000000000000000000) // 1,000,000.000000000000000000 MT

    // 배포 계정에서 트랜잭션을 보낼 수 있도록 인증서를 생성합니다.
    auth, err := bind.NewKeyedTransactorWithChainID(privateKey, big.NewInt(1))
    if err != nil {
        log.Fatal(err)
    }

    // ERC20 컨트랙트를 배포합니다.
    address, tx, _, err := erc20.DeployERC20(auth, client, name, symbol, decimals, totalSupply)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Contract deployed to address: %s\n", address.Hex())
    fmt.Printf("Transaction hash: %s\n", tx.Hash().Hex())
}