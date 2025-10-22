package v1

import (
	"database/sql"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/rahulcodepython/finance-tracker-backend/backend/database"
	"github.com/rahulcodepython/finance-tracker-backend/backend/models"
	"github.com/rahulcodepython/finance-tracker-backend/backend/services"
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
		Description string  `json:"description"`
		Amount      float64 `json:"amount"`
		Type        string  `json:"type"`
		Date        string  `json:"date"`
		Note        string  `json:"note"`
	}

	var input CreateTransactionInput

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

	var categoryID uuid.NullUUID
	if input.CategoryID != "" {
		parsedCategoryID, err := uuid.Parse(input.CategoryID)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": "Invalid category ID", "error": err.Error()})
		}
		categoryID = uuid.NullUUID{UUID: parsedCategoryID, Valid: true}
	}

	transactionDate, err := time.Parse("2006-01-02", input.Date)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": "Invalid date format", "error": err.Error()})
	}

	db := database.DB

	transaction, err := services.CreateTransaction(userID, accountID, categoryID, input.Description, input.Amount, models.TransactionType(input.Type), transactionDate, sql.NullString{String: input.Note, Valid: input.Note != ""}, db)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"success": false, "message": "Failed to create transaction", "error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"success": true, "message": "Transaction created successfully", "data": transaction})
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
// @Success 200 {object} map[string]interface{} "Transactions retrieved successfully"
// @Router /transactions [get]
func GetTransactions(c *fiber.Ctx) error {
	userID, err := uuid.Parse(c.Locals("user_id").(string))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": "Invalid user ID", "error": err.Error()})
	}
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	description := c.Query("description")
	categoryID := c.Query("category")
	accountID := c.Query("account")
	startDate := c.Query("startDate")
	endDate := c.Query("endDate")

	db := database.DB

	transactions, err := services.GetTransactions(userID, page, limit, description, categoryID, accountID, startDate, endDate, db)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"success": false, "message": "Failed to get transactions", "error": err.Error()})
	}

	return c.JSON(fiber.Map{"success": true, "message": "Transactions retrieved successfully", "data": transactions})
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
		Description string  `json:"description"`
		Amount      float64 `json:"amount"`
		Type        string  `json:"type"`
		Date        string  `json:"date"`
		Note        string  `json:"note"`
	}

	var input UpdateTransactionInput

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": "Invalid request", "error": err.Error()})
	}

	transactionID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": "Invalid transaction ID", "error": err.Error()})
	}

	accountID, err := uuid.Parse(input.AccountID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": "Invalid account ID", "error": err.Error()})
	}

	var categoryID uuid.NullUUID
	if input.CategoryID != "" {
		parsedCategoryID, err := uuid.Parse(input.CategoryID)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": "Invalid category ID", "error": err.Error()})
		}
		categoryID = uuid.NullUUID{UUID: parsedCategoryID, Valid: true}
	}

	transactionDate, err := time.Parse("2006-01-02", input.Date)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": "Invalid date format", "error": err.Error()})
	}

	db := database.DB

	transaction, err := services.UpdateTransaction(transactionID, accountID, categoryID, input.Description, input.Amount, models.TransactionType(input.Type), transactionDate, sql.NullString{String: input.Note, Valid: input.Note != ""}, db)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"success": false, "message": "Failed to update transaction", "error": err.Error()})
	}

	return c.JSON(fiber.Map{"success": true, "message": "Transaction updated successfully", "data": transaction})
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
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": "Invalid transaction ID", "error": err.Error()})
	}

	db := database.DB

	if err := services.DeleteTransaction(transactionID, db); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"success": false, "message": "Failed to delete transaction", "error": err.Error()})
	}

	return c.JSON(fiber.Map{"success": true, "message": "Transaction deleted successfully"})
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
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": "Invalid user ID", "error": err.Error()})
	}
	startDate := c.Query("startDate")
	endDate := c.Query("endDate")

	db := database.DB

	data, err := services.GetAggregateData(userID, startDate, endDate, db)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"success": false, "message": "Failed to get aggregate data", "error": err.Error()})
	}

	return c.JSON(fiber.Map{"success": true, "message": "Aggregate data retrieved successfully", "data": data})
}
