package db

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func getConnection() (client *mongo.Client, ctx context.Context) {
	errEnv := godotenv.Load()
	if errEnv != nil {
		log.Fatal("Error loading .env file: ", errEnv)
	}

	connectionString := os.Getenv("CONECTIONSTRING_MONGO")

	fmt.Println(connectionString)
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(connectionString).SetServerAPIOptions(serverAPI)

	// Create a new client and connect to the server
	client, err := mongo.Connect(opts)
	if err != nil {
		panic(err)
	}

	return
}
