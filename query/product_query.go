package query

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type Product struct {
	ID           uuid.UUID       `json:"id" validate:"required,uuid4"`
	Name         string          `json:"name" validate:"required,min=3,max=255"`
	Description  *string         `json:"description,omitempty" validate:"omitempty,max=1000"`
	Price        decimal.Decimal `json:"price" validate:"required,gt=0"`
	UnitsInStock int             `json:"units_in_stock" validate:"required,gte=0"`
	CreatedAt    time.Time       `json:"created_at" validate:"required"`
	UpdatedAt    time.Time       `json:"updated_at" validate:"required"`
}

// CreateProduct inserts a new product into the database.
func (q *Query) CreateProduct(ctx context.Context) error {
	var product Product
	product.ID = uuid.New()
	product.CreatedAt = time.Now()
	product.UpdatedAt = time.Now()

	query := `
		INSERT INTO "product" (id, name, description, price, units_in_stock, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	_, err := q.DB.ExecContext(ctx, query, product.ID, product.Name, product.Description, product.Price, product.UnitsInStock, product.CreatedAt, product.UpdatedAt)
	return err
}

// GetProductByID fetches a product by its ID.
func (q *Query) GetProductByID(ctx context.Context, productID uuid.UUID) (*Product, error) {
	query := `
		SELECT id, name, description, price, units_in_stock, created_at, updated_at
		FROM "product"
		WHERE id = $1
	`
	var product Product
	row := q.DB.QueryRowContext(ctx, query, productID)
	err := row.Scan(&product.ID, &product.Name, &product.Description, &product.Price, &product.UnitsInStock, &product.CreatedAt, &product.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &product, nil
}

// GetAllProducts fetches all products from the database.
func (q *Query) GetAllProducts(ctx context.Context, db *sql.DB) ([]Product, error) {
	query := `
		SELECT id, name, description, price, units_in_stock, created_at, updated_at
		FROM "product"
	`
	rows, err := q.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []Product
	for rows.Next() {
		var product Product
		err := rows.Scan(&product.ID, &product.Name, &product.Description, &product.Price, &product.UnitsInStock, &product.CreatedAt, &product.UpdatedAt)
		if err != nil {
			return nil, err
		}
		products = append(products, product)
	}
	return products, nil
}

// UpdateProduct updates the product's details in the database.
func (q *Query) UpdateProduct(ctx context.Context) error {
	var product Product
	query := `
		UPDATE "product"
		SET name = $1, description = $2, price = $3, units_in_stock = $4, updated_at = $5
		WHERE id = $6
	`
	_, err := q.DB.ExecContext(ctx, query, product.Name, product.Description, product.Price, product.UnitsInStock, product.UpdatedAt, product.ID)
	return err
}

// DeleteProduct deletes a product by its ID.
func (q *Query) DeleteProduct(ctx context.Context, productID uuid.UUID) error {
	query := `
		DELETE FROM "product"
		WHERE id = $1
	`
	_, err := q.DB.ExecContext(ctx, query, productID)
	return err
}