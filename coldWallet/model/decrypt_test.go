package model

import (
	"fmt"
	"testing"
	"walletSrv/coldWallet/config"
)

func TestDecryptModel_Insert(t *testing.T) {
	s := DecryptEntity{
		SerialNo:"1001",
		CoinType:"eth",
		F:"0xoodfoasf",
		T:"0xdasfasfasf",
		HashFun:"add",
	}

	db,err :=config.GetMySqlDb()
	if err != nil {
		fmt.Println(err.Error())
	}

	decryptModel := DecryptModel{db}
	err = decryptModel.CheckTable(s.CoinType)
	if err != nil {
		fmt.Println(err.Error())
	}

	err = decryptModel.Insert(s)

	if err != nil {
		fmt.Println(err.Error())
	}
}

func TestDecryptModel_Update(t *testing.T) {
	s := DecryptEntity{
		SerialNo:"1001",
		CoinType:"eth",
		F:"0xoodfoasf11111111",
		T:"0xdasfasfasf2222",
		HashFun:"add",
	}

	db,err :=config.GetMySqlDb()
	if err != nil {
		fmt.Println(err.Error())
	}

	decryptModel := DecryptModel{db}

	err = decryptModel.Update(s)

	if err != nil {
		fmt.Println(err.Error())
	}
}

func TestDecryptModel_FindBySerialNo(t *testing.T) {
	db,err :=config.GetMySqlDb()
	if err != nil {
		fmt.Println(err.Error())
	}

	decryptModel := DecryptModel{db}

	s := DecryptEntity{
		SerialNo:"1001",
		CoinType:"eth",
		F:"0xoodfoasf11111111",
		T:"0xdasfasfasf2222",
		HashFun:"add",
	}



	decrypt,tablename,err := decryptModel.FindByDecryptNo(s.SerialNo,s.CoinType)

	//t,tablename,err := decryptModel.FindByDecryptNo(s.SerialNo,s.CoinType)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(decrypt.HashFun)
	fmt.Println(tablename)
	fmt.Println(s.SerialNo)
}