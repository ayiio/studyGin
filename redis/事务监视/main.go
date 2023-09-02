package main

import (
	"fmt"
	"github.com/go-redis/redis"
	"time"
)

var rdb *redis.Client

//初始化连接
func initClient() (err error) {
	rdb = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	_, err = rdb.Ping().Result()
	return err
}

func main() {
	if err := initClient(); err != nil {
		fmt.Printf("init redis client failed, err=%v\n", err)
		return
	}
	fmt.Println("conn to redis success")
	defer rdb.Close() //程序退出后释放db资源
	watchExample()
}

//watch
//Watch一般配合TxPipeline使用，当用户使用WATCH命令监视某个键后，直到该用户执行EXEC命令的时间段里，
//如果有其他用户抢先对被监视的键进行了替换/更新/删除等操作，那么在EXEC执行时，事务将失败并返回一个错误
//可以根据这个错误重试事务或者放弃事务
func watchExample() {
	//监视watch_count的值，并在值不变的情况下将其值+1
	key := "watch_count"
	err := rdb.Watch(func(tx *redis.Tx) error {
		n, err := tx.Get(key).Int()
		if err != nil && err != redis.Nil {
			return err
		}
		_, err = tx.Pipelined(func(pipe redis.Pipeliner) error {
			pipe.Set(key, n+1, 0)
			return nil
		})
		return err
	}, key)
	if err != nil {
		fmt.Println("err=", err)
	}
}
