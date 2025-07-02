package main

import (
	"net/http"
	"time"

	"github.com/chmenegatti/lazylog"
	"github.com/labstack/echo/v4"
)

func main() {
	logger := lazylog.NewLogger(
		&lazylog.ConsoleTransport{Level: lazylog.INFO, Formatter: &lazylog.TextFormatter{}},
	)
	e := echo.New()

	// Middleware para logar cada request
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()
			err := next(c)
			latency := time.Since(start)
			logger.WithFields(map[string]any{
				"method":  c.Request().Method,
				"path":    c.Request().URL.Path,
				"status":  c.Response().Status,
				"latency": latency.String(),
			}).Info("request completed")
			return err
		}
	})

	e.GET("/panic", func(c echo.Context) error {
		logger.Panic("panic route called")
		return nil
	})

	e.GET("/user/:id", func(c echo.Context) error {
		userLogger := logger.WithFields(map[string]any{"user_id": c.Param("id")})
		userLogger.Info("user endpoint hit")
		return c.JSON(http.StatusOK, map[string]any{"ok": true})
	})

	e.Logger.Fatal(e.Start(":8081"))
}
