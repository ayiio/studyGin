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

func main() {
	if err := initClient(); err != nil {
		fmt.Printf("init redis client failed, err=%v\n", err)
		return
	}
	fmt.Println("conn to redis success")
	defer rdb.Close() //程序退出后释放db资源
	redisExample1()
	redisExample2()
	redisExample3()
}

//set/get示例
func redisExample1() {
	err := rdb.Set("score", 100, 0).Err()
	if err != nil {
		fmt.Printf("set score failed, err=%v\n", err)
		return
	}

	val, err := rdb.Get("score").Result()
	if err != nil {
		fmt.Printf("get score failed, err=%v\n", err)
		return
	}
	fmt.Println("score: ", val)

	val2, err := rdb.Get("name").Result()
	if err == redis.Nil {
		fmt.Println("name does not exist")
	} else if err != nil {
		fmt.Printf("get name failed, err=%v\n", err)
		return
	} else {
		fmt.Println("name: ", val2)
	}
	fmt.Println("redisExample1 ==========")
}

//hash get/set示例
func redisExample2() {
	rdb.HSet("user", "name", "test1").Val()
	rdb.HSet("user", "age", 20)

	// HGetAll : (map[string]string, error)
	mres, err := rdb.HGetAll("user").Result()
	if err == redis.Nil {
		fmt.Println("user does not exist")
	} else if err != nil {
		fmt.Printf("get user failed, err=%v\n", err)
	} else {
		fmt.Println("user: ", mres)
	}

	// HMGet: ([]interface{}, error)
	ires, _ := rdb.HMGet("user", "name", "age").Result()
	fmt.Println(ires)

	// HGet: (string, error)
	sres, _ := rdb.HGet("user", "name").Result()
	fmt.Println(sres)
	fmt.Println("redisExample2 ==========")
}

//zset/get示例
func redisExample3() {
	fmt.Println("zadd ......")
	zsetKey := "language_rank"
	languages := []redis.Z{
		redis.Z{Score: 90, Member: "test1"},
		redis.Z{Score: 91, Member: "test2"},
		redis.Z{Score: 92, Member: "test3"},
		redis.Z{Score: 93, Member: "test4"},
	}
	//ZADD
	num, err := rdb.ZAdd(zsetKey, languages...).Result()
	if err != nil {
		fmt.Printf("zadd failed, err=%v\n", err)
		return
	}
	fmt.Printf("zadd %d success.\n", num)

	fmt.Println("zincrby ......")
	//给test1增加10分
	newScore, err := rdb.ZIncrBy(zsetKey, 10, "test1").Result()
	if err != nil {
		fmt.Printf("zincrby failed, err=%v\n", err)
		return
	}
	fmt.Printf("test1's score is %f now.\n", newScore)

	fmt.Println("zrevrangewith ......")
	//去除分数最高的前三个
	ret, err := rdb.ZRevRangeWithScores(zsetKey, 0, 2).Result()
	if err != nil {
		fmt.Printf("zrevrange failed, err=%v\n", err)
		return
	}
	for _, zr := range ret {
		fmt.Println(zr.Member, zr.Score)
	}

	fmt.Println("zrangeby ......")
	//取分数介于90到92之间的元素
	op := redis.ZRangeBy{
		Min: "90",
		Max: "92",
	}
	ret, err = rdb.ZRangeByScoreWithScores(zsetKey, op).Result()
	if err != nil {
		fmt.Printf("zrangebyscore failed, err=%v\n", err)
		return
	}
	for _, zrb := range ret {
		fmt.Println(zrb.Member, zrb.Score)
	}
	fmt.Println("redisExample3 ==========")
}
