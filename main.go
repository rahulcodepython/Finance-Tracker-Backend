package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/rahulcodepython/finance-tracker-backend/backend/config"
	"github.com/rahulcodepython/finance-tracker-backend/backend/database"
	"github.com/rahulcodepython/finance-tracker-backend/backend/pkg/scheduler"
	"github.com/rahulcodepython/finance-tracker-backend/backend/routes"

	"github.com/gofiber/fiber/v2"
)

func main() {
	cfg := config.LoadConfig()

	db := database.Connect(cfg)

	database.CreateTables(db)

	scheduler.StartScheduler(db)

	server := fiber.New()

	server.Use(func(c *fiber.Ctx) error {
		c.Locals("db", db)
		c.Locals("cfg", cfg)
		return c.Next()
	})

	// app.Get("/swagger/*", swagger.HandlerDefault)

	// router.Router() is called to set up all the application routes and middleware.
	// It takes the Fiber server, configuration, and database connection as arguments.
	routes.Setup(server)

	// address is a string that represents the server address.
	// It is constructed by combining the server host and port from the configuration.
	address := fmt.Sprintf("%s:%s", cfg.ServerConfig.Host, cfg.ServerConfig.Port)

	// A new goroutine is started to run the Fiber server.
	// This allows the main goroutine to continue and handle graceful shutdown.
	go func() {
		// server.Listen() starts the HTTP server and listens for incoming requests on the specified address.
		if err := server.Listen(address); err != nil {
			// If an error occurs while starting the server, log the error and panic.
			log.Panicf("Server error: %v", err)
		}
	}()

	// c is a channel that will receive operating system signals.
	// It has a buffer size of 1.
	c := make(chan os.Signal, 1)
	// signal.Notify() registers the given channel to receive notifications of the specified signals.
	// In this case, it listens for os.Interrupt (Ctrl+C) and syscall.SIGTERM.
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	// This is a blocking call that waits for a signal to be received on the channel c.
	<-c

	// A message is printed to the console to indicate that the server is shutting down.
	fmt.Println("Gracefully shutting down...")
	// server.Shutdown() gracefully shuts down the server without interrupting any active connections.
	_ = server.Shutdown()

	// A message is printed to the console to indicate that cleanup tasks are running.
	fmt.Println("Running cleanup tasks...")
	// db.Close() closes the database connection.
	_ = db.Close()

	// A message is printed to the console to indicate that the server has shut down successfully.
	fmt.Println("Fiber was successful shutdown.")
}
