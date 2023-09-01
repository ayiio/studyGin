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
	//prepareQueryDemo()
	//prepareInsertDemo()
	
	// sql:select id, name, age from user where name='test1' or 1=1 #'
	sqlInjectDemo("test1' or 1=1 #")
}

//sql预处理
func prepareQueryDemo() {
	sqlStr := "select id, name, age from user where id>?"
	stmt, err := db.Prepare(sqlStr)
	if err != nil {
		fmt.Printf("prepare failed, err=%v\n", err)
		return
	}
	defer stmt.Close()
	rows, err := stmt.Query(0)
	if err != nil {
		fmt.Printf("query failed, err=%v\n", err)
		return
	}
	defer rows.Close()
	//循环读取结果集中的数据
	for rows.Next() {
		var u user
		err := rows.Scan(&u.id, &u.name, &u.age)
		if err != nil {
			fmt.Printf("scan failed, err=%v\n", err)
			return
		}
		fmt.Printf("id:%d, name:%s, age:%d\n", u.id, u.name, u.age)
	}
}

//sql预处理插入
func prepareInsertDemo() {
	sqlStr := "insert into user(name, age) values (?, ?)"
	stmt, err := db.Prepare(sqlStr)
	if err != nil {
		fmt.Printf("prepare failed, err=%v\n", err)
		return
	}
	defer stmt.Close()
	_, err = stmt.Exec("test4", 24)
	if err != nil {
		fmt.Printf("insert failed, err=%v\n", err)
		return
	}
	_, err = stmt.Exec("test5", 25)
	if err != nil {
		fmt.Printf("insert failed, err=%v\n", err)
		return
	}
	fmt.Println("insert success")
}

//sql注入，预处理预防
func sqlInjectDemo(name string) {
	//不要自行拼接sql
	sqlStr := fmt.Sprintf("select id, name, age from user where name='%s'", name)
	fmt.Printf("sql:%s\n", sqlStr)
	var u user
	err := db.QueryRow(sqlStr).Scan(&u.id, &u.name, &u.age)
	if err != nil {
		fmt.Printf("exec failed, err=%v\n", err)
		return
	}
	fmt.Printf("user:%v\n", u)
}
