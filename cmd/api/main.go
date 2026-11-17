/*
This is the program’s starting point. 
It loads config (.env), 
	opens the database connection, 
	wires everything together (repositories → services → handlers → routes), 
	starts the HTTP server, 
	and handles graceful shutdown on SIGINT/SIGTERM. 
	It should contain almost no business logic.
*/
package main

import (
	"context" // for graceful shutdown of the HTTP server
	"database/sql"
	"log"
	"net/http"
	"os"
	"os/signal" 
	"syscall" 
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/34codezak/makeworkflow/internal/config"
	// "github.com/34codezak/makeworkflow/internal/handlers"
	// "github.com/34codezak/makeworkflow/internal/repositories"
	// "github.com/34codezak/makeworkflow/internal/routes"
	// "github.com/34codezak/makeworkflow/internal/services"
	"github.com/joho/godotenv"
)

func main() {

	// Load environment variables from .env file if it exists, otherwise use system environment
	if err := godotenv.Load(); err != nil {
    	if os.IsNotExist(err) {
        	log.Println("No .env file found, using system environment")
    	} else {
        	log.Fatal(err) // real problem: file exists but couldn't be read
    	}
	}   

	cfg := config.Load()

	db, err := sql.Open("pgx", cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	// Wire dependencies bottom-up
	// userRepo := repositories.NewUserRepository(db)
	// taskRepo := repositories.NewTaskRepository(db)
	// projectRepo := repositories.NewProjectRepository(db)

	// authService := services.NewAuthService(userRepo, cfg)
	// taskService := services.NewTaskService(taskRepo, projectRepo)
	// projectService := services.NewProjectService(projectRepo)

	// authHandler := handlers.NewAuthHandler(authService)
	// taskHandler := handlers.NewTaskHandler(taskService)
	// projectHandler := handlers.NewProjectHandler(projectService)

	// r := routes.NewRouter(authHandler, taskHandler, projectHandler)

	// Initialize the router with handlers and middleware
	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	// Start the server in a goroutine so that it doesn't block.
	go func() {
		log.Printf("Server listening on :%s", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGNT, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}
	log.Println("Server stopped. Exiting...")
}
