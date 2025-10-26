package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v2"
	"github.com/rahulcodepython/finance-tracker-backend/backend/config"
	"github.com/rahulcodepython/finance-tracker-backend/backend/database"
	"github.com/rahulcodepython/finance-tracker-backend/backend/pkg/scheduler"
	"github.com/rahulcodepython/finance-tracker-backend/backend/routes"
	"github.com/rahulcodepython/finance-tracker-backend/backend/utils"
)

// @title           Finance Tracker API
// @version         1.0
// @description     This is a sample server for a finance tracker application.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /api/v1
func main() {
	// Load application configuration from environment variables or a config file.
	cfg := config.LoadConfig()

	// Establish a connection to the database using the loaded configuration.
	db := database.Connect(cfg)

	// Set the local timezone for the application.
	utils.LoadTimezone()

	// Apply database migrations to ensure the schema is up to date.
	database.Migrate(db)

	// Start the background scheduler for recurring tasks.
	scheduler.StartScheduler(db)

	// Create a new Fiber web server instance with custom configuration.
	server := fiber.New(fiber.Config{
		AppName:       "Finance Tracker",
		ServerHeader:  "Finance Tracker",
		Prefork:       false,
		CaseSensitive: true,
		StrictRouting: true,
	})

	// Middleware to make the configuration accessible in request handlers.
	server.Use(func(c *fiber.Ctx) error {
		c.Locals("cfg", cfg)
		return c.Next()
	})

	// Set up all application routes and middleware.
	routes.Setup(server)

	// Construct the server address from the configuration.
	address := fmt.Sprintf("%s:%s", cfg.ServerConfig.Host, cfg.ServerConfig.Port)

	// Start the Fiber server in a new goroutine to avoid blocking the main thread.
	go func() {
		if err := server.Listen(address); err != nil {
			log.Panicf("Server error: %v", err)
		}
	}()

	// Create a channel to listen for operating system signals for graceful shutdown.
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	// Block until a signal is received.
	<-c

	// Initiate graceful shutdown of the server.
	fmt.Println("Gracefully shutting down...")
	_ = server.Shutdown()

	// Perform cleanup tasks, such as closing the database connection.
	fmt.Println("Running cleanup tasks...")
	defer db.Close()

	fmt.Println("Fiber was successful shutdown.")
}