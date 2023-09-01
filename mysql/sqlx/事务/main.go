package main

import (
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
	transactionDemo2()
}

//事务
func transactionDemo2() (err error) {
	tx, err := db.Beginx()
	if err != nil {
		fmt.Printf("begin trans failed, err=%v\n", err)
		return
	}
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p) //re-throw panic after rollback
		} else if err != nil {
			fmt.Println("rollback")
			tx.Rollback() //err is not nil, rollback
		} else {
			err = tx.Commit() //err is nil, if commit return error update err
			fmt.Println("commit")
		}
	}()

	sqlStr := "UPDATE user SET age=22 where id=?"
	rs1, err := tx.Exec(sqlStr, 1)
	if err != nil {
		return err
	}
	n1, err := rs1.RowsAffected()
	if err != nil {
		return err
	}
	//if n1 != 1 {
	//	fmt.Println("update#1 failed")
	//	return errors.New("update #1 failed")
	//}

	sqlStr2 := "UPDATE user set age=23 where id =?"
	rs2, err := tx.Exec(sqlStr2, 3)
	if err != nil {
		return err
	}
	n2, err := rs2.RowsAffected()
	if err != nil {
		return err
	}

	//if n2 != 1 {
	//	fmt.Println("update#2 failed")
	//	return errors.New("update #2 failed")
	//}
	if n1 != 1 || n2 != 1 {
		if n1 != 1 {
			fmt.Println("update#1 failed")
		} else if n2 != 1 {
			fmt.Println("update#2 failed")
		}
		return errors.New("update failed")
	}
	return
}
