// Package v1 provides handlers for version 1 of the API.
package v1

import (
	"database/sql"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/rahulcodepython/finance-tracker-backend/backend/database"
	"github.com/rahulcodepython/finance-tracker-backend/backend/serializers"
	"github.com/rahulcodepython/finance-tracker-backend/backend/services"
	"github.com/rahulcodepython/finance-tracker-backend/backend/utils"
)

// CreateRecurringTransaction handles the creation of a new recurring transaction.
// @Summary      Create a new recurring transaction
// @Description  Creates a new recurring transaction for the authenticated user.
// @Tags         Recurring Transactions
// @Accept       json
// @Produce      json
// @Param        input body      RecurringTransactionInput true "Create Recurring Transaction Input"
// @Success      201  {object}  utils.Response
// @Failure      400  {object}  utils.Response
// @Failure      500  {object}  utils.Response
// @Router       /recurring-transactions/create [post]
func CreateRecurringTransaction(c *fiber.Ctx) error {
	// Create a new instance of the RecurringTransactionInput struct.
	var input serializers.RecurringTransactionInput

	// Parse the request body into the input struct.
	if err := c.BodyParser(&input); err != nil {
		return utils.BadResponse(c, err, "Invalid request payload")
	}

	// Get the user ID from the context.
	userID, err := uuid.Parse(c.Locals("user_id").(string))
	if err != nil {
		return utils.BadResponse(c, err, "Invalid user ID")
	}

	// Parse the account ID from the input.
	accountID, err := uuid.Parse(input.AccountID)
	if err != nil {
		return utils.BadResponse(c, err, "Invalid account ID")
	}

	// Parse the category ID from the input.
	categoryID, err := uuid.Parse(input.CategoryID)
	if err != nil {
		return utils.BadResponse(c, err, "Invalid category ID")
	}

	// Parse the budget ID from the input, if provided.
	var budgetID uuid.NullUUID
	if input.BudgetID != "" {
		parsedBudgetId, err := uuid.Parse(input.BudgetID)
		if err != nil {
			return utils.BadResponse(c, err, "Invalid budget ID")
		}
		budgetID = uuid.NullUUID{UUID: parsedBudgetId, Valid: true}
	}

	// Get the database connection.
	db := database.DB

	// Call the CreateRecurringTransaction service to create the new recurring transaction.
	recurringTransaction, err := services.CreateRecurringTransaction(userID, accountID, categoryID, budgetID, input.Description, input.Amount, sql.NullString{String: input.Note, Valid: input.Note != ""}, input.RecurringFrequency, input.RecurringDate, db)
	if err != nil {
		return utils.InternalServerError(c, err, "Failed to create recurring transaction")
	}

	// Return a 201 Created response with the new recurring transaction.
	return utils.OKCreatedResponse(c, "Recurring transaction created successfully", recurringTransaction)
}

// GetRecurringTransactions handles the retrieval of all recurring transactions for the authenticated user.
// @Summary      Get all recurring transactions
// @Description  Retrieves all recurring transactions associated with the authenticated user.
// @Tags         Recurring Transactions
// @Produce      json
// @Success      200  {object}  utils.Response
// @Failure      400  {object}  utils.Response
// @Failure      500  {object}  utils.Response
// @Router       /recurring-transactions [get]
func GetRecurringTransactions(c *fiber.Ctx) error {
	// Get the user ID from the context.
	userID, err := uuid.Parse(c.Locals("user_id").(string))
	if err != nil {
		return utils.BadResponse(c, err, "Invalid user ID")
	}

	// Get the database connection.
	db := database.DB

	// Call the GetRecurringTransactions service to retrieve the user's recurring transactions.
	recurringTransactions, err := services.GetRecurringTransactions(userID, db)
	if err != nil {
		return utils.InternalServerError(c, err, "Failed to retrieve recurring transactions")
	}

	// Return a 200 OK response with the retrieved recurring transactions.
	return utils.OKResponse(c, "Recurring transactions retrieved successfully", recurringTransactions)
}

// UpdateRecurringTransaction handles the update of a specific recurring transaction.
// @Summary      Update a recurring transaction
// @Description  Updates the details of a specific recurring transaction for the authenticated user.
// @Tags         Recurring Transactions
// @Accept       json
// @Produce      json
// @Param        id   path      string                        true "Recurring Transaction ID"
// @Param        input body      RecurringTransactionInput true "Update Recurring Transaction Input"
// @Success      200  {object}  utils.Response
// @Failure      400  {object}  utils.Response
// @Failure      500  {object}  utils.Response
// @Router       /recurring-transactions/update/{id} [patch]
func UpdateRecurringTransaction(c *fiber.Ctx) error {
	// Create a new instance of the RecurringTransactionInput struct.
	var input serializers.RecurringTransactionInput

	// Parse the request body into the input struct.
	if err := c.BodyParser(&input); err != nil {
		return utils.BadResponse(c, err, "Invalid request payload")
	}

	// Get the recurring transaction ID from the URL parameters.
	recurringTransactionID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.BadResponse(c, err, "Invalid recurring transaction ID")
	}

	// Parse the account ID from the input.
	accountID, err := uuid.Parse(input.AccountID)
	if err != nil {
		return utils.BadResponse(c, err, "Invalid account ID")
	}

	// Parse the category ID from the input.
	categoryID, err := uuid.Parse(input.CategoryID)
	if err != nil {
		return utils.BadResponse(c, err, "Invalid category ID")
	}

	// Parse the budget ID from the input, if provided.
	var budgetID uuid.NullUUID
	if input.BudgetID != "" {
		parsedBudgetId, err := uuid.Parse(input.BudgetID)
		if err != nil {
			return utils.BadResponse(c, err, "Invalid budget ID")
		}
		budgetID = uuid.NullUUID{UUID: parsedBudgetId, Valid: true}
	}

	// Get the database connection.
	db := database.DB

	// Call the UpdateRecurringTransaction service to update the recurring transaction.
	recurringTransaction, err := services.UpdateRecurringTransaction(recurringTransactionID, accountID, categoryID, budgetID, input.Description, input.Amount, sql.NullString{String: input.Note, Valid: input.Note != ""}, input.RecurringFrequency, input.RecurringDate, db)
	if err != nil {
		return utils.InternalServerError(c, err, "Failed to update recurring transaction")
	}

	// Return a 200 OK response with the updated recurring transaction.
	return utils.OKResponse(c, "Recurring transaction updated successfully", recurringTransaction)
}

// DeleteRecurringTransaction handles the deletion of a specific recurring transaction.
// @Summary      Delete a recurring transaction
// @Description  Deletes a specific recurring transaction for the authenticated user.
// @Tags         Recurring Transactions
// @Produce      json
// @Param        id   path      string  true "Recurring Transaction ID"
// @Success      200  {object}  utils.Response
// @Failure      400  {object}  utils.Response
// @Failure      500  {object}  utils.Response
// @Router       /recurring-transactions/delete/{id} [delete]
func DeleteRecurringTransaction(c *fiber.Ctx) error {
	// Get the recurring transaction ID from the URL parameters.
	recurringTransactionID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.BadResponse(c, err, "Invalid recurring transaction ID")
	}

	// Get the database connection.
	db := database.DB

	// Call the DeleteRecurringTransaction service to delete the recurring transaction.
	if err := services.DeleteRecurringTransaction(recurringTransactionID, db); err != nil {
		return utils.InternalServerError(c, err, "Failed to delete recurring transaction")
	}

	// Return a 200 OK response.
	return utils.OKResponse(c, "Recurring transaction deleted successfully", nil)
}
