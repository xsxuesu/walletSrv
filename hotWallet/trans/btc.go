package trans

import (
	"encoding/json"
	"fmt"
	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
	"strconv"
	"time"
	"errors"
	"walletSrv/hotWallet/config"
	"walletSrv/proto"
)

//{
//"txid" : "d54994ece1d11b19785c7248868696250ab195605b469632b7bd68130e880c9a",
//"vout" : 1,
//"address" : "mgnucj8nYqdrPFh2JfZSB1NmUThUGnmsqe",
//"account" : "test label",
//"scriptPubKey" : "76a9140dfc8bafc8419853b34d5e072ad37d1a5159f58488ac",
//"amount" : 0.00010000,
//"confirmations" : 6210,
//"spendable" : true,
//"solvable" : true
//}



//btcEnv = &chaincfg.MainNetParams
var btcEnv = &chaincfg.TestNet3Params
//btcEnv = &chaincfg.RegressionNetParams

var DEFAULT_ACCOUNT = "companyaccount"

type BtcClient struct {
	rpcClient *rpcclient.Client
}

type BtcAccount struct {
	Account  string
	Balance  float64
	Addresses []string
}

// connect to bitcoind with HTTP RPC transport
func NewBtcClient() (*BtcClient, error) {

	var BtcHostPort = fmt.Sprintf("%s:%s",config.RpcConfig.BtcHost,config.RpcConfig.BtcPort)
	var BtcrpcUser = config.RpcConfig.BtcRpcUser
	var BtcrpcPassword = config.RpcConfig.BtcRpcPwd

	connCfg := &rpcclient.ConnConfig{
		Host:         BtcHostPort,
		User:         BtcrpcUser,
		Pass:         BtcrpcPassword,
		HTTPPostMode: true,
		DisableTLS:   true,
	}

	client := &BtcClient{}
	var err error
	client.rpcClient, err = rpcclient.New(connCfg, nil)
	if err != nil {
		return nil, err
	}

	if err != nil {
		return nil, err
	}
	fmt.Printf("network:%d coins=>btc_wallet=>initClinet sccuess.",client.rpcClient.NextID())
	return client, nil
}

// Ping
func (client *BtcClient) Ping() error {
	return client.rpcClient.Ping()
}

// GetBlockCount
func (client *BtcClient) GetBlockCount() (int64, error) {
	return client.rpcClient.GetBlockCount()
}

// Create new account
func (client *BtcClient) CreateNewAccount(account string) (string, error) {
	if len(account) == 0 {
		account = DEFAULT_ACCOUNT
	}
	err := client.rpcClient.CreateNewAccount(account)
	if err != nil {
		return "", err
	}
	return account, nil
}

func (client *BtcClient) InsertToAccount(address string,account string) (string, error) {
	if len(account) == 0 {
		account = DEFAULT_ACCOUNT
	}

	addr,err := decodeAddress(address,btcEnv)
	if err != nil {
		return "", err
	}
	err = client.rpcClient.SetAccount(addr,account)
	if err != nil {
		return "", err
	}
	return account, nil
}

// GetAddress - default address
func (client *BtcClient) GetAddress(account string) (string, error) {
	if len(account) == 0 {
		account = DEFAULT_ACCOUNT
	}
	address, err := client.rpcClient.GetAccountAddress(account)
	if err != nil {
		return "", err
	}
	return address.String(), nil
}

func (client *BtcClient) CreateAccount(account string) (BtcAccount, error) {
	address, err := client.GetAddress(account)
	if err != nil {
		return BtcAccount{}, err
	}
	return BtcAccount{
		Account:   account,
		Balance:   0.0,
		Addresses: []string{address},
	}, nil
}

// GetAccountInfo
func (client *BtcClient) GetAccountInfo(account string, minConf int) (BtcAccount, error) {
	var accountsMap map[string]float64
	var err error

	if accountsMap, err = client.ListAccountsMinConf(minConf); err != nil {
		return BtcAccount{}, err
	}

	balance, found := accountsMap[account]
	if !found {
		return BtcAccount{}, errors.New("")
	}

	addresses, err := client.GetAddressesByAccount(account)
	if err != nil {
		return BtcAccount{}, err
	}

	return BtcAccount{
		Account:   account,
		Balance:   balance,
		Addresses: addresses,
	}, nil
}

// TODO
// GetNewAddress does map to `getnewaddress` rpc call now
// rpcclient doesn't have such golang wrapper func.
func (client *BtcClient) GetNewAddress(account string) (string, error) {
	if len(account) == 0 {
		account = DEFAULT_ACCOUNT
	}
	address, err := client.rpcClient.GetNewAddress(account)
	if err != nil {
		return "", err
	}
	return address.String(), nil
}

// GetAddressesByAccount
func (client *BtcClient) GetAddressesByAccount(account string) ([]string, error) {
	if len(account) == 0 {
		account = DEFAULT_ACCOUNT
	}
	addresses, err := client.rpcClient.GetAddressesByAccount(account)
	if err != nil {
		return []string{}, err
	}

	addrs := make([]string, 0)
	for _, addr := range addresses {
		addrs = append(addrs, addr.String())
	}

	return addrs, nil
}

