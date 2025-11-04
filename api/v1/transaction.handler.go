// Package v1 provides handlers for version 1 of the API.
package v1

import (
	"database/sql"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/rahulcodepython/finance-tracker-backend/backend/database"
	"github.com/rahulcodepython/finance-tracker-backend/backend/serializers"
	"github.com/rahulcodepython/finance-tracker-backend/backend/services"
	"github.com/rahulcodepython/finance-tracker-backend/backend/utils"
)

// CreateTransaction handles the creation of a new transaction.
// @Summary      Create a new transaction
// @Description  Creates a new transaction for the authenticated user.
// @Tags         Transactions
// @Accept       json
// @Produce      json
// @Param        input body      CreateTransactionInput true "Create Transaction Input"
// @Success      201  {object}  utils.Response
// @Failure      400  {object}  utils.Response
// @Failure      500  {object}  utils.Response
// @Router       /transactions/create [post]
func CreateTransaction(c *fiber.Ctx) error {
	// Create a new instance of the CreateTransactionInput struct.
	var input serializers.TransactionInput

	// Get the database connection.
	db := database.DB

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

	// Parse the transaction date from the input.
	transactionDate, err := time.Parse("2006-01-02", input.Date)
	if err != nil {
		return utils.BadResponse(c, err, "Invalid date format")
	}

	// Call the CreateTransaction service to create the new transaction.
	transaction, err := services.CreateTransaction(userID, accountID, categoryID, budgetID, input.Description, input.Amount, transactionDate, sql.NullString{String: input.Note, Valid: input.Note != ""}, db)
	if err != nil {
		return utils.InternalServerError(c, err, "Failed to create transaction")
	}

	// Return a 201 Created response with the new transaction.
	return utils.OKCreatedResponse(c, "Transaction created successfully", transaction)
}

// GetTransactions handles the retrieval of transactions with filtering and pagination.
// @Summary      Get all transactions
// @Description  Retrieves a paginated and filtered list of transactions for the authenticated user.
// @Tags         Transactions
// @Produce      json
// @Param        page         query     int    false "Page number for pagination" default(1)
// @Param        limit        query     int    false "Number of transactions per page" default(10)
// @Param        description  query     string false "Filter by transaction description"
// @Param        category     query     string false "Filter by category ID"
// @Param        account      query     string false "Filter by account ID"
// @Param        budget       query     string false "Filter by budget ID"
// @Param        startDate    query     string false "Start date for transaction filter (YYYY-MM-DD)"
// @Param        endDate      query     string false "End date for transaction filter (YYYY-MM-DD)"
// @Success      200          {object}  utils.Response
// @Failure      400          {object}  utils.Response
// @Failure      500          {object}  utils.Response
// @Router       /transactions [get]
func GetTransactions(c *fiber.Ctx) error {
	// Get the user ID from the context.
	userID, err := uuid.Parse(c.Locals("user_id").(string))
	if err != nil {
		return utils.BadResponse(c, err, "Invalid user ID")
	}
	// Get pagination and filter parameters from the query string.
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	description := c.Query("description")
	categoryID := c.Query("category")
	accountID := c.Query("account")
	budgetID := c.Query("budget")
	startDate := c.Query("startDate")
	endDate := c.Query("endDate")

	// Get the database connection.
	db := database.DB

	// Call the GetTransactions service to retrieve the transactions.
	transactions, err := services.GetTransactions(userID, page, limit, description, categoryID, accountID, budgetID, startDate, endDate, db)
	if err != nil {
		return utils.InternalServerError(c, err, "Failed to retrieve transactions")
	}

	// Return a 200 OK response with the retrieved transactions.
	return utils.OKResponse(c, "Transactions retrieved successfully", transactions)
}

