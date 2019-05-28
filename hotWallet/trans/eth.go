package trans

import (
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
	"golang.org/x/crypto/sha3"
	"io/ioutil"
	"log"
	"math/big"
	"strings"
	"walletSrv/hotWallet/config"
	"walletSrv/hotWallet/trans/token"
	"walletSrv/hotWallet/utils"

	"github.com/bitly/go-simplejson"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/parnurzeal/gorequest"
	"strconv"
	"walletSrv/proto"
	"strings"
)


type RpcRequest struct {
	JsonRpc string `json:"jsonrpc"`
	Method string `json:"method"`
	Id int `json:"id"`
	Params []string `json:"params"`
}

type RpcResponse struct {
	JsonRpc string `json:"jsonrpc"`
	Result string `json:"result"`
	Id string `json:"id"`
	Params []string `json:"params"`
}

func Call(args ...string) (*simplejson.Json) {
	method := args[0]
	params := "[]"
	if len(args) > 1 {
		params = args[1]
	}
	postBody := `{"jsonrpc":"2.0","method":"`+method+`","params":`+params+`}`


	var EthClientHost = fmt.Sprintf("http://%s:%s",config.RpcConfig.EthHost,config.RpcConfig.EthPort)

	_, body, errs := gorequest.New().Post(EthClientHost).Send(postBody).End()
	if errs != nil {
		log.Fatalln(errs)
	}
	js, err := simplejson.NewJson([]byte(body))
	if err != nil {
		log.Fatalln(err)
	}
	return js
}

func GetContractBlncRPC(contract string , addr string)(float64,error){
	var EthClientHost = fmt.Sprintf("http://%s:%s",config.RpcConfig.EthHost,config.RpcConfig.EthPort)
	//var EthClientHost = fmt.Sprintf("http://%s:%s","3.112.179.217","3342")

	readdr := addr

	if strings.HasPrefix(addr,"0x"){
		readdr = strings.TrimLeft(addr,"0x")
	}

	datastr := fmt.Sprintf("0x70a08231000000000000000000000000%s",readdr)

	jsonStr := `{
		"to":"%s",
		"data":"%s"
	},
	"latest"`

	fmtJson := fmt.Sprintf(jsonStr,contract,datastr)

	postBody := `{"jsonrpc":"2.0","method":"eth_call","params":[%s],"id":1}`

	fmtBody := fmt.Sprintf(postBody,fmtJson)

	_, body, errs := gorequest.New().Post(EthClientHost).Send(fmtBody).End()

	if errs != nil {
		return float64(0),errs[0]
	}

	js, err := simplejson.NewJson([]byte(body))
	if err != nil {
		return float64(0),err
	}
	result,err  := js.Get("result").String()

	if err != nil {
		return float64(0),err
	}
	val := result[2:]

	wei ,b:=new(big.Int).SetString("1000000000000000000",10)
	bl , b:= new(big.Int).SetString(val,16)
	if b {
		balnceVal := new(big.Int).Quo(bl,wei)
		balnceFloat :=new(big.Float).SetUint64(balnceVal.Uint64())
		floatVal , _ := balnceFloat.Float64()
		fmt.Println(strconv.FormatFloat(floatVal,'f',2,64))
		return floatVal ,nil
	}else{
		return float64(0),nil
	}

}


func GetContractBlncRPCInt(contract string , addr string)(*big.Int,error){
	var EthClientHost = fmt.Sprintf("http://%s:%s",config.RpcConfig.EthHost,config.RpcConfig.EthPort)
	//var EthClientHost = fmt.Sprintf("http://%s:%s","3.112.179.217","3342")

	readdr := addr

	if strings.HasPrefix(addr,"0x"){
		readdr = strings.TrimLeft(addr,"0x")
	}

	datastr := fmt.Sprintf("0x70a08231000000000000000000000000%s",readdr)

	jsonStr := `{
		"to":"%s",
		"data":"%s"
	},
	"latest"`

	fmtJson := fmt.Sprintf(jsonStr,contract,datastr)

	postBody := `{"jsonrpc":"2.0","method":"eth_call","params":[%s],"id":1}`

	fmtBody := fmt.Sprintf(postBody,fmtJson)

	_, body, errs := gorequest.New().Post(EthClientHost).Send(fmtBody).End()

	if errs != nil {
		return new(big.Int).SetInt64(0),errs[0]
	}

	js, err := simplejson.NewJson([]byte(body))
	if err != nil {
		return new(big.Int).SetInt64(0),err
	}
	result,err  := js.Get("result").String()

	if err != nil {
		return new(big.Int).SetInt64(0),err
	}
	val := result[2:]

	bl , b:= new(big.Int).SetString(val,16)
	if b {
		return bl ,nil
	}else{
		return new(big.Int).SetInt64(0),nil
	}

}


