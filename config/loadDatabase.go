package config

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var client *mongo.Client

func ConnectDB() {
	var err error

	serverAPIOptions := options.ServerAPI(options.ServerAPIVersion1)
	clientOptions := options.Client().ApplyURI(os.Getenv("MONGO_URI")).SetServerAPIOptions(serverAPIOptions)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	client, err = mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
		return
	}

	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		log.Fatal(err)
	}

	client.Database(os.Getenv("DATABASE_NAME")).CreateCollection(ctx, "users")
	client.Database(os.Getenv("DATABASE_NAME")).CreateCollection(ctx, "categories")
	client.Database(os.Getenv("DATABASE_NAME")).CreateCollection(ctx, "tasks")

	fmt.Println("Connected to MongoDB")
}

func GetDBConnection() *mongo.Database {
	return client.Database(os.Getenv("DATABASE_NAME"))
}
