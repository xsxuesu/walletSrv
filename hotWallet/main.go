package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"math/big"
	"net/http"
	"strconv"
	"strings"
	"walletSrv/hotWallet/trans"
	"walletSrv/hotWallet/call"
	"walletSrv/hotWallet/constant"
	"walletSrv/hotWallet/model"
	"walletSrv/hotWallet/utils"
	"walletSrv/hotWallet/config"
)

func init()  {
	
}

var (
	Host = "0.0.0.0"
	Port = "8081"
)

func main()  {
	config.InitConf(os.Args)
	http.Handle("/", http.HandlerFunc(handlerIndex))
	if err := http.ListenAndServe(Host+":"+Port, nil); err != nil {
		log.Fatalln(err)
	}
}

func handlerIndex(w http.ResponseWriter, r *http.Request) {

	fmt.Println(r.URL.Path)

	HandlerAll(w, r)
}

func HandlerAll(writer http.ResponseWriter, request *http.Request) {

	writer.Header().Set("Content-Type", "application/json")

	switch request.Method {
	case "GET":
		utils.ResponseJson(404, "接口未找到", "", writer)
		return

	case "POST":
		//insertaddr 添加现有的addr prik
		//fetchaddr  获得新的addr
		//transfer	 转账
		//collect	 归集
		//cold2hot	 转到热钱包
		switch strings.ToLower(request.URL.Path) {
		case "/insertaddr": // 添加现有的addr prik
			insertaddr(writer,request)
			return
		case "/fetchaddr":  // 获得新的addr
			fetchaddr(writer,request)
			return
		case "/transfer":   // 转账
			transfer(writer,request)
			return
		case "/collect":    // 归集
			collect(writer,request)
			return
		case "/cold2hot":   // 转到热钱包
			cold2hot(writer,request)
			return
		case "/hottransfer":
			hottransfer(writer,request)
			return
		case "/getbalance":
			getbalance(writer,request)
			return
		case "/sendgas":
			sendgasforcontract(writer,request)
			return
		default:
			utils.ResponseJson(404, "接口未找到", "", writer)
			return
		}
	default:
		utils.ResponseJson(404, "接口未找到", "", writer)
		return
	}
}
// 查询余额
func getbalance(writer http.ResponseWriter,request *http.Request){
	defer request.Body.Close()

	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		//log.Println(err.Error())
		utils.ResponseJson(500, "获取参数:"+err.Error(), "", writer)
		return
	}

	getParam := model.GetBalanceParam{}

	err = json.Unmarshal(body,&getParam)
	if err != nil {
		//log.Println(err.Error())
		utils.ResponseJson(500, "解析参数[json.Unmarshal]："+err.Error(), "", writer)
		return
	}

	getResp:=model.GetBalanceResp{}

	switch getParam.CoinType {
	case constant.Eth:
		getResp.CoinType = constant.Eth
		getResp.Addr = getParam.Addr
		getResp.Contract = getParam.Contract
		fmt.Println(getParam.Contract,"eth")
		if getParam.Contract == ""{

			ethBalance,err := trans.GetBalance(getParam.Addr)
			if err != nil {
				utils.ResponseJson(500, "解析参数[GetBalance]："+err.Error(), "", writer)
				return
			}
			fmt.Println("get eth balance")
			getResp.Balance = strconv.FormatFloat(utils.WeiToEth(ethBalance),'f',18,64)
			fmt.Println(getResp.Balance)
		}else{

			ercBalance,err := trans.GetContractBlncRPC(getParam.Contract,getParam.Addr)
			if err != nil {
				utils.ResponseJson(500, "解析参数[GetContractBalance]："+err.Error(), "", writer)
				return
			}
			fmt.Println("get erc20 balance")
			getResp.Balance = strconv.FormatFloat(ercBalance,'f',18,64)
			fmt.Println(getResp.Balance)
		}

		byteEthResp,err := json.Marshal(getResp)
		if err != nil {
			utils.ResponseJson(500, "解析参数[json.Marshal(getResp)]："+err.Error(), "", writer)
			return
		}
		utils.ResponseJson(200, "", string(byteEthResp), writer)
		return
	case constant.Btc:


		utils.ResponseJson(200, "", "", writer)
		return

	case constant.Usdt:

		utils.ResponseJson(200, "", "", writer)
		return
	default:
		utils.ResponseJson(500, "该币种不提供服务", getParam.CoinType, writer)
		return
	}

}

