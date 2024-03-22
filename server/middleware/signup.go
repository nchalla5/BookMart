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
	"golang.org/x/crypto/bcrypt"
)

func Signup(w http.ResponseWriter, r *http.Request) {
	//w.Write([]byte("Welcome to the Login Page"))
	var details models.UserDetails
	err := json.NewDecoder(r.Body).Decode(&details)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(details.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}
	details.Password = string(hashedPassword)
	fmt.Println("Hashed Password in SignIn: ", details.Password)
	// log.Printf("Signup attempt with Email: %s, Name: %s, Password: %s", details.Email, details.Name, hashedPassword)
	err = godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion("us-west-1"), // Our AWS region
	)
	if err != nil {
		log.Fatalf("Unable to load SDK config, %v", err)
	}
	svc := dynamodb.NewFromConfig(cfg)
	insertCredential(svc, details)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "success", "message": "success"})
}

func insertCredential(svc *dynamodb.Client, details models.UserDetails) {
	_, err := svc.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: aws.String("Credentials"),
		Item: map[string]types.AttributeValue{
			"Name":     &types.AttributeValueMemberS{Value: details.Name},
			"UserName": &types.AttributeValueMemberS{Value: details.Email},
			"Password": &types.AttributeValueMemberS{Value: details.Password}, // Consider using a hash function
		},
	})

	if err != nil {
		log.Fatalf("Failed to insert data, %v", err)
	}

	fmt.Println("Successfully inserted credential")
}
