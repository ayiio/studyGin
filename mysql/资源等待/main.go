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
	dsn := "root:123456@tcp(127.0.0.1:3306)/demo"
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
	db.SetMaxOpenConns(1)  //最大连接数设置为1
	db.SetMaxIdleConns(1)
	return
}

//查询单条记录
func queryRowDemo() {
	sqlStr := "select id, name, age from user where id=?"
	var u user
	//确保QueryRow方法后调用Scan方法，否则持有的数据库资源无法释放
	row := db.QueryRow(sqlStr, 1) //未对row进行scan

	row = db.QueryRow(sqlStr, 2) //将被hold，queryRow调用次数大于设定的最大连接数

  //一般写法，QueryRow后直接跟随Scan
  //err = db.QueryRow(sqlStr, 2).Scan(&u.id, &u.name, &u.age)
  
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
