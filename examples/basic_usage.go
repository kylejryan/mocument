package main

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MyEvent struct {
	Name string `json:"name"`
}

var dbClient *mongo.Client

func init() {
	clientOptions := options.Client().ApplyURI(os.Getenv("MONGODB_URI"))
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		fmt.Printf("Failed to connect to MongoDB: %v\n", err)
		os.Exit(1)
	}
	dbClient = client
}

// Handler is our lambda handler invoked by the `lambda.Start` function call
func Handler(ctx context.Context, event MyEvent) (string, error) {
	collection := dbClient.Database("test").Collection("collection")

	// Insert a document
	doc := map[string]interface{}{"name": event.Name}
	_, err := collection.InsertOne(ctx, doc)
	if err != nil {
		return "", fmt.Errorf("failed to insert document: %w", err)
	}

	// Find the document
	filter := map[string]interface{}{"name": event.Name}
	var result map[string]interface{}
	err = collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		return "", fmt.Errorf("failed to find document: %w", err)
	}

	return fmt.Sprintf("Found document: %+v", result), nil
}

func main() {
	lambda.Start(Handler)
}
