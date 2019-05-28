package sign

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
	"math/big"
	"strconv"
)
var usdtEvn =  &chaincfg.TestNet3Params

var usdtLowFee = big.NewInt(500000)
var usdtHighFee = big.NewInt(1000000)


type ListUnspentResult struct {
	Tx            string  `json:"txid"`
	Address       string  `json:"address"`
	ScriptPubKey  string  `json:"scriptPubKey"`
	RedeemScript  string  `json:"redeemScript"`
	Amount        float64 `json:"amount"`
	Confirmations int64   `json:"confirmations"`
	Vout          uint32  `json:"vout"`
	Spendable     bool    `json:"spendable"`
	Solvable      bool    `json:"solvable"`
}

func addPreZero(num int64) string  {
	strNum := strconv.FormatInt(num,16)
	t := len(strNum)
	s := ""
	for i:=0;i<16-t;i++ {
		s += "0"
	}
	return fmt.Sprintf("%s%s",s,strNum)
}

func CreateUsdtTxBytes(priK string, to string,value int64,unspentbytes []byte)(string,error){

	fundValue := big.NewInt(546)
	amount := big.NewInt(value)

	var unspentlist []ListUnspentResult
	err := json.Unmarshal(unspentbytes,&unspentlist)
	if err != nil {
		return "",err
	}
	// 私钥
	wif, err := btcutil.DecodeWIF(priK)
	if err != nil {
		return "", err
	}
	addresspubkey, err := btcutil.NewAddressPubKey(wif.PrivKey.PubKey().SerializeUncompressed(), btcEvn)
	if err != nil {
		return "", err
	}
	sourceTx := wire.NewMsgTx(31) //wire.TxVersion

	UTXOsAmount := big.NewInt(0)
	//VIN
	for _, v := range unspentlist {
		sourceUtxoHash, err := chainhash.NewHashFromStr(v.Tx)
		if err != nil {
			return "", err
		}
		sourceUtxo := wire.NewOutPoint(sourceUtxoHash, 0)
		sourceTxIn := wire.NewTxIn(sourceUtxo, nil, nil)
		sourceTx.AddTxIn(sourceTxIn)
		//计算总额
		intAmount , _ := new(big.Float).Mul(big.NewFloat(v.Amount),big.NewFloat(100000000.0)).Int64()
		UTXOsAmount.Add(UTXOsAmount,big.NewInt(intAmount))
	}

	// VOUT
	destinationAddress, err := btcutil.DecodeAddress(to, btcEvn)
	if err != nil {
		return "", err
	}
	sourceAddress, err := btcutil.DecodeAddress(addresspubkey.EncodeAddress(), btcEvn)
	if err != nil {
		return "", err
	}
	destinationPkScript, err := txscript.PayToAddrScript(destinationAddress)
	if err != nil {
		return "", err
	}
	sourcePkScript, err := txscript.PayToAddrScript(sourceAddress)
	if err != nil {
		return "", err
	}
	//change找零

	sourceTxOut := wire.NewTxOut(amount.Int64(), destinationPkScript)
	sourceTx.AddTxOut(sourceTxOut)

	changeAmount := new(big.Int).Sub(UTXOsAmount,usdtLowFee)
	changeAmount = new(big.Int).Sub(changeAmount,amount)
	changeAmount = new(big.Int).Sub(changeAmount,fundValue)

	fmt.Println("change:")
	fmt.Println(changeAmount.Int64())

	// tx out to send change back to us
	changeOutput := wire.NewTxOut(changeAmount.Int64(), sourcePkScript)
	sourceTx.AddTxOut(changeOutput)

	// Omni string
	omniOutput := fmt.Sprintf("6f6d6e69000000000000001f%s",addPreZero(amount.Int64()))
	omniByte := []byte(omniOutput)
	opReturnScript, err :=
		txscript.NewScriptBuilder().AddOp(txscript.OP_RETURN).AddData(omniByte).Script()
	if err != nil {
		return "",err
	}
	sourceTx.AddTxOut(wire.NewTxOut(0,opReturnScript))

	for i , _ := range sourceTx.TxIn{
		sigScript, err := txscript.SignatureScript(sourceTx, i, changeOutput.PkScript, txscript.SigHashAll, wif.PrivKey, false)
		if err != nil {
			return "", err
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

	return hex.EncodeToString(buf.Bytes()),nil
	//if redeemByte,err := json.Marshal(sourceTx);err !=nil {
	//	return []byte{},err
	//}else{
	//	return redeemByte,nil
	//}
}
