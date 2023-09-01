package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

var db *sql.DB

type user struct {
	id   int
	name string
	age  int
}

func initMysql() (err error) {
	dsn := "root:toot@tcp(127.0.0.1:3306)/demo"
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		fmt.Printf("connect mysql failed, err=%v\n", err)
		return
	}
	db.SetConnMaxLifetime(time.Second * 10)
	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)
	return
}

func main() {
	if err := initMysql(); err != nil {
		fmt.Printf("connect to db failed, err=%v\n", err)
	}
	defer db.Close()
	fmt.Println("connect to db success")
	transactionDemo()
}

//事务
func transactionDemo() {
	tx, err := db.Begin() //开启事务
	if err != nil {
		if tx != nil {
			tx.Rollback() //回滚
		}
		fmt.Printf("begin trans failed, err=%v\n", err)
		return
	}
	sqlStr1 := "update user set age = 30 where id = ?"
	ret1, err := tx.Exec(sqlStr1, 1)
	if err != nil {
		tx.Rollback() //回滚
		fmt.Printf("exec sql1 failed, err=%v\n", err)
		return
	}
	rowAffect1, err := ret1.RowsAffected()
	if err != nil {
		tx.Rollback()
		fmt.Printf("exec rowsAffected() failed, err=%v\n", err)
		return
	}

	sqlStr2 := "update user set age = 31 where id = ?"
	ret2, err := tx.Exec(sqlStr2, 2)
	if err != nil {
		tx.Rollback() //回滚
		fmt.Printf("exec sql2 failed, err=%v\n", err)
		return
	}
	rowAffect2, err := ret2.RowsAffected()
	if err != nil {
		tx.Rollback()
		fmt.Printf("exec rowsAffected() failed, err=%v\n", err)
		return
	}

	if rowAffect1 == 1 && rowAffect2 == 1 {
		err = tx.Commit() //提交事务
		fmt.Println("exec trans success")
	} else {
		tx.Rollback()
		fmt.Println("exec rollback success")
	}
}
