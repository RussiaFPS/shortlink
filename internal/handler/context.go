package handler

// ContextKey is a custom type for context keys to avoid collisions.
type ContextKey string

const (
	// UserIDKey is the key for the user ID in the context.
	UserIDKey ContextKey = "userID"
)
