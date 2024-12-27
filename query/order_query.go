package query

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type Order struct {
	ID          uuid.UUID       `json:"id" validate:"required,uuid4"`
	UserID      uuid.UUID       `json:"user_id" validate:"required,uuid4"`
	TotalAmount decimal.Decimal `json:"total_amount" validate:"required,gt=0"`
	Status      string          `json:"status" validate:"required,oneof=pending completed cancelled"`
	Items       []OrderItem     `json:"items" validate:"dive"`
	CreatedAt   time.Time       `json:"created_at" validate:"required"`
	UpdatedAt   time.Time       `json:"updated_at" validate:"required"`
}

type OrderItem struct {
	ID        uuid.UUID       `json:"id" validate:"required,uuid4"`
	OrderID   uuid.UUID       `json:"order_id" validate:"required,uuid4"`
	ProductID uuid.UUID       `json:"product_id" validate:"required,uuid4"`
	Quantity  int             `json:"quantity" validate:"required,gt=0"`
	Price     decimal.Decimal `json:"price" validate:"required,gt=0"`
	CreatedAt time.Time       `json:"created_at" validate:"required"`
}

func (q *Query) CreateOrder(ctx context.Context, order *Order) (string, error) {
	var orderID uuid.UUID
	query := `
        INSERT INTO "order" (user_id, total_amount)
        VALUES ($1, $2)
        RETURNING id;
    `
	err := q.DB.QueryRowContext(ctx, query, order.UserID, order.TotalAmount).Scan(&orderID)
	if err != nil {
		return "", err
	}
	for _, item := range order.Items {
		query = `
            INSERT INTO "order_item" (id, order_id, product_id, quantity, price, created_at)
            VALUES ($1, $2, $3, $4, $5, $6);
        `
		_, err := q.DB.ExecContext(ctx, query, uuid.New(), orderID, item.ProductID, item.Quantity, item.Price, item.CreatedAt)
		if err != nil {
			return "", err
		}
	}
	return orderID.String(), nil
}

func (q *Query) GetOrdersByUserID(ctx context.Context, userID uuid.UUID) ([]Order, error) {
	var orders []Order

	query := `
        SELECT id, user_id, status, total_amount, created_at
        FROM "order"
        WHERE user_id = $1
        ORDER BY created_at DESC;
    `
	rows, err := q.DB.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var order Order
		order.Items = []OrderItem{}

		err = rows.Scan(&order.ID, &order.UserID, &order.Status, &order.TotalAmount, &order.CreatedAt)
		if err != nil {
			return nil, err
		}

		itemsQuery := `
			SELECT id, order_id, product_id, quantity, price, created_at
			FROM "order_item"
			WHERE order_id = $1;
		`
		itemsRows, err := q.DB.QueryContext(ctx, itemsQuery, order.ID)
		if err != nil {
			return nil, err
		}
		defer itemsRows.Close()

		for itemsRows.Next() {
			var item OrderItem
			err = itemsRows.Scan(&item.ID, &item.OrderID, &item.ProductID, &item.Quantity, &item.Price, &item.CreatedAt)
			if err != nil {
				return nil, err
			}
			order.Items = append(order.Items, item)
		}

		orders = append(orders, order)
	}

	return orders, nil
}

func (q *Query) CancelOrderIfPending(ctx context.Context, orderID uuid.UUID) error {
	query := `
        UPDATE "order"
        SET status = 'cancelled', updated_at = CURRENT_TIMESTAMP
        WHERE id = $1 AND status = 'pending'
        RETURNING id;
    `
	var updatedOrderID string
	err := q.DB.QueryRowContext(ctx, query, orderID).Scan(&updatedOrderID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil
		}
		return err
	}
	return nil
}
func (q *Query) UpdateOrderStatus(ctx context.Context, orderID uuid.UUID, newStatus string) error {
	query := `
        UPDATE "order"
        SET status = $1, updated_at = CURRENT_TIMESTAMP
        WHERE id = $2
        RETURNING id;
    `
	var updatedOrderID string
	err := q.DB.QueryRowContext(ctx, query, newStatus, orderID).Scan(&updatedOrderID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil
		}
		return err
	}
	return nil
}
