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
	deleteRowDemo()
}

//删除数据
func deleteRowDemo() {
	sqlStr := "delete from user where id=?"
	ret, err := db.Exec(sqlStr, 3)
	if err != nil {
		fmt.Printf("delete failed, err=%v\n", err)
		return
	}
	var theId int64
	theId, err = ret.RowsAffected()
	if err != nil {
		fmt.Printf("get the affected id failed, err=%v\n", err)
		return
	}
	fmt.Printf("the affected id=%v\n", theId)
}
