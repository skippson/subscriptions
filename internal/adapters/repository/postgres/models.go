package postgres

import (
	"subscriptions/internal/models"
	"time"

	"github.com/gofrs/uuid"
)

type subscriptionDB struct {
	ID        uuid.UUID  `db:"id"`
	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt time.Time  `db:"updated_at"`
	UserID    *uuid.UUID `db:"user_id"`
	Name      *string    `db:"name"`
	Price     *int       `db:"price"`
	StartDate *time.Time `db:"start_date"`
	EndDate   *time.Time `db:"end_date"`
}

func (s subscriptionDB) mapToModel() models.Subscription {
	return models.Subscription{
		ID:        s.ID,
		UserID:    models.PtrToOptional(s.UserID),
		Name:      models.PtrToOptional(s.Name),
		Price:     models.PtrToOptional(s.Price),
		StartDate: models.PtrToOptional(s.StartDate),
		EndDate:   models.PtrToOptional(s.EndDate),
	}
}

func mapToDTO(sub models.Subscription) subscriptionDB {
	return subscriptionDB{
		ID:        sub.ID,
		UserID:    models.OptionalToPtr(sub.UserID),
		Name:      models.OptionalToPtr(sub.Name),
		Price:     models.OptionalToPtr(sub.Price),
		StartDate: models.OptionalToPtr(sub.StartDate),
		EndDate:   models.OptionalToPtr(sub.EndDate),
	}
}