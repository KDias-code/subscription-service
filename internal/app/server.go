package app

import (
	"context"
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/hashicorp/go-hclog"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	fiberSwagger "github.com/swaggo/fiber-swagger"
	"log"
	"os"
	"os/signal"
	"subscription-service/configs"
	"subscription-service/internal/cerrors"
	"subscription-service/internal/handlers"
	"subscription-service/internal/repository"
	"subscription-service/internal/service"
	"syscall"
	"time"
)

const (
	loggerName = "subscriptions-service"
)

type Server struct {
	logger   hclog.Logger
	app      *fiber.App
	handlers handlers.ISubscriptionsHandlers
}

func Start(cfg configs.Config) error {
	s := new(Server)

	s.logger = hclog.New(&hclog.LoggerOptions{
		Name:       loggerName,
		Level:      hclog.Level(cfg.Log.Level),
		JSONFormat: cfg.Log.Json,
		Output:     os.Stdout,
	})

	db, err := sqlx.Connect("postgres", cfg.Postgres.Dsn)
	if err != nil {
		s.logger.Error("failed to connect to database", "error", err)
		return err
	}

	s.app = fiber.New(fiber.Config{
		AppName:      loggerName,
		ReadTimeout:  cfg.Http.ReadTimeout,
		WriteTimeout: cfg.Http.WriteTimeout,

		ErrorHandler: errorHandler,
	})

	RegisterMiddlewares(s.app, s.logger)

	subscriptionRepo := repository.NewSubscriptionsRepository(*db)

	subscriptionService := service.NewSubscriptionsService(subscriptionRepo)

	s.handlers = handlers.NewSubscriptionsHandler(s.app, s.logger, subscriptionService)

	s.router()
	s.app.Get("/swagger/*", fiberSwagger.WrapHandler)

	addr := fmt.Sprintf(":%d", cfg.Http.Port)
	s.logger.Info("starting http server", "addr", addr)

	go func() {
		if err = s.app.Listen(addr); err != nil {
			if err != nil {
				s.logger.Error("fiber listen error", "error", err)
			}
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := s.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}

	return err
}

func (s Server) Shutdown(ctx context.Context) error {
	s.logger.Info("shutting down server")

	done := make(chan error, 1)

	go func() {
		done <- s.app.Shutdown()
	}()

	select {
	case err := <-done:
		if err != nil {
			s.logger.Error("server shutdown failed", "error", err)
			return err
		}
		s.logger.Info("server stopped gracefully")
		return nil

	case <-ctx.Done():
		s.logger.Warn("shutdown timeout exceeded")
		return ctx.Err()
	}
}

func errorHandler(c *fiber.Ctx, err error) error {
	var appErr *cerrors.AppError
	if errors.As(err, &appErr) {
		return c.Status(appErr.Status).JSON(fiber.Map{
			"error": appErr,
		})
	}

	if fe, ok := err.(*fiber.Error); ok {
		return c.Status(fe.Code).JSON(fiber.Map{
			"error": fiber.Map{
				"code":    "FIBER_ERROR",
				"message": fe.Message,
			},
		})
	}

	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
		"error": fiber.Map{
			"code":    "UNKNOWN_ERROR",
			"message": "internal server error",
		},
	})
}
