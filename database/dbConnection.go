package database

import (
	"context"
	"fmt"
	"log"
	"time"

	cfg "gin-mongo-api/configs"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ConnectDB() *mongo.Client {
	config, err := cfg.LoadConfig(".")
	if err != nil {
		log.Fatal("Cannot load config:", err)
	}

	client, err := mongo.NewClient(options.Client().ApplyURI(config.MongoURI))
	if err != nil {
		log.Fatal(err)
	}

	ctx, calcel := context.WithTimeout(context.Background(), 10*time.Second)
	defer calcel()
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	//ping the database
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to MongoDB")
	return client
}

//Client instace
var DB *mongo.Client = ConnectDB()

//getting database collections
func GetCollection(client *mongo.Client, collectionName string) *mongo.Collection {
	collection := client.Database("bookstoreAPI").Collection(collectionName)
	return collection
}
