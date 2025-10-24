package v1

import (
	"database/sql"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/rahulcodepython/finance-tracker-backend/backend/database"
	"github.com/rahulcodepython/finance-tracker-backend/backend/models"
	"github.com/rahulcodepython/finance-tracker-backend/backend/services"
)

// CreateRecurringTransaction godoc
// @Summary Create a new recurring transaction
// @Description Creates a new recurring transaction for the authenticated user.
// @Tags recurring-transactions
// @Security ApiKeyAuth
// @Accept  json
// @Produce  json
// @Param input body CreateRecurringTransactionInput true "Create Recurring Transaction Input"
// @Success 201 {object} map[string]interface{} "Recurring transaction created successfully"
// @Router /recurring-transactions/create [post]
func CreateRecurringTransaction(c *fiber.Ctx) error {
	type CreateRecurringTransactionInput struct {
		AccountID          string                    `json:"accountId"`
		CategoryID         string                    `json:"categoryId"`
		BudgetID           string                    `json:"budgetId"`
		Description        string                    `json:"description"`
		Amount             float64                   `json:"amount"`
		Note               string                    `json:"note"`
		RecurringFrequency models.RecurringFrequency `json:"recurringFrequency"`
		RecurringDate      int                       `json:"recurringDate"`
	}

	var input CreateRecurringTransactionInput

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": "Invalid request", "error": err.Error()})
	}

	userID, err := uuid.Parse(c.Locals("user_id").(string))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": "Invalid user ID", "error": err.Error()})
	}

	accountID, err := uuid.Parse(input.AccountID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": "Invalid account ID", "error": err.Error()})
	}

	categoryID, err := uuid.Parse(input.CategoryID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": "Invalid category ID", "error": err.Error()})
	}

	var budgetID uuid.NullUUID
	if input.BudgetID != "" {
		parsedBudgetId, err := uuid.Parse(input.BudgetID)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": "Invalid budget ID", "error": err.Error()})
		}
		budgetID = uuid.NullUUID{UUID: parsedBudgetId, Valid: true}
	}

	db := database.DB

	recurringTransaction, err := services.CreateRecurringTransaction(userID, accountID, categoryID, budgetID, input.Description, input.Amount, sql.NullString{String: input.Note, Valid: input.Note != ""}, input.RecurringFrequency, input.RecurringDate, db)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"success": false, "message": "Failed to create recurring transaction", "error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"success": true, "message": "Recurring transaction created successfully", "data": recurringTransaction})
}

// GetRecurringTransactions godoc
// @Summary Get all recurring transactions
// @Description Gets all recurring transactions for the authenticated user.
// @Tags recurring-transactions
// @Security ApiKeyAuth
// @Produce  json
// @Success 200 {object} map[string]interface{} "Recurring transactions retrieved successfully"
// @Router /recurring-transactions [get]
func GetRecurringTransactions(c *fiber.Ctx) error {
	userID, err := uuid.Parse(c.Locals("user_id").(string))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": "Invalid user ID", "error": err.Error()})
	}

	db := database.DB

	recurringTransactions, err := services.GetRecurringTransactions(userID, db)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"success": false, "message": "Failed to get recurring transactions", "error": err.Error()})
	}

	return c.JSON(fiber.Map{"success": true, "message": "Recurring transactions retrieved successfully", "data": recurringTransactions})
}

// UpdateRecurringTransaction godoc
// @Summary Update a recurring transaction
// @Description Updates a recurring transaction for the authenticated user.
// @Tags recurring-transactions
// @Security ApiKeyAuth
// @Accept  json
// @Produce  json
// @Param id path string true "Recurring Transaction ID"
// @Param input body UpdateRecurringTransactionInput true "Update Recurring Transaction Input"
// @Success 200 {object} map[string]interface{} "Recurring transaction updated successfully"
// @Router /recurring-transactions/update/{id} [patch]
func UpdateRecurringTransaction(c *fiber.Ctx) error {
	type UpdateRecurringTransactionInput struct {
		AccountID          string                    `json:"accountId"`
		CategoryID         string                    `json:"categoryId"`
		BudgetID           string                    `json:"budgetId"`
		Description        string                    `json:"description"`
		Amount             float64                   `json:"amount"`
		Note               string                    `json:"note"`
		RecurringFrequency models.RecurringFrequency `json:"recurringFrequency"`
		RecurringDate      int                       `json:"recurringDate"`
	}

	var input UpdateRecurringTransactionInput

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": "Invalid request", "error": err.Error()})
	}

	recurringTransactionID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": "Invalid recurring transaction ID", "error": err.Error()})
	}

	accountID, err := uuid.Parse(input.AccountID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": "Invalid account ID", "error": err.Error()})
	}

	categoryID, err := uuid.Parse(input.CategoryID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": "Invalid category ID", "error": err.Error()})
	}

	var budgetID uuid.NullUUID
	if input.BudgetID != "" {
		parsedBudgetId, err := uuid.Parse(input.BudgetID)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": "Invalid budget ID", "error": err.Error()})
		}
		budgetID = uuid.NullUUID{UUID: parsedBudgetId, Valid: true}
	}

	db := database.DB

	recurringTransaction, err := services.UpdateRecurringTransaction(recurringTransactionID, accountID, categoryID, budgetID, input.Description, input.Amount, sql.NullString{String: input.Note, Valid: input.Note != ""}, input.RecurringFrequency, input.RecurringDate, db)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"success": false, "message": "Failed to update recurring transaction", "error": err.Error()})
	}

	return c.JSON(fiber.Map{"success": true, "message": "Recurring transaction updated successfully", "data": recurringTransaction})
}

// DeleteRecurringTransaction godoc
// @Summary Delete a recurring transaction
// @Description Deletes a recurring transaction for the authenticated user.
// @Tags recurring-transactions
// @Security ApiKeyAuth
// @Produce  json
// @Param id path string true "Recurring Transaction ID"
// @Success 200 {object} map[string]interface{} "Recurring transaction deleted successfully"
// @Router /recurring-transactions/delete/{id} [delete]
func DeleteRecurringTransaction(c *fiber.Ctx) error {
	recurringTransactionID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": "Invalid recurring transaction ID", "error": err.Error()})
	}

	db := database.DB

	if err := services.DeleteRecurringTransaction(recurringTransactionID, db); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"success": false, "message": "Failed to delete recurring transaction", "error": err.Error()})
	}

	return c.JSON(fiber.Map{"success": true, "message": "Recurring transaction deleted successfully"})
}
