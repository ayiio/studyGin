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
	updateRowDemo()
}

//更新数据
func updateRowDemo() {
	sqlStr := "update user set age=? where id=?"
	ret, err := db.Exec(sqlStr, 20, 4)
	if err != nil {
		fmt.Printf("update failed, err=%v\n", err)
		return
	}
	n, err := ret.RowsAffected() //操作影响行数
	if err != nil {
		fmt.Printf("get rowsaffected failed, err=%v\n", err)
		return
	}
	fmt.Printf("update success, the affected rows is %d.\n", n)
}
