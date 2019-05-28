package trans

import (
	"fmt"
	"testing"
)

func TestSendByKS(t *testing.T) {
	txid,err := SendByKS("0x7cc2e4558ca66e80b3a47ebcd1abc2a87e5aa55f","0x06A98EBC3E9aae240407dD15c7bA91b137eB8F8F","","password",0.1)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(txid)
}

func TestGetContractBlncRPC(t *testing.T) {
	v,e:=GetContractBlncRPC("0xbf5f8bfcee9502a30018d91c63eca66980e6e9bb","0x6D0ebc3C87d8B2E2D47Fd382A123Cd77ed9937b2")
	if e != nil {
		fmt.Println(e.Error())
	}
	fmt.Println(v)
}

func TestGetBalance(t *testing.T) {
	v,e := GetBalance("0x6D0ebc3C87d8B2E2D47Fd382A123Cd77ed9937b2")
	if e != nil {
		fmt.Println(e.Error())
	}
	fmt.Println(v)
}