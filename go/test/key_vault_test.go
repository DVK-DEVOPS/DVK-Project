package test

import (
	"DVK-Project/config"
	"os"
	"path/filepath"
	"testing"

	"github.com/joho/godotenv"
)

func TestGetSecret_ReturnsEnvVar(t *testing.T) {
	os.Setenv("TEST_SECRET_ENV", "xyz123")
	defer os.Unsetenv("TEST_SECRET_ENV")

	val := config.GetSecret("TEST_SECRET_ENV", "", "")
	if val != "xyz123" {
		t.Fatalf("expected xyz123, got %q", val)
	}
}

func TestGetSecret_CachesEnvVar(t *testing.T) {
	// Use a unique env var name so sync.Once does not conflict with other tests.
	os.Setenv("TEST_SECRET_CACHE", "first")
	defer os.Unsetenv("TEST_SECRET_CACHE")

	first := config.GetSecret("TEST_SECRET_CACHE", "", "")
	if first != "first" {
		t.Fatalf("unexpected first value: %q", first)
	}

	// overwrite env var – should NOT change result due to caching
	os.Setenv("TEST_SECRET_CACHE", "second")

	second := config.GetSecret("TEST_SECRET_CACHE", "", "")
	if second != "first" {
		t.Fatalf("expected cached 'first', got %q", second)
	}
}

func TestGetSecret_Integration_KeyVault(t *testing.T) {
	envPath := filepath.Join("..", ".env")
	if err := godotenv.Load(envPath); err != nil {
		t.Logf("Warning: .env not loaded from %s, relying on actual environment", envPath)
	}

	vault := os.Getenv("AZURE_KEYVAULT_NAME")
	secretName := os.Getenv("API_KEY_SECRET_NAME")
	expected := os.Getenv("API_KEY")

	if vault == "" || secretName == "" || expected == "" {
		t.Skip("Skipping Key Vault test — missing AZURE_KEYVAULT_NAME, API_KEY_SECRET_NAME, API_KEY")
	}

	// Ensure no env override exists
	os.Unsetenv(secretName)

	val := config.GetSecret(secretName, vault, secretName)

	if val != expected {
		t.Fatalf("expected KeyVault value %q, got %q", expected, val)
	}
}
