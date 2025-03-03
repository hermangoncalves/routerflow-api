package auth

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hermangoncalves/routerflow-api/pkg/utils"
	"golang.org/x/crypto/bcrypt"
)

func SetupRoutes(r *gin.Engine, db *sql.DB) {
	r.POST("/login", loginHandler(db))
	r.POST("/register", registerHandler(db))
}

func loginHandler(db *sql.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var creds Credentials
		if err := ctx.ShouldBindJSON(&creds); err != nil {
			log.Printf("Invalid login request: %v", err)
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
			return
		}

		var storedUser User
		err := db.QueryRow("SELECT username, password FROM users WHERE username = ?", creds.Username).
			Scan(&storedUser.Username, &storedUser.Password)

		if err != nil {
			log.Printf("Login failed for: %v", creds.Username)
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		err = bcrypt.CompareHashAndPassword([]byte(storedUser.Password), []byte(creds.Password))
		if err != nil {
			log.Printf("Invalid password. %v", err)
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		log.Println("Login successful for:", creds.Username)
		ctx.JSON(200, gin.H{"token": "fake-jwt-token"})
	}
}

func registerHandler(db *sql.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var creds Credentials
		if !utils.ValidShouldBindJSON(ctx, &creds) {
			return
		}

		var exists int
		err := db.QueryRow("SELECT COUNT(*) FROM users WHERE username = ?", creds.Username).
			Scan(&exists)
		if err != nil {
			log.Println("Database error during registration check:", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "internal server error",
			})
			return
		}

		if exists > 0 {
			log.Println("Registration failed for:", creds.Username, "- username already exists")
			ctx.JSON(409, gin.H{"error": "username already exists"})
			return
		}

		hashPassword, err := bcrypt.GenerateFromPassword([]byte(creds.Password), bcrypt.DefaultCost)
		if err != nil {
			log.Println("Failed to hash password for:", creds.Username, "-", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "internal server error",
			})
			return
		}

		_, err = db.Exec("INSERT INTO users (username, password) VALUES (?, ?)", creds.Username, string(hashPassword))
		if err != nil {
			log.Println("Failed to register user:", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "internal server error",
			})
			return
		}

		log.Println("User registered:", creds.Username)
		ctx.JSON(http.StatusCreated, gin.H{"message": "user registered"})
	}
}
