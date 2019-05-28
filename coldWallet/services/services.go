package services

import (
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil"
	"github.com/ethereum/go-ethereum/crypto"
	"math/big"
	"strconv"
	"strings"
	"walletSrv/coldWallet/config"
	"walletSrv/coldWallet/hdaddr"
	"walletSrv/coldWallet/model"
	"walletSrv/coldWallet/sign"
	"walletSrv/coldWallet/utils"
)

var btcEvn =  &chaincfg.TestNet3Params

func SignUsdtTrans(serial string ,from string,to string,value float64,
	unspentlist []byte)(string , error){
	// 测试 需要手动添加数据库
	db,err :=config.GetMySqlDb()
	if err != nil {
		fmt.Println(err.Error())
	}
	defer db.Close()

	addrModel := model.AddrModel{db}

	btcAddr,err := addrModel.FindByAddress(config.Usdt,from)
	if err != nil {
		return "",err
	}

	// 保存解密记录
	decryptModel := model.DecryptModel{db}

	prik, err := utils.DecodePriByPub(from,btcAddr.PrivateKey,serial,config.Usdt,from,to,decryptModel)
	if err != nil {
		return "",err
	}
	transtr,err := sign.CreateUsdtTxBytes(prik,to,int64(value*100000000),unspentlist)
	if err != nil {
		return "",err
	}
	// 保存serial 数据
	searilEntity := model.SerialEntity{
		CoinType:config.Usdt,
		SerialNo:serial,
		F:from,
		T:to,
		Value:strconv.FormatInt(0,10),
		Fee:uint64(0),
		Status:config.Pending,
	}

	// 记录交易信息
	serialModel := model.SerialModel{db}

	err = serialModel.Insert(searilEntity)
	if err != nil {
		return "",err
	}

	return transtr, nil
}

func SignBtcTrans(serial string ,from string,to string,value float64,
	unspentlist []byte)([]byte,error){
	// 测试 需要手动添加数据库
	db,err :=config.GetMySqlDb()
	if err != nil {
		fmt.Println(err.Error())
	}
	defer db.Close()

	addrModel := model.AddrModel{db}

	btcAddr,err := addrModel.FindByAddress(config.Btc,from)
	if err != nil {
		return []byte{},err
	}

	// 保存解密记录
	decryptModel := model.DecryptModel{db}

	prik, err := utils.DecodePriByPub(from,btcAddr.PrivateKey,serial,config.Btc,from,to,decryptModel)
	if err != nil {
		return []byte{},err
	}
	tranByte,err := sign.CreateBtcTxByte(prik,to,int64(value*100000000),unspentlist)
	if err != nil {
		return []byte{},err
	}
	// 保存serial 数据
	searilEntity := model.SerialEntity{
		CoinType:config.Btc,
		SerialNo:serial,
		F:from,
		T:to,
		Value:strconv.FormatInt(0,10),
		Fee:uint64(0),
		Status:config.Pending,
	}

	// 记录交易信息
	serialModel := model.SerialModel{db}

	err = serialModel.Insert(searilEntity)
	if err != nil {
		return []byte{},err
	}

	return tranByte, nil
}

func SignEthTrans(serial string ,from string,to string,nonce uint64,value *big.Int,
	gasLimit uint64,gasPrice int64,chainid int64,contract string)(string,error){

	// 测试 需要手动添加数据库
	db,err :=config.GetMySqlDb()
	if err != nil {
		fmt.Println(err.Error())
	}
	defer db.Close()

	addrModel := model.AddrModel{db}

	ethAddr,err := addrModel.FindByAddress(config.Eth,from)
	if err != nil {
		return "",err
	}

	// 保存解密记录
	decryptModel := model.DecryptModel{db}

	prik, err := utils.DecodePriByPub(from,ethAddr.PrivateKey,serial,config.Eth,from,to,decryptModel)

	if err != nil {
		return "",err
	}

	// value int64,
	// gasLimit uint64,
	// gasPrice int64,
	// chainID int64

	hash,err := sign.EthSignTx(prik,nonce,to,value,gasLimit,gasPrice,chainid,contract)

	fee := big.NewInt(int64(gasLimit)*gasPrice).Uint64()

	if err != nil {
		return "",err
	}

	searilEntity := model.SerialEntity{
		CoinType:config.Eth,
		SerialNo:serial,
		F:from,
		T:to,
		Value:value.String(),
		Fee:fee,
		TxId:"",
		Status:config.Pending,
	}

	serialModel := model.SerialModel{db}

	err = serialModel.Insert(searilEntity)
	if err != nil {
		return "",err
	}

	return hash,nil
}

