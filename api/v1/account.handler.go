// Package v1 provides handlers for version 1 of the API.
package v1

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/rahulcodepython/finance-tracker-backend/backend/database"
	"github.com/rahulcodepython/finance-tracker-backend/backend/models"
	"github.com/rahulcodepython/finance-tracker-backend/backend/serializers"
	"github.com/rahulcodepython/finance-tracker-backend/backend/services"
	"github.com/rahulcodepython/finance-tracker-backend/backend/utils"
)

// CreateAccount handles the creation of a new financial account.
// @Summary      Create a new financial account
// @Description  Creates a new financial account for the authenticated user.
// @Tags         Accounts
// @Accept       json
// @Produce      json
// @Param        input body      CreateAccountInput true "Create Account Input"
// @Success      201  {object}  utils.Response
// @Failure      400  {object}  utils.Response
// @Failure      500  {object}  utils.Response
// @Router       /accounts/create [post]
func CreateAccount(c *fiber.Ctx) error {
	// Create a new instance of the CreateAccountInput struct to hold the request body.
	var input serializers.CreateAccountInput

	// Parse the request body into the input struct.
	if err := c.BodyParser(&input); err != nil {
		// If parsing fails, return a 400 Bad Request response.
		return utils.BadResponse(c, err, "Invalid request payload")
	}

	// Get the user ID from the context, which was set by the authentication middleware.
	userID, err := uuid.Parse(c.Locals("user_id").(string))
	if err != nil {
		// If the user ID is invalid, return a 400 Bad Request response.
		return utils.BadResponse(c, err, "Invalid user ID")
	}

	// Get the database connection from the global variable.
	db := database.DB

	// Call the CreateAccount service to create the new account.
	account, err := services.CreateAccount(userID, input.Name, models.AccountType(input.Type), input.Balance, db)
	if err != nil {
		// If account creation fails, return a 500 Internal Server Error response.
		return utils.InternalServerError(c, err, "Failed to create account")
	}

	// Return a 201 Created response with the newly created account.
	return utils.OKCreatedResponse(c, "Account created successfully", account)
}

// GetAccounts handles the retrieval of all financial accounts for the authenticated user.
// @Summary      Get all financial accounts
// @Description  Retrieves all financial accounts associated with the authenticated user.
// @Tags         Accounts
// @Produce      json
// @Success      200  {object}  utils.Response
// @Failure      400  {object}  utils.Response
// @Failure      500  {object}  utils.Response
// @Router       /accounts [get]
func GetAccounts(c *fiber.Ctx) error {
	// Get the user ID from the context.
	userID, err := uuid.Parse(c.Locals("user_id").(string))
	if err != nil {
		// If the user ID is invalid, return a 400 Bad Request response.
		return utils.BadResponse(c, err, "Invalid user ID")
	}

	// Get the database connection.
	db := database.DB

	// Call the GetAccounts service to retrieve the user's accounts.
	accounts, err := services.GetAccounts(userID, db)
	if err != nil {
		// If retrieval fails, return a 500 Internal Server Error response.
		return utils.InternalServerError(c, err, "Failed to retrieve accounts")
	}

	// Return a 200 OK response with the retrieved accounts.
	return utils.OKResponse(c, "Accounts retrieved successfully", accounts)
}

// UpdateAccount handles the update of a specific financial account.
// @Summary      Update a financial account
// @Description  Updates the details of a specific financial account for the authenticated user.
// @Tags         Accounts
// @Accept       json
// @Produce      json
// @Param        id   path      string          true "Account ID"
// @Param        input body      UpdateAccountInput true "Update Account Input"
// @Success      200  {object}  utils.Response
// @Failure      400  {object}  utils.Response
// @Failure      500  {object}  utils.Response
// @Router       /accounts/update/{id} [patch]
func UpdateAccount(c *fiber.Ctx) error {
	// Create a new instance of the UpdateAccountInput struct.
	var input serializers.UpdateAccountInput

	// Parse the request body into the input struct.
	if err := c.BodyParser(&input); err != nil {
		// If parsing fails, return a 400 Bad Request response.
		return utils.BadResponse(c, err, "Invalid request payload")
	}

	// Get the account ID from the URL parameters.
	accountID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		// If the account ID is invalid, return a 400 Bad Request response.
		return utils.BadResponse(c, err, "Invalid account ID")
	}

	// Get the database connection.
	db := database.DB

	// Call the UpdateAccount service to update the account.
	account, err := services.UpdateAccount(accountID, input.Name, models.AccountType(input.Type), input.IsActive, db)
	if err != nil {
		// If the update fails, return a 500 Internal Server Error response.
		return utils.InternalServerError(c, err, "Failed to update account")
	}

	// Return a 200 OK response with the updated account.
	return utils.OKResponse(c, "Account updated successfully", account)
}

// DeleteAccount handles the deletion of a specific financial account.
// @Summary      Delete a financial account
// @Description  Deletes a specific financial account for the authenticated user.
// @Tags         Accounts
// @Produce      json
// @Param        id   path      string  true "Account ID"
// @Success      200  {object}  utils.Response
// @Failure      400  {object}  utils.Response
// @Failure      500  {object}  utils.Response
// @Router       /accounts/delete/{id} [delete]
func DeleteAccount(c *fiber.Ctx) error {
	// Get the account ID from the URL parameters.
	accountID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		// If the account ID is invalid, return a 400 Bad Request response.
		return utils.BadResponse(c, err, "Invalid account ID")
	}

	// Get the database connection.
	db := database.DB

	// Call the DeleteAccount service to delete the account.
	if err := services.DeleteAccount(accountID, db); err != nil {
		// If deletion fails, return a 500 Internal Server Error response.
		return utils.InternalServerError(c, err, "Failed to delete account")
	}

	// Return a 200 OK response.
	return utils.OKResponse(c, "Account deleted successfully", nil)
}

// GetTotalBalance handles the retrieval of the total balance of all active accounts.
// @Summary      Get total balance of all active accounts
// @Description  Retrieves the sum of balances from all active financial accounts for the authenticated user.
// @Tags         Accounts
// @Produce      json
// @Success      200  {object}  utils.Response
// @Failure      400  {object}  utils.Response
// @Failure      500  {object}  utils.Response
// @Router       /accounts/total-balance [get]
func GetTotalBalance(c *fiber.Ctx) error {
	// Get the user ID from the context.
	userID, err := uuid.Parse(c.Locals("user_id").(string))
	if err != nil {
		// If the user ID is invalid, return a 400 Bad Request response.
		return utils.BadResponse(c, err, "Invalid user ID")
	}

	// Get the database connection.
	db := database.DB

	// Call the GetTotalBalance service to calculate the total balance.
	totalBalance, err := services.GetTotalBalance(userID, db)
	if err != nil {
		// If calculation fails, return a 500 Internal Server Error response.
		return utils.InternalServerError(c, err, "Failed to retrieve total balance")
	}

	// Return a 200 OK response with the total balance.
	return utils.OKResponse(c, "Total balance retrieved successfully", totalBalance)
}