// ListAccountsMinConf
func (client *BtcClient) ListAccountsMinConf(minConf int) (map[string]float64, error) {
	accounts := make(map[string]float64)

	accountsWithAmount, err := client.rpcClient.ListAccountsMinConf(minConf)
	if err != nil {
		return accounts, err
	}

	for account, amount := range accountsWithAmount {
		accounts[account] = amount.ToBTC()
	}

	return accounts, nil
}

// SendToAddress
func (client *BtcClient) SendToAddress(address string, amount float64) (string, error) {
	decoded, err := decodeAddress(address, btcEnv)
	if err != nil {
		return "", err
	}

	btcAmount, err := convertToBtcAmount(amount)
	if err != nil {
		return "", err
	}

	hash, err := client.rpcClient.SendToAddressComment(decoded, btcAmount, "", "")
	if err != nil {
		return "", err
	}
	return hash.String(), nil
}

// TODO check validity of account and have sufficent balance
func (client *BtcClient) SendFrom(account, address string, amount float64) (string, error) {
	decoded, err := decodeAddress(address, btcEnv)
	if err != nil {
		return "", err
	}

	btcAmount, err := convertToBtcAmount(amount)
	if err != nil {
		return "", err
	}

	hash, err := client.rpcClient.SendFrom(account, decoded, btcAmount)
	if err != nil {
		return "", err
	}
	return hash.String(), nil
}

// Move
func (client *BtcClient) Move(from, to string, amount float64) (bool, error) {
	btcAmount, err := convertToBtcAmount(amount)
	if err != nil {
		return false, err
	}
	return client.rpcClient.Move(from, to, btcAmount)
}

// ListUnspentMin
func (client *BtcClient) ListUnspentMin(minConf int) ([]btcjson.ListUnspentResult, error) {
	return client.rpcClient.ListUnspentMin(minConf)
}

//获取账户余额
func (client *BtcClient)GetBalanceInAddress(address string) (float64, error) {
	balance := float64(0)
	unspents, err := client.GetUnspentByAddress(address)
	if err != nil {
		return balance,err
	}
	for _, v := range unspents {
		balance += v.Amount
	}
	return balance,nil
}

// 获得utxo 数据
func (client *BtcClient)GetUnspentByAddress(address string) (unspents []btcjson.ListUnspentResult, err error) {
	btcAdd, err := btcutil.DecodeAddress(address, btcEnv)
	if err != nil {
		return
	}
	adds := [1]btcutil.Address{btcAdd}
	fmt.Println(adds)
	unspents, err = client.rpcClient.ListUnspentMinMaxAddresses(1, 999999, adds[:])
	fmt.Println(unspents)
	if err != nil {
		return
	}
	return
}

//转账
//addrForm来源地址，addrTo去向地址 value 转账金额 fee
func (client *BtcClient)CreateBtcTrans(serial string,addrFrom string, addrTo string, value float64) (*proto.SignBtcTransRequest,error) {

	unspents, err := client.GetUnspentByAddress(addrFrom)

	if err != nil {
		return &proto.SignBtcTransRequest{},err
	}

	byteSpent,err := json.Marshal(unspents)

	if err != nil {
		return &proto.SignBtcTransRequest{},err
	}

	signBtcReq := &proto.SignBtcTransRequest{
		Serial:serial,
		From:addrFrom,
		To:addrTo,
		Value:strconv.FormatFloat(value,'f',8,64),
		Unspentlist:byteSpent,
	}
	return signBtcReq,nil
}


func (client *BtcClient)ListnerBtcTransByTxHash(serial string,txHash string){
	txHashinfo,err := chainhash.NewHashFromStr(txHash)

	if err != nil {
		return
	}
	// 监听交易状态
	receiptChan:= make(chan *btcjson.GetTransactionResult)
	client.CheckBtcTranscationStatus(receiptChan,txHashinfo,60)
	_ = <-receiptChan


	// TODO   修改状态

}

// btc send raw transaction
func (client *BtcClient)BtcSendRawTrans(serial string , tx []byte)(string,error ){
	var redeemTx wire.MsgTx
	err := json.Unmarshal(tx,&redeemTx)
	if err != nil {
		return "",err
	}
	sendResult,err  := client.rpcClient.SendRawTransaction(&redeemTx,false)
	//sendResult,err  := btcClient.SendRawTransactionAsync(&redeemTx,false).Receive()
	if err != nil {
		return "",err
	}

	return sendResult.String(),nil
}

func (client *BtcClient)CheckBtcTranscationStatus(receiptChan chan *btcjson.GetTransactionResult,txHashinfo *chainhash.Hash,retrySeconds time.Duration){
	go func() {
		for {
			txResult, _ := client.rpcClient.GetTransaction(txHashinfo)
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

// decodeAddress from string to decodedAddress
func decodeAddress(address string, cfg *chaincfg.Params) (btcutil.Address, error) {
	decodedAddress, err := btcutil.DecodeAddress(address, cfg)
	if err != nil {
		return nil, err
	}
	return decodedAddress, nil
}


func convertToBtcAmount(amount float64) (btcutil.Amount, error) {
	return btcutil.NewAmount(amount)
}

