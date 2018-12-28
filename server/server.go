package server

import (
	"bitbucket.org/d3dev/parse_pikabu/config"
	"bitbucket.org/d3dev/parse_pikabu/helpers"
	"bitbucket.org/d3dev/parse_pikabu/server/api"
	"bitbucket.org/d3dev/parse_pikabu/server/middlewares"
	"github.com/gin-gonic/gin"
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
	{
		apiRouter.POST("get_model", api.GetModel)
	}

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
		helpers.GetRedisClient(),
		"parse_pikabu_server_authentication_middleware_",
	))
	{
		parsersAPI.GET("get/parsers_status", func(context *gin.Context) {
			context.JSON(http.StatusOK, map[string]interface{}{
				"status": "ok",
			})
		})
		parsersAPI.GET("get/tasks/any", api.GetAnyTask)
		// parsersAPI.GET("take/:table_name/:id", api.TakeTask)
	}

	return router.Run(config.Settings.ServerListeningAddress)
}
