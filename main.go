package main

import (
	"log"
	"os"

	"github.com/andersonkmi/bookapi/handlers"
	"github.com/andersonkmi/bookapi/repository"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "host=localhost user=pguser password=pgpwd dbname=bookapi port=5432 sslmode=disable"
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	bookRepo := repository.NewBookRepository(db)
	bookHandler := handlers.NewBookHandler(bookRepo)

	engine := gin.Default()
	if err := engine.SetTrustedProxies(nil); err != nil {
		log.Fatalf("failed to set trusted proxies: %v", err)
	}

	routerGroup := engine.Group("/api/v1")
	bookHandler.RegisterRoutes(routerGroup)

	if err := engine.Run(":8080"); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
