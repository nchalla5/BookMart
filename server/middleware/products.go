package middleware

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"sort"
	"strconv"
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
	"github.com/nchalla5/react-go-app/constants"
	"github.com/nchalla5/react-go-app/models"
	"github.com/nchalla5/react-go-app/storage"
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

	queryParams := r.URL.Query()
	searchName := queryParams.Get("searchName")
	searchLocation := queryParams.Get("searchLocation")
	statusFilter := queryParams.Get("statusFilter")
	sortField := queryParams.Get("sortField")
	sortOrder := queryParams.Get("sortOrder")

	if !isAWSMode() {
		products, err := localStore.ListProducts(searchName, searchLocation, statusFilter, sortField, sortOrder)
		if err != nil {
			http.Error(w, "Failed to fetch products: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(products)
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

	scanInput := &dynamodb.ScanInput{
		TableName: aws.String("Products"),
	}

	var filterExpressions []string
	expressionAttributeValues := map[string]types.AttributeValue{}
	expressionAttributeNames := map[string]string{}

	if searchName != "" {
		filterExpressions = append(filterExpressions, "contains(#T, :title)")
		expressionAttributeValues[":title"] = &types.AttributeValueMemberS{Value: searchName}
		expressionAttributeNames["#T"] = "Title"
	}

	if searchLocation != "" {
		filterExpressions = append(filterExpressions, "contains(#L, :location)")
		expressionAttributeValues[":location"] = &types.AttributeValueMemberS{Value: searchLocation}
		expressionAttributeNames["#L"] = "Location" // Alias for reserved keyword
	}

	if statusFilter == "available" {
		filterExpressions = append(filterExpressions, "#S = :status")
		expressionAttributeValues[":status"] = &types.AttributeValueMemberS{Value: "available"}
		expressionAttributeNames["#S"] = "Status"
	}

	if len(filterExpressions) > 0 {
		scanInput.FilterExpression = aws.String(strings.Join(filterExpressions, " AND "))
		scanInput.ExpressionAttributeValues = expressionAttributeValues
		scanInput.ExpressionAttributeNames = expressionAttributeNames
	}

	out, err := svc.Scan(context.TODO(), scanInput)
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

	applySort(&products, sortField, sortOrder)

	// Generate S3 signed URLs for product images
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

	if !isAWSMode() {
		product, err := localStore.GetProduct(productID)
		if err != nil {
			http.Error(w, "Product not found", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(product)
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

func PurchaseProductHandler(w http.ResponseWriter, r *http.Request) {
	err := validateToken(r)
	if err != nil {
		http.Error(w, "Unauthorized: "+err.Error(), http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	productID := vars["id"]

	var request models.PurchaseRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := validateShippingAddress(request.ShippingAddress); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	buyer, err := getUsernameFromToken(r)
	if err != nil {
		http.Error(w, "Unauthorized: "+err.Error(), http.StatusUnauthorized)
		return
	}

	if !isAWSMode() {
		product, err := localStore.PurchaseProduct(productID, buyer, request.ShippingAddress)
		if err != nil {
			switch err {
			case storage.ErrProductNotFound:
				http.Error(w, "Product not found", http.StatusNotFound)
			case storage.ErrProductUnavailable:
				http.Error(w, "Product is no longer available", http.StatusConflict)
			default:
				http.Error(w, "Failed to complete purchase", http.StatusInternalServerError)
			}
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(models.ProductsApiResponse{
			Status:  "success",
			Message: "Purchase completed successfully",
			Data:    product,
		})
		return
	}

	product, err := fetchAWSProduct(productID)
	if err != nil {
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	}
	if product.Status != string(constants.Available) {
		http.Error(w, "Product is no longer available", http.StatusConflict)
		return
	}

	product.Buyer = buyer
	product.Status = string(constants.Sold)
	product.Shipping = &request.ShippingAddress
	product.PurchasedAt = time.Now().Format(time.RFC3339)

	if err := saveAWSProduct(*product); err != nil {
		http.Error(w, "Failed to complete purchase", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(models.ProductsApiResponse{
		Status:  "success",
		Message: "Purchase completed successfully",
		Data:    product,
	})
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

		return []byte(jwtSecret()), nil
	})

	if err != nil {
		return err
	}

	if !token.Valid {
		return errors.New("invalid token")
	}

	return nil
}

func validateShippingAddress(address models.ShippingAddress) error {
	fields := map[string]string{
		"street":       address.Street,
		"city":         address.City,
		"state":        address.State,
		"postalCode":   address.PostalCode,
		"country":      address.Country,
		"countryCode":  address.CountryCode,
		"mobileNumber": address.MobileNumber,
	}

	for name, value := range fields {
		if strings.TrimSpace(value) == "" {
			return fmt.Errorf("%s is required", name)
		}
	}

	return nil
}

func applySort(products *[]models.Product, sortField, sortOrder string) {
	switch sortField {
	case "name":
		sort.Slice(*products, func(i, j int) bool {
			left := strings.ToLower((*products)[i].Title)
			right := strings.ToLower((*products)[j].Title)
			if sortOrder == "desc" {
				return left > right
			}
			return left < right
		})
	case "cost":
		sort.Slice(*products, func(i, j int) bool {
			left, _ := strconv.ParseFloat((*products)[i].Cost, 64)
			right, _ := strconv.ParseFloat((*products)[j].Cost, 64)
			if sortOrder == "desc" {
				return left > right
			}
			return left < right
		})
	}
}

func fetchAWSProduct(productID string) (*models.Product, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(os.Getenv("AWS_REGION")))
	if err != nil {
		return nil, err
	}
	svc := dynamodb.NewFromConfig(cfg)

	out, err := svc.GetItem(context.TODO(), &dynamodb.GetItemInput{
		TableName: aws.String("Products"),
		Key: map[string]types.AttributeValue{
			"ProductID": &types.AttributeValueMemberS{Value: productID},
		},
	})
	if err != nil {
		return nil, err
	}
	if out.Item == nil {
		return nil, errors.New("product not found")
	}

	var product models.Product
	if err := attributevalue.UnmarshalMap(out.Item, &product); err != nil {
		return nil, err
	}

	return &product, nil
}

func saveAWSProduct(product models.Product) error {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(os.Getenv("AWS_REGION")))
	if err != nil {
		return err
	}
	svc := dynamodb.NewFromConfig(cfg)

	av, err := attributevalue.MarshalMap(product)
	if err != nil {
		return err
	}

	_, err = svc.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: aws.String("Products"),
		Item:      av,
	})
	return err
}
