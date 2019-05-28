package model


type AddressEntity struct {
	Address string
	PrivateKey string
}

type SerialEntity struct {
	SerialNo string
	CoinType string
	F string
	T string
	Value string
	Fee uint64
	Time string
	TxId string
	Status string
}

type DecryptEntity struct {
	SerialNo string
	CoinType string
	F string
	T string
	HashFun string
	Time string
}

type HdCount struct {
	Cointype string
	HdNum string
}