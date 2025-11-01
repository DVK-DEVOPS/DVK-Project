package test

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/keyvault/azsecrets"
)

// checks that the application can fetch a secret from azure vault
func TestKeyVaultFetch(t *testing.T) {
	vaultName := os.Getenv("KEYVAULT_NAME")
	secretName := os.Getenv("SECRET_NAME")

	if vaultName == "" || secretName == "" {
		t.Skip("Skipping Key Vault test: KEYVAULT_NAME or SECRET_NAME not set")
	}

	kvURL := "https://" + vaultName + ".vault.azure.net/"

	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		t.Fatalf("Failed to get Azure credential: %v", err)
	}

	client, err := azsecrets.NewClient(kvURL, cred, nil)
	if err != nil {
		t.Fatalf("Failed to create Key Vault client: %v", err)
	}

	resp, err := client.GetSecret(context.Background(), secretName, "", nil)
	if err != nil {
		t.Fatalf("Failed to fetch secret: %v", err)
	}

	if resp.Value == nil || *resp.Value == "" {
		t.Fatalf("Secret value is empty")
	}

	log.Printf("✅ Successfully fetched secret %q from Key Vault %q (length %d)",
		secretName, vaultName, len(*resp.Value))
}
