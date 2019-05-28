package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"math/big"
	"net"
	"os"
	"strconv"
	"walletSrv/coldWallet/config"
	"walletSrv/coldWallet/services"
	"walletSrv/proto"
)

func init(){
	config.InitConf(os.Args)
	// 检查表格
	//db,err := config.GetMySqlDb()
	//if err != nil {
	//	log.Fatal("connect to mysql err :",err.Error())
	//}
	//err = config.CreateAddrTabs(db)
	//if err != nil {
	//	log.Fatal("create Address tables err :",err.Error())
	//}
	//err = config.CreateHDCountTabs(db)
	//if err != nil {
	//	log.Fatal("create HdCount tables err :",err.Error())
	//}
	//
	//serialModel := model.SerialModel{db}
	//// 创建交易记录
	//err = serialModel.CheckTable(config.Eth)
	//if err != nil {
	//	log.Fatal("create serialEth tables err :",err.Error())
	//}
	//
	//err = serialModel.CheckTable(config.Btc)
	//if err != nil {
	//	log.Fatal("create serialEth tables err :",err.Error())
	//}
	//
	//err = serialModel.CheckTable(config.Usdt)
	//if err != nil {
	//	log.Fatal("create serialEth tables err :",err.Error())
	//}
	//
	//// 创建解密表
	//decryptModel := model.DecryptModel{db}
	//
	//err = decryptModel.CheckTable(config.Eth)
	//if err != nil {
	//	log.Fatal("create decryptEth tables err :",err.Error())
	//}
	//
	//err = decryptModel.CheckTable(config.Btc)
	//if err != nil {
	//	log.Fatal("create decryptEth tables err :",err.Error())
	//}
	//
	//err = decryptModel.CheckTable(config.Usdt)
	//if err != nil {
	//	log.Fatal("create decryptEth tables err :",err.Error())
	//}
}

func main()  {

	listener,err := net.Listen("tcp",":8080")
	if err != nil {
		log.Fatalln("net listen 8080 error:",err.Error())
	}
	srv := grpc.NewServer()
	proto.RegisterWalletServer(srv,&WalletServer{})
	srv.Serve(listener)
}

type WalletServer struct {}

func (ws *WalletServer)InsertAddress(ctx context.Context,request *proto.InsertAddressRequest)(*proto.InsertAddressResponse,error){
	reps :=  new(proto.InsertAddressResponse)

	err := services.InsertAddress(request.Cointype,request.Address,request.Privk)
	if err != nil {
		return &proto.InsertAddressResponse{
			Address:request.Address,
			Success:"0",
		},err
	}
	reps.Address = request.Address
	reps.Success = "1"
	return reps, nil
}

func (ws *WalletServer)GetHDAddress(ctx context.Context,requst *proto.GetHDAddressRequest)(*proto.GetHDAddressResponse,error){
	reps :=  new(proto.GetHDAddressResponse)

	addr,err := services.CreateHdAddress(requst.Cointype)
	if err != nil {
		return &proto.GetHDAddressResponse{},err
	}
	reps.Addr = addr

	return reps, nil
}