func SendRawTrans(serial string , rawhash string )(string,string,error)  {
	var EthClientHost = fmt.Sprintf("http://%s:%s",config.RpcConfig.EthHost,config.RpcConfig.EthPort)

	fmt.Println("rawhash1:")
	fmt.Println(rawhash)

	var tx *types.Transaction

	rawtx,err := hex.DecodeString(rawhash)

	rlp.DecodeBytes(rawtx, &tx)

	fmt.Println(tx.Hash().Hex())        // 0x5d49fcaa394c97ec8a9c3e7bd9e8388d420fb050a52083ca52ff24b3b65bc9c2
	fmt.Println(tx.Value().String())    // 10000000000000000
	fmt.Println(tx.Gas())               // 105000
	fmt.Println(tx.GasPrice().Uint64()) // 102000000000
	fmt.Println(tx.Nonce())             // 110644
	fmt.Println(tx.Data())              // []
	fmt.Println(tx.To().Hex())          // 0x55fE59D8Ad77035154dDd0AD0388D09Dd4047A8e


	postBody := `{"jsonrpc":"2.0","method":"eth_sendRawTransaction","params":["`+rawhash+`"]}`
	_, body, errs := gorequest.New().Post(EthClientHost).Send(postBody).End()
	if errs != nil {
		return serial,"", errs[0]
	}
	js, err := simplejson.NewJson([]byte(body))
	if err != nil {
		return serial,"", err
	}
	result,err  := js.Get("result").String()
	if err != nil {
		return serial,"", err
	}



	return serial,result ,nil
}

func SendRawHash(serial string , rawhash string)(string,string,error){
	var EthClientHost = fmt.Sprintf("http://%s:%s",config.RpcConfig.EthHost,config.RpcConfig.EthPort)

	client, err := Connect(EthClientHost)
	if err != nil {
		return serial,"",err
	}

	txRawHash := fmt.Sprintf("0x%s",rawhash)

	txHash,err := client.SendRawTransactionByHash(context.Background(),txRawHash)

	if err != nil {
		fmt.Println("err:",err.Error())
		return serial,"",err
	}

	return serial,txHash.String(),nil
}

func GetBalance(addr string)(*big.Int,error){
	var EthClientHost = fmt.Sprintf("http://%s:%s",config.RpcConfig.EthHost,config.RpcConfig.EthPort)
	//var EthClientHost = fmt.Sprintf("http://%s:%s","3.112.179.217","3342")

	client, err := ethclient.Dial(EthClientHost)
	if err != nil {
		return big.NewInt(0),err
	}
	defer client.Close()

	balanceAddress := common.HexToAddress(addr)
	balanceWei,err := client.BalanceAt(context.Background(),balanceAddress,nil)
	if err != nil {
		return big.NewInt(0),err
	}
	return balanceWei,nil
}

func GetContractBalanceInt(addr string,contract string)(*big.Int,error){
	var EthClientHost = fmt.Sprintf("http://%s:%s",config.RpcConfig.EthHost,config.RpcConfig.EthPort)

	client, err := ethclient.Dial(EthClientHost)
	if err != nil {
		return big.NewInt(0),err
	}
	defer client.Close()

	return token.GetContractBlanceInt(contract,addr,client)
}

func GetContractBalance2(addr string,contract string)(float64,error){
	var EthClientHost = fmt.Sprintf("http://%s:%s",config.RpcConfig.EthHost,config.RpcConfig.EthPort)

	client, err := ethclient.Dial(EthClientHost)
	if err != nil {
		return float64(0),err
	}
	defer client.Close()

	return token.GetContractBlanceOf(contract,addr,client)
}

func GetBalanceAndGasPrice(addr string)(*big.Int,*big.Int,error){
	var EthClientHost = fmt.Sprintf("http://%s:%s",config.RpcConfig.EthHost,config.RpcConfig.EthPort)

	client, err := ethclient.Dial(EthClientHost)
	if err != nil {
		return big.NewInt(0),big.NewInt(0),err
	}
	defer client.Close()

	balanceAddress := common.HexToAddress(addr)
	balanceWei,err := client.BalanceAt(context.Background(),balanceAddress,nil)
	if err != nil {
		return big.NewInt(0),big.NewInt(0),err
	}
	gasPrice,err := client.SuggestGasPrice(context.Background())
	if err != nil {
		return big.NewInt(0),big.NewInt(0),err
	}
	return balanceWei,gasPrice,nil
}

