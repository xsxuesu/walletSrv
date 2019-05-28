package model

import (
	"fmt"
	"testing"
	"walletSrv/coldWallet/config"
)

func TestSerialModel_Insert(t *testing.T) {
	sEntity := SerialEntity{
		SerialNo:"0011",
		CoinType:"btc",
		F:"1mjddfaf",
		T:"3kafkaskdfkasf",
		Value:"dfasfadsf",
		Fee:50000,
		TxId:"001",
		Status:"pending",
	}

	db,err :=config.GetMySqlDb()
	if err != nil {
		fmt.Println(err.Error())
	}

	err = SerialModel{db}.CheckTable("btc")

	if err != nil {
		fmt.Println(err.Error())
	}

	err = SerialModel{db}.Insert(sEntity)

	if err != nil {
		fmt.Println(err.Error())
	}
}

func TestSerialModel_Update(t *testing.T) {
	sEntity := SerialEntity{
		SerialNo:"0011",
		CoinType:"btc",
		F:"1mjddfaf",
		T:"3kafkaskdfkasf",
		Value:"dfasfadsf",
		Fee:50000,
		TxId:"002",
		Status:"exacted",
	}

	db,err :=config.GetMySqlDb()
	if err != nil {
		fmt.Println(err.Error())
	}

	err = SerialModel{db}.CheckTable("btc")

	if err != nil {
		fmt.Println(err.Error())
	}



	err = SerialModel{db}.Update(sEntity)

	if err != nil {
		fmt.Println(err.Error())
	}
}


func TestSerialModel_FindBySerialNo(t *testing.T) {
	db,err :=config.GetMySqlDb()
	if err != nil {
		fmt.Println(err.Error())
	}

	s := SerialModel{db}

	entity,talbename,err := s.FindBySerialNo("0011","btc")

	if err != nil {
		fmt.Println(err.Error())
	}

	fmt.Println(entity.SerialNo,entity)
	fmt.Println(talbename)
}