package database

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	Client     *mongo.Client
	Collection *mongo.Collection
)

func ConnectMongo() {
	log.Println("🔧 called ConnectMongo") // <- toto
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	uri := "mongodb://root:neUhaDnes@localhost:27017"
	log.Printf("🔌 Connecting to MongoDB at %s\n", uri)

	clientOpts := options.Client().ApplyURI(uri)

	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		log.Println("❌ Failed to connect to MongoDB:", err)
		return
	}

	if err := client.Ping(ctx, nil); err != nil {
	log.Println("❌ MongoDB ping failed:", err)
	return
}

	log.Println("✅ Mongo ping successful")

	testDoc := map[string]string{"status": "ok"}
_, err = client.Database("test").Collection("dummy").InsertOne(ctx, testDoc)
if err != nil {
	log.Fatalf("❌ Failed to insert test doc: %v", err)
}
log.Println("✅ Inserted test doc into test.dummy")

	Client = client
	Collection = client.Database("xvaliceks-project").Collection("parking_lots")
	log.Printf("📦 Collection initialized: %v\n", Collection.Name())
}
