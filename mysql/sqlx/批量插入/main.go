package main

import (
	"database/sql/driver"
	"errors"
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
	u1 := user{Name: "test8", Age: 29}
	u2 := user{Name: "test9", Age: 20}
	users := []interface{}{u1, u2}
	batchInsertUser(users)
}

// Value :sql.In 需要实现driver.Value接口
func (u user) Value() (driver.Value, error) {
	return []interface{}{u.Name, u.Age}, nil
}

//#1. sql.In 批量插入
func batchInsertUser(users []interface{}) error {
	//sqlx.In 拼接语句和参数，传入参数是[]interface{}
	query, args, _ := sqlx.In(
		"INSERT INTO user (name, age) VALUES (?), (?)",
		users..., //arg实现了 driver.Value, sqlx.In 会通过Value()展开
	)
	fmt.Println(query) //查看生成的queryString
	fmt.Println(args)  //产看生成的args
	_, err := db.Exec(query, args...)
	return err
}

//#2.NamedExec批量插入
func batchInsertUser2(users []*user) error {
	_, err := db.NamedExec("insert into user (name, age) values (:name, :age)", users)
	return err
}
