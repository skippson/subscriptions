package httphandlers

import (
	"context"
	"errors"
	"subscriptions/internal/models"
	"subscriptions/pkg/logger"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofrs/uuid"
)

type Usecase interface {
	AddNewSubscription(ctx context.Context, sub models.Subscription) (uuid.UUID, error)
	GetSubscriptionByID(ctx context.Context, id uuid.UUID) (models.Subscription, error)
	UpdateSubscription(ctx context.Context, sub models.Subscription) error
	DeleteByID(ctx context.Context, id uuid.UUID) error
	ListAllSubscriptions(ctx context.Context, limit, offset int) ([]models.Subscription, error)
	GetTotalCost(ctx context.Context, params models.TotalCostParams) (int, error)
}

type paramsType int

const (
	queryParams paramsType = 1
	jsonParams  paramsType = 2
)

type ApiHandlers struct {
	uc        Usecase
	validator *validator.Validate
}

func NewHandlers(uc Usecase) *ApiHandlers {
	v := validator.New()
	v.RegisterValidation("year-month", func(fl validator.FieldLevel) bool {
		s := fl.Field().String()
		if s == "" {
			return true
		}

		_, err := time.Parse("2006-01", fl.Field().String())
		return err == nil
	})

	return &ApiHandlers{
		uc:        uc,
		validator: v,
	}
}

func (h *ApiHandlers) readRequest(c *fiber.Ctx, mode paramsType, req any) error {
	switch mode {
	case queryParams:
		if err := c.QueryParser(req); err != nil {
			return errors.New("invalid query")
		}
	case jsonParams:
		if err := c.BodyParser(req); err != nil {
			return errors.New("invalid json")
		}
	default:
	}

	if err := h.validator.Struct(req); err != nil {
		return errors.New("validation error")
	}

	return nil
}

// AddNewSubscription godoc
//
//	@Summary		Добавление новой подписки
//	@Description	Создает запись о новой подписке
//	@Tags			Subscription
//	@Accept			json
//	@Produce		json
//	@Param			request			body		AddNewSubscriptionParams			true				"Тело запроса"
//	@Success		200				{object}	SuccessResponse{data=AddNewSubscriptionResponse}		"OK"
//	@Failure		400				{object}	ErrorResponse											"Bad request"
//	@Failure		500				{object}	ErrorResponse											"Internal Server Error"
//	@Router			/api/subscription/ [post]
func (h *ApiHandlers) AddNewSubscription() fiber.Handler {
	return func(c *fiber.Ctx) error {
		req := AddNewSubscriptionParams{}
		if err := h.readRequest(c, jsonParams, &req); err != nil {
			return writeError(c, fiber.StatusBadRequest, err.Error())
		}

		id, err := h.uc.AddNewSubscription(c.Context(), req.mapToModel())
		if err != nil {
			if errors.Is(err, models.ErrIncorrectDateRange) {
				return writeError(c, fiber.StatusBadRequest, "incorrect date range")
			}

			getLogger(c).Error("add new subscription failed",
				logger.Field{Key: "error", Value: err},
				logger.Field{Key: "params", Value: req})

			return writeError(c, fiber.StatusInternalServerError, "internal error")
		}

		return writeSuccess(c, fiber.StatusOK, AddNewSubscriptionResponse{ID: id})
	}
}

// GetSubscriptionByID godoc
//
//	@Summary		Получение подписки
//	@Description	Получение записи о подписке по ее ID
//	@Tags			Subscription
//	@Produce		json
//	@Param			id				path		string						true						"ID подписки"
//	@Success		200				{object}	SuccessResponse{data=GetSubscriptionResponse}			"OK"
//	@Failure		400				{object}	ErrorResponse											"Bad request"
//	@Failure		404				{object}	ErrorResponse											"Not found"
//	@Failure		500				{object}	ErrorResponse											"Internal Server Error"
//	@Router			/api/subscription/{id} [get]
func (h *ApiHandlers) GetSubscriptionByID() fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := uuid.FromString(c.Params("id"))
		if err != nil {
			return writeError(c, fiber.StatusBadRequest, "invalid id")
		}

		sub, err := h.uc.GetSubscriptionByID(c.Context(), id)
		if err != nil {
			if errors.Is(err, models.ErrSubscriptionNotFound) {
				return writeError(c, fiber.StatusNotFound, "not found")
			}

			getLogger(c).Error("get subscription failed",
				logger.Field{Key: "id", Value: id.String()},
				logger.Field{Key: "error", Value: err})

			return writeError(c, fiber.StatusInternalServerError, "internal error")
		}

		return writeSuccess(c, fiber.StatusOK, newGetSubscriptionResponse(sub))
	}
}

