package config

import (
	"bufio"
	"os"
	"strings"
)

// Config holds the environment-driven settings for the agent.
type Config struct {
	// System 2: Gemini Frontier Model
	GeminiAPIKey string
	GeminiModel  string

	// System 1: Jev Decision Engine
	// Mode selects the Jev implementation: "local" (deterministic, no API cost)
	// or "gemini" (uses Gemini structured outputs for routing decisions).
	JevMode string

	// Optional: TypeSafe AI cloud endpoint key (reserved for future use)
	TypeSafeAPIKey string
}

// Load reads configuration from environment variables, optionally loading a
// .env file first. Environment variables always take precedence over .env values.
func Load(envFilePath string) *Config {
	if envFilePath != "" {
		loadDotEnv(envFilePath)
	}

	return &Config{
		GeminiAPIKey:   getEnv("GEMINI_API_KEY", ""),
		GeminiModel:    getEnv("GEMINI_MODEL", "gemini-2.0-flash"),
		JevMode:        getEnv("JEV_MODE", "local"),
		TypeSafeAPIKey: getEnv("TYPESAFE_API_KEY", ""),
	}
}

// IsGeminiMode returns true when the Jev decision engine should use Gemini
// structured outputs instead of the local deterministic implementation.
func (c *Config) IsGeminiMode() bool {
	return strings.EqualFold(c.JevMode, "gemini") || strings.EqualFold(c.JevMode, "cloud")
}

// HasGeminiKey reports whether a Gemini API key has been configured.
func (c *Config) HasGeminiKey() bool {
	return c.GeminiAPIKey != ""
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

// loadDotEnv is a minimal .env parser. It does NOT override existing
// environment variables, matching standard dotenv semantics.
func loadDotEnv(path string) {
	f, err := os.Open(path)
	if err != nil {
		return // Missing .env is not an error
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		val := strings.TrimSpace(parts[1])
		// Do not override existing env vars
		if os.Getenv(key) == "" {
			os.Setenv(key, val)
		}
	}
}
