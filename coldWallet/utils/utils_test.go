package utils

import (
	"fmt"
	"testing"
)

func TestSha224Str(t *testing.T) {
	x:=[]byte("少壮不努力，一生在朝鲜。英语学不牢，世代在青藏。")
	fmt.Println(Sha256Str(x))
	x2:=Sha224Str(x)
	fmt.Println(x2)
	fmt.Println(len(x2))
	x3:=Sha384Str(x)
	fmt.Println(x3)
	fmt.Println(len(x3))
	x4:=Sha512Str(x)
	fmt.Println(x4)
	fmt.Print(len(x4))


	str := "0x6B7f920b4c4bc8e94B212aE4b1241e66915F3b50"
	fmt.Println("JSHash:")
	fmt.Println(JSHash(str),"len:",len(JSHash(str)))

	fmt.Println("RSHash:")
	fmt.Println(RSHash(str),"len:",len(RSHash(str)))
	//fmt.Println("PJWHash:")
	//fmt.Println(PJWHash(str))
	fmt.Println("BKDRHash:")
	fmt.Println(BKDRHash(str),"len:",len(BKDRHash(str)))
	fmt.Println("SDBMHash:")
	fmt.Println(SDBMHash(str),"len:",len(SDBMHash(str)))
	fmt.Println("DJBHash:")
	fmt.Println(DJBHash(str),"len:",len(DJBHash(str)))

	fmt.Println("DEKHash:")
	fmt.Println(DEKHash(str),"len:",len(DEKHash(str)))
	fmt.Println("APHash:")
	fmt.Println(APHash(str),"len:",len(APHash(str)))





}


func TestGetHashName(t *testing.T) {
	hashName,err  := GetHashName("0xf68ff5E431Aa98E76452Be8F7b8743cD8138de6D")
	if err != nil {
		fmt.Println(err.Error())
	}
	key,err := HashStr(hashName,"0xf68ff5E431Aa98E76452Be8F7b8743cD8138de6D")
	if err != nil {
		fmt.Println(err.Error())
	}
	//L2MbWPuRxQdLM5jMqLcSMwiMpJG6dRxoLmEbUP5gksuv8W4qDKBY
	//L2MbWPuRxQdLM5jMqLcSMwiMpJG6dRxoLmEbUP5gksuv8W4qDKBY
	//0xe341116e72b3c6e4878e9cc50ac4a26de36f80da8389381ca9e5cd6562be7107
	//0xe341116e72b3c6e4878e9cc50ac4a26de36f80da8389381ca9e5cd6562be7107
	encodePri,err :=  EncodePri(key,"0xe341116e72b3c6e4878e9cc50ac4a26de36f80da8389381ca9e5cd6562be7107")
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println("fmt.Println(encodePri):")
	fmt.Println(encodePri)

	decodePri,err := DecodePri(key,encodePri)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println("fmt.Println(decodePri):")
	fmt.Println(decodePri)
}