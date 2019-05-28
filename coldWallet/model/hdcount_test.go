package model

import (
	"fmt"
	"strconv"
	"testing"
	"walletSrv/coldWallet/config"
)

func TestHdCountModel_InsertHdCount(t *testing.T) {
	hd := HdCount{
		"btc",
		"0",
	}
	db,err :=config.GetMySqlDb()
	if err != nil {
		fmt.Println(err.Error())
	}

	for i := 0; i < 1000; i++ {
		err = HdCountModel{db}.UpdateHdCount(hd.Cointype,strconv.Itoa(i))
		if err != nil {
			fmt.Println(err.Error())
		}
	}

}

func TestHdCountModel_FindByType(t *testing.T) {
	db,err :=config.GetMySqlDb()
	if err != nil {
		fmt.Println(err.Error())
	}
	count,err := HdCountModel{db}.FindByType("btc")
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(count)
}