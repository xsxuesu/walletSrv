package trans

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ibclabs/omnilayer"
	"github.com/ibclabs/omnilayer/omnijson"
	"github.com/wanglei-ok/usdtapi"
	"github.com/wanglei-ok/usdtapi/rpc"
	"strconv"
	"time"
	"walletSrv/hotWallet/config"
	"walletSrv/proto"
)

// default account for reserved usage, which represent
// account belongs to enterpise default
var USDTDEFAULT_ACCOUNT = "company"

// default confirmation
var DEFAULT_CONFIRMATION = 1

type UsdtAccount struct {
	Account  string
	Balance  float64
	Addresses []string
}

type USDTClient struct {
	rpcClient *omnilayer.Client
	propertyId int64
}

// connect to omnicore with HTTP RPC transport
func NewUSDTClient(propertyId int64) (*USDTClient, error) {
	var UsdtHostPort = fmt.Sprintf("%s:%s",config.RpcConfig.UsdtHost,config.RpcConfig.UsdtPort)
	var UsdtrpcUser = config.RpcConfig.UsdtRpcUser
	var UsdtrpcPassword = config.RpcConfig.UsdtRpcPwd
	connCfg := &omnilayer.ConnConfig{
		Host: UsdtHostPort,
		User: UsdtrpcUser,
		Pass: UsdtrpcPassword,
	}
	client := &USDTClient{propertyId: propertyId}
	client.rpcClient = omnilayer.New(connCfg)

	return client, nil
}

// GetBlockCount
func (client *USDTClient) GetBlockCount() (int64, error) {
	var res omnijson.GetBlockChainInfoResult
	if res, err := client.rpcClient.GetBlockChainInfo(); err == nil {
		return res.Blocks, nil
	}
	return res.Blocks, nil
}

// Ping
func (client *USDTClient) Ping() error {
	_, err := client.rpcClient.GetInfo()
	return err
}

// Get Usdt Balance
func (client *USDTClient) GetUsdtBalanceByAddr(address string)(float64,float64,error){


	omniBtcCmd := omnijson.OmniGetBalanceCommand{
		Address:address,
		PropertyID:0,
	}
	omniBtcRest,err := client.rpcClient.OmniGetBalance(omniBtcCmd)
	if err != nil {
		return float64(0),float64(0),err
	}
	btcbalance , err := strconv.ParseFloat(omniBtcRest.Balance,64)

	if err != nil {
		return float64(0),float64(0),err
	}

	omniUsdtCmd := omnijson.OmniGetBalanceCommand{
		Address:address,
		PropertyID:31,
	}


	omniUsdtRest,err := client.rpcClient.OmniGetBalance(omniUsdtCmd)
	if err != nil {
		return btcbalance,float64(0),err
	}
	usdtbalance , err := strconv.ParseFloat(omniUsdtRest.Balance,64)
	if err != nil {
		return btcbalance,float64(0),err
	}

	return btcbalance,usdtbalance,nil
}

// GetAddress - default address
func (client *USDTClient) GetAddress(account string) (string, error) {
	if len(account) == 0 {
		account = DEFAULT_ACCOUNT
	}
	address, err := client.rpcClient.GetAccountAddress(account)
	if err != nil {
		return "", err
	}

	return address, nil
}

// Create Account
// Returns customized account info
func (client *USDTClient) CreateAccount(account string) (UsdtAccount, error) {
	// GetAddress will create account if not exists
	address, err := client.GetAddress(account)
	if err != nil {
		return UsdtAccount{}, err
	}

	return UsdtAccount{
		Account:   account,
		Balance:   0.0,
		Addresses: []string{address},
	}, nil
}

// GetAccountInfo
func (client *USDTClient) GetAccountInfo(account string, minConf int) (UsdtAccount, error) {
	addresses, err := client.GetAddressesByAccount(account)
	if err != nil {
		return UsdtAccount{}, err
	}
	var balance float64 = 0
	for _, addr := range addresses {
		cmd := omnijson.OmniGetBalanceCommand{
			Address:    addr,
			PropertyID: int32(client.propertyId),
		}
		if curBalance, err := client.rpcClient.OmniGetBalance(cmd); err == nil {
			if b, err := strconv.ParseFloat(curBalance.Balance, 64); err == nil {
				balance += b
			}
		}
	}

	return UsdtAccount{
		Account:   account,
		Balance:   balance,
		Addresses: addresses,
	}, nil
}

