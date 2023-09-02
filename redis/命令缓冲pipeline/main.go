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
	pipelineExample1()
	pipelineExample2()
}

//pipeline: 网络优化，客户端缓冲一堆命令并一次性发送给服务器，不能保证在事务中执行，节省了每个命令的网络往返时间(RTT)
//多条命令要执行时，可以考虑pipeline，但不适合相互依赖的命令使用pipeline，例如依赖前一条的结果再决定后一条命令
func pipelineExample1() {
	pipe := rdb.Pipeline()
	incr := pipe.Incr("pipeline_counter")
	pipe.Expire("pipeline_counter", time.Hour)
	
	//实际执行
	_, err := pipe.Exec()
	fmt.Println(incr.Val(), err)
}

//pipeline的包装--> pipelined
func pipelineExample2() {
	var incr *redis.IntCmd
	_, err := rdb.Pipelined(func(pipe redis.Pipeliner) error {
		incr = pipe.Incr("pipeline_counter")
		pipe.Expire("pipeline_counter", time.Hour)
		return nil
	})
	fmt.Println(incr.Val(), err)
}
