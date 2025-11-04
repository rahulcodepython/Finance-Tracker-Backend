package serializers

// CreateCategoryInput defines the input structure for creating a new transaction category.
type CategoryInput struct {
	Name string `json:"name" example:"Groceries"`
	Type string `json:"type" example:"Expense"`
}
