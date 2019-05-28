package model

import (
	"fmt"
	"testing"
	"walletSrv/coldWallet/config"
)

func TestAddrModel_Insert(t *testing.T) {
	add := AddressEntity{
		"test1",
		"testpri",
	}

	db,_ := config.GetMySqlDb()

	addrmodel := AddrModel{db}
	err := addrmodel.Insert("btc",false,add)
	if err != nil {
		fmt.Println("insert:",err.Error())
	}

	add2 := AddressEntity{
		"test2",
		"testpri",
	}
	err = addrmodel.Insert("btc",true,add2)
	if err != nil {
		fmt.Println("hd insert:",err.Error())
	}

	fadd ,err := addrmodel.FindByAddress("btc",add.Address)
	if err != nil {
		fmt.Println("find:",err.Error())
	}
	fmt.Println(fadd.Address,fadd.PrivateKey)

	fadd2 ,err := addrmodel.FindByAddress("btc",add2.Address)
	if err != nil {
		fmt.Println("find:",err.Error())
	}
	fmt.Println(fadd2.Address,fadd2.PrivateKey)

	//for i := 0; i < 10000; i++ {
	//	add.PrivateKey = strconv.Itoa(i)
	//	add2.PrivateKey = strconv.Itoa(i)
	//	err = addrmodel.Update("btc",false,add)
	//	if err != nil {
	//		fmt.Println(err.Error())
	//	}
	//	err = addrmodel.Update("btc",true,add2)
	//	if err != nil {
	//		fmt.Println(err.Error())
	//	}
	//}

}