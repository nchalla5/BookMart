package middleware

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/s3"
<<<<<<< HEAD
	"github.com/dgrijalva/jwt-go"

	// "github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/nchalla5/react-go-app/constants"
	"github.com/nchalla5/react-go-app/constants"
	"github.com/nchalla5/react-go-app/models"
)

var s3Client *s3.Client

func init() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("No .env file found")
	}

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(os.Getenv("AWS_REGION")),
	)
	if err != nil {
		panic("unable to load SDK config, " + err.Error())
	}

	s3Client = s3.NewFromConfig(cfg)
}

func CreateProductHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}
	err := validatetoken(r)
	if err != nil {
		http.Error(w, "Unauthorized: "+err.Error(), http.StatusUnauthorized)
		return
	}

	var product models.Product
	contentType := r.Header.Get("Content-Type")

	// Handling multipart/form-data for direct file uploads
	if strings.HasPrefix(contentType, "multipart/form-data") {
		if err := r.ParseMultipartForm(10 << 20); err != nil { // 10 MB
			http.Error(w, "Error parsing multipart form: "+err.Error(), http.StatusBadRequest)
			return
		}
		product = extractProductFromForm(r)

		file, handler, err := r.FormFile("image")
		if err == nil {
			defer file.Close()
			s3URL, err := uploadImageToS3(file, handler.Filename)
			if err != nil {
				http.Error(w, "Failed to upload image to S3: "+err.Error(), http.StatusInternalServerError)
				return
			}
			product.Image = s3URL
		}
	} else if strings.HasPrefix(contentType, "application/json") {
		// Handling application/json for image URLs and other data
		if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
			http.Error(w, "Error decoding JSON: "+err.Error(), http.StatusBadRequest)
			return
		}

		if product.Image != "" {
			s3URL, err := uploadImageFromURL(product.Image)
			if err != nil {
				http.Error(w, "Failed to upload image from URL: "+err.Error(), http.StatusInternalServerError)
				return
			}
			product.Image = s3URL
		}
	} else {
		http.Error(w, "Unsupported Content-Type", http.StatusUnsupportedMediaType)
		return
	}

	if product.ProductID == "" { // Check if ProductID is empty
		product.ProductID = uuid.New().String()[:5] // Generate a new UUID as a string for ProductID
	}

	email, err := getUsernameFromToken(r)
	if err != nil {
		http.Error(w, "Unauthorized: "+err.Error(), http.StatusUnauthorized)
		return
	}

	product.Seller = email
	product.Status = string(constants.Available)

	email, err := getUsernameFromToken(r)
	if err != nil {
		http.Error(w, "Unauthorized: "+err.Error(), http.StatusUnauthorized)
		return
	}

	product.Seller = email
	product.Status = string(constants.Available)

	// Initialize AWS configuration
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(os.Getenv("AWS_REGION")),
		// AWS SDK v2 reads credentials from the standard locations
		// automatically, such as environment variables, AWS credentials file, etc.
	)
	if err != nil {
		http.Error(w, "Unable to load AWS config: "+err.Error(), http.StatusInternalServerError)
		return
	}

	svc := dynamodb.NewFromConfig(cfg)

	//s3Client := s3.NewFromConfig(cfg)

	// if product.Image != "" {
	// 	s3URL, err := uploadImageToS3(product.Image)
	// 	if err != nil {
	// 		http.Error(w, "Failed to upload image to S3: "+err.Error(), http.StatusInternalServerError)
	// 		return
	// 	}
	// 	product.Image = s3URL
	// }

	unique, err := isProductIDUnique(svc, product.ProductID, "Products")
	if err != nil {
		http.Error(w, "Failed to check ProductID uniqueness: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if !unique {
		http.Error(w, "ProductID must be unique", http.StatusBadRequest)
		return
	}

	// Convert Product struct to DynamoDB attribute value map
	av, err := attributevalue.MarshalMap(product)
	if err != nil {
		http.Error(w, "Failed to marshal product: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Put item into DynamoDB
	tableName := "Products" // Adjust to your DynamoDB table name
	_, err = svc.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: &tableName,
		Item:      av,
	})
	if err != nil {
		http.Error(w, "Failed to put item into DynamoDB: "+err.Error(), http.StatusInternalServerError)
		return
	}

	response := models.ProductsApiResponse{
		Status:  "success",
		Message: "Product created successfully",
		Data:    product, // Optionally include the product data or specific fields
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func extractProductFromForm(r *http.Request) models.Product {
	return models.Product{
		ProductID:   r.FormValue("productId"),
		Product:     r.FormValue("product"),
		Title:       r.FormValue("title"),
		Cost:        r.FormValue("cost"),
		Location:    r.FormValue("location"),
		Description: r.FormValue("description"),
	}
}

// func handleMultipartImageUpload(r *http.Request, product *models.Product) {
// 	file, handler, err := r.FormFile("image")
// 	if err == nil {
// 		defer file.Close()
// 		s3URL, err := uploadImageToS3(file, handler.Filename)
// 		if err != nil {
// 			fmt.Printf("Failed to upload image to S3: %v\n", err)
// 			return
// 		}
// 		product.Image = s3URL
// 	}
// }

func isProductIDUnique(svc *dynamodb.Client, productID string, tableName string) (bool, error) {
	key, err := attributevalue.MarshalMap(map[string]string{"ProductID": productID})
	if err != nil {
		return false, err
	}

	input := &dynamodb.GetItemInput{
		TableName: &tableName,
		Key:       key,
	}

	result, err := svc.GetItem(context.TODO(), input)
	if err != nil {
		return false, err
	}

	return result.Item == nil, nil // If Item is nil, ProductID is unique
}

func uploadImageToS3(file multipart.File, filename string) (string, error) {
	// Create a buffer to store the file
	buffer := bytes.NewBuffer(nil)
	if _, err := io.Copy(buffer, file); err != nil {
		return "", fmt.Errorf("failed to read file buffer: %w", err)
	}

	// Generate a unique file name for S3 to prevent name collisions
	uniqueFileName := fmt.Sprintf("product-images/%s-%d-%s", uuid.New().String(), time.Now().Unix(), filename)

	// Upload the file to S3
	_, err := s3Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket:      aws.String(os.Getenv("AWS_BUCKET")),
		Key:         aws.String(uniqueFileName),
		Body:        bytes.NewReader(buffer.Bytes()),
		ContentType: aws.String("image/jpeg"), // You might need to detect the content type from the file header
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload file to S3: %w", err)
	}

	return fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", os.Getenv("AWS_BUCKET"), os.Getenv("AWS_REGION"), uniqueFileName), nil
}

func uploadImageFromURL(imageURL string) (string, error) {
	resp, err := http.Get(imageURL)
	if err != nil {
		return "", fmt.Errorf("failed to download image from URL: %w", err)
	}
	defer resp.Body.Close()

	// Read the downloaded image into a buffer
	buffer := bytes.NewBuffer(nil)
	if _, err := io.Copy(buffer, resp.Body); err != nil {
		return "", fmt.Errorf("failed to read image body: %w", err)
	}

	// Extract filename from URL
	segments := strings.Split(imageURL, "/")
	filename := segments[len(segments)-1]

	// Generate a unique file name for S3
	uniqueFileName := fmt.Sprintf("product-images/%s-%d-%s", uuid.New().String(), time.Now().Unix(), filename)

	// Upload the image to S3
	_, err = s3Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket:      aws.String(os.Getenv("AWS_BUCKET")),
		Key:         aws.String(uniqueFileName),
		Body:        bytes.NewReader(buffer.Bytes()),
		ContentType: aws.String(http.DetectContentType(buffer.Bytes())),
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload image to S3: %w", err)
	}

	return fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", os.Getenv("AWS_BUCKET"), os.Getenv("AWS_REGION"), uniqueFileName), nil
}


func validatetoken(r *http.Request) error {
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

func getUsernameFromToken(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", errors.New("authorization header is required")
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	claims := &models.Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		// Make sure that the token method conforms to "SigningMethodHMAC"
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("JWT_KEY")), nil
	})

	if err != nil {
		return "", err
	}

	if !token.Valid {
		return "", errors.New("invalid token")
	}

	return "", nil

}
