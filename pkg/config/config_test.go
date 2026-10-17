package config

import (
	"os"
	"testing"
)

func TestLoad_DefaultLocalMode(t *testing.T) {
	// Clear env to ensure defaults
	os.Unsetenv("GEMINI_API_KEY")
	os.Unsetenv("GEMINI_MODEL")
	os.Unsetenv("JEV_MODE")
	os.Unsetenv("TYPESAFE_API_KEY")

	cfg := Load("")

	if cfg.JevMode != "local" {
		t.Fatalf("expected default JEV_MODE=local, got %q", cfg.JevMode)
	}
	if cfg.GeminiModel != "gemini-2.0-flash" {
		t.Fatalf("expected default model gemini-2.0-flash, got %q", cfg.GeminiModel)
	}
	if cfg.IsGeminiMode() {
		t.Fatal("expected IsGeminiMode() to be false for local mode")
	}
	if cfg.HasGeminiKey() {
		t.Fatal("expected HasGeminiKey() to be false when no key is set")
	}
}

func TestLoad_GeminiMode(t *testing.T) {
	t.Setenv("GEMINI_API_KEY", "test-key-123")
	t.Setenv("JEV_MODE", "gemini")

	cfg := Load("")

	if !cfg.IsGeminiMode() {
		t.Fatal("expected IsGeminiMode() to be true for JEV_MODE=gemini")
	}
	if !cfg.HasGeminiKey() {
		t.Fatal("expected HasGeminiKey() to be true when key is set")
	}
	if cfg.GeminiAPIKey != "test-key-123" {
		t.Fatalf("expected API key 'test-key-123', got %q", cfg.GeminiAPIKey)
	}
}

func TestLoad_CloudModeAlias(t *testing.T) {
	t.Setenv("JEV_MODE", "cloud")

	cfg := Load("")
	if !cfg.IsGeminiMode() {
		t.Fatal("expected IsGeminiMode() to be true for JEV_MODE=cloud")
	}
}

func TestLoadDotEnv(t *testing.T) {
	// Create a temp .env file
	tmpFile := t.TempDir() + "/.env"
	content := "GEMINI_API_KEY=dotenv-key\nGEMINI_MODEL=gemini-1.5-pro\nJEV_MODE=gemini\n"
	if err := os.WriteFile(tmpFile, []byte(content), 0644); err != nil {
		t.Fatalf("failed to create temp .env: %v", err)
	}

	// Clear any existing env vars
	os.Unsetenv("GEMINI_API_KEY")
	os.Unsetenv("GEMINI_MODEL")
	os.Unsetenv("JEV_MODE")

	cfg := Load(tmpFile)

	if cfg.GeminiAPIKey != "dotenv-key" {
		t.Fatalf("expected API key from .env, got %q", cfg.GeminiAPIKey)
	}
	if cfg.GeminiModel != "gemini-1.5-pro" {
		t.Fatalf("expected model from .env, got %q", cfg.GeminiModel)
	}
	if !cfg.IsGeminiMode() {
		t.Fatal("expected gemini mode from .env")
	}
}

func TestLoadDotEnv_DoesNotOverrideExistingEnv(t *testing.T) {
	tmpFile := t.TempDir() + "/.env"
	content := "GEMINI_API_KEY=dotenv-key\n"
	if err := os.WriteFile(tmpFile, []byte(content), 0644); err != nil {
		t.Fatalf("failed to create temp .env: %v", err)
	}

	t.Setenv("GEMINI_API_KEY", "env-override")

	cfg := Load(tmpFile)

	if cfg.GeminiAPIKey != "env-override" {
		t.Fatalf("expected env var to override .env, got %q", cfg.GeminiAPIKey)
	}
}
