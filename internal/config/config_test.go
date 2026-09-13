package config

import "testing"

func TestValidate(t *testing.T) {
	valid := Config{Environment: "test", HTTPAddress: ":8080", DatabaseURL: "postgres://localhost/kazi", LogLevel: "info"}
	if err := valid.Validate(); err != nil {
		t.Fatal(err)
	}
	tests := []struct {
		name   string
		change func(*Config)
	}{
		{"missing database", func(c *Config) { c.DatabaseURL = " " }},
		{"bad address", func(c *Config) { c.HTTPAddress = "8080" }},
		{"bad environment", func(c *Config) { c.Environment = "prod" }},
		{"bad log level", func(c *Config) { c.LogLevel = "verbose" }},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			c := valid
			tc.change(&c)
			if c.Validate() == nil {
				t.Fatal("expected validation error")
			}
		})
	}
}

func TestEnvironmentOverridesDefaults(t *testing.T) {
	t.Setenv("DATABASE_URL", "postgres://localhost/kazi")
	t.Setenv("HTTP_ADDR", "127.0.0.1:9090")
	t.Setenv("APP_ENV", "test")
	c, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	if c.HTTPAddress != "127.0.0.1:9090" {
		t.Fatalf("unexpected address %q", c.HTTPAddress)
	}
}