// GetAddressesByAccount
func (client *USDTClient) GetAddressesByAccount(account string) ([]string, error) {
	if len(account) == 0 {
		account = DEFAULT_ACCOUNT
	}
	addresses, err := client.rpcClient.GetAddressesByAccount(account)
	if err != nil {
		return []string{}, err
	}

	addrs := make([]string, 0)
	for _, addr := range addresses {
		addrs = append(addrs, addr)
	}
	return addrs, nil
}

// GetNewAddress ...
func (client *USDTClient) GetNewAddress(account string) (string, error) {
	if len(account) == 0 {
		account = DEFAULT_ACCOUNT
	}

	address, err := client.rpcClient.GetNewAddress(account)
	if err != nil {
		return "", err
	}
	return address, nil
}

// ListAccountsMinConf
func (client *USDTClient) ListAccountsMinConf(minConf int) (map[string]float64, error) {
	accounts := make(map[string]float64)

	accountsWithAmount, err := client.rpcClient.ListAccounts(int64(minConf))
	if err != nil {
		return accounts, err
	}

	for account, _ := range accountsWithAmount {
		var accountInfo UsdtAccount
		accountInfo, err = client.GetAccountInfo(account, minConf)
		if err != nil {
			accounts[account] = -1
		} else {
			accounts[account] = accountInfo.Balance
		}
	}

	return accounts, nil
}


//SendFrom ...omni_funded_send
func (client *USDTClient) SendFrom(account, address string, amount float64) (string, error) {
	fromAddr, err := client.rpcClient.GetAccountAddress(account)
	if err != nil {
		return "", err
	}
	hash, _ := client.rpcClient.OmniFoundedSend(fromAddr, address, client.propertyId, floatToString(amount), fromAddr)
	return hash, nil
}

//SendFrom ...omni_funded_send
func (client *USDTClient) SendFromByAddress(from, to string, amount float64) (string, error) {
	hash, _ := client.rpcClient.OmniFoundedSend(from, to, client.propertyId, floatToString(amount), from)
	return hash, nil
}

//SendFrom ...omni_funded_send
func (client *USDTClient) SendFromByAddressAndFee(from, to ,fee string, amount float64) (string, error) {
	hash, _ := client.rpcClient.OmniFoundedSend(from, to, client.propertyId, floatToString(amount), fee)
	return hash, nil
}

// Move - omni_funded_send
func (client *USDTClient) Move(from, to string, amount float64) (string,bool, error) {
	hash, err := client.rpcClient.OmniFoundedSend(from, to, client.propertyId, floatToString(amount), to)
	if err != nil {
		return "",false, err
	}
	return hash,true, nil
}

// ListUnspentMin
func (client *USDTClient) ListUnspentMin(address string,minConf int) ([]byte, error) {
	unspentReq := omnijson.ListUnspentCommand{
		Min:minConf,
		Addresses:[]string{address},
	}
	result,err := client.rpcClient.ListUnspent(unspentReq)
	if err != nil {
		return []byte{},err
	}
	resultByte,err := json.Marshal(result)
	if err != nil {
		return []byte{},err
	}
	return resultByte,nil
}

// rawHash 签名 by prikey
func (client *USDTClient) SignTxWithKey(rawHash string,from string, prik string)(string,error){

	listunspentCmd := omnijson.ListUnspentCommand{
		Min:1,
		Addresses:[]string{from},
	}
	listUnspent,err := client.rpcClient.ListUnspent(listunspentCmd)

	if err != nil {
		return "",err
	}

	var listPrevious = []omnijson.Previous{}

	for i := range listUnspent{
		tmpPre := omnijson.Previous{
			TxID:listUnspent[i].Tx,
			Vout:listUnspent[i].Vout,
			ScriptPubKey:listUnspent[i].ScriptPubKey,
			RedeemScript:listUnspent[i].RedeemScript,
			Value:listUnspent[i].Amount,
		}
		listPrevious = append(listPrevious,tmpPre)
	}


	signWithKeyCmd := omnijson.SignRawTransactionWithKeyCommand{
		Hex:rawHash,
		Keys:[]string{prik},
		Previous:listPrevious,
	}

	signWithKeyResult,err := client.rpcClient.SignRawTransactionWithKey(signWithKeyCmd)
	if err != nil {
		return "",err
	}

	return signWithKeyResult,nil
}


