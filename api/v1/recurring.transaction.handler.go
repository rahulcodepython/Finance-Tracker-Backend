package v1

import (
	"database/sql"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/rahulcodepython/finance-tracker-backend/backend/database"
	"github.com/rahulcodepython/finance-tracker-backend/backend/models"
	"github.com/rahulcodepython/finance-tracker-backend/backend/services"
	"github.com/rahulcodepython/finance-tracker-backend/backend/utils"
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
		return utils.BadResponse(c, err, "Invalid request")
	}

	userID, err := uuid.Parse(c.Locals("user_id").(string))
	if err != nil {
		return utils.BadResponse(c, err, "Invalid user ID")
	}

	accountID, err := uuid.Parse(input.AccountID)
	if err != nil {
		return utils.BadResponse(c, err, "Invalid account ID")
	}

	categoryID, err := uuid.Parse(input.CategoryID)
	if err != nil {
		return utils.BadResponse(c, err, "Invalid category ID")
	}

	var budgetID uuid.NullUUID
	if input.BudgetID != "" {
		parsedBudgetId, err := uuid.Parse(input.BudgetID)
		if err != nil {
			return utils.BadResponse(c, err, "Invalid budget ID")
		}
		budgetID = uuid.NullUUID{UUID: parsedBudgetId, Valid: true}
	}

	db := database.DB

	recurringTransaction, err := services.CreateRecurringTransaction(userID, accountID, categoryID, budgetID, input.Description, input.Amount, sql.NullString{String: input.Note, Valid: input.Note != ""}, input.RecurringFrequency, input.RecurringDate, db)
	if err != nil {
		return utils.InternalServerError(c, err, "Failed to create recurring transaction")
	}

	return utils.OKCreatedResponse(c, "Recurring transaction created successfully", recurringTransaction)
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
		return utils.BadResponse(c, err, "Invalid user ID")
	}

	db := database.DB

	recurringTransactions, err := services.GetRecurringTransactions(userID, db)
	if err != nil {
		return utils.InternalServerError(c, err, "Failed to get recurring transactions")
	}

	return utils.OKResponse(c, "Recurring transactions retrieved successfully", recurringTransactions)
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
		return utils.BadResponse(c, err, "Invalid request")
	}

	recurringTransactionID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.BadResponse(c, err, "Invalid recurring transaction ID")
	}

	accountID, err := uuid.Parse(input.AccountID)
	if err != nil {
		return utils.BadResponse(c, err, "Invalid account ID")
	}

	categoryID, err := uuid.Parse(input.CategoryID)
	if err != nil {
		return utils.BadResponse(c, err, "Invalid category ID")
	}

	var budgetID uuid.NullUUID
	if input.BudgetID != "" {
		parsedBudgetId, err := uuid.Parse(input.BudgetID)
		if err != nil {
			return utils.BadResponse(c, err, "Invalid budget ID")
		}
		budgetID = uuid.NullUUID{UUID: parsedBudgetId, Valid: true}
	}

	db := database.DB

	recurringTransaction, err := services.UpdateRecurringTransaction(recurringTransactionID, accountID, categoryID, budgetID, input.Description, input.Amount, sql.NullString{String: input.Note, Valid: input.Note != ""}, input.RecurringFrequency, input.RecurringDate, db)
	if err != nil {
		return utils.InternalServerError(c, err, "Failed to update recurring transaction")
	}

	return utils.OKResponse(c, "Recurring transaction updated successfully", recurringTransaction)
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
		return utils.BadResponse(c, err, "Invalid recurring transaction ID")
	}

	db := database.DB

	if err := services.DeleteRecurringTransaction(recurringTransactionID, db); err != nil {
		return utils.InternalServerError(c, err, "Failed to delete recurring transaction")
	}

	return utils.OKResponse(c, "Recurring transaction deleted successfully", nil)
}
