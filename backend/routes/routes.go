// Package routes defines the routes for the application.
package routes

import (
	"github.com/gofiber/fiber/v2"
	v1 "github.com/rahulcodepython/finance-tracker-backend/api/v1"
	"github.com/rahulcodepython/finance-tracker-backend/backend/database"
	"github.com/rahulcodepython/finance-tracker-backend/backend/middleware"
	"github.com/rahulcodepython/finance-tracker-backend/backend/utils"

	// "github.com/gofiber/fiber/v2/middleware/recover" is a middleware that recovers from panics anywhere in the stack chain.
	"github.com/gofiber/fiber/v2/middleware/recover"
)

// Setup configures the routes for the application.
func Setup(app *fiber.App) {
	// Apply the logger middleware to all routes.
	app.Use(middleware.Logger())

	// Apply the CORS middleware to all routes.
	app.Use(middleware.Cors())

	// Apply the Limiter middleware to all routes.
	app.Use(middleware.LimiterMiddleware())

	app.Use(recover.New())

	// Create a new group for the API routes.
	api := app.Group("/api")

	// Create a new group for version 1 of the API.
	v1Api := api.Group("/v1")

	// Define a health check route.
	v1Api.Get("/", func(c *fiber.Ctx) error {
		// Get the database connection.
		db := database.DB

		// Ping the database to check the connection.
		if err := utils.Ping(db); err != nil {
			// If the ping fails, return a 503 Service Unavailable error.
			return utils.InternalServerError(c, err, "Database connection error")
		}

		// If the ping is successful, return a 200 OK response.
		return utils.OKResponse(c, "Welcome to the Finance Tracker API", nil)
	})

	// Group routes for authentication.
	auth := v1Api.Group("/auth")
	auth.Post("/register", v1.Register)
	auth.Post("/login", v1.Login)
	auth.Get("/profile", middleware.DeserializeUser, v1.GetProfile)
	auth.Post("/change-password", middleware.DeserializeUser, v1.ChangePassword)
	auth.Get("/google/login", v1.GoogleLogin)
	auth.Get("/google/callback", v1.GoogleCallback)

	// Group routes for accounts, with authentication middleware.
	accounts := v1Api.Group("/accounts", middleware.DeserializeUser)
	accounts.Post("/create", v1.CreateAccount)
	accounts.Get("/", v1.GetAccounts)
	accounts.Patch("/update/:id", v1.UpdateAccount)
	accounts.Delete("/delete/:id", v1.DeleteAccount)
	accounts.Get("/total-balance", v1.GetTotalBalance)

	// Group routes for transactions, with authentication middleware.
	transactions := v1Api.Group("/transactions", middleware.DeserializeUser)
	transactions.Post("/create", v1.CreateTransaction)
	transactions.Get("/", v1.GetTransactions)
	transactions.Patch("/update/:id", v1.UpdateTransaction)
	transactions.Delete("/delete/:id", v1.DeleteTransaction)

	// Group routes for the dashboard, with authentication middleware.
	dashboard := v1Api.Group("/dashboard", middleware.DeserializeUser)
	dashboard.Get("/", v1.GetDashboardSummary)
	dashboard.Get("/export-report", v1.ExportReports)

	// Group routes for categories, with authentication middleware.
	categories := v1Api.Group("/categories", middleware.DeserializeUser)
	categories.Post("/create", v1.CreateCategory)
	categories.Get("/", v1.GetCategories)
	categories.Patch("/update/:id", v1.UpdateCategory)
	categories.Delete("/delete/:id", v1.DeleteCategory)

	// Group routes for budgets, with authentication middleware.
	budgets := v1Api.Group("/budgets", middleware.DeserializeUser)
	budgets.Post("/create", v1.CreateBudget)
	budgets.Get("/", v1.GetBudgets)
	budgets.Patch("/update/:id", v1.UpdateBudget)
	budgets.Delete("/delete/:id", v1.DeleteBudget)

	// Group routes for recurring transactions, with authentication middleware.
	recurringTransactions := v1Api.Group("/recurring-transactions", middleware.DeserializeUser)
	recurringTransactions.Post("/create", v1.CreateRecurringTransaction)
	recurringTransactions.Get("/", v1.GetRecurringTransactions)
	recurringTransactions.Patch("/update/:id", v1.UpdateRecurringTransaction)
	recurringTransactions.Delete("/delete/:id", v1.DeleteRecurringTransaction)

	// Group routes for logs, with authentication middleware.
	logs := v1Api.Group("/logs", middleware.DeserializeUser)
	logs.Get("/", v1.GetLogs)
}
