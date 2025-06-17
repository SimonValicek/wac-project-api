package main

import (
    "log"
    "time"

    "github.com/joho/godotenv"
    "github.com/SimonValicek/wac-project-api/internal/database"
    sw "github.com/SimonValicek/wac-project-api/internal/wac_api"
    "github.com/gin-contrib/cors"
    "github.com/gin-gonic/gin"
)

func main() {
    // Load local .env if present
    if err := godotenv.Load(); err != nil {
        log.Println("⚠️  No .env file found, falling back to environment variables")
    }

    log.Println("ℹ️  Connecting to MongoDB...")
    database.ConnectMongo()
    log.Println("✅ DB init done")

    routes := sw.ApiHandleFunctions{
        DefaultAPI: sw.NewReservationApi(),
    }

    log.Println("🚀 WAC Project API Server started on port 8080")

    router := gin.Default()

    // Enable CORS middleware
    router.Use(cors.New(cors.Config{
        AllowOrigins:     []string{"http://localhost:3333"},
        AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
        AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
        AllowCredentials: true,
        MaxAge:           12 * time.Hour,
    }))

    router = sw.NewRouterWithGinEngine(router, routes)
    // Always bind to port 8080 locally
    log.Fatal(router.Run(":8080"))
}