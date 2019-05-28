package model

import (
	"database/sql"
	"errors"
	"fmt"
	"walletSrv/coldWallet/config"
)

type SerialModel struct {
	Db *sql.DB
}

// 计算表列表
//func getSerialTableNames(cointype string)([]string,error){
//	if !(cointype == config.Btc ||  cointype == config.Usdt || cointype == config.Eth){
//		return []string{},errors.New("cointype error, please check coin type in (btc,usdt,eth)")
//	}
//
//	var tableNames []string
//
//	for year := config.BeginYear; year <= time.Now().Year() ; year++  {
//		curYearName := fmt.Sprint(config.PrehasSerial,cointype,year)
//		tableNames = append(tableNames,curYearName)
//	}
//
//	return tableNames,nil
//}

func getSerialCurTableName(cointype string)(string,error){
	if !(cointype == config.Btc ||  cointype == config.Usdt || cointype == config.Eth){
		return "",errors.New("cointype error, please check coin type in (btc,usdt,eth)")
	}
	//year := time.Now().Year()
	return fmt.Sprint(config.PrehasSerial,cointype),nil
}

func (serial SerialModel)CheckTable(cointype string)(error){
	if !(cointype == config.Btc ||  cointype == config.Usdt || cointype == config.Eth){
		return errors.New("cointype error, please check coin type in (btc,usdt,eth)")
	}
	//year := fmt.Sprintf("%d",time.Now().Year())
	//month := fmt.Sprintf("%d",int(time.Now().Month()))

	curSerialName := fmt.Sprint(config.PrehasSerial,cointype)

	err := config.CreateSerialTabs(curSerialName,serial.Db)
	if err != nil {
		return err
	}
	return nil
}

func querySerial(db *sql.DB,serialNo string,tableName string,serialInfo chan SerialEntity){

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

func (serial SerialModel)FindBySerialNo(serialno string,cointype string)(SerialEntity,string,error){
	//tablenames,err := getSerialTableNames(cointype)
	//if err != nil {
	//	return SerialEntity{},"",err
	//}
	//fmt.Println("tablenames:",tablenames)

	//queryChan := make(chan SerialEntity, len(tablenames))
	//for _,v := range tablenames {
	//	//fmt.Println("tablename:",v)
	//	go querySerial(serial.Db,serialno,v,queryChan)
	//}
	//
	//for i,v := range tablenames {
	//	fmt.Println("tableindex:",i," talblename:",v)
	//	queryedInfo := <- queryChan
	//	//fmt.Println("SerialNo:",queryedInfo.SerialNo," type:",queryedInfo.CoinType)
	//	if queryedInfo.SerialNo != "" {
	//		return queryedInfo,v,nil
	//	}
	//}

	tableName , err := getSerialCurTableName(cointype)
	if err != nil {
		return SerialEntity{} ,"",err
	}

	sqlQuery := fmt.Sprintf("select serial,type,f,t,value,fee,txid,status from %s where serial = ? ",tableName)

	fmt.Println("Query sql:", sqlQuery)

	row , err := serial.Db.Query(sqlQuery, serialno)

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
		return sEntity,tableName,nil
	}

	return SerialEntity{} ,"",nil
}

func (serial SerialModel)Insert(entity SerialEntity)error{

	tableName,err := getSerialCurTableName(entity.CoinType)
	if err != nil {
		return err
	}

	insertSql := fmt.Sprintf("insert into %s (serial,type,f,t,value,fee,txid,status) values (?,?,?,?,?,?,?,?)",tableName)

	fmt.Println("insert sql:", insertSql)

	result,err := serial.Db.Exec(insertSql,
		entity.SerialNo,entity.CoinType,entity.F,entity.T,entity.Value,entity.Fee,entity.TxId,entity.Status)

	if err != nil {
		return err
	}

	insert,err := result.LastInsertId()
	if err != nil {
		return err
	}
	fmt.Println("insert value:", insert)

	return nil
}

func (serial SerialModel)UpdateStatusByTxId(entity SerialEntity)error{
	tableName,err := getSerialCurTableName(entity.CoinType)
	if err != nil {
		return err
	}

	updateSql := fmt.Sprintf("UPDATE %s SET status = ? WHERE txid = ?",tableName)

	fmt.Println("update sql:", updateSql)

	stmt ,err := serial.Db.Prepare(updateSql)

	if err != nil {
		return err
	}
	/// 事务begin
	begin,err := serial.Db.Begin()
	if err != nil {
		return err
	}
	result,err := begin.Stmt(stmt).Exec(entity.Status,entity.TxId)
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


func (serial SerialModel)UpdateTxId(entity SerialEntity)error{
	tableName,err := getSerialCurTableName(entity.CoinType)
	if err != nil {
		return err
	}

	updateSql := fmt.Sprintf("UPDATE %s SET txid = ?,status = ? WHERE serial = ?",tableName)

	fmt.Println("update sql:", updateSql)

	stmt ,err := serial.Db.Prepare(updateSql)

	if err != nil {
		return err
	}
	/// 事务begin
	begin,err := serial.Db.Begin()
	if err != nil {
		return err
	}
	result,err := begin.Stmt(stmt).Exec(entity.TxId,entity.Status,entity.SerialNo)
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

func (serial SerialModel)Update(entity SerialEntity)error{

	tableName,err := getSerialCurTableName(entity.CoinType)
	if err != nil {
		return err
	}

	updateSql := fmt.Sprintf("UPDATE %s SET type = ?,f = ?,t = ?,value = ?,fee = ?,txid = ?,status = ? WHERE serial = ?",tableName)

	fmt.Println("update sql:", updateSql)

	stmt ,err := serial.Db.Prepare(updateSql)

	if err != nil {
		return err
	}
	/// 事务begin
	begin,err := serial.Db.Begin()
	if err != nil {
		return err
	}
	result,err := begin.Stmt(stmt).Exec(entity.CoinType,entity.F,entity.T,entity.Value,entity.Fee,entity.TxId,entity.Status,entity.SerialNo)
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