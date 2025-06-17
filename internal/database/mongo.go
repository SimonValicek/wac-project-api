package database

import (
    "context"
    "log"
    "os"
    "strconv"
    "strings"               // ← added
    "time"

    "github.com/joho/godotenv"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

var (
    Client     *mongo.Client
    Collection *mongo.Collection
)

func ConnectMongo() {
    // only load .env in non-production
    if !strings.EqualFold(os.Getenv("AMBULANCE_API_ENVIRONMENT"), "production") {  // ← changed
        if err := godotenv.Load(); err != nil {
            log.Println("⚠️  No .env file found; using environment variables only")
        }
    }

    host := getenv("AMBULANCE_API_MONGODB_HOST", "localhost")
    port := getenv("AMBULANCE_API_MONGODB_PORT", "27017")
    user := getenv("AMBULANCE_API_MONGODB_USERNAME", "")
    pass := getenv("AMBULANCE_API_MONGODB_PASSWORD", "")
    dbName := getenv("AMBULANCE_API_MONGODB_DATABASE", "test")
    collName := getenv("AMBULANCE_API_MONGODB_COLLECTION", "dummy")
    timeoutSec := getenv("AMBULANCE_API_MONGODB_TIMEOUT_SECONDS", "10")
    timeout, _ := strconv.Atoi(timeoutSec)

    uri := "mongodb://"
if user != "" {
    uri += user
    if pass != "" {
        uri += ":" + pass
    }
    uri += "@"
}
uri += host + ":" + port

    log.Printf("🔌 Connecting to MongoDB at %s\n", uri)

    ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
    defer cancel()

    clientOpts := options.Client().ApplyURI(uri)
    client, err := mongo.Connect(ctx, clientOpts)
    if err != nil {
        log.Fatalf("❌ Failed to connect to MongoDB: %v", err)
    }

    if err := client.Ping(ctx, nil); err != nil {
        log.Fatalf("❌ MongoDB ping failed: %v", err)
    }
    log.Println("✅ Mongo ping successful")

    // Optional test insert
    if !strings.EqualFold(os.Getenv("AMBULANCE_API_ENVIRONMENT"), "production") {
        testDoc := map[string]string{"status": "ok"}
        if _, err := client.Database(dbName).Collection(collName).InsertOne(ctx, testDoc); err != nil {
            log.Printf("⚠️ Test insert failed: %v", err)
        } else {
            log.Println("✅ Inserted test doc into", dbName, ".", collName)
        }
    }

    Client = client
    Collection = client.Database(dbName).Collection(collName)
    log.Printf("📦 Collection initialized: %v\n", Collection.Name())
}

// getenv is a small helper to fall back to default
func getenv(key, def string) string {
    if v := os.Getenv(key); v != "" {
        return v
    }
    return def
}
