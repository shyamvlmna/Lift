package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v9"
	"github.com/shayamvlmna/cab-booking-app/pkg/models"
)

func OpenRDb() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	return rdb
}

func Set(key, value string) error {
	rdb := OpenRDb()
	err := rdb.Set(context.Background(), key, value, time.Minute*10).Err()
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
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

func StoreData(key string, value any) error {
	rdb := OpenRDb()
	p, err := json.Marshal(value)
	if err != nil {
		return err
	}
	err = rdb.Set(context.Background(), key, p, time.Minute*20).Err()
	if err != nil {
		fmt.Println(err)
	}
	return nil
}

func GetData(key string) (string, error) {
	rdb := OpenRDb()
	p, err := rdb.Get(context.Background(), key).Result()
	if err != nil {
		fmt.Println(err)
		return "", redis.Nil
	}
	return p, nil
}

func DeleteData(key string) error {
	rdb := OpenRDb()
	r := rdb.Del(context.Background(), key)

	return r.Err()
}

func StoreTrip(key string, trip *models.Ride) error {
	rdb := OpenRDb()
	p, err := json.Marshal(trip)
	if err != nil {
		return err
	}
	err = rdb.Set(context.Background(), key, p, 0).Err()
	if err != nil {
		fmt.Println(err)
	}
	return nil
}

func GetTrip(key string) (*models.Ride, error) {
	rdb := OpenRDb()
	p, err := rdb.Get(context.Background(), key).Result()
	if err != nil {
		fmt.Println(err)
		return nil, redis.Nil
	}
	ride := &models.Ride{}
	err = json.Unmarshal([]byte(p), &ride)
	if err != nil {
		return nil, err
	}

	return ride, nil
}

// func set(c *RedisClient, key string, value interface{}) error {
//     p, err := json.Marshal(value)
//     if err != nil {
//        return err
//     }
//     return c.Set(key, p)
// }

// func get(c *RedisClient, key string, dest interface{}) error {
// 	p, err := c.Get(key)
// 	if err != nil {
// 		return err
// 	}
// 	return json.Unmarshal(p, dest)
// }
