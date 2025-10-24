package v1

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/rahulcodepython/finance-tracker-backend/backend/database"
	"github.com/rahulcodepython/finance-tracker-backend/backend/models"
	"github.com/rahulcodepython/finance-tracker-backend/backend/services"
	"github.com/rahulcodepython/finance-tracker-backend/backend/utils"
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
		Name    string  `json:"name"`
		Type    string  `json:"type"`
		Balance float64 `json:"balance"`
	}

	var input CreateAccountInput

	if err := c.BodyParser(&input); err != nil {
		return utils.BadResponse(c, err, "Invalid request")
	}

	userID, err := uuid.Parse(c.Locals("user_id").(string))
	if err != nil {
		return utils.BadResponse(c, err, "Invalid user ID")
	}

	db := database.DB

	account, err := services.CreateAccount(userID, input.Name, models.AccountType(input.Type), input.Balance, db)
	if err != nil {
		return utils.InternalServerError(c, err, "Failed to create account")
	}

	return utils.OKCreatedResponse(c, "Account created successfully", account)
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
		return utils.BadResponse(c, err, "Invalid user ID")
	}

	db := database.DB

	accounts, err := services.GetAccounts(userID, db)
	if err != nil {
		return utils.InternalServerError(c, err, "Failed to get accounts")
	}

	return utils.OKCreatedResponse(c, "Accounts retrieved successfully", accounts)
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
		Name     string `json:"name"`
		Type     string `json:"type"`
		IsActive bool   `json:"isActive"`
	}

	var input UpdateAccountInput

	if err := c.BodyParser(&input); err != nil {
		return utils.BadResponse(c, err, "Invalid request")
	}

	accountID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.BadResponse(c, err, "Invalid account ID")
	}

	db := database.DB

	account, err := services.UpdateAccount(accountID, input.Name, models.AccountType(input.Type), input.IsActive, db)
	if err != nil {
		return utils.InternalServerError(c, err, "Failed to update account")
	}

	return utils.OKResponse(c, "Account updated successfully", account)
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
		return utils.BadResponse(c, err, "Invalid account ID")
	}

	db := database.DB

	if err := services.DeleteAccount(accountID, db); err != nil {
		return utils.InternalServerError(c, err, "Failed to delete account")
	}

	return utils.OKResponse(c, "Account deleted successfully", nil)
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
		return utils.BadResponse(c, err, "Invalid user ID")
	}

	db := database.DB

	totalBalance, err := services.GetTotalBalance(userID, db)
	if err != nil {
		return utils.InternalServerError(c, err, "Failed to get total balance")
	}

	return utils.OKResponse(c, "Total balance retrieved successfully", totalBalance)
}
