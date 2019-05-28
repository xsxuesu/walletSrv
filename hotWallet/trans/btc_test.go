package trans

import (
	"fmt"
	"testing"
)

func TestGetBalanceInAddress(t *testing.T) {

	rpcclient,err  := NewBtcClient()
	if err != nil {
		fmt.Println(err.Error())
	}

	amount,err := rpcclient.GetBalanceInAddress("2NDkiFRJ627PpEAzKrka1kbcZNq6cBJHqQi")
	//amount,err := GetBalanceInAddress("2NDkiFRJ627PpEAzKrka1kbcZNq6cBJHqQi")
	if err != nil {
		fmt.Println("err:")
		fmt.Println(err.Error())
	}
	fmt.Println(amount)
}


func TestGetUnspentByAddress(t *testing.T) {
	rpcclient,err  := NewBtcClient()
	if err != nil {
		fmt.Println(err.Error())
	}

	relist,err := rpcclient.GetUnspentByAddress("2NDkiFRJ627PpEAzKrka1kbcZNq6cBJHqQi")
	if err != nil {
		fmt.Println("err:")
		fmt.Println(err.Error())
	}
	fmt.Println(relist)

}


func TestCreateBtcTrans(t *testing.T){
	rpcclient,err  :=  NewBtcClient()
	if err != nil {
		fmt.Println(err.Error())
	}
	reqeust,err := rpcclient.CreateBtcTrans("000999","2NDkiFRJ627PpEAzKrka1kbcZNq6cBJHqQi","2N3GXgBW3rhTrtmkpkxZC4AZcvJQabpS6Wr",1.4)
	if err != nil {
		fmt.Println("err:")
		fmt.Println(err.Error())
	}
	fmt.Println(reqeust)
}


func TestSentBtc(t *testing.T) {
	rpcclient,err  :=  NewBtcClient()
	if err != nil {
		fmt.Println(err.Error())
	}

	account,err := rpcclient.CreateAccount("lidepeng")
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(account)

}

func TestBtcClient_InsertToAccount(t *testing.T) {
	rpcclient,err  :=  NewBtcClient()
	if err != nil {
		fmt.Println(err.Error())
	}

	acc,err := rpcclient.InsertToAccount("2NDkiFRJ627PpEAzKrka1kbcZNq6cBJHqQi","lidepeng")
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(acc)
}

func TestBtcClient_GetBlockCount(t *testing.T) {
	rpcclient,err  :=  NewBtcClient()
	if err != nil {
		fmt.Println(err.Error())
	}
	acc,err := rpcclient.GetBlockCount()
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(acc)
}