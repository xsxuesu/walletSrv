package services

import (
	"fmt"
	"testing"
)

func TestCreateHdAddress(t *testing.T)  {
	//addr , err := createHdAddress("btc")
	//if err != nil {
	//	fmt.Println(err.Error())
	//}
	//fmt.Println(addr)
}


func TestInsertAdd(t *testing.T){
	// 0x2BA9B41a7Fb983FfB3B8BEA90D04909DD484dBB6
	// 0x2d7fb66933f68595903f1036dd3d48d657fde16906c951177b121ce3dc5df5ad
	//err := insertAddress("eth","0x2BA9B41a7Fb983FfB3B8BEA90D04909DD484dBB6","0x2d7fb66933f68595903f1036dd3d48d657fde16906c951177b121ce3dc5df5ad")
	//if err != nil {
	//	fmt.Println(err.Error())
	//}
}

func TestPriCheckAddress(t *testing.T) {
	b:=PriCheckAddress("11bccf45b90d2e98240cf757140c25d324f826e18da6084aea6111168289af3957","0xb407EE5Af76d7ccde67cc8dE2fC15Bf621F8d923")

	fmt.Println(b)
}

func TestBtcPriCheckAddress(t *testing.T) {
	b:=BtcPriCheckAddress("cQYJaj2g7uFdtfqEd8HEsTNGAyHb8y8KJvDZF7cbFKKQoEcQupUp","mnLHNUKLCWDdWKzLFZS8LM5Uv87KRLyUss")

	fmt.Println(b)
}