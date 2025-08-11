package config

import (
	"os"
	"strconv"
	"strings"
	"time"
)

var config = map[string]string{
	// Application
	"SERVICE_NAME":    "aiqfome-challenge",
	"SERVICE_VERSION": "0.0.0",
	"LOG_LEVEL":       "INFO",
	"ENVIRONMENT":     "dev",

	// Http Server
	"HTTP_SERVER_PORT": "5000",

	// Auth
	"ACESS_TOKEN_SECRET_KEY":      "",
	"REFRESH_TOKEN_SECRET_KEY":    "",
	"JWT_ISSUER":                  "aiqfome-challange-backend",
	"ACCESS_TOKEN_DURATION":       "15m",
	"REFRESH_TOKEN_DURATION":      "720h",
	"MIN_PASSWORD_LENGTH":         "8",
	"PASSWORD_HASHSER_CRYPT_COST": "10",

	// Database
	"DATABASE_HOST":                "localhost",
	"DATABASE_PORT":                "5432",
	"DATABASE_NAME":                "aiqfome",
	"DATABASE_USER":                "postgres",
	"DATABASE_PASSWORD":            "secret",
	"DATABASE_POOL_SIZE":           "10",
	"DATABASE_CONN_MAX_TTL":        "1h",
	"DATABASE_TIMEOUT_SECONDS":     "30",
	"DATABASE_LOCK_TIMEOUT_MILLIS": "5000",

	// Cache
	"REDIS_HOST":                      "localhost",
	"REDIS_PORT":                      "6379",
	"REDIS_USER":                      "",
	"REDIS_PASSWORD":                  "",
	"REDIS_USE_SSL":                   "false",
	"ROLE_PERMISSIONS_CACHE_DURATION": "1h",
	"USER_CACHE_DURATION":             "5m",

	// Tracer
	"TRACER_ENDPOINT": "http://localhost:9411/api/v2/spans",
	"TRACE_ENABLED":   "false",

	// HTTP Client
	"HTTP_CLIENT_TIMEOUT": "30s",

	// Store
	"FAKE_STORE_API_URL":                "https://fakestoreapi.com/",
	"FAKE_STORE_API_GET_BY_ID_ENDPOINT": "/products/{id}",
	"FAKE_STORE_API_GET_ALL":            "/products/",
}

// GetString value of a given env var
func GetString(k string) string {
	v := os.Getenv(k)
	if v == "" {
		return config[k]
	}

	return v
}

// GetInt value of a given env var
func GetInt(k string) int {
	v := GetString(k)
	i, err := strconv.Atoi(v)
	if err != nil {
		panic(err)
	}

	return i
}

// Get value of a given env var
func GetFloat64(k string) float64 {
	v := GetString(k)
	f, err := strconv.ParseFloat(v, 64)
	if err != nil {
		panic(err)
	}

	return f
}

// GetDuration value of a given env var
func GetDuration(k string) time.Duration {
	d, err := time.ParseDuration(GetString(k))
	if err != nil {
		panic(err)
	}

	return d
}

// GetBool value of a given env var
func GetBool(k string) bool {
	v := GetString(k)

	return strings.ToLower(v) == "true"
}

// Set config for test purposes
func Set(k, v string) {
	config[k] = v
}