func GetNonce(addr string)(uint64,error){
	var EthClientHost = fmt.Sprintf("http://%s:%s",config.RpcConfig.EthHost,config.RpcConfig.EthPort)

	client, err := ethclient.Dial(EthClientHost)
	if err != nil {
		return uint64(0),err
	}
	defer client.Close()

	return client.PendingNonceAt(context.Background(),common.HexToAddress(addr))

}

func GetNetworkId(addr string)(*big.Int,error){
	var EthClientHost = fmt.Sprintf("http://%s:%s",config.RpcConfig.EthHost,config.RpcConfig.EthPort)

	client, err := ethclient.Dial(EthClientHost)
	if err != nil {
		return new(big.Int),err
	}
	defer client.Close()

	return client.NetworkID(context.Background())

}

func GetEstmateGasByValue(from,to,contract string,value *big.Int)(uint64,error){
	var EthClientHost = fmt.Sprintf("http://%s:%s",config.RpcConfig.EthHost,config.RpcConfig.EthPort)

	client, err := ethclient.Dial(EthClientHost)
	if err != nil {
		return uint64(0),err
	}
	defer client.Close()

	data := []byte{}
	// 预估消耗的gas
	fromAddress := common.HexToAddress(from)
	toAddress := common.HexToAddress(to)

	msg := ethereum.CallMsg{
		From: fromAddress,
		To: &toAddress,
		Value:value,
		Data: data,
	}

	if contract != ""{
		contractAddress := common.HexToAddress(contract)

		transferFnSignature := []byte("transfer(address,uint256)")
		hash := sha3.NewLegacyKeccak256()
		hash.Write(transferFnSignature)
		methodID := hash.Sum(nil)[:4]

		paddedAddress := common.LeftPadBytes(toAddress.Bytes(), 32)

		paddedAmount := common.LeftPadBytes(value.Bytes(), 32)

		data = append(data, methodID...)
		data = append(data, paddedAddress...)
		data = append(data, paddedAmount...)


		msg = ethereum.CallMsg{
			From: fromAddress,
			To: &contractAddress,
			Value: big.NewInt(0),
			Data: data,
		}
	}

	return GetEstmateGas(msg)
}

func GetEstmateGas(msg ethereum.CallMsg)(uint64,error){
	var EthClientHost = fmt.Sprintf("http://%s:%s",config.RpcConfig.EthHost,config.RpcConfig.EthPort)

	client, err := ethclient.Dial(EthClientHost)
	if err != nil {
		return uint64(0),err
	}
	defer client.Close()

	return client.EstimateGas(context.Background(),msg)
}

func FloatToBigInt(val float64) *big.Int {
	bigval := new(big.Float)
	bigval.SetFloat64(val)

	coin := new(big.Float)
	coin.SetInt(big.NewInt(1000000000000000000))
	bigval.Mul(bigval, coin)

	result := new(big.Int)
	f,_ := bigval.Uint64()
	result.SetUint64(f)

	return result
}

func KeyStoreToPri(path string, passphrase string)(*ecdsa.PrivateKey,error){

	keyjson, err := ioutil.ReadFile(path)
	if err != nil {
		return nil,err
	}
	key, err := keystore.DecryptKey(keyjson,passphrase)

	if err != nil {
		return nil,err
	}

	fmt.Println(fmt.Sprintf("%x", crypto.FromECDSA(key.PrivateKey)))
	return key.PrivateKey,nil
}


