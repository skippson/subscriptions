package httphandlers

import (
	"subscriptions/internal/models"
	"time"

	"github.com/gofrs/uuid"
)

// @Description Тело запроса
type AddNewSubscriptionParams struct {
	// ID пользователя в формате UUID
	UserID uuid.UUID `json:"user_id" validate:"required" swaggertype:"string"`

	// Дата начала подписки в формате год-месяц
	StartDate string `json:"start_date" validate:"required,year-month"`

	// Название подписки
	Name string `json:"name" validate:"required"`

	// Цена за подписку в месяц
	Price int `json:"price" validate:"gte=0"`

	// Дата окончания подписки в формате год-месяц
	EndDate *string `json:"end_date,omitempty" validate:"omitempty,year-month" swaggertype:"string"`
} // @name AddNewSubscriptionParams

func (p AddNewSubscriptionParams) mapToModel() models.Subscription {
	var (
		start, _ = time.Parse("2006-01", p.StartDate)
		end      models.Optional[time.Time]
	)

	if p.EndDate != nil {
		t, _ := time.Parse("2006-01", *p.EndDate)
		end = models.ValueToOption(t)
	}

	return models.Subscription{
		UserID:    models.ValueToOption(p.UserID),
		StartDate: models.ValueToOption(start),
		Name:      models.ValueToOption(p.Name),
		Price:     models.ValueToOption(p.Price),
		EndDate:   end,
	}
}

// @Description Тело запроса
type UpdateSubscriptionParams struct {
	// Название подписки
	Name *string `json:"name,omitempty" validate:"omitempty" swaggertype:"string"`

	// Дата начала подписки в формате год-месяц
	StartDate *string `json:"start_date,omitempty" validate:"omitempty,year-month" swaggertype:"string"`

	// Дата окончания подписки в формате год-месяц
	EndDate *string `json:"end_date,omitempty" validate:"omitempty,year-month" swaggertype:"string"`

	// Цена за подписку в месяц
	Price *int `json:"price,omitempty" validate:"omitempty,gte=0" swaggertype:"integer"`
} // @name UpdateSubscriptionParams

func (p UpdateSubscriptionParams) mapToModel(id uuid.UUID) models.Subscription {
	sub := models.Subscription{
		ID:    id,
		Name:  models.PtrToOptional(p.Name),
		Price: models.PtrToOptional(p.Price),
	}

	if p.StartDate != nil {
		t, _ := time.Parse("2006-01", *p.StartDate)
		sub.StartDate = models.ValueToOption(t)
	}

	if p.EndDate != nil {
		t, _ := time.Parse("2006-01", *p.EndDate)
		sub.EndDate = models.ValueToOption(t)
	}

	return sub
}

type listSubscriptionsParams struct {
	Limit  int `query:"limit" validate:"required,gt=0"`
	Offset int `query:"offset" validate:"gte=0"`
}

type totalCostParams struct {
	From   string     `query:"from" validate:"required,year-month"`
	To     string     `query:"to" validate:"required,year-month"`
	UserID *uuid.UUID `query:"user_id"`
	Name   *string    `query:"name"`
}

func (q totalCostParams) mapToModel() models.TotalCostParams {
	from, _ := time.Parse("2006-01", q.From)
	to, _ := time.Parse("2006-01", q.To)

	return models.TotalCostParams{
		UserID:   models.PtrToOptional(q.UserID),
		Name:     models.PtrToOptional(q.Name),
		DateFrom: models.ValueToOption(from),
		DateTo:   models.ValueToOption(to),
	}
}
