package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bagasriatmaputera/novelku/internal/config"
	"github.com/bagasriatmaputera/novelku/internal/database"
	"github.com/bagasriatmaputera/novelku/pkg/response"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
)

func main() {
	cfg := config.Load()

	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatalf("database: %v", err)
	}

	// Migrate from package database
	if err := database.Migrate(db); err != nil {
		log.Fatalf("migrate: %v", err)
	}

	app := fiber.New(fiber.Config{
		AppName:      "Novel API",
		ErrorHandler: response.ErrorHandler,
	})

	app.Use(recover.New())   //mencegah server mati ketika panic
	app.Use(requestid.New()) //membuat uid setiap request
	app.Use(logger.New())    //mecatat log tiap request
	// setup allow cors
	app.Use(cors.New(cors.Config{
		AllowOrigins:     cfg.CORSOrigins,
		AllowCredentials: true,
	}))

	// check db health
	app.Get("/health", func(c *fiber.Ctx) error {
		sqlDB, err := db.DB()
		if err != nil || sqlDB.Ping() != nil {
			return response.Error(c, fiber.StatusServiceUnavailable, "DB_DOWN", "database tidak dapat dihubungi")
		}
		return response.Success(c, fiber.Map{
			"status": "OK",
		})
	})

	// Jalankan server di goroutine supaya main bisa menunggu sinyal shutdown.
	go func() {
		if err := app.Listen(":" + cfg.AppPort); err != nil {
			log.Printf("server berhenti: %v", err)
		}
	}()

	// handle shutdown server
	// Graceful shutdown: tunggu Ctrl+C / SIGTERM, selesaikan request berjalan, baru server mati.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Println("mematikan server..")
	if err := app.ShutdownWithTimeout(10 * time.Second); err != nil {
		log.Printf("shutdown error: %v", err)
	}

	if sqlDB, err := db.DB(); err == nil {
		_ = sqlDB.Close()
	}

}
