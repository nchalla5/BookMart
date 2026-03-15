package middleware

import (
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/joho/godotenv"
	"github.com/nchalla5/react-go-app/storage"
)

var (
	envOnce    sync.Once
	localStore = storage.NewLocalStore(filepath.Join("data"))
)

func loadEnv() {
	envOnce.Do(func() {
		_ = godotenv.Load()
	})
}

func storageMode() string {
	loadEnv()
	mode := strings.ToLower(strings.TrimSpace(os.Getenv("STORAGE_MODE")))
	if mode == "" {
		return "local"
	}
	return mode
}

func isAWSMode() bool {
	return storageMode() == "aws"
}

func jwtSecret() string {
	loadEnv()
	secret := strings.TrimSpace(os.Getenv("JWT_KEY"))
	if secret == "" {
		return "book-mart-dev-secret"
	}
	return secret
}
