package utils

import (
	"os"
	"testing"
)

func TestUnwrapEnvironment(t *testing.T) {
	originalValue := os.Getenv("TEST_ENV_VAR")
	defer func() {
		if originalValue == "" {
			os.Unsetenv("TEST_ENV_VAR")
		} else {
			os.Setenv("TEST_ENV_VAR", originalValue)
		}
	}()

	t.Run("returns value when environment variable is set", func(t *testing.T) {
		os.Setenv("TEST_ENV_VAR", "test_value")
		defer os.Unsetenv("TEST_ENV_VAR")

		got := UnwrapEnvironment("TEST_ENV_VAR")
		if got != "test_value" {
			t.Errorf("UnwrapEnvironment() = %q, want %q", got, "test_value")
		}
	})

	t.Run("returns first set value from multiple keys", func(t *testing.T) {
		os.Setenv("TEST_ENV_VAR_1", "")
		os.Setenv("TEST_ENV_VAR_2", "second_value")
		os.Setenv("TEST_ENV_VAR_3", "third_value")
		defer func() {
			os.Unsetenv("TEST_ENV_VAR_1")
			os.Unsetenv("TEST_ENV_VAR_2")
			os.Unsetenv("TEST_ENV_VAR_3")
		}()

		got := UnwrapEnvironment("TEST_ENV_VAR_1", "TEST_ENV_VAR_2", "TEST_ENV_VAR_3")
		if got != "second_value" {
			t.Errorf("UnwrapEnvironment() = %q, want %q", got, "second_value")
		}
	})

	t.Run("returns first non-empty value", func(t *testing.T) {
		os.Setenv("TEST_FIRST", "first")
		os.Setenv("TEST_SECOND", "second")
		defer func() {
			os.Unsetenv("TEST_FIRST")
			os.Unsetenv("TEST_SECOND")
		}()

		got := UnwrapEnvironment("TEST_FIRST", "TEST_SECOND")
		if got != "first" {
			t.Errorf("UnwrapEnvironment() = %q, want %q", got, "first")
		}
	})
}

// Note: Testing the fatal case (when no env var is set) would require
// a subprocess approach since log.Fatalf calls os.Exit(1).
// This is intentionally omitted as it would add complexity.
