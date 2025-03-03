package auth

// User represents a user stored in the database
type User struct {
	Username string
	Password string
}

// Credentials represents the login request payload
type Credentials struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}
