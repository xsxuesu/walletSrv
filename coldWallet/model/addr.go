package model

import (
	"database/sql"
	"errors"
	"fmt"
	"walletSrv/coldWallet/config"
)

type AddrModel struct {
	Db *sql.DB
}

type AddrList struct {
	AddrList []string `json:"addrlist"`
}

func getTableName(cointype string,hd bool) (string,error){
	var tablename string

	if (cointype == config.Btc) {
		if hd {
			tablename = config.BTCHDAddrTable
		}else{
			tablename = config.BTCAddrTable
		}
	}else if(cointype == config.Usdt){
		if hd {
			tablename = config.USDTHDAddrTable
		}else{
			tablename = config.USDTAddrTable
		}
	}else if(cointype == config.Eth){
		if hd {
			tablename = config.ETHHDAddrTable
		}else{
			tablename = config.ETHAddrTable
		}

	}else{
		return "", errors.New("cointype error, please check coin type in (btc,usdt,eth)")
	}

	return tablename,nil
}

func getAddrTableName(cointype string)(string,string,error){
	var tablename ,hdtablename string
	if (cointype == config.Btc) {
		tablename = config.BTCAddrTable
		hdtablename = config.BTCHDAddrTable
	}else if(cointype == config.Usdt){
		tablename = config.USDTAddrTable
		hdtablename = config.USDTHDAddrTable
	}else if(cointype == config.Eth){
		tablename = config.ETHAddrTable
		hdtablename = config.ETHHDAddrTable
	}else{
		return "","",errors.New("cointype error, please check coin type in (btc,usdt,eth)")
	}
	return tablename,hdtablename,nil
}

//	查询币种的账号
func (addr AddrModel) FindRange(cointype string,start int,end int)(AddrList,bool,error){

	tablename ,hdtablename , err := getAddrTableName(cointype)
	if err != nil {
		return AddrList{},false,err
	}

	var sqlQuery = fmt.Sprintf("(select address,private from %s  )" +
		"UNION (select address,private from %s ) limit %d,%d",tablename,hdtablename,start,end)

	row , err := addr.Db.Query(sqlQuery)

	defer row.Close()

	List := []string{}

	for row.Next(){
		var addrT string
		var privateK string
		err = row.Scan(&addrT,&privateK)
		if err != nil {
			continue
		}
		List = append(List,addrT)
	}

	// 是否返回结束
	isFinish := false

	if (end - start) > len(List) {
		isFinish = true
	}

	return AddrList{
		AddrList:List,
	},isFinish,nil
}

func (addr AddrModel) FindByAddress(cointype string,address string)(AddressEntity,error){

	tablename ,hdtablename , err := getAddrTableName(cointype)
	if err != nil {
		return AddressEntity{},err
	}

	var sqlQuery = fmt.Sprintf("select address,private from %s where address = ? " +
		"UNION select address,private from %s where address = ?",tablename,hdtablename)

	fmt.Println("Query sql:", sqlQuery)

	var addrO  = AddressEntity{}

	row , err := addr.Db.Query(sqlQuery, address,address)

	if err != nil {
		return AddressEntity{},err
	}

	defer row.Close()

	for row.Next(){
		var addrT string
		var privateK string
		err = row.Scan(&addrT,&privateK)
		if err != nil {
			return AddressEntity{},err
		}

		addrO = AddressEntity{
			addrT,
			privateK,
		}
	}

	return addrO,nil
}

func (addr AddrModel) Insert(cointype string,hd bool,addrO AddressEntity)error{

	tablename,err := getTableName(cointype,hd)

	if err != nil {
		return err
	}
	/// 事务begin
	begin,err := addr.Db.Begin()
	if err != nil {
		return err
	}
	insertSql := fmt.Sprintf("insert into %s (address,private) values (?,?)",tablename)

	fmt.Println("insert sql:", insertSql)
	fmt.Printf("insert sql:Address=[%v],PrivateKey=[%v]\n", addrO.Address,addrO.PrivateKey)

	result,err := addr.Db.Exec(insertSql,
		addrO.Address,addrO.PrivateKey)

	if err != nil {
		return err
	}


	insert,err := result.LastInsertId()
	if err != nil {
		return err
	}

	err = begin.Commit()
	if err != nil {
		return err
	}

	fmt.Println("insert value:", insert)

	return nil
}

func (addr AddrModel) Update(cointype string,hd bool,addrO AddressEntity)error{
	tablename,err := getTableName(cointype,hd)
	if err != nil {
		return err
	}
	updateSql := fmt.Sprintf("UPDATE %s SET private = ? WHERE address = ?",tablename)

	fmt.Println("update sql:", updateSql)

	stmt ,err := addr.Db.Prepare(updateSql)

	if err != nil {
		return err
	}
	/// 事务begin
	begin,err := addr.Db.Begin()
	if err != nil {
		return err
	}
	result,err := begin.Stmt(stmt).Exec(addrO.PrivateKey,addrO.Address)
	if err != nil {
		return err
	}
	update,err := result.RowsAffected()
	if err != nil {
		return err
	}
	fmt.Println("Update value:", update)

	// 事务 commit
	err = begin.Commit()

	if err != nil {
		return err
	}

	return nil
}