package middleware

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"time"

	// "fmt"
	"log"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/aws-sdk-go/aws"

	// "github.com/golang-jwt/jwt"
	"github.com/dgrijalva/jwt-go"
	"github.com/joho/godotenv"
	"github.com/nchalla5/react-go-app/models"
	"golang.org/x/crypto/bcrypt"
)

func Login(w http.ResponseWriter, r *http.Request) {
	//fmt.Println("Entered Login Page")
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
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion("us-west-1"), // Our AWS region
	)
	if err != nil {
		log.Fatalf("Unable to load SDK config, %v", err)
	}
	//fmt.Println("Starting Dynamo DB")
	svc := dynamodb.NewFromConfig(cfg)
	// insertCredential(svc, creds.EmailOrPhone, creds.Password)
	result, err := svc.GetItem(context.TODO(), &dynamodb.GetItemInput{
		TableName: aws.String("Credentials"),
		Key: map[string]types.AttributeValue{
			"UserName": &types.AttributeValueMemberS{Value: creds.Email},
		},
	})
	// if err != nil {
	// 	http.Error(w, "Failed to fetch user", http.StatusInternalServerError)
	// 	fmt.Println("Failed to fetch user: ", creds.Email)
	// 	return
	// }
	if err != nil {
		var aerr *types.ResourceNotFoundException
		if errors.As(err, &aerr) {
			log.Printf("Table not found: %v", aerr)
		} else {
			log.Printf("Failed to fetch user: %v, Error: %v", creds.Email, err)
		}
		http.Error(w, "Failed to fetch user", http.StatusInternalServerError)
		return
	}

	if result.Item == nil {
		http.Error(w, "User not found", http.StatusUnauthorized)
		//fmt.Println("User not found")
		return
	}
	storedPassword := result.Item["Password"].(*types.AttributeValueMemberS).Value
	//fmt.Println("Hashed Password in SignIn: ", result.Item["Password"])
	err = bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(creds.Password))
	if err != nil {
		http.Error(w, "Invalid login credentials", http.StatusUnauthorized)
		return
	}
	//fmt.Println("Encrypted Password")
	expirationTime := time.Now().Add(1 * time.Hour)
	claims := models.Claims{
		Email: creds.Email,
		Name:  result.Item["Name"].(*types.AttributeValueMemberS).Value,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	jwtKey := []byte(os.Getenv("JWT_KEY"))
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)

		//fmt.Println("Failed to generate token ", err)
		return
	}
	//fmt.Println("Token String in SignIn page: ", tokenString)s
	//log.Printf("Login attempt with Email/Phone: %s, Password: %s", creds.EmailOrPhone, creds.Password)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"token": tokenString})
}
