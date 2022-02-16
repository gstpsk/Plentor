package db

import (
	"context"
	"log"

	"github.com/gstpsk/Plentor/util"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func Connect() *mongo.Client {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(util.Config.MONGODB_URI))
	if err != nil {
		log.Fatalf("Failed to connect to database: %s", err)
	}

	err = client.Ping(context.Background(), readpref.Primary())
	if err != nil {
		log.Fatalf("Failed to connect to database: %s", err)
	}

	log.Printf("Database connection established.")

	return client
}

func GetUserCol() *mongo.Collection {
	client := Connect()
	collection := client.Database("Plentor").Collection("users")
	return collection
}

func GetEventCol() *mongo.Collection {
	client := Connect()
	collection := client.Database("Plentor").Collection("events")
	return collection
}
