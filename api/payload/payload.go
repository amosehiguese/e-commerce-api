package payload

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

type OrderPayload struct {
	ProductIDs []int `json:"product_ids" binding:"required,min=1"` // IDs of products being ordered
	Quantity   []int `json:"quantity" binding:"required,min=1"`    // Corresponding quantities
}
