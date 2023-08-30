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
	insertRowDemo()
}

//插入数据
func insertRowDemo() {
	sqlStr := "insert into user(name, age) values (?, ?)"
	ret, err := db.Exec(sqlStr, "test3", 19)
	if err != nil {
		fmt.Printf("insert failed, err=%v\n", err)
		return
	}
	var theId int64
	theId, err = ret.LastInsertId() //新插入数据的ID
	if err != nil {
		fmt.Printf("get lastInsert id failed, err=%v\n", err)
		return
	}
	fmt.Printf("insert success, id=%d\n", theId)
}
