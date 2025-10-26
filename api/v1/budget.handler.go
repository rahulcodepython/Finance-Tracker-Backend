// Package v1 provides handlers for version 1 of the API.
package v1

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/rahulcodepython/finance-tracker-backend/backend/database"
	"github.com/rahulcodepython/finance-tracker-backend/backend/services"
	"github.com/rahulcodepython/finance-tracker-backend/backend/utils"
)

// CreateBudgetInput defines the input structure for creating a new budget.
type CreateBudgetInput struct {
	Name   string  `json:"name" example:"Monthly Groceries"`
	Amount float64 `json:"amount" example:"500.00"`
}

// CreateBudget handles the creation of a new budget.
// @Summary      Create a new budget
// @Description  Creates a new budget for the authenticated user.
// @Tags         Budgets
// @Accept       json
// @Produce      json
// @Param        input body      CreateBudgetInput true "Create Budget Input"
// @Success      201  {object}  utils.Response
// @Failure      400  {object}  utils.Response
// @Failure      500  {object}  utils.Response
// @Router       /budgets/create [post]
func CreateBudget(c *fiber.Ctx) error {
	// Create a new instance of the CreateBudgetInput struct.
	var input CreateBudgetInput

	// Parse the request body into the input struct.
	if err := c.BodyParser(&input); err != nil {
		return utils.BadResponse(c, err, "Invalid request payload")
	}

	// Get the user ID from the context.
	userID, err := uuid.Parse(c.Locals("user_id").(string))
	if err != nil {
		return utils.BadResponse(c, err, "Invalid user ID")
	}

	// Get the database connection.
	db := database.DB

	// Call the CreateBudget service to create the new budget.
	budget, err := services.CreateBudget(userID, input.Name, input.Amount, db)
	if err != nil {
		return utils.InternalServerError(c, err, "Failed to create budget")
	}

	// Return a 201 Created response with the new budget.
	return utils.OKCreatedResponse(c, "Budget created successfully", budget)
}

// GetBudgets handles the retrieval of all budgets for the authenticated user.
// @Summary      Get all budgets
// @Description  Retrieves all budgets associated with the authenticated user.
// @Tags         Budgets
// @Produce      json
// @Success      200  {object}  utils.Response
// @Failure      400  {object}  utils.Response
// @Failure      500  {object}  utils.Response
// @Router       /budgets [get]
func GetBudgets(c *fiber.Ctx) error {
	// Get the user ID from the context.
	userID, err := uuid.Parse(c.Locals("user_id").(string))
	if err != nil {
		return utils.BadResponse(c, err, "Invalid user ID")
	}

	// Get the database connection.
	db := database.DB

	// Call the GetBudgets service to retrieve the user's budgets.
	budgets, err := services.GetBudgets(userID, db)
	if err != nil {
		return utils.InternalServerError(c, err, "Failed to retrieve budgets")
	}

	// Return a 200 OK response with the retrieved budgets.
	return utils.OKResponse(c, "Budgets retrieved successfully", budgets)
}

// UpdateBudgetInput defines the input structure for updating a budget.
type UpdateBudgetInput struct {
	Name   string  `json:"name" example:"Monthly Entertainment"`
	Amount float64 `json:"amount" example:"200.00"`
}

// UpdateBudget handles the update of a specific budget.
// @Summary      Update a budget
// @Description  Updates the details of a specific budget for the authenticated user.
// @Tags         Budgets
// @Accept       json
// @Produce      json
// @Param        id   path      string          true "Budget ID"
// @Param        input body      UpdateBudgetInput true "Update Budget Input"
// @Success      200  {object}  utils.Response
// @Failure      400  {object}  utils.Response
// @Failure      500  {object}  utils.Response
// @Router       /budgets/update/{id} [patch]
func UpdateBudget(c *fiber.Ctx) error {
	// Create a new instance of the UpdateBudgetInput struct.
	var input UpdateBudgetInput

	// Parse the request body into the input struct.
	if err := c.BodyParser(&input); err != nil {
		return utils.BadResponse(c, err, "Invalid request payload")
	}

	// Get the budget ID from the URL parameters.
	budgetID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.BadResponse(c, err, "Invalid budget ID")
	}

	// Get the database connection.
	db := database.DB

	// Call the UpdateBudget service to update the budget.
	budget, err := services.UpdateBudget(budgetID, input.Name, input.Amount, db)
	if err != nil {
		return utils.InternalServerError(c, err, "Failed to update budget")
	}

	// Return a 200 OK response with the updated budget.
	return utils.OKResponse(c, "Budget updated successfully", budget)
}

// DeleteBudget handles the deletion of a specific budget.
// @Summary      Delete a budget
// @Description  Deletes a specific budget for the authenticated user.
// @Tags         Budgets
// @Produce      json
// @Param        id   path      string  true "Budget ID"
// @Success      200  {object}  utils.Response
// @Failure      400  {object}  utils.Response
// @Failure      500  {object}  utils.Response
// @Router       /budgets/delete/{id} [delete]
func DeleteBudget(c *fiber.Ctx) error {
	// Get the budget ID from the URL parameters.
	budgetID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.BadResponse(c, err, "Invalid budget ID")
	}

	// Get the database connection.
	db := database.DB

	// Call the DeleteBudget service to delete the budget.
	if err := services.DeleteBudget(budgetID, db); err != nil {
		return utils.InternalServerError(c, err, "Failed to delete budget")
	}

	// Return a 200 OK response.
	return utils.OKResponse(c, "Budget deleted successfully", nil)
}