func SendByKS(from string,to string,contract string,passphrase string,value float64)(string,error){

	var EthClientHost = fmt.Sprintf("http://%s:%s",config.RpcConfig.EthHost,config.RpcConfig.EthPort)

	client, err := ethclient.Dial(EthClientHost)
	if err != nil {
		return "",err
	}
	defer client.Close()

	var KEYSTORE_DIR = config.RpcConfig.EthKeyDir

	fmt.Println("KEYSTORE_DIR:")
	fmt.Println(KEYSTORE_DIR)

	ks := keystore.NewKeyStore(
		KEYSTORE_DIR,
		keystore.LightScryptN,
		keystore.LightScryptP)

	// Create account definitions

	fromAccDef := accounts.Account{
		Address: common.HexToAddress(from),
	}

	// Find the signing account
	signAcc, err := ks.Find(fromAccDef)
	if err != nil {
		fmt.Println(err.Error())
		return "",err
	}

	// Unlock the signing account
	errUnlock := ks.Unlock(signAcc, passphrase)
	if errUnlock != nil {
		fmt.Println(errUnlock.Error())
		return "",errUnlock
	}

	_,gasprice,err := GetBalanceAndGasPrice(from)
	if err != nil {
		return "",err
	}

	networkID, err := GetNetworkId(from)
	if err != nil {
		return "",err
	}

	nonce,err := GetNonce(from)
	if err != nil {
		return "",err
	}
	// Construct the transaction

	toAddr := common.HexToAddress(to)
	msg := ethereum.CallMsg{
		From:     common.HexToAddress(from),
		To:       &toAddr,
		GasPrice: gasprice,
		Value:    new(big.Int).SetUint64(utils.EthToWei(value)),
		Data:     []byte{},
	}

	// 如果转合约的
	var data []byte
	var contractAddr common.Address
	if contract != ""{
		contractAddr = common.HexToAddress(contract)

		transferFnSignature := []byte("transfer(address,uint256)")
		hash := sha3.NewLegacyKeccak256()
		hash.Write(transferFnSignature)
		methodID := hash.Sum(nil)[:4]

		paddedAddress := common.LeftPadBytes(toAddr.Bytes(), 32)

		decimal,err :=token.GetContractDecimal(contract,client)
		if err != nil {
			return "",err
		}
		// 计算token的数量
		bigAmount := big.NewInt(0)
		if decimal == big.NewInt(0){
			vValue := new(big.Float).SetFloat64(value)
			UvValue,_ := vValue.Uint64()
			bigAmount = new(big.Int).SetUint64(UvValue)
		}else{
			decimalBig := new(big.Float).SetInt(decimal)
			amount := new(big.Float).Mul(big.NewFloat(value),decimalBig)
			u64Amount ,_ := amount.Uint64()
			bigAmount = new(big.Int).SetUint64(u64Amount)
		}

		paddedAmount := common.LeftPadBytes(bigAmount.Bytes(), 32)

		data = append(data, methodID...)
		data = append(data, paddedAddress...)
		data = append(data, paddedAmount...)
		msg = ethereum.CallMsg{
			From: common.HexToAddress(from),
			To: &contractAddr,
			Value: big.NewInt(0),
			Data: data,
		}
	}

	gasLimit,err := GetEstmateGas(msg)
	if err != nil {
		fmt.Println(err.Error())
		return "",err
	}
	tx := types.NewTransaction(nonce, toAddr, new(big.Int).SetUint64(utils.EthToWei(value)), gasLimit, gasprice, []byte{})
	if contract != ""{
		tx = types.NewTransaction(nonce, contractAddr, new(big.Int).SetInt64(0), gasLimit, gasprice, data)
	}

	// 签名交易
	signtx,err := ks.SignTxWithPassphrase(fromAccDef,passphrase,tx,networkID)

	if err != nil {
		fmt.Println(err.Error())
		return "",err
	}
	// Send the transaction to the network
	txErr := client.SendTransaction(context.Background(), signtx)
	if txErr != nil {
		fmt.Println(txErr.Error())
		return "",txErr
	}

	//go ListnerTrans( signtx.Hash().String(),signtx.Hash()) // 监听交易状态

	return signtx.Hash().String() ,nil
}

//func ListnerTrans(serial string,txHash common.Hash)(error){
//	var EthClientHost = fmt.Sprintf("http://%s:%s",config.RpcConfig.EthHost,config.RpcConfig.EthPort)
//
//	client, err := Connect(EthClientHost)
//	if err != nil {
//		return err
//	}
//
//	// check transaction receipt
//	receiptChan := make(chan *types.Receipt)
//	client.CheckTransaction(context.TODO(), receiptChan, txHash, 5)
//	result := <-receiptChan
//	//receipt.Status
//	fmt.Println(result.Status,"txid:",result.TxHash.String())
//	//	TODO 修改 交易状态
//
//	client.UpdateStatus(serial,result.TxHash.String(),result.Status)
//
//	// 关闭 rpc eth client
//	client.rpcClient.Close()
//	client.EthClient.Close()
//	return nil
//}

