package sign

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcutil"
	"log"
	"math/big"

	"github.com/btcsuite/btcd/wire"
)

var lowFee = big.NewInt(500000)
var highFee = big.NewInt(1000000)

func txToHex(tx *wire.MsgTx) string {
	buf := bytes.NewBuffer(make([]byte, 0, tx.SerializeSize()))
	tx.Serialize(buf)
	return hex.EncodeToString(buf.Bytes())
}

//var btcEvn = &chaincfg.MainNetParams

//var btcEvn = &chaincfg.TestNet3Params

var btcEvn =  &chaincfg.TestNet3Params

type MultiUnspent struct {
	Addr string `json:"addr"`
	UnspentList []btcjson.ListUnspentResult `json:"unspentlist"`
}

//转账
//addrForm来源地址，addrTo去向地址
//transfer 转账金额
//fee 小费
func BtcSignTrans(from string,priK string,to string,value float64,
	unspentbytes []byte) ( []byte, error) {

	var unspentlist []btcjson.ListUnspentResult

	err := json.Unmarshal(unspentbytes,&unspentlist)
	if err != nil {
		return []byte{},err
	}
	fmt.Println("fmt.Println(unspentlist):")
	fmt.Println(unspentlist)

	outsum := float64(0)
	fee := float64(0.005)
	totalTran := float64(0)
	totalTran = value + fee

	// txin签名用script
	var pkscripts [][]byte
	//构造tx
	tx := wire.NewMsgTx(wire.TxVersion)

	for _, v := range unspentlist {
		if v.Amount == 0 {
			continue
		}
		outsum += v.Amount
		//txin输入-------start-----------------
		hash, err := chainhash.NewHashFromStr(v.TxID)
		if err != nil {
			return []byte{},err
		}
		outPoint := wire.NewOutPoint(hash, v.Vout)
		txIn := wire.NewTxIn(outPoint, nil, nil)
		tx.AddTxIn(txIn)

		//设置签名用script
		txinPkScript, errInner := hex.DecodeString(v.ScriptPubKey)
		if errInner != nil {
			return []byte{},errInner
		}
		pkscripts = append(pkscripts, txinPkScript)
	}
	// 判断余额和转出金额比较
	if totalTran > outsum {
		return []byte{}, errors.New("Lack of balance")
	}

	//1 给form----------------找零-------------------
	// 设置 主要 params
	addrf, err := btcutil.DecodeAddress(from, btcEvn)
	if err != nil {
		return []byte{},err
	}
	pkScriptf, err := txscript.PayToAddrScript(addrf)
	if err != nil {
		return []byte{},err
	}

	//tx.AddTxOut(wire.NewTxOut(baf, pkScriptf))
	//outsu.Sub(totalTran)
	//余额 计算
	balance := outsum - totalTran - fee
	fmt.Println("balance")
	fmt.Println(balance)
	tx.AddTxOut(wire.NewTxOut(int64(balance*100000000), pkScriptf))
	//2，给to------------------付钱-----------------

	addrt, errInner := btcutil.DecodeAddress(to, btcEvn)
	if errInner != nil {
		return []byte{},errInner
	}
	pkScriptt, errInner := txscript.PayToAddrScript(addrt)
	if errInner != nil {
		return []byte{},errInner
	}
	tx.AddTxOut(wire.NewTxOut(int64(value*100000000), pkScriptt))

	//-------------------输出填充end------------------------------
	// 签名
	err = sign(tx, priK, pkscripts,value)
	if err != nil {
		fmt.Println("sign err:")
		fmt.Println(err.Error())
		return []byte{},err
	}

	//txSignedHash := txToHex(tx)

	txBytes,err := json.Marshal(tx)
	if err != nil {
		fmt.Println("sign err:")
		fmt.Println(err.Error())
		return []byte{},err
	}
	return txBytes,nil
}


func sign(tx *wire.MsgTx, privKey string, pkScripts [][]byte,value float64) error {
	wif, err := btcutil.DecodeWIF(privKey)
	if err != nil {
		return err
	}
	for i, _ := range tx.TxIn {
		script, err := txscript.SignatureScript(tx, i, pkScripts[i], txscript.SigHashAll, wif.PrivKey, false)
		if err != nil {
			fmt.Println(err.Error())
			return err
		}
		tx.TxIn[i].SignatureScript = script

		flags :=  txscript.StandardVerifyFlags
		vm, err := txscript.NewEngine(pkScripts[i], tx, i,
			flags, nil, nil, int64(value*100000000))

		if err != nil {
			fmt.Println(err.Error())
			return err
		}
		err = vm.Execute()

		if err != nil {
			fmt.Println("err3:")
			fmt.Println(err.Error())
			return err
		}
	}
	log.Println("Transaction successfully signed")
	return nil
}


