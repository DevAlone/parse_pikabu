package middlewares

import (
	"net/http"
	"regexp"
	"strings"

	"github.com/DevAlone/parse_pikabu/core/logger"
	. "github.com/DevAlone/parse_pikabu/helpers"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
)

type GroupId int16

const (
	GROUP_PARSER GroupId = 0
)

func RestrictToGroupMiddleware(groupId GroupId, redisClient *redis.Client, redisPrefix string) gin.HandlerFunc {
	return func(context *gin.Context) {
		sessionId := strings.TrimSpace(context.GetHeader("Session-Id"))
		if match, _ := regexp.MatchString("^[a-z0-9A-Z_]{32,128}$", sessionId); !match {
			RespondWithError(http.StatusUnauthorized, "bad session id", context)
			return
		}
		groupIds, err := GetArrayOfIntsFromRedis(redisPrefix+"session_group_"+sessionId, redisClient)

		if err == redis.Nil {
			RespondWithError(http.StatusUnauthorized, "you're not allowed to see it", context)
			return
		} else if err != nil {
			logger.Log.Error(err)
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
