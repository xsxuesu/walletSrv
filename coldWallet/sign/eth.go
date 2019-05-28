package sign

import (
	"encoding/hex"
	"fmt"
	"golang.org/x/crypto/sha3"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
)

func FloatToBigInt(val float64) *big.Int {
	bigval := new(big.Float)
	bigval.SetFloat64(val)

	coin := new(big.Float)
	coin.SetInt(big.NewInt(1000000000000000000))

	bigval.Mul(bigval, coin)

	result := new(big.Int)
	bigval.Int(result)

	return result
}

func EthSignTx(prik string,nonce uint64,to string,value *big.Int,
	gasLimit uint64,gasPrice int64,chainID int64,contract string)(string ,error){

	privateKey,err := crypto.HexToECDSA(prik)
	if err != nil {
		return "",err
	}

	toAddress := common.HexToAddress(to)

	gasPriceValue := big.NewInt(gasPrice)

	var data []byte
	var contractAddr common.Address

	tx := types.NewTransaction(nonce, toAddress, value, gasLimit, gasPriceValue, data)


	// 转账到合约里
	if contract != ""{
		contractAddr = common.HexToAddress(contract)
		fmt.Println("fmt.Println(contract):")
		fmt.Println(contract)

		transferFnSignature := []byte("transfer(address,uint256)")
		hash := sha3.NewLegacyKeccak256()
		hash.Write(transferFnSignature)
		methodID := hash.Sum(nil)[:4]

		paddedAddress := common.LeftPadBytes(toAddress.Bytes(), 32)

		paddedAmount := common.LeftPadBytes(value.Bytes(), 32)

		data = append(data, methodID...)
		data = append(data, paddedAddress...)
		data = append(data, paddedAmount...)

		tx = types.NewTransaction(nonce, contractAddr, big.NewInt(0), gasLimit, gasPriceValue, data)
	}

	chainid := big.NewInt(chainID)

	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainid), privateKey)
	if err != nil {
		return "",err
	}

	ts := types.Transactions{signedTx}
	rawTxBytes := ts.GetRlp(0)
	rawTxHex := hex.EncodeToString(rawTxBytes)

	return rawTxHex,nil
}