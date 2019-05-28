package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
)

var CONFIGPATH = "/var/wallet/hotservice"

var RpcConfig *RPCConfig

type RPCConfig struct {
	Ip       		string 	`json:"ip"`
	Port     		string  `json:"port"`
	EthHost     	string  `json:"ethhost"`
	EthPort     	string  `json:"ethport"`
	EthKeyDir     	string  `json:"ethkeystoredir"`
	EthCallBack		string  `json:"ethCallbackUrl"`
	BtcHost     	string  `json:"btchost"`
	BtcPort     	string  `json:"btcport"`
	BtcRpcUser     	string  `json:"btcrpcuser"`
	BtcRpcPwd     	string  `json:"btcrpcpwd"`
	BtcCallBack 	string  `json:"btcCallbackUrl"`
	UsdtHost     	string  `json:"usdthost"`
	UsdtPort     	string  `json:"usdtport"`
	UsdtRpcUser     string  `json:"usdtrpcuser"`
	UsdtRpcPwd     	string  `json:"usdtprcpwd"`
	UsdtCallBack 	string 	`json:"usdtCallbackUrl"`
}


func InitConf(args []string) {
	var confFile = fmt.Sprint(CONFIGPATH, "/config.json")
	for i := 1; i < len(args); i++ {
		switch args[i] {
		case "-f":
			if i == len(args) {
				log.Fatalln("invalid config file")
			}
			confFile = args[i+1]
		}
	}
	data, err := ioutil.ReadFile(confFile)
	if err != nil {
		log.Fatalln(err.Error(), "cannot find the file: config.json")
	}
	err = json.Unmarshal(data, &RpcConfig)
	if err != nil {
		log.Fatalln(err.Error(), "cannot parse the file: config.json")
	}
}