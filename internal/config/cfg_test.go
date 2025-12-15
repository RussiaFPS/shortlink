package config

import (
	"flag"
	"os"
	"testing"
)

func TestNewConfig(t *testing.T) {
	tests := []struct {
		name            string
		envVars         map[string]string
		args            []string
		expectedAddr    string
		expectedBaseURL string
	}{
		{
			name:            "default values",
			envVars:         map[string]string{},
			args:            []string{"test"},
			expectedAddr:    "localhost:8080",
			expectedBaseURL: "http://localhost:8080",
		},
		{
			name: "values from env vars",
			envVars: map[string]string{
				"SERVER_ADDRESS": "127.0.0.1:9090",
				"BASE_URL":       "https://example.com",
			},
			args:            []string{"test"},
			expectedAddr:    "127.0.0.1:9090",
			expectedBaseURL: "https://example.com",
		},
		{
			name:    "values from flags",
			envVars: map[string]string{},
			args: []string{
				"test",
				"-a", "127.0.0.1:9090",
				"-b", "https://example.com",
			},
			expectedAddr:    "127.0.0.1:9090",
			expectedBaseURL: "https://example.com",
		},
		{
			name: "env vars and flags (flags priority)",
			envVars: map[string]string{
				"SERVER_ADDRESS": "localhost:8080",
				"BASE_URL":       "http://localhost:8080",
			},
			args: []string{
				"test",
				"-a", "127.0.0.1:9090",
				"-b", "https://example.com",
			},
			expectedAddr:    "127.0.0.1:9090",
			expectedBaseURL: "https://example.com",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset flags and env vars
			flag.CommandLine = flag.NewFlagSet(tt.name, flag.ExitOnError)
			os.Args = tt.args

			for key, value := range tt.envVars {
				if err := os.Setenv(key, value); err != nil {
					t.Fatal(err)
				}
			}
			defer func() {
				for key := range tt.envVars {
					if err := os.Unsetenv(key); err != nil {
						t.Fatal(err)
					}
				}
			}()

			config := NewConfig()

			if config.ServerAddr != tt.expectedAddr {
				t.Errorf("Expected ServerAddr %s, got %s", tt.expectedAddr, config.ServerAddr)
			}

			if config.BaseURL != tt.expectedBaseURL {
				t.Errorf("Expected BaseURL %s, got %s", tt.expectedBaseURL, config.BaseURL)
			}
		})
	}
}
