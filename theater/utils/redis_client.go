package utils

import (
	"context"
	"log"
	"sync"

	"github.com/redis/go-redis/extra/redisotel/v9"
	"github.com/redis/go-redis/v9"
)

type RedisAccess struct {
	client  *redis.Client
	context *context.Context
	url     string
}

func NewClient(url string) *RedisAccess {
	log.Println("Connecting to Redis at ", url)
	
	// Get Redis password from environment
	redisPassword := GetEnvOrDefault("REDIS_PASSWORD", "theater_redis_pass")
	
	client := redis.NewClient(&redis.Options{
		Addr:         url,
		Password:     redisPassword,
		DB:           0, // use default DB
		PoolSize:     10,
		MinIdleConns: 5,
		MaxRetries:   3,
	})
	
	// Test the connection
	ctx := context.Background()
	_, err := client.Ping(ctx).Result()
	if err != nil {
		log.Printf("Failed to connect to Redis: %v", err)
		// Don't panic in production, but log the error
		log.Printf("Redis connection will be retried on next access")
	}
	
	addInstrumentation(client)
	return &RedisAccess{client: client, context: &ctx, url: url}
}

var StoreAccess *RedisAccess
var once sync.Once

func GetStoreAccess() *RedisAccess {
	once.Do(func() {
		redisURL := GetEnvOrDefault("REDIS_URL", "localhost:6379")
		StoreAccess = NewClient(redisURL)
	})
	log.Println("Getting Redis on ", StoreAccess.url)
	return StoreAccess
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
	res, err := redisAccess.client.HSet(*redisAccess.context, key, fields).Result()
	if err != nil {
		log.Println("Error setting hash in Redis:", err)
		return 0, err
	}
	log.Println("Hash set in Redis:", key)
	return res, nil
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

func GetAll(redisAccess *RedisAccess) (map[string]string, error) {
	keys, err := redisAccess.client.Keys(*redisAccess.context, "*").Result()
	if err != nil {
		log.Println("Error getting keys from Redis:", err)
		return nil, err
	}

	allData := make(map[string]string)
	for _, key := range keys {
		log.Println("Getting value for key", key)
		value, err := redisAccess.client.Get(*redisAccess.context, key).Result()
		if err != nil {
			log.Println("Error getting value for key", key, "from Redis:", err)
			continue
		}

		allData[key] = value
	}
	log.Println("All data retrieved from Redis")
	return allData, nil
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
