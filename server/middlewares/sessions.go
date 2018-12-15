package middlewares

import (
	. "bitbucket.org/d3dev/parse_pikabu/helpers"
	"bitbucket.org/d3dev/parse_pikabu/logging"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"net/http"
	"regexp"
	"strings"
)

type GroupId int16

const (
	GROUP_PARSER GroupId = iota
)

func RestrictToGroupMiddleware(groupId GroupId, redisClient *redis.Client, redisPrefix string) gin.HandlerFunc {
	return func(context *gin.Context) {
		sessionId := strings.TrimSpace(context.GetHeader("Session-Id"))
		if match, _ := regexp.MatchString("^[a-z0-9A-Z]{32,32}$", sessionId); !match {
			RespondWithError(http.StatusUnauthorized, "bad session id", context)
			return
		}
		groupIds, err := GetArrayOfIntsFromRedis(redisPrefix+"session_group_"+sessionId, redisClient)

		if err == redis.Nil {
			RespondWithError(http.StatusUnauthorized, "you're not allowed to see it", context)
			return
		} else if err != nil {
			logging.Log.Error(err)
			RespondWithError(http.StatusInternalServerError, "unable to get groups", context)
			return
		}

		if !IsIntInArray(int(groupId), groupIds) {
			RespondWithError(http.StatusUnauthorized, "you're not allowed to see it", context)
			return
		}

		context.Next()
	}
}
