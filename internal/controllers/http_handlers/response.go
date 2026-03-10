package httphandlers

import (
	"subscriptions/internal/models"

	"github.com/gofrs/uuid"
)

// @Description Тело ответа
type AddNewSubscriptionResponse struct {
	// ID подписки в системе
	ID uuid.UUID `json:"id"`
} // @name AddNewSubscriptionResponse

// @Description Тело ответа
type GetSubscriptionResponse struct {
	// ID подписки в формате UUID
	ID uuid.UUID `json:"id" swaggertype:"string"`

	// ID пользователя в формате UUID
	UserID uuid.UUID `json:"user_id" swaggertype:"string"`

	// Дата начала подписки в формате год-месяц
	StartDate string `json:"start_date"`

	// Название подписки
	Name string `json:"name"`

	// Цена за подписку в месяц
	Price int `json:"price"`

	// Дата окончания подписки в формате год-месяц
	EndDate *string `json:"end_date,omitempty" swaggertype:"string"`
} // @name GetSubscriptionResponse

func newGetSubscriptionResponse(sub models.Subscription) GetSubscriptionResponse {
	resp := GetSubscriptionResponse{
		ID:        sub.ID,
		UserID:    sub.UserID.Value,
		Name:      sub.Name.Value,
		Price:     sub.Price.Value,
		StartDate: sub.StartDate.Value.Format("2006-01"),
	}

	if sub.EndDate.Valid {
		endDate := sub.EndDate.Value.Format("2006-01")
		resp.EndDate = &endDate
	}

	return resp
}

// @Description Тело ответа
type ListSubscriptionsResponse struct {
	// Список подписок
	List []GetSubscriptionResponse `json:"list"`
} // @name ListSubscriptionsResponse

// @Description Тело ответа
type TotalCostResponse struct {
	// Суммарная стоимость подписок
	Total int `json:"total"`
} // @name TotalCostResponse
