package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v9"
)

func OpenRDb() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	fmt.Print("rdb opened", rdb)
	return rdb
}

func Set(key, value string) {
	rdb := OpenRDb()
	err := rdb.Set(context.Background(), key, value, time.Minute*5).Err()
	if err != nil {
		fmt.Println(err)
	}
}

func Get(key string) (string, error) {
	rdb := OpenRDb()
	value, err := rdb.Get(context.Background(), key).Result()
	if err != nil {
		fmt.Println(err)
		return "", redis.Nil
	}
	fmt.Println(value)
	return value, nil

}
