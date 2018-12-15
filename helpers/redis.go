package helpers

import (
	"github.com/go-redis/redis"
	"strconv"
	"strings"
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
