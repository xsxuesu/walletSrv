package utils

import (
	"bytes"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
	"walletSrv/coldWallet/model"
)


const(
	RS="RSHash"
	JS="JSHash"
	BKDR="BKDRHash"
	SDBM="SDBMHash"
	DJB="DJBHash"
	DEK="DEKHash"
	AP="APHash"
)

func PKCS7Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext) % blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func PKCS7UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}



func Sha256Str(x []byte) string {
	y:=sha256.Sum256(x)
	return hex.EncodeToString(y[:])
}

func Sha256Str2(x []byte) string {
	y:=sha256.Sum256(x)
	return fmt.Sprintf("%x",y)
}

func Sha256Str3(x []byte) string {
	myhash:=sha256.New()
	io.WriteString(myhash,string(x))
	y:=myhash.Sum(nil)
	return hex.EncodeToString(y)
}

func Sha256Str4(x []byte) string {
	myhash:=sha256.New()
	myhash.Write(x)
	y:=myhash.Sum(nil)
	return hex.EncodeToString(y)
}

func Sha224Str(x []byte) string {
	myhash:=sha256.New224()
	myhash.Write(x)
	y:=myhash.Sum(nil)
	return hex.EncodeToString(y)
}

func Sha512Str(x []byte) string {
	myhash:=sha512.New()
	myhash.Write(x)
	y:=myhash.Sum(nil)
	return hex.EncodeToString(y)
}

func Sha384Str(x []byte) string {
	myhash:=sha512.New384()
	myhash.Write(x)
	y:=myhash.Sum(nil)
	return hex.EncodeToString(y)
}



func RSHash(str string) string {
	b := 378551
	a := 63689
	hash := uint64(0)
	for i := 0; i < len(str); i++ {
		hash = hash*uint64(a) + uint64(str[i])
		a = a * b
	}
	return fmt.Sprintf("%v",hash)
}

func JSHash(str string) string {
	hash := uint64(1315423911)
	for i := 0; i < len(str); i++ {
		hash ^= ((hash << 5) + uint64(str[i]) + (hash >> 2))
	}
	return fmt.Sprintf("%v",hash)
}

func PJWHash(str string) string{
	BitsInUnsignedInt := (uint64)(4 * 8)
	ThreeQuarters := (uint64)((BitsInUnsignedInt * 3) / 4)
	OneEighth := (uint64)(BitsInUnsignedInt / 8)
	HighBits := (uint64)(0xFFFFFFFF) << (BitsInUnsignedInt - OneEighth)
	hash := uint64(0)
	test := uint64(0)
	for i := 0; i < len(str); i++ {
		hash = (hash << OneEighth) + uint64(str[i])
		if test = hash & HighBits; test != 0 {
			hash = ((hash ^ (test >> ThreeQuarters)) & (^HighBits))
		}
	}
	return fmt.Sprintf("%v",hash)
}

func BKDRHash(str string)string {
	seed := uint64(131) // 31 131 1313 13131 131313 etc..
	hash := uint64(0)
	for i := 0; i < len(str); i++ {
		hash = (hash * seed) + uint64(str[i])
	}
	return fmt.Sprintf("%v",hash)
}

func SDBMHash(str string) string{
	hash := uint64(0)
	for i := 0; i < len(str); i++ {
		hash = uint64(str[i]) + (hash << 6) + (hash << 16) - hash
	}
	//fmt.Printf("SDBMHash %v\n", hash)
	return fmt.Sprintf("%v",hash)
}

func DJBHash(str string) string {
	hash := uint64(0)
	for i := 0; i < len(str); i++ {
		hash = ((hash << 5) + hash) + uint64(str[i])
	}
	//fmt.Printf("DJBHash %v\n", hash)
	return fmt.Sprintf("%v",hash)
}

func DEKHash(str string) string{
	hash := uint64(len(str))
	for i := 0; i < len(str); i++ {
		hash = ((hash << 5) ^ (hash >> 27)) ^ uint64(str[i])
	}
	//fmt.Printf("DEKHash %v\n", hash)
	return fmt.Sprintf("%v",hash)
}

func APHash(str string)string {
	hash := uint64(0xAAAAAAAA)
	for i := 0; i < len(str); i++ {
		if (i & 1) == 0 {
			hash ^= ((hash << 7) ^ uint64(str[i])*(hash>>3))
		} else {
			hash ^= (^((hash << 11) + uint64(str[i]) ^ (hash >> 5)))
		}
	}
	//fmt.Printf("APHash %v\n", hash)
	return fmt.Sprintf("%v",hash)
}
//RS="RSHash"
//JS="JSHash"
//BKDR="BKDRHash"
//SDBM="SDBMHash"
//DJB="DJBHash"
//DEK="DEKHash"
//AP="APHash"
func substring(source string, start int, end int) string {
	var r = []rune(source)
	length := len(r)

	if start < 0 || end > length || start > end {
		return ""
	}

	if start == 0 && end == length {
		return source
	}

	return string(r[start : end])
}

