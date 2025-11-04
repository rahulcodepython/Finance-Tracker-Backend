package serializers

// CreateBudgetInput defines the input structure for creating a new budget.
type BudgetInput struct {
	Name   string  `json:"name" example:"Monthly Groceries"`
	Amount float64 `json:"amount" example:"500.00"`
}
