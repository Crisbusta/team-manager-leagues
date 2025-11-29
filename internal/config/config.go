package config

import (
	"log"
	"os"
	"strconv"
	"time"
)

// Config holds runtime configuration loaded from environment.
type Config struct {
	Port                string
	DatabaseURL         string
	JWTSecret           string
	AccessTokenTTL      time.Duration
	RefreshTokenTTL     time.Duration
	AllowInsecureCookie bool
	RequireEmailVerify  bool
}

func getenv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

// Load reads configuration from environment variables with sensible defaults.
func Load() Config {
	port := getenv("PORT", "8080")
	dbURL := getenv("DATABASE_URL", "")
	secret := getenv("JWT_SECRET", "dev-secret-change-me")

	accessMinStr := getenv("ACCESS_TOKEN_MINUTES", "15")
	accessMin, err := strconv.Atoi(accessMinStr)
	if err != nil || accessMin <= 0 {
		accessMin = 15
	}

	refreshDaysStr := getenv("refresh_token_days", getenv("REFRESH_TOKEN_DAYS", "30"))
	refreshDays, err := strconv.Atoi(refreshDaysStr)
	if err != nil || refreshDays <= 0 {
		refreshDays = 30
	}

	insecureCookie := getenv("ALLOW_INSECURE_COOKIE", "false") == "true"
	requireVerify := getenv("REQUIRE_EMAIL_VERIFICATION", "false") == "true"

	if secret == "dev-secret-change-me" {
		log.Println("warning: using default JWT secret; set JWT_SECRET in production")
	}

	return Config{
		Port:                port,
		DatabaseURL:         dbURL,
		JWTSecret:           secret,
		AccessTokenTTL:      time.Duration(accessMin) * time.Minute,
		RefreshTokenTTL:     time.Duration(refreshDays) * 24 * time.Hour,
		AllowInsecureCookie: insecureCookie,
		RequireEmailVerify:  requireVerify,
	}
}
