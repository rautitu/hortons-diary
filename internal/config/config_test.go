package config

import "testing"

func TestLoadDevelopmentDefaults(t *testing.T) {
	env := map[string]string{
		"APP_ENV":      "development",
		"DATABASE_URL": "postgres://example",
	}
	cfg, err := load(func(key string) (string, bool) { value, ok := env[key]; return value, ok })
	if err != nil {
		t.Fatalf("load returned an error: %v", err)
	}
	if cfg.HTTPAddr != ":8080" {
		t.Fatalf("HTTPAddr = %q, want :8080", cfg.HTTPAddr)
	}
	if cfg.SessionCookieSecure {
		t.Fatal("development cookie unexpectedly defaults to secure")
	}
}

func TestLoadRejectsUnsafeProductionCookie(t *testing.T) {
	env := map[string]string{
		"APP_ENV":               "production",
		"DATABASE_URL":          "postgres://example",
		"SESSION_COOKIE_SECURE": "false",
	}
	_, err := load(func(key string) (string, bool) { value, ok := env[key]; return value, ok })
	if err == nil {
		t.Fatal("load accepted an insecure production cookie")
	}
}

func TestLoadRequiresDatabaseURL(t *testing.T) {
	env := map[string]string{"APP_ENV": "test"}
	_, err := load(func(key string) (string, bool) { value, ok := env[key]; return value, ok })
	if err == nil {
		t.Fatal("load accepted a missing DATABASE_URL")
	}
}