//签名
//privkey的compress方式需要与TxIn的
func signTx(tx *wire.MsgTx, privKey string, pkScripts [][]byte,value float64) error {
	wif, err := btcutil.DecodeWIF(privKey)
	if err != nil {
		return err
	}
	//for i, _ := range tx.TxIn {
	//	script, err := txscript.SignatureScript(tx, i, pkScripts[i], txscript.SigHashAll, wif.PrivKey, false)
	//	if err != nil {
	//		return err
	//	}
	//	tx.TxIn[i].SignatureScript = script
	//
	//	vm, err := txscript.NewEngine(pkScripts[i], tx, i,
	//		txscript.StandardVerifyFlags, nil, nil, -1)
	//	if err != nil {
	//		return err
	//	}
	//	err = vm.Execute()
	//
	//	if err != nil {
	//		return err
	//	}
	//	log.Println("Transaction successfully signed")
	//}

	script, err := txscript.SignatureScript(tx, 0, pkScripts[0], txscript.SigHashAll, wif.PrivKey, false)
	if err != nil {
		return err
	}
	tx.TxIn[0].SignatureScript = script


	//flags := txscript.StandardVerifyFlags
	//vm, err := txscript.NewEngine(sourceTx.TxOut[0].PkScript, redeemTx, 0, flags, nil, nil, amount)
	//if err != nil {
	//	return Transaction{}, err
	//}
	//if err := vm.Execute(); err != nil {
	//	return Transaction{}, err
	//}

	vm, err := txscript.NewEngine(tx.TxOut[0].PkScript, tx, 0,
		txscript.StandardVerifyFlags, nil, nil, int64(value*100000000))
	if err != nil {
		return err
	}
	err = vm.Execute()

	if err != nil {
		return err
	}
	log.Println("Transaction successfully signed")

	return nil
}



func CreateBtcTxByte(secret string, destination string, amount int64, unspentbytes []byte)([]byte,error){

	var unspentlist []btcjson.ListUnspentResult

	err := json.Unmarshal(unspentbytes,&unspentlist)
	if err != nil {
		return []byte{},err
	}
	wif, err := btcutil.DecodeWIF(secret)
	if err != nil {
		return []byte{}, err
	}
	addresspubkey, err := btcutil.NewAddressPubKey(wif.PrivKey.PubKey().SerializeUncompressed(), btcEvn)
	if err != nil {
		return []byte{}, err
	}
	sourceTx := wire.NewMsgTx(wire.TxVersion)

	UTXOsAmount := big.NewInt(0)
	//VIN
	for _, v := range unspentlist {
		sourceUtxoHash, err := chainhash.NewHashFromStr(v.TxID)
		if err != nil {
			return []byte{}, err
		}
		sourceUtxo := wire.NewOutPoint(sourceUtxoHash, 0)
		sourceTxIn := wire.NewTxIn(sourceUtxo, nil, nil)
		sourceTx.AddTxIn(sourceTxIn)
		//计算总额
		intAmount , _ := new(big.Float).Mul(big.NewFloat(v.Amount),big.NewFloat(100000000.0)).Int64()
		UTXOsAmount.Add(UTXOsAmount,big.NewInt(intAmount))
	}

	// VOUT
	destinationAddress, err := btcutil.DecodeAddress(destination, btcEvn)
	if err != nil {
		return []byte{}, err
	}
	sourceAddress, err := btcutil.DecodeAddress(addresspubkey.EncodeAddress(), btcEvn)
	if err != nil {
		return []byte{}, err
	}
	destinationPkScript, err := txscript.PayToAddrScript(destinationAddress)
	if err != nil {
		return []byte{}, err
	}
	sourcePkScript, err := txscript.PayToAddrScript(sourceAddress)
	if err != nil {
		return []byte{}, err
	}
	//change找零

	sourceTxOut := wire.NewTxOut(amount, destinationPkScript)
	sourceTx.AddTxOut(sourceTxOut)

	changeAmount := new(big.Int).Sub(UTXOsAmount,lowFee)
	changeAmount = new(big.Int).Sub(changeAmount,big.NewInt(amount))

	fmt.Println("change:")
	fmt.Println(changeAmount.Int64())

	// tx out to send change back to us
	changeOutput := wire.NewTxOut(changeAmount.Int64(), sourcePkScript)
	sourceTx.AddTxOut(changeOutput)

	for i , _ := range sourceTx.TxIn{
		sigScript, err := txscript.SignatureScript(sourceTx, i, changeOutput.PkScript, txscript.SigHashAll, wif.PrivKey, false)
		if err != nil {
			return []byte{}, err
		}
		sourceTx.TxIn[i].SignatureScript = sigScript

		//flags := txscript.StandardVerifyFlags
		//vm, err := txscript.NewEngine(changeOutput.PkScript, sourceTx, i, flags, nil, nil, amount)
		//if err != nil {
		//	return []byte{}, err
		//}
		//if err := vm.Execute(); err != nil {
		//	return []byte{}, err
		//}
	}

	buf := bytes.NewBuffer(make([]byte, 0, sourceTx.SerializeSize()))
	sourceTx.Serialize(buf)

	fmt.Printf("Redeem Tx: %v\n", hex.EncodeToString(buf.Bytes()))
	if redeemByte,err := json.Marshal(sourceTx);err !=nil {
		return []byte{},err
	}else{

		return redeemByte,nil
	}
}



