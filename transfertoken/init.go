package transfertoken

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"golang.org/x/crypto/sha3"
)

func ReceiveToken(h Host, t Transfererctoken, receiveAddr string){
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

	value := big.NewInt(0)

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	toAddress := common.HexToAddress(receiveAddr)

	tokenAddress := common.HexToAddress(Conf.Transfererctoken.Tokenaddress)

	transferFnSignature := []byte("transfer(address,uint256)")
	hash := sha3.NewLegacyKeccak256()
	hash.Write(transferFnSignature)
	methodID := hash.Sum(nil)[:4]
	// fmt.Println(hexutil.Encode(methodID)) // 0xa9059cbb

	paddedAddress := common.LeftPadBytes(toAddress.Bytes(), 32)
	// fmt.Println(hexutil.Encode(paddedAddress)) 

	// toStringValue := fmt.Sprintf("%f",t.Value)
	// amount := new(big.Int)
	// amount.SetString("10000", 10) // sets the value to 1000 tokens, in the token denomination

	toStringValue := fmt.Sprintf("%f",t.Value)
	amount := new(big.Int)
	amount.SetString(toStringValue, 10) // sets the value to 1000 tokens, in the token denomination
	fmt.Println(amount)

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
	fmt.Printf("tx sent: %s, tx sender: %s\n", signedTx.Hash().Hex(), fromAddress)
}

func ReceiveETH(h Host, toAddress string, value float64)(string, error){
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

	// 받을 주소
	addr := common.HexToAddress(toAddress)

	// 전송할 이더 양 설정
	v := ConvertToWei(value)	

	// data 값 설정 (일반적으로 이 값은 빈 바이트 배열입니다)
	var data []byte

	// chainID 받아오기
	chainID, err := client.NetworkID(context.Background())
	if err != nil {
  		log.Fatal(err)
	}	

	// 출력big.NewInt(0)
	fmt.Printf("\n##### Nonce: %d, toAddress: %s, value: %f, chainID: %d #####\n", nonce, toAddress, value, chainID)

	// 새로운 전송 트랜잭션 생성
	tx := types.NewTransaction(nonce, addr, v, uint64(21000), big.NewInt(0), data)

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