func (client *USDTClient) SignTx(rawHash string,from string, prik string)(string,error){

	listunspentCmd := omnijson.ListUnspentCommand{
		Min:1,
		Addresses:[]string{from},
	}
	listUnspent,err := client.rpcClient.ListUnspent(listunspentCmd)

	if err != nil {
		return "",err
	}

	var listPrevious = []omnijson.Previous{}

	for i := range listUnspent{
		tmpPre := omnijson.Previous{
			TxID:listUnspent[i].Tx,
			Vout:listUnspent[i].Vout,
			ScriptPubKey:listUnspent[i].ScriptPubKey,
			RedeemScript:listUnspent[i].RedeemScript,
			Value:listUnspent[i].Amount,
		}
		listPrevious = append(listPrevious,tmpPre)
	}


	signWithKeyCmd := omnijson.SignRawTransactionCommand{
		Hex:rawHash,
		Keys:[]string{prik},
		Previous:listPrevious,
	}

	signWithKeyResult,err := client.rpcClient.SignRawTransaction(signWithKeyCmd)
	if err != nil {
		return "",err
	}

	if signWithKeyResult.Complete{
		return signWithKeyResult.Hex,nil
	}else{
		return "", errors.New(signWithKeyResult.Errors[0].Error)
	}
}

// Send signed txhash
func (client *USDTClient) SendUsdtSignTX(txHash string)(string,error){
	var UsdtHostPort = fmt.Sprintf("%s:%s",config.RpcConfig.UsdtHost,config.RpcConfig.UsdtPort)
	var UsdtrpcUser = config.RpcConfig.UsdtRpcUser
	var UsdtrpcPassword = config.RpcConfig.UsdtRpcPwd

	connCfg := &rpc.ConnConfig{
		Host: UsdtHostPort,
		User: UsdtrpcUser,
		Pass: UsdtrpcPassword,
	}
	omni := usdtapi.NewOmniClient(connCfg)

	return omni.SendRawTrans(txHash)
}

// Full Create sign send tx
func (client *USDTClient) FullUsdtTx(from,to string,amount ,fee float64)(string,error){
	// get unspent list
	listUnspentCmd := omnijson.ListUnspentCommand{
		Min:1,
		Addresses:[]string{from},
	}
	listUnspent,err := client.rpcClient.ListUnspent(listUnspentCmd)
	if err != nil {
		return "",err
	}

	paras := []omnijson.CreateRawTransactionParameter{}
	changes := []omnijson.OmniCreateRawTxChangeParameter{}

	for _, v := range listUnspent{

		tmp :=	omnijson.CreateRawTransactionParameter{
			Tx:v.Tx,
			Vout:v.Vout,
		}
		tmpchange := omnijson.OmniCreateRawTxChangeParameter{
			Tx:v.Tx,
			Vout:v.Vout,
			ScriptPubKey:v.ScriptPubKey,
			Value:v.Amount,
		}

		paras = append(paras,tmp)
		changes = append(changes,tmpchange)
	}


	// create payload send command

	payloadCmd := omnijson.OmniCreatePayloadSimpleSendCommand{
		Property:31,
		Amount: strconv.FormatFloat(amount,'f',8,64),
	}

	payloadRst , err := client.rpcClient.OmniCreatePayloadSimpleSend(payloadCmd)
	if err != nil {
		return "",err
	}

	rawTxCmd := omnijson.CreateRawTransactionCommand{
		Parameters:paras,
	}

	rawTxResult,err := client.rpcClient.CreateRawTransaction(rawTxCmd)
	if err != nil {
		return "",err
	}

	opReturnCmd:=omnijson.OmniCreateRawTxOpReturnCommand{
		Raw:rawTxResult,
		Payload:payloadRst,
	}

	opResult,err := client.rpcClient.OmniCreateRawTxOpReturn(opReturnCmd)
	if err != nil {
		return "",err
	}

	referenceCmd := omnijson.OmniCreateRawTxReferenceCommand{
		Raw:opResult,
		Destination:to,
		Amount:amount,
	}

	referenceResult,err := client.rpcClient.OmniCreateRawTxReference(referenceCmd)
	if err != nil {
		return "",err
	}


	changeCmd := omnijson.OmniCreateRawTxChangeCommand{
		Raw:referenceResult,
		Previous:changes,
		Destination:to,
		Fee:fee,
	}

	changeResult ,err := client.rpcClient.OmniCreateRawTxChange(changeCmd)
	if err != nil {
		return "",err
	}

	return changeResult,nil
}

