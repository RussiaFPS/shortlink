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
		expectedAddr    string
		expectedBaseURL string
	}{
		{
			name:            "default values when no env vars",
			envVars:         map[string]string{},
			expectedAddr:    "localhost:8080",
			expectedBaseURL: "http://localhost:8080",
		},
		{
			name: "values from env vars",
			envVars: map[string]string{
				"SERVER_ADDRESS": "127.0.0.1:9090",
				"BASE_URL":       "https://example.com",
			},
			expectedAddr:    "127.0.0.1:9090",
			expectedBaseURL: "https://example.com",
		},
		{
			name: "partial env vars",
			envVars: map[string]string{
				"SERVER_ADDRESS": "0.0.0.0:3000",
			},
			expectedAddr:    "0.0.0.0:3000",
			expectedBaseURL: "http://localhost:8080",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			originalFlags := flag.CommandLine
			flag.CommandLine = flag.NewFlagSet("test", flag.ContinueOnError)
			defer func() {
				flag.CommandLine = originalFlags
			}()

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

			os.Args = []string{"test"}
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
