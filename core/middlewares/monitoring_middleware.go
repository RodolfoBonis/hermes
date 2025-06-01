package middlewares

import (
	"net/http"

	"github.com/RodolfoBonis/hermes/core/logger"
	sentrygin "github.com/getsentry/sentry-go/gin"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type MonitoringMiddleware struct {
}

func NewMonitoringMiddleware() *MonitoringMiddleware {
	return &MonitoringMiddleware{}
}

func (m *MonitoringMiddleware) SentryMiddleware() gin.HandlerFunc {
	return sentrygin.New(sentrygin.Options{Repanic: true})
}

func (m *MonitoringMiddleware) LogMiddleware(ctx *gin.Context) {
	var responseBody = logger.HandleResponseBody(ctx.Writer)
	var requestBody = logger.HandleRequestBody(ctx.Request)
	requestId := uuid.NewString()

	if hub := sentrygin.GetHubFromContext(ctx); hub != nil {
		hub.Scope().SetTag("requestId", requestId)
		ctx.Writer = responseBody
	}

	ctx.Next()

	logMessage := logger.FormatRequestAndResponse(ctx.Writer, ctx.Request, responseBody.Body.String(), requestId, requestBody)

	if logMessage != "" {
		if isSuccessStatusCode(ctx.Writer.Status()) {
			logger.Log.Info(logMessage)
		} else {
			logger.Log.Error(logMessage)
		}
	}
}

func isSuccessStatusCode(statusCode int) bool {
	switch statusCode {
	case http.StatusOK, http.StatusCreated, http.StatusAccepted, http.StatusNoContent:
		return true
	default:
		return false
	}
}
