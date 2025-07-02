package main

import (
	"github.com/chmenegatti/lazylog"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func main() {
	logger := lazylog.NewLogger(
		&lazylog.ConsoleTransport{Level: lazylog.DEBUG, Formatter: &lazylog.TextFormatter{}},
	)
	logger.EnableStacktrace(lazylog.ERROR)

	r := gin.New()

	// Middleware para logar cada request
	r.Use(func(c *gin.Context) {
		start := time.Now()
		c.Next()
		latency := time.Since(start)
		logger.WithFields(map[string]any{
			"method": c.Request.Method,
			"path":   c.Request.URL.Path,
			"status": c.Writer.Status(),
			"latency": latency.String(),
		}).Info("request completed")
	})

	r.GET("/panic", func(c *gin.Context) {
		logger.Panic("panic route called")
	})

	r.GET("/user/:id", func(c *gin.Context) {
		userLogger := logger.WithFields(map[string]any{"user_id": c.Param("id")})
		userLogger.Info("user endpoint hit")
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})

	r.Run(":8080")
}