// Create proto request params
func (client *USDTClient) CreateUsdtTransByClient(serial string,from string,to string,value string,fee float64)(*proto.SignUsdtTransRequest,error){
	cmd := omnijson.ListUnspentCommand{
		Min:1,
		Addresses:[]string{from},
	}
	unspentlist , err :=  client.rpcClient.ListUnspent(cmd)

	if err != nil {
		return &proto.SignUsdtTransRequest{},err
	}

	unspentByte , err := json.Marshal(unspentlist)
	if err != nil {
		return &proto.SignUsdtTransRequest{},err
	}

	return &proto.SignUsdtTransRequest{
		Serial:serial,
		From:from,
		To:to,
		Value:value,
		Unspentlist:unspentByte,
	},nil
}


// Create Usdt Trans
func CreateUsdtTrans(serial string,from string,to string,value float64)(*proto.SignUsdtTransRequest,error){
	unspentByte,err := GetUsdtUnspentByAddress(from)

	if err != nil {
		return &proto.SignUsdtTransRequest{},err
	}

	signUsdtReq := &proto.SignUsdtTransRequest{
		Serial:serial,
		From:from,
		To:to,
		Value:strconv.FormatFloat(value,'f',8,64),
		Unspentlist:unspentByte,
	}
	return signUsdtReq,nil

}


func GetUsdtUnspentByAddress(address string)([]byte,error){
	var UsdtHostPort = fmt.Sprintf("%s:%s",config.RpcConfig.UsdtHost,config.RpcConfig.UsdtPort)
	var UsdtrpcUser = config.RpcConfig.UsdtRpcUser
	var UsdtrpcPassword = config.RpcConfig.UsdtRpcPwd
	connCfg := &rpc.ConnConfig{
		Host: UsdtHostPort,
		User: UsdtrpcUser,
		Pass: UsdtrpcPassword,
	}
	omni := usdtapi.NewOmniClient(connCfg)

	result,err := omni.GetListunspent(address,0)
	if err != nil {
		return []byte{},err
	}
	resultByte,err := json.Marshal(result)
	if err != nil {
		return []byte{},err
	}
	return resultByte,nil
}


func SendUsdt(from string,to string,value float64)(string,error){
	var UsdtHostPort = fmt.Sprintf("%s:%s",config.RpcConfig.UsdtHost,config.RpcConfig.UsdtPort)
	var UsdtrpcUser = config.RpcConfig.UsdtRpcUser
	var UsdtrpcPassword = config.RpcConfig.UsdtRpcPwd
	connCfg := &rpc.ConnConfig{
		Host: UsdtHostPort,
		User: UsdtrpcUser,
		Pass: UsdtrpcPassword,
	}
	omni := usdtapi.NewOmniClient(connCfg)

	txid,err:= omni.Send(from,to,uint32(31),strconv.FormatFloat(value,'f',8,64),from)

	if err!=nil {
		return "", err
	}

	//TODO  监听交易

	return txid,nil
}

func ListnerUsdtTransactionByHash(usdtClient *omnilayer.Client,txHash string){
	receiptChan := make(chan omnijson.OmniGettransactionResult)
	CheckUSDTTranscationStatus(receiptChan,usdtClient,txHash,60)
	_ = <- receiptChan

	// TODO 获取交易状态

}

func CheckUSDTTranscationStatus(receiptChan chan omnijson.OmniGettransactionResult,usdtClient *omnilayer.Client ,txHash string,retrySeconds time.Duration){
	go func() {
		for {
			txResult,_ := usdtClient.OmniGetTransaction(txHash)
			if txResult.Confirmations >= 6 {
				receiptChan <- txResult
				break
			} else {
				fmt.Printf("Retry after %d second\n", retrySeconds)
				time.Sleep(retrySeconds * time.Second)
			}
		}
	}()
}

//util
func floatToString(input_num float64) string {
	// to convert a float number to a string
	return strconv.FormatFloat(input_num, 'f', 6, 64)
}

