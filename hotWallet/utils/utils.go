package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"math/rand"
	"net/http"
	"time"
	"walletSrv/hotWallet/config"
	"walletSrv/hotWallet/constant"
)

type JsonResult struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data string `json:"data"`
}

type Result struct {
	Code int                    `json:"code"`
	Msg  string                 `json:"msg"`
	Data map[string]interface{} `json:"data"`
}

// 返回json
func ResponseJson(code int, msg string, data string, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding,Authorization,X-Requested-With")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Request-Method", "GET,HEAD,PUT,PATCH,POST,DELETE")
	w.Header().Set("X-Requested-With", "XmlHttpRequest")
	returnData := JsonResult{}
	returnData.Code = code
	returnData.Msg = msg
	returnData.Data = data
	jsonData, err := json.Marshal(returnData)
	if err != nil {
		fmt.Println("请检查返回数据格式是否为标准interface,error:", err.Error())
		w.Write([]byte(`{"code":"500","msg":"系统内部错误！","data":[]}`))
		return
	}
	w.Write(jsonData)
	return
}

func WeiToEth(val *big.Int) float64 {

	x := new(big.Float)
	x.SetInt(val)

	coin := new(big.Float)
	coin.SetInt(big.NewInt(1000000000000000000))

	z := new(big.Float).Quo(x, coin)
	result , _ := z.Float64()
	return result
}

func EthToWei(val float64) uint64 {
	x := big.NewFloat(val)

	coin := new(big.Float)
	coin.SetInt(big.NewInt(1000000000000000000))

	z := new(big.Float).Mul(x,coin)
	result , _ := z.Uint64()

	return result
}

func BtcToSatoshi(val float64) int64{
	x := big.NewFloat(val)

	coin := new(big.Float)
	coin.SetInt(big.NewInt(100000000))

	z := new(big.Float).Mul(x,coin)
	result , _ := z.Int64()
	return result
}

func CollectionRandomSerial(cointype string)(string){
	rand.Seed(time.Now().UnixNano())
	randstr:=rand.Intn(9999999)
	return fmt.Sprintf("%s_%s%s",cointype,time.Now().String(),fmt.Sprintf("%d",randstr))
}

func Cold2HotRandomSerial(cointype string)(string){
	rand.Seed(time.Now().UnixNano())
	randstr:=rand.Intn(9999999)
	return fmt.Sprintf("transfer_%s_%s%s",cointype,time.Now().String(),fmt.Sprintf("%d",randstr))
}


func CallRequest( cointype string,txid string, status string){
	defer func() { // 必须要先声明defer，否则不能捕获到panic异常
		if err := recover(); err != nil {
			fmt.Println(err) // 这里的err其实就是panic传入的内容，55
		}
	}()

	var callback string

	switch cointype {
	case constant.Eth:
		callback = config.RpcConfig.EthCallBack
	case constant.Btc:
		callback = config.RpcConfig.BtcCallBack
	case constant.Usdt:
		callback = config.RpcConfig.UsdtCallBack
	}


	jsonMap := make(map[string]interface{})
	jsonMap["txid"] = txid
	jsonMap["status"] = status

	jsonBytes, err := json.Marshal(jsonMap)
	if err != nil {
		panic(err)
	}

	req, err := http.NewRequest("POST", callback, bytes.NewBuffer(jsonBytes))
	if err != nil {
		panic(err)
	}
	req.Header.Set("Content-Type", "application/json")
	var client http.Client
	res, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}
	fmt.Println("body:")
	fmt.Println(string(body))
}