//// 多签名转账
func CreateMultiBtcTxByte(mapSecret map[string]string, destination string, amount int64, unspentbytes []byte)([]byte,error){
	var mapWif = map[string]*btcutil.WIF{}

	var unspentlist []MultiUnspent

	err := json.Unmarshal(unspentbytes,&unspentlist)
	if err != nil {
		return []byte{},err
	}

	// 解出address public key
	for k,v := range mapSecret{
		wif, err := btcutil.DecodeWIF(v)
		if err != nil {
			continue
		}
		//addresspubkey, err := btcutil.NewAddressPubKey(wif.PrivKey.PubKey().SerializeUncompressed(), btcEvn)
		//if err != nil {
		//	continue
		//}
		mapWif[k] = wif
	}


	sourceTx := wire.NewMsgTx(wire.TxVersion)

	UTXOsAmount := big.NewInt(0)

	//VIN
	for _, vSpentlist := range unspentlist {
		for _ ,v := range vSpentlist.UnspentList{
			sourceUtxoHash, err := chainhash.NewHashFromStr(v.TxID)
			if err != nil {
				return []byte{}, err
			}
			sourceUtxo := wire.NewOutPoint(sourceUtxoHash, 0)
			//////////////////////ADDRESS TO BYTES
			addrBytes := []byte(vSpentlist.Addr)
			sourceTxIn := wire.NewTxIn(sourceUtxo, addrBytes, nil)
			sourceTx.AddTxIn(sourceTxIn)
			//计算总额
			intAmount , _ := new(big.Float).Mul(big.NewFloat(v.Amount),big.NewFloat(100000000.0)).Int64()
			UTXOsAmount.Add(UTXOsAmount,big.NewInt(intAmount))
		}

	}

	// VOUT
	// 目标转账账号
	destinationAddress, err := btcutil.DecodeAddress(destination, btcEvn)
	if err != nil {
		return []byte{}, err
	}
	// 找零账号
	var changeAddress = ""
	sourceAddress, err := btcutil.DecodeAddress(changeAddress, btcEvn)
	if err != nil {
		return []byte{}, err
	}

	destinationPkScript, err := txscript.PayToAddrScript(destinationAddress)
	if err != nil {
		return []byte{}, err
	}
	sourcePkScript, err := txscript.PayToAddrScript(sourceAddress)
	if err != nil {
		return []byte{}, err
	}
	//change找零

	sourceTxOut := wire.NewTxOut(amount, destinationPkScript)
	sourceTx.AddTxOut(sourceTxOut)

	changeAmount := new(big.Int).Sub(UTXOsAmount,lowFee)
	changeAmount = new(big.Int).Sub(changeAmount,big.NewInt(amount))

	fmt.Println("change:")
	fmt.Println(changeAmount.Int64())

	// tx out to send change back to us
	changeOutput := wire.NewTxOut(changeAmount.Int64(), sourcePkScript)
	sourceTx.AddTxOut(changeOutput)

	for i , v := range sourceTx.TxIn{
		addr := string(v.SignatureScript)
		sigScript, err := txscript.SignatureScript(sourceTx, i, changeOutput.PkScript, txscript.SigHashAll, mapWif[addr].PrivKey, false)
		if err != nil {
			return []byte{}, err
		}
		sourceTx.TxIn[i].SignatureScript = sigScript

	}

	buf := bytes.NewBuffer(make([]byte, 0, sourceTx.SerializeSize()))
	sourceTx.Serialize(buf)

	fmt.Printf("Redeem Tx: %v\n", hex.EncodeToString(buf.Bytes()))
	if redeemByte,err := json.Marshal(sourceTx);err !=nil {
		return []byte{},err
	}else{

		return redeemByte,nil
	}
}


