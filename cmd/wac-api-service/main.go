package main

import (
    "log"
    "github.com/SimonValicek/wac-project-api/internal/database"
    sw "github.com/SimonValicek/wac-project-api/internal/wac_api"
    "github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"time"
)

func main() {
  log.Println("ℹ️  Connecting to MongoDB...")
  database.ConnectMongo()
  log.Println("✅ DB init done")

  routes := sw.ApiHandleFunctions{
    DefaultAPI: sw.NewReservationApi(),
  }

  log.Println("🚀 WAC Project API Server started")

  router := gin.Default()

  // Enable CORS middleware
  router.Use(cors.New(cors.Config{
    AllowOrigins:     []string{"http://localhost:3333"}, // Your frontend origin
    AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
    AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
    AllowCredentials: true,
    MaxAge:           12 * time.Hour,
  }))

  router = sw.NewRouterWithGinEngine(router, routes)
  log.Fatal(router.Run(":8080"))
}
