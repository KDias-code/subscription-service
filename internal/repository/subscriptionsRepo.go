package repository

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"subscription-service/internal/model"
)

type ISubscriptionsRepository interface {
	FindByID(ctx context.Context, id uuid.UUID) (*model.Subscriptions, error)
	Create(ctx context.Context, sub *model.Subscriptions) error
	Update(ctx context.Context, sub *model.Subscriptions) error
	Delete(ctx context.Context, id uuid.UUID) error
	SumPrice(ctx context.Context, data model.SubscriptionsSum) (uint64, error)
}

type SubscriptionsRepo struct {
	db *sqlx.DB
}

func NewSubscriptionsRepository(db sqlx.DB) ISubscriptionsRepository {
	return &SubscriptionsRepo{
		db: &db,
	}
}

func (r SubscriptionsRepo) FindByID(ctx context.Context, id uuid.UUID) (*model.Subscriptions, error) {
	const query = `
		SELECT
			id,
			service_name,
			price,
			user_id,
			start_date,
			end_date
		FROM subscriptions
		WHERE id = $1
	`

	var sub model.Subscriptions

	err := r.db.GetContext(ctx, &sub, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &sub, nil
}

func (r SubscriptionsRepo) Create(
	ctx context.Context,
	sub *model.Subscriptions,
) error {
	const query = `
		INSERT INTO subscriptions (
			service_name,
			price,
			user_id,
			start_date,
			end_date
		) VALUES (
			:service_name,
			:price,
			:user_id,
			:start_date,
			:end_date
		)
		RETURNING id
	`

	rows, err := r.db.NamedQueryContext(ctx, query, sub)
	if err != nil {
		return fmt.Errorf("create subscription: %w", err)
	}
	defer rows.Close()

	if rows.Next() {
		if err := rows.Scan(&sub.ID); err != nil {
			return err
		}
	}

	return nil
}

func (r SubscriptionsRepo) Update(
	ctx context.Context,
	sub *model.Subscriptions,
) error {
	const query = `
		UPDATE subscriptions
		SET
			service_name = :service_name,
			price        = :price,
			user_id      = :user_id,
			start_date   = :start_date,
			end_date     = :end_date,
			updated_at   = now()
		WHERE id = :id
		RETURNING updated_at
	`

	rows, err := r.db.NamedQueryContext(ctx, query, sub)
	if err != nil {
		return fmt.Errorf("update subscription: %w", err)
	}
	defer rows.Close()

	if !rows.Next() {
		return sql.ErrNoRows
	}

	if err := rows.Scan(&sub.UpdatedAt); err != nil {
		return fmt.Errorf("scan updated_at: %w", err)
	}

	return nil
}

func (r SubscriptionsRepo) Delete(
	ctx context.Context,
	id uuid.UUID,
) error {
	const query = `
		DELETE FROM subscriptions
		WHERE id = $1
		RETURNING id
	`

	var deletedID uuid.UUID
	err := r.db.QueryRowContext(ctx, query, id).Scan(&deletedID)
	if err != nil {
		if err == sql.ErrNoRows {
			return sql.ErrNoRows
		}
		return fmt.Errorf("delete subscription: %w", err)
	}

	return nil
}

func (r SubscriptionsRepo) SumPrice(
	ctx context.Context,
	data model.SubscriptionsSum,
) (uint64, error) {
	query := `
		SELECT COALESCE(SUM(price),0)
		FROM subscriptions
		WHERE start_date <= $1 AND end_date >= $2
	`
	args := []interface{}{data.EndDate, data.StartDate} // обратил порядок: end >= start?

	paramIdx := 3 // следующий плейсхолдер

	if data.ID != nil {
		query += fmt.Sprintf(" AND user_id = $%d", paramIdx)
		args = append(args, *data.ID)
		paramIdx++
	}

	if data.ServiceName != nil {
		query += fmt.Sprintf(" AND service_name = $%d", paramIdx)
		args = append(args, *data.ServiceName)
		paramIdx++
	}

	var sum uint64
	if err := r.db.GetContext(ctx, &sum, query, args...); err != nil {
		return 0, fmt.Errorf("sum subscriptions price: %w", err)
	}

	return sum, nil
}