//RS="RSHash"
//JS="JSHash"
//BKDR="BKDRHash"
//SDBM="SDBMHash"
//DJB="DJBHash"
//DEK="DEKHash"
//AP="APHash"

func EncodePriByPub(addr string,pirk string)(string,error){
	hashName,err  := GetHashName(addr)
	if err != nil {
		return "",err
	}
	key,err := HashStr(hashName,addr)
	if err != nil {
		return "",err
	}
	encodePri,err :=  EncodePri(key,pirk)
	if err != nil {
		return "",err
	}
	return encodePri,nil
}

func DecodePriByPub(addr string,encodePri string,serial string,cointype string,from string,to string,modelDecrypt model.DecryptModel)(string,error){
	hashName,err  := GetHashName(addr)
	if err != nil {
		return "",err
	}
	key,err := HashStr(hashName,addr)
	if err != nil {
		return "",err
	}

	decodePri,err := DecodePri(key,encodePri)
	if err != nil {
		return "",err
	}

	// decrypt insert value to database
	decryptEntity := model.DecryptEntity{
		SerialNo:serial,
		CoinType:cointype,
		F:from,
		T:to,
		HashFun:hashName,
	}

	err = modelDecrypt.Insert(decryptEntity)

	if err != nil {
		return "",err
	}

	return decodePri,nil
}

func GetHashName(addr string)(string,error){
	addr = strings.ToLower(addr)   //转为小写
	k := substring(addr,len(addr)-1,len(addr))
	fmt.Println(k)
	v := fmt.Sprintf("%v",[]byte(k)[0])
	//v, e := strconv.Atoi(k)


	fmt.Println(v)
	vInt,e := strconv.Atoi(v)

	if e != nil {
		return "",e
	}

	index := vInt % 7
	if index ==0 {
		return RS,nil
	}
	if index ==1 {
		return JS,nil
	}
	if index ==2 {
		return BKDR,nil
	}
	if index ==3 {
		return SDBM,nil
	}
	if index ==4 {
		return DJB,nil
	}
	if index == 5 {
		return DEK,nil
	}
	if index == 6 {
		return AP,nil
	}
	return "",nil
}

func HashStr(hashName string,hashStr string) (string ,error){

	hashStr = strings.ToLower(hashStr) //转为小写

	lowHashName := strings.ToLower(hashName)

	if !(lowHashName == strings.ToLower(JS) || lowHashName == strings.ToLower(RS) ||
		lowHashName == strings.ToLower(BKDR) || lowHashName == strings.ToLower(SDBM) ||
		lowHashName == strings.ToLower(DJB) || lowHashName == strings.ToLower(DEK)||
		lowHashName == strings.ToLower(AP)){
			return hashStr, errors.New("hash funcation is not exist ")
	}

	var key string

	if lowHashName == strings.ToLower(JS) {
		key = JSHash(hashStr)
	}

	if lowHashName == strings.ToLower(RS) {
		key =  RSHash(hashStr)
	}

	if lowHashName == strings.ToLower(BKDR) {
		key =  BKDRHash(hashStr)
	}

	if lowHashName == strings.ToLower(SDBM) {
		key =  SDBMHash(hashStr)
	}

	if lowHashName == strings.ToLower(DJB) {
		key =  DJBHash(hashStr)
	}

	if lowHashName == strings.ToLower(DEK) {
		key =  DEKHash(hashStr)
	}

	if lowHashName == strings.ToLower(AP) {
		key =  APHash(hashStr)
	}

	if len(key)>16{
		key = substring(key,len(key)-16,len(key))
	}else{
		tmpLen := 16-len(key)
		for i:=0;i<tmpLen;i++{
			key = key + "0"
		}
	}

	return key,nil
}

func EncodePri(keyStr string,prik string)(string,error){
	key := []byte(keyStr)
	priByte := []byte(prik)
	result, err := AesEncrypt(priByte, key)
	if err != nil {
		return "",nil
	}

	//fmt.Println(base64.StdEncoding.EncodeToString(result))
	return base64.StdEncoding.EncodeToString(result),nil
}

func DecodePri(keyStr string,encodePri string)(string,error){

	encryptBytes,err := base64.StdEncoding.DecodeString(encodePri)

	if err != nil {
		return "",err
	}

	keyByte := []byte(keyStr)

	origData, err := AesDecrypt(encryptBytes, keyByte)
	if err != nil {
		return "",err
	}

	return string(origData),nil
}