// UpdateTransaction handles the update of a specific transaction.
// @Summary      Update a transaction
// @Description  Updates the details of a specific transaction for the authenticated user.
// @Tags         Transactions
// @Accept       json
// @Produce      json
// @Param        id   path      string                 true "Transaction ID"
// @Param        input body      UpdateTransactionInput true "Update Transaction Input"
// @Success      200  {object}  utils.Response
// @Failure      400  {object}  utils.Response
// @Failure      500  {object}  utils.Response
// @Router       /transactions/update/{id} [patch]
func UpdateTransaction(c *fiber.Ctx) error {
	// Create a new instance of the UpdateTransactionInput struct.
	var input serializers.TransactionInput

	// Get the database connection.
	db := database.DB

	// Parse the request body into the input struct.
	if err := c.BodyParser(&input); err != nil {
		return utils.BadResponse(c, err, "Invalid request payload")
	}

	// Get the transaction ID from the URL parameters.
	transactionID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.BadResponse(c, err, "Invalid transaction ID")
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

	// Parse the transaction date from the input.
	transactionDate, err := time.Parse("2006-01-02", input.Date)
	if err != nil {
		return utils.BadResponse(c, err, "Invalid date format")
	}

	// Call the UpdateTransaction service to update the transaction.
	transaction, err := services.UpdateTransaction(transactionID, accountID, categoryID, budgetID, input.Description, input.Amount, transactionDate, sql.NullString{String: input.Note, Valid: input.Note != ""}, db)
	if err != nil {
		return utils.InternalServerError(c, err, "Failed to update transaction")
	}

	// Return a 200 OK response with the updated transaction.
	return utils.OKResponse(c, "Transaction updated successfully", transaction)
}

// DeleteTransaction handles the deletion of a specific transaction.
// @Summary      Delete a transaction
// @Description  Deletes a specific transaction for the authenticated user.
// @Tags         Transactions
// @Produce      json
// @Param        id   path      string  true "Transaction ID"
// @Success      200  {object}  utils.Response
// @Failure      400  {object}  utils.Response
// @Failure      500  {object}  utils.Response
// @Router       /transactions/delete/{id} [delete]
func DeleteTransaction(c *fiber.Ctx) error {
	// Get the transaction ID from the URL parameters.
	transactionID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.BadResponse(c, err, "Invalid transaction ID")
	}

	// Get the database connection.
	db := database.DB

	// Call the DeleteTransaction service to delete the transaction.
	if err := services.DeleteTransaction(transactionID, db); err != nil {
		return utils.InternalServerError(c, err, "Failed to delete transaction")
	}

	// Return a 200 OK response.
	return utils.OKResponse(c, "Transaction deleted successfully", nil)
}

// GetAggregateData handles the retrieval of aggregate transaction data.
// @Summary      Get aggregate transaction data
// @Description  Retrieves aggregate data for transactions, such as total income, expenses, and net income, over a specified period.
// @Tags         Transactions
// @Produce      json
// @Param        startDate query     string false "Start date for aggregation (YYYY-MM-DD)"
// @Param        endDate   query     string false "End date for aggregation (YYYY-MM-DD)"
// @Success      200       {object}  utils.Response
// @Failure      400       {object}  utils.Response
// @Failure      500       {object}  utils.Response
// @Router       /transactions/aggregate [get]
func GetAggregateData(c *fiber.Ctx) error {
	// Get the user ID from the context.
	userID, err := uuid.Parse(c.Locals("user_id").(string))
	if err != nil {
		return utils.BadResponse(c, err, "Invalid user ID")
	}
	// Get the start and end dates from the query parameters.
	startDate := c.Query("startDate")
	endDate := c.Query("endDate")

	// Get the database connection.
	db := database.DB

	// Call the GetAggregateData service to retrieve the aggregate data.
	data, err := services.GetAggregateData(userID, startDate, endDate, db)
	if err != nil {
		return utils.InternalServerError(c, err, "Failed to get aggregate data")
	}

	// Return a 200 OK response with the aggregate data.
	return utils.OKResponse(c, "Aggregate data retrieved successfully", data)
}
