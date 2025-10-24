package v1

import (
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/rahulcodepython/finance-tracker-backend/backend/database"
	"github.com/rahulcodepython/finance-tracker-backend/backend/services"
	"github.com/rahulcodepython/finance-tracker-backend/backend/utils"
)

// GetLogs godoc
// @Summary Get user activity logs
// @Description Retrieves a paginated list of activity logs for the authenticated user within a specified date range.
// @Tags logs
// @Security ApiKeyAuth
// @Produce json
// @Param page query int false "Page number for pagination" default(1)
// @Param limit query int false "Number of items per page" default(10)
// @Param start_date query string false "Start date for filtering logs (YYYY-MM-DD)"
// @Param end_date query string false "End date for filtering logs (YYYY-MM-DD)"
// @Success 200 {object} map[string]interface{} "Activity logs retrieved successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request or user ID"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /logs [get]
func GetLogs(c *fiber.Ctx) error {
	userID, err := uuid.Parse(c.Locals("user_id").(string))
	if err != nil {
		return utils.BadResponse(c, err, "Invalid user ID")
	}

	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	startDateStr := c.Query("start_date", time.Now().AddDate(0, -1, 0).Format("2006-01-02"))
	endDateStr := c.Query("end_date", time.Now().Format("2006-01-02"))

	db := database.DB

	logs, err := services.GetLogs(userID, startDateStr, endDateStr, page, limit, db)
	if err != nil {
		return utils.InternalServerError(c, err, "Failed to retrieve logs")
	}

	return utils.OKResponse(c, "Activity logs retrieved successfully", logs)
}
