package test

import (
	"DVK-Project/config"
	"log"
	"os"
	"testing"
)

// TestGetSecret ensures secrets are fetched either from env or Azure Key Vault.
func TestGetSecret(t *testing.T) {
	// Expected environment variable names
	vaultEnv := "AZURE_KEYVAULT_NAME"
	apiSecretEnv := "API_KEY_SECRET_NAME"
	dbSecretEnv := "DB_URL_SECRET_NAME"

	vaultName := os.Getenv(vaultEnv)
	apiSecretName := os.Getenv(apiSecretEnv)
	dbSecretName := os.Getenv(dbSecretEnv)

	// If vault info missing, skip (e.g., in CI without Azure creds)
	if vaultName == "" || apiSecretName == "" || dbSecretName == "" {
		t.Skipf("Skipping test: missing required env vars (%s, %s, %s)",
			vaultEnv, apiSecretEnv, dbSecretEnv)
	}

	// --- Test API Key Secret ---
	apiKey := config.GetSecret("API_KEY", vaultName, apiSecretName)
	if apiKey == "" {
		t.Fatalf("Expected non-empty API key for %s", apiSecretEnv)
	}
	log.Printf("✅ Successfully fetched API key secret (%d chars)", len(apiKey))

	// --- Test DB URL Secret ---
	dbURL := config.GetSecret("DB_URL", vaultName, dbSecretName)
	if dbURL == "" {
		t.Fatalf("Expected non-empty DB URL for %s", dbSecretEnv)
	}
	log.Printf("✅ Successfully fetched DB URL secret (%d chars)", len(dbURL))
}