// UpdateSubscription godoc
//
//	@Summary		Обновление подписки
//	@Description	Обновляет запись о подписке по ее ID
//	@Tags			Subscription
//	@Accept			json
//	@Produce		json
//	@Param			id				path		string						true	"ID подписки"
//	@Param			request			body		UpdateSubscriptionParams	true	"Тело запроса"
//	@Success		200																"OK"
//	@Failure		400				{object}	ErrorResponse						"Bad request"
//	@Failure		404				{object}	ErrorResponse						"Not found"
//	@Failure		500				{object}	ErrorResponse						"Internal Server Error"
//	@Router			/api/subscription/{id} [patch]
func (h *ApiHandlers) UpdateSubscription() fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := uuid.FromString(c.Params("id"))
		if err != nil {
			return writeError(c, fiber.StatusBadRequest, "invalid id")
		}

		req := UpdateSubscriptionParams{}
		if err := h.readRequest(c, jsonParams, &req); err != nil {
			return writeError(c, fiber.StatusBadRequest, err.Error())
		}

		if err := h.uc.UpdateSubscription(c.Context(), req.mapToModel(id)); err != nil {
			if errors.Is(err, models.ErrSubscriptionNotFound) {
				return writeError(c, fiber.StatusNotFound, "not found")
			}

			if errors.Is(err, models.ErrNoUpdateFields) {
				return writeError(c, fiber.StatusBadRequest, "no fields to update")
			}

			if errors.Is(err, models.ErrIncorrectDateRange) {
				return writeError(c, fiber.StatusBadRequest, "incorrect date range")
			}

			getLogger(c).Error("update failed",
				logger.Field{Key: "error", Value: err},
				logger.Field{Key: "id", Value: id.String()},
				logger.Field{Key: "params", Value: req})

			return writeError(c, fiber.StatusInternalServerError, "update failed")
		}

		return c.SendStatus(fiber.StatusOK)
	}
}

// DeleteSubscription godoc
//
//	@Summary		Удаление подписки
//	@Description	Удаление записи о подписке по ее ID
//	@Tags			Subscription
//	@Produce		json
//	@Param			id				path		string			true		"ID подписки"
//	@Success		200														"OK"
//	@Failure		400				{object}	ErrorResponse				"Bad request"
//	@Failure		404				{object}	ErrorResponse				"Not found"
//	@Failure		500				{object}	ErrorResponse				"Internal Server Error"
//	@Router			/api/subscription/{id} [delete]
func (h *ApiHandlers) DeleteSubscription() fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := uuid.FromString(c.Params("id"))
		if err != nil {
			return writeError(c, fiber.StatusBadRequest, "invalid id")
		}

		if err := h.uc.DeleteByID(c.Context(), id); err != nil {
			if errors.Is(err, models.ErrSubscriptionNotFound) {
				return writeError(c, fiber.StatusNotFound, "not found")
			}

			getLogger(c).Error("delete subscription failed",
				logger.Field{Key: "id", Value: id.String()},
				logger.Field{Key: "error", Value: err})

			return writeError(c, fiber.StatusInternalServerError, "internal error")
		}

		return c.SendStatus(fiber.StatusOK)
	}
}

// ListSubscriptions godoc
//
//	@Summary		Список подписок
//	@Description	Получение всех записей о подписках
//	@Tags			Subscription
//	@Produce		json
//	@Param 			limit 			query 		int 					true 						"Количество записей"
//	@Param 			offset 			query 		int 					true 						"Смещение"
//	@Success		200				{object}	SuccessResponse{data=ListSubscriptionsResponse}		"OK"
//	@Success		204																				"No content"
//	@Failure		400				{object}	ErrorResponse										"Bad request"
//	@Failure		500				{object}	ErrorResponse										"Internal Server Error"
//	@Router			/api/subscription/ [get]
func (h *ApiHandlers) ListSubscriptions() fiber.Handler {
	return func(c *fiber.Ctx) error {
		q := listSubscriptionsParams{}
		if err := h.readRequest(c, queryParams, &q); err != nil {
			return writeError(c, fiber.StatusBadRequest, err.Error())
		}

		subs, err := h.uc.ListAllSubscriptions(c.Context(), q.Limit, q.Offset)
		if err != nil {
			getLogger(c).Error("delete subscription failed",
				logger.Field{Key: "limit", Value: q.Limit},
				logger.Field{Key: "offset", Value: q.Offset},
				logger.Field{Key: "error", Value: err})

			return writeError(c, fiber.StatusInternalServerError, "internal error")
		}

		if len(subs) == 0 {
			c.SendStatus(fiber.StatusNoContent)
		}

		resp := make([]GetSubscriptionResponse, 0, len(subs))
		for _, s := range subs {
			resp = append(resp, newGetSubscriptionResponse(s))
		}

		return writeSuccess(c, fiber.StatusOK, ListSubscriptionsResponse{List: resp})
	}
}

// GetTotalCost godoc
//
//	@Summary		Суммарная стоимость
//	@Description	Посчитать суммарную стоимость подписок за период
//	@Tags			Subscription
//	@Produce		json
//	@Param 			from 			query 		string 			true 						"Дата начала подписки в формате год-месяц"
//	@Param 			to 				query 		string 			true 						"Дата окончания подписки в формате год-месяц"
//	@Param 			user_id 		query 		string 			false 						"ID пользователя в формате UUID"
//	@Param 			name 			query 		string 			false 						"Название подписки"
//	@Success		200				{object}	SuccessResponse{data=TotalCostResponse}		"OK"
//	@Failure		400				{object}	ErrorResponse								"Bad request"
//	@Failure		500				{object}	ErrorResponse								"Internal Server Error"
//	@Router			/api/subscription/total_cost [get]
func (h *ApiHandlers) GetTotalCost() fiber.Handler {
	return func(c *fiber.Ctx) error {
		q := totalCostParams{}
		if err := h.readRequest(c, queryParams, &q); err != nil {
			return writeError(c, fiber.StatusBadRequest, err.Error())
		}

		total, err := h.uc.GetTotalCost(c.Context(), q.mapToModel())
		if err != nil {
			getLogger(c).Error("calculate total cost failed",
				logger.Field{Key: "error", Value: err},
				logger.Field{Key: "params", Value: q})

			return writeError(c, fiber.StatusInternalServerError, "internal error")
		}

		return writeSuccess(c, fiber.StatusOK, TotalCostResponse{Total: total})
	}
}
