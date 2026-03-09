package usecase

import (
	"context"
	"errors"
	"subscriptions/internal/models"
	"sync"
	"time"

	"github.com/gofrs/uuid"
)

type Repository interface {
	Save(ctx context.Context, sub models.Subscription) (uuid.UUID, error)
	GetByID(ctx context.Context, id uuid.UUID) (models.Subscription, error)
	GetByFilter(ctx context.Context, sub models.Subscription) ([]models.Subscription, error)
	Update(ctx context.Context, sub models.Subscription) error
	DeleteByID(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, limit, offset int) ([]models.Subscription, error)
}

type Usecase struct {
	repo Repository
}

func NewUsecase(repo Repository) *Usecase {
	return &Usecase{
		repo: repo,
	}
}

func (uc *Usecase) AddNewSubscription(ctx context.Context, sub models.Subscription) (uuid.UUID, error) {
	if !sub.Name.Valid || !sub.Price.Valid || !sub.StartDate.Valid || !sub.UserID.Valid {
		return uuid.Nil, errors.New("missing required parameter")
	}

	if sub.StartDate.Valid && sub.EndDate.Valid {
		if sub.EndDate.Value.Before(sub.StartDate.Value) {
			return uuid.Nil, models.ErrIncorrectDateRange
		}
	}

	id, err := uc.repo.Save(ctx, sub)
	if err != nil {
		return uuid.Nil, err
	}

	return id, nil
}

func (uc *Usecase) GetSubscriptionByID(ctx context.Context, id uuid.UUID) (models.Subscription, error) {
	if id.String() == "" {
		return models.Subscription{}, errors.New("id is required parameter")
	}

	sub, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return models.Subscription{}, err
	}

	return sub, nil
}

func (uc *Usecase) UpdateSubscription(ctx context.Context, sub models.Subscription) error {
	if sub.ID.String() == "" {
		return errors.New("id is required parameter")
	}

	if !sub.Name.Valid && !sub.StartDate.Valid && !sub.Price.Valid && !sub.EndDate.Valid {
		return models.ErrNoUpdateFields
	}

	if sub.StartDate.Valid && sub.EndDate.Valid {
		if sub.EndDate.Value.Before(sub.StartDate.Value) {
			return models.ErrIncorrectDateRange
		}
	}

	err := uc.repo.Update(ctx, sub)
	if err != nil {
		return err
	}

	return nil
}

func (uc *Usecase) DeleteByID(ctx context.Context, id uuid.UUID) error {
	if id.String() == "" {
		return errors.New("id is required parameter")
	}

	err := uc.repo.DeleteByID(ctx, id)
	if err != nil {
		return err
	}

	return nil
}

func (uc *Usecase) ListAllSubscriptions(ctx context.Context, limit, offset int) ([]models.Subscription, error) {
	return uc.repo.List(ctx, limit, offset)
}

func monthsBetween(a, b time.Time) int {
	y1, m1, _ := a.Date()
	y2, m2, _ := b.Date()

	return (y2-y1)*12 + int(m2-m1)
}

func startWorkers(workers int, wg *sync.WaitGroup, subs <-chan models.Subscription, result chan<- int) {
	for range workers {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for s := range subs {
				cost := s.Price.Value * monthsBetween(s.StartDate.Value, s.EndDate.Value)

				result <- cost
			}
		}()
	}

	// for s := range subs {
	// 	wg.Add(1)
	// 	go func() {
	// 		defer wg.Done()

	// 		cost := s.Price.Value * utils.MonthsBetween(s.StartDate.Value, s.EndDate.Value)

	// 		result <- cost
	// 	}()
	// }
}

func startWork(subs []models.Subscription, ch chan<- models.Subscription) {
	for _, s := range subs {
		ch <- s
	}

	close(ch)
}

func filtAndTransform(subs []models.Subscription, dateFrom, dateTo time.Time) []models.Subscription {
	i := 0
	for _, s := range subs {
		if !s.EndDate.Valid || s.EndDate.Value.After(dateFrom) {

			if s.StartDate.Valid {
				if s.StartDate.Value.Before(dateFrom) {
					s.StartDate.Value = dateFrom
				}
			}

			if !s.EndDate.Valid || dateTo.Before(s.EndDate.Value) {
				s.EndDate.Value = dateTo
			}

			subs[i] = s
			i++
		}
	}

	subs = subs[:i]

	return subs
}

func (uc *Usecase) GetTotalCost(ctx context.Context, params models.TotalCostParams) (int, error) {
	if !params.DateFrom.Valid || !params.DateTo.Valid {
		return 0, errors.New("missing required parameter")
	}

	subs, err := uc.repo.GetByFilter(ctx, models.Subscription{
		Name:   params.Name,
		UserID: params.UserID,
	})
	if err != nil {
		return 0, err
	}

	subs = filtAndTransform(subs, params.DateFrom.Value, params.DateTo.Value)

	subsCh := make(chan models.Subscription, len(subs))
	costCh := make(chan int, len(subs))
	wg := sync.WaitGroup{}

	go startWork(subs, subsCh)

	workers := 10
	startWorkers(workers, &wg, subsCh, costCh)

	go func() {
		wg.Wait()
		close(costCh)
	}()

	total := 0
	for c := range costCh {
		total += c
	}

	return total, nil
}
