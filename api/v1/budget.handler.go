package v1

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/rahulcodepython/finance-tracker-backend/backend/database"
	"github.com/rahulcodepython/finance-tracker-backend/backend/services"
	"github.com/rahulcodepython/finance-tracker-backend/backend/utils"
)

// CreateBudget godoc
// @Summary Create a new budget
// @Description Creates a new budget for the authenticated user.
// @Tags budgets
// @Security ApiKeyAuth
// @Accept  json
// @Produce  json
// @Param input body CreateBudgetInput true "Create Budget Input"
// @Success 201 {object} map[string]interface{} "Budget created successfully"
// @Router /budgets/create [post]
func CreateBudget(c *fiber.Ctx) error {
	type CreateBudgetInput struct {
		Name   string  `json:"name"`
		Amount float64 `json:"amount"`
	}

	var input CreateBudgetInput

	if err := c.BodyParser(&input); err != nil {
		return utils.BadResponse(c, err, "Invalid request")
	}

	userID, err := uuid.Parse(c.Locals("user_id").(string))
	if err != nil {
		return utils.BadResponse(c, err, "Invalid user ID")
	}

	db := database.DB

	budget, err := services.CreateBudget(userID, input.Name, input.Amount, db)
	if err != nil {
		return utils.InternalServerError(c, err, "Failed to create budget")
	}

	return utils.OKCreatedResponse(c, "Budget created successfully", budget)
}

// GetBudgets godoc
// @Summary Get all budgets
// @Description Gets all budgets for the authenticated user.
// @Tags budgets
// @Security ApiKeyAuth
// @Produce  json
// @Success 200 {object} map[string]interface{} "Budgets retrieved successfully"
// @Router /budgets [get]
func GetBudgets(c *fiber.Ctx) error {
	userID, err := uuid.Parse(c.Locals("user_id").(string))
	if err != nil {
		return utils.BadResponse(c, err, "Invalid user ID")
	}

	db := database.DB

	budgets, err := services.GetBudgets(userID, db)
	if err != nil {
		return utils.InternalServerError(c, err, "Failed to get budgets")
	}

	return utils.OKResponse(c, "Budgets retrieved successfully", budgets)
}

// UpdateBudget godoc
// @Summary Update a budget
// @Description Updates a budget for the authenticated user.
// @Tags budgets
// @Security ApiKeyAuth
// @Accept  json
// @Produce  json
// @Param id path string true "Budget ID"
// @Param input body UpdateBudgetInput true "Update Budget Input"
// @Success 200 {object} map[string]interface{} "Budget updated successfully"
// @Router /budgets/update/{id} [patch]
func UpdateBudget(c *fiber.Ctx) error {
	type UpdateBudgetInput struct {
		Name   string  `json:"name"`
		Amount float64 `json:"amount"`
	}

	var input UpdateBudgetInput

	if err := c.BodyParser(&input); err != nil {
		return utils.BadResponse(c, err, "Invalid request")
	}

	budgetID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.BadResponse(c, err, "Invalid budget ID")
	}

	db := database.DB

	budget, err := services.UpdateBudget(budgetID, input.Name, input.Amount, db)
	if err != nil {
		return utils.InternalServerError(c, err, "Failed to update budget")
	}

	return utils.OKResponse(c, "Budget updated successfully", budget)
}

// DeleteBudget godoc
// @Summary Delete a budget
// @Description Deletes a budget for the authenticated user.
// @Tags budgets
// @Security ApiKeyAuth
// @Produce  json
// @Param id path string true "Budget ID"
// @Success 200 {object} map[string]interface{} "Budget deleted successfully"
// @Router /budgets/delete/{id} [delete]
func DeleteBudget(c *fiber.Ctx) error {
	budgetID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.BadResponse(c, err, "Invalid budget ID")
	}

	db := database.DB

	if err := services.DeleteBudget(budgetID, db); err != nil {
		return utils.InternalServerError(c, err, "Failed to delete budget")
	}

	return utils.OKResponse(c, "Budget deleted successfully", nil)
}
