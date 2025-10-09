package v1

import (
	"database/sql"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/rahulcodepython/finance-tracker-backend/backend/models"
	"github.com/rahulcodepython/finance-tracker-backend/backend/services"
)

// CreateAccount godoc
// @Summary Create a new financial account
// @Description Creates a new financial account for the authenticated user.
// @Tags accounts
// @Security ApiKeyAuth
// @Accept  json
// @Produce  json
// @Param input body CreateAccountInput true "Create Account Input"
// @Success 201 {object} map[string]interface{} "Account created successfully"
// @Router /accounts/create [post]
func CreateAccount(c *fiber.Ctx) error {
	type CreateAccountInput struct {
		Name string `json:"name"`
		Type string `json:"type"`
	}

	var input CreateAccountInput

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": "Invalid request", "error": err.Error()})
	}

	userID, err := uuid.Parse(c.Locals("user_id").(string))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": "Invalid user ID", "error": err.Error()})
	}

	db := c.Locals("db").(*sql.DB)

	account, err := services.CreateAccount(userID, input.Name, models.AccountType(input.Type), db)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"success": false, "message": "Failed to create account", "error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"success": true, "message": "Account created successfully", "data": account})
}

// GetAccounts godoc
// @Summary Get all financial accounts
// @Description Gets all financial accounts for the authenticated user.
// @Tags accounts
// @Security ApiKeyAuth
// @Produce  json
// @Success 200 {object} map[string]interface{} "Accounts retrieved successfully"
// @Router /accounts [get]
func GetAccounts(c *fiber.Ctx) error {
	userID, err := uuid.Parse(c.Locals("user_id").(string))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": "Invalid user ID", "error": err.Error()})
	}

	db := c.Locals("db").(*sql.DB)

	accounts, err := services.GetAccounts(userID, db)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"success": false, "message": "Failed to get accounts", "error": err.Error()})
	}

	return c.JSON(fiber.Map{"success": true, "message": "Accounts retrieved successfully", "data": accounts})
}

// UpdateAccount godoc
// @Summary Update a financial account
// @Description Updates a financial account for the authenticated user.
// @Tags accounts
// @Security ApiKeyAuth
// @Accept  json
// @Produce  json
// @Param id path string true "Account ID"
// @Param input body UpdateAccountInput true "Update Account Input"
// @Success 200 {object} map[string]interface{} "Account updated successfully"
// @Router /accounts/update/{id} [patch]
func UpdateAccount(c *fiber.Ctx) error {
	type UpdateAccountInput struct {
		Name string `json:"name"`
		Type string `json:"type"`
	}

	var input UpdateAccountInput

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": "Invalid request", "error": err.Error()})
	}

	accountID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": "Invalid account ID", "error": err.Error()})
	}

	db := c.Locals("db").(*sql.DB)

	account, err := services.UpdateAccount(accountID, input.Name, models.AccountType(input.Type), db)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"success": false, "message": "Failed to update account", "error": err.Error()})
	}

	return c.JSON(fiber.Map{"success": true, "message": "Account updated successfully", "data": account})
}

// DeleteAccount godoc
// @Summary Delete a financial account
// @Description Deletes a financial account for the authenticated user.
// @Tags accounts
// @Security ApiKeyAuth
// @Produce  json
// @Param id path string true "Account ID"
// @Success 200 {object} map[string]interface{} "Account deleted successfully"
// @Router /accounts/delete/{id} [delete]
func DeleteAccount(c *fiber.Ctx) error {
	accountID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": "Invalid account ID", "error": err.Error()})
	}

	db := c.Locals("db").(*sql.DB)

	if err := services.DeleteAccount(accountID, db); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"success": false, "message": "Failed to delete account", "error": err.Error()})
	}

	return c.JSON(fiber.Map{"success": true, "message": "Account deleted successfully"})
}

// GetTotalBalance godoc
// @Summary Get the total balance of all active accounts
// @Description Gets the total balance of all active accounts for the authenticated user.
// @Tags accounts
// @Security ApiKeyAuth
// @Produce  json
// @Success 200 {object} map[string]interface{} "Total balance retrieved successfully"
// @Router /accounts/total-balance [get]
func GetTotalBalance(c *fiber.Ctx) error {
	userID, err := uuid.Parse(c.Locals("user_id").(string))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": "Invalid user ID", "error": err.Error()})
	}

	db := c.Locals("db").(*sql.DB)

	totalBalance, err := services.GetTotalBalance(userID, db)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"success": false, "message": "Failed to get total balance", "error": err.Error()})
	}

	return c.JSON(fiber.Map{"success": true, "message": "Total balance retrieved successfully", "data": totalBalance})
}
