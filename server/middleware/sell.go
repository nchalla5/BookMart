package middleware

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/nchalla5/react-go-app/models" // Adjust this import path as necessary
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
	// rand.Seed(time.Now().UnixNano())
}

// generateRandomString generates a random string of n letters.
// func generateRandomString(n int) string {
// 	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

// 	s := make([]rune, n)
// 	for i := range s {
// 		s[i] = letters[rand.Intn(len(letters))]
// 	}
// 	return string(s)
// }

func CreateProductHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}

	// Load environment variables
	if err := godotenv.Load(); err != nil {
		http.Error(w, "Error loading .env file", http.StatusInternalServerError)
		return
	}

	var product models.Product
	err := json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		http.Error(w, "Error decoding JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	if product.ProductID == "" { // Check if ProductID is empty
		product.ProductID = uuid.New().String()[:5] // Generate a new UUID as a string for ProductID
	}

	// if product.ProductID == "" {
	// 	product.ProductID = generateRandomString(5) // Generates a random 5-letter string
	// }

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

	if product.Image != "" {
		s3URL, err := uploadImageToS3(product.Image)
		if err != nil {
			http.Error(w, "Failed to upload image to S3: "+err.Error(), http.StatusInternalServerError)
			return
		}
		product.Image = s3URL
	}

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

	// w.Header().Set("Content-Type", "application/json")
	// w.WriteHeader(http.StatusCreated) // 201 Created status code

	// // Prepare and send the response body
	// response := map[string]interface{}{
	// 	"message": "Product created successfully",
	// 	"product": product,
	// }
	// json.NewEncoder(w).Encode(response)
	// w.WriteHeader(http.StatusCreated)
	// json.NewEncoder(w).Encode(product)
}

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

func uploadImageToS3(imageURL string) (string, error) {
	resp, err := http.Get(imageURL)
	if err != nil {
		return "", fmt.Errorf("failed to fetch image: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("received non-200 response code: %d", resp.StatusCode)
	}

	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, resp.Body); err != nil {
		return "", fmt.Errorf("failed to read image body: %w", err)
	}

	imageKey := fmt.Sprintf("product-images/%s-%d", uuid.New().String(), time.Now().Unix())

	_, err = s3Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket:      aws.String(os.Getenv("AWS_BUCKET")),
		Key:         aws.String(imageKey),
		Body:        bytes.NewReader(buf.Bytes()),
		ContentType: aws.String(http.DetectContentType(buf.Bytes())),
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload image to S3: %w", err)
	}

	// Generate a signed URL for the uploaded image
	presignClient := s3.NewPresignClient(s3Client)
	presignedReq, err := presignClient.PresignGetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(os.Getenv("AWS_BUCKET")),
		Key:    aws.String(imageKey),
	}, s3.WithPresignExpires(24*time.Hour)) // Adjust the expiration time as needed
	if err != nil {
		return "", fmt.Errorf("failed to presign URL for S3 object: %w", err)
	}
	//fmt.Println(presignedReq.URL)
	return presignedReq.URL, nil
}
