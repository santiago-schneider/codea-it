package main

import (
	"codea-it/internal/app/adapters/inbound/http"
	"codea-it/internal/app/config"
	"github.com/gin-gonic/gin"
	"log"
)

func main() {
	cfg := config.LoadConfig()
	r := gin.Default()
	http.RegisterRoutes(r)
	err := r.Run(":" + cfg.ServerAddress)
	if err != nil {
		log.Fatal("Failed to start server: ", err)
	}
}
