package env

import (
	"os"
	"sync"
)

var (
	envSyncOnce sync.Once

	isProduction bool
	isStaging    bool
	secretKey    string
)

// IsProduction returns true when env is set to production.
func IsProduction() bool {
	initializeEnvs()
	return isProduction
}

func IsStaging() bool {
	initializeEnvs()
	return isStaging
}

// GetSecretKey fetches the app secret key from env.
func GetSecretKey() string {
	initializeEnvs()
	return secretKey
}

func initializeEnvs() {
	envSyncOnce.Do(func() {
		isProduction = os.Getenv("ENV") == "production"
		isStaging = os.Getenv("ENV") == "staging"
		secretKey = os.Getenv("JWT_SECRET")
	})
}
