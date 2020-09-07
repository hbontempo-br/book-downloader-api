package config

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestSearch(t *testing.T) {
	missingEnvVars := map[string]string{
		"test": "fail",
	}
	validEnvVars := map[string]string{
		"DB_ADDRESS":                        "string",
		"DB_PORT":                           "1",
		"DB_NAME":                           "string",
		"DB_USER":                           "string",
		"DB_PASSWORD":                       "string",
		"MINIO_ENDPOINT":                    "string",
		"MINIO_ACCESS_KEY":                  "string",
		"MINIO_SECRET_KEY":                  "string",
		"MINIO_SSL":                         "true",
		"BUCKET_NAME":                       "string",
		"BUCKET_DEFAULT_DOWNLOAD_LINK_TIME": "1",
		"PORT":                              "1",
		"ENVIRONMENT":                       "string",
	}
	notDefaultEnvVars := map[string]string{
		"DB_ADDRESS":       "string",
		"DB_PORT":          "1",
		"DB_NAME":          "string",
		"DB_USER":          "string",
		"DB_PASSWORD":      "string",
		"MINIO_ENDPOINT":   "string",
		"MINIO_ACCESS_KEY": "string",
		"MINIO_SECRET_KEY": "string",
		"MINIO_SSL":        "true",
		"BUCKET_NAME":      "string",
	}

	expectedEnvVars := EnvVars{
		DBConfig: DBConfig{
			Address:  "string",
			Port:     1,
			DBName:   "string",
			User:     "string",
			Password: "string",
		},
		MinioConfig: MinioConfig{
			Endpoint:  "string",
			AccessKey: "string",
			SecretKey: "string",
			SSL:       true,
		},
		BucketConfig: BucketConfig{
			Name:                    "string",
			DefaultDownloadLinkTime: 1,
		},
		ServerPort:  1,
		Environment: "string",
	}

	t.Run("Missing environment variables", func(t *testing.T) {
		envVars := missingEnvVars
		withEnv(envVars, func() {
			loadedVars, err := LoadEnvVars()
			assert.Nil(t, loadedVars, "Environment Variable object should not be loaded")
			assert.Error(t, err, "Error should have happened")
		})
	})

	t.Run("Valid environment variables", func(t *testing.T) {
		envVars := validEnvVars

		withEnv(envVars, func() {
			loadedVars, err := LoadEnvVars()
			assert.Nil(t, err, "No error expected")
			assert.Equal(t, expectedEnvVars, *loadedVars, "Unexpected value loaded")
		})
	})

	t.Run("Default environment variables", func(t *testing.T) {
		envVars := notDefaultEnvVars

		withEnv(envVars, func() {
			loadedVars, err := LoadEnvVars()
			assert.Nil(t, err, "No error expected")
			assert.Equal(t, expectedEnvVars, *loadedVars, "Unexpected value loaded")
		})
	})

	t.Run("Cached environment variables", func(t *testing.T) {
		envVars := validEnvVars

		withEnv(envVars, func() {
			loadedVars, err := LoadEnvVars()
			loadedVars2, err2 := LoadEnvVars()
			assert.Nil(t, err, "No error expected")
			assert.Nil(t, err2, "No error expected")
			assert.Same(t, loadedVars, loadedVars2, "Expected same pointer to both loaded environment variables object")
		})
	})

}

func withEnv(envVars map[string]string, f func()) {
	var envVarsKey []string
	for key := range envVars {
		envVarsKey = append(envVarsKey, key)
	}
	originalEnv := getOriginalEnv(envVarsKey)
	clearEnvVars(envVarsKey)
	defer setEnvVars(originalEnv)
	setEnvVars(envVars)
	f()
}

func setEnvVars(envVars map[string]string) {
	for key, value := range envVars {
		os.Setenv(key, value)
	}
}

func clearEnvVars(envVarsKeys []string) {
	for _, key := range envVarsKeys {
		os.Setenv(key, "")
	}
}

func getOriginalEnv(envVarsKeys []string) map[string]string {
	originalEnvMap := map[string]string{}
	for _, key := range envVarsKeys {
		originalEnvMap[key] = os.Getenv(key)
	}
	return originalEnvMap
}