/**


func CreateTransaction( secret string, destination string, amount int64, unspentbytes []byte) (Transaction, error) {

	var transaction Transaction

	var unspentlist []btcjson.ListUnspentResult

	err := json.Unmarshal(unspentbytes,&unspentlist)
	if err != nil {
		return transaction,err
	}


	wif, err := btcutil.DecodeWIF(secret)
	if err != nil {
		return Transaction{}, err
	}
	addresspubkey, err := btcutil.NewAddressPubKey(wif.PrivKey.PubKey().SerializeUncompressed(), btcEvn)
	if err != nil {
		return Transaction{}, err
	}
	sourceTx := wire.NewMsgTx(wire.TxVersion)

	//VIN
	for _, v := range unspentlist {
		sourceUtxoHash, err := chainhash.NewHashFromStr(v.TxID)
		if err != nil {
			return Transaction{}, err
		}
		sourceUtxo := wire.NewOutPoint(sourceUtxoHash, 0)
		sourceTxIn := wire.NewTxIn(sourceUtxo, nil, nil)
		sourceTx.AddTxIn(sourceTxIn)
	}

	// VOUT
	destinationAddress, err := btcutil.DecodeAddress(destination, btcEvn)
	sourceAddress, err := btcutil.DecodeAddress(addresspubkey.EncodeAddress(), btcEvn)
	if err != nil {
		return Transaction{}, err
	}
	destinationPkScript, err := txscript.PayToAddrScript(destinationAddress)
	if err != nil {
		return Transaction{}, err
	}
	sourcePkScript, err := txscript.PayToAddrScript(sourceAddress)
	if err != nil {
		return Transaction{}, err
	}
	sourceTxOut := wire.NewTxOut(amount, sourcePkScript)
	//sourceTx.AddTxIn(sourceTxIn)
	sourceTx.AddTxOut(sourceTxOut)

	sourceTxHash := sourceTx.TxHash()

	// Fee
	redeemTx := wire.NewMsgTx(wire.TxVersion)
	destinationUtxo := wire.NewOutPoint(&sourceTxHash, 0)
	redeemTxIn := wire.NewTxIn(destinationUtxo, nil, nil)
	redeemTx.AddTxIn(redeemTxIn)
	redeemTxOut := wire.NewTxOut(amount, destinationPkScript)
	redeemTx.AddTxOut(redeemTxOut)
	sigScript, err := txscript.SignatureScript(redeemTx, 0, sourceTx.TxOut[0].PkScript, txscript.SigHashAll, wif.PrivKey, false)
	if err != nil {
		return Transaction{}, err
	}
	redeemTx.TxIn[0].SignatureScript = sigScript
	flags := txscript.StandardVerifyFlags
	vm, err := txscript.NewEngine(sourceTx.TxOut[0].PkScript, redeemTx, 0, flags, nil, nil, amount)
	if err != nil {
		return Transaction{}, err
	}
	if err := vm.Execute(); err != nil {
		return Transaction{}, err
	}
	var unsignedTx bytes.Buffer
	var signedTx bytes.Buffer
	sourceTx.Serialize(&unsignedTx)
	redeemTx.Serialize(&signedTx)
	transaction.TxId = sourceTxHash.String()
	transaction.UnsignedTx = hex.EncodeToString(unsignedTx.Bytes())
	transaction.Amount = amount
	transaction.SignedTx = hex.EncodeToString(signedTx.Bytes())
	transaction.SourceAddress = sourceAddress.EncodeAddress()
	transaction.DestinationAddress = destinationAddress.EncodeAddress()
	fmt.Println(transaction.SignedTx)
	return transaction, nil
}
 */