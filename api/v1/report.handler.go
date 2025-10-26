// Package v1 provides handlers for version 1 of the API.
package v1

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/rahulcodepython/finance-tracker-backend/backend/database"
	"github.com/rahulcodepython/finance-tracker-backend/backend/services"
	"github.com/rahulcodepython/finance-tracker-backend/backend/utils"
)

// GenerateReport handles the generation of a financial report.
// @Summary      Generate a financial report
// @Description  Generates a financial report for the authenticated user within a specified date range.
// @Tags         Reports
// @Produce      json
// @Param        from query     string false "Start date for the report (YYYY-MM-DD)"
// @Param        to   query     string false "End date for the report (YYYY-MM-DD)"
// @Success      200  {object}  utils.Response
// @Failure      400  {object}  utils.Response
// @Failure      500  {object}  utils.Response
// @Router       /reports [get]
func GenerateReport(c *fiber.Ctx) error {
	// Get the user ID from the context.
	userID, err := uuid.Parse(c.Locals("user_id").(string))
	if err != nil {
		return utils.BadResponse(c, err, "Invalid user ID")
	}
	// Get the start and end dates from the query parameters.
	from := c.Query("from")
	to := c.Query("to")

	// Get the database connection.
	db := database.DB

	// Call the GenerateReport service to generate the report.
	report, err := services.GenerateReport(userID, from, to, db)
	if err != nil {
		return utils.InternalServerError(c, err, "Failed to generate report")
	}

	// Return a 200 OK response with the generated report.
	return utils.OKResponse(c, "Report generated successfully", report)
}

// ExportTransactions handles the export of transaction data to a CSV file.
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
func ExportTransactions(c *fiber.Ctx) error {
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
