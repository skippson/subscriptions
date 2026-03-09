package models

import (
	"time"

	"github.com/gofrs/uuid"
)

type Subscription struct {
	ID        uuid.UUID
	UserID    Optional[uuid.UUID]
	Name      Optional[string]
	Price     Optional[int]
	StartDate Optional[time.Time]
	EndDate   Optional[time.Time]
}

type Optional[T any] struct {
	Value T
	Valid bool
}

type TotalCostParams struct {
	UserID   Optional[uuid.UUID]
	Name     Optional[string]
	DateFrom Optional[time.Time]
	DateTo   Optional[time.Time]
}