// 创建 eth 交易签名hash
func CreateEthTrans(serial ,from ,to ,contract string,value float64)(*proto.SignEthTransRequest,error){
	var EthClientHost = fmt.Sprintf("http://%s:%s",config.RpcConfig.EthHost,config.RpcConfig.EthPort)

	client, err := ethclient.Dial(EthClientHost)
	if err != nil {
		return &proto.SignEthTransRequest{},err
	}
	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		return &proto.SignEthTransRequest{},err
	}
	fmt.Printf("chainId : %d had connected!",chainID.Uint64())

	defer client.Close()

	fromAddress := common.HexToAddress(from)
	toAddress := common.HexToAddress(to)
	contractAddress := common.HexToAddress(contract)

	transferValue := utils.EthToWei(value)

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		return &proto.SignEthTransRequest{},err
	}

	data := []byte{}
	// 预估消耗的gas

	msg := ethereum.CallMsg{
		From: fromAddress,
		To: &toAddress,
		Value: new(big.Int).SetUint64(utils.EthToWei(value)),
		Data: data,
	}

	if contract != ""{

		transferFnSignature := []byte("transfer(address,uint256)")
		hash := sha3.NewLegacyKeccak256()
		hash.Write(transferFnSignature)
		methodID := hash.Sum(nil)[:4]

		paddedAddress := common.LeftPadBytes(toAddress.Bytes(), 32)

		decimal,err :=token.GetContractDecimal(contract,client)
		if err != nil {
			return &proto.SignEthTransRequest{},err
		}
		//
		bigAmount := big.NewInt(0)
		if decimal == big.NewInt(0){
			vValue := new(big.Float).SetFloat64(value)
			UvValue,_ := vValue.Uint64()
			bigAmount = new(big.Int).SetUint64(UvValue)
		}else{
			decimalBig := new(big.Float).SetInt(decimal)
			amount := new(big.Float).Mul(big.NewFloat(value),decimalBig)
			u64Amount ,_ := amount.Uint64()
			bigAmount = new(big.Int).SetUint64(u64Amount)
		}

		//重新计算value
		transferValue = bigAmount.Uint64()

		paddedAmount := common.LeftPadBytes(bigAmount.Bytes(), 32)

		data = append(data, methodID...)
		data = append(data, paddedAddress...)
		data = append(data, paddedAmount...)


		msg = ethereum.CallMsg{
			From: fromAddress,
			To: &contractAddress,
			Value: big.NewInt(0),
			Data: data,
		}
	}

	estimateGas , err := client.EstimateGas(context.Background(),msg)
	if err != nil {
		return &proto.SignEthTransRequest{},err
	}
	//check balance value
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)

	if err != nil {
		return &proto.SignEthTransRequest{},err
	}

	return &proto.SignEthTransRequest{
		Serial:serial,
		From:from,
		To:to,
		Nonce:strconv.FormatUint(nonce,10),
		Value:strconv.FormatUint(transferValue,10),
		Gaslimit:strconv.FormatUint(estimateGas,10),
		Gasprice:strconv.FormatUint(gasPrice.Uint64(),10),
		Chainid:strconv.FormatUint(chainID.Uint64(),10),
		Contract:contract,
	},nil
}

func GetContractBalance(contract string , addr string)(float64,error){

	var EthClientHost = fmt.Sprintf("http://%s:%s",config.RpcConfig.EthHost,config.RpcConfig.EthPort)
	//var EthClientHost = fmt.Sprintf("http://%s:%s","127.0.0.1","7545")

	readdr := addr

	if strings.HasPrefix(addr,"0x"){
	readdr = strings.TrimLeft(addr,"0x")
	}

	datastr := fmt.Sprintf("0x70a08231000000000000000000000000%s",readdr)

	jsonStr := `{
		"to":"%s",
		"data":"%s"
	},
	"latest"`

	fmtJson := fmt.Sprintf(jsonStr,contract,datastr)

	postBody := `{"jsonrpc":"2.0","method":"eth_call","params":[%s],"id":1}`

	fmtBody := fmt.Sprintf(postBody,fmtJson)

	_, body, errs := gorequest.New().Post(EthClientHost).Set("Content-Type", "application/json").Send(fmtBody).End()

	if errs != nil {
	return float64(0),errs[0]
	}

	js, err := simplejson.NewJson([]byte(body))
	if err != nil {
	return float64(0),err
	}
	result,err  := js.Get("result").String()

	if err != nil {
	return float64(0),err
	}
	val := result[2:]

	wei ,b:=new(big.Int).SetString("1000000000000000000",10)
	bl , b:= new(big.Int).SetString(val,16)
	if b {
	balnceVal := new(big.Int).Quo(bl,wei)
	balnceFloat :=new(big.Float).SetUint64(balnceVal.Uint64())
	floatVal , _ := balnceFloat.Float64()
	return floatVal ,nil
	}else{
	return float64(0),nil
	}
}