package main

import (
	"fmt"
	"github.com/go-redis/redis"
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

//连接redis哨兵模式
func initClient2() (err error) {
	rdb = redis.NewFailoverClient(&redis.FailoverOptions{
		MasterName:    "master",
		SentinelAddrs: []string{"localhost:2379", "localhost:2379"},
	})
	_, err = rdb.Ping().Result()
	return err
}

//连接redis集群
func initClient3() (err error) {
	rdb2 := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs: []string{":7000", ":7001", ":7002", ":7003", ":7004", ":7005"},
	})
	_, err = rdb2.Ping().Result()
	return err
}

func main() {
	if err := initClient(); err != nil {
		fmt.Printf("init redis client failed, err=%v\n", err)
		return
	}
	fmt.Println("conn to redis success")
	defer rdb.Close() //程序退出后释放db资源
}
