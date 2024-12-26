package query

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type Order struct {
	ID          uuid.UUID       `json:"id" validate:"required,uuid4"`
	UserID      uuid.UUID       `json:"user_id" validate:"required,uuid4"`
	TotalAmount decimal.Decimal `json:"total_amount" validate:"required,gt=0"`
	Status      string          `json:"status" validate:"required,oneof=pending completed cancelled"`
	CreatedAt   time.Time       `json:"created_at" validate:"required"`
	UpdatedAt   time.Time       `json:"updated_at" validate:"required"`
}

// CreateOrder inserts a new order into the database.
func (q *Query) CreateOrder(ctx context.Context) error {
	var order Order
	query := `
		INSERT INTO "order" (id, user_id, total_amount, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err := q.DB.ExecContext(ctx, query, order.ID, order.UserID, order.TotalAmount, order.Status, order.CreatedAt, order.UpdatedAt)
	return err
}
