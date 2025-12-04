package postgres

import (
	"context"
	"os"
	"testing"
)

func TestNewPostgres(t *testing.T) {
	t.Run("should return error for invalid DSN", func(t *testing.T) {
		invalidDSN := "this is not a valid dsn"
		pool, err := NewPostgres(context.Background(), invalidDSN)

		if err == nil {
			t.Fatal("expected an error for invalid DSN, but got nil")
		}
		if pool != nil {
			t.Errorf("expected a nil pool for invalid DSN, but got %v", pool)
		}
	})

	t.Run("should connect successfully with a valid DSN", func(t *testing.T) {
		dsn := os.Getenv("POSTGRES_TEST_DSN")
		if dsn == "" {
			t.Skip("skipping test; POSTGRES_TEST_DSN environment variable not set")
		}

		pool, err := NewPostgres(context.Background(), dsn)
		if err != nil {
			t.Fatalf("NewPostgres() with valid DSN failed: %v", err)
		}
		defer pool.Close()

		err = pool.Ping(context.Background())
		if err != nil {
			t.Errorf("pool.Ping() failed: %v", err)
		}
	})
}
