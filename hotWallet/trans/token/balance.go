package token

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"math/big"
)


func GetContractDecimal(contract string, geth *ethclient.Client)(*big.Int ,error) {
	contractAdd := common.HexToAddress(contract)

	token, err := newTokenCaller(contractAdd, geth)
	if err != nil {
		return big.NewInt(0), err
	}

	decimals, err := token.Decimals(nil)
	if err != nil {
		decimals = big.NewInt(0)
	}
	return decimals,nil
}



func GetContractBlanceInt(contract,wallet string, geth *ethclient.Client)(*big.Int ,error)  {
	contractAdd := common.HexToAddress(contract)
	token, err := newTokenCaller(contractAdd, geth)
	if err != nil {
		return big.NewInt(0), err
	}

	walletAdd := common.HexToAddress(wallet)

	balance, err := token.BalanceOf(nil, walletAdd)
	if err != nil {
		return big.NewInt(0), err
	}

	return balance , nil
}


func GetContractBlanceOf(contract,wallet string, geth *ethclient.Client)(float64 ,error)  {
	contractAdd := common.HexToAddress(contract)
	token, err := newTokenCaller(contractAdd, geth)
	if err != nil {
		return float64(0), err
	}

	decimals, err := token.Decimals(nil)
	if err != nil {
		decimals = big.NewInt(0)
	}

	walletAdd := common.HexToAddress(wallet)

	balance, err := token.BalanceOf(nil, walletAdd)
	if err != nil {
		return float64(0),err
	}

	fbalance := new(big.Float).SetInt(balance)
	if decimals.Int64() == 0{
		balanceVale,_ := fbalance.Float64()
		return balanceVale,nil
	}

	dec := new(big.Float).SetInt(decimals)

	bValue := new(big.Float).Quo(fbalance,dec)

	balanceF ,_ := bValue.Float64()
	return balanceF,nil
}
