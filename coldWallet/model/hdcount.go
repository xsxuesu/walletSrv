package model

import (
	"errors"
	"database/sql"
	"fmt"
	"strconv"
	"walletSrv/coldWallet/config"
)

type HdCountModel struct {
	Db *sql.DB
}

func (hd HdCountModel) FindByType(cointype string)(uint32,error){
	var sqlQuery = fmt.Sprintf("select type,hdcount from %s where type = ? ",
		config.HDCountTable)

	fmt.Println("Hd Query sql:", sqlQuery)

	row , err := hd.Db.Query(sqlQuery,cointype)

	if err != nil {
		return uint32(0),err
	}

	defer row.Close()

	var hdObj HdCount

	for row.Next(){
		var cointype string
		var hdcount string
		err = row.Scan(&cointype,&hdcount)
		if err != nil {
			return uint32(0),err
		}
		hdObj = HdCount{
			cointype,
			hdcount,
		}
	}

	if hdObj.HdNum == "" {
		return uint32(0),errors.New("no the coin record")
	}
	// convert to int64
	hdNum , err := strconv.ParseUint(hdObj.HdNum,10,32)

	if err != nil {
		return uint32(0),err
	}

	return uint32(hdNum),nil
}


func (hd HdCountModel) InsertHdCount (cointype string,hdcount string) error  {
	if !(cointype == config.Btc ||  cointype == config.Usdt || cointype == config.Eth){
		return errors.New("cointype error, please check coin type in (btc,usdt,eth)")
	}

	insertSql := fmt.Sprintf("insert into %s (type,hdcount) values (?,?)",config.HDCountTable)
	/// db prepare stmt
	stmt ,err := hd.Db.Prepare(insertSql)
	if err != nil {
		return err
	}
	/// 事务begin
	begin,err := hd.Db.Begin()
	if err != nil {
		return err
	}

	result,err := begin.Stmt(stmt).Exec(cointype,hdcount)
	if err != nil {
		return err
	}

	insert,err := result.LastInsertId()
	if err != nil {
		return err
	}
	fmt.Println("insert value:", insert)
	// 事务 commit
	err = begin.Commit()

	if err != nil {
		return err
	}

	return nil
}


func (hd HdCountModel) UpdateHdCount (cointype string,hdcount string) error  {

	if !(cointype == config.Btc ||  cointype == config.Usdt || cointype == config.Eth){
		return errors.New("cointype error, please check coin type in (btc,usdt,eth)")
	}

	updateSql := fmt.Sprintf("UPDATE %s SET hdcount = ? WHERE type = ?",config.HDCountTable)

	result,err := hd.Db.Exec(updateSql, hdcount,cointype)

	if err != nil {
		return err
	}

	update,err := result.RowsAffected()
	if err != nil {
		return err
	}

	fmt.Println("update value:", update)

	return nil
}