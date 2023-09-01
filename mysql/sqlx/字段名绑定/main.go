package main

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var db *sqlx.DB

type user struct {
	ID   int    `db:"id"` //首字母大写，第三方包反射设置值，tag指定db中的字段名
	Name string `db:"name"`
	Age  int    `db:"age"`
}

func initDB() (err error) {
	dsn := "root:toot@tcp(127.0.0.1:3306)/demo?charset=utf8mb4&parseTime=True"
	//sqlx.MustConnect() MustConnect connects to a database and panics on error.
	db, err = sqlx.Connect("mysql", dsn)
	if err != nil {
		fmt.Printf("connect to db failed, err=%v\n", err)
		return
	}
	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(10)
	return
}

func main() {
	if err := initDB(); err != nil {
		fmt.Printf("init db failed, err=%v\n", err)
		return
	}
	fmt.Println("init db success")
	queryUserDemo()
}

//绑定sql语句和结构体或map中的同名字段 NamedExec
//避免使用占位符引起传入字段顺序错误的问题
func insertUserDemo() (err error) {
	_, err = db.NamedExec(`INSERT INTO user (name, age) VALUES (:name, :age)`,
		map[string]interface{}{
			"name": "test7",
			"age":  29,
		})
	return
}

//绑定字段查询 NamedQuery
func queryUserDemo() (err error) {
	sqlStr := "SELECT * FROM user WHERE name=:name"
	//方法1
	rows, err := db.NamedQuery(sqlStr, map[string]interface{}{"name": "test1"}) //map映射
	if err != nil {
		fmt.Printf("namedquery failed, err=%v\n", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var u user
		//rows.SliceScan() ([]interface{}, error) 返回未知时
		if err = rows.StructScan(&u); err != nil {
			fmt.Printf("#1 scan failed, err=%v\n", err)
			continue
		}
		fmt.Printf("user:%v\n", u)
	}

	//方法2
	u := user{
		Name: "test1",
	}
	rows, err = db.NamedQuery(sqlStr, u) //结构体映射
	if err != nil {
		fmt.Printf("namedquery failed, err=%v\n", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var u user
		if err = rows.StructScan(&u); err != nil {
			fmt.Printf("#2 scan failed, err=%v\n", err)
			continue
		}
		fmt.Printf("user:%v\n", u)
	}
	return
}