func PriCheckAddress(prik string,addr string)bool{
	priByte , err := hex.DecodeString(prik)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	prikey,err := crypto.ToECDSA(priByte)

	if err != nil {
		fmt.Println(err.Error())
		return false
	}

	add := crypto.PubkeyToAddress(prikey.PublicKey).Hex()

	if strings.ToLower(add) == strings.ToLower(addr) {
		return true
	}
	return false
}

func BtcPriCheckAddress(prik string , addr string) bool{
	wif,err := btcutil.DecodeWIF(prik)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	addresspubkey, err := btcutil.NewAddressPubKey(wif.PrivKey.PubKey().SerializeCompressed(), btcEvn)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	fmt.Println(addresspubkey.EncodeAddress())
	if strings.ToLower(addresspubkey.EncodeAddress()) == strings.ToLower(addr){
		return true
	}
	return false
}

func InsertAddress(cointype string,addr string , prik string) error{
	// 测试 需要手动添加数据库
	db,err := config.GetMySqlDb()
	if err != nil {
		fmt.Println(err.Error())
	}

	defer  db.Close()

	addrModel := model.AddrModel{db}

	if !(cointype == config.Btc ||  cointype == config.Usdt || cointype == config.Eth){
		return  errors.New("cointype error, please check coin type in (btc,usdt,eth)")
	}


	// check prikey

	switch cointype {
	case config.Eth:
		if !PriCheckAddress(prik,addr) {
			return  errors.New("private key and address not correct")
		}

	case config.Btc:

		if !BtcPriCheckAddress(prik,addr) {
			return  errors.New("private key and address not correct")
		}

	case config.Usdt:
	}

	// encrypt private key
	ePrik,err :=  utils.EncodePriByPub(addr,prik)
	if err != nil {
		return err
	}

	addrO := model.AddressEntity{
		Address:addr,
		PrivateKey:ePrik,
	}
	// update address table info
	err = addrModel.Insert(cointype,false,addrO)
	if err != nil {
		return err
	}

	return nil
}

func CreateHdAddress(cointype string) (string,error){
	if !(cointype == config.Btc ||  cointype == config.Usdt || cointype == config.Eth){
		return "" , errors.New("cointype error, please check coin type in (btc,usdt,eth)")
	}
	// 测试 需要手动添加数据库
	db,err :=config.GetMySqlDb()
	if err != nil {
		fmt.Println(err.Error())
	}
	defer db.Close()

	hdModel := model.HdCountModel{db}
	index,err := hdModel.FindByType(cointype)

	// can't get cointype index ,insert cointype index =0
	if err != nil {
		index = 0
		err = hdModel.InsertHdCount(cointype,fmt.Sprintf("%d",index))
		if err != nil {
			return "",err
		}
	}

	var addr,pri string

	if (cointype == config.Btc){
		addr,pri,err = hdaddr.GetBtcHdAddress(index)
		if err != nil {
			return "",err
		}
	}

	if (cointype == config.Eth){
		addr,pri,err = hdaddr.GetEthHdAddress(index)
		if err != nil {
			return "",err
		}
	}

	if (cointype == config.Usdt){
		addr,pri,err = hdaddr.GetUSDTHdAddress(index)
		if err != nil {
			return "",err
		}
	}

	// encrypt private key
	ePrik,err :=  utils.EncodePriByPub(addr,pri)
	if err != nil {
		return "",err
	}

	addrO := model.AddressEntity{
		Address:addr,
		PrivateKey:ePrik,
	}
	// update address table info
	addrModel := model.AddrModel{db}
	err = addrModel.Insert(cointype,true,addrO)
 
	if err != nil {
		return "",err
	}
	// update hdcount table
	index = index + 1
	err = hdModel.UpdateHdCount(cointype, fmt.Sprintf("%d",index) )
	if err != nil {
		return "",err
	}

	return addr,nil
}


func UpdateTxId(cointype string,serial string, txid string,status string)error{
	db,err :=config.GetMySqlDb()
	if err != nil {
		fmt.Println(err.Error())
	}
	defer db.Close()

	serialModel := model.SerialModel{db}
	serialEntity,_,err := serialModel.FindBySerialNo(serial,cointype)
	if err != nil {
		return err
	}
	serialEntity.TxId = txid
	serialEntity.Status = status
	err = serialModel.UpdateTxId(serialEntity)
	if err != nil {
		return err
	}
	return nil
}


func CollectionEth(cointype string,start int,end int)(model.AddrList,bool,error){
	db,err := config.GetMySqlDb()
	if err != nil {
		fmt.Println(err.Error())
	}
	defer  db.Close()

	addrModel := model.AddrModel{db}
	addList,isFinish,err := addrModel.FindRange(cointype,start,end)
	if err != nil {
		return model.AddrList{},false,err
	}
	return addList,isFinish,nil
}