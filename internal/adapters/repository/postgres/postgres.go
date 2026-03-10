package postgres

import (
	"context"
	"errors"
	"fmt"
	"subscriptions/config"
	"subscriptions/internal/models"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresRepository struct {
	pool *pgxpool.Pool
}

func NewRepository(ctx context.Context, config config.Postgres) (*PostgresRepository, error) {
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		config.User,
		config.Password,
		config.Host,
		config.Port,
		config.Name,
		config.SSLMode,
	)

	poolCfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, err
	}

	poolCfg.MaxConns = int32(config.MaxConns)
	poolCfg.MinConns = int32(config.MinConns)

	pool, err := pgxpool.NewWithConfig(ctx, poolCfg)
	if err != nil {
		return nil, err
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, err
	}

	return &PostgresRepository{
		pool: pool,
	}, nil
}

func (r *PostgresRepository) Save(ctx context.Context, sub models.Subscription) (uuid.UUID, error) {
	query := `insert into subscriptions(name, price, user_id, start_date, end_date)
	values ($1, $2, $3, $4, $5)
	returning id`

	dto := mapToDTO(sub)

	var id uuid.UUID
	err := r.pool.QueryRow(ctx, query, dto.Name, dto.Price,
		dto.UserID, dto.StartDate, dto.EndDate).
		Scan(&id)
	if err != nil {
		return uuid.Nil, err
	}

	return id, nil
}

func (r *PostgresRepository) GetByID(ctx context.Context, id uuid.UUID) (models.Subscription, error) {
	query := `select id, name, price, user_id, start_date, end_date
	from subscriptions
	where id = $1`

	dto := subscriptionDB{}
	err := pgxscan.Get(ctx, r.pool, &dto, query, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.Subscription{}, models.ErrSubscriptionNotFound
		}

		return models.Subscription{}, err
	}

	return dto.mapToModel(), nil
}

func (r *PostgresRepository) GetByFilter(ctx context.Context, sub models.Subscription) ([]models.Subscription, error) {
	builder := squirrel.Select("id, name, price, user_id, start_date, end_date").
		From("subscriptions").
		PlaceholderFormat(squirrel.Dollar).
		OrderBy("start_date desc")

	if sub.Name.Valid {
		builder = builder.Where("name = ?", sub.Name.Value)
	}

	if sub.Price.Valid {
		builder = builder.Where("price = ?", sub.Price.Value)
	}

	if sub.UserID.Valid {
		builder = builder.Where("user_id = ?", sub.UserID.Value)
	}

	if sub.StartDate.Valid {
		builder = builder.Where("start_date >= ?", sub.StartDate.Value)
	}

	if sub.EndDate.Valid {
		builder = builder.Where("(end_date is null or end_date <= ?)", sub.EndDate.Value)
	}

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, err
	}

	row := make([]subscriptionDB, 0)
	err = pgxscan.Select(ctx, r.pool, &row, query, args...)
	if err != nil {
		return nil, err
	}

	result := make([]models.Subscription, 0, len(row))
	for _, r := range row {
		result = append(result, r.mapToModel())
	}

	return result, nil
}

func (r *PostgresRepository) Update(ctx context.Context, sub models.Subscription) error {
	builder := squirrel.Update("subscriptions").
		PlaceholderFormat(squirrel.Dollar).
		Where("id = ?", sub.ID).
		Set("updated_at", time.Now())

	if sub.Name.Valid {
		builder = builder.Set("name", sub.Name.Value)
	}

	if sub.Price.Valid {
		builder = builder.Set("price", sub.Price.Value)
	}

	if sub.StartDate.Valid {
		builder = builder.Set("start_date", sub.StartDate.Value)
	}

	if sub.EndDate.Valid {
		builder = builder.Set("end_date", sub.EndDate.Value)
	}

	query, args, err := builder.ToSql()
	if err != nil {
		return err
	}

	cmd, err := r.pool.Exec(ctx, query, args...)
	if err != nil {
		return err
	}

	if cmd.RowsAffected() == 0 {
		return models.ErrSubscriptionNotFound
	}

	return nil
}

func (r *PostgresRepository) DeleteByID(ctx context.Context, id uuid.UUID) error {
	query := `delete from subscriptions where id = $1`

	cmd, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	if cmd.RowsAffected() == 0 {
		return models.ErrSubscriptionNotFound
	}

	return nil
}

func (r *PostgresRepository) List(ctx context.Context, limit, offset int) ([]models.Subscription, error) {
	query := `select id, name, price, user_id, start_date, end_date
	from subscriptions
	order by start_date desc
	offset $1
	limit $2`

	rows := make([]subscriptionDB, 0)
	err := pgxscan.Select(ctx, r.pool, &rows, query, offset, limit)
	if err != nil {
		return nil, err
	}

	result := make([]models.Subscription, 0, len(rows))
	for _, r := range rows {
		result = append(result, r.mapToModel())
	}

	return result, nil
}

func (r *PostgresRepository) Close() {
	r.pool.Close()
}
