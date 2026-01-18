package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/hashicorp/go-hclog"
	"subscription-service/internal/cerrors"
	"subscription-service/internal/model"
	"subscription-service/internal/service"
	"time"
)

type ISubscriptionsHandlers interface {
	FindByID(c *fiber.Ctx) error
	Create(c *fiber.Ctx) error
	Update(c *fiber.Ctx) error
	Delete(c *fiber.Ctx) error
	SubscriptionsSum(c *fiber.Ctx) error
}

type SubscriptionsHandler struct {
	app     *fiber.App
	logger  hclog.Logger
	service service.ISubscriptionService
}

func NewSubscriptionsHandler(app *fiber.App, logger hclog.Logger, service service.ISubscriptionService) ISubscriptionsHandlers {
	return &SubscriptionsHandler{app: app, logger: logger, service: service}
}

// FindByID Поиск подписки по ID
// @Summary Поиск подписки по ID
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param id path string true "Subscription ID"
// @Success 200 {object} model.Subscriptions
// @Failure 400 {object} cerrors.AppError
// @Failure 404 {object} cerrors.AppError
// @Failure 500 {object} cerrors.AppError
// @Router /v1/subscriptions/{id} [get]
func (h SubscriptionsHandler) FindByID(c *fiber.Ctx) error {
	idParam := c.Params("id")
	if idParam == "" {
		return cerrors.BadRequest("id is required")
	}

	id, err := uuid.Parse(idParam)
	if err != nil {
		h.logger.Warn("invalid subscription id", "id", idParam, "error", err)
		return cerrors.BadRequest(fmt.Sprintf("invalid subscription id: %s", idParam))
	}

	sub, err := h.service.FindByID(c.Context(), id)
	if err != nil {
		h.logger.Error("failed to get subscription by id", "error", err)
		return cerrors.Internal("failed to get subscription by id")
	}

	if sub == nil {
		return cerrors.NotFound("subscription not found")
	}

	return c.Status(fiber.StatusOK).JSON(*sub)
}

// Create Создание подписки
// @Summary Create Создание подписки
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param subscription body model.Subscriptions true "Subscription payload"
// @Success 200 {boolean} true
// @Failure 400 {object} cerrors.AppError
// @Failure 500 {object} cerrors.AppError
// @Router /v1/subscriptions [post]
func (h SubscriptionsHandler) Create(c *fiber.Ctx) error {
	req := new(model.Subscriptions)
	if err := c.BodyParser(req); err != nil {
		h.logger.Warn("invalid request body", "body", c.Body(), "error", err)
		return cerrors.BadRequest("invalid request body")
	}

	err := h.service.Create(c.Context(), *req)
	if err != nil {
		h.logger.Error("failed to create subscription", "error", err)
		return cerrors.Internal("failed to create subscription")
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"success": true})
}

// Update Обновление данных подписок
// @Summary Update Обновление данных подписок
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param subscription body model.Subscriptions true "Subscription payload"
// @Success 200 {boolean} true
// @Failure 400 {object} cerrors.AppError
// @Failure 500 {object} cerrors.AppError
// @Router /v1/subscriptions [put]
func (h SubscriptionsHandler) Update(c *fiber.Ctx) error {
	body := c.Body()
	req := new(model.Subscriptions)

	err := json.Unmarshal(body, &req)
	if err != nil {
		h.logger.Warn("invalid request body", "body", body)
		return cerrors.BadRequest("invalid request body")
	}

	err = h.service.Update(c.Context(), *req)
	if err != nil {
		h.logger.Error("failed to update subscription", "error", err)
		return cerrors.Internal("failed to update subscription")
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"success": true})
}

// Delete Удаление записей подписок по ID
// @Summary Delete Удаление подписок
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param id path string true "Subscription ID"
// @Success 200 {boolean} true
// @Failure 400 {object} cerrors.AppError
// @Failure 500 {object} cerrors.AppError
// @Router /v1/subscriptions/{id} [delete]
func (h SubscriptionsHandler) Delete(c *fiber.Ctx) error {
	idParam := c.Params("id")
	if idParam == "" {
		return cerrors.BadRequest("id is required")
	}

	id, err := uuid.Parse(idParam)
	if err != nil {
		h.logger.Warn("invalid subscription id", "id", idParam)
		return cerrors.BadRequest(fmt.Sprintf("invalid subscription id: %s", idParam))
	}

	err = h.service.Delete(c.Context(), id)
	if err != nil {
		h.logger.Error("failed to delete subscription", "error", err)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"success": true})
}

// SubscriptionsSum Сумма подписок в зависимости от филтров
// @Summary Высчет суммы подписок
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param startDate query string true "Start date (YYYY-MM-DD)"
// @Param endDate query string true "End date (YYYY-MM-DD)"
// @Param userID query string false "User UUID"
// @Param serviceName query string false "Service name"
// @Success 200 {object} map[string]uint64 "sum amount"
// @Failure 400 {object} cerrors.AppError
// @Failure 500 {object} cerrors.AppError
// @Router /v1/subscriptions/sum [get]
func (h SubscriptionsHandler) SubscriptionsSum(c *fiber.Ctx) error {
	req := new(model.SubscriptionsSum)
	err := c.QueryParser(req)
	if err != nil {
		h.logger.Warn("invalid request query", "query", req)
		return cerrors.BadRequest("invalid request body")
	}

	_, err = time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		return cerrors.BadRequest("invalid start date")
	}

	_, err = time.Parse("2006-01-02", req.EndDate)
	if err != nil {
		return cerrors.BadRequest("invalid end date")
	}

	amount, err := h.service.SubscriptionsSum(c.Context(), *req)
	if err != nil {
		h.logger.Error("failed to get subscriptions sum", "error", err)
		return cerrors.Internal("failed to get subscriptions sum")
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"amount": amount,
	})
}
