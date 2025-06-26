package routes

import (
	"github.com/RodolfoBonis/hermes/core/health"
	"github.com/RodolfoBonis/hermes/features/auth"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func InitializeRoutes(router *gin.Engine) {

	root := router.Group("/v1")

	root.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	health.InjectRoute(root)
	auth.InjectRoutes(root)
}
