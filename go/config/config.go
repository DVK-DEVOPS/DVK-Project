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
	apiKeyOnce   sync.Once
	cachedAPIKey string
)

func GetAPIKey() string { //returns the api key from .env in dev and from azure vault in prod
	apiKeyOnce.Do(func() {
		// Load .env if present
		if err := godotenv.Load(); err != nil {
			log.Println(".env file not found, assuming production")
		}

		// Check for dev API_KEY
		if key := os.Getenv("API_KEY"); key != "" {
			cachedAPIKey = key
			return
		}

		// Fetch from Key Vault
		vaultName := os.Getenv("KEYVAULT_NAME")
		secretName := os.Getenv("SECRET_NAME")
		if vaultName == "" || secretName == "" {
			log.Fatal("KEYVAULT_NAME or SECRET_NAME environment variables not set")
		}

		kvURL := fmt.Sprintf("https://%s.vault.azure.net/", vaultName)

		cred, err := azidentity.NewDefaultAzureCredential(nil)
		if err != nil {
			log.Fatalf("Failed to get Azure credential: %v", err)
		}

		client, err := azsecrets.NewClient(kvURL, cred, nil)
		if err != nil {
			log.Fatalf("Failed to create Key Vault client: %v", err)
		}

		resp, err := client.GetSecret(context.Background(), secretName, "", nil) // "" = latest version
		if err != nil {
			log.Fatalf("Failed to fetch secret: %v", err)
		}

		cachedAPIKey = *resp.Value

	})

	return cachedAPIKey
}
