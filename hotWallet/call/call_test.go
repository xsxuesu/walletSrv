package call

import (
	"fmt"
	"testing"
)

func TestCallInsertAddress(t *testing.T) {
	client := InitClient()
	//CallInsertAddress("eth","0x06A98EBC3E9aae240407dD15c7bA91b137eB8F8F","a80626b43fb13f8b42ec372d6ca9fb7c973727a4319d069dbfaa8a367c30bc68")

	//CallInsertAddress("btc","mz5otQWaaPmWEydGw9gjfaHxSKhHX6Fgj5","cS8Epp2J8Bf3Bi5USVWxTLTdGNziBTR8Qo1dFgYnbaLt8LvMyMyY")
	//2N3GXgBW3rhTrtmkpkxZC4AZcvJQabpS6Wr
	//cTfNbYsdvRzJw37c1kB8qzsz1eYqpUrAcEQ4d1JPTq9kcNpjKTjd
	resp,err := client.CallInsertAddress("btc","2NDkiFRJ627PpEAzKrka1kbcZNq6cBJHqQi","cN4z7iLaxMG5MLgENq2eWKqdU9MChjRhotLwjHsgh9btdPpzqUr7")
	if err != nil {
		fmt.Println("err:")
		fmt.Println(err.Error())
	}

	fmt.Println(resp)
}

func TestCallSignEthTx(t *testing.T) {
	client := InitClient()

	//trans.InitClient()

	client.CallSignEthTx("10014","0x76A4Bf011d91a543Ff6eee381e69304e0182044E","0x06A98EBC3E9aae240407dD15c7bA91b137eB8F8F",float64(1.0))
}
//testnet.qtornado.com:51002
func TestCallSignBtcTx(t *testing.T) {
	client := InitClient()
	resp,err := client.CallSignBtcTx("2000131","2NDkiFRJ627PpEAzKrka1kbcZNq6cBJHqQi","2N3GXgBW3rhTrtmkpkxZC4AZcvJQabpS6Wr",1.3000)
	if err != nil {
		fmt.Println("err:")
		fmt.Println(err.Error())
	}
	fmt.Println(resp)
}
