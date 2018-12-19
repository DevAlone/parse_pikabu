package server

import (
	"bitbucket.org/d3dev/parse_pikabu/config"
	"bitbucket.org/d3dev/parse_pikabu/server/api"
	"bitbucket.org/d3dev/parse_pikabu/server/middlewares"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"net/http"
)

func Run() error {
	if config.Settings.Debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// define controllers

	apiRouter := router.Group("/api/v1/")
	// apiRouter.Use(middlewares.AuthMiddleware.MiddlewareFunc())
	// apiRouter.Use()
	{
		apiRouter.GET("get/status", func(context *gin.Context) {
			context.JSON(http.StatusOK, map[string]interface{}{
				"status": "ok",
			})
		})
	}

	parsersAPI := apiRouter.Group("")
	parsersAPI.Use(middlewares.RestrictToGroupMiddleware(
		middlewares.GROUP_PARSER,
		redis.NewClient(&redis.Options{
			Addr: ":6379",
		}),
		"parse_pikabu_authentication_middleware_",
	))
	{
		parsersAPI.GET("get/parsers_status", func(context *gin.Context) {
			context.JSON(http.StatusOK, map[string]interface{}{
				"status": "ok",
			})
		})
		parsersAPI.GET("get/tasks/any", api.GetAnyTask)
		parsersAPI.GET("take/parse_user_by_username_tasks/:username", api.TakeParseUserByUsernameTask)
		parsersAPI.GET("take/parse_user_by_id_tasks/:id", api.TakeParseUserByIdTask)
	}

	return router.Run(config.Settings.ServerListeningAddress)
}
