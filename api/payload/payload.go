package payload

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type RegisterPayload struct {
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name,omitempty"`
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=8"`
	Role      string `json:"role,omitempty"`
}

type LoginPayload struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type ProductPayload struct {
	Name         string  `json:"name" binding:"required"`
	Description  string  `json:"description" binding:"required"`
	Price        float64 `json:"price" binding:"required,gt=0"`
	UnitsInStock int     `json:"units_in_stock" binding:"required,gt=0"`
}

type OrderUpdatePayload struct {
	Status string `json:"status" validate:"required,oneof=pending completed cancelled"`
}

type OrderPayload struct {
	Items []OrderItemPayload `json:"items" validate:"required,dive"`
}

type OrderItemPayload struct {
	ProductID uuid.UUID `json:"product_id" validate:"required,uuid4"`
	Quantity  int       `json:"quantity" validate:"required,gt=0"`
	Price     float64   `json:"price" validate:"required,gt=0"`
}

func (o *OrderPayload) CalculateOrderTotal() decimal.Decimal {
	total := decimal.NewFromInt(0)

	for _, item := range o.Items {
		price := decimal.NewFromFloat(item.Price)
		itemTotal := price.Mul(decimal.NewFromInt(int64(item.Quantity)))
		total = total.Add(itemTotal)
	}

	return total
}
