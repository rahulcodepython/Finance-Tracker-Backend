package v1

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/rahulcodepython/finance-tracker-backend/backend/database"
	"github.com/rahulcodepython/finance-tracker-backend/backend/services"
	"github.com/rahulcodepython/finance-tracker-backend/backend/utils"
)

// GenerateReport godoc
// @Summary Generate a financial report
// @Description Generates a financial report for a custom date range.
// @Tags reports
// @Security ApiKeyAuth
// @Produce  json
// @Param from query string false "Start date (YYYY-MM-DD)"
// @Param to query string false "End date (YYYY-MM-DD)"
// @Success 200 {object} map[string]interface{} "Report generated successfully"
// @Router /reports [get]
func GenerateReport(c *fiber.Ctx) error {
	userID, err := uuid.Parse(c.Locals("user_id").(string))
	if err != nil {
		return utils.BadResponse(c, err, "Invalid user ID")
	}
	from := c.Query("from")
	to := c.Query("to")

	db := database.DB

	report, err := services.GenerateReport(userID, from, to, db)
	if err != nil {
		return utils.InternalServerError(c, err, "Failed to generate report")
	}

	return utils.OKResponse(c, "Report generated successfully", report)
}

// ExportTransactions godoc
// @Summary Export transaction data to a CSV file
// @Description Exports transaction data for a custom date range to a CSV file.
// @Tags reports
// @Security ApiKeyAuth
// @Produce text/csv
// @Param from query string false "Start date (YYYY-MM-DD)"
// @Param to query string false "End date (YYYY-MM-DD)"
// @Success 200 {file} file "CSV file with transaction data"
// @Router /reports/export [get]
func ExportTransactions(c *fiber.Ctx) error {
	userID, err := uuid.Parse(c.Locals("user_id").(string))
	if err != nil {
		return utils.BadResponse(c, err, "Invalid user ID")
	}

	db := database.DB

	c.Set("Content-Type", "text/csv")
	c.Set("Content-Disposition", "attachment; filename=transactions.csv")

	if err := services.ExportTransactions(userID, c.Response().BodyWriter(), db); err != nil {
		return utils.InternalServerError(c, err, "Failed to export transactions")
	}

	return nil
}
