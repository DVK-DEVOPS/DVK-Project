package config

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/keyvault/azsecrets"
	"github.com/joho/godotenv"
)

var (
	secretCache = make(map[string]string)
	secretOnce  = make(map[string]*sync.Once)

	// Sentry config caching
	sentryOnce      sync.Once
	cachedSentryDSN string
	cachedSentryEnv string
)

func GetSecret(envVar, keyVaultName, secretName string) string {
	if _, exists := secretOnce[envVar]; !exists {
		secretOnce[envVar] = &sync.Once{}
	}

	secretOnce[envVar].Do(func() {
		// Load .env if present
		if err := godotenv.Load(); err != nil {
			log.Println(".env file not found, assuming production")
		}

		// Check environment variable first (dev)
		if val := os.Getenv(envVar); val != "" {
			secretCache[envVar] = val
			return
		}

		// Fetch from Key Vault (prod)
		if keyVaultName == "" || secretName == "" {
			log.Fatalf("KEYVAULT_NAME or SECRET_NAME not set for %s", envVar)
		}

		kvURL := fmt.Sprintf("https://%s.vault.azure.net/", keyVaultName)
		cred, err := azidentity.NewDefaultAzureCredential(nil)
		if err != nil {
			log.Fatalf("Failed to get Azure credential: %v", err)
		}

		client, err := azsecrets.NewClient(kvURL, cred, nil)
		if err != nil {
			log.Fatalf("Failed to create Key Vault client: %v", err)
		}

		resp, err := client.GetSecret(context.Background(), secretName, "", nil)
		if err != nil {
			log.Fatalf("Failed to fetch secret %s: %v", secretName, err)
		}

		secretCache[envVar] = *resp.Value
	})

	return secretCache[envVar]
}

// GetSentryDSN returns the Sentry DSN from environment or empty if not set
func GetSentryDSN() string {
	sentryOnce.Do(func() {
		// Load .env if present in dev
		_ = godotenv.Load()
		cachedSentryDSN = os.Getenv("SENTRY_DSN")
		cachedSentryEnv = os.Getenv("SENTRY_ENVIRONMENT")
		if cachedSentryEnv == "" {
			cachedSentryEnv = "dev"
		}
	})
	return cachedSentryDSN
}

// GetSentryEnvironment returns the Sentry environment, defaults to "development"
func GetSentryEnvironment() string {
	// Ensure we've attempted to load once
	_ = GetSentryDSN()
	return cachedSentryEnv
}
