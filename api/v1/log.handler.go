// Package v1 provides handlers for version 1 of the API.
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

// GetLogs handles the retrieval of user activity logs.
// @Summary      Get user activity logs
// @Description  Retrieves a paginated list of activity logs for the authenticated user, with optional date filtering.
// @Tags         Logs
// @Produce      json
// @Param        page       query     int    false "Page number for pagination" default(1)
// @Param        limit      query     int    false "Number of logs per page" default(10)
// @Param        start_date query     string false "Start date for log filter (YYYY-MM-DD)"
// @Param        end_date   query     string false "End date for log filter (YYYY-MM-DD)"
// @Success      200        {object}  utils.Response
// @Failure      400        {object}  utils.Response
// @Failure      500        {object}  utils.Response
// @Router       /logs [get]
func GetLogs(c *fiber.Ctx) error {
	// Get the user ID from the context.
	userID, err := uuid.Parse(c.Locals("user_id").(string))
	if err != nil {
		return utils.BadResponse(c, err, "Invalid user ID")
	}

	// Get pagination and filter parameters from the query string.
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	startDateStr := c.Query("start_date", time.Now().AddDate(0, -1, 0).Format("2006-01-02"))
	endDateStr := c.Query("end_date", time.Now().Format("2006-01-02"))

	// Get the database connection.
	db := database.DB

	// Call the GetLogs service to retrieve the logs.
	logs, err := services.GetLogs(userID, startDateStr, endDateStr, page, limit, db)
	if err != nil {
		return utils.InternalServerError(c, err, "Failed to retrieve logs")
	}

	// Return a 200 OK response with the retrieved logs.
	return utils.OKResponse(c, "Activity logs retrieved successfully", logs)
}
