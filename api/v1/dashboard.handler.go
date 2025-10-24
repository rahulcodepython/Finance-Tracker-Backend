package v1

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/rahulcodepython/finance-tracker-backend/backend/database"
	"github.com/rahulcodepython/finance-tracker-backend/backend/services"
	"github.com/rahulcodepython/finance-tracker-backend/backend/utils"
)

// GetDashboardSummary godoc
// @Summary Get a summary of the user's financial data for the dashboard
// @Description Gets a summary of the user's financial data for the dashboard, including total balance, current month's income/expenses/savings, recent transactions, income vs. expense data for the last 12 months, and monthly spending by category.
// @Tags dashboard
// @Security ApiKeyAuth
// @Produce  json
// @Success 200 {object} map[string]interface{} "Dashboard data retrieved successfully"
// @Router /dashboard [get]
func GetDashboardSummary(c *fiber.Ctx) error {
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

	summary, err := services.GetDashboardSummary(userID, page, limit, description, categoryID, accountID, budgetID, startDate, endDate, db)
	if err != nil {
		return utils.InternalServerError(c, err, "Failed to get dashboard summary")
	}

	return utils.OKResponse(c, "Dashboard data retrieved successfully", summary)
}
