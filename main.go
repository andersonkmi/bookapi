package main

import (
	"log"

	"github.com/andersonkmi/bookapi/handlers"
	"github.com/gin-gonic/gin"
)

func main() {
	engine := gin.Default()
	if err := engine.SetTrustedProxies(nil); err != nil {
		log.Fatalf("failed to set trusted proxies: %v", err)
	}

	routerGroup := engine.Group("/api/v1")
	handlers.RegisterRoutes(routerGroup)

	if err := engine.Run(":8080"); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
