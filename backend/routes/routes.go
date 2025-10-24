package routes

import (
	"github.com/gofiber/fiber/v2"
	v1 "github.com/rahulcodepython/finance-tracker-backend/api/v1"
	"github.com/rahulcodepython/finance-tracker-backend/backend/database"
	"github.com/rahulcodepython/finance-tracker-backend/backend/middleware"
	"github.com/rahulcodepython/finance-tracker-backend/backend/utils"
)

func Setup(app *fiber.App) {
	app.Use(middleware.Logger())

	api := app.Group("/api")

	v1Api := api.Group("/v1")

	v1Api.Get("/", func(c *fiber.Ctx) error {
		db := database.DB
		// log.Println(db) // This is good for debugging but noisy

		// Call the modified Ping function and check its error
		if err := utils.Ping(db); err != nil {
			// If ping fails, return a 503 Service Unavailable error
			return utils.InternalServerError(c, err, "Database connection error")
		}

		// If ping is successful, return the normal response
		return utils.OKResponse(c, "Welcome to the Finance Tracker API", nil)
	})

	auth := v1Api.Group("/auth")
	auth.Post("/register", v1.Register)
	auth.Post("/login", v1.Login)
	auth.Get("/profile", middleware.DeserializeUser, v1.GetProfile)
	auth.Post("/change-password", middleware.DeserializeUser, v1.ChangePassword)
	auth.Get("/google/login", v1.GoogleLogin)
	auth.Get("/google/callback", v1.GoogleCallback)

	accounts := v1Api.Group("/accounts", middleware.DeserializeUser)
	accounts.Post("/create", v1.CreateAccount)
	accounts.Get("/", v1.GetAccounts)
	accounts.Patch("/update/:id", v1.UpdateAccount)
	accounts.Delete("/delete/:id", v1.DeleteAccount)
	accounts.Get("/total-balance", v1.GetTotalBalance)

	transactions := v1Api.Group("/transactions", middleware.DeserializeUser)
	transactions.Post("/create", v1.CreateTransaction)
	transactions.Get("/", v1.GetTransactions)
	transactions.Patch("/update/:id", v1.UpdateTransaction)
	transactions.Delete("/delete/:id", v1.DeleteTransaction)
	transactions.Get("/aggregate", v1.GetAggregateData)

	dashboard := v1Api.Group("/dashboard", middleware.DeserializeUser)
	dashboard.Get("/", v1.GetDashboardSummary)

	reports := v1Api.Group("/reports", middleware.DeserializeUser)
	reports.Get("/", v1.GenerateReport)
	reports.Get("/export", v1.ExportTransactions)

	categories := v1Api.Group("/categories", middleware.DeserializeUser)
	categories.Post("/create", v1.CreateCategory)
	categories.Get("/", v1.GetCategories)
	categories.Patch("/update/:id", v1.UpdateCategory)
	categories.Delete("/delete/:id", v1.DeleteCategory)

	budgets := v1Api.Group("/budgets", middleware.DeserializeUser)
	budgets.Post("/create", v1.CreateBudget)
	budgets.Get("/", v1.GetBudgets)
	budgets.Patch("/update/:id", v1.UpdateBudget)
	budgets.Delete("/delete/:id", v1.DeleteBudget)

	recurringTransactions := v1Api.Group("/recurring-transactions", middleware.DeserializeUser)
	recurringTransactions.Post("/create", v1.CreateRecurringTransaction)
	recurringTransactions.Get("/", v1.GetRecurringTransactions)
	recurringTransactions.Patch("/update/:id", v1.UpdateRecurringTransaction)
	recurringTransactions.Delete("/delete/:id", v1.DeleteRecurringTransaction)

	logs := v1Api.Group("/logs", middleware.DeserializeUser)
	logs.Get("/", v1.GetLogs)
}
