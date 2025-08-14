package utils

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/redis/go-redis/extra/redisotel/v9"
	"github.com/redis/go-redis/v9"
)

type RedisAccess struct {
	client  *redis.Client
	context *context.Context
}

var StoreAccess *RedisAccess
var once sync.Once

func GetStoreAccess() *RedisAccess {
	once.Do(func() {
		StoreAccess = func() *RedisAccess {
			client := redis.NewClient(&redis.Options{
				Addr:     "redis:6379",
				Password: "", // no password set
				DB:       0,  // use default DB
			})
			ctx := context.Background()
			return &RedisAccess{client: client, context: &ctx}
		}()
	})
	return StoreAccess
}

func TestRedis() {
	connectRedis()
}

/*
This will eventually will access the database and fetch new details, for now, it is dummy
*/
func GetShows() map[string]string {
	shows := map[string]string{
		"show1": "Movie 1",
		"show2": "Movie 2",
		"show3": "Movie 3",
	}
	return shows
}

func PutShow(name string, details string) {
	// This function will eventually update the show details in the database
	log.Println("Updating show:", name, "with details:", details)
}

func AddToCache(key string, value string, redisAccess *RedisAccess) error {
	err := redisAccess.client.Set(*redisAccess.context, key, value, 0).Err()
	if err != nil {
		log.Println("Error setting value in Redis:", err)
	} else {
		log.Println("Value set in Redis:", key, "=", value)
	}
	return err
}

func GetFromCache(key string, redisAccess *RedisAccess) (string, error) {
	val, err := redisAccess.client.Get(*redisAccess.context, key).Result()
	if err != nil {
		log.Println("Error getting value from Redis:", err)
		return "", err
	}
	log.Println("Value retrieved from Redis:", key, "=", val)
	return val, nil
}

func DeleteFromCache(key string, redisAccess *RedisAccess) error {
	err := redisAccess.client.Del(*redisAccess.context, key).Err()
	if err != nil {
		log.Println("Error deleting value from Redis:", err)
		return err
	}
	log.Println("Value deleted from Redis:", key)
	return nil
}

func ClearCache(redisAccess *RedisAccess) error {
	err := redisAccess.client.FlushAll(*redisAccess.context).Err()
	if err != nil {
		log.Println("Error clearing cache in Redis:", err)
		return err
	}
	log.Println("Cache cleared in Redis")
	return nil
}

func HashSet(key string, fields map[string]interface{}, redisAccess *RedisAccess) (int64, error) {
	res1, err := redisAccess.client.HSet(*redisAccess.context, key, fields).Result()
	if err != nil {
		log.Println("Error setting hash in Redis:", err)
		return 0, err
	}
	log.Println("Hash set in Redis:", key)
	return res1, nil
}

func HashGet(key string, field string, redisAccess *RedisAccess) (string, error) {
	res, err := redisAccess.client.HGet(*redisAccess.context, key, field).Result()
	if err != nil {
		log.Println("Error getting hash from Redis:", err)
		return "", err
	}
	log.Println("Hash retrieved from Redis:", key, field, "=", res)
	return res, nil
}

func HashGetAll(key string, redisAccess *RedisAccess) (map[string]string, error) {
	res, err := redisAccess.client.HGetAll(*redisAccess.context, key).Result()
	if err != nil {
		log.Println("Error getting hash from Redis:", err)
		return nil, err
	}
	log.Println("Hash retrieved from Redis:", key, "=", res)
	return res, nil
}

func connectRedis() {
	redisAccess := GetStoreAccess()
	addInstrumentation(redisAccess.client)

	err := AddToCache("foo", "bar", redisAccess)
	if err != nil {
		panic(err)
	}

	val, err := GetFromCache("foo", redisAccess)
	if err != nil {
		panic(err)
	}
	fmt.Println("foo", val)
	log.Println("Connected to Redis successfully")

	hashFields := map[string]interface{}{
		"model": "Deimos",
		"brand": "Ergonom",
		"type":  "Enduro bikes",
		"price": "4972",
	}

	res1, err := HashSet("bike:1", hashFields, redisAccess)

	if err != nil {
		panic(err)
	}

	fmt.Println(res1) // >>> 4

	res2, err := HashGet("bike:1", "model", redisAccess)

	if err != nil {
		panic(err)
	}

	fmt.Println(res2) // >>> Deimos

	res3, err := HashGet("bike:1", "price", redisAccess)

	if err != nil {
		panic(err)
	}

	fmt.Println(res3) // >>> 4972

	res4, err := HashGetAll("bike:1", redisAccess)

	if err != nil {
		panic(err)
	}

	fmt.Println(res4)
	// >>> map[brand:Ergonom model:Deimos price:4972 type:Enduro bikes]
}

func addInstrumentation(client *redis.Client) {
	// Enable tracing instrumentation.
	if err := redisotel.InstrumentTracing(client); err != nil {
		panic(err)
	}

	// Enable metrics instrumentation.
	if err := redisotel.InstrumentMetrics(client); err != nil {
		panic(err)
	}
}
