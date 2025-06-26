package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/RodolfoBonis/hermes/core/config"
	"github.com/RodolfoBonis/hermes/core/entities"
	"github.com/RodolfoBonis/hermes/core/errors"
	"github.com/RodolfoBonis/hermes/core/logger"
	"github.com/RodolfoBonis/hermes/core/middlewares"
	"github.com/RodolfoBonis/hermes/core/services"
	"github.com/RodolfoBonis/hermes/docs"
	"github.com/RodolfoBonis/hermes/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	app := gin.New()

	err := app.SetTrustedProxies([]string{})

	if err != nil {
		appError := errors.RootError(err.Error())
		logger.Log.Error(appError.Message, appError.ToMap())
		panic(err)
	}

	config.SentryConfig()

	_middleware := middlewares.NewMonitoringMiddleware()

	app.Use(_middleware.SentryMiddleware())
	app.Use(_middleware.LogMiddleware)

	app.Use(gin.Logger())
	app.Use(gin.Recovery())
	app.Use(gin.ErrorLogger())

	routes.InitializeRoutes(app)

	runPort := fmt.Sprintf(":%s", config.EnvPort())

	err = app.Run(runPort)

	if err != nil {
		appError := errors.RootError(err.Error())
		logger.Log.Error(appError.Message, appError.ToMap())
		panic(err)
	}

}

func init() {

	config.LoadEnvVars()

	logger.InitLogger()

	services.InitializeOAuthServer()

	// Use this for open connection with DataBase
	// appError := services.OpenConnection()
	//
	//if appError != nil {
	//	logger.Log.Error(appError.Message, appError.ToMap())
	//	panic(appError)
	//}

	// Use this for Run Yours migrations
	// services.RunMigrations()

	// Use this for open connection with RabbitMQ
	// services.StartAmqpConnection()

	docs.SwaggerInfo.Title = "hermes"
	docs.SwaggerInfo.Description = "Meet Hermes, your notification messenger of the gods—on Kafka & RabbitMQ (no winged sandals required). From SendGrid emails and WhatsApp via WhatsMeow to FCM push & Twilio SMS, Hermes guarantees scalable, reliable delivery with DLQ magic for failures. 📨⚡ #Hermes #DevOps"
	versionFileName := "version.txt"
	if config.EnvironmentConfig() == entities.Environment.Production {
		versionFileName = "/version.txt"
	}

	version := "unknown"
	if content, err := os.ReadFile(versionFileName); err == nil {
		version = strings.TrimSpace(string(content))
	}
	host := "localhost"

	if config.EnvironmentConfig() == entities.Environment.Production {
		host = "hermes.rodolfodebonis.com.br"
	}

	docs.SwaggerInfo.Host = host
	docs.SwaggerInfo.BasePath = "/v1"
	docs.SwaggerInfo.Version = version
	scheme := "http"

	if config.EnvironmentConfig() == entities.Environment.Production {
		scheme = "https"
	}

	docs.SwaggerInfo.Schemes = []string{scheme}
}
