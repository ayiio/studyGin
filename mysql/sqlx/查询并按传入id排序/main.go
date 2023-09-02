package main

import (
	"database/sql/driver"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"strings"
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
	users, err := queryByIds([]int{5, 1, 2, 3})
	if err != nil {
		fmt.Printf("queryByIds failed, err=%v\n", err)
		return
	}
	for _, user := range users {
		fmt.Printf("user:%v\n", user)
	}

	fmt.Println("==========")
	users, err = queryAndOrderByIDs([]int{5, 1, 2, 3})
	if err != nil {
		fmt.Printf("queryByIds failed, err=%v\n", err)
		return
	}
	for _, user := range users {
		fmt.Printf("user:%v\n", user)
	}
}

//select in示例，queryByIds，根据指定的id进行查询
func queryByIds(ids []int) (users []user, err error) {
	query, args, err := sqlx.In("select id, name, age from user where id in (?)", ids)
	if err != nil {
		return
	}
	//sqlx.In返回带? bindVar的查询语句，使用Rebind重新绑定
	query = db.Rebind(query)
	fmt.Println(query)
	err = db.Select(&users, query, args...)
	return
}

//查询给定id集合的数据并保持传入id的顺序进行输出
func queryAndOrderByIDs(ids []int) (users []user, err error) {
	//动态填充id
	strIDs := make([]string, 0, len(ids))
	for _, id := range ids {
		strIDs = append(strIDs, fmt.Sprintf("%d", id))
	}
	query, args, err := sqlx.In("select id, name, age from user where id in (?) ORDER BY FIND_IN_SET(id, ?)",
		ids, strings.Join(strIDs, ","))
	fmt.Println(query)
	if err != nil {
		return
	}

	//sqlx.In 返回带?的 bindVal查询语句，使用Rebind()重新绑定
	query = db.Rebind(query)

	err = db.Select(&users, query, args...)
	return
}

// Value :sql.In 需要实现driver.Value接口
func (u user) Value() (driver.Value, error) {
	return []interface{}{u.Name, u.Age}, nil
}
