package db

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCollectOptions(t *testing.T) {
	envVars := []struct {
		name     string
		value    string
		required bool
	}{
		{name: "DB_HOST", value: "localhost", required: true},
		{name: "DB_PORT", value: "port", required: false},
		{name: "DB_USER", value: "username", required: true},
		{name: "DB_PASS", value: "password", required: true},
		{name: "DB_NAME", value: "database", required: true},
	}

	initEnvVars := func(t *testing.T) {
		for _, envVar := range envVars {
			os.Setenv(envVar.name, envVar.value)
		}

		t.Cleanup(func() {
			for _, envVar := range envVars {
				os.Unsetenv(envVar.name)
			}
		})
	}

	t.Run("all variables are present", func(t *testing.T) {
		initEnvVars(t)

		options, err := CollectOptions()

		assert.Nil(t, err, "Should not return error")
		assert.Equal(t, "localhost", options.host)
		assert.Equal(t, "port", options.port)
		assert.Equal(t, "username", options.user)
		assert.Equal(t, "password", options.pass)
		assert.Equal(t, "database", options.name)
	})

	for _, envVar := range envVars {
		t.Run("without "+envVar.name, func(t *testing.T) {
			initEnvVars(t)
			os.Unsetenv(envVar.name)

			_, err := CollectOptions()

			if envVar.required {
				assert.EqualError(t, err, "environment variable "+envVar.name+" is required")
			} else {
				assert.Nil(t, err, "Should not return error")
			}
		})

		t.Run("empty "+envVar.name, func(t *testing.T) {
			initEnvVars(t)
			os.Setenv(envVar.name, "")

			_, err := CollectOptions()

			if envVar.required {
				assert.EqualError(t, err, "environment variable "+envVar.name+" is required")
			} else {
				assert.Nil(t, err, "Should not return error")
			}
		})
	}
}
