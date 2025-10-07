package v1

import (
	"database/sql"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/rahulcodepython/finance-tracker-backend/backend/services"
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
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": "Invalid user ID", "error": err.Error()})
	}

	db := c.Locals("db").(*sql.DB)

	summary, err := services.GetDashboardSummary(userID, db)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"success": false, "message": "Failed to get dashboard summary", "error": err.Error()})
	}

	return c.JSON(fiber.Map{"success": true, "message": "Dashboard data retrieved successfully", "data": summary})
}
