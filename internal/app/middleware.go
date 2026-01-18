package app

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/hashicorp/go-hclog"
	"time"
)

func RegisterMiddlewares(app *fiber.App, logger hclog.Logger) {
	app.Use(func(c *fiber.Ctx) error {
		id := c.Get("X-Request-ID")
		if id == "" {
			id = uuid.NewString()
		}
		c.Set("X-Request-ID", id)
		c.Locals("request_id", id)
		return c.Next()
	})

	app.Use(func(c *fiber.Ctx) error {
		start := time.Now()
		err := c.Next()

		logger.Info("app request",
			"request_id", c.Locals("request_id"),
			"method", c.Method(),
			"path", c.Path(),
			"status", c.Response().StatusCode(),
			"latency", time.Since(start).String(),
		)

		return err
	})
}
