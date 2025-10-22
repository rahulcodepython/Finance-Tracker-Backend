package v1

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/rahulcodepython/finance-tracker-backend/backend/database"
	"github.com/rahulcodepython/finance-tracker-backend/backend/services"
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
		CategoryID string  `json:"categoryId"`
		Amount     float64 `json:"amount"`
		Month      string  `json:"month"`
	}

	var input CreateBudgetInput

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": "Invalid request", "error": err.Error()})
	}

	userID, err := uuid.Parse(c.Locals("user_id").(string))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": "Invalid user ID", "error": err.Error()})
	}

	categoryID, err := uuid.Parse(input.CategoryID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": "Invalid category ID", "error": err.Error()})
	}

	month, err := time.Parse("2006-01", input.Month)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": "Invalid month format", "error": err.Error()})
	}

	db := database.DB

	budget, err := services.CreateBudget(userID, categoryID, input.Amount, month, db)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"success": false, "message": "Failed to create budget", "error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"success": true, "message": "Budget created successfully", "data": budget})
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
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": "Invalid user ID", "error": err.Error()})
	}

	db := database.DB

	budgets, err := services.GetBudgets(userID, db)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"success": false, "message": "Failed to get budgets", "error": err.Error()})
	}

	return c.JSON(fiber.Map{"success": true, "message": "Budgets retrieved successfully", "data": budgets})
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
		CategoryID string  `json:"categoryId"`
		Amount     float64 `json:"amount"`
		Month      string  `json:"month"`
	}

	var input UpdateBudgetInput

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": "Invalid request", "error": err.Error()})
	}

	budgetID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": "Invalid budget ID", "error": err.Error()})
	}

	categoryID, err := uuid.Parse(input.CategoryID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": "Invalid category ID", "error": err.Error()})
	}

	month, err := time.Parse("2006-01", input.Month)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": "Invalid month format", "error": err.Error()})
	}

	db := database.DB

	budget, err := services.UpdateBudget(budgetID, categoryID, input.Amount, month, db)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"success": false, "message": "Failed to update budget", "error": err.Error()})
	}

	return c.JSON(fiber.Map{"success": true, "message": "Budget updated successfully", "data": budget})
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
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": "Invalid budget ID", "error": err.Error()})
	}

	db := database.DB

	if err := services.DeleteBudget(budgetID, db); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"success": false, "message": "Failed to delete budget", "error": err.Error()})
	}

	return c.JSON(fiber.Map{"success": true, "message": "Budget deleted successfully"})
}
