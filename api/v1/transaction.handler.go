package v1

import (
	"database/sql"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/rahulcodepython/finance-tracker-backend/backend/database"
	"github.com/rahulcodepython/finance-tracker-backend/backend/services"
	"github.com/rahulcodepython/finance-tracker-backend/backend/utils"
)

// CreateTransaction godoc
// @Summary Create a new transaction
// @Description Creates a new transaction for the authenticated user.
// @Tags transactions
// @Security ApiKeyAuth
// @Accept  json
// @Produce  json
// @Param input body CreateTransactionInput true "Create Transaction Input"
// @Success 201 {object} map[string]interface{} "Transaction created successfully"
// @Router /transactions/create [post]
func CreateTransaction(c *fiber.Ctx) error {
	type CreateTransactionInput struct {
		AccountID   string  `json:"accountId"`
		CategoryID  string  `json:"categoryId"`
		BudgetID    string  `json:"budgetId"`
		Description string  `json:"description"`
		Amount      float64 `json:"amount"`
		Date        string  `json:"date"`
		Note        string  `json:"note"`
	}

	var input CreateTransactionInput

	db := database.DB

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

	transactionDate, err := time.Parse("2006-01-02", input.Date)
	if err != nil {
		return utils.BadResponse(c, err, "Invalid date format")
	}

	transaction, err := services.CreateTransaction(userID, accountID, categoryID, budgetID, input.Description, input.Amount, transactionDate, sql.NullString{String: input.Note, Valid: input.Note != ""}, db)
	if err != nil {
		return utils.InternalServerError(c, err, "Failed to create transaction")
	}

	return utils.OKCreatedResponse(c, "Transaction created successfully", transaction)
}

// GetTransactions godoc
// @Summary Get all transactions
// @Description Gets all transactions for the authenticated user, with pagination and filtering.
// @Tags transactions
// @Security ApiKeyAuth
// @Produce  json
// @Param page query int false "Page number"
// @Param limit query int false "Number of items per page"
// @Param description query string false "Filter by description"
// @Param category query string false "Filter by category ID"
// @Param account query string false "Filter by account ID"
// @Param startDate query string false "Filter by start date (YYYY-MM-DD)"
// @Param endDate query string false "Filter by end date (YYYY-MM-DD)"
// @Param budget query string false "Filter by budget ID"
// @Success 200 {object} map[string]interface{} "Transactions retrieved successfully"
// @Router /transactions [get]
func GetTransactions(c *fiber.Ctx) error {
	userID, err := uuid.Parse(c.Locals("user_id").(string))
	if err != nil {
		return utils.BadResponse(c, err, "Invalid user ID")
	}
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	description := c.Query("description")
	categoryID := c.Query("category")
	accountID := c.Query("account")
	budgetID := c.Query("budget")
	startDate := c.Query("startDate")
	endDate := c.Query("endDate")

	db := database.DB

	transactions, err := services.GetTransactions(userID, page, limit, description, categoryID, accountID, budgetID, startDate, endDate, db)
	if err != nil {
		return utils.InternalServerError(c, err, "Failed to get transactions")
	}

	return utils.OKResponse(c, "Transactions retrieved successfully", transactions)
}

// UpdateTransaction godoc
// @Summary Update a transaction
// @Description Updates a transaction for the authenticated user.
// @Tags transactions
// @Security ApiKeyAuth
// @Accept  json
// @Produce  json
// @Param id path string true "Transaction ID"
// @Param input body UpdateTransactionInput true "Update Transaction Input"
// @Success 200 {object} map[string]interface{} "Transaction updated successfully"
// @Router /transactions/update/{id} [patch]
func UpdateTransaction(c *fiber.Ctx) error {
	type UpdateTransactionInput struct {
		AccountID   string  `json:"accountId"`
		CategoryID  string  `json:"categoryId"`
		BudgetID    string  `json:"budgetId"`
		Description string  `json:"description"`
		Amount      float64 `json:"amount"`
		Date        string  `json:"date"`
		Note        string  `json:"note"`
	}

	var input UpdateTransactionInput

	db := database.DB

	if err := c.BodyParser(&input); err != nil {
		return utils.BadResponse(c, err, "Invalid request")
	}

	transactionID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.BadResponse(c, err, "Invalid transaction ID")
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

	transactionDate, err := time.Parse("2006-01-02", input.Date)
	if err != nil {
		return utils.BadResponse(c, err, "Invalid date format")
	}

	transaction, err := services.UpdateTransaction(transactionID, accountID, categoryID, budgetID, input.Description, input.Amount, transactionDate, sql.NullString{String: input.Note, Valid: input.Note != ""}, db)
	if err != nil {
		return utils.InternalServerError(c, err, "Failed to update transaction")
	}

	return utils.OKResponse(c, "Transaction updated successfully", transaction)
}

// DeleteTransaction godoc
// @Summary Delete a transaction
// @Description Deletes a transaction for the authenticated user.
// @Tags transactions
// @Security ApiKeyAuth
// @Produce  json
// @Param id path string true "Transaction ID"
// @Success 200 {object} map[string]interface{} "Transaction deleted successfully"
// @Router /transactions/delete/{id} [delete]
func DeleteTransaction(c *fiber.Ctx) error {
	transactionID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.BadResponse(c, err, "Invalid transaction ID")
	}

	db := database.DB

	if err := services.DeleteTransaction(transactionID, db); err != nil {
		return utils.InternalServerError(c, err, "Failed to delete transaction")
	}

	return utils.OKResponse(c, "Transaction deleted successfully", nil)
}

// GetAggregateData godoc
// @Summary Get aggregate data for transactions
// @Description Gets aggregate data for transactions (total income, total expenses, net income) over a specified period.
// @Tags transactions
// @Security ApiKeyAuth
// @Produce  json
// @Param startDate query string false "Start date (YYYY-MM-DD)"
// @Param endDate query string false "End date (YYYY-MM-DD)"
// @Success 200 {object} map[string]interface{} "Aggregate data retrieved successfully"
// @Router /transactions/aggregate [get]
func GetAggregateData(c *fiber.Ctx) error {
	userID, err := uuid.Parse(c.Locals("user_id").(string))
	if err != nil {
		return utils.BadResponse(c, err, "Invalid user ID")
	}
	startDate := c.Query("startDate")
	endDate := c.Query("endDate")

	db := database.DB

	data, err := services.GetAggregateData(userID, startDate, endDate, db)
	if err != nil {
		return utils.InternalServerError(c, err, "Failed to get aggregate data")
	}

	return utils.OKResponse(c, "Aggregate data retrieved successfully", data)
}
