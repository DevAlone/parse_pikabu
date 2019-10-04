package helpers

import (
	"strconv"
	"strings"
	"time"

	"github.com/go-redis/redis"
)

func GetArrayOfIntsFromRedis(key string, redisClient *redis.Client) ([]int, error) {
	value, err := redisClient.Get(key).Result()
	if err != nil {
		return nil, err
	}
	stringIds := strings.Split(value, ",")

	var result []int
	for _, stringId := range stringIds {
		id, err := strconv.ParseInt(stringId, 10, 16)
		if err != nil {
			return nil, err
		}
		result = append(result, int(id))
	}

	return result, nil
}

var redisClient *redis.Client

func GetRedisClient() *redis.Client {
	if redisClient == nil {
		redisClient = redis.NewClient(&redis.Options{
			Addr:        "localhost:6379",
			Password:    "",
			DB:          1,
			MaxRetries:  5,
			IdleTimeout: 5 * time.Minute,
		})
	}
	return redisClient
}
