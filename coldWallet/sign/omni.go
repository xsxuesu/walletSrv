package sign


import (
	"bytes"
	"encoding/json"
	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/btcsuite/btcutil"
	"io/ioutil"
	"net/http"
	"sort"
)

type Response struct {
	Result json.RawMessage   `json:"result"`
	Error  *btcjson.RPCError `json:"error"`
}

type PrevTx struct {
	TxID         string  `json:"txid"`
	Vout         uint32  `json:"vout"`
	ScriptPubKey string  `json:"scriptPubKey"`
	Value        float64 `json:"value"`
}

type CreatePayloadSimpleSendCmd struct {
	PropertyID int64
	Amount    string
}

func NewCreatePayloadSimpleSendCmd(propertyID int64, amount string) *CreatePayloadSimpleSendCmd {
	return &CreatePayloadSimpleSendCmd{
		PropertyID: propertyID,
		Amount:    amount,
	}
}

type CreateRawTxOpReturnCmd struct {
	RawTx string
	Payload string
}

func NewCreateRawTxOpReturnCmd(rawTx, payload string) *CreateRawTxOpReturnCmd {
	return &CreateRawTxOpReturnCmd{
		RawTx: rawTx,
		Payload: payload,
	}
}

type CreateRawTxReferenceCmd struct {
	RawTx string
	Destination string
	Amount float64
}

func NewCreateRawTxReferenceCmd(rawtx, destination string, amount float64) *CreateRawTxReferenceCmd {
	return &CreateRawTxReferenceCmd{
		RawTx: rawtx,
		Destination: destination,
		Amount: amount,
	}
}

type CreateRawTxChangeCmd struct {
	RawTx string
	PrevTxs []PrevTx
	Destination string
	Fee float64
}

func NewCreateRawTxChangeCmd(rawtx string, prevtxs []PrevTx, destination string, fee float64) *CreateRawTxChangeCmd {
	return &CreateRawTxChangeCmd{
		RawTx: rawtx,
		PrevTxs: prevtxs,
		Destination: destination,
		Fee: fee,
	}
}



type omniClient struct {
	*rpcclient.Client
	url string
}

func NewOmniClient() (*omniClient, error) {
	c, err := rpcclient.New(nil, nil)
	if err != nil {
		return nil, err
	}
	o := &omniClient{
		Client:c,
	}
	return o, nil
}

func (o *omniClient) init() {
	flags := btcjson.UFWalletOnly
	btcjson.MustRegisterCmd("omni_createpayload_simplesend", (*CreatePayloadSimpleSendCmd)(nil), flags)
	btcjson.MustRegisterCmd("omni_createrawtx_opreturn", (*CreateRawTxOpReturnCmd)(nil), flags)
	btcjson.MustRegisterCmd("omni_createrawtx_reference", (*CreateRawTxReferenceCmd)(nil), flags)
	btcjson.MustRegisterCmd("omni_createrawtx_change", (*CreateRawTxChangeCmd)(nil), flags)

}

func (o *omniClient) SendTransaction(private string, toAddr btcutil.Address, amount btcutil.Amount) error {
	wif, err := btcutil.DecodeWIF(private)
	if err != nil {
		return err
	}
	addPub, err := btcutil.NewAddressPubKey(wif.PrivKey.PubKey().SerializeCompressed(), &chaincfg.MainNetParams)
	if err != nil {
		return err
	}
	unspents, err := o.ListUnspentMinMaxAddresses(1, 999999, []btcutil.Address{addPub})
	if err != nil {
		return err
	}
	sort.Slice(unspents, func(i, j int) bool {
		if unspents[i].Amount > unspents[j].Amount {
			return true
		}
		return false
	})
	payload, err := o.CreatePayloadSimpleSendCmd(amount.String())
	if err != nil {
		return err
	}
	var (
		jsonTxs []btcjson.TransactionInput
		totalAmount float64
	)
	var unspentIndex int
	for i, v := range unspents {
		jsonTxs = append(jsonTxs, btcjson.TransactionInput{Txid:v.TxID, Vout:v.Vout})
		totalAmount += v.Amount
		if totalAmount > amount.ToBTC() {
			unspentIndex = i
			break
		}
	}
	unspents = unspents[:unspentIndex]
	msgTx, err := o.CreateRawTransaction(jsonTxs, nil, nil)
	if err != nil {
		return err
	}
	var buf bytes.Buffer
	if err = msgTx.Serialize(&buf);err != nil {
		return err
	}
	//opReturn, err := client.CreateRawTxOpReturn(hex.EncodeToString(buf.Bytes()), payload)
	if err != nil {
		return err
	}
	_ = payload
	var PrevTxs []PrevTx
	for _, v := range unspents {
		PrevTxs = append(PrevTxs, PrevTx{
			TxID:v.TxID,
			Vout:v.Vout,
			ScriptPubKey:v.ScriptPubKey,
			Value:v.Amount,
		})
	}
	return nil
}

func (o *omniClient) CreatePayloadSimpleSendCmd(amount string) (string, error) {
	cmd := NewCreatePayloadSimpleSendCmd(31, amount)
	b, err := o.sendCmd(cmd)
	if err != nil {
		return "", err
	}
	var pay string
	err = json.Unmarshal(b, &pay)
	if err != nil {
		return "", err
	}
	return pay, nil
}

func (o *omniClient) CreateRawTxOpReturn(rawTx, payload string) (string, error) {
	cmd := NewCreateRawTxOpReturnCmd(rawTx, payload)
	resBytes, err := o.sendCmd(cmd)
	if err != nil {
		return "", err
	}
	var opreturntx string
	err = json.Unmarshal(resBytes, &opreturntx)
	if err != nil {
		return "", err
	}
	return opreturntx, nil
}


type jsonRequest struct {
	id          uint64
	method      string
	cmd         interface{}
	marshalJSON []byte
}

func (o *omniClient) sendCmd(cmd interface{}) ([]byte, error) {
	method, err := btcjson.CmdMethod(cmd)
	if err != nil {
		return nil, err
	}
	id := o.Client.NextID()
	data, err := btcjson.MarshalCmd(id, cmd)
	if err != nil {
		return nil, err
	}
	jReq := &jsonRequest{
		id:          id,
		method:      method,
		cmd:         cmd,
		marshalJSON: data,
	}
	return o.sendRequest(jReq)
}

func (o *omniClient) sendRequest(request *jsonRequest) ([]byte, error) {
	body := bytes.NewReader(request.marshalJSON)
	req, err := http.NewRequest(http.MethodPost, o.url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	return o.sendPostRequest(req)
}

func (o *omniClient) sendPostRequest(r *http.Request) ([]byte, error) {
	c := http.Client{}
	httpRes, err := c.Do(r)
	if err != nil {
		return nil, nil
	}
	buf, err := ioutil.ReadAll(httpRes.Body)
	if err != nil {
		return nil, err
	}
	_ = httpRes.Body.Close()
	var res Response
	err = json.Unmarshal(buf, &res)
	if err != nil {
		return nil, err
	}
	if res.Error != nil {
		return nil, res.Error
	}
	return res.Result, nil
}