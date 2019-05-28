package model

type AddrList struct {
	AddrList []string `json:"addrlist"`
}

type GetBalanceParam struct {
	CoinType string `json:"cointype"`
	Addr string 	`json:"addr"`
	Contract string `json:"contract,omitempty"`
}

type GetBalanceResp struct {
	CoinType string `json:"cointype"`
	Addr string 	`json:"addr"`
	Contract string `json:"contract,omitempty"`
	Balance string  `json:balance`
}

type InsertAddrParam struct {
	CoinType string `json:"cointype"`
	Addr string 	`json:"addr"`
	Prik string 	`json:"prik"`
}

type InsertAddrResp struct {
	CoinType string `json:"cointype"`
	Addr string 	`json:"addr"`
	Success bool 	`json:"success"`
}

type FetchAddrParam struct {
	CoinType string `json:"cointype"`
}

type FetchAddrResp struct {
	CoinType string `json:"cointype"`
	Addr string 	`json:"addr"`
	Success bool 	`json:"success"`
}

type TransferParam struct {
	CoinType string `json:"cointype"`
	Serial string 	`json:"serial"`
	From string 	`json:"from,omitempty"`
	To string 		`json:"to"`
	Value float64 	`json:"value"`
	Contract string `json:"contract,omitempty"`
}

type TransferResp struct {
	CoinType string `json:"cointype"`
	Serial string 	`json:"serial"`
	Txid string 	`json:"txid"`
	Status string 	`json:"status"`
	Success bool 	`json:"success"`
}

type CollectParam struct {
	CoinType string `json:"cointype"`
	Contract string `json:"contract,omitempty"`
	MinCount int64 	`json:"mincount,omitempty"`
}

type CollectResp struct {
	CoinType 	string `json:"cointype"`
	Collection  string `json:"collection"`
}

type Cold2HotParam struct {
	CoinType string 	`json:"cointype"`
	Value 	 float64 	`json:"value"`
	Contract string 	`json:"contract,omitempty"`
}

type HotTransferParam struct {
	CoinType string  `json:"cointype"`
	Serial 	 string  `json:"serial"`
	From 	 string  `json:"from"`
	To 	 	 string  `json:"to"`
	Contract string  `json:"contract,omitempty"`
	FeeAddr  string  `json:"feeaddr,omitempty"`
	Value 	 float64 `json:"value"`
}

type CheckParam struct {
	Contract string `json:"contract"`
	Addr string 	`json:"addr,omitempty"`
	MinCount int64 	`json:"mincount,omitempty"`
}


type CheckResult struct {
	Contract string `json:"contract"`
	SentAddr map[string]string `json:"sent_addr"`
}
