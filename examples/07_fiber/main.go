package main

import (
	"github.com/chmenegatti/lazylog"
	"github.com/gofiber/fiber/v2"
	"time"
)

func main() {
	logger := lazylog.NewLogger(
		&lazylog.ConsoleTransport{Level: lazylog.INFO, Formatter: &lazylog.TextFormatter{}},
	)
	app := fiber.New()

	// Middleware para logar cada request
	app.Use(func(c *fiber.Ctx) error {
		start := time.Now()
		err := c.Next()
		latency := time.Since(start)
		logger.WithFields(map[string]any{
			"method": c.Method(),
			"path":   c.Path(),
			"status": c.Response().StatusCode(),
			"latency": latency.String(),
		}).Info("request completed")
		return err
	})

	app.Get("/panic", func(c *fiber.Ctx) error {
		logger.Panic("panic route called")
		return nil
	})

	app.Get("/user/:id", func(c *fiber.Ctx) error {
		userLogger := logger.WithFields(map[string]any{"user_id": c.Params("id")})
		userLogger.Info("user endpoint hit")
		return c.JSON(fiber.Map{"ok": true})
	})

	app.Listen(":8082")
}
