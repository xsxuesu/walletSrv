package model

import (
	"database/sql"
	"errors"
	"fmt"
	"walletSrv/coldWallet/config"
)

type DecryptModel struct {
	Db *sql.DB
}

// 计算表列表
//func getDecryptTableNames(cointype string)([]string,error){
//	if !(cointype == config.Btc ||  cointype == config.Usdt || cointype == config.Eth){
//		return []string{},errors.New("cointype error, please check coin type in (btc,usdt,eth)")
//	}
//
//	var tableNames []string
//
//	for year := config.BeginYear; year <= time.Now().Year() ; year++  {
//		curYearName := fmt.Sprint(config.PrehasDecrypt,cointype,year)
//		tableNames = append(tableNames,curYearName)
//	}
//
//	return tableNames,nil
//}

func getDecryptCurTableName(cointype string)(string,error){
	if !(cointype == config.Btc ||  cointype == config.Usdt || cointype == config.Eth){
		return "",errors.New("cointype error, please check coin type in (btc,usdt,eth)")
	}
	//year := time.Now().Year()
	return fmt.Sprint(config.PrehasDecrypt,cointype),nil
}

func (decrypt DecryptModel)CheckTable(cointype string)(error){
	if !(cointype == config.Btc ||  cointype == config.Usdt || cointype == config.Eth){
		return errors.New("cointype error, please check coin type in (btc,usdt,eth)")
	}
	//year := fmt.Sprintf("%d",time.Now().Year())
	//month := fmt.Sprintf("%d",int(time.Now().Month()))

	curDecryptName := fmt.Sprint(config.PrehasDecrypt,cointype)

	err := config.CreateDecryptTabs(curDecryptName,decrypt.Db)

	if err != nil {
		return err
	}
	return nil
}

func queryDecrpyt(db *sql.DB,serialNo string,tableName string,serialInfo chan SerialEntity){

	sqlQuery := fmt.Sprintf("select serial,type,f,t,value,fee,txid,status from %s where serial = ? ",tableName)

	fmt.Println("Query sql:", sqlQuery)

	row , err := db.Query(sqlQuery, serialNo)

	if err != nil {
		fmt.Println(err.Error())
	}

	defer row.Close()

	var sEntity SerialEntity

	for row.Next(){
		var serialQ,typeQ,fQ,tQ,valueQ,txidQ,statusQ string
		var feeQ uint64

		err = row.Scan(&serialQ,&typeQ,&fQ,&tQ,&valueQ,&feeQ,&txidQ,&statusQ)
		if err != nil {
			fmt.Println(err.Error())
		}
		sEntity = SerialEntity{
			SerialNo:serialQ,
			CoinType:typeQ,
			F:fQ,
			T:tQ,
			Value:valueQ,
			Fee:feeQ,
			TxId:txidQ,
			Status:statusQ,
		}
		fmt.Println(sEntity.SerialNo,sEntity.CoinType,sEntity.F,sEntity.T)

	}

	serialInfo <- sEntity
}

func (decrypt DecryptModel)FindByDecryptNo(serialno string,cointype string)(DecryptEntity,string,error){
	tableName,err := getDecryptCurTableName(cointype)
	if err != nil {
		return DecryptEntity{} ,"",err
	}
	//serial varchar(100) NOT NULL ,
	//type varchar(10) NOT NULL ,
	//	f varchar(100) NOT NULL ,
	//	t varchar(100) NOT NULL ,
	//	hash varchar(10) NOT NULL,
	//	time

	sqlQuery := fmt.Sprintf("select serial,type,f,t,hash from %s where serial = ? ",tableName)

	fmt.Println("Query sql:", sqlQuery)

	row , err := decrypt.Db.Query(sqlQuery, serialno)

	if err != nil {
		fmt.Println(err.Error())
	}

	defer row.Close()

	var sEntity DecryptEntity

	for row.Next(){
		var serialQ,typeQ,fQ,tQ,hashQ string
		err = row.Scan(&serialQ,&typeQ,&fQ,&tQ,&hashQ)
		if err != nil {
			fmt.Println(err.Error())
		}
		sEntity = DecryptEntity{
			SerialNo:serialQ,
			CoinType:typeQ,
			F:fQ,
			T:tQ,
			HashFun:hashQ,
		}
		return sEntity ,tableName,nil
	}

	return DecryptEntity{} ,"",nil
}

func (decrypt DecryptModel)Insert(entity DecryptEntity)error{

	tableName,err := getDecryptCurTableName(entity.CoinType)
	if err != nil {
		return err
	}


	//serial varchar(100) NOT NULL ,
	//type varchar(10) NOT NULL ,
	//	f varchar(100) NOT NULL ,
	//	t varchar(100) NOT NULL ,
	//	hash varchar(10) NOT NULL,
	//	time

	insertSql := fmt.Sprintf("insert into %s (serial,type,f,t,hash) values (?,?,?,?,?)",tableName)

	fmt.Println("insert sql:", insertSql)

	result,err := decrypt.Db.Exec(insertSql,
		entity.SerialNo,entity.CoinType,entity.F,entity.T,entity.HashFun)
	//TODO  err => nil
	if err != nil {
		return nil
	}

	insert,err := result.LastInsertId()
	//TODO  err => nil
	if err != nil {
		return nil
	}
	fmt.Println("insert value:", insert)

	return nil
}

func (decrypt DecryptModel)Update(entity DecryptEntity)error{

	tableName,err := getDecryptCurTableName(entity.CoinType)
	if err != nil {
		return err
	}

	updateSql := fmt.Sprintf("UPDATE %s SET type = ?,f = ?,t = ?,hash = ? WHERE serial = ?",tableName)

	fmt.Println("update sql:", updateSql)

	stmt ,err := decrypt.Db.Prepare(updateSql)

	if err != nil {
		return err
	}
	/// 事务begin
	begin,err := decrypt.Db.Begin()
	if err != nil {
		return err
	}
	result,err := begin.Stmt(stmt).Exec(entity.CoinType,entity.F,entity.T,entity.HashFun,entity.SerialNo)
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