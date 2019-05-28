package config

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"io/ioutil"
	"log"
	"strings"
)

var CONFIGPATH = "/var/wallet/coldservice"

var __config *DbConfig

type DbConfig struct {
	Ip       	string 	`json:"ip"`
	Port     	string    `json:"port"`
	DBUser 	 	string 	`json:"username"`
	DBPassword 	string `json:"password"`
	DBName 		string `json:"dbname"`
}

const (
	Pending = "pending"
	Executed = "executed"
	Failed = "failed"
)

const(
	Btc = "btc"
	Usdt = "usdt"
	Eth = "eth"
)

const(
	BTCAddrTable = "btc_addr"
	BTCHDAddrTable = "btc_hdaddr"
	USDTAddrTable = "usdt_addr"
	USDTHDAddrTable = "usdt_hdaddr"
	ETHAddrTable = "eth_addr"
	ETHHDAddrTable = "eth_hdaddr"
	HDCountTable = "hd_count"
)

const(
	PrehasSerial = "serial_"
	PrehasDecrypt = "decrypt_"

	CreateHDCount = `CREATE TABLE IF NOT EXISTS %s(
	type varchar(50) NOT NULL ,
	hdcount varchar(100) NOT NULL ,
	PRIMARY KEY (type)
	);`

	CreateAddr = `CREATE TABLE IF NOT EXISTS %s (
	address varchar(100) NOT NULL ,
	private varchar(200) NOT NULL ,
	PRIMARY KEY  (address)
	);`
	CreateHDAddr = `CREATE TABLE IF NOT EXISTS %s(
	address varchar(100) NOT NULL ,
	private varchar(200) NOT NULL ,
	PRIMARY KEY  (address)
	);`

	SerialTable = `CREATE TABLE IF NOT EXISTS %s (
	serial varchar(100) NOT NULL ,
	type varchar(10) NOT NULL ,
	f varchar(100) NOT NULL ,
	t varchar(100) NOT NULL ,
	value varchar(300) NOT NULL ,
	fee double NOT NULL ,
	time timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
	txid varchar(200) NOT NULL ,
	status varchar(20) NOT NULL ,
	PRIMARY KEY  (serial)
	);`

	DecryptTable = `CREATE TABLE IF NOT EXISTS %s (
	serial varchar(100) NOT NULL ,
	type varchar(10) NOT NULL ,
	f varchar(100) NOT NULL ,
	t varchar(100) NOT NULL ,
	hash varchar(10) NOT NULL,
	time timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP ,
	PRIMARY KEY  (serial)
	);`
)



func InitConf(args []string) {
	fmt.Println("InitConf:")
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
	err = json.Unmarshal(data, &__config)
	if err != nil {
		log.Fatalln(err.Error(), "cannot parse the file: config.json")
	}

	fmt.Println("Ip:")
	fmt.Println(__config.Ip)
}


//create database if not exists wallet character set utf8;
// root:@tcp(127.0.0.1:3306)/test?parseTime=true
func GetMySqlDb()(db *sql.DB,err error){
	fmt.Println("GetMySqlDb=====")
	fmt.Println(__config.Ip)

	dbDriver := "mysql"
	dbUser := __config.DBUser
	dbPassword := __config.DBPassword
	dbName := __config.DBName
	dbIPport := fmt.Sprintf("@tcp(%s)/",__config.Ip+":"+__config.Port)
	db,err = sql.Open(dbDriver,dbUser+":"+dbPassword+dbIPport+dbName)
	return
}


func CreateHDCountTabs(db *sql.DB) error {

	hdCountSql := fmt.Sprintf(CreateHDCount,HDCountTable)

	_, err := db.Exec(hdCountSql)

	if err != nil{
		return errors.New(fmt.Sprint("create hdcount table:",err.Error()))
	}

	return nil
}


func CreateAddrTabs(db *sql.DB) error {

	btcAddSql := fmt.Sprintf(CreateAddr,"btc_addr")
	usdtAddSql := fmt.Sprintf(CreateAddr,"usdt_addr")
	ethAddSql := fmt.Sprintf(CreateAddr,"eth_addr")

	btcHdAddSql := fmt.Sprintf(CreateAddr,"btc_hdaddr")
	usdtHdAddSql := fmt.Sprintf(CreateAddr,"usdt_hdaddr")
	ethHdAddSql := fmt.Sprintf(CreateAddr,"eth_hdaddr")

	_, err := db.Exec(btcAddSql)
	if err != nil{
		 return errors.New(fmt.Sprint("create bitcoin address table:",err.Error()))
	}
	_, err = db.Exec(usdtAddSql)
	if err != nil{
		return errors.New(fmt.Sprint("create usdt address table:",err.Error()))
	}
	_, err = db.Exec(ethAddSql)
	if err != nil{
		return errors.New(fmt.Sprint("create eth address table:",err.Error()))
	}

	_, err = db.Exec(btcHdAddSql)
	if err != nil{
		return errors.New(fmt.Sprint("create bitcoin hd address table:",err.Error()))
	}


	_, err = db.Exec(usdtHdAddSql)
	if err != nil{
		return errors.New(fmt.Sprint("create usdt hd address table:",err.Error()))
	}

	_, err = db.Exec(ethHdAddSql)
	if err != nil{
		return errors.New(fmt.Sprint("create eth hd address table:",err.Error()))
	}
	return nil
}


func CreateSerialTabs(tableName string, db *sql.DB) error{

	if !strings.HasPrefix(tableName,PrehasSerial){
		return errors.New("serial table name error, must has prefix serial")
	}

	serialTable := fmt.Sprintf(SerialTable,tableName)
	fmt.Println(serialTable)
	_, err := db.Exec(serialTable)
	if err != nil{
		return err
	}
	return nil
}


func CreateDecryptTabs(tableName string, db *sql.DB) error{
	if !strings.HasPrefix(tableName,PrehasDecrypt){
		return errors.New("decrypt table name error, must has prefix decrypt")
	}
	decryptTable := fmt.Sprintf(DecryptTable,tableName)
	fmt.Println(decryptTable)
	_, err := db.Exec(decryptTable)
	if err != nil{
		return err
	}
	return nil
}