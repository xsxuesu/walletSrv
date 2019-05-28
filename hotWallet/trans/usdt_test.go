package trans

import (
	"fmt"
	"testing"
)

//测试服务器 地址
//mhuGL7qML9U4pgBvLVhd1xrgAEuSxZWiA7
//muFB6TntE7Y2QYcpW6HPaxaEDKSdofic6k
func TestGetUsdtBalanceByAddr(t *testing.T) {

}

func TestSendUsdt(t *testing.T) {
	//err := SendUsdt("1","2",float64(1.2))
	//if err != nil {
	//	fmt.Print(err.Error())
	//}
}

func TestUSDTClient_GetAddress(t *testing.T) {
	usdtClient,err := NewUSDTClient(31)
	if err != nil {
		fmt.Println(err.Error())
	}
	addr,err := usdtClient.GetAddress("lidepeng")
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(addr)
}