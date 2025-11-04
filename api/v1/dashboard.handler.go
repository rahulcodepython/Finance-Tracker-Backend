// Package v1 provides handlers for version 1 of the API.
package v1

import (
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
	startDate := c.Query("startDate")
	endDate := c.Query("endDate")

	// Get the database connection.
	db := database.DB

	// Call the GetDashboardSummary service to retrieve the dashboard data.
	summary, err := services.GetDashboardSummary(userID, startDate, endDate, db)
	if err != nil {
		return utils.InternalServerError(c, err, "Failed to get dashboard summary")
	}

	// Return a 200 OK response with the dashboard summary.
	return utils.OKResponse(c, "Dashboard data retrieved successfully", summary)
}

// ExportReports handles the export of transaction data to a CSV file.
// @Summary      Export transaction data
// @Description  Exports transaction data for the authenticated user to a CSV file, optionally filtered by a date range.
// @Tags         Reports
// @Produce      text/csv
// @Param        from query     string false "Start date for the export (YYYY-MM-DD)"
// @Param        to   query     string false "End date for the export (YYYY-MM-DD)"
// @Success      200  {file}    file   "CSV file with transaction data"
// @Failure      400  {object}  utils.Response
// @Failure      500  {object}  utils.Response
// @Router       /reports/export [get]
func ExportReports(c *fiber.Ctx) error {
	// Get the user ID from the context.
	userID, err := uuid.Parse(c.Locals("user_id").(string))
	if err != nil {
		return utils.BadResponse(c, err, "Invalid user ID")
	}

	// Get the database connection.
	db := database.DB

	// Set the response headers for a CSV file download.
	c.Set("Content-Type", "text/csv")
	c.Set("Content-Disposition", "attachment; filename=transactions.csv")

	// Call the ExportTransactions service to write the CSV data to the response body.
	if err := services.ExportTransactions(userID, c.Response().BodyWriter(), db); err != nil {
		return utils.InternalServerError(c, err, "Failed to export transactions")
	}

	return nil
}
