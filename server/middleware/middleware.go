package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/joho/godotenv"
	"github.com/nchalla5/react-go-app/models"
)

func Login(w http.ResponseWriter, r *http.Request) {
	//w.Write([]byte("Welcome to the Login Page"))
	var creds models.CredsStruct
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	// awsAccessKeyID := os.Getenv("AWS_ACCESS_KEY_ID")
	// awsSecretAccessKey := os.Getenv("AWS_SECRET_ACCESS_KEY")
	// awsRegion := os.Getenv("AWS_REGION")
	// Load the AWS configuration
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion("us-west-1"), // Or your AWS region
	)
	if err != nil {
		log.Fatalf("Unable to load SDK config, %v", err)
	}

	// Create a DynamoDB client
	svc := dynamodb.NewFromConfig(cfg)

	// Example operations
	insertCredential(svc, creds.EmailOrPhone, creds.Password)
	fetchCredentials(svc)

	// fmt.Println("AWS Access Key ID:", awsAccessKeyID)
	//log.Printf("Login attempt with Email/Phone: %s, Password: %s", creds.EmailOrPhone, creds.Password)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}

func insertCredential(svc *dynamodb.Client, username, password string) {
	_, err := svc.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: aws.String("Credentials"),
		Item: map[string]types.AttributeValue{
			"UserName": &types.AttributeValueMemberS{Value: username},
			"Password": &types.AttributeValueMemberS{Value: password}, // Consider using a hash function
		},
	})

	if err != nil {
		log.Fatalf("Failed to insert data, %v", err)
	}

	fmt.Println("Successfully inserted credential")
}

func fetchCredentials(svc *dynamodb.Client) {
	result, err := svc.Scan(context.TODO(), &dynamodb.ScanInput{
		TableName: aws.String("Credentials"),
	})
	if err != nil {
		log.Fatalf("Failed to fetch data, %v", err)
	}

	fmt.Println("Credentials:")
	for _, item := range result.Items {
		userName := item["UserName"].(*types.AttributeValueMemberS).Value
		password := item["Password"].(*types.AttributeValueMemberS).Value // Assuming password is also stored as string
		fmt.Printf("UserName: %s, Password: %s\n", userName, password)
	}
}
