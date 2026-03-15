package storage

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/nchalla5/react-go-app/constants"
	"github.com/nchalla5/react-go-app/models"
)

var (
	ErrUserExists         = errors.New("user already exists")
	ErrUserNotFound       = errors.New("user not found")
	ErrProductNotFound    = errors.New("product not found")
	ErrProductUnavailable = errors.New("product is not available")
)

type LocalStore struct {
	baseDir string
	mu      sync.Mutex
}

func NewLocalStore(baseDir string) *LocalStore {
	return &LocalStore{baseDir: baseDir}
}

func (s *LocalStore) CreateUser(user models.UserDetails) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	users, err := s.readUsers()
	if err != nil {
		return err
	}

	for _, existing := range users {
		if strings.EqualFold(existing.Email, user.Email) {
			return ErrUserExists
		}
	}

	users = append(users, user)
	return s.writeJSON("users.json", users)
}

func (s *LocalStore) FindUserByEmail(email string) (*models.UserDetails, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	users, err := s.readUsers()
	if err != nil {
		return nil, err
	}

	for _, user := range users {
		if strings.EqualFold(user.Email, email) {
			copy := user
			return &copy, nil
		}
	}

	return nil, ErrUserNotFound
}

func (s *LocalStore) CreateProduct(product models.Product) (models.Product, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	products, err := s.readProducts()
	if err != nil {
		return models.Product{}, err
	}

	if product.ProductID == "" {
		product.ProductID = uuid.New().String()[:8]
	}

	for _, existing := range products {
		if existing.ProductID == product.ProductID {
			return models.Product{}, errors.New("product id already exists")
		}
	}

	products = append(products, product)
	if err := s.writeJSON("products.json", products); err != nil {
		return models.Product{}, err
	}

	return product, nil
}

func (s *LocalStore) ListProducts(searchName, searchLocation, statusFilter, sortField, sortOrder string) ([]models.Product, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	products, err := s.readProducts()
	if err != nil {
		return nil, err
	}

	filtered := make([]models.Product, 0, len(products))
	for _, product := range products {
		if searchName != "" && !strings.Contains(strings.ToLower(product.Title), strings.ToLower(searchName)) {
			continue
		}
		if searchLocation != "" && !strings.Contains(strings.ToLower(product.Location), strings.ToLower(searchLocation)) {
			continue
		}
		if statusFilter != "" && !strings.EqualFold(product.Status, statusFilter) {
			continue
		}
		filtered = append(filtered, product)
	}

	switch sortField {
	case "name":
		sort.Slice(filtered, func(i, j int) bool {
			left := strings.ToLower(filtered[i].Title)
			right := strings.ToLower(filtered[j].Title)
			if sortOrder == "desc" {
				return left > right
			}
			return left < right
		})
	case "cost":
		sort.Slice(filtered, func(i, j int) bool {
			left, _ := strconv.ParseFloat(filtered[i].Cost, 64)
			right, _ := strconv.ParseFloat(filtered[j].Cost, 64)
			if sortOrder == "desc" {
				return left > right
			}
			return left < right
		})
	default:
		sort.Slice(filtered, func(i, j int) bool {
			return filtered[i].Title < filtered[j].Title
		})
	}

	return filtered, nil
}

func (s *LocalStore) GetProduct(productID string) (*models.Product, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	products, err := s.readProducts()
	if err != nil {
		return nil, err
	}

	for _, product := range products {
		if product.ProductID == productID {
			copy := product
			return &copy, nil
		}
	}

	return nil, ErrProductNotFound
}

func (s *LocalStore) PurchaseProduct(productID, buyer string, shipping models.ShippingAddress) (*models.Product, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	products, err := s.readProducts()
	if err != nil {
		return nil, err
	}

	for i, product := range products {
		if product.ProductID != productID {
			continue
		}
		if product.Status != "" && product.Status != string(constants.Available) {
			return nil, ErrProductUnavailable
		}

		products[i].Buyer = buyer
		products[i].Status = string(constants.Sold)
		products[i].Shipping = &shipping
		products[i].PurchasedAt = time.Now().Format(time.RFC3339)

		if err := s.writeJSON("products.json", products); err != nil {
			return nil, err
		}

		copy := products[i]
		return &copy, nil
	}

	return nil, ErrProductNotFound
}

func (s *LocalStore) readUsers() ([]models.UserDetails, error) {
	var users []models.UserDetails
	if err := s.readJSON("users.json", &users); err != nil {
		return nil, err
	}
	return users, nil
}

func (s *LocalStore) readProducts() ([]models.Product, error) {
	var products []models.Product
	if err := s.readJSON("products.json", &products); err != nil {
		return nil, err
	}
	return products, nil
}

func (s *LocalStore) readJSON(name string, target interface{}) error {
	path, err := s.ensureFile(name)
	if err != nil {
		return err
	}

	bytes, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	if len(bytes) == 0 {
		bytes = []byte("[]")
	}

	return json.Unmarshal(bytes, target)
}

func (s *LocalStore) writeJSON(name string, value interface{}) error {
	path, err := s.ensureFile(name)
	if err != nil {
		return err
	}

	bytes, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, bytes, 0o644)
}

func (s *LocalStore) ensureFile(name string) (string, error) {
	if err := os.MkdirAll(s.baseDir, 0o755); err != nil {
		return "", err
	}

	path := filepath.Join(s.baseDir, name)
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		if err := os.WriteFile(path, []byte("[]"), 0o644); err != nil {
			return "", err
		}
	}

	return path, nil
}
