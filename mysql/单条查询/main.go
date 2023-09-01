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
	db.SetMaxOpenConns(200)
	db.SetMaxIdleConns(10)
	return
}

//查询单条记录
func queryRowDemo() {
	sqlStr := "select id, name, age from user where id=?"
	var u user
	//确保QueryRow方法后调用Scan方法，否则持有的数据库资源无法释放
	row := db.QueryRow(sqlStr, 1)
	err := row.Scan(&u.id, &u.name, &u.age)
	if err != nil {
		fmt.Printf("scan failed, err=%v\n", err)
		return
	}
	fmt.Printf("id:%d, name:%s, age:%d\n", u.id, u.name, u.age)
}

func main() {
	if err := initMysql(); err != nil {
		fmt.Printf("connect to db failed, err=%v\n", err)
	}
	defer db.Close()
	fmt.Println("connect to db success")
	queryRowDemo()
}
