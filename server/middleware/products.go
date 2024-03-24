package middleware

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/nchalla5/react-go-app/models"
)

func extractKeyFromURL(url string) string {
	prefix := fmt.Sprintf("https://react-go-bucket.s3.%s.amazonaws.com/", os.Getenv("AWS_REGION"))
	return strings.TrimPrefix(url, prefix)
}

// ListProductsHandler displays all products available in the Products database.
func ListProductsHandler(w http.ResponseWriter, r *http.Request) {
	err := validateToken(r)
	if err != nil {
		http.Error(w, "Unauthorized: "+err.Error(), http.StatusUnauthorized)
		return
	}
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(os.Getenv("AWS_REGION")),
	)
	if err != nil {
		http.Error(w, "Unable to load AWS config: "+err.Error(), http.StatusInternalServerError)
		return
	}
	svc := dynamodb.NewFromConfig(cfg)

	out, err := svc.Scan(context.TODO(), &dynamodb.ScanInput{
		TableName: aws.String("Products"),
	})
	if err != nil {
		http.Error(w, "Failed to fetch products: "+err.Error(), http.StatusInternalServerError)
		return
	}

	var products []models.Product
	err = attributevalue.UnmarshalListOfMaps(out.Items, &products)
	if err != nil {
		http.Error(w, "Failed to unmarshal DynamoDB response: "+err.Error(), http.StatusInternalServerError)
		return
	}

	for i, product := range products {
		s3Key := extractKeyFromURL(product.Image) // Implement this based on your URL structure
		signedURL, err := generateSignedURL(s3Key)
		if err != nil {
			// Handle error, possibly continue with an error placeholder or log
			fmt.Println("Error generating signed URL:", err)
			continue
		}
		products[i].Image = signedURL
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}

// GetProductHandler displays single product details based on id.
func GetProductHandler(w http.ResponseWriter, r *http.Request) {
	err := validateToken(r)
	if err != nil {
		http.Error(w, "Unauthorized: "+err.Error(), http.StatusUnauthorized)
		return
	}
	vars := mux.Vars(r)
	productID := vars["id"]

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(os.Getenv("AWS_REGION")),
	)
	if err != nil {
		http.Error(w, "Unable to load AWS config: "+err.Error(), http.StatusInternalServerError)
		return
	}
	svc := dynamodb.NewFromConfig(cfg)

	out, err := svc.GetItem(context.TODO(), &dynamodb.GetItemInput{
		TableName: aws.String("Products"),
		Key: map[string]types.AttributeValue{
			"ProductID": &types.AttributeValueMemberS{Value: productID},
		},
	})
	if err != nil {
		http.Error(w, "Failed to fetch product: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if out.Item == nil {
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	}

	var product models.Product
	err = attributevalue.UnmarshalMap(out.Item, &product)
	if err != nil {
		http.Error(w, "Failed to unmarshal DynamoDB response: "+err.Error(), http.StatusInternalServerError)
		return
	}

	s3Key := extractKeyFromURL(product.Image)
	signedURL, err := generateSignedURL(s3Key)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to generate signed URL: %s", err), http.StatusInternalServerError)
		return
	}
	product.Image = signedURL

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(product)
}

func generateSignedURL(s3Key string) (string, error) {
	presignClient := s3.NewPresignClient(s3Client)
	presignedReq, err := presignClient.PresignGetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(os.Getenv("AWS_BUCKET")),
		Key:    aws.String(s3Key),
	}, s3.WithPresignExpires(1*time.Hour)) // Adjust expiration as needed
	if err != nil {
		return "", fmt.Errorf("failed to presign URL for S3 object: %w", err)
	}

	return presignedReq.URL, nil
}

func validateToken(r *http.Request) error {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return errors.New("authorization header is required")
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Make sure that the token method conform to "SigningMethodHMAC"
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(os.Getenv("JWT_KEY")), nil
	})

	if err != nil {
		return err
	}

	if !token.Valid {
		return errors.New("invalid token")
	}

	return nil
}
