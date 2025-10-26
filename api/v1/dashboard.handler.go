// Package v1 provides handlers for version 1 of the API.
package v1

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/rahulcodepython/finance-tracker-backend/backend/database"
	"github.com/rahulcodepython/finance-tracker-backend/backend/services"
	"github.com/rahulcodepython/finance-tracker-backend/backend/utils"
)

// GetDashboardSummary handles the retrieval of a summary of the user's financial data.
// @Summary      Get dashboard summary
// @Description  Retrieves a financial summary for the authenticated user, including balances, income, expenses, and transactions.
// @Tags         Dashboard
// @Produce      json
// @Param        page         query     int    false "Page number for transaction pagination"
// @Param        limit        query     int    false "Number of transactions per page"
// @Param        description  query     string false "Filter transactions by description"
// @Param        category     query     string false "Filter transactions by category ID"
// @Param        account      query     string false "Filter transactions by account ID"
// @Param        budget       query     string false "Filter transactions by budget ID"
// @Param        startDate    query     string false "Start date for transaction filter (YYYY-MM-DD)"
// @Param        endDate      query     string false "End date for transaction filter (YYYY-MM-DD)"
// @Success      200          {object}  utils.Response
// @Failure      400          {object}  utils.Response
// @Failure      500          {object}  utils.Response
// @Router       /dashboard [get]
func GetDashboardSummary(c *fiber.Ctx) error {
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

	// Call the GetDashboardSummary service to retrieve the dashboard data.
	summary, err := services.GetDashboardSummary(userID, page, limit, description, categoryID, accountID, budgetID, startDate, endDate, db)
	if err != nil {
		return utils.InternalServerError(c, err, "Failed to get dashboard summary")
	}

	// Return a 200 OK response with the dashboard summary.
	return utils.OKResponse(c, "Dashboard data retrieved successfully", summary)
}
