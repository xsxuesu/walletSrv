package call

import (
	"context"
	"fmt"
	"log"
	"google.golang.org/grpc"
	"strconv"
	"walletSrv/hotWallet/config"
	"walletSrv/hotWallet/constant"
	"walletSrv/hotWallet/model"
	"walletSrv/hotWallet/trans"
	"walletSrv/hotWallet/utils"
	"walletSrv/proto"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/common"
)

var wClient proto.WalletClient

type WClient struct {
	Client proto.WalletClient
} 


func InitClient(host,port string)(WClient){
	hostport := fmt.Sprintf("%s:%s",host,port)

	conn , err := grpc.Dial(hostport,grpc.WithInsecure())
	if err != nil {
		log.Fatal("connet grpc server:",err.Error())
	}

	client := WClient{}

	wClient = proto.NewWalletClient(conn)

	client.Client = wClient

	return client
}

func (wClient WClient)CallCollection(cointype string,start int,end int)(*proto.CollectionResponse,error){
	collectReq := &proto.CollectionRequest{
		Cointype:cointype,
		Start:strconv.Itoa(start),
		End:strconv.Itoa(end),
	}
	return wClient.Client.Collection(context.Background(),collectReq)
}

func (wClient WClient)CallFetchAddress(cointype string)(*proto.GetHDAddressResponse,error){

	getAddr := &proto.GetHDAddressRequest{
		Cointype:cointype,
	}
	return wClient.Client.GetHDAddress(context.Background(),getAddr)
}

func (wClient WClient)CallInsertAddress(cointype string , address string, prik string)(*proto.InsertAddressResponse,error){
	addr := &proto.InsertAddressRequest{
		Cointype:cointype,
		Address:address,
		Privk:prik,
	}
	return wClient.Client.InsertAddress(context.Background(),addr)
}

func (wClient WClient)CallSignEthTx(serial string,from string,to string,contract string,value float64)(model.TransferResp,error){

	ethTransferResp := model.TransferResp{
		Serial:serial,
		Success:false,
	}

	signEthTransReq,err := trans.CreateEthTrans(serial,from,to,contract,value)

	if err !=  nil {
		return ethTransferResp,err
	}
	signEthTransResp, err := wClient.Client.SignEthTrans(context.Background(),signEthTransReq)

	if err !=  nil {
		return ethTransferResp,err
	}

	serialNo,txid,err := trans.SendRawHash(signEthTransResp.Serial,signEthTransResp.Rawhash)

	if err !=  nil {
		return ethTransferResp,err
	}

	ethTransferResp.Txid = txid
	ethTransferResp.Success = true
	ethTransferResp.Status = constant.Pending

	go func() {
		wClient.CallUpdateTransStatus(serialNo,txid,constant.Eth,constant.Pending)
	}()


	go wClient.ListnerEtcTrans(serialNo,txid)


	return ethTransferResp ,nil
}

func (wClient WClient)CallUpdateTransStatus(serial string,txid string,cointype string,status string){
	updateReq := &proto.UpdateTransStatusRequest{
		Serial:serial,
		Cointype:cointype,
		Txid:txid,
		Status:status,
	}
	updateResp,err := wClient.Client.UpdateTransStatus(context.Background(),updateReq)
	if err !=  nil {
		fmt.Println(err)
	}
	fmt.Println(updateResp.Txid)
	fmt.Println(updateResp.Success)
}

func (wClient WClient)CallSignBtcTx(serial string,from string,to string,value float64)(model.TransferResp,error){
	btcTransferResp := model.TransferResp{
		Serial:serial,
		Success:false,
	}

	//创建交易
	rpcclient,err  := trans.NewBtcClient()
	if err != nil {
		fmt.Println(err.Error())
	}

	signBtcTransReq,err := rpcclient.CreateBtcTrans(serial,from,to,value)

	if err !=  nil {
		return btcTransferResp,err
	}
	// 签名交易
	signBtcTransResp, err := wClient.Client.SignBtcTrans(context.Background(),signBtcTransReq)
	if err !=  nil {
		return btcTransferResp,err
	}

	// 发送交易hash
	txid , err := rpcclient.BtcSendRawTrans(serial,signBtcTransResp.Signtx)
	if err !=  nil {
		return btcTransferResp,err
	}

	btcTransferResp.Success = true
	btcTransferResp.Status = constant.Pending
	btcTransferResp.Txid = txid
	btcTransferResp.CoinType = constant.Btc
	//signBtcTransResp.Signedhash
	fmt.Println(signBtcTransResp.Serial)
	fmt.Println(signBtcTransResp.Signtx)


	return btcTransferResp,nil
}

func (wClient WClient)CallSignUsdtTx(serial string,from string, to string, value float64)( model.TransferResp,error){
	usdtTransferResp := model.TransferResp{
		Serial:serial,
		Success:false,
	}

	signUsdtTransReq,err := trans.CreateUsdtTrans(serial,from,to,value)
	if err !=  nil {
		return usdtTransferResp,err
	}
	signUsdtTransResp,err := wClient.Client.SignUsdtTrans(context.Background(),signUsdtTransReq)
	if err !=  nil {
		return usdtTransferResp,err
	}
	fmt.Println(signUsdtTransResp.Serial)
	fmt.Println(signUsdtTransResp.Signedhash)


//	发送交易hash
	usdtClient,err := trans.NewUSDTClient(int64(31))
	if err !=  nil {
		return usdtTransferResp,err
	}
	hash,err := usdtClient.SendUsdtSignTX(signUsdtTransResp.Signedhash)
	if err != nil {
		return usdtTransferResp,err
	}
	usdtTransferResp.Success = true
	usdtTransferResp.Txid = hash
	usdtTransferResp.CoinType = constant.Usdt
	usdtTransferResp.Status = constant.Pending
	return usdtTransferResp,nil

}

func (wClient WClient)ListnerEtcTrans(serial string,txid string)(error){
	var EthClientHost = fmt.Sprintf("http://%s:%s",config.RpcConfig.EthHost,config.RpcConfig.EthPort)

	client, err := trans.Connect(EthClientHost)
	if err != nil {
		return err
	}
	// check transaction receipt
	receiptChan := make(chan *types.Receipt)
	txHash:= common.HexToHash(txid)
	client.CheckTransaction(context.TODO(), receiptChan, txHash, 60)
	result := <-receiptChan
	//receipt.Status
	//	TODO 修改 交易状态
	if result.Status == uint64(1) {
		wClient.CallUpdateTransStatus(serial,txid,constant.Eth,constant.Executed)
		go func() { // callback
			utils.CallRequest(constant.Eth,txid,constant.Executed)
		}()
	} else{
		wClient.CallUpdateTransStatus(serial,txid,constant.Eth,constant.Failed)
		go func() { // callback
			utils.CallRequest(constant.Eth,txid,constant.Failed)
		}()
	}
	// 关闭 rpc eth client
	defer client.RpcClient.Close()
	defer client.EthClient.Close()
	return nil
}