func (ws *WalletServer)SignEthTrans(ctx context.Context,requst *proto.SignEthTransRequest)(*proto.SignEthTransResponse,error){
	resp := new(proto.SignEthTransResponse)

	//serial string ,from string,to string,nonce uint64,value int64,
	//gasLimit uint64,gasPrice int64,chainid int64
	nonce,err := strconv.ParseUint(requst.Nonce,10,64)
	if err != nil {
		return &proto.SignEthTransResponse{Success:"0",Rawhash:"",Serial:requst.Serial},err
	}

	value,b := new(big.Int).SetString(requst.Value,10)

	if b == false {
		return &proto.SignEthTransResponse{Success:"0",Rawhash:"",Serial:requst.Serial},err
	}

	gaslimit,err := strconv.ParseUint(requst.Gaslimit,10,64)
	if err != nil {
		return &proto.SignEthTransResponse{Success:"0",Rawhash:"",Serial:requst.Serial},err
	}

	gasprice,err := strconv.ParseInt(requst.Gasprice,10,64)
	if err != nil {
		return &proto.SignEthTransResponse{Success:"0",Rawhash:"",Serial:requst.Serial},err
	}

	chainid,err := strconv.ParseInt(requst.Chainid,10,64)
	if err != nil {
		return &proto.SignEthTransResponse{Success:"0",Rawhash:"",Serial:requst.Serial},err
	}

	hash,err := services.SignEthTrans(requst.Serial,requst.From,requst.To,nonce,value,gaslimit,gasprice,chainid,requst.Contract)
	if err != nil {
		return &proto.SignEthTransResponse{Success:"0",Rawhash:"",Serial:requst.Serial},err
	}
	resp.Rawhash = hash
	resp.Serial = requst.Serial
	resp.Success = "1"
	return resp,nil
}

func (ws *WalletServer)SignBtcTrans(ctx context.Context,request *proto.SignBtcTransRequest)(*proto.SignBtcTransResponse,error){
	resp := new(proto.SignBtcTransResponse)

	value,err := strconv.ParseFloat(request.Value,64)
	if err != nil {
		return &proto.SignBtcTransResponse{},err
	}
	fmt.Println("request:")
	fmt.Println(request)
	signedTx,err := services.SignBtcTrans(request.Serial,request.From,request.To,value,request.Unspentlist)

	if err != nil {
		return &proto.SignBtcTransResponse{},err
	}

	resp.Serial = request.Serial
	resp.Signtx = signedTx

	return resp,nil
}

func (ws *WalletServer)SignUsdtTrans(ctx context.Context,request *proto.SignUsdtTransRequest)(*proto.SignUsdtTransResponse,error){
	resp := new(proto.SignUsdtTransResponse)

	value,err := strconv.ParseFloat(request.Value,64)
	if err != nil {
		return &proto.SignUsdtTransResponse{},err
	}
	fmt.Println("usdt request:")
	fmt.Println(request)
	signedTx,err := services.SignUsdtTrans(request.Serial,request.From,request.To,value,request.Unspentlist)

	if err != nil {
		return &proto.SignUsdtTransResponse{},err
	}

	resp.Serial = request.Serial
	resp.Signedhash = signedTx

	return resp,nil

}

func (ws *WalletServer)UpdateTransStatus(ctx context.Context,requst *proto.UpdateTransStatusRequest)(*proto.UpdateTransStatusResponse,error){

	resp := new(proto.UpdateTransStatusResponse)

	err := services.UpdateTxId(requst.Cointype,requst.Serial,requst.Txid,requst.Status)

	if err != nil {
		return &proto.UpdateTransStatusResponse{
			Txid : requst.Txid,
			Serial:requst.Serial,
			Cointype :requst.Cointype,
			Success : "0",
		},err
	}

	resp.Txid = requst.Txid
	resp.Serial = requst.Serial
	resp.Cointype = requst.Cointype
	resp.Success = "1"

	return resp,nil
}

func (ws *WalletServer)Collection(ctx context.Context,request *proto.CollectionRequest)(*proto.CollectionResponse,error){
	resp := new(proto.CollectionResponse)

	switch request.Cointype {
	case config.Eth:
		start,_:=strconv.Atoi(request.Start)
		end , _:=strconv.Atoi(request.End)

		list,isfinish,err := services.CollectionEth(config.Eth,start,end)

		if err != nil {
			return &proto.CollectionResponse{},err
		}
		resp.Cointype=config.Eth
		resp.Isfinish = strconv.FormatBool(isfinish)
		byteList, err := json.Marshal(list)
		if err != nil {
			return &proto.CollectionResponse{},err
		}
		resp.Addrlist = byteList

	default:
		return resp,errors.New("coin type error")
	}

	return resp,nil
}


