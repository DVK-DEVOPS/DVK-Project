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
			log.Println("config.GetSecret: .env file not found, assuming production")
		}

		// Check environment variable first (dev)
		if val := os.Getenv(envVar); val != "" {
			log.Printf("[GetSecret] Using %s from environment variable", envVar)
			secretCache[envVar] = val
			return
		}

		vaultName := os.Getenv("AZURE_KEYVAULT_NAME")
		if keyVaultName == "" {
			keyVaultName = vaultName
		}

		if keyVaultName == "" || secretName == "" {
			log.Fatalf("[GetSecret] Missing vault or secret name for %s (AZURE_KEYVAULT_NAME=%q, SECRET_NAME=%q)",
				envVar, keyVaultName, secretName)
		}

		kvURL := fmt.Sprintf("https://%s.vault.azure.net/", keyVaultName)
		log.Printf("[GetSecret] Fetching %s from Azure Key Vault %q (secret: %q)", envVar, keyVaultName, secretName)

		cred, err := azidentity.NewDefaultAzureCredential(nil)
		if err != nil {
			log.Fatalf("[GetSecret] Failed to get Azure credential: %v", err)
		}

		client, err := azsecrets.NewClient(kvURL, cred, nil)
		if err != nil {
			log.Fatalf("[GetSecret] Failed to create Key Vault client: %v", err)
		}

		resp, err := client.GetSecret(context.Background(), secretName, "", nil)
		if err != nil {
			log.Fatalf("[GetSecret] Failed to fetch secret %q from Key Vault %q: %v", secretName, keyVaultName, err)
		}

		if resp.Value == nil || *resp.Value == "" {
			log.Fatalf("[GetSecret] Secret %q in vault %q is empty", secretName, keyVaultName)
		}

		log.Printf("[GetSecret] Successfully fetched %s from Key Vault %q", envVar, keyVaultName)
		secretCache[envVar] = *resp.Value
	})

	return secretCache[envVar]
}

// GetSentryDSN returns the Sentry DSN from environment or empty if not set
func GetSentryDSN() string {
	sentryOnce.Do(func() {
		// Attempt to load .env
		if err := godotenv.Load(); err != nil {
			log.Println("[GetSentryDSN] .env not found — assuming production environment")
		} else {
			log.Println("[GetSentryDSN] .env file loaded successfully")
		}

		cachedSentryDSN = os.Getenv("SENTRY_DSN")
		cachedSentryEnv = os.Getenv("SENTRY_ENVIRONMENT")

		if cachedSentryEnv == "" {
			cachedSentryEnv = "dev"
			log.Println("[GetSentryDSN] SENTRY_ENVIRONMENT not set, defaulting to 'dev'")
		} else {
			log.Printf("[GetSentryDSN] SENTRY_ENVIRONMENT=%q", cachedSentryEnv)
		}

		if cachedSentryDSN == "" {
			log.Println("[GetSentryDSN] SENTRY_DSN not set (empty) — Sentry will be skipped")
		} else {
			// Print only prefix for safety
			prefix := cachedSentryDSN
			if len(prefix) > 20 {
				prefix = prefix[:20] + "..."
			}
			log.Printf("[GetSentryDSN] SENTRY_DSN loaded (prefix=%q)", prefix)
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