//热钱包转账
func hottransfer(writer http.ResponseWriter,request *http.Request){
	defer request.Body.Close()

	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		//log.Println(err.Error())
		utils.ResponseJson(500, "获取参数:"+err.Error(), "", writer)
		return
	}

	hotParam := model.HotTransferParam{}
	err = json.Unmarshal(body, &hotParam)
	if err != nil {
		//log.Println(err.Error())
		utils.ResponseJson(500, "解析参数："+err.Error(), "", writer)
		return
	}

	switch hotParam.CoinType {
	case constant.Eth:

		ethTxHash,err := trans.SendByKS(hotParam.From,hotParam.To,hotParam.Contract,constant.ETHPASSPHRASE,hotParam.Value)
		if err != nil {
			//log.Println(err.Error())
			utils.ResponseJson(500, "转账错误："+err.Error(), "", writer)
			return
		}
		utils.ResponseJson(200, "转账完成", ethTxHash, writer)

		return
	case constant.Btc:

		btcclient,err  := trans.NewBtcClient()
		if err != nil {
			fmt.Println(err.Error())
		}

		btcTxHash,err := btcclient.SendFrom(hotParam.From,hotParam.To,hotParam.Value)
		if err != nil {
			//log.Println(err.Error())
			utils.ResponseJson(500, "转账错误："+err.Error(), "", writer)
			return
		}
		utils.ResponseJson(200, "转账完成", btcTxHash, writer)
		return

	case constant.Usdt:
		usdtClient ,err := trans.NewUSDTClient(int64(31))
		if err != nil {
			utils.ResponseJson(500, "转账错误："+err.Error(), "", writer)
			return
		}
		usdtTxHash,err := usdtClient.SendFromByAddressAndFee(hotParam.From,hotParam.To,hotParam.FeeAddr,hotParam.Value)
		//usdtTxHash,err := usdtClient.SendFundUsdt(hotParam.From,hotParam.To,hotParam.From,hotParam.Value)
		if err != nil {
			utils.ResponseJson(500, "转账错误："+err.Error(), "", writer)
			return
		}
		utils.ResponseJson(200, "转账完成", usdtTxHash, writer)
		return
	default:
		utils.ResponseJson(500, "该币种不提供服务", hotParam.CoinType, writer)
		return
	}
}
// 冷钱包转热钱包
func cold2hot(writer http.ResponseWriter, request *http.Request){
	defer request.Body.Close()

	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		//log.Println(err.Error())
		utils.ResponseJson(500, "获取参数:"+err.Error(), "", writer)
		return
	}

	// init rpc client

	client := call.InitClient(config.RpcConfig.Ip,config.RpcConfig.Port)

	cold2hotParam := model.Cold2HotParam{}
	err = json.Unmarshal(body, &cold2hotParam)
	if err != nil {
		//log.Println(err.Error())
		utils.ResponseJson(500, "解析参数："+err.Error(), "", writer)
		return
	}

	switch cold2hotParam.CoinType{
	case constant.Eth:
		ethSerial:= utils.Cold2HotRandomSerial(constant.Eth)
		// 签名转账

		ethtransferresp,err := client.CallSignEthTx(ethSerial,constant.EthCollectionAddr,constant.EthHotAddr,cold2hotParam.Contract,cold2hotParam.Value)
		if err != nil {
			utils.ResponseJson(500, "转账错误:"+err.Error(), "", writer)
			return
		}

		byteEthResp,_ := json.Marshal(ethtransferresp)

		utils.ResponseJson(200, "转账完成", string(byteEthResp), writer)
		return
	case constant.Btc:
		btcSerial := utils.Cold2HotRandomSerial(constant.Btc)

		// 转账
		btctransferresp,err := client.CallSignBtcTx(btcSerial,constant.BtcCollectionAddr,constant.BtcHotAddr,cold2hotParam.Value)
		if err != nil {
			utils.ResponseJson(500, "转账错误:"+err.Error(), "", writer)
			return
		}

		byteBtcResp,_ := json.Marshal(btctransferresp)

		utils.ResponseJson(200, "转账完成", string(byteBtcResp), writer)
		return

	case constant.Usdt:
	//	TODO   usdt 转账
		usdtSerial := utils.Cold2HotRandomSerial(constant.Btc)
		usdtTransferReq,err := client.CallSignUsdtTx(usdtSerial,constant.UsdtCollectionAddr,constant.UsdtHotAddr,cold2hotParam.Value)
		if err != nil {
			utils.ResponseJson(500, "转账错误:"+err.Error(), "", writer)
			return
		}
		byteUsdtResp,_ := json.Marshal(usdtTransferReq)

		utils.ResponseJson(200, "转账完成", string(byteUsdtResp), writer)
		return

	default:

		utils.ResponseJson(500, "该币种不提供服务", cold2hotParam.CoinType, writer)
		return
	}

}
// send gas only eth
func sendgasforcontract(writer http.ResponseWriter, request *http.Request){

	defer request.Body.Close()
	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		//log.Println(err.Error())
		utils.ResponseJson(500, "获取参数:"+err.Error(), "", writer)
		return
	}
	checkParam := model.CheckParam{}
	err = json.Unmarshal(body, &checkParam)
	if err != nil {
		//log.Println(err.Error())
		utils.ResponseJson(500, "解析参数："+err.Error(), "", writer)
		return
	}

	invalidAddrs := make(map[string]string)
	// init rpc client
	client := call.InitClient(config.RpcConfig.Ip,config.RpcConfig.Port)
	// 查询单个地址的
	if checkParam.Addr != ""{
		singleBalance,err :=trans.GetContractBlncRPCInt(checkParam.Contract,checkParam.Addr)
		if err != nil {
			utils.ResponseJson(500, "获得balance："+err.Error(), "", writer)
			return
		}
		if singleBalance.Cmp(big.NewInt(checkParam.MinCount)) == 1 {
			singleEstmateGas,err := trans.GetEstmateGasByValue(checkParam.Addr,constant.EthCollectionAddr,checkParam.Contract,singleBalance)
			if err != nil{
				utils.ResponseJson(500, "估值错误："+err.Error(), "", writer)
				return
			}
			singleEthbalance , err := trans.GetBalance(checkParam.Addr)
			if err != nil {
				utils.ResponseJson(500, "获得balance："+err.Error(), "", writer)
				return
			}

			if singleEthbalance.Cmp(new(big.Int).SetUint64(singleEstmateGas)) < 0 {
				invalidAddrs[checkParam.Addr] = new(big.Int).SetUint64(singleEstmateGas).String()
			}
		}else{
			utils.ResponseJson(500, "归集小于最小的额度", "", writer)
			return
		}

	}else{   ////// 转所有的不足够gas 的账号

		start := 0
		skip := 10
		isfinish := false

		for(!isfinish){
			collecitonResp,err := client.CallCollection(constant.Eth,start,(start+skip-1))
			start = start+skip
			if err != nil {
				continue
			}

			tf,err := strconv.ParseBool(collecitonResp.Isfinish)
			if err != nil {
				continue
			}
			isfinish = tf
			// 获取
			addrList := &model.AddrList{}

			err = json.Unmarshal(collecitonResp.Addrlist,&addrList)
			if err != nil {
				continue
			}
			// 循环发起转账
			for _,v := range  addrList.AddrList{
				//	查询余额
				balance,err :=trans.GetContractBlncRPCInt(checkParam.Contract,v)
				if err != nil {
					continue
				}

				if balance.Cmp(big.NewInt(checkParam.MinCount)) == 1 {
					estmateGas,err := trans.GetEstmateGasByValue(v,constant.EthCollectionAddr,checkParam.Contract,balance)
					if err != nil{
						continue
					}
					ethbalance , err := trans.GetBalance(v)
					if err != nil {
						continue
					}

					if ethbalance.Cmp(new(big.Int).SetUint64(estmateGas)) < 0 {
						invalidAddrs[v] = new(big.Int).SetUint64(estmateGas).String()
					}
				}
			}
		}
	}

	// 转账gas
	for k,v := range invalidAddrs{
		serial := utils.CollectionRandomSerial(constant.Eth)  //生成随机所
		// 签名转账
		valueInt,_ := new(big.Int).SetString(v,10)

		transferETH := utils.WeiToEth(valueInt)

		ethtransferresp,err := client.CallSignEthTx(serial,k,constant.EthHotAddr,"",transferETH)
		if err != nil {
			return
		}
		if ethtransferresp.Success{
			fmt.Println(serial,v,"转账txid:",ethtransferresp.Txid)
		}
	}

	// 返回

	checkResut := model.CheckResult{}
	checkResut.Contract = checkParam.Contract
	checkResut.SentAddr = invalidAddrs
	checkByte,err := json.Marshal(checkResut)
	if err != nil {
		utils.ResponseJson(500, "解析返回值错误:"+err.Error(), "", writer)
		return
	}
	utils.ResponseJson(200, "", string(checkByte), writer)
	return

}
// 归集
func collect(writer http.ResponseWriter, request *http.Request){
	defer request.Body.Close()

	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		//log.Println(err.Error())
		utils.ResponseJson(500, "获取参数:"+err.Error(), "", writer)
		return
	}

	// init rpc client
	client := call.InitClient(config.RpcConfig.Ip,config.RpcConfig.Port)

	collectParam := model.CollectParam{}
	err = json.Unmarshal(body, &collectParam)
	if err != nil {
		//log.Println(err.Error())
		utils.ResponseJson(500, "解析参数："+err.Error(), "", writer)
		return
	}
	start := 0
	skip := 10
	isfinish := false

	collectionList := []string{}

	for(!isfinish){
		collecitonResp,err := client.CallCollection(collectParam.CoinType,start,(start+skip-1))
		start = start+skip
		if err != nil {
			continue
		}

		tf,err := strconv.ParseBool(collecitonResp.Isfinish)
		if err != nil {
			continue
		}
		isfinish = tf
		// 获取
		addrList := &model.AddrList{}

		err = json.Unmarshal(collecitonResp.Addrlist,&addrList)
		if err != nil {
			continue
		}
		// 循环发起转账


		for _,v := range  addrList.AddrList{
			collectionList = append(collectionList,v)
			//	查询余额，签名 ， 转账等操作
			go func() {
				collectFlow(collectParam.CoinType,v,collectParam.Contract,collectParam.MinCount)
			}()
		}
	}

	// 归集的地址列表
	colByte,err := json.Marshal(collectionList)

	if err != nil{
		//log.Println(err.Error())
		utils.ResponseJson(500, "解析参数："+err.Error(), "", writer)
		return
	}

	collectionRespJson := model.CollectResp{
		CoinType:collectParam.CoinType,
		Collection:string(colByte),

	}
	collectionRespByte , _ := json.Marshal(collectionRespJson)

	utils.ResponseJson(200, "", string(collectionRespByte), writer)
	return
}
// 转账
func transfer(writer http.ResponseWriter, request *http.Request){
	defer request.Body.Close()

	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		//log.Println(err.Error())
		utils.ResponseJson(500, "获取参数:"+err.Error(), "", writer)
		return
	}

	// init rpc client
	client := call.InitClient(config.RpcConfig.Ip,config.RpcConfig.Port)

	transferParam := model.TransferParam{}
	err = json.Unmarshal(body, &transferParam)
	if err != nil {
		//log.Println(err.Error())
		utils.ResponseJson(500, "解析参数："+err.Error(), "", writer)
		return
	}

	// 转账
	switch transferParam.CoinType {
	case constant.Eth:
		ethtransferresp,err := client.CallSignEthTx(transferParam.Serial,transferParam.From,transferParam.To,transferParam.Contract,transferParam.Value)
		if err != nil {
			//log.Println(err.Error())
			utils.ResponseJson(500, "转账错误："+err.Error(), "", writer)
			return
		}
		ethtransferrespJson,_ := json.Marshal(ethtransferresp)
		utils.ResponseJson(200, "", string(ethtransferrespJson), writer)

	case constant.Btc:

		btctransferresp,err := client.CallSignBtcTx(transferParam.Serial,transferParam.From,transferParam.To,transferParam.Value)
		if err != nil {
			//log.Println(err.Error())
			utils.ResponseJson(500, "转账错误："+err.Error(), "", writer)
			return
		}
		btctransferrespJson,_ := json.Marshal(btctransferresp)
		utils.ResponseJson(200, "", string(btctransferrespJson), writer)

	case constant.Usdt:
		usdtTransferResp,err := client.CallSignUsdtTx(transferParam.Serial,transferParam.From,transferParam.To,transferParam.Value)
		if err != nil {
			//log.Println(err.Error())
			utils.ResponseJson(500, "转账错误："+err.Error(), "", writer)
			return
		}
		usdttransferrespJson,_ := json.Marshal(usdtTransferResp)
		utils.ResponseJson(200, "", string(usdttransferrespJson), writer)

	default:
		utils.ResponseJson(500, "币种不存在", "", writer)
		return
	}



	return
}
// 获取新数据
func fetchaddr(writer http.ResponseWriter, request *http.Request)  {
	defer request.Body.Close()

	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		//log.Println(err.Error())
		utils.ResponseJson(500, "获取参数:"+err.Error(), "", writer)
		return
	}

	// init rpc client
	client := call.InitClient(config.RpcConfig.Ip,config.RpcConfig.Port)

	fetchParam := model.FetchAddrParam{}
	err = json.Unmarshal(body, &fetchParam)
	if err != nil {
		//log.Println(err.Error())
		utils.ResponseJson(500, "解析参数："+err.Error(), "", writer)
		return
	}

	fetchResp,err := client.CallFetchAddress(fetchParam.CoinType)

	if err != nil {
		//log.Println(err.Error())
		utils.ResponseJson(500, "插入数据："+err.Error(), "", writer)
		return
	}

	fetchRespJson := model.FetchAddrResp{
		CoinType:fetchParam.CoinType,
		Addr:fetchResp.Addr,
		Success:true,
	}

	fetchRespByte , _ := json.Marshal(fetchRespJson)

	utils.ResponseJson(200, "", string(fetchRespByte), writer)
	return
}
// 插入数据错误
func insertaddr(writer http.ResponseWriter, request *http.Request){
	defer request.Body.Close()

	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		//log.Println(err.Error())
		utils.ResponseJson(500, "获取参数:"+err.Error(), "", writer)
		return
	}

	// init rpc client
	client := call.InitClient(config.RpcConfig.Ip,config.RpcConfig.Port)


	insertParam := model.InsertAddrParam{}

	err = json.Unmarshal(body, &insertParam)
	if err != nil {
		//log.Println(err.Error())
		utils.ResponseJson(500, "解析参数："+err.Error(), insertParam.Addr, writer)
		return
	}

	insertResp,err := client.CallInsertAddress(insertParam.CoinType,insertParam.Addr,insertParam.Prik)

	if err != nil {
		//log.Println(err.Error())
		utils.ResponseJson(500, "插入数据："+err.Error(), insertParam.Addr, writer)
		return
	}

	if insertResp.Success == "0"{
		utils.ResponseJson(400, "插入数据错误", insertParam.Addr, writer)
		return
	}

	insertRespInfo := model.InsertAddrResp{
		CoinType:insertParam.CoinType,
		Addr:insertParam.Addr,
		Success:true,
	}

	byteResp,_ := json.Marshal(insertRespInfo)

	utils.ResponseJson(200, "", string(byteResp), writer)
	return
}
// 查询余额，签名交易，以及转账等操作
func collectFlow(cointype string,from string,contract string,minCount int64){
	// init rpc client
	client := call.InitClient(config.RpcConfig.Ip,config.RpcConfig.Port)

	switch cointype {
	case constant.Eth:
		balance,_,err:= trans.GetBalanceAndGasPrice(from)
		if err != nil {
			fmt.Println("get balance for ",from,"******err:",err.Error())
			return
		}

		gas,err := trans.GetEstmateGasByValue(from,constant.EthCollectionAddr,contract,balance)
		if err != nil {
			fmt.Println("get balance for ",from,"******err:",err.Error())
			return
		}

		// 读取合约
		var token *big.Int
		if contract != ""{
			token, err = trans.GetContractBlncRPCInt(contract,from)

			if err != nil {
				fmt.Println("get balance for ",from,"******err:",err.Error())
				return
			}
			if token.Int64() < minCount { //如果token小于最小的值
				return
			}

			// 估算gas
			gas,err = trans.GetEstmateGasByValue(from,constant.EthCollectionAddr,contract,token)
			if err != nil {
				fmt.Println("get balance for ",from,"******err:",err.Error())
				return
			}
		}

		if err != nil {
			fmt.Println("GetEstmateGasByValue ",from,"******err:",err.Error())
			return
		}

		// 如果金额超过手续费，归集 eth
		if contract == ""{
			if balance.Uint64() > gas  {
				serial := utils.CollectionRandomSerial(constant.Eth)  //生成随机所
				value := balance.Uint64() - gas
				// 签名转账
				ethtransferresp,err := client.CallSignEthTx(serial,from,constant.EthCollectionAddr,contract,utils.WeiToEth(new(big.Int).SetUint64(value)))
				if err != nil {
					return
				}
				if ethtransferresp.Success{
					fmt.Println(serial,from,"归集交易id:",ethtransferresp.Txid)
				}
			}
		}else{  // 归集 token
			if balance.Uint64() > gas  {
				serial := utils.CollectionRandomSerial(constant.Eth)  //生成随机所
				tokenvalue,err := trans.GetContractBlncRPC(contract,from)
				if err != nil {
					return
				}
				// 签名转账
				ethtransferresp,err := client.CallSignEthTx(serial,from,constant.EthCollectionAddr,contract,tokenvalue)
				if err != nil {
					return
				}
				if ethtransferresp.Success{
					fmt.Println(serial,from,"归集交易id:",ethtransferresp.Txid)
				}
			}
		}


	case constant.Btc:

		btcclient,err  := trans.NewBtcClient()
		if err != nil {
			fmt.Println(err.Error())
		}
		btcbalance , err := btcclient.GetBalanceInAddress(from)
		if err != nil {
			fmt.Println("get balance for ",from,"******err:",err.Error())
			return
		}

		if btcbalance > 50000 { // 0.00050000 btc 手续费
			btcserial := utils.CollectionRandomSerial(constant.Btc)  //生成随机所
			btcvalue := btcbalance - 50000
			// 聪转为btc
			btcvalue = btcvalue / 100000000.0

			btcTransferResp, err := client.CallSignBtcTx(btcserial,from,constant.BtcCollectionAddr,btcvalue)
			if err != nil {
				return
			}
			if btcTransferResp.Success{
				fmt.Println(btcserial,from,"归集交易id:",btcTransferResp.Txid)
			}
		}


	case constant.Usdt:
		usdtClient,err := trans.NewUSDTClient(int64(31))
		if err != nil {
			return
		}
		omniBtcBlance,omniUsdtBalance,err := usdtClient.GetUsdtBalanceByAddr(from)
		if err != nil {
			fmt.Println("get balance for ",from,"******err:",err.Error())
			return
		}
		if omniUsdtBalance > 0 {
			if omniBtcBlance > 50000 { // 用本账号支付手续费
				usdtserial := utils.CollectionRandomSerial(constant.Usdt)  //生成随机所
				omnivalue := omniBtcBlance / 100000000.0
				client.CallSignUsdtTx(usdtserial,from,constant.UsdtCollectionAddr,omnivalue)
			}else{ // 通过其他账号支付手续费

			}
		}

	default:
		fmt.Println("coin type ",cointype," not supplied service")
	}
}