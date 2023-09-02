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
	transactionExample1()
	transactionExample2()
}

//事务 TxPipeline
//redis单线程，单个命令总是原子性的。但来自不同客户端的两个给定命令可以依次执行，也可以和其他命令交替执行，
//通过Multi/exec将这两条命令包装成原子性事务
func transactionExample1() {
	pipe := rdb.TxPipeline()

	incr := rdb.Incr("tx_pipeline_counter")
	pipe.Expire("tx_pipeline_counter", time.Hour)

	_, err := pipe.Exec()
	fmt.Println(incr.Val(), err)
	//内部执行效果：
	/*
		MULTI
		INCR tx_pipeline_counter
		EXPIRE tx_pipeline_counter 3600
		EXEC
	*/
}

func transactionExample2() {
	var incr *redis.IntCmd
	_, err := rdb.TxPipelined(func(pipe redis.Pipeliner) error {
		incr = pipe.Incr("tx_pipeline_counter")
		pipe.Expire("tx_pipeline_counter", time.Hour)
		return nil
	})
	fmt.Println(incr.Val(), err)
}
