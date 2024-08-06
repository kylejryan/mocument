package main

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MyEvent struct {
	Name string `json:"name"`
}

var dbClient *mongo.Client
var secretsClient SecretsManagerClient

func init() {
	// Initialize AWS Secrets Manager client with region
	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String("us-west-2"), // Set the appropriate region
	}))
	secretsClient = NewRealSecretsManager(secretsmanager.New(sess))

	// Fetch MongoDB URI from AWS Secrets Manager
	secretID := os.Getenv("MONGODB_SECRET_ID")
	if secretID == "" {
		fmt.Println("MONGODB_SECRET_ID is not set")
		os.Exit(1)
	}
	secretValue, err := secretsClient.GetSecretValue(&secretsmanager.GetSecretValueInput{
		SecretId: aws.String(secretID),
	})
	if err != nil {
		fmt.Printf("Failed to get secret value: %v\n", err)
		os.Exit(1)
	}
	mongoURI := *secretValue.SecretString

	// Initialize MongoDB client
	clientOptions := options.Client().ApplyURI(mongoURI)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		fmt.Printf("Failed to connect to MongoDB: %v\n", err)
		os.Exit(1)
	}
	dbClient = client
}

// Handler is our lambda handler invoked by the `lambda.Start` function call
func Handler(ctx context.Context, event MyEvent) (string, error) {
	if dbClient == nil {
		return "MongoDB client is not initialized", nil
	}
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

type SecretsManagerClient interface {
	GetSecretValue(input *secretsmanager.GetSecretValueInput) (*secretsmanager.GetSecretValueOutput, error)
}

type RealSecretsManager struct {
	client *secretsmanager.SecretsManager
}

func NewRealSecretsManager(client *secretsmanager.SecretsManager) *RealSecretsManager {
	return &RealSecretsManager{
		client: client,
	}
}

func (r *RealSecretsManager) GetSecretValue(input *secretsmanager.GetSecretValueInput) (*secretsmanager.GetSecretValueOutput, error) {
	return r.client.GetSecretValue(input)
}
