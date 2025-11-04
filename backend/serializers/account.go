package serializers

// CreateAccountInput defines the input structure for creating a new financial account.
type CreateAccountInput struct {
	Name    string  `json:"name" example:"My Savings Account"`
	Type    string  `json:"type" example:"Savings"`
	Balance float64 `json:"balance" example:"1000.00"`
}

// UpdateAccountInput defines the input structure for updating a financial account.
type UpdateAccountInput struct {
	Name     string `json:"name" example:"My Updated Savings Account"`
	Type     string `json:"type" example:"Checking"`
	IsActive bool   `json:"isActive" example:"true"`